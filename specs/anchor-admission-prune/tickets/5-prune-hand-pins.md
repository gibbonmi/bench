# Prune the hand-coded prose pins

Blocked by: 4a-resolve-earlier-deletions.md
Writes: internal/conformance/docs_workflow_helpers_test.go, internal/conformance/docs_workflow_checks_test.go, internal/conformance/validity_checks_test.go, internal/conformance/skills_index_checks_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/recurrence_maintenance_contract_test.go, internal/conformance/fixture_bite_test.go, internal/conformance/checks_test.go, internal/conformance/commitment_guidance_test.go, internal/conformance/tier_live_tree_test.go, tests/canary/workflow-guidance-anchors, tests/canary/skills-index-command-adapters
Covers: AP28, AP29, AP30, AP32

## What to build

The hand-coded prose pins in the conformance tests keep only admitted facts. A pin stays when its text is a Forbid form, a parser or template label, a routing reference, or a `.bench/BENCH.md` invariant. A pin also stays when a program reads its text. Every other pin is a rule-sentence pin, so it leaves with its test cases and its fixtures.

The structured-phase clause check, the shared-rule check, the skills-index checks, and the `BENCH_MODEL` pin stay. The pin on "implementation-line recommendation contract" leaves. Remove a helper that loses its last caller, but keep the anchor-bite helpers that other tests share.

The ticket record carries a pin census. It lists each pin site with its class and its disposition, and each fixture that leaves. `assets/fixture-census.md` lists the expected fixtures for this ticket, and the fixture-bite test decides each borderline case. No touched file may grow.

## Acceptance

- [ ] No hand-coded conformance pin requires rule wording outside the admitted classes.
- [ ] Each kept pin still turns red when its text goes.
- [ ] Each `.bench/BENCH.md` marker has exactly one enforcing check.
- [ ] `bench test --package ./internal/conformance` and `bench test --check docs-currency-workflow` pass.
