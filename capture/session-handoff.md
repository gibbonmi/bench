# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5a9f3a5`, 3 dirty paths, 51 unpushed commits
Spec: `specs/worktree-exec-run-binary/spec.md` (Status: staged)
Gate: green at `18401c8` — stale, work tree `6737ac9`

## State

FT229 is landed and retired. Published commit `c8a3fad2`, gate green, spec
flipped to `Status: implemented` by that commit; `specs/ft229-hygiene-batch/`
is gone and its two durable hostile-input classes are promoted into
`projects/benchkit.md`. The integration source was released and removed. Zero
tickets-only folders remain.

Review found two fail-opens at the enforcement boundary that the gate could not
reach, both introduced by this spec's own narrowing of the degraded guard rim.
Both are closed and graded. That is the result worth carrying forward, not the
seven features.

Pending capture, uncommitted by design, all waiting on the drain:

- `capture/retros/ft229-hygiene-batch.md` — untracked.
- `capture/agent-performance/claude-models.md` — refreshed, uncommitted.
- `capture/learnings.md` — three open entries from this build.
- `capture/IDEAS.md` — three parked CLI gaps.
- `ROADMAP.md` carries three FT229 references, including the FT174 dependency
  row that FT229 unblocks. Reconcile is the drain's, deliberately left.

`specs/worktree-exec-run-binary/` is a staged spec authored in a parallel
session. It is not this session's work and nothing here has touched it.

## Next command

`/bench-drain`

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
