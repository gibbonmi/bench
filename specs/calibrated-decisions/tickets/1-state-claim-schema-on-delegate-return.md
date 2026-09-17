# State the claim schema on the delegate return

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-delegate/SKILL.md, internal/anchors/registry_calibration.go (new), internal/anchors/registry_calibration_test.go (new), internal/anchors/registry_data.go, reviews/calibrated-decisions.md (new), tests/canary/workflow-guidance-anchors/calibration-claim-schema (new), tests/canary/workflow-guidance-anchors/calibration-abstained-row (new), tests/canary/workflow-guidance-anchors/calibration-coordinator-label (new), tests/canary/workflow-guidance-anchors/calibration-malformed-row (new), tests/canary/workflow-guidance-anchors/calibration-no-free-text (new), tests/canary/workflow-guidance-anchors/calibration-delegate-pointer (new), tests/canary/claude-agent-definitions/agent-unnamed-in-skill, tests/canary/claude-agent-definitions/skill-names-missing-agent, tests/canary/workflow-guidance-anchors/delegate-cap-change-pinning-package, tests/canary/workflow-guidance-anchors/delegate-charge-effort-cap, tests/canary/workflow-guidance-anchors/delegate-coverage-row-charge, tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green, tests/canary/workflow-guidance-anchors/delegate-exec-only-every-caller, tests/canary/workflow-guidance-anchors/delegate-model-id-escalation, tests/canary/workflow-guidance-anchors/delegate-own-family-native-surface, tests/canary/workflow-guidance-anchors/delegate-release-at-acceptance, tests/canary/workflow-guidance-anchors/delegate-resume-handoff-contents, tests/canary/workflow-guidance-anchors/delegate-self-probe-missing-row, tests/canary/workflow-guidance-anchors/delegated-per-ticket-author, tests/canary/workflow-guidance-anchors/dg-6-owner-reference, tests/canary/workflow-guidance-anchors/prepared-review-delegate-handoff-route, tests/canary/workflow-guidance-anchors/prepared-triage-bounds, tests/canary/workflow-guidance-anchors/delegate-anchor-probe-owning-check, tests/canary/workflow-guidance-anchors/delegate-grammar-fence-inventory, tests/canary/workflow-guidance-anchors/delegate-live-tree-inventory-fence, tests/canary/workflow-guidance-anchors/delegate-out-of-fence-write, tests/canary/workflow-guidance-anchors/delegate-probe-mutated-bytes, tests/canary/workflow-guidance-anchors/delegate-root-conformance-pass, tests/canary/workflow-guidance-anchors/delegate-serial-ceiling-fence, tests/canary/workflow-guidance-anchors/delegate-skip-ownership-check, tests/canary/workflow-guidance-anchors/delegated-cap-exhaustion-trigger, tests/canary/workflow-guidance-anchors/delegated-lost-session-trigger, tests/canary/workflow-guidance-anchors/delegated-no-progress-trigger, tests/canary/workflow-guidance-anchors/delegated-terminal-failure-trigger, tests/canary/workflow-guidance-anchors/delegated-writer-termination, tests/canary/workflow-guidance-anchors/prepared-triage-authority, tests/canary/workflow-guidance-anchors/prepared-triage-cli-boundary, tests/canary/workflow-guidance-anchors/prepared-triage-input-contract, tests/canary/claude-agent-definitions/model-declared, tests/canary/claude-agent-definitions/name-mismatch, tests/canary/claude-agent-definitions/shell-tool-absent, tests/canary/claude-agent-definitions/spawning-tool, tests/canary/claude-agent-definitions/tools-absent, tests/canary/workflow-guidance-anchors/delegate-cross-harness-reviewer-pointer, tests/canary/workflow-guidance-anchors/delegate-parallel-route-anchor, tests/canary/workflow-guidance-anchors/delegate-stash-refusal-anchor, tests/canary/workflow-guidance-anchors/fix-pass-sentinel-anchor, tests/canary/workflow-guidance-anchors/shared-worktree-path-pin, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: CR1, CR2, CR3, CR13, CR14, CR29, CR32

Chunk: CD1.

## What to build

A coordinator reads a delegate return and finds each done-claim row in the claim schema.
Add one `Claim schema` section to the delegation discipline reference.
The section states the three status values, the integer confidence range, and the no-free-text rule.
It states that an `abstained` row carries no confidence and that no later probe turns an abstention into a refuted claim.
It states that the coordinator's probe of the exact tree labels each row `held` or `refuted`, whatever its status.
It states the three malformed shapes and that a malformed row returns to the author before any label.

Add one pointer sentence to the delegate skill inside its two lines of headroom.
The sentence points the charge at the `Claim schema` section, so a delegate returns the exact shape.

Scenario: a delegate returns a row with status `claimed` and confidence `7`.
The coordinator probes the row's named check on the exact tree and labels the row `held`.
A second row arrives with status `abstained` and confidence `9`.
The coordinator returns that row to the author as malformed and labels nothing.

Create the calibration anchors registry file beside the ticket-passes precedent, with its own harness test.
Append its group to the registry composition line.
Each new rule sentence becomes one `require-in-section` anchor with one omission canary fixture.
The six fixtures follow the `BASE`, `EXPECT`, `MUTATE.json` shape of `delegate-self-probe-missing-row`.
Successor tickets append their anchors to this same registry file and its test.

Preserve every co-named fixture pin without changing its planted diagnostic.
Reclaim no line: the delegate skill has two lines of headroom, and this ticket uses at most one.

## Acceptance

- [ ] The delegation discipline reference carries a `Claim schema` section with the status, confidence, and no-free-text sentences.
- [ ] The section carries the abstained-row sentence, the coordinator-label sentence, and the malformed-row sentence.
- [ ] The delegate skill carries one sentence that points the charge at the `Claim schema` section, and the file holds at most 126 lines.
- [ ] Each of the six new sentences is a registry anchor, and its omission canary fails the fixture-bite test when the sentence is absent.
- [ ] Every existing fixture that pins the two edited files still plants its diagnostic.

Run `bench test --check docs-currency-workflow`, `bench test --check guidance-prose-budgets`, and `bench test --package ./internal/anchors/... --run 'TestCalibration|TestEveryRetainedFixtureBites'`.
Run the prose lane on both edited Markdown files.
Record the two scenario rows with their confidence and label in the review pickup.
