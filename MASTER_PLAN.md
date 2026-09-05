# MASTER PLAN: CTXbank (The Memory Bank & Project Lifecycle Engine)

**Project:** CTXbank  
**Document Type:** Architecture Master Plan & Autonomous Lifecycle Execution Specification  
**Status:** Under Review — Pending Human Gate 1 Approval  
**Specification Reference:** `CTXbank.md`  

---

## 1. Executive Summary & Problem Formulation

Modern autonomous AI coding agents (Cline, Cursor, Claude Code, Antigravity, Roo Code) suffer from context drift, compaction amnesia, and prompt-adherence degradation. While Cline established the de facto standard human-readable Memory Bank layout (`activeContext.md`, `progress.md`, `systemPatterns.md`, etc.), it remains a **pure prompt-engineering convention**. It lacks:
1. **Deterministic enforcement:** The agent decides when and if it reads memory files.
2. **Cryptographic hashing & delta sync:** Entire multi-kilobyte files are re-read repeatedly, wasting context window tokens.
3. **Crash-safe atomic storage:** Mid-write agent crashes or context cancellations leave truncated files.
4. **Tool/CLI interop:** Developers working outside agent loops cannot check status, pause, or resume.

**CTXbank** solves this not by inventing an incompatible proprietary file standard, but by adopting the Cline file naming schema as a human-readable projection layer, and deploying a high-performance, single-binary Go engine and Model Context Protocol (MCP) server beneath it.

---

## 2. System Architecture & Module Boundaries

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│                                 CLI (cmd/ctx)                                │
│       init   │   status   │   pause   │   resume   │   audit   │   ingest    │
└──────────────────────────────────────┬───────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┴───────────────────────────────────────┐
│                          Core Deterministic Engine                           │
│  ┌────────────────────────┐  ┌─────────────────────────┐  ┌────────────────┐ │
│  │    internal/core       │  │      internal/git       │  │ internal/ast   │ │
│  │ - Atomic fsync/rename  │  │ - Porcelain git wrapper │  │ - Tree-sitter  │ │
│  │ - SHA-256 Manifest     │  │ - Rebase/head guards    │  │ - Symbol scan  │ │
│  │ - Volatility Router    │  │ - Velocity tracker      │  │ - Route parser │ │
│  └────────────────────────┘  └─────────────────────────┘  └────────────────┘ │
└──────────────────────────────────────┬───────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┴───────────────────────────────────────┐
│                          Integration & Service Layer                         │
│  ┌────────────────────────┐  ┌─────────────────────────┐  ┌────────────────┐ │
│  │      internal/mcp      │  │     internal/rules      │  │  internal/llm  │ │
│  │ - MCP Tools & Resources│  │ - Vendor stub generator │  │ - Ollama HTTP  │ │
│  │ - Session Read Cursors │  │ - Minimal rule pointers │  │ - Tier-1 model │ │
│  └────────────────────────┘  └─────────────────────────┘  └────────────────┘ │
└──────────────────────────────────────┬───────────────────────────────────────┘
                                       │
┌──────────────────────────────────────▼───────────────────────────────────────┐
│                               Filesystem State                               │
│  memory-bank/                 .context/                 research/            │
│  ├── *.md                     └── config.toml           ├── inbox/           │
│  └── .state/manifest.json                               └── archive/         │
└──────────────────────────────────────────────────────────────────────────────┘
```

### Module Boundary Guarantees
- **`internal/core`**: Zero external dependencies outside Go standard library. Handles cryptographic hashing, atomic file writing with `.tmp` staging, and manifest state management.
- **`internal/git`**: Shells out to system `git` strictly via `exec.Command` with arguments as string slices. No CGO or `libgit2` links. Never performs destructive write operations (read-only porcelain queries).
- **`internal/mcp`**: Adheres to the official Model Context Protocol. Exposes tools and resources over standard IO (stdio) and optional HTTP. Tracks per-agent-session file hashes to return lightweight `{"unchanged": true}` payloads.
- **`internal/llm`**: Isolated behind a pluggable interface. Defaults to local Ollama (`localhost:11434`) or llama.cpp. If offline, gracefully falls back to deterministic Tier-0 parsing without throwing errors or failing workflows.

---

## 3. Explicit Tech Stack & Runtime Environment Validation

### 3.1 Technology Stack Matrix
| Component | Selected Technology | Rationale & Tradeoffs |
|---|---|---|
| **Core Language** | Go (v1.22+) | Single static binary compilation, cross-platform portability, zero runtime dependencies, fast startup (<10ms). |
| **Git Operations** | System Git via porcelain CLI | Shelling out avoids CGO and `libgit2` cross-compilation hurdles while guaranteeing 100% git correctness. |
| **AST / Symbol Parsing** | `smacker/go-tree-sitter` | Incremental AST parsing for TypeScript, JavaScript, Go, Python, and Rust. |
| **TUI & Terminal Output** | `charmbracelet/bubbletea` + `lipgloss` | Modern interactive TUI for long-running audits; clean ASCII fallback for CI and pipes. |
| **MCP Protocol** | `github.com/mark3labs/mcp-go` | Native Go implementation of the Model Context Protocol stdio transport. |
| **Local LLM Backend** | Ollama HTTP REST API | Local-first, zero cloud credentials required, standard developer adoption. |
| **Configuration** | TOML (`BurntSushi/toml`) | Simple, human-readable config at `.context/config.toml`. API keys loaded strictly from env vars. |

### 3.2 Host Environment Audit (Current System)
- **OS:** Windows 10/11 x86_64
- **Git Version:** `2.51.0.windows.1` (Present & fully supported)
- **Node Version:** `v24.19.0` (Present)
- **Python Version:** `3.12.10` (Present)
- **Go Status:** Not currently installed in PATH.
  - *Resolution:* During execution bootstrap, install Go 1.22+ via:
    ```powershell
    winget install GoLang.Go
    # or download official MSI from https://go.dev/dl/
    ```

---

## 4. UI/UX Theme Integration, Layout & Design Tokens

CTXbank prioritizes an uncluttered, high-information-density developer experience. Terminal output must look pristine in modern terminals (Windows Terminal, iTerm2, Alacritty) while degrading cleanly to pure 7-bit ASCII when piped to agents, CI runners, or files.

### 4.1 Terminal Color Palette & Lipgloss Tokens
```text
Color Tokens:
  Primary Accent:   #7C3AED (Deep Violet / Electric Indigo)
  Secondary Accent: #06B6D4 (Cyan / Clean Highlighting)
  Success State:    #10B981 (Emerald Green)
  Warning / Stale:  #F59E0B (Amber Gold)
  Danger / Dirty:   #EF4444 (Crimson Red)
  Subtle / Muted:   #64748B (Slate Gray)
  Background Tint:  #1E293B (Dark Slate Box Fill)
```

### 4.2 Terminal Layout Specifications

#### Single-Shot Status Box (`ctx status`)
```text
┌─ khata-app · feature/ocr-retry ───────────────────────────────┐
│ ● 2 files dirty   ✗ 1 failing test   ⏱ idle 3d                │
│ Focus: OCR retry backoff                                       │
│ Next:  clamp jitter in retry.ts:44                             │
└───────────────────────────────────────────────────────────────┘
```
- **Terminal Detection Rule:**
  - If `isatty.IsTerminal(os.Stdout.Fd()) == true`: Use UTF-8 rounded box characters (`╭─`, `╰─`, `│`), colored badge indicators (`●` red for dirty, `✓` green for clean).
  - If output is piped or on non-UTF-8 console: Use standard ASCII (`+---`, `|`, `* dirty`).

#### Interactive Audit Progress View (`ctx audit`)
- Powered by `bubbletea`:
  - Step 1: `[✔]` Manifest scan (`package.json`, `go.mod`)
  - Step 2: `[✔]` Directory heuristics (`src/api`, `prisma`)
  - Step 3: `[⠋]` Symbol scan (tree-sitter analyzing functions...)
  - Step 4: `[ ]` Git velocity scan (commit frequency)

---

## 5. End-State Deliverables & Verification Criteria

Upon completing all phases, CTXbank must provide:
1. **Single Static Binary (`ctx.exe` / `ctx`):**
   - Size < 15 MB.
   - Zero dynamic library dependencies (compiled with `CGO_ENABLED=0`).
   - Execution time for `ctx status`: < 15 milliseconds.
2. **Zero-Fluff Memory Bank Generation:**
   - Strict Cline compatibility across `projectbrief.md`, `productContext.md`, `systemPatterns.md`, `techContext.md`, `decisionLog.md`, `activeContext.md`, `progress.md`.
   - `activeContext.md` maintained under 150 lines.
3. **Crash-Safe Operations:**
   - Simulated process kill during file write leaves zero corrupt files or broken manifests.
4. **Token-Saving Delta Synchronization:**
   - MCP `read_active_context` returns `{"unchanged": true}` when file hashes match caller's session cursor.
5. **Multi-Vendor Rule Integration:**
   - Auto-generated pointer stubs in `.cursorrules`, `.cursor/rules/`, `CLAUDE.md`, `.agentrules`, and `.github/copilot-instructions.md`.

---

## 6. Milestone Checkpoints & Human Gatekeeping

| Milestone | Scope | Deliverables | Verification Gate |
|---|---|---|---|
| **M1: Core CLI Foundation** | Phases 01, 02, 03 | `pkg/types`, `internal/core`, `internal/git`, `cmd/ctx` | `ctx init`, `ctx status`, `ctx pause`, `ctx resume` fully operational. Atomic test pass. |
| **M2: MCP & Agent Protocols** | Phases 08, 09 | `internal/mcp`, `internal/rules`, `cmd/ctx-mcp` | MCP stdio tool invocation verified via test harness. Minimal stubs written. |
| **M3: Reconnaissance & CI Lint** | Phases 04, 05, 11 | `internal/ast`, `internal/audit`, `internal/linter`, TUI | `ctx audit` generates system map. `ctx lint-memory` enforces token budget in CI. |
| **M4: Research & Intelligence** | Phases 06, 07 | `internal/ingest`, `internal/llm` | SimHash deduplication, interactive diff confirmation, Ollama classification. |
| **M5: Workspace & Hardening** | Phase 10 | Cross-project `ctx list`, release packaging | Multi-repo scan validated, static binary packaged for release. |

---

## 7. Human Gate 1 Sign-Off Protocol

> [!IMPORTANT]
> In accordance with the Autonomous Project Bootstrapper specifications:
> **Execution is halted at this gate.** The agent will not generate application feature code, phase markdown files, or mutate codebase files until the user reviews this Master Plan and gives explicit approval to proceed with Milestone 1 scaffolding and execution.
