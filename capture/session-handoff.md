# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5fbb241`, clean tree, 16 unpushed commits
Spec: `specs/per-component-gate-scoping/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged), `specs/spec-build-lifecycle-preconditions/spec.md` (Status: staged)
Gate: green at `7630c36` — current

## State

- **Phase reached: `/bench-implement-spec --full` implement, wave 1 of the
  `per-component-gate-scoping` spec build.** The lifecycle run is active at
  subject `5fbb241`; assignments `pcgs-t1-expose` (freshness exposure) and
  `pcgs-t2-fixture` (kit-shaped fixture root) are with write-delegates. The
  17 tickets were normalized at `3972744` to the lifecycle parser's shape
  (bracketed row ids, single-line path-only ownership fences).
- A complete reference implementation of the whole spec sits on the local branch
  `per-component-gate-scoping` (20 commits, base `acf02e8`); delegates port
  ticket-by-ticket from its mapped commits. The 20 dirty `../bench-pcgs*`
  worktrees are that build's leftovers — verify subsumption before cleaning.
- **A second concurrent session builds FT176** (`spec-build-lifecycle-preconditions`)
  in `.bench` assignment worktrees, landing per-ticket gated commits on main. Its
  marker ticket (`pass-the-found-green-marker-as-the-expected-prior-tip`) was landed
  by *this* session as `5fbb241` to undeadlock `spec build start`; the other
  session's copy in worktree `f52d997e…` is stale and should be dropped. Serialize
  gate runs and landings between the sessions.
- Inert, leave in place: the stuck `reduced-gate-phase-set` run record and the five
  refused recovery refs (FT176's acceptance fixtures).
- Nothing has been pushed; push is the reviewer's call.

## Next command

`/bench-implement-spec --full specs/per-component-gate-scoping/spec.md`

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
