package ingest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ctxbank/ctx/internal/core"
)

// IngestionProposal holds the proposed changes and destination.
type IngestionProposal struct {
	SourceFile   string   `json:"source_file"`
	TargetFile   string   `json:"target_file"`
	ProposedDiff string   `json:"proposed_diff"`
	NewItems     []string `json:"new_items"`
	Deduplicated int      `json:"deduplicated_count"`
}

// PrepareIngestion processes an inbox file against the existing memory bank.
func PrepareIngestion(bankDir, sourcePath string) (*IngestionProposal, error) {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read research file %s: %w", sourcePath, err)
	}

	chunks := ParseResearchMarkdown(string(data))
	targetDoc := "productContext.md"
	targetPath := filepath.Join(bankDir, targetDoc)
	existingData, _ := os.ReadFile(targetPath)
	existingStr := string(existingData)

	var filteredChunks []ExtractedChunk
	var candidateItems []string
	dedupCount := 0

	normalizedStr := strings.ReplaceAll(existingStr, "\r\n", "\n")
	existingParagraphs := strings.Split(normalizedStr, "\n\n")

	for _, chunk := range chunks {
		cleanChunk := strings.TrimSpace(strings.ReplaceAll(chunk.Content, "\r\n", "\n"))
		if cleanChunk == "" && len(chunk.Items) == 0 {
			continue
		}

		isDuplicate := false
		for _, existingPara := range existingParagraphs {
			trimmedPara := stripMarkdownHeaders(existingPara)
			if trimmedPara == "" {
				continue
			}
			if IsNearDuplicate(cleanChunk, trimmedPara, 3) {
				isDuplicate = true
				break
			}
		}

		if isDuplicate {
			dedupCount++
			continue
		}
		filteredChunks = append(filteredChunks, chunk)
		candidateItems = append(candidateItems, chunk.Items...)
	}

	if len(filteredChunks) == 0 {
		return &IngestionProposal{
			SourceFile:   sourcePath,
			TargetFile:   targetDoc,
			Deduplicated: dedupCount,
		}, nil
	}

	// Generate proposed addition block
	var sb strings.Builder
	timestamp := time.Now().UTC().Format("2006-01-02 15:04")
	baseName := filepath.Base(sourcePath)
	sb.WriteString(fmt.Sprintf("\n\n## Research Ingestion: %s (%s)\n", baseName, timestamp))

	for _, chunk := range filteredChunks {
		sb.WriteString(fmt.Sprintf("### %s\n%s\n\n", chunk.Heading, chunk.Content))
	}

	proposalText := sb.String()

	// Unified diff preview
	diff := fmt.Sprintf("--- a/%s\n+++ b/%s\n@@ proposed addition @@\n%s", targetDoc, targetDoc, proposalText)

	return &IngestionProposal{
		SourceFile:   sourcePath,
		TargetFile:   targetDoc,
		ProposedDiff: diff,
		NewItems:     candidateItems,
		Deduplicated: dedupCount,
	}, nil
}

// CommitIngestion applies the proposal and moves the research file to archive.
func CommitIngestion(bankDir, repoDir string, prop *IngestionProposal) error {
	targetPath := filepath.Join(bankDir, prop.TargetFile)
	existing, _ := os.ReadFile(targetPath)

	// Extract proposed addition lines from diff
	lines := strings.Split(prop.ProposedDiff, "\n")
	var additionLines []string
	capture := false
	for _, l := range lines {
		if strings.HasPrefix(l, "@@") {
			capture = true
			continue
		}
		if capture {
			additionLines = append(additionLines, l)
		}
	}

	updated := append(existing, []byte("\n"+strings.Join(additionLines, "\n"))...)
	if err := core.WriteAtomic(targetPath, updated, 0644); err != nil {
		return fmt.Errorf("failed to commit update to %s: %w", prop.TargetFile, err)
	}
	_ = core.UpdateFileMeta(bankDir, prop.TargetFile, "")

	// Move source file to research/archive/
	archiveDir := filepath.Join(repoDir, "research", "archive")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	baseName := filepath.Base(prop.SourceFile)
	archiveName := fmt.Sprintf("%s_%s", time.Now().Format("20060102_150405"), baseName)
	archivePath := filepath.Join(archiveDir, archiveName)

	if err := os.Rename(prop.SourceFile, archivePath); err != nil {
		// Fallback to copy and remove if across partitions
		data, err := os.ReadFile(prop.SourceFile)
		if err != nil {
			return err
		}
		if err := core.WriteAtomic(archivePath, data, 0644); err != nil {
			return err
		}
		_ = os.Remove(prop.SourceFile)
	}

	return nil
}

func stripMarkdownHeaders(text string) string {
	var lines []string
	for _, l := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return strings.Join(lines, " ")
}

