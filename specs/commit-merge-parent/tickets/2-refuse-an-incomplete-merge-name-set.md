# Refuse an incomplete name set during a pending merge

Blocked by: 1-read-the-pending-merge-through-one-owner.md
Writes: internal/commit/commitment_test.go, internal/commit/dry_run_test.go, internal/commit/landing_test.go, internal/commit/merge_parent_test.go (new), internal/commit/refusal_route_test.go, internal/landing/attribution.go, internal/landing/composition_test.go, internal/landing/gitexec.go, internal/landing/landing_helpers_test.go, internal/landing/landing_test.go, internal/landing/pending_merge.go (new), internal/landing/pending_merge_fixture_test.go (new), internal/landing/pending_merge_test.go (new), internal/landing/state_test.go, internal/refusalroute/faces_commit.go
Covers: MP24, MP25, MP32, MP33, MP34, MP35, MP36, MP37, MP38, MP39, MP40, MP41, MP42, MP43, MP62

## What to build

During a pending merge, `bench commit` refuses a name set that omits a path that the merge changed, and the refusal names each omitted path.
An unreadable merge state refuses the commit, so a fault cannot pass as "no merge".

Ticket 1 supplies `git.PendingMerge`, and its What to build states the three answers that this ticket reads.
Change `attributedPaths` in `internal/landing/attribution.go` as the complete-name-set decisions of the spec state:

- After the name checks, it reads `git.PendingMerge` for the root.
- A reader refusal returns as a plain error, not as a `NamedPathError`, so `bench commit` prints the `commit-handback` face.
- Without a merge head, it returns the names as it does now, so a staged unnamed path stays allowed.
- With a merge head, it lists the merged paths through `git diff-index --cached --name-only --no-renames -z <expected>`.
- It lists the covered paths through the same command, limited by one `:(literal)` pathspec for each name.
- The omitted paths are the merged paths that the covered listing does not print, in listing order.
- An omitted path raises `NamedPathError`, so `bench commit` prints the `commit-named-path` face.
- The refusal starts with `pending merge changed paths that the commit does not name:`, and each omitted path follows Go-quoted, with one space between paths.

The NUL framing keeps each name byte for byte, so no C-quoting reaches the comparison.
Read the listing through `outputRaw` in `gitexec.go`, which keeps every byte.
With `--no-renames`, a merged rename lists its deleted path and its added path, so the caller names both.
Git's pathspec decides what a name covers, as it does for `compose` and `reconcile`.

`ResolveAttributedPaths`, `CandidateTree`, and `Owner.composeAuthorized` call `attributedPaths`, so all three share the rule.
`bench commit` calls `ResolveAttributedPaths` before `AuthorizeCandidate`, so the D1 refusal prints before a commitment refusal.
`Owner.DryRun` reads the same rule through `composeAuthorized`.
Keep the signature of `attributedPaths` in this ticket, so `internal/landing/landing.go` does not change.
`internal/commit/commit.go` stays unchanged, because it is at its line budget.

Keep the D1 complement, the merged listing less the covered listing, in `attribution.go`.
The coordinator probe of the spec's validation plan swaps the merged listing for the covered listing at that site, and MP32 must then fail.
Put a listing helper in `attribution.go` or in the new `internal/landing/pending_merge.go`.

Update the doc comment of the `commit-named-path` face in `internal/refusalroute/faces_commit.go`, so that it names the pending-merge cause.
The refusal registry gains no face.
No new refusal literal holds `next=`, `next[`, `run bench `, or `; then `.

Put MP24, MP35, MP36, MP37, MP38, and MP42 in the planned tests that their seam cells name in `internal/landing/pending_merge_test.go`.
Each test builds a real merge in a linked worktree and reads the merge state through `git.AdminPath`.
A conflicted `git merge` exits 1, and the `git` helper of the landing tests fails on a nonzero exit.
Give the fixture its own Git call that expects the conflict, as `stageHandMerge` in `internal/worktree/land_refusal_fixture_test.go` does.
Put the linked-worktree merge fixture in the new `internal/landing/pending_merge_fixture_test.go`, so that `pending_merge_test.go` keeps room for ticket 3.

`compositionState` in `landing_helpers_test.go` reads `.git/MERGE_HEAD` directly, and that read fails in a linked worktree.
If a test reads that snapshot, change the read to `git.AdminPath` in place and keep the signature, so `composition_test.go` needs no edit.

In this ticket, a publication during a pending merge still records one parent and leaves MERGE_HEAD in place.
Ticket 3 adds the merge parents and the completion step.
The publishing rows MP37 and MP42 therefore assert the publication and the published tree, not the parent count.

Put MP25, MP32, MP33, MP34, MP39, MP40, MP41, and MP62 in the planned tests that their seam cells name in `internal/commit/merge_parent_test.go`.
These journeys run real Git through the public `Command` seam.
`landingRepo` builds the bound linked worktree, and MP40 uses `planningCommitRepo` for the unbound one.
Each test first resets the staged `tracked.txt` that `prepareLandingCheckout` leaves.
The conflicted merge mirrors `stageHandMerge`, which no shared helper exports.
MP40 asserts only that the D1 refusal prints and that no commitment refusal prints.

The fixture can make the incoming commit on a side branch of the linked worktree.
If the fixture needs the primary checkout instead, extend `landingRepo` in `internal/commit/commitment_test.go`.
`internal/commit/landing_test.go` holds the shared journey helpers, and a helper for the conflicted merge can join them.
`dry_run_test.go` holds the agent-clearable refusal table, and `refusal_route_test.go` holds the commit face fixtures.
Each table can take a pending-merge row, and each current row stays as it is.

MP43 is the current `TestLandComposesOnlyNamedPathsAndCASesExpectedBase` in `internal/landing/landing_test.go`.
That test holds no pending merge, so it stays green without an edit.
`TestLandPreservesTrackedPathCompositionFailureBeforeAuthorization` in `state_test.go` sets `GIT_OBJECT_DIRECTORY` and expects the composition error.
The new merge-state read runs before `compose`, so that test must still receive the composition error.
Keep each changed and new file at 400 lines or fewer.

## Acceptance

- [ ] With a directory at MERGE_HEAD, `Owner.Land` refuses with a plain error, and no ref moves.
- [ ] With a directory at MERGE_HEAD, `bench commit` exits 1, prints the `commit-handback` route, and moves no ref.
- [ ] `bench commit` refuses a pending merge whose name set omits a modified merged path.
- [ ] `bench commit` refuses a pending merge whose name set omits an added merged path.
- [ ] The refusal names each of two omitted paths after the prefix `pending merge changed paths that the commit does not name:`.
- [ ] `Owner.Land` refuses a pending merge whose name set omits the conflicted path.
- [ ] `Owner.Land` refuses a pending merge whose name set omits a path that the merge deleted.
- [ ] `Owner.Land` publishes when a named directory covers every merged path below it.
- [ ] With merged paths `pkg/b.go` and `pkgx/a.go`, `Owner.Land` refuses the name set `pkg` and names `pkgx/a.go`.
- [ ] `bench commit --dry-run` refuses an incomplete merge name set with exit 1.
- [ ] In an unbound planning worktree, the incomplete-merge refusal prints, and no commitment refusal prints.
- [ ] An omitted path with a space, a glob character, an ESC byte, and a newline prints Go-quoted on one refusal line.
- [ ] `Owner.Land` publishes when the name set spells a merged path with a space, a glob character, a newline, and a non-ASCII byte.
- [ ] Without a pending merge, `Owner.Land` publishes while an unnamed path stays staged.
- [ ] The incomplete-merge refusal prints the `commit-named-path` route, and the route ends with the caller's own commit.
- [ ] Each refusal leaves HEAD at its previous commit.
