# Document the test-only confirmation

Blocked by: 04-write-confirmations-with-bench-record.md
Writes: .agents/skills/bench-craft-line/references/bounded-repair-policy.md, .agents/commands/bench-review-implementation.md, .bench/BENCH-reference.md, projects/benchkit.md, CHANGELOG.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/conformance/implementation_continuation_test.go, internal/conformance/retained_workflow_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RC46, RC47

## What to build

The guidance states the test-only exception after the form exists.

- The bounded repair policy states the test-only exception beside the current-results rule. It names the confirmation's finding list as the coordinator's disposition of those findings.
- The review phase keeps its confirming-round sentence. A new sentence after it states that a test-only repair whose coordinator probe bites closes through `bench record confirmation` and takes no confirming round.
- The review phase's record step names the new form. Add a separate sentence and leave the pinned step 6 sentence that starts "Write each completed chunk, verification, review, and amendment entry" unchanged.
- The reference's checkpoint sentence and the project profile's completion-checkpoint seam name the confirmation.
- The changelog gains one `[Unreleased]` entry for the form.
- The implementation phase file stays unchanged, because it sits at its exact prose budget.

Put the new anchor rows in `internal/anchors/registry_chunk_chain.go`, because the retained-workflow registry has no headroom left after ticket 01.
A policy row keeps the `implementation continuation: ` prefix, so its twin in `wantPredicates` still selects it.
A review-phase row takes its twin in the chunk-chain anchor test.
Record a demonstrated red for each twin.
Run `bench preflight build repair-confirmation --propose-writes` for this ticket, and add each listed canary fixture directory to this `Writes:` line.

## Acceptance

- [ ] RC46: the policy states the test-only exception and names the confirmation's finding list as the coordinator's disposition. The anchors and twins red when either sentence is absent.
- [ ] RC47: the review phase states that a test-only repair with a biting coordinator probe takes no confirming round. The anchor and twin red when the sentence is absent.
- [ ] The prose-mechanics check passes on every edited guidance file.
