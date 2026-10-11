# Record each merge head as a parent and end the merge

Blocked by: 1-read-the-pending-merge-through-one-owner.md, 2-refuse-an-incomplete-merge-name-set.md
Writes: internal/commit/merge_parent_test.go (new), internal/landing/attribution.go, internal/landing/composition_test.go, internal/landing/landing.go, internal/landing/landing_helpers_test.go, internal/landing/landing_reviewed_test.go, internal/landing/landing_test.go, internal/landing/lane_owner.go, internal/landing/pending_merge.go (new), internal/landing/pending_merge_fixture_test.go (new), internal/landing/pending_merge_test.go (new), internal/landing/state_test.go
Covers: MP19, MP20, MP21, MP22, MP23, MP26, MP27, MP28, MP29, MP30, MP31, MP44, MP45, MP46, MP47, MP48, MP49, MP50, MP51, MP52, MP53, MP54, MP56

## What to build

During a pending merge, `bench commit` publishes a true merge commit.
Its parent 1 is the previous tip, and each merge head follows in file order.
After the ref update, Git ends the merge, so the checkout holds no pending merge and no unmerged entry.

Ticket 1 supplies `git.PendingMerge`.
Ticket 2 supplies the merge-state read and the D1 rule in `attributedPaths`, so only a complete name set reaches the publication.

Change `attributedPaths` so that it also returns the merge heads that its D1 check read.
`ResolveAttributedPaths` and `CandidateTree` keep their signatures, so `internal/commit/commit.go` stays unchanged.
`Owner.composeAuthorized` keeps those merge heads, so `Owner.Land` and the recheck use that one list.

Change `internal/landing/landing.go` as the parents, recheck, and completion decisions of the spec state:

- A pending merge waives the `nothing to commit` refusal, because the merge parent is the change.
- Without a pending merge, the `nothing to commit` refusal stays as it is.
- `Owner.DryRun` shares the waiver through `composeAuthorized`, and it still publishes nothing.
- After authorization and before `commit-tree`, `Owner.Land` reads `git.PendingMerge` again and compares the answer with the kept list.
- A changed answer refuses with the plain error `pending merge changed during the commit`, and no ref moves.
- A reader refusal at the recheck also refuses before publication.
- `Owner.Land` publishes `commit-tree <tree> -p <expected> -p <merge head>...`, with the merge heads in file order.
- After the ref update, `Owner.Land` runs `git merge --quit` when the commit recorded a merge head.
- The reconcile of the named paths runs after `git merge --quit`, so a reconcile failure leaves no pending merge.
- Every failure after the ref update returns `PublishedUnreconciledError` with the published commit, so `bench commit` exits 3.

Keep the `commit-tree` call and its merge-head parents in `Owner.Land` in `landing.go`.
The central author probe of the spec's validation plan omits the merge heads there, and MP20 must then fail through `bench probe`.
Git's `merge --quit` removes MERGE_HEAD, MERGE_MSG, MERGE_MODE, and AUTO_MERGE, and Bench removes no admin entry file itself.
The MP2 review reads that completion step for MP51.
Update the doc comments of `composeAuthorized` and `Owner.Land`, so that they state the waiver, the recheck, the parents, and the completion step.

Keep each changed and new file at 400 lines or fewer.
Put a new helper in `internal/landing/pending_merge.go` when `landing.go` would pass 400 lines.
`landing_reviewed_test.go` is over its line budget, so add no line there.
`NewLane` and `NewForLane` in `lane_owner.go` start from `New`, so a new `Owner` field reaches the lane owner without an edit.
`composition_test.go`, `landing_reviewed_test.go`, `landing_test.go`, and `state_test.go` call `Owner.Land`, and its signature stays.

Observe MP19 red at the public `Command` seam before the first production edit of this ticket, as the spec's validation plan states.
Put MP19, MP44, MP46, and MP54 in the planned `TestCommitRecordsTheMergeParent` in `internal/commit/merge_parent_test.go`.
Reuse the conflicted-merge fixture of ticket 2, and take the expected parents from `git rev-parse` of the fixture commits.

Put the landing rows of this ticket in the planned tests that their seam cells name in `internal/landing/pending_merge_test.go`.
Those rows are MP20, MP21, MP23, MP26, MP27, MP28, MP29, MP30, MP45, MP47, MP48, MP49, MP50, MP52, and MP53.
Reuse the linked-worktree merge fixture of ticket 2 in `pending_merge_fixture_test.go`.
Read the parents with `git rev-list --parents -n 1`, the merge state through `git.AdminPath`, and the index with `git ls-files -u`.
The parent assertion mirrors `TestLandDetachedHEADHasOneExpectedParent` and `TestMergePublishesPreviousTipAsFirstParent`.

The recheck rows replace `o.authorize` with a stub that changes the merge state, as `TestLandReconcilesAuthorizedNamedSnapshotAfterAuthorizationMutation` does.
One more case of the recheck test plants a directory at MERGE_HEAD during authorization, because a reader refusal at the recheck also refuses.
The remainder rows replace `o.reconcile`, as `TestLandReportsPublishedCommitWhenReconciliationFails` does.
Before the commit, the completion test asserts that each of the four merge-state files exists, so each absence assertion can fail.
MP29 resolves the conflict with the ours side, so the composed tree equals the base tree.
MP30 makes a pending merge with no merged path, for example through `git merge --no-commit --no-ff -s ours <commit>`.

MP22 is the current `TestLandDetachedHEADHasOneExpectedParent` in `internal/landing/landing_test.go`.
MP31 is the tracked-symlink row of `TestLandPreAuthorizationRefusalTable` in `internal/landing/state_test.go`.
Both stay green without an edit.

Ticket 4 states this merge contract in `.bench/BENCH-reference.md`, so this ticket changes no guidance.
The delivery changes neither `ROADMAP.md` nor `roadmap/FT258.md`, and the next `/bench-drain` reconciles that row.
The fence review reads MP56: the fence names no roadmap path and does not name `internal/commit/commit.go`.

## Acceptance

- [ ] After a conflicted merge, `bench commit` with the complete name set publishes the previous tip and the incoming commit as the two parents.
- [ ] `Owner.Land` publishes the expected base as parent 1 and the merge head as parent 2.
- [ ] A pending octopus merge publishes both merge heads as parents 2 and 3 in file order.
- [ ] Without a pending merge, `Owner.Land` publishes exactly one parent.
- [ ] After the merge commit, `Owner.Compose` of the source onto the incoming commit reports no conflict.
- [ ] A MERGE_HEAD that appears during authorization makes `Owner.Land` refuse, and no ref moves.
- [ ] A MERGE_HEAD that disappears during authorization makes `Owner.Land` refuse, and no ref moves.
- [ ] A MERGE_HEAD that names another commit after authorization makes `Owner.Land` refuse, and no ref moves.
- [ ] A directory at MERGE_HEAD that appears during authorization makes `Owner.Land` refuse before publication.
- [ ] After an ours resolution, `Owner.Land` publishes a merge commit whose tree equals the base tree.
- [ ] A pending merge with no merged path publishes a merge commit with an unchanged tree.
- [ ] Without a pending merge, an unchanged tree still refuses with `nothing to commit`.
- [ ] After the merge commit, a second `bench commit` of a new change publishes one parent.
- [ ] After the merge commit, `git ls-files -u` prints nothing.
- [ ] After the merge commit, `git status --porcelain=v1` prints nothing in the worktree.
- [ ] After the merge commit, MERGE_HEAD, MERGE_MSG, MERGE_MODE, and AUTO_MERGE are absent from the checkout administration directory.
- [ ] The completion step runs `git merge --quit`, and Bench removes no merge-state file itself.
- [ ] When the reconcile fails during a pending merge, `Owner.Land` returns `PublishedUnreconciledError` with the published merge commit.
- [ ] After that reconcile failure, MERGE_HEAD is absent.
- [ ] The last stdout line of a merge commit is `committed{published_commit,path_count}` with the count of the name set.
- [ ] The delivery changes no roadmap path and leaves `internal/commit/commit.go` unchanged.
