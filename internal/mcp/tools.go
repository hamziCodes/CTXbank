package mcp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ctxbank/ctx/internal/checkpoint"
	"github.com/ctxbank/ctx/internal/core"
)

var (
	ErrLineBudgetExceeded = errors.New("line count budget exceeded (> 150 lines)")
	ErrInvalidFile        = errors.New("invalid or unauthorized memory bank file")
)

// AllowedFiles restricts writes strictly to memory-bank markdown documents.
var AllowedFiles = map[string]bool{
	"activeContext.md":  true,
	"progress.md":       true,
	"decisionLog.md":    true,
	"systemPatterns.md": true,
	"techContext.md":    true,
	"projectbrief.md":   true,
	"productContext.md": true,
}

// ToolReadActiveContext handles read_active_context(session_id).
func ToolReadActiveContext(bankDir, sessionID string) (map[string]interface{}, error) {
	if sessionID == "" {
		sessionID = "default"
	}

	activeChanged, _, activeContent, err := CheckAndAdvanceCursor(bankDir, sessionID, "activeContext.md")
	if err != nil {
		return nil, err
	}

	progressChanged, _, progressContent, err := CheckAndAdvanceCursor(bankDir, sessionID, "progress.md")
	if err != nil {
		return nil, err
	}

	var unchangedFiles []string
	if !activeChanged {
		unchangedFiles = append(unchangedFiles, "activeContext.md")
	}
	if !progressChanged {
		unchangedFiles = append(unchangedFiles, "progress.md")
	}

	if !activeChanged && !progressChanged {
		return map[string]interface{}{
			"unchanged":       true,
			"unchanged_files": unchangedFiles,
		}, nil
	}

	res := map[string]interface{}{
		"unchanged":       false,
		"unchanged_files": unchangedFiles,
	}

	if activeChanged {
		res["activeContext"] = string(activeContent)
	} else {
		// Read content without cursor advancement
		data, _ := os.ReadFile(filepath.Join(bankDir, "activeContext.md"))
		res["activeContext"] = string(data)
	}

	if progressChanged {
		res["progress"] = string(progressContent)
	} else {
		data, _ := os.ReadFile(filepath.Join(bankDir, "progress.md"))
		res["progress"] = string(data)
	}

	return res, nil
}

// ToolReadStaticContext handles read_static_context(session_id, files).
func ToolReadStaticContext(bankDir, sessionID string, requestedFiles []string) (map[string]interface{}, error) {
	if sessionID == "" {
		sessionID = "default"
	}

	if len(requestedFiles) == 0 {
		requestedFiles = []string{"projectbrief.md", "productContext.md", "systemPatterns.md", "techContext.md"}
	}

	result := make(map[string]interface{})
	var unchangedFiles []string

	for _, filename := range requestedFiles {
		clean := filepath.Base(filename)
		if !AllowedFiles[clean] {
			continue
		}
		changed, _, content, err := CheckAndAdvanceCursor(bankDir, sessionID, clean)
		if err != nil {
			return nil, err
		}
		if !changed {
			unchangedFiles = append(unchangedFiles, clean)
		} else {
			result[clean] = string(content)
		}
	}

	result["unchanged_files"] = unchangedFiles
	return result, nil
}

// ToolAppendMemoryDelta handles append_memory_delta(session_id, file, section, content).
func ToolAppendMemoryDelta(bankDir, sessionID, filename, section, content string) (map[string]interface{}, error) {
	cleanFile := filepath.Base(filename)
	if !AllowedFiles[cleanFile] {
		return nil, fmt.Errorf("%w: %s", ErrInvalidFile, cleanFile)
	}

	targetPath := filepath.Join(bankDir, cleanFile)
	existingBytes, err := os.ReadFile(targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read target file: %w", err)
	}

	existingStr := string(existingBytes)
	existingLines := len(strings.Split(existingStr, "\n"))
	deltaLines := len(strings.Split(content, "\n"))

	// Enforce 150-line budget on activeContext.md
	if cleanFile == "activeContext.md" && (existingLines+deltaLines > 150) {
		return nil, fmt.Errorf("%w: current %d lines + new %d lines exceeds 150 lines limit",
			ErrLineBudgetExceeded, existingLines, deltaLines)
	}

	// Append under named section
	header := fmt.Sprintf("## %s", section)
	var updated string
	if strings.Contains(existingStr, header) {
		idx := strings.Index(existingStr, header) + len(header)
		updated = existingStr[:idx] + "\n" + content + existingStr[idx:]
	} else {
		updated = strings.TrimRight(existingStr, "\r\n") + "\n\n" + header + "\n" + content + "\n"
	}

	if err := core.WriteAtomic(targetPath, []byte(updated), 0644); err != nil {
		return nil, fmt.Errorf("failed to write updated content: %w", err)
	}

	_ = core.UpdateFileMeta(bankDir, cleanFile, "")

	return map[string]interface{}{
		"success": true,
		"file":    cleanFile,
		"section": section,
		"lines":   len(strings.Split(updated, "\n")),
	}, nil
}

// ToolUpdateMilestone handles update_milestone(session_id, milestone, status, next_steps).
func ToolUpdateMilestone(bankDir, sessionID, milestone, status string, nextSteps []string) (map[string]interface{}, error) {
	progressPath := filepath.Join(bankDir, "progress.md")
	data, err := os.ReadFile(progressPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read progress.md: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	updated := false
	targetCheck := "[ ] " + milestone
	newCheck := fmt.Sprintf("[%s] %s", status, milestone)

	for i, l := range lines {
		if strings.Contains(l, targetCheck) {
			lines[i] = strings.Replace(l, targetCheck, newCheck, 1)
			updated = true
			break
		}
	}

	if !updated {
		// Append milestone to list
		lines = append(lines, fmt.Sprintf("- [%s] %s", status, milestone))
	}

	newProgress := strings.Join(lines, "\n")
	if err := core.WriteAtomic(progressPath, []byte(newProgress), 0644); err != nil {
		return nil, err
	}
	_ = core.UpdateFileMeta(bankDir, "progress.md", "")

	// Update Next steps in activeContext.md if provided
	if len(nextSteps) > 0 {
		activePath := filepath.Join(bankDir, "activeContext.md")
		if activeData, err := os.ReadFile(activePath); err == nil {
			activeLines := strings.Split(string(activeData), "\n")
			var rewritten []string
			skipOldNext := false

			for _, line := range activeLines {
				if strings.HasPrefix(strings.TrimSpace(line), "## Next steps") {
					rewritten = append(rewritten, "## Next steps")
					for stepIdx, step := range nextSteps {
						rewritten = append(rewritten, fmt.Sprintf("%d. %s", stepIdx+1, step))
					}
					skipOldNext = true
					continue
				} else if strings.HasPrefix(strings.TrimSpace(line), "## ") {
					skipOldNext = false
				}

				if !skipOldNext {
					rewritten = append(rewritten, line)
				}
			}
			newActive := strings.Join(rewritten, "\n")
			_ = core.WriteAtomic(activePath, []byte(newActive), 0644)
			_ = core.UpdateFileMeta(bankDir, "activeContext.md", "")
		}
	}

	return map[string]interface{}{
		"success":   true,
		"milestone": milestone,
		"status":    status,
	}, nil
}

// ToolReportManualChanges handles report_manual_changes(session_id, files, description).
func ToolReportManualChanges(bankDir, repoDir, sessionID string, files []string, description string) (map[string]interface{}, error) {
	decisionPath := filepath.Join(bankDir, "decisionLog.md")
	entry := fmt.Sprintf("\n## [AGENT-HANDOFF] Mid-Task Manual Changes\n- **Files:** %s\n- **Description:** %s\n",
		strings.Join(files, ", "), description)

	existing, err := os.ReadFile(decisionPath)
	if err != nil {
		existing = []byte("# Decision Log\n")
	}

	updated := append(existing, []byte(entry)...)
	if err := core.WriteAtomic(decisionPath, updated, 0644); err != nil {
		return nil, err
	}
	_ = core.UpdateFileMeta(bankDir, "decisionLog.md", "")

	ckpt, err := checkpoint.CreateSnapshot(bankDir, repoDir, "Agent Handoff Snapshot", nil, []string{description})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":       true,
		"checkpoint_id": ckpt.ID,
		"description":   description,
	}, nil
}

// ToolRequestCheckpoint handles request_checkpoint(session_id, focus).
func ToolRequestCheckpoint(bankDir, repoDir, sessionID, focus string) (map[string]interface{}, error) {
	if focus == "" {
		focus = "Agent Checkpoint Request"
	}

	ckpt, err := checkpoint.CreateSnapshot(bankDir, repoDir, focus, nil, nil)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":       true,
		"checkpoint_id": ckpt.ID,
		"branch":        ckpt.Branch,
		"dirty_count":   len(ckpt.DirtyFiles),
	}, nil
}
