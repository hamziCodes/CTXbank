# Contributing to CTXbank

First off, thank you for considering contributing to CTXbank! We are building a deterministic, crash-safe memory layer for AI coding agents and human developers.

---

## 1. Core Engineering Principles

Before writing code, please review our core architectural tenets:
1. **Crash-Safety & Atomic Writes:** Never write directly to target files. Always write to a temporary file (`.tmp`), flush to storage (`fsync`), and perform an atomic rename.
2. **Deterministic & Token-Efficient:** We prioritize zero-token waste. Read operations must support delta caching and return `{"unchanged": true}` if files have not drifted.
3. **Zero Stubs or Placeholders:** Every method or function must be completely implemented with error handling and unit test coverage. No `// TODO` or empty placeholders in production code.
4. **Line Budgets:** `memory-bank/activeContext.md` has a hard budget of **< 150 lines**. Keep hot files concise.
5. **No Broken Builds:** Every PR must pass all unit tests (`go test -v -race ./...`) and lint checks before review.

---

## 2. Development Setup

### Prerequisites
- **Go 1.22+** installed on your system.
- **Git** (version 2.30+).
- *(Optional)* **Ollama** running locally at `http://localhost:11434` for testing local LLM features.

### Building Locally
```bash
# Clone your fork
git clone https://github.com/<your-username>/CTXbank.git
cd CTXbank

# Run tests
go test -v -race ./...

# Build binaries
go build -ldflags="-s -w" -o ./bin/ctx.exe ./cmd/ctx
go build -ldflags="-s -w" -o ./bin/ctx-mcp.exe ./cmd/ctx-mcp
```

---

## 3. Project Structure

```text
├── cmd/
│   ├── ctx/          # Main CLI entrypoint (init, status, pause, resume, audit, ingest, list)
│   └── ctx-mcp/      # Standalone Model Context Protocol (MCP) server
├── internal/
│   ├── audit/        # 4-stage brownfield reconnaissance engine
│   ├── checkpoint/   # Deterministic state snapshots
│   ├── config/       # BYOK configuration (.context/config.toml)
│   ├── core/         # Atomic file writes, manifest SHA-256 tracking, Cline scaffold
│   ├── git/          # Porcelain git integration and safety guards
│   ├── ingest/       # SimHash deduplication & research note parser
│   ├── linter/       # Memory bank token budget validator
│   ├── llm/          # Local Ollama REST client with offline fallback
│   ├── mcp/          # JSON-RPC 2.0 stdio protocol handler
│   ├── rules/        # Multi-vendor rule stub generator
│   ├── tui/          # Terminal card rendering & ASCII fallbacks
│   └── workspace/    # Cross-project workspace scanner
├── memory-bank/      # Self-hosting memory bank for CTXbank
└── pkg/types/        # Canonical structs (Manifest, Checkpoint, etc.)
```

---

## 4. How to Submit a Contribution

### Step 1: Branching
Create a feature branch from `main`:
```bash
git checkout -b feat/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### Step 2: Implement & Test
- Write tests alongside your feature in `*_test.go`.
- Ensure all tests pass:
  ```bash
  go test -v -race ./...
  ```
- Check memory bank line budgets:
  ```bash
  ./bin/ctx.exe lint-memory
  ```

### Step 3: Commit Messages
Write clear, conventional commit messages:
- `feat: add support for custom SimHash shingle sizes`
- `fix: resolve windows path normalization in workspace scanner`
- `docs: improve quickstart instructions in README`

### Step 4: Open a Pull Request
Push your branch to your fork and open a Pull Request against `main`. Fill out the Pull Request template completely.

---

## 5. Community Guidelines
Be respectful, constructive, and helpful to fellow contributors. We welcome improvements to documentation, test coverage, bug fixes, and new features!
