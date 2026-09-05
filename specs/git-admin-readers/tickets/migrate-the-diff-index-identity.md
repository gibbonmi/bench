# Migrate the diff index identity to the reader

Blocked by: add-the-named-git-admin-readers.md
Writes: internal/diff/snapshot.go, internal/diff/identity_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: GR26

## What to build

Verify the premise first. Read `AdminPath` in internal/git/worktree_admin.go, which
the blocker ticket adds. Read `indexIdentity` at internal/diff/snapshot.go line
271. Read `capturedIdentity`, which calls it. The site spells `rev-parse
--git-path index` today, and it joins a relative answer onto the root.

Change `indexIdentity` to call `git.AdminPath` with the name `index`. Return the
typed failure of the reader instead of the join. Keep the current answer for an
absent index file, because absence stays the caller's fact.

Add one new test to internal/diff/identity_test.go, beside the existing identity
tests. Name it `TestIndexIdentityRefusesUnresolvedAnswer`. The five command
registries in the `Writes:` line are closure headroom. They close the diff
binding, and this ticket edits none of them. Install the `fail-git-path` stub on
`PATH` with `gittest.StubGit`. The stub exits nonzero on the file query and
passes every other invocation to the real `git`.

## Acceptance

- [ ] `indexIdentity` calls the file reader and spells no Git flag.
- [ ] `indexIdentity` returns a typed resolution failure under the `fail-git-path` stub.
- [ ] The pre-existing `internal/diff` suite passes with its test logic unchanged.
- [ ] Self-probe: restore the join onto the root, and report the new test red.
