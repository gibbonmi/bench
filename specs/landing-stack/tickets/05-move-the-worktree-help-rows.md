# Move the worktree help rows beside the worktree leaves

Blocked by: none
Writes: cmd/bench/main.go, cmd/bench/worktree_leaves.go
Covers: none

## What to build

This ticket is a prefactor. It changes no output.

`cmd/bench/main.go` is over its 400-line limit, so the lane's growth ratchet reds any growth of that file. Tickets 07 and 11 each add one `bench help` row to the `worktree` family.

Move the `helpRow` list of the `worktree` command in `commandRegistry` into one function in `cmd/bench/worktree_leaves.go`, beside `worktreeLeaves` and `worktreeSuffix`. `commandRegistry` calls that function. Each row keeps its order, its suffix, its gap, and its description, so `bench help` and `bench worktree --help` print the same bytes.

After the move, a later help row of the `worktree` family is an edit of `cmd/bench/worktree_leaves.go` alone.

## Acceptance

- [ ] `TestHelpInventoryIsComplete` passes with no change of its expected text.
- [ ] `cmd/bench/main.go` is shorter than before the move, and no worktree `helpRow` stays in it.

## Checkpoint verification

- `bench test --package ./cmd/bench`
- `bench test --check root-conformance`
