# Phase 08: Model Context Protocol (MCP) Server

**Domain:** Integration / Agent Interop  
**Status:** Ready to Implement  
**Specification Source:** `CTXbank.md` (§1, §2.3, §4)  

---

## 1. Phase Objective & Boundary
Implement the full Model Context Protocol (MCP) server conforming to standard MCP over stdio (with optional HTTP transport):
- Provide tools: `read_active_context`, `read_static_context`, `append_memory_delta`, `update_milestone`, `report_manual_changes`, `request_checkpoint`.
- Expose resources: `memory://activeContext`, `memory://progress`, `memory://checkpoint/latest`, `memory://decisionLog`.
- Implement per-session delta caching using `agent_read_cursors` in `manifest.json`. Return lightweight `{"unchanged": true}` payload when hashes match.

*Exclusions:* Direct UI rendering.

---

## 2. Prerequisites & Dependencies
- Phases 01, 02, and 05 complete.
- MCP Go SDK (`github.com/mark3labs/mcp-go`).

---

## 3. Step-by-Step Implementation Steps
1. **Implement Agent Session Tracker (`internal/mcp/session.go`):**
   - Manage read cursor hashes per agent session ID in `manifest.json`.
   - Implement `HasFileChanged(sessionID, filename string) (bool, string, error)`.
2. **Implement MCP Tools (`internal/mcp/tools.go`):**
   - `read_active_context(session_id)`: returns `activeContext.md` and `progress.md` or `{"unchanged": true}`.
   - `read_static_context(session_id, files)`: returns requested static/semi-static files, hash-gated.
   - `append_memory_delta(session_id, file, section, content)`: validates line budget, performs atomic append under section.
   - `update_milestone(session_id, milestone, status, next_steps)`: updates `progress.md` and `activeContext.md`.
   - `report_manual_changes(session_id, files, description)`: automated checkpointing.
   - `request_checkpoint(session_id)`: full deterministic pause.
3. **Implement MCP Resources (`internal/mcp/resources.go`):**
   - Register URI templates for `memory://*`.
4. **Implement MCP Stdio Server (`internal/mcp/server.go`):**
   - Bind tools, resources, and JSON-RPC stdio transport.
5. **Implement Server Entrypoints (`cmd/ctx-mcp/main.go`, `cmd/ctx/cmd_serve.go`):**
   - Allow invocation via `ctx serve --mcp` or standalone `ctx-mcp`.
6. **Develop Protocol Tests (`internal/mcp/mcp_test.go`):**
   - Send JSON-RPC payloads; assert tool discovery and tool invocation responses.
   - Verify second call with same session ID yields `unchanged: true`.

---

## 4. Concrete Deliverables
- `internal/mcp/session.go`, `internal/mcp/tools.go`, `internal/mcp/resources.go`, `internal/mcp/server.go`
- `cmd/ctx-mcp/main.go` & `cmd/ctx/cmd_serve.go`
- `internal/mcp/mcp_test.go`

---

## 5. Validation & Verification Routine
```bash
go test -v ./internal/mcp/...
```
- Assert `read_active_context` saves tokens by returning `unchanged: true` when file is untouched.

---

## 6. Security & Vulnerability Audit
- Sanitize tool parameters: constrain all file reads/writes strictly within `memory-bank/`.
