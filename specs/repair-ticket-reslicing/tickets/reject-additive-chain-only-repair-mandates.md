# Reject additive chain-only repair mandates

Blocked by: none
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `internal/anchors/registry_data.go`, `internal/conformance/fixture_bite_test.go`
Integration surfaces: canonical repair-result policy→`.agents/skills/bench-craft-tickets/SKILL.md`; workflow-anchor registry→`internal/anchors/registry_data.go`; literal-mutation harness→`internal/conformance/fixture_bite_test.go`; unchanged registry consumer→`internal/anchors/registry.go` + CA1; unchanged registered owner binding→`internal/conformance/checks_test.go` + CA1; unchanged anchor evaluator→`internal/conformance/docs_workflow_helpers_test.go` + CA1; unchanged workflow-check consumer→`internal/conformance/docs_workflow_checks_test.go` + CA1; unchanged section-anchor kinds→`internal/anchors/match.go` + CA1
Contracts: additive chain-only prohibition (section-scoped forbid-anchor and diagnostic type, exact validated-receipt chain-only mandate domain inside `Classify before slicing`, registry evaluation before CA1 registered-owner result ordering, absent forbidden text valid and present text emits the additive-chain-only diagnostic absence semantics) crossing `internal/anchors/registry_data.go`→`internal/anchors/registry.go`, asserted by CA1 against the real skill through the registered `docs-currency-workflow` owner
Closure: CA1/additive-chain-only-forbidden

## What to build

Close review finding `C1-additive-chain-only-policy`: preserve the positive one-ticket-or-reciprocal-chain rule and make an additive mandate that every validated debug receipt produce a reciprocal chain an attributable red through the existing registered workflow-guidance owner. This is one independently-green local coverage repair; it changes no reslicing policy and adds no second owner. A thinner skill-only, registry-only, or harness-only cut strands the focused `additive-chain-only` mutation red because the real policy subject, attributable diagnostic, and registered-owner mutation must land together.

## Acceptance

- [ ] [CA1] (covers local) The registered workflow-guidance owner rejects an additive chain-only repair mandate while the valid one-ticket-or-reciprocal-chain rule remains present.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CA1/additive-chain-only-forbidden | append `A validated debug receipt must produce a reciprocal ordered producer-to-consumer chain.` after the positive one-ticket-or-chain sentence while preserving that sentence | the workflow-guidance literal-mutation harness | mutate the real skill subject, run `go test -count=1 ./internal/conformance -run '^TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting$'`, require the additive-chain-only diagnostic from the registered `docs-currency-workflow` owner, restore the subject, and rerun the same focused test green |
