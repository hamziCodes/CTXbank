# System Patterns

## Core Architecture
- **Layer 1: CLI (`cmd/ctx`)**: Static binary for terminal workflows.
- **Layer 2: Core Storage (`internal/core`)**: Atomic writes, SHA-256 manifest engine.
- **Layer 3: Git Porcelain (`internal/git`)**: Non-destructive status and diff analysis.
- **Layer 4: MCP Protocol (`internal/mcp`)**: Stdio-based agent protocol with session delta caching.

## Critical Design Rules
- Atomic write pattern: write to `<file>.tmp.<rand>` -> `fsync()` -> `rename()`.
- Never store API keys in plaintext config files (resolve strictly via environment variables).
