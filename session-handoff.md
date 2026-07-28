# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean tree apart from this commit, ahead of origin
Spec: none active — `specs/ft91-gate-phase-split.md` (Status: implemented) is
held unretired on purpose; its stories-4/5/9 ruling now rides map ticket #4
Gate: green for all code at `64635b8`; the pin reads stale only because
doc-only commits landed after it

## State

- **FT91's eighth arm is mapped: `decisions/gate-critical-path.md`, opened
  2026-07-28.** Five tickets. The frontier is two agent-alone Research
  tickets: #1 diagnoses why the gate absorbed only 24 s of the seventh arm's
  131 s suite win (phase-timeline capture; names the current critical path and
  the artifact-suite inflation mechanism), #2 inventories the artifact suite
  by build-vs-inspect. The grills hang behind them: #3 prepared-artifact
  sharing (blocked by #2), #4 gate-pipeline reopen plus the unretired spec's
  stories 4/5/9 (blocked by #1), #5 the closing ticket — FT91's stop
  condition and the arm selection (blocked by #1, #3, #4).
- **Two scope rulings taken at bootstrap (reviewer, 2026-07-28):** the
  stories-4/5/9 spec-retirement ruling folds into ticket #4, and FT91's stop
  condition is decided in-map against #1's evidence rather than set now.
- **`decisions/cost-follows-project-size.md` stays open by design.** Ticket #6
  (cheap-tier retest, opportunistic Task) and the parked `-count=1` decision
  keep its `bench status` row; unrelated to this map.
- **Both inboxes are empty; the drain closed 2026-07-28** and the roadmap
  parses clean.

## Next command

`/bench-shape-idea` resume on `decisions/gate-critical-path.md` ticket #1 —
Research, agent-alone: capture the gate phase timeline and diagnose the
parallelism gap. Ticket #2 (artifact-suite inventory) is likewise unblocked
and independent, so either order works; #1 first, since #4 and #5 both grill
against its numbers.

Behind the map: `/bench-write-spec` for FT98 (the one preserve-then-discard
primitive), then FT71 (versioned local shift evidence, the remaining HIGH
bank-track row).

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
