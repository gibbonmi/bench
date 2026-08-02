# Enumerate the operation-named refusals through the real CLI

Blocked by: Exempt abandon from the recomposition refusal; Pass the found green
marker as the expected prior tip; Reach recomposition before the gates it
discards
Ownership fence: `internal/contract/runtime/runtime_spec_build_test.go` and any
new sibling file under `internal/contract/runtime`
Assumptions: this suite runs the real `bin/bench.sh` as a subprocess against a
temporary repository, so it is the only seam that observes what an operator
actually reads. The Go-side ticket owns the `mutation` token; this ticket owns
the operator-facing spelling and is where the two are reconciled. Re-derive the
suite's existing spec-build helpers from the tree at pickup rather than adding a
second fixture family.

## What to build

This is the consumer end of three contracts that no single upstream ticket can
assert alone, and it asserts them against the real producer — the shipped CLI —
rather than against a package-level double.

**The operation names, enumerated.** Six precondition-gated operations reach an
operator: `assign`, `checkpoint`, `integrate`, `review`, `promote`, and
`abandon --apply`. The internal `mutation` token and the operator-facing
spelling are two different vocabularies, and nothing upstream reconciles them —
`abandon --apply` is one operation to an operator and one `mutation` constant
plus a flag inside the process. Assert over the whole family. One shared literal
satisfies a single-operation assertion, which is why enumerating every operation
is the only version of this row that bites.

**The `review` case needs a real receipt.** `review` validates its receipt
*before* it reaches the precondition layer, so a stub receipt returns the
receipt refusal and the row passes while testing nothing. Drive it with a valid
three-axis receipt bound to the current candidate.

**Two messages, not one.** The resolver produces both a dirty-checkout refusal
and a no-working-branch refusal. Fixing one is the cheap half-fix.

**The process boundary is the point.** Unit-level green has hidden serialization
defects in this exact package before, so every row here reloads the run record
through a fresh process rather than reusing an in-process service.

## Acceptance

- [ ] CT1 — each of the six operations names itself, and not `start`, in the dirty-checkout refusal through a real CLI process.
- [ ] CT2 — each of the same six operations names itself in the no-working-branch refusal.
- [ ] CT3 — the `review` row reaches the precondition refusal rather than the receipt refusal, driven by a valid three-axis receipt bound to the current candidate.
- [ ] CT4 — `start` on a branch whose green marker is divergent refuses end to end, and `start` on a branch whose marker is a strict ancestor succeeds end to end.
- [ ] CT5 — `abandon --apply` retires a run whose branch tip advanced, end to end, and its recovery refs survive.
- [ ] CT6 — a refusal carrying a branch name with a control byte reaches the operator intact and escaped, with no raw control byte in the emitted bytes.
- [ ] CT7 — each row above observes a fresh process reloading the run record, not a reused in-process service.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CT1 | name the operation for `assign` and leave the other five borrowing `start` | `TestDirtyCheckoutRefusalNamesEachOperation` | start a run, dirty the checkout, run each of the six CLI invocations as a subprocess, assert each refusal names its own operation |
| CT2 | fix the dirty-checkout literal only | `TestNoWorkingBranchRefusalNamesEachOperation` | detach HEAD, run the same six invocations, assert each refusal names its own operation |
| CT3 | drive the `review` row with a stub receipt | `TestReviewReachesPreconditionRefusal` | write a valid three-axis receipt bound to the current candidate, dirty the checkout, run `review`, assert the precondition refusal rather than the receipt refusal |
| CT4 | pass the found marker unconditionally, without the owner's ancestor rule | `TestStartMarkerAncestryEndToEnd` | two subtests planting an ancestor marker and a divergent marker, running `start` as a subprocess, expecting success and refusal respectively |
| CT5 | exempt the whole precondition call for abandon rather than the recomposition branch | `TestAbandonRetiresRunOnMovedTipEndToEnd` | compose a run, commit on the branch, run `abandon` then `abandon --apply`, assert terminal state and surviving recovery refs |
| CT6 | interpolate the operation into the message with plain concatenation, bypassing the control escaper | `TestRefusalSurvivesControlByteBranchName` | create a branch whose name carries a control byte, trigger a refusal, assert the emitted bytes carry no raw control byte |
| CT7 | assert the rows against an in-process service reusing a loaded record | `TestSpecBuildRefusalsCrossProcessBoundary` | for each row, invoke through a fresh subprocess and assert the run record on disk is re-read rather than inherited |
