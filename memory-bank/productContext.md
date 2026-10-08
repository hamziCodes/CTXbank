# Product Context — CTXbank

## Why CTXbank Exists
Modern AI coding agents (Cursor, Claude Code, Cline, Antigravity, GitHub Copilot) suffer from memory amnesia across multi-day tasks, sessions, and branch switches. Traditional prompt-engineering approaches rely on massive prompt pasting or loose markdown files that drift over time, blow past context token budgets, and corrupt files during unexpected interrupts.
CTXbank solves this by providing a deterministic, crash-safe, local-first memory layer. It combines cryptographic SHA-256 verification, atomic filesystem operations, AST code reconnaissance, and standard Model Context Protocol (MCP) delta-syncing.

## Target Audience & Personas
- **Autonomous AI Coding Agents:** Cursor, Claude Code, Windsurf, Cline, and Antigravity needing reliable state persistence and token-conscious context briefs (< 150 lines).
- **Solo Developers & Technical Leads:** Engineers switching between multi-repo passion projects who need instantaneous pickup briefs (`ctx resume`) without re-reading commits or notes.
- **Non-Technical & Product Collaborators:** Team members using the embedded Web UI to inspect project architecture, review changes, and track milestones without command-line intimidation.

## User Experience Goals
- **Instantaneous Zero-Friction Setup:** Works immediately with `ctx init` or in 1-click via the Web UI; requires zero cloud credentials or accounts.
- **Extreme Token Discipline:** Delta-caching returns a tiny 20-byte payload (`{"unchanged": true}`) when context has not drifted, saving thousands of tokens per turn.
- **Crash-Safe Operations:** Every write executes atomically via temporary sibling files (`.tmp`), physical disk synchronization (`fsync`), and atomic replacement (`rename`).
- **Prompt Sync Workflow:** Seamless collaboration where CTXbank provides deterministic AST facts and safety snapshots, while the user's coding agent writes domain-specific documentation.
