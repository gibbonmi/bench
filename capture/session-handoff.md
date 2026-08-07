# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `/home/devuser/workspace/bench`
Branch: `main` — pre-commit HEAD `c8ea95ea1a9f`, 51 unpushed commits
Spec: `specs/pre-push-guard-visibility/spec.md` (Status: staged; no active build run)
Gate: unavailable on the pre-commit tree; this handoff lands only through the roadmap batch's green `bench commit`

## State

**Phase reached: roadmap reconciled and capture drained.**

`conformance-harness-scope` and `gate-decision-test-seam` are promoted and
retired. Their recommendations are folded into the existing lifecycle,
review-evidence, ticket-grammar, gate-attribution, and residual-verdict owners;
FT199 carries the distinct lifecycle-native repair-entry gap. FT135 remains the
sole staged spec. The ideas inbox, learnings journal, occurrence ledger, and
retros directory are empty.

Closed decisions stay closed: the gate-decision slice preserves the read-only
decision seam and representative full-engine composition proofs; the
conformance-harness slice exports the registry-resolved fixture check and keeps
representative shell-ordering journeys; pre-push guard visibility still follows
its compiled decision map. The refreshed sequence drains FT135 before new spec
authoring.

## Next command

`$bench-implement-spec`

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
