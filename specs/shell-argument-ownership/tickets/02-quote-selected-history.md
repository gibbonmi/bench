# Quote selected-history command operands

Blocked by: 01-quote-tree-build-arguments.md
Writes: internal/spec/history_selected.go, internal/spec/history_selected_test.go
Covers: S13

## What to build

Migrate historyDetail from axi.ShellQuote to sanitize.ShellQuote.
Keep LineSafe before rendering and retain selected history, operand order, and response envelope.

Review chunk: S-B1.
The predecessor supplies the accepted contract named above.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

The real printed detail command executes through a POSIX shell and reads the same selected history.
Safe and hostile accepted operands use the approved always-quoted spelling.

- [ ] Selected-history detail operands use the surviving quote owner. (S13).

## Checkpoint verification

- `bench test --package ./internal/spec`

## Required omission evidence

Bypass the owner with former bare-token output.
The independent exact fixture must fail even if argv still round-trips.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.
