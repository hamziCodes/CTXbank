# Phase 11: Terminal UX & Interactive Bubbletea TUI

**Domain:** Frontend / Terminal UI  
**Status:** Ready to Implement  
**Specification Source:** `CTXbank.md` (§3.2, §5)  

---

## 1. Phase Objective & Boundary
Implement high-polish interactive terminal user interfaces:
- Live progress and component tree viewer for long-running `ctx audit` using `charmbracelet/bubbletea` and `lipgloss`.
- Sparse, single-shot, pipeable cards for `ctx status`, `ctx pause`, and `ctx resume`.
- Automatic detection of TTY color support; zero ANSI escape codes when redirected to a pipe or on CI.

*Exclusions:* Web GUI or browser-based frontends.

---

## 2. Prerequisites & Dependencies
- `github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/lipgloss`.
- Phases 03 & 04 complete.

---

## 3. Step-by-Step Implementation Steps
1. **Implement TTY Detection & Formatting (`internal/tui/format.go`):**
   - Check if stdout is an interactive terminal (`mattn/go-isatty`).
   - Define fallback ASCII border styles (`+---`, `|`) for non-TTY environments.
2. **Implement Bubbletea Audit UI (`internal/tui/audit_view.go`):**
   - Bubbletea Model/Update/View cycle.
   - Live spinners for passes (manifest, structure, symbol, git velocity).
   - Interactive tree component to browse discovered architecture layers.
3. **Style Cards with Lipgloss (`internal/tui/styles.go`):**
   - Electric Indigo border (`#7C3AED`), Cyan highlights (`#06B6D4`), Emerald indicators (`#10B981`).
4. **Develop Tests (`internal/tui/tui_test.go`):**
   - Assert pipe output contains clean plaintext without ANSI escape codes.

---

## 4. Concrete Deliverables
- `internal/tui/format.go`, `internal/tui/styles.go`, `internal/tui/audit_view.go`
- `internal/tui/tui_test.go`

---

## 5. Validation & Verification Routine
```bash
go test -v ./internal/tui/...
./bin/ctx status | cat
```
- Verify zero control characters in piped output.

---

## 6. Security & Vulnerability Audit
- Ensure clean terminal reset on exit/panic (restoring cursor and normal buffer).
