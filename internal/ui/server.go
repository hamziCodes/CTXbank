package ui

import (
	"crypto/subtle"
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

	"github.com/hamziCodes/CTXbank/internal/audit"
	"github.com/hamziCodes/CTXbank/internal/checkpoint"
	"github.com/hamziCodes/CTXbank/internal/core"
	"github.com/hamziCodes/CTXbank/internal/git"
	"github.com/hamziCodes/CTXbank/internal/ingest"
	"github.com/hamziCodes/CTXbank/internal/linter"
	syncpkg "github.com/hamziCodes/CTXbank/internal/sync"
	"github.com/hamziCodes/CTXbank/pkg/types"
)

//go:embed web/*
var embeddedFiles embed.FS

// Version is the ctx binary version, set by cmd/ctx at startup
// (release builds stamp it via ldflags). Surfaced in /api/status so
// the hosted Connect page can tell stale dashboards apart.
var Version = "dev"

// writableMemoryFiles is the whitelist for /api/file/save: only the 7
// Cline-compatible memory-bank files may be written through the dashboard.
var writableMemoryFiles = map[string]bool{
	"projectbrief.md":  true,
	"productContext.md": true,
	"systemPatterns.md": true,
	"techContext.md":    true,
	"activeContext.md":  true,
	"progress.md":       true,
	"decisionLog.md":    true,
}

// Server hosts the CTXbank local interactive web dashboard.
// One server = one project. There is no project switching: the dashboard
// serves the repository it was launched in, gated by that project's
// persistent dashboard token (`ctx token`).
type Server struct {
	mu      sync.RWMutex
	RepoDir string
	BankDir string
	Port    int
	// AuthToken is this project's persistent dashboard token, required by
	// every /api/* endpoint (query ?token= or X-CTX-Token header). It
	// prevents other local processes — and drive-by websites hitting
	// 127.0.0.1 — from reading or mutating the memory bank.
	AuthToken string
}

// NewServer initializes a dashboard server for the target repository,
// loading (or creating) the project's persistent dashboard token.
func NewServer(repoDir string, port int) (*Server, error) {
	if port <= 0 {
		port = 4242
	}
	bankDir := filepath.Join(repoDir, core.MemoryBankDir)
	token, err := core.GetProjectToken(bankDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load project token: %w", err)
	}
	return &Server{
		RepoDir:   repoDir,
		BankDir:   bankDir,
		Port:      port,
		AuthToken: token,
	}, nil
}

// withCORS allows the hosted site (e.g. ctxbank.vertexdevstudio.tech) to
// probe the local dashboard from its Connect page. Auth still rests on the
// project token — CORS alone grants nothing, since every /api/* call
// without a valid token answers 401.
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "X-CTX-Token, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// requireAuth is middleware enforcing the per-launch dashboard token on
// every /api/* route. Static UI assets stay public; they are useless
// without API access.
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			token = r.Header.Get("X-CTX-Token")
		}
		if subtle.ConstantTimeCompare([]byte(token), []byte(s.AuthToken)) != 1 {
			respondError(w, http.StatusUnauthorized, "missing or invalid dashboard token")
			return
		}
		next(w, r)
	}
}

func (s *Server) getPaths() (string, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.RepoDir, s.BankDir
}

// Start binds to localhost, serves the API and embedded UI, and opens the default browser.
func (s *Server) Start(openBrowser bool) error {
	mux := http.NewServeMux()

	// API Routes — all guarded by the project token, with CORS so the
	// hosted Connect page can probe this local dashboard.
	api := func(h http.HandlerFunc) http.HandlerFunc {
		return withCORS(s.requireAuth(h))
	}
	mux.HandleFunc("/api/status", api(s.handleStatus))
	mux.HandleFunc("/api/files", api(s.handleFiles))
	mux.HandleFunc("/api/file", api(s.handleFileGet))
	mux.HandleFunc("/api/file/save", api(s.handleFileSave))
	mux.HandleFunc("/api/graph", api(s.handleGraph))
	mux.HandleFunc("/api/checkpoints", api(s.handleCheckpoints))
	mux.HandleFunc("/api/checkpoint/create", api(s.handleCheckpointCreate))
	mux.HandleFunc("/api/audit", api(s.handleAudit))
	mux.HandleFunc("/api/lint", api(s.handleLint))
	mux.HandleFunc("/api/ingest/preview", api(s.handleIngestPreview))
	mux.HandleFunc("/api/ingest/commit", api(s.handleIngestCommit))
	mux.HandleFunc("/api/sync", api(s.handleSync))
	mux.HandleFunc("/api/prompt-sync", api(s.handlePromptSync))
	mux.HandleFunc("/api/prompt-sync/status", api(s.handlePromptSyncStatus))
	mux.HandleFunc("/api/prompt-sync/verify", api(s.handlePromptSyncVerify))

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

	url := fmt.Sprintf("http://localhost:%d?token=%s", s.Port, s.AuthToken)
	fmt.Printf("\n[+] CTXbank dashboard live at %s\n", url)
	fmt.Println("    Project token is required for all API access (embedded in the URL above).")
	fmt.Println("    Keep it private — it grants read/write access to this project's memory bank.")
	fmt.Println("    Or connect via the hosted site: https://ctxbank.vertexdevstudio.tech/#connect")
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
		"ctx_version":   Version,
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

	// Only the 7 known memory-bank files are writable through the dashboard.
	// The name is sanitized with filepath.Base AND checked against the
	// whitelist so the manifest can never drift from what was written.
	cleanName := filepath.Base(req.Filename)
	if !writableMemoryFiles[cleanName] {
		respondError(w, http.StatusBadRequest, "only memory-bank files may be saved through the dashboard")
		return
	}

	if cleanName == "activeContext.md" {
		lines := strings.Split(req.Content, "\n")
		if len(lines) > 150 {
			respondError(w, http.StatusBadRequest, "activeContext.md strictly exceeds the 150-line budget ceiling")
			return
		}
	}

	_, bankDir := s.getPaths()
	targetPath := filepath.Join(bankDir, cleanName)
	if err := core.WriteAtomic(targetPath, []byte(req.Content), 0644); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = core.UpdateFileMeta(bankDir, cleanName, "")
	respondJSON(w, http.StatusOK, map[string]string{"status": "saved", "filename": cleanName})
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

func (s *Server) handlePromptSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	repoDir, bankDir := s.getPaths()
	res, err := syncpkg.GenerateSyncPrompt(repoDir, bankDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, res)
}

func (s *Server) handlePromptSyncStatus(w http.ResponseWriter, r *http.Request) {
	_, bankDir := s.getPaths()
	ledger, err := syncpkg.GetLedger(bankDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, ledger)
}

func (s *Server) handlePromptSyncVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	_, bankDir := s.getPaths()
	res, err := syncpkg.VerifySyncCompletion(bankDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, res)
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
