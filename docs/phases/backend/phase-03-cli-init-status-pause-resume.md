# Phase 03: Primary CLI Surface (`init`, `status`, `pause`, `resume`)

**Domain:** Backend / CLI Engine  
**Status:** Ready to Implement  
**Specification Source:** `CTXbank.md` (§1.1, §3, §3.2)  

---

## 1. Phase Objective & Boundary
Implement the static CLI binary commands for day-to-day developer and agent lifecycle management:
- `ctx init`: Scaffold memory bank, state files, and minimal vendor stubs.
- `ctx status`: Output a single-shot, zero-token terminal status brief or `--json` structured envelope.
- `ctx pause`: Execute a safe repository checkpoint and interactively capture manual out-of-band edits to update `activeContext.md` and `decisionLog.md`.
- `ctx resume`: Render an instant pickup brief for a human developer or `--agent` context pre-loader.

*Exclusions:* AST analysis, Ollama LLM queries, and full-screen Bubbletea TUI views.

---

## 2. Prerequisites & Dependencies
- Phase 01 (`internal/core`) and Phase 02 (`internal/git`, `internal/checkpoint`) complete.
- CLI argument handling (`spf13/cobra` or Go `flag`).
- Terminal formatting (`internal/tui/format.go`).

---

## 3. Step-by-Step Implementation Steps
1. **Build CLI Entrypoint (`cmd/ctx/main.go`):**
   - Initialize command root `ctx` with flags (`--version`, `--help`, `--dir`).
   - Wire subcommands: `init`, `status`, `pause`, `resume`.
2. **Implement `ctx init` (`cmd/ctx/cmd_init.go`):**
   - Verify if `memory-bank/` already exists; prevent accidental clobbering unless `--force` is passed.
   - Scaffold directories and templates via `core.InitBank`.
   - Call vendor rules stub generator (Phase 09).
3. **Implement `ctx status` (`cmd/ctx/cmd_status.go`, `internal/tui/status.go`):**
   - Inspect git state, dirty file counts, branch name.
   - Read hot files (`activeContext.md`, `progress.md`).
   - Calculate staleness (days and commits since last checkpoint).
   - If `--json` flag is provided, emit clean JSON payload to stdout.
   - Otherwise, print ASCII/Unicode formatted status card.
4. **Implement `ctx pause` (`cmd/ctx/cmd_pause.go`):**
   - Check git safety (abort if mid-rebase).
   - Create new checkpoint snapshot.
   - Identify files modified outside active agent session.
   - Prompt developer interactively: `What did you change here? (one line, Enter to skip)`.
   - Atomically update `activeContext.md` (under `## Recent`) and append to `decisionLog.md`.
5. **Implement `ctx resume` (`cmd/ctx/cmd_resume.go`):**
   - Load latest checkpoint and active context.
   - Format pickup brief (current focus, last changes, upcoming next steps).
6. **Develop Integration Tests (`cmd/ctx/cli_test.go`):**
   - Execute CLI commands using test harnesses and assert stdout outputs.

---

## 4. Concrete Deliverables
- `cmd/ctx/main.go`
- `cmd/ctx/cmd_init.go`, `cmd/ctx/cmd_status.go`, `cmd/ctx/cmd_pause.go`, `cmd/ctx/cmd_resume.go`
- `internal/tui/format.go` & `internal/tui/status.go`
- `cmd/ctx/cli_test.go`

---

## 5. Validation & Verification Routine
```bash
go build -o ./bin/ctx.exe ./cmd/ctx
./bin/ctx init
./bin/ctx status
./bin/ctx status --json
./bin/ctx pause
./bin/ctx resume
```
- Ensure zero ANSI escape codes when stdout is redirected to a pipe.
- Ensure `--json` output validates against JSON schema.

---

## 6. Security & Vulnerability Audit
- Handle EOF and SIGINT cleanly during interactive `ctx pause` prompts without leaving dangling lock files or corrupt state.
- Compile static binary with `CGO_ENABLED=0` to ensure size < 15MB.
