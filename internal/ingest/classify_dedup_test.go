package ingest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hamziCodes/CTXbank/internal/core"
)

func TestClassifyTarget(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    string
	}{
		{"tech decision note", "We decided to use Postgres for storage. The database schema and SQL migration use JWT auth on every API endpoint.", "decisionLog.md"},
		{"pure tech note", "The API endpoint uses Redis cache, Docker, and a SQL database with low latency.", "techContext.md"},
		{"architecture note", "The architecture uses a layered module pattern; each component keeps coupling low.", "systemPatterns.md"},
		{"product note", "Our user onboarding experience should feel fast for a first-time customer.", "productContext.md"},
	}
	for _, c := range cases {
		if got := classifyTarget(c.content); got != c.want {
			t.Errorf("%s: classifyTarget = %s, want %s", c.name, got, c.want)
		}
	}
}

func TestPrepareIngestionDedupesWithinDocument(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}
	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)

	src := filepath.Join(tempWorkspace, "note.md")
	content := "# Research note\nThe sky survey uses a wide field camera.\nThe sky survey uses a wide field camera.\nCalibration runs nightly at the observatory.\n"
	if err := os.WriteFile(src, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	prop, err := PrepareIngestion(bankDir, src)
	if err != nil {
		t.Fatalf("PrepareIngestion failed: %v", err)
	}
	if n := strings.Count(prop.ProposedDiff, "The sky survey uses a wide field camera."); n != 1 {
		t.Errorf("duplicate line committed %d times, want 1:\n%s", n, prop.ProposedDiff)
	}
	if prop.Deduplicated < 1 {
		t.Errorf("expected deduplicated count >= 1, got %d", prop.Deduplicated)
	}
}
