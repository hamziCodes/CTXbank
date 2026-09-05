package core

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestWriteAtomicBasic(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "test.txt")
	content := []byte("hello atomic world")

	if err := WriteAtomic(targetPath, content, 0644); err != nil {
		t.Fatalf("WriteAtomic failed: %v", err)
	}

	readBack, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("failed to read back file: %v", err)
	}

	if string(readBack) != string(content) {
		t.Errorf("expected %q, got %q", string(content), string(readBack))
	}
}

func TestWriteAtomicOverwrite(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "overwrite.txt")

	if err := WriteAtomic(targetPath, []byte("version 1"), 0644); err != nil {
		t.Fatalf("WriteAtomic v1 failed: %v", err)
	}

	if err := WriteAtomic(targetPath, []byte("version 2 updated"), 0644); err != nil {
		t.Fatalf("WriteAtomic v2 failed: %v", err)
	}

	readBack, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("failed to read back file: %v", err)
	}

	if string(readBack) != "version 2 updated" {
		t.Errorf("expected 'version 2 updated', got %q", string(readBack))
	}
}

func TestWriteAtomicConcurrent(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "concurrent.txt")

	var wg sync.WaitGroup
	workers := 10

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			payload := []byte(string(rune('A' + idx)))
			_ = WriteAtomic(targetPath, payload, 0644)
		}(i)
	}

	wg.Wait()

	// Ensure final file exists and is valid length
	info, err := os.Stat(targetPath)
	if err != nil {
		t.Fatalf("concurrent target missing: %v", err)
	}

	if info.Size() == 0 {
		t.Errorf("file was left empty after concurrent writes")
	}

	// Verify no stray .tmp files are left in the directory
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("failed to read temp dir: %v", err)
	}

	for _, entry := range entries {
		if entry.Name() != "concurrent.txt" {
			t.Errorf("uncleaned temporary file remaining: %s", entry.Name())
		}
	}
}
