# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `0721d9a`, 1 dirty path, 10 unpushed commits
Spec: `specs/ft91-gate-fastpath/spec.md` (Status: staged)
Gate: green at `543ca92` — stale, work tree `fe3ed27`

## State

- **`--full` run phase: implement (in progress).** Tickets derived under
  `specs/ft91-gate-fastpath/tickets/` — ten tickets, two chains: canary lever
  (guard-marker-reader → carry-per-test-scoping → migrate-fixture-test-bindings
  → enforce-test-bindings) and gate lever (extend-freshness-close-tools →
  reuse-gate-verdict → plumb-gate-fresh-flag → collapse-commit-reuse-check),
  then revise-adr-0002 and record-fastpath-evidence.
- **No ticket has landed yet.** Frontier: guard-marker-reader and
  extend-freshness-close-tools, dispatched to worktree-isolated write
  delegates at the mid binding.
- **Closed decisions stay closed:** lever 3 (`-count=1` removal) refused;
  scoped baselines rejected (wide-baseline principle); `bench commit` gains no
  `--fresh`; `shellcheck` stays undeclared in the closure.

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
