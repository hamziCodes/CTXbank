# The Memory Bank & Project Lifecycle Engine — Architecture & Implementation Spec

## 0. One Critical Piece of Prior Art You Should Know Before Building This

Before the architecture: the file layout you've specified (`activeContext.md`, `progress.md`, an architecture doc, a research/ingestion flow) is **not a blank canvas — it's already a de facto standard**. Cline's **Memory Bank** pattern (`memory-bank/projectbrief.md`, `productContext.md`, `activeContext.md`, `systemPatterns.md`, `techContext.md`, `progress.md`) has existed since 2024, is documented in Cline's own docs, and has been forked into community MCP servers (e.g. `cline-mcp-memory-bank`) with `decisionLog.md` and `projectContext.md` extensions. Multiple agent tools (Cline, Roo Code, and prompt-adapted Cursor/Windsurf setups) already know to look for this exact hierarchy when told to "initialize memory bank."

**What Cline's Memory Bank does *not* have — and what your spec is actually asking for — is everything deterministic:** it is a pure prompt-engineering convention. The agent itself decides when to read it, when to write it, and what "current" means; there's no git-diffing, no delta hashing, no crash-safe atomic writes, no CLI, no MCP tool contract, no cross-agent enforcement. It works only as well as the current agent's discipline about following its custom instructions.

**Implication for this build**: don't invent a sixth incompatible file-naming scheme. **Adopt the Cline file names as the human-readable projection layer**, and add underneath them exactly what's missing — a deterministic binary that writes them, hashes them, diffs them, and serves them over MCP so *any* agent gets consistent state regardless of whether it "remembers" to follow memory-bank instructions. This is the same design principle as the prior report's `AGENTS.md` compatibility call: interoperate with the standard that already has adoption, don't fragment it.

---

## 1. System Architecture & Data Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              REPO ROOT                                       │
│                                                                               │
│  .cursorrules / .cursor/rules/   CLAUDE.md   .agentrules                     │
│  .github/copilot-instructions.md      (auto-generated stubs, see §1.2)       │
│         │ all point to →                                                     │
│         ▼                                                                    │
│  memory-bank/                                                                │
│    projectbrief.md      systemPatterns.md      techContext.md               │
│    productContext.md    activeContext.md  ◄── changes every checkpoint      │
│    progress.md          decisionLog.md                                      │
│    .state/                                                                   │
│      manifest.json      ◄── hash + version per file, delta-sync source of truth
│      checkpoints/ckpt_*.json                                                 │
│  research/inbox/        ◄── raw brainstorm dumps land here                  │
└─────────────────────────────────────────────────────────────────────────────┘
        ▲                          ▲                              ▲
        │ writes                  │ reads/writes                 │ reads/writes
┌───────┴────────┐        ┌────────┴─────────┐          ┌─────────┴──────────┐
│  `ctx` binary   │◄──────►│  ctx MCP server   │◄────────►│  Coding agents      │
│  (Go, static)   │  IPC   │  (stdio/HTTP)      │  MCP     │  Cursor / Claude    │
│                 │        │                    │          │  Code / Antigravity │
│  git2 / tree-   │        │  tools:            │          │  / Qoder / VS Code  │
│  sitter / diff  │        │  read_active_ctx   │          │  Agent / Blackbox   │
│  engine         │        │  append_delta      │          └─────────────────────┘
└─────────────────┘        │  update_milestone  │
        ▲                  │  report_manual_    │
        │                  │  changes           │
   ┌────┴─────┐            └────────────────────┘
   │ Developer │  direct CLI use — never requires an agent to be running
   └───────────┘
```

**Why this avoids being "just an SDK":** the binary is the *only* required component and works with zero agents in the loop — `ctx status`/`ctx pause`/`ctx resume` are useful to a developer who never opens an AI tool. The MCP server is a thin, optional wrapper that exposes the *same* filesystem state, so agent-native and human-only workflows read one source of truth instead of forking into "the SDK's model of the project" vs. "the real project." An SDK requires you to write code against it; this requires nothing but running a binary and optionally pointing an agent at a stdio MCP process — the integration surface is a file format and a protocol, not an API contract a developer has to code to.

### 1.2 The rules-file problem, and why to keep it minimal

Auto-writing five different vendor rule files is a maintenance trap (each vendor changes its format independently — `.cursor/rules/` moved from single-file to directory-based MDC rules; Copilot's location and precedence have shifted). Solve it the way `AGENTS.md` itself solves it: **write one canonical stub into each vendor location that says almost nothing except "read `memory-bank/activeContext.md` and `memory-bank/progress.md` before doing anything; update them before finishing."** Never let the CLI regenerate vendor-specific formatting logic beyond that one-line pointer — the content lives in `memory-bank/`, not duplicated five ways.

---

## 2. Memory Bank Schema, Token Efficiency, and Delta Sync

### 2.1 Directory layout

```
memory-bank/
├── projectbrief.md       # Rarely changes — scope, constraints, non-goals
├── productContext.md     # Rarely changes — why, for whom, UX goals
├── systemPatterns.md     # Changes on architectural decisions only
├── techContext.md        # Changes when stack/deps change
├── decisionLog.md        # Append-only — one entry per significant decision, never rewritten
├── activeContext.md      # Changes every checkpoint — the hot file
├── progress.md           # Changes every checkpoint — milestone ledger
└── .state/
    ├── manifest.json      # hash + byte-size + last-modified-checkpoint-id per file
    └── checkpoints/
        └── ckpt_<id>.json # full deterministic snapshot (from the resumption-engine schema)
```

Each file has a **volatility class** the sync engine uses to decide default read scope:

| Class | Files | Re-read policy |
|---|---|---|
| Static | `projectbrief.md`, `productContext.md` | Read once per agent session, cache by hash |
| Semi-static | `systemPatterns.md`, `techContext.md`, `decisionLog.md` | Read if manifest hash changed since last read |
| Hot | `activeContext.md`, `progress.md` | Always read in full — they're small by design (target <150 lines each) |

### 2.2 Token discipline: format rules, not just intent

"Strip conversational fluff" needs an enforced format, or every agent re-introduces prose over time. `activeContext.md` template (agent- and human-legible, no narrative):

```markdown
# Active Context — updated ckpt_20260903_2137
## Focus
OCR retry backoff logic (src/ocr/retry.ts)
## Recent (last 3 checkpoints)
- Added jitter to retry delay (ckpt_2136)
- Fixed test flake in ocr.test.ts (ckpt_2135)
## Next steps
1. Clamp jitter to [0, maxDelay] — retry.ts:44
2. Re-run `npm test -- ocr`
3. Wire retry count into telemetry
## Open decisions
- Max-retry policy: waiting on product spec confirmation
```

Bullet-only, fixed sections, no adjectives, hard line-count budgets enforced by `ctx lint-memory` (fails CI or a pre-commit hook if `activeContext.md` exceeds ~150 lines — forces summarization instead of accretion, the same failure mode that makes hand-kept `NOTES.md` files unusable after month three).

### 2.3 Delta-sync algorithm — never re-read everything

```
manifest.json:
{
  "files": {
    "activeContext.md":  { "sha256": "9f2a...", "bytes": 812,  "last_ckpt": "ckpt_2137" },
    "progress.md":        { "sha256": "3bc1...", "bytes": 2044, "last_ckpt": "ckpt_2130" },
    "systemPatterns.md":  { "sha256": "aa02...", "bytes": 4110, "last_ckpt": "ckpt_1990" }
  },
  "agent_read_cursors": {
    "claude-code:session-88f": { "activeContext.md": "9f2a...", "progress.md": "aa11..." }
  }
}
```

- On every write, `ctx` recomputes SHA-256 of the file body (frontmatter excluded) and bumps `manifest.json`.
- On `read_active_context()` (MCP call), the server compares the *caller's last known hash* (tracked per agent session ID in `agent_read_cursors`) against the current hash. If unchanged, it returns a 20-byte `{"unchanged": true}` instead of the file body — this is the actual token savings mechanism, not chunking within a file.
- For files that do change, only **that file's** body is returned — never a full-corpus re-index. Cross-file consistency isn't needed because the volatility classes above mean static/semi-static files change rarely enough that per-file hashing alone gives >90% of the token savings without needing a line-level diff protocol.
- `decisionLog.md` is append-only by construction (never rewritten, only appended), so its delta is trivially "everything after byte offset N," tracked the same way `tail -f` tracks a log file — no hashing needed for that file, just a stored offset.

This is intentionally simpler than a true CRDT or block-hash content-addressed store (used by tools like restic/Perkeep) — at the file sizes a project's memory bank actually reaches (a few KB to low hundreds of KB), whole-file hashing is sufficient, and reaching for chunk-level rolling hashes would be solving a scale problem this product doesn't have.

### 2.4 Research ingestion pipeline (`ctx ingest`)

```
research/inbox/*.md  ──►  ctx ingest  ──►  Tier-0: strip markdown boilerplate,
                                            dedupe near-identical paragraphs
                                            (simhash), extract headers as
                                            candidate milestones
                          │
                          ▼ (only if a Tier-1/2 model is configured)
                          Tier-1/2: classify each extracted chunk into
                          {projectbrief | productContext | systemPatterns |
                           techContext | discard} and propose a diff against
                          the existing file, shown to the user for confirm
                          before writing
                          │
                          ▼
                          research/inbox/*.md moved to research/archive/
                          (never deleted — ingestion is a proposal, not a
                          silent overwrite)
```

`ctx ingest` **always** ends in a diff-review step (`git diff`-style, in-terminal), never an autonomous silent rewrite of the memory bank — the failure mode of "the tool guessed wrong and now the architecture doc lies" is worse than the friction of a confirmation prompt.

---

## 3. CLI Command Specification

| Command | Purpose | LLM required? |
|---|---|---|
| `ctx init` | Scaffolds `memory-bank/`, `.state/`, writes minimal rule stubs to detected vendor locations | No |
| `ctx audit` | Brownfield deep scan (see §3.1) | Optional (Tier-0 structural map is free; narrative summary is Tier-1/2) |
| `ctx ingest <path>` | Distill `research/inbox/` into memory-bank files, always with review | Tier-1/2 for classification; Tier-0 fallback = dump under a dated heading in `productContext.md` untriaged |
| `ctx status` | Zero-token terminal brief or `--json` for agents | No |
| `ctx pause` / `ctx checkpoint` | Snapshot + interactive manual-edit capture | No (interactive prompt is plain terminal I/O) |
| `ctx resume` | Render pickup brief | No (Tier-1/2 only enriches prose) |
| `ctx serve --mcp` | Launch MCP server over stdio or HTTP | N/A |
| `ctx lint-memory` | Enforce size/format budgets, fail CI if violated | No |

### 3.1 `ctx audit` — brownfield reconnaissance, concretely

Deterministic passes, in order, all before any LLM call:
1. **Manifest scan**: `package.json`/`Cargo.toml`/`pyproject.toml`/`go.mod` → language, deps, scripts → seeds `techContext.md`.
2. **Structural scan**: directory heuristics (`src/api`, `prisma/schema.*`, `**/*.controller.ts`, `**/migrations/`) → seeds a draft component map in `systemPatterns.md`.
3. **Symbol scan** (tree-sitter): exported functions/classes per file, entrypoints (`main`, route handlers) → candidate "core business logic" list.
4. **Git velocity scan**: commit frequency per path over the last 90 days → which subsystem is "active" (feeds `activeContext.md`'s initial guess), plus current branch/dirty-state exactly as in the resumption-engine checkpoint schema.

Only after those four produce a structured JSON intermediate does an optional Tier-1/2 pass turn it into prose for `productContext.md`'s "why this exists" section — a purely narrative field with no deterministic source. If no model is configured, that field is left as a TODO placeholder rather than fabricated, since guessing business intent from code is exactly the kind of hallucination-prone task that should degrade gracefully to "ask the human," not produce confident fiction.

### 3.2 Terminal UX mockups

```
$ ctx status
┌─ khata-app · feature/ocr-retry ───────────────────────────────┐
│ ● 2 files dirty   ✗ 1 failing test   ⏱ idle 3d                │
│ Focus: OCR retry backoff                                       │
│ Next:  clamp jitter in retry.ts:44                              │
└─────────────────────────────────────────────────────────────┘

$ ctx pause
Snapshotting… done (created ckpt_20260903_2137)
Detected 1 file changed outside any agent session: src/ocr/retry.ts
  → What did you change here? (one line, Enter to skip)
  > tried clamping jitter manually, still flaky
Saved to activeContext.md → decisionLog.md
```

Styling: sparse, ASCII-safe by default (auto-detect TTY color support), degrading cleanly on CI logs — a bubbletea/charm-style TUI is reserved for `ctx audit`'s progress view (long-running, benefits from a live tree render) but `status`/`pause`/`resume` stay single-shot text output so they're pipeable and agent-parseable with `--json`.

---

## 4. MCP Server Specification

```
Server: memory-bank-mcp

Tools:
  read_active_context(session_id)
      → { activeContext, progress, unchanged_files: [...] }
      Uses agent_read_cursors to skip unchanged hot files.

  read_static_context(session_id, files?: string[])
      → returns projectbrief/productContext/systemPatterns/techContext,
        each individually hash-gated as in §2.3

  append_memory_delta(session_id, file, section, content)
      → appends under a named section (never a raw overwrite);
        validated against §2.2 line-budget before write

  update_milestone(session_id, milestone, status, next_steps[])
      → structured write to progress.md + activeContext.md's "Next steps"

  report_manual_changes(session_id, files[], description)
      → the automated equivalent of the ctx pause interactive prompt,
        for an agent that detects it's about to hand off mid-task

  request_checkpoint(session_id)
      → triggers a full deterministic `ctx pause` from the agent side
        (e.g. before an agent context-compaction event)

Resources:
  memory://activeContext
  memory://progress
  memory://checkpoint/latest
  memory://decisionLog (append-only stream)

Prompts:
  "resume-session" → pre-loads read_active_context() into the agent's
                      first turn, mirroring the resumption-engine's
                      generate_pickup_brief("agent")
```

**How Cursor/Antigravity/Claude Code invoke this autonomously**: each already supports MCP tool discovery at session start. The one-line rule stub from §1.2 is what makes invocation *reliable* rather than *possible* — MCP exposes the capability, but an agent that never gets told "you must call `read_active_context` before touching code" won't reliably do so unprompted. The stub is the missing enforcement layer that a bare MCP server doesn't provide on its own; this mirrors why Cline's Memory Bank works today purely off custom instructions even without any binary or protocol behind it — the instruction-to-consult-memory is doing real work, not the storage format.

---

## 5. Technology Stack

| Concern | Recommendation | Why |
|---|---|---|
| Core binary language | **Go** | Single static binary, trivial cross-compilation (`GOOS`/`GOARCH` matrix) with zero runtime, good-enough performance for git/filesystem-bound work, simpler contributor onboarding than Rust for an OSS project soliciting external PRs |
| Git operations | Shell out to system `git` for porcelain commands (`diff --stat`, `log`, `status --porcelain`), reserve `go-git` (pure Go) only where shelling out is awkward (e.g. reading blobs for hashing) | System `git` is always more correct/up-to-date than any reimplementation, and shelling out avoids `libgit2`/cgo cross-compilation pain entirely — a strong reason to prefer this over Rust+`git2` (which pulls in a C dependency and complicates static linking) |
| Symbol/AST diffing | `tree-sitter` via Go bindings (`smacker/go-tree-sitter` or the official Go grammar bindings) | Battle-tested, incremental parsing, grammars exist for every language this audience uses (TS/JS, Python, Go, Rust, Swift/Kotlin for the mobile work mentioned) |
| TUI (audit progress view only) | `charmbracelet/bubbletea` + `lipgloss` | De facto Go TUI standard, already familiar to contributors coming from `gh`, `lazygit`, etc. |
| MCP server | Official `modelcontextprotocol/go-sdk` (or the community Go SDK if the official one lags) over stdio transport primarily, HTTP/SSE as a secondary transport for remote/daemon use | Matches the ecosystem's own reference implementations, minimizes protocol-drift risk as MCP's spec version increments |
| Local LLM tier | Ollama's HTTP API (`localhost:11434`) as the default local backend, since it's the most widely pre-installed local runtime; `llama.cpp` server as a fallback for GPU-less minimal footprints | Avoids embedding an inference engine in the binary itself — keeps the core `ctx` binary small and the LLM tier fully pluggable |
| Config/BYOK | Single `.context/config.toml`, keys read from environment variables only, never written to disk in plaintext by the tool itself | Standard practice; keeps the zero-cost/local-first posture honest — no key material persisted where a `git add -A` could leak it |

**Binary size target**: <15MB static binary with no CGO, which rules out `libgit2`-based bindings (a strong secondary argument for shelling out to system git over `git2`-style bindings, beyond the correctness point above).

---

## 6. Roadmap, Failure Modes, and Launch Strategy

### 6.1 Failure modes and safeguards

| Scenario | Safeguard |
|---|---|
| Merge conflict in `memory-bank/*.md` across branches | Treat `activeContext.md`/`progress.md` as **branch-local by default** (gitignored, regenerated per-branch from `.state/checkpoints/`); only `projectbrief.md`/`productContext.md`/`systemPatterns.md` are meant to be shared/committed, and those change rarely enough that conflicts are rare and resolve like any other doc conflict |
| Detached HEAD / mid-rebase | `ctx status`/`ctx pause` detect `.git/rebase-merge` or `.git/rebase-apply` and refuse to write a checkpoint claiming a stable state — surfaces "resolve rebase first" instead of silently checkpointing a transient tree |
| Sudden crash mid-write | Atomic write pattern for every file: write to `<name>.md.tmp`, `fsync`, `rename()` — the same pattern used by `git` itself for its own object writes; a crash never leaves a half-written file visible to a reader |
| Manifest/file drift (someone hand-edits `activeContext.md` outside the tool) | `ctx status` recomputes the hash on every run; a mismatch against `manifest.json` is surfaced as "externally modified since last checkpoint" rather than silently trusted or silently overwritten |
| Ingestion misclassification | §2.4's mandatory diff-review step — never autonomous |
| Stale memory bank (weeks old, code has moved on) | `ctx status` computes days-since-last-checkpoint vs. commits-since-last-checkpoint and flags staleness explicitly, same mechanism as the prior resumption-engine report |

### 6.2 Contribution roadmap

1. **v0.1** — `init`, `status`, `pause`, `resume` against the Cline-compatible file layout; zero LLM dependency; ships as a single binary via `go install` + GitHub Releases.
2. **v0.2** — MCP server (`serve --mcp`), rule-stub auto-writer for the four major vendor formats.
3. **v0.3** — `audit` (deterministic passes only), `lint-memory` for CI enforcement.
4. **v0.4** — `ingest` with Tier-0 heuristics; Ollama-backed Tier-1 classification/summarization behind a config flag.
5. **v0.5** — BYOK Tier-2, cross-project `ctx list` dashboard for the polymath workflow.
6. **Community phase** — publish the MCP tool contract and file schema as a short spec document (mirroring how `AGENTS.md` published its convention) so other implementations (not just this Go binary) can interoperate; explicitly invite Cline/Roo maintainers to comment, since compatibility with their existing user base is the fastest adoption path available.

**Positioning for launch**: not "yet another AI memory tool" — pitch it as *"Cline's Memory Bank, minus the part where it only works if the agent remembers to follow instructions."* That framing is honest about prior art, gives early adopters (who already know the file names) zero migration cost, and differentiates on the one thing genuinely missing: a deterministic, crash-safe, cross-agent-enforced layer underneath a convention people already trust.
