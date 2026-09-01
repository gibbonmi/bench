# Add the both-ends rule to craft-gate

Blocked by: none
Writes: .agents/skills/bench-craft-gate/SKILL.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors, tests/canary/workflow-guidance-anchors/craft-gate-indirected-value-both-ends (new), tests/canary/workflow-guidance-anchors/craft-gate-single-edit-defeat (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: KG1, KG2, KG3

## What to build

A gate author reads two new rules in `craft-gate`, at the two places the author already
opens. Each rule ships as a sentence in the skill file, one anchor tuple, one entry in a
registry test, and one live-mirror fixture.

Rule (a) lands as a standalone sentence at the end of the first paragraph under
`Run the real path`. The sentence names three indirected values: a workflow output, a
config key, and an environment variable. It states that the check grades the producer and
the consumer, or their binding, in the same change. The mutation deletes the binding
clause, so hold that clause inside the needle.

Rule (b) lands as a standalone sentence at the end of the second paragraph under
`Prove it bites`. The sentence states that the author asks which single edit defeats a new
check while the gate stays green. The mutation swaps `single edit` for `change`, so hold
the words `single edit` inside the needle.

Give rule (a) one `RequireInSection` tuple over the section `Run the real path` in
`internal/anchors/registry_data.go`. Give rule (b) one `RequireInSection` tuple over the
section `Prove it bites`. Both tuples take `AfterImplementSpec`, the group the newest
tuples in that file use. Keep each needle on one physical line, because a needle that
wraps across two lines never matches.

Add `TestCraftGateBothEndsAnchorsRedOnRemoval` to
`internal/anchors/registry_data_test.go`. The function writes each needle and each
diagnostic independently of the registry. It proves that each tuple reds when a synthetic
tree drops the needle. It also proves that the live root stays silent.

Add two fixtures under `tests/canary/workflow-guidance-anchors/` in the live-mirror shape.
`BASE` names the live file path `.agents/skills/bench-craft-gate/SKILL.md`. `MUTATE.json`
holds one old-and-new needle swap. `EXPECT` holds the diagnostic the tuple raises. Do not
use the older `files/` snapshot shape for either fixture.

The skill file holds 118 lines today and lands at exactly 120. Each rule costs one net
line after a reflow of the ragged break in the `Run the real path` paragraph. The glob row
of 120 in `projects/benchkit.md` holds, and no budget row moves for this ticket.

This ticket appends two tuples and one test function to the two registry files. The next
ticket appends beside them, and it edits nothing this ticket added.

## Acceptance

- [ ] KG1 — the fixture `craft-gate-indirected-value-both-ends` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG2 — the fixture `craft-gate-single-edit-defeat` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG3 — `TestGuidanceProseBudgetsHoldOnTheLiveTree` stays green with `craft-gate` SKILL.md at 120 lines.
