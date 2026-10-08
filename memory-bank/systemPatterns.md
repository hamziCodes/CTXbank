# System Patterns — generated from deterministic audit

## Architectural Component Map
### CLI / Application Entrypoints (`cmd`)
Command-line binary targets and entrypoint wiring.

### CLI / Application Entrypoints (`cmd/ctx`)
Command-line binary targets and entrypoint wiring.

### CLI / Application Entrypoints (`cmd/ctx-mcp`)
Command-line binary targets and entrypoint wiring.

### Core Domain Engine (`internal/core`)
Deterministic domain logic, atomic storage, and invariant enforcement.

### Git Integration Service (`internal/git`)
Porcelain git query abstraction and repository safety checks.

### Model Context Protocol (MCP) Server (`internal/mcp`)
Stdio and HTTP protocol handlers, tool endpoints, and session delta caching.

### Terminal Presentation Layer (`internal/tui`)
TTY formatting, status card views, and interactive Bubbletea UI.

## Core Entrypoints & Symbols
- **entrypoint** (`cmd/ctx/main.go`:30) — `main`
- **entrypoint** (`cmd/ctx-mcp/main.go`:12) — `main`
- **entrypoint** (`internal/audit/audit_test.go`:60) — `main`
- **route** (`internal/ui/server.go`:71) — `HandleFunc /api/status`
- **route** (`internal/ui/server.go`:72) — `HandleFunc /api/files`
- **route** (`internal/ui/server.go`:73) — `HandleFunc /api/file`
- **route** (`internal/ui/server.go`:74) — `HandleFunc /api/file/save`
- **route** (`internal/ui/server.go`:75) — `HandleFunc /api/graph`
- **route** (`internal/ui/server.go`:76) — `HandleFunc /api/checkpoints`
- **route** (`internal/ui/server.go`:77) — `HandleFunc /api/checkpoint/create`
- **route** (`internal/ui/server.go`:78) — `HandleFunc /api/audit`
- **route** (`internal/ui/server.go`:79) — `HandleFunc /api/lint`
- **route** (`internal/ui/server.go`:80) — `HandleFunc /api/ingest/preview`
- **route** (`internal/ui/server.go`:81) — `HandleFunc /api/ingest/commit`
- **route** (`internal/ui/server.go`:82) — `HandleFunc /api/projects`
- **route** (`internal/ui/server.go`:83) — `HandleFunc /api/project/inspect`
- **route** (`internal/ui/server.go`:84) — `HandleFunc /api/project/switch`
- **route** (`internal/ui/server.go`:85) — `HandleFunc /api/project/init`
- **route** (`internal/ui/server.go`:86) — `HandleFunc /api/sync`
- **route** (`internal/ui/server.go`:87) — `HandleFunc /api/prompt-sync`
- **route** (`internal/ui/server.go`:88) — `HandleFunc /api/prompt-sync/status`
- **route** (`internal/ui/server.go`:89) — `HandleFunc /api/prompt-sync/verify`
- **route** (`internal/ui/server.go`:103) — `Handle /`
