package audit

import (
	"os"
	"path/filepath"
	"strings"
)

// ComponentArchitecture represents discovered architectural tiers.
type ComponentArchitecture struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Files       []string `json:"files,omitempty"`
}

// IgnoredDirectories are skipped during directory scanning.
var IgnoredDirectories = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	"bin":          true,
	"dist":         true,
	"build":        true,
	".context":     true,
	".gemini":      true,
}

// ScanStructuralHeuristics identifies architectural layers and component relationships.
func ScanStructuralHeuristics(repoDir string) []ComponentArchitecture {
	var components []ComponentArchitecture

	_ = filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if IgnoredDirectories[name] {
				return filepath.SkipDir
			}

			relPath, _ := filepath.Rel(repoDir, path)
			normPath := filepath.ToSlash(relPath)

			// Check structural heuristics
			switch {
			case normPath == "cmd" || strings.HasPrefix(normPath, "cmd/"):
				components = append(components, ComponentArchitecture{
					Name:        "CLI / Application Entrypoints",
					Path:        normPath,
					Description: "Command-line binary targets and entrypoint wiring.",
				})
			case normPath == "internal/core" || normPath == "src/core" || normPath == "core":
				components = append(components, ComponentArchitecture{
					Name:        "Core Domain Engine",
					Path:        normPath,
					Description: "Deterministic domain logic, atomic storage, and invariant enforcement.",
				})
			case normPath == "internal/git" || normPath == "src/git":
				components = append(components, ComponentArchitecture{
					Name:        "Git Integration Service",
					Path:        normPath,
					Description: "Porcelain git query abstraction and repository safety checks.",
				})
			case normPath == "internal/mcp" || normPath == "src/mcp":
				components = append(components, ComponentArchitecture{
					Name:        "Model Context Protocol (MCP) Server",
					Path:        normPath,
					Description: "Stdio and HTTP protocol handlers, tool endpoints, and session delta caching.",
				})
			case normPath == "internal/tui" || normPath == "src/tui" || normPath == "src/ui":
				components = append(components, ComponentArchitecture{
					Name:        "Terminal Presentation Layer",
					Path:        normPath,
					Description: "TTY formatting, status card views, and interactive Bubbletea UI.",
				})
			case strings.Contains(normPath, "api") || strings.Contains(normPath, "routes") || strings.Contains(normPath, "controllers"):
				components = append(components, ComponentArchitecture{
					Name:        "API & Routing Layer",
					Path:        normPath,
					Description: "External request routing, validation, and serialization.",
				})
			case strings.Contains(normPath, "models") || strings.Contains(normPath, "schema") || strings.Contains(normPath, "migrations"):
				components = append(components, ComponentArchitecture{
					Name:        "Data Persistence & Schema",
					Path:        normPath,
					Description: "Database entity schemas, migrations, and storage drivers.",
				})
			}
		}
		return nil
	})

	return deduplicateComponents(components)
}

func deduplicateComponents(comps []ComponentArchitecture) []ComponentArchitecture {
	seen := make(map[string]bool)
	var unique []ComponentArchitecture
	for _, c := range comps {
		if !seen[c.Path] {
			seen[c.Path] = true
			unique = append(unique, c)
		}
	}
	return unique
}
