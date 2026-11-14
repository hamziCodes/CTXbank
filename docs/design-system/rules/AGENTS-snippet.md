# Wiring the toolkit to your agent

Three ways to link this toolkit so an agent picks it up automatically. Do at least one; doing two is not harmful.

Adjust `VERTEX-universal-DESIGN/` in every snippet below if the toolkit sits elsewhere in the repository.

---

## 1. Cursor — project rule (recommended)

Copy the rule file into the target repository:

```bash
mkdir -p .cursor/rules
cp VERTEX-universal-DESIGN/rules/vertex-design.mdc .cursor/rules/
```

It carries `alwaysApply: true`, so it loads on every request without anyone remembering to mention it.

---

## 2. Any agent — root `AGENTS.md`

Append this to the repository's root `AGENTS.md`, or create the file if it does not exist. Most coding agents read it automatically.

```markdown
## Design system

All UI in this repository is governed by the VERTEX Universal Design toolkit in
`VERTEX-universal-DESIGN/`.

Before writing or changing any interface code, read in this order:

1. `VERTEX-universal-DESIGN/AGENTS.md` — the contract and adaptation procedure
2. `VERTEX-universal-DESIGN/project-profile.md` — what this project already decided
3. `VERTEX-universal-DESIGN/spec/SURFACES.md` — which archetype you are building
4. `VERTEX-universal-DESIGN/spec/DESIGN.md` — any visual decision
5. `VERTEX-universal-DESIGN/spec/PRE-SHIP.md` — before reporting work finished

Hard rules: one navy accent (`var(--brand-primary)` fills, `var(--accent-text)` text
and focus); never hardcode a colour, radius, shadow or duration; light and dark both
required; one easing curve `cubic-bezier(0.16, 1, 0.3, 1)`; every list declares empty,
loading and error; WCAG AA; never invent facts.

Reusable components: `VERTEX-universal-DESIGN/code/react/primitives.tsx`.
```

Claude Code users can point `CLAUDE.md` at it with a single line: `@AGENTS.md`.

---

## 3. Per-session prompt

When neither of the above is available, paste this at the start of a session:

```
This repository uses the VERTEX Universal Design toolkit in VERTEX-universal-DESIGN/.
Read AGENTS.md and project-profile.md there before touching any UI, follow the
archetype rules in spec/SURFACES.md, and walk spec/PRE-SHIP.md before telling me
you are done.
```

---

## Which task prompt to use

| You want to | Use |
| :--- | :--- |
| Build a new screen, page or component | [`../prompts/new-surface.md`](../prompts/new-surface.md) |
| Bring existing UI onto the system | [`../prompts/restyle-existing.md`](../prompts/restyle-existing.md) |
| Find out how far off the system a codebase is | [`../prompts/audit-existing.md`](../prompts/audit-existing.md) |

---

## Verifying the link works

Ask the agent, in a fresh session, without naming any file:

> What accent colour does this project use for text on a dark background, and why is it not the brand colour?

A correctly linked agent answers `--accent-text` / `#609abe`, and explains that `--brand-primary` (`#1a537c`) measures 2.2:1 against the `#161719` canvas and fails WCAG. If it guesses, invents a hex, or asks you which design system you mean, the link is not working — check the rule file location and the paths inside it.
