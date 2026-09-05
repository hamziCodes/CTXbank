package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ctxbank/ctx/internal/checkpoint"
	"github.com/ctxbank/ctx/internal/core"
	"github.com/ctxbank/ctx/internal/rules"
)

func TestCLIInitAndStatus(t *testing.T) {
	tempWorkspace := t.TempDir()

	// Test Init
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	stubs, err := rules.GenerateStubs(tempWorkspace)
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}
	if len(stubs) == 0 {
		t.Errorf("expected at least one stub generated")
	}

	// Verify activeContext exists
	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)
	activeContextPath := filepath.Join(bankDir, "activeContext.md")
	content, err := os.ReadFile(activeContextPath)
	if err != nil {
		t.Fatalf("failed to read activeContext: %v", err)
	}

	if !strings.Contains(string(content), "## Focus") {
		t.Errorf("activeContext missing ## Focus section")
	}

	// Test Checkpoint snapshot creation
	ckpt, err := checkpoint.CreateSnapshot(bankDir, tempWorkspace, "Testing CLI", []string{"Next step 1"}, []string{"manual note"})
	if err != nil {
		t.Fatalf("CreateSnapshot failed: %v", err)
	}
	if ckpt.ID == "" {
		t.Errorf("expected non-empty checkpoint ID")
	}

	// Verify snapshot JSON is valid
	ckptPath := filepath.Join(bankDir, ".state", "checkpoints", ckpt.ID+".json")
	ckptData, err := os.ReadFile(ckptPath)
	if err != nil {
		t.Fatalf("failed to read snapshot file: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(ckptData, &parsed); err != nil {
		t.Fatalf("corrupt checkpoint JSON: %v", err)
	}
	if parsed["id"] != ckpt.ID {
		t.Errorf("mismatched parsed checkpoint ID: %v", parsed["id"])
	}
}
