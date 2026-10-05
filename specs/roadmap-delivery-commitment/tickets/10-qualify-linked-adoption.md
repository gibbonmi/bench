# Qualify installed adoption and prepare the finite milestone

Blocked by: 09-project-commitment-guidance.md
Writes: internal/systemtest/adoption_test.go, internal/systemtest/owner_landing_fixture_test.go, internal/systemtest/land_route_test.go, internal/commitment/repository/repository.go, internal/commitment/commitcmd/command.go, internal/commitment/store_test.go, internal/roadmap/tree.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, specs/roadmap-delivery-commitment/adoption-proposal.md (new), reviews/roadmap-delivery-commitment.md, specs/shared-test-fixtures
Covers: DC51, DC52, DC53, DC63, DC85, DC86

## What to build

Extend the existing installed-owner adoption journey to prove the complete prerequisite in a linked project.
Ticket 09 supplies all commands and reader surfaces. Use BENCH_KIT through the existing sealed system test runner and disposable kit-install fixture.
Show that the prerequisite can publish before its policy exists. After installation, planning and explicit adoption enable the approved delivery; absent adoption still refuses new delivery.

Remove the candidate's admission call in the system fixture and attempt displaced publication. The trusted installed broker must refuse it before publication.
Restore the call and retain the green journey. Do not create a second nested test runner or system fixture framework.

Prepare an adoption proposal for the reviewer after the prerequisite publishes. Verify each named quality owner against delivered history and retain only its remaining obligation.
The finite set is FT376, FT373, and FT349, in that order. Do not reopen already delivered work.
Record exact source identities, criteria, ordered remaining outcomes, and any approved legacy continuation evidence. Findings outside this set remain uncommitted.

The proposal is a review artifact, not an active policy. Include replayable inventory, plan, and approve commands for the newly installed version.
Actual kit adoption follows prerequisite publication and explicit reviewer direction. This ticket must not write `.bench/commitment.json` or grant itself authority.
Record final coverage and review pickup in the existing feature review record.

Approval stages the policy alone in a project with no board. The inventory reports the identity that the plan binds and each run's bound deliverable or scope.

Read adoption_test.go, its installed-owner fixture helper, the reviewed wrapper land_route contract, and the three named roadmap owners. Reuse existing fixture callers. These two system files own the journey and its helper; add any newly discovered relocation destination to this ticket before use.

## Acceptance

- [ ] The installed wrapper in a linked project refuses an uncommitted production start (DC51).
- [ ] The prerequisite installs before policy adoption, after which only explicitly approved adoption admits delivery (DC52).
- [ ] The proposal contains only verified remaining obligations from the exact named set, with source evidence for each inclusion or omission (DC53).
- [ ] A candidate that removes its guard still cannot publish displaced work through the installed broker; demonstrate the red-capable omission probe (DC63).
- [ ] The proposal identifies explicit post-publication adoption commands and does not activate policy during the build.
- [ ] Approval in a linked project with no board stages the policy and writes no ROADMAP.md (DC85).
- [ ] Inventory identities feed the plan input unchanged, and each run row names its deliverable or scope (DC86).

## Checkpoint verification

Run `bench test --check system` with BENCH_KIT supplied by that verb. Run `bench test --package ./internal/commitment`, `bench test --package ./internal/commitment/repository`, `bench test --package ./internal/roadmap`, and `bench test --package ./cmd/bench`. Use the existing adoption journey. Review DC53 against decision 14 and the source delivery history; record that manual check with exact references.
