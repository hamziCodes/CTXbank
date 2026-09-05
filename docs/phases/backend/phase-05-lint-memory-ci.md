# Phase 05: Memory Bank Budget Linter (`ctx lint-memory`)

**Domain:** Backend / Quality Assurance & CI  
**Status:** Ready to Implement  
**Specification Source:** `CTXbank.md` (§2.2, §3)  

---

## 1. Phase Objective & Boundary
Implement `ctx lint-memory` to enforce strict token budgets and structure across the memory bank:
- Hard line limit: `activeContext.md` must not exceed 150 lines; `progress.md` must not exceed 200 lines.
- Mandatory sections check for `activeContext.md`: `## Focus`, `## Recent`, `## Next steps`, `## Open decisions`.
- Formatting discipline: enforce bullet-only, adjective-sparse notes; reject conversational paragraphs.
- Manifest integrity check: compute live SHA-256 hashes and flag unrecorded manual edits.
- Return exit code `0` on clean pass; return `1` with actionable diagnostics on violation for CI/pre-commit integration.

*Exclusions:* Automated summarization (delegated to Phase 07).

---

## 2. Prerequisites & Dependencies
- Phase 01 complete (`internal/core`).

---

## 3. Step-by-Step Implementation Steps
1. **Define Lint Rule Engine (`internal/linter/rules.go`):**
   - Rule `LineCountLimit`: verify line counts within budget.
   - Rule `RequiredHeaders`: verify required markdown headings exist.
   - Rule `FormatDiscipline`: detect multi-line prose blocks exceeding threshold.
   - Rule `ManifestDrift`: compare live hashes against `manifest.json`.
2. **Implement Linter Runner (`internal/linter/linter.go`):**
   - Execute all rules against `memory-bank/`.
   - Aggregate diagnostic errors and warnings.
   - Support `--fix` flag to automatically prune completed historical checkpoints.
3. **Integrate CLI Command (`cmd/ctx/cmd_lint.go`):**
   - Connect to `ctx lint-memory`.
   - Print human-readable diagnostics or JSON output.
4. **Develop Unit Tests (`internal/linter/linter_test.go`):**
   - Test against compliant memory banks.
   - Test against oversized `activeContext.md` (>150 lines); verify exact failure message and non-zero exit code.

---

## 4. Concrete Deliverables
- `internal/linter/rules.go`, `internal/linter/linter.go`
- `cmd/ctx/cmd_lint.go`
- `internal/linter/linter_test.go`

---

## 5. Validation & Verification Routine
```bash
go test -v ./internal/linter/...
./bin/ctx lint-memory
```
- Assert exit code is `0` when valid, `1` when budget is breached.

---

## 6. Security & Vulnerability Audit
- Read-only execution guarantee: ensure CI runs never mutate files on disk when `--fix` is not specified.
