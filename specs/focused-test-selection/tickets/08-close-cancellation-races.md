# Close the cancellation races

Blocked by: 07-close-graph-environment-and-edge-proofs.md
Writes: internal/testreport/command.go, internal/testreport/selection.go, internal/testreport/cancel_test.go, reviews/focused-test-selection.md (delete)

## What to build

Observe cancellation while the decoder and the process complete in either order. Drain both Go process groups after a graceful signal.

Use one parking fixture source. Make its descendant survive the graceful signal so removal of the final drain turns both tests red.

Delete the review pickup after every repair predicate is green.

## Acceptance

- [x] N02 — `go test` cancellation cannot block after decode completes.
- [x] N02 — the `go list` oracle requires the final group drain.
- [x] N02 — the `go test` oracle requires the final group drain.
- [x] N02 — interrupted Go hops return no partial package table.
- [x] N02 — the parking fixture has one source.

Delivered outcome: cancellation handles either decode and process completion
order. Both Go-hop oracles require the final group drain.
