# Quote selected-history command operands

Blocked by: 01-quote-tree-build-arguments.md
Writes: internal/spec/history_selected.go, internal/spec/history_selected_test.go, internal/spec/widen.go, internal/spec/widen_test.go, internal/spec/widen_rerun_test.go (new)
Covers: S13, S40

## What to build

Migrate historyDetail from axi.ShellQuote to sanitize.ShellQuote.
Migrate the spec widen rerun line in widen.go the same way, and keep its `--decision <reference>` placeholder.
Keep LineSafe before rendering and retain selected history, operand order, and response envelope.

Review chunk: S-B1.
The predecessor supplies the accepted contract named above.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

The real printed detail command executes through a POSIX shell and reads the same selected history.
Safe and hostile accepted operands use the approved always-quoted spelling.

- [ ] Selected-history detail operands use the surviving quote owner. (S13).
- [ ] The spec widen rerun line quotes the slug, the ticket, and each entry through the surviving owner. (S40).

## Checkpoint verification

- `bench test --package ./internal/spec`

## Required omission evidence

Bypass the owner with former bare-token output at historyDetail and at the widen rerun line.
The independent exact fixture must fail even if argv still round-trips.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.
