# Name the clean route for an unlanded release and accept a cleanup-pending clean target

Blocked by: none
Writes: internal/worktree/list.go, internal/worktree/ownership.go, internal/worktree/path.go, internal/worktree/clean_set.go, internal/worktree/identity_component.go, internal/worktree/unlanded_route_test.go, specs/ft345-release-route/tickets/ft345-release-route.md
Covers: none

## What to build

A release refuses an assignment whose branch did not land, and the release leaves
that row in the `cleanup-pending` state. Each retirement surface must then name a
route that can succeed. The release refusal names
`bench worktree clean --discard-branch <path>` as its next route. The
`bench worktree list` help row for a `cleanup-pending` row with a proven unlanded
branch names the same route. One route type renders both surfaces, so the two
surfaces cannot name different commands.

The `clean --target <label>` form accepts a `cleanup-pending` row, as the path form
does. A release refusal for a landed branch keeps its current route. A row whose
landed state is unknown keeps the release route.

## Acceptance

- [ ] A release of an unlanded assignment exits 1, and its `next=` line names
  `bench worktree clean --discard-branch <path>`. That clean plans a removal, and
  its apply removes the tree.
- [ ] After that refusal, the `bench worktree list` help row for the
  `cleanup-pending` row names `bench worktree clean --discard-branch <path>`, and
  that route removes the tree.
- [ ] `bench worktree clean --discard-branch --target <label>` on that
  `cleanup-pending` row plans a removal, and its apply removes the tree.
