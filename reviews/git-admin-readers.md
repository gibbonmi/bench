# Review pickup: git-admin-readers

Frozen base: `9d1a59e011f07b78ab1ab3afccff2596fbec69af`. Reviewed tip: `6f5fd12a959515e6c3f745e7d44c5b7f9ff35b9e`.
Round 1 ran three axes through `codex exec` on `gpt-6-astra` at medium effort on 2026-09-05.
The repair ticket `repair-review-round-1.md` carries the accepted findings.

## Standards

- Count: 8. Worst: `Worktrees` repeats the bounded query, trim, and validation that `CommonDir` owns.
- ask-user | internal/git/worktree_admin.go `Worktrees` | one source per fact | `Worktrees` repeats the common-directory resolution that `CommonDir` now owns. Not repaired in this round.
- ask-user | internal/git/worktree_admin.go `AdminPath` | one source per fact | The empty-answer refusal repeats the predicate inside `validateCommonDir`. Not repaired in this round.
- ask-user | internal/conformance/git_plumbing_owner_test.go `stringLiteralValue` | one source per fact | The helper repeats `stringLiteral` in subcommand_routing_test.go. Not repaired in this round.
- ask-user | internal/git/admin_readers_test.go | one source per fact | The hostile table re-types the stub's `missing-admin`, `symlink-admin`, and `file-admin` joins. Not repaired in this round.
- ask-user | internal/gittest/gittest.go `StubGitDir` comment | comment register | The comment inventories the stub modes the script below encodes. Not repaired in this round.
- ask-user | five new test files | comment register | Test doc comments cite coverage row ids and narrate the migration. Not repaired in this round.
- ask-user | internal/conformance/git_plumbing_owner_test.go, checks_test.go | comment register | Two comments carry the red record the spec owns. Not repaired in this round.
- auto-fix | internal/git/worktree_admin.go `AdminPath` comment | comments describe current code | "Every caller lstats the answer" is false for `indexIdentity`. Repaired with the `AdminPath` change.

## Spec

- Count: 2. Worst: the check walks only `cmd/` and `internal/`, so a root-level Go file escapes it.
- auto-fix | internal/conformance/git_plumbing_owner_test.go `checkGitPlumbingOwner` | "parses each non-test Go file outside `internal/git/` and `internal/gittest/`" | The walker skips `consumer_payload.go` at the module root. Row GR27 gains a root-level case.
- auto-fix | GR11 | "holds the repository-level lock file" | The test asserts only a non-nil release function. Row GR11 gains a held-lock assertion.

## Coverage

- Count: 3. Worst: a root with a symlink followed by `..` makes `AdminPath` answer a path in the wrong directory.
- auto-fix | root `/base/jump/..` where `jump` is a symlink | `filepath.Abs` cleans the `..` lexically. The reader answers `/base/.git/index`, and Git means `/base/physical/.git/index`. | Row GR5 gains that root shape.
- auto-fix | a symlinked temporary parent | the independent absolute-format run canonicalizes the parent, and the reader keeps the alias. GR2 reds for the primary and bare cases. | Row GR2 gains a symlinked parent.
- no-op | a repository root whose path text holds `--git-dir` or `--git-path` | the three pass-through stub modes match the flag as a substring of `$*`. | The stub is a test harness, and no test root carries flag text. Reported, not repaired.
