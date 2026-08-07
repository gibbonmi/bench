# Route all matrix captures through the recorder

Blocked by: repair-generation-recorder-and-raw-slot-bytes.md
Ownership fence: `internal/gate/check_slots_test.go`
Integration surfaces: exhaustive matrix generation capture→`internal/gate/check_slots_test.go` + OP1; accepted finding `SPEC-OP1-001`→OP1
Contracts: every real working-tree generation captured by the exhaustive changed/restored-state matrix crosses the recorder owned by `internal/gate/check_slots_test.go`, asserted by OP1 at the actual capture seam
Closure: OP1/no-raw-alias, OP1/second-real-capture-red

## What to build

Repair the integrated matrix so no raw `captureWorkingTree` alias can bypass its generation recorder. Preserve the test-only ownership fence, the literal public-document expectations, one seeded full-engine green, direct decision calls, raw slot-store comparison, and the existing representative controls. The accepted review finding is `SPEC-OP1-001`: the candidate's raw capture alias can perform a second real capture without incrementing the wrapper-owned count.

## Acceptance

- [ ] [OP1] Every real working-tree generation capture used by each exhaustive changed or restored state is observed at the actual capture seam, and adding a second real capture through the candidate's former raw-alias route makes the focused matrix fail with the generation bound before the candidate is restored byte-exact.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OP1/no-raw-alias | retain or recreate a callable raw `captureWorkingTree` alias beside the recorded capture path and use it for an additional real capture | the seam-owned matrix generation bound | against candidate `aff74e24b08b34801c70058dbcf130b72efc5b61`, add the second raw-alias capture, run the focused matrix, observe the current false green; implement the repair, repeat the same mutation, require an exact generation-bound red, then restore the candidate byte-exact |
| OP1/second-real-capture-red | invoke the matrix's real working-tree capture path twice for one changed state | the seam-owned matrix generation bound | add the second capture at the repaired seam, run the focused matrix, require the capture count to report two rather than one, then restore the test and rerun it green |
