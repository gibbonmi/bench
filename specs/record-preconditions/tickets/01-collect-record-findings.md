# Collect every record finding from the checkpoint predicates

Blocked by: none
Writes: internal/reviewrecord/coverage.go, internal/reviewrecord/check.go, internal/reviewrecord/confirmation.go, internal/reviewrecord/fold.go, internal/reviewrecord/repair.go, internal/reviewrecord/ticket.go, internal/reviewrecord/collect.go (new), internal/reviewrecord/collect_test.go (new)
Covers: RP22

## What to build

Review chunk: RP-C1.

Add a collector to package `reviewrecord` that walks the predicates of `checkSource` (`internal/reviewrecord/coverage.go:17-173`) and yields every finding. A finding carries the chunk ID, the entry, the state (`missing` or `stale`), the error, and the `Repair` value when a form repairs it. Within one chunk, the identity group stops at its first failure: ancestry, source digest, plan identity, and acceptance rows. The collector then evaluates the verification, review, chain, and mapping groups of that chunk, and it continues with the next chunk. It takes the same `complete` and `verification` modes as `checkSource`.

Rewrite `checkSource` to return the error of the first collected finding, and keep its signature and its refusal order. Every caller keeps its text: `CheckSource`, `CheckTrees` (`check.go:16`), and `RecordCompletion` (`completion.go:37`). The collector reaches `checkVerification` and `checkCompletion` (`check.go`), `checkChunkReviews` (`confirmation.go:212`), `currentPlan` (`repair.go:55`), and `coveredAfter`, `uncoveredFold`, and `digestFollows` (`fold.go`). Change those callees only where a callee must yield more than its first error.

For a chain gap, the collector lists each commit on the gap range whose source digest does not follow its first parent's digest under `digestFollows`. For an uncovered reviewed source, it lists each such commit on the range from the covered tip to the requested tip. A default-branch merge that no recorded fold covers takes the reason `unrecorded fold` from `uncoveredFold`. No other code derives a chain, staleness, or base predicate.

Export one read-only accessor for the required base of a chunk from `chunkBase` (`ticket.go:174-188`), so that `recordcmd` reads the same derivation that `bench record ticket` uses. Keep `chunkBase` as the one derivation.

## Acceptance

- [ ] For each refusal fixture of the `reviewrecord` package tests, the first collected finding's error equals the `CheckSource` error (RP22).
- [ ] The existing tests of `internal/reviewrecord` and `internal/reviewrecord/recordcmd` stay green with no edit to their expectations.
- [ ] A fixture with two stale entries in one chunk yields two findings from the collector.

## Verification

Run `bench test --package ./internal/reviewrecord/...` whole, not through `--run`. Run `bench test --package ./internal/gate` and `bench test --package ./internal/preflight/ownership`, because `CheckTrees` has callers there. Observe RP22 red before the production edit, and record the route of the red. Self-probe: make the collector return only its first finding, and report that the two-finding test goes red.
