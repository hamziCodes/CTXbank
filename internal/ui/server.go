package ui

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ctxbank/ctx/internal/audit"
	"github.com/ctxbank/ctx/internal/checkpoint"
	"github.com/ctxbank/ctx/internal/core"
	"github.com/ctxbank/ctx/internal/git"
	"github.com/ctxbank/ctx/internal/ingest"
	"github.com/ctxbank/ctx/internal/linter"
	"github.com/ctxbank/ctx/pkg/types"
)

//go:embed web/*
var embeddedFiles embed.FS

// Server hosts the CTXbank local interactive web dashboard.
type Server struct {
	RepoDir string
	BankDir string
	Port    int
}

// NewServer initializes a dashboard server for the target repository.
func NewServer(repoDir string, port int) *Server {
	if port <= 0 {
		port = 4242
	}
	bankDir := filepath.Join(repoDir, core.MemoryBankDir)
	return &Server{
		RepoDir: repoDir,
		BankDir: bankDir,
		Port:    port,
	}
}

// Start binds to localhost, serves the API and embedded UI, and opens the default browser.
func (s *Server) Start(openBrowser bool) error {
	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/files", s.handleFiles)
	mux.HandleFunc("/api/file", s.handleFileGet)
	mux.HandleFunc("/api/file/save", s.handleFileSave)
	mux.HandleFunc("/api/graph", s.handleGraph)
	mux.HandleFunc("/api/checkpoints", s.handleCheckpoints)
	mux.HandleFunc("/api/checkpoint/create", s.handleCheckpointCreate)
	mux.HandleFunc("/api/audit", s.handleAudit)
	mux.HandleFunc("/api/lint", s.handleLint)
	mux.HandleFunc("/api/ingest/preview", s.handleIngestPreview)
	mux.HandleFunc("/api/ingest/commit", s.handleIngestCommit)

	// Static Web Assets: Prefer local disk in dev mode for hot reload; fallback to embedded in release
	var fileSystem http.FileSystem
	localWebDir := filepath.Join(s.RepoDir, "internal", "ui", "web")
	if fi, err := os.Stat(localWebDir); err == nil && fi.IsDir() {
		fileSystem = http.Dir(localWebDir)
	} else {
		subFS, err := fs.Sub(embeddedFiles, "web")
		if err != nil {
			return fmt.Errorf("failed to locate embedded web directory: %w", err)
		}
		fileSystem = http.FS(subFS)
	}
	mux.Handle("/", http.FileServer(fileSystem))

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", s.Port))
	if err != nil {
		// Fallback to random free port
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return err
		}
	}
	s.Port = listener.Addr().(*net.TCPAddr).Port

	url := fmt.Sprintf("http://localhost:%d", s.Port)
	fmt.Printf("\n[+] CTXbank Interactive Dashboard live at %s\n", url)
	fmt.Println("    Adheres to VERTEX Universal Design (Solid Objects, Tabular Numbers, 4-State Containers)")
	fmt.Println("    Press Ctrl+C to stop the dashboard server.")

	if openBrowser {
		go func() {
			time.Sleep(200 * time.Millisecond)
			openURL(url)
		}()
	}

	return http.Serve(listener, mux)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	branch, _ := git.GetCurrentBranch(s.RepoDir)
	dirtyFiles, _ := git.GetStatus(s.RepoDir)

	activePath := filepath.Join(s.BankDir, "activeContext.md")
	activeFocus := "No active focus declared"
	idleHours := 0.0
	var lineCount int

	if data, err := os.ReadFile(activePath); err == nil {
		lines := strings.Split(string(data), "\n")
		lineCount = len(lines)
		inFocus := false
		for _, l := range lines {
			trimmed := strings.TrimSpace(l)
			if strings.HasPrefix(trimmed, "## Focus") {
				inFocus = true
				continue
			} else if strings.HasPrefix(trimmed, "## ") {
				inFocus = false
			}
			if inFocus && trimmed != "" && !strings.HasPrefix(trimmed, "#") {
				activeFocus = trimmed
				break
			}
		}
	}

	if info, err := os.Stat(activePath); err == nil {
		idleHours = time.Since(info.ModTime()).Hours()
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"project_name":  filepath.Base(s.RepoDir),
		"repo_dir":      s.RepoDir,
		"branch":        branch,
		"dirty_count":   len(dirtyFiles),
		"dirty_files":   dirtyFiles,
		"active_focus":  activeFocus,
		"idle_hours":    idleHours,
		"active_lines":  lineCount,
		"budget_status": getBudgetStatus(lineCount),
	})
}

func (s *Server) handleFiles(w http.ResponseWriter, r *http.Request) {
	manifest, err := core.LoadManifest(s.BankDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to load manifest")
		return
	}

	type FileDetail struct {
		Name         string `json:"name"`
		Volatility   string `json:"volatility"`
		LineCount    int    `json:"line_count"`
		ByteSize     int64  `json:"byte_size"`
		Content      string `json:"content"`
		BudgetStatus string `json:"budget_status"`
	}

	var results []FileDetail
	for name, meta := range manifest.Files {
		filePath := filepath.Join(s.BankDir, name)
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		results = append(results, FileDetail{
			Name:         name,
			Volatility:   string(meta.Volatility),
			LineCount:    len(lines),
			ByteSize:     meta.Bytes,
			Content:      string(data),
			BudgetStatus: getBudgetStatus(len(lines)),
		})
	}

	respondJSON(w, http.StatusOK, results)
}

func (s *Server) handleFileGet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		respondError(w, http.StatusBadRequest, "file name required")
		return
	}
	cleanName := filepath.Base(name)
	filePath := filepath.Join(s.BankDir, cleanName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		respondError(w, http.StatusNotFound, "file not found")
		return
	}
	lines := strings.Split(string(data), "\n")
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"filename": cleanName,
		"content":  string(data),
		"lines":    len(lines),
	})
}

func (s *Server) handleFileSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}

	var req struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Filename == "activeContext.md" {
		lines := strings.Split(req.Content, "\n")
		if len(lines) > 150 {
			respondError(w, http.StatusBadRequest, "activeContext.md strictly exceeds the 150-line budget ceiling")
			return
		}
	}

	targetPath := filepath.Join(s.BankDir, req.Filename)
	if err := core.WriteAtomic(targetPath, []byte(req.Content), 0644); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = core.UpdateFileMeta(s.BankDir, req.Filename, "")
	respondJSON(w, http.StatusOK, map[string]string{"status": "saved", "filename": req.Filename})
}

func (s *Server) handleGraph(w http.ResponseWriter, r *http.Request) {
	type Node struct {
		ID       string `json:"id"`
		Label    string `json:"label"`
		Type     string `json:"type"` // core, context, hot, checkpoint
		Subtitle string `json:"subtitle"`
		Lines    int    `json:"lines"`
		Status   string `json:"status"` // normal, warning, danger
	}

	type Edge struct {
		Source string `json:"source"`
		Target string `json:"target"`
		Label  string `json:"label"`
	}

	nodes := []Node{
		{ID: "projectbrief", Label: "projectbrief.md", Type: "core", Subtitle: "Foundation & Scope"},
		{ID: "productContext", Label: "productContext.md", Type: "context", Subtitle: "UX Goals & Ingested Notes"},
		{ID: "systemPatterns", Label: "systemPatterns.md", Type: "context", Subtitle: "Architecture & Rules"},
		{ID: "techContext", Label: "techContext.md", Type: "context", Subtitle: "Stack & Dependencies"},
		{ID: "activeContext", Label: "activeContext.md", Type: "hot", Subtitle: "Active Focus & Next Steps"},
		{ID: "progress", Label: "progress.md", Type: "hot", Subtitle: "Milestone Ledger"},
		{ID: "decisionLog", Label: "decisionLog.md", Type: "core", Subtitle: "Audit & Architecture ADRs"},
	}

	// Read line counts
	for i, n := range nodes {
		p := filepath.Join(s.BankDir, n.Label)
		if data, err := os.ReadFile(p); err == nil {
			nodes[i].Lines = len(strings.Split(string(data), "\n"))
			if n.ID == "activeContext" && nodes[i].Lines > 120 {
				nodes[i].Status = "warning"
			}
			if n.ID == "activeContext" && nodes[i].Lines >= 150 {
				nodes[i].Status = "danger"
			}
		}
	}

	edges := []Edge{
		{Source: "projectbrief", Target: "productContext", Label: "informs"},
		{Source: "projectbrief", Target: "systemPatterns", Label: "guides"},
		{Source: "systemPatterns", Target: "techContext", Label: "constrains"},
		{Source: "productContext", Target: "activeContext", Label: "prioritizes"},
		{Source: "techContext", Target: "activeContext", Label: "grounds"},
		{Source: "activeContext", Target: "progress", Label: "advances"},
		{Source: "activeContext", Target: "decisionLog", Label: "records"},
	}

	// Add recent checkpoints to the tree
	ckptIDs, _ := checkpoint.ListSnapshots(s.BankDir)
	for i, id := range ckptIDs {
		if i >= 5 {
			break
		}
		nodes = append(nodes, Node{
			ID:       id,
			Label:    id,
			Type:     "checkpoint",
			Subtitle: "snapshot",
		})
		edges = append(edges, Edge{
			Source: "activeContext",
			Target: id,
			Label:  "snapshot",
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"nodes": nodes,
		"edges": edges,
	})
}

func (s *Server) handleCheckpoints(w http.ResponseWriter, r *http.Request) {
	ids, err := checkpoint.ListSnapshots(s.BankDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list checkpoints")
		return
	}

	var fullCheckpoints []types.Checkpoint
	for _, id := range ids {
		ckpt, err := checkpoint.LoadSnapshot(s.BankDir, id)
		if err == nil {
			fullCheckpoints = append(fullCheckpoints, *ckpt)
		}
	}

	respondJSON(w, http.StatusOK, fullCheckpoints)
}

func (s *Server) handleCheckpointCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}

	var req struct {
		Focus string `json:"focus"`
		Notes string `json:"notes"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	var manualEdits []string
	if req.Notes != "" {
		manualEdits = append(manualEdits, req.Notes)
	}

	ckpt, err := checkpoint.CreateSnapshot(s.BankDir, s.RepoDir, req.Focus, nil, manualEdits)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, ckpt)
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	report, err := audit.RunAudit(s.RepoDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, report)
}

func (s *Server) handleLint(w http.ResponseWriter, r *http.Request) {
	fix := r.URL.Query().Get("fix") == "true"
	res, err := linter.LintMemoryBank(s.BankDir, fix)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, res)
}

func (s *Server) handleIngestPreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}

	var req struct {
		Content  string `json:"content"`
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	// Write temp file to inspect
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("inbox_%d.md", time.Now().UnixNano()))
	_ = os.WriteFile(tempFile, []byte(req.Content), 0644)
	defer os.Remove(tempFile)

	prop, err := ingest.PrepareIngestion(s.BankDir, tempFile)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, prop)
}

func (s *Server) handleIngestCommit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}

	var prop ingest.IngestionProposal
	if err := json.NewDecoder(r.Body).Decode(&prop); err != nil {
		respondError(w, http.StatusBadRequest, "invalid proposal payload")
		return
	}

	// Write temporary dummy source if not present
	if _, err := os.Stat(prop.SourceFile); os.IsNotExist(err) {
		tempSrc := filepath.Join(s.RepoDir, "research", "inbox", "web_upload.md")
		_ = os.MkdirAll(filepath.Dir(tempSrc), 0755)
		_ = os.WriteFile(tempSrc, []byte(prop.ProposedDiff), 0644)
		prop.SourceFile = tempSrc
	}

	if err := ingest.CommitIngestion(s.BankDir, s.RepoDir, &prop); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "committed", "target": prop.TargetFile})
}

func getBudgetStatus(lines int) string {
	if lines >= 150 {
		return "danger"
	}
	if lines > 120 {
		return "warning"
	}
	return "normal"
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func openURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
