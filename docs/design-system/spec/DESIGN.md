# VERTEX Universal Design — Core Specification

**What never changes, in any VERTEX product.** Colour, type, depth, density, motion, voice and accessibility. How to *compose* these into a given surface is in [`SURFACES.md`](./SURFACES.md); before shipping, walk [`PRE-SHIP.md`](./PRE-SHIP.md).

Canonical values live in [`../tokens/tokens.css`](../tokens/tokens.css). If this document and that file disagree, **the CSS wins** and this document is the bug — fix it in the same change.

---

## 1. Identity

### 1.1 The philosophy: The Engineered Slab

VERTEX rejects the cold, homogenous aesthetic of modern AI-powered SaaS. Not translucent glass, not neon glow, not abstract 3D shapes. VERTEX is **precision-machined navy slabs on clean paper**.

It should read like engineering documentation rather than a marketing deck: confident type, exactly one accent colour, solid surfaces with crisp hairline outlines, hard offset shadows that behave like real objects, and numbers that line up in columns.

One sentence to settle arguments: **a VERTEX surface is an object with an edge and a weight, not a pane of frosted glass.**

```
+------------------------------------------+
|  SLAB (--app-surface)                    |
|  1.5px solid --app-line-strong           |
|  radius --radius-slab                    |
+------------------------------------------+
  \________________________________________\   3px 4px 0 --app-shadow-color
                                                zero blur, always
```

The system is deliberately narrow. One accent, one typeface, one easing curve, one shadow shape. Narrowness is what lets ten products look like one company, and what lets an agent make a correct decision without a design review.

### 1.2 Hard constraints

Prohibited everywhere, in every archetype, with no per-project override.

| Category | Banned | Use instead |
| :--- | :--- | :--- |
| **Accent** | A second accent of any kind. Cobalt `#2B5EF5`, orange, purple, per-feature hues. | One navy family. `--brand-primary` fills, `--accent-text` text and focus. |
| **Gradients** | Purple-to-blue, cyan-to-violet, mesh blobs, neon washes, gradient text. | Solid fills only. |
| **Materials** | Glassmorphism as a card material. Frosted panels, translucent cards, `backdrop-filter` on content. | Solid surfaces, 1.5px outline, hard offset shadow. **One exception:** fixed navigation chrome and sticky kickers may use `backdrop-blur` at ≥80% background opacity. Chrome only, never content. |
| **Shadows** | Blurred drop shadows, coloured glows, `text-shadow` on body copy, layered elevation systems. | `--shadow-slab` and its strong and pressed variants. `--shadow-sheet` is the only blurred shadow in the system. |
| **Icons** | 3D blobs, robot avatars, sparkle icons, gradient duotone sets, mixed icon families. | One family, 2px-equivalent stroke (`strokeWidth={1.75}`), `currentColor`, 16–22px. |
| **Type** | A second display family. Monospace in UI chrome. Letter-spaced body copy. Justified text. | Inter Tight at 900/700/500/300. Monospace only for code, logs, IDs and diffs. |
| **Contrast** | `#ffffff` on `#000000` in dark mode. `--brand-primary` as text on a dark canvas (2.2:1). Pure black canvases. | Light: `#000000` on `#ffffff`. Dark: `#f3f4f6` on `#161719`. Accent text on dark: `--accent-text`. |
| **Theming** | Shipping one mode. Neon dark mode. A dark-mode flash on load. | Both modes required. Charcoal `#161719`, never pure black. Pre-paint theme script. |
| **Chrome** | Pill badges as decoration. Generic 3-up bordered card grids. Template admin shells. Skeuomorphic textures. | Editorial and asymmetric layout. Badges only when they carry a status value. |
| **Motion** | Mixed easing curves. Scanlines, glitch, flicker, matrix rain. Parallax above 0.3. Multiple simultaneous animations per viewport. Re-animating on every state change. | One curve. Fade plus 16–32px travel. One animated element per viewport. `prefers-reduced-motion` honoured globally. |
| **Content** | Invented reviews, client logos, SLAs, metrics, legal entities, addresses. Lorem ipsum in a reviewable build. | Real copy, or a visibly unfilled slot. |
| **Copy** | `leverage`, `synergy`, `seamless`, `revolutionize`, `cutting-edge`, `state-of-the-art`, `game-changing`, `next-gen`, `world-class`. Title Case headings. Oxford commas. Generic CTAs. | Sentence case, numerals, contextual CTAs. See §8. |

---

## 2. Colour

Every colour is a CSS custom property declared twice — `:root` for light, `html.dark` for dark. Both modes are mandatory.

### 2.1 The accent split

This is the single most misunderstood rule in the system, so it comes first.

`--brand-primary` is a **fill** token. In dark mode it is `#1a537c`, which against the `#161719` canvas measures **2.2:1** — a failure even for large text. Accent text and focus rings therefore bind to a different token:

| Token | Light | Dark | Use |
| :--- | :--- | :--- | :--- |
| `--brand-primary` | `#042940` | `#1a537c` | **Fills only.** Buttons, filled slabs, selected states, the wordmark slash. |
| `--brand-hover` | `#0b3b5b` | `#236798` | Hover state of a filled surface. |
| `--accent-text` | `#042940` | `#609abe` | **Text and icons.** Links, kickers, active nav, inline emphasis. |
| `--accent-text-hover` | `#0b3b5b` | `#90bbd5` | Hover state of accent text. |
| `--focus-ring` | `#042940` | `#609abe` | `outline: 2px solid var(--focus-ring); outline-offset: 2px`. |

White on `--brand-primary` is 15.0:1 in light mode, so a filled navy button is always safe. Navy *text* on a dark canvas never is.

### 2.2 The navy ramp

One accent family, hue locked to 203–205°. The brand values are members of this ramp, not separate colours.

| Step | Hex | Role |
| :--- | :--- | :--- |
| `--navy-50` | `#eff6fb` | Info wells, hover tint on paper |
| `--navy-100` | `#dbebf5` | Selected rows, chart band fill |
| `--navy-200` | `#bad6e8` | Chart series 5, disabled accent, delta text on navy |
| `--navy-300` | `#90bbd5` | Accent text hover on dark (8.8:1), chart series 3 |
| `--navy-400` | `#609abe` | Accent text and focus ring on dark (5.8:1). Decorative on light (3.1:1). |
| `--navy-500` | `#236798` | `--brand-hover` dark, chart series 2 |
| `--navy-600` | `#1a537c` | `--brand-primary` dark — **fill only** |
| `--navy-700` | `#0b3b5b` | `--brand-hover` light |
| `--navy-800` | `#042940` | `--brand-primary` light — the signature VERTEX navy |
| `--navy-900` | `#021c2c` | Deep slab, chart axis on dark |
| `--navy-950` | `#01111b` | Maximum depth. Not a canvas — use `#161719`. |

### 2.3 Canvases

Two canvas systems. The archetype picks one; see `SURFACES.md`.

**Marketing canvas** — white or charcoal, no paper tint:

| Token | Light | Dark |
| :--- | :--- | :--- |
| `--background` | `#ffffff` | `#161719` |
| `--foreground` | `#000000` | `#f3f4f6` |
| `--card` | `#ffffff` | `#212326` |
| `--border` | `rgb(0 0 0 / 0.08)` | `#2e3238` |
| `--surface-dark` | `#000000` | `#212326` |
| `--surface-hover` | `rgb(0 0 0 / 0.04)` | `#2d3035` |

**Application canvas** — paper, so slabs read as objects:

| Token | Light | Dark |
| :--- | :--- | :--- |
| `--app-canvas` | `#f2f6f9` | `#161719` |
| `--app-surface` | `#ffffff` | `#212326` |
| `--app-surface-sunken` | `#e7eef4` | `#101113` |
| `--app-surface-inverse` | `#042940` | `#f3f4f6` |
| `--app-on-inverse` | `#ffffff` | `#161719` |
| `--app-line` | `#dbe4ec` | `#2e3238` |
| `--app-line-strong` | `#c2d0dd` | `#3d434b` |
| `--app-ink` | `#000000` | `#f3f4f6` |
| `--app-ink-soft` | `rgb(0 0 0 / 0.64)` | `rgb(243 244 246 / 0.72)` |
| `--app-ink-muted` | `rgb(0 0 0 / 0.52)` | `#9ca3af` |

> **The `--muted` trap.** In light mode `--muted` is literally `#000000`. Light-mode de-emphasis works through **alpha** (`text-foreground/70`), not through a grey. Only in dark mode is `--muted` a real grey. For application UI use `--app-ink-soft` and `--app-ink-muted`, which behave correctly in both modes.

Footer chrome stays charcoal in both modes and does **not** flip with the canvas: `--footer-bg` `#161719`, `--footer-elevated` `#212326`, `--footer-foreground` `#f3f4f6`, `--footer-muted` `#9ca3af`, `--footer-border` `#2e3238`.

### 2.4 Semantic colour

| Token | Light | Light soft | Dark | Dark soft | Meaning |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `--positive` | `#0f6b4f` | `#e2f2ec` | `#3fae86` | 16% tint | Money in, collected, healthy, success |
| `--negative` | `#b02a25` | `#fbe9e7` | `#f0766a` | 16% tint | Money out, overdue, destructive, failure |
| `--attention` | `#8a5b00` | `#fdf1dc` | `#e0a94a` | 16% tint | Needs review, pacing warning, low stock |
| `--info` | `#042940` | `#eff6fb` | `#609abe` | 16% tint | Neutral status, informational wells |

Every text-weight value clears 4.5:1 on its own canvas, and `--positive` also clears it on `--positive-soft` (5.6:1).

**Colour is never the only signal.** Every semantic value is paired with a sign, an icon or a word, so meaning survives colour-blindness, greyscale printing and a bad screen in bright sunlight.

### 2.5 Charts

Charts break design systems faster than anything else, so the series order is fixed. Same data, same reading, every product.

| Series | Light | Dark |
| :--- | :--- | :--- |
| 1 | `--navy-800` | `--navy-300` |
| 2 | `--navy-500` | `--navy-400` |
| 3 | `--navy-300` | `--navy-500` |
| 4 | `--navy-600` | `--navy-200` |
| 5 | `--navy-200` | `--navy-600` |

Rules: never reorder per chart · above five series, switch chart type rather than adding colours · gridlines are `--app-line` at 1px, axis labels `--app-ink-muted` at 11px · semantic colours appear only when the series genuinely means positive or negative · every chart has a text alternative, either a data table or a summary sentence · no gradient fills, no drop shadows, no 3D.

### 2.6 The one escape hatch

`--product-signal` lets a product line carry a single identity tint. It defaults to `var(--brand-primary)` and is permitted **only** on the product identity chip and the nav mark, on areas under roughly 40×40px. Banned on CTAs, links, focus rings, charts, text and large fills. Using it requires an entry in the project's deviation register.

---

## 3. Type

One typeface, four weights. Hierarchy comes from **weight and scale contrast**, not from pairing families.

**Inter Tight.** Not Geist, not Outfit, not DM Sans, not JetBrains Mono for UI.

```ts
// next/font
import { Inter_Tight } from "next/font/google";
const font = Inter_Tight({
  subsets: ["latin"],
  weight: ["300", "500", "700", "900"],
  variable: "--font-inter-tight",
});
```

```css
/* portable fallback */
@import url('https://fonts.googleapis.com/css2?family=Inter+Tight:wght@300;500;700;900&display=swap');
```

### 3.1 Weight roles

| Weight | Role |
| :--- | :--- |
| `900` | Wordmark, the `/`, display type, hero anchors, section lockups |
| `700` | Headings, labels, buttons, nav links, **all numeric figures** |
| `500` | Body copy, list descriptions, metadata — **the default body weight** |
| `300` | The `studio` suffix, lockup suffixes, large decorative numerals |

Body copy is `500`, not `400`. Inter Tight at 400 reads thin on these surfaces.

### 3.2 Scale

| Role | Size | Weight | Tracking |
| :--- | :--- | :--- | :--- |
| Display wordmark | `--fs-display` | 900 / 300 | `-0.01em` |
| Hero anchor | `--fs-hero` | 900 | tight |
| Section title | `--fs-title` | 900 | tight |
| Slab heading | 28px → `clamp(2rem, 3.6vw, 3.25rem)` | 700 | tight, `leading-[1.05]` |
| Section heading | 22px → 28px md → `clamp(1.75rem, 2.5vw, 2.25rem)` lg | 700 | tight |
| Lead | `--fs-lead` | 500 | — |
| Subheading | 18px → 20px | 700 | — |
| Body | `var(--text-body)` (15/14/13 by density) | 500 | — |
| Label | 13px | 700 | `0.04em` |
| Nav link | 13px | 500 | `0.04em` |
| Kicker | 13px | 700 | `0.08em`, `--accent-text` |
| Parenthetical kicker | 14px → 16px | 700 | — |
| Eyebrow | 11px | 700 | `0.12em`, uppercase |
| Bottom-nav label | 10px | 700 | `0.02em` |

**Measure:** body caps at `46ch`, prose columns at `--prose-max` (720px). Never run a paragraph across 1440px.

### 3.3 Numerals

| Role | Size | Weight | Tracking |
| :--- | :--- | :--- | :--- |
| Primary metric | 30px → 40px lg | 700 | `-0.02em` |
| Secondary metric | 22px → 26px | 700 | `-0.015em` |
| Inline row figure | `var(--text-row)` | 700 | `-0.01em` |
| Delta badge | 12px | 700 | `0` |
| Decorative numeral | `clamp(5.5rem, 18vw, 12rem)` | 300 | `leading-none`, `aria-hidden` |

**Every figure uses `tabular-nums`** (the `.vx-tabular` class). Columns must align across rows, and a value that changes must not reflow its row.

### 3.4 Money

- **Format** currency code, space, grouped thousands: `PKR 184,250`. Never `Rs.`, `₨`, an unspaced code or a bare number.
- **In** prefixed `+ ` in `--positive`. **Out** prefixed `- ` in `--negative`. **Neutral** no sign, `--app-ink`.
- **One formatter.** Every figure goes through `formatMoney()`. No inline `toLocaleString` in components — that is how two screens end up formatting the same number differently.
- **Direction is data, not styling.** Pass a positive magnitude plus a sign, so colour and prefix can never disagree.
- **Precision** whole units in lists and metrics; two decimals only in reconciliation and invoice lines, consistent down the whole column.
- **Locale and currency** come from the project profile. `PKR` and `en-PK` are defaults, not constraints.

### 3.5 Monospace

`--font-mono` is permitted for code, log lines, IDs, hashes, diffs, JSON and stack traces. It is banned in labels, buttons, navigation, headings and body copy. Expected in the `console` archetype, rare elsewhere.

---

## 4. Shape and depth

### 4.1 Radius

| Token | Value | Applies to |
| :--- | :--- | :--- |
| `--radius-control` | `12px` | Buttons, inputs, textareas, choice rows, chips |
| `--radius-slab` | `18px` | Cards, list rows, dropdowns, icon tiles |
| `--radius-panel` | `24px` | Metric slabs, hero cards, inverse panels |
| `--radius-sheet` | `28px` | Bottom sheets (`rounded-t-sheet`) |
| `--radius-pill` | `9999px` | Toggles, avatars, status dots, segmented thumbs |
| `--radius-marketing` | `60px` → `clamp(40px, 4.167vw, 60px)` lg | **`marketing` archetype only** |

Application surfaces never exceed 28px. A 60px radius inside a dashboard reads as a mistake.

### 4.2 Borders

| Weight | Token | Use |
| :--- | :--- | :--- |
| `1px` | `--app-line` | Dividers, table rules, input rest state |
| `1.5px` | `--app-line-strong` | The standard slab outline — the default |
| `2px` | `--brand-primary` | Selected state, focus ring, brand panel edge |

### 4.3 Shadows

| Token | Use |
| :--- | :--- |
| `--shadow-slab` | Every resting slab |
| `--shadow-slab-strong` | FAB, floating bar, dragged item |
| `--shadow-slab-pressed` | Pressed state, paired with `translate(2px, 2px)` |
| `--shadow-sheet` | Bottom-sheet lift — the only blurred shadow in the system |

Zero blur everywhere else. Brand-filled slabs carry **no** shadow: on navy the offset is invisible in light mode and muddy in dark. `marketing` slabs carry no shadow either — their depth comes from scale and overlay.

The press is load-bearing on touch, where there is no hover: `translate(2px, 2px)` while the shadow collapses from `3px 4px` to `1px 2px`, so the object appears to move into its own shadow.

---

## 5. Density and layout

### 5.1 Density

Set `data-density` on any container. It nests — a dense table inside a comfortable page is valid and common.

| Variable | `comfortable` | `compact` | `dense` |
| :--- | :--- | :--- | :--- |
| `--control-h` | 48px | 40px | 34px |
| `--control-pad-x` | 20px | 16px | 12px |
| `--slab-pad` | 24px | 18px | 14px |
| `--stack-gap` | 18px | 14px | 10px |
| `--row-pad-y` | 12px | 10px | 7px |
| `--field-gap` | 20px | 16px | 12px |
| `--icon-tile` | 36px | 32px | 28px |
| `--text-body` | 15px | 14px | 13px |

**Touch floor.** On `@media (pointer: coarse)`, `compact` and `dense` clamp `--control-h` back to 44px and `--row-pad-y` to 12px. This is what keeps a dense console usable when someone opens it on a tablet.

Components read these variables rather than hardcoding heights, which is what lets one component library serve a marketing page and an admin console without forking.

### 5.2 Layout

| Token | Value | Use |
| :--- | :--- | :--- |
| `--page-max` | `1440px` | Marketing page width |
| `--app-max` | `740px` | Single-column app width |
| `--prose-max` | `720px` | Form and prose columns |
| `--nav-height` | `72px` | Top bar |
| `--bottom-nav-height` | `70px` | Mobile bottom nav |
| `--sidebar-width` | `264px` | Dashboard sidebar (72px collapsed) |
| `--page-pad` | `24px` → `32px` md → `clamp(24px, 2.778vw, 40px)` lg | Horizontal page padding |
| `--page-inset` | `24px` → `clamp(64px, 11.111vw, 160px)` lg | Deep editorial inset |

Vertical rhythm uses `--stack-gap` between major blocks, roughly `0.75×` between a heading and its content, and `--row-gap` between sibling rows.

---

## 6. Motion

**One curve, system-wide, never mixed:** `cubic-bezier(0.16, 1, 0.3, 1)` — `--ease-vertex`, or `[0.16, 1, 0.3, 1]` in JS.

| Interaction | Spec |
| :--- | :--- |
| Section reveal | `opacity 0→1`, `y 24→0`, `--dur-reveal`, once, `amount: 0.6` |
| Hero reveal | `y 28→0`, `--dur-hero`, `0.04s` stagger |
| Form panel reveal | `y 32→0`, `0.9s`, `amount: 0.2` |
| Step change | wait-mode swap, in `y 16`, out `y -12`, `--dur-panel` |
| Icon swap | wait-mode swap, `y 5/-5`, `--dur-swap` |
| Disclosure | height `0 → auto`, `--dur-panel` |
| Sticky stack | `scale 1→0.94`, overlay `0→0.42`, scroll-linked |
| Bottom sheet | backdrop `opacity 0→1`; panel `y 100%→0%`, `--dur-panel` |
| Press | `translate(2px, 2px)` + shadow collapse, `--dur-micro` |
| Row removal | height and opacity to `0` over `0.3s`, then filter from state |
| Loop | only the scroll cue: `y [0, -24]`, `1.15s`, linear — the one correct use of linear easing |

**Budgets.** `full` allows everything. `reduced` drops scroll-linked transforms and loops. `minimal` is opacity-only. The archetype sets a default; the project profile can lower it. Lowering the budget is a legitimate engineering decision, not a failure.

**Travel budget** 16–32px. A 50px+ slide reads as a template transition.

**Density budget** one animated element per viewport. A reveal and a sticky scale in the same band is one too many.

**Never re-animate on state change.** Reveals fire on first mount. Re-animating a list on every filter change is actively hostile to someone doing their job.

**Reduced motion** is not an afterthought. Every animated component reads the preference and passes `initial={reduce ? false : …}`, with scroll transforms skipped entirely. `tokens.css` also collapses all transitions globally as a backstop.

---

## 7. Components

Working implementations: [`../code/react/primitives.tsx`](../code/react/primitives.tsx).

Every interactive component declares all of: **rest, hover, focus-visible, active, disabled, loading, invalid** — wherever each applies. A component missing a state is not finished.

### 7.1 Buttons

| Variant | Fill | Hover | Use |
| :--- | :--- | :--- | :--- |
| `primary` | `--brand-primary`, white text | `--brand-hover` | The one main action per view |
| `inverse` | `--app-surface-inverse` | press offset | Secondary emphasis, FAB |
| `ghost` | transparent, `--app-line` border | `--app-line-strong` | Tertiary, toolbars |
| `danger` | `--negative`, white text | `opacity .9` | Destructive, always confirmed |

Geometry: `min-height: var(--control-h)`, `padding-inline: var(--control-pad-x)`, `--radius-control`, 700 weight.

- **Focus** `outline: 2px solid var(--focus-ring); outline-offset: 2px` via `focus-visible`.
- **Active** colour darkens on flat buttons; `translate(2px, 2px)` plus shadow collapse on shadowed ones.
- **Disabled** `opacity: 0.4`, `pointer-events: none`, `aria-disabled="true"`. Never a grey fill — greys read as a second brand colour.
- **Loading** label stays, a spinner replaces the icon, width is preserved, `aria-busy="true"`. The label never disappears without a replacement.
- **One primary per view.** Two primaries means the hierarchy is unresolved.

Also: **text link** (`--accent-text`, 700, colour change on hover, never opacity) · **filter chip** (pill, `aria-pressed`, inverse fill when active) · **FAB** (inverse, `--shadow-slab-strong`, press offset, above the bottom nav).

### 7.2 Slabs

The single container primitive. Four tones, one geometry.

| Tone | Fill | Border | Text |
| :--- | :--- | :--- | :--- |
| `surface` | `--app-surface` | 1.5px `--app-line-strong` | `--app-ink` |
| `brand` | `--brand-primary` | none | `#ffffff` |
| `inverse` | `--app-surface-inverse` | none | `--app-on-inverse` |
| `attention` | `--attention-soft` | 1.5px `--attention-line` | `--app-ink` |

**Metric slab** — eyebrow, figure at 30→40px/700/`-0.02em`, optional delta badge, optional icon tile at `--icon-tile`, and a footer row divided by a hairline holding a caption and one text action.

A delta without a stated comparison period is not shippable. `+4.8%` alone is unreadable; `+4.8% vs last week` is information.

**Stat card** — icon tile, eyebrow, figure, caption, in that order top to bottom. Grid `grid-cols-2` on tablet, 3 or 4-up on desktop.

**List row** — optional avatar or initials tile, title at 700, detail at `--text-detail` in `--app-ink-muted`, trailing figure. **Hover changes the border, not the fill** — a fill change across 40 rows flickers.

**Sticky stack slab** — `marketing` only. One per viewport, scroll-linked scale and overlay, stepped sticky offsets.

### 7.3 Data table

The `console` and `dashboard` workhorse.

- **Header** sticky, `--app-surface-sunken`, 11px/700/`0.12em` uppercase, `--app-ink-muted`. Sortable headers are buttons with `aria-sort`.
- **Rows** `--row-pad-y` vertical padding, 1px `--app-line` rules, `--app-surface-hover` on hover. Never zebra striping — the rules already do that job.
- **Numeric columns** right-aligned with `tabular-nums`. Text columns left-aligned. Never centre a column of data.
- **Selection** checkbox column, header checkbox for select-all, a sticky action bar showing the count when anything is selected.
- **Sort and filter state** lives in the URL so a view can be shared.
- **Above 100 rows** virtualise, and keep the header sticky.
- **Long values** truncate with the full value available on hover or expansion, never silently cut.
- **Empty, loading and error** states are required, and "empty by filter" is a different state from "empty by data".
- **Below `800px`** the table becomes a card list. A horizontally scrolling table on a phone is a failure, not a fallback.

### 7.4 Forms

**Input geometry:** `min-height: var(--control-h)`, `--radius-control`, 1px `--app-line`, `--app-surface` fill, **16px font on mobile** or iOS Safari zooms on focus.

| State | Treatment |
| :--- | :--- |
| Rest | `--app-line` |
| Hover | `--app-line-strong` |
| Focus | `--brand-primary` border plus the focus ring |
| Filled | no change — a filled input must not look disabled |
| Invalid | `--negative` border, `aria-invalid`, message below with `role="alert"` |
| Disabled | `--app-surface-sunken` fill, `--app-ink-muted` text |

Rules: every field has a real `<label for>` — a placeholder is never a label · required is marked with `*` in `--accent-text`, optional is marked "· Optional", and the form does not mix conventions · validate on blur and on submit, never on every keystroke · the error names the fix, not the failure · focus moves to the first invalid field on a failed submit · radio and checkbox groups live in a `<fieldset>` with a `<legend>`.

Also: **choice button** (radio semantics, `aria-pressed`, brand fill when selected) · **segmented switch** (sunken rail, `--app-surface` thumb with a shared layout animation, labels never move) · **search row** (icon, borderless input, clear button that mounts only when there is a query) · **toggle** (`role="switch"`, `aria-checked`, animated thumb).

### 7.5 Navigation

**Top bar** — `--nav-height`, transparent at rest, solid with a border on scroll. Three-column grid: mark, links, controls. Sanctioned `backdrop-blur`.

**Sidebar** — `--sidebar-width`, collapsible to 72px with tooltips. Active item takes `--accent-text` and a 2px leading edge. Groups get 11px/700 uppercase eyebrows.

**Bottom nav** — `--bottom-nav-height`, 3–5 tabs, `aria-current="page"` on the active tab, a centred 600px bar at `≥800px`. Body reserves bottom padding so the last row is never trapped.

**Stepper** — `01 / 04` at 13px/700/`0.08em` in `--accent-text`. Back always available, forward validates.

**Breadcrumbs** — past two levels of nesting. Current page is not a link.

### 7.6 Overlays

**Bottom sheet** (mobile-first) — `--scrim` backdrop, panel at `min(100%, 560px)`, `rounded-t-sheet`, `--shadow-sheet`, a handle bar, `max-height: 92svh` with internal scroll.

**Dialog** (desktop) — centred, `--radius-panel`, `--app-surface`, max 560px.

Both require: `role="dialog"`, `aria-modal="true"`, a labelled heading, focus trapped inside, focus returned to the trigger on close, `Escape` to dismiss, and body scroll locked while open.

**Tooltip** — supplementary only. Never the sole carrier of information, never on touch-only affordances.

**Toast** — one at a time, 5s minimum, dismissible, `role="status"` for confirmation and `role="alert"` for failure. Never for anything the user must act on.

### 7.7 Feedback states

**Never leave a container blank.** Every list, table and panel declares four states.

**Empty** — dashed `--app-line-strong` border, one icon in a `--app-surface-sunken` circle, a title at 700, a body under `40ch` in `--app-ink-muted`, and one action whose label names the action: "Add your first patient", not "Get started". Empty-by-filter offers "Clear filters" instead.

**Loading** — skeletons at the exact height of the real row, so nothing shifts. Spinners only in buttons and for indeterminate work. Under 400ms show nothing; a flash of skeleton is worse than a brief pause. Skeleton rows are `aria-hidden`, the container is `aria-busy`.

**Error** — `--negative-soft` fill, `--negative-line` border, `role="alert"`, a title in `--negative`, a body in plain language, and a retry action. Never a stack trace, a bare error code, or the word "unexpected". An error without a next step is a dead end.

**Success** — inline and quiet. A toast for a transient confirmation, a state change for a persistent one. Never a full-screen celebration.

---

## 8. Voice and content

- **Sentence case** everywhere: headings, buttons, labels, meta titles.
- **No Oxford comma.** "web, mobile and AI".
- **Numerals always.** `34%`, not "thirty-four percent".
- **Contextual CTAs.** "Book a Care demo", never "Submit" or "Learn More".
- **Banned words:** `leverage`, `synergy`, `seamless`, `revolutionize`, `cutting-edge`, `state-of-the-art`, `game-changing`, `next-gen`, `world-class`. Their presence fails the build.
- **Brand lockup:** `VERTEX` (900) + `/` (900, brand colour) + `studio` (300), `whitespace-nowrap`, from one config token. Never `VertexStudio`, never lowercase, never recolour the slash, never break across lines.
- **Product naming:** `VERTEX/name`, lowercase suffix.
- **Error copy names the fix.** "Add your name and an email or phone number", not "Validation failed".
- **Empty-state copy names the action.** "Add your first receipt", not "No data available".
- **Never invent** reviews, logos, SLAs, metrics, legal entities or addresses. If real content does not exist, ask for it or leave the slot visibly unfilled.

---

## 9. Accessibility

WCAG 2.1 AA is a build-time requirement, not a later pass.

### 9.1 Measured contrast

| Pair | Ratio | Verdict |
| :--- | :--- | :--- |
| `#000000` on `#ffffff` | 21.0:1 | Pass |
| `#042940` on `#ffffff` | 15.0:1 | Pass |
| `#ffffff` on `#042940` | 15.0:1 | Pass |
| `#0f6b4f` on `#ffffff` | 6.5:1 | Pass |
| `#0f6b4f` on `#e2f2ec` | 5.6:1 | Pass |
| `#b02a25` on `#ffffff` | 6.6:1 | Pass |
| `#8a5b00` on `#ffffff` | 5.9:1 | Pass |
| `#f3f4f6` on `#161719` | 15.4:1 | Pass |
| `#609abe` on `#161719` | 5.8:1 | Pass |
| `#90bbd5` on `#161719` | 8.8:1 | Pass |
| `#3fae86` on `#161719` | 6.5:1 | Pass |
| `#f0766a` on `#161719` | 6.4:1 | Pass |
| `#e0a94a` on `#161719` | 8.5:1 | Pass |
| **`#1a537c` on `#161719`** | **2.2:1** | **Fail — fill only** |
| `#609abe` on `#ffffff` | 3.1:1 | Decorative only |

### 9.2 Standing requirements

1. **Focus is always visible.** `focus-visible`, 2px `--focus-ring`, 2px offset. Never `outline: none` without a replacement.
2. **Targets are at least 44×44px,** including in `dense` on a coarse pointer.
3. **Colour is never the only signal.** Pair with sign, icon or text.
4. **State is announced.** `aria-pressed` on toggles, `aria-current` on active nav, `aria-checked` with `role="switch"`, `aria-expanded` with `aria-controls` on disclosures, `aria-sort` on sortable headers, `aria-busy` while loading, `role="alert"` on errors.
5. **Keyboard reaches everything.** Logical tab order, arrow keys within composite widgets, `Escape` closes overlays, no keyboard traps outside intentional modals.
6. **Decorative content is hidden.** Giant numerals and repeated wordmarks get `aria-hidden` and `select-none`, with the accessible name in an `sr-only` span.
7. **Motion is optional,** globally, with no per-component exceptions.
8. **Theme does not flash.** Resolved pre-paint; components read it via a subscription with a stable server snapshot.
9. **Zoom to 200% does not break layout,** and text reflows rather than clipping.
10. **Forced-colors mode works** — borders survive, shadows drop out.

---

## 10. Extending the system

A missing token is not a licence to write a hex.

1. **Check first.** Most "missing" tokens exist under a different name. Search `tokens.css`.
2. **Add in both modes.** A token that exists only in light mode is a dark-mode bug waiting to be filed.
3. **Measure the contrast** of any colour that will carry text, in both modes, and record the ratio.
4. **Mirror it** into `tokens.ts` and `tokens.json`.
5. **Document it** in the relevant table here.
6. **Log it** in `CHANGELOG.md` with a version bump.

A token in fewer than all four places will resurface as a hardcoded value in someone's component within a month.

**Adding a component:** confirm no primitive already covers it, define all seven states, define both themes, define keyboard and screen-reader behaviour, then add it to `primitives.tsx` and §7 here.

**Adding an archetype:** see the last section of `SURFACES.md`. The answer is almost always no.
