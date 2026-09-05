# Changelog

Version history for VERTEX Universal Design. Every product records the version it uses in its own `project-profile.md`, so upgrades can be diffed rather than guessed at.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Versioning is semantic, with a design-system reading of each level:

- **Major** — a token is removed or changes meaning, or a rule changes such that existing compliant UI becomes non-compliant. Requires a migration pass.
- **Minor** — tokens, components, archetypes or rules are added. Existing UI stays compliant.
- **Patch** — documentation corrections, clarifications, measured-value fixes. No visual change.

---

## [1.0.0] — 2026-08-24

First universal release. Generalises the single-product design system into a portable toolkit that any VERTEX app, dashboard, website or internal tool can adopt.

### Added — the adaptation layer

- **`AGENTS.md`** — the agent contract: load order, a six-step first-run procedure, standing never/always rules, an explicit table of what is fixed versus what adapts, and guidance for retrofitting a project that already has a design.
- **`project-profile.md`** — per-project template covering stack detection, surfaces in scope, substitutions, locale and currency, motion and performance budgets, accessibility commitment, adoption status, and a **deviation register** that turns exceptions into recorded decisions rather than silent drift.
- **`spec/SURFACES.md`** — five surface archetypes (`marketing`, `dashboard`, `mobile-app`, `form-flow`, `console`), each fixing canvas, density, radius ceiling, width, navigation, type lead, motion budget and permitted components. This is the mechanism that lets one token layer serve a landing page and an admin console.
- **`spec/PRE-SHIP.md`** — the gate: blockers that cannot ship, required items that must be fixed or registered, release-level checks, and a report template that asks *how* each item was verified.
- **`rules/vertex-design.mdc`** — Cursor project rule with `alwaysApply: true`.
- **`rules/AGENTS-snippet.md`** — wiring for any agent, plus a single question that verifies the link actually works.
- **`prompts/`** — task prompts for building a new surface, auditing an existing codebase, and converting one in stages, each explaining why it is shaped the way it is.
- **`prompts/standalone-system-prompt.md`** — a self-contained prompt for tools that cannot read the repository (v0, Lovable, Figma Make, ChatGPT).
- **`CHANGELOG.md`** — this file.

### Added — tokens

- **Density scale.** `comfortable` / `compact` / `dense` as `data-density`, exposing `--control-h`, `--control-pad-x`, `--slab-pad`, `--stack-gap`, `--row-pad-y`, `--field-gap`, `--icon-tile` and `--text-body`. Nests freely, so a dense table inside a comfortable page is valid.
- **Touch floor.** On `@media (pointer: coarse)`, `compact` and `dense` clamp `--control-h` back to 44px and `--row-pad-y` to 12px. A dense console stays usable on a tablet.
- **Archetype canvases.** `[data-vertex-surface]` selects the marketing or application canvas.
- **`--product-signal`** — the one scoped escape hatch. Defaults to `var(--brand-primary)`; permitted only on the product identity chip and nav mark, under ~40×40px; requires a deviation register entry.
- **`--font-mono`** — a system monospace stack, permitted only for code, logs, IDs and diffs. Previously there was no legitimate way to render a log line, which meant every console surface invented one.
- **Chart series order** — fixed light and dark sequences in `tokens.ts` and `tokens.json`, so the same data reads the same way in every product.
- **Forced-colors support** — Windows high-contrast mode keeps borders and drops shadows.
- **`tokens/tokens.json`** — platform-neutral mirror for Figma, native, email and slides, so nobody hand-transcribes hex values.
- **`tokens/theme-init.ts`** — pre-paint theme bootstrap plus density and surface helpers, dependency-free.
- **`tokens/tailwind-v3.cjs`** — preset for projects still on Tailwind v3.

### Changed

- **Restructured** from six flat files into `spec/`, `tokens/`, `code/`, `rules/` and `prompts/`, with `AGENTS.md` and `project-profile.md` as the entry points.
- **`tokens.css` is now framework-neutral.** The Tailwind v4 `@theme inline` block and `@custom-variant dark` moved into `tokens/tailwind-v4.css`, so plain-CSS, CSS Modules and CSS-in-JS projects can consume the tokens without Tailwind. Four documented delivery paths, and a rule never to mix two.
- **`DESIGN.md` is now surface-agnostic.** The one-pager section catalogue and website-specific content contracts were removed in favour of the archetype rules in `SURFACES.md`. Added a data table specification, chart rules, a density section, a monospace policy, and a colour-token extension procedure.
- **`primitives.tsx` is now density-aware.** Components read `var(--control-h)`, `var(--control-pad-x)`, `var(--slab-pad)`, `var(--stack-gap)`, `var(--row-pad-y)` and `var(--text-body)` instead of hardcoding sizes, which is what lets one component library serve every archetype without forking.
- **Money formatting is now locale-parameterised.** `formatMoney()` takes `currency`, `precision` and `locale`; `PKR` and `en-PK` are defaults rather than constraints.
- **`agent-prompt.md` was split** into `rules/vertex-design.mdc` (repo-linked), `rules/AGENTS-snippet.md` (wiring) and `prompts/standalone-system-prompt.md` (self-contained).

### Added — components

- **`Surface`** and **`DensityScope`** — declare archetype and density for a subtree.
- **`DataTable`** — sticky header, `aria-sort` on sortable headers, right-aligned tabular numeric columns, hover on the row, no zebra striping, loading and empty states built in.
- **`SidebarNav`** — collapsible dashboard navigation with a 2px active leading edge and tooltips when collapsed.
- **`ThemeToggle`** and **`useVertexTheme`** — read the theme from the DOM class so state can never disagree with what is painted.
- **`Sheet`** now traps focus, restores focus to the trigger on close, closes on `Escape`, and locks body scroll.
- **`Reveal`** takes a `disabled` prop for projects on the `minimal` motion budget.

### Fixed

- **Accent text contrast.** `--brand-primary` in dark mode (`#1a537c`) measures 2.2:1 against the `#161719` canvas and fails WCAG even for large text. `--accent-text` (`#609abe`, 5.8:1) and `--accent-text-hover` (`#90bbd5`, 8.8:1) are now the tokens for text and focus rings in both modes, and the split is documented everywhere it can be misread.
- **The `--muted` trap.** In light mode `--muted` is literally `#000000`, so light-mode de-emphasis works through alpha rather than through a grey. Documented, with `--app-ink-soft` and `--app-ink-muted` as the tokens that behave correctly in both modes.
- **`--radius-sheet`** documented consistently as a single `28px` value applied via `rounded-t-sheet`, rather than as a compound corner value.

### Notes for adopters

- Nothing here is breaking for the original design-system folder; this is a superset with a new entry point. Existing token names are unchanged.
- If you are adopting into a project that already has UI, run `prompts/audit-existing.md` **before** converting anything. Knowing the scope is worth one session.
- Record the toolkit version in your `project-profile.md` so the next upgrade is a diff rather than an archaeology exercise.
