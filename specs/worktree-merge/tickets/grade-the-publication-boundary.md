# Grade the publication boundary

Blocked by: refuse-unsafe-merge-inputs.md
Writes: internal/worktree/merge.go, internal/worktree/merge_test.go, internal/landing/merge.go, internal/gate/authorization/authorization.go

## What to build

The lane and the boundary after it. A failing lane check refuses at exit 1
with the check named and nothing changed. A fast-forward runs the lane on the
incoming tree, so an ungraded commit never reaches the branch.

The lane's
prose subject is the exact set of Markdown paths that differ between the
previous tip's tree and the composed tree. The lane authority gains one input, the
previous tip. It derives that set by a tree diff against the graded tree, so
the verb composes once. A
checkout edited between the pre-check and the ref update refuses before the
ref moves. The reconcile is `git reset --merge <published>` at the target alone, with no
residue bracket. A reconcile failure after the ref update exits 3 with
`merged{...,next=git -C <path> reset --merge <published-commit>}`; a path that
is not line-safe takes the `bench worktree exec <id> --` pointer form.

## Acceptance

- [ ] WM19: a failing lane check prints `lane{outcome=fail,check=<name>}` at
      exit 1, and the tip and the checkout stay unchanged.
- [ ] WM21: a reconcile failure after the ref update exits 3 with
      `merged{...,next=git -C <path> reset --merge <published-commit>}` naming
      the published commit.
- [ ] WM22: a checkout edited between the pre-check and the ref update refuses
      before the ref moves, and the edit survives.
- [ ] WM23: a fast-forward from a commit whose tree fails the lane refuses,
      and the tip stays.
- [ ] WM24: the lane's prose placeholder resolves to exactly the Markdown
      paths that differ between the previous tip's tree and the composed tree.
