package tui

import (
	"fmt"
	"strings"

	"github.com/ctxbank/ctx/internal/audit"
)

// RenderAuditReport formats the full reconnaissance audit report for terminal output.
func RenderAuditReport(report *audit.AuditReport) string {
	box := GetBox()
	var sb strings.Builder

	header := fmt.Sprintf("%s BROWNFIELD RECONNAISSANCE AUDIT %s", box.TopLeft, box.TopRight)
	sb.WriteString(header + "\n")

	// Pass 1: Manifests
	sb.WriteString(fmt.Sprintf("%s [Pass 1/4] Manifest & Dependency Scan\n", box.Check))
	for _, m := range report.Manifests {
		sb.WriteString(fmt.Sprintf("  • Language: %s (Build: %s)\n", m.Language, m.BuildSystem))
		if len(m.Dependencies) > 0 {
			depsStr := strings.Join(m.Dependencies, ", ")
			if len(depsStr) > 60 {
				depsStr = depsStr[:57] + "..."
			}
			sb.WriteString(fmt.Sprintf("  • Dependencies (%d): %s\n", len(m.Dependencies), depsStr))
		}
	}

	// Pass 2: Architecture Components
	sb.WriteString(fmt.Sprintf("\n%s [Pass 2/4] Structural Component Heuristics\n", box.Check))
	for _, c := range report.Components {
		sb.WriteString(fmt.Sprintf("  • %-30s -> %s\n", c.Name, c.Path))
	}

	// Pass 3: Symbols & Entrypoints
	sb.WriteString(fmt.Sprintf("\n%s [Pass 3/4] Symbols & Core Logic Extraction\n", box.Check))
	entrypoints := 0
	funcs := 0
	for _, s := range report.Symbols {
		if s.Kind == "entrypoint" || s.Kind == "route" {
			entrypoints++
		} else {
			funcs++
		}
	}
	sb.WriteString(fmt.Sprintf("  • Discovered %d entrypoints/routes, %d exported types/functions\n", entrypoints, funcs))

	// Pass 4: Git Velocity
	sb.WriteString(fmt.Sprintf("\n%s [Pass 4/4] Git Commit Velocity (Last 90 Days)\n", box.Check))
	if len(report.Velocity) > 0 {
		for i, v := range report.Velocity {
			if i >= 5 {
				break
			}
			sb.WriteString(fmt.Sprintf("  • %-30s : %d commits\n", v.Path, v.CommitCount))
		}
	} else {
		sb.WriteString("  • No recent commits (clean baseline)\n")
	}

	footer := fmt.Sprintf("%s%s%s", box.BottomLeft, strings.Repeat(box.Horizontal, 45), box.BottomRight)
	sb.WriteString("\n" + footer + "\n")

	return sb.String()
}
