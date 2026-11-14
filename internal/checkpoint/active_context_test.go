package checkpoint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hamziCodes/CTXbank/internal/core"
	"github.com/hamziCodes/CTXbank/pkg/types"
)

// The pause -> resume loop: a checkpoint must be reflected in
// activeContext.md, because resume/status read Focus/Recent from that file.
func TestUpdateActiveContextClosesLoop(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}
	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)

	ckpt := &types.Checkpoint{
		ID:          "ckpt_20261114_120000",
		ActiveFocus: "wiring CLI",
		ManualEdits: []string{"built the adder"},
	}
	updateActiveContext(bankDir, ckpt)

	data, err := os.ReadFile(filepath.Join(bankDir, "activeContext.md"))
	if err != nil {
		t.Fatalf("read activeContext: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "# Active Context — updated ckpt_20261114_120000") {
		t.Errorf("title not stamped with checkpoint ID:\n%s", content)
	}
	if !strings.Contains(content, "## Focus\n\nwiring CLI\n") {
		t.Errorf("focus not written:\n%s", content)
	}
	if strings.Contains(content, "<!-- TODO: what are you working on right now?") {
		t.Errorf("template placeholder survived:\n%s", content)
	}
	if !strings.Contains(content, "- ckpt_20261114_120000: built the adder") {
		t.Errorf("recent bullet missing:\n%s", content)
	}
	if strings.Contains(content, "- Initialized memory bank layout.") {
		t.Errorf("bootstrap placeholder bullet survived:\n%s", content)
	}
	if !strings.Contains(content, "## Open decisions") {
		t.Errorf("open decisions section lost:\n%s", content)
	}

	// Three more checkpoints: Recent must cap at 3, newest first.
	for _, id := range []string{"ckpt_2", "ckpt_3", "ckpt_4"} {
		updateActiveContext(bankDir, &types.Checkpoint{ID: id, ActiveFocus: "f-" + id})
	}
	data, _ = os.ReadFile(filepath.Join(bankDir, "activeContext.md"))
	content = string(data)
	recent := content[strings.Index(content, "## Recent"):strings.Index(content, "## Next steps")]
	bullets := 0
	for _, l := range strings.Split(recent, "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), "- ckpt_") {
			bullets++
		}
	}
	if bullets != 3 {
		t.Errorf("expected exactly 3 recent bullets, got %d:\n%s", bullets, recent)
	}
	if !strings.Contains(recent, "- ckpt_4:") {
		t.Errorf("newest checkpoint not first:\n%s", recent)
	}
}

// CreateSnapshot on a non-git temp dir must still work in tests that use a
// plain folder (the helper is best-effort and must not corrupt the file
// when sections are hand-edited).
func TestUpdateActiveContextMissingFile(t *testing.T) {
	bankDir := t.TempDir()
	// No activeContext.md: must be a silent no-op, not a panic.
	updateActiveContext(bankDir, &types.Checkpoint{ID: "ckpt_x", ActiveFocus: "f"})
	if _, err := os.Stat(filepath.Join(bankDir, "activeContext.md")); !os.IsNotExist(err) {
		t.Errorf("missing file must not be created by the helper")
	}
}
