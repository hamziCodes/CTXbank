package ui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/ctxbank/ctx/internal/core"
)

func TestUIServerAPIRoutes(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	server := NewServer(tempWorkspace, 0)

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
}

