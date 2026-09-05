# Phase 09: Minimal Multi-Vendor Rule Stub Engine

**Domain:** Integration / Agent Compatibility  
**Status:** Ready to Implement  
**Specification Source:** `CTXbank.md` (§1.1, §1.2)  

---

## 1. Phase Objective & Boundary
Implement automatic generation and maintenance of minimal vendor rule stubs:
- Target formats: `.cursorrules`, `.cursor/rules/memory-bank.mdc`, `CLAUDE.md`, `.agentrules`, `.github/copilot-instructions.md`.
- Keep rules minimal: emit a single canonical directive telling agents to read `memory-bank/activeContext.md` and `memory-bank/progress.md` before doing anything and update them before finishing.
- Never clobber or duplicate large rule files across 5 formats. Preserve pre-existing custom rules.

*Exclusions:* Complex IDE plugin packaging.

---

## 2. Prerequisites & Dependencies
- Phase 01 complete (`internal/core`).

---

## 3. Step-by-Step Implementation Steps
1. **Define Vendor Stub Templates (`internal/rules/templates.go`):**
   - Cursor single file (`.cursorrules`).
   - Cursor MDC (`.cursor/rules/memory-bank.mdc`).
   - Claude Code (`CLAUDE.md`).
   - GitHub Copilot (`.github/copilot-instructions.md`).
   - Generic (`.agentrules`).
2. **Implement Stub Generator (`internal/rules/generator.go`):**
   - Detect existing files.
   - If file exists, check whether memory-bank directive is already present. If missing, append cleanly without destroying user rules.
   - If file is absent, create using `core.WriteAtomic`.
3. **Wire into `ctx init` (`cmd/ctx/cmd_init.go`):**
   - Execute automatically on initialization.
4. **Develop Unit Tests (`internal/rules/rules_test.go`):**
   - Verify non-destructive append to existing files.
   - Verify clean creation of missing files.

---

## 4. Concrete Deliverables
- `internal/rules/templates.go`, `internal/rules/generator.go`
- `internal/rules/rules_test.go`

---

## 5. Validation & Verification Routine
```bash
go test -v ./internal/rules/...
```
- Assert stubs point directly to `memory-bank/activeContext.md`.

---

## 6. Security & Vulnerability Audit
- Atomic writes: ensure partial writes cannot wipe existing developer rules files.
