# Declare persisted fields on one spec header line

Blocked by: none
Writes: internal/spec/spec.go, internal/spec/persisted_test.go (new), .bench/commitment.json
Covers: CB1, CB2, CB3, CB4, CB5, CB6, CB7, CB66, CB67

## What to build

Deliver CB-C1: a spec declares each persisted field that its build adds on one `Persisted fields:` header line, and one exported scan reads that line.

Add the header scan to `internal/spec` beside the scan of `Status:` and `Roadmap:`. Put the scan in `spec.go`, which has 140 lines of headroom. The scan takes the spec bytes and returns the parsed `<record>:<path>` items, or an empty list for a spec with no line. It ignores a line inside a code fence. The parser refuses an empty item, a record name other than `intent-ledger` or `review-record`, an empty path segment, a path segment that holds whitespace or a colon, a duplicate item, and a second unfenced `Persisted fields:` line. Each refusal is an error whose text starts with `persisted-fields-invalid`. This ticket adds no refusal-route face, and it adds no field predicate.

Tickets 2, 3, 5, and 6 consume this scan. Each consumer prints the `persisted-fields-invalid` error when the scan refuses.

The `.bench/commitment.json` path is on this line for the spec-stage close, which commits the deliverable binding inside the range that the plan-only preflight grades. The author makes no other change to it.

## Acceptance

- [ ] `TestPersistedFieldsHeader` returns `intent-ledger:commitment.gaps` and `review-record:folds` from one line, returns an empty list for a spec without the line, and ignores a fenced line (CB1, CB2, CB3).
- [ ] The parser refuses `intent-ledger:a,,`, `gate-evidence:x`, `intent-ledger:commitment..gaps`, and a duplicate item, and each error starts with `persisted-fields-invalid` (CB4, CB5, CB6, CB7).
- [ ] The parser refuses `intent-ledger:commitment. gaps`, `intent-ledger:a:b`, and a spec with two unfenced `Persisted fields:` lines, and each error starts with `persisted-fields-invalid` (CB66, CB67).

## Checkpoint verification

Run CB-C1's verification commands from the spec's completion plan. Probe the fence state with one omission mutant and record the red.
