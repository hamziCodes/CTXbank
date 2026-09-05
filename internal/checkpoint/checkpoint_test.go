package checkpoint

import (
	"path/filepath"
	"testing"

	"github.com/ctxbank/ctx/internal/core"
)

func TestCreateAndLoadSnapshot(t *testing.T) {
	tempWorkspace := t.TempDir()

	// Initialize memory bank first
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)

	focus := "Implementing Test Coverage"
	nextSteps := []string{"Step A", "Step B"}
	manualEdits := []string{"modified something manually"}

	ckpt, err := CreateSnapshot(bankDir, tempWorkspace, focus, nextSteps, manualEdits)
	if err != nil {
		t.Fatalf("CreateSnapshot failed: %v", err)
	}

	if ckpt.ID == "" {
		t.Errorf("expected non-empty checkpoint ID")
	}
	if ckpt.ActiveFocus != focus {
		t.Errorf("expected focus %q, got %q", focus, ckpt.ActiveFocus)
	}

	// Load snapshot back
	loaded, err := LoadSnapshot(bankDir, ckpt.ID)
	if err != nil {
		t.Fatalf("LoadSnapshot failed: %v", err)
	}

	if loaded.ID != ckpt.ID || loaded.ActiveFocus != focus {
		t.Errorf("mismatched loaded snapshot: %+v", loaded)
	}

	// List snapshots
	snapshots, err := ListSnapshots(bankDir)
	if err != nil {
		t.Fatalf("ListSnapshots failed: %v", err)
	}
	if len(snapshots) != 1 || snapshots[0] != ckpt.ID {
		t.Errorf("expected snapshots [%s], got %+v", ckpt.ID, snapshots)
	}
}
