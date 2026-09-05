package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanManifestsGo(t *testing.T) {
	tempDir := t.TempDir()
	goModContent := `module example.com/testmod

go 1.22

require (
	github.com/stretchr/testify v1.8.4
	golang.org/x/sys v0.15.0
)
`
	if err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	manifests := ScanManifests(tempDir)
	if len(manifests) == 0 {
		t.Fatalf("expected at least 1 manifest found")
	}

	m := manifests[0]
	if m.Language != "Go" {
		t.Errorf("expected language Go, got %s", m.Language)
	}
	if len(m.Dependencies) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(m.Dependencies))
	}
}

func TestScanStructuralHeuristics(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(tempDir, "cmd", "app"), 0755)
	_ = os.MkdirAll(filepath.Join(tempDir, "internal", "core"), 0755)
	_ = os.MkdirAll(filepath.Join(tempDir, "internal", "mcp"), 0755)

	components := ScanStructuralHeuristics(tempDir)
	if len(components) < 3 {
		t.Errorf("expected at least 3 architectural components, got %d", len(components))
	}
}

func TestScanSymbols(t *testing.T) {
	tempDir := t.TempDir()
	sampleGo := `package main

import "fmt"

func ProcessData() {
	fmt.Println("test")
}

func main() {
	ProcessData()
}
`
	if err := os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(sampleGo), 0644); err != nil {
		t.Fatalf("failed to write main.go: %v", err)
	}

	symbols := ScanSymbols(tempDir)
	foundMain := false
	foundProcess := false

	for _, s := range symbols {
		if s.Name == "main" && s.Kind == "entrypoint" {
			foundMain = true
		}
		if s.Name == "ProcessData" && s.Kind == "exported_func" {
			foundProcess = true
		}
	}

	if !foundMain {
		t.Errorf("failed to detect main entrypoint")
	}
	if !foundProcess {
		t.Errorf("failed to detect exported function ProcessData")
	}
}
