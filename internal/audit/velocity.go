package audit

import (
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// PathVelocity tracks change frequency for a given path.
type PathVelocity struct {
	Path        string `json:"path"`
	CommitCount int    `json:"commit_count"`
}

// ScanGitVelocity analyzes commit history over the last 90 days to discover hot subsystems.
func ScanGitVelocity(repoDir string) ([]PathVelocity, error) {
	cmd := exec.Command("git", "log", "--since=90 days ago", "--name-only", "--format=")
	cmd.Dir = repoDir

	out, err := cmd.Output()
	if err != nil {
		// If git log fails (e.g. no commits), return empty list gracefully
		return nil, nil
	}

	frequency := make(map[string]int)
	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			continue
		}
		norm := filepath.ToSlash(trimmed)
		frequency[norm]++

		// Also track parent directory
		dir := filepath.ToSlash(filepath.Dir(norm))
		if dir != "." && dir != "" {
			frequency[dir]++
		}
	}

	var results []PathVelocity
	for p, count := range frequency {
		results = append(results, PathVelocity{Path: p, CommitCount: count})
	}

	// Sort descending by commit count
	sort.Slice(results, func(i, j int) bool {
		return results[i].CommitCount > results[j].CommitCount
	})

	if len(results) > 15 {
		results = results[:15]
	}

	return results, nil
}
