# Render the calibration table in the retro scaffold

Blocked by: 3-declare-expected-repair-rounds-with-score.md
Writes: internal/retros/retros.go, internal/roadmap/retro_scaffold.go, internal/roadmap/retro_scaffold_test.go, .agents/commands/bench-final-check.md, internal/anchors/registry_calibration.go (new), internal/anchors/registry_calibration_test.go (new), reviews/calibrated-decisions.md (new), tests/canary/workflow-guidance-anchors/calibration-retro-table-duty (new), tests/canary/workflow-guidance-anchors/calibration-retro-aggregate-duty (new), tests/canary/workflow-guidance-anchors/delegated-account-inventory, tests/canary/workflow-guidance-anchors/delegated-account-reconciliation, tests/canary/workflow-guidance-anchors/delegated-final-verification, tests/canary/workflow-guidance-anchors/delegated-paid-comparison, tests/canary/workflow-guidance-anchors/final-check-bare-leftover-clean-retired, tests/canary/workflow-guidance-anchors/final-check-census-read-before-land, tests/canary/workflow-guidance-anchors/final-check-landed-worktree-sweep, tests/canary/workflow-guidance-anchors/final-check-light-path-changelog-heading, tests/canary/workflow-guidance-anchors/final-check-scratch-branch-clean, tests/canary/docs-currency-token-diet/introduces-undeclared-command, tests/canary/docs-currency-token-diet/stale-command-reference, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: CR17, CR18, CR19, CR20, CR21, CR28, CR35

Chunk: CD4.

## What to build

A retro author runs `bench retro <slug> --scaffold` and finds the calibration table shape in the draft.
Export the calibration table header as one constant in the retros package, beside the two derived-section headings.
The header names six cells: surface, claim, status, confidence, label, and model, effort, and role.
The required heading list does not change.

The scaffold renders the header and one row of six `unknown` cells under the existing delegate-performance heading.
It renders them as a third case in `scaffoldSection`, rendered in the delegate-performance section that precedes the repair table.
No second spelling of the header exists in the renderer.
The scaffold renders the table in a repository with no `capture/agent-performance` directory.

In the final-check command, the retro capture section states two duties.
The retro fills the calibration table with one row per labeled claim holding its surface, claim, status, confidence, label, and model, effort, and role.
The retro states the Brier mean, the pair count, and the abstention count below the table, with `unknown` for a mean over zero pairs.

Scenario: a temporary repository holds a spec with two tickets and no scorecard directory.
The scaffold output's delegate-performance section ends with the header row and exactly one `unknown` row.
`retros.Parse` still accepts `testdata/eligible.md`.

Add `TestRetroScaffoldRendersCalibrationTable` to the scaffold test file.
It builds the repository through `newScaffoldRepo`, asserts the scorecard directory is absent, and reads the section through `sectionOf`.
This ticket is the last that touches `internal/retros` and `internal/roadmap`, so it carries the package invariant.
Append two anchors to the calibration registry file and its harness test, each with one omission canary fixture.

Preserve every co-named fixture pin without changing its planted diagnostic.

## Acceptance

- [ ] The scaffold renders the calibration header row under the delegate-performance heading in a repository with no scorecard directory.
- [ ] The scaffold renders exactly one row of six `unknown` cells beneath the header.
- [ ] `TestRetroScaffoldRendersCalibrationTable` compares the rendered header with the exported constant, and the review greps `internal/roadmap` for a second spelling.
- [ ] The header spelling lives in one exported retros constant, and the renderer has no second spelling.
- [ ] `retros.Parse` accepts the nine-heading body in `testdata/eligible.md` after the change.
- [ ] The final-check command carries the table-duty sentence and the aggregate-duty sentence, and each omission canary bites.
- [ ] Every existing fixture that pins the final-check command still plants its diagnostic.

Run `bench test --package ./internal/roadmap/... ./internal/retros/...`, `bench test --check docs-currency-workflow`, and `bench test --package ./internal/anchors/... --run 'TestCalibration|TestEveryRetainedFixtureBites'`.
Run the prose lane on the edited command file.
Record the scaffold scenario output in the review pickup.
