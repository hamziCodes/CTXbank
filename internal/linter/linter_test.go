package linter

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/ctxbank/ctx/internal/core"
)

func TestLintMemoryCompliant(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)
	res, err := LintMemoryBank(bankDir, false)
	if err != nil {
		t.Fatalf("LintMemoryBank returned error: %v", err)
	}

	if !res.Passed {
		t.Errorf("freshly scaffolded memory bank should pass lint: violations=%v", res.Violations)
	}
}

func TestLintMemoryOversized(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)
	activePath := filepath.Join(bankDir, "activeContext.md")

	// Blow past 150-line limit
	bloatedContent := `# Active Context

## Focus
Bloat test

## Recent
` + strings.Repeat("- extra bullet item\n", 160) + `
## Next steps
1. Fix bloat

## Open decisions
- None
`
	_ = core.WriteAtomic(activePath, []byte(bloatedContent), 0644)

	res, err := LintMemoryBank(bankDir, false)
	if err != nil {
		t.Fatalf("LintMemoryBank error: %v", err)
	}

	if res.Passed {
		t.Errorf("oversized activeContext should fail linting")
	}

	// Now test auto-fix
	fixRes, err := LintMemoryBank(bankDir, true)
	if err != nil {
		t.Fatalf("Lint with auto-fix failed: %v", err)
	}
	if len(fixRes.Warnings) == 0 {
		t.Errorf("expected warning indicating auto-fix pruned bullets")
	}
}
