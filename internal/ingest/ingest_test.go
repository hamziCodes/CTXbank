package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ctxbank/ctx/internal/core"
)

func TestSimHashDeduplication(t *testing.T) {
	textA := "This is a deterministic test paragraph for calculating 64-bit simhash fingerprints."
	textB := "This is a deterministic test paragraph for calculating 64-bit simhash fingerprints."
	textC := "Completely unrelated text discussing cloud database deployment strategies."

	hashA := SimHash64(textA)
	hashB := SimHash64(textB)
	hashC := SimHash64(textC)

	if hashA != hashB {
		t.Errorf("identical text must yield identical hash: %x vs %x", hashA, hashB)
	}

	distAC := HammingDistance(hashA, hashC)
	if distAC <= 5 {
		t.Errorf("unrelated text should have high Hamming distance: %d", distAC)
	}

	if !IsNearDuplicate(textA, textB, 3) {
		t.Errorf("expected near duplicate detection to be true")
	}
	if IsNearDuplicate(textA, textC, 3) {
		t.Errorf("unrelated text should not be detected as duplicate")
	}
}

func TestPrepareAndCommitIngestion(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)
	inboxDir := filepath.Join(tempWorkspace, "research", "inbox")
	_ = os.MkdirAll(inboxDir, 0755)

	researchFile := filepath.Join(inboxDir, "notes.md")
	notesContent := `# Candidate Features
- [ ] Implement OCR backoff jitter
- [ ] Add Prometheus metric counters

## High Level Objectives
We must ensure that OCR retries do not synchronize across clients.
`
	if err := os.WriteFile(researchFile, []byte(notesContent), 0644); err != nil {
		t.Fatalf("failed to write notes: %v", err)
	}

	prop, err := PrepareIngestion(bankDir, researchFile)
	if err != nil {
		t.Fatalf("PrepareIngestion failed: %v", err)
	}

	if len(prop.NewItems) < 2 {
		t.Errorf("expected candidate items extracted: %+v", prop.NewItems)
	}

	if err := CommitIngestion(bankDir, tempWorkspace, prop); err != nil {
		t.Fatalf("CommitIngestion failed: %v", err)
	}

	// Verify original file was moved from inbox to archive
	if _, err := os.Stat(researchFile); !os.IsNotExist(err) {
		t.Errorf("research file should have been moved from inbox")
	}

	archiveEntries, _ := os.ReadDir(filepath.Join(tempWorkspace, "research", "archive"))
	if len(archiveEntries) != 1 {
		t.Errorf("expected 1 file in archive, got %d", len(archiveEntries))
	}
}
