# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `acf5c9d`, 3 dirty paths, 36 unpushed commits
Spec: `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `7c61ed8` — stale, work tree `d9a0656`

## State

**Phase reached: gate-decision-test-seam promoted, retired, and retro captured.**

The spec-build run `527a22c5` is terminal: promotion commit `3cce8aa` published
candidate `f7f0ea8c` and the spec-retire commit `acf5c9d` removed
`specs/gate-decision-test-seam/` and its ROADMAP row. The interrupted-gate
incident is resolved — plain `bench gate` (never `--fresh`) recovered the
pending record; the original `ComposedGreen` refusal cause is recorded as
unprovable (overwritten record, no transcript copy). Review receipt
`80069545…` holds seven risk-accepted findings as the veto surface.

Uncommitted capture awaiting the drain: `capture/IDEAS.md` (four new parked
ideas: refusal-advice differentiation, pending-record preservation,
absolute-receipt-path error, check_slots_test follow-ups),
`capture/learnings.md` (one open 2026-08-07 entry on the `--fresh` misroute),
`capture/retros/gate-decision-test-seam.md`, and this file. FT135
(`specs/pre-push-guard-visibility/spec.md`) is the one staged spec.

Closed decisions stay closed: risk-accepted review dispositions were taken
under the reviewer's standing batch approval and reopen only by reviewer veto;
the decision-seam matrix's reload-coverage narrowing is the spec's own recorded
trade.

## Next command

`/bench-what-next` — drain the ideas, the open learning, and the pending retro.

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
