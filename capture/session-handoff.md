# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `f8f90cc`, clean tree, 18 unpushed commits
Spec: `specs/land-spec-amendments/spec.md` (Status: staged)
Gate: green at `97e6b15` — stale, work tree `01d3cd8`

## State

FT225 is implemented and awaiting semantic review. `/bench-implement-spec --full`
is mid-run: all three tickets of `specs/land-spec-amendments/spec.md` committed
green, serially, on one retained integration source.

Retained integration source: assignment `5cb7d3bb7fd94540461c8951ba854102`,
request id `ft225-land-spec-amendments`, label `land-spec-amendments`.
Frozen review base `f8f90cc5`; source tip `7664600c`. The destination (`main`)
is still at `f8f90cc5` and clean. Both preflight modes are green on that source.

Landed on the source, in order: `12194544` implicit `specs/<slug>/`
authorization in `pathsAuthorizedCheck` (LS7-LS11); `cb1462a6` the landing
publishes the source's spec bytes and neutralizes the spec path before the
merge (LS1-LS6, LS12, LS13); `7664600c` the review-phase guidance names the
landing source as where a spec amendment commits (story 14).

Each ticket was probed twice — once by its author, once by the coordinator at a
different site and mutation kind — and every probe came back red. The published
commit's parents are asserted at both the unit and end-to-end seams.

Closed decisions that stay closed (also recorded in the spec): the landing
publishes `Implemented(source-tip spec bytes)` and the composition neutralizes
the spec path so a destination-side amendment — overlapping included — never
conflicts; `paths-authorized` implicitly authorizes `specs/<slug>/` via the
existing spec-path fact; the reviewer's typed `--base` and `--source-tip` are
the acceptance; the in-range fence-widening exposure is a named Won't handle
with the typed range as the control.

Two calls this run made, both flagged for veto: review axes run at the cheap
tier per standing preference rather than the profile's mid, and that drift is
queued for `capture/learnings.md`; the landing runs without a further pause
because `--full` over a signed-off spec is its standing approval. Nothing is
pushed — that stays reviewer-owned.

## Next command

`/bench-review-implementation`

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
