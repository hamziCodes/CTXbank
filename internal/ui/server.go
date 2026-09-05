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
	"sync"
	"time"

	"github.com/ctxbank/ctx/internal/audit"
	"github.com/ctxbank/ctx/internal/checkpoint"
	"github.com/ctxbank/ctx/internal/core"
	"github.com/ctxbank/ctx/internal/git"
	"github.com/ctxbank/ctx/internal/ingest"
	"github.com/ctxbank/ctx/internal/linter"
	"github.com/ctxbank/ctx/internal/workspace"
	"github.com/ctxbank/ctx/pkg/types"
)

//go:embed web/*
var embeddedFiles embed.FS

// Server hosts the CTXbank local interactive web dashboard.
type Server struct {
	mu      sync.RWMutex
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

func (s *Server) getPaths() (string, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.RepoDir, s.BankDir
}

func (s *Server) setRepoDir(newRepo string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.RepoDir = newRepo
	s.BankDir = filepath.Join(newRepo, core.MemoryBankDir)
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
	mux.HandleFunc("/api/projects", s.handleProjects)
	mux.HandleFunc("/api/project/inspect", s.handleProjectInspect)
	mux.HandleFunc("/api/project/switch", s.handleProjectSwitch)
	mux.HandleFunc("/api/project/init", s.handleProjectInit)
	mux.HandleFunc("/api/sync", s.handleSync)

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
	repoDir, bankDir := s.getPaths()
	branch, _ := git.GetCurrentBranch(repoDir)
	dirtyFiles, _ := git.GetStatus(repoDir)

	manifestPath := filepath.Join(bankDir, core.StateDirname, core.ManifestFilename)
	hasBank := false
	if _, err := os.Stat(manifestPath); err == nil {
		hasBank = true
	}

	activePath := filepath.Join(bankDir, "activeContext.md")
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
		"project_name":  filepath.Base(repoDir),
		"repo_dir":      repoDir,
		"branch":        branch,
		"has_bank":      hasBank,
		"dirty_count":   len(dirtyFiles),
		"dirty_files":   dirtyFiles,
		"active_focus":  activeFocus,
		"idle_hours":    idleHours,
		"active_lines":  lineCount,
		"budget_status": getBudgetStatus(lineCount),
	})
}

func (s *Server) handleFiles(w http.ResponseWriter, r *http.Request) {
	_, bankDir := s.getPaths()
	manifest, err := core.LoadManifest(bankDir)

	type FileDetail struct {
		Name         string `json:"name"`
		Filename     string `json:"filename"`
		Volatility   string `json:"volatility"`
		LineCount    int    `json:"line_count"`
		Lines        int    `json:"lines"`
		ByteSize     int64  `json:"byte_size"`
		Bytes        int64  `json:"bytes"`
		Content      string `json:"content"`
		BudgetStatus string `json:"budget_status"`
	}

	var results []FileDetail
	if err != nil || manifest == nil {
		// If manifest doesn't exist yet, return empty list gracefully
		respondJSON(w, http.StatusOK, results)
		return
	}

	for name, meta := range manifest.Files {
		filePath := filepath.Join(bankDir, name)
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		results = append(results, FileDetail{
			Name:         name,
			Filename:     name,
			Volatility:   string(meta.Volatility),
			LineCount:    len(lines),
			Lines:        len(lines),
			ByteSize:     meta.Bytes,
			Bytes:        meta.Bytes,
			Content:      string(data),
			BudgetStatus: getBudgetStatus(len(lines)),
		})
	}

	respondJSON(w, http.StatusOK, results)
}

func (s *Server) handleFileGet(w http.ResponseWriter, r *http.Request) {
	_, bankDir := s.getPaths()
	name := r.URL.Query().Get("name")
	if name == "" {
		respondError(w, http.StatusBadRequest, "file name required")
		return
	}
	cleanName := filepath.Base(name)
	filePath := filepath.Join(bankDir, cleanName)
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

	_, bankDir := s.getPaths()
	targetPath := filepath.Join(bankDir, req.Filename)
	if err := core.WriteAtomic(targetPath, []byte(req.Content), 0644); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = core.UpdateFileMeta(bankDir, req.Filename, "")
	respondJSON(w, http.StatusOK, map[string]string{"status": "saved", "filename": req.Filename})
}

func (s *Server) handleGraph(w http.ResponseWriter, r *http.Request) {
	_, bankDir := s.getPaths()
	type Node struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
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
		{ID: "projectbrief", Name: "projectbrief.md", Label: "projectbrief.md", Type: "core", Subtitle: "Foundation & Scope"},
		{ID: "productContext", Name: "productContext.md", Label: "productContext.md", Type: "context", Subtitle: "UX Goals & Ingested Notes"},
		{ID: "systemPatterns", Name: "systemPatterns.md", Label: "systemPatterns.md", Type: "context", Subtitle: "Architecture & Rules"},
		{ID: "techContext", Name: "techContext.md", Label: "techContext.md", Type: "context", Subtitle: "Stack & Dependencies"},
		{ID: "activeContext", Name: "activeContext.md", Label: "activeContext.md", Type: "hot", Subtitle: "Active Focus & Next Steps"},
		{ID: "progress", Name: "progress.md", Label: "progress.md", Type: "hot", Subtitle: "Milestone Ledger"},
		{ID: "decisionLog", Name: "decisionLog.md", Label: "decisionLog.md", Type: "core", Subtitle: "Audit & Architecture ADRs"},
	}

	// Read line counts
	for i, n := range nodes {
		p := filepath.Join(bankDir, n.Label)
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
	ckptIDs, _ := checkpoint.ListSnapshots(bankDir)
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
	_, bankDir := s.getPaths()
	ids, err := checkpoint.ListSnapshots(bankDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list checkpoints")
		return
	}

	var fullCheckpoints []types.Checkpoint
	for _, id := range ids {
		ckpt, err := checkpoint.LoadSnapshot(bankDir, id)
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

	repoDir, bankDir := s.getPaths()
	ckpt, err := checkpoint.CreateSnapshot(bankDir, repoDir, req.Focus, nil, manualEdits)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, ckpt)
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	repoDir, _ := s.getPaths()
	report, err := audit.RunAudit(repoDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, report)
}

func (s *Server) handleLint(w http.ResponseWriter, r *http.Request) {
	_, bankDir := s.getPaths()
	fix := r.URL.Query().Get("fix") == "true"
	res, err := linter.LintMemoryBank(bankDir, fix)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, res)
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	repoDir, bankDir := s.getPaths()
	manifestPath := filepath.Join(bankDir, core.StateDirname, core.ManifestFilename)
	if _, err := os.Stat(manifestPath); err != nil {
		if err := core.InitBank(repoDir, false); err != nil {
			respondError(w, http.StatusInternalServerError, "failed to initialize memory bank: "+err.Error())
			return
		}
	}

	report, err := audit.RunAudit(repoDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to run reconnaissance audit: "+err.Error())
		return
	}

	if err := audit.ApplyAudit(bankDir, report); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to apply audit: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"message":  "Memory bank successfully synchronized with current codebase",
		"packages": len(report.Components),
		"symbols":  len(report.Symbols),
	})
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

	_, bankDir := s.getPaths()
	prop, err := ingest.PrepareIngestion(bankDir, tempFile)
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

	repoDir, bankDir := s.getPaths()
	// Write temporary dummy source if not present
	if _, err := os.Stat(prop.SourceFile); os.IsNotExist(err) {
		tempSrc := filepath.Join(repoDir, "research", "inbox", "web_upload.md")
		_ = os.MkdirAll(filepath.Dir(tempSrc), 0755)
		_ = os.WriteFile(tempSrc, []byte(prop.ProposedDiff), 0644)
		prop.SourceFile = tempSrc
	}

	if err := ingest.CommitIngestion(bankDir, repoDir, &prop); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "committed", "target": prop.TargetFile})
}

func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
	repoDir, _ := s.getPaths()
	parentDir := filepath.Dir(repoDir)
	discovered, _ := workspace.ScanWorkspaces(parentDir, 2)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"current_project": filepath.Base(repoDir),
		"current_path":    repoDir,
		"discovered":      discovered,
	})
}

func (s *Server) handleProjectInspect(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Query().Get("path")
	if targetPath == "" {
		respondError(w, http.StatusBadRequest, "path parameter required")
		return
	}
	cleanPath := filepath.Clean(targetPath)
	fi, err := os.Stat(cleanPath)
	if err != nil || !fi.IsDir() {
		respondError(w, http.StatusNotFound, "directory does not exist")
		return
	}

	manifestPath := filepath.Join(cleanPath, core.MemoryBankDir, core.StateDirname, core.ManifestFilename)
	hasBank := false
	if _, err := os.Stat(manifestPath); err == nil {
		hasBank = true
	}

	branch, _ := git.GetCurrentBranch(cleanPath)
	dirty, _ := git.GetStatus(cleanPath)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"name":        filepath.Base(cleanPath),
		"path":        cleanPath,
		"has_bank":    hasBank,
		"branch":      branch,
		"dirty_count": len(dirty),
	})
}

func (s *Server) handleProjectSwitch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}

	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Path == "" {
		respondError(w, http.StatusBadRequest, "invalid path")
		return
	}
	cleanPath := filepath.Clean(req.Path)
	fi, err := os.Stat(cleanPath)
	if err != nil || !fi.IsDir() {
		respondError(w, http.StatusNotFound, "directory does not exist")
		return
	}

	manifestPath := filepath.Join(cleanPath, core.MemoryBankDir, core.StateDirname, core.ManifestFilename)
	if _, err := os.Stat(manifestPath); err != nil {
		respondError(w, http.StatusBadRequest, "project has not been initialized with CTXbank yet")
		return
	}

	s.setRepoDir(cleanPath)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"name":    filepath.Base(cleanPath),
		"path":    cleanPath,
	})
}

func (s *Server) handleProjectInit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}

	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Path == "" {
		respondError(w, http.StatusBadRequest, "invalid path")
		return
	}
	cleanPath := filepath.Clean(req.Path)
	fi, err := os.Stat(cleanPath)
	if err != nil || !fi.IsDir() {
		respondError(w, http.StatusNotFound, "directory does not exist")
		return
	}

	// 1. Initialize Memory Bank
	if err := core.InitBank(cleanPath, false); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to initialize memory bank: "+err.Error())
		return
	}

	// 2. Run Reconnaissance Audit & apply
	bankDir := filepath.Join(cleanPath, core.MemoryBankDir)
	report, err := audit.RunAudit(cleanPath)
	if err == nil && report != nil {
		_ = audit.ApplyAudit(bankDir, report)
	}

	// 3. Switch active workspace to this project
	s.setRepoDir(cleanPath)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"name":    filepath.Base(cleanPath),
		"path":    cleanPath,
	})
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
