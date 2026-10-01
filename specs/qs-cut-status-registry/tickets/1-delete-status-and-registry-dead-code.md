# Delete the dead status, registry, and test-repository code

Blocked by: none
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main.go, internal/status/status.go, internal/status/route_test.go, internal/status/status_producible_test.go, internal/testrepo/working_tree.go, capture/restructure-backlog.md
Covers: none

## What to build

The quality survey of 2026-09-29, "Small certain cuts" table, named three items with no production caller. A sweep of the module confirmed the facts below:

- The command registry gives 30 of its 65 rows a process attachment class. Only `commandDispositions` reads the class, and only `TestCommandDispositionsAreComplete` calls `commandDispositions`. The test skips each row with no class, so it does not prove that every command has a disposition. `TestCommandRegistryAXIDispositionsAreComplete` holds the completeness guarantee, and it stays.
- `status.IsInvocable` has no production caller. Two tests call it.
- `testrepo.TwoHopRelativeSymlink` has no caller.

Delete the `processAttachment` type, its three constants, the `commandDisposition` type, the `Attachment` field, and `commandDispositions`. Remove the `Attachment` value from each registry row in `cmd/bench/main.go`. Delete `TestCommandDispositionsAreComplete`.

Delete `IsInvocable`. The survey also named `parseAction`, but it is not dead. The test helper `testSignal` uses it in 36 places to give each expected signal its typed action. So keep `parseAction`. Rename `TestIsInvocable` to `TestParseActionInvocable` and make it call `parseAction(text).invocable()`, so the grammar keeps its edge cases. Make the producible-signal test call `parseAction` in place of `IsInvocable`, so its check does not become weaker.

Delete `TwoHopRelativeSymlink`.

Remove `TestCommandDispositionsAreComplete` and the word "disposition" from the `command_registry_test.go` row in `capture/restructure-backlog.md`. Move `parseAction`, `match`, `matchOptionalPath`, and the step-separator guard from `internal/status/status.go` into `internal/status/route_test.go`, because only status tests call them. Add a `TestParseActionInvocable` case that only the step-separator guard rejects.

This change gives no change to CLI output, help text, exit codes, or status rows.

## Acceptance

- [ ] No Go source names `processAttachment`, `attachmentDirect`, `attachmentSystem`, `attachmentShip`, `commandDisposition`, or `commandDispositions`.
- [ ] No registry row in `cmd/bench/main.go` sets an `Attachment` field.
- [ ] No Go source names `IsInvocable` or `TwoHopRelativeSymlink`.
- [ ] `TestParseActionInvocable` holds each case that `TestIsInvocable` held.
- [ ] The producible-signal test requires each produced action to parse as invocable.
- [ ] `go vet ./...` and the changed-package tests pass.
- [ ] The `subcommand-routing`, `axi-query-registry`, and `package-shipped-surface` checks pass.
- [ ] The `command_registry_test.go` row in `capture/restructure-backlog.md` does not name `TestCommandDispositionsAreComplete` or a disposition assertion.
- [ ] No production Go file in `internal/status` defines `parseAction`, `match`, `matchOptionalPath`, or the step separator.
- [ ] `TestParseActionInvocable` fails when the step-separator guard is removed.
