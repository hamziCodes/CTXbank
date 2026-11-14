# CTXbank: Orchestration Blueprint & Architectural Execution Engine

**Document Version:** 1.0.0  
**Project:** CTXbank (The Memory Bank & Project Lifecycle Engine)  
**Specification Source:** `CTXbank.md`  
**Classification:** Core System Blueprint & Autonomous Orchestration Specification  

---

## 1. System Architecture & Directory Hierarchy

CTXbank is a deterministic, crash-safe, local-first binary and Model Context Protocol (MCP) server that establishes a single source of truth for coding agent memory and developer project resumption. It replaces loose prompt-engineering conventions with cryptographic file hashing, atomic filesystem operations, porcelain git integration, AST analysis, and delta-synchronized MCP resources.

### 1.1 Complete Project Directory Structure

```text
CTXbank/
├── .cursor/
│   └── rules/
│       └── memory-bank.mdc             # Auto-generated minimal vendor stub
├── .cursorrules                        # Minimal vendor stub pointing to memory-bank/
├── .rules                              # Autonomous Agent Governance & Behavioral Constraints
├── CLAUDE.md                           # Claude Code minimal vendor stub
├── .agentrules                         # Generic agent minimal stub
├── .github/
│   └── copilot-instructions.md         # GitHub Copilot minimal vendor stub
├── .context/
│   └── config.toml                     # User config (Ollama URL, model overrides, keys via env)
├── cmd/
│   ├── ctx/
│   │   └── main.go                     # Primary static CLI entrypoint
│   └── ctx-mcp/
│       └── main.go                     # Standalone/Stdio/HTTP MCP Server entrypoint
├── internal/
│   ├── config/
│   │   ├── config.go                   # Config parser (.context/config.toml) & env resolution
│   │   └── config_test.go
│   ├── core/
│   │   ├── atomic.go                   # Write to <file>.tmp -> fsync -> rename() crash-safety
│   │   ├── atomic_test.go              # Crash & interruption simulation tests
│   │   ├── manifest.go                 # SHA-256 calculation, manifest.json manager, delta calculator
│   │   ├── manifest_test.go
│   │   ├── bank.go                     # Memory bank loader, volatility router, template generator
│   │   └── bank_test.go
│   ├── git/
│   │   ├── porcelain.go                # Git CLI wrapper (status --porcelain, diff, log, branch)
│   │   ├── safeguards.go               # Detect detached HEAD, rebase-merge, rebase-apply, dirty files
│   │   └── git_test.go
│   ├── checkpoint/
│   │   ├── engine.go                   # Snapshot creation, ckpt_<id>.json writer, resume generator
│   │   ├── schema.go                   # Checkpoint schema definition & serializers
│   │   └── checkpoint_test.go
│   ├── ast/
│   │   ├── scanner.go                  # Tree-sitter AST parser for TS/JS, Go, Python, Rust
│   │   ├── symbols.go                  # Exported symbols, route handlers, entrypoints extractor
│   │   └── ast_test.go
│   ├── audit/
│   │   ├── reconnaissance.go           # 4-stage deterministic scan (manifest, structure, symbol, git)
│   │   ├── heuristics.go               # Directory layout heuristics (Prisma, controllers, routes)
│   │   └── audit_test.go
│   ├── linter/
│   │   ├── linter.go                   # CI linter enforcing 150-line activeContext limit & schema
│   │   └── linter_test.go
│   ├── ingest/
│   │   ├── pipeline.go                 # Research inbox reader, triage, diff generator, archive mover
│   │   ├── simhash.go                  # Near-duplicate paragraph detection & deduplication
│   │   ├── parser.go                   # Markdown header & candidate milestone extractor
│   │   └── ingest_test.go
│   ├── llm/
│   │   ├── client.go                   # Pluggable LLM interface (Ollama HTTP / llama.cpp / mock)
│   │   ├── ollama.go                   # Ollama REST client (localhost:11434)
│   │   ├── prompts.go                  # Zero-hallucination structured extraction prompts
│   │   └── llm_test.go
│   ├── mcp/
│   │   ├── server.go                   # MCP server implementation using modelcontextprotocol/go-sdk
│   │   ├── tools.go                    # MCP Tools: read_active_ctx, append_delta, update_milestone, etc.
│   │   ├── resources.go                # MCP Resources: memory://activeContext, memory://progress, etc.
│   │   ├── prompts.go                  # MCP Prompts: "resume-session"
│   │   ├── session.go                  # Agent session read-cursor tracker for zero-token unchanged responses
│   │   └── mcp_test.go
│   ├── rules/
│   │   ├── generator.go                # Minimal vendor stub writer (Cursor, Claude, Copilot, Antigravity)
│   │   └── rules_test.go
│   └── tui/
│       ├── format.go                   # ASCII-safe terminal formatter & ANSI detector
│       ├── status.go                   # Zero-token brief & JSON terminal emitter
│       ├── audit_view.go               # Charmbracelet/Bubbletea live tree audit UI
│       └── tui_test.go
├── pkg/
│   └── types/
│       ├── manifest.go                 # Exported Manifest & FileMeta structs
│       ├── checkpoint.go               # Exported Checkpoint format structs
│       └── memory.go                   # Memory bank domain structs & constants
├── docs/
│   └── phases/
│       ├── backend/
│       │   ├── phase-01-core-manifest-store.md
│       │   ├── phase-02-git-atomic-checkpoint.md
│       │   ├── phase-03-cli-init-status-pause-resume.md
│       │   ├── phase-04-audit-recon-engine.md
│       │   ├── phase-05-lint-memory-ci.md
│       │   ├── phase-06-research-ingestion-tier0.md
│       │   └── phase-07-llm-tier1-ollama.md
│       ├── integration/
│       │   ├── phase-08-mcp-server-protocol.md
│       │   ├── phase-09-vendor-rules-stubs.md
│       │   └── phase-10-cross-project-dashboard.md
│       └── frontend/
│           └── phase-11-terminal-ux-bubbletea.md
├── memory-bank/
│   ├── projectbrief.md                 # Static: Scope, constraints, non-goals
│   ├── productContext.md               # Static: User experience goals, why this exists
│   ├── systemPatterns.md               # Semi-static: Architectural patterns, component relationships
│   ├── techContext.md                  # Semi-static: Tech stack, dependencies, runtime constraints
│   ├── decisionLog.md                  # Semi-static: Append-only ledger of technical decisions
│   ├── activeContext.md                # Hot: Current task focus, recent changes, next steps (<150 lines)
│   ├── progress.md                     # Hot: Milestone checklist, status matrix
│   └── .state/
│       ├── manifest.json               # Source of truth: SHA-256 hashes, sizes, last checkpoints
│       └── checkpoints/                # Full deterministic snapshots (ckpt_<timestamp>.json)
├── research/
│   ├── inbox/                          # Raw brainstorm dumps, notes, RFCs
│   └── archive/                        # Ingested and processed research documents
├── go.mod                              # Go module definition
├── go.sum                              # Cryptographic checksums of dependencies
├── Makefile                            # Build, test, lint, and cross-compilation pipeline
├── MASTER_PLAN.md                      # Operational master plan
└── ORCHESTRATION_BLUEPRINT.md          # This orchestration blueprint
```

---

## 2. Chronological Execution Strategy & Milestone Roadmap

The implementation progresses across five distinct milestones matching the specification:

```text
[ Milestone 1: Core Foundation & Deterministic CLI ]
   ├── Phase 01: Core Manifest & Atomic Store
   ├── Phase 02: Git Porcelain Integration & Safe Checkpointing
   └── Phase 03: Primary CLI Surface (init, status, pause, resume)
          │
          ▼ [ Human Gate 1: CLI Smoke Tests & Atomic Verification ]
[ Milestone 2: MCP Server & Multi-Agent Interop ]
   ├── Phase 08: MCP Server Protocol Implementation
   └── Phase 09: Multi-Vendor Minimal Rule Stub Engine
          │
          ▼ [ Human Gate 2: Agent Tool Invocation & Cursor Sync Verification ]
[ Milestone 3: Brownfield Reconnaissance & CI Enforcement ]
   ├── Phase 04: 4-Stage Reconnaissance Audit Engine
   ├── Phase 05: Memory Bank Budget Linter (ctx lint-memory)
   └── Phase 11: Terminal UX & Bubbletea Interactive Audit TUI
          │
          ▼ [ Human Gate 3: Live Repo Audit & Budget Lint Validation ]
[ Milestone 4: Research Ingestion & Local Intelligence ]
   ├── Phase 06: Tier-0 Research Ingestion Pipeline & SimHash
   └── Phase 07: Tier-1 Local LLM Pipeline (Ollama Integration)
          │
          ▼ [ Human Gate 4: Ingestion Diff-Review & Summarization Verification ]
[ Milestone 5: Workspace Aggregation & Production Hardening ]
   └── Phase 10: Cross-Project Dashboard (`ctx list`) & BYOK Hardening
          │
          ▼ [ Final Human Gate: End-to-End System Audit & Release Packaging ]
```

---

## 3. Detailed Phase Blueprints

---

### Phase 01: Core Manifest, Atomic File Operations & Memory Bank Layout
**File Location:** `/docs/phases/backend/phase-01-core-manifest-store.md`

#### 1. Phase Objective & Boundary
Implement the deterministic storage backbone: atomic file write engine (`.tmp` -> `fsync` -> `rename()`), cryptographic SHA-256 hashing (excluding frontmatter), `manifest.json` state tracking, and the standard Cline-compatible `memory-bank/` directory layout generator. Excludes Git CLI commands and LLM logic.

#### 2. Prerequisites & Dependencies
- Go runtime environment (Go 1.22+).
- Initialized `go.mod` (`module github.com/ctxbank/ctx`).
- Standard library dependencies: `os`, `path/filepath`, `crypto/sha256`, `encoding/json`, `sync`.

#### 3. Step-by-Step Implementation Steps
1. Create `pkg/types/manifest.go` defining `FileMeta` (`SHA256`, `Bytes`, `LastCheckpoint`, `UpdatedAt`) and `Manifest` (`Version`, `Files map[string]FileMeta`, `AgentReadCursors map[string]map[string]string`).
2. Create `internal/core/atomic.go` implementing `WriteAtomic(path string, data []byte, perm os.FileMode) error`.
   - Write to `<path>.tmp.<rand>`.
   - Call `file.Sync()` on disk descriptor.
   - Atomic rename via `os.Rename`.
3. Create `internal/core/manifest.go` providing `ComputeFileHash(path string) (string, int64, error)`, `LoadManifest(bankDir string)`, and `SaveManifest(bankDir string, m *Manifest)`.
4. Create `internal/core/bank.go` defining templates for all standard Cline files (`projectbrief.md`, `productContext.md`, `systemPatterns.md`, `techContext.md`, `decisionLog.md`, `activeContext.md`, `progress.md`).
5. Write comprehensive unit tests in `atomic_test.go` and `manifest_test.go` including a simulated crash/kill test verifying no corrupted or half-written files remain.

#### 4. Concrete Deliverables
- `pkg/types/manifest.go`
- `internal/core/atomic.go` and `internal/core/atomic_test.go`
- `internal/core/manifest.go` and `internal/core/manifest_test.go`
- `internal/core/bank.go` and `internal/core/bank_test.go`

#### 5. Validation & Verification Routine
```bash
go test -v -race -cover ./internal/core/...
```
- Assert atomic write replaces target seamlessly.
- Assert SHA-256 hash updates upon content mutation.
- Assert `manifest.json` persists file sizes and hashes deterministically.

#### 6. Security & Vulnerability Audit
- Verify file permissions are strictly `0644` for markdown/state files and `0755` for directories.
- Sanitize file paths against directory traversal attacks (`filepath.Clean` validation).

---

### Phase 02: Git Porcelain Integration & Safe Checkpoint Engine
**File Location:** `/docs/phases/backend/phase-02-git-atomic-checkpoint.md`

#### 1. Phase Objective & Boundary
Implement read-only system Git integration via shell execution (`status --porcelain`, `diff --stat`, `log -n 1`, `branch --show-current`), repository safety checks (detecting rebase, detached HEAD, merge conflicts), and the deterministic snapshot engine creating `ckpt_<timestamp>.json`. Excludes interactive prompts and MCP protocol.

#### 2. Prerequisites & Dependencies
- Phase 01 complete (`pkg/types`, `internal/core`).
- System `git` available in PATH.

#### 3. Step-by-Step Implementation Steps
1. Create `pkg/types/checkpoint.go` defining `Checkpoint` struct (ID, Timestamp, Branch, CommitHash, DirtyFiles, DiffStat, ActiveFocus, NextSteps, ManualEdits).
2. Create `internal/git/porcelain.go` implementing safe `exec.Command("git", ...)` wrappers:
   - `GetStatus() ([]GitFileStatus, error)`
   - `GetCurrentBranch() (string, error)`
   - `GetHeadCommit() (string, string, error)`
   - `GetDiffStat() (string, error)`
3. Create `internal/git/safeguards.go` implementing:
   - `IsMidRebase(gitDir string) bool` (checks `.git/rebase-merge` and `.git/rebase-apply`).
   - `IsDetachedHead(gitDir string) bool`.
   - `CheckRepositorySafety(repoRoot string) error` (returns descriptive error if git is in unsafe transient state).
4. Create `internal/checkpoint/engine.go`:
   - Generate monotonic checkpoint ID `ckpt_<YYYYMMDD_HHMMSS>`.
   - Capture repo state and dirty files.
   - Write snapshot atomically to `memory-bank/.state/checkpoints/ckpt_<id>.json`.
   - Update `memory-bank/.state/manifest.json` with checkpoint references.

#### 4. Concrete Deliverables
- `pkg/types/checkpoint.go`
- `internal/git/porcelain.go` and `internal/git/safeguards.go`
- `internal/checkpoint/engine.go`
- Unit tests with mock git repositories in `git_test.go` and `checkpoint_test.go`.

#### 5. Validation & Verification Routine
```bash
go test -v -race ./internal/git/... ./internal/checkpoint/...
```
- Simulate mid-rebase state by creating `.git/rebase-merge` folder; assert checkpoint creation is rejected with "resolve rebase first".
- Assert checkpoint snapshot JSON matches schema and contains exact dirty file list and commit hash.

#### 6. Security & Vulnerability Audit
- Never pass unescaped user inputs into shell commands; use parameterized `exec.Command("git", args...)`.
- Ensure untracked sensitive files (e.g., `.env`, credentials) are ignored according to `.gitignore`.

---

### Phase 03: Primary CLI Surface (`init`, `status`, `pause`, `resume`)
**File Location:** `/docs/phases/backend/phase-03-cli-init-status-pause-resume.md`

#### 1. Phase Objective & Boundary
Build the static CLI binary `ctx` with commands: `init`, `status` (sparse ASCII-safe brief or `--json`), `pause` / `checkpoint` (snapshotting + interactive manual edit prompt), and `resume` (pickup brief). No local LLM dependencies required.

#### 2. Prerequisites & Dependencies
- Phase 01 & Phase 02 completed.
- CLI argument parsing (`spf13/cobra` or Go `flag` package).

#### 3. Step-by-Step Implementation Steps
1. Create `cmd/ctx/main.go` and wire command hierarchy:
   - `ctx init`: Scaffolds `memory-bank/`, `.state/`, and invokes Phase 09 rule stub creation.
   - `ctx status [--json]`: Renders single-shot status card (branch, dirty count, test state, idle time, focus, next steps).
   - `ctx pause`: Runs git safety check, triggers checkpoint, detects files changed outside agent sessions, prompts user for explanation (1 line), and updates `activeContext.md` & `decisionLog.md`.
   - `ctx resume [--agent]`: Renders human pickup brief or machine-formatted agent startup context.
2. Implement terminal formatting in `internal/tui/status.go` with auto-detection for TTY color support vs. CI plain ASCII.
3. Add staleness detection: calculate elapsed days and commit delta since the last checkpoint, flagging warnings if stale.

#### 4. Concrete Deliverables
- `cmd/ctx/main.go`
- `internal/tui/format.go` & `internal/tui/status.go`
- Integration tests in `cmd/ctx/cli_test.go`.

#### 5. Validation & Verification Routine
```bash
go build -o ./bin/ctx.exe ./cmd/ctx
./bin/ctx init
./bin/ctx status
./bin/ctx status --json
./bin/ctx pause
./bin/ctx resume
```
- Verify CLI outputs sparse ASCII card without color codes when stdout is redirected to a pipe.
- Verify `--json` output conforms strictly to machine-readable JSON schema.

#### 6. Security & Vulnerability Audit
- Verify stdin handling during `ctx pause` handles EOF, SIGINT, and Ctrl+C gracefully without writing partial corrupt data.
- Ensure binary compiles with zero CGO (`CGO_ENABLED=0`) for a static binary under 15MB.

---

### Phase 04: Brownfield Reconnaissance & AST Analysis Engine
**File Location:** `/docs/phases/backend/phase-04-audit-recon-engine.md`

#### 1. Phase Objective & Boundary
Implement `ctx audit` deterministic brownfield analysis consisting of four strict passes:
1. Manifest scan (`package.json`, `go.mod`, `Cargo.toml`, `pyproject.toml`).
2. Structural directory heuristics (`src/api`, `controllers`, `migrations`, `prisma`).
3. Symbol/AST scan via tree-sitter (exported functions, route handlers, entrypoints).
4. Git velocity scan (commit frequency per path over the last 90 days).
Generates draft `techContext.md` and component map in `systemPatterns.md`. Excludes LLM calls.

#### 2. Prerequisites & Dependencies
- Go tree-sitter bindings (`smacker/go-tree-sitter`).
- Completed Phase 01 and Phase 02.

#### 3. Step-by-Step Implementation Steps
1. Create `internal/ast/scanner.go` configuring tree-sitter grammars for TypeScript/JavaScript, Go, Python, and Rust.
2. Create `internal/ast/symbols.go` extracting exported types, functions, and API entrypoints.
3. Create `internal/audit/heuristics.go` identifying architecture patterns (e.g. Next.js App Router, Express, Gin, Django).
4. Create `internal/audit/reconnaissance.go` orchestrating the 4 passes into an intermediate `AuditResult` struct.
5. Populate draft `techContext.md` and `systemPatterns.md` without fabricating business requirements.

#### 4. Concrete Deliverables
- `internal/ast/scanner.go` & `internal/ast/symbols.go`
- `internal/audit/reconnaissance.go` & `internal/audit/heuristics.go`
- Audit engine test suite in `internal/audit/audit_test.go`.

#### 5. Validation & Verification Routine
```bash
go test -v ./internal/ast/... ./internal/audit/...
./bin/ctx audit --dry-run
```
- Verify scan completes in < 3 seconds on standard codebases.
- Verify detected dependencies match `package.json`/`go.mod`.

#### 6. Security & Vulnerability Audit
- Ensure AST parser is resource-bounded to prevent memory exhaustion on giant or minified files (skip files > 1MB or vendor/node_modules).

---

### Phase 05: Memory Bank Budget Linter (`ctx lint-memory`)
**File Location:** `/docs/phases/backend/phase-05-lint-memory-ci.md`

#### 1. Phase Objective & Boundary
Implement `ctx lint-memory` to enforce strict token and formatting budgets: hard line limit (< 150 lines) on `activeContext.md` and `progress.md`, required markdown header structures, prohibited conversational fluff, and manifest hash integrity. Designed to fail CI or pre-commit hooks upon violation.

#### 2. Prerequisites & Dependencies
- Phase 01 complete (`internal/core`).

#### 3. Step-by-Step Implementation Steps
1. Create `internal/linter/linter.go` defining validation rules:
   - Line count budget: `activeContext.md` <= 150 lines; `progress.md` <= 200 lines.
   - Mandatory sections check for `activeContext.md`: `## Focus`, `## Recent`, `## Next steps`, `## Open decisions`.
   - Bullet-only formatting validation (reject conversational multi-paragraph prose).
   - Manifest drift check: verify actual file SHA-256 matches `manifest.json`.
2. Implement CLI command `ctx lint-memory` returning exit code `0` on pass, exit code `1` with actionable failure list on violation.
3. Add `--fix` flag to automatically trim completed older checkpoint items if requested.

#### 4. Concrete Deliverables
- `internal/linter/linter.go` and `internal/linter/linter_test.go`
- CLI command integration in `cmd/ctx/main.go`.

#### 5. Validation & Verification Routine
```bash
go test -v ./internal/linter/...
./bin/ctx lint-memory
```
- Test against an artificial 200-line `activeContext.md`; verify command fails with exact line count error.
- Test against drifted file; verify manifest mismatch is flagged.

#### 6. Security & Vulnerability Audit
- Ensure linter never mutates files when executed in CI mode (without `--fix`).

---

### Phase 06: Tier-0 Research Ingestion Pipeline & SimHash
**File Location:** `/docs/phases/backend/phase-06-research-ingestion-tier0.md`

#### 1. Phase Objective & Boundary
Build deterministic Tier-0 research ingestion (`ctx ingest <path>`): processes raw markdown dumps from `research/inbox/`, strips boilerplate, performs near-duplicate paragraph deduplication via SimHash, extracts headers as candidate milestones, proposes diffs, and archives ingested files to `research/archive/`. Must never perform silent unconfirmed overwrites.

#### 2. Prerequisites & Dependencies
- Phase 01 complete.

#### 3. Step-by-Step Implementation Steps
1. Create `internal/ingest/simhash.go` implementing 64-bit SimHash with Hamming distance thresholding for text deduplication.
2. Create `internal/ingest/parser.go` to parse markdown headers, action items, and candidate milestones.
3. Create `internal/ingest/pipeline.go`:
   - Read files from `research/inbox/`.
   - Deduplicate incoming text against existing memory-bank files.
   - Format proposed diff using unified diff syntax.
   - Prompt developer for terminal confirmation `[Y/n]`.
   - On confirmation, atomically write changes and move source file to `research/archive/<timestamp>_<filename>`.

#### 4. Concrete Deliverables
- `internal/ingest/simhash.go` & `internal/ingest/simhash_test.go`
- `internal/ingest/parser.go` & `internal/ingest/pipeline.go`
- CLI command `ctx ingest` in `cmd/ctx/main.go`.

#### 5. Validation & Verification Routine
```bash
go test -v ./internal/ingest/...
```
- Verify duplicate paragraphs are detected and pruned.
- Verify source files in `research/inbox/` are preserved if user declines diff.

#### 6. Security & Vulnerability Audit
- Sanitize file names before moving to `research/archive/` to prevent path injection.

---

### Phase 07: Tier-1 Local LLM Integration (Ollama / llama.cpp)
**File Location:** `/docs/phases/backend/phase-07-llm-tier1-ollama.md`

#### 1. Phase Objective & Boundary
Implement optional Tier-1 AI summarization and classification engine using local Ollama (`http://localhost:11434`) or `llama.cpp` server. Classifies ingested research chunks into memory-bank target files and enriches pickup briefs. If Ollama is absent or unconfigured, system degrades gracefully to deterministic Tier-0 with zero errors.

#### 2. Prerequisites & Dependencies
- Phase 06 complete.
- `.context/config.toml` parser.

#### 3. Step-by-Step Implementation Steps
1. Create `internal/llm/client.go` defining interface `LLMClient` (`ClassifyChunk(text string) (TargetFile, error)`, `SummarizeNarrative(data AuditResult) (string, error)`).
2. Create `internal/llm/ollama.go` implementing Ollama REST API client (`/api/generate` and `/api/chat`).
3. Implement connection timeout (2 seconds) and healthcheck to detect local model availability.
4. Implement fallback: if offline or model missing, return explicit `ErrLLMUnavailable`, falling back to deterministic headers without hallucinating.
5. Create structured JSON output prompts enforcing strict response boundaries.

#### 4. Concrete Deliverables
- `internal/llm/client.go`, `internal/llm/ollama.go`, `internal/llm/prompts.go`
- Unit tests with mock HTTP server in `internal/llm/llm_test.go`.

#### 5. Validation & Verification Routine
```bash
go test -v ./internal/llm/...
```
- Verify client gracefully reports status when Ollama is not running.
- Verify JSON parsing recovers from malformed model output.

#### 6. Security & Vulnerability Audit
- Bind Ollama requests strictly to localhost or user-configured endpoint.
- Never log raw prompts containing sensitive code or keys.

---

### Phase 08: Model Context Protocol (MCP) Server
**File Location:** `/docs/phases/integration/phase-08-mcp-server-protocol.md`

#### 1. Phase Objective & Boundary
Implement the full MCP server conforming to the official Model Context Protocol specification over stdio and optional HTTP/SSE. Exposes deterministic tools (`read_active_context`, `read_static_context`, `append_memory_delta`, `update_milestone`, `report_manual_changes`, `request_checkpoint`), resources (`memory://activeContext`, `memory://progress`), and session-based delta caching via `agent_read_cursors`.

#### 2. Prerequisites & Dependencies
- `github.com/mark3labs/mcp-go` or official `modelcontextprotocol/go-sdk`.
- Phases 01, 02, and 05 complete.

#### 3. Step-by-Step Implementation Steps
1. Create `internal/mcp/session.go` managing agent session read cursors stored in `manifest.json`.
2. Create `internal/mcp/tools.go` implementing MCP tools:
   - `read_active_context(session_id)`: returns active context, progress, or `{"unchanged": true}` if caller's session hash matches manifest.
   - `read_static_context(session_id, files)`: hash-gated static/semi-static retrieval.
   - `append_memory_delta(session_id, file, section, content)`: validates line budget, appends under named section.
   - `update_milestone(session_id, milestone, status, next_steps)`: structured update to `progress.md` and `activeContext.md`.
   - `report_manual_changes(session_id, files, description)`: automated checkpointing of mid-task file alterations.
   - `request_checkpoint(session_id)`: full deterministic pause.
3. Create `internal/mcp/resources.go` registering URI handlers for `memory://activeContext`, `memory://progress`, `memory://checkpoint/latest`, and `memory://decisionLog`.
4. Create `internal/mcp/server.go` binding tools, resources, and stdio transport.
5. Create standalone entrypoint `cmd/ctx-mcp/main.go` and CLI flag `ctx serve --mcp`.

#### 4. Concrete Deliverables
- `internal/mcp/server.go`, `internal/mcp/tools.go`, `internal/mcp/resources.go`, `internal/mcp/session.go`
- `cmd/ctx-mcp/main.go`
- MCP protocol test suite in `internal/mcp/mcp_test.go`.

#### 5. Validation & Verification Routine
```bash
go test -v ./internal/mcp/...
```
- Simulate two sequential `read_active_context` calls for same session ID; verify second call returns `{"unchanged": true}`.
- Test `append_memory_delta` with content exceeding line budget; verify rejection with budget error.

#### 6. Security & Vulnerability Audit
- Sanitize file arguments in tool invocations to restrict access strictly within the project's `memory-bank/` directory.

---

### Phase 09: Minimal Multi-Vendor Rule Stub Engine
**File Location:** `/docs/phases/integration/phase-09-vendor-rules-stubs.md`

#### 1. Phase Objective & Boundary
Implement rule stub generator writing canonical, minimal one-line pointers into detected vendor locations (`.cursorrules`, `.cursor/rules/memory-bank.mdc`, `CLAUDE.md`, `.agentrules`, `.github/copilot-instructions.md`). Ensures zero format duplication while guaranteeing agent memory-bank compliance.

#### 2. Prerequisites & Dependencies
- Phase 01 complete.

#### 3. Step-by-Step Implementation Steps
1. Create `internal/rules/generator.go` defining vendor stub templates.
2. Standard stub content:
   ```markdown
   # AGENT INSTRUCTION: MANDATORY MEMORY BANK SYNC
   Before executing tasks or modifying code, read `memory-bank/activeContext.md` and `memory-bank/progress.md` (or invoke MCP tool `read_active_context`).
   Before completing tasks, record updates to `activeContext.md` or invoke `append_memory_delta` / `update_milestone`.
   ```
3. Detect existing vendor files before creating; if existing file contains instructions, append or preserve existing content without clobbering.
4. Integrate with `ctx init` to write stubs automatically upon project initialization.

#### 4. Concrete Deliverables
- `internal/rules/generator.go` and `internal/rules/rules_test.go`

#### 5. Validation & Verification Routine
```bash
go test -v ./internal/rules/...
```
- Verify stubs are created in correct directories.
- Verify existing custom rules are preserved.

#### 6. Security & Vulnerability Audit
- Ensure file writes use atomic safe write routines to avoid zero-byte stub files.

---

### Phase 10: Cross-Project Workspace Dashboard & BYOK Hardening
**File Location:** `/docs/phases/integration/phase-10-cross-project-dashboard.md`

#### 1. Phase Objective & Boundary
Implement cross-project overview command `ctx list` scanning parent directories for active `memory-bank/` installations, displaying aggregated status cards (focus, dirty files, days idle). Implement secure configuration parsing (`.context/config.toml`) resolving API keys strictly from environment variables.

#### 2. Prerequisites & Dependencies
- Phase 03 complete.

#### 3. Step-by-Step Implementation Steps
1. Create `internal/config/config.go` parsing `.context/config.toml` (Ollama base URL, default models, scan paths).
2. Enforce zero plaintext credentials: API keys must be loaded from `os.Getenv` and never written to disk.
3. Implement `ctx list [dir]` command:
   - Walk target directory depth 2.
   - Detect `memory-bank/.state/manifest.json`.
   - Read latest checkpoint and active context.
   - Print consolidated multi-project status table.

#### 4. Concrete Deliverables
- `internal/config/config.go` & `internal/config/config_test.go`
- `cmd/ctx/list.go`

#### 5. Validation & Verification Routine
```bash
go test -v ./internal/config/...
./bin/ctx list ..
```
- Assert no API keys are stored in `config.toml`.
- Verify `ctx list` displays table with accurate project metrics.

#### 6. Security & Vulnerability Audit
- Reject scanning outside user home/workspace roots.
- Mask any sensitive project folder names if configured.

---

### Phase 11: Terminal UX & Interactive Bubbletea TUI
**File Location:** `/docs/phases/frontend/phase-11-terminal-ux-bubbletea.md`

#### 1. Phase Objective & Boundary
Implement the interactive terminal user experience using `charmbracelet/bubbletea` and `lipgloss` for long-running commands (`ctx audit`), while keeping `status`, `pause`, and `resume` sparse, single-shot, pipeable, and ASCII-safe.

#### 2. Prerequisites & Dependencies
- `github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/lipgloss`.
- Phases 03 & 04 complete.

#### 3. Step-by-Step Implementation Steps
1. Create `internal/tui/audit_view.go` using Bubbletea:
   - Live progress indicator across the 4 audit passes.
   - Interactive tree view of discovered modules and symbols.
   - Keyboard navigation (`j`, `k`, `Enter`, `q`).
2. Style terminal cards in `internal/tui/status.go` using Lipgloss with fallback to plain ASCII borders (`+---`, `|`) if non-TTY or on CI.
3. Validate pipeability: ensure `ctx status | grep Focus` works with zero ANSI noise when piped.

#### 4. Concrete Deliverables
- `internal/tui/audit_view.go`
- `internal/tui/format.go` & `internal/tui/status.go`

#### 5. Validation & Verification Routine
```bash
go test -v ./internal/tui/...
```
- Verify pipe redirection outputs clean ASCII text.
- Verify interactive model handles quit signals cleanly.

#### 6. Security & Vulnerability Audit
- Clean terminal state on exit/crash (restore terminal cursor and alternate screen buffer).

---

## 4. Verification & Self-Healing Audit Checklists

For every phase executed by an autonomous coding agent, the following audit checklist must be executed and 100% green before marking the phase complete:

| Checkpoint Category | Verification Command / Metric | Acceptance Threshold |
|---|---|---|
| **Compilation & Build** | `go build ./...` | Zero compilation errors or warnings |
| **Static Analysis & Vet** | `go vet ./...` | Zero vet issues |
| **Automated Test Coverage** | `go test -v -race -cover ./...` | > 85% branch coverage on core packages |
| **Memory / Race Safety** | `go test -race ./internal/...` | Zero race conditions detected |
| **Static Binary Footprint** | `go build -ldflags="-s -w" -o ctx.exe ./cmd/ctx` | Binary size < 15 MB, zero CGO bindings |
| **Atomic Write Integrity** | Run simulated SIGKILL during write test | Zero corrupted or zero-byte target files |
| **Token Budget Compliance** | `ctx lint-memory` | `activeContext.md` < 150 lines, format valid |
| **MCP Protocol Validation** | In-memory stdio test harness | Session caching returns `{"unchanged": true}` |
| **Rebase & Git Safeguard** | Simulated `.git/rebase-merge` test | Refuses snapshot with "resolve rebase first" |
