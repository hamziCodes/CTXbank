# Phase 02: Git Porcelain Integration & Safe Checkpoint Engine

**Domain:** Backend / Git & Snapshots  
**Status:** Ready to Implement  
**Specification Source:** `CTXbank.md` (§1, §2.1, §3, §6.1)  

---

## 1. Phase Objective & Boundary
Implement safe, read-only Git integration and the deterministic snapshot engine:
- Shell out to system `git` strictly via porcelain commands (`status --porcelain`, `diff --stat`, `log -n 1`, `branch --show-current`).
- Safeguard checks: detect detached HEAD, active rebase (`.git/rebase-merge` or `.git/rebase-apply`), and merge conflicts.
- Checkpoint creation engine: produce deterministic `ckpt_<YYYYMMDD_HHMMSS>.json` snapshots capturing repository dirty state, branch, commit, and diff statistics.

*Exclusions:* Git write mutations (`commit`, `checkout`, `rebase`), MCP endpoints, and interactive CLI prompts.

---

## 2. Prerequisites & Dependencies
- Phase 01 complete (`pkg/types`, `internal/core`).
- System Git (`git` command) available in PATH.

---

## 3. Step-by-Step Implementation Steps
1. **Define Checkpoint Data Models (`pkg/types/checkpoint.go`):**
   - Define `GitFileStatus` (`Path string`, `StagingStatus string`, `WorktreeStatus string`).
   - Define `Checkpoint` (`ID string`, `Timestamp time.Time`, `Branch string`, `HeadCommit string`, `CommitMessage string`, `DirtyFiles []GitFileStatus`, `DiffStat string`, `ActiveFocus string`, `NextSteps []string`, `ManualEdits []string`).
2. **Implement Porcelain Git Wrapper (`internal/git/porcelain.go`):**
   - Implement `GetStatus(repoDir string) ([]types.GitFileStatus, error)`.
   - Implement `GetCurrentBranch(repoDir string) (string, error)`.
   - Implement `GetHeadCommit(repoDir string) (hash string, message string, err error)`.
   - Implement `GetDiffStat(repoDir string) (string, error)`.
   - Ensure commands are invoked with structured argument slices (no shell interpolation).
3. **Implement Safety Guards (`internal/git/safeguards.go`):**
   - Implement `IsMidRebase(gitDir string) bool`.
   - Implement `IsDetachedHead(gitDir string) bool`.
   - Implement `ValidateRepoSafety(repoDir string) error` (returns descriptive error preventing checkpointing in unsafe states).
4. **Implement Checkpoint Engine (`internal/checkpoint/engine.go`):**
   - Implement `GenerateCheckpointID(t time.Time) string`.
   - Implement `CreateCheckpoint(bankDir, repoDir string, focus string, nextSteps []string, manualEdits []string) (*types.Checkpoint, error)`.
   - Store checkpoint in `memory-bank/.state/checkpoints/ckpt_<id>.json` using `core.WriteAtomic`.
   - Update `manifest.json` with checkpoint references.
5. **Develop Unit & Integration Tests (`internal/git/git_test.go`, `internal/checkpoint/checkpoint_test.go`):**
   - Create temporary git repositories using `git init`.
   - Test porcelain queries against clean and dirty working trees.
   - Simulate `.git/rebase-merge` directory; assert checkpoint aborts with clear error message.

---

## 4. Concrete Deliverables
- `pkg/types/checkpoint.go`
- `internal/git/porcelain.go` & `internal/git/safeguards.go` & `internal/git/git_test.go`
- `internal/checkpoint/engine.go` & `internal/checkpoint/checkpoint_test.go`

---

## 5. Validation & Verification Routine
```bash
go test -v -race ./internal/git/... ./internal/checkpoint/...
```
- Assert dirty files match `git status --porcelain`.
- Assert rebase detection prevents snapshots.
- Verify snapshot JSON is valid and loadable.

---

## 6. Security & Vulnerability Audit
- Zero shell injection risk: strictly use `exec.Command("git", arg1, arg2...)` without shell interpreters.
- Exclude tracked credentials and files matching `.gitignore`.
