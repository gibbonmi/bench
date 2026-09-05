# Add the named Git administration readers

Blocked by: none
Writes: internal/git/worktree_admin.go, internal/git/checkout.go, internal/git/admin_readers_test.go (new), internal/git/worktree_admin_hostile_test.go, internal/git/worktree_admin_enum_test.go, internal/gittest/gittest.go
Covers: GR1, GR2, GR3, GR4, GR5, GR6, GR7, GR8, GR9, GR36, GR38, GR39, GR40

## What to build

Verify the premise first. Read `CommonDir`, `commonDirArgs`, `validateCommonDir`,
`boundedGit`, and `Worktrees` in internal/git/worktree_admin.go. Read
`ResolutionError`, `WorktreeFailure`, and `investigateGitFailureAction` in the
same file. Read `IsPrimaryCheckout` in internal/git/checkout.go. Read
`worktreeListTimeout` and `SetWorktreeListTimeoutForTest` in internal/git/git.go.
Read `StubGit` and `StubGitDir` in internal/gittest/gittest.go.

Read `TestWorktreesRejectsBadCommonDirBeforePorcelain` and
`TestWorktreesTypesStartFailures` in internal/git/worktree_admin_hostile_test.go.
Read `TestCommonDirReturnsUnvalidatedOutput` and
`TestCommonDirKeepsPlainOutputFailure` in internal/git/worktree_admin_enum_test.go.
Read `TestAppendWorktreeSurfacesClassifyFailure` in
internal/status/status_render_test.go.

Export two readers under the exact names `AdminDir(root string) (string, error)`
and `AdminPath(root, name string) (string, error)`. Every sibling ticket calls
these two names, so do not rename them. `AdminDir` runs the arguments `rev-parse
--path-format=absolute --git-dir` through `boundedGit`. `AdminPath` runs
`rev-parse --git-path <name>` in Git's default path format the same way. It
joins a relative answer onto the absolute root. The default format never
resolves a symlink. Route
`CommonDir` through `boundedGit` too.

Validate the answer of `AdminDir` and the answer of `CommonDir` with
`validateCommonDir`. Refuse an empty answer, a missing path, a symlink, and a
non-directory. Refuse an empty answer in `AdminPath`. Run
no existence check in `AdminPath`, because absence is the caller's fact. Add a
`Subject` field to `ResolutionError` with the three nouns `git common directory`,
`checkout administration directory`, and `checkout administration path`. Open the
error text with the subject noun, and keep the action `investigate the git
failure`.

Change `IsPrimaryCheckout` to call `AdminDir` and `CommonDir`. Keep no
`--absolute-git-dir` spelling in the adapter. Add the stub modes `bad-git-dir`,
`empty-git-dir`, `symlink-git-dir`, `file-git-dir`, `block-git-dir`,
`block-git-path`, `fail-git-dir`, `fail-git-path`, and `relative-git-path` to `StubGitDir`.
Answer the directory query and the file query in the stub.

In `fail-git-dir` mode exit nonzero on the directory query, and in `fail-git-path`
mode exit nonzero on the file query. In `relative-git-path` mode answer the
file query with `.git/<name>`. In those three modes pass every other invocation to
the real `git`, which the stub locates before the test replaces `PATH`. Five
sibling tickets need that pass-through, so prove it with one adapter test.

Write the new adapter tests in internal/git/admin_readers_test.go. Replace
`TestCommonDirReturnsUnvalidatedOutput` and `TestCommonDirKeepsPlainOutputFailure`
with the validated posture. Reuse the package's one repository constructor,
because the ordinary-build census allows one.

## Acceptance

- [ ] The directory reader answers the independent `rev-parse` path over a primary checkout, a linked worktree, and a bare repository.
- [ ] The file reader for `index` answers the independent `rev-parse` path over the same three repositories.
- [ ] The directory reader and `CommonDir` each refuse a missing path, an empty answer, and a symlink.
- [ ] The directory reader refuses a regular-file answer with the text `non-directory`.
- [ ] The file reader joins a relative answer onto the absolute root.
- [ ] The file reader returns a symlinked lease's own path, and `Lstat` reports a symlink.
- [ ] The directory reader returns a `timed out` typed failure under a 50-millisecond test bound.
- [ ] The file reader and `CommonDir` each return a `timed out` typed failure under the same bound.
- [ ] With an empty `PATH`, the three readers each return a typed failure that holds `rev-parse` and `executable file not found`.
- [ ] The error text opens with the subject noun for each of the three subjects.
- [ ] `CommonDir` returns a typed error that holds `rev-parse` after a nonzero exit.
- [ ] Over a bare repository the directory reader and `CommonDir` answer one path.
- [ ] `IsPrimaryCheckout` answers true over a bare repository.
- [ ] `TestAppendWorktreeSurfacesClassifyFailure` passes with the words `git common directory`.
- [ ] In `relative-git-path` mode the stub answers a real `rev-parse --show-toplevel` query.
- [ ] Self-probe: drop the non-directory check from the validator, and report the hostile reader test red.
