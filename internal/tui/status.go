package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/ctxbank/ctx/pkg/types"
)

// StatusCardData contains all fields needed to render the status card.
type StatusCardData struct {
	ProjectName string
	Branch      string
	DirtyCount  int
	IdleTime    time.Duration
	Focus       string
	NextStep    string
	TestsPass   bool
}

// RenderStatusCard formats the status card into a string according to CTXbank §3.2.
func RenderStatusCard(data StatusCardData) string {
	box := GetBox()

	headerText := fmt.Sprintf(" %s · %s ", data.ProjectName, data.Branch)
	totalWidth := 65
	rightDashes := totalWidth - len(headerText) - len(box.TopLeft) - len(box.TopRight)
	if rightDashes < 2 {
		rightDashes = 2
	}
	topBorder := fmt.Sprintf("%s%s%s%s", box.TopLeft, headerText, strings.Repeat(box.Horizontal, rightDashes), box.TopRight)

	// Line 1: indicators
	var dirtyBadge string
	if data.DirtyCount > 0 {
		dirtyBadge = fmt.Sprintf("%s %d files dirty", box.Bullet, data.DirtyCount)
	} else {
		dirtyBadge = fmt.Sprintf("%s working tree clean", box.Check)
	}

	idleDays := int(data.IdleTime.Hours() / 24)
	idleStr := fmt.Sprintf("idle %dd", idleDays)
	if idleDays == 0 {
		hours := int(data.IdleTime.Hours())
		if hours == 0 {
			idleStr = "idle <1h"
		} else {
			idleStr = fmt.Sprintf("idle %dh", hours)
		}
	}

	metricLine := fmt.Sprintf(" %s   ⏱ %s", dirtyBadge, idleStr)
	metricPadding := totalWidth - len(metricLine) - 3
	if metricPadding < 0 {
		metricPadding = 0
	}
	line1 := fmt.Sprintf("%s%s%s %s", box.Vertical, metricLine, strings.Repeat(" ", metricPadding), box.Vertical)

	// Line 2: Focus
	focusLine := fmt.Sprintf(" Focus: %s", data.Focus)
	focusPadding := totalWidth - len(focusLine) - 3
	if focusPadding < 0 {
		focusPadding = 0
	}
	line2 := fmt.Sprintf("%s%s%s %s", box.Vertical, focusLine, strings.Repeat(" ", focusPadding), box.Vertical)

	// Line 3: Next
	nextLine := fmt.Sprintf(" Next:  %s", data.NextStep)
	nextPadding := totalWidth - len(nextLine) - 3
	if nextPadding < 0 {
		nextPadding = 0
	}
	line3 := fmt.Sprintf("%s%s%s %s", box.Vertical, nextLine, strings.Repeat(" ", nextPadding), box.Vertical)

	// Bottom border
	bottomBorder := fmt.Sprintf("%s%s%s", box.BottomLeft, strings.Repeat(box.Horizontal, totalWidth-4), box.BottomRight)

	return strings.Join([]string{topBorder, line1, line2, line3, bottomBorder}, "\n")
}

// FormatCheckpointSummary formats a brief string for CLI output.
func FormatCheckpointSummary(ckpt *types.Checkpoint) string {
	return fmt.Sprintf("Snapshot created: %s (Branch: %s, Dirty files: %d)", ckpt.ID, ckpt.Branch, len(ckpt.DirtyFiles))
}
