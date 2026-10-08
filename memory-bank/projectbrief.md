# Project Brief — CTXbank

## Overview
CTXbank is a high-performance, deterministic, local-first memory bank and project lifecycle engine for AI-assisted software engineering. Built in idiomatic Go, it provides an immutable state layer, atomic crash-safe file updates, real-time token budget governance (<150 lines per file), Model Context Protocol (MCP) tool integration, brownfield AST reconnaissance, and an embedded modern Web UI designed according to VERTEX Universal Design guidelines.

## Core Deliverables & Capabilities
- **Deterministic CLI (`ctx`):** Commands for initializing (`init`), tracking status (`status`), pausing (`pause`), resuming (`resume`), brownfield AST auditing (`audit`), token budget checking (`lint-memory`), and AI directive prompting (`prompt-sync`).
- **Prompt Sync Engine (`internal/sync`):** Automated safety snapshot creation, codebase AST analysis, compilation of structured AI directive prompts, persistent sync count ledger, and quality verification.
- **Embedded Web Dashboard (`internal/ui`):** Modern VERTEX dashboard with Overview Mission Control, Architecture Dependency Graph, Real-Time Memory Editor, Research Notes Ingestion with SimHash diffs, and Project Switcher / Workspace Ingestion modal.
- **Model Context Protocol (MCP) Server (`internal/mcp`):** Standardized stdio protocol exposing deterministic tools (`read_active_context`, `append_memory_delta`, `update_milestone`, `create_checkpoint`) with SHA-256 caching for zero token waste.
- **Multi-Track Extensions & Installers:** Native VS Code/Cursor Activity Bar extension (`extensions/vscode/`), cross-platform GitHub Actions release matrix, and 1-click install scripts.

## Core Constraints
- Single static binary footprint strictly < 15MB with zero CGO (`CGO_ENABLED=0`).
- Atomic crash-safe writes (`.tmp` -> `fsync` -> `rename`) on all state mutations.
- Strict token discipline: keep hot context files condensed (< 150 lines).

## Explicit Non-Goals
- Cloud-hosted proprietary databases (local-first Git/filesystem design).
- Destructive git commands or unconfirmed file overwrites.
