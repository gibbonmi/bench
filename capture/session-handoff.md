# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `e7b8569`, clean tree, 23 unpushed commits
Spec: `specs/per-component-gate-scoping/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged), `specs/spec-build-lifecycle-preconditions/spec.md` (Status: staged)
Gate: green at `9c254c8` — current

## State

- **FT176 (`spec-build-lifecycle-preconditions`): phase reached =
  `/bench-implement-spec --full` review.** All nine tickets are landed as
  per-ticket gated commits on main (light-path; no lifecycle run — a second
  concurrent run would recompose-churn against the pcgs run). This session's
  five: `5fc1b55` exempt-abandon, `9166960` docs recomposition, `dfcc71d`
  plan-absent-target, `eec8b22` identity/liveness split, `e7b8569` CLI refusal
  enumeration. A fresh-context three-axis review delegate over those nine
  commits is in flight; findings route as repair fix-and-gate commits, then
  `bench spec implemented` + `/bench-final-check` close the build. All FT176
  `.bench` assignment worktrees are retired (recovery refs kept).
- **pcgs (`per-component-gate-scoping`): `/bench-implement-spec --full`
  implement, wave 1** in the other session. Lifecycle run active; assignments
  `pcgs-t1-expose` and `pcgs-t2-fixture` with write-delegates; tickets
  normalized at `3972744`. Reference implementation on local branch
  `per-component-gate-scoping` (20 commits, base `acf02e8`); the 20 dirty
  `../bench-pcgs*` worktrees are its leftovers — verify subsumption before
  cleaning. Serialize gate runs and landings between the sessions.
- Inert, leave in place: the stuck `reduced-gate-phase-set` run record and the five
  refused recovery refs (FT176's acceptance fixtures).
- Nothing has been pushed; push is the reviewer's call.

## Next command

`/bench-implement-spec --full specs/spec-build-lifecycle-preconditions/spec.md`

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
