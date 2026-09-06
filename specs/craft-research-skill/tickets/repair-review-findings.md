# Repair the accepted review findings

Blocked by: point-the-callers-at-craft-research.md
Writes: .agents/commands/bench-assess.md, .agents/skills/bench-craft-research/SKILL.md, tests/canary/workflow-guidance-anchors/craft-research-read-side-boundary, .agents/skills/bench-craft-spec/references/bootstrap-authority.md, internal/anchors/registry_craft_research.go, internal/anchors/registry_craft_research_test.go, internal/conformance/registry_test.go, tests/canary/workflow-guidance-anchors/craft-research-bootstrap-probe-pointer (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: CR10, CR16, CR20

## What to build

Repair the five accepted findings in `reviews/craft-research-skill.md`: S1, S2, C2, F1, and F2.

Cut the duplicated unknowns sentence from step 3 of the assessment command. Trim the doc comment on the research anchor test to the independence sentence. Reword the skill's write-access sentence so it names a write delegate and a done-claim.

Reword the bootstrap-authority pointer so it points at the skill's probe rule. The pointer must not charge the skill with the probe. Update its needle in the registry and in the mutation table. Add the new registry file to the fixture owner list in the conformance registry test. Add one fixture for the bootstrap pointer under the `workflow-guidance-anchors` family.

## Acceptance

- [ ] [R1] `TestCraftResearchAnchorsRedOnRemoval` passes, and the reworded bootstrap needle reds on omission through `bench probe --check docs-currency-workflow`.
- [ ] [R2] `bench gate-prose` passes on the three edited Markdown files, and the skill stays at or under 122 lines.
- [ ] [R3] The fixture-bite owner test runs the new fixture green, and the fixture owner list names `internal/anchors/registry_craft_research.go`.
- [ ] [R4] Step 3 of the assessment command keeps its anchored phrase and states the unknowns rule zero times.
