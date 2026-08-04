# Collapse the spent-set restatement in the reclaim refusal

Blocked by: none
Ownership fence: `internal/specbuild`
Contracts: the spent-set account crosses `internal/specbuild/reclaim.go`→`cmd/bench/specbuild.go`, which now renders the returned receipt, and is asserted by SS1 against the real `ApplyReclaim` return rather than against the error string alone; absence semantics are that a plan-time failure returns no receipt and its refusal text is unchanged, asserted by SS3
Assumptions: `cmd/bench` already renders the interrupted-apply receipt, so removing the restatement costs the operator no information; the refusal's drifted-ref name and its convergence hint stay, because neither is carried by the receipt; claims re-derived from the tree at pickup

## What to build

This closes row IR5 of `surface-the-interrupted-reclaim-receipt.md`, which could not be
built from that ticket's fence — the restatement it names lives in this package, and the
delegate correctly stopped rather than reaching across. That fence was mis-drawn at
slicing time; this ticket is the correction, not new scope.

`ApplyReclaim`'s interrupted-apply refusal formats the deleted refs into its error text:

```go
return deleted, fmt.Errorf("reclaim %s: %w; deleted %s; request a fresh plan to finish the remainder", ref.Name, err, deletedNames(deleted))
```

The returned receipt already carries exactly that set, with each ref's object and
disposition, and the command surface now renders it. So the operator reads the same list
twice, from two derivations that can disagree — and when they do, the one that is not the
receipt is the one that lies. That is the duplicated-knowledge defect the code standard
forbids.

Drop `; deleted %s` and its argument. The refusal keeps the ref that drifted and the
convergence hint, because neither is in the receipt. `deletedNames` then has no caller and
goes with it, rather than surviving as a helper nothing derives anything from.

## Acceptance

- [ ] [SS1] an interrupted apply's error text names the ref that drifted and does not restate the deleted set, while the returned receipt still carries that set in full.
- [ ] [SS2] the refusal keeps its convergence hint, so a maintainer still learns a fresh plan covers the remainder.
- [ ] [SS3] a plan-time failure — a drifted or malformed fingerprint — keeps its existing refusal text and still returns no receipt.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SS1 | restore `; deleted %s` and its `deletedNames` argument in the interrupted-apply refusal | the spent-set single-source test | re-add the clause, run `go test ./internal/specbuild -run Reclaim -count=1 -timeout 900s`, expect the no-restatement assertion to fail |
| SS2 | drop the convergence hint along with the restatement | the refusal-hint test | remove the hint text, run `go test ./internal/specbuild -run Reclaim -count=1 -timeout 900s`, expect the hint assertion to fail |
| SS3 | return an empty receipt from the drifted-fingerprint path | the plan-time refusal test | return a zero receipt with the error, run `go test ./internal/specbuild -run Reclaim -count=1 -timeout 900s`, expect the no-receipt assertion to fail |
