package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ctxbank/ctx/internal/core"
)

func TestPromptSyncEngine(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctxbank-sync-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	bankDir := filepath.Join(tempDir, core.MemoryBankDir)

	// 1. Test Initial Ledger
	ledger, err := GetLedger(bankDir)
	if err != nil {
		t.Fatalf("unexpected error getting initial ledger: %v", err)
	}
	if ledger.TotalSyncs != 0 {
		t.Errorf("expected 0 total syncs, got %d", ledger.TotalSyncs)
	}
	if ledger.Status != "idle" {
		t.Errorf("expected 'idle' status, got %q", ledger.Status)
	}

	// 2. Test GenerateSyncPrompt
	res, err := GenerateSyncPrompt(tempDir, bankDir)
	if err != nil {
		t.Fatalf("GenerateSyncPrompt failed: %v", err)
	}
	if res.TotalSyncs != 1 {
		t.Errorf("expected TotalSyncs to be 1, got %d", res.TotalSyncs)
	}
	if res.Status != "awaiting_agent" {
		t.Errorf("expected status 'awaiting_agent', got %q", res.Status)
	}
	if !strings.Contains(res.PromptText, "AI DIRECTIVE: Populate CTXbank Memory Bank") {
		t.Errorf("prompt text missing expected header")
	}

	// Verify prompt file was written to disk
	promptData, err := os.ReadFile(res.PromptFile)
	if err != nil {
		t.Fatalf("failed to read written prompt file %s: %v", res.PromptFile, err)
	}
	if len(promptData) == 0 {
		t.Errorf("prompt file is empty")
	}

	// 3. Test VerifySyncCompletion - should detect template/placeholders in fresh bank
	vRes, err := VerifySyncCompletion(bankDir)
	if err != nil {
		t.Fatalf("VerifySyncCompletion failed: %v", err)
	}
	// Initial bank contains generic template phrases, so issues are expected
	if vRes.Verified {
		t.Errorf("expected initial bank to have issues due to generic template phrases, but got verified=true")
	}
	if len(vRes.Issues) == 0 {
		t.Errorf("expected verification issues, got none")
	}

	// 4. Overwrite files with realistic custom documentation
	files := []string{
		"projectbrief.md",
		"productContext.md",
		"systemPatterns.md",
		"techContext.md",
		"activeContext.md",
		"progress.md",
		"decisionLog.md",
	}
	for _, f := range files {
		var sampleContent string
		if f == "activeContext.md" {
			sampleContent = `# Active Context
## Focus
Developing the real checkout system.
## Recent
- Setup checkout API.
## Next steps
- Integrate payment gateway.
## Open decisions
- None.
`
		} else {
			sampleContent = `# ` + f + `
## Real System Architecture
This project is an advanced e-commerce microservices platform.
- Service A handles billing.
- Service B handles inventory.
- Real production configuration verified.
- Database: PostgreSQL 16.
`
		}
		if err := core.WriteAtomic(filepath.Join(bankDir, f), []byte(sampleContent), 0644); err != nil {
			t.Fatalf("failed to write test file %s: %v", f, err)
		}
	}

	// 5. Test VerifySyncCompletion with cleaned files
	vRes2, err := VerifySyncCompletion(bankDir)
	if err != nil {
		t.Fatalf("VerifySyncCompletion failed: %v", err)
	}
	if !vRes2.Verified {
		t.Errorf("expected verified=true after custom content, got issues: %v", vRes2.Issues)
	}
	if vRes2.Status != "verified" {
		t.Errorf("expected status 'verified', got %q", vRes2.Status)
	}

	// 6. Check updated ledger
	updatedLedger, err := GetLedger(bankDir)
	if err != nil {
		t.Fatalf("failed to reload ledger: %v", err)
	}
	if updatedLedger.Status != "verified" {
		t.Errorf("expected ledger status 'verified', got %q", updatedLedger.Status)
	}
	if updatedLedger.TotalSyncs != 1 {
		t.Errorf("expected TotalSyncs to be 1, got %d", updatedLedger.TotalSyncs)
	}
}
