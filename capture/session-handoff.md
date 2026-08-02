# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `17af7e4`, 1 dirty path, 1 unpushed commit
Spec: `specs/per-component-gate-scoping/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `25a558c` — stale, work tree `0d7abbd`

## State

- **FT176 (`spec-build-lifecycle-preconditions`): complete.** All nine tickets
  landed as light-path gated commits on main; fresh-context three-axis review
  found no blocker; the four concrete findings landed as `5cb645b`; the spec is
  `Status: implemented`. Retro at `capture/retros/spec-build-lifecycle-preconditions.md`.
  Reviewer-call findings await the drain in `capture/IDEAS.md`: restart-after-terminal
  green-marker refusal (C1), husk/dangling-symlink liveness classification
  (C2/C3), the resume.go Planned-phase crash window, and the prepared-abandon
  identity-softening exemption. One spec-map precision note: the "drifted review
  binding" case in the story-1 five-fault row is unreachable in the promotion
  path and its test case was deleted — the spec row should drop to four cases
  (reviewer approval pending). All FT176 assignment worktrees retired,
  recovery refs kept.
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

`/bench-what-next` — the board's leading invocable signal (`drain`).

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
