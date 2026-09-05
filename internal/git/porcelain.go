package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ctxbank/ctx/pkg/types"
)

// Runner abstracts git command execution for testing.
type Runner interface {
	Run(dir string, args ...string) (string, error)
}

// SystemRunner executes git commands by shelling out to the system binary.
type SystemRunner struct{}

func (s *SystemRunner) Run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s failed: %w (stderr: %s)", strings.Join(args, " "), err, stderr.String())
	}
	return stdout.String(), nil
}

var defaultRunner Runner = &SystemRunner{}

// SetRunner allows substituting a mock git runner during unit tests.
func SetRunner(r Runner) {
	defaultRunner = r
}

// ResetRunner restores the default system git runner.
func ResetRunner() {
	defaultRunner = &SystemRunner{}
}

// GetStatus executes `git status --porcelain` and parses dirty/untracked files.
func GetStatus(repoDir string) ([]types.GitFileStatus, error) {
	out, err := defaultRunner.Run(repoDir, "status", "--porcelain")
	if err != nil {
		return nil, err
	}

	var results []types.GitFileStatus
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		staged := string(line[0])
		worktree := string(line[1])
		filePath := strings.TrimSpace(line[3:])

		results = append(results, types.GitFileStatus{
			Path:           filePath,
			StagingStatus:  staged,
			WorktreeStatus: worktree,
		})
	}
	return results, nil
}

// GetCurrentBranch returns the active git branch or "(detached)" if detached HEAD.
func GetCurrentBranch(repoDir string) (string, error) {
	out, err := defaultRunner.Run(repoDir, "branch", "--show-current")
	if err != nil {
		return "", err
	}
	branch := strings.TrimSpace(out)
	if branch == "" {
		// Possibly detached HEAD
		return "(detached)", nil
	}
	return branch, nil
}

// GetHeadCommit returns the SHA hash and commit subject of HEAD.
func GetHeadCommit(repoDir string) (string, string, error) {
	out, err := defaultRunner.Run(repoDir, "log", "-n", "1", "--format=%H%x00%s")
	if err != nil {
		// New repository without commits
		return "", "initial commit", nil
	}

	parts := strings.Split(strings.TrimSpace(out), "\x00")
	if len(parts) >= 2 {
		return parts[0], parts[1], nil
	} else if len(parts) == 1 && parts[0] != "" {
		return parts[0], "", nil
	}
	return "", "", nil
}

// GetDiffStat returns the summary output of `git diff --stat`.
func GetDiffStat(repoDir string) (string, error) {
	out, err := defaultRunner.Run(repoDir, "diff", "--stat")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// GetGitDir returns the root .git directory path.
func GetGitDir(repoDir string) (string, error) {
	out, err := defaultRunner.Run(repoDir, "rev-parse", "--git-dir")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
