package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ctxbank/ctx/pkg/types"
)

func TestStripFrontmatter(t *testing.T) {
	yamlInput := []byte(`---
title: Sample
volatility: hot
---
# Main Content
Line 1
Line 2`)

	stripped := StripFrontmatter(yamlInput)
	expected := "# Main Content\nLine 1\nLine 2"
	if string(stripped) != expected {
		t.Errorf("expected %q, got %q", expected, string(stripped))
	}

	noFrontmatter := []byte("# Direct Markdown")
	if string(StripFrontmatter(noFrontmatter)) != "# Direct Markdown" {
		t.Errorf("content without frontmatter should remain untouched")
	}
}

func TestManifestSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()

	manifest := types.NewManifest()
	manifest.Files["activeContext.md"] = types.FileMeta{
		SHA256:         "abc123hash",
		Bytes:          1024,
		LastCheckpoint: "ckpt_test_01",
		Volatility:     types.VolatilityHot,
	}
	manifest.AgentReadCursors["agent-1"] = map[string]string{
		"activeContext.md": "abc123hash",
	}

	if err := SaveManifest(tempDir, manifest); err != nil {
		t.Fatalf("SaveManifest failed: %v", err)
	}

	loaded, err := LoadManifest(tempDir)
	if err != nil {
		t.Fatalf("LoadManifest failed: %v", err)
	}

	meta, ok := loaded.Files["activeContext.md"]
	if !ok {
		t.Fatalf("activeContext.md missing from loaded manifest")
	}

	if meta.SHA256 != "abc123hash" || meta.Bytes != 1024 {
		t.Errorf("mismatched metadata: %+v", meta)
	}

	cursorHash, ok := loaded.AgentReadCursors["agent-1"]["activeContext.md"]
	if !ok || cursorHash != "abc123hash" {
		t.Errorf("mismatched cursor: %+v", loaded.AgentReadCursors)
	}
}

func TestComputeFileHash(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sample.md")
	content := []byte("---\nmeta: true\n---\nBody text")

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hashStr, bytesCount, err := ComputeFileHash(filePath)
	if err != nil {
		t.Fatalf("ComputeFileHash failed: %v", err)
	}

	if hashStr == "" {
		t.Errorf("hash string should not be empty")
	}
	if bytesCount != int64(len(content)) {
		t.Errorf("expected byte count %d, got %d", len(content), bytesCount)
	}
}
