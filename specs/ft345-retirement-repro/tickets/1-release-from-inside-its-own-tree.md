# A release that runs inside its own tree completes

Blocked by: none
Writes: internal/git/worktree_admin.go, internal/worktree/path.go, internal/worktree/worktree.go, internal/worktree/release_inside_test.go, internal/worktree/parallel_census_test.go, internal/treetarget/run.go, internal/harness/worktree.go, CHANGELOG.md
Covers: none

## What to build

FT345 has a release face. A `bench worktree release` that runs inside its own tree
removed the tree and then exited 1. The repro runs the release through
`bench worktree exec` inside a throwaway assignment. The release removes the tree,
then fails with `resolve git common directory` on the removed directory. The
assignment stays `cleanup-pending`, and a second release from the primary checkout
completes it.

The cause is the release root. The CLI roots a call at the tree that holds the
working directory. The release reads the intent ledger through that root after the
removal, to write its terminal receipt. The root is gone at that time.

The release work now runs from the primary checkout. One `git` function answers the
primary checkout's path from any checkout. The release span keeps the caller's root,
and the structure budget keeps the release file at its current length. The release, the `primary` tree target,
and the harness remove hook use that function. Before this change, the tree target
and the hook each read the first registration themselves.

## Acceptance

- [x] A release whose root is its own tree exits 0, removes the tree, and retires the assignment.
- [x] The release, the `primary` tree target, and the harness remove hook read the primary checkout through one function.
- [x] Removal of the re-root turns the new test red.

## Verification

`TestReleaseFromInsideItsOwnTreeCompletes` was red before the fix with the reported
`resolve git common directory` failure, and it is green after the fix. The
`bench test --changed` run passes; its skips are capability skips for sockets and
devices on this host. Two `bench probe` swaps bit the test, and each probe restored
its file. The first swap removed the re-root in `releaseRoot`. The second swap made
`git.PrimaryCheckout` return the caller's root. This worktree's build released a
throwaway assignment from inside its own tree with exit 0 and a `removed` row.
