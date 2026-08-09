# Trace bootstrap authority before execution

Blocked by: none
Ownership fence: `.agents/skills/bench-craft-spec/SKILL.md`, `.agents/commands/bench-write-spec.md`, `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/fixture_bite_test.go`, `tests/canary/workflow-guidance-anchors/bootstrap-authority-rule-deletion/BASE`, `tests/canary/workflow-guidance-anchors/bootstrap-authority-rule-deletion/MUTATE.json`, `tests/canary/workflow-guidance-anchors/bootstrap-authority-rule-deletion/EXPECT`, `tests/canary/workflow-guidance-anchors/bootstrap-authority-before-after-softening/BASE`, `tests/canary/workflow-guidance-anchors/bootstrap-authority-before-after-softening/MUTATE.json`, `tests/canary/workflow-guidance-anchors/bootstrap-authority-before-after-softening/EXPECT`, `CHANGELOG.md`, `specs/light-path-bootstrap-authority-spec-prose/tickets/trace-bootstrap-authority.md`
Integration surfaces: `craft-spec` rule owner→`.agents/skills/bench-craft-spec/SKILL.md`; `/bench-write-spec` edge-walk and falsification consumer→`.agents/commands/bench-write-spec.md`; workflow-anchor checker→`internal/conformance/docs_workflow_helpers_test.go`; direct mutation owner→`internal/conformance/fixture_bite_test.go`; deletion canary base→`tests/canary/workflow-guidance-anchors/bootstrap-authority-rule-deletion/BASE`; deletion canary mutation→`tests/canary/workflow-guidance-anchors/bootstrap-authority-rule-deletion/MUTATE.json`; deletion canary expectation→`tests/canary/workflow-guidance-anchors/bootstrap-authority-rule-deletion/EXPECT`; softening canary base→`tests/canary/workflow-guidance-anchors/bootstrap-authority-before-after-softening/BASE`; softening canary mutation→`tests/canary/workflow-guidance-anchors/bootstrap-authority-before-after-softening/MUTATE.json`; softening canary expectation→`tests/canary/workflow-guidance-anchors/bootstrap-authority-before-after-softening/EXPECT`; user-visible behavior→`CHANGELOG.md`
Contracts: none crosses
Closure: BA1/rule-anchor, BA1/before-launch-validator, BA2/edge-walk-rule-reference, BA2/falsification-rule-reference, BA3/canary-rule-deletion, BA3/canary-before-after-softening, BA3/additive-after-launch-refusal

## What to build

Specs that claim trusted or authenticated execution, or refusal before execution, expose the bootstrap authority path before implementation. The one tracer keeps the rule in `craft-spec`, makes `/bench-write-spec` name it at the two consuming steps, and gives the workflow-anchor conformance owner direct and canary deletion/order bites. The trace's conditional trigger, self-authentication prohibition, absence semantics, raw-entrypoint marker probe, and slicing ownership are semantic reviewer-owned claims; conformance protects the rule anchor and rejects missing, softened, or additive after-launch validation without claiming to parse the other meanings.

## Acceptance

- [ ] [BA1] `craft-spec` publishes the named bootstrap-authority rule with a before-launch validator; review owns its conditional trace, self-authentication, absence, marker-probe, and slicing semantics.
- [ ] [BA2] `/bench-write-spec` names the `craft-spec` bootstrap-authority rule while walking edges and charging falsification, without restating the rule.
- [ ] [BA3] The real workflow-anchor checker rejects deletion, before-to-after softening, and an additive after-launch instruction with attributable diagnostics; real canary fixtures cover deletion and softening.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BA1/rule-anchor | delete the bootstrap-authority rule heading | the workflow-anchor fixture bite | remove the heading from `craft-spec`, run `go test -count=1 ./internal/conformance -run '^TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting$'`, expect the bootstrap-authority deletion diagnostic |
| BA1/before-launch-validator | replace the before-launch validator requirement with an after-launch check | the workflow-anchor fixture bite | soften the bootstrap-authority order in `craft-spec`, run the focused conformance test, expect the after-launch diagnostic |
| BA2/edge-walk-rule-reference | delete the named bootstrap-authority reference from the edge-walk step | the workflow-anchor fixture bite | remove the edge-walk reference from `/bench-write-spec`, run the focused conformance test, expect the command-reference diagnostic |
| BA2/falsification-rule-reference | delete the named bootstrap-authority reference from the falsification step | the workflow-anchor fixture bite | remove the falsification reference from `/bench-write-spec`, run the focused conformance test, expect the command-reference diagnostic |
| BA3/canary-rule-deletion | delete the bootstrap-authority rule heading | the workflow-guidance-anchors canary | run the deletion canary through `TestDocsCurrencyTokenDietAndWorkflowFixturesBite`, expect its distinct deletion diagnostic |
| BA3/canary-before-after-softening | replace the before-launch validator requirement with an after-launch check | the workflow-guidance-anchors canary | run the softening canary through `TestDocsCurrencyTokenDietAndWorkflowFixturesBite`, expect its distinct after-launch diagnostic |
| BA3/additive-after-launch-refusal | append an instruction permitting validation after launch while retaining the before-launch requirement | the workflow-anchor unit mutation | append the after-launch instruction in `craft-spec`, run the focused conformance test, expect the after-launch diagnostic |
