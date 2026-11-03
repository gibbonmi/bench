# Contract to one traversal owner and one fingerprint input owner

Blocked by: 04-migrate-unclaimed-and-contract.md
Writes: internal/worktree/clean_set.go, internal/worktree/clean_discard.go, internal/worktree/clean_landed.go, internal/worktree/clean_unclaimed.go, internal/worktree/cleanup_members.go, internal/worktree/cleanup_members_census_test.go (new), internal/worktree/cleanup_members_baseline_test.go, internal/worktree/cleanup_members_baseline_record_test.go, internal/worktree/cleanup_members_modes_test.go, internal/worktree/parallel_census_test.go, internal/worktree/testdata/cleanup-member-traversal/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: R25, R26

## What to build

Deliver the second half of C4: the complete tree has one member-traversal owner and one owner of each fingerprint input.
This ticket changes no cleanup behavior.

Write TestCleanupTraversalCopySurvival in the new census test file.
It reads the four adapter bodies and proves that each delegates its visits to the shared owner.
It must fail if any original member-iteration loop remains.
A harmless identifier check alone is insufficient.
Remove any remaining obsolete traversal support without retiring public entry points.

The explicit and landed fingerprints each build the options input, and each repeats one hash loop.
Give the options input and the hash loop one production owner, and make both fingerprints use it.
If the unclaimed baseline holds a copy of the unclaimed fingerprint inputs, move that list into production by the same route.
The digest bytes and the frozen records must not change.
The baseline renders each production list and checks each digest in the same run.

Complete the R26 comparison over every manifest case.
A changed digest relation, input omission, ref tip, or recovery destination must fail the comparison in each mode.

The mandatory registry paths preserve current dispatch, help, AXI membership, and routing.
This ticket introduces no command or rendered schema change.

## Acceptance

- [ ] TestCleanupTraversalCopySurvival reds on an original loop restored at each of the four adapter sites, then passes after each restoration (R25).
- [ ] One production owner builds the options input and runs the hash loop for each fingerprint, and no digest byte changes.
- [ ] TestCleanupTraversalBaseline compares every manifest case with meaningful digest relations, OIDs, recovery destinations, output rows, and durable state (R26).
- [ ] An omitted unclaimed fingerprint input produces an observed baseline red before restoration.

## Checkpoint verification

Run C4's verification commands from the spec's completion plan.
Record the named red, restoration, and green in the review pickup.
Use bench probe for each behavioral omission or restored loop after implementation.
A compile failure does not prove a preservation property.

This ticket's author freezes C4's delta for independent review.
