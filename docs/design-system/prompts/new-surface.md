# Prompt — build a new surface

Use this to build a new screen, page, view or component on the VERTEX system. Fill in the brief, paste the whole thing, and the agent has everything it needs to make correct decisions without a design review.

---

## The brief

```markdown
Build a new surface on the VERTEX Universal Design system.

## What
[ e.g. A patient list for the clinic dashboard: searchable, filterable by status,
  showing name, last visit, outstanding balance and next appointment. ]

## Who uses it and on what
[ e.g. A clinic receptionist, on a desktop, all day, alongside a phone. ]

## What success looks like
[ e.g. She can find any patient in under 5 seconds and see whether they owe money
  without opening the record. ]

## Data available
[ Paste the type, schema or endpoint shape. If it does not exist yet, say so and
  propose one. ]

## Copy
[ Paste real copy, or write "needs copy" against each slot. Do not accept invented
  marketing lines or placeholder testimonials. ]

## Constraints
[ e.g. Must work offline. Must fit in the existing sidebar layout. No new dependencies. ]

## Procedure

1. Read VERTEX-universal-DESIGN/AGENTS.md and project-profile.md.
2. Pick the archetype from spec/SURFACES.md and tell me which one and why, before
   writing code.
3. Tell me which existing primitives and project components you will reuse, and what
   genuinely needs to be new.
4. Build it. Bind every colour, radius, shadow and duration to a token. Read the
   density variables rather than hardcoding sizes.
5. Implement every state: ready, empty, loading, error, and per-control rest, hover,
   focus-visible, active, disabled, invalid.
6. Walk spec/PRE-SHIP.md and report using the template in section D.

## Before you write code, tell me

- The archetype and density you chose, and why
- The composition order, block by block
- Anything in the brief that is underspecified, with your recommended default
- Anything you need from me that you cannot determine from the repository

Then stop and wait for my go-ahead.
```

---

## Why it is shaped this way

**"Who uses it and on what" settles more arguments than any other line.** A receptionist on a desktop all day wants `compact` density and keyboard shortcuts. A field technician on a phone in the rain wants `comfortable`, big targets and an offline state. Both are compliant; the user decides which.

**"What success looks like" prevents feature soup.** Without it, an agent adds every plausible affordance. With it, there is a test for whether an element earns its place.

**Asking for the plan before the code** costs one exchange and saves a rebuild. Archetype and composition order are cheap to correct in prose and expensive to correct in components.

**Naming what to reuse** is what stops the fifth slightly-different button from being born.

---

## Shorter version

For a small, well-understood addition:

```
Build [ what ] for [ who ] on [ device ].
Follow VERTEX-universal-DESIGN — read AGENTS.md and project-profile.md first.
Archetype: [ one of marketing | dashboard | mobile-app | form-flow | console ].
Reuse primitives from code/react/primitives.tsx.
All four container states required. Both themes required.
Walk spec/PRE-SHIP.md before telling me it is done.
```
