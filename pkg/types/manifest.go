package types

import "time"

// VolatilityClass defines how frequently a file changes and its re-read policy.
type VolatilityClass string

const (
	VolatilityStatic     VolatilityClass = "static"      // projectbrief.md, productContext.md
	VolatilitySemiStatic VolatilityClass = "semi-static" // systemPatterns.md, techContext.md, decisionLog.md
	VolatilityHot        VolatilityClass = "hot"         // activeContext.md, progress.md
)

// FileMeta holds cryptographic and volatility metadata for a managed memory bank file.
type FileMeta struct {
	SHA256         string          `json:"sha256"`
	Bytes          int64           `json:"bytes"`
	LastCheckpoint string          `json:"last_ckpt"`
	UpdatedAt      time.Time       `json:"updated_at,omitempty"`
	Volatility     VolatilityClass `json:"volatility,omitempty"`
}

// Manifest represents the single source of truth (.state/manifest.json)
// for file hashes, byte counts, and per-agent session read cursors.
type Manifest struct {
	Version          string                       `json:"version"`
	Files            map[string]FileMeta          `json:"files"`
	AgentReadCursors map[string]map[string]string `json:"agent_read_cursors"`
}

// NewManifest initializes an empty, ready-to-use Manifest struct.
func NewManifest() *Manifest {
	return &Manifest{
		Version:          "1.0.0",
		Files:            make(map[string]FileMeta),
		AgentReadCursors: make(map[string]map[string]string),
	}
}
