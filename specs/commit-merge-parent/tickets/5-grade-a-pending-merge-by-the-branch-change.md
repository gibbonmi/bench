# Grade a pending merge by the branch's own change

Blocked by: 1-read-the-pending-merge-through-one-owner.md, 2-refuse-an-incomplete-merge-name-set.md, 3-record-each-merge-head-and-end-the-merge.md, 4-document-the-merge-contract-of-bench-commit.md
Writes: CHANGELOG.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/commitment/repository/candidate.go, internal/commitment/repository/stale_base_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary
Covers: MP57, MP58, MP59, MP60, MP61

## What to build

During a pending merge, the commitment grading of a `bench commit` candidate reads only the branch's own change.
A merged commit of the default branch counts as the default branch's own.
A change that a merged sibling made still counts as the branch's change.

Ticket 1 supplies `git.PendingMerge`, and its What to build states the three answers that this ticket reads.
Tickets 2 and 3 supply the merge commit of `bench commit`, and the `CHANGELOG.md` entry of this ticket states that behavior.
Ticket 4 also names the command registry files, so this ticket starts after ticket 4 commits.

Change `branchReference` in `internal/commitment/repository/candidate.go` as the merge-grading decisions of the spec state:

- It reads `git.PendingMerge` for the store root before it runs `git merge-base`.
- A reader refusal returns the default tip as the reference, which is the stricter grading.
- With merge heads, it runs `git merge-base <default tip> HEAD <merge head>...`.
- Without a merge head, it keeps `git merge-base HEAD <default tip>`.
- A base equal to the tip still gives the tip, and a policy identity mismatch still gives the tip.

Update the doc comment of `branchReference`, so that it states the pending-merge reference and the reader-refusal fallback.
`candidate.go` reads the merge state only through `git.PendingMerge`, so it spells no `--git-path` and holds no MERGE_HEAD read of its own.
Keep `candidate.go` and `stale_base_test.go` at 400 lines or fewer.

`authorizeCandidate` calls `branchReference` only for a candidate that is not a publication (`candidate.go:52-55`).
The landing publication at `internal/commitment/repository/publication.go:51` passes `publication` as true, so its grading does not change.
`bench commit` refuses an unreadable merge state in `attributedPaths` first, as ticket 2 builds it.
The reader-refusal fallback therefore serves only a direct caller of `AuthorizeCandidate`.
`internal/commit/commit.go` stays unchanged.

Put MP57, MP58, MP59, and MP60 in the planned `TestPendingMergeCandidate` in `internal/commitment/repository/stale_base_test.go`.
Each row grades an unbound assignment of `lightPathRepo`, as the testing decisions of the spec state.
The candidate of each row is the `git write-tree` of its index.
`TestStaleBaseCandidate` is the precedent, and `landProduction` and the step helpers of `light_path_test.go` serve the rows without an edit.
A merge with `--no-commit --no-ff` and no conflict exits zero, so `gittest.Output` can make each merge.
No row injects a merge state, and each expected path comes from the fixture writes.

The rows build these states:

- MP57 commits a light-path ticket and its production path, moves main through `landProduction`, and merges main.
- MP58 commits only a planning path, moves main through `landProduction`, and merges main.
- MP59 commits a light-path ticket, then merges a sibling commit that main does not contain and that adds an uncovered production path.
- MP60 commits a light-path ticket and its production path, moves main through `landProduction`, and plants a directory at MERGE_HEAD through `git.AdminPath`.

MP60 merges nothing, so its candidate lacks the production change of main.

Before the production edit, MP57, MP58, and MP60 fail, because the current reference reads no merge state.
MP59 passes before the edit, so its red comes from a behavioral swap.
One such swap puts the merge head in place of the default tip as the first `git merge-base` operand.
MP61 is the current `TestStaleBaseCandidate`, which stays green without an edit.

The binding registry binds the command registry files to `internal/commitment`, so they join this `Writes:` line.
The anchor registry files name `CHANGELOG.md`, and two canary units pin it, so these files also join this `Writes:` line.
Each canary unit replaces one older entry, so the new section changes no unit.
This ticket does not change the registry files, the anchor registry files, or the canary units.

Add one section under `## [Unreleased]` in `CHANGELOG.md`, above `### Commit success record`.
The section holds one `Fixed` entry, because the spec calls this delivery a fix.
The entry follows the form of the current entries, and it states these four behaviors:

- The merge parents: during a pending merge, a worktree `bench commit` records the previous tip and then each merge head as parents.
- The merge completion: after the publication, the commit ends the merge through Git.
- The name-set refusal: the commit refuses a name set that omits a merged path, and the refusal names each omitted path.
- The merge grading: the commitment grading counts the merged commits of the default branch as the default branch's own.

A last sentence that starts with `Before,` states the old behavior: one parent, and a merge that stayed pending.
Keep the entry under `### Commit success record` unchanged.

## Acceptance

- [ ] In a light-path worktree, `AuthorizeCandidate` admits a complete merge of main that carries a production change of main.
- [ ] In an unbound planning worktree, `AuthorizeCandidate` admits a complete merge of main that carries a production change of main.
- [ ] In a light-path worktree, `AuthorizeCandidate` refuses a pending merge of a sibling commit that adds an uncovered production path, and the refusal names that path.
- [ ] With a directory at MERGE_HEAD, `AuthorizeCandidate` refuses a light-path candidate without the later production change of main, and names that production path.
- [ ] Without a pending merge, `AuthorizeCandidate` keeps every stale-base admission and refusal.
- [ ] `candidate.go` holds no read of MERGE_HEAD of its own.
- [ ] One typed `CHANGELOG.md` entry under `## [Unreleased]` states the merge parents, the merge completion, the name-set refusal, and the merge grading.
