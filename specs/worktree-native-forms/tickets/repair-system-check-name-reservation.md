# Repair: reserve the system check name

Blocked by: repair-worktree-single-resolution-and-replay.md
Writes: internal/testreport/check_test.go

## What to build

Review finding C2 in `reviews/worktree-native-forms.md`. The `--check` parser
routes `gate.SystemPhaseName` past the registry lookup, so a registry check of
the same name would run the wrong suite with no refusal. One test asserts that
the conformance registry names no check equal to `gate.SystemPhaseName`, so a
later registration turns the gate red.

## Acceptance

- [ ] WF46: `registry.Find(gate.SystemPhaseName)` reports not found, proved red by a throwaway registry entry named `system` that the test author removes.
- [ ] The gate `test` phase stays green for the whole `internal/testreport` package.
