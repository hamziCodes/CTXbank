# Decision Log — CTXbank Architecture Decision Records (ADR)

All architectural and design decisions are logged here sequentially. This file is append-only.

## [ADR 001] Cline Memory Bank Schema with Deterministic Engine
- **Decision:** Adopt the standard 7-file Cline Memory Bank schema (`projectbrief.md`, `productContext.md`, `systemPatterns.md`, `techContext.md`, `activeContext.md`, `progress.md`, `decisionLog.md`) backed by a deterministic Go core.
- **Rationale:** Preserves ecosystem compatibility with existing AI agent workflows while eliminating amnesia, race conditions, and unstructured drift.

## [ADR 002] Atomic File Mutations & Crash Safety
- **Decision:** Enforce atomic write semantics (`<file>.tmp.<salt>` &rarr; `fsync()` &rarr; `os.Rename`) on all file mutations.
- **Rationale:** Ensures that sudden process kills, OS crashes, or tool timeouts never leave files in a truncated or corrupted state.

## [ADR 003] VERTEX Universal Design for Embedded Web UI
- **Decision:** Implement embedded web dashboard with VERTEX Universal Design (smooth 12-16px curves, dark/light themes, 4 container states, SVG node-link architecture graph) using zero-dependency HTML/CSS/JS embedded via `embed.FS`.
- **Rationale:** Delivers a premium, state-of-the-art visual experience accessible to non-technical users while preserving the single-static-binary constraint (<15MB).

## [ADR 004] Native VS Code & Cursor Extension Integration
- **Decision:** Package a standalone VS Code extension providing an Activity Bar view, status bar budget meter, and interactive webview bridging to `ctx`.
- **Rationale:** Allows developers to inspect memory and trigger checkpoints directly from their editor without context switching to a browser or terminal.

## [ADR 005] Prompt Sync Two-Tier Model (AST Fact Extraction + AI Agent Directive)
- **Decision:** Implement `ctx prompt-sync` which automatically takes a safety snapshot, extracts AST facts deterministically without external LLM dependencies, and compiles an actionable markdown directive prompt for the developer's AI coding agent.
- **Rationale:** Solves the brownfield bootstrapping problem. Local deterministic AST extracts indisputable facts, while the user's active AI agent populates deep domain context and is verified for compliance via `ctx prompt-sync --verify`.
