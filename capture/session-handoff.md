# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5300bfa`, clean tree, 3 unpushed commits
Spec: `specs/ft164-ticket-contracts/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `5300bfa` (reduced run; full-tree ancestor evidence inherited)

## State

- **FT164 spec staged and signed off** at `specs/ft164-ticket-contracts/spec.md`
  (commit `5300bfa`): the ticket-contract core split, 7 stories, 24-needle
  anchor inventory, 30 coverage rows, falsification pass folded in. The
  reviewer approved the spec and its four flagged calls on 2026-08-02:
  ticket-local requirement IDs, basename-keyed `Blocked by:`, the
  process-boundary edge class landing in `craft-spec` plus the profile
  checklist, and the contracts-discovery spec-side half as one pointer
  sentence in `craft-spec`'s slicing section. These decisions stay closed.
  The repair riders (model-comparison, inventory-currency, shared-cache)
  stay on the roadmap under FT164's row for the next drain to re-home.
- **craft-research boundaries are closed** and recorded in
  `decisions/craft-research.md` (committed by its own session): read-side
  research discipline only, nothing carved from `craft-delegate`, no
  `craft-tickets` edits, no tier statements. FT164's implementation and the
  craft-research build are disjoint and land in either order.
- Implementation notes the FT164 build must honor (from the spec): the one
  `internal/specbuild` edit is an export-only parser entry point; anchors are
  section-scoped with each needle getting a mutation-table row; completion
  evidence is a fresh-session dogfood (cold-written ticket through
  `bench spec build assign` in a fixture repo), since the build touches the
  gate and forfeits `craft-synthesis`'s prose-only substitute.
- pcgs residuals ride their existing roadmap rows (FT183 review residuals,
  FT181); the 21 out-of-pool `../bench-pcgs*` worktrees remain cleanup-ready
  on reviewer say-so, per the prior handoff.
- Nothing has been pushed; push is the reviewer's call.

## Next command

`/bench-implement-spec` — the staged frontier drains before anything else:
`specs/pre-push-guard-visibility/spec.md` (FT135) first per the roadmap's
recommended sequence, then `specs/ft164-ticket-contracts/spec.md` in a fresh
mid-tier session.

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
