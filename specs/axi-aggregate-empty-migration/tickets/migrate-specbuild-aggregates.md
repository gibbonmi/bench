# Migrate spec-build abandonment and reclamation aggregates

Blocked by: none
Ownership fence: `internal/specbuild`, `cmd/bench/specbuild.go`, `cmd/bench/specbuild_test.go`
Integration surfaces: ordered aggregate carrier→`internal/axi/aggregate.go` exercised by SA1; abandonment plan producer and fingerprint preimage→`internal/specbuild/lifecycle.go` exercised by SA1; reclamation plan producer, class inventory, and fingerprint preimage→`internal/specbuild/reclaim.go` exercised by SA1; abandonment and reclamation renderers→`cmd/bench/specbuild.go` exercised by SA1; empty-state sibling→migrate-specbuild-empty-state.md; legacy carrier contraction→contract-aggregate-empty-routes.md
Contracts: the abandonment counts (`worktrees`, `provisional_refs`, `checkpoints`, `recovery_refs`) with the plan `Fingerprint`, and the `ReclamationPlan.Classes()` name/count pairs with the reclamation `Fingerprint`, cross `internal/specbuild`→`internal/axi/aggregate.go` and on to `cmd/bench/specbuild.go`; counts are non-negative decimals, class names are the `refReclaimable`/`refActive`/`refUnclassified`/`refAmbiguous` literals, and fingerprints are lowercase 64-hex digests of the `abandonmentFacts`/`reclamationFacts` preimages; order is the declared class order with every zero retained; asserted by SA1 against the real `Abandon` and `Reclaim` producers
Closure: SA1/abandon-worktrees, SA1/abandon-provisional, SA1/abandon-checkpoints, SA1/abandon-recovery, SA1/abandon-fingerprint-preimage, SA1/class-reclaimable, SA1/class-active, SA1/class-unclassified, SA1/class-ambiguous, SA1/class-order, SA1/reclaim-fingerprint-preimage, SA1/route

## What to build

`bench spec build abandon` and `bench spec build reclaim` supply their already-derived
lifecycle counts to the shared ordered aggregate carrier and render identical bytes.
`AbandonmentPlan` keeps ownership of its four counts and of the fingerprint
`abandonmentFacts` commits to; `ReclamationPlan.Classes()` keeps ownership of the complete
class inventory — including the classes whose count is zero — and of the fingerprint
`reclamationFacts` commits to. The carrier transports ordered typed facts and never
recomputes a count from the rendered `reclaim_refs`/`abandon_worktrees` rows, and never
re-derives either fingerprint.

AE5 is an already-covered row: `TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes`
stays exactly as it is and remains the named existing control. This ticket adds the subject
mutation the row lacks — a new `TestSpecBuildReclamationAggregateRouteCarriesOwnerClasses`
in `internal/specbuild` asserting the rendered receipt was produced from the plan's own
class inventory through the shared carrier.

This ticket is deliberately separate from `migrate-specbuild-empty-state.md`: the
abandonment/reclamation count aggregates and the `state=empty` status projection land green
independently, and no project-gate red spans the two.

Tree condition that must hold when this ticket is refreshed: `internal/axi/aggregate.go`
exists and declares the exported ordered-aggregate type `Aggregate` with its typed fact
entry `Fact`. If that path or either symbol is absent, stop and report rather than build —
the prerequisite `axi-carriers-and-registry` build has not landed.

## Acceptance

- [ ] [SA1] (covers AE5) abandonment and reclamation render every owner-derived count, including each zero-valued reclamation class, in the owner's order through the shared aggregate carrier, with both plan fingerprints still derived from their own preimages.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SA1/abandon-worktrees | render the abandonment `worktrees` count from `len(plan.Worktrees)` after the retained-worktree filter instead of the planned set | `TestSpecBuildAbandonmentAggregateRouteCarriesOwnerCounts` (`internal/specbuild`, new) | run `go test ./internal/specbuild -run TestSpecBuildAbandonmentAggregateRouteCarriesOwnerCounts -count=1 -timeout 300s`; the assertion that the `abandon` row's `worktrees` cell equals the plan's own worktree count fails with the reduced number; the fixture is a `t.TempDir()` repository with two planned worktrees and every git call runs under the service's `ctx` deadline, so a wedged git subprocess ends the test |
| SA1/abandon-provisional | render `0` for `provisional_refs` whenever `plan.ProvisionalRefs` is non-empty but every ref is already deleted | `TestSpecBuildAbandonmentAggregateRouteCarriesOwnerCounts` (`internal/specbuild`, new) | run `go test ./internal/specbuild -run TestSpecBuildAbandonmentAggregateRouteCarriesOwnerCounts -count=1 -timeout 300s`; the assertion that `provisional_refs` equals `len(plan.ProvisionalRefs)` fails with `0`; bounded by the service context deadline on each `for-each-ref` call |
| SA1/abandon-checkpoints | render the `checkpoints` cell from the integrated checkpoint set instead of `plan.UnintegratedCheckpoints` | `TestSpecBuildAbandonmentAggregateRouteCarriesOwnerCounts` (`internal/specbuild`, new) | run `go test ./internal/specbuild -run TestSpecBuildAbandonmentAggregateRouteCarriesOwnerCounts -count=1 -timeout 300s`; the assertion that `checkpoints` equals the unintegrated-checkpoint count fails with the integrated count; bounded by the service context deadline |
| SA1/abandon-recovery | omit the `recovery_refs` fact from the abandonment aggregate when its count is `0` | `TestSpecBuildAbandonmentAggregateRouteCarriesOwnerCounts` (`internal/specbuild`, new) | run `go test ./internal/specbuild -run TestSpecBuildAbandonmentAggregateRouteCarriesOwnerCounts -count=1 -timeout 300s`; the header assertion on `abandon[1]{fingerprint,worktrees,provisional_refs,checkpoints,recovery_refs}` fails against a four-field header; bounded by the service context deadline |
| SA1/abandon-fingerprint-preimage | drop `worktree.OwnerFingerprint` from the `abandonmentFacts` preimage in `internal/specbuild/lifecycle.go` | `TestAbandonmentFingerprintCommitsToEveryPlannedFact` (`internal/specbuild`, new) | run `go test ./internal/specbuild -run TestAbandonmentFingerprintCommitsToEveryPlannedFact -count=1 -timeout 300s`; the assertion that two plans differing only in a worktree's owner fingerprint produce different `plan.Fingerprint` values fails with two equal digests, so an apply authorized against one owner would satisfy the other; bounded by the service context deadline |
| SA1/class-reclaimable | drop the `refReclaimable` entry from `ReclamationPlan.Classes()` | `TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes` (`internal/specbuild`) | run `go test ./internal/specbuild -run TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes -count=1 -timeout 300s`; the `reflect.DeepEqual` comparison against the four-element `want` slice fails with a three-element `got` missing `{reclaimable, 2}`; the reclamation fixture creates a fixed ref set and every git call runs under the service context deadline |
| SA1/class-active | drop the `refActive` entry from `ReclamationPlan.Classes()` when its count is `0` | `TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes` (`internal/specbuild`) | run `go test ./internal/specbuild -run TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes -count=1 -timeout 300s`; the `reflect.DeepEqual` comparison fails with `{active, 0}` absent — the exact zero-class omission the row exists to catch; bounded by the service context deadline |
| SA1/class-unclassified | drop the `refUnclassified` entry from `ReclamationPlan.Classes()` | `TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes` (`internal/specbuild`) | run `go test ./internal/specbuild -run TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes -count=1 -timeout 300s`; the `reflect.DeepEqual` comparison fails with `{unclassified, 1}` absent; bounded by the service context deadline |
| SA1/class-ambiguous | drop the `refAmbiguous` entry from `ReclamationPlan.Classes()` | `TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes` (`internal/specbuild`) | run `go test ./internal/specbuild -run TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes -count=1 -timeout 300s`; the `reflect.DeepEqual` comparison fails with `{ambiguous, 2}` absent; bounded by the service context deadline |
| SA1/class-order | return `Classes()` sorted by descending count instead of the declared class order | `TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes` (`internal/specbuild`) | run `go test ./internal/specbuild -run TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes -count=1 -timeout 300s`; the `reflect.DeepEqual` comparison fails on element order even though the same four pairs are present; bounded by the service context deadline |
| SA1/reclaim-fingerprint-preimage | drop `ref.Disposition` from the `reclamationFacts` preimage in `internal/specbuild/reclaim.go` | `TestReclamationFingerprintCommitsToEveryRefDisposition` (`internal/specbuild`, new) | run `go test ./internal/specbuild -run TestReclamationFingerprintCommitsToEveryRefDisposition -count=1 -timeout 300s`; the assertion that two plans over the same ref names and objects but different dispositions produce different `plan.Fingerprint` values fails with two equal digests, so an apply planned against a retained ref would authorize deleting it; bounded by the service context deadline |
| SA1/route | keep the pre-migration local `toon.Table("reclaim", headers, ...)` header/summary build in `renderReclamation` and never construct the shared aggregate | `TestSpecBuildReclamationAggregateRouteCarriesOwnerClasses` (`internal/specbuild`, new) | run `go test ./internal/specbuild -run TestSpecBuildReclamationAggregateRouteCarriesOwnerClasses -count=1 -timeout 300s`; the assertion that the receipt's class facts were carried by `axi.Aggregate` from `plan.Classes()` fails with no aggregate observed, even though the rendered receipt bytes are unchanged; bounded by the service context deadline |
