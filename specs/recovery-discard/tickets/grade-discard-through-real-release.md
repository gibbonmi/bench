# Grade recovery plan and discard through a real release

Blocked by: add-recovery-discard.md
Ownership fence: `internal/contract/runtime/runtime_worktree_test.go`
Contracts: the recovery ref crosses `bench worktree release`'s production→`bench worktree recovery`'s plan and discard inside `internal/contract/runtime/runtime_worktree_test.go` and is asserted by GD1 against the real producing command rather than a hand-built ref
Assumptions: `add-recovery-discard.md` has landed so `--discard` exists and emits its own terminal action; the existing landedness fixtures stay as unit coverage of the comparison itself and are not replaced; the runtime fixture's real-`bench` preconditions are honoured through the package's existing helpers; every claim is re-derived from the tree at pickup

## What to build

The existing runtime recovery row builds its refs by hand — a cherry-pick, a
`commit-tree`, an `update-ref`, and a directly written assignment row — after
cleaning a foreign worktree. Every unit row can pass against a ref of that shape
while the real `bench worktree release` produces a ref shape the plan rejects.
That is the composition degenerate for this whole spec, and this session's
diagnosis hit it twice.

Compose the two commands that already exist separately in this package: drive a
real `bench worktree create`, dirty the checkout, run a real `bench worktree
release` to produce the preserved recovery ref, then plan that ref and discard it
through the real CLI. Assert the observable end state — the ref is gone, the
assignment row is closed, and the receipt names discard rather than retire.

The point is that the ref under grade was produced by production code. A row that
would stay green against a hand-built ref is not this row.

## Acceptance

- [ ] [GD1] a real `bench worktree release` over a dirty checkout produces a recovery ref, and `bench worktree recovery` plans that exact ref without error.
- [ ] [GD2] `--discard` with the fingerprint that plan just reported exits zero, and the recovery ref no longer resolves.
- [ ] [GD3] the assignment row the released worktree owned is closed after the discard.
- [ ] [GD4] the discard receipt names the discard action, distinguishable from the retire path's terminal action.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GD1 | make `release` skip writing the recovery ref for a dirty checkout | the real-producer composition test | return before the recovery-ref write in the release path, rebuild the subject binary, run `go test ./internal/contract/runtime -run Recovery -timeout 300s`, expect the plan step to fail on the missing ref |
| GD2 | ignore the supplied fingerprint in the discard arm and exit zero without deleting | the real-producer composition test | drop the delete call, rebuild the subject binary, run `go test ./internal/contract/runtime -run Recovery -timeout 300s`, expect the ref-gone assertion to fail |
| GD3 | skip the assignment compaction after a discard | the real-producer composition test | return before the compaction, rebuild the subject binary, run `go test ./internal/contract/runtime -run Recovery -timeout 300s`, expect the closed-row assertion to fail |
| GD4 | emit the retire terminal action for a discard | the real-producer composition test | reuse the retired action string, rebuild the subject binary, run `go test ./internal/contract/runtime -run Recovery -timeout 300s`, expect the receipt-action assertion to fail |
