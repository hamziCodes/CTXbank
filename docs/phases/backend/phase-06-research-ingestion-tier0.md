# Phase 06: Tier-0 Research Ingestion Pipeline & SimHash

**Domain:** Backend / Ingestion Engine  
**Status:** Ready to Implement  
**Specification Source:** `CTXbank.md` (§2.4, §3)  

---

## 1. Phase Objective & Boundary
Implement the deterministic Tier-0 research ingestion pipeline (`ctx ingest <path>`):
- Ingest raw markdown brainstorms and notes from `research/inbox/`.
- Strip web/markdown boilerplate and noise.
- Perform near-duplicate paragraph detection using 64-bit SimHash with Hamming distance thresholding.
- Extract markdown headers as candidate milestones and architectural components.
- Present a unified diff preview to the user in the terminal for mandatory human confirmation (`[Y/n]`).
- On confirmation, atomically update target memory bank files and move source files to `research/archive/`. Never delete or overwrite silently.

*Exclusions:* LLM classification (delegated to Tier-1 in Phase 07).

---

## 2. Prerequisites & Dependencies
- Phase 01 complete (`internal/core`).

---

## 3. Step-by-Step Implementation Steps
1. **Implement SimHash Deduplication (`internal/ingest/simhash.go`):**
   - Tokenize text into 3-word shingles.
   - Hash tokens with FNV-1a into 64-bit hashes.
   - Calculate Hamming distance between incoming paragraphs and existing memory bank paragraphs.
   - Discard paragraphs with Hamming distance <= 3 (near-duplicates).
2. **Implement Markdown Content Parser (`internal/ingest/parser.go`):**
   - Extract headers (`#`, `##`, `###`) and categorize text chunks.
   - Identify candidate tasks (checklists `- [ ]`, `TODO:`).
3. **Implement Diff Generation & Review Loop (`internal/ingest/pipeline.go`):**
   - Compute proposed additions to `memory-bank/productContext.md` or `systemPatterns.md`.
   - Render in-terminal colorized unified diff.
   - Interactively prompt user for confirmation.
4. **Implement Archival Mover (`internal/ingest/archive.go`):**
   - Move ingested file to `research/archive/<YYYYMMDD_HHMMSS>_<filename>`.
5. **Develop Unit Tests (`internal/ingest/ingest_test.go`, `internal/ingest/simhash_test.go`):**
   - Verify deduplication filters duplicated text.
   - Verify archival move preserves file integrity.

---

## 4. Concrete Deliverables
- `internal/ingest/simhash.go` & `internal/ingest/simhash_test.go`
- `internal/ingest/parser.go`, `internal/ingest/pipeline.go`, `internal/ingest/archive.go`
- `cmd/ctx/cmd_ingest.go`

---

## 5. Validation & Verification Routine
```bash
go test -v ./internal/ingest/...
```
- Ingest test file; verify diff is displayed and source file safely archived upon approval.

---

## 6. Security & Vulnerability Audit
- Sanitize target file paths to ensure archived files remain within `research/archive/`.
