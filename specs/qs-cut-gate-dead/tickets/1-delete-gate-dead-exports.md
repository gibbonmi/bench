# Delete the dead gate entry points and snapshot helpers

Blocked by: none
Writes: internal/gate/gate.go, internal/gate/tree_snapshot.go, internal/gate/prospective.go, capture/restructure-backlog.md
Covers: none

## What to build

The quality survey of 2026-09-29, "Small certain cuts" table, found six functions in `internal/gate` that have a definition and no caller. A search of the module confirms the claim for each function. No production code, no test, no census, no `.bench/` script, no ADR, and no project profile names them:

- `Run`, `RunContext`, `RunAndRecord`, and `ExecuteReusingFreshGreen` in `internal/gate/gate.go`.
- `captureWorkingTree` in `internal/gate/tree_snapshot.go`.
- `buildProspectiveSubjectFor` in `internal/gate/prospective.go`.

The live entry points stay: `Command`, `RunCommand`, `RunAndRecordContext`, and `Execute`. Each helper that a deleted function calls has a different live caller, so each helper stays.

Delete the six functions. Do not change a live function. In `capture/restructure-backlog.md`, the `gate.go` row starts its moved range at `Run`. Start that range at `RunAndRecordContext`, the first entry point that stays.

## Acceptance

- [ ] No Go source defines or calls `gate.Run`, `RunContext`, `RunAndRecord`, `ExecuteReusingFreshGreen`, `captureWorkingTree`, or `buildProspectiveSubjectFor`.
- [ ] `RunAndRecordContext`, `RunCommand`, `Command`, and `Execute` are unchanged.
- [ ] The restructure backlog names no deleted function.
- [ ] `go vet ./...` and the changed-package tests pass.
