# Repair the review round 1 findings

Blocked by: add-the-git-plumbing-owner-check.md
Writes: internal/git/worktree_admin.go, internal/git/admin_readers_test.go, internal/worktree/admin_readers_test.go, internal/conformance/git_plumbing_owner_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: GR2, GR5, GR11, GR27

## What to build

Read `reviews/git-admin-readers.md` for the accepted findings. Read `AdminPath` in
internal/git/worktree_admin.go, `TestAdminPathJoinsRelativeAnswerOntoRoot` and
`TestAdminPathMatchesIndependentRevParse` in internal/git/admin_readers_test.go,
`TestLockCleanupRegistrationFallsBackForNonWorktree` in
internal/worktree/admin_readers_test.go, and `checkGitPlumbingOwner` with
`TestGitPlumbingOwnerRedsARetypedFlag` in
internal/conformance/git_plumbing_owner_test.go.

Change `AdminPath` to join a relative answer onto the symlink-resolved absolute root.
A root that reaches the repository through a symlink followed by `..` then answers
the directory Git meant. Keep the final component unresolved, so GR40 holds. Rewrite the
doc comment so that it states the constraint without a claim about every caller.

Extend GR5's test with a root shaped `<base>/jump/..`, where `jump` is a symlink to
`<base>/physical/child`. Extend GR2's test with a primary checkout reached through a
symlinked parent directory. Extend GR11's test so that a second non-blocking lock on
the planted lock file fails while the release is pending and succeeds after it runs.

Change the walker of `checkGitPlumbingOwner` so that it also grades the non-test Go
files at the module root. Keep `cmd/` and `internal/` and skip every other
directory. Add one root-level case to GR27's test. The pickup file leaves in the
same commit that closes these findings. The five command registries in the
`Writes:` line are closure headroom for the worktree binding, and this ticket edits
none of them.

## Acceptance

- [ ] `AdminPath` answers the physical directory over a root with a symlink followed by `..`.
- [ ] GR2 passes over a primary checkout reached through a symlinked parent.
- [ ] GR11 proves the repository-level lock is held until the release runs.
- [ ] The check reds a retyped flag in a Go file at the module root.
- [ ] `reviews/git-admin-readers.md` is deleted in the closing commit.
- [ ] Self-probe: restore the lexical `filepath.Abs` join, and report the GR5 symlink case red.
