package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitBankScaffold(t *testing.T) {
	tempWorkspace := t.TempDir()

	if err := InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	bankDir := filepath.Join(tempWorkspace, MemoryBankDir)

	// Check that all 7 standard files exist
	expectedFiles := []string{
		"projectbrief.md",
		"productContext.md",
		"systemPatterns.md",
		"techContext.md",
		"decisionLog.md",
		"activeContext.md",
		"progress.md",
	}

	for _, filename := range expectedFiles {
		path := filepath.Join(bankDir, filename)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s does not exist: %v", filename, err)
		}
	}

	// Check that manifest exists and has all 7 files recorded
	manifest, err := LoadManifest(bankDir)
	if err != nil {
		t.Fatalf("failed to load manifest: %v", err)
	}

	for _, filename := range expectedFiles {
		meta, ok := manifest.Files[filename]
		if !ok {
			t.Errorf("manifest missing entry for %s", filename)
		}
		if meta.SHA256 == "" || meta.Bytes == 0 {
			t.Errorf("invalid metadata for %s: %+v", filename, meta)
		}
	}

	// Re-initializing without force should fail
	if err := InitBank(tempWorkspace, false); err == nil {
		t.Errorf("expected error re-initializing without force")
	}

	// Re-initializing with force should succeed
	if err := InitBank(tempWorkspace, true); err != nil {
		t.Errorf("InitBank with force failed: %v", err)
	}
}
