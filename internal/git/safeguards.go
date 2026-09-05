package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrMidRebase     = errors.New("repository is currently mid-rebase; resolve or abort rebase before checkpointing")
	ErrDetachedHead  = errors.New("repository is in detached HEAD state; switch to a named branch before checkpointing")
	ErrUnmergedFiles = errors.New("repository contains unresolved merge conflicts")
)

// IsMidRebase checks if git is in the middle of an interactive or standard rebase.
func IsMidRebase(gitDir string) bool {
	rebaseMerge := filepath.Join(gitDir, "rebase-merge")
	if info, err := os.Stat(rebaseMerge); err == nil && info.IsDir() {
		return true
	}
	rebaseApply := filepath.Join(gitDir, "rebase-apply")
	if info, err := os.Stat(rebaseApply); err == nil && info.IsDir() {
		return true
	}
	return false
}

// ValidateRepoSafety asserts that the repository is in a safe, non-transient state.
func ValidateRepoSafety(repoDir string) error {
	gitDir, err := GetGitDir(repoDir)
	if err != nil {
		// Not a git repository, safe to proceed in non-git mode
		return nil
	}

	// Resolve relative gitDir
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(repoDir, gitDir)
	}

	if IsMidRebase(gitDir) {
		return ErrMidRebase
	}

	branch, err := GetCurrentBranch(repoDir)
	if err == nil && branch == "(detached)" {
		return ErrDetachedHead
	}

	statuses, err := GetStatus(repoDir)
	if err == nil {
		for _, s := range statuses {
			if (s.StagingStatus == "U" || s.WorktreeStatus == "U") ||
				(s.StagingStatus == "A" && s.WorktreeStatus == "A") ||
				(s.StagingStatus == "D" && s.WorktreeStatus == "D") {
				return fmt.Errorf("%w: conflict detected in %s", ErrUnmergedFiles, s.Path)
			}
		}
	}

	return nil
}
