# Project Profile: CTXbank

> **Canonical Design Authority & Adaptation Profile**  
> Governs all UI/UX implementations across the CTXbank Web Dashboard (`ctx ui`) and the VS Code Extension.

**Toolkit version in use:** v1.0.0  
**Profile last reviewed:** 2026-09-06 by CTXbank Core Architecture Team  

---

## 1. What this project is

| Field | Value |
| :--- | :--- |
| Product name | **CTXbank** |
| One-line purpose | Deterministic memory bank and project lifecycle engine for AI coding assistants and developers. |
| Primary user | Developers and prompt engineers collaborating with AI coding agents (Cursor, Claude Code, Antigravity, VS Code) needing zero context drift and instant pickup briefs. |
| Secondary user | Technical leads and open-source contributors managing multi-project workspaces and enforcing CI token budgets. |
| Repository | `https://github.com/hamziCodes/CTXbank` |

---

## 2. Surfaces in scope

| Surface | Archetype | Density | Canvas | Notes |
| :--- | :--- | :--- | :--- | :--- |
| **Local Web Dashboard (`ctx ui`)** | `dashboard` | `compact` | `--app-canvas` | Embedded Go HTTP server on `localhost:4242`. Interactive node-link dependency graphs, visual memory cards, drag & drop research ingestion. |
| **VS Code Extension Sidebar** | `console` | `dense` | `--app-canvas` | Activity bar container with Active Focus tree, Memory Bank file list, and Checkpoint timeline. |
| **VS Code Visual Webview Tab** | `dashboard` | `compact` | `--app-canvas` | Full-tab interactive node-link architecture graph inside the IDE. |

---

## 3. Stack

| Field | Value | How it was determined |
| :--- | :--- | :--- |
| Framework and version | Go 1.23 embedded server + Vanilla HTML5/ES6 | Native single-binary architecture with zero npm dependencies for `ctx ui`; VS Code Extension API for IDE extension. |
| Language | Go 1.23, Vanilla ES6, TypeScript | `go.mod`, `tsconfig.json` |
| Styling | Plain CSS with VERTEX design tokens | Direct consumption of `var(--token)` from `tokens/tokens.css`. |
| Token delivery path | `plain-css` | Direct import, zero preprocessors needed. |
| Motion library | CSS only (`cubic-bezier(0.16, 1, 0.3, 1)`) | Lightweight, zero runtime overhead. |
| Icon library | Custom clean SVGs (1.75px stroke) | Ported directly to match VERTEX visual density. |
| Font loading | Inter Tight with system fallback | `font-family: 'Inter Tight', system-ui, -apple-system, Segoe UI, Roboto, sans-serif;` |
| Theme mechanism | `html.dark` class + pre-paint script | Instant dark/light mode toggle with zero first-paint flash. |
| Target platforms | Desktop browsers (Edge, Chrome, Safari, Firefox) + VS Code / Cursor IDE | Desktop workstation focus. |
| Browser floor | Modern evergreen browsers (Chrome 110+, Safari 16+, Edge 110+) | Local development environment. |

### Substitutions

| Toolkit assumes | This project uses | Consequence |
| :--- | :--- | :--- |
| `motion/react` | CSS transitions | State transitions use `cubic-bezier(0.16, 1, 0.3, 1)` with `150ms-200ms` durations; zero JS bundle bloat. |
| `lucide-react` | Inline clean SVGs | Icons render at 1.75 stroke equivalent, preserving binary footprint `< 15MB`. |
| Tailwind | Plain CSS | Consumes `var(--token)` directly without requiring Node.js build steps. |

---

## 4. Content and locale

| Field | Value |
| :--- | :--- |
| Primary locale | `en-US` |
| Additional locales | None |
| Text direction | `ltr` |
| Number format | Tabular figures (`font-variant-numeric: tabular-nums`) for line counts, byte sizes, timestamps, and commit hashes. |
| Date format | `YYYY-MM-DD HH:MM:SS UTC` |

---

## 5. Motion and performance budget

| Field | Value |
| :--- | :--- |
| Motion budget | `reduced` (snappy, purposeful, no gratuitous physics) |
| Reason | High-efficiency developer tool; instantaneous response required. |
| Embedded bundle ceiling | `< 350 KB` total embedded frontend assets (HTML + CSS + JS) |
| Binary size limit | Static binary must remain strictly under `15.0 MB`. |

---

## 6. Accessibility commitment

| Field | Value |
| :--- | :--- |
| Conformance target | WCAG 2.1 AA |
| Contrast minimum | 4.5:1 on normal text, 3:1 on large text/icons. On dark canvas, accent text uses `--accent-text` (`#609abe`). |
| Minimum touch/click target | `36x36px` in compact dashboard, `28x28px` in dense console tree. |

---

## 7. Deviation register

No deviations. CTXbank strictly implements the 5 VERTEX standing rules:
1. **One accent, navy:** `--brand-primary` fills; `--accent-text` for text and focus rings on dark canvas.
2. **Solid objects with edges:** 1.5px solid border, hard zero-blur offset shadow (`box-shadow: 2px 2px 0 var(--border)`). No blurry drop shadows, no glassmorphism.
3. **Numbers are the interface:** Tabular numbers on all line budget gauges (`85 / 150 lines`), byte sizes, and checkpoints.
4. **Every container has four states:** Ready, empty, loading, error. Never an unhandled blank panel.
5. **Both themes, always:** Seamless light and dark mode support with instant pre-paint bootstrap.

---

## 8. Adoption status

| Surface | Tokens installed | Converted | Both themes verified | Gate passed | Notes |
| :--- | :--- | :--- | :--- | :--- | :--- |
| Track 2: Embedded Web UI (`ctx ui`) | Ready | In Progress | Planned | Pending | Embedding `tokens.css` into Go binary. |
| Track 3: VS Code Extension Sidebar | Ready | Planned | Planned | Pending | Adapting VS Code CSS tokens to VERTEX standards. |

---

## 9. Project-specific decisions

- **Visual Architecture Graph:** Interactive SVG node-link diagram rendering nodes as VERTEX solid cards (1.5px solid borders, hard 2px offset shadow, navy accent connectors).
- **Line Budget Gauges:** Token counts render with tabular figures and clear visual indicator bars (`green` when `<120 lines`, `amber` when `120-149 lines`, `red/error` at `150+ lines`).
- **Drag & Drop Research Ingestion:** Research note dropzone styled with solid dashed 1.5px border; on drop, reveals side-by-side unified diff with distinct green/red delta lines.
