# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `883aa362`, clean tree, 3 unpushed commits
Spec: `specs/ft230-release-through-bench/spec.md` — Status: staged, reviewer-approved 2026-08-20, two tickets.
Gate: green at `883aa362`.

## State

FT230's spec is staged, reviewed (2 iterations to accept), signed off, and
committed. The build has not started. Two tickets under
`specs/ft230-release-through-bench/tickets/`: `wire-adapter-selection.md`
(frontier) then `swap-workflow-and-flip-conformance.md` (blocked by the first).
Lines: group A opus/high, groups B and C opus/medium.

Decisions that stay closed: `--adapter npm|fixture` defaults to `fixture` with
no environment twin; the tag push is the reviewer's attended act, so CI submit
does not violate the runbook's presence rule; promotion stays out of CI. The
build retires two step-name byte contracts (the platform-first diagnostic in
`native_workflow_test.go` and the `preflight-publish-order-bypassed` canary)
after ticket 1 lands their replacement ordering assertion — retirements are
approved, recorded in the spec.

Earlier this session, the reviewed IntenTIC assessment folded into the board:
FT231 and FT106 amended, FT239–FT241 created, FT239 cross-referenced to audit
items A3/A12. The audit's P0/P1 spine is landed except A7, which is this FT230
spec. `capture/learnings.md` holds one undrained entry (the FT230 review-miss
learning with a proposed `craft-spec` rule change).

## Next command

`/bench-implement-spec ft230-release-through-bench` — fresh mid-tier session,
tickets serial on one retained integration source.

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
