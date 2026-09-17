# Declare expected repair rounds and own the score rule

Blocked by: 2-state-finding-confidence-in-review.md
Writes: .agents/skills/bench-craft-line/references/calibration-score.md (new), .agents/skills/bench-craft-line/SKILL.md, internal/anchors/registry_calibration.go (new), internal/anchors/registry_calibration_test.go (new), reviews/calibrated-decisions.md (new), tests/canary/workflow-guidance-anchors/calibration-line-expected-rounds (new), tests/canary/workflow-guidance-anchors/calibration-brier-rule (new), tests/canary/workflow-guidance-anchors/calibration-abstention-score (new), tests/canary/workflow-guidance-anchors/calibration-label-sources (new), tests/canary/workflow-guidance-anchors/calibration-expectation-label (new), tests/canary/workflow-guidance-anchors/calibration-declaration-line (new), tests/canary/workflow-guidance-anchors/delegated-author-limit, tests/canary/workflow-guidance-anchors/delegated-mid-review-route, tests/canary/workflow-guidance-anchors/delegated-tier-authorization, tests/canary/workflow-guidance-anchors/delegated-unbound-model-stop, tests/canary/workflow-guidance-anchors/ticket-stage-routing-anchor, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: CR7, CR8, CR9, CR10, CR11, CR12, CR31

Chunk: CD3.

## What to build

A retained author declares the line with an expected repair-round count and a confidence, and one reference owns the score.
Create the calibration score reference under the line skill.
It states that the line declaration states the expected repair-round count and a stated confidence as an integer from 0 to 10.
It states that the calibration score of one claim is `(p - label)^2`, where `p = n / 10` and the label is 1 for `held` and 0 for `refuted`.
It states that an abstention scores 0, stays out of the Brier mean, and is counted apart.

It states that a label source is the gate, the coordinator's probe of the exact tree, or the reviewer's disposition.
It states that a model judgment is never a label source.
It states that the repair-attribution table's actual round count labels the expectation `held` when it equals the expected count and `refuted` otherwise.

Add one `Expected repair rounds:` line to the declaration block in the line skill.
The skill sits at its 130-line budget, so reclaim one line elsewhere in the same file.
Name the reclaimed line in the review pickup, and keep every anchored sentence's bytes.
The declaration section also points at the calibration score reference.

Scenario: an author declares `Expected repair rounds: 1 / confidence 8`.
The build's repair-attribution table records two rounds.
The retro labels the expectation `refuted` and scores it `(0.8 - 0)^2 = 0.64`.

Append six anchors to the calibration registry file and its harness test, each with one omission canary fixture.
The declaration-line anchor uses `require-in-section` on the `The declaration` section.

Preserve every co-named fixture pin without changing its planted diagnostic.

## Acceptance

- [ ] The calibration score reference carries the expected-rounds, Brier-rule, abstention, label-source, and expectation-label sentences.
- [ ] The line skill's declaration block carries the `Expected repair rounds:` line, and the file holds at most 130 lines.
- [ ] The review pickup names the line the skill reclaimed, and no existing anchor needle in the skill changed.
- [ ] Each of the six new sentences is a registry anchor, and its omission canary fails the fixture-bite test when the sentence is absent.
- [ ] Every existing fixture that pins the line skill still plants its diagnostic.

Run `bench test --check docs-currency-workflow`, `bench test --check guidance-prose-budgets`, and `bench test --package ./internal/anchors/... --run 'TestCalibration|TestEveryRetainedFixtureBites'`.
Run the prose lane on both edited Markdown files.
Record this build's own line declaration, with its expected rounds and confidence, in the review pickup.
