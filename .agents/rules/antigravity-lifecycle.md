# Antigravity Rule Discovery, Access & Lifecycle Specification

## 1. Antigravity Customization Discovery Flow
Antigravity automatically discovers customizations by traversing specific directory paths starting from the active working file up to the repository root:

```text
[Current Working Directory]
       │
       ▼ (walk up to repository root)
[Repository Root: CTXbank/]
   ├── AGENTS.md                  <── Unconditionally loaded canonical agent rules
   ├── GEMINI.md                  <── Directory-level rule mirror
   ├── .agents/rules/*.md         <── Granular workspace rules catalog
   └── .agentrules / .cursorrules <── Vendor-compatible pointer stubs
```

## 2. Rule Access & Deduplication Mechanics
- **Unconditional Application:** `AGENTS.md` and `GEMINI.md` are always-on rules loaded at the root. They require no frontmatter and are applied to all child files and subdirectories.
- **Deduplication:** When Antigravity traverses the workspace, it maps rules by their resolved canonical paths. Even if reached via multiple paths or subagent calls, rules are injected exactly once per conversation turn.
- **Separation of Rules vs. Skills:** 
  - Rules (`AGENTS.md`, `.agents/rules/`) enforce invariant behavioral boundaries, coding standards, and safety gates.
  - Skills (`skills/<skill_name>/SKILL.md`) are progressively disclosed on-demand when relevant task triggers occur.

## 3. How Antigravity Interacts with the Memory Bank
Antigravity does not rely on transient memory across sessions or context window compactions. Instead, it accesses the local deterministic store:
1. **Pre-Execution Hook:** Antigravity reads `memory-bank/activeContext.md` and `memory-bank/progress.md` to establish ground truth.
2. **Execution Boundary:** Antigravity executes tasks conforming to `ORCHESTRATION_BLUEPRINT.md` and logs progress in `/docs/phases/`.
3. **Post-Execution Hook:** Antigravity updates `memory-bank/activeContext.md` with recent focus, next steps, and increments the progress ledger.
