# Prune the core registry to the admitted rows

Blocked by: 1-fold-retained-workflow.md
Writes: internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_craft_research_test.go, internal/anchors/registry_retained_workflow.go, internal/conformance/ft311_preparation_test.go, internal/conformance/docs_workflow_helpers_test.go, internal/conformance/fixture_bite_test.go, internal/conformance/retained_workflow_test.go, cmd/bench/testdata/anchors/fixture-repo/AGENTS.md, cmd/bench/anchor_help_test.go, cmd/bench/anchors_dir_test.go, tests/canary/workflow-guidance-anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: AP02, AP03, AP04, AP05, AP06, AP07, AP08, AP10, AP18, AP20, AP21, AP22, AP23, AP24, AP46, AP47

## What to build

The core registry file keeps only its admitted rows: Forbid, LABEL, ROUTE, COMMAND, and the 8 core INVARIANT rows on `.bench/BENCH.md`. The Roles marker row also leaves, because the shared-rule check is its single source. Each removed row leaves together with its test rule and the fixture directories that `assets/fixture-census.md` lists for this ticket.

Add one INVARIANT row for each of invariants 1, 2, and 3, with these needles:

- `Never edit, skip, weaken, or delete a test or a check to make it pass.`
- `Before any multi-cycle stage, state the model, the effort, the iteration cap, and a one-clause justification in one line.`
- `Put no file paths and no code snippets in an ADR; they rot.`

Each new row gets one test rule that proves a deletion of its clause raises its diagnostic. The integration-source stale-phrase check loses its file list with its Require rows. So move each of its stale phrases into the registry as a `ForbidCaseFoldedEmphasis` row on each file that it scanned. That kind folds case, as the check did. Then remove the check and its count pin.

The check also stripped backticks and `<code>` tags before it matched. That matters only for a phrase of more than one word with an inner code span. For each such phrase, add a row for each spelling that the old normalization matched at the ticket base.

Remove the entries for removed rules from the handoff-fixture list test and the cadence mutation test. This is the reviewer-decided retirement of those rules. After this ticket, `AGENTS.md` keeps one Require row. So rework the anchors help test and its fixture repository to prove an absent needle with that row.

The new case-folded rows come before the debug-loop rows in registry order. So make the anchors help test select its case-folded emphasis row by its needle, not by position. Re-cut each COMMAND needle in this file that is over the bound under the `internal/prose` word count. No touched file may grow.

## Acceptance

- [ ] Deletion of each of the three invariant clauses raises its INVARIANT diagnostic.
- [ ] The core tuple test and the shared-rule check still prove the 8 core rows and the Roles marker.
- [ ] Each fixture of a removed row is gone, and each kept fixture still bites.
- [ ] The ticket record carries a removal census for the removed rows, and for each moved integration-source phrase with the kind of its new row.
- [ ] `bench test --package ./internal/anchors`, `bench test --package ./internal/conformance`, `bench test --check docs-currency-workflow`, `bench test --check conformance-canary-families`, and `bench test --package ./cmd/bench` pass.
