# Fold the retained-workflow check into the registry

Blocked by: none
Writes: internal/anchors/registry_retained_workflow.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/conformance/retained_workflow_test.go, internal/conformance/implementation_continuation_test.go, internal/conformance/docs_workflow_checks_test.go, internal/conformance/tier_live_tree_test.go, cmd/bench/testdata/anchors/fixture-repo/AGENTS.md, cmd/bench/anchor_help_test.go, tests/canary/workflow-guidance-anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: AP02, AP03, AP04, AP10, AP25, AP26, AP27, AP31, AP46

## What to build

The retained-workflow check stops being a second registry. It reads its expected diagnostics from the main registry, by its group and its diagnostic prefix. Each stale phrase that it forbids becomes a registry Forbid row on the same file. Its hard-coded model identifiers leave, because the line-binding check already compares the profile prose with the binding.

The retained-workflow registry file keeps only its admitted rows: every Forbid row, the one INVARIANT row on `.bench/BENCH.md`, and the COMMAND rows. Each removed row leaves with its test rule and its entries in the implementation-continuation pin list. Its fixture directories, which `assets/fixture-census.md` lists for this ticket, leave too. Keep each constant that a row in another registry file still reads.

This ticket runs first, because the check hard-codes diagnostics of rows in the core and topic registries. The fixture repository of the anchors help test plants three `AGENTS.md` needles that this ticket removes, so update that fixture in the same commit. Re-cut each COMMAND needle in this file that is over `prose.MaxSentenceWords` under the `internal/prose` word count. Put any new test in an existing test file, because `internal/conformance` has no file-count headroom.

## Acceptance

- [ ] The retained-workflow check holds no diagnostic literal that the registry also holds, and a test proves it from `Entries()`.
- [ ] The retained-workflow test derives its expected row count from `Entries()`.
- [ ] The retained-workflow check holds no model identifier, and the ticket record shows the identifier list gone.
- [ ] Each stale phrase that the check forbade at the ticket base is a registry Forbid row on the same file.
- [ ] The ticket record carries a removal census: each removed row is DROP in `assets/classified.md`, and each stale phrase is paired with its new Forbid row.
- [ ] `bench test --check docs-currency-workflow`, `bench test --package ./internal/conformance`, `bench test --package ./internal/anchors`, and `bench test --package ./cmd/bench` pass.
