# Declare persisted fields and resolve them against the record decoders

Blocked by: none
Writes: internal/spec/spec.go, internal/spec/persisted_test.go (new), internal/intent/ledger/fields.go (new), internal/intent/ledger/fields_test.go (new), internal/reviewrecord/parse.go, internal/reviewrecord/fields_test.go (new), .bench/commitment.json
Covers: CB1, CB2, CB3, CB4, CB5, CB6, CB7, CB8, CB9, CB10, CB11, CB12

## What to build

Deliver CB-C1: a spec declares each persisted field on one `Persisted fields:` header line, and each record owner answers whether its own decoder knows a field path.

Add the header scan to `internal/spec` beside the scan of `Status:` and `Roadmap:`. Put the scan in `spec.go`, which has 140 lines of headroom. The scan returns the parsed `<record>:<path>` items, or an empty list for a spec with no line. It ignores a line inside a code fence. The parser refuses an empty item, a record name other than `intent-ledger` or `review-record`, an empty path segment, and a duplicate item. Each refusal is an error whose text starts with `persisted-fields-invalid`. This ticket adds no refusal-route face.

Add one exported predicate to `internal/intent/ledger` and one to `internal/reviewrecord`. Each predicate takes a dot-separated path and walks the `json` tags of its root type by reflection. It descends through pointers and slice elements, and it refuses a path through a non-struct. The struct tags are the one source of field knowledge. `internal/intent/ledger` is a leaf package under the purity census: its new code and test import only the standard library, and the test does not call `t.Parallel`. Put the review-record predicate in `parse.go`, which has 82 lines of headroom.

The `.bench/commitment.json` path is on this line for the spec-stage close, which commits the deliverable binding inside the range that the plan-only preflight grades. The author makes no other change to it.

Later tickets consume these three exported values: the header scan, the ledger predicate, and the review-record predicate.

## Acceptance

- [ ] `TestPersistedFieldsHeader` returns `intent-ledger:commitment.gaps` and `review-record:folds` from one line, returns an empty list for a spec without the line, and ignores a fenced line (CB1, CB2, CB3).
- [ ] The parser refuses `intent-ledger:a,,`, `gate-evidence:x`, `intent-ledger:commitment..gaps`, and a duplicate item, and each error starts with `persisted-fields-invalid` (CB4, CB5, CB6, CB7).
- [ ] `TestLedgerKnownField` resolves `commitment.gaps` and `commitment.gaps.chunk`, and refuses `commitment.unknown` and `schema.extra` (CB8, CB9, CB10, CB11).
- [ ] `TestReviewRecordKnownField` resolves `folds` and refuses `foldz` (CB12).
- [ ] The purity census of `internal/intent/ledger` stays green with the new files.

## Checkpoint verification

Run CB-C1's verification commands from the spec's completion plan. Probe each predicate with one omission mutant (for example, a walker that stops at a slice) and record the red.
