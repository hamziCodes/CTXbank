# Phase 01: Core Manifest, Atomic File Operations & Memory Bank Layout

**Domain:** Backend / Core Storage  
**Status:** Ready to Implement  
**Specification Source:** `CTXbank.md` (§1, §2.1, §2.3, §6.1)  

---

## 1. Phase Objective & Boundary
Implement the deterministic, crash-safe storage layer for CTXbank:
- Atomic file write operations using temporary file staging (`<file>.tmp.<rand>`), file sync (`fsync`), and atomic replacement (`rename`).
- Cryptographic SHA-256 calculation for file contents (excluding frontmatter).
- State persistence and tracking in `memory-bank/.state/manifest.json`.
- Standard Cline-compatible `memory-bank/` directory layout generator with predefined templates.

*Exclusions:* Git CLI porcelain commands, MCP server protocol, AST parsing, and terminal UI.

---

## 2. Prerequisites & Dependencies
- Go 1.22+ runtime.
- Go standard library packages: `os`, `path/filepath`, `crypto/sha256`, `encoding/hex`, `encoding/json`, `sync`, `time`.

---

## 3. Step-by-Step Implementation Steps
1. **Define Manifest Data Models (`pkg/types/manifest.go`):**
   - Define `FileMeta` containing `SHA256 string`, `Bytes int64`, `LastCheckpoint string`, `UpdatedAt time.Time`.
   - Define `Manifest` containing `Version string`, `Files map[string]FileMeta`, `AgentReadCursors map[string]map[string]string`.
2. **Implement Atomic Write Semantics (`internal/core/atomic.go`):**
   - Create `WriteAtomic(path string, data []byte, perm os.FileMode) error`.
   - Generate unique temp file in same directory: `<path>.tmp.<rand>`.
   - Write data and invoke `file.Sync()`.
   - Close descriptor and rename atomically over target destination.
3. **Implement Cryptographic Manifest Engine (`internal/core/manifest.go`):**
   - Create `ComputeFileHash(path string) (string, int64, error)`.
   - Implement frontmatter-stripping hash calculation.
   - Implement `LoadManifest(bankDir string) (*types.Manifest, error)`.
   - Implement `SaveManifest(bankDir string, m *types.Manifest) error` using `WriteAtomic`.
   - Implement `UpdateFileMeta(bankDir string, relPath string, ckptID string) error`.
4. **Implement Memory Bank Generator & Templates (`internal/core/bank.go`):**
   - Generate all 7 Cline memory bank files (`projectbrief.md`, `productContext.md`, `systemPatterns.md`, `techContext.md`, `decisionLog.md`, `activeContext.md`, `progress.md`).
   - Create `.state/` directory and initialize empty `manifest.json`.
   - Categorize files by volatility class (Static, Semi-static, Hot).
5. **Develop Deterministic Unit Tests:**
   - `internal/core/atomic_test.go`: test concurrency, overwrite safety, error propagation on permission denial.
   - `internal/core/manifest_test.go`: test hashing accuracy, manifest serialization, delta detection.
   - `internal/core/bank_test.go`: verify full initialization scaffold and template accuracy.

---

## 4. Concrete Deliverables
- `pkg/types/manifest.go`
- `internal/core/atomic.go` & `internal/core/atomic_test.go`
- `internal/core/manifest.go` & `internal/core/manifest_test.go`
- `internal/core/bank.go` & `internal/core/bank_test.go`

---

## 5. Validation & Verification Routine
```bash
go test -v -race -cover ./pkg/types/... ./internal/core/...
```
- Verify 100% test passing rate.
- Coverage threshold: > 85% on `internal/core`.
- Simulated crash test: ensure interrupted writes leave zero corrupted target files.

---

## 6. Security & Vulnerability Audit
- Sanitize all file paths with `filepath.Clean` and check for directory traversal attempts (`..`).
- Restrict file permissions strictly to `0644` for markdown/JSON and `0755` for directories.
- Ensure temporary files are cleaned up on write failure via deferred `os.Remove`.
