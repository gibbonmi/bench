# Close the required argv classes over every census member

Blocked by: derive-nested-grammar-membership.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`
Integration surfaces: root, wrapper, and nested census→derive-nested-grammar-membership.md; class matrix and not-applicable reasons→`internal/axi/compatibility`; declared class marks per command→`decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md` exercised unchanged by every CA1 row; comparator consumer→compare-four-observations.md
Contracts: the class matrix — one `{case ID, class, argv, applicability}` row per member and class — crosses `decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md`→`internal/axi/compatibility`; its type is a matrix keyed by case ID, membership is every census member crossed with the five declared case classes, order is case ID then the inventory's class order, and an inapplicable pair carries a non-empty reason string rather than being absent; asserted by CA1 against the real inventory tables and the derived census
Closure: CA1/class-success, CA1/class-empty, CA1/class-refusal, CA1/class-stale, CA1/class-recovery, CA1/not-applicable-reason, CA1/inventory-row-totality, CA1/census-member-totality

## What to build

The census now names every member; this ticket says what must be compared for each
one. The required argv classes are the five case marks the FT173 command-help
inventory declares — `S` success, `E` empty or no-op, `R` refusal or usage, `St`
stale or drift, and `Rc` recovery — and every member carries, for each class,
either a concrete argv vector under case ID `<member>-<class>` or a non-empty
reason naming why that class cannot exist for it (`bench tree-hash` has no empty
projection; `bench version` has no recovery state).

Closure runs both ways. Every row of every inventory table must resolve to a census
member, so a command the inventory names cannot be silently dropped, and every
census member must carry a class or a reason, so a member the registry declares
cannot sit in the index with nothing to compare. Both directions are asserted from
the parsed inventory tables and the derived census, never from a hand-written list.

## Acceptance

- [ ] [CA1] (covers CO3) every census member carries an argv vector for each declared case class it marks or a non-empty not-applicable reason, every inventory row resolves to a census member, and every census member resolves to at least one class or reason.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CA1/class-success | drop the `S` class from the matrix rows built for `bench version`, whose inventory row marks `S` | the class-matrix test | run `go test ./cmd/bench -run TestCaseIndexBindsEveryApplicableClass/success -timeout 300s`; it fails at the missing-pair assertion naming member `root-version` and class `S` against the inventory mark, and reporting no case ID `root-version-s`; the inventory parse and each dispatcher call run under a 15s `context.WithTimeout` |
| CA1/class-empty | drop the `E` class from the matrix rows built for `bench learnings`, whose inventory row marks `S,E,R` | the class-matrix test | run `go test ./cmd/bench -run TestCaseIndexBindsEveryApplicableClass/empty -timeout 300s`; it fails at the missing-pair assertion naming member `root-learnings` and class `E`, reporting the inventory marks `S,E,R` against the derived `S,R`; bounded by the 15s deadline |
| CA1/class-refusal | drop the `R` class from the matrix rows built for `bench spec build` with no operation, whose inventory row marks `R` only | the class-matrix test | run `go test ./cmd/bench -run TestCaseIndexBindsEveryApplicableClass/refusal -timeout 300s`; it fails at the missing-pair assertion naming member `nested-spec-build-home` and class `R`, leaving that member with an empty class set; bounded by the 15s deadline |
| CA1/class-stale | drop the `St` class from the matrix rows built for `bench diff`, whose inventory row marks `S,E,R,St` | the class-matrix test | run `go test ./cmd/bench -run TestCaseIndexBindsEveryApplicableClass/stale -timeout 300s`; it fails at the missing-pair assertion naming member `root-diff` and class `St`, reporting the mid-read drift argv that no longer has a case ID; bounded by the 15s deadline |
| CA1/class-recovery | drop the `Rc` class from the matrix rows built for `bench worktree recovery`, whose inventory row marks `S,E,R,St,Rc` | the class-matrix test | run `go test ./cmd/bench -run TestCaseIndexBindsEveryApplicableClass/recovery -timeout 300s`; it fails at the missing-pair assertion naming member `nested-worktree-recovery` and class `Rc`; bounded by the 15s deadline |
| CA1/not-applicable-reason | accept an empty string as the not-applicable reason for `bench tree-hash` and the `E` class | the class-matrix test | run `go test ./cmd/bench -run TestInapplicableClassCarriesAConcreteReason -timeout 300s`; it fails at the non-empty-reason assertion naming member `root-tree-hash` and class `E`, reporting an empty reason where the matrix requires prose naming the missing projection; bounded by the 15s deadline |
| CA1/inventory-row-totality | build the matrix from the first inventory table only, skipping the adoption, gate/worktree, spec-build, release, and plumbing tables | the inventory-totality test | run `go test ./cmd/bench -run TestEveryInventoryRowResolvesToACensusMember -timeout 300s`; it fails at the unresolved-row assertion listing the first unresolved inventory row `bench setup [--plan or --yes]` and the count of remaining unresolved rows; the inventory parse runs under the 15s deadline |
| CA1/census-member-totality | allow a census member to carry neither a class nor a reason | the census-totality test | run `go test ./cmd/bench -run TestEveryCensusMemberCarriesAClassOrReason -timeout 300s`; it fails at the empty-matrix-entry assertion naming member `root-freshness-publish`, which the inventory marks `S,R,St` while the wrapper routes it only through the unrecognized-token arm; bounded by the 15s deadline |
