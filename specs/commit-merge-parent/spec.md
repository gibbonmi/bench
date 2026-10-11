# Record the merge parent in bench commit

Status: staged

Related roadmap: FT258, partial. This spec settles the merge-parent rule. The candidate-set forms `--changed`, `--dry-run --changed`, and `--ticket`, and their derivation decision, stay on FT258.

Decision source: named reviewed artifact `roadmap/FT258.md` at `1a947e86`, with its 2026-10-10 reviewer decision from drain commit `0e3eedc0`, the `FT258-merge-parent` criteria, and the reviewer decisions D1 and D2 of 2026-10-10

Audience: every repository that links the kit

Verification log: 2 iteration(s) to accept — the spec review and the ticket review each accepted in their first iteration

## Problem

A conflict in `bench worktree merge` or in `bench worktree land` sends the reviewer to a hand `git merge` and then to `bench commit`.
The `handMerge` route states those two steps (`internal/refusalroute/registry.go:189-192`).
`bench commit` publishes through `landing.Owner.Land`, which gives `commit-tree` one parent (`internal/landing/landing.go:163`).
The resolution commit therefore drops the merged commit from its history, and MERGE_HEAD stays pending.

Each recorded occurrence then needed a raw commit that recorded the second parent.
The 2026-08-27 build, the 2026-09-09 landing, the 2026-09-12 light-path batch, and the 2026-09-13 pilot show the same pattern.
The 2026-09-09 commit named only the conflicted file, so its other merged paths stayed staged outside the commit.
The destructive-git guard now denies a raw commit and a raw merge in every form (`internal/gitguard/verdict.go:26-27`, `:52-55`).
No agent route records the merge parent today.

## Solution

During a pending merge, `bench commit` publishes a true merge commit.
The commit keeps the previous tip as its first parent, and it records each merge head after that parent.
The commit refuses a name set that omits a path that the merge changed, and the refusal names each omitted path.
After the commit publishes, it ends the merge through Git, so the checkout holds no pending merge.
The commitment grading of a merge commit reads the merged commits of the default branch as the default branch's own.

## User stories

Line: opus / xhigh
Implementation-line reason: MP2 is the hardest material chunk. It changes the composition, publication, and completion steps of the landing owner, and real Git merges prove each step. Its tickets cross owner seams, and a repair round there costs over 100k tokens, so the scorecard routes them to Opus at xhigh. Each ticket author starts on this line, and `craft-line` owns any effort change inside that session.
Harder chunks: MP2
Run line, by reviewer direction on 2026-10-11: opus / xhigh for each ticket author, opus / high for each review axis, and opus / xhigh for each consultation. An opus / max shadow author on tickets 1, 2, and 3 supplies the effort assessment, and its work lands nothing.

### The merge parent

1. As an agent, I want `bench commit` to record the merge head as parent 2, so that the landing reads the merged commit as an ancestor.
2. As an agent, I want the first parent to stay the previous tip of the checkout, so that first-parent readers keep their chain.
3. As a reviewer, I want `bench commit` never to author a one-parent commit during a pending merge, so that no resolution drops the incoming history.
4. As an agent, I want a commit without a pending merge to keep exactly one parent, so that ordinary commits do not change.
5. As an agent who finishes an octopus merge, I want each merge head recorded in file order, so that no pending parent goes missing.
6. As an agent, I want the landing after my merge commit to compose without the resolved conflict, so that I resolve each conflict once.
7. As an agent with an ours resolution, I want the merge commit to publish with the base tree, so that the history records the merge.
8. As an agent whose merge changed no path, I want the merge commit to publish, so that I need no raw commit for the ancestry.
9. As an agent, I want the empty-diff refusal to stay when no merge is pending, so that a commit with no change still refuses.
10. As a reviewer, I want the commit to refuse when the merge state changes during authorization, so that the published parents match the graded state.

### The merge state reader

11. As an agent in a linked worktree, I want the commit to read MERGE_HEAD in its checkout administration directory, so that the merge stays visible.
12. As a reviewer, I want an unreadable merge state to refuse the commit, so that a fault cannot pass as "no merge".
13. As a reviewer, I want an empty MERGE_HEAD to refuse, so that absence and emptiness stay two different states.
14. As a reviewer, I want a MERGE_HEAD line that is not a full commit ID to refuse, so that only Git's record becomes a parent.
15. As a reviewer, I want a link or a special file at MERGE_HEAD to refuse unread, so that the commit neither blocks nor follows links.
16. As an agent, I want a MERGE_HEAD without a final newline to give its head, so that a hand edit does not break the commit.
17. As a maintainer, I want one reader to own the pending-merge fact, so that the commit, the grading, and the landing route cannot drift.
18. As a reviewer, I want the landing conflict route to treat a malformed merge state as undecided, so that it obeys the one reader.
19. As a reviewer, I want the pending-merge route of the landing to keep `git merge --continue`, so that this delivery changes no landing route.

### The complete name set

20. As an agent, I want `bench commit` to refuse when I omit a modified merged path, so that the landing cannot revert that change.
21. As an agent, I want an omitted added path to refuse, so that the landing cannot drop the incoming file.
22. As an agent, I want an omitted conflicted path to refuse, so that an unresolved conflict cannot stay out of the merge commit.
23. As an agent, I want an omitted path that the merge deletes to refuse, so that an incoming deletion does not go missing.
24. As an agent, I want the refusal to name every omitted path, so that one rerun corrects the name set.
25. As an agent, I want a named directory to cover the merged paths below it, so that one name covers a merged subtree.
26. As an agent, I want a directory name not to cover a sibling that only shares its prefix, so that the coverage rule stays exact.
27. As an agent, I want `bench commit --dry-run` to refuse an incomplete name set, so that a rehearsal predicts the commit.
28. As an agent, I want the incomplete-merge refusal to print before a commitment refusal, so that I see the merge cause first.
29. As an agent, I want the refusal to quote each omitted path, so that hostile bytes in a path cannot break the refusal line.
30. As an agent, I want a merged path with hostile bytes to match its literal name, so that a correct name set does not refuse.
31. As an agent, I want staged unnamed paths to stay allowed when no merge is pending, so that the current commit contract does not change.
32. As a reviewer, I want the caller to keep the duty to name every path, so that status discovery grants no write authority.

### The merge completion

33. As an agent, I want the commit to end the pending merge after it publishes, so that a later commit records no stale merge head.
34. As an agent, I want no unmerged index entry after the merge commit, so that the worktree status is clean.
35. As an agent, I want each merge-state file of Git to be gone after the merge commit, so that Git reports no merge in progress.
36. As a maintainer, I want Git itself to end the merge, so that Bench writes no worktree admin entry.
37. As an agent, I want a failure after publication to report the published merge commit with exit 3, so that the current restore route applies.
38. As an agent, I want no pending merge after a reconcile failure, so that the restore route repairs only the checkout.

### The merge grading

39. As an agent in a light-path worktree, I want a complete merge of main to pass the commitment grading, so that Bench completes it.
40. As a spec author in an unbound worktree, I want a complete merge of main to pass the commitment grading, so that planning continues.
41. As a reviewer, I want a merged sibling's own changes to grade as the branch's own, so that a merge carries no unauthorized change.
42. As a reviewer, I want an unreadable merge state to grade against the default tip, so that a fault takes the stricter grading.
43. As an agent, I want the grading without a pending merge to stay unchanged, so that ordinary commits keep their admission.

### The unchanged surfaces and the scope

44. As an agent, I want the record of a merge commit to keep its current fields, so that no reader of the commit output changes.
45. As a teammate who walks in cold, I want the working reference to state the commit's merge contract, so that I name every merged path.
46. As the reviewer, I want the delivery to close no roadmap row and to keep the candidate-set forms on FT258, so that they stay open.

## Implementation decisions

### Terms

- **pending merge**: the state that a MERGE_HEAD file records in the checkout administration directory.
- **merge head**: the commit ID on one MERGE_HEAD line.
- **merged path**: one path that `git diff-index --cached --name-only --no-renames -z <base>` lists during a pending merge.
- **name set**: the paths that the caller names on one `bench commit` call.

The spec uses these terms only in its own text. It adds no CONTEXT.md entry.

D1 names the merged-path listing.
Git aborts a merge when the index differs from HEAD (`git help merge`, PRE-MERGE CHECKS, Git 2.43.0), so the listing starts as the record of the merge.
A path that the caller stages after the merge also joins the listing, and the caller then names that path too.

### Closed decisions

The `FT258-merge-parent` criteria (`.bench/commitment.json:452-465`) state the outcome:

- `.parent`: "When MERGE_HEAD exists, bench commit records the merge parent as the second parent of the commit."
- `.pending`: "bench commit never authors a one-parent commit while merge state is pending."
- `.open`: "The delivery closes no roadmap row, and the candidate-set forms --changed, --dry-run --changed, and --ticket stay on FT258."

The reviewer decided two forks in this phase on 2026-10-10:

- D1: During a pending merge, `bench commit` refuses unless the named paths include every path that the merge changed, as `git diff-index --cached HEAD` lists them. The refusal lists each omitted path through the existing `commit-named-path` face. The caller still names every path, so status discovery grants no write.
- D2: During a pending merge, the commitment grading of the commit reads the branch's own change from `git merge-base <default tip> HEAD <each merge head>`. Only commits on the default branch count as the default branch's own, so the changes that a merged sibling made still grade as the branch's.

The row sentence at `roadmap/FT258.md:17` ("A conflict-resolution lane commit does not complete that merge state.") describes the defect before this fix.
The repair-collection-pilot retro in commit `f90d3a28` asked to "Explain that a conflict-resolution lane commit does not complete `MERGE_HEAD`".
The roadmap-light-path retro in commit `4717f4c7` (lines 67-69) reports a single-parent commit that "left MERGE_HEAD pending".
This spec makes the lane commit a true merge commit that ends the merge.

### The pending-merge owner

`internal/git` owns the one pending-merge reader, beside `IsPrimaryCheckout` in `checkout.go`.
The reader resolves MERGE_HEAD through `git.AdminPath`, because the `git-plumbing-owner` check reds a hand-spelled `--git-path` outside `internal/git` (`internal/conformance/git_plumbing_owner_test.go:36-54`).
The reader runs `os.Lstat` on the path before any read.
An absent path gives no merge head and no error.
A path that is not a regular file refuses, and no read occurs.
That rule covers a directory, a FIFO, a socket, a device, and a symbolic link.

The reader splits the content into lines, and a final newline is optional.
Each line holds only lowercase hexadecimal digits, and `git.ResolveCommit` must answer the line unchanged (`internal/git/git.go:220-223`).
That rule refuses an abbreviated ID, a branch name, an unknown object, and an object that is not a commit.
An empty file refuses.
The reader returns the merge heads in file order.
Each refusal carries the sentence prefix `read pending merge:`.

Three consumers read through this owner: the landing owner, `branchReference` (`internal/commitment/repository/candidate.go:128-137`), and `sourceMergePending` (`internal/worktree/land_refusal.go:82-91`).
`sourceMergePending` answers true only when the reader returns a merge head.
A reader refusal answers false, which keeps the undecided policy of its comment at `internal/worktree/land_refusal.go:78-81`.
The old file read in `sourceMergePending` goes away.

### The complete name set

`landing.attributedPaths` applies the D1 rule after it resolves the names.
Its three callers then share the rule: `ResolveAttributedPaths`, `CandidateTree`, and `Owner.composeAuthorized`.
`bench commit` calls `ResolveAttributedPaths` first (`internal/commit/commit.go:136`), before `AuthorizeCandidate` (`internal/commit/commit.go:144`).
The D1 refusal therefore prints before a commitment refusal.

The merged paths come from `git diff-index --cached --name-only --no-renames -z <base>`, where the base is the expected commit.
The NUL framing keeps each name byte for byte, so no C-quoting reaches the comparison (observation 7).
The listing names an unmerged path once, and it names a deleted path (observations 1 and 7).
With `--no-renames`, a merged rename appears as its deleted path and its added path.

The covered paths come from the same listing, limited by one `:(literal)` pathspec for each name.
Git's pathspec therefore decides what a name covers, as it does for `compose` and `reconcile` (`internal/landing/attribution.go:154`, `:177`).
A directory name covers the merged paths below it, and `pkg` does not cover `pkgx/a.go` (observation 7).
The omitted paths are the merged paths that the limited listing does not print.
D1 runs after the name checks of `attributedPaths` and before `compose`.

An uncovered merged path raises `landing.NamedPathError`, so the commit prints the `commit-named-path` face.
The sentence starts with `pending merge changed paths that the commit does not name:`.
Each uncovered path follows in Go-quoted form, in listing order, with one space between paths.
The rule applies only during a pending merge, so a staged unnamed path stays allowed otherwise.

### The parents and the empty tree

`Owner.composeAuthorized` keeps the merge heads that its D1 check reads, so `Owner.Land` and the recheck use that one list.
A pending merge waives the `nothing to commit` refusal, because the merge parent is the change.
Without a pending merge, that refusal stays as it is (`internal/landing/landing.go:139-141`).
`Owner.Land` publishes `commit-tree <tree> -p <expected> -p <merge head>...`, with the merge heads in file order.
`Owner.DryRun` applies the same D1 rule and the same waiver, and it publishes nothing.

### The recheck

After authorization and before `commit-tree`, `Owner.Land` reads the pending merge again and compares it with the list that the composition used.
A changed answer refuses with the plain error `pending merge changed during the commit`, and no ref moves.
A changed answer is a merge that appears, a merge that disappears, or another merge head list.
A reader refusal at the recheck also refuses before publication.
The authorization of a commit can take minutes, so this recheck guards the `.pending` criterion.
The `LandReviewed` recheck of its moving identities is the precedent (`internal/landing/landing.go:247-257`).

### The merge completion

After the ref update, `Owner.Land` runs `git merge --quit` when the commit recorded a merge head.
The reconcile of the named paths runs after that step.
`git merge --quit` removes MERGE_HEAD, MERGE_MSG, MERGE_MODE, and AUTO_MERGE, and it keeps the index and the worktree.
The reconcile then restores each named path from the composed tree, which also clears its unmerged index entries.
Git stays the only writer of its worktree admin entry, as CONTEXT.md states for that term.
Every failure after the ref update returns `landing.PublishedUnreconciledError`, so `bench commit` exits 3 with the current route.

### The merge grading

`branchReference` reads the pending merge of the store root.
With merge heads, it runs `git merge-base <default tip> HEAD <merge head>...`.
`git help merge-base` defines that form as the base of the first commit and a hypothetical merge of the other commits (Git 2.43.0).

A reader refusal leaves the default tip as the reference, which is the stricter grading that the function already applies to a failure.
Without a merge head, it keeps `git merge-base HEAD <default tip>`.
The rest of `branchReference` stays: a base equal to the tip gives the tip, and a policy identity mismatch gives the tip.
`bench commit` refuses an unreadable merge state in `attributedPaths`, so the fallback serves only a direct caller of `AuthorizeCandidate`.

### The unchanged surfaces

The commit refusals travel the current faces, so the refusal registry gains no face.
D1 raises `commit-named-path`.
A reader refusal and a recheck refusal raise a plain error, which `landingFace` maps to `commit-handback` (`internal/commit/commit.go:215-229`).
A failure after publication raises `commit-published-unreconciled`.
No new refusal literal holds `next=`, `next[`, `run bench `, or `; then ` (`internal/conformance/refusal_route_bypass_test.go`).

The `committed{published_commit,path_count}` record, the dry-run lines, the help text, and the `bench help` row stay unchanged.
`internal/commit/commit.go` stays unchanged, because it is at its line budget, and the growth ratchet reds a file that grows past its limit (`internal/structure/structure.go:283-339`).
The landing faces stay unchanged: `composition-conflict-pending` keeps `git merge --continue`, and `handMerge` keeps its bare `bench commit` step.

`.bench/BENCH-reference.md` gains two sentences after the fast-lane note in its Command Notes.
The proposed text is: "During a pending merge, a worktree `bench commit` requires every path that the merge changed. It records each merge head as a further parent and then ends the merge."
The anchors of that file stay intact, including the Command Notes anchor at line 280 and the route anchor at line 307.

`CHANGELOG.md` gains one `Fixed` entry for this delivery under `## [Unreleased]`, and the commit record entry stays unchanged.

### The roadmap

The delivery edits neither `ROADMAP.md` nor `roadmap/FT258.md`.
The next `/bench-drain` reconciles the merge-parent sentences of `roadmap/FT258.md:15-21` after the delivery lands.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| MP1 / 1-read-the-pending-merge-through-one-owner.md | One owner reads the pending merge of any checkout, and the landing conflict route reads it through that owner | MP01-MP18 | the `internal/git` checkout tests, the worktree admin reader tests, and the LRS6 and LRS22 journey | no |
| MP2 / 2-refuse-an-incomplete-merge-name-set.md, 3-record-each-merge-head-and-end-the-merge.md, 4-document-the-merge-contract-of-bench-commit.md | `bench commit` publishes a true merge commit for a complete name set, refuses an incomplete name set, and ends the merge | MP19-MP56, MP62 | the `internal/landing` pending-merge tests, the `internal/commit` merge tests, the current `Land` and refusal tables, the `docs-currency-workflow` and `skills-index-command-adapters` checks, and the `internal/anchors` package tests | yes |
| MP3 / 5-grade-a-pending-merge-by-the-branch-change.md | The commitment grading of a pending merge reads the merged commits of the default branch as the default branch's own | MP57-MP61 | the repository pending-merge candidate test, `TestStaleBaseCandidate`, the `docs-currency-workflow` check, and the `internal/anchors` package tests | no |

MP1 creates the seam that MP2 and MP3 consume, so it stays small, and its chunk review closes before a consumer ticket starts.
MP2 and MP3 change disjoint files, but tickets 4 and 5 both list the command registry files and two anchor registry files in `Writes:`.
Preflight requires those closure files, so the MP3 ticket waits for the last MP2 ticket.
The serial order is MP1, MP2, and then MP3, because MP2 delivers the `.parent` and `.pending` criteria.
The orchestrator reconciles the final acceptance after the MP3 checkpoint.

### Completion plan

This version 2 completion plan declares the delegated run and one verification inventory for each ticket.
Its commands are future obligations, not executed proof.
Each probe field names the mutation kind, its site, and the test that turns red, and the author derives the exact mutated text.

```bench-completion-plan
{
  "version": 2,
  "execution": {
    "mode": "delegate",
    "run_id": "commit-merge-parent-build-20261011",
    "orchestrator_session": "claude:cmp-orchestrator-20261011",
    "author_limit": 1,
    "assignments": {
      "1-read-the-pending-merge-through-one-owner.md": [],
      "2-refuse-an-incomplete-merge-name-set.md": [],
      "3-record-each-merge-head-and-end-the-merge.md": [],
      "4-document-the-merge-contract-of-bench-commit.md": [],
      "5-grade-a-pending-merge-by-the-branch-change.md": []
    }
  },
  "chunks": [
    {
      "id": "MP1",
      "tickets": ["1-read-the-pending-merge-through-one-owner.md"],
      "verification": [
        {"id": "t1-reader-focused", "ticket": "1-read-the-pending-merge-through-one-owner.md", "command": "bench test --package ./internal/git --run 'TestPendingMerge'", "probe": "Omit the regular-file refusal after os.Lstat in git.PendingMerge in internal/git/checkout.go, so that a directory at MERGE_HEAD reaches the read. TestPendingMergeRefusesUnreadableState turns red."},
        {"id": "t1-route-focused", "ticket": "1-read-the-pending-merge-through-one-owner.md", "command": "bench test --package ./internal/worktree --run 'TestSourceMergePending|TestLandCommandPublicConflictRepairRequiresNewReviewedTip|TestPackageTestCountPin'"},
        {"id": "t1-floor", "ticket": "1-read-the-pending-merge-through-one-owner.md", "command": "bench test --package ./internal/git --package ./internal/worktree --package ./internal/conformance --package ./cmd/bench"},
        {"id": "t1-root-conformance", "ticket": "1-read-the-pending-merge-through-one-owner.md", "command": "bench test --check root-conformance"},
        {"id": "t1-git-plumbing-owner", "ticket": "1-read-the-pending-merge-through-one-owner.md", "command": "bench test --check git-plumbing-owner"},
        {"id": "t1-ordinary-build-census", "ticket": "1-read-the-pending-merge-through-one-owner.md", "command": "bench test --check ordinary-build-census"},
        {"id": "t1-skip-ownership", "ticket": "1-read-the-pending-merge-through-one-owner.md", "command": "bench test --check skip-ownership"}
      ]
    },
    {
      "id": "MP2",
      "tickets": ["2-refuse-an-incomplete-merge-name-set.md", "3-record-each-merge-head-and-end-the-merge.md", "4-document-the-merge-contract-of-bench-commit.md"],
      "verification": [
        {"id": "t2-landing-focused", "ticket": "2-refuse-an-incomplete-merge-name-set.md", "command": "bench test --package ./internal/landing --run 'TestLandRefusesAnIncompleteMergeNameSet|TestLandRefusesAnUnreadableMergeState|TestLandComposesOnlyNamedPathsAndCASesExpectedBase'", "probe": "Omit the NamedPathError refusal that attributedPaths raises for the omitted merged paths in internal/landing/attribution.go, so that an incomplete name set composes. TestLandRefusesAnIncompleteMergeNameSet turns red."},
        {"id": "t2-commit-focused", "ticket": "2-refuse-an-incomplete-merge-name-set.md", "command": "bench test --package ./internal/commit --run 'TestCommitRefusesAnIncompleteMergeNameSet|TestCommitRefusesAnUnreadableMergeState|TestDryRunRefusesAnIncompleteMergeNameSet|TestMergeNameSetRefusalPrecedesCommitment'"},
        {"id": "t2-floor", "ticket": "2-refuse-an-incomplete-merge-name-set.md", "command": "bench test --package ./internal/git --package ./internal/landing --package ./internal/commit --package ./internal/refusalroute"},
        {"id": "t2-root-conformance", "ticket": "2-refuse-an-incomplete-merge-name-set.md", "command": "bench test --check root-conformance"},
        {"id": "t2-git-plumbing-owner", "ticket": "2-refuse-an-incomplete-merge-name-set.md", "command": "bench test --check git-plumbing-owner"},
        {"id": "t2-ordinary-build-census", "ticket": "2-refuse-an-incomplete-merge-name-set.md", "command": "bench test --check ordinary-build-census"},
        {"id": "t3-landing-focused", "ticket": "3-record-each-merge-head-and-end-the-merge.md", "command": "bench test --package ./internal/landing --run 'TestLandRecordsEachMergeHeadAsAParent|TestComposeAfterMergeCommitReadsTheMergedBase|TestLandRefusesAMergeStateChangeDuringAuthorization|TestLandPublishesAMergeWithAnUnchangedTree|TestLandEndsThePendingMerge|TestLandReportsPublishedMergeWhenReconciliationFails|TestLandDetachedHEADHasOneExpectedParent|TestLandPreAuthorizationRefusalTable'", "probe": "Omit the merge-head -p arguments from the commit-tree call of Owner.Land in internal/landing/landing.go, so that the publication records one parent during a pending merge. TestLandRecordsEachMergeHeadAsAParent turns red."},
        {"id": "t3-commit-focused", "ticket": "3-record-each-merge-head-and-end-the-merge.md", "command": "bench test --package ./internal/commit --run 'TestCommitRecordsTheMergeParent'"},
        {"id": "t3-floor", "ticket": "3-record-each-merge-head-and-end-the-merge.md", "command": "bench test --package ./internal/git --package ./internal/landing --package ./internal/commit"},
        {"id": "t3-root-conformance", "ticket": "3-record-each-merge-head-and-end-the-merge.md", "command": "bench test --check root-conformance"},
        {"id": "t3-git-plumbing-owner", "ticket": "3-record-each-merge-head-and-end-the-merge.md", "command": "bench test --check git-plumbing-owner"},
        {"id": "t3-ordinary-build-census", "ticket": "3-record-each-merge-head-and-end-the-merge.md", "command": "bench test --check ordinary-build-census"},
        {"id": "t4-guidance", "ticket": "4-document-the-merge-contract-of-bench-commit.md", "command": "bench test --check docs-currency-workflow"},
        {"id": "t4-skills-index", "ticket": "4-document-the-merge-contract-of-bench-commit.md", "command": "bench test --check skills-index-command-adapters"},
        {"id": "t4-floor", "ticket": "4-document-the-merge-contract-of-bench-commit.md", "command": "bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench"},
        {"id": "t4-root-conformance", "ticket": "4-document-the-merge-contract-of-bench-commit.md", "command": "bench test --check root-conformance"}
      ]
    },
    {
      "id": "MP3",
      "tickets": ["5-grade-a-pending-merge-by-the-branch-change.md"],
      "verification": [
        {"id": "t5-grading-focused", "ticket": "5-grade-a-pending-merge-by-the-branch-change.md", "command": "bench test --package ./internal/commitment/repository --run 'TestPendingMergeCandidate|TestStaleBaseCandidate'", "probe": "Swap the default tip for the merge head as the first git merge-base operand in branchReference in internal/commitment/repository/candidate.go, so that the merged sibling's change grades as the default branch's own. TestPendingMergeCandidate turns red."},
        {"id": "t5-floor", "ticket": "5-grade-a-pending-merge-by-the-branch-change.md", "command": "bench test --package ./internal/git --package ./internal/landing --package ./internal/commit --package ./internal/commitment/repository --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench"},
        {"id": "t5-guidance", "ticket": "5-grade-a-pending-merge-by-the-branch-change.md", "command": "bench test --check docs-currency-workflow"},
        {"id": "t5-root-conformance", "ticket": "5-grade-a-pending-merge-by-the-branch-change.md", "command": "bench test --check root-conformance"},
        {"id": "t5-git-plumbing-owner", "ticket": "5-grade-a-pending-merge-by-the-branch-change.md", "command": "bench test --check git-plumbing-owner"},
        {"id": "t5-ordinary-build-census", "ticket": "5-grade-a-pending-merge-by-the-branch-change.md", "command": "bench test --check ordinary-build-census"}
      ]
    }
  ],
  "final_verification": [
    {"id": "coverage", "command": "bench coverage --check specs/commit-merge-parent/spec.md"},
    {"id": "owner-packages", "command": "bench test --package ./internal/git --package ./internal/landing --package ./internal/commit"},
    {"id": "grading-package", "command": "bench test --package ./internal/commitment/repository"},
    {"id": "route-package", "command": "bench test --package ./internal/worktree"},
    {"id": "refusal-registry", "command": "bench test --package ./internal/refusalroute"},
    {"id": "anchors-package", "command": "bench test --package ./internal/anchors"},
    {"id": "command-package", "command": "bench test --package ./cmd/bench"},
    {"id": "conformance-package", "command": "bench test --package ./internal/conformance"},
    {"id": "guidance", "command": "bench test --check docs-currency-workflow"},
    {"id": "skills-index", "command": "bench test --check skills-index-command-adapters"},
    {"id": "root-conformance", "command": "bench test --check root-conformance"},
    {"id": "git-plumbing-owner", "command": "bench test --check git-plumbing-owner"},
    {"id": "ordinary-build-census", "command": "bench test --check ordinary-build-census"},
    {"id": "skip-ownership", "command": "bench test --check skip-ownership"}
  ]
}
```

## Testing decisions

A good test publishes through real Git and reads the result through Git.
It reads the parents with `git rev-list --parents -n 1`, the merge state through `git.AdminPath`, and the index with `git ls-files -u`.
A refusal row reads the exit code, the stderr sentence, the printed face, and an unmoved HEAD.
No test parses the code under test for its expectation.
The parents come from `git rev-parse` of the fixture commits, and the omitted paths come from the fixture writes.

These seams receive tests:

- `internal/git/checkout_test.go`: the reader over a real linked worktree. Every Git call goes through `runGit` (`internal/git/testhelpers_test.go:16-23`), because the branch-native census allows one repository or process constructor in the `internal/git` tests (`internal/conformance/ordinary_build_census_test.go:68-80`). A planted file covers each malformed and special shape. `TestIsPrimaryCheckoutSeparatesLinkedWorktree` is the precedent.
- `internal/landing/pending_merge_test.go` (new): `Owner.Land` and `Owner.Compose` over a real conflicted merge in a linked worktree. The `o.authorize` stub changes the merge state for the recheck rows, as `TestLandReconcilesAuthorizedNamedSnapshotAfterAuthorizationMutation` does. The `o.reconcile` stub drives the failure after publication, as `TestLandReportsPublishedCommitWhenReconciliationFails` does. The parent assertion mirrors `TestLandDetachedHEADHasOneExpectedParent` and `TestMergePublishesPreviousTipAsFirstParent`.
- `internal/commit/merge_parent_test.go` (new): real Git journeys through the public `Command` seam. `landingRepo` builds the bound linked worktree, and `planningCommitRepo` builds the unbound one. Each test first resets the staged `tracked.txt` that `prepareLandingCheckout` leaves. The conflicted merge mirrors `stageHandMerge` (`internal/worktree/land_refusal_fixture_test.go:253-261`), which no shared helper exports.
- `internal/commitment/repository/stale_base_test.go`: `AuthorizeCandidate` over `lightPathRepo`, an unbound assignment, a moved main, and a real `git merge --no-commit --no-ff`. The candidate is the `git write-tree` of the merge index. `TestStaleBaseCandidate` is the precedent.
- `internal/worktree/admin_readers_test.go`: `sourceMergePending` over a planted MERGE_HEAD. The worktree parallel census requires `t.Parallel()` in a new top-level test there (`internal/worktree/parallel_census_test.go:1-14`).

No system test joins the build.
The command-level journeys already cross the public seam with real Git, and the system suite caps its shared repositories at three.
A new test file in `internal/commit` or `internal/landing` adds only advisory crowding debt, because the growth ratchet grades files and no directory count.
`internal/worktree/land_journey_test.go` is over its line budget, so no ticket adds a line there.
The gate observes the feature through the package suites, and the fast lane runs the focused packages at each ticket commit.

### Seam diagram

    trigger: bench commit -m <msg> -- <path>... in an assignment worktree
        │
        ▼
    name set ──▶ [ landing.attributedPaths ── git.PendingMerge ── merged-path coverage (D1) ]
        │
        ├──▶ [ commitrepo.AuthorizeCandidate ── branchReference (D2) ── git.PendingMerge ]
        │
        ▼
    [ landing.Owner.Land: compose ── authorize ── recheck ── commit-tree -p base -p head... ── update-ref ── git merge --quit ── reconcile ]
        │
        ▼
    a published merge commit with a clean checkout, an exit 1 refusal, or an exit 3 remainder
        ◀ tests attach at git.PendingMerge, Owner.Land and Owner.Compose with real Git merges,
          the commit Command seam, AuthorizeCandidate, and sourceMergePending

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| MP01 | 11 | `git.PendingMerge` returns the merge head that `git merge --no-commit --no-ff` writes in a linked worktree | planned TestPendingMergeReadsEveryHead in internal/git/checkout_test.go | A hand-joined `.git/MERGE_HEAD` reads nothing in a linked worktree, so it answers no merge. |
| MP02 | 4 | `git.PendingMerge` returns no merge head and no error when MERGE_HEAD is absent | planned TestPendingMergeReadsEveryHead in internal/git/checkout_test.go | A reader that refuses absence refuses every ordinary commit. |
| MP03 | 5 | `git.PendingMerge` returns both merge heads of a pending octopus merge in file order | planned TestPendingMergeReadsEveryHead in internal/git/checkout_test.go | A `git rev-parse MERGE_HEAD` reader returns only the first head (observation E). |
| MP04 | 16 | `git.PendingMerge` returns the merge head of a MERGE_HEAD whose last line has no newline | planned TestPendingMergeReadsEveryHead in internal/git/checkout_test.go | A reader that requires a final newline refuses a hand-edited file. |
| MP05 | 12 | `git.PendingMerge` refuses a root outside every repository | planned TestPendingMergeRefusesUnreadableState in internal/git/checkout_test.go | A reader that reads a resolution failure as no merge lets a one-parent commit through. |
| MP06 | 12 | `git.PendingMerge` refuses a directory at MERGE_HEAD | planned TestPendingMergeRefusesUnreadableState in internal/git/checkout_test.go | A reader that reads each read error as no merge lets a one-parent commit through. |
| MP07 | 15 | `git.PendingMerge` refuses a FIFO at MERGE_HEAD within the bounded test deadline | planned TestPendingMergeRefusesUnreadableState in internal/git/checkout_test.go | A plain read of a FIFO blocks the commit. |
| MP08 | 15 | `git.PendingMerge` refuses a dangling symbolic link at MERGE_HEAD | planned TestPendingMergeRefusesUnreadableState in internal/git/checkout_test.go | A plain read reports a dangling link as absent, so the commit sees no merge. |
| MP09 | 15 | `git.PendingMerge` refuses a live symbolic link at MERGE_HEAD | planned TestPendingMergeRefusesUnreadableState in internal/git/checkout_test.go | A followed link makes bytes outside the checkout administration directory a parent. |
| MP10 | 13 | `git.PendingMerge` refuses an empty MERGE_HEAD | planned TestPendingMergeRefusesMalformedState in internal/git/checkout_test.go | A reader that reads an empty file as absent lets a one-parent commit through during a pending merge. |
| MP11 | 14 | `git.PendingMerge` refuses an abbreviated commit ID line | planned TestPendingMergeRefusesMalformedState in internal/git/checkout_test.go | A revision-resolving reader accepts a shorthand that Git never writes. |
| MP12 | 14 | `git.PendingMerge` refuses the line `main` | planned TestPendingMergeRefusesMalformedState in internal/git/checkout_test.go | A revision-resolving reader turns a branch name into a parent. |
| MP13 | 14 | `git.PendingMerge` refuses a full object ID that names no object | planned TestPendingMergeRefusesMalformedState in internal/git/checkout_test.go | A syntax-only reader gives an unknown parent to `commit-tree`. |
| MP14 | 14 | `git.PendingMerge` refuses a full object ID that names a tree | planned TestPendingMergeRefusesMalformedState in internal/git/checkout_test.go | An existence-only reader gives a non-commit parent to `commit-tree`. |
| MP15 | 17, 18 | `sourceMergePending` answers false for an empty MERGE_HEAD in a source worktree | planned TestSourceMergePendingReadsTheSharedOwner in internal/worktree/admin_readers_test.go | The current file-presence read answers true for an empty file, so the row stays red while that copy survives. |
| MP16 | 18 | `sourceMergePending` answers false for a source outside every repository | `internal/worktree/admin_readers_test.go` (`TestSourceMergePendingIsFalseWhenUndecided`) | A migration that reads a resolution failure as a pending merge fails this test. |
| MP17 | 19 | The landing conflict route names `merge --continue` for a source that holds a valid MERGE_HEAD | `internal/worktree/land_journey_test.go` (`TestLandCommandPublicConflictRepairRequiresNewReviewedTip`) | LRS6 fails when the migration loses the pending state. |
| MP18 | 18 | The landing conflict route names the commit-and-review form for a directory at MERGE_HEAD | `internal/worktree/land_journey_test.go` (`TestLandCommandPublicConflictRepairRequiresNewReviewedTip`) | LRS22 fails when an unreadable state reads as a pending merge. |
| MP19 | 1 | After a conflicted merge of an incoming commit, `bench commit` with the complete name set publishes the previous tip and the incoming commit as the two parents | planned TestCommitRecordsTheMergeParent in internal/commit/merge_parent_test.go | The current single `-p` publication records only the previous tip. |
| MP20 | 2 | `Owner.Land` publishes the expected base as parent 1 and the merge head as parent 2 | planned TestLandRecordsEachMergeHeadAsAParent in internal/landing/pending_merge_test.go | A parent list that puts the merge head first breaks every first-parent reader. |
| MP21 | 5 | `Owner.Land` publishes both merge heads of a pending octopus merge as parents 2 and 3 in file order | planned TestLandRecordsEachMergeHeadAsAParent in internal/landing/pending_merge_test.go | A first-head reader drops the second merge head. |
| MP22 | 4 | `Owner.Land` without a pending merge publishes exactly one parent | `internal/landing/landing_test.go` (`TestLandDetachedHEADHasOneExpectedParent`) | An implementation that always adds a parent fails this pinned one-parent shape. |
| MP23 | 6 | After the merge commit, `Owner.Compose` of the source onto the incoming commit reports no conflict | planned TestComposeAfterMergeCommitReadsTheMergedBase in internal/landing/pending_merge_test.go | A one-parent resolution keeps the old merge base, so the textual conflict comes back. |
| MP24 | 3, 12 | `Owner.Land` refuses a directory at MERGE_HEAD and moves no ref | planned TestLandRefusesAnUnreadableMergeState in internal/landing/pending_merge_test.go | A consumer that reads the reader refusal as no merge publishes a one-parent commit. |
| MP25 | 12 | `bench commit` prints the `commit-handback` route for an unreadable merge state | planned TestCommitRefusesAnUnreadableMergeState in internal/commit/merge_parent_test.go | A typed path error sends a reviewer's fault to the agent as a path correction. |
| MP26 | 3, 10 | `Owner.Land` refuses and moves no ref when MERGE_HEAD appears during authorization | planned TestLandRefusesAMergeStateChangeDuringAuthorization in internal/landing/pending_merge_test.go | Without the recheck, the commit publishes one parent while a merge is pending. |
| MP27 | 10 | `Owner.Land` refuses and moves no ref when MERGE_HEAD disappears during authorization | planned TestLandRefusesAMergeStateChangeDuringAuthorization in internal/landing/pending_merge_test.go | Without the recheck, the commit records a parent for an ended merge. |
| MP28 | 10 | `Owner.Land` refuses and moves no ref when MERGE_HEAD names another commit after authorization | planned TestLandRefusesAMergeStateChangeDuringAuthorization in internal/landing/pending_merge_test.go | Without the recheck, the published parent differs from the graded state. |
| MP29 | 7 | `Owner.Land` publishes a merge commit whose tree equals the base tree after an ours resolution | planned TestLandPublishesAMergeWithAnUnchangedTree in internal/landing/pending_merge_test.go | The current empty-diff refusal blocks the only Bench route to that merge parent. |
| MP30 | 8 | `Owner.Land` publishes a merge commit with an unchanged tree for a pending merge with no merged path | planned TestLandPublishesAMergeWithAnUnchangedTree in internal/landing/pending_merge_test.go | The empty true merges `8cd1e045` and `93c0e93b` needed a raw commit for this shape. |
| MP31 | 9 | `Owner.Land` without a pending merge refuses an unchanged tree with `nothing to commit` | `internal/landing/state_test.go` (`TestLandPreAuthorizationRefusalTable`) | A blanket removal of the empty-diff refusal fails the tracked-symlink row of this table. |
| MP32 | 20, 32 | `bench commit` refuses a pending merge whose name set omits a modified merged path | planned TestCommitRefusesAnIncompleteMergeNameSet in internal/commit/merge_parent_test.go | A partial merge commit makes the later landing revert the incoming change (observation B). |
| MP33 | 21 | `bench commit` refuses a pending merge whose name set omits an added merged path | planned TestCommitRefusesAnIncompleteMergeNameSet in internal/commit/merge_parent_test.go | A partial merge commit makes the later landing drop the incoming file (observation B). |
| MP34 | 24 | The incomplete-merge refusal names each of two omitted paths | planned TestCommitRefusesAnIncompleteMergeNameSet in internal/commit/merge_parent_test.go | A refusal that names only the first path costs one rerun for each omitted path. |
| MP35 | 22 | `Owner.Land` refuses a pending merge whose name set omits the conflicted path | planned TestLandRefusesAnIncompleteMergeNameSet in internal/landing/pending_merge_test.go | A merged-path list built from stage-0 index entries alone misses the unmerged path. |
| MP36 | 23 | `Owner.Land` refuses a pending merge whose name set omits a path that the merge deleted | planned TestLandRefusesAnIncompleteMergeNameSet in internal/landing/pending_merge_test.go | A merged-path list built from worktree files misses a deleted path. |
| MP37 | 25 | `Owner.Land` publishes when a named directory covers every merged path below it | planned TestLandRefusesAnIncompleteMergeNameSet in internal/landing/pending_merge_test.go | An exact-match coverage rule refuses a merged subtree that its directory names. |
| MP38 | 26 | With merged paths `pkg/b.go` and `pkgx/a.go`, `Owner.Land` refuses the name set `pkg` and names `pkgx/a.go` | planned TestLandRefusesAnIncompleteMergeNameSet in internal/landing/pending_merge_test.go | A prefix match without the separator reads `pkg` as a cover of `pkgx/a.go`. |
| MP39 | 27 | `bench commit --dry-run` refuses an incomplete merge name set with exit 1 | planned TestDryRunRefusesAnIncompleteMergeNameSet in internal/commit/merge_parent_test.go | A rule that runs only in `Owner.Land` lets the rehearsal pass a commit that refuses. |
| MP40 | 28 | In an unbound planning worktree, the incomplete-merge refusal prints in place of the commitment refusal | planned TestMergeNameSetRefusalPrecedesCommitment in internal/commit/merge_parent_test.go | A rule in `Owner.composeAuthorized` alone runs after `AuthorizeCandidate`, so the commitment refusal wins. |
| MP41 | 29 | The refusal prints an omitted path that holds a space, a glob character, an ESC byte, and a newline in Go-quoted form on one line | planned TestCommitRefusesAnIncompleteMergeNameSet in internal/commit/merge_parent_test.go | A raw path splits or corrupts the one-line refusal. |
| MP42 | 30 | `Owner.Land` publishes when the name set spells a merged path that holds a space, a glob character, a newline, and a non-ASCII byte | planned TestLandRefusesAnIncompleteMergeNameSet in internal/landing/pending_merge_test.go | A newline-framed `git diff-index` read C-quotes that name, so the correct name set refuses. |
| MP43 | 31 | `Owner.Land` without a pending merge publishes while an unnamed path stays staged | `internal/landing/landing_test.go` (`TestLandComposesOnlyNamedPathsAndCASesExpectedBase`) | A coverage rule that applies without a pending merge refuses this staged foreign path. |
| MP44 | 33 | After the merge commit, a second `bench commit` of a new change publishes one parent | planned TestCommitRecordsTheMergeParent in internal/commit/merge_parent_test.go | A commit that leaves MERGE_HEAD in place records the merge head again. |
| MP45 | 34 | After the merge commit, `git ls-files -u` prints nothing | planned TestLandEndsThePendingMerge in internal/landing/pending_merge_test.go | A reconcile that skips unmerged entries leaves stage 1 to 3 entries behind. |
| MP46 | 34 | After the merge commit, `git status --porcelain=v1` prints nothing in the worktree | planned TestCommitRecordsTheMergeParent in internal/commit/merge_parent_test.go | A partial reconcile leaves merged content staged against the new tip. |
| MP47 | 35 | After the merge commit, MERGE_HEAD is absent from the checkout administration directory | planned TestLandEndsThePendingMerge in internal/landing/pending_merge_test.go | The current publication leaves MERGE_HEAD pending, as the 2026-09-12 retro reports. |
| MP48 | 35 | After the merge commit, MERGE_MSG is absent from the checkout administration directory | planned TestLandEndsThePendingMerge in internal/landing/pending_merge_test.go | A removal of MERGE_HEAD alone leaves MERGE_MSG behind. |
| MP49 | 35 | After the merge commit, MERGE_MODE is absent from the checkout administration directory | planned TestLandEndsThePendingMerge in internal/landing/pending_merge_test.go | A removal of MERGE_HEAD alone leaves MERGE_MODE behind. |
| MP50 | 35 | After the merge commit, AUTO_MERGE is absent from the checkout administration directory | planned TestLandEndsThePendingMerge in internal/landing/pending_merge_test.go | A removal of MERGE_HEAD alone leaves AUTO_MERGE behind. |
| MP51 | 36 | The commit ends the merge through `git merge --quit` and removes no admin entry file itself | review-owned: the MP2 review reads the completion step | A hand removal writes Git-owned state that CONTEXT.md reserves for Git. |
| MP52 | 37 | `Owner.Land` returns `PublishedUnreconciledError` with the published merge commit when reconcile fails during a pending merge | planned TestLandReportsPublishedMergeWhenReconciliationFails in internal/landing/pending_merge_test.go | A merge branch that returns a plain error after publication hides the published commit. |
| MP53 | 38 | MERGE_HEAD is absent after a reconcile failure during a pending merge | planned TestLandReportsPublishedMergeWhenReconciliationFails in internal/landing/pending_merge_test.go | A completion step after the reconcile leaves the merge pending when the reconcile fails. |
| MP54 | 44 | The last stdout line of a merge commit is `committed{published_commit,path_count}` with the count of the name set | planned TestCommitRecordsTheMergeParent in internal/commit/merge_parent_test.go | A merge field in the record breaks its pinned readers. |
| MP55 | 45 | `.bench/BENCH-reference.md` states the merge contract of `bench commit` in its Command Notes | review-owned: the MP2 review reads the two sentences and the `bench anchors .bench/BENCH-reference.md` output | Without the note, a cold reader learns the name-set rule only from a refusal. |
| MP56 | 46 | The ownership fence names no roadmap path and does not name `internal/commit/commit.go` | review-owned: the fence review, which the `source-not-fenced` landing refusal enforces | A fence that admits a roadmap or grammar edit lets the delivery close FT258 or add a candidate-set flag. |
| MP57 | 39 | In a light-path worktree, `AuthorizeCandidate` admits a complete merge of main that carries a production change of main | planned TestPendingMergeCandidate in internal/commitment/repository/stale_base_test.go | The current reference counts the change of main as the branch's, so `lightPathCover` refuses it. |
| MP58 | 40 | In an unbound planning worktree, `AuthorizeCandidate` admits a complete merge of main that carries a production change of main | planned TestPendingMergeCandidate in internal/commitment/repository/stale_base_test.go | The current reference raises the unbound refusal for the change of main. |
| MP59 | 41 | In a light-path worktree, `AuthorizeCandidate` refuses a pending merge of a sibling commit that adds an uncovered production path | planned TestPendingMergeCandidate in internal/commitment/repository/stale_base_test.go | A reference that exempts every merge head admits the change of the sibling. |
| MP60 | 42 | With a directory at MERGE_HEAD, `AuthorizeCandidate` refuses a light-path candidate that lacks a later production change of main | planned TestPendingMergeCandidate in internal/commitment/repository/stale_base_test.go | A fallback to the old merge base admits that candidate on a fault. |
| MP61 | 43 | Without a pending merge, `AuthorizeCandidate` keeps every stale-base admission and refusal | `internal/commitment/repository/stale_base_test.go` (`TestStaleBaseCandidate`) | A reference change that ignores the pending condition fails these rows. |
| MP62 | 24 | The incomplete-merge refusal prints the `commit-named-path` route | planned TestCommitRefusesAnIncompleteMergeNameSet in internal/commit/merge_parent_test.go | A plain error sends the agent's own name-set fix to the reviewer. |

### Edge inventory

The profile's hostile-input checklist applies to the merged-path names and to the MERGE_HEAD file.
MP41 and MP42 cover spaces, glob characters, control bytes, a newline, and a non-ASCII byte.
The `-z` listing avoids the C-quoting that a patch header or a newline-framed listing applies.
MP02 and MP10 cover the absent-versus-empty pair, and MP04 covers a last line without a newline.
MP06 to MP09 cover special files and both symbolic link shapes before any read.
MP01 covers the relative `gitdir:` pointer of a linked worktree, because `git.AdminPath` lets Git resolve it.

No test swaps a package variable.
The recheck rows and the remainder rows replace the `authorize` and `reconcile` fields of one `Owner` value.
The refusal rows cover each fault before publication, inside authorization, and after the ref update.

- **Won't handle**: a socket or a device at MERGE_HEAD — the non-regular-file branch of `git.PendingMerge` that MP07 proves also refuses them before any read.
- **Won't handle**: a swap of MERGE_HEAD between its `os.Lstat` and its read — an assignment worktree has one writer, and the `Owner.Land` recheck still runs.
- **Won't handle**: a conflict marker in a named non-Go file — a name is the resolution, as `git add` is, and `formatNamedGoFiles` refuses a Go parse failure.
- **Won't handle**: a pending cherry-pick or revert (CHERRY_PICK_HEAD or REVERT_HEAD) — the decision covers a pending merge only, and FT254 slice 3 owns those continuation routes.
- **Won't handle**: a planted MERGE_HEAD that names a valid commit — the parent rule trusts the merge record, and the `bench worktree land` review grades the result.
- **Won't handle**: a MERGE_AUTOSTASH entry — `git merge --quit` saves it to the stash list (`git help merge`), and the reviewer settles it.
- **Won't handle**: a `git merge --quit` failure after the ref update — no seam reaches it, and `Owner.Land` still returns the exit-3 remainder.
- **Won't handle**: a change of the `composition-conflict-pending` route — `landingConflictRefusal` keeps `git merge --continue`, and FT254 slice 3 owns the resolve verb.
- **Won't handle**: a complete `bench commit` step in `handMerge` — the `merge-conflict` and `composition-conflict` faces keep the bare step, and the commit refusal names the missing paths.
- **Won't handle**: a case-variant name on a case-folding file system — the rule of `landing.attributedPaths` compares names byte for byte, so the commit refuses.
- **Won't handle**: a merged gitlink — the listing names it as one path, and the coverage rule of `landing.attributedPaths` treats it as any other path.
- **Won't handle**: a type change of a merged path — the name-only listing carries no status letter, so the coverage rule of MP32 serves it.
- **Won't handle**: a merge in the primary checkout — `bench commit` refuses the primary checkout before any merge read (`internal/commit/commit.go:116-123`).

## Ownership fences

The ticket slice derives this fence, so it equals the union of the ticket Writes lines, plus the review pickup.
It holds the owner and test files, the command registry closure, the anchor registry closure, and each canary unit that pins an edited guidance file.
No ticket changes a closure file or a canary unit.
The fence names no roadmap path and does not name the production file of the commit command (MP56).
The coordinator owns the review pickup.

- `.bench/BENCH-reference.md`
- `CHANGELOG.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/commit/commitment_test.go`
- `internal/commit/dry_run_test.go`
- `internal/commit/landing_test.go`
- `internal/commit/merge_parent_test.go`
- `internal/commit/refusal_route_test.go`
- `internal/commitment/repository/candidate.go`
- `internal/commitment/repository/stale_base_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/git/checkout.go`
- `internal/git/checkout_test.go`
- `internal/landing/attribution.go`
- `internal/landing/composition_test.go`
- `internal/landing/gitexec.go`
- `internal/landing/landing.go`
- `internal/landing/landing_helpers_test.go`
- `internal/landing/landing_reviewed_test.go`
- `internal/landing/landing_test.go`
- `internal/landing/lane_owner.go`
- `internal/landing/pending_merge.go`
- `internal/landing/pending_merge_fixture_test.go`
- `internal/landing/pending_merge_test.go`
- `internal/landing/state_test.go`
- `internal/refusalroute/faces_commit.go`
- `internal/worktree/admin_readers_test.go`
- `internal/worktree/land_journey_test.go`
- `internal/worktree/land_refusal.go`
- `internal/worktree/parallel_census_test.go`
- `reviews/commit-merge-parent.md`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/agents-handoff-section-rule`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`

## Out of scope

- The candidate-set forms `--changed`, `--dry-run --changed`, and `--ticket`, with their derivation decision: 10 edits, 3 gate runs. FT258 keeps them.
- The `bench worktree resolve` verb and its continuation routes: 8 edits, 2 gate runs. FT254 slice 3 owns them.
- A commit that completes a pending cherry-pick or revert: 6 edits, 2 gate runs.
- A conflict-marker scan of named files that are not Go: 4 edits, 1 gate run.
- The reconcile of the delivered sentences in `roadmap/FT258.md:15-21`: 1 edit, 1 gate run. The next `/bench-drain` owns it.

These estimates are planning figures, not approved priorities.

## Further notes

Evidence status: source-backed proposal, with observed Git runs and no Bench runtime claim.
Remaining uncertainty: future red-to-green evidence and the independent spec review.

### Flagged additions

Each item below goes beyond the literal decision source, and the reviewer can veto it at sign-off:

- E1: the commit ends the merge with `git merge --quit` after the ref update and before the reconcile (MP44-MP53).
- E2: a pending merge waives the `nothing to commit` refusal (MP29, MP30), and the refusal stays otherwise (MP31).
- E3: each MERGE_HEAD line becomes a parent in file order (MP03, MP21).
- E4: `Owner.Land` rechecks the merge state after authorization and refuses a change (MP26-MP28).
- E5: one reader owns the pending-merge fact, and the landing route reads a malformed state as undecided (MP15-MP18).
- E5 hardening: the reader refuses a symbolic link, a special file, an empty file, and each line that is not a full commit ID (MP05-MP14).
- E7: the record, the help, and `internal/commit/commit.go` stay unchanged (MP54, MP56), and `.bench/BENCH-reference.md` gains two sentences (MP55).
- D1 order: the D1 refusal prints before a commitment refusal (MP40), and its sentence prefix is fixed in this spec.

E6 is an exclusion, not an addition: the landing faces keep their routes (MP17, MP18, and two Won't handle lines).

### Source sentence to row

| source sentence | rows |
| --- | --- |
| `roadmap/FT258.md:15`: "When `MERGE_HEAD` exists, the sanctioned commit records the merge parent." | MP19, MP20, MP21 |
| `roadmap/FT258.md:15-16`: "It never authors a one-parent commit while merge state is pending." | MP10, MP24, MP25, MP26, MP27, MP28 |
| `roadmap/FT258.md:17`: "A conflict-resolution lane commit does not complete that merge state." | the defect before this fix; MP44 and MP47 prove the fix |
| `roadmap/FT258.md:19`: "The reviewer decided on 2026-10-10 that a pending merge records its second parent." | MP19, MP20 |
| `roadmap/FT258.md:19-20`: "The merge-parent rule ships first as its own outcome." | the scope of this spec and MP56 |
| `roadmap/FT258.md:20-21`: "The candidate-set forms stay on this row as a later feature, and their derivation decision stays open." | MP56 and Out of scope |
| `roadmap/FT258.md:3-13`: the candidate-set paragraphs | Out of scope |
| criterion `.parent` | MP19, MP20 |
| criterion `.pending` | MP24, MP26 |
| criterion `.open` | MP56 |
| D1 | MP32-MP43, MP62 |
| D2 | MP57-MP61 |

The occurrences at `roadmap/FT258.md:24-29` support the problem statement.
The parent occurrences are 2026-08-27 (commits `2a793789` and `8cd1e045`) and 2026-09-09 (commits `c7f85df3` and `7c9517d1`).
The row does not list two more occurrences: 2026-09-12 (retro commit `4717f4c7`) and 2026-09-13 (commits `39b7ceff` and `93c0e93b`).
The other occurrences concern the candidate set, which stays out of scope.

### Pre-review proof checklist

- `Cited symbols`: `landing.Owner.Land`, `Owner.DryRun`, `Owner.composeAuthorized`, `attributedPaths`, `ResolveAttributedPaths`, `CandidateTree`, `compose`, `reconcile`, `NamedPathError`, and `PublishedUnreconciledError` resolve in `internal/landing/landing.go` and `internal/landing/attribution.go`. `Owner.Compose` resolves at `internal/landing/composition.go:55`. `git.AdminPath`, `git.AdminDir`, and `git.ResolveCommit` resolve at `internal/git/worktree_admin.go:204-226` and `internal/git/git.go:220`. `Store.AuthorizeCandidate`, `branchReference`, and `lightPathCover` resolve in `internal/commitment/repository`. `sourceMergePending` and `landingConflictRefusal` resolve at `internal/worktree/land_refusal.go:82-107`.
- `Import edges`: none new. `internal/landing`, `internal/commitment/repository`, and `internal/worktree` already import `internal/git`.
- `Source-row clauses and occurrences`: the source table above quotes each clause, and the paragraph after it lists every occurrence.
- `Promised field labels`: `pending merge changed paths that the commit does not name:`, `read pending merge:`, and `pending merge changed during the commit`.
- `Changed-function callers`: `bench consumers` lists each caller below.
  - `attributedPaths`: `internal/landing/attribution.go:55`, `:60`, and `internal/landing/landing.go:127`
  - `ResolveAttributedPaths` and `CandidateTree`: `internal/commit/commit.go:136` and `:140`
  - `Owner.composeAuthorized`: `internal/landing/landing.go:152` (`Owner.DryRun`) and `:158` (`Owner.Land`)
  - `Owner.Land` and `Owner.DryRun`: `internal/commit/commit.go:187` and `:175`, plus the `Land` test callers in `internal/landing`
  - `branchReference`: `internal/commitment/repository/candidate.go:54`
  - `AuthorizeCandidate`: `internal/commit/commit.go:144`, plus five test callers in `internal/commitment`
  - `sourceMergePending`: `internal/worktree/land_refusal.go:99` and `internal/worktree/admin_readers_test.go:65`
- `Copy survival`: MP15 stays red while the file read of `sourceMergePending` survives.
- `Rendered-shape readers`: none. No current rendered shape changes, and the three new sentences have no reader.
- `Pin operators`: none.
- `Entry reads`: `git.PendingMerge` runs below `Owner.Land`, below `AuthorizeCandidate`, and below the landing verb. Real repositories grade each entry: MP19-MP54 for the landing owner, MP57-MP61 for the grading, and MP15-MP18 for the route. No entry takes an injected merge state.
- `Derived expectations`: the expected parents come from `git rev-parse` of the fixture commits. The expected omitted paths come from the fixture writes. The expected merge-state files come from observation 1. No expectation calls the code under test.
- `Consolidated rules`: the table below lists the one consolidated read.
- `Quantified obligations`: "each merge head" maps to MP03 and MP21. "Every merged path" maps to MP32-MP37. "Each omitted path" maps to MP34. "Every failure after the ref update" maps to MP52 and one Won't handle line.
- `Workflow-step writes`: none.

| site | old rule | new rule | row |
| --- | --- | --- | --- |
| `sourceMergePending` | Any readable MERGE_HEAD answers true, and an error answers false. | Only a valid merge head answers true, and a reader refusal answers false. | MP15, MP16 |
| `sourceMergePending` on a FIFO | The plain read blocks. | The reader refuses before any read. | MP07 |

### Reader inventory

These readers change their result after a two-parent worktree commit, and each change is intended or harmless:

- `internal/diff/range.go:322` and `:287`: the merge-base fallback moves to the merged main tip, and a recorded base still wins (`internal/diff/range.go:341-352`).
- `internal/diff/range.go:151-156`: a `--base` that only the merge parent reaches now passes the ancestry check.
- `internal/landing/composition.go:66`: the landing merge base advances (MP23).
- `internal/preflight/default_branch.go:20-25` and `internal/preflight/decision.go:264-286`: the `base-current` row can turn green.
- `internal/commitment/repository/candidate.go:128-137`: D2 changes the merge commit's grading, and later commits read the new merge base.
- `internal/git/status.go:155` and `:197`: the behind count drops.
- `internal/status/handoff.go:200` and `:264`: the counts past a recorded tip include the merged commits.
- `internal/diff/snapshot.go:90` and `:120`: the `bench diff` log lists the merged commits.
- `internal/git/git.go:189-202`: the `--merges` probe of `LandedInDefault` retains a branch that main does not contain, as `internal/worktree/resume_test.go:274-287` pins.
- `internal/landing/merge.go:65`, `:76`, and `internal/worktree/merge.go:367`: a later worktree merge of the same main commit reads as current.

`bench worktree merge` already publishes two-parent worktree commits (`internal/landing/merge.go:115`), so these readers meet no new shape.

These readers stay unchanged:

- `internal/worktree/landed.go:60-64` reads the first-parent chain.
- `internal/worktree/land_resume.go:246-254` and `internal/worktree/landingpolicy/landingpolicy.go:232-239` read the landing commit, not a worktree commit.
- `internal/worktree/reset_envelope.go:35` reads a reset envelope.
- `internal/worktree/land.go:274`, `:279`, `internal/reviewrecord/coverage.go:46`, `:85`, `internal/gate/authorization/authorization.go:86-104`, and `internal/worktree/reset.go:191`, `:195` keep their first-parent ancestry.
- No script, hook, or workflow file reads a commit's parents.

The pinned strings stay unchanged by design:

- The commit record: `internal/commit/landing_test.go:117-119` and `:256-260`, `CHANGELOG.md:11`, and `.agents/commands/bench-final-check.md:172-174`.
- The dry-run lines: `internal/commit/dry_run_test.go:24` and `:81`.
- The help: `internal/commit/commit_test.go:92`, `internal/commit/landing_test.go:75` and `:90`, `cmd/bench/commit_chain_test.go:150`, and `cmd/bench/command_registry_test.go:451`.
- The help row: `cmd/bench/help_inventory_test.go:144`.
- The route texts: `internal/worktree/land_surface_test.go:297` and `:356`, `internal/worktree/merge_route_test.go:156`, `internal/worktree/land_journey_test.go:304` and `:316`, and `internal/worktree/refusal_route_test.go:122`.

The guidance stays true after the fix: `.bench/BENCH-reference.md:306-311`, `docs/adr/0014-main-receives-writes-only-through-landings.md:9` and `:25`, `.bench/BENCH-reference.md:524-525`, and the comment at `internal/gitguard/verdict.go:52-55`.
No anchor pins `internal/commit/commit.go`, `internal/landing/landing.go`, or the three refusal-route files.

### Observations

The coordinator ran these observations on 2026-10-10 with Git 2.43.0.
Each run used a throwaway repository with a linked worktree under the session scratch directory.
The runs executed raw `git merge`, `git commit-tree`, `git update-ref`, and `git merge --quit` inside a script, which the destructive-git guard does not scan.
The guard refused a later top-level raw `git reset` in the same repository, and the coordinator did not retry it.

1. A conflicted merge wrote MERGE_HEAD, MERGE_MSG, MERGE_MODE, and AUTO_MERGE in the admin entry of the worktree. `git rev-parse --git-path MERGE_HEAD` named that file. `git diff-index --cached --name-status HEAD` listed `A`, `M`, and `U` paths. `git merge --quit` removed the four files, kept the index and the worktree, and exited 0.
2. Observation A: the coordinator composed the complete name set in a temp index. A two-parent `commit-tree`, the `update-ref` compare-and-swap, and a `git restore --source` of each name followed. The run left no unmerged entry and a clean status. After `git merge --quit`, MERGE_HEAD was absent. `git merge-base main HEAD MERGE_HEAD` answered the main tip, while `git merge-base HEAD main` answered the old base.
3. Observation B: a two-parent commit that named only the conflicted file left the other merged paths staged. A three-way `git merge-tree` of main and that commit reverted the modified file and dropped the added file.
4. Observation C: an ours resolution gave a tree equal to the HEAD tree, and the two-parent `commit-tree` succeeded.
5. Observation D: `git merge --quit` with no pending merge exited 0.
6. Observation E: an octopus merge wrote one MERGE_HEAD line for each head, and `git rev-parse MERGE_HEAD` answered only the first.
7. Observation F ran at the top level with only `git add` and `git diff-index`, which the guard allows. It used the observation-1 worktree, which still held the unmerged entries of `conflict.txt` after `git merge --quit`. `git diff-index --cached --name-only --no-renames -z HEAD` printed `conflict.txt` once, beside the added and the modified path. The coordinator then staged `pkg/b.go`, `pkgx/a.go`, and a name that holds a space, a newline, and `*`. The `-z` listing printed the raw newline byte, and the newline-framed listing printed that name C-quoted. The listing limited by `:(literal)pkg` printed `pkg/b.go` and not `pkgx/a.go`, and `:(literal)name*` printed nothing.

### Not verified

- Whether `git reset --hard`, which the reset plan applies (`internal/worktree/reset_apply.go:123`), removes MERGE_HEAD. The guard refused the observation. The completion order makes MP53 independent of that fact.
- The `git merge --quit` failure path, which no seam of this build reaches.
- Whether the fast lane or the gate authorization reads the merge state or the real index. The 2026-09-12 retro reports that `bench commit` graded a resolution while a merge was pending.
- The effect on `internal/worktree/clean_classes.go:161` and `:194`.
- The tests that pin the handoff counts of `internal/status/handoff.go`.
- Whether `refusalroute.New` changes the cause sentence. The D1 sentence quotes each path itself.
- The exact commands of the 2026-08-27 and 2026-09-09 occurrences. Their commit graphs are verified.
- Git behavior on a version other than 2.43.0.

### Validation plan

The build first observes MP19 red at the public `Command` seam, before the first production edit.
Each independent expectation then records one red through a named omission or a behavioral swap.

The central author probe omits the merge heads from the `commit-tree` arguments in `internal/landing/landing.go`. MP20 must fail through `bench probe`.
The coordinator probe swaps the merged listing for the covered listing in the D1 complement in `internal/landing/attribution.go`, so no merged path reads as omitted. MP32 must fail.
The two probes differ in kind and in site, and each one must report a restored subject.

Run `bench test --package ./internal/git --package ./internal/landing --package ./internal/commit` for MP1 and MP2.
Run `bench test --package ./internal/commitment/repository` for MP3, and `bench test --package ./internal/worktree` for the route rows.
Run `bench test --check root-conformance`, `bench test --check git-plumbing-owner`, and `bench test --check ordinary-build-census` at each ticket.
Run the `docs-currency-workflow` check and the `internal/anchors` package tests for the `.bench/BENCH-reference.md` and `CHANGELOG.md` changes, and the `skills-index-command-adapters` check for the reference change.
Run `bench test --check skip-ownership` when a new test can skip, such as a FIFO fixture.
The coordinator runs the whole gate only at the landing.

### Review chronology

An independent Fable high-effort review accepted the spec in iteration 1 on 2026-10-10 at commit `58295efb`.
It returned seven non-blocking findings, and this acceptance folds six of them:

- Git's literal pathspec now decides name coverage, so the coverage rule has one source with `compose`.
- The Terms section states that the D1 listing starts as the record of the merge.
- Observation F records the exact listing flags over an unmerged path and a hostile name.
- The spec states two line budgets without copied line counts.
- MP38 names its fixture paths, and the spec places D1 between the name checks and `compose`.
- The line reason leaves each effort change to `craft-line`.

The seventh finding reads `roadmap/FT258.md:17` as the reviewer's veto point, and the flagged addition E1 carries that reading.
The review recommended that each flagged addition stay, and it found no unlisted addition.

An independent Fable high-effort review accepted the ticket set in iteration 1 on 2026-10-10 at commit `24b6991b`.
It returned seven non-blocking findings, and this acceptance folds three of them:

- The unchanged surfaces now state the `CHANGELOG.md` entry that ticket 5 writes.
- The MP3 chunk of the completion plan runs the `internal/git`, `internal/landing`, and `internal/commit` package tests.
- Ticket 5 names the production path of main in its MP60 refusal criterion.

The other four findings need no change.
One is a line-budget risk in `internal/landing/pending_merge_test.go`, which an in-scope plan expansion covers.
Two correct coordinator citations that no ticket cites.
One is a fence overlap with the staged `specs/shared-test-fixtures/spec.md`, and no check compares the fences of two specs (`internal/coverage/coverage.go:189`, `internal/preflight/gather.go:56`).
