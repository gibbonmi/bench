# 3. Rotate and retain the record

Blocked by: none
Writes: internal/otelrecord/writer.go, internal/otelrecord/reader.go, internal/otelrecord/reader_test.go, internal/otelrecord/retention_test.go (new), internal/bounds/bounds.go, tests/canary/package-core-guard/bounds-duplicate-owner
Covers: LE13, LE92, LE94, LE14, LE93, LE15, LE16, LE17, LE18, LE19, LE20, LE21

## What to build

Chunk: LE-A.

Add `bounds.RecordSegmentLimit` and `bounds.RecordSegmentsRetained` to the bounds registry. Propose 16 MiB and 8, as the spec's reviewer-owned sizes.

Add one helper in the record package that formats and parses the sealed name `traces-<sequence>.jsonl`. The sequence is a zero-padded 20-digit decimal, so name order equals sequence order. The writer, the prune, and both readers use this one helper.

When an append would take the live segment `traces.jsonl` past the limit, the writer tries a non-blocking exclusive lock on one lock file. With the lock, it checks the size again and reads the sealed names. The next sequence is one more than the highest sequence present. If a file already holds that name, the writer takes the next free sequence. It then renames the live segment to that name and removes the lowest sequences until the retained count remains.

After the rotation the writer appends to a new live segment. A writer that cannot take the lock appends to the live segment and rotates on a later append. Each append stays one synchronous `O_APPEND` write. No clock enters the sealed name.

Give the tests an unexported constructor that takes a small limit and a small retained count.

Make `ReadSpans` and `ReadSelected` read the sealed segments in sequence order and then the live segment. Each segment passes the writer's grade, `gradeRecordPath`, before the open. `ReadSpans` has no grade today, so this ticket adds it. `NewestLanding` reads through `ReadSpans`, so it sees every segment.

## Acceptance

- [ ] With a 1 KiB limit, an append past the limit seals the live segment under the next sequence and starts a new live segment.
- [ ] Two rotations in a row leave two sealed segments with consecutive sequences, and the first keeps every line it sealed.
- [ ] A planted file at the next sequence name stays unchanged, and the rotation seals under the following free sequence.
- [ ] After the retained count plus two rotations, exactly the retained count of sealed segments remains, and the two lowest sequences are gone.
- [ ] The prune removes the lowest sequences even when their modification times are the newest in the directory.
- [ ] `ReadSpans` returns the spans of two sealed segments in sequence order and then the spans of the live segment.
- [ ] `NewestLanding` returns the stages of a landing whose phase spans sit in a sealed segment, and `ReadSelected` returns a selected span from a sealed segment.
- [ ] Two writers with separate handles that append across forced rotations leave every line whole, and the line total equals the append total.
- [ ] While another open file holds the rotation lock, an append past the limit writes to the live segment and creates no sealed segment.
- [ ] A FIFO at a sealed segment name makes `ReadSpans` return an error within the test deadline, and a symlink there makes it return an error.
