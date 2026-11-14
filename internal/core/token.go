package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectTokenFilename is the file inside memory-bank/.state/ holding this
// project's dashboard token. One token per project: it is what the hosted
// site's "Connect" page asks for, and what `ctx ui` requires on every
// /api/* call. It is never committed (see ensureStateGitignored).
const ProjectTokenFilename = "dashboard_token"

func projectTokenPath(bankDir string) string {
	return filepath.Join(bankDir, StateDirname, ProjectTokenFilename)
}

// generateTokenValue returns a fresh 128-bit random hex token.
func generateTokenValue() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// GetProjectToken returns this project's dashboard token, creating and
// persisting one (0600) on first use.
func GetProjectToken(bankDir string) (string, error) {
	path := projectTokenPath(bankDir)
	if data, err := os.ReadFile(path); err == nil {
		if token := strings.TrimSpace(string(data)); token != "" {
			return token, nil
		}
	}
	return RegenerateProjectToken(bankDir)
}

// RegenerateProjectToken replaces the project token with a fresh value.
// The old token stops working immediately.
func RegenerateProjectToken(bankDir string) (string, error) {
	token, err := generateTokenValue()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(projectTokenPath(bankDir)), 0755); err != nil {
		return "", err
	}
	// 0600: the token is a secret; never world-readable.
	if err := os.WriteFile(projectTokenPath(bankDir), []byte(token+"\n"), 0600); err != nil {
		return "", err
	}
	return token, nil
}

// EnsureStateGitignored makes sure memory-bank/.state/ (manifest,
// checkpoints, dashboard token, sync ledger) is never committed.
// It appends to the repo .gitignore without clobbering or duplicating.
func EnsureStateGitignored(repoDir string) {
	ignorePath := filepath.Join(repoDir, ".gitignore")
	entry := MemoryBankDir + "/" + StateDirname + "/"

	var existing string
	if data, err := os.ReadFile(ignorePath); err == nil {
		existing = string(data)
		for _, line := range strings.Split(existing, "\n") {
			if strings.TrimSpace(line) == entry {
				return // already ignored
			}
		}
	}
	if existing != "" && !strings.HasSuffix(existing, "\n") {
		existing += "\n"
	}
	_ = os.WriteFile(ignorePath, []byte(existing+entry+"\n"), 0644)
}
