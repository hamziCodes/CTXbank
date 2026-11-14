package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hamziCodes/CTXbank/internal/core"
	"github.com/hamziCodes/CTXbank/internal/git"
	"github.com/hamziCodes/CTXbank/pkg/types"
)

// GenerateID produces a deterministic, human-readable monotonic checkpoint identifier.
func GenerateID(t time.Time) string {
	return fmt.Sprintf("ckpt_%s", t.UTC().Format("20060102_150405"))
}

// CheckpointDir returns the directory path for storing checkpoint snapshots.
func CheckpointDir(bankDir string) string {
	return filepath.Join(bankDir, core.StateDirname, "checkpoints")
}

// uniqueCheckpointID guards against same-second collisions (e.g. an MCP
// request_checkpoint racing a manual `ctx pause`): it appends a monotonic
// suffix until the snapshot path is free. Checkpoints must never silently
// overwrite each other.
func uniqueCheckpointID(bankDir, baseID string) string {
	ckptID := baseID
	for n := 2; ; n++ {
		ckptPath := filepath.Join(CheckpointDir(bankDir), fmt.Sprintf("%s.json", ckptID))
		if _, err := os.Stat(ckptPath); os.IsNotExist(err) {
			return ckptID
		}
		ckptID = fmt.Sprintf("%s-%d", baseID, n)
	}
}

// CreateSnapshot captures current git and workspace state and writes it atomically to disk.
func CreateSnapshot(bankDir, repoDir string, focus string, nextSteps []string, manualEdits []string) (*types.Checkpoint, error) {
	if err := git.ValidateRepoSafety(repoDir); err != nil {
		return nil, fmt.Errorf("safety check failed: %w", err)
	}

	now := time.Now().UTC()
	ckptID := uniqueCheckpointID(bankDir, GenerateID(now))

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

	// Close the pause -> resume loop: resume/status/list/dashboard all read
	// Focus/Recent/Next steps from activeContext.md, so every checkpoint
	// must be reflected there. Best-effort: a checkpoint must never fail
	// just because activeContext.md is missing or hand-edited.
	updateActiveContext(bankDir, ckpt)

	// Update manifest file references. A corrupt or missing manifest is
	// rebuilt from the current files instead of staying broken forever
	// (`ctx doctor` still reports the corruption until a checkpoint runs).
	manifest, err := core.LoadManifest(bankDir)
	if err != nil || manifest == nil {
		manifest = types.NewManifest()
		// Seed the rebuild with the bank files that actually exist on disk;
		// a fresh manifest starts with an empty file map.
		for _, name := range []string{"projectbrief.md", "productContext.md", "systemPatterns.md", "techContext.md", "decisionLog.md", "activeContext.md", "progress.md"} {
			if _, statErr := os.Stat(filepath.Join(bankDir, name)); statErr == nil {
				manifest.Files[name] = types.FileMeta{}
			}
		}
	}
	{
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

// updateActiveContext reflects a checkpoint into memory-bank/activeContext.md:
// the title stamp, the Focus section, the Recent bullet list (newest first,
// capped at 3), and the Next steps list when the checkpoint carries one.
// Other sections (e.g. Open decisions) are preserved verbatim. The bootstrap
// placeholder bullet is dropped once a real checkpoint exists.
func updateActiveContext(bankDir string, ckpt *types.Checkpoint) {
	path := filepath.Join(bankDir, "activeContext.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	type section struct {
		heading string // "" for the preamble (title block)
		body    []string
	}
	var sections []section
	cur := section{}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "## ") {
			sections = append(sections, cur)
			cur = section{heading: line}
			continue
		}
		cur.body = append(cur.body, line)
	}
	sections = append(sections, cur)

	note := ckpt.ActiveFocus
	if len(ckpt.ManualEdits) > 0 && ckpt.ManualEdits[0] != "" {
		note = ckpt.ManualEdits[0]
	}
	newBullet := fmt.Sprintf("- %s: %s", ckpt.ID, note)

	for i, sec := range sections {
		switch {
		case sec.heading == "":
			for j, l := range sec.body {
				if strings.HasPrefix(l, "# ") {
					sec.body[j] = fmt.Sprintf("# Active Context — updated %s", ckpt.ID)
					break
				}
			}
		case strings.HasPrefix(sec.heading, "## Focus"):
			if ckpt.ActiveFocus != "" {
				sec.body = []string{"", ckpt.ActiveFocus, ""}
			}
		case strings.HasPrefix(sec.heading, "## Recent"):
			var kept []string
			for _, l := range sec.body {
				t := strings.TrimSpace(l)
				if strings.HasPrefix(t, "- ") && t != "- Initialized memory bank layout." && !strings.HasPrefix(t, "- "+ckpt.ID) {
					kept = append(kept, t)
				}
			}
			bullets := append([]string{newBullet}, kept...)
			if len(bullets) > 3 {
				bullets = bullets[:3]
			}
			sec.body = append([]string{""}, bullets...)
			sec.body = append(sec.body, "")
		case strings.HasPrefix(sec.heading, "## Next steps"):
			if len(ckpt.NextSteps) > 0 {
				sec.body = []string{""}
				for n, step := range ckpt.NextSteps {
					sec.body = append(sec.body, fmt.Sprintf("%d. %s", n+1, step))
				}
				sec.body = append(sec.body, "")
			}
		}
		sections[i] = sec
	}

	var out []string
	for _, sec := range sections {
		if sec.heading != "" {
			out = append(out, sec.heading)
		}
		out = append(out, sec.body...)
	}
	// Trim trailing blank lines, keep exactly one final newline.
	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	_ = core.WriteAtomic(path, []byte(strings.Join(out, "\n")+"\n"), 0644)
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
