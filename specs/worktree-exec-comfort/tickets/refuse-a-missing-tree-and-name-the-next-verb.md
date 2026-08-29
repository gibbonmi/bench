# Refuse a missing tree and name the next verb

Blocked by: none
Writes: internal/worktree/path.go, internal/worktree/list.go, internal/worktree/identifier_operand_test.go, internal/worktree/list_actions_test.go, internal/worktree/testdata/

## What to build

The shared target refusal printer prints its detail line and then one line
`next=<verb>`. Every refusal before the target resolves names
`bench worktree list`, so `exec` and `path` describe one failure one way. No
refusal before resolution prints a `worktree:` line, because no path exists
yet.

The target resolver checks the tree after the active-state check and before the
creation bundle. A not-exist stat error refuses with the reason
`worktree tree is missing`. Any other stat error keeps today's path, so a
permission error still reads as the creation bundle's refusal.

One recovery producer names the route out of a missing tree. The producer
returns `bench worktree clean --landed` for a landed assignment branch, and the
assignment's release form otherwise. The missing-tree refusal's `next=` line
carries that producer's output.

`list` builds an assignment row's actions from the same tree fact. A row whose
tree cell is `missing` gets one action from the recovery producer, and it gets
no `path` action and no `exec` action. A row whose tree cell is `present` keeps
its `path` and `exec` actions unchanged.

The recovery producer is one function that both `list` and the refusal printer
call. A second copy of the landed rule in `list` drifts from the refusal.

## Acceptance

- [ ] F5: `bench worktree exec no-such-label -- true` prints
      `bench worktree exec: target is unassigned` and then
      `next=bench worktree list` on stderr at exit 1.
- [ ] F6: `bench worktree path no-such-label` prints the same two lines with
      the `bench worktree path` prefix.
- [ ] F7: an active assignment whose worktree directory is removed refuses exec
      with `bench worktree exec: worktree tree is missing` at exit 1.
- [ ] F8: that refusal's second line is `next=bench worktree clean --landed`
      when the assignment branch has landed.
- [ ] F9: that refusal's second line names the assignment's release form when
      the branch has not landed.
- [ ] F10: `actionsForRows` yields no `path` action and no `exec` action for an
      active assignment row whose tree cell is `missing`.
- [ ] F11: the same row yields exactly one action whose text equals the
      recovery producer's output.
- [ ] F12: an active row whose tree cell is `present` keeps its `path` and
      `exec` actions unchanged.
