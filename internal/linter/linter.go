package linter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ctxbank/ctx/internal/core"
)

const (
	MaxActiveContextLines = 150
	MaxProgressLines      = 200
)

// LintResult contains the output of the memory bank linting checks.
type LintResult struct {
	Passed     bool     `json:"passed"`
	Violations []string `json:"violations"`
	Warnings   []string `json:"warnings"`
}

// LintMemoryBank inspects memory-bank/ against formatting, size, and manifest integrity rules.
func LintMemoryBank(bankDir string, autoFix bool) (*LintResult, error) {
	res := &LintResult{Passed: true}

	// 1. Validate activeContext.md
	activePath := filepath.Join(bankDir, "activeContext.md")
	data, err := os.ReadFile(activePath)
	if err != nil {
		res.Violations = append(res.Violations, "activeContext.md is missing from memory-bank/")
		res.Passed = false
	} else {
		content := string(data)
		lines := strings.Split(content, "\n")
		lineCount := len(lines)

		// Rule: Line budget <= 150 lines
		if lineCount > MaxActiveContextLines {
			res.Violations = append(res.Violations, fmt.Sprintf("activeContext.md budget exceeded: %d lines (limit: %d lines)", lineCount, MaxActiveContextLines))
			res.Passed = false

			if autoFix {
				// Auto-fix: trim older recent bullets
				trimmed := trimOlderBullets(content, MaxActiveContextLines)
				_ = core.WriteAtomic(activePath, []byte(trimmed), 0644)
				_ = core.UpdateFileMeta(bankDir, "activeContext.md", "")
				res.Warnings = append(res.Warnings, "Auto-fixed activeContext.md by pruning older recent checkpoints.")
			}
		}

		// Rule: Mandatory sections
		requiredSections := []string{"## Focus", "## Recent", "## Next steps", "## Open decisions"}
		for _, sec := range requiredSections {
			if !strings.Contains(content, sec) {
				res.Violations = append(res.Violations, fmt.Sprintf("activeContext.md missing mandatory section: %s", sec))
				res.Passed = false
			}
		}
	}

	// 2. Validate progress.md
	progressPath := filepath.Join(bankDir, "progress.md")
	if progData, err := os.ReadFile(progressPath); err == nil {
		progLines := len(strings.Split(string(progData), "\n"))
		if progLines > MaxProgressLines {
			res.Violations = append(res.Violations, fmt.Sprintf("progress.md budget exceeded: %d lines (limit: %d lines)", progLines, MaxProgressLines))
			res.Passed = false
		}
	}

	// 3. Validate Manifest Drift
	manifest, err := core.LoadManifest(bankDir)
	if err == nil && manifest != nil {
		for filename, meta := range manifest.Files {
			fullPath := filepath.Join(bankDir, filename)
			if actualHash, actualBytes, err := core.ComputeFileHash(fullPath); err == nil {
				if actualHash != meta.SHA256 {
					res.Warnings = append(res.Warnings, fmt.Sprintf("Manifest drift: %s has unrecorded modifications (recorded hash %s vs actual %s)", filename, meta.SHA256[:8], actualHash[:8]))
					if autoFix {
						_ = core.UpdateFileMeta(bankDir, filename, meta.LastCheckpoint)
					}
				}
				_ = actualBytes
			}
		}
	}

	return res, nil
}

func trimOlderBullets(content string, maxLines int) string {
	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines {
		return content
	}

	inRecent := false
	var result []string
	recentBullets := 0

	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "## Recent") {
			inRecent = true
			result = append(result, l)
			continue
		} else if strings.HasPrefix(trimmed, "## ") {
			inRecent = false
		}

		if inRecent && strings.HasPrefix(trimmed, "- ") {
			recentBullets++
			if recentBullets > 3 {
				// Keep only last 3 bullets in Recent section
				continue
			}
		}
		result = append(result, l)
	}

	return strings.Join(result, "\n")
}
