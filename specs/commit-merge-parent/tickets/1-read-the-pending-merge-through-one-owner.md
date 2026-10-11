# Read the pending merge through one owner

Blocked by: none
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/git/checkout.go, internal/git/checkout_test.go, internal/worktree/admin_readers_test.go, internal/worktree/land_journey_test.go, internal/worktree/land_refusal.go, internal/worktree/parallel_census_test.go
Covers: MP01, MP02, MP03, MP04, MP05, MP06, MP07, MP08, MP09, MP10, MP11, MP12, MP13, MP14, MP15, MP16, MP17, MP18

## What to build

One reader in `internal/git` answers whether a checkout holds a pending merge, and the landing conflict route reads that answer.
A valid MERGE_HEAD gives its merge heads in file order, an absent MERGE_HEAD gives none, and every other state refuses.

Add `git.PendingMerge` to `internal/git/checkout.go`, beside `IsPrimaryCheckout`.
It takes the checkout root and returns the merge heads and an error.
The reader follows the pending-merge owner decisions of the spec:

- It resolves the MERGE_HEAD path through `git.AdminPath`, so a linked worktree reads its own admin entry.
- It runs `os.Lstat` on that path before any read.
- An absent path gives no merge head and no error.
- A path that is not a regular file refuses, and no read occurs.
- The reader splits the content into lines, and the final newline is optional.
- Each line holds only lowercase hexadecimal digits, and `git.ResolveCommit` must answer the line unchanged.
- An empty file refuses.
- Each refusal starts with `read pending merge:`.

The reader calls the current `git.ResolveCommit`, so `internal/git/git.go` does not change.
The MP2 chunk reads `git.PendingMerge` in the landing owner, and the MP3 chunk reads it in `branchReference`.
Both chunks read this contract:

- No error and no merge head mean no pending merge.
- No error and one or more merge heads mean a pending merge with those heads in file order.
- An error means an unreadable merge state.

Change `sourceMergePending` in `internal/worktree/land_refusal.go` to read through `git.PendingMerge`.
It answers true only when the reader returns a merge head, and a reader refusal answers false.
Remove the old `os.ReadFile` of MERGE_HEAD.
That read is the only use of the `os` and `path/filepath` imports (`land_refusal.go:87`), so remove both imports.
Update the doc comment at `land_refusal.go:78-81`, so that it names the one reader and keeps the undecided policy.

Put MP01 to MP14 in the planned tests that their seam cells name in `internal/git/checkout_test.go`.
`TestIsPrimaryCheckoutSeparatesLinkedWorktree` is the precedent for the linked worktree.
Send every Git call through `runGit`.
The branch-native census allows one repository or process constructor in the `internal/git` tests (`internal/git/testhelpers_test.go:18`).
The census also counts a call named `Run` or `Output` with an `init`, `clone`, or `worktree` literal (`internal/conformance/ordinary_build_census_test.go:291-300`).
A merge with `--no-commit --no-ff` and no conflict exits zero, so `runGit` can make the MP01 and MP03 merges.

The FIFO row takes its deadline from `bounds.TestDeadline`, as `worktreesWithin` does.
It skips only through `capability.Capability` with `capability.Fifo`, as `internal/git/refs_test.go:51` does.

Put MP15 in the planned `TestSourceMergePendingReadsTheSharedOwner` in `internal/worktree/admin_readers_test.go`.
The test plants an empty MERGE_HEAD in a source worktree and expects false.
The current file-presence read answers true, so MP15 stays red while that copy survives.
The test calls `t.Parallel()` and binds no environment, PATH, or working directory, so the serial ceiling stays.
Keep `internal/git/checkout_test.go` and `internal/worktree/admin_readers_test.go` at 400 lines or fewer.

Each new top-level test in `internal/worktree` reds `TestPackageTestCountPin`.
Raise `worktreeTestCount` by one for each such test in this ticket.
`internal/worktree/parallel_census_test.go` is over its line budget, so change the pin line in place and add no line.

MP16 is the current `TestSourceMergePendingIsFalseWhenUndecided`, which stays green without an edit.
MP17 and MP18 are LRS6 and LRS22 in `TestLandCommandPublicConflictRepairRequiresNewReviewedTip`.
LRS6 plants a full commit ID and a newline, and LRS22 plants a directory, so both stay green without an edit.
`internal/worktree/land_journey_test.go` is over its line budget, so add no line there.
If the new reader makes a comment there false, correct that comment in place.

`internal/worktree/land_refusal_fixture_test.go:160` also plants a full commit ID and a newline.
The `composition-conflict-pending` fixture of the landing refusal walk therefore stays green without an edit.
The binding registry binds the command registry files to `internal/worktree`, so they join this `Writes:` line.
This ticket does not change them.

## Acceptance

- [ ] `git.PendingMerge` returns the merge head that `git merge --no-commit --no-ff` writes in a linked worktree.
- [ ] Without a MERGE_HEAD, `git.PendingMerge` returns no merge head and no error.
- [ ] A pending octopus merge returns both merge heads in file order.
- [ ] A MERGE_HEAD whose last line has no newline returns its merge head.
- [ ] A root outside every repository refuses.
- [ ] A directory, a FIFO, a dangling symbolic link, and a live symbolic link at MERGE_HEAD each refuse.
- [ ] The FIFO refusal returns within the test deadline.
- [ ] An empty file, an abbreviated commit ID, the line `main`, an unknown full object ID, and a full tree ID each refuse.
- [ ] Each refusal error starts with `read pending merge:`.
- [ ] `sourceMergePending` answers false for an empty MERGE_HEAD and for a source outside every repository.
- [ ] The landing conflict route names `merge --continue` for a source with a valid MERGE_HEAD.
- [ ] The landing conflict route names the commit-and-review form for a directory at MERGE_HEAD.
- [ ] `land_refusal.go` holds no read of MERGE_HEAD of its own.
