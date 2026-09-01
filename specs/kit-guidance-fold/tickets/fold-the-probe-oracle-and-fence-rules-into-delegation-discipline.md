# Fold the probe oracle and fence rules into delegation discipline

Blocked by: make-the-falsification-pass-standing-for-guidance-diffs.md
Writes: .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-spec/references/map-discipline.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors, tests/canary/workflow-guidance-anchors/delegate-anchor-probe-owning-check (new), tests/canary/workflow-guidance-anchors/delegate-skip-ownership-check (new), tests/canary/workflow-guidance-anchors/delegate-root-conformance-pass (new), tests/canary/workflow-guidance-anchors/delegate-probe-mutated-bytes (new), tests/canary/workflow-guidance-anchors/delegate-serial-ceiling-fence (new), tests/canary/workflow-guidance-anchors/delegate-live-tree-inventory-fence (new), tests/canary/workflow-guidance-anchors/delegate-grammar-fence-inventory (new), tests/canary/workflow-guidance-anchors/delegate-out-of-fence-write (new), tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: KG18, KG19, KG20, KG21, KG22, KG23, KG24, KG25, KG26, KG27, KG28, KG29, KG42

## What to build

A coordinator writes a charge whose probe grades the live file, and whose fence holds
every file the ticket writes. The rules land in two reference files, and neither skill
file takes a net line.

Five bullets land under `In the charge` in
`.agents/skills/bench-craft-delegate/references/delegation-discipline.md`. The first
bullet requires an anchor-adding charge to name `bench test --check <owning-check>` as its
probe. The second bullet requires the capability-skip conformance check in the focused
checks of a charge that adds a test that can skip. The third bullet fences the
serial-census ceiling file, `internal/worktree/parallel_census_test.go`, for a charge that
binds `PATH` or the process environment. The
fourth bullet fences `internal/conformance/tier_test.go` for a charge that adds a
live-tree test.

The fifth bullet requires a grammar charge to enumerate the shared fixture owners and the
exact-record assertion families in its fence. Keep the term "exact-record assertion
families", because the charge-side rule needs a category name that does not rot.

Two bullets land under `Probes` in the same file. The first bullet states that
`bench test --package ./internal/conformance` is not the root conformance pass, and it
extends the existing live-tree probe bullet. The second bullet requires the coordinator to
confirm the mutated bytes against the copy aside, before the coordinator reads the
verdict.

The out-of-fence bullet takes an in-place generalization from "registry" to "write". The
generalized bullet requires the delegate to report an out-of-fence write before the
delegate edits. Write no second bullet on that fact, because a second bullet is a second
source.

The literal-bytes bullet under `Before the map locks` in
`.agents/skills/bench-craft-spec/references/map-discipline.md` takes an in-place
extension. The bullet now covers a deleted or a moved sentence, and it names `tests/` and
`internal/conformance`. The mutation drops `or moves`.

The mutations fix the needle words for the other rows. Row KG18 swaps the owning check for
the anchors package. Row KG19 drops the focused-checks clause. Row KG20 swaps `is not` for
`is`. Row KG21 drops the before-the-verdict clause. Row KG22 weakens `includes` to
`may include`, and row KG24 keeps only the fixture conjunct.

Hold each of those words inside its needle.

Leave the `Repair-charge template` section untouched, because FT164 owns it. The existing
test `TestRepairChargeTemplateAnchorsRedOnRemoval` pins its five field needles. Each rule
this ticket adds appears once across the two reference files, and the Standards axis reads
the diff for a second copy.

Give each new bullet one `RequireInSection` tuple in `internal/anchors/registry_data.go`.
Rows KG18, KG19, KG22, KG23, KG24, and KG25 scope to the section `In the charge`. Rows
KG20 and KG21 scope to the section `Probes`. Row KG27 scopes to the section
`Before the map locks`. Every tuple takes `AfterImplementSpec`, the group the newest
tuples in that file use.

Add `TestChargeProbeOracleAnchorsRedOnRemoval` to
`internal/anchors/registry_data_test.go`, and give it rows KG18 through KG25. Extend the
existing `TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval` with the row KG27
needle. Each function writes its needles and its diagnostics independently of the
registry. Each function proves a red on a synthetic tree that drops the needle, and
silence on the live root. Keep each fixture needle on one physical line, because the fixture materializer refuses an `old` value that spans a line wrap.

Add nine fixtures under `tests/canary/workflow-guidance-anchors/` in the live-mirror
shape. `BASE` names the live file path the tuple reads. `MUTATE.json` holds one
old-and-new needle swap. `EXPECT` holds the diagnostic the tuple raises. No fixture targets
`delegation-discipline.md` today, so these fixtures are its first readers. Do not use the
older `files/` snapshot shape for any fixture.

Both reference files carry no budget row. `.agents/skills/bench-craft-delegate/SKILL.md`
holds 122 lines at its exact row, and `.agents/skills/bench-craft-spec/SKILL.md` holds 151
lines under its row of 152. Add no line to either skill file, and move no budget row.

This ticket appends its tuples and its new test function to the two registry files. It
edits nothing the previous three tickets added, except the row KG27 needle inside the
existing map-discipline test function.

## Acceptance

- [ ] KG18 — the fixture `delegate-anchor-probe-owning-check` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG19 — the fixture `delegate-skip-ownership-check` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG20 — the fixture `delegate-root-conformance-pass` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG21 — the fixture `delegate-probe-mutated-bytes` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG22 — the fixture `delegate-serial-ceiling-fence` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG23 — the fixture `delegate-live-tree-inventory-fence` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG24 — the fixture `delegate-grammar-fence-inventory` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG25 — the fixture `delegate-out-of-fence-write` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG26 — `TestRepairChargeTemplateAnchorsRedOnRemoval` stays green on its five field needles.
- [ ] KG27 — `TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval` reds when a synthetic tree drops the moved-sentence needle, and stays silent on the live root.
- [ ] KG28 — `TestGuidanceProseBudgetsHoldOnTheLiveTree` stays green with `craft-delegate` SKILL.md at 122 lines or fewer.
- [ ] KG42 — `TestGuidanceProseBudgetsHoldOnTheLiveTree` stays green with `craft-spec` SKILL.md at 152 lines or fewer.
- [ ] KG29 — each rule this ticket adds appears once across the two reference files.
