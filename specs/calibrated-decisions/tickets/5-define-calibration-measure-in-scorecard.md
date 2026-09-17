# Define the calibration measure in the scorecard

Blocked by: 4-render-calibration-table-in-retro-scaffold.md
Writes: capture/agent-performance/README.md, capture/agent-performance/claude-models.md, capture/agent-performance/open-ai-models.md, internal/anchors/registry_calibration.go (new), internal/anchors/registry_calibration_test.go (new), reviews/calibrated-decisions.md (new), tests/canary/workflow-guidance-anchors/calibration-scorecard-measure (new), tests/canary/workflow-guidance-anchors/calibration-routing-input (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: CR24, CR25, CR26

Chunk: CD4.

## What to build

A retro author reads the scorecard and finds one calibration measure with one definition.
Add one `calibration` row to the README's Measures table.
The row defines the Brier mean, the pair count, and the abstention count.
In the update contract, state that a provider with no labeled pair shows `unknown` in the cell.
State that the calibration measure is one input to the two-run routing rule, obeys the ten-assignment cap, and never moves a tier on its own.

Add a `calibration` column to each provider routing table.
Write `unknown` in every cell, because no labeled pair exists before this build's retro.
Each provider file stays at or below 120 lines under the README's own contract.

Scenario: the Claude scorecard's six routing rows each show `unknown` in the new column.
A later retro replaces one cell with a Brier mean, a pair count, and an abstention count.

Append two anchors to the calibration registry file and its harness test, each with one omission canary fixture.
The README has no existing anchor, so these are its first two.
`docs-currency-workflow` is a Dev-tier check, so a repo-only anchor path under `capture/` follows the `projects/benchkit.md` precedent.

## Acceptance

- [ ] The README's Measures table carries a `calibration` row that defines the Brier mean, the pair count, and the abstention count.
- [ ] The README's update contract carries the two-run, ten-assignment-cap, no-tier-move sentence, and its omission canary bites.
- [ ] Each provider routing table carries a `calibration` column, and every cell reads `unknown`.
- [ ] Both provider files hold at most 120 lines.

Run `bench test --check docs-currency-workflow` and `bench test --package ./internal/anchors/... --run 'TestCalibration|TestEveryRetainedFixtureBites'`.
Run the prose lane on the three edited Markdown files.
Record the column addition and the `unknown` cell rule in the review pickup.
