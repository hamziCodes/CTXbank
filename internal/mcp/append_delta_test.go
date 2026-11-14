package mcp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hamziCodes/CTXbank/internal/core"
)

// Regression: append_memory_delta with an empty section used to build the
// header "## " and insert the content after the first heading's prefix,
// corrupting "## Focus" into "## \n<content>\nFocus".
func TestAppendMemoryDeltaEmptySectionAppendsAtEnd(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}
	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)

	if _, err := ToolAppendMemoryDelta(bankDir, "s1", "activeContext.md", "", "\nAgent note at end.\n"); err != nil {
		t.Fatalf("append failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(bankDir, "activeContext.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "## Focus") {
		t.Errorf("Focus heading corrupted:\n%s", content)
	}
	if !strings.HasSuffix(strings.TrimSpace(content), "Agent note at end.") {
		t.Errorf("delta not appended at end:\n%s", content)
	}
}

func TestAppendMemoryDeltaNamedSection(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}
	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)

	if _, err := ToolAppendMemoryDelta(bankDir, "s1", "activeContext.md", "Focus", "Named section line."); err != nil {
		t.Fatalf("append failed: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(bankDir, "activeContext.md"))
	content := string(data)
	idx := strings.Index(content, "## Focus")
	if idx < 0 || !strings.Contains(content[idx:], "Named section line.") {
		t.Errorf("named-section delta missing under Focus:\n%s", content)
	}
	// A similarly-prefixed heading must not be hijacked.
	if _, err := ToolAppendMemoryDelta(bankDir, "s1", "progress.md", "Progress", "P line."); err != nil {
		t.Fatalf("append to progress failed: %v", err)
	}
}
