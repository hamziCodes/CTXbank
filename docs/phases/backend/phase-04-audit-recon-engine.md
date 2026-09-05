# Phase 04: Brownfield Reconnaissance & AST Analysis Engine

**Domain:** Backend / Codebase Reconnaissance  
**Status:** Ready to Implement  
**Specification Source:** `CTXbank.md` (§3, §3.1)  

---

## 1. Phase Objective & Boundary
Implement `ctx audit` deterministic brownfield analysis with 4 sequential passes:
1. **Manifest Scan:** Inspect `package.json`, `go.mod`, `Cargo.toml`, `pyproject.toml` for language, dependencies, and build scripts to seed `techContext.md`.
2. **Structural Scan:** Apply directory heuristics (`src/api`, `prisma/schema.*`, `**/*.controller.ts`, `**/migrations/`) to seed a draft component map in `systemPatterns.md`.
3. **Symbol Scan (Tree-Sitter):** Extract exported functions, classes, entrypoints (`main`, route handlers) to identify candidate core business logic.
4. **Git Velocity Scan:** Calculate commit frequency per path over the last 90 days to identify active subsystems for `activeContext.md`.

*Exclusions:* LLM-generated narrative prose, which is relegated to Tier-1. If no LLM is configured, narrative sections are left as human TODO prompts rather than hallucinated.

---

## 2. Prerequisites & Dependencies
- Phase 01 & Phase 02 complete.
- `smacker/go-tree-sitter` bindings for Go, TypeScript, Python, and Rust.

---

## 3. Step-by-Step Implementation Steps
1. **Implement Manifest Scanner (`internal/audit/manifest.go`):**
   - Parse `package.json` (Node), `go.mod` (Go), `Cargo.toml` (Rust), `pyproject.toml` (Python).
   - Extract primary dependencies, dev dependencies, and execution scripts.
2. **Implement Structural Heuristics (`internal/audit/heuristics.go`):**
   - Detect web frameworks (Next.js, Express, Gin, FastAPI, Actix).
   - Detect ORM / database schemas (Prisma, Gorm, Drizzle, SQLAlchemy).
   - Group files into architectural layers (presentation, business logic, persistence).
3. **Implement Tree-Sitter Symbol Scanner (`internal/ast/scanner.go`):**
   - Initialize language parsers.
   - Extract exported functions, method receivers, interface definitions, and route declarations.
   - Impose size limits: skip minified files or files > 1MB.
4. **Implement Git Velocity Tracker (`internal/git/velocity.go`):**
   - Execute `git log --since="90 days ago" --name-only --format=""`.
   - Aggregate change counts per directory and file to identify hot areas.
5. **Implement Reconnaissance Coordinator (`internal/audit/reconnaissance.go`):**
   - Aggregate outputs into `AuditReport`.
   - Seed `memory-bank/techContext.md` and `memory-bank/systemPatterns.md` atomically.
6. **Write Unit Tests (`internal/audit/audit_test.go`, `internal/ast/ast_test.go`).**

---

## 4. Concrete Deliverables
- `internal/ast/scanner.go` & `internal/ast/symbols.go`
- `internal/audit/manifest.go`, `internal/audit/heuristics.go`, `internal/audit/reconnaissance.go`
- `internal/git/velocity.go`
- Test fixtures and suites in `internal/audit/audit_test.go`.

---

## 5. Validation & Verification Routine
```bash
go test -v ./internal/ast/... ./internal/audit/...
./bin/ctx audit --dry-run
```
- Verify scanning completes in < 3 seconds on standard repositories.
- Verify generated `techContext.md` and `systemPatterns.md` contain accurate, non-hallucinated facts.

---

## 6. Security & Vulnerability Audit
- Resource bounding on AST parsing: prevent memory spikes by enforcing a 10-second timeout and 1MB file threshold.
- Directory traversal protection on repository walking.
