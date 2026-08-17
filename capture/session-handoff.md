# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: retained integration worktree
`/home/devuser/.bench/worktrees/bench-3325222104/faa5f7ec961cc1595f594f518e2bf419-629a0b333f46f0815b7adc9143d0c7ba`
(branch `bench/assign/faa5f7ec961cc1595f594f518e2bf419/629a0b333f46f0815b7adc9143d0c7ba`),
base `main` HEAD `1a8b5a6f`
Spec: `specs/progressive-roadmap/spec.md` — Status: staged, implemented, reviewed;
not yet landed to `main`
Gate: green at source tip `63c4c97a` (this file's own commit will follow it)

## State

All five tickets of `specs/progressive-roadmap/spec.md` (FT198, the split
roadmap board) landed serially on the retained source, each independently
gate-verified and mutation-probed by the coordinator: `5ffb0470` (parser +
this repo's 67-row migration), `aebb0540` (roadmap-detail-integrity
conformance check + 8-fixture canary family), `cd9c60e0` (status
reconcile + stripped journey), `7dde37ad` (CLI rendering), `181c86f8`
(docs), plus two small fence-correction commits (`d4d0f998`, `25fa73a6`)
truing the spec's Ownership fences list to a migration-necessitated edit in
`internal/conformance/recurrence_maintenance_contract_test.go`.

`/bench-review-implementation` ran three axes in parallel (`opus`/medium):
Standards (7 raw / 4 repair targets), Spec (4 raw / 2 repair targets),
Coverage (8 raw / 6 repair targets, one instance shared with a Spec finding).
Findings and dispositions are in `reviews/progressive-roadmap.md`, committed
at `63c4c97a`. Worst issues: a test
fixture harness pasted across 5 sites with a same-package name shadow
(Standards); a coverage-row pin landed in the wrong package, transitively
covered but not directly (Spec); a degraded `roadmap/` directory state is
unexercised and the diff's own new sequence-trust entry for it is provably
dead code by test (Coverage).

**Not yet done: reviewer disposition on the findings.** Several are `ask-user`
(judgment, not mechanical) — notably: two canary fixtures with no
`MUTATE.json` (spec text vs. established tree convention, non-behavioral
contradiction, flag for veto); a wrapped-heading-with-present-row-file case
that emits both a wrapped-heading and a spurious orphan diagnostic (is the
double diagnostic intended?); a live symlink row file being accepted as
authoritative despite the spec's Won't-handle text claiming the classifier
already refuses it. The `auto-fix` findings (fixture-harness duplication,
diagnostic-format double derivation, PR-talk comments, the dashboard-pin
location, two coverage-test gaps) are repairable inside already-approved
scope without further sign-off.

## Next command

On reviewer disposition of the `ask-user` findings: route the accepted
findings (all `auto-fix`, plus whichever `ask-user` ones are accepted) as
slim repair tickets and resume `/bench-implement-spec progressive-roadmap`
on this same retained source. When repairs land green and review is
re-satisfied (or residual risk is explicitly accepted), the destination runs
`bench worktree land --request <opaque-id> --base 1a8b5a6f --source-tip <tip>
--spec progressive-roadmap -m <msg> <worktree-path>`, then `/bench-final-check`.

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

## Machine note

The globally installed `bench` on PATH resolves a stale 0.2.0 release,
untrusted by the gate's freshness seal against this tree's current source.
Working pattern used throughout this build: `bash scripts/go-build.sh
/home/devuser/workspace/bench dist/bench` from the main tree (gitignored,
safe to rebuild repeatedly), then every worktree invocation as
`BENCH_RUN_BINARY=/home/devuser/workspace/bench/dist/bench bench <args>`.
`bench commit`'s internal gate call refuses an inherited binary for a never-
before-graded prospective (uncommitted) tree state — run `bench gate --fresh`
once on the exact tree state first (caches a reusable green verdict), then
`bench commit` immediately after on the unchanged tree.
