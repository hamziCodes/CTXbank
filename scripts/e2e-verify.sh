#!/bin/bash
# End-to-end verification for the kai branch fixes.
# Builds ctx from source and exercises every fix. Fails fast on any regression.
set -u
command -v go >/dev/null && : || { echo "go not in PATH"; exit 1; }

REPO="$(cd "$(dirname "$0")/.." && pwd)"
WORK="$(mktemp -d)"
export HOME="$WORK/fakehome"   # isolate the projects registry
mkdir -p "$HOME"
PASS=0; FAIL=0

ok()   { PASS=$((PASS+1)); echo "  PASS: $1"; }
bad()  { FAIL=$((FAIL+1)); echo "  FAIL: $1"; }

echo "== building ctx =="
(cd "$REPO" && go build -o "$WORK/ctx" ./cmd/ctx) || { echo "BUILD FAILED"; exit 1; }
CTX="$WORK/ctx"

echo "== Fix 1: canonical module path =="
"$REPO" >/dev/null 2>&1
(cd "$REPO" && go list -m | grep -q "github.com/hamziCodes/CTXbank") && ok "go list -m == repo path" || bad "module path"
! grep -rq "github.com/ctxbank/ctx" "$REPO" --include="*.go" && ok "no stale import paths" || bad "stale imports remain"

echo "== Fix 2: neutral init templates =="
PROJ="$WORK/demo-proj"; mkdir -p "$PROJ"; cd "$PROJ" && git init -q .
"$CTX" init >/dev/null
! grep -qi "ctxbank" memory-bank/projectbrief.md && ok "projectbrief has no CTXbank self-description" || bad "projectbrief mentions CTXbank"
grep -q "TODO" memory-bank/projectbrief.md && ok "projectbrief has TODO placeholders" || bad "no TODOs"
grep -q "ctx audit --apply" memory-bank/activeContext.md && ok "activeContext points at audit/prompt-sync" || bad "activeContext guidance"

echo "== core flows still work =="
"$CTX" status >/dev/null && ok "ctx status" || bad "ctx status"
"$CTX" resume | grep -q "Active Context" && ok "ctx resume" || bad "ctx resume"
"$CTX" audit >/dev/null && ok "ctx audit" || bad "ctx audit"
"$CTX" lint-memory >/dev/null && ok "ctx lint-memory" || bad "ctx lint-memory"

echo "== Fix 7: pause flags + non-TTY =="
echo "change" > notes.txt
"$CTX" pause --note "e2e note" --focus "e2e focus" </dev/null >/dev/null \
  && ok "pause --note/--focus non-interactive" || bad "pause flags"
ID1=$(ls -t memory-bank/.state/checkpoints/ | head -1)
grep -q "e2e focus" "memory-bank/.state/checkpoints/$ID1" && ok "focus recorded in checkpoint" || bad "focus missing"
echo "piped" | "$CTX" pause >/dev/null 2>&1 && ok "pause with piped stdin does not hang" || bad "pause piped"

echo "== Fix 4: checkpoint collision =="
N_BEFORE=$(ls memory-bank/.state/checkpoints/ | wc -l)
"$CTX" pause --note "a" </dev/null >/dev/null; "$CTX" pause --note "b" </dev/null >/dev/null
N_AFTER=$(ls memory-bank/.state/checkpoints/ | wc -l)
[ "$N_AFTER" -ge "$((N_BEFORE+2))" ] && ok "rapid checkpoints never overwrite ($N_BEFORE -> $N_AFTER)" || bad "checkpoint overwrite"

echo "== Fix 8: ctx doctor =="
"$CTX" doctor | grep -q "all critical checks passed" && ok "doctor passes in healthy project" || bad "doctor"
mkdir -p "$WORK/empty" && (cd "$WORK/empty" && "$CTX" doctor >/dev/null 2>&1) \
  && bad "doctor should fail without bank" || ok "doctor fails cleanly without bank"

echo "== Fix 9: completions =="
"$CTX" completion bash | grep -q "pause" && ok "bash completion lists commands" || bad "bash completion"
"$CTX" completion zsh | grep -q "_ctx" && ok "zsh completion" || bad "zsh completion"
"$CTX" completion fish | grep -q "complete -c ctx" && ok "fish completion" || bad "fish completion"
"$CTX" completion powershell | grep -q "Register-ArgumentCompleter" && ok "powershell completion" || bad "ps completion"
"$CTX" completion bogus >/dev/null 2>&1 && bad "bad shell should fail" || ok "bad shell rejected"

echo "== Fix 10: recent-projects registry =="
[ -f "$HOME/.ctxbank/projects.json" ] && ok "registry file created on init" || bad "no registry"
"$CTX" list | grep -q "demo-proj" && ok "ctx list shows project from registry" || bad "ctx list registry"
"$CTX" list --json | grep -q "demo-proj" && ok "ctx list --json" || bad "ctx list --json"

echo "== Fix 5: dashboard auth token =="
PORT=$(python3 -c "import socket; s=socket.socket(); s.bind(('127.0.0.1',0)); print(s.getsockname()[1]); s.close()")
"$CTX" ui --no-open --port "$PORT" >"$WORK/ui.log" 2>&1 &
UIPID=$!
sleep 2
URL=$(grep -o "http://localhost:$PORT?token=[a-f0-9]*" "$WORK/ui.log" | head -1)
TOKEN=$(echo "$URL" | sed 's/.*token=//')
if [ -z "$TOKEN" ]; then bad "no token printed"; else ok "token printed in terminal"; fi
CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:$PORT/api/status")
[ "$CODE" = "401" ] && ok "API without token -> 401" || bad "no-token gave $CODE"
CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:$PORT/api/status?token=$TOKEN")
[ "$CODE" = "200" ] && ok "API with token -> 200" || bad "token gave $CODE"
CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "X-CTX-Token: $TOKEN" "http://localhost:$PORT/api/files")
[ "$CODE" = "200" ] && ok "API with header token -> 200" || bad "header token gave $CODE"
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "http://localhost:$PORT/api/file/save?token=$TOKEN" \
  -d '{"filename":"evil.json","content":"x"}')
[ "$CODE" = "400" ] && ok "evil.json save -> 400" || bad "evil.json gave $CODE"
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "http://localhost:$PORT/api/file/save?token=$TOKEN" \
  -d '{"filename":"progress.md","content":"# Progress\n- e2e\n"}')
[ "$CODE" = "200" ] && ok "progress.md save -> 200" || bad "progress.md gave $CODE"
kill $UIPID 2>/dev/null; wait 2>/dev/null

echo "== Fix 6: ctx-mcp removed =="
[ ! -d "$REPO/cmd/ctx-mcp" ] && ok "cmd/ctx-mcp deleted" || bad "cmd/ctx-mcp still exists"
! grep -rq "ctx-mcp" "$REPO/.github" "$REPO/Makefile" "$REPO/scripts/install.sh" && ok "no ctx-mcp in workflows/Makefile/install" || bad "ctx-mcp refs remain"

echo "== Fix 3: version injection =="
(cd "$REPO" && go build -ldflags="-s -w -X main.Version=v0.0-e2e" -o "$WORK/ctx-ver" ./cmd/ctx)
"$WORK/ctx-ver" --version | grep -q "v0.0-e2e" && ok "ldflags version stamp works" || bad "version stamp"
grep -q 'main.Version=${{ github.ref_name }}' "$REPO/.github/workflows/release.yml" && ok "release.yml injects tag" || bad "release.yml"

echo "== MCP serve --mcp =="
MCP_OUT=$(printf '%s\n' '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' \
  '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | (cd "$PROJ" && "$CTX" serve --mcp))
echo "$MCP_OUT" | grep -q "read_active_context" && ok "MCP tools/list works" || bad "MCP tools/list"
echo "$MCP_OUT" | grep -q "memory-bank-mcp" && ok "MCP initialize works" || bad "MCP initialize"

echo ""
echo "==================================="
echo "PASS: $PASS  FAIL: $FAIL"
echo "==================================="
rm -rf "$WORK"
[ "$FAIL" -eq 0 ]
