# Capture the pre-change cleanup baseline

Blocked by: none
Writes: internal/worktree/cleanup_members_baseline_test.go (new), internal/worktree/testdata/cleanup-member-traversal/ (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: none

## What to build

Deliver the first half of C1: a frozen baseline of the original cleanup behavior for all four modes.
This ticket changes no production file.
It runs before any traversal migration, so each record comes from an original run.

Write the case manifest that the spec's Testing decisions section enumerates.
Give each applicable mode its common cases: empty, retained-only, mixed, all-success, first failure, and later failure.
Add the recorded and landed outer-preflight and locked-drift cases.
Add the unrecorded holder-removal and recovery-ref cases.
Add the unclaimed full-selection drift case.
Add one landed missing-tree release row, which landed preflight skips and the landed visits apply.

Capture each case's response rows and durable-state outcome once, and store them under the testdata directory.
Obey the spec's identity and digest mapping section for normalization.
Write TestCleanupTraversalBaseline as the permanent comparison that drives the current adapters against those records.
The comparison passes on the unchanged production code.
No ordinary test restores source or rewrites the records.

Use the existing fixture owners in clean_set_apply_test.go and the existing unclaimed and discarded-ref helpers.
Do not paste their repository construction into the new file.
The post-replan unclaimed witnesses of R18 and R21 stay with ticket 04.

The mandatory registry paths preserve current dispatch, help, AXI membership, and routing.
This ticket introduces no command or rendered schema change.

## Acceptance

- [ ] The manifest names each mode and case above, with no unchecked Cartesian product.
- [ ] TestCleanupTraversalBaseline passes on the unchanged production code and reads only the frozen records.
- [ ] The normalization keeps actions, row order, digest relations, errors, recovery refs, and deletion outcomes meaningful.
- [ ] A dropped unreached suffix in the explicit stop projection produces an observed baseline red before restoration.

## Checkpoint verification

Run this ticket's requirement from the spec's completion plan.
Record the named red, restoration, and green in the review pickup.
Use bench probe for the behavioral mutation after the ticket commit.
A compile failure does not prove a preservation property.

Ticket 01 is the successor in C1, and its author freezes the chunk delta.
