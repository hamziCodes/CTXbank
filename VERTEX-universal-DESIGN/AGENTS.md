# VERTEX Universal Design — Agent Contract

**You are reading the design authority for this repository.** Every UI decision in this project is governed by this toolkit. Read this file completely before writing or changing any interface code.

This toolkit is product-agnostic on purpose. It carries the parts of VERTEX that must never vary — colour, type, depth, motion, voice, accessibility — and a documented procedure for adapting everything else to the project you are actually in. Following that procedure is what produces one coherent design language across every VERTEX product instead of ten cousins that merely rhyme.

---

## 1. Load order

Read in this order. Stop as soon as you have what the task needs.

| Step | File | Read it when |
| :--- | :--- | :--- |
| 1 | **This file** | Always. First. |
| 2 | [`project-profile.md`](./project-profile.md) | Always. This is the project's own adaptation record. |
| 3 | [`spec/SURFACES.md`](./spec/SURFACES.md) | Always. Tells you which archetype you are building. |
| 4 | [`spec/DESIGN.md`](./spec/DESIGN.md) | Any visual decision — colour, type, depth, motion, a component's states. |
| 5 | [`tokens/tokens.css`](./tokens/tokens.css) | Wiring the project up, or when you need a token's exact value. |
| 6 | [`code/react/primitives.tsx`](./code/react/primitives.tsx) | Building React UI. Copy from it rather than reinventing. |
| 7 | [`spec/PRE-SHIP.md`](./spec/PRE-SHIP.md) | Before you report any UI work as finished. Mandatory. |

Task-shaped prompts live in [`prompts/`](./prompts/): [`new-surface.md`](./prompts/new-surface.md), [`restyle-existing.md`](./prompts/restyle-existing.md), [`audit-existing.md`](./prompts/audit-existing.md).

**Precedence when sources disagree:** the project's live CSS wins over `tokens/tokens.css`, which wins over `tokens/tokens.ts` and `tokens/tokens.json`, which win over the prose in `spec/`. Prose is documentation; if it contradicts a token file, the prose is the bug — fix it in the same change.

---

## 2. First-run procedure

Run this the first time you touch UI in a repository that contains this toolkit. It takes one pass and every later session reads the result instead of re-deriving it.

### Step 1 — Detect the stack

Inspect the repository. Do not ask the user what you can read yourself.

- **Framework** — `package.json`, `next.config.*`, `vite.config.*`, `app/` vs `pages/`, or a non-web target
- **Styling** — Tailwind v4 (`@import "tailwindcss"` and `@theme`), Tailwind v3 (`tailwind.config.*`), CSS Modules, styled-components, vanilla CSS
- **Motion** — is `motion` / `framer-motion` present?
- **Icons** — is `lucide-react` present?
- **Theme** — does a light/dark mechanism already exist? A `class="dark"` strategy, `next-themes`, a media query?

### Step 2 — Fill in the project profile

Open `project-profile.md` and complete every `[ ]` field from what you found. Commit it. This file is the contract for the project: every later session, human or agent, reads it instead of guessing, which is what keeps two agents six months apart from making two different sets of reasonable-but-divergent choices.

If a field genuinely cannot be determined from the repository — the currency, the primary locale, whether an existing product accent must be honoured — ask the user those specific questions. Ask once, record the answers, and never ask again.

### Step 3 — Install the tokens

Pick exactly one delivery path from the table below and follow it. Never mix two.

| Project styling | Do this |
| :--- | :--- |
| Tailwind v4 | Import `tokens/tokens.css`, then `tokens/tailwind-v4.css`, in that order, after `@import "tailwindcss"`. |
| Tailwind v3 | Import `tokens/tokens.css` in the global stylesheet. Merge `tokens/tailwind-v3.cjs` into `tailwind.config.js` via `presets: [require("./path/to/tailwind-v3.cjs")]`. |
| Plain CSS, CSS Modules, or CSS-in-JS | Import `tokens/tokens.css` only. Reference `var(--token)` directly. Skip both Tailwind adapters. |
| Non-web (native, email, slides, Figma) | Read `tokens/tokens.json`. Do not hand-transcribe hex values. |

Then add the no-flash theme bootstrap from `tokens/theme-init.ts`. A dark-mode flash on first paint is a shipping defect, not a rough edge.

### Step 4 — Choose the surface archetype

Read `spec/SURFACES.md` and pick the archetype for what you are building: `marketing`, `dashboard`, `mobile-app`, `form-flow`, or `console`. The archetype fixes the canvas, radius set, density, navigation pattern, motion budget and permitted components. Record it in the profile.

A project may contain several. A marketing page and its customer dashboard are different archetypes in one repository, and that is expected — they share tokens and diverge in composition.

### Step 5 — Build

Compose from `code/react/primitives.tsx` where the stack allows. Where it does not, port the same states and semantics rather than approximating them. Bind every colour, radius, shadow and duration to a token.

### Step 6 — Gate

Walk `spec/PRE-SHIP.md` before reporting the work finished. Report which items you verified and how. "Looks right" is not verification.

---

## 3. Standing rules

These apply to every task in every VERTEX repository, with no exceptions and no per-project overrides.

### Never

1. **Never hardcode a colour, radius, shadow or duration.** If you are typing a hex, a `px` radius or a `cubic-bezier`, you are meant to be typing a token name. The one exception is adding a new token to `tokens/tokens.css` itself.
2. **Never introduce a second accent colour.** Navy is the only accent. Not cobalt `#2B5EF5`, not orange, not purple, not a per-feature hue.
3. **Never ship one theme.** Light and dark are both required before any surface is considered done.
4. **Never use `--brand-primary` as text on a dark canvas.** It measures 2.2:1 and fails. Accent text and focus rings use `--accent-text`.
5. **Never leave a container in one state.** Every list, table and panel declares ready, empty, loading and error before it ships.
6. **Never invent facts.** No placeholder reviews, client logos, SLAs, metrics, legal entities or addresses. If real copy does not exist, ask for it or leave the slot visibly unfilled.
7. **Never mix easing curves.** One curve system-wide: `cubic-bezier(0.16, 1, 0.3, 1)`.
8. **Never remove a focus indicator** without replacing it with a visible one.

### Always

1. **Always read `project-profile.md` first.** It tells you what this project already decided.
2. **Always extend rather than fork.** A missing token is a token to add in both modes, mirror into `tokens.ts` and `tokens.json`, and document. It is not a reason to write a one-off hex.
3. **Always record a deviation.** If the project genuinely needs to depart from this spec, write it in the profile's deviation register with a reason and a review date. An undocumented deviation is drift; a documented one is a decision.
4. **Always prefer an existing project component** over a new one, and an existing primitive over a new component.
5. **Always state what you verified** when reporting UI work, naming the checklist items you actually checked.

---

## 4. What adapts and what does not

This is the heart of the toolkit. Get this distinction right and everything else follows.

### Fixed — identical in every VERTEX product

| Dimension | The invariant |
| :--- | :--- |
| Accent | One navy family. `--brand-primary` fills, `--accent-text` for text and focus. |
| Canvas and ink | White or paper on light, charcoal `#161719` on dark. Never pure black. |
| Semantic colour | `--positive` / `--negative` / `--attention` / `--info`, always paired with a non-colour signal. |
| Typeface | Inter Tight at 900 / 700 / 500 / 300. Monospace only for code, logs, IDs and diffs. |
| Depth | Solid surfaces, crisp outlines, hard zero-blur offset shadows. |
| Easing | `cubic-bezier(0.16, 1, 0.3, 1)`. |
| Voice | Sentence case, no Oxford comma, numerals, contextual CTAs, the banned-word list. |
| Accessibility | WCAG AA, visible focus, 44px touch targets, colour never the sole signal. |
| Banned patterns | Gradients, glassmorphism as a material, blurred shadows, 3D blobs, sparkle icons, template SaaS chrome. |

### Adapts — per project, per archetype, recorded in the profile

| Dimension | Range | Driven by |
| :--- | :--- | :--- |
| Canvas choice | `--background` or `--app-canvas` | Archetype |
| Radius set | Marketing up to 60px · application 12–28px | Archetype |
| Density | `comfortable` · `compact` · `dense` | Archetype and data density |
| Type scale | Display-led or UI-led | Archetype |
| Navigation | Top bar · bottom nav · sidebar · stepper | Archetype and platform |
| Motion budget | Full · reduced · minimal | Archetype and performance target |
| Layout width | 1440px page · 740px column · fluid console | Archetype |
| Locale and currency | Any; `PKR` is the default, not a constraint | Profile |
| Component set | Only what the archetype permits | Archetype |

The rule in one line: **tokens never change, composition always does.**

---

## 5. Handling a project that already has a design

Most adoptions are retrofits, not greenfield. Do not rewrite the whole interface in one pass.

1. **Audit first.** Use [`prompts/audit-existing.md`](./prompts/audit-existing.md) to produce a written inventory of violations before changing anything.
2. **Land the tokens first.** Install the token layer and map the project's existing colour variables onto VERTEX tokens. This alone usually resolves most of the drift.
3. **Convert by surface, not by file.** One screen at a time, each landing in a coherent state. Never leave the app half-converted across a release.
4. **Keep a deviation register.** Anything you cannot convert yet goes in the profile with a reason, so the next session inherits the reasoning instead of rediscovering the problem.
5. **Report honestly.** Say what converted, what did not, and what it would take. A partial conversion described accurately is far more useful than a claim of completeness.

---

## 6. When you are unsure

In priority order:

1. **Read.** `spec/DESIGN.md` covers almost every question about colour, type, depth, motion, component states and voice.
2. **Copy the nearest sibling.** If VERTEX already solved a similar surface, match it rather than inventing a second answer.
3. **Choose the more conservative option.** Fewer effects, fewer colours, more contrast, more whitespace. The VERTEX failure mode is never "too plain."
4. **Ask a specific question.** One decision, two or three concrete options, your recommendation first. Do not ask the user to design for you.
5. **Never guess silently.** An undocumented guess becomes the next project's precedent.
