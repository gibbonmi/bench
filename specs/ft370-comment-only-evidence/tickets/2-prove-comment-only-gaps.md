# Prove a comment-only Go gap between two trees

Blocked by: 1-move-tree-change-reader.md
Writes: internal/commentgap/ (new), internal/git/tree.go
Covers: CG9, CG10, CG11, CG12, CG13, CG14, CG15, CG16, CG17, CG18, CG19, CG20, CG21, CG22, CG23, CG24, CG25, CG26, CG27, CG28, CG29, CG30, CG31, CG32, CG42, CG43, CG46

## What to build

Create `internal/commentgap`, a new package that owns the comment-only
verdict. Its one export is `Prove(root, reviewed, later string) error`, and the
two operands are Git tree IDs. The spec's `The tree rule` and `The Go rule`
decisions state each rule and its order. The first rule that fails wins.

Obey these exact contracts:

- `Prove` returns nil for two equal tree IDs, and for two trees in which each change is a proven Go file.
- A refusal wraps exactly one rule sentinel and names the first path that the rule refused.
- The package exports the ten sentinels `ErrUnreadable`, `ErrEmptyChanges`, `ErrStatus`, `ErrMode`, `ErrNotGo`, `ErrScan`, `ErrTokens`, `ErrDirective`, `ErrCgo`, and `ErrExampleOutput`.
- `Prove` reads the configured change list through `git.TreeChanges` and keeps its empty-list refusal.
- A nonempty list takes a complete read through `git.TreeChangesIncludingSubmodules`, then the existing per-path rules. Each blob uses `git.ReadTreeFile`.
- Add `TreeChangesIncludingSubmodules` in `internal/git/tree.go`. It shares the existing reader and parser, with `--ignore-submodules=none` as its only policy change.
- The Go rule compares the token kind and the literal of each token, and not the positions.
- The Go rule keeps the automatic semicolon as a token, so a newline that a block comment gains can change the list.
- The directive check compares two ordered lists of directive-shaped comment texts, not two sets or two counts.

The spec states the three directive shapes and the example-output prefixes.
Use `go/build/constraint.IsPlusBuild` for the legacy build line, and
`go/parser` with `ImportsOnly` for the `"C"` import. A parse failure wraps
`ErrScan`. The package imports `internal/git` and the standard library only.

Ticket 3 calls `Prove` from the checkpoint and does not print the sentinel. So
the error text can name the path and the rule in any form.

Put the Go rule rows over byte pairs in `TestProveGoGap`, and the tree rows over
real tree pairs in `TestProveTreeGap`. Both tests go in
`internal/commentgap/commentgap_test.go`. Each refusal row asserts its sentinel
with `errors.Is`, so a row cannot pass through another rule. The tree fixture
helpers can go in a second test file of the package, so each file stays under
400 lines.

## Acceptance

- [ ] `TestProveGoGap` holds CG9 to CG21, CG42, and CG43, and each refusal row asserts its own sentinel.
- [ ] `TestProveTreeGap` holds CG22 to CG32, and each refusal row asserts its own sentinel.
- [ ] CG31 proves a gap in `a b*.go` and in a Go file whose name holds a tab.
- [ ] CG30 sets `diff.ignoreSubmodules=all` in the fixture repository and gets `ErrEmptyChanges`.
- [ ] A directive check that compares only the count of the directive texts makes the CG15, CG17, and CG18 rows fail.
- [ ] CG46 refuses a hidden gitlink beside a visible Go comment edit with `ErrMode`.
- [ ] `bench test --package ./internal/git` and `bench test --package ./internal/gate` pass.
- [ ] `go list -deps` shows that `internal/commentgap` imports no Bench package other than `internal/git` and its dependencies.
- [ ] Each file in `internal/commentgap` stays under 400 lines.
