# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `ea9dcaf`, clean tree, 21 unpushed commits
Spec: `specs/ft91-gate-fastpath/spec.md` (Status: staged)
Gate: green at `fcafeb2` — stale, work tree `e26358a`

## State

- **The `--full` run is complete.** Ten build tickets, a semantic review
  (fresh delegate, mid binding), a Codex falsification pass (top binding), and
  one reviewer-signed repair pass all landed as gate-green commits. Measured:
  solo canary 25.2 s (≤60 s stop rule met), forced full gate 128.0 s,
  unchanged-tree reuse 0.57 s.
- **Reviewer-ruled during the run:** commit reuses pre-lock via the gate home
  (`ExecuteReusingFreshGreen`); `npm` dropped from the subject closure (node
  stays); freshness 60 min; ADR 0002 posture 5 revised.
- **Open reviewer items (flagged, not applied):** the load-coupled
  `TestExecuteDeadlineRecordsDistinctTimeout` flake (attributed in
  `.bench/learnings.md`; a fix edits a gate test); the ADR reopen trigger
  stays commit-shaped; minor fail-closed observations from review (empty
  `-test.list` output, live-symlink markers, `--frsh`-as-positional plumbing).
- **Unpushed commits await the reviewer's push; ship-tier (`bench
  prep-release`) has not run.** Retirement of `specs/ft91-gate-fastpath/`
  waits for the retirement signal after merge; retro pending drain.

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
