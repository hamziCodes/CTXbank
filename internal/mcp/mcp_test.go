package mcp

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ctxbank/ctx/internal/core"
)

func TestMCPServerHandshakeAndTools(t *testing.T) {
	tempWorkspace := t.TempDir()

	// Initialize memory bank
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)
	server := NewServer(bankDir, tempWorkspace)

	// 1. Test initialize
	initReq := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}` + "\n"
	var out bytes.Buffer
	if err := server.Serve(strings.NewReader(initReq), &out); err != nil {
		t.Fatalf("Serve failed: %v", err)
	}

	var initResp JSONRPCResponse
	if err := json.Unmarshal(out.Bytes(), &initResp); err != nil {
		t.Fatalf("corrupt response JSON: %v", err)
	}
	if initResp.Error != nil {
		t.Fatalf("initialize returned error: %+v", initResp.Error)
	}

	// 2. Test tools/list
	out.Reset()
	listReq := `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}` + "\n"
	_ = server.Serve(strings.NewReader(listReq), &out)

	var listResp struct {
		Result struct {
			Tools []struct {
				Name string `json:"name"`
			} `json:"tools"`
		} `json:"result"`
	}
	if err := json.Unmarshal(out.Bytes(), &listResp); err != nil {
		t.Fatalf("failed to parse tools/list: %v", err)
	}
	if len(listResp.Result.Tools) < 6 {
		t.Errorf("expected 6 tools, got %d", len(listResp.Result.Tools))
	}
}

func TestMCPDeltaCaching(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)
	server := NewServer(bankDir, tempWorkspace)
	sessionID := "test-agent-session-42"

	// First call to read_active_context
	callReq1 := `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"read_active_context","arguments":{"session_id":"` + sessionID + `"}}}` + "\n"
	var out1 bytes.Buffer
	_ = server.Serve(strings.NewReader(callReq1), &out1)

	var resp1 struct {
		Result struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
	}
	_ = json.Unmarshal(out1.Bytes(), &resp1)
	if len(resp1.Result.Content) == 0 {
		t.Fatalf("expected content in first call: %s", out1.String())
	}

	var payload1 map[string]interface{}
	_ = json.Unmarshal([]byte(resp1.Result.Content[0].Text), &payload1)
	if payload1["unchanged"] == true {
		t.Errorf("first call should not be unchanged")
	}

	// Second immediate call with same session ID must return unchanged: true (Zero-token savings!)
	callReq2 := `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"read_active_context","arguments":{"session_id":"` + sessionID + `"}}}` + "\n"
	var out2 bytes.Buffer
	_ = server.Serve(strings.NewReader(callReq2), &out2)

	var resp2 struct {
		Result struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
	}
	_ = json.Unmarshal(out2.Bytes(), &resp2)

	var payload2 map[string]interface{}
	_ = json.Unmarshal([]byte(resp2.Result.Content[0].Text), &payload2)
	if payload2["unchanged"] != true {
		t.Errorf("second call with identical hash MUST return unchanged=true, got: %+v", payload2)
	}

	// Append delta to activeContext.md
	deltaReq := `{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"append_memory_delta","arguments":{"session_id":"` + sessionID + `","file":"activeContext.md","section":"Focus","content":"- Newly Added Line\n"}}}` + "\n"
	var out3 bytes.Buffer
	_ = server.Serve(strings.NewReader(deltaReq), &out3)

	// Third call to read_active_context must detect the change!
	callReq3 := `{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"read_active_context","arguments":{"session_id":"` + sessionID + `"}}}` + "\n"
	var out4 bytes.Buffer
	_ = server.Serve(strings.NewReader(callReq3), &out4)

	var resp3 struct {
		Result struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
	}
	_ = json.Unmarshal(out4.Bytes(), &resp3)
	var payload3 map[string]interface{}
	_ = json.Unmarshal([]byte(resp3.Result.Content[0].Text), &payload3)

	if payload3["unchanged"] == true {
		t.Errorf("call after delta mutation should NOT be unchanged")
	}
}

func TestMCPLineBudgetEnforcement(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)
	server := NewServer(bankDir, tempWorkspace)

	// Try to append 200 lines to activeContext.md
	giantChunk := strings.Repeat("- excessive line\n", 200)
	payload, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      10,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "append_memory_delta",
			"arguments": map[string]interface{}{
				"session_id": "budget-test",
				"file":       "activeContext.md",
				"section":    "Focus",
				"content":    giantChunk,
			},
		},
	})

	var out bytes.Buffer
	_ = server.Serve(strings.NewReader(string(payload)+"\n"), &out)

	var resp struct {
		Result struct {
			IsError bool `json:"isError"`
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
	}
	_ = json.Unmarshal(out.Bytes(), &resp)

	if !resp.Result.IsError {
		t.Errorf("expected error when line budget > 150 lines is exceeded")
	}
}

func TestMCPResourcesAndPrompts(t *testing.T) {
	tempWorkspace := t.TempDir()
	if err := core.InitBank(tempWorkspace, false); err != nil {
		t.Fatalf("InitBank failed: %v", err)
	}

	bankDir := filepath.Join(tempWorkspace, core.MemoryBankDir)
	server := NewServer(bankDir, tempWorkspace)

	// 1. Read resource memory://activeContext
	resReq := `{"jsonrpc":"2.0","id":20,"method":"resources/read","params":{"uri":"memory://activeContext"}}` + "\n"
	var resOut bytes.Buffer
	_ = server.Serve(strings.NewReader(resReq), &resOut)

	var resResp struct {
		Result struct {
			Contents []struct {
				URI  string `json:"uri"`
				Text string `json:"text"`
			} `json:"contents"`
		} `json:"result"`
	}
	_ = json.Unmarshal(resOut.Bytes(), &resResp)
	if len(resResp.Result.Contents) == 0 || !strings.Contains(resResp.Result.Contents[0].Text, "Active Context") {
		t.Errorf("failed to read memory://activeContext: %s", resOut.String())
	}

	// 2. Get prompt resume-session
	promptReq := `{"jsonrpc":"2.0","id":21,"method":"prompts/get","params":{"name":"resume-session"}}` + "\n"
	var promptOut bytes.Buffer
	_ = server.Serve(strings.NewReader(promptReq), &promptOut)

	var promptResp struct {
		Result struct {
			Messages []struct {
				Content struct {
					Text string `json:"text"`
				} `json:"content"`
			} `json:"messages"`
		} `json:"result"`
	}
	_ = json.Unmarshal(promptOut.Bytes(), &promptResp)
	if len(promptResp.Result.Messages) == 0 || !strings.Contains(promptResp.Result.Messages[0].Content.Text, "Agent Pickup Brief") {
		t.Errorf("failed to get resume-session prompt: %s", promptOut.String())
	}
}
