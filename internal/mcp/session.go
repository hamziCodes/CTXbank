package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ctxbank/ctx/internal/core"
)

var sessionMu sync.Mutex

// CheckAndAdvanceCursor checks if the specified file has changed relative to the caller's session cursor.
// If unchanged, returns changed=false and empty content.
// If changed or unread, updates the session cursor in manifest.json and returns changed=true with the file content.
func CheckAndAdvanceCursor(bankDir, sessionID, filename string) (bool, string, []byte, error) {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	fullPath := filepath.Join(bankDir, filename)
	currentHash, _, err := core.ComputeFileHash(fullPath)
	if err != nil {
		return false, "", nil, fmt.Errorf("failed to compute hash for %s: %w", filename, err)
	}

	manifest, err := core.LoadManifest(bankDir)
	if err != nil {
		return false, "", nil, fmt.Errorf("failed to load manifest: %w", err)
	}

	if manifest.AgentReadCursors == nil {
		manifest.AgentReadCursors = make(map[string]map[string]string)
	}

	sessionCursors, exists := manifest.AgentReadCursors[sessionID]
	if !exists {
		sessionCursors = make(map[string]string)
		manifest.AgentReadCursors[sessionID] = sessionCursors
	}

	lastKnownHash := sessionCursors[filename]
	if lastKnownHash == currentHash {
		// File is untouched since this agent session last read it
		return false, currentHash, nil, nil
	}

	// Content has changed or first read; update cursor
	sessionCursors[filename] = currentHash
	manifest.AgentReadCursors[sessionID] = sessionCursors

	if err := core.SaveManifest(bankDir, manifest); err != nil {
		return false, "", nil, fmt.Errorf("failed to save manifest cursor: %w", err)
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return false, "", nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return true, currentHash, content, nil
}
