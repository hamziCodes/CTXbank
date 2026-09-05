package ingest

import (
	"bufio"
	"strings"
)

// ExtractedChunk represents a parsed section from a research document.
type ExtractedChunk struct {
	Heading string   `json:"heading"`
	Content string   `json:"content"`
	Items   []string `json:"items,omitempty"`
}

// ParseResearchMarkdown parses raw notes, extracts headings and candidate milestones.
func ParseResearchMarkdown(markdown string) []ExtractedChunk {
	var chunks []ExtractedChunk
	scanner := bufio.NewScanner(strings.NewReader(markdown))

	currentHeading := "General Notes"
	var currentLines []string
	var currentItems []string

	flush := func() {
		content := strings.TrimSpace(strings.Join(currentLines, "\n"))
		if content != "" || len(currentItems) > 0 {
			chunks = append(chunks, ExtractedChunk{
				Heading: currentHeading,
				Content: content,
				Items:   currentItems,
			})
		}
		currentLines = nil
		currentItems = nil
	}

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "#") {
			flush()
			currentHeading = strings.TrimLeft(trimmed, "# ")
			continue
		}

		if strings.HasPrefix(trimmed, "- [ ]") || strings.HasPrefix(trimmed, "* [ ]") ||
			strings.HasPrefix(trimmed, "TODO:") || strings.HasPrefix(trimmed, "- ") {
			currentItems = append(currentItems, trimmed)
		}

		currentLines = append(currentLines, line)
	}

	flush()
	return chunks
}
