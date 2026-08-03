# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `7546e54` before the roadmap-drain commit
Spec: `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: stale green before this reduced-scope roadmap-and-capture batch

## State

- The 2026-08-03 roadmap drain removed FT183 after `bench spec history check-level-conformance-scoping` confirmed its retirement at `7546e54`; no live spec was retired. Six ideas, the tracked-file set-aside learning, and the check-level-conformance-scoping retro are fully dispositioned; both capture inboxes and the retro queue are empty.
- FT185 owns structured `bench gate` output as an independent spec, preserving the gate-pipeline map's closed ruling that no output redesign rides the pipeline build. FT186 owns the two `internal/gate` structural refactors.
- FT180 now owns spec-optional `$bench-implement-spec --full`; FT174 owns the complete real-ticket dependency, ownership, Red-mutations, and acceptance-ID grammar; FT98 owns the recurring tracked-file set-aside gap.
- The retro's terminal-record work is in FT162, its repair-slice rule is in FT164, and its public-CLI cross-harness review rule is in FT158.
- Decisions that stay closed: all rulings in the six active decision maps, including the 2026-08-03 gate-pipeline and gate-structure amendments.

## Next command

`$bench-shape-idea`

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
