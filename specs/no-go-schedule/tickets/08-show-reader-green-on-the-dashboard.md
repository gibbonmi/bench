# Show a reader-record tip as reader green on the ambient dashboard

Blocked by: 05-plan-the-reader-set-at-a-landing.md
Writes: internal/gate/reader_inspect.go (new), internal/gate/verdict.go, internal/status/status.go, internal/status/status_reader_test.go (new)
Covers: NG65, NG66, NG67

## What to build

`gate.Inspect` reports two more facts. The first says whether the current tree holds a reader record. The second is the count of reader-record commits after the green marker up to HEAD, from the history reader of ticket 05.

`GateVerdict` (`internal/status/status.go:437-458`) carries both facts. When the gate row would read stale and the current tree holds a reader record, `appendGateInfo` prints `reader green (<n> of 5)` with no command, at the stale row's rank. No stale row prints for that tip. A stale tip with no reader record keeps the stale row and its `bench gate` route (`internal/status/status.go:463-468`).

`internal/gate/verdict.go` holds 402 of its 457 granted lines. Put the record reads in `reader_inspect.go`.

## Acceptance

- [ ] A tip whose tree holds a reader record, with two reader-record commits after the marker, shows `reader green (2 of 5)` with no command.
- [ ] That tip prints no `stale (` row.
- [ ] A stale tip with no reader record keeps the stale row and its `bench gate` route.
