# Compose a merge in the landing owner

Blocked by: none
Writes: internal/landing/merge.go (new), internal/landing/merge_test.go (new)

## What to build

The landing owner gains one merge operation beside the reviewed landing. It
takes the repository root, the branch ref, the expected previous tip, and the
incoming commit. It also takes the worktree path, the worktree fingerprint,
and the subject line. It resolves both commits exactly and decides the kind by
ancestry: `current` when the tip already contains the incoming commit,
`fast-forward` when the tip is an ancestor of it, and `merge` otherwise.

A `merge` composes the pair through the existing merge-tree composition with
the capture rule table. A conflict outside the table returns the landing's
conflict error with its paths. The owner authorizes the graded tree through
its worktree-commit policy; on a fast-forward the graded tree is the incoming
tree.

It rechecks the branch tip and the worktree fingerprint. It creates the
two-parent commit with the previous tip first, and it updates the branch ref
by compare-and-swap on the previous tip. A `current` result returns before
authorization and changes nothing. The result carries the kind, the previous
tip, the new tip, the tree, and the resolved capture paths.

The command ticket consumes this operation and passes the subject line; this
ticket owns the parent order, the kind decision, and the compare-and-swap.

## Acceptance

- [ ] WM4: the published commit's first parent is the previous tip and its
      second parent is the incoming commit.
- [ ] WM20: a branch tip moved after the lane makes the ref update refuse, and
      the moved tip survives.
- [ ] A fixture whose tip is an ancestor of the incoming commit returns
      `fast-forward` and moves the ref to the incoming commit with no new
      object.
- [ ] A fixture whose tip contains the incoming commit returns `current` and
      calls neither authorize nor update-ref.
