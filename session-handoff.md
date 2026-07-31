# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `ef0cf50`, 1 dirty path, 43 unpushed commits
Spec: `specs/ft126-recurrence-tallying/spec.md` (Status: staged), `specs/ft128-agent-line-binding/spec.md` (Status: staged)
Gate: green at `a640496` — stale, work tree `7a3fb15`

## State

- **Decision-map integrity and phase ownership is implemented and retired on
  `main` at `ef0cf50`.** The current format has one schema owner, graph,
  readiness, exact source-record and tree validation, a read-only five-column
  AXI query and template, ambient distinct-map counting, and a gate-bound
  49-fixture mutation family.
- **Workflow ownership is closed:** decision maps are situational and use
  decision tickets for reviewer choices; spec authoring accepts exactly one of
  three reviewed sources and owns engineering seams, tests, coverage, hostile
  inputs, and gate attachment. Implementation tickets remain independently
  green build units.
- **Fresh-session dogfood and mid-binding three-axis review completed.** The
  concrete source-locator, cycle-diagnostic, and silent-shaping projection
  findings are repaired and gated. The remaining review observation is
  historical: two implementation commits crossed the spec's declared slice
  fences; no current-tree defect or rewrite is pending.
- **The last oracle was green for the retirement tree.** The ambient gate now
  reads stale only because the required retro and this handoff were written
  after that boundary; final-check explicitly does not re-gate capture
  artifacts. Ship-tier `bench prep-release` has not run, and nothing was pushed.
- **Pending capture:** `.bench/retros/decision-map-integrity-and-phase-ownership.md`.
  FT126 and FT128 remain staged, and the remaining active shaping maps and
  roadmap work were outside this phase.

## Next command

`$bench-what-next`

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
