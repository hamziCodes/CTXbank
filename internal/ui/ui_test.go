package ui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hamziCodes/CTXbank/internal/core"
)

func TestUIServerAPIRoutes(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	server, srvErr := NewServer(tempWorkspace, 0)
	if srvErr != nil {
		t.Fatalf("NewServer failed: %v", srvErr)
	}

	// Test /api/status handler
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	w := httptest.NewRecorder()
	server.handleStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK from /api/status, got %d", w.Code)
	}

	var statusResp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&statusResp); err != nil {
		t.Fatalf("failed to decode status response: %v", err)
	}

	if statusResp["project_name"] != filepath.Base(tempWorkspace) {
		t.Errorf("expected project_name %s, got %v", filepath.Base(tempWorkspace), statusResp["project_name"])
	}

	// Test /api/graph handler
	reqGraph := httptest.NewRequest(http.MethodGet, "/api/graph", nil)
	wGraph := httptest.NewRecorder()
	server.handleGraph(wGraph, reqGraph)

	if wGraph.Code != http.StatusOK {
		t.Errorf("expected 200 OK from /api/graph, got %d", wGraph.Code)
	}

	var graphResp map[string]interface{}
	if err := json.NewDecoder(wGraph.Body).Decode(&graphResp); err != nil {
		t.Fatalf("failed to decode graph response: %v", err)
	}

	nodes, ok := graphResp["nodes"].([]interface{})
	if !ok || len(nodes) < 7 {
		t.Errorf("expected at least 7 core nodes, got: %v", graphResp)
	}

	// Test /api/file handler
	reqFile := httptest.NewRequest(http.MethodGet, "/api/file?name=projectbrief.md", nil)
	wFile := httptest.NewRecorder()
	server.handleFileGet(wFile, reqFile)

	if wFile.Code != http.StatusOK {
		t.Errorf("expected 200 OK from /api/file, got %d", wFile.Code)
	}

	var fileResp map[string]interface{}
	if err := json.NewDecoder(wFile.Body).Decode(&fileResp); err != nil {
		t.Fatalf("failed to decode file response: %v", err)
	}
	if fileResp["filename"] != "projectbrief.md" {
		t.Errorf("expected filename projectbrief.md, got %v", fileResp["filename"])
	}

	// Test /api/prompt-sync handler
	reqPromptSync := httptest.NewRequest(http.MethodPost, "/api/prompt-sync", nil)
	wPromptSync := httptest.NewRecorder()
	server.handlePromptSync(wPromptSync, reqPromptSync)
	if wPromptSync.Code != http.StatusOK {
		t.Errorf("expected 200 OK from /api/prompt-sync, got %d (body: %s)", wPromptSync.Code, wPromptSync.Body.String())
	}

	// Test /api/prompt-sync/status handler
	reqSyncStatus := httptest.NewRequest(http.MethodGet, "/api/prompt-sync/status", nil)
	wSyncStatus := httptest.NewRecorder()
	server.handlePromptSyncStatus(wSyncStatus, reqSyncStatus)
	if wSyncStatus.Code != http.StatusOK {
		t.Errorf("expected 200 OK from /api/prompt-sync/status, got %d", wSyncStatus.Code)
	}

	// Test /api/prompt-sync/verify handler
	reqSyncVerify := httptest.NewRequest(http.MethodPost, "/api/prompt-sync/verify", nil)
	wSyncVerify := httptest.NewRecorder()
	server.handlePromptSyncVerify(wSyncVerify, reqSyncVerify)
	if wSyncVerify.Code != http.StatusOK {
		t.Errorf("expected 200 OK from /api/prompt-sync/verify, got %d", wSyncVerify.Code)
	}
}


func TestRequireAuthMiddleware(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}
	server, srvErr := NewServer(tempWorkspace, 0)
	if srvErr != nil {
		t.Fatalf("NewServer failed: %v", srvErr)
	}
	if server.AuthToken == "" {
		t.Fatal("expected a generated auth token")
	}
	guarded := server.requireAuth(server.handleStatus)

	// No token -> 401
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	w := httptest.NewRecorder()
	guarded(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", w.Code)
	}

	// Wrong token -> 401
	reqBad := httptest.NewRequest(http.MethodGet, "/api/status?token=wrong", nil)
	wBad := httptest.NewRecorder()
	guarded(wBad, reqBad)
	if wBad.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 with wrong token, got %d", wBad.Code)
	}

	// Correct query token -> 200
	reqOK := httptest.NewRequest(http.MethodGet, "/api/status?token="+server.AuthToken, nil)
	wOK := httptest.NewRecorder()
	guarded(wOK, reqOK)
	if wOK.Code != http.StatusOK {
		t.Errorf("expected 200 with query token, got %d", wOK.Code)
	}

	// Correct header token -> 200
	reqHdr := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	reqHdr.Header.Set("X-CTX-Token", server.AuthToken)
	wHdr := httptest.NewRecorder()
	guarded(wHdr, reqHdr)
	if wHdr.Code != http.StatusOK {
		t.Errorf("expected 200 with header token, got %d", wHdr.Code)
	}
}

func TestFileSaveWhitelist(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}
	server, srvErr := NewServer(tempWorkspace, 0)
	if srvErr != nil {
		t.Fatalf("NewServer failed: %v", srvErr)
	}

	postSave := func(filename, content string) *httptest.ResponseRecorder {
		body := `{"filename":` + quoteJSON(filename) + `,"content":` + quoteJSON(content) + `}`
		req := httptest.NewRequest(http.MethodPost, "/api/file/save?token="+server.AuthToken, strings.NewReader(body))
		w := httptest.NewRecorder()
		server.requireAuth(server.handleFileSave)(w, req)
		return w
	}

	// Arbitrary filename rejected
	if w := postSave("evil.json", "x"); w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for evil.json, got %d", w.Code)
	}
	// Path traversal rejected (Base -> "passwd", not whitelisted)
	if w := postSave("../../etc/passwd", "x"); w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for traversal, got %d", w.Code)
	}
	// Over-budget activeContext rejected
	big := strings.Repeat("line\n", 151)
	if w := postSave("activeContext.md", big); w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for 151-line activeContext, got %d", w.Code)
	}
	// Valid save works and updates the manifest under the sanitized name
	if w := postSave("progress.md", "# Progress\n- done\n"); w.Code != http.StatusOK {
		t.Errorf("expected 200 for progress.md, got %d (%s)", w.Code, w.Body.String())
	}
	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)
	manifest, err := core.LoadManifest(bankDir)
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}
	if _, ok := manifest.Files["progress.md"]; !ok {
		t.Errorf("manifest missing progress.md after save: %+v", manifest.Files)
	}
	if _, ok := manifest.Files["./progress.md"]; ok {
		t.Errorf("manifest recorded unsanitized filename")
	}
}

func quoteJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
