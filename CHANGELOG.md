# Changelog

## v0.2.0 — 2026-11-14

The "loop actually closes" release. v0.2.0 is the first version where the
workflow the README sells — pause at night, resume in the morning — works
end to end through the CLI alone, and where ingestion is safe to run from
scripts and agents.

### Fixed
- **Pause → resume loop closed.** `ctx pause` (and MCP `request_checkpoint`,
  `prompt-sync`, and the dashboard snapshot button) now write the checkpoint
  into `activeContext.md` — title stamp, Focus, and the Recent list (newest
  first, capped at 3) — so `ctx resume`, `ctx status`, `ctx list`, and the
  dashboard show what you actually paused on instead of the init template.
- **Ingest can no longer auto-confirm.** EOF or non-interactive stdin
  without `--yes` now aborts with a non-zero exit instead of silently
  applying the proposal. Scripts can gate on the exit code.
- **`ctx ingest <file> --yes` works in the documented order** (flags are
  now scanned position-independently; Go's flag package previously stopped
  parsing at the filename).
- **Ingest files notes where they belong.** A deterministic classifier
  routes decision-led notes to `decisionLog.md`, tech notes to
  `techContext.md`, architecture to `systemPatterns.md`, and product notes
  to `productContext.md` (previously everything landed in productContext).
- **Deduplication that actually dedupes.** Repeated lines inside one note,
  lines already recorded anywhere in the bank, and whole re-ingested notes
  are all dropped; re-ingesting the same file now reports "already present"
  and changes nothing.

- **MCP `append_memory_delta` corrupted the first heading** when called
  with no section: it built the header `"## "`, matched the prefix of
  `## Focus`, and inserted the delta between them. An empty section now
  appends at the end of the file, and named sections match whole heading
  lines only.
- **Corrupt manifest self-heals.** A checkpoint now rebuilds a corrupt or
  missing `manifest.json` from the files on disk instead of leaving it
  broken forever (`ctx doctor` still reports the corruption until then).

### Changed
- **Ingest never moves your original file.** It archives a *copy* under
  `research/archive/` and leaves the source in place.
- Audit `--apply` no longer writes empty sections; it states plainly when
  no manifests or components were detected.

### Repository
- Internal planning docs moved to `docs/internal/`; the VERTEX design
  system moved to `docs/design-system/`.
- `memory-bank/.state/` (checkpoints, manifest, dashboard token) is no
  longer tracked in git.

## v0.1.0 — 2026-09-05

Initial public release: 7-file memory bank, checkpoint engine, MCP agent
bridge with delta responses, local dashboard with per-project tokens,
brownfield audit, prompt-sync, and the VS Code extension.
