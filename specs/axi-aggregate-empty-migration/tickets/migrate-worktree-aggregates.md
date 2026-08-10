# Migrate worktree aggregates

Blocked by: none
Ownership fence: `internal/worktree`
Integration surfaces: ordered aggregate carrier→`internal/axi/aggregate.go` exercised by WA1; ignored-inventory producer and digest preimage→`internal/worktree/subshell.go` exercised by WA1; inventory summary renderer→`internal/worktree/classifier.go` exercised by WA1; legacy carrier contraction→contract-aggregate-empty-routes.md
Contracts: the `IgnoredInventory` count/bytes/shown/truncated facts and the `Digest` authority value cross `internal/worktree/subshell.go`→`internal/axi/aggregate.go`; count is a decimal or the literal `at-least=1001`, bytes and shown are decimals, truncated is a Go-formatted bool, and the digest is the lowercase hex sha256 of the `canonicalParts` preimage of `(name, mode, size)` per ignored path; order is count, bytes, shown, truncated; an empty safe inventory renders `count=0 bytes=0 shown=0 truncated=false`, asserted by WA1 against the real `inventoryIgnored` producer
Closure: WA1/count, WA1/at-least, WA1/bytes, WA1/shown, WA1/truncated, WA1/digest-preimage, WA1/digest-authority, WA1/order, WA1/route

## What to build

The worktree ignored inventory supplies its already-derived facts to the shared ordered
aggregate carrier and renders the identical `count=… bytes=… shown=… truncated=…` summary.
`inventoryIgnored` keeps every derivation: the `ignoredEntryLimit` (1000) at-least ceiling,
the `ignoredByteLimit` (1 GiB) over-limit flag, the 20-entry default display cap versus the
`--full` entry cap, and — above all — the `Digest` it computes with `fingerprintParts` over
the `(name, mode, size)` preimage. The carrier never recomputes or re-derives the digest;
it transports the value the owner produced.

The existing inventory tests are positive controls only, so this ticket adds the
shared-route mutation the row lacks: a new `TestIgnoredInventoryAggregateRouteCarriesOwnerFacts`
in `internal/worktree` drives the real `PlanExplicitWithOptions` and asserts the summary
facts reached the renderer through the shared carrier holding the producer's own values.

Tree condition that must hold when this ticket is refreshed: `internal/axi/aggregate.go`
exists and declares the exported ordered-aggregate type `Aggregate` with its typed fact
entry `Fact`. If that path or either symbol is absent, stop and report rather than build —
the prerequisite `axi-carriers-and-registry` build has not landed.

## Acceptance

- [ ] [WA1] (covers AE4) the ignored inventory renders the owner's count, bytes, shown, and truncated facts in that order through the shared aggregate carrier, and the inventory digest stays derived from the owner's own preimage.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| WA1/count | set `inventory.Count = inventory.Shown` after the display cap is applied | `TestIgnoredInventoryEntryAndByteBoundaries` (`internal/worktree`) | run `go test ./internal/worktree -run 'TestIgnoredInventoryEntryAndByteBoundaries/entries-21' -count=1 -timeout 300s`; the `plan.Ignored.Count != count` check fails with `20` against the fixture's 21 ignored entries; the fixture writes at most 1001 one-byte files and the walk stops at `ignoredEntryLimit`, so the enumeration cannot run away |
| WA1/at-least | render the numeric count instead of `at-least=1001` when `inventory.AtLeast` is set | `TestIgnoredInventoryEntryAndByteBoundaries` (`internal/worktree`) | run `go test ./internal/worktree -run 'TestIgnoredInventoryEntryAndByteBoundaries/entries-1001' -count=1 -timeout 300s`; the assertion that the over-ceiling plan takes `ActionRetain` with an `at-least=1001` count cell fails against a bare `1001`; the walk breaks at `ignoredEntryLimit` so the fixture never enumerates beyond 1001 entries |
| WA1/bytes | accumulate `inventory.Bytes` only for the entries below the 20-entry display cap | `TestIgnoredInventoryEntryAndByteBoundaries` (`internal/worktree`) | run `go test ./internal/worktree -run 'TestIgnoredInventoryEntryAndByteBoundaries/bytes-1073741825' -count=1 -timeout 300s`; the over-`ignoredByteLimit` case fails because `OverLimit` stays false when the visible-entry sum falls under 1 GiB; the fixture writes one sparse file at the byte boundary and the entry walk is capped at `ignoredEntryLimit` |
| WA1/shown | set `inventory.Shown = inventory.Count` before applying the `show` cap | `TestIgnoredInventoryEntryAndByteBoundaries` (`internal/worktree`) | run `go test ./internal/worktree -run 'TestIgnoredInventoryEntryAndByteBoundaries/entries-21' -count=1 -timeout 300s`; the `plan.Ignored.Shown != wantShown` check fails with `21` against the expected `20`; bounded by `ignoredEntryLimit` |
| WA1/truncated | compute `inventory.Truncated` as `inventory.AtLeast` alone, dropping the `Count > show` disjunct | `TestIgnoredInventoryEntryAndByteBoundaries` (`internal/worktree`) | run `go test ./internal/worktree -run 'TestIgnoredInventoryEntryAndByteBoundaries/entries-21' -count=1 -timeout 300s`; the `plan.Ignored.Truncated != (count > 20)` check fails with `false` at 21 entries; bounded by `ignoredEntryLimit` |
| WA1/digest-preimage | drop `strconv.FormatInt(info.Size(), 10)` from the `parts` triple appended per ignored path, so the digest preimage carries only name and mode | `TestIgnoredInventoryDigestCommitsToEverySizedEntry` (`internal/worktree`, new) | run `go test ./internal/worktree -run TestIgnoredInventoryDigestCommitsToEverySizedEntry -count=1 -timeout 300s`; the assertion that two inventories with identical path names and modes but different file sizes produce different `Ignored.Digest` values fails with two equal digests; the fixture holds two ignored files of a few bytes each and the walk is capped by `ignoredEntryLimit` |
| WA1/digest-authority | in the migrated route, recompute the digest inside the shared carrier from the rendered count/bytes/shown facts instead of transporting `inventory.Digest` | `TestIgnoredInventoryDigestCommitsToEverySizedEntry` (`internal/worktree`, new) | run `go test ./internal/worktree -run TestIgnoredInventoryDigestCommitsToEverySizedEntry -count=1 -timeout 300s`; the assertion that two inventories with equal count/bytes/shown but different path names produce different digests fails with two equal digests; bounded by `ignoredEntryLimit` |
| WA1/order | supply the aggregate facts as count, shown, bytes, truncated | `TestIgnoredInventoryAggregateRouteCarriesOwnerFacts` (`internal/worktree`, new) | run `go test ./internal/worktree -run TestIgnoredInventoryAggregateRouteCarriesOwnerFacts -count=1 -timeout 300s`; the exact-summary assertion on `count=1 bytes=1 shown=1 truncated=false` fails against the reordered `count=1 shown=1 bytes=1 truncated=false`; the fixture holds one ignored file and the walk is capped by `ignoredEntryLimit` |
| WA1/route | keep the pre-migration local `IgnoredInventory.Summary()` string build and never construct the shared aggregate | `TestIgnoredInventoryAggregateRouteCarriesOwnerFacts` (`internal/worktree`, new) | run `go test ./internal/worktree -run TestIgnoredInventoryAggregateRouteCarriesOwnerFacts -count=1 -timeout 300s`; the assertion that the summary was carried by `axi.Aggregate` from the producer's four owner facts fails with no aggregate observed, even though the rendered summary bytes are unchanged; bounded by `ignoredEntryLimit` |
