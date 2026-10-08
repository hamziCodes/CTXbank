package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ctxbank/ctx/internal/audit"
	"github.com/ctxbank/ctx/internal/checkpoint"
	"github.com/ctxbank/ctx/internal/core"
	"github.com/ctxbank/ctx/internal/git"
	"github.com/ctxbank/ctx/internal/linter"
)

const (
	PromptSyncStateFilename = "prompt_sync.json"
	PromptRelPath           = ".context/prompts/agent_sync_prompt.md"
)

// PromptSyncLedger tracks the persistent history and state of Prompt Sync executions.
type PromptSyncLedger struct {
	TotalSyncs       int       `json:"total_syncs"`
	LastSyncTime     time.Time `json:"last_sync_time"`
	LastCheckpointID string    `json:"last_checkpoint_id"`
	Status           string    `json:"status"` // "idle", "awaiting_agent", "verified", "issues_detected"
	LastAgentResult  string    `json:"last_agent_result,omitempty"`
	LastPromptPath   string    `json:"last_prompt_path,omitempty"`
}

// SyncPromptResult contains the output of a prompt generation run.
type SyncPromptResult struct {
	CheckpointID string `json:"checkpoint_id"`
	PromptFile   string `json:"prompt_file"`
	PromptText   string `json:"prompt_text"`
	TotalSyncs   int    `json:"total_syncs"`
	Status       string `json:"status"`
}

// VerificationResult contains the audit outcome of checking the agent's populated memory bank.
type VerificationResult struct {
	Verified    bool           `json:"verified"`
	Status      string         `json:"status"`
	Issues      []string       `json:"issues"`
	FileDetails map[string]int `json:"file_details"` // filename -> line count
}

// GetLedger loads the prompt sync ledger from the memory bank state directory.
func GetLedger(bankDir string) (*PromptSyncLedger, error) {
	statePath := filepath.Join(bankDir, core.StateDirname, PromptSyncStateFilename)
	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &PromptSyncLedger{
				TotalSyncs: 0,
				Status:     "idle",
			}, nil
		}
		return nil, fmt.Errorf("failed to read prompt sync ledger: %w", err)
	}

	var ledger PromptSyncLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		return &PromptSyncLedger{
			TotalSyncs: 0,
			Status:     "idle",
		}, nil
	}
	return &ledger, nil
}

// SaveLedger atomically saves the prompt sync ledger to disk.
func SaveLedger(bankDir string, ledger *PromptSyncLedger) error {
	statePath := filepath.Join(bankDir, core.StateDirname, PromptSyncStateFilename)
	data, err := json.MarshalIndent(ledger, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal prompt sync ledger: %w", err)
	}
	return core.WriteAtomic(statePath, data, 0644)
}

// GenerateSyncPrompt creates an automated safety checkpoint, scans the codebase, compiles a tailored
// AI Directive Prompt, updates the sync ledger, and saves the prompt for the user's AI coding agent.
func GenerateSyncPrompt(repoDir, bankDir string) (*SyncPromptResult, error) {
	// 1. Ensure Memory Bank is initialized if missing
	if _, err := os.Stat(bankDir); os.IsNotExist(err) {
		if err := core.InitBank(repoDir, false); err != nil {
			return nil, fmt.Errorf("failed to initialize memory bank: %w", err)
		}
	}

	// 2. Automated Safety Pre-Sync Checkpoint
	ckptID := "none"
	focus := "Pre-Prompt-Sync Snapshot: Automated safety backup before AI agent populates memory bank"
	nextSteps := []string{
		"Provide generated prompt to AI coding agent (Cursor, Claude Code, Antigravity)",
		"Verify memory bank updates via 'ctx prompt-sync --verify' or Web UI",
	}
	ckpt, err := checkpoint.CreateSnapshot(bankDir, repoDir, focus, nextSteps, nil)
	if err == nil && ckpt != nil {
		ckptID = ckpt.ID
	}

	// 3. Reconnaissance Audit
	report, _ := audit.RunAudit(repoDir)
	if report != nil {
		_ = audit.ApplyAudit(bankDir, report)
	}

	// 4. Gather Git Context
	branch, _ := git.GetCurrentBranch(repoDir)
	headHash, headMsg, _ := git.GetHeadCommit(repoDir)

	// 5. Build AI Directive Prompt
	projectName := filepath.Base(repoDir)
	promptText := compileAgentPrompt(projectName, repoDir, branch, headHash, headMsg, report)

	// 6. Write prompt to disk in repoDir/.context/prompts/ and bankDir/.state/
	promptRelPath := filepath.Join(".context", "prompts", "agent_sync_prompt.md")
	promptFullPath := filepath.Join(repoDir, promptRelPath)
	if err := core.WriteAtomic(promptFullPath, []byte(promptText), 0644); err != nil {
		// Fallback to bankDir/.state
		promptFullPath = filepath.Join(bankDir, core.StateDirname, "agent_sync_prompt.md")
		_ = core.WriteAtomic(promptFullPath, []byte(promptText), 0644)
	}

	// Also keep a copy inside memory bank state
	_ = core.WriteAtomic(filepath.Join(bankDir, core.StateDirname, "agent_sync_prompt.md"), []byte(promptText), 0644)

	// 7. Update Ledger
	ledger, _ := GetLedger(bankDir)
	ledger.TotalSyncs++
	ledger.LastSyncTime = time.Now().UTC()
	ledger.LastCheckpointID = ckptID
	ledger.Status = "awaiting_agent"
	ledger.LastPromptPath = promptFullPath
	_ = SaveLedger(bankDir, ledger)

	return &SyncPromptResult{
		CheckpointID: ckptID,
		PromptFile:   promptFullPath,
		PromptText:   promptText,
		TotalSyncs:   ledger.TotalSyncs,
		Status:       ledger.Status,
	}, nil
}

// VerifySyncCompletion checks whether the AI agent has populated the memory bank with non-generic,
// valid markdown adhering to the strict <150 lines token budget.
func VerifySyncCompletion(bankDir string) (*VerificationResult, error) {
	files := []string{
		"projectbrief.md",
		"productContext.md",
		"systemPatterns.md",
		"techContext.md",
		"activeContext.md",
		"progress.md",
		"decisionLog.md",
	}

	fileDetails := make(map[string]int)
	var issues []string

	// Check with linter
	lintRes, err := linter.LintMemoryBank(bankDir, false)
	if err == nil && lintRes != nil {
		for _, v := range lintRes.Violations {
			issues = append(issues, fmt.Sprintf("[Violation] %s", v))
		}
	}

	// Scan each file for content substance and generic templates
	for _, f := range files {
		path := filepath.Join(bankDir, f)
		data, err := os.ReadFile(path)
		if err != nil {
			issues = append(issues, fmt.Sprintf("Missing required file: %s", f))
			continue
		}

		content := string(data)
		lines := strings.Split(content, "\n")
		lineCount := len(lines)
		fileDetails[f] = lineCount

		if lineCount < 6 {
			issues = append(issues, fmt.Sprintf("%s is too short (%d lines); expected full architectural documentation", f, lineCount))
		}

		if lineCount > 150 {
			issues = append(issues, fmt.Sprintf("%s exceeds the 150-line token budget (%d lines)", f, lineCount))
		}

		// Check for unreplaced generic template phrases
		genericPhrases := []string{
			"High-level description of the project vision",
			"Replaces brittle prompt-engineering conventions",
			"Generic framework template",
			"Replace this with your project details",
		}
		for _, phrase := range genericPhrases {
			if strings.Contains(content, phrase) {
				issues = append(issues, fmt.Sprintf("%s still contains generic placeholder: %q", f, phrase))
				break
			}
		}
	}

	ledger, _ := GetLedger(bankDir)
	verified := len(issues) == 0
	status := "verified"
	if !verified {
		status = "issues_detected"
	}

	ledger.Status = status
	if verified {
		ledger.LastAgentResult = fmt.Sprintf("All %d memory files verified compliant at %s", len(files), time.Now().UTC().Format(time.RFC3339))
	} else {
		ledger.LastAgentResult = fmt.Sprintf("%d issues detected during verification at %s", len(issues), time.Now().UTC().Format(time.RFC3339))
	}
	_ = SaveLedger(bankDir, ledger)

	return &VerificationResult{
		Verified:    verified,
		Status:      status,
		Issues:      issues,
		FileDetails: fileDetails,
	}, nil
}

func compileAgentPrompt(projectName, repoDir, branch, headHash, headMsg string, report *audit.AuditReport) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# AI DIRECTIVE: Populate CTXbank Memory Bank for `%s`\n\n", projectName))
	sb.WriteString("> **Context:** This project uses **CTXbank** as its deterministic, crash-safe memory layer.\n")
	sb.WriteString("> You have been provided with this directive prompt because a **Prompt Sync** was executed.\n")
	sb.WriteString("> Your mission is to explore the codebase and populate all 7 files in `memory-bank/` with deep, accurate, project-specific technical knowledge.\n\n")

	sb.WriteString("## 1. Project Reconnaissance Findings\n")
	sb.WriteString(fmt.Sprintf("- **Repository Name:** `%s`\n", projectName))
	sb.WriteString(fmt.Sprintf("- **Root Directory:** `%s`\n", repoDir))
	if branch != "" {
		sb.WriteString(fmt.Sprintf("- **Git Branch:** `%s`\n", branch))
	}
	if headHash != "" {
		sb.WriteString(fmt.Sprintf("- **Head Commit:** `%s` (%s)\n", headHash, headMsg))
	}

	if report != nil {
		if len(report.Manifests) > 0 {
			sb.WriteString("\n### Detected Ecosystem & Tech Stack\n")
			for _, m := range report.Manifests {
				sb.WriteString(fmt.Sprintf("- **Language:** %s | **Build System:** %s\n", m.Language, m.BuildSystem))
				if len(m.Dependencies) > 0 {
					limit := 10
					if len(m.Dependencies) < limit {
						limit = len(m.Dependencies)
					}
					sb.WriteString(fmt.Sprintf("  - Key Dependencies: %s\n", strings.Join(m.Dependencies[:limit], ", ")))
				}
				if len(m.Scripts) > 0 {
					sb.WriteString("  - Commands:\n")
					for k, v := range m.Scripts {
						sb.WriteString(fmt.Sprintf("    - `%s`: `%s`\n", k, v))
					}
				}
			}
		}

		if len(report.Components) > 0 {
			sb.WriteString("\n### Discovered Architecture Components\n")
			for _, c := range report.Components {
				sb.WriteString(fmt.Sprintf("- **%s** (`%s`): %s\n", c.Name, c.Path, c.Description))
			}
		}

		if len(report.Symbols) > 0 {
			sb.WriteString("\n### Discovered Symbols & Entrypoints\n")
			count := 0
			for _, s := range report.Symbols {
				if s.Kind == "entrypoint" || s.Kind == "route" {
					sb.WriteString(fmt.Sprintf("- `%s` (%s:%d) [%s]\n", s.Name, s.File, s.Line, s.Kind))
					count++
					if count >= 12 {
						break
					}
				}
			}
		}
	}

	sb.WriteString("\n## 2. Memory Bank Structure & Schema Requirements\n")
	sb.WriteString("You MUST update or populate each of the following 7 files in the `memory-bank/` directory. Each file must be concise, factual, and strictly under 150 lines:\n\n")

	sb.WriteString("### 1. `memory-bank/projectbrief.md` (< 100 lines)\n")
	sb.WriteString("- Core problem statement and purpose of this specific application.\n")
	sb.WriteString("- Primary deliverables and architectural goals.\n")
	sb.WriteString("- Key constraints and explicit non-goals.\n\n")

	sb.WriteString("### 2. `memory-bank/productContext.md` (< 120 lines)\n")
	sb.WriteString("- Why this product exists from a user and business perspective.\n")
	sb.WriteString("- Target audience, user personas, and primary use cases.\n")
	sb.WriteString("- Core user experience (UX) flows and design expectations.\n\n")

	sb.WriteString("### 3. `memory-bank/systemPatterns.md` (< 150 lines)\n")
	sb.WriteString("- System architecture diagrams and structural patterns (MVC, Hexagonal, Clean, Event-driven, etc.).\n")
	sb.WriteString("- Key components, boundaries, data flow, and package relationships.\n")
	sb.WriteString("- Critical algorithms or invariants that must never be broken.\n\n")

	sb.WriteString("### 4. `memory-bank/techContext.md` (< 150 lines)\n")
	sb.WriteString("- Complete technologies used, runtime versions, and core packages.\n")
	sb.WriteString("- Build, test, lint, and run commands.\n")
	sb.WriteString("- Local environment prerequisites and developer configuration.\n\n")

	sb.WriteString("### 5. `memory-bank/activeContext.md` (CRITICAL HOT CONTEXT — STRICT < 150 lines)\n")
	sb.WriteString("- Current sprint / active development focus.\n")
	sb.WriteString("- Recent milestones or completed features (last 3 items).\n")
	sb.WriteString("- Immediate next 3-5 technical steps.\n")
	sb.WriteString("- Active architectural decisions or dilemmas.\n\n")

	sb.WriteString("### 6. `memory-bank/progress.md` (< 150 lines)\n")
	sb.WriteString("- Complete roadmap status: what is working vs what is pending.\n")
	sb.WriteString("- Milestone completion checklist.\n")
	sb.WriteString("- Known issues, bugs, or technical debt.\n\n")

	sb.WriteString("### 7. `memory-bank/decisionLog.md` (< 150 lines)\n")
	sb.WriteString("- Architecture Decision Records (ADRs) explaining key design choices.\n")
	sb.WriteString("- Context, options considered, and rationale for decisions made.\n\n")

	sb.WriteString("## 3. Strict Operating Rules\n")
	sb.WriteString("1. **No Placeholders:** Remove any remaining generic placeholder text or template stubs. Every line must describe THIS codebase.\n")
	sb.WriteString("2. **Budget Discipline:** No file may exceed 150 lines. Keep descriptions dense, structured, and information-rich.\n")
	sb.WriteString("3. **Atomic Writes:** Directly write/update files inside `memory-bank/`.\n")
	sb.WriteString("4. **Zero Drift:** Ensure consistency across all 7 files.\n\n")
	sb.WriteString("Proceed now by analyzing the codebase and updating the 7 memory bank files.\n")

	return sb.String()
}
