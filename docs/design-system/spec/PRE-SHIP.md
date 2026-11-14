# Pre-Ship Gate

**Walk this before reporting any UI work as finished.** Not a suggestion — the gate is what stops small compromises from becoming the house style.

Rules for using it:

- **Verify, do not assume.** "Looks right" is not verification. Each item below says how to check it.
- **Report what you checked.** Name the items you verified and how. A claim of completeness without evidence is worth less than an honest partial.
- **A failed blocker stops the ship.** Fix it, or record it in the profile's deviation register with a reason and a review date. Silence is not an option.

---

## A. Blockers — none of these may ship

Any single unchecked item here means the work is not done.

### A1 — Tokens

- [ ] **Zero hardcoded colours.** Grep the diff for `#`, `rgb(`, `hsl(`. Every hit is either a token definition in `tokens.css` or a violation.
- [ ] **Zero hardcoded radii, shadows or durations.** Grep for `border-radius:`, `box-shadow:`, `cubic-bezier`, `transition:` with a literal.
- [ ] **No second accent.** No cobalt `#2B5EF5`, no orange, no purple, no per-feature hue.
- [ ] **No gradients** of any kind.

### A2 — Both themes

- [ ] **Every new surface renders correctly in light mode.** Look at it.
- [ ] **Every new surface renders correctly in dark mode.** Actually toggle it. Do not infer from the token names.
- [ ] **No dark-mode flash on first load,** including a hard refresh with dark set.
- [ ] **`--brand-primary` is never text on a dark canvas.** Grep for `text-brand-primary` and confirm every hit sits on a light canvas or a filled surface.

### A3 — Accessibility

- [ ] **Every interactive element has a visible focus indicator.** Tab through the whole surface and watch.
- [ ] **Every target is at least 44×44px,** verified on a touch viewport, including `compact` and `dense` regions.
- [ ] **Text contrast clears 4.5:1** (3:1 for large text), in both themes. Use the measured table in `DESIGN.md` §9.1, or a contrast checker for anything new.
- [ ] **Colour is not the only signal** for any status, direction or error.
- [ ] **Keyboard reaches and operates everything.** No mouse. Overlays close on `Escape`. No traps outside intentional modals.
- [ ] **Every form control has a real label.** A placeholder is not a label.
- [ ] **State is announced:** `aria-pressed`, `aria-current`, `aria-checked`, `aria-expanded`, `aria-sort`, `aria-busy`, `role="alert"` where each applies.
- [ ] **Decorative content is hidden** from assistive tech, with an `sr-only` name where the visual carried meaning.

### A4 — States

- [ ] **Every list, table and panel has an empty state.** Trigger it.
- [ ] **Every one has a loading state.** Throttle the network and look.
- [ ] **Every one has an error state.** Force a failure and look.
- [ ] **Empty-by-filter is distinct from empty-by-data.**
- [ ] **Every interactive element has hover, focus-visible, active and disabled** where each applies.
- [ ] **Every async action has a loading state** that preserves layout and sets `aria-busy`.

### A5 — Motion

- [ ] **`prefers-reduced-motion` is honoured.** Enable it at OS level and confirm the surface is still fully usable and nothing is stuck mid-animation.
- [ ] **One easing curve.** Grep for `cubic-bezier` and `ease-` and confirm only `--ease-vertex` appears.
- [ ] **Nothing re-animates on a state change.** Change a filter or a tab and confirm the list does not replay its entrance.
- [ ] **At most one animated element per viewport.**

### A6 — Content

- [ ] **No invented facts.** No placeholder reviews, logos, metrics, SLAs, entities or addresses.
- [ ] **No banned words.** Grep for: `leverage`, `synergy`, `seamless`, `revolutionize`, `cutting-edge`, `state-of-the-art`, `game-changing`, `next-gen`, `world-class`.
- [ ] **Sentence case** on every heading, button and label.
- [ ] **No lorem ipsum** in anything reviewable.

---

## B. Required — fix now or register the deviation

### B1 — Archetype fidelity

- [ ] The surface declares `data-vertex-surface` and, where it differs from the default, `data-density`.
- [ ] Radii stay within the archetype's ceiling. Application surfaces do not exceed 28px.
- [ ] The canvas matches the archetype.
- [ ] The navigation pattern matches the archetype.
- [ ] The motion budget matches the archetype, or the profile records a lower one.
- [ ] No component banned by the archetype appears.

### B2 — Layout

- [ ] Nothing shifts on load. CLS under 0.1.
- [ ] Body copy is capped near `46ch`; no paragraph runs the full page width.
- [ ] Works at 320px width without horizontal scroll.
- [ ] Works at 200% zoom — text reflows rather than clipping.
- [ ] Tables become card lists below `800px`.
- [ ] Fixed navigation does not cover content; bottom padding is reserved.

### B3 — Numbers and money

- [ ] Every figure uses `tabular-nums`.
- [ ] Every money value routes through `formatMoney()`; no inline number formatting.
- [ ] Direction is carried by a sign plus a colour, never colour alone.
- [ ] Decimal precision is consistent down each column.
- [ ] Every delta states its comparison period.

### B4 — Type

- [ ] Inter Tight only; weights limited to 900/700/500/300.
- [ ] Body copy is weight 500, not 400.
- [ ] Monospace appears only in code, logs, IDs and diffs.
- [ ] No text below 12px, and inputs are at least 16px on mobile.

### B5 — Code quality

- [ ] Complete TypeScript. No `any`, no TODOs, no unimplemented handlers, no commented-out blocks.
- [ ] Composed from existing primitives where they exist.
- [ ] No new dependency without a note in the profile.
- [ ] Linter and type-checker pass.
- [ ] Comments explain constraints, not what the code does.

---

## C. Verify before a release

Not per-change, but before anything reaches users.

- [ ] Lighthouse meets the profile targets — default Performance ≥90 mobile / ≥95 desktop, SEO ≥95, Accessibility ≥95.
- [ ] Core Web Vitals: LCP under 2.5s, CLS under 0.1, INP under 200ms.
- [ ] Tested on a real mid-range Android phone, not just a desktop emulator.
- [ ] Tested in Safari — iOS and macOS.
- [ ] Tested with a screen reader on the primary flow.
- [ ] Tested in Windows forced-colors mode.
- [ ] Fonts are preloaded and subset; no flash of unstyled text.
- [ ] Images are sized, lazy-loaded below the fold, and carry real `alt` text or `alt=""` when decorative.
- [ ] The project profile's adoption table is updated.

---

## D. Report template

Paste this, filled in, when reporting UI work.

```markdown
### Pre-ship report

Surface: [ name ] · Archetype: [ marketing | dashboard | mobile-app | form-flow | console ]
Density: [ comfortable | compact | dense ] · Motion budget: [ full | reduced | minimal ]

Blockers (A): [ all pass | list failures ]
- Tokens: [ how verified, e.g. grepped diff for # and rgb( — 0 hits outside tokens.css ]
- Both themes: [ how verified, e.g. toggled both, screenshots attached ]
- Accessibility: [ how verified, e.g. tabbed the full surface, contrast from DESIGN.md §9.1 ]
- States: [ how verified, e.g. forced empty, throttled to 3G, killed the endpoint ]
- Motion: [ how verified, e.g. enabled reduced motion at OS level ]
- Content: [ how verified, e.g. grepped the banned list — 0 hits ]

Required (B): [ all pass | list failures ]

Not verified: [ anything you could not check, and why ]
Deviations registered: [ profile entry numbers, or none ]
```

**On honesty.** An accurate partial report is more useful than an inflated complete one. If dark mode is untested, say so — that sentence costs less than the bug it prevents.
