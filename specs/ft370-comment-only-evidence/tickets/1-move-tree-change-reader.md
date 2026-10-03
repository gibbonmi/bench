# Move the raw tree-change reader to internal/git

Blocked by: none
Writes: internal/git/tree.go, internal/gate/lane_select.go
Covers: CG40, CG41

## What to build

Move the one `git diff --raw --no-renames -z` reader from
`gate.ComposedChanges` to `internal/git`, so that the lane and the classifier
share one reader. The spec's `The shared raw tree-change reader` decision states
the move. Obey these exact contracts:

- `internal/git` exports `TreeChanges(root, from, to string) ([]TreeChange, error)` in `tree.go`.
- `TreeChange` keeps the four fields `Status`, `SrcMode`, `DstMode`, and `Path`, with the meanings that `ComposedChange` has today.
- `TreeChanges` gives its two operands to Git as they are. A tree ID that no commit names is a valid operand.
- The reader keeps the `-z` framing, so each path keeps its bytes, a space, a tab, and `*` included.
- The reader keeps Git's path order. A Git failure or an unreadable entry returns an error.

Make `gate.ComposedChange` an alias of `git.TreeChange`. Make
`gate.ComposedChanges` call `git.TreeChanges` with `base^{tree}` and the tree.
The gate wrapper keeps its `gate: composed change list unavailable` text.
Delete `parseComposedChange` and `composedChangeFields` from `internal/gate`.
No test matches the moved entry-parse error text.

Ticket 2 consumes `git.TreeChanges` with two tree IDs that
`git.TreeWithoutFile` wrote. `internal/git` gains no import, so it still
depends only on `internal/bounds` and `internal/canonicalpath`. Run
`go list -deps` on `internal/git` to confirm this edge set.

## Acceptance

- [ ] `TestComposedChangesExpandsANamedDirectory` gives the two changed files with their modes, without an edit to the test.
- [ ] `TestComposedChangesRepresentsARenameAsDeletionAndAddition` and `TestComposedChangesCarriesTheSymlinkMode` stay green without an edit.
- [ ] A reader that swaps the source mode and the destination mode makes `TestComposedChangesRepresentsARenameAsDeletionAndAddition` fail.
- [ ] `internal/gate` holds no raw-diff entry parser, and `gate.ComposedChanges` calls `git.TreeChanges`.
- [ ] `bench test --package ./internal/git` and `bench test --package ./internal/gate` pass.
- [ ] `internal/git/tree.go` stays under 400 lines, and `internal/gate/lane_select.go` does not grow.
