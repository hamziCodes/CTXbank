package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ctxbank/ctx/internal/checkpoint"
)

// ResourceDescriptor defines an MCP resource item.
type ResourceDescriptor struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

// ListResources returns the supported static and dynamic memory bank resources.
func ListResources() []ResourceDescriptor {
	return []ResourceDescriptor{
		{
			URI:         "memory://activeContext",
			Name:        "Active Context",
			Description: "The hot memory bank context file describing active focus and next steps.",
			MimeType:    "text/markdown",
		},
		{
			URI:         "memory://progress",
			Name:        "Milestone Progress",
			Description: "Current status ledger of project milestones and sub-phases.",
			MimeType:    "text/markdown",
		},
		{
			URI:         "memory://checkpoint/latest",
			Name:        "Latest Checkpoint Snapshot",
			Description: "JSON snapshot of latest repository state, branch, commit, and dirty files.",
			MimeType:    "application/json",
		},
		{
			URI:         "memory://decisionLog",
			Name:        "Decision Log",
			Description: "Append-only architectural and operational decision ledger.",
			MimeType:    "text/markdown",
		},
	}
}

// ReadResource reads the content of a memory:// URI.
func ReadResource(bankDir, uri string) (string, string, error) {
	switch uri {
	case "memory://activeContext":
		data, err := os.ReadFile(filepath.Join(bankDir, "activeContext.md"))
		return string(data), "text/markdown", err

	case "memory://progress":
		data, err := os.ReadFile(filepath.Join(bankDir, "progress.md"))
		return string(data), "text/markdown", err

	case "memory://decisionLog":
		data, err := os.ReadFile(filepath.Join(bankDir, "decisionLog.md"))
		return string(data), "text/markdown", err

	case "memory://checkpoint/latest":
		snapshots, err := checkpoint.ListSnapshots(bankDir)
		if err != nil || len(snapshots) == 0 {
			return "{}", "application/json", nil
		}
		latestID := snapshots[len(snapshots)-1]
		ckpt, err := checkpoint.LoadSnapshot(bankDir, latestID)
		if err != nil {
			return "", "", err
		}
		data, err := json.MarshalIndent(ckpt, "", "  ")
		return string(data), "application/json", err

	default:
		return "", "", fmt.Errorf("unknown resource URI: %s", uri)
	}
}
