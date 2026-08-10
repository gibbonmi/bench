# Introduce the bounded projection carrier

Blocked by: none
Ownership fence: `internal/axi/projection.go`, `internal/axi/projection_test.go`
Integration surfaces: projection API→`internal/axi/projection.go`; registry declarations→declare-axi-query-root-metadata.md; compatibility oracle package root→`internal/axi/compatibility` exercised unchanged (the sibling spec owns it; this ticket adds only root-package files beside it)
Contracts: selected content, integer total/emitted/omitted, boolean truncated, completeness enum, and counting-unit enum cross caller→`internal/axi/projection.go`, the owner declares every value, field order is fixed, and unknown completeness is explicit, asserted by PC1 against real owner-supplied facts rather than a recomputing stub
Closure: PC1/content, PC1/total, PC1/emitted, PC1/omitted, PC1/truncated, PC1/completeness, PC1/unit

## What to build

`internal/axi` gains a bounded projection carrier that stores what the owner
supplies and infers nothing: the selected content, the total, the emitted count,
the omitted count, the truncated flag, a completeness enum with an explicit
unknown, and the counting unit (bytes, code points, rows). Bench has four
independent truncation policies today with different units, so any inference —
total from emitted, omitted from `total - emitted`, truncated from omitted, or a
normalized unit — is exactly the drift this carrier exists to prevent.

Tree condition at refresh time: this spec follows `axi-compatibility-oracle`, so
`internal/axi/compatibility` already exists as a sibling package directory. This
ticket writes only the two root-package files on its fence line and must not
touch anything under `internal/axi/compatibility`.

## Acceptance

- [ ] [PC1] (covers CR3) projections return the owner's selected content, total, emitted, omitted, truncated flag, completeness, and counting unit without inferring any of them.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC1/content | return a recomputed prefix of the source instead of the owner's selected content | `TestProjectionReturnsOwnerSelectedContent` in `internal/axi` (this ticket authors it) | run `go test ./internal/axi -run TestProjectionReturnsOwnerSelectedContent -timeout 60s`; expect the equality assertion `Content() = "abc", want "a…c"` for the owner-selected elided value; bound is the `-timeout 60s` binary deadline over in-memory values |
| PC1/total | compute `total` from the emitted count | `TestProjectionKeepsOwnerTotalIndependentOfEmitted` in `internal/axi` | run `go test ./internal/axi -run TestProjectionKeepsOwnerTotalIndependentOfEmitted -timeout 60s`; expect the equality assertion `Total() = 3, want 97` for the projection whose owner declares 97 total and 3 emitted; bound is the `-timeout 60s` binary deadline |
| PC1/emitted | compute `emitted` as `len(content)` in bytes | `TestProjectionKeepsOwnerEmittedUnderDeclaredUnit` in `internal/axi` | run `go test ./internal/axi -run TestProjectionKeepsOwnerEmittedUnderDeclaredUnit -timeout 60s`; expect the equality assertion `Emitted() = 9, want 3` for the three-code-point multibyte content the owner counted in code points; bound is the `-timeout 60s` binary deadline |
| PC1/omitted | compute `omitted` as `total - emitted` | `TestProjectionKeepsOwnerOmittedCount` in `internal/axi` | run `go test ./internal/axi -run TestProjectionKeepsOwnerOmittedCount -timeout 60s`; expect the equality assertion `Omitted() = 94, want 12` for the owner that omits only the twelve rows it filtered; bound is the `-timeout 60s` binary deadline |
| PC1/truncated | derive `truncated` from `omitted > 0` | `TestProjectionKeepsOwnerTruncatedFlagIndependent` in `internal/axi` | run `go test ./internal/axi -run TestProjectionKeepsOwnerTruncatedFlagIndependent -timeout 60s`; expect the equality assertion `Truncated() = true, want false` for the projection that omitted filtered rows without hitting any cap; bound is the `-timeout 60s` binary deadline |
| PC1/completeness | map the unknown completeness value onto `complete=false` | `TestProjectionKeepsUnknownCompletenessExplicit` in `internal/axi` | run `go test ./internal/axi -run TestProjectionKeepsUnknownCompletenessExplicit -timeout 60s`; expect the equality assertion `Completeness() = axi.CompletenessPartial, want axi.CompletenessUnknown`; bound is the `-timeout 60s` binary deadline |
| PC1/unit | collapse the byte and code-point units onto one normalized `characters` value | `TestProjectionKeepsCountingUnitsDistinct` in `internal/axi` | run `go test ./internal/axi -run TestProjectionKeepsCountingUnitsDistinct -timeout 60s`; expect the inequality assertion `UnitBytes == UnitCodePoints` to report both projections reporting `characters`, want `bytes` and `code-points`; bound is the `-timeout 60s` binary deadline |
