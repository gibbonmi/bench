# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `209b4d6` at handoff write; commit staging `specs/covers-traceability` lands immediately after
Spec: `specs/covers-traceability/spec.md` (Status: staged, reviewer-approved 2026-08-05), `specs/ft187-communication-surface-cut/spec.md` (staged), `specs/go-build-cache-footprint/spec.md` (staged), `specs/pre-push-guard-visibility/spec.md` (staged)
Gate: stale (reduced-scope drift) — the staging commit's gate run refreshes it

## State

**Phase reached: covers-traceability spec staged and signed off; implementation starts in a fresh mid-tier session.**

The spec adds coverage-map row IDs, per-row `covers` annotations in tickets,
assign refusals, and a pre-gate promote totality check computed over the
integrated assignments' digest-verified tickets. A mid-tier falsification pass
returned six findings; all are folded in (decoy-ticket degenerate,
digest-mismatch refusal, example-agreement covers literal, corrected checker
red, range rows refuse as unannotated, pre-gate ordering row).
`bench coverage --check` is green at 15 rows. Closed decisions: 1:1 row
mapping, spec-opt-in rollout, `covers: local` accepted at assign and graded by
review, inline row grammar. A thin-slice slicing-gradient flip for
`craft-tickets` is parked in `capture/IDEAS.md`, deliberately after this build
produces rows-per-ticket data.

Separate workflow, untouched here: lifecycle run FT195
(`go-build-cache-footprint`) has repairs integrated and a fresh exact-candidate
review pending; its retro should name ticket-shape/checkpoint evidence as the
next priority — which this spec now is.

## Next command

`/bench-implement-spec specs/covers-traceability/spec.md`

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
