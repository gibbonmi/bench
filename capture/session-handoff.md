# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `7c2684b`, clean tree, 4 unpushed commits
Spec: `specs/per-component-gate-scoping/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `ce22983` — stale, work tree `ce22983`

## State

- **The `/bench-what-next` drain landed at `7c2684b`, reviewer-approved.** Capture is
  empty: no parked ideas, no open learnings, no pending retros. The pass added FT180
  (spec-optional route decided at shape-idea's exit), dropped the check-level
  gate-scoping idea as already covered by the staged spec, and reworded FT113 —
  its stale-verdict face is resolved by the shipped reduced-gate scope (`specs/`
  is allowlisted, evidence is content-addressed), leaving only the
  flip-counts-as-a-path and one-flip-author residuals.
- The 16 ticket files under `specs/per-component-gate-scoping/tickets/` — authored
  at spec staging but never committed — landed in the same gated commit.
- Two specs sit staged and unimplemented: `specs/per-component-gate-scoping/spec.md`
  and `specs/pre-push-guard-visibility/spec.md`. The refreshed
  `## Recommended sequence` in `ROADMAP.md` orders the board: FT176 spec first,
  then the two staged implementations.
- **FT176 is the board's HIGH**: the spec-build lifecycle's preconditions deadlock
  mid-repair runs. Its permanently-active run record (the stuck
  `reduced-gate-phase-set` run, six registered worktrees, retained provisional
  refs) is inert and is the fix's acceptance fixture — leave it in place.
- Also inert: five recovery refs `bench worktree recovery` refuses because it
  cannot prove their payloads landed. The refusal is correct.
- Nothing has been pushed; `main` is ahead of `origin/main` (push is the
  reviewer's call).

## Next command

`/bench-write-spec` — FT176, the top line of the refreshed recommended sequence.

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
