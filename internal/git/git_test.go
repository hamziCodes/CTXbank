package git

import (
	"os"
	"path/filepath"
	"testing"
)

type mockRunner struct {
	responses map[string]string
}

func (m *mockRunner) Run(dir string, args ...string) (string, error) {
	cmd := args[0]
	if val, ok := m.responses[cmd]; ok {
		return val, nil
	}
	return "", nil
}

func TestGetStatusParsing(t *testing.T) {
	mock := &mockRunner{
		responses: map[string]string{
			"status": " M internal/git/porcelain.go\n?? newfile.txt\nA  staged.txt\n",
		},
	}
	SetRunner(mock)
	defer ResetRunner()

	statuses, err := GetStatus("/fake/repo")
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if len(statuses) != 3 {
		t.Fatalf("expected 3 items, got %d", len(statuses))
	}

	if statuses[0].Path != "internal/git/porcelain.go" || statuses[0].WorktreeStatus != "M" {
		t.Errorf("unexpected status[0]: %+v", statuses[0])
	}
	if statuses[1].Path != "newfile.txt" || statuses[1].StagingStatus != "?" {
		t.Errorf("unexpected status[1]: %+v", statuses[1])
	}
	if statuses[2].Path != "staged.txt" || statuses[2].StagingStatus != "A" {
		t.Errorf("unexpected status[2]: %+v", statuses[2])
	}
}

func TestSafeguardsRebase(t *testing.T) {
	tempGitDir := t.TempDir()

	if IsMidRebase(tempGitDir) {
		t.Errorf("empty dir should not be mid-rebase")
	}

	// Create rebase-merge
	rebaseDir := filepath.Join(tempGitDir, "rebase-merge")
	if err := os.Mkdir(rebaseDir, 0755); err != nil {
		t.Fatalf("failed to create rebase dir: %v", err)
	}

	if !IsMidRebase(tempGitDir) {
		t.Errorf("expected IsMidRebase to be true when rebase-merge exists")
	}
}
