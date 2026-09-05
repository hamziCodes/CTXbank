package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// JSONRPCRequest models an incoming JSON-RPC 2.0 message.
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse models an outgoing JSON-RPC 2.0 response.
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError models a standard JSON-RPC 2.0 error.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Server implements the Model Context Protocol over stdio.
type Server struct {
	bankDir string
	repoDir string
	mu      sync.Mutex
}

// NewServer initializes an MCP Server bound to target bank and repository directories.
func NewServer(bankDir, repoDir string) *Server {
	return &Server{
		bankDir: bankDir,
		repoDir: repoDir,
	}
}

// Serve reads JSON-RPC line-delimited requests from in and writes responses to out.
func (s *Server) Serve(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	// Support payloads up to 10MB
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	encoder := json.NewEncoder(out)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			s.sendError(encoder, nil, -32700, "Parse error")
			continue
		}

		// Notifications (ID is nil) do not send responses unless required
		if req.ID == nil && req.Method == "notifications/initialized" {
			continue
		}

		resp := s.handleRequest(req)
		s.mu.Lock()
		_ = encoder.Encode(resp)
		s.mu.Unlock()
	}

	return scanner.Err()
}

func (s *Server) sendError(enc *json.Encoder, id interface{}, code int, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = enc.Encode(JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &RPCError{Code: code, Message: msg},
	})
}

func (s *Server) handleRequest(req JSONRPCRequest) JSONRPCResponse {
	switch req.Method {
	case "initialize":
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"serverInfo": map[string]string{
					"name":    "memory-bank-mcp",
					"version": "0.1.0",
				},
				"capabilities": map[string]interface{}{
					"tools":     map[string]interface{}{},
					"resources": map[string]interface{}{},
					"prompts":   map[string]interface{}{},
				},
			},
		}

	case "ping":
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  map[string]interface{}{},
		}

	case "tools/list":
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": s.getToolsList(),
			},
		}

	case "tools/call":
		return s.handleToolCall(req)

	case "resources/list":
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"resources": ListResources(),
			},
		}

	case "resources/read":
		return s.handleResourceRead(req)

	case "prompts/list":
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"prompts": ListPrompts(),
			},
		}

	case "prompts/get":
		return s.handlePromptGet(req)

	default:
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: -32601, Message: fmt.Sprintf("Method not found: %s", req.Method)},
		}
	}
}

func (s *Server) getToolsList() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "read_active_context",
			"description": "Reads activeContext.md and progress.md. Returns {'unchanged': true} if unchanged since session cursor.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"session_id": map[string]string{"type": "string", "description": "Unique agent session identifier"},
				},
			},
		},
		{
			"name":        "read_static_context",
			"description": "Reads static memory bank files individually hash-gated.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"session_id": map[string]string{"type": "string"},
					"files": map[string]interface{}{
						"type":  "array",
						"items": map[string]string{"type": "string"},
					},
				},
			},
		},
		{
			"name":        "append_memory_delta",
			"description": "Appends delta under named markdown section; strictly enforces 150-line budget.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"session_id": map[string]string{"type": "string"},
					"file":       map[string]string{"type": "string"},
					"section":    map[string]string{"type": "string"},
					"content":    map[string]string{"type": "string"},
				},
				"required": []string{"file", "section", "content"},
			},
		},
		{
			"name":        "update_milestone",
			"description": "Performs structured update to progress.md checklist and activeContext.md next steps.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"session_id": map[string]string{"type": "string"},
					"milestone":  map[string]string{"type": "string"},
					"status":     map[string]string{"type": "string"},
					"next_steps": map[string]interface{}{
						"type":  "array",
						"items": map[string]string{"type": "string"},
					},
				},
				"required": []string{"milestone", "status"},
			},
		},
		{
			"name":        "report_manual_changes",
			"description": "Captures out-of-band changes before an agent handoff or context compaction.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"session_id":  map[string]string{"type": "string"},
					"files":       map[string]interface{}{"type": "array", "items": map[string]string{"type": "string"}},
					"description": map[string]string{"type": "string"},
				},
				"required": []string{"description"},
			},
		},
		{
			"name":        "request_checkpoint",
			"description": "Triggers a full deterministic checkpoint snapshot from the agent.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"session_id": map[string]string{"type": "string"},
					"focus":      map[string]string{"type": "string"},
				},
			},
		},
	}
}

func (s *Server) handleToolCall(req JSONRPCRequest) JSONRPCResponse {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: -32602, Message: "Invalid tool call params"},
		}
	}

	sessionID, _ := params.Arguments["session_id"].(string)
	if sessionID == "" {
		sessionID = "default"
	}

	var result interface{}
	var err error

	switch params.Name {
	case "read_active_context":
		result, err = ToolReadActiveContext(s.bankDir, sessionID)

	case "read_static_context":
		var files []string
		if rawFiles, ok := params.Arguments["files"].([]interface{}); ok {
			for _, f := range rawFiles {
				if str, ok := f.(string); ok {
					files = append(files, str)
				}
			}
		}
		result, err = ToolReadStaticContext(s.bankDir, sessionID, files)

	case "append_memory_delta":
		file, _ := params.Arguments["file"].(string)
		section, _ := params.Arguments["section"].(string)
		content, _ := params.Arguments["content"].(string)
		result, err = ToolAppendMemoryDelta(s.bankDir, sessionID, file, section, content)

	case "update_milestone":
		milestone, _ := params.Arguments["milestone"].(string)
		status, _ := params.Arguments["status"].(string)
		var nextSteps []string
		if rawSteps, ok := params.Arguments["next_steps"].([]interface{}); ok {
			for _, st := range rawSteps {
				if str, ok := st.(string); ok {
					nextSteps = append(nextSteps, str)
				}
			}
		}
		result, err = ToolUpdateMilestone(s.bankDir, sessionID, milestone, status, nextSteps)

	case "report_manual_changes":
		var files []string
		if rawFiles, ok := params.Arguments["files"].([]interface{}); ok {
			for _, f := range rawFiles {
				if str, ok := f.(string); ok {
					files = append(files, str)
				}
			}
		}
		desc, _ := params.Arguments["description"].(string)
		result, err = ToolReportManualChanges(s.bankDir, s.repoDir, sessionID, files, desc)

	case "request_checkpoint":
		focus, _ := params.Arguments["focus"].(string)
		result, err = ToolRequestCheckpoint(s.bankDir, s.repoDir, sessionID, focus)

	default:
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: -32601, Message: fmt.Sprintf("Tool not found: %s", params.Name)},
		}
	}

	if err != nil {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"isError": true,
				"content": []map[string]string{
					{"type": "text", "text": fmt.Sprintf("Error: %v", err)},
				},
			},
		}
	}

	resJSON, _ := json.Marshal(result)
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": string(resJSON)},
			},
		},
	}
}

func (s *Server) handleResourceRead(req JSONRPCRequest) JSONRPCResponse {
	var params struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: -32602, Message: "Invalid resource params"},
		}
	}

	content, mime, err := ReadResource(s.bankDir, params.URI)
	if err != nil {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: -32000, Message: err.Error()},
		}
	}

	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"contents": []map[string]string{
				{
					"uri":      params.URI,
					"mimeType": mime,
					"text":     content,
				},
			},
		},
	}
}

func (s *Server) handlePromptGet(req JSONRPCRequest) JSONRPCResponse {
	var params struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(req.Params, &params)

	if params.Name == "resume-session" {
		activeData, _ := os.ReadFile(filepath.Join(s.bankDir, "activeContext.md"))
		progressData, _ := os.ReadFile(filepath.Join(s.bankDir, "progress.md"))
		promptText := GetResumePrompt(string(activeData), string(progressData))

		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"messages": []map[string]interface{}{
					{
						"role": "user",
						"content": map[string]string{
							"type": "text",
							"text": promptText,
						},
					},
				},
			},
		}
	}

	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Error:   &RPCError{Code: -32601, Message: fmt.Sprintf("Prompt not found: %s", params.Name)},
	}
}
