package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ctxbank/ctx/internal/core"
	"github.com/ctxbank/ctx/internal/git"
)

// ProjectSummary represents status metrics for an initialized CTXbank workspace.
type ProjectSummary struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Branch      string    `json:"branch"`
	DirtyCount  int       `json:"dirty_count"`
	ActiveFocus string    `json:"active_focus"`
	IdleDays    int       `json:"idle_days"`
	LastUpdated time.Time `json:"last_updated"`
}

// ScanWorkspaces discovers all repositories within rootPath containing an active memory-bank.
func ScanWorkspaces(rootPath string, maxDepth int) ([]ProjectSummary, error) {
	var projects []ProjectSummary
	cleanRoot := filepath.Clean(rootPath)

	_ = filepath.Walk(cleanRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || name == "bin" {
				return filepath.SkipDir
			}

			// Check depth
			rel, _ := filepath.Rel(cleanRoot, path)
			depth := len(strings.Split(filepath.ToSlash(rel), "/"))
			if depth > maxDepth && rel != "." {
				return filepath.SkipDir
			}

			// Check for memory-bank directory
			manifestPath := filepath.Join(path, core.MemoryBankDir, core.StateDirname, core.ManifestFilename)
			if _, err := os.Stat(manifestPath); err == nil {
				// Found a managed CTXbank project!
				proj := inspectProject(path)
				projects = append(projects, proj)
				return filepath.SkipDir // Do not descend into memory-bank subdirs
			}
		}
		return nil
	})

	return projects, nil
}

func inspectProject(repoPath string) ProjectSummary {
	name := filepath.Base(repoPath)
	bankDir := filepath.Join(repoPath, core.MemoryBankDir)
	branch, _ := git.GetCurrentBranch(repoPath)
	dirty, _ := git.GetStatus(repoPath)

	activeFocus := "No focus declared"
	idleDays := 0
	lastUpdated := time.Now().UTC()

	activeContextPath := filepath.Join(bankDir, "activeContext.md")
	if info, err := os.Stat(activeContextPath); err == nil {
		lastUpdated = info.ModTime()
		idleDays = int(time.Since(info.ModTime()).Hours() / 24)

		if data, err := os.ReadFile(activeContextPath); err == nil {
			lines := strings.Split(string(data), "\n")
			inFocus := false
			for _, l := range lines {
				trimmed := strings.TrimSpace(l)
				if strings.HasPrefix(trimmed, "## Focus") {
					inFocus = true
					continue
				} else if strings.HasPrefix(trimmed, "## ") {
					inFocus = false
				}
				if inFocus && trimmed != "" && !strings.HasPrefix(trimmed, "#") {
					activeFocus = trimmed
					break
				}
			}
		}
	}

	return ProjectSummary{
		Name:        name,
		Path:        repoPath,
		Branch:      branch,
		DirtyCount:  len(dirty),
		ActiveFocus: activeFocus,
		IdleDays:    idleDays,
		LastUpdated: lastUpdated,
	}
}
