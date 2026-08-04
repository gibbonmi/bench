# Surface the interrupted reclaim receipt at the command surface

Blocked by: none
Ownership fence: `cmd/bench/specbuild.go`, `internal/contract/runtime/runtime_spec_build_test.go`
Contracts: the partial `ReclamationPlan` receipt crosses `internal/specbuild/reclaim.go`→`cmd/bench/specbuild.go` carrying the refs an interrupted apply spent, their objects, their dispositions, and the receipt's own fingerprint, and is asserted by IR1 and IR3 against the real `bench spec build reclaim --apply` rather than against the core function alone; absence semantics are that a plan-time failure carries no receipt and must stay an error alone, asserted by IR4
Assumptions: `ApplyReclaim` already returns the spent-set receipt alongside its error and this ticket does not change that contract; the fault seam `reclaim/apply` already exists for injecting the interruption; claims re-derived from the tree at pickup

## What to build

`ApplyReclaim` is deliberate about an interrupted apply. Git deletes one ref at a time with
no transaction across the set, so the function makes the outcome honest rather than atomic:
it returns a receipt reporting exactly the refs that are gone, carrying the fingerprint of
that set alone — never the spent plan's, which no longer describes anything — so the
maintainer can see what was spent and converge on a fresh plan.

The command surface throws that receipt away. In `cmd/bench/specbuild.go`'s `reclaim` arm,
a non-nil error routes to `buildError(err, hint)` and the `plan` value is discarded
unread. What survives is the ref names that `ApplyReclaim` formatted into its error string
with `deletedNames`. So the structured receipt — each ref's object, its disposition, and
the receipt fingerprint of the spent set — never reaches the operator, and the only
surviving account of what was deleted is prose inside an error message. That prose is a
second derivation of a fact the receipt already holds, which is exactly the shape this
repo's code standard forbids: when the two disagree, the operator is reading the one that
is not the receipt.

Render the receipt at the surface when an interrupted apply returns one, alongside the
error and its hint, so the maintainer sees the spent set as structured output rather than
as a sentence. The error string then no longer needs to restate the ref names, but that
restatement lives in `internal/specbuild/reclaim.go` and is `collapse-the-spent-set-restatement.md`'s
to remove — this ticket may not reach it. A failure that produced no receipt
at all — a refused fingerprint, a drifted plan, a run that does not exist — stays an error
alone, because inventing an empty receipt would report a deletion set for an operation that
deleted nothing.

## Acceptance

- [ ] [IR1] an apply interrupted partway through deletion prints the receipt naming exactly the refs it spent, driven through the real `bench spec build reclaim --apply` with the fault injected at the existing `reclaim/apply` seam.
- [ ] [IR2] that output also carries the failure and its convergence hint, so the receipt does not read as success.
- [ ] [IR3] the printed receipt carries the fingerprint of the spent set, not the fingerprint of the plan that was requested.
- [ ] [IR4] a failure that produced no receipt — a malformed or drifted fingerprint — prints the error alone and no receipt table.
- [ ] [IR6] the command's exit status stays non-zero for every interrupted or refused apply, so the receipt never reads as a completed operation.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| IR1 | discard the receipt on error, as the tree does today | the interrupted-apply runtime test | drop the receipt render from the reclaim arm, run `go test ./internal/contract/runtime -run SpecBuild -count=1 -timeout 600s`, expect the spent-ref assertion to fail |
| IR2 | render the receipt and swallow the error | the interrupted-apply runtime test | return the receipt with exit zero and no error text, run `go test ./internal/contract/runtime -run SpecBuild -count=1 -timeout 600s`, expect the failure-text assertion to fail |
| IR3 | render the requested plan's fingerprint instead of the receipt's | the spent-set fingerprint test | substitute the plan fingerprint, run `go test ./internal/contract/runtime -run SpecBuild -count=1 -timeout 600s`, expect the fingerprint assertion to fail |
| IR4 | render an empty receipt table for every failure | the no-receipt refusal test | render unconditionally, run `go test ./internal/contract/runtime -run SpecBuild -count=1 -timeout 600s`, expect the no-table assertion to fail |
| IR6 | return exit zero when a receipt is present | the exit-status test | return zero with the receipt, run `go test ./internal/contract/runtime -run SpecBuild -count=1 -timeout 600s`, expect the non-zero-exit assertion to fail |
