# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `ad613ef`, 63 dirty paths, 2 unpushed commits
Spec: `specs/pre-push-guard-visibility/spec.md` (Status: staged), `specs/reduced-gate-phase-set/spec.md` (Status: staged)
Gate: green at `4fabf09` — stale, work tree `fd59740`

## State

- The `/bench-what-next` drain landed at `ad613ef`. All three capture sources are
  empty and the roadmap is reconciled; nothing from that pass is outstanding.
- Uncommitted here: the capture co-location migration, reviewer-approved as a light
  path. `IDEAS.md`, `session-handoff.md`, `.bench/learnings.md`, and `.bench/retros/`
  moved to `capture/`; `ROADMAP.md` stays at the root and `.bench-notes.md` stays
  put because it is per-worktree shift scratch, not repository capture. The four
  path constants in `internal/roadmap`, `internal/status`, `internal/learnings`, and
  `internal/retros` are the single source, and ~278 references across 57 files were
  swept to match.
- Two defects the migration exposed and fixed in the same diff: `bench init`
  constructed the journal path from parts instead of the constant, so it would have
  kept writing to the old location; and `bench idea` and `bench handoff` both assumed
  their target directory existed, which was free at the repository root and is not
  under `capture/`.
- Also uncommitted: `specs/reduced-gate-phase-set/spec.md` (Status: staged, FT168's
  second face) awaiting reviewer sign-off, plus the FT166→FT168 roadmap fold and the
  refreshed recommended sequence. The spec has been through a top-tier falsification
  pass whose three fatal findings are all addressed in the current draft.
- The spec's central design: a changeset confined to the allowlist runs a reduced
  phase set, and excludable phases run against a stripped worktree with capabilities
  required, so a check that would silently degrade to a capability skip reds instead.
  That posture is the fix for the first draft's central error — absence alone does
  not produce a red in this tree.
- `specs/pre-push-guard-visibility/spec.md` (Status: staged) is unstarted and now
  second in the sequence.
- Both merged specs remain safe to retire; six `recovered` worktree assignments with
  missing trees remain in the pool from the cadence build.
- Nothing has been pushed; main is well ahead of `origin/main`.

## Next command

`/bench-shape-idea` — the board's leading invocable signal (`decisions`).

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
