# VERTEX Universal Design

**The default UI/UX toolkit for every VERTEX product.** Drop this folder into any app, dashboard, website or internal tool, link it to your agent, and you get one coherent design language across everything VERTEX ships.

It is deliberately not a component library you install and forget. It is a **token layer plus an adaptation procedure**: the parts of VERTEX that must never vary, and a documented method for fitting them to the project actually in front of you. That is the difference between ten products that look like one company and ten products that merely rhyme.

> **New here?** Read [`AGENTS.md`](./AGENTS.md). It is the contract, and it is short.

---

## Install in 3 steps

### 1. Copy the folder in

```bash
cp -r VERTEX-universal-DESIGN /path/to/your-project/
```

Or add it as a git submodule so every product tracks the same version:

```bash
git submodule add <toolkit-repo-url> VERTEX-universal-DESIGN
```

### 2. Link it to your agent

```bash
mkdir -p .cursor/rules
cp VERTEX-universal-DESIGN/rules/vertex-design.mdc .cursor/rules/
```

That rule carries `alwaysApply: true`, so every agent request in the repository is governed by the system without anyone remembering to mention it. For non-Cursor agents, or to also wire a root `AGENTS.md`, see [`rules/AGENTS-snippet.md`](./rules/AGENTS-snippet.md).

### 3. Tell the agent to run the first-run procedure

> Read `VERTEX-universal-DESIGN/AGENTS.md` and run the first-run procedure on this repository.

The agent detects your stack, fills in [`project-profile.md`](./project-profile.md), installs the tokens on the one path that matches your setup, and wires the no-flash theme bootstrap. One pass, then every later session reads the profile instead of guessing.

**Verify the link worked.** In a fresh session, without naming any file, ask: *"What accent colour does this project use for text on a dark background, and why is it not the brand colour?"* A correctly linked agent answers `--accent-text` / `#609abe` and explains that `--brand-primary` fails at 2.2:1 on the dark canvas. If it guesses or asks which design system you mean, the link is not working.

---

## What's in here

```
VERTEX-universal-DESIGN/
├── AGENTS.md                  ← the contract. Load order, first-run procedure,
│                                standing rules, what adapts and what does not
├── project-profile.md         ← per-project template: stack, surfaces, locale,
│                                budgets, deviation register, adoption status
├── CHANGELOG.md               ← version history, so products can diff and upgrade
│
├── spec/
│   ├── DESIGN.md              ← the universal core. Colour, type, shape, depth,
│   │                            density, motion, components, voice, accessibility
│   ├── SURFACES.md            ← the adaptation layer. Five archetypes, each fixing
│   │                            canvas, radius, density, nav, motion, components
│   └── PRE-SHIP.md            ← the gate. Blockers, required items, report template
│
├── tokens/
│   ├── tokens.css             ← THE SOURCE OF TRUTH. Framework-neutral CSS
│   │                            variables, density scale, archetype canvases,
│   │                            base styles, utilities, reduced motion
│   ├── tailwind-v4.css        ← @theme inline adapter for Tailwind v4
│   ├── tailwind-v3.cjs        ← preset for Tailwind v3
│   ├── tokens.ts              ← typed mirror for charts, canvas, native, scripts
│   ├── tokens.json            ← platform-neutral mirror for Figma, native, email
│   └── theme-init.ts          ← pre-paint theme bootstrap, density and surface helpers
│
├── code/react/
│   └── primitives.tsx         ← every component in DESIGN.md §7, density-aware
│
├── rules/
│   ├── vertex-design.mdc      ← copy into .cursor/rules/
│   └── AGENTS-snippet.md      ← wiring for any agent, plus how to verify it
│
└── prompts/
    ├── new-surface.md         ← build a new screen or component
    ├── restyle-existing.md    ← bring existing UI onto the system, in stages
    ├── audit-existing.md      ← measure how far off the system a codebase is
    └── standalone-system-prompt.md
                               ← self-contained prompt for tools that cannot read
                                 the repo: v0, Lovable, Figma Make, ChatGPT
```

---

## How adaptation works

The whole design of this toolkit rests on one line:

> **Tokens never change. Composition always does.**

**Fixed in every VERTEX product:** one navy accent, the canvas and ink pairs, semantic colour, Inter Tight at four weights, solid surfaces with hard offset shadows, `cubic-bezier(0.16, 1, 0.3, 1)`, the voice rules, WCAG AA, and the banned-pattern list.

**Adapted per project:** which canvas, which radius ceiling, which density, which navigation, which motion budget, which layout width, which locale and currency, and which components are even permitted.

That second list is not freeform. It is driven by the **surface archetype** — `marketing`, `dashboard`, `mobile-app`, `form-flow` or `console` — and each archetype fixes all of it. Pick the archetype first and the rest follows, which is exactly why an agent can make correct decisions without a design review.

One repository commonly holds several archetypes. A marketing site, its customer dashboard and a staff phone app share one token layer and diverge in composition. That is the intended shape.

---

## Working with it day to day

| Task | Start with |
| :--- | :--- |
| New screen, page or component | [`prompts/new-surface.md`](./prompts/new-surface.md) |
| Existing project onto the system | [`prompts/audit-existing.md`](./prompts/audit-existing.md), then [`prompts/restyle-existing.md`](./prompts/restyle-existing.md) |
| A tool that cannot read the repo | [`prompts/standalone-system-prompt.md`](./prompts/standalone-system-prompt.md) |
| A question about a specific value | [`tokens/tokens.css`](./tokens/tokens.css), then [`spec/DESIGN.md`](./spec/DESIGN.md) |
| Choosing a layout | [`spec/SURFACES.md`](./spec/SURFACES.md) |
| Finishing anything | [`spec/PRE-SHIP.md`](./spec/PRE-SHIP.md) |

**Retrofits convert one surface at a time.** Audit first, install the tokens, then convert surface by surface with each landing in a coherent state. A half-converted app looks more broken than a consistently old one.

---

## Precedence

When two sources disagree, in this order:

1. The **project's live CSS** — what actually ships
2. **`tokens/tokens.css`** — the canonical token definitions
3. **`tokens.ts`** and **`tokens.json`** — the mirrors
4. The prose in **`spec/`** — documentation

If prose contradicts a token file, the prose is the bug. Fix it in the same change.

A new token goes in all four places — `tokens.css` in both modes, `tokens.ts`, `tokens.json`, and the relevant table in `DESIGN.md` — plus a `CHANGELOG.md` entry. A token in fewer than four places will resurface as a hardcoded hex in someone's component within a month.

---

## The five rules, if you remember nothing else

1. **One accent, navy.** `--brand-primary` fills. `--accent-text` for text and focus rings, because the brand navy fails contrast on the dark canvas.
2. **Solid objects with edges.** 1.5px outline, hard zero-blur offset shadow. No gradients, no glass.
3. **Numbers are the interface.** Tabular figures, one money formatter, every delta states its period.
4. **Every container has four states.** Ready, empty, loading, error. Never a blank panel.
5. **Both themes, always.** Light and dark are one feature, not two.

---

## Deviations

Real projects need exceptions. The system handles that — just not silently.

Record every departure in the deviation register in [`project-profile.md`](./project-profile.md), with a reason, a scope and a review date. A documented deviation is a decision the next session can trust or revisit; an undocumented one is drift.

Some things are never deviations, only violations: a second accent, a single-theme surface, `--brand-primary` as text on dark, a removed focus indicator, a missing empty/loading/error state, invented facts, or a banned word.

---

**Version 1.0.0** · See [`CHANGELOG.md`](./CHANGELOG.md) for history and upgrade notes.
