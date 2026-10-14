# Extract the record change functions

Blocked by: none
Writes: internal/reviewrecord/write.go, internal/reviewrecord/coverage.go, internal/reviewrecord/change.go (new), internal/reviewrecord/completion.go, internal/reviewrecord/export_test.go, internal/reviewrecord/replacement_test.go, internal/reviewrecord/fold_test.go, internal/reviewrecord/repair_test.go, internal/reviewrecord/recordcmd/amendment.go, internal/reviewrecord/recordcmd/chunk.go, internal/reviewrecord/recordcmd/verification.go, internal/reviewrecord/recordcmd/ordering_test.go
Covers: RT37, RT39

## What to build

The exported amendment, chunk, and verification writers each keep one change function over the record.
The chunk writer's create function, which owns the version 2 refusal for a new record, becomes one shared function too.
Each exported writer calls `write` with its own functions, and no refusal text or rule changes.

| site | old rule | new rule |
|---|---|---|
| amendment writer | The change body sits inside `RecordAmendment`. | One change function owns the body, and `RecordAmendment` calls it. |
| chunk writer | The change body sits inside `RecordChunk`. | One change function owns the body, and `RecordChunk` calls it. |
| verification writer | The change body sits inside `RecordVerification`. | One change function owns the body, and `RecordVerification` calls it. |
| chunk create | The create closure sits inside `RecordChunk`. | One create function owns the version 2 refusal, and `RecordChunk` calls it. |

`write`, `absent`, and `coveredAfter` are unexported in `internal/reviewrecord`, so the functions stay in that package.
`internal/reviewrecord/write.go` has 9 lines of headroom, so the extracted functions move to a new file.
The signature of `write` stays the same where possible, so `export_test.go` and `replacement_test.go` stay unchanged.

This ticket changes no behavior.
The ticket form of ticket 05 composes these functions.

## Acceptance

- [ ] Each of the four functions has one definition, and its exported writer calls it.
- [ ] A chunk entry with a stale plan digest still refuses through the chunk form, and `TestChunkNeedsCurrentAmendment` passes.
- [ ] The existing amendment, chunk, and verification tests pass with no change to their expectations.
