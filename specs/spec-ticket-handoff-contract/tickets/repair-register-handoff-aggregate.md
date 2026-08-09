# Register the handoff aggregate as a live-tree assertion

Blocked by: none
Ownership fence: `internal/conformance/tier_test.go`
Integration surfaces: `TestSpecTicketHandoffWorkflowFixturesAreComplete`→`classifiedLiveTreeTests`
Contracts: the live-tree aggregate crosses `internal/conformance/fixture_bite_test.go`→`internal/conformance/tier_test.go`; type is an exact Go test name, domain is classified live-tree assertions, order is lexical with neighboring `TestSpec...` members, absence makes `TestConformanceMetaBites` red on the prospective implemented tree
Closure: LR1/live-tree-classification

## What to build

Close the candidate-attributed promotion red by registering `TestSpecTicketHandoffWorkflowFixturesAreComplete` in the existing classified live-tree test inventory. Add no new runner, parser, or authority.

## Acceptance

- [ ] [LR1] (covers SH11) the exact prospective `Status: implemented` subject passes `TestConformanceMetaBites` because the new handoff aggregate is explicitly classified.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LR1/live-tree-classification | omit the aggregate's exact test name from `classifiedLiveTreeTests` | `TestConformanceMetaBites` | materialize the prospective implemented subject, run the named meta test, require `conformance meta unregistered live-tree assertion TestSpecTicketHandoffWorkflowFixturesAreComplete`, restore the member |
