package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hamziCodes/CTXbank/internal/core"
)

func TestUniqueCheckpointID(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}
	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)

	// Free base ID passes through unchanged.
	if got := uniqueCheckpointID(bankDir, "ckpt_free"); got != "ckpt_free" {
		t.Errorf("expected ckpt_free, got %s", got)
	}

	// Occupy ckpt_busy.json -> expect -2 suffix.
	dir := CheckpointDir(bankDir)
	if err := os.WriteFile(filepath.Join(dir, "ckpt_busy.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	if got := uniqueCheckpointID(bankDir, "ckpt_busy"); got != "ckpt_busy-2" {
		t.Errorf("expected ckpt_busy-2, got %s", got)
	}

	// Occupy -2 as well -> expect -3.
	if err := os.WriteFile(filepath.Join(dir, "ckpt_busy-2.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	if got := uniqueCheckpointID(bankDir, "ckpt_busy"); got != "ckpt_busy-3" {
		t.Errorf("expected ckpt_busy-3, got %s", got)
	}
}

func TestCreateSnapshotSameSecondNoOverwrite(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}
	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)

	// Force both snapshots into the same second by pre-creating the file
	// for the current second's ID.
	fixed := GenerateID(time.Now().UTC())
	dir := CheckpointDir(bankDir)
	if err := os.WriteFile(filepath.Join(dir, fixed+".json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	ckpt, err := CreateSnapshot(bankDir, tempWorkspace, "collision test", nil, nil)
	if err != nil {
		t.Fatalf("CreateSnapshot failed: %v", err)
	}
	if ckpt.ID == fixed {
		t.Errorf("expected suffixed ID on collision, got bare %s (would overwrite)", ckpt.ID)
	}
	// The pre-existing file must be untouched.
	if data, _ := os.ReadFile(filepath.Join(dir, fixed+".json")); string(data) != "{}" {
		t.Errorf("pre-existing snapshot file was overwritten")
	}
	if _, err := LoadSnapshot(bankDir, ckpt.ID); err != nil {
		t.Errorf("new snapshot not loadable under %s: %v", ckpt.ID, err)
	}
}
