# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `6a69cdd9` before the drain commit, clean tree otherwise, unpushed backlog (~66 commits; push remains reviewer-owned)
Spec: `specs/single-build-serial-gate/spec.md` (Status: staged); `axi-query-disclosure` retired this pass
Gate: green at `d360fc0` cached tree (matches work tree)

## State

The 2026-08-12 `/bench-what-next` drain is drafted as one uncommitted batch
awaiting reviewer approval: FT173 collapsed to its R11 residual and the spec
retired; new rows FT202 (fence + census scope decision), FT203 (worktree list
flake, `/bench-debug`), FT204 (transcript query, decision), FT205
(`craft-delegate` release-path clause); the inherited-refusal learning
re-parked in FT6 pending a real `bench commit` repro; inbox, journal, and the
axi-query-disclosure retro drained to empty; sequence refreshed (FT171 spec,
FT203 debug, FT175 shape). If the drain commit (`... spec-retire:
axi-query-disclosure`) is absent from history, the batch is still awaiting
approval — review the working diff, then commit it whole on green.

Decisions staying closed: `internal/gittest` fence amendment; guards
real-stale fixture; coverage `why` rewording; FT175 shape unblocked by the
capstone landing per the 2026-08-02 reviewer ruling.

## Next command

`/bench-implement-spec` on `specs/single-build-serial-gate/spec.md`

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
