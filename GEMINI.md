# GEMINI.md: Antigravity Workspace Directives

> **Antigravity Workspace Rules Link**  
> See [AGENTS.md](file:///d:/PROJECTS/PASSION%20PROJECTS/CTXbank/AGENTS.md) for the complete, canonical rules and governance specification.

## Core Antigravity Directives:
1. **Memory Bank Consult:** Before modifying code, read `memory-bank/activeContext.md` and `memory-bank/progress.md` (or call MCP tool `read_active_context`).
2. **Quality & Completeness:** Zero empty stubs, zero `// TODO` placeholders, and 100% test-passing verification (`go test -v -race -cover ./...`).
3. **Immutability & Safety:** Atomic file writes only (`.tmp` -> `fsync` -> `rename`). Never modify git working tree destructively.
4. **Token Budget:** Keep `activeContext.md` under 150 lines. No conversational filler.
5. **Phase Logging:** Store all phase definitions and execution records in `/docs/phases/`.
