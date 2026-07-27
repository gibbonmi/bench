# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — shape pass committed (`f56b522`), **unpushed**; `git push` is
hook-blocked for the worker, so the push is the reviewer's own command.

## State

- **`decisions/gate-pipeline.md` is nearly closed.** Resolved this pass: #1
  (manifest: five fields, `.bench/phases.json`, gate.sh stays the entry,
  absent=built-in / empty+malformed=red), #2 (runner survey: extension not
  rewrite; seams are `benchkitPhasesForCommand` and `splitSerialPhases`), #3
  (`checkGoCore` split: gofmt/vet/test/race out, build collapses into the
  serial build phase, structural residual stays; universal phases in the
  built-in table, kit-specific in the kit manifest; test enumeration via a
  bench plumbing subcommand), #4 (unweighted width, cancel-dependents red
  posture, gate-deadline-plus-name-stragglers timeouts), #5 (the prerequisite
  canary check-scoping slice: family→check registry table + CHECK override,
  `BENCH_CONFORMANCE_CHECK` env, all-loud fail posture, per-check shared
  vacuity baseline), #7 (fixture-backed parity + recorded reds).
- **Open: #6 and #8.** #6 is Research, agent-alone, now unblocked: bucket the
  39 package-core-guard fixtures by destination phase and verify the
  family→check assumption across all nine conformance families. #8 is the
  reviewer's task: make regroup-app's toolchain shape readable so the v1
  schema can be walked against it.
- The map's `## Handoff` section is written at close, after #6 and #8.
- **Decisions that stay closed:** ship is a superset of dev; diff-scoped
  gating stays ruled unsound; no check weakening for wall-clock;
  cross-language incrementality stays behind FT91's revive trigger; no weight
  field, no per-phase timeouts, no manifest tier field in v1.
- `bench prep-release` stays shelved — blocked by FT116's race and FT142's
  ship-track findings; both are board rows, not handoff state.
- The branch/worktree sweep (23 non-`main` branches, 19 worktrees, work
  verified present in `main`) remains proposed, not executed — reviewer's call.

## Next command

`/bench-shape-idea` resume on `decisions/gate-pipeline.md` #6 (Research,
agent-alone). #8 waits on the reviewer supplying regroup-app's build/test/lint
commands. Reviewer-run `git push origin main` still pending.

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
