# Repair the review findings

Blocked by: grade-the-publication-boundary.md
Writes: internal/worktree/merge.go, internal/worktree/merge_test.go, internal/worktree/land_identity.go, internal/worktree/ownership.go, internal/gate/authorization/authorization.go, internal/gate/authorization/lane_test.go, specs/worktree-merge/spec.md, reviews/worktree-merge.md (delete)

## What to build

The six code repairs `reviews/worktree-merge.md` names, S1, S2, P2, C1, C2,
and C3, each red-first. The status-clean predicate and the checked-out-branch
predicate each collapse to one source in the worktree package. A failed
ancestry query or an unresolved default branch in the `--from` commit lookup
refuses and names the failure. An ambiguous `--from` prefix refuses naming
every matching assignment id.

The lane's prose subject reads `git diff --name-only -z` and splits on NUL,
so a non-ASCII Markdown path reaches the prose check as itself. The sibling
lookup runs over active assignments only. The same commit deletes
`reviews/worktree-merge.md`.

## Acceptance

- [ ] WM35: a `--from` prefix that matches two assignments refuses naming
      both ids, with and without a branch of the same spelling.
- [ ] WM36: an incoming Markdown path with a non-ASCII byte reaches the lane's
      prose check as its own bytes, not C-quoted.
- [ ] WM37: a `--from` spelling that names only a non-active assignment
      resolves through the commit lookup, or refuses naming the value when no
      commit matches.
- [ ] WM38: a failed ancestry query in the `--from` lookup refuses naming the
      query, not the history.
- [ ] WM39: a `--from` spelling the target grammar rejects as a path, such as
      `refs/heads/main`, resolves through the commit lookup; only the
      ambiguity outcome refuses.
- [ ] `landingDestination` and `mergeTargetTip` refuse a dirty checkout
      through one predicate. `landingSource`, `validateCreationBundle`, and
      `mergeOnAssignmentBranch` read the checked-out branch through one
      predicate. The existing landing and merge tests stay green.
