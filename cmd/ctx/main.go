package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ctxbank/ctx/internal/audit"
	"github.com/ctxbank/ctx/internal/checkpoint"
	"github.com/ctxbank/ctx/internal/core"
	"github.com/ctxbank/ctx/internal/git"
	"github.com/ctxbank/ctx/internal/ingest"
	"github.com/ctxbank/ctx/internal/linter"
	"github.com/ctxbank/ctx/internal/llm"
	"github.com/ctxbank/ctx/internal/mcp"
	"github.com/ctxbank/ctx/internal/rules"
	"github.com/ctxbank/ctx/internal/tui"
	"github.com/ctxbank/ctx/internal/workspace"
)

const Version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	command := os.Args[1]
	switch command {
	case "init":
		runInit(os.Args[2:])
	case "status":
		runStatus(os.Args[2:])
	case "list":
		runList(os.Args[2:])
	case "pause", "checkpoint":
		runPause(os.Args[2:])
	case "resume":
		runResume(os.Args[2:])
	case "audit":
		runAudit(os.Args[2:])
	case "ingest":
		runIngest(os.Args[2:])
	case "lint-memory", "lint":
		runLintMemory(os.Args[2:])
	case "serve":
		runServe(os.Args[2:])
	case "--version", "-v", "version":
		fmt.Printf("ctx version %s\n", Version)
	case "--help", "-h", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`CTXbank — The Deterministic Memory Bank & Project Lifecycle Engine (v%s)

Usage:
  ctx <command> [arguments]

Commands:
  init            Initialize memory-bank/ layout and minimal vendor rules
  status [--json] Display single-shot project health and active context
  list [dir]      List and inspect all CTXbank workspaces under directory
  pause           Safe checkpoint + capture manual out-of-band changes
  resume          Display instant pickup brief for human or agent
  audit [--apply] Run 4-stage brownfield reconnaissance scan
  ingest <path>   Ingest research notes into memory-bank with diff review
  lint-memory     Enforce line budgets (< 150 lines) and memory integrity
  serve --mcp     Run the Model Context Protocol (MCP) server over stdio

Flags:
  --help, -h      Show this help message
  --version, -v   Print version information
`, Version)
}

func runInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	force := fs.Bool("force", false, "Force re-initialization even if memory-bank exists")
	_ = fs.Parse(args)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error determining working directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Initializing CTXbank memory bank...")
	if err := core.InitBank(cwd, *force); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing memory bank: %v\n", err)
		os.Exit(1)
	}

	stubs, err := rules.GenerateStubs(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to generate some vendor stubs: %v\n", err)
	} else if len(stubs) > 0 {
		fmt.Printf("Generated minimal vendor stubs: %s\n", strings.Join(stubs, ", "))
	}

	fmt.Println("Initialization complete: memory-bank/ scaffolded with crash-safe atomic engine.")
}

func runStatus(args []string) {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	jsonOutput := fs.Bool("json", false, "Output machine-readable JSON status")
	_ = fs.Parse(args)

	isJSON := *jsonOutput
	for _, a := range args {
		if a == "--json" {
			isJSON = true
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	bankDir := filepath.Join(cwd, core.MemoryBankDir)
	projectName := filepath.Base(cwd)
	branch, _ := git.GetCurrentBranch(cwd)
	dirtyFiles, _ := git.GetStatus(cwd)

	// Determine idle time and active focus from activeContext.md
	activeContextPath := filepath.Join(bankDir, "activeContext.md")
	var focus = "No active focus declared"
	var nextStep = "Define next milestone"
	var idleDuration = time.Duration(0)

	if info, err := os.Stat(activeContextPath); err == nil {
		idleDuration = time.Since(info.ModTime())
		if content, err := os.ReadFile(activeContextPath); err == nil {
			lines := strings.Split(string(content), "\n")
			inFocus := false
			inNext := false
			for _, l := range lines {
				trimmed := strings.TrimSpace(l)
				if strings.HasPrefix(trimmed, "## Focus") {
					inFocus = true
					inNext = false
					continue
				} else if strings.HasPrefix(trimmed, "## Next steps") {
					inNext = true
					inFocus = false
					continue
				} else if strings.HasPrefix(trimmed, "## ") {
					inFocus = false
					inNext = false
				}

				if inFocus && trimmed != "" && !strings.HasPrefix(trimmed, "#") {
					focus = trimmed
					inFocus = false
				}
				if inNext && trimmed != "" && !strings.HasPrefix(trimmed, "#") {
					nextStep = strings.TrimLeft(trimmed, "0123456789.- ")
					inNext = false
				}
			}
		}
	}

	if isJSON {
		payload := map[string]interface{}{
			"project":       projectName,
			"branch":        branch,
			"dirty_count":   len(dirtyFiles),
			"dirty_files":   dirtyFiles,
			"idle_seconds":  int(idleDuration.Seconds()),
			"focus":         focus,
			"next_step":     nextStep,
			"has_active_ctx": idleDuration > 0,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(payload)
		return
	}

	cardData := tui.StatusCardData{
		ProjectName: projectName,
		Branch:      branch,
		DirtyCount:  len(dirtyFiles),
		IdleTime:    idleDuration,
		Focus:       focus,
		NextStep:    nextStep,
	}

	fmt.Println(tui.RenderStatusCard(cardData))
}

func runPause(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	bankDir := filepath.Join(cwd, core.MemoryBankDir)
	if _, err := os.Stat(bankDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: memory-bank not found. Run 'ctx init' first.\n")
		os.Exit(1)
	}

	// Prompt for manual change description
	fmt.Print("Snapshotting... detecting changes...\n")
	dirty, _ := git.GetStatus(cwd)
	var manualNotes []string
	if len(dirty) > 0 {
		fmt.Printf("Detected %d modified file(s).\n", len(dirty))
		fmt.Print("  -> What did you change here? (one line, Enter to skip): ")
		reader := bufio.NewReader(os.Stdin)
		note, _ := reader.ReadString('\n')
		note = strings.TrimSpace(note)
		if note != "" {
			manualNotes = append(manualNotes, note)
		}
	}

	ckpt, err := checkpoint.CreateSnapshot(bankDir, cwd, "Manual Pause Checkpoint", nil, manualNotes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Checkpoint failed: %v\n", err)
		os.Exit(1)
	}

	// Append note to decisionLog if provided
	if len(manualNotes) > 0 {
		decisionPath := filepath.Join(bankDir, "decisionLog.md")
		logEntry := fmt.Sprintf("\n## [%s] Developer Manual Pause\n- **Changes:** %s\n", ckpt.ID, manualNotes[0])
		if existing, err := os.ReadFile(decisionPath); err == nil {
			updated := append(existing, []byte(logEntry)...)
			_ = core.WriteAtomic(decisionPath, updated, 0644)
			_ = core.UpdateFileMeta(bankDir, "decisionLog.md", ckpt.ID)
		}
	}

	fmt.Printf("Snapshot created: %s (saved to activeContext.md -> decisionLog.md)\n", ckpt.ID)
}

func runResume(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	bankDir := filepath.Join(cwd, core.MemoryBankDir)
	activeContextPath := filepath.Join(bankDir, "activeContext.md")
	content, err := os.ReadFile(activeContextPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No activeContext.md found at %s. Run 'ctx init' first.\n", bankDir)
		os.Exit(1)
	}

	fmt.Printf("=== CTXbank Pickup Brief ===\n\n%s\n", string(content))
}

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	mcpFlag := fs.Bool("mcp", false, "Start Model Context Protocol (MCP) server over stdio")
	_ = fs.Parse(args)

	if !*mcpFlag {
		fmt.Fprintf(os.Stderr, "Usage: ctx serve --mcp\n")
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	bankDir := filepath.Join(cwd, core.MemoryBankDir)
	if _, err := os.Stat(bankDir); os.IsNotExist(err) {
		_ = core.InitBank(cwd, false)
	}

	server := mcp.NewServer(bankDir, cwd)
	if err := server.Serve(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
		os.Exit(1)
	}
}

func runAudit(args []string) {
	fs := flag.NewFlagSet("audit", flag.ExitOnError)
	applyFlag := fs.Bool("apply", false, "Apply discovered architecture and tech stack to memory-bank/")
	_ = fs.Parse(args)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	report, err := audit.RunAudit(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Audit failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(tui.RenderAuditReport(report))

	if *applyFlag {
		bankDir := filepath.Join(cwd, core.MemoryBankDir)
		if err := audit.ApplyAudit(bankDir, report); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to apply audit to memory bank: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Applied audit results to techContext.md and systemPatterns.md.")
	}
}

func runLintMemory(args []string) {
	fs := flag.NewFlagSet("lint-memory", flag.ExitOnError)
	fixFlag := fs.Bool("fix", false, "Automatically prune older bullets if line budget is exceeded")
	_ = fs.Parse(args)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	bankDir := filepath.Join(cwd, core.MemoryBankDir)
	res, err := linter.LintMemoryBank(bankDir, *fixFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Linter error: %v\n", err)
		os.Exit(1)
	}

	box := tui.GetBox()
	if len(res.Warnings) > 0 {
		for _, w := range res.Warnings {
			fmt.Printf("! Warning: %s\n", w)
		}
	}

	if !res.Passed {
		fmt.Printf("%s Memory bank budget violations found:\n", box.Cross)
		for _, v := range res.Violations {
			fmt.Printf("  - %s\n", v)
		}
		os.Exit(1)
	}

	fmt.Printf("%s All memory bank token and format budgets satisfied (< 150 lines).\n", box.Check)
}

func runIngest(args []string) {
	fs := flag.NewFlagSet("ingest", flag.ExitOnError)
	yesFlag := fs.Bool("yes", false, "Automatically accept proposed diff without interactive confirmation")
	_ = fs.Parse(args)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	remaining := fs.Args()
	if len(remaining) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: ctx ingest <path-to-research-note.md> [--yes]\n")
		os.Exit(1)
	}

	sourcePath := remaining[0]
	if !filepath.IsAbs(sourcePath) {
		sourcePath = filepath.Join(cwd, sourcePath)
	}

	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: research file not found at %s\n", sourcePath)
		os.Exit(1)
	}

	bankDir := filepath.Join(cwd, core.MemoryBankDir)
	prop, err := ingest.PrepareIngestion(bankDir, sourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ingestion preparation failed: %v\n", err)
		os.Exit(1)
	}

	if prop.ProposedDiff == "" {
		fmt.Println("All content in source file is already present (SimHash near-duplicate detected). No changes proposed.")
		return
	}

	// Check optional Tier-1 LLM classification
	ollama := llm.NewOllamaClient("", "")
	if ollama.IsAvailable() {
		fmt.Println("Tier-1 Local LLM (Ollama) active: classifying research chunk...")
		if targetDoc, err := ollama.ClassifyChunk(prop.ProposedDiff); err == nil && targetDoc != llm.DocDiscard {
			prop.TargetFile = string(targetDoc)
		}
	} else {
		fmt.Println("Tier-0 Ingestion active (Ollama offline; using deterministic classification).")
	}

	fmt.Printf("\nProposed Diff for %s:\n%s\n", prop.TargetFile, prop.ProposedDiff)

	if !*yesFlag {
		fmt.Print("\nApply this ingestion proposal to memory-bank and archive source? [Y/n]: ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "" && answer != "y" && answer != "yes" {
			fmt.Println("Ingestion aborted. Source file remains in inbox.")
			return
		}
	}

	if err := ingest.CommitIngestion(bankDir, cwd, prop); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to commit ingestion: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Ingestion committed to %s! Source moved to research/archive/.\n", prop.TargetFile)
}

func runList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	jsonOutput := fs.Bool("json", false, "Output machine-readable JSON array of projects")
	depthFlag := fs.Int("depth", 3, "Maximum directory depth to search for memory-bank")
	_ = fs.Parse(args)

	// Position-independent --json check
	isJSON := *jsonOutput
	targetDir := "."
	for _, arg := range args {
		if arg == "--json" {
			isJSON = true
		} else if !strings.HasPrefix(arg, "-") {
			targetDir = arg
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	searchRoot := targetDir
	if !filepath.IsAbs(searchRoot) {
		searchRoot = filepath.Join(cwd, targetDir)
	}

	projects, err := workspace.ScanWorkspaces(searchRoot, *depthFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning workspaces: %v\n", err)
		os.Exit(1)
	}

	if isJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(projects)
		return
	}

	box := tui.GetBox()
	if len(projects) == 0 {
		fmt.Printf("No active CTXbank projects found under %s\n", targetDir)
		return
	}

	fmt.Printf("%s CTXbank Multi-Project Dashboard (%d active workspaces)\n\n", box.Check, len(projects))
	fmt.Printf("%-20s %-15s %-10s %-10s %s\n", "PROJECT", "BRANCH", "DIRTY", "IDLE", "ACTIVE FOCUS")
	fmt.Printf("%-20s %-15s %-10s %-10s %s\n", "-------", "------", "-----", "----", "------------")

	for _, p := range projects {
		dirtyStr := "clean"
		if p.DirtyCount > 0 {
			dirtyStr = fmt.Sprintf("%d files", p.DirtyCount)
		}
		idleStr := fmt.Sprintf("%dd", p.IdleDays)
		if p.IdleDays == 0 {
			idleStr = "<1d"
		}

		focus := p.ActiveFocus
		if len(focus) > 35 {
			focus = focus[:32] + "..."
		}

		fmt.Printf("%-20s %-15s %-10s %-10s %s\n", p.Name, p.Branch, dirtyStr, idleStr, focus)
	}
}


