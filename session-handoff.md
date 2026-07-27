# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean, 14 unpushed commits ahead of `origin/main`
Spec: `specs/ft91-gate-phase-split.md` (Status: implemented, deliberately unretired)
Gate: green at `08482c6`, all ten phases

## State

- **FT146 is fixed and gate-green.** `scripts/build-offline-archives.sh` now
  enumerates its output directory before doing any archive work and refuses
  unless every entry is a regular file named `redbench-*.tar.gz` or
  `redbench-*.tgz`. Reproduced first through the script itself — a git worktree
  with a committed file came back as four release archives, exit 0 — and that
  repro no longer reproduces. Regression test is
  `TestOfflineArchiveBuildRefusesOutputItCannotAccountFor`.
- **Half of the FT146 row was a mis-attribution, and it needs your verdict.** The
  row's second half says the artifact contract tests can resolve their output
  directory to the graded root. They cannot: both call sites in
  `artifact_offline_test.go` already use `t.TempDir()`, and the production caller
  passes a fresh path inside its own mktemp stage, so neither ever reaches the
  destructive branch. Nothing is left to build there. **The row should come off
  the board** — that removal is yours, not something the fix commit took.
- **The original destroying invocation was never pinned.** The defect is proved
  through the accused script, which is what the fix answers for, but which caller
  passed a live worktree as `<output-dir>` on 2026-07-27 is still unknown. The
  capture's "nine tarballs" means it ran in same-output mode. If you want that
  chased, it is a separate question from the fix.
- **A stale pool entry still needs your `rm -rf`.**
  `~/.bench/worktrees/bench-2826441890/220aa857…-72b9811f…` holds the tarballs and
  no git repo; `bench worktree release` and `bench worktree clean` both fail
  closed on it. Until it goes, `bench status`'s git row reads
  `git state unavailable` — that blinding is FT145.
- **`ft91-gate-phase-split` is still unretired on purpose.** Stories 4 and 5
  shipped as probed phases instead of the kit-owned `.bench/phases.json` they
  named, and story 9 is unbuilt as unsatisfiable. Retiring the spec destroys your
  veto surface.
- **Push needs `bench gate pin` first** — `.bench/` changed in slice C and the
  pre-push hook wants a fresh pin. That command needs an interactive TTY, so it
  is yours.
- One open `.bench/learnings.md` entry: the drain promoted a capture's diagnosis
  without checking it, which is what produced the false half of FT146.

## Next command

`/bench-what-next`

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
