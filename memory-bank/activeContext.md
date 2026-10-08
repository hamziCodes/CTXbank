# Active Context — prompt_sync_engine_and_dropdown_fix

## Focus
Implemented the Prompt Sync Engine (automated pre-sync snapshots, AST code reconnaissance, structured AI directive prompt compilation, persistent sync counter ledger, and quality verification) and resolved the Project Selector dropdown modal hierarchy.

## Recent (last 3 checkpoints)
- Prompt Sync Engine (`internal/sync/`): Automatic pre-sync snapshot capture, brownfield AST reconnaissance, AI directive compilation to `.context/prompts/agent_sync_prompt.md`, sync ledger persistence, and quality verification.
- CLI & Web API Integration: Added `ctx prompt-sync [--verify|--status]` commands and REST endpoints (`/api/prompt-sync`, `/api/prompt-sync/status`, `/api/prompt-sync/verify`).
- Web UI & Dropdown Fix: Fixed modal HTML nesting bug so the Project Selector dropdown (`📁 CTXbank ▾`) opens cleanly, added header Prompt Sync button with real-time spin animation and counter badge, and integrated interactive prompt modal with 1-click clipboard copy.

## Next steps
1. Restart `ctx ui` dev server to bind newly compiled `/api/prompt-sync` Go REST handlers.
2. Provide generated directive prompts to coding agents when tapping into fresh brownfield repositories.
3. Track and verify agent-populated memory banks with `ctx prompt-sync --verify`.

## Open decisions
- None. All unit and API tests pass; static binary is verified at 11.33 MB (< 15MB limit).
