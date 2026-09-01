# Name the repair ticket owner before re-review

Blocked by: add-the-both-ends-rule-to-craft-gate.md
Writes: .agents/commands/bench-review-implementation.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors, tests/canary/workflow-guidance-anchors/review-repair-ticket-owner (new), tests/canary/workflow-guidance-anchors/review-repair-ticket-covers (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: KG4, KG5, KG6, KG7

## What to build

A coordinator reads one rule that keeps `rows-owned` green at the final check. The rule
ships as a sentence pair in the review phase file, two anchor tuples, one entry in a
registry test, and two live-mirror fixtures.

The sentence pair lands as one paragraph under `Review modes` in
`.agents/commands/bench-review-implementation.md`. The paragraph sits between the
re-review paragraph and the landing paragraph. The first sentence states that the
coordinator writes one repair ticket before the repair-scoped re-review, when accepted
repairs amend the coverage map. The mutation swaps `before` for `after`, so hold the word
`before` inside the needle.

The second sentence states that the ticket cites each amended row in `Covers:` and records
the accepted repairs. The mutation drops the `Covers:` clause, so hold that clause inside
the needle. The ticket template already carries `Covers:` and `Blocked by:`, so a repair
ticket is an ordinary ticket to the parser.

Insert the paragraph, and rewrite no existing sentence in the `Review modes` section. A
live-tree test, `TestReviewConvergenceContractCurrentDocs`, pins six normalized substrings
of this file. A reflow of the re-review paragraph breaks one of those substrings.

The rule has one home. Add no line to `.agents/commands/bench-implement-spec.md`, and add
no line to `.agents/skills/bench-craft-tickets/SKILL.md`. The Standards axis reads the
diff for a second copy of the rule.

Give each sentence one `RequireInSection` tuple over the section `Review modes` in
`internal/anchors/registry_data.go`. Both tuples take `AfterImplementSpec`, the group the
newest tuples in that file use. Keep each needle on one physical line, because a needle
that wraps across two lines never matches.

Add `TestRepairTicketOwnerAnchorsRedOnRemoval` to
`internal/anchors/registry_data_test.go`. The function writes each needle and each
diagnostic independently of the registry. It proves that each tuple reds when a synthetic
tree drops the needle. It also proves that the live root stays silent.

Add two fixtures under `tests/canary/workflow-guidance-anchors/` in the live-mirror shape.
`BASE` names the live file path `.agents/commands/bench-review-implementation.md`.
`MUTATE.json` holds one old-and-new needle swap. `EXPECT` holds the diagnostic the tuple
raises. Do not use the older `files/` snapshot shape for either fixture.

The review phase file carries no budget row, so no budget row moves for this ticket.

This ticket appends two tuples and one test function to the two registry files. It edits
nothing the previous ticket added, and the next ticket appends beside both sets.

## Acceptance

- [ ] KG4 — the fixture `review-repair-ticket-owner` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG5 — the fixture `review-repair-ticket-covers` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG6 — the diff carries the repair ticket rule in the review phase file alone.
- [ ] KG7 — `TestReviewConvergenceContractCurrentDocs` stays green on the live tree.
