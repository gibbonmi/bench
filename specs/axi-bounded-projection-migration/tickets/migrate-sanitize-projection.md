# Migrate sanitize projection

Blocked by: none
Ownership fence: `internal/sanitize`
Integration surfaces: shared bounded-projection carrier→`internal/axi` exported symbol `axi.Projection` (existing at refresh time), exercised by SN1; final contraction→contract-projection-routes.md
Contracts: preview content string, original byte integer, truncated boolean, code-point unit, and controls disposition cross `internal/sanitize/sanitize.go`→`internal/axi`, membership is the two entry points `Preview` and `Controls`, order is select-escape-suffix, and absence is empty input returning empty content with total 0 and truncated false, asserted by SN1 against the real `internal/axi` carrier
Closure: SN1/cap-120, SN1/code-points, SN1/original-bytes, SN1/controls, SN1/backslash, SN1/suffix, SN1/uncapped-controls, SN1/route-preview

## What to build

`sanitize.Preview` selects its first 120 code points through the shared `internal/axi` bounded-projection carrier instead of its local `runes[:120]` slice, while every observable byte it emits stays identical: the same `writeEscaped` control/backslash escaping, the same `"… (%d bytes)"` suffix naming `len(value)` of the original string, and the same uncapped `Controls`/`Preformatted` paths that never project at all.

Tree condition at refresh time: the package directory `internal/axi` exists and exports `axi.Projection`, the CR3 carrier holding selected content plus owner-supplied total, emitted, omitted, truncated, completeness, and counting unit without inferring any of them (the symbol name is fixed by the carriers spec). Confirm with `go doc ./internal/axi Projection` before starting; if the package or symbol is absent, stop and report rather than inventing a local shim.

The two entry points stay in one ticket because there is no thinner cut to attempt: `Preview` is a single function body, and `Controls` changes not at all — a per-fact split would produce tickets with no independently landable production change, so the split test has no candidate rather than a stranded red.

## Acceptance

- [ ] [SN1] (covers BP1) `Preview` selects 120 code points through the shared `internal/axi` projection while its escaping, byte suffix, and uncapped `Controls` path emit byte-identical output.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SN1/cap-120 | change the owner cap `Preview` supplies to the shared projection from 120 to 60 | existing `internal/sanitize`.`TestPreviewBoundariesAndControls` | run `go test ./internal/sanitize -run TestPreviewBoundariesAndControls -count=1 -timeout 60s`; the equality `Preview(strings.Repeat("é",121)) != strings.Repeat("é",120)+"… (242 bytes)"` fires because the result carries a 60-rune prefix; `-timeout 60s` bounds the in-process case |
| SN1/code-points | supply the shared projection a byte counting unit instead of code points, so the cap slices `[]byte(value)[:120]` | existing `internal/sanitize`.`TestPreviewBoundariesAndControls` | run `go test ./internal/sanitize -run TestPreviewBoundariesAndControls -count=1 -timeout 60s`; the equality `Preview(strings.Repeat("é",120)) != strings.Repeat("é",120)` fires because a 240-byte string is cut at 120 bytes; `-timeout 60s` bounds the in-process case |
| SN1/original-bytes | supply the projection total as the emitted selection length rather than `len(value)`, so the suffix reports the truncated length | existing `internal/sanitize`.`TestPreviewBoundariesAndControls` | run `go test ./internal/sanitize -run TestPreviewBoundariesAndControls -count=1 -timeout 60s`; the same 121-rune equality fires because the suffix reads `… (240 bytes)` instead of `… (242 bytes)`; `-timeout 60s` bounds the in-process case |
| SN1/controls | delete the `unicode.IsControl(r)` arm from `writeEscaped` so U+001B and U+0007 pass through raw | existing `internal/sanitize`.`TestPreviewBoundariesAndControls` | run `go test ./internal/sanitize -run TestPreviewBoundariesAndControls -count=1 -timeout 60s`; `containsControl(got)` reports a raw control byte and the `bs+"u001b"` and `bs+"u0007"` token-presence assertions fire; `-timeout 60s` bounds the in-process case |
| SN1/backslash | delete the `case r == '\\'` arm from `writeEscaped` so a literal backslash is not doubled | existing `internal/sanitize`.`TestPreviewBoundariesAndControls` | run `go test ./internal/sanitize -run TestPreviewBoundariesAndControls -count=1 -timeout 60s`; the equality `Preview("back"+bs+"slash") != "back"+bs+bs+"slash"` fires; `-timeout 60s` bounds the in-process case |
| SN1/suffix | emit the truncation suffix as `"..."` without the `(%d bytes)` clause | existing `internal/sanitize`.`TestPreviewBoundariesAndControls` | run `go test ./internal/sanitize -run TestPreviewBoundariesAndControls -count=1 -timeout 60s`; the 121-rune equality fires because the tail is `...` instead of `… (242 bytes)`; `-timeout 60s` bounds the in-process case |
| SN1/uncapped-controls | route `Controls` through the same shared projection with the 120-code-point cap `Preview` supplies | existing `internal/sanitize`.`TestControlsEscapesWithoutCapping` | run `go test ./internal/sanitize -run TestControlsEscapesWithoutCapping -count=1 -timeout 60s`; the equality `Controls(strings.Repeat("é",200)) != strings.Repeat("é",200)` fires with `len=120`; `-timeout 60s` bounds the in-process case |
| SN1/route-preview | restore the local `runes = runes[:120]` slice in `Preview` and drop the shared-projection call | authored `internal/sanitize`.`TestPreviewRoutesThroughSharedProjection` (this ticket writes it in `internal/sanitize/projection_route_test.go`) | run `go test ./internal/sanitize -run TestPreviewRoutesThroughSharedProjection -count=1 -timeout 60s`; the assertion that the projection carrier returned for a 121-rune input reports unit=code-points, total=242, emitted=120, truncated=true fires because no carrier is constructed at all (nil carrier / zero unit); `-timeout 60s` bounds the in-process case |
