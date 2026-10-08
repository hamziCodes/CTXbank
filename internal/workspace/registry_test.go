package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hamziCodes/CTXbank/internal/core"
)

func withTempHome(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	// os.UserHomeDir caches on some platforms only after first call; on
	// linux it reads $HOME every time, which t.Setenv controls.
}

func initBankAt(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := core.InitBank(dir, false); err != nil {
		t.Fatalf("InitBank(%s) failed: %v", dir, err)
	}
}

func TestRegistryRecordAndRecent(t *testing.T) {
	withTempHome(t)
	base := t.TempDir()
	projA := filepath.Join(base, "proj-a")
	projB := filepath.Join(base, "proj-b")
	initBankAt(t, projA)
	initBankAt(t, projB)

	if err := RecordProject(projA); err != nil {
		t.Fatalf("RecordProject A: %v", err)
	}
	if err := RecordProject(projB); err != nil {
		t.Fatalf("RecordProject B: %v", err)
	}

	recent, err := RecentProjects()
	if err != nil {
		t.Fatalf("RecentProjects: %v", err)
	}
	if len(recent) != 2 {
		t.Fatalf("expected 2 recent projects, got %d", len(recent))
	}
	// Most recently recorded first.
	if recent[0].Path != projB || recent[1].Path != projA {
		t.Errorf("wrong MRU order: %+v", recent)
	}

	// Re-recording moves to front without duplicating.
	if err := RecordProject(projA); err != nil {
		t.Fatal(err)
	}
	recent, _ = RecentProjects()
	if len(recent) != 2 || recent[0].Path != projA {
		t.Errorf("re-record should move to front without dup: %+v", recent)
	}

	// A project whose manifest vanished is filtered out.
	os.RemoveAll(filepath.Join(projB, core.MemoryBankDir))
	recent, _ = RecentProjects()
	if len(recent) != 1 || recent[0].Path != projA {
		t.Errorf("expected only proj-a after manifest removal: %+v", recent)
	}
}
