package audit

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DiscoveredSymbol represents an entrypoint or exported business logic identifier.
type DiscoveredSymbol struct {
	File string `json:"file"`
	Kind string `json:"kind"` // entrypoint, exported_func, struct_type, route
	Name string `json:"name"`
	Line int    `json:"line"`
}

var (
	reGoFunc    = regexp.MustCompile(`^func\s+([A-Z][a-zA-Z0-9_]*)\s*\(`)
	reGoMain    = regexp.MustCompile(`^func\s+main\s*\(`)
	reGoType    = regexp.MustCompile(`^type\s+([A-Z][a-zA-Z0-9_]*)\s+(struct|interface)`)
	reGoRoute   = regexp.MustCompile(`\.(GET|POST|PUT|DELETE|Handle|HandleFunc)\s*\(\s*["']([^"']+)["']`)
	reTsExport  = regexp.MustCompile(`^export\s+(function|class|interface|type)\s+([a-zA-Z0-9_]+)`)
	rePyDef     = regexp.MustCompile(`^def\s+([a-zA-Z0-9_]+)\s*\(`)
)

// ScanSymbols inspects code files across the repository to identify candidate core logic.
func ScanSymbols(repoDir string) []DiscoveredSymbol {
	var symbols []DiscoveredSymbol

	_ = filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if IgnoredDirectories[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip non-code files or files > 1MB to avoid resource exhaustion
		if info.Size() > 1024*1024 {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".go" && ext != ".ts" && ext != ".js" && ext != ".py" && ext != ".rs" {
			return nil
		}

		relPath, _ := filepath.Rel(repoDir, path)
		normPath := filepath.ToSlash(relPath)

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 0

		for scanner.Scan() {
			lineNum++
			line := strings.TrimSpace(scanner.Text())

			switch ext {
			case ".go":
				if reGoMain.MatchString(line) {
					symbols = append(symbols, DiscoveredSymbol{File: normPath, Kind: "entrypoint", Name: "main", Line: lineNum})
				} else if m := reGoFunc.FindStringSubmatch(line); len(m) > 1 {
					symbols = append(symbols, DiscoveredSymbol{File: normPath, Kind: "exported_func", Name: m[1], Line: lineNum})
				} else if m := reGoType.FindStringSubmatch(line); len(m) > 1 {
					symbols = append(symbols, DiscoveredSymbol{File: normPath, Kind: "type_def", Name: m[1], Line: lineNum})
				} else if m := reGoRoute.FindStringSubmatch(line); len(m) > 2 {
					symbols = append(symbols, DiscoveredSymbol{File: normPath, Kind: "route", Name: m[1] + " " + m[2], Line: lineNum})
				}

			case ".ts", ".js":
				if m := reTsExport.FindStringSubmatch(line); len(m) > 2 {
					symbols = append(symbols, DiscoveredSymbol{File: normPath, Kind: "exported_" + m[1], Name: m[2], Line: lineNum})
				}

			case ".py":
				if m := rePyDef.FindStringSubmatch(line); len(m) > 1 && !strings.HasPrefix(m[1], "_") {
					symbols = append(symbols, DiscoveredSymbol{File: normPath, Kind: "function", Name: m[1], Line: lineNum})
				}
			}
		}

		return nil
	})

	return symbols
}
