# Add the create --from sibling start

Blocked by: none
Writes: internal/worktree/worktree.go, internal/worktree/worktree_test.go, internal/worktree/merge.go, internal/worktree/merge_test.go, internal/usage/worktree.go, cmd/bench/command_registry_test.go

## What to build

An agent runs `bench worktree create --from <target>` and gets a new assignment
that starts at the sibling's committed tip. The flag value resolves through the
sibling lookup that `mergeSiblingTip` owns, with no target to exclude. The
resolved tip becomes the `requestedStart` operand of `createAt`, so the ledger
`start` records that tip. Every `--from` refusal prints through
`printTargetRefusal` with the verb name `bench worktree create`. A `--from` with
`--refresh` is a usage refusal at exit 2, and it runs before `refreshop.Consume`.

This ticket crosses one shared surface. It extends the `usage.WorktreeCreate`
grammar in `internal/usage/worktree.go` with `[--from <target>]`, and that const
is a member of the `worktreeCommands` list the build ticket also edits. The full create
grammar joins the kept-grammar list in `cmd/bench/command_registry_test.go`, so
`bench worktree --help` pins the flag. The exact `bench help` fixture gains no
create row, because that inventory never listed `create`.

The build ticket, this ticket, and the gate-form ticket share
`internal/usage/worktree.go` and `cmd/bench/command_registry_test.go`. That is a
write conflict, not a capability dependency, so the coordinator serializes the
three in the order build, create, gate form. In that order this ticket is the
last one to write `internal/worktree`, so it carries that package's gate
invariant. The sibling lookup proves the creation bundle before it reads the tip,
by the precedent of `TestMergeRefusesAFailedTargetIdentityComponent`.

## Acceptance

- [ ] WF23: `create --from <sibling>` yields a worktree whose `HEAD` and ledger `start` equal the sibling tip, on the branch `intent.AssignmentBranchRef(owner, id)` names.
- [ ] WF24: `--from no-such-label` prints `--from names no active assignment` and `next=bench worktree list` at exit 1, and the ledger gains no record.
- [ ] WF25: a sibling with an uncommitted edit refuses with `sibling checkout is not clean` and `next=bench worktree exec <id> -- bench commit`.
- [ ] WF26: a sibling whose checkout is detached refuses with `sibling is not on its assignment branch`.
- [ ] WF27: `--refresh --from <sibling>` returns the usage line that names `--from` at exit 2, and `refreshop.Consume` runs nothing.
- [ ] WF28: two siblings that share a 9-character prefix make `--from <prefix>` refuse with `target is ambiguous: ` and both ids.
- [ ] WF29: `--from $'a\x01b'` refuses with `--from contains control characters` before a malformed ledger can refuse.
- [ ] WF30: `--from <landed-sibling>` refuses with `--from names no active assignment`.
- [ ] WF31: `bench worktree create --help` prints the grammar with `[--from <target>]` at exit 0.
- [ ] WF42: `bench worktree --help` names the full create grammar with `[--from <target>]`.
- [ ] WF43: a mutated owner marker, lock, or assignment state on the sibling refuses with that component's detail, and the ledger gains no record.
- [ ] The gate `test` phase stays green for the whole `internal/worktree` package.
