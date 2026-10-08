# AI DIRECTIVE: Populate CTXbank Memory Bank for `CTXbank`

> **Context:** This project uses **CTXbank** as its deterministic, crash-safe memory layer.
> You have been provided with this directive prompt because a **Prompt Sync** was executed.
> Your mission is to explore the codebase and populate all 7 files in `memory-bank/` with deep, accurate, project-specific technical knowledge.

## 1. Project Reconnaissance Findings
- **Repository Name:** `CTXbank`
- **Root Directory:** `D:\PROJECTS\PASSION PROJECTS\CTXbank`
- **Git Branch:** `main`
- **Head Commit:** `2e845e5b846be67c3764fa0578ebfbc53fb5062b` (feat(ui): add 1-click Sync Context button, automatic uninitialized project detection, and fix property mappings for graph and files)

### Detected Ecosystem & Tech Stack
- **Language:** Go | **Build System:** Go Modules
  - Commands:
    - `test`: `go test -v -race -cover ./...`
    - `build`: `go build ./...`

### Discovered Architecture Components
- **CLI / Application Entrypoints** (`cmd`): Command-line binary targets and entrypoint wiring.
- **CLI / Application Entrypoints** (`cmd/ctx`): Command-line binary targets and entrypoint wiring.
- **CLI / Application Entrypoints** (`cmd/ctx-mcp`): Command-line binary targets and entrypoint wiring.
- **Core Domain Engine** (`internal/core`): Deterministic domain logic, atomic storage, and invariant enforcement.
- **Git Integration Service** (`internal/git`): Porcelain git query abstraction and repository safety checks.
- **Model Context Protocol (MCP) Server** (`internal/mcp`): Stdio and HTTP protocol handlers, tool endpoints, and session delta caching.
- **Terminal Presentation Layer** (`internal/tui`): TTY formatting, status card views, and interactive Bubbletea UI.

### Discovered Symbols & Entrypoints
- `main` (cmd/ctx/main.go:30) [entrypoint]
- `main` (cmd/ctx-mcp/main.go:12) [entrypoint]
- `main` (internal/audit/audit_test.go:60) [entrypoint]
- `HandleFunc /api/status` (internal/ui/server.go:71) [route]
- `HandleFunc /api/files` (internal/ui/server.go:72) [route]
- `HandleFunc /api/file` (internal/ui/server.go:73) [route]
- `HandleFunc /api/file/save` (internal/ui/server.go:74) [route]
- `HandleFunc /api/graph` (internal/ui/server.go:75) [route]
- `HandleFunc /api/checkpoints` (internal/ui/server.go:76) [route]
- `HandleFunc /api/checkpoint/create` (internal/ui/server.go:77) [route]
- `HandleFunc /api/audit` (internal/ui/server.go:78) [route]
- `HandleFunc /api/lint` (internal/ui/server.go:79) [route]

## 2. Memory Bank Structure & Schema Requirements
You MUST update or populate each of the following 7 files in the `memory-bank/` directory. Each file must be concise, factual, and strictly under 150 lines:

### 1. `memory-bank/projectbrief.md` (< 100 lines)
- Core problem statement and purpose of this specific application.
- Primary deliverables and architectural goals.
- Key constraints and explicit non-goals.

### 2. `memory-bank/productContext.md` (< 120 lines)
- Why this product exists from a user and business perspective.
- Target audience, user personas, and primary use cases.
- Core user experience (UX) flows and design expectations.

### 3. `memory-bank/systemPatterns.md` (< 150 lines)
- System architecture diagrams and structural patterns (MVC, Hexagonal, Clean, Event-driven, etc.).
- Key components, boundaries, data flow, and package relationships.
- Critical algorithms or invariants that must never be broken.

### 4. `memory-bank/techContext.md` (< 150 lines)
- Complete technologies used, runtime versions, and core packages.
- Build, test, lint, and run commands.
- Local environment prerequisites and developer configuration.

### 5. `memory-bank/activeContext.md` (CRITICAL HOT CONTEXT — STRICT < 150 lines)
- Current sprint / active development focus.
- Recent milestones or completed features (last 3 items).
- Immediate next 3-5 technical steps.
- Active architectural decisions or dilemmas.

### 6. `memory-bank/progress.md` (< 150 lines)
- Complete roadmap status: what is working vs what is pending.
- Milestone completion checklist.
- Known issues, bugs, or technical debt.

### 7. `memory-bank/decisionLog.md` (< 150 lines)
- Architecture Decision Records (ADRs) explaining key design choices.
- Context, options considered, and rationale for decisions made.

## 3. Strict Operating Rules
1. **No Placeholders:** Remove any remaining generic placeholder text or template stubs. Every line must describe THIS codebase.
2. **Budget Discipline:** No file may exceed 150 lines. Keep descriptions dense, structured, and information-rich.
3. **Atomic Writes:** Directly write/update files inside `memory-bank/`.
4. **Zero Drift:** Ensure consistency across all 7 files.

Proceed now by analyzing the codebase and updating the 7 memory bank files.
