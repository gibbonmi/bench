# Repair: one target resolution, one --from guard, replay before --from

Blocked by: record-new-forms-in-prose.md
Writes: internal/worktree/build.go, internal/worktree/path.go, internal/worktree/merge.go, internal/worktree/worktree.go, internal/worktree/worktree_test.go, internal/worktree/build_test.go

## What to build

Review findings S1, S2, J2, and C3 in `reviews/worktree-native-forms.md`.
`resolveWorktree` returns the selected assignment record beside the path, or a
new `resolveAssignment` does and `resolveWorktree` wraps it, so
`resolveBuildTarget` reads the ledger once. The `path.go` printer comment names
the four verbs that share it. One helper owns the `--from` control-byte guard
and its sentence, and the merge verb and the create verb both call it.

`CreateCommand` runs the request replay lookup before it resolves `--from`.
So a replay with the same `--request` returns the existing record whatever
the sibling's state is. The lookup stays single-sourced with `createAt`; if the
order needs `createAt` to take a start resolver instead of a start value,
that is the shape.

## Acceptance

- [ ] S2: `resolveBuildTarget` calls `intent.Assignments` at most once per invocation, and every WF1 to WF11 row stays green.
- [ ] S1: `rg -- "--from contains control characters" internal/worktree/*.go` matches one production site, and WF29 stays green.
- [ ] J2: the `printTargetRefusal` comment names path, exec, show, and build.
- [ ] WF45: a `create --from <sibling>` replay with the same `--request` after the sibling takes an uncommitted edit exits 0 with the existing record. The ledger gains no second record.
- [ ] The gate `test` phase stays green for the whole `internal/worktree` package.
