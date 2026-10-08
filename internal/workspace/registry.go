package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/hamziCodes/CTXbank/internal/core"
)

const maxRegistryEntries = 20

// ProjectRecord is one remembered CTXbank workspace.
type ProjectRecord struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	LastOpened time.Time `json:"last_opened"`
}

func registryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".ctxbank", "projects.json"), nil
}

func loadRegistry() ([]ProjectRecord, error) {
	path, err := registryPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var records []ProjectRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, nil // corrupt registry: start fresh rather than fail
	}
	return records, nil
}

func saveRegistry(records []ProjectRecord) error {
	path, err := registryPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// RecordProject remembers a workspace (called on `ctx init` and `ctx ui`).
// Most-recently-used first, capped at maxRegistryEntries.
func RecordProject(dir string) error {
	clean := filepath.Clean(dir)
	abs, err := filepath.Abs(clean)
	if err != nil {
		abs = clean
	}
	records, _ := loadRegistry()
	now := time.Now().UTC()
	kept := make([]ProjectRecord, 0, maxRegistryEntries)
	kept = append(kept, ProjectRecord{
		Name:       filepath.Base(abs),
		Path:       abs,
		LastOpened: now,
	})
	for _, r := range records {
		if r.Path == abs {
			continue // re-inserted at front above
		}
		kept = append(kept, r)
		if len(kept) >= maxRegistryEntries {
			break
		}
	}
	return saveRegistry(kept)
}

// RecentProjects returns remembered workspaces that still exist and still
// contain a memory-bank manifest, most-recent first.
func RecentProjects() ([]ProjectRecord, error) {
	records, err := loadRegistry()
	if err != nil {
		return nil, err
	}
	var live []ProjectRecord
	for _, r := range records {
		manifestPath := filepath.Join(r.Path, core.MemoryBankDir, core.StateDirname, core.ManifestFilename)
		if _, err := os.Stat(manifestPath); err == nil {
			live = append(live, r)
		}
	}
	return live, nil
}
