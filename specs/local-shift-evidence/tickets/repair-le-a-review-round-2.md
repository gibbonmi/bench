# Repair the LE-A confirming round findings

Blocked by: repair-le-a-review-round-1.md
Writes: internal/otelrecord/writer.go, internal/otelrecord/retention_test.go, internal/otelrecord/reader_test.go, internal/otelrecord/encode_test.go, specs/local-shift-evidence/spec.md
Covers: LE106

## What to build

Chunk: LE-A.

The reviewer decided that the writer refuses a rotation when the highest sealed sequence is the largest one. That refusal keeps a wrapped sequence from renaming over a sealed segment. Add row LE106 and its test.

Make the `ReadSelected` test hold the per-segment line count: a malformed live line after a sealed segment reads as its own line 1. Restate the LE94 test comment with the row's current reason. Remove the review ID from the problem-address test comment.

## Acceptance

- [ ] A planted file at the largest sequence name stays unchanged, and an append past the limit returns an error and seals no segment.
- [ ] A malformed live line after a sealed segment reads as `line 1 malformed`.
