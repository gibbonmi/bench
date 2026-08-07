# Add the breakdown-review moment to the pre-build pipeline

Blocked by: none
Ownership fence: `.agents/commands/bench-implement-spec.md`, `.agents/skills/bench-craft-tickets/SKILL.md`, `projects/benchkit.md`, `internal/anchors/registry_data.go`
Integration surfaces: anchors Require rows→`internal/anchors/registry_data.go`; conformance grading→existing internal/conformance TestRootConformance path exercised by AB1; registry append tail→add-retro-repair-attribution.md
Contracts: none crosses — every needle and its Require row land inside this ticket's own fence; sibling tickets append disjoint registry rows, sequenced by their `Blocked by:` chain
Closure: AB1/br1-needle, AB2/section-heading, AB2/witness-needle, AB3/routing-row

## What to build

The coordinator gets a named pre-assignment review moment: `/bench-implement-spec`
gains a step between ticket derivation and `bench spec build start`,
`craft-tickets` gains the section that owns the charge, and the benchkit profile
caches the pass's routing. Grouped as one ticket because the three edits are one
story's single capability (spec story 1) and share the registry fence with every
sibling; no thinner cut strands a red — the grouping is the spec's per-story
slice, and the sequencing cost of thinner cuts buys no independent value.

Each needle is an exact quoted string from the spec's coverage map. Land each
anchors `Require` row first, observe the missing-anchor conformance red, then
land the prose — per the map's red signals.

- BR1 (`.agents/commands/bench-implement-spec.md`, in `## First derive the tickets`
  between ticket derivation and lifecycle start): needle
  `one fresh read-only delegate grades the ticket breakdown before any assignment; a harness that cannot spawn one runs the pass inline and flags it in the build plan`.
  Findings are reslices repaired before `bench spec build start`.
- BR2 (`.agents/skills/bench-craft-tickets/SKILL.md`, new final section after
  `## Write one file per ticket`): heading needle
  `## Review the breakdown before assignment` and witness needle
  `every named producer path exists in the tree and holds the named value`.
  The section names the moment (ticket files written, nothing assigned), the
  owner (one fresh read-only delegate), and the charge pointing by name at the
  six per-field obligations the skill already states: the `Integration surfaces:`
  `none` re-search, the keep-together split attempt, label-shaped rows and
  closure tokens, mutation subject-not-assertion honesty, `covers` honesty, and
  fence-versus-advertisement — plus the one new item, the producer-path witness.
  When the harness cannot spawn a delegate, `craft-delegate`'s capability-aware
  policy applies and the pass runs inline, flagged in the build plan.
- BR3 (`projects/benchkit.md` Lines section): needle `Ticket-breakdown review pass`
  as a cached-routing bullet — mid model, medium effort, one iteration,
  read-only, standing grant like the falsification pass. Mirror the existing
  "Spec falsification pass" bullet's structure; the binding table is untouched.

New registry rows: append at the true end of the registry slice in
`internal/anchors/registry_data.go` with no `Group:` key (defaults to
`BeforeStructured`), mirroring the file's ungrouped rows. Diagnostic format:
`<file> missing acceptance coverage anchor: <needle>`.

Semantic exception (explicit, per craft-tickets): the six-obligation enumeration
inside BR2's charge and the routing row's field values are review-graded prose;
the gate proves only the quoted needles.

## Acceptance

- [ ] [AB1] (covers BR1) `/bench-implement-spec` carries the BR1 needle between ticket derivation and lifecycle start, with its Require row in the registry.
- [ ] [AB2] (covers BR2) `craft-tickets` carries the `## Review the breakdown before assignment` section with the witness needle and the six-obligation charge, with both Require rows in the registry.
- [ ] [AB3] (covers BR3) `projects/benchkit.md` Lines carries the `Ticket-breakdown review pass` cached-routing bullet, with its Require row in the registry.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AB1/br1-needle | delete the BR1 needle sentence from `.agents/commands/bench-implement-spec.md` (Require row retained) | conformance workflow-anchors check | `BENCH_CONFORMANCE_ROOT=$PWD go test -count=1 ./internal/conformance -run '^TestRootConformance$'`, expect `.agents/commands/bench-implement-spec.md missing acceptance coverage anchor: one fresh read-only delegate grades…` |
| AB2/section-heading | delete the `## Review the breakdown before assignment` heading from the skill (Require row retained) | conformance workflow-anchors check | same command, expect the heading's missing-anchor diagnostic |
| AB2/witness-needle | delete the producer-path witness sentence from the skill (Require row retained) | conformance workflow-anchors check | same command, expect the witness needle's missing-anchor diagnostic |
| AB3/routing-row | delete the `Ticket-breakdown review pass` bullet from `projects/benchkit.md` (Require row retained) | conformance workflow-anchors check | same command, expect the profile's missing-anchor diagnostic |
