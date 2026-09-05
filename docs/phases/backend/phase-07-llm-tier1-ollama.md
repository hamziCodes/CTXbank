# Phase 07: Tier-1 Local LLM Pipeline (Ollama / llama.cpp)

**Domain:** Backend / AI Integration  
**Status:** Ready to Implement  
**Specification Source:** `CTXbank.md` (§2.4, §3.1, §5)  

---

## 1. Phase Objective & Boundary
Implement optional Tier-1 AI intelligence for research classification and narrative generation:
- Connect to local Ollama server (`http://localhost:11434`) or `llama.cpp` server.
- Classify ingested research chunks into `{projectbrief, productContext, systemPatterns, techContext, discard}`.
- Enrich brownfield audit reports with concise business intent descriptions in `productContext.md`.
- Graceful Degradation: If Ollama is offline or unconfigured, the system degrades silently to deterministic Tier-0 mode with zero errors.

*Exclusions:* Cloud API dependencies or hardcoded API keys.

---

## 2. Prerequisites & Dependencies
- Phase 06 complete.
- `.context/config.toml` parser.

---

## 3. Step-by-Step Implementation Steps
1. **Define LLM Interface (`internal/llm/client.go`):**
   - Interface `Client`: `Classify(ctx context.Context, text string) (TargetDoc, error)`, `Summarize(ctx context.Context, audit AuditResult) (string, error)`.
2. **Implement Ollama REST Client (`internal/llm/ollama.go`):**
   - Implement HTTP calls to `/api/generate` with low temperature (`0.1`) and structured JSON schema outputs.
   - Enforce a 3-second timeout for local availability checks.
3. **Implement Structured Prompts (`internal/llm/prompts.go`):**
   - Zero-hallucination system prompt: instruct model to output only strict JSON matching requested taxonomy.
4. **Wire Fallback Mechanism (`internal/llm/fallback.go`):**
   - If `ollama` is unreachable, return `ErrOffline`, allowing callers to execute pure Tier-0 operations.
5. **Develop Unit Tests (`internal/llm/llm_test.go`):**
   - Test against a mock HTTP server simulating valid, malformed, and timeout responses.

---

## 4. Concrete Deliverables
- `internal/llm/client.go`, `internal/llm/ollama.go`, `internal/llm/prompts.go`, `internal/llm/fallback.go`
- `internal/llm/llm_test.go`

---

## 5. Validation & Verification Routine
```bash
go test -v ./internal/llm/...
```
- Verify fallback occurs cleanly with mock server down.

---

## 6. Security & Vulnerability Audit
- Local-only by default: restrict endpoints to loopback (`127.0.0.1`, `localhost`).
- Never leak repository file contents to unverified external URLs.
