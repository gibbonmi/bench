# State the finding confidence in review

Blocked by: 1-state-claim-schema-on-delegate-return.md
Writes: .agents/skills/bench-craft-review/references/finding-discipline.md, .agents/skills/bench-craft-review/SKILL.md, .agents/commands/bench-review-implementation.md, internal/anchors/registry_calibration.go (new), internal/anchors/registry_calibration_test.go (new), reviews/calibrated-decisions.md (new), tests/canary/workflow-guidance-anchors/calibration-finding-confidence (new), tests/canary/workflow-guidance-anchors/calibration-finding-blocking (new), tests/canary/workflow-guidance-anchors/calibration-disposition-label (new), tests/canary/workflow-guidance-anchors/calibration-advice-no-confidence (new), tests/canary/workflow-guidance-anchors/calibration-review-pointer (new), tests/canary/workflow-guidance-anchors/calibration-pickup-confidence (new), tests/canary/workflow-guidance-anchors/craft-review-coverage-row-projection, tests/canary/workflow-guidance-anchors/dg-33, tests/canary/workflow-guidance-anchors/dg-33-reference, tests/canary/workflow-guidance-anchors/dg-37, tests/canary/workflow-guidance-anchors/dg-38, tests/canary/workflow-guidance-anchors/dg-43, tests/canary/workflow-guidance-anchors/dg-43-replay, tests/canary/workflow-guidance-anchors/review-finding-discipline-pointer, tests/canary/workflow-guidance-anchors/dg-35, tests/canary/workflow-guidance-anchors/dg-35-unavailable, tests/canary/workflow-guidance-anchors/dg-36, tests/canary/workflow-guidance-anchors/review-citation-location, tests/canary/workflow-guidance-anchors/review-env-var-producer, tests/canary/workflow-guidance-anchors/review-seam-amendment, tests/canary/workflow-guidance-anchors/review-string-expectation-catch, tests/canary/workflow-guidance-anchors/review-strong-finding-run, tests/canary/workflow-guidance-anchors/review-test-deletion-coverage, tests/canary/workflow-guidance-anchors/delegated-axis-exclusions, tests/canary/workflow-guidance-anchors/delegated-chunk-tip-review, tests/canary/workflow-guidance-anchors/prepared-review-axis-returns, tests/canary/workflow-guidance-anchors/prepared-review-blast-evidence, tests/canary/workflow-guidance-anchors/prepared-review-capable-handoff, tests/canary/workflow-guidance-anchors/prepared-review-inline-axis-route, tests/canary/workflow-guidance-anchors/prepared-review-legacy-entry-points, tests/canary/workflow-guidance-anchors/prepared-review-native-dispatch, tests/canary/workflow-guidance-anchors/prepared-review-runtime-capability, tests/canary/workflow-guidance-anchors/prepared-review-shared-evidence, tests/canary/workflow-guidance-anchors/review-base-merged-main-tip, tests/canary/workflow-guidance-anchors/review-clean-terminal-result, tests/canary/workflow-guidance-anchors/review-cross-harness-opt-in, tests/canary/workflow-guidance-anchors/review-falsification-accept-routing, tests/canary/workflow-guidance-anchors/review-falsification-dispositions, tests/canary/workflow-guidance-anchors/review-repair-ticket-covers, tests/canary/workflow-guidance-anchors/review-repair-ticket-owner, tests/canary/workflow-guidance-anchors/review-standing-falsification, tests/canary/workflow-guidance-anchors/coverage-axis-anchor, tests/canary/workflow-guidance-anchors/review-persistence-anchor, tests/canary/workflow-guidance-anchors/review-preflight-explicit-base, tests/canary/workflow-guidance-anchors/review-universal-claim-bar, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: CR4, CR5, CR6, CR15, CR16, CR22, CR23, CR30

Chunk: CD2.

## What to build

A review axis writes each finding with a stated confidence, and the reviewer's disposition labels it.
Add the finding confidence rule to the finding discipline reference.
The rule states that a finding carries a stated confidence as an integer from 0 to 10.
It states that the confidence never changes whether a finding blocks.
It states the fixed mapping: `auto-fix` and `ask-user` label a finding `held`, and `no-op` labels it `refuted`.
It states that optional advice carries no confidence.

Keep the review skill's existing pointer sentence byte for byte, period included. Three pins hold it: the registry row, its harness test, and the pointer fixture.
Add one second sentence on the same line: "A finding also states its confidence as an integer from 0 to 10."
Anchor the new sentence on its own words, and keep the file at zero net growth.
In the review-implementation command, the pickup step states that each actionable finding line carries its stated confidence.
The review record JSON types in `internal/reviewrecord` do not change; the confidence rides on the pickup line.

Scenario: an axis writes a behavior finding at confidence `3` with a citation.
The finding blocks the round, because its kind and citation decide.
The reviewer disposes it `no-op`, and the pickup labels it `refuted`.

Append six anchors to the calibration registry file and its harness test.
Each rule sentence takes one omission canary fixture in the `BASE`, `EXPECT`, `MUTATE.json` shape.
The review-pointer fixture removes only the new sentence and leaves the existing needle and its period intact.

Preserve every co-named fixture pin without changing its planted diagnostic.

## Acceptance

- [ ] The finding discipline reference carries the confidence, blocking, disposition-mapping, and advice sentences.
- [ ] The review skill's pointer line keeps the existing needle's bytes, carries the new confidence sentence, and the file holds at most 122 lines.
- [ ] The review-implementation pickup step states that each actionable finding line carries its stated confidence.
- [ ] The exported types in `internal/reviewrecord` are byte-identical to the base commit.
- [ ] Each of the six new sentences is a registry anchor, and its omission canary fails the fixture-bite test when the sentence is absent.
- [ ] Every existing fixture that pins the three edited files still plants its diagnostic.

Run `bench test --check docs-currency-workflow`, `bench test --check guidance-prose-budgets`, and `bench test --package ./internal/anchors/... --run 'TestCalibration|TestEveryRetainedFixtureBites'`.
Run the prose lane on the three edited Markdown files.
Record the scenario finding with its confidence, disposition, and label in the review pickup.
