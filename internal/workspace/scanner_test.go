package workspace

import (
	"path/filepath"
	"testing"

	"github.com/ctxbank/ctx/internal/core"
)

func TestScanWorkspaces(t *testing.T) {
	parentDir := t.TempDir()

	proj1 := filepath.Join(parentDir, "project-alpha")
	proj2 := filepath.Join(parentDir, "project-beta")
	nonProj := filepath.Join(parentDir, "regular-folder")

	if err := core.InitBank(proj1, false); err != nil {
		t.Fatalf("InitBank proj1 failed: %v", err)
	}
	if err := core.InitBank(proj2, false); err != nil {
		t.Fatalf("InitBank proj2 failed: %v", err)
	}
	_ = core.WriteAtomic(filepath.Join(nonProj, "readme.txt"), []byte("hello"), 0644)

	results, err := ScanWorkspaces(parentDir, 2)
	if err != nil {
		t.Fatalf("ScanWorkspaces failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected exactly 2 managed projects discovered, got %d", len(results))
	}

	names := map[string]bool{
		results[0].Name: true,
		results[1].Name: true,
	}

	if !names["project-alpha"] || !names["project-beta"] {
		t.Errorf("expected project-alpha and project-beta found, got %+v", results)
	}
}
