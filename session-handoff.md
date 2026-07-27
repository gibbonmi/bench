# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `19d1e5b`, 8 dirty paths, 9 unpushed commits
Spec: none staged.
Gate: green at `ad7607c` — current

## State

- **The FT91 interim-defect direct pass is committed at `19d1e5b`,
  gate-green.** `conformanceSubprocessEnv` now scrubs all three conformance
  control vars (`ROOT`/`TIER`/`CHECK`), and `formatProbeFailure` spills full
  probe output to `<git-dir>/bench-conformance-probe.log`, naming the path in
  the diagnostic (additive line; canary EXPECTs pin bare labels and are
  unaffected). Both defects can leave the FT91 row at the next drain.
- **`decisions/gate-pipeline.md` is closed and its dependency order is
  reviewer-confirmed (2026-07-27):** slice B next (manifest + DAG runner,
  tickets #1/#4/#9), then slice C (`checkGoCore` split + fixture migration,
  #3/#6/#7). FT143's cheap kit-root family→check binding assertion is
  bundled into slice C's spec (recorded in the map's Handoff); FT143's
  roadmap row stays until that ships.
- **Decisions that stay closed:** baseline grouping key is the resolved check
  name alone; the live sweep's did-not-bite verdict is the binding's
  enforcement; no fixture merging; the family→check table stays in
  `internal/conformance/registry`. FT144's workflow decision ("intent
  stands, seam moves") is still the reviewer's, unmade.
- Codex CLI note: `codex exec` must run with stdin closed (`</dev/null`) or
  it blocks reading the pipe forever.
- `bench prep-release` stays shelved — blocked by FT116's race and FT142's
  ship-track findings; both are board rows, not handoff state.
- The branch/worktree sweep (23 non-`main` branches, 19 worktrees) remains
  proposed, not executed — reviewer's call.

## Next command

`/bench-write-spec` — slice B of `decisions/gate-pipeline.md` (the manifest +
DAG runner), fresh mid-tier session per the profile's spec-authoring line.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
