package ingest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hamziCodes/CTXbank/internal/core"
)

// IngestionProposal holds the proposed changes and destination.
type IngestionProposal struct {
	SourceFile   string   `json:"source_file"`
	TargetFile   string   `json:"target_file"`
	ProposedDiff string   `json:"proposed_diff"`
	NewItems     []string `json:"new_items"`
	Deduplicated int      `json:"deduplicated_count"`
}

// classifyTarget picks the memory file a research note belongs in, using
// deterministic keyword scoring over the whole document. Tech content goes
// to techContext.md, explicit decisions to decisionLog.md, architecture to
// systemPatterns.md, and product/user content to productContext.md (the
// historical default when nothing scores).
func classifyTarget(content string) string {
	lower := strings.ToLower(content)
	// A note that LEADS with an explicit decision is a decision, however
	// technical its subject is: file it in the decision log. Headings are
	// skipped; the first two real content lines get this fast path.
	contentLines := 0
	for _, line := range strings.Split(lower, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		contentLines++
		if strings.Contains(t, "decided") || strings.Contains(t, "decision:") || strings.Contains(t, "chose") {
			return "decisionLog.md"
		}
		if contentLines >= 2 {
			break
		}
	}
	score := func(words []string) int {
		n := 0
		for _, w := range words {
			n += strings.Count(lower, w)
		}
		return n
	}
	decision := score([]string{"we decided", "decision:", "decided to", "chose", "chosen", "adopted", "rejected", "trade-off", "tradeoff", "rationale"})
	tech := score([]string{"postgres", "mysql", "sqlite", "database", "sql", "redis", "jwt", "oauth", "api", "sdk", "cli", "docker", "kubernetes", "golang", "python", "typescript", "javascript", "framework", "library", "latency", "cache", "queue", "server", "endpoint", "migration", "schema"})
	arch := score([]string{"architecture", "pattern", "layer", "module", "component", "microservice", "monolith", "coupling", "invariant"})
	product := score([]string{"user", "customer", "audience", "experience", "onboarding", "pricing", "market", "persona"})

	best, target := product, "productContext.md"
	if arch > best {
		best, target = arch, "systemPatterns.md"
	}
	if tech > best {
		best, target = tech, "techContext.md"
	}
	if decision > best {
		target = "decisionLog.md"
	}
	return target
}

// PrepareIngestion processes an inbox file against the existing memory bank.
func PrepareIngestion(bankDir, sourcePath string) (*IngestionProposal, error) {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read research file %s: %w", sourcePath, err)
	}

	chunks := ParseResearchMarkdown(string(data))
	targetDoc := classifyTarget(string(data))
	targetPath := filepath.Join(bankDir, targetDoc)
	existingData, _ := os.ReadFile(targetPath)
	existingStr := string(existingData)

	var filteredChunks []ExtractedChunk
	var candidateItems []string
	dedupCount := 0

	normalizedStr := strings.ReplaceAll(existingStr, "\r\n", "\n")
	existingParagraphs := strings.Split(normalizedStr, "\n\n")

	// Dedup state, seeded with every line already recorded anywhere in the
	// bank: re-ingesting the same note (or a line already captured in
	// another file) must not commit the same content twice. The same set
	// then also dedupes repeats within this document itself.
	seenLines := map[string]bool{}
	if entries, err := os.ReadDir(bankDir); err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			if b, err := os.ReadFile(filepath.Join(bankDir, e.Name())); err == nil {
				for _, l := range strings.Split(string(b), "\n") {
					if t := strings.TrimSpace(l); t != "" {
						seenLines[t] = true
					}
				}
			}
		}
	}
	var keptTexts []string

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
		if !isDuplicate {
			for _, kept := range keptTexts {
				if IsNearDuplicate(cleanChunk, kept, 3) {
					isDuplicate = true
					break
				}
			}
		}

		if isDuplicate {
			dedupCount++
			continue
		}

		// Drop exact duplicate lines already kept from this same document.
		var dedupedLines []string
		for _, line := range strings.Split(cleanChunk, "\n") {
			key := strings.TrimSpace(line)
			if key != "" {
				if seenLines[key] {
					dedupCount++
					continue
				}
				seenLines[key] = true
			}
			dedupedLines = append(dedupedLines, line)
		}
		chunk.Content = strings.TrimSpace(strings.Join(dedupedLines, "\n"))
		if chunk.Content == "" && len(chunk.Items) == 0 {
			continue
		}
		keptTexts = append(keptTexts, chunk.Content)
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

// CommitIngestion applies the proposal and archives a COPY of the research
// file under research/archive/. The original is left in place: ingestion
// must never delete or move the user's only copy of a note.
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

	// Archive a copy of the source file; the original stays where it was.
	archiveDir := filepath.Join(repoDir, "research", "archive")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	baseName := filepath.Base(prop.SourceFile)
	archiveName := fmt.Sprintf("%s_%s", time.Now().Format("20060102_150405"), baseName)
	archivePath := filepath.Join(archiveDir, archiveName)

	data, err := os.ReadFile(prop.SourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file for archiving: %w", err)
	}
	if err := core.WriteAtomic(archivePath, data, 0644); err != nil {
		return fmt.Errorf("failed to archive source copy: %w", err)
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
