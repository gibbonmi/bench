# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `8436d41`, 20 unpushed commits
Spec: `specs/ft91-phase-manifest-dag.md` — **implemented**, gate-green, unpushed.
Gate: green at `78c917e` — stale only by the spec status-flip doc commit.

## State

- **Slice B built and landed (2026-07-27), gate-green, not yet pushed.** Six
  commits: the DAG scheduler (`5a009ee`), the `.bench/phases.json` loader
  (`d09635f`), the deadline straggler report (`37039a3`), the canary fixture
  (`bd761db`), the review artifact (`d8fb574`), and the review fix pass
  (`6d1bb80`). All 29 coverage rows realized; `bench coverage --check` valid.
- **Pushing is the reviewer's call and has not happened.** Everything below the
  push is done.
- **Four items awaiting a reviewer verdict**, all surfaced in the build's exit
  report and none of them blocking:
  1. **Coverage row 7's red signal is wrong in the spec.** `os/exec` dedups
     `cmd.Env` keeping the last value, so the row's "append hands the child two
     values" cannot happen and its mapped test is vacuous. The unmapped
     `TestMergeEnvStripsThenSets` is the only non-vacuous coverage story 7 has.
     Amending the map is the reviewer's edit; the build left the spec alone.
  2. **Coverage row 6 names a test that no longer exists.**
     `TestRunnerPhaseDirIsRelativeToRoot` became
     `TestRunnerPhaseDirIsAbsoluteOrRoot` when `Phase.Dir` was given one
     anchoring authority; the real graded-root semantic is pinned by
     `TestManifestDirResolvesAgainstGradedRoot`.
  3. **`internal/gate/` holds 21 source files against a granted 16** (19 before
     this spec). Advisory — `bench structure` is not a gate check — so it is a
     re-grant or a split, not a defect.
  4. **One residual fail-open left open deliberately:** a deadline firing with no
     phase in flight can still print `gate: green` at the `runPhases` seam.
     `Execute` overrides it to 124 via `context.Cause` and a bare
     `bench gate-phases` carries no deadline, so it is unreachable as a false
     green.
- **The semantic review found four fail-open defects, all now closed and
  verified through the built binary:** a duplicate JSON key silently shadowed a
  phase and the gate went green; an unsatisfiable graph reported green having run
  nothing; a phase exiting 130 on its own was read as cancellation, leaving a red
  gate with a `Pending` verdict; an optional phase with an unusable working
  directory read as `skipped (not installed)`.
- **Next capability is slice C** (`checkGoCore` split + fixture migration +
  parity, decision map #3/#6/#7), whose spec also carries FT143's kit-root
  family→check binding assertion. Slice C is what first consumes the manifest —
  the kit ships no `.bench/phases.json` yet, and consumer-facing manifest docs
  land with it. FT143's roadmap row stays until that ships. FT144's workflow
  decision remains the reviewer's, unmade.
- `.bench/learnings.md` carries one open entry (a returned delegate is not a
  drained one — its background test sweep flaked a load-sensitive gate).
- `bench prep-release` stays shelved — blocked by FT116's race and FT142's
  ship-track findings; both are board rows, not handoff state.
- The branch/worktree sweep remains proposed, not executed — reviewer's call.
  The six worktrees this build cut were all released.

## Next command

Push is yours. After that, `/bench-what-next` — `bench status` flags one open
learning to drain, and the reconcile will want the FT91 row updated for slice B.

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
