# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `e027228`, clean tree, 19 unpushed commits
Spec: `specs/axi-aggregate-empty-migration/spec.md` (Status: staged), `specs/axi-bounded-projection-migration/spec.md` (Status: staged), `specs/axi-carriers-and-registry/spec.md` (Status: staged), `specs/axi-compatibility-oracle/spec.md` (Status: staged), `specs/axi-outcome-action-migration/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `03e184f` — current

## State

FT173's first spec, `specs/axi-compatibility-oracle/spec.md` (pin `8ae1512f`,
12 tickets), has an active `bench spec build` run: candidate `a6af041a`, four
tickets integrated (`authenticate-baseline-manifest`, `derive-root-registry-membership`,
`derive-wrapper-surface-membership`, `capture-pinned-baseline`), each with a
checkpointed receipt carrying delegate red→green logs plus an independent
coordinator mutation probe — keep that per-ticket probe discipline.
Ready frontier (serial: every fence shares `internal/axi/compatibility`):
`derive-nested-grammar-membership` next, then `close-required-argv-classes`,
`compare-four-observations`, `pin-default-full-and-empty-classes`,
`pin-truncation-bound-edges` + `pin-toon-byte-classes`, then the two hostile
tickets. Reviewer-set lines: sonnet/medium for the nested-grammar census and
both pin-* extension tickets; opus/high for the rest; write delegates only,
one at a time. Reviewer-ordered: before promotion, a cross-harness review of
the composed diff via `codex exec` on `gpt-5.6-terra`, high reasoning, yolo
approval posture, charged to refute "spec implemented". Standing decisions:
no committed manifest under testdata (seal binds the live tree; capture is
test-time from a `git archive` build of the pin) — flagged for reviewer veto
alongside the BC1 file-identity rewording landed at `e027228d`; orphaned
assignment `166ea6331ceb8e84bb0fc650735eaf75` awaits post-run
`bench spec build reclaim`; `specs/axi-bounded-projection-migration/tickets/contract-projection-routes.md`
still pins retired subject `974020e4` — repair when that build starts.

## Next command

`/bench-implement-spec`

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
