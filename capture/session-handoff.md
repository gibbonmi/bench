# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: the `parallel-landings` integration worktree, base `1a135f1b` (the spec commit); `main` is at that commit
Spec: `specs/parallel-landings/spec.md` — `Status: staged`, reviewer-approved 2026-08-23
Gate: green on the spec commit (2026-08-23)

## State

`/bench-implement-spec --full parallel-landings --reviewer opus medium` is in
the build phase. The integration worktree (label `parallel-landings`) holds the
build; `bench preflight build parallel-landings` was green before the first
edit. Tickets commit serially in `Blocked by:` order on the worktree. The
reviewer pre-approved the worker's judgement on build-time decisions and capped
the review at two loops (opus / medium).

Ticket order: `make-spec-optional-on-the-landing`, then
`close-a-tickets-only-folder-on-the-landing`, then
`union-merge-the-phase-owned-journals`, then
`name-the-source-repair-in-the-conflict-refusal`, then
`route-every-phase-through-a-worktree` (fable / high; the other four opus /
medium).

Closed decisions: merge composition stays the landing primitive, no rebase;
the journal union and the destination default; the light path joins the
worktree rule; the worktree rule is guidance, not a hook. The three parked
ideas in `capture/IDEAS.md` close by implementation at the next drain.

## Next command

`/bench-implement-spec --full parallel-landings --reviewer opus medium` (resumes from the worktree's committed tickets)

## Shape

Rewritten in full at every phase close, pruned rather than accreted. A fresh
session pays for every line it reads cold; drop anything it would not act on.

Operational gotchas are placed by lifetime, not copied here. One that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes. One scoped to a build
belongs instead in that spec's coverage rows.

This file names at most when you'll hit one, never the command — a second copy
drifts from the source.

Keep the three sections above. **State** holds what is true now, including anything
uncommitted. **Next command** holds the exact harness-native invocation, not a
description of it. This section is the third.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
