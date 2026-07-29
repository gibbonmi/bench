# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `c481e2b`, clean tree, 20 unpushed commits
Spec: `specs/ft91-gate-fastpath/spec.md` (Status: staged)
Gate: green at `fcafeb2` — stale, work tree `8bbfbc5`

## State

- **`--full` run phase: implement complete; review not yet run.** All ten
  tickets landed as gate-green commits from the staging commit forward; every
  TDD-able coverage row went red→green; measurements recorded in the spec's
  map (canary 25.2 s — stop rule met; forced full gate 128.0 s; unchanged-tree
  reuse 0.57 s).
- **Reviewer items pending:** escalation/falsification questions for the
  review pass; the load-coupled `TestExecuteDeadlineRecordsDistinctTimeout`
  timing flake (attributed in `.bench/learnings.md`, fix is a gate-test edit
  and therefore a reviewer call); the ADR reopen-trigger wording (commit-shaped
  vs work-shaped); the npm closure-hash cost (drop-npm lever if it ever bites).
- **Closed decisions stay closed:** lever 3 refused; scoped baselines
  rejected; `bench commit` gains no `--fresh`; `shellcheck` undeclared.

## Next command

`/bench-implement-spec --full ft91-gate-fastpath`

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
