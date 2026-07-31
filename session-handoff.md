# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `977f772`, 1 dirty path, 44 unpushed commits
Spec: `specs/ft126-recurrence-tallying/spec.md` (Status: staged), `specs/ft128-agent-line-binding/spec.md` (Status: staged)
Gate: green at `977f772`; this handoff-only rewrite has not been re-gated

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
- **The diff-visual product decisions remain closed and its two approved HTML
  assets remain committed.** The report form and sections, three-drawn/two-
  relational edge kinds, dependency-free layout, declared context nodes,
  three-valued coverage badge, and justified heat are settled. The current map
  is honestly `Status: shaping` under the new schema because embedded-diff size,
  before-chain provenance, opt-in details, command/schema naming, and ad-hoc
  report placement have not been reclassified from `Not yet specified`.
  Preserve the settled product decisions when resolving that workflow state.
  The prior engineering outline was schema core (emitter and validator), then
  renderer, then final-check integration; current spec authoring owns the final
  seams and contracts.
- **Fresh-session dogfood and mid-binding three-axis review completed.** The
  concrete source-locator, cycle-diagnostic, and silent-shaping projection
  findings are repaired and gated. The remaining review observation is
  historical: two implementation commits crossed the spec's declared slice
  fences; no current-tree defect or rewrite is pending.
- **The last oracle was green for `977f772`.** This handoff-only rewrite does
  not re-gate capture. Ship-tier `bench prep-release` has not run, and nothing
  was pushed. The decision-map retro is committed and pending drain. FT126 and
  FT128 remain staged; the other active shaping maps and roadmap work were
  outside these phases.

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
