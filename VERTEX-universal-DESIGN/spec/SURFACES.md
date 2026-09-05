# Surface Archetypes

**This is the adaptation layer.** [`DESIGN.md`](./DESIGN.md) defines what never changes. This file defines the five legitimate ways to compose it.

The principle: **tokens never change, composition always does.** A marketing page and an admin console draw from the same colour, type, depth and motion system, then compose it differently because they serve different tasks. That is why VERTEX products look related without looking identical, and why a new product can be built in a week without a design review.

Pick the archetype first. It settles the canvas, radius range, density, navigation, motion budget and permitted components before you write a line of layout — which is exactly the set of decisions where projects otherwise drift apart.

---

## Choosing

```
Is it public-facing, selling or explaining?
├─ YES → marketing
└─ NO  → It is an authenticated tool.
         │
         Is the main job one linear task with inputs?
         ├─ YES → form-flow
         └─ NO
              │
              Is a phone the primary device?
              ├─ YES → mobile-app
              └─ NO
                   │
                   Is it operator tooling — logs, queues, bulk edit, tables over 20 rows?
                   ├─ YES → console
                   └─ NO  → dashboard
```

One repository commonly holds several. A marketing site, its customer dashboard and a staff phone app are three archetypes sharing one token layer. That is the intended shape, not a compromise.

---

## At a glance

| | `marketing` | `dashboard` | `mobile-app` | `form-flow` | `console` |
| :--- | :--- | :--- | :--- | :--- | :--- |
| Canvas | `--background` | `--app-canvas` | `--app-canvas` | `--app-canvas` | `--app-canvas` |
| Density | `comfortable` | `compact` | `comfortable` | `comfortable` | `dense` |
| Radius ceiling | `--radius-marketing` (60px) | `--radius-panel` (24px) | `--radius-sheet` (28px) | `--radius-panel` (24px) | `--radius-slab` (18px) |
| Width | `--page-max` 1440px | fluid, sidebar + content | `--app-max` 740px | `--prose-max` 720px | fluid, full bleed |
| Navigation | Fixed top bar | Sidebar | Bottom nav | Stepper | Top bar + sidebar tree |
| Type lead | Display (900 at `--fs-hero`+) | UI (700 at 22px) | UI (700 at 22–25px) | UI (700 at 18–22px) | UI (700 at 15–18px) |
| Motion budget | Full | Reduced | Reduced | Reduced | Minimal |
| Shadows | None (scale + overlay) | `--shadow-slab` | `--shadow-slab` | `--shadow-slab` | None (borders only) |
| Sections per viewport | 1 | 3–6 blocks | 3–5 blocks | 1 step | As many as fit |

---

## 1. `marketing`

Public pages: home, product, pricing, case study, about, campaign landing.

**The job:** route a visitor to the right next step in under two clicks. Everything else is secondary.

- **Canvas** `--background` — white in light, charcoal `#161719` in dark. Not the paper canvas.
- **Density** `comfortable`. Never compact; a dense marketing page reads as an admin panel.
- **Radius** `rounded-4xl` (32px) on mobile, `--radius-marketing` from `lg`. This is the only archetype permitted above 28px.
- **Layout** `--page-max` (1440px), padding `--page-pad`, deep insets `--page-inset` from `lg`. Asymmetric and editorial. Never a centred 3-up card grid.
- **Type** display-led. Hero anchors at `--fs-hero` in 900. Section lockups at `--fs-title`. Body at 15→17px in 500, capped at `46ch`.
- **Navigation** fixed top bar at `--nav-height`, transparent at rest (`bg-background/40`), solid on scroll (`bg-background/80` with `border-border`). Full-height mobile panel with the wordmark pinned at the bottom.
- **Depth** no offset shadows. Depth comes from a scroll-linked `scale: 1 → 0.94` plus a black overlay ramping `0 → 0.42` as the next slab arrives.
- **Motion** full budget. One animated element per viewport. Sticky offsets step down (`top-24`, `top-28`, `top-32`, `top-36`) so every card edge stays visible.
- **Sections** one slab per viewport. Big, confident, singular.

**Signature moves:** the full-bleed wordmark lockup · sticky stacking slabs alternating brand navy and `--surface-dark` · parenthetical kickers `(HOW WE WORK)` · hover-reveal links that translate a clipped two-line column · the looping scroll cue.

**Never here:** stat grids, data tables, dense list rows, sidebars, bottom navs, offset shadows, `compact` density.

---

## 2. `dashboard`

Authenticated overview for an operator on a desktop or tablet: clinic day view, sales summary, project status.

**The job:** answer "what needs my attention right now" above the fold, then let the user drill in.

- **Canvas** `--app-canvas` (paper), so `--app-surface` slabs read as objects sitting on it.
- **Density** `compact`. This is a working tool, not a brochure.
- **Radius** `--radius-slab` (18px) for cards and rows, `--radius-panel` (24px) for the primary metric slab, `--radius-control` (12px) for controls. Never above 24px.
- **Layout** `--sidebar-width` (264px) plus fluid content. Content column capped near 1280px so tables do not stretch to absurd line lengths on ultrawide displays.
- **Type** UI-led. Section headings 22px in 700. Metrics 30→40px in 700 at `-0.02em`. Body at `var(--text-body)`.
- **Navigation** persistent left sidebar, collapsible to `--sidebar-width-collapsed` (72px) with icons and tooltips. Active item takes `--accent-text` and a 2px leading edge.
- **Depth** `--shadow-slab` on every resting slab. Brand-filled slabs carry none.
- **Motion** reduced. Reveals on first mount only, never on filter or tab changes — re-animating a list every time a filter changes is actively annoying. No scroll-linked transforms.

**Composition order** — the same on every VERTEX dashboard, so an operator moving between products already knows where to look:

```
1. Eyebrow (date or scope) + greeting or page title
2. Context switcher (segmented) if the screen has modes
3. ONE primary metric slab, tone="brand"
4. Stat grid, 2-up on tablet, 3 or 4-up on desktop
5. Attention slab — only when there is something to act on
6. List or table preview with a "View all" text link
```

Two brand-navy slabs never stack adjacently. The second takes `tone="inverse"`.

**Never here:** display type above `--fs-title`, marketing radii, sticky stacking slabs, full-bleed wordmarks, decorative animation.

---

## 3. `mobile-app`

Phone-first authenticated app, native or web: staff tools, field entry, customer app.

**The job:** one thumb, one hand, often poor light and worse signal.

- **Canvas** `--app-canvas`.
- **Density** `comfortable`. Touch demands it — this is the density decision that is not negotiable.
- **Radius** up to `--radius-sheet` (28px) for bottom sheets, `--radius-panel` (24px) for hero slabs, `--radius-slab` (18px) for rows.
- **Layout** single column at `--app-max` (740px), `px-5`, and `pb-[86px]` so the last row clears the bottom nav.
- **Type** UI-led. Screen title 22–25px in 700. Body 15px in 500. Never below 12px.
- **Navigation** bottom nav at `--bottom-nav-height` (70px), 3–5 tabs, becoming a centred 600px bar at `≥800px`. A sixth tab means the information architecture is wrong; an overflow menu is not the fix.
- **Depth** `--shadow-slab` on slabs, `--shadow-slab-strong` on the FAB. The sticker press (`translate(2px, 2px)`) matters most here — on touch there is no hover, so the press is the only feedback the user gets.
- **Motion** reduced. Bottom sheets slide `y: 100% → 0%`. No scroll-linked transforms; they stutter on mid-range Android.

**Required patterns:** every input at least 16px or iOS Safari zooms on focus · every target at least 44×44px · bottom sheets rather than centred modals · destructive actions confirmed · a visible offline state when the product works offline.

**Never here:** hover-only affordances, sidebars, dense tables (use cards), multi-column forms.

---

## 4. `form-flow`

One linear task: contact wizard, onboarding, checkout, booking, KYC.

**The job:** finish. Every decision serves completion rate.

- **Canvas** `--app-canvas`, or `--background` when embedded in a marketing page.
- **Density** `comfortable`. Never compress a form; every pixel saved costs completion.
- **Radius** `--radius-panel` (24px) on the panel, `--radius-control` (12px) on fields.
- **Layout** single column at `--prose-max` (720px). Panel padding `p-5` → `md:p-8` → `lg:p-10`. Never two columns of inputs — it doubles the eye path and breaks tab order expectations.
- **Type** UI-led. Step legend 18–20px in 700. Labels 13px in 700 at `0.04em`. Inputs at 16px minimum.
- **Navigation** a stepper. Show `01 / 04` at 13px/700, `0.08em`, in `--accent-text`. Back is always available; forward validates.
- **Motion** reduced. Step transitions use `AnimatePresence mode="wait"`, in `y: 16`, out `y: -12`, `--dur-panel`.

**Required patterns:** validate per step, never only at the end · block forward movement, never silently discard input · the error message names the fix ("Add your name and an email or phone number"), carries `role="alert"`, and focus moves to the first invalid field · a review step before anything irreversible · state survives a back-navigation · one primary action per step, bottom-right.

**Never here:** more than one primary action per step, hidden required fields, a submit button that does not say what it will do.

---

## 5. `console`

Internal operator tooling: log viewers, queues, bulk edit, audit trails, admin tables.

**The job:** scan and act on a lot of rows fast. Information density is the feature.

- **Canvas** `--app-canvas`.
- **Density** `dense`. On a coarse pointer this clamps back to a 44px control height automatically, which is what keeps it usable on a tablet.
- **Radius** `--radius-slab` (18px) ceiling, `--radius-control` (12px) on controls. Table rows are square-edged.
- **Layout** full-bleed fluid. Sticky table headers. Virtualised rows above 100 items.
- **Type** UI-led and small: headings 15–18px in 700, body 13px in 500. **This is the one archetype where `--font-mono` is expected** — for IDs, hashes, log lines, diffs, JSON and stack traces. Never for chrome, labels or buttons.
- **Navigation** top bar plus a sidebar tree. Breadcrumbs when nesting goes past two levels.
- **Depth** borders only. No offset shadows — at this density they turn into visual noise.
- **Motion** minimal. Opacity-only transitions. No reveals, no layout animation, no scroll-linked transforms. An animated table is an unreadable table.

**Required patterns:** every column sortable or explicitly not · filters reflected in the URL so a view can be shared · bulk selection shows a count and requires confirmation · every destructive action is undoable or confirmed with the record name typed out · timestamps carry a timezone · long values truncate with the full value available on demand, never silently cut.

**Never here:** decorative animation, marketing radii, display type, offset shadows, hover-only affordances on rows that also need keyboard access.

---

## Composing across archetypes

**A marketing page that embeds a form.** The section stays `marketing`; the form panel inside it switches to `form-flow` rules. Set `data-density="comfortable"` on the panel and use application radii inside it. The seam should be invisible.

**A dashboard with a dense table.** Keep the page `dashboard` / `compact` and set `data-density="dense"` on the table container only. Density nests; this is the mechanism working as designed.

**A mobile view of a dashboard.** Below `800px`, a `dashboard` becomes a `mobile-app`: the sidebar becomes a bottom nav, tables become cards, density goes to `comfortable`. This is a real archetype switch at a breakpoint, not a squeeze — plan for it rather than discovering it.

**An email or a PDF.** Neither is an archetype. Read `tokens/tokens.json`, use table layout and inline styles, assume no custom font, and keep to light mode.

---

## Adding an archetype

Do not, unless a real surface genuinely fits none of the five. Six archetypes is how a design system starts becoming a suggestion.

If you must: define all nine columns of the "at a glance" table, list its signature moves and its bans, add it to `tokens.css` section 5 and `tokens.ts` `archetypeCanvas`, then record the addition in `CHANGELOG.md` with a note about what existing surface forced it.
