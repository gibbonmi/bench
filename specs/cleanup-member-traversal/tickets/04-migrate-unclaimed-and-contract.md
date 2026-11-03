# Migrate unclaimed traversal

Blocked by: 03-migrate-landed-traversal.md
Writes: internal/worktree/clean_unclaimed.go, internal/worktree/clean_set.go, internal/worktree/clean_discard.go, internal/worktree/clean_landed.go, internal/worktree/cleanup_members_census_test.go (new), internal/worktree/clean_set_outcomes_test.go, internal/worktree/clean_set_wiring_test.go, internal/worktree/cleanup_members_baseline_test.go (new), internal/worktree/testdata/cleanup-member-traversal/ (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/parallel_census_test.go, internal/worktree/cleanup_members_modes_test.go, internal/worktree/cleanup_members_baseline_record_test.go
Covers: R01, R02, R04, R05, R14, R18, R19, R20, R21, R24, R27, R28

## What to build

Deliver the first half of C4: unclaimed cleanup consumes the accepted C1 traversal.
Remove applyUnclaimedAssignmentSet's original loop in this ticket.
Consume the predecessor baseline records and completed adapter migrations.
Ticket 00 owns the frozen records, and this ticket adds only the post-replan unclaimed witnesses.
applyUnclaimedAssignmentSet keeps its current signature with no run parameter.

Keep StepUnlockedReplan before full fingerprint and length replan.
Keep each tuple equality check immediately before that member's retained or deletion visit.
Keep direct exact-OID deletion without recorded-worktree authority.
After deletion failure return only the completed prefix and failed row.
Add no production lifecycle hook and append no unreached suffix.

Implement the accepted post-replan witnesses from the spec.
A ref lock preserves successful selection replan but forces first or later deletion failure.
A real reference-transaction hook moves a later ref after an earlier deletion commits.
The exact-OID guard must preserve the moved ref's new tip.
The hook records execution and ignores its nested update.

Ticket 05 owns the census, the fingerprint input owner, and the full R26 comparison.
Each mode-agnostic row that this ticket covers needs one witness that drives applyUnclaimedAssignmentSet.

The mandatory registry paths preserve current dispatch, help, AXI membership, and routing.
This ticket introduces no command or rendered schema change.

## Acceptance

- [ ] Empty, retained-only, mixed, and successful unclaimed cases preserve ordered rows and skip retained qualification and effects (R01, R02, R04, R05).
- [ ] TestCleanupUnclaimedSelectionDrift refuses additions, removals, and tuple changes before deletion (R20).
- [ ] TestCleanupUnclaimedDeletionFailureAfterReplan proves successful full replan before first and later deletion faults (R18).
- [ ] Those fault cases return exactly the completed prefix and failed row, with no suffix rows or rollback (R14, R18, R19).
- [ ] TestCleanupUnclaimedMovedRefAfterReplan proves the hook ran and preserves the moved later ref while the earlier deletion remains complete (R21).
- [ ] The failed and unvisited refs retain their expected OIDs, including head, tip, and checked-out-ref preservation cases.
- [ ] Unclaimed stale remedies, original boundary trace, and authority refusals match frozen responses and durable state (R24, R27, R28).
- [ ] Changed unclaimed suffix projection, an omitted old-OID guard, and a continued visit after failure produce observed behavioral reds.

## Checkpoint verification

Run C4's verification commands from the spec's completion plan.
The focused producer checks run at this checkpoint while successors remain unbuilt.
Record the named red, restoration, green, and baseline comparison in the review pickup.
Use bench probe for each behavioral omission or swap after implementation.
A compile failure does not prove a preservation property.

Ticket 05 is the successor in C4, and its author freezes the chunk delta.
