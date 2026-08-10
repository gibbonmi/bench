# Migrate outline aggregates

Blocked by: none
Ownership fence: `internal/outline`
Integration surfaces: ordered aggregate carrier→`internal/axi/aggregate.go` exercised by OA1; outline scan producer→`internal/outline/outline.go` exercised by OA1; `outline_meta` renderer→`internal/outline/outline.go` exercised by OA1; legacy carrier contraction→contract-aggregate-empty-routes.md
Contracts: the seven `outline_meta` facts cross `internal/outline/outline.go`→`internal/axi/aggregate.go`; six are non-negative decimal counts and `truncated` is a Go-formatted bool; order is tracked_files, scanned_files, skipped_files, total_symbols, emitted_symbols, omitted_symbols, truncated; a zero count is always emitted rather than dropped, asserted by OA1 against the real `Command` producer
Closure: OA1/tracked, OA1/scanned, OA1/skipped, OA1/total, OA1/emitted, OA1/omitted, OA1/truncated, OA1/order, OA1/route

## What to build

`bench outline` supplies its already-derived scan facts to the shared ordered aggregate
carrier and renders the identical `outline_meta[1]{...}` block. `Command` keeps every
derivation: `len(files)`, `scanned`, `len(skips)`, `totalSymbols`, the `bounds.OutlineRowLimit`
row bound that fixes `emitted`, `omitted = totalSymbols - emitted`, and
`truncated = omitted > 0 || len(skips) > 0`.

AE2 is an already-covered row: `TestCommandBoundsRowsAndFullRetainsMetadata` and
`TestCommandNamesSizeBinaryAndNonregularSkips` stay exactly as they are and remain the
named existing controls. This ticket adds the subject mutation the row lacks — a new
`TestOutlineMetaAggregateRouteCarriesOwnerFacts` in `internal/outline` that drives the
real `Command` over a fixture repository and asserts the metadata block was produced from
the owner facts through the shared carrier, so a byte-equal local re-render is red.

Tree condition that must hold when this ticket is refreshed: `internal/axi/aggregate.go`
exists and declares the exported ordered-aggregate type `Aggregate` with its typed fact
entry `Fact`. If that path or either symbol is absent, stop and report rather than build —
the prerequisite `axi-carriers-and-registry` build has not landed.

## Acceptance

- [ ] [OA1] (covers AE2) `outline_meta` renders the owner's tracked, scanned, skipped, total, emitted, omitted, and truncated facts in that order through the shared aggregate carrier.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OA1/tracked | supply `scanned` in place of `len(files)` as the tracked_files fact | `TestOutlineMetaAggregateRouteCarriesOwnerFacts` (`internal/outline`, new) | run `go test ./internal/outline -run TestOutlineMetaAggregateRouteCarriesOwnerFacts -count=1 -timeout 120s`; the assertion that tracked_files is `2` for a repository holding one scanned `.go` file and one skipped binary file fails with `1`; the fixture is a `t.TempDir()` git repository of four small files and every file read is capped by `bounds.OutlineFileLimit` (2 MiB), so no read can grow unbounded |
| OA1/scanned | supply `len(files)` in place of the `scanned` counter | `TestCommandNamesSizeBinaryAndNonregularSkips` (`internal/outline`) | run `go test ./internal/outline -run TestCommandNamesSizeBinaryAndNonregularSkips -count=1 -timeout 120s`; extend nothing — the existing skip fixture has tracked files that are never scanned, so the scanned_files cell renders the tracked count and the metadata comparison fails; each of the four fixture files is bounded by `bounds.OutlineFileLimit` |
| OA1/skipped | supply `0` for skipped_files whenever the emitted rows are unbounded | `TestCommandNamesSizeBinaryAndNonregularSkips` (`internal/outline`) | run `go test ./internal/outline -run TestCommandNamesSizeBinaryAndNonregularSkips -count=1 -timeout 120s`; the `outline_skips` rows still list `large.go,oversized` while skipped_files renders `0`, so the assertion pairing the skip count with the emitted skip rows fails; bounded by `bounds.OutlineFileLimit` on the 2 MiB fixture files |
| OA1/total | supply `emitted` (the post-bound row count) as total_symbols | `TestCommandBoundsRowsAndFullRetainsMetadata` (`internal/outline`) | run `go test ./internal/outline -run TestCommandBoundsRowsAndFullRetainsMetadata -count=1 -timeout 120s`; the exact metadata row assertion `"1","1","0","201","200","1","true"` fails because total_symbols renders `200`; the fixture writes one 201-symbol file under `bounds.OutlineFileLimit` and the default run is bounded to `bounds.OutlineRowLimit` rows |
| OA1/emitted | supply `totalSymbols` as emitted_symbols | `TestCommandBoundsRowsAndFullRetainsMetadata` (`internal/outline`) | run `go test ./internal/outline -run TestCommandBoundsRowsAndFullRetainsMetadata -count=1 -timeout 120s`; the same exact metadata row assertion fails because emitted_symbols renders `201` while the table header is `outline[200]{`; bounded by `bounds.OutlineRowLimit` |
| OA1/omitted | compute omitted as `totalSymbols - len(rows)` before the row bound is applied | `TestCommandBoundsRowsAndFullRetainsMetadata` (`internal/outline`) | run `go test ./internal/outline -run TestCommandBoundsRowsAndFullRetainsMetadata -count=1 -timeout 120s`; the bounded-run assertion fails because omitted_symbols renders `0` instead of `1`; bounded by `bounds.OutlineRowLimit` |
| OA1/truncated | compute truncated as `omitted > 0` alone, dropping the `len(skips) > 0` disjunct | `TestCommandNamesSizeBinaryAndNonregularSkips` (`internal/outline`) | run `go test ./internal/outline -run TestCommandNamesSizeBinaryAndNonregularSkips -count=1 -timeout 120s`; with every symbol emitted but three files skipped, the truncated cell renders `false`, failing the assertion that a skipped input keeps truncated `true`; bounded by `bounds.OutlineFileLimit` |
| OA1/order | supply the aggregate facts as tracked, scanned, total, skipped, emitted, omitted, truncated | `TestCommandBoundsRowsAndFullRetainsMetadata` (`internal/outline`) | run `go test ./internal/outline -run TestCommandBoundsRowsAndFullRetainsMetadata -count=1 -timeout 120s`; the header assertion `tracked_files,scanned_files,skipped_files,total_symbols,emitted_symbols,omitted_symbols,truncated` fails against the reordered header; bounded by `bounds.OutlineRowLimit` |
| OA1/route | keep the pre-migration local `toon.Table("outline_meta", ...)` call and never construct the shared aggregate | `TestOutlineMetaAggregateRouteCarriesOwnerFacts` (`internal/outline`, new) | run `go test ./internal/outline -run TestOutlineMetaAggregateRouteCarriesOwnerFacts -count=1 -timeout 120s`; the assertion that the metadata block was constructed through `axi.Aggregate` from the seven owner facts fails with no aggregate observed, even though the rendered bytes are unchanged; bounded by `bounds.OutlineFileLimit` and `bounds.OutlineRowLimit` on the fixture repository |
