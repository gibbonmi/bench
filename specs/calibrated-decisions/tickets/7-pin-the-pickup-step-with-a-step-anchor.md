# Pin the pickup step with a step-scoped anchor kind

Blocked by: 2-state-finding-confidence-in-review.md
Writes: internal/anchors/registry.go, internal/anchors/match.go, internal/anchors/locate.go, internal/anchors/anchor_harness_test.go, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_data_test.go, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, cmd/bench/anchors_command.go, cmd/bench/anchor_help_test.go, reviews/calibrated-decisions.md, tests/canary/workflow-guidance-anchors/calibration-pickup-step-move (new), tests/canary/workflow-guidance-anchors/calibration-pickup-confidence, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: CR22, CR38

Chunk: CD2b.

## What to build

A maintainer registers an anchor that requires a sentence inside one numbered step of an H2 section.
A move of that sentence to another step reds `docs-currency-workflow`.
The reviewer decided this seam during the CD2 review, because `require-in-section` cannot pin step 6 of the review-implementation command.

Add one `RequireInStep` kind to the anchors package and one `Step int` field beside `Section` on `Anchor`.
A zero `Step` means the anchor is not step-scoped.
The kind implies section scope: the locator resolves the H2 section first, then narrows the body to the step.

A step starts at an unindented line outside a fenced block. That line opens with a decimal number, a period, and a space or tab at column zero.
The step body runs to the next such line or to the section end.
An indented continuation line stays inside its step.
The locator matches the literal step digits the reader sees, not the ordinal position.

A section with no such step raises its own diagnostic. A section with two lines that open the same step raises its own diagnostic. Model both on the missing-section and duplicate-section diagnostics.

Model the step helper on `sectionRunesMapped` in `locate.go`, including its fence toggle and its origin mapping, so the reported line stays exact.
Add the kind to the evaluator's require arm and to the scoped-kind predicate in `match.go`.
Cache the step result beside the section cache in `registry.go`, keyed by file, section, and step.

Extend the `bench anchors` projection with a fifth `step` column after `section`, and add `require-in-step` to the kind names.
A non-step anchor prints `0` in that column.
Update every hard-coded table expectation in `cmd/bench/anchor_help_test.go`, including the empty-table literal.

Extend the anchor harness so a rule can carry a step. For a step-scoped rule, the harness writes numbered lines inside the section body.
Add `TestAnchorHarnessStepRules` beside the existing harness tests: a positive rule, a move-to-another-step red, a missing-step diagnostic, and a duplicate-step diagnostic.
Add `RequireInStep` to the kind table in `TestEvaluatePathAnchorKinds`.

Re-pin CR22: change the `calibration-pickup-confidence` anchor in `registry_calibration.go` to `RequireInStep` with `Section: "Process"` and `Step: 6`, and keep its needle and diagnostic bytes.
Update its harness test row for the new kind and step.
Keep the existing `calibration-pickup-confidence` fixture as the omission canary.
Add one `calibration-pickup-step-move` fixture. Its `MUTATE.json` deletes the sentence from step 6 and inserts the identical bytes into step 5. Its `EXPECT` is the diagnostic of the anchor.

Scenario: the pickup sentence sits in step 6 of the `Process` section, and `docs-currency-workflow` passes.
An author moves the sentence into step 5 and keeps every other byte.
The check fails with the `calibration-pickup-confidence` diagnostic, and `bench anchors .agents/commands/bench-review-implementation.md` prints `require-in-step,Process,6` for the row.

Preserve every co-named fixture pin without changing its planted diagnostic.
Keep `internal/anchors/registry_data.go` at 490 lines; this ticket adds no anchor group.

## Acceptance

- [ ] A harness rule with a step whose needle sits in that step raises no diagnostic.
- [ ] The same needle under another step raises the anchor's diagnostic.
- [ ] A section with no such step, and a section with two lines that open the same step, each raise their own diagnostic.
- [ ] `bench anchors <path>` prints a `step` column, `0` for a section anchor and the step number for a step anchor, and names the kind `require-in-step`.
- [ ] The `calibration-pickup-confidence` anchor is step-scoped to step 6 of `Process`, and the `calibration-pickup-step-move` canary bites through `docs-currency-workflow`.
- [ ] Every existing fixture that pins the review-implementation command still plants its diagnostic.

Run `bench test --package ./internal/anchors/...`, `bench test --package ./cmd/bench/... --run 'TestAnchors'`, `bench test --check docs-currency-workflow`, and `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
Run `go vet ./...` and the prose lane on the pickup.
Record the step-move probe and the done-claim rows with their confidence in the review pickup.
