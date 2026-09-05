package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// WriteAtomic writes data to targetPath in a crash-safe, atomic manner:
// 1. Writes to a temporary sibling file: <path>.tmp.<randHex>
// 2. Flushes buffers to physical disk via fsync()
// 3. Atomically renames over targetPath
func WriteAtomic(targetPath string, data []byte, perm os.FileMode) error {
	cleanPath := filepath.Clean(targetPath)
	dir := filepath.Dir(cleanPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directories for %s: %w", cleanPath, err)
	}

	// Generate random suffix to avoid collisions in concurrent executions
	randBytes := make([]byte, 8)
	if _, err := rand.Read(randBytes); err != nil {
		return fmt.Errorf("failed to generate random salt: %w", err)
	}
	tmpPath := fmt.Sprintf("%s.tmp.%s", cleanPath, hex.EncodeToString(randBytes))

	tmpFile, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("failed to create temporary file %s: %w", tmpPath, err)
	}

	// Ensure cleanup in case of panic or premature failure
	success := false
	defer func() {
		if !success {
			_ = tmpFile.Close()
			_ = os.Remove(tmpPath)
		}
	}()

	// Write payload
	if _, err := tmpFile.Write(data); err != nil {
		return fmt.Errorf("failed to write data to %s: %w", tmpPath, err)
	}

	// Fsync to guarantee bytes are physically committed to storage
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("fsync failed on %s: %w", tmpPath, err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close %s: %w", tmpPath, err)
	}

	// Atomic rename to target path
	if err := os.Rename(tmpPath, cleanPath); err != nil {
		return fmt.Errorf("atomic rename from %s to %s failed: %w", tmpPath, cleanPath, err)
	}

	success = true
	return nil
}
