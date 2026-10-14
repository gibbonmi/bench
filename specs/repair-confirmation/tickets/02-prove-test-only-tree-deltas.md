# Prove test-only tree deltas

Blocked by: none
Writes: internal/testgap/testgap.go (new), internal/testgap/testgap_test.go (new), internal/gittest/gittest.go, internal/gittest/trees.go (new), internal/gittest/gittest_test.go, internal/commentgap/trees_test.go, internal/commentgap/commentgap_test.go
Covers: RC13, RC14, RC15, RC16, RC17, RC18, RC19, RC52, RC53

## What to build

A new package `internal/testgap` proves that a tree delta changes only Go test files.
Its `Prove(root, reviewed, later)` returns nil for equal trees and for a delta in which every change obeys all of these conditions:

- The status is `A` or `M`.
- The destination mode is a regular file mode, and an `M` change keeps its source mode.
- `outline.IsGoTestFile` accepts the path.
- No path element is `testdata`.

It reads the delta through `TreeChangesIncludingSubmodules`, which runs with `--no-renames`.
Each refusal returns one named sentinel: `ErrUnreadable`, `ErrEmptyChanges`, `ErrStatus`, `ErrMode`, `ErrNotTest`, or `ErrTestdata`.
The error names the path that failed.
Mirror the shape of `commentgap.Prove`.

Move the tree fixture builder out of the comment-only prover tests into `internal/gittest` as an exported helper.
The comment-only prover tests and the test-only prover tests both call that helper, and the old unexported helpers go away.
The comment-only prover tests keep every case and every assertion.
Before you lock the import edges, run `go list -deps` and confirm that `internal/testgap` imports no `internal/reviewrecord` package.

## Acceptance

- [ ] RC13: the prover refuses a `_test.go` path with a `testdata` element with `ErrTestdata`.
- [ ] RC14: the prover refuses a deleted `_test.go` file with `ErrStatus`.
- [ ] RC15: the prover refuses a `_test.go` path whose destination mode is a symbolic link with `ErrMode`.
- [ ] RC16: the prover refuses a modified `_test.go` file whose mode changes with `ErrMode`.
- [ ] RC17: the prover refuses a non-Go path with `ErrNotTest`.
- [ ] RC18: the prover refuses an empty change list between different trees with `ErrEmptyChanges`.
- [ ] RC19: the prover accepts a test path that contains a space.
- [ ] RC52: the prover refuses a `_test.go` gitlink with mode `160000` with `ErrMode`.
- [ ] RC53: the prover refuses a renamed `_test.go` file with `ErrStatus`.
- [ ] The comment-only prover tests pass unchanged in behavior through the shared tree builder.
