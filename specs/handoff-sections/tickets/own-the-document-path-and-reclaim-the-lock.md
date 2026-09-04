# Own the document path and reclaim the lock

Blocked by: remove-the-section-at-retirement.md, resolve-the-callers-section-in-bench-handoff.md
Writes: internal/handoffdoc, internal/status/handoff.go, internal/worktree/lifecycle.go, internal/worktree/worktree_test.go, internal/handoff/text.go, internal/preflight/gather.go, internal/preflight/gather_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: HS4

## What to build

Verify the premise first: `capture/session-handoff.md` is spelled in
`internal/status/handoff.go` and again in `internal/worktree/lifecycle.go`,
and the leaf's `acquire` leaves `capture/session-handoff.md.lock` on disk
after every write. The residue is an undeclared ignored path, and a landing
refuses it. Then give the leaf one exported `DocumentPath` constant, and make
both spellings read it. Reclaim the lock on release with the safe-unlink
shape: after the flock succeeds, compare the open file's inode with the
path's inode, and reopen when they differ; on release, unlink the lock file
before the flock drops. A writer that dies mid-hold leaves the lock file,
and the next writer reclaims it.

Retire the dead `canonicalRoot` in `internal/preflight/gather.go` and its
test, and correct the package comment in `internal/handoff/text.go`, which
still names the old parser.

## Acceptance

- [ ] After `Update` returns, no `.lock` file sits beside the document.
- [ ] The two-writer parallel test still leaves every section present.
- [ ] One spelling of the document path remains outside the leaf, in a test that counts it.
- [ ] Self-probe: drop the inode comparison after the flock, and report the parallel test result.
