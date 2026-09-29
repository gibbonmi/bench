# Remove two copies from the spec test helpers

Blocked by: none
Writes: internal/spec/history_selected_test.go, internal/spec/spec_test.go
Covers: none

## What to build

Two spec test helpers hold a second copy of a fact that another source already owns. Remove each copy:

- `TestSelectedHistoryTrueBytes` has a `ParseFloat` guard. The 31-byte row check already enforces the TOON quoting rule, because a quoted hash adds two bytes to the row. The guard is an incomplete copy of that rule. Remove the guard.
- `writeSpec` is an alias of `writeFolderSpec`. Point each caller at `writeFolderSpec`, and remove the alias.

## Acceptance

- [ ] `TestSelectedHistoryTrueBytes` contains no `strconv.ParseFloat` call, and the test passes.
- [ ] No `writeSpec(` call or definition remains in `internal/spec`, and `bench test --package internal/spec` is green.
- [ ] The commit passes its lane.
