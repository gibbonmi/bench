# State the probe-only rule in the guidance

Blocked by: 01-accept-probe-only-plans-and-retire-check-floor.md
Writes: .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-spec/references/completion-plan.md, internal/anchors/registry_build_cost.go (new), internal/anchors/registry_build_cost_test.go (new), internal/anchors/registry_data.go
Covers: BC34, BC35, BC36, BC37

## What to build

The Build section of `.agents/commands/bench-implement-spec.md` states the probe-only rule with this sentence: "After the ticket commit, the author runs only the probes of the plan, because the chunk checkpoint grades every other check."
The section loses the clause "then runs each verification requirement once at that commit" and the sentence about the `check-floor` row. The file keeps its 97-line budget, because each sentence edit stays inside its paragraph.
Three other sentences name the floor, and each one changes:

- The "Declare the line" section names a red `check-floor` row as a route into the staleness pass. Remove that name, and keep the `completion-plan` route.
- The Land section states that a red `check-floor` row blocks the next charge and the review charge. Delete that sentence.
- The Land section names "any floor edit" before a repair dispatch. Delete that clause.

The completion-plan reference states: "Each chunk requirement names a `probe`, and a chunk can list no requirement." It loses its `check-floor` sentence.
Its template shows the `probe` key on the chunk requirement. Its parser sentence says that the parser refuses an empty final list, because an empty chunk list is now valid.

Create `internal/anchors/registry_build_cost.go` with the anchors of this spec, and append its slice to the registry in `registry_data.go`. This ticket adds the Require anchors of BC34 and BC36 and the Forbid anchors of BC35 and BC37. `TestBuildCostAnchors` pins each entry.

## Acceptance

- [ ] The Build section holds the probe-only sentence (BC34).
- [ ] `bench-implement-spec.md` holds neither the old verification clause nor `check-floor` (BC35).
- [ ] The completion-plan reference holds the probe-only sentence (BC36).
- [ ] The completion-plan reference holds no `check-floor` (BC37).
- [ ] `bench test --check docs-currency-workflow` passes on the edited files.
