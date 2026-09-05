# Phase 10: Cross-Project Workspace Dashboard & BYOK Hardening

**Domain:** Integration / Multi-Project Orchestration  
**Status:** Ready to Implement  
**Specification Source:** `CTXbank.md` (§5, §6.2)  

---

## 1. Phase Objective & Boundary
Implement cross-project workspace aggregation (`ctx list [dir]`) and secure configuration management:
- Recursively discover projects containing initialized `memory-bank/` folders.
- Render aggregated multi-project status cards (dirty file counts, idle time, active focus).
- Parse configuration from `.context/config.toml` with strict Bring-Your-Own-Key (BYOK) environment resolution (never storing plaintext keys on disk).

*Exclusions:* Cloud-synced multi-tenant databases.

---

## 2. Prerequisites & Dependencies
- Phases 01, 02, and 03 complete.
- `BurntSushi/toml` parser.

---

## 3. Step-by-Step Implementation Steps
1. **Implement Config Loader (`internal/config/config.go`):**
   - Read `.context/config.toml`.
   - Resolve keys dynamically via `os.Getenv("OPENAI_API_KEY")`, `os.Getenv("ANTHROPIC_API_KEY")`.
   - Prohibit saving key strings to disk.
2. **Implement Workspace Scanner (`internal/workspace/scanner.go`):**
   - Scan directory tree up to depth 3.
   - Detect `memory-bank/.state/manifest.json`.
   - Extract project name, branch, active focus, and days idle.
3. **Implement CLI Command `ctx list` (`cmd/ctx/cmd_list.go`):**
   - Render multi-project overview table in terminal.
4. **Develop Unit Tests (`internal/config/config_test.go`, `internal/workspace/scanner_test.go`):**
   - Assert env resolution works properly.
   - Assert scan correctly enumerates multiple workspaces.

---

## 4. Concrete Deliverables
- `internal/config/config.go` & `internal/config/config_test.go`
- `internal/workspace/scanner.go`
- `cmd/ctx/cmd_list.go`

---

## 5. Validation & Verification Routine
```bash
go test -v ./internal/config/... ./internal/workspace/...
./bin/ctx list .
```

---

## 6. Security & Vulnerability Audit
- Zero credential leakage: audit config serializers to verify credentials cannot be marshaled into `.context/config.toml`.
