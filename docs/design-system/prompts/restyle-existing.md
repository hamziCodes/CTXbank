# Prompt — bring existing UI onto the system

Use this after [`audit-existing.md`](./audit-existing.md). Convert **one surface at a time**, each landing in a coherent state. Never leave an app half-converted across a release — a mixed interface looks more broken than a consistently old one.

---

## Stage 1 — land the tokens

Do this once per repository, before converting any surface. It is usually the change with the highest visible payoff per line touched.

```markdown
Install the VERTEX Universal Design token layer in this repository. Change no
component logic and no layout — tokens and theme wiring only.

Read VERTEX-universal-DESIGN/AGENTS.md first, then follow its first-run procedure.

1. Detect the stack: framework, styling approach, motion library, icon library,
   existing theme mechanism. Report what you found.
2. Fill in VERTEX-universal-DESIGN/project-profile.md completely and commit it.
3. Install the tokens using the one delivery path that matches this stack. Never
   mix two.
4. Add the pre-paint theme bootstrap from tokens/theme-init.ts so dark mode does
   not flash.
5. Map the project's existing colour variables onto VERTEX tokens. Give me the
   mapping table before you apply it, including anything with no clean equivalent.
6. Do not delete the old variables yet — alias them to the new tokens so nothing
   breaks mid-migration.

Report: what you installed, the mapping table, anything that had no equivalent, and
what visibly changed.
```

---

## Stage 2 — convert one surface

Repeat per surface, in the order the audit recommended.

```markdown
Convert [ surface name and path ] to the VERTEX Universal Design system.

Read AGENTS.md, project-profile.md and spec/SURFACES.md first.

## Procedure

1. Tell me the archetype this surface should be, and its density. Justify both.
2. List what changes: colours, radii, shadows, type, spacing, motion, states,
   component swaps. Group by "mechanical" versus "needs a judgement call".
3. Show me the plan and stop. Do not write code yet.
4. On my go-ahead, convert it:
   - Replace every hardcoded value with a token
   - Swap bespoke components for primitives from code/react/primitives.tsx where
     one exists
   - Add every missing state: empty, loading, error, and per-control states
   - Bring it within the archetype's radius ceiling and density
   - Reduce to one easing curve
   - Fix accessibility gaps found in the audit
5. Walk spec/PRE-SHIP.md and report using the template in section D.

## Rules

- Preserve behaviour. This is a restyle, not a rewrite. If you believe the behaviour
  is wrong, say so separately — do not silently change it.
- Preserve every existing test. If a test asserts a hardcoded colour, update the test
  and say that you did.
- Do not convert anything outside the named surface, however tempting the adjacent
  file looks.
- Anything you cannot convert goes in the deviation register with a reason, a scope
  and a review date — not into a silent compromise.
- If the conversion turns out to be larger than the audit suggested, stop and tell me
  before continuing.
```

---

## Stage 3 — retire the aliases

Once every surface is converted:

```markdown
Remove the legacy colour variable aliases added in stage 1.

1. Find every remaining reference to the old variables and report them by file.
2. Replace each with its VERTEX token.
3. Delete the alias block.
4. Confirm both themes still render correctly on every surface.
5. Update the adoption table in project-profile.md.

If any reference cannot be cleanly replaced, leave the alias, tell me which one and
why, and register the deviation.
```

---

## Common traps

**Converting colours but not composition.** Correct tokens on a template card grid is still off-system. The archetype rules in `SURFACES.md` are the other half of the job.

**Missing dark mode until the end.** Convert both themes together, per surface. Retrofitting dark mode across an already-converted app means touching everything twice.

**Silently dropping features to fit the spec.** If a surface needs something the system does not cover, that is a deviation to register or a token to add — not a feature to quietly remove.

**Converting the biggest surface first.** Start with a high-traffic, low-effort one. It proves the system, builds the muscle memory, and surfaces stack-specific surprises while they are still cheap.

**Leaving the profile stale.** The adoption table and deviation register are what the next session reads. An out-of-date profile is worse than no profile, because it will be trusted.
