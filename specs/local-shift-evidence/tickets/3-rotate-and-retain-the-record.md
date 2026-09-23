# 3. Rotate and retain the record

Blocked by: none
Writes: internal/otelrecord/writer.go, internal/otelrecord/reader.go, internal/otelrecord/reader_test.go, internal/otelrecord/retention_test.go (new), internal/bounds/bounds.go, tests/canary/package-core-guard/bounds-duplicate-owner
Covers: LE13, LE14, LE15, LE16, LE17, LE18, LE19, LE20, LE21

## What to build

Chunk: LE-A.

Add `bounds.RecordSegmentLimit` and `bounds.RecordSegmentsRetained` to the bounds registry. Propose 16 MiB and 8, as the spec's reviewer-owned sizes.

When an append would take the live segment `traces.jsonl` past the limit, the writer tries a non-blocking exclusive lock on one lock file. With the lock, it checks the size again, renames the live segment to `traces-<UTC stamp>-<pid>.jsonl`, and removes the oldest sealed segments until the retained count remains. It then appends to a new live segment. A writer that cannot take the lock appends to the live segment and rotates on a later append. Each append stays one synchronous `O_APPEND` write.

Give the tests an unexported constructor that takes a small limit and a small retained count.

Make `ReadSpans` and `ReadSelected` read the sealed segments in name order and then the live segment. Each segment passes the existing symlink and regular-file grade before the open. `NewestLanding` reads through `ReadSpans`, so it sees every segment.

## Acceptance

- [ ] With a 1 KiB limit, an append past the limit leaves one sealed `traces-<stamp>-<pid>.jsonl` segment and a new live segment that holds the line.
- [ ] After the retained count plus two rotations, exactly the retained count of sealed segments remains, and the two oldest names are gone.
- [ ] `ReadSpans` returns the spans of two sealed segments in name order and then the spans of the live segment.
- [ ] `NewestLanding` returns the stages of a landing whose phase spans sit in a sealed segment.
- [ ] `ReadSelected` returns a selected span from a sealed segment.
- [ ] Two writers with separate handles that append across forced rotations leave every line whole, and the line total equals the append total.
- [ ] While another open file holds the rotation lock, an append past the limit writes to the live segment and creates no sealed segment.
- [ ] A FIFO at a sealed segment name makes `ReadSpans` return an error within the test deadline, and a symlink there makes it return an error.
