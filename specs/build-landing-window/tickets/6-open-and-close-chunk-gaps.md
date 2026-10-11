# Open a gap at a chunk checkpoint and close it at the next dispatch

Blocked by: 5-hold-main-for-an-open-build.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/record_operations_test.go, internal/commitment/repository/gap.go (new), internal/commitment/repository/window.go (new), internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gate/checkpoint.go, internal/gate/chunk_gap.go (new), internal/gate/chunk_gap_test.go (new), internal/gate/complete_checkpoint_test.go, internal/gate/gate.go, internal/gate/refusal_route_test.go, internal/gate/review_checkpoint_commits_test.go, internal/gate/review_checkpoint_test.go, internal/reviewrecord/assignment.go, internal/reviewrecord/plan.go, internal/reviewrecord/recordcmd/assignment.go, internal/reviewrecord/recordcmd/assignment_test.go, internal/reviewrecord/recordcmd/dispatch_gap_test.go (new), internal/reviewrecord/recordtest/fixture.go, internal/worktree/delegated_integration_test.go, internal/worktree/land_window_test.go (new), internal/worktree/parallel_census_test.go
Covers: BW60, BW61, BW62, BW63, BW64, BW65, BW66, BW67, BW68, BW69, BW70, BW95

## What to build

A green chunk checkpoint opens a gap for the bound build, and the first dispatch of a later chunk closes it.
That dispatch refuses while `main` is ahead of the source, so the build folds `main` once before the charge.
Ticket 5 supplies the window that reads the gap record, and ticket 1 supplies the `ChunkGap` record type.

The gate verb opens a gap.
After `RunCommand` (`internal/gate/gate.go:185-210`) returns a green result for `--checkpoint <spec> --chunk <id>`, it reads the plan at HEAD.
It writes the gap record `{Chunk: <id>, Open: true, Last: false}` when two conditions hold.
The plan names a chunk after `<id>`, and the root's active assignment holds a binding of an open build.
Otherwise it writes nothing, so a red checkpoint, a last-chunk checkpoint, and an unbound worktree write no gap.

A write failure prints `chunk gap not recorded: <cause>` and exits 1, and a rerun reuses the green verdict and writes the record.
`internal/gate` already imports `internal/commitment/repository`, so put the gap decision in `internal/commitment/repository/gap.go`.

`bench record assignment` closes a gap.
The record command (`internal/reviewrecord/recordcmd/assignment.go:40`) resolves the dispatched ticket's chunk from the version 2 plan at HEAD.
When the outcome has no gap record, or the dispatched chunk follows the record's chunk in plan order, the dispatch is the chunk's first dispatch.
A first dispatch runs one ledger decision before the plan write:

1. Refuse when the default branch tip is not an ancestor of HEAD.
2. Write `{Chunk: <dispatched chunk>, Open: false, Last: <the chunk is the plan's last>}`.

The refusal prints `dispatch refused: the default branch tip <commit> is not an ancestor of the source; fold it with bench worktree merge --from <branch> <assignment-id> and record it with bench record fold <slug> --source <fold-commit> before the dispatch`, with exit 1.
A dispatch to the record's chunk or to an earlier chunk writes no ledger change.
So a repair session keeps the recorded chunk and its `Last` value.

A dispatch from a worktree without a binding of an open build writes no ledger change.
If the plan write fails after the ledger write, a rerun finds the chunk current and writes the plan.
`internal/reviewrecord` sits below `internal/commitment/repository`, so the ledger decision lives in the record command, not in `internal/reviewrecord`.
Raise `worktreeTestCount` in `internal/worktree/parallel_census_test.go` in place for each new top-level test.

## Acceptance

- [ ] A green `bench gate --checkpoint <spec> --chunk C1` in the bound integration worktree writes an open gap for chunk `C1`.
- [ ] During that open gap, an unbound landing publishes.
- [ ] A green checkpoint of the plan's last chunk, a red checkpoint, and a checkpoint in an unbound worktree each write no gap.
- [ ] The first `bench record assignment` of chunk `C2` writes a closed gap for `C2`.
- [ ] The first dispatch of the first chunk writes a closed gap for that chunk.
- [ ] That dispatch refuses with `dispatch refused: the default branch tip` when `main` moved past HEAD, names `bench worktree merge --from main` and `bench record fold`, and leaves the plan fence unchanged.
- [ ] A second dispatch inside chunk `C2` leaves the ledger bytes unchanged.
- [ ] A repair dispatch to a chunk-`C1` ticket while the record holds chunk `C3` with `Last` true leaves the ledger bytes unchanged.
- [ ] A first dispatch that closes the gap during an unbound landing's gate refuses that landing at publish.
