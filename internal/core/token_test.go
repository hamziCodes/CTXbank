package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectTokenCreateGetRegenerate(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}
	bankDir := filepath.Join(tempWorkspace, MemoryBankDir)

	// First call creates and persists.
	tok1, err := GetProjectToken(bankDir)
	if err != nil {
		t.Fatalf("GetProjectToken: %v", err)
	}
	if len(tok1) != 32 {
		t.Errorf("expected 32-char hex token, got %q", tok1)
	}

	// Second call returns the same persisted token.
	tok2, err := GetProjectToken(bankDir)
	if err != nil {
		t.Fatalf("GetProjectToken (2nd): %v", err)
	}
	if tok1 != tok2 {
		t.Errorf("token not persisted: %q != %q", tok1, tok2)
	}

	// File must not be world-readable.
	fi, err := os.Stat(filepath.Join(bankDir, StateDirname, ProjectTokenFilename))
	if err != nil {
		t.Fatalf("token file missing: %v", err)
	}
	if perm := fi.Mode().Perm(); perm != 0600 {
		t.Errorf("expected 0600 token file, got %o", perm)
	}

	// Regenerate rotates.
	tok3, err := RegenerateProjectToken(bankDir)
	if err != nil {
		t.Fatalf("RegenerateProjectToken: %v", err)
	}
	if tok3 == tok1 {
		t.Errorf("regenerated token identical to old one")
	}
	if tok4, _ := GetProjectToken(bankDir); tok4 != tok3 {
		t.Errorf("GetProjectToken did not return regenerated token")
	}
}

func TestEnsureStateGitignored(t *testing.T) {
	tempWorkspace := t.TempDir()

	// Fresh repo: creates the entry.
	EnsureStateGitignored(tempWorkspace)
	data, err := os.ReadFile(filepath.Join(tempWorkspace, ".gitignore"))
	if err != nil {
		t.Fatalf("expected .gitignore to be created: %v", err)
	}
	if got := string(data); got != "memory-bank/.state/\n" {
		t.Errorf("unexpected .gitignore content: %q", got)
	}

	// Existing content: appended, not clobbered, not duplicated.
	_ = os.WriteFile(filepath.Join(tempWorkspace, ".gitignore"), []byte("node_modules/\n"), 0644)
	EnsureStateGitignored(tempWorkspace)
	EnsureStateGitignored(tempWorkspace) // idempotent
	data, _ = os.ReadFile(filepath.Join(tempWorkspace, ".gitignore"))
	want := "node_modules/\nmemory-bank/.state/\n"
	if string(data) != want {
		t.Errorf("expected %q, got %q", want, string(data))
	}
}
