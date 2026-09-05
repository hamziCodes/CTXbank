package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ctxbank/ctx/pkg/types"
)

const (
	ManifestFilename = "manifest.json"
	StateDirname     = ".state"
)

// StripFrontmatter removes YAML (`--- ... ---`) or TOML (`+++ ... +++`) frontmatter
// from markdown before computing content-based SHA-256 hashes.
func StripFrontmatter(content []byte) []byte {
	trimmed := bytes.TrimSpace(content)
	if bytes.HasPrefix(trimmed, []byte("---")) {
		rest := trimmed[3:]
		idx := bytes.Index(rest, []byte("---"))
		if idx != -1 {
			return bytes.TrimLeft(rest[idx+3:], "\r\n")
		}
	} else if bytes.HasPrefix(trimmed, []byte("+++")) {
		rest := trimmed[3:]
		idx := bytes.Index(rest, []byte("+++"))
		if idx != -1 {
			return bytes.TrimLeft(rest[idx+3:], "\r\n")
		}
	}
	return content
}

// ComputeFileHash calculates the SHA-256 checksum (excluding frontmatter) and byte size of a file.
func ComputeFileHash(path string) (string, int64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	body := StripFrontmatter(data)
	hasher := sha256.New()
	hasher.Write(body)
	hashStr := hex.EncodeToString(hasher.Sum(nil))

	return hashStr, int64(len(data)), nil
}

// ManifestPath returns the absolute path to manifest.json given the memory bank directory.
func ManifestPath(bankDir string) string {
	return filepath.Join(bankDir, StateDirname, ManifestFilename)
}

// LoadManifest reads and decodes the manifest.json file from the specified memory-bank directory.
func LoadManifest(bankDir string) (*types.Manifest, error) {
	path := ManifestPath(bankDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return types.NewManifest(), nil
		}
		return nil, fmt.Errorf("failed to read manifest at %s: %w", path, err)
	}

	var m types.Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("corrupt manifest JSON at %s: %w", path, err)
	}
	if m.Files == nil {
		m.Files = make(map[string]types.FileMeta)
	}
	if m.AgentReadCursors == nil {
		m.AgentReadCursors = make(map[string]map[string]string)
	}
	return &m, nil
}

// SaveManifest writes the manifest atomically to .state/manifest.json.
func SaveManifest(bankDir string, m *types.Manifest) error {
	path := ManifestPath(bankDir)
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize manifest: %w", err)
	}
	// Append trailing newline
	data = append(data, '\n')
	return WriteAtomic(path, data, 0644)
}

// UpdateFileMeta recalculates hash and byte size for relPath and records it in manifest.json.
func UpdateFileMeta(bankDir, relPath, lastCkpt string) error {
	fullPath := filepath.Join(bankDir, relPath)
	hashStr, bytesCount, err := ComputeFileHash(fullPath)
	if err != nil {
		return err
	}

	m, err := LoadManifest(bankDir)
	if err != nil {
		return err
	}

	// Infer volatility class
	volatility := types.VolatilitySemiStatic
	baseName := strings.ToLower(filepath.Base(relPath))
	switch baseName {
	case "projectbrief.md", "productcontext.md":
		volatility = types.VolatilityStatic
	case "activecontext.md", "progress.md":
		volatility = types.VolatilityHot
	case "decisionlog.md", "systempatterns.md", "techcontext.md":
		volatility = types.VolatilitySemiStatic
	}

	m.Files[relPath] = types.FileMeta{
		SHA256:         hashStr,
		Bytes:          bytesCount,
		LastCheckpoint: lastCkpt,
		UpdatedAt:      time.Now().UTC(),
		Volatility:     volatility,
	}

	return SaveManifest(bankDir, m)
}
