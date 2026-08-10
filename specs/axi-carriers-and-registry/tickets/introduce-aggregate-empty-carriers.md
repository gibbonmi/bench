# Introduce aggregate and empty carriers

Blocked by: none
Ownership fence: `internal/axi/aggregate.go`, `internal/axi/aggregate_test.go`, `internal/axi/empty.go`, `internal/axi/empty_test.go`
Integration surfaces: aggregate API→`internal/axi/aggregate.go`; empty-class API→`internal/axi/empty.go`; registry declarations→declare-axi-query-root-metadata.md; compatibility oracle package root→`internal/axi/compatibility` exercised unchanged (the sibling spec owns it; this ticket adds only root-package files beside it)
Contracts: ordered name plus typed scalar and the empty-class enum cross caller→`internal/axi/aggregate.go` and `internal/axi/empty.go`, scalar domain is string/number/bool/unknown, order is insertion order, and absence is an explicit no-semantic-empty class, asserted by AE1 against real owner-supplied facts rather than a normalizing stub
Closure: AE1/name, AE1/order, AE1/type, AE1/unknown, AE1/zero, AE1/empty-class, AE1/no-empty

## What to build

`internal/axi` gains two carriers. The aggregate carrier holds an ordered list of
owner-named typed scalars: it keeps the owner's name, keeps insertion order,
keeps the scalar's Go type distinct (number stays a number, not a rendered
string), keeps `unknown` distinct from zero, and never drops a fact whose value
is zero. The empty carrier holds an owner-declared empty class — Bench today has
a zero-row table, a one-row `state=empty` projection, and prose empties such as
`No gate cache`, which are different declarations, not variants to normalize —
plus an explicit `no semantic empty` class for a surface that has none.

Tree condition at refresh time: this spec follows `axi-compatibility-oracle`, so
`internal/axi/compatibility` already exists as a sibling package directory. This
ticket writes only the four root-package files on its fence line and must not
touch anything under `internal/axi/compatibility`.

## Acceptance

- [ ] [AE1] (covers CR4) aggregates keep owner names, insertion order, scalar types, unknowns, and explicit zeros, and empty declarations keep each owner-declared class distinct from an explicit no-semantic-empty.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AE1/name | return the aggregate's index position in place of its owner-supplied name | `TestAggregatePreservesOwnerName` in `internal/axi` (this ticket authors it) | run `go test ./internal/axi -run TestAggregatePreservesOwnerName -timeout 60s`; expect the equality assertion `Facts()[0].Name = "1", want "open_learnings"`; bound is the `-timeout 60s` binary deadline over in-memory values |
| AE1/order | sort the fact list by name before returning it | `TestAggregatePreservesInsertionOrder` in `internal/axi` | run `go test ./internal/axi -run TestAggregatePreservesInsertionOrder -timeout 60s`; expect the equality assertion `Facts()[0].Name = "alpha", want "zulu"` for facts inserted in reverse-alphabetical order; bound is the `-timeout 60s` binary deadline |
| AE1/type | store every scalar as its `fmt.Sprint` string | `TestAggregateKeepsScalarTypesDistinct` in `internal/axi` | run `go test ./internal/axi -run TestAggregateKeepsScalarTypesDistinct -timeout 60s`; expect the type assertion `Facts()[0].Value.(int)` to fail with `value has dynamic type string`; bound is the `-timeout 60s` binary deadline |
| AE1/unknown | coerce the unknown scalar to the zero value of its declared type | `TestAggregateKeepsUnknownDistinctFromZero` in `internal/axi` | run `go test ./internal/axi -run TestAggregateKeepsUnknownDistinctFromZero -timeout 60s`; expect the assertion `Facts()[0].Kind = axi.ScalarNumber, want axi.ScalarUnknown` while the sibling zero fact still reports `axi.ScalarNumber`; bound is the `-timeout 60s` binary deadline |
| AE1/zero | skip any fact whose numeric value is zero when building the list | `TestAggregateKeepsExplicitZeroFact` in `internal/axi` | run `go test ./internal/axi -run TestAggregateKeepsExplicitZeroFact -timeout 60s`; expect the length assertion `len(Facts()) = 1, want 2` for the pair `{"blocked": 0, "open": 3}`; bound is the `-timeout 60s` binary deadline |
| AE1/empty-class | map every declared empty class onto one `zero-rows` value | `TestEmptyClassesStayDistinct` in `internal/axi` | run `go test ./internal/axi -run TestEmptyClassesStayDistinct -timeout 60s`; expect the inequality assertion to report the zero-row-table and prose-empty declarations both reading `zero-rows`, want `zero-rows` and `prose`; bound is the `-timeout 60s` binary deadline |
| AE1/no-empty | default a surface that declares no semantic empty to the zero-row class | `TestEmptyRequiresExplicitNoSemanticEmpty` in `internal/axi` | run `go test ./internal/axi -run TestEmptyRequiresExplicitNoSemanticEmpty -timeout 60s`; expect the assertion `Class() = axi.EmptyZeroRows, want axi.EmptyNone` for the surface constructed with the explicit no-semantic-empty declaration, and the `errors.Is(err, axi.ErrEmptyClassMissing)` assertion to report a nil error for the unset one; bound is the `-timeout 60s` binary deadline |
