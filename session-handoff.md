# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — spec staged and approved (`92968c5`), **unpushed**; `git push`
is hook-blocked for the worker, so the push is the reviewer's own command.

## State

- **`specs/ft91-canary-check-scoping.md` is staged and reviewer-approved** —
  slice A of the gate-pipeline map (canary check-scoping: registry
  family→check table, `BENCH_CONFORMANCE_CHECK` env, all-loud fail postures,
  per-check shared vacuity baselines keyed on check name alone, seven stray
  CHECK files). Coverage map valid at 16 rows. A Sol falsification pass
  (Codex CLI, reviewer-granted) blocked the first draft on seven findings;
  all were verified against the tree and folded in before approval.
- **`decisions/gate-pipeline.md` stays closed** — its Handoff carries the
  seams for slices B (manifest + DAG runner) and C (`checkGoCore` split +
  fixture migration), which spec after A ships.
- **Decisions that stay closed:** baseline grouping key is the resolved check
  name alone (unscoped fixtures share today's single full baseline); the
  live sweep's did-not-bite verdict is the binding's enforcement; no
  fixture merging; story lines as approved in the spec (six mid, one cheap).
- Codex CLI note: `codex exec` must run with stdin closed (`</dev/null`) or
  it blocks reading the pipe forever — cost two dead attempts this session.
- `bench prep-release` stays shelved — blocked by FT116's race and FT142's
  ship-track findings; both are board rows, not handoff state.
- The branch/worktree sweep (23 non-`main` branches, 19 worktrees) remains
  proposed, not executed — reviewer's call.

## Next command

`/bench-implement-spec specs/ft91-canary-check-scoping.md` on a fresh
mid-tier session. Reviewer-run `git push origin main` when convenient
(1 commit waiting).

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
