package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ProjectManifestInfo holds data extracted from dependency and build manifests.
type ProjectManifestInfo struct {
	Language     string            `json:"language"`
	BuildSystem  string            `json:"build_system"`
	Dependencies []string          `json:"dependencies"`
	DevDeps      []string          `json:"dev_dependencies,omitempty"`
	Scripts      map[string]string `json:"scripts,omitempty"`
}

// ScanManifests inspects the repository root for standard build/dependency manifests.
func ScanManifests(repoDir string) []ProjectManifestInfo {
	var results []ProjectManifestInfo

	// 1. Go (go.mod)
	goModPath := filepath.Join(repoDir, "go.mod")
	if data, err := os.ReadFile(goModPath); err == nil {
		info := ProjectManifestInfo{
			Language:    "Go",
			BuildSystem: "Go Modules",
			Scripts: map[string]string{
				"build": "go build ./...",
				"test":  "go test -v -race -cover ./...",
			},
		}
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		inRequire := false
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "require (") {
				inRequire = true
				continue
			}
			if inRequire {
				if line == ")" {
					inRequire = false
					continue
				}
				parts := strings.Fields(line)
				if len(parts) >= 1 {
					info.Dependencies = append(info.Dependencies, parts[0])
				}
			} else if strings.HasPrefix(line, "require ") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					info.Dependencies = append(info.Dependencies, parts[1])
				}
			}
		}
		results = append(results, info)
	}

	// 2. Node.js (package.json)
	pkgJSONPath := filepath.Join(repoDir, "package.json")
	if data, err := os.ReadFile(pkgJSONPath); err == nil {
		var parsed struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
			Scripts         map[string]string `json:"scripts"`
		}
		if err := json.Unmarshal(data, &parsed); err == nil {
			info := ProjectManifestInfo{
				Language:    "JavaScript / TypeScript",
				BuildSystem: "npm / pnpm / yarn",
				Scripts:     parsed.Scripts,
			}
			for dep := range parsed.Dependencies {
				info.Dependencies = append(info.Dependencies, dep)
			}
			for dep := range parsed.DevDependencies {
				info.DevDeps = append(info.DevDeps, dep)
			}
			results = append(results, info)
		}
	}

	// 3. Rust (Cargo.toml)
	cargoPath := filepath.Join(repoDir, "Cargo.toml")
	if data, err := os.ReadFile(cargoPath); err == nil {
		info := ProjectManifestInfo{
			Language:    "Rust",
			BuildSystem: "Cargo",
			Scripts: map[string]string{
				"build": "cargo build",
				"test":  "cargo test",
			},
		}
		lines := strings.Split(string(data), "\n")
		inDeps := false
		for _, l := range lines {
			t := strings.TrimSpace(l)
			if strings.HasPrefix(t, "[dependencies]") {
				inDeps = true
				continue
			} else if strings.HasPrefix(t, "[") {
				inDeps = false
			}
			if inDeps && strings.Contains(t, "=") {
				parts := strings.SplitN(t, "=", 2)
				info.Dependencies = append(info.Dependencies, strings.TrimSpace(parts[0]))
			}
		}
		results = append(results, info)
	}

	// 4. Python (pyproject.toml / requirements.txt)
	reqPath := filepath.Join(repoDir, "requirements.txt")
	if data, err := os.ReadFile(reqPath); err == nil {
		info := ProjectManifestInfo{
			Language:    "Python",
			BuildSystem: "pip",
		}
		lines := strings.Split(string(data), "\n")
		for _, l := range lines {
			t := strings.TrimSpace(l)
			if t != "" && !strings.HasPrefix(t, "#") {
				info.Dependencies = append(info.Dependencies, strings.Split(t, "==")[0])
			}
		}
		results = append(results, info)
	}

	return results
}
