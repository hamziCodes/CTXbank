package types

import "time"

// GitFileStatus represents the working tree and index status of a file.
type GitFileStatus struct {
	Path           string `json:"path"`
	StagingStatus  string `json:"staging_status"`  // Index status (e.g. M, A, D, ?)
	WorktreeStatus string `json:"worktree_status"` // Worktree status (e.g. M, D, ?)
}

// Checkpoint represents a deterministic, point-in-time snapshot of repository state
// stored in memory-bank/.state/checkpoints/ckpt_<id>.json.
type Checkpoint struct {
	ID            string          `json:"id"`
	Timestamp     time.Time       `json:"timestamp"`
	Branch        string          `json:"branch"`
	HeadCommit    string          `json:"head_commit"`
	CommitMessage string          `json:"commit_message"`
	DirtyFiles    []GitFileStatus `json:"dirty_files"`
	DiffStat      string          `json:"diff_stat"`
	ActiveFocus   string          `json:"active_focus"`
	NextSteps     []string        `json:"next_steps"`
	ManualEdits   []string        `json:"manual_edits"`
}

// CheckpointSummary contains high-level metrics for quick terminal and status cards.
type CheckpointSummary struct {
	ID         string    `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	Branch     string    `json:"branch"`
	DirtyCount int       `json:"dirty_count"`
	Focus      string    `json:"focus"`
}
