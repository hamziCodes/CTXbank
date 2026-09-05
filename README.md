# CTXbank

> **A simple, bulletproof memory bank for AI coding assistants and developers.**  
> Never lose your project context, never let your AI hallucinate rules, and never waste tokens re-reading files that haven't changed.

---

## What is CTXbank? (The Simple Explanation)

If you use AI coding assistants like **Cursor, Claude Code, Antigravity, or VS Code**, you've probably noticed three frustrating problems:

1. **AI Amnesia:** When a chat gets too long, the AI forgets decisions you made an hour ago.
2. **Context Switching:** You take a break for a weekend, come back, and have no idea where you or the AI left off.
3. **Token Burn:** Every time an AI reads your entire codebase just to fix a single bug, you burn thousands of tokens on things that haven't changed.

**CTXbank solves this.** Think of it as a **universal save-state system and brain backup** for your project. It keeps a clean folder of notes called `memory-bank/`, tracks changes with safe snapshots, and tells your AI exactly what it needs to know—and nothing more.

---

## How It Works (In 3 Plain-English Concepts)

### 1. The Memory Bank (The Knowledge Base)
When you run `ctx init`, CTXbank creates a folder named `memory-bank/` with 7 simple markdown files:
- **`projectbrief.md`** — What are we building? (The core mission).
- **`productContext.md`** — Why does this exist and how should users experience it?
- **`systemPatterns.md`** — How is the code structured? (Architecture rules).
- **`techContext.md`** — What tech stack, libraries, and tools are we using?
- **`activeContext.md`** — What are we working on *right now*? (Kept short: under 150 lines).
- **`progress.md`** — What's done, what's in progress, and what's next?
- **`decisionLog.md`** — Why did we make key decisions? (Prevents circular debates).

### 2. The Checkpoint Engine (The Save Button)
Whenever you finish a task, want to switch git branches, or log off for the night, run:
```bash
ctx pause
```
CTXbank takes a snapshot of your git state, records your recent changes, and creates a restore point. When you return, just run:
```bash
ctx resume
```
It immediately prints a quick summary of what you were doing so you (or your AI) can pick up work in 5 seconds.

### 3. The Agent Bridge (Zero Token Waste)
Instead of your AI re-reading hundreds of files on every message, CTXbank acts as a local server (using the Model Context Protocol / MCP). 
- If no files changed, CTXbank sends a tiny **72-byte** response: `{"unchanged": true}`.
- Your AI reads only what is new, saving you context window space and money.

---

## Quickstart (Under 60 Seconds)

### Step 1: Install CTXbank

Choose the easiest install method for your system:

#### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/hamziCodes/CTXbank/main/scripts/install.ps1 | iex
```

#### macOS & Linux (Terminal)
```bash
curl -fsSL https://raw.githubusercontent.com/hamziCodes/CTXbank/main/scripts/install.sh | sh
```

#### Go Developers
```bash
go install github.com/hamziCodes/CTXbank/cmd/ctx@latest
```

*(Or download the standalone `.zip` / `.tar.gz` for your OS directly from [GitHub Releases](https://github.com/hamziCodes/CTXbank/releases)).*

### Step 2: Initialize Your Project
Open a terminal in your project root and run:
```bash
ctx init
```
This generates your `memory-bank/` folder and setup rules for Cursor, Claude, and GitHub Copilot.

### Step 3: Scan an Existing Project (Optional)
If you already have code in your project, run:
```bash
ctx audit --apply
```
CTXbank will scan your dependencies, project structure, and exported functions, and automatically fill out `techContext.md` and `systemPatterns.md` with factual information.

---

## Everyday Commands Cheat Sheet

| Command | What it does |
|---|---|
| `ctx ui` (or `ctx dashboard`) | Launches the embedded VERTEX web dashboard with interactive architecture graphs. |
| `ctx status` | Shows a clean terminal card with your git branch, dirty files, and current focus. |
| `ctx pause` | Saves a snapshot, asks for any quick notes on what you did, and logs it. |
| `ctx resume` | Prints an instant summary of what you were working on so you can pick up immediately. |
| `ctx audit` | Scans your project structure, code symbols, and git history without changing code. |
| `ctx ingest <file.md>` | Ingests raw brainstorm notes, removes duplicate text automatically, and proposes updates. |
| `ctx list` | Shows a dashboard of all your active CTXbank projects across your computer. |
| `ctx lint-memory` | Checks that your memory bank files aren't bloated (keeps activeContext under 150 lines). |
| `ctx serve --mcp` | Starts the background server so AI tools can talk to CTXbank directly. |

---

## Interactive Visual Interfaces

CTXbank offers two rich, zero-learning-curve visual interfaces built to the high-precision **VERTEX Universal Design System**:

### 1. Embedded Web Dashboard (`ctx ui`)
Run a single command in your terminal:
```bash
ctx ui
```
This instantly boots an embedded local web dashboard (served directly from the static binary on `http://localhost:4242`) with:
- **Interactive SVG Architecture Graph:** Clickable nodes representing core components, test suites, and memory files with live line counts and status rings.
- **Visual Memory Cards:** Live preview and tabbed editor for all 7 memory bank files with real-time budget meter warning before hitting line limits.
- **Drag-and-Drop Ingestion:** Drop meeting notes, markdown brainstorms, or research documents into the browser to auto-deduplicate (SimHash) and merge into memory files.
- **Checkpoint Manager:** 1-click snapshot creation, rollback preview, and branch audit timeline.

### 2. Native VS Code & Cursor Extension
Install the official extension directly in VS Code or Cursor:
```bash
code --install-extension https://github.com/hamziCodes/CTXbank/releases/download/v0.1.0/ctxbank-0.1.0.vsix
```
*(Or download `ctxbank-0.1.0.vsix` from [Releases](https://github.com/hamziCodes/CTXbank/releases/tag/v0.1.0) and run `Extensions -> Install from VSIX...`).*

**Features inside your IDE:**
- **Activity Bar Sidebar:** Dedicated CTXbank icon with 3 collapsible views:
  - **Active Focus:** Current task objective, token budget meter, and active branch.
  - **Memory Bank Files:** Direct file explorer with line counts and quick-edit actions.
  - **Checkpoints:** Snapshot history with 1-click restore points.
- **Status Bar Integration:** Shows current memory token health (`CTX: 48/150 lines`) right in your bottom status bar.
- **Visual Webview Panel:** Run `CTXbank: Open Visual Architecture Tree` to inspect your project's component tree right inside an editor tab.


---

## How to Connect Your AI Tools

### Cursor
Add CTXbank to your Cursor MCP settings (`Cursor Settings -> Features -> MCP`):
- **Name:** `ctxbank`
- **Type:** `command`
- **Command:** `d:/path/to/CTXbank/bin/ctx.exe serve --mcp`

### Claude Code / Antigravity / Other Agents
Add this to your project or global `mcpServers` configuration:
```json
{
  "mcpServers": {
    "ctxbank": {
      "command": "d:/path/to/CTXbank/bin/ctx.exe",
      "args": ["serve", "--mcp"]
    }
  }
}
```

---

## Why CTXbank Won't Break Your Files

Many automated tools break code by writing halfway through a crash. CTXbank uses **Atomic Writes**:
1. It writes your updates to a temporary draft file (`.tmp`).
2. It forces the operating system to flush the data safely to physical storage.
3. It instantly swaps the new file in place.

If your laptop loses power mid-save, your files will never be corrupted or half-written.

---

## License
MIT License. Free and open source for individual developers and teams.
