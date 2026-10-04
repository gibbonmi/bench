# Close verified spec delivery in its publication

Blocked by: 05-authorize-current-publication.md
Writes: internal/commitment (new), internal/intent, internal/worktree, internal/landing, internal/roadmap, internal/spec, internal/gate/completion.go, internal/gate/completion_test.go, internal/gate/commitment_completion_test.go (new), docs/adr/0015-the-landing-verb-is-the-one-author-of-the-spec-flip.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DC34, DC35, DC36, DC37, DC38, DC39, DC40, DC41, DC42, DC75, DC76

## What to build

Extend the broker's verified spec completion transform to close exactly the fulfilled bound obligations in the same candidate tree.
Ticket 05 supplies current publication authority. Existing source-bound completion and review records supply delivery evidence.
The roadmap owner removes the row, its detail owner, every completed sequence reference, and fulfilled dependency references. Preserve unrelated and residual obligations.

Compute one exact allowed transform. The completion oracle compares the candidate with that transform, including spec status, policy delivery facts, roadmap, and sequence effects.
A list of allowed paths is insufficient. Published facts identify the reviewed source and evidence without a circular reference to their own commit.
Update the durable landing ADR for this resulting ownership contract.

Keep persistence failures before the oracle, interruption within it, and post-publication local failures distinguishable. Resume preserves the original publication and does not repeat it.
Use the existing publisher and recovery fixture seams. No independent closure commit or later drain is part of successful delivery.

Read the existing completion transform, completion oracle, roadmap row and sequence parsers, and source-bound review record contract. This is the largest ticket; its fresh context is limited to that one publication journey. Extract large owner files into focused siblings within their declared directories before growing them.

## Acceptance

- [ ] Verified A closes its row and detail, and every sequence reference, in the same published commit (DC34, DC35).
- [ ] A partial delivery retains its residual obligation and unrelated B (DC36, DC39).
- [ ] A red gate or interrupted oracle publishes neither delivery nor closure (DC37, DC41).
- [ ] Omitting one sequence removal fails the exact-transform check; demonstrate that mutation (DC38).
- [ ] A pre-oracle persistence failure leaves authority and destination unchanged (DC40).
- [ ] A terminal local-write failure resumes the original publication once (DC42).
- [ ] Verified delivery removes satisfied dependency references from the board dependency tables (DC75).
- [ ] A listed legacy run closes its delivered scope, and a partly delivered scope stays open (DC76).

## Checkpoint verification

Run `bench test --package ./internal/worktree`, `bench test --package ./internal/landing`, `bench test --package ./internal/gate`, `bench test --package ./internal/roadmap`, `bench test --package ./internal/commitment`, and `bench test --package ./internal/commitment/repository`. Record the sequence-omission red and restored green. Tests must compare the real composed tree.
