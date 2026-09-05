# Prompt — audit existing UI

Run this **before** converting anything. An audit is cheap, reversible and tells you whether a retrofit is an afternoon or a fortnight. Converting first and discovering the scope halfway through is how half-migrated codebases happen.

Output is a written report. No code changes.

---

## The prompt

```markdown
Audit this codebase against the VERTEX Universal Design system. Produce a written
report only — change no code.

Read VERTEX-universal-DESIGN/AGENTS.md, spec/DESIGN.md and spec/SURFACES.md first.

## Scope
[ e.g. all of src/, or src/features/billing/** ]

## Report these sections

### 1. Inventory
- Every distinct surface, and the archetype it most closely matches
- Component count, and how many are duplicates of each other
- Styling approach in use, and whether more than one is in play

### 2. Colour violations
- Every hardcoded colour value, with file and line, grouped by how often it appears
- Any second accent, gradient, or banned hue
- Whether dark mode exists, and which surfaces are missing it
- Any accent used as text on a dark canvas
- Contrast failures, with measured ratios

### 3. Type violations
- Font families in use, and any beyond Inter Tight
- Weights in use, and any outside 900/700/500/300
- Body copy at weight 400
- Monospace outside code, logs, IDs and diffs
- Text below 12px, and inputs below 16px on mobile

### 4. Depth and shape violations
- Blurred drop shadows and coloured glows
- Glassmorphism used as a card material
- Radii outside the archetype's ceiling
- Border weights other than 1px, 1.5px and 2px

### 5. Motion violations
- Every distinct easing curve found, and where
- Whether prefers-reduced-motion is honoured
- Animations that replay on state change
- Viewports with more than one animated element

### 6. State gaps
- Lists, tables and panels missing an empty state
- Missing loading states
- Missing error states
- Interactive elements missing hover, focus-visible, active or disabled

### 7. Accessibility gaps
- Missing or removed focus indicators
- Targets under 44x44px
- Form controls without a real label
- Missing aria state: aria-pressed, aria-current, aria-checked, aria-expanded,
  aria-sort, aria-busy, role="alert"
- Elements unreachable by keyboard
- Status conveyed by colour alone

### 8. Content violations
- Every banned word, with location
- Title Case headings
- Generic CTAs
- Invented reviews, logos, metrics, SLAs, entities or addresses
- Lorem ipsum or placeholder copy

### 9. Effort estimate
For each surface: Low (tokens only), Medium (tokens plus recomposition),
High (needs a rethink). Say what makes each one what it is.

### 10. Recommended order
Which surface to convert first and why. Favour the one that is highest traffic and
lowest effort, so the system proves itself early.

## Rules for this audit
- Cite file and line for every finding. An uncited finding is unusable.
- Count occurrences. "47 hardcoded hex values across 12 files" tells me the shape of
  the problem; "some hardcoded colours" does not.
- Separate what is genuinely wrong from what is merely different.
- Do not fix anything. Do not propose a diff. Report only.
- If the scope is too large to audit fully, audit a representative slice, say exactly
  what you covered, and extrapolate with your reasoning shown.
```

---

## Reading the report

**Hardcoded colours concentrated in a few files** is the good outcome. Installing the token layer and replacing those values usually resolves most of the visible drift on its own.

**Hardcoded colours spread across every component** means the conversion is a real project. Plan it surface by surface and keep the deviation register honest in the meantime.

**Missing dark mode** is usually the largest single line item, because it forces every hardcoded colour into the open at once. It is also the highest-value fix, since it makes the rest mechanical.

**Missing states** matter more than they look. They are invisible in a demo and the top source of production bug reports.

**A high effort estimate on a low-traffic surface** is a reason to leave it, register the deviation and move on. Not everything needs converting this quarter.
