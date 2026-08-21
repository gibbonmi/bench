# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `60d1154`, clean tree, 2 unpushed commits
Spec: `specs/land-executable-freshness/spec.md` (Status: staged)
Gate: green at `9bda612` — stale, work tree `71c000e`

## State

`/bench-implement-spec --full` is mid-run. The one ticket
(`01-refuse-stale-landing-executable`) is built and committed green on the
retained integration source: worktree request `land-exec-freshness-01`, label
`land-executable-freshness integration`. Frozen review base
`60d1154b`, source tip `afcc2a3b`. The diff touches seven files: the freshness
owner gains a presence-only `DeclaresBuildInputs`, `LandCommand` gains the
invoked-executable parameter and proves it before any repository proof, the
registry closure forwards `Command.Executable`, and `roadmap/FT242.md` plus its
`ROADMAP.md` index and sequence lines carry the re-scope. All eight acceptance
rows have tests. Six mutation probes were run and each turned its row red,
including the coordinator's own probe of the check moved before the resume
dispatch.

One contestable call to veto or keep: the delegate also rewrote the
`ROADMAP.md` recommended-sequence item 1, which still pointed at the ask that
closed as FT225. That is beyond the ticket's literal wording.

The phase reached is review. Nothing is landed to the destination yet.

## Next command

`/bench-review-implementation` — base `60d1154b`, source tip `afcc2a3b`, in the
integration worktree; then `bench worktree land` from `~/workspace/bench`, then
`/bench-final-check`.

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
