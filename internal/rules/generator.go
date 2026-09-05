package rules

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ctxbank/ctx/internal/core"
)

const CanonicalDirective = `# AGENT INSTRUCTION: MANDATORY MEMORY BANK SYNC
Before executing tasks or modifying code, read memory-bank/activeContext.md and memory-bank/progress.md (or invoke MCP tool read_active_context).
Before completing tasks, record updates to activeContext.md or invoke append_memory_delta / update_milestone.
`

// GenerateStubs writes minimal pointer stubs to standard vendor rule locations.
func GenerateStubs(workspaceDir string) ([]string, error) {
	targets := []string{
		".cursorrules",
		".agentrules",
		"CLAUDE.md",
		filepath.Join(".cursor", "rules", "memory-bank.mdc"),
		filepath.Join(".github", "copilot-instructions.md"),
	}

	var generated []string
	for _, relPath := range targets {
		targetPath := filepath.Join(workspaceDir, relPath)
		if existing, err := os.ReadFile(targetPath); err == nil {
			// If file already exists, only append if directive is missing
			if !strings.Contains(string(existing), "memory-bank/activeContext.md") {
				updated := strings.TrimRight(string(existing), "\r\n") + "\n\n" + CanonicalDirective
				if err := core.WriteAtomic(targetPath, []byte(updated), 0644); err == nil {
					generated = append(generated, relPath)
				}
			}
		} else {
			// Create file cleanly
			if err := core.WriteAtomic(targetPath, []byte(CanonicalDirective), 0644); err == nil {
				generated = append(generated, relPath)
			}
		}
	}
	return generated, nil
}
