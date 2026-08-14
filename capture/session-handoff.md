# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — spec staged at the latest `spec:` commit, capture close follows it
Spec: `specs/spec-authoring-and-light-path/spec.md` — `Status: staged`, reviewer-approved 2026-08-14, verification log `spec 5 + tickets 3`
Gate: green at the staging commit

## State

`spec-authoring-and-light-path` is staged and signed off with its eight tickets
under `specs/spec-authoring-and-light-path/tickets/`. Initial frontier (all
serial — every ticket shares `internal/anchors/registry_data.go` and
`tests/canary/workflow-guidance-anchors/`): lighten-light-path,
move-slicing-into-write-spec, cross-harness-reviewer-recipes,
drain-time-light-path, repair-stale-anchors; then main-session-authorship →
verification-loops → reviewer-override-flag.

Closed reviewer decisions the build must not reopen: uncapped verification
loops; light-path ticket file kept with inline implementation; own-family
reviewers via native agent surface (CLI only cross-family or when no native
surface exists); one bundled spec; transitional slicing-step trigger (worded
against the current falsification review, re-pointed to loop 1 by
verification-loops); drain implement-now lands as its own commit via a second
named batch-rule exception.

Two light-path commits already landed the enforcement-surface-read rule in
`bench-write-spec.md` step 1 — the build must preserve it through the step-9
rewrite. An A/B rerun prompt for that rule (Codex form) was handed to the
reviewer; if it ran, `capture/learnings.md` carries an "A/B rerun" entry for
the next drain.

## Next command

`/bench-implement-spec spec-authoring-and-light-path` — fresh mid-tier session.

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
