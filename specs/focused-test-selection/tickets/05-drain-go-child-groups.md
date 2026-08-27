# Drain the Go child process groups

Blocked by: 04-repair-selection-and-environment.md
Writes: internal/testreport/selection.go, internal/testreport/command.go, internal/testreport/cancel_test.go

## What to build

Put the changed-mode `go list` child in a private process group. Drain its descendants after an interruption.

Keep the existing process-group behavior for the focused `go test` child. Add independent cancellation proofs for both Go hops.

## Acceptance

- [ ] N02 — cancellation of `go list` leaves no child or descendant alive.
- [ ] N02 — cancellation of `go test` leaves no child or descendant alive.
- [ ] N02 — an interrupted Go hop emits no partial package table.
