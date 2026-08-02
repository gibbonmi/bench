# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `627be9c`, 2 dirty paths, 7 unpushed commits
Spec: `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `9adeae0` — stale, work tree `1072ea8`

## State

- **pcgs (`per-component-gate-scoping`): promoted and retired.** The full
  `bench spec build` lifecycle ran end to end: promotion published `f341800`
  (candidate `2e1a61c`), the spec was marked implemented by promotion and
  retired at `627be9c` (roadmap row removed; the field-set-slice residual
  recorded in `projects/benchkit.md`). Five review findings stand flagged for
  reviewer veto — S1 reduced-run reachability comment, S2 provenance comments,
  S3 dead ReusableGreen branch in prep-release Refusal, Sp1 derivation-source
  check does not bind an entry to its named derivation, Sp5 stale superseded
  repair-ticket text — detailed in the review receipt retained by the run and
  in this session's exit report.
- The 2 dirty paths are pending capture by design: the pcgs retro
  (`capture/retros/per-component-gate-scoping.md`) and a new learnings entry
  (delegate scratchpad collision); the drain owns their landing.
- **Cleanup ready when the reviewer says so:** the 21 out-of-pool
  `../bench-pcgs*` worktrees and their `pcgs/*` + `per-component-gate-scoping`
  branches are fully subsumed by the promotion (every ticket ported byte-exact
  or superseded by review repairs) — `bench worktree clean` per path retires
  them. FT176's inert fixtures (the stuck `reduced-gate-phase-set` run record,
  five refused recovery refs) stay in place.
- FT176 remains complete per the prior handoff; its reviewer-call findings
  await the drain in `capture/IDEAS.md`.
- Nothing has been pushed; push is the reviewer's call.

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
