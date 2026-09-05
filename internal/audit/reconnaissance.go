package audit

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ctxbank/ctx/internal/core"
)

// AuditReport aggregates all 4 deterministic reconnaissance passes.
type AuditReport struct {
	Manifests  []ProjectManifestInfo   `json:"manifests"`
	Components []ComponentArchitecture `json:"components"`
	Symbols    []DiscoveredSymbol      `json:"symbols"`
	Velocity   []PathVelocity          `json:"velocity"`
}

// RunAudit executes all 4 deterministic reconnaissance passes in order.
func RunAudit(repoDir string) (*AuditReport, error) {
	manifests := ScanManifests(repoDir)
	components := ScanStructuralHeuristics(repoDir)
	symbols := ScanSymbols(repoDir)
	velocity, _ := ScanGitVelocity(repoDir)

	return &AuditReport{
		Manifests:  manifests,
		Components: components,
		Symbols:    symbols,
		Velocity:   velocity,
	}, nil
}

// GenerateTechContextMarkdown turns the audit results into a deterministic techContext.md document.
func (r *AuditReport) GenerateTechContextMarkdown() string {
	var sb strings.Builder
	sb.WriteString("# Tech Context — generated from deterministic audit\n\n")

	sb.WriteString("## Discovered Tech Stack\n")
	for _, m := range r.Manifests {
		sb.WriteString(fmt.Sprintf("- **Language:** %s\n", m.Language))
		sb.WriteString(fmt.Sprintf("- **Build System:** %s\n", m.BuildSystem))
		if len(m.Dependencies) > 0 {
			sb.WriteString(fmt.Sprintf("- **Key Dependencies:** %s\n", strings.Join(m.Dependencies, ", ")))
		}
		if len(m.Scripts) > 0 {
			sb.WriteString("- **Build & Test Commands:**\n")
			for k, v := range m.Scripts {
				sb.WriteString(fmt.Sprintf("  - `%s`: `%s`\n", k, v))
			}
		}
	}

	sb.WriteString("\n## Active Repository Velocity (Last 90 Days)\n")
	if len(r.Velocity) > 0 {
		for _, v := range r.Velocity {
			sb.WriteString(fmt.Sprintf("- `%s` (%d commits)\n", v.Path, v.CommitCount))
		}
	} else {
		sb.WriteString("- No recent commit history found (fresh repository or shallow clone).\n")
	}

	return sb.String()
}

// GenerateSystemPatternsMarkdown turns audit results into a deterministic systemPatterns.md component map.
func (r *AuditReport) GenerateSystemPatternsMarkdown() string {
	var sb strings.Builder
	sb.WriteString("# System Patterns — generated from deterministic audit\n\n")

	sb.WriteString("## Architectural Component Map\n")
	for _, c := range r.Components {
		sb.WriteString(fmt.Sprintf("### %s (`%s`)\n", c.Name, c.Path))
		sb.WriteString(fmt.Sprintf("%s\n\n", c.Description))
	}

	sb.WriteString("## Core Entrypoints & Symbols\n")
	entryCount := 0
	for _, s := range r.Symbols {
		if s.Kind == "entrypoint" || s.Kind == "route" {
			sb.WriteString(fmt.Sprintf("- **%s** (`%s`:%d) — `%s`\n", s.Kind, s.File, s.Line, s.Name))
			entryCount++
		}
	}
	if entryCount == 0 {
		sb.WriteString("- Standard library / modular package layout.\n")
	}

	return sb.String()
}

// ApplyAudit updates memory-bank/techContext.md and systemPatterns.md with audited results.
func ApplyAudit(bankDir string, report *AuditReport) error {
	techPath := filepath.Join(bankDir, "techContext.md")
	if err := core.WriteAtomic(techPath, []byte(report.GenerateTechContextMarkdown()), 0644); err != nil {
		return err
	}
	_ = core.UpdateFileMeta(bankDir, "techContext.md", "")

	sysPath := filepath.Join(bankDir, "systemPatterns.md")
	if err := core.WriteAtomic(sysPath, []byte(report.GenerateSystemPatternsMarkdown()), 0644); err != nil {
		return err
	}
	_ = core.UpdateFileMeta(bankDir, "systemPatterns.md", "")

	return nil
}
