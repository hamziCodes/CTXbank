package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ctxbank/ctx/internal/core"
	"github.com/ctxbank/ctx/internal/git"
	"github.com/ctxbank/ctx/pkg/types"
)

// GenerateID produces a deterministic, human-readable monotonic checkpoint identifier.
func GenerateID(t time.Time) string {
	return fmt.Sprintf("ckpt_%s", t.UTC().Format("20060102_150405"))
}

// CheckpointDir returns the directory path for storing checkpoint snapshots.
func CheckpointDir(bankDir string) string {
	return filepath.Join(bankDir, core.StateDirname, "checkpoints")
}

// CreateSnapshot captures current git and workspace state and writes it atomically to disk.
func CreateSnapshot(bankDir, repoDir string, focus string, nextSteps []string, manualEdits []string) (*types.Checkpoint, error) {
	if err := git.ValidateRepoSafety(repoDir); err != nil {
		return nil, fmt.Errorf("safety check failed: %w", err)
	}

	now := time.Now().UTC()
	ckptID := GenerateID(now)

	branch, _ := git.GetCurrentBranch(repoDir)
	headHash, headMsg, _ := git.GetHeadCommit(repoDir)
	dirtyFiles, _ := git.GetStatus(repoDir)
	diffStat, _ := git.GetDiffStat(repoDir)

	ckpt := &types.Checkpoint{
		ID:            ckptID,
		Timestamp:     now,
		Branch:        branch,
		HeadCommit:    headHash,
		CommitMessage: headMsg,
		DirtyFiles:    dirtyFiles,
		DiffStat:      diffStat,
		ActiveFocus:   focus,
		NextSteps:     nextSteps,
		ManualEdits:   manualEdits,
	}

	ckptPath := filepath.Join(CheckpointDir(bankDir), fmt.Sprintf("%s.json", ckptID))
	data, err := json.MarshalIndent(ckpt, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to serialize checkpoint: %w", err)
	}
	data = append(data, '\n')

	if err := core.WriteAtomic(ckptPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write checkpoint snapshot: %w", err)
	}

	// Update manifest file references
	manifest, err := core.LoadManifest(bankDir)
	if err == nil && manifest != nil {
		for filename := range manifest.Files {
			fullPath := filepath.Join(bankDir, filename)
			if hashStr, bytesCount, err := core.ComputeFileHash(fullPath); err == nil {
				meta := manifest.Files[filename]
				meta.SHA256 = hashStr
				meta.Bytes = bytesCount
				meta.LastCheckpoint = ckptID
				meta.UpdatedAt = now
				manifest.Files[filename] = meta
			}
		}
		_ = core.SaveManifest(bankDir, manifest)
	}

	return ckpt, nil
}

// LoadSnapshot reads a specific checkpoint snapshot from disk.
func LoadSnapshot(bankDir, ckptID string) (*types.Checkpoint, error) {
	path := filepath.Join(CheckpointDir(bankDir), fmt.Sprintf("%s.json", ckptID))
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read checkpoint at %s: %w", path, err)
	}

	var ckpt types.Checkpoint
	if err := json.Unmarshal(data, &ckpt); err != nil {
		return nil, fmt.Errorf("failed to parse checkpoint JSON: %w", err)
	}
	return &ckpt, nil
}

// ListSnapshots returns all stored checkpoint IDs ordered chronologically.
func ListSnapshots(bankDir string) ([]string, error) {
	dir := CheckpointDir(bankDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var ids []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			name := entry.Name()
			id := name[:len(name)-len(".json")]
			ids = append(ids, id)
		}
	}
	return ids, nil
}
