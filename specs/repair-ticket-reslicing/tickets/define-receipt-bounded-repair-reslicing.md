# Define receipt-bounded repair reslicing

Blocked by: none
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `internal/anchors/registry_data.go`, `internal/conformance/fixture_bite_test.go`
Integration surfaces: canonical repair-reslicing owner→orchestrate-receipt-derived-repair-chains.md; workflow-anchor registry→`internal/anchors/registry_data.go`; unchanged registry consumer→`internal/anchors/registry.go` + RS1; literal-mutation harness→`internal/conformance/fixture_bite_test.go`; unchanged anchor evaluator→`internal/conformance/docs_workflow_helpers_test.go` + RS1; unchanged workflow-check consumer→`internal/conformance/docs_workflow_checks_test.go` + RS1; unchanged graded root→`internal/conformance/gate_entry_test.go` + RS1; unchanged section-anchor kinds→`internal/anchors/match.go` + RS1
Contracts: canonical repair-reslicing owner reference (skill-path type, exact `.agents/skills/bench-craft-tickets/SKILL.md` domain, producer ticket before consumer ticket ordering, missing pointer invalid) crossing `.agents/skills/bench-craft-tickets/SKILL.md`→orchestrate-receipt-derived-repair-chains.md, asserted by dependent row OC1 against the real skill producer
Closure: RE1/maximum-envelope, RR1/one-ticket-result, RR1/reciprocal-chain-result, RU1/union-contained, RU1/escape-forbidden

## What to build

Make `craft-tickets` the single owner of repair reslicing. A validated debug receipt's required fence is the maximum repair envelope, not one indivisible ticket fence. Apply the ordinary independently-green split rule inside it: the result may be the common one-ticket repair or a reciprocal ordered producer-to-consumer chain, and the union of every resulting ticket fence stays inside the receipt's required fence. Extend the existing section-scoped workflow anchors and real-source literal-mutation harness so omission and additive contradiction both fail through the graded conformance root.

## Acceptance

- [ ] [RE1] (covers RS1) The receipt's required fence is the maximum envelope in which ordinary independently-green repair reslicing runs.
- [ ] [RR1] (covers RS2) Repair reslicing permits both one repair ticket and a reciprocal ordered repair chain.
- [ ] [RU1] (covers RS3) Every repair-ticket fence is inside the receipt's required fence, so their union cannot escape it, and contradictory escape permission is rejected while the positive rule remains.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RE1/maximum-envelope | remove the maximum-envelope clause while preserving the one-or-chain and union clauses | the workflow-guidance literal-mutation harness | mutate the real skill subject, run `go test -count=1 ./internal/conformance -run '^TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting$'`, require the missing-envelope diagnostic, restore the subject, and rerun green |
| RR1/one-ticket-result | replace the one-or-chain result with chain-only prose while preserving the envelope and union clauses | the workflow-guidance literal-mutation harness | mutate the real skill subject, run `go test -count=1 ./internal/conformance -run '^TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting$'`, require the missing-result diagnostic, restore the subject, and rerun green |
| RR1/reciprocal-chain-result | replace the one-or-chain result with one-ticket-only prose while preserving the envelope and union clauses | the workflow-guidance literal-mutation harness | mutate the real skill subject, run `go test -count=1 ./internal/conformance -run '^TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting$'`, require the missing-result diagnostic, restore the subject, and rerun green |
| RU1/union-contained | remove the positive union-containment clause while preserving the envelope and result clauses | the workflow-guidance literal-mutation harness | mutate the real skill subject, run `go test -count=1 ./internal/conformance -run '^TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting$'`, require the missing-union diagnostic, restore the subject, and rerun green |
| RU1/escape-forbidden | add permission for one chain ticket to escape the receipt's required fence while preserving the positive union clause | the workflow-guidance literal-mutation harness | mutate the real skill subject, run `go test -count=1 ./internal/conformance -run '^TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting$'`, require the escaped-union diagnostic, restore the subject, and rerun green |
