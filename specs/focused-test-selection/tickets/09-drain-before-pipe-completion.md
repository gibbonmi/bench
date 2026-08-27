# Drain before pipe completion

Blocked by: 08-close-cancellation-races.md
Writes: internal/testreport/command.go, internal/testreport/selection.go, internal/testreport/cancel_test.go, reviews/focused-test-selection.md (delete)

## What to build

Drain a Go process group before a resistant descendant can block pipe or process completion. Keep one interruption posture for both Go hops.

Make the parking descendant inherit stdout. Delete the review pickup after the retained-pipe predicate is green.

## Acceptance

- [x] N02 — a SIGINT-resistant descendant can retain stdout without blocking cancellation.
- [x] N02 — both Go hops drain before they await final pipe or process completion.
- [x] N02 — interrupted Go hops return no partial package table.

Delivered outcome: the central canceller drains the full process group before
either Go hop waits for final pipe or process completion.
