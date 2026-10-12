# Build landing window

Status: staged

Roadmap: FT342

Related roadmap: FT253. This spec settles the whole commitment outcome FT342-FT253, so it also delivers the FT253 landing lease.

Decision source: named reviewed artifact — the approved commitment outcome FT342-FT253 in `.bench/commitment.json` with its six criteria (lease, window, lifetime, override, fold, route), its sources `roadmap/FT342.md` (reviewer decisions of 2026-10-09) and `roadmap/FT253.md`, and the `/bench-write-spec FT342 —reviewer fable high` invocation of 2026-10-11

Verification log: 3 iteration(s) to accept — the fable/high spec round took 2 iterations. Iteration 1 rejected two predicates that admitted unintended callers, and iteration 2 accepted the repairs. The fable/high ticket round accepted in 1 iteration.

## Problem

A chunked build has no sanctioned route for a moved `main`.
An unrelated landing can move `main` at any time, and each move costs the build a fold, a gate run, and often a review round.
Fifteen occurrences on `roadmap/FT342.md` record this cost from 2026-09-25 to 2026-10-09.

Four surfaces fail today:

- A bare `bench preflight build` grades `base-current` red between chunks. Its route names a `main` merge that the review chain refuses (`docs/known-issues.md:14`, `internal/reviewrecord/coverage.go:86`).
- A move after the last chunk checkpoint makes `bench record completion` refuse the stale source (`internal/reviewrecord/coverage.go:127-134`). The route is a refreeze, a full re-record, and another review round.
- A completion landing over a moved destination refuses with `land-infrastructure`, whose route names `bench doctor` (`internal/gate/completion.go:71`, `internal/gate/authorization/authorization.go:164`, `internal/refusalroute/faces_land.go:116-119`).
- The final checks need not include the package of a test file that the fold brings in. The ft336 build found such a red only in its third Spec round.

Two landings can also run at once.
The gate lock covers only the gate phase (`internal/gate/run_transaction.go:76-84`).
A second landing meets that lock as an infrastructure outcome that names `bench doctor` and a bare PID.

## Solution

`bench worktree land` holds one lease on its destination in the intent ledger, from composition through publish.
A second landing refuses before composition and names the owner's assignment id.

While a chunked build is open, the landing refuses every landing that is not the build's own, unless the build is in a chunk gap.
`bench commitment start` opens the build.
Completion, `bench commitment block`, and the release of every bound assignment close it.
A green chunk checkpoint opens a gap, and the next chunk's first dispatch closes it after the build folds `main`.

A reviewer can grant one landing a recorded override with a reason, from a terminal in the primary checkout.
No agent route grants it.

`bench record fold` records a fold of `main` into the source.
A clean fold needs no review, and only the resolution delta gets review.
The review chain accepts a recorded fold between chunks and after the last chunk.
The final checks include the package of each test file that a fold brings in.

A refusal for a moved destination names `bench worktree merge` and the landing rerun, never `bench doctor`.

## Reviewer decisions

Each decision below is the author's call inside the 2026-10-09 bounds. Each one waits for reviewer veto.

- R1 (scope): the spec covers all six criteria of the outcome FT342-FT253, because the commitment binds both rows into one outcome. FT253's lease is chunk BW-A.
- R2 (window): the open build is derived, not stored. An open build is an outcome with a claim, with a bound staged-spec deliverable, and with at least one active bound assignment. Abandon is the state where no bound assignment stays active. A tickets-only deliverable opens no window. The build's own landing lands the bound deliverable from a bound assignment. A spec-less landing from a bound worktree waits for a gap.
- R3 (window start): the window opens at `bench commitment start`, as the decision states. A landing between the start and the first chunk checkpoint therefore waits.
- R4 (gaps): the gap state lives in the ledger. The green `bench gate --checkpoint <spec> --chunk <id>` verb opens it. The first `bench record assignment` of a chunk after the recorded chunk closes it. That dispatch refuses while the default branch tip is not an ancestor of the source. A repair dispatch to an earlier chunk leaves the gap record unchanged.
- R5 (override trust root): the override verb needs the primary checkout, a terminal on stdin, and the typed assignment id. The follow-on guard denies the verb in an agent tool call. The trust assumption is reviewer-visible: a fabricated pseudo-terminal or a direct ledger edit defeats these honest-mistake layers.
- R6 (fold review): a non-empty resolution delta takes the same three independent axes as a chunk, scoped to the resolution paths.
- R7 (final checks): each merged test package becomes a derived final verification requirement, `merged-test:<package>`, with the command `bench test --package <package>`.
- R8 (fence): the review preflight authorizes a path that a recorded fold lists in its resolution. This is a flagged addition. Without it, a resolved `main` path needs a fence expansion, as the FT337 occurrence records. The skip holds only while the path keeps its bytes and mode from the fold commit.
- R9 (fold red): `merge-fold-red` keeps reviewer authority. Its route becomes the hand merge with the repair, `bench commit`, and `bench record fold`.
- R10 (proof fault): a completion proof fault over a source that already holds the destination routes to `land-handback`, not to `bench doctor`.
- R11 (records): the lease, gap, and override records join `CommitmentState`, so the `DATA_HANDLING.md` inventory check covers them.

## User stories

Line: opus / high.
Implementation-line reason: BW-D is the hardest chunk, because three verbs meet under one ledger lock with two race windows. Spec precision is high, seam uncertainty is moderate, and the leverage override applies to the BW-F guidance. The FIFO gate barrier and the real-process race test give strong tests.
Harder chunks: BW-C, BW-D, BW-E.

### Landing lease

1. As an operator, I want a landing to hold a destination lease from composition through publish, so that two landings never compose at once.
2. As an operator, I want a second landing refused before composition with the owner's assignment id, so that I know whom to wait for.
3. As an operator, I want that refusal to route me to wait and rerun, so that one paste resumes the landing.
4. As an operator, I want the lease released in the publish decision, so that the next landing can start at once.
5. As an operator, I want the lease released when the landing refuses after composition, so that a refusal never blocks the destination.
6. As an operator, I want the next landing to take over a dead owner's lease, so that a killed landing never freezes the destination.
7. As an operator, I want a landing resume to take no lease, so that recovery of a published landing never waits.
8. As an operator, I want a concurrent landing refused by the lease, not the gate lock, so that it never names `bench doctor`.
9. As a maintainer, I want the lease in the intent ledger under its shared lock, so that every checkout sees one lease state.
10. As a maintainer, I want the ledger to refuse two leases on one destination, so that a hand edit cannot create a second owner.
11. As a maintainer, I want `DATA_HANDLING.md` to list each new ledger field, so that the record inventory stays complete.
12. As a maintainer, I want a tolerant ledger write to keep the new records, so that a recovery pass never drops them.
13. As an operator, I want a lease lock fault refused before composition, so that no landing runs without a lease.

### Open-build window

14. As an orchestrator, I want an unrelated landing refused while my chunked build is open, so that `main` never moves under a chunk.
15. As an operator, I want the window refusal to name the build and its next gap, so that I know when to rerun.
16. As an orchestrator, I want my build's own landing admitted through its window, so that the completion landing publishes.
17. As an orchestrator, I want the window to open at `bench commitment start`, so that `main` holds still from the first chunk.
18. As an orchestrator, I want the window to close when the completion landing publishes, so that later landings flow.
19. As an orchestrator, I want `bench commitment block` to close the window, so that a stopped build cannot freeze `main`.
20. As an orchestrator, I want the window closed when no bound assignment stays active, so that an abandoned build cannot freeze `main`.
21. As a reviewer, I want a tickets-only deliverable to open no window, so that a one-ticket outcome does not hold `main`.
22. As an operator, I want the window checked again at publish under the ledger lock, so that a newly opened build still counts.
23. As an orchestrator, I want my own light-path landings refused outside a gap, so that they wait like any unrelated landing.
24. As an operator, I want a landing refused when another open build is outside its gap, so that parallel builds both hold `main`.
25. As an operator, I want the window refusal to win over the lease refusal, so that a lease wait never ends at the window.

### Chunk gaps

26. As an orchestrator, I want a green chunk checkpoint with a successor chunk to open a gap, so that waiting landings land between chunks.
27. As an orchestrator, I want the last chunk to open no gap, so that a move after the last freeze needs the override.
28. As an orchestrator, I want a red checkpoint to open no gap, so that only an accepted chunk lets `main` move.
29. As an operator, I want an unrelated landing admitted while the build's gap is open, so that the gap does its job.
30. As an orchestrator, I want the next chunk's first dispatch to close the gap, so that `main` holds still during the chunk.
31. As an orchestrator, I want that dispatch refused while `main` is ahead of the source, so that the charge pins a folded base.
32. As an orchestrator, I want a same-chunk or earlier-chunk dispatch to keep the gap state, so that its tickets dispatch freely.
33. As an operator, I want a landing refused at publish when its gap closed during its gate, so that the chunk keeps its base.
34. As an orchestrator, I want a checkpoint in an unbound worktree to write no gap, so that only a bound build opens one.

### Reviewer override

35. As a reviewer, I want to grant one landing a recorded override with a reason, so that I can move `main` after the freeze.
36. As a reviewer, I want the override to admit that one landing through every open window, so that a back-out lands.
37. As a reviewer, I want the override consumed when its landing publishes, so that it admits no second publication.
38. As a reviewer, I want the override kept when its landing refuses before publish, so that the rerun still lands.
39. As a reviewer, I want the override verb refused without a terminal on stdin, so that an agent tool call cannot grant it.
40. As a reviewer, I want the verb to require the typed assignment id, so that a null or scripted stdin cannot grant it.
41. As a reviewer, I want the verb refused outside the primary checkout, so that the landed build of Bench writes the record.
42. As a reviewer, I want the follow-on guard to deny the verb in an agent tool call, so that an honest mistake stops early.
43. As a reviewer, I want an override for an inactive assignment refused, so that a typo grants nothing.
44. As a reviewer, I want an empty reason or a reason with a control byte refused, so that the record stays one readable line.
45. As a reviewer, I want the back-out documented as a revert landed from a fresh worktree with the override, so that nobody resets `main`.

### Fold record

46. As an author, I want `bench record fold` to record a fold of `main`, so that the chain accepts a fold between or after chunks.
47. As an author, I want a clean fold to need no review, so that a fold of unrelated changes costs no round.
48. As an author, I want the fold record to list the paths that differ from the automatic merge, so that only they get review.
49. As a reviewer, I want a non-empty resolution to need independent reviews bound to the fold, so that a hand repair gets review.
50. As an author, I want a recorded fold to serve as the next chunk's base, so that the chain accepts a fold between chunks.
51. As an author, I want fold records to cover a completion source past the last chunk tip, so that a late fold needs no refreeze.
52. As an author, I want a chain refusal over an unrecorded fold to name `bench record fold`, so that it names its repair.
53. As an author, I want any commit except a default-branch fold over the covered source refused, so that no other commit passes.
54. As an author, I want the final checks to include each merged test file's package, so that a merged test red never hides.
55. As a maintainer, I want the testreport package index to derive the merged packages, so that the fold and `--changed` agree.
56. As an author, I want the review preflight to authorize a recorded resolution path, so that it needs no fence expansion.
57. As an author, I want a red fold routed to a hand merge, `bench commit`, and the fold record, so that the face names a route.
58. As an author, I want a bare `bench preflight build` between chunks to name a fold that the chain accepts, so that the known issue retires.

### Moved destination

59. As an operator, I want a moved destination at completion routed to `bench worktree merge` and the rerun, so that I never run `bench doctor`.
60. As an operator, I want a ref move during my gate routed the same way, so that one face names one repair.
61. As an operator, I want a compare-and-swap refusal at publish routed the same way, so that the last race names the fold.
62. As an operator, I want a checkout change without a ref move to keep the handback, so that no fold is prescribed for it.
63. As an operator, I want a proof fault over a current source routed to the reviewer, so that it never reaches `bench doctor`.
64. As a maintainer, I want the proof cause to reach the landing as a type, so that no reason text decides the route.

### Guidance

65. As an orchestrator, I want the build guidance to state the gap fold and the fold record, so that I follow the enforced route.
66. As a reviewer, I want the review guidance to state the fold record and the resolution review, so that I know what to review.
67. As a reader, I want the reference guide to state the lease, window, gap, override, and back-out, so that one place names them.
68. As a reader, I want the glossary to define the new terms with Avoid lists, so that sessions use one word for each thing.
69. As a maintainer, I want the base-current known-issue row removed and a changelog entry added, so that the docs state current behavior.

### Reviewed exclusions

70. As a reviewer, I want a `--wait` queue left out, so that the landing keeps one refusal and one rerun.
71. As a reviewer, I want a non-landing gate run to keep its lock refusal, so that this spec changes only landing paths.
72. As a reviewer, I want reviewer identity left to the honest-mistake layer, so that this spec adds no signing key.
73. As a reviewer, I want a pty-backed terminal test left to a dependency decision, so that the dependency standard holds.

## Implementation decisions

### Records

`CommitmentState` (`internal/intent/ledger/commitment.go:11`) gains three lists, each `omitempty`.

- `leases`: one `LandingLease{Destination, Assignment}` for each leased destination ref.
- `gaps`: one `ChunkGap{Milestone, Outcome, Chunk, Open, Last}` for each outcome that has dispatched or checkpointed a chunk.
- `overrides`: one `LandingOverride{Assignment, Reason}` for each assignment that holds an unconsumed override.

`ValidateCommitment` refuses each of these records:

- a second lease on one destination
- a second gap for one outcome
- a second override for one assignment
- an invalid assignment identity
- an empty reason or a reason with a control byte

`readPathTolerant` (`internal/intent/transaction.go:89-135`) copies the three lists, so a tolerant write keeps them.
`internal/intent/ledger_aliases.go` re-exports the three types.
`DATA_HANDLING.md` lists each leaf field between the commitment record markers.

The strict ledger decode uses `DisallowUnknownFields` (`internal/jsonfile/decode.go:51`).
An older binary therefore refuses a ledger that holds one of the new records.
The Edge inventory states this residual.

### Landing lease

The lease owner is `internal/commitment/repository`, beside `PublishAdmitted`, because the publication admission already decides under the intent lock.
`landAttributed` acquires the lease after the admission proofs and before `j.landReviewed`, which composes first (`internal/worktree/land.go:200-204`).

Liveness is a kernel record lock, as the gate lock does (`internal/gate/run_transaction.go:231-252`).
The lock file is `<git-common-dir>/bench-landing-<first 16 hex of sha256(destination ref)>.lock`, created at mode 0600, and never removed.
An in-process holder map answers for the calling process, because POSIX record locks belong to a process (`internal/gocache/lock.go:33-40`).
The lease adds no kill-0 probe.

The acquire decision runs in one `intent.Transact` with a compensation:

1. Decide the window for the publication (see Window). A window refusal ends the decision.
2. Take `gocache.RecordLock(F_WRLCK)` with `F_SETLK` on the lock file. If another holder has it, refuse with the lease face, and name the assignment of the ledger's lease record for that destination.
3. Write the lease record. A record whose lock had no holder belonged to a dead owner, and the new record replaces it.

The compensation releases the record lock when the ledger write fails.
A lock file that cannot be opened or locked for another reason refuses with `landing lease unavailable: <reason>` under `land-infrastructure`.

The release has two sites.
`PublishAdmitted` removes the lease record in the same decision that publishes, and it then reports the change so the transaction writes it.
The pre-publish error branch at `internal/worktree/land.go:205-211` removes the record in its own transaction.
Each site then unlocks and closes the lock file.
A kill skips both sites, and the kernel drops the record lock, so the next acquire takes the record over.
A ledger write that fails after the ref update leaves the lease record, and the same takeover clears it.

`--resume` takes no lease, because it never composes or publishes.

The lease refusal prints `refused{detail=landing destination refs/heads/<branch> is leased by assignment <id>,next=...}`.
Its face is `land-lease-held` with agent authority.
Its route is the instruction `wait until the landing of assignment <id> ends`, then the caller's rerun.

### Window

`internal/commitment/repository` derives the open builds from one ledger read and the policy.
An open build is a claim `(milestone, outcome)` with at least one qualifying binding of that outcome.
A qualifying binding names a staged-spec deliverable and an assignment that is active in `ledger.Assignments`.
The staged-spec test reuses the deliverable classification `validateDeliverable` (`internal/commitment/repository/sources.go:90`), which `bench commitment start` calls at `internal/commitment/repository/admission.go:41`.
A fault that resolves the default-branch tip refuses the landing with `land-infrastructure` and the cause.

The build's assignments are the active bound assignments of that outcome.
A sibling that `bench worktree create --from` binds is one of them (`internal/commitment/repository/readiness.go:121-150`).
The build's own landing is a publication from one of the build's assignments whose deliverable equals the build's bound deliverable.
A spec-less or tickets-only landing from a bound assignment is not the build's own landing.

A publication passes the window when, for each open build, one of these holds:

- The publication is the build's own landing.
- The build's gap record has `Open` true.
- The ledger holds an override for the publication assignment.

The landing decides the window in the acquire decision and again in `PublishAdmitted`.
The publish decision is the authoritative one, because it holds the intent lock across the ref update.

The window refusal prints `refused{detail=landing destination refs/heads/<branch> is held by open build <outcome> (<deliverable>); wait for <gap>,next=...}`.
The `<gap>` text derives from the build's gap record:

- No record: `its first chunk checkpoint`.
- A closed record with `Last` false: `the checkpoint of chunk <chunk>`.
- A closed record with `Last` true: `its completion landing`.

When several open builds refuse, the refusal names the first in ledger claim order.
Its face is `land-build-window` with agent authority.
Its route is the instruction `wait for the gap that the refusal names, or ask the reviewer for an override`, then the caller's rerun.

`ReconcileDelivered` and `SetBlocker` drop the gap record of the outcome whose claim they drop.
`bench commitment start` drops a stale gap record of the outcome that it claims.

### Gaps

The gate verb opens a gap.
After `RunCommand` (`internal/gate/gate.go:185-210`) returns a green result for `--checkpoint <spec> --chunk <id>`, it reads the plan at HEAD.
The verb writes the gap record `{Chunk: <id>, Open: true, Last: false}` when two conditions hold.
The plan names a chunk after `<id>`, and the active assignment of the root holds a binding of an open build.
Otherwise it writes nothing.

The verb also writes nothing when the gap record names a chunk after `<id>` in plan order.
So a repair checkpoint of an earlier chunk keeps the recorded chunk and its `Last` value.
A write failure prints `chunk gap not recorded: <cause>` and exits 1, and a rerun reuses the green verdict and writes the record.

`bench record assignment` closes a gap.
The record command (`internal/reviewrecord/recordcmd/assignment.go:40`) resolves the dispatched ticket's chunk from the version 2 plan at HEAD.
When the outcome has no gap record, or the dispatched chunk follows the record's chunk in plan order, the dispatch is the chunk's first dispatch.
A first dispatch runs one ledger decision before the plan write:

1. Refuse when the default branch tip is not an ancestor of HEAD.
2. Write `{Chunk: <dispatched chunk>, Open: false, Last: <the chunk is the plan's last>}`.

The refusal prints `dispatch refused: the default branch tip <commit> is not an ancestor of the source; fold it with bench worktree merge --from <branch> <assignment-id> and record it with bench record fold <slug> --source <fold-commit> before the dispatch`, with exit 1.
A dispatch to the record's chunk or to an earlier chunk writes no ledger change.
A repair session for a ticket of an earlier chunk is such a dispatch, so it keeps the recorded chunk and its `Last` value.
A dispatch from a worktree without a binding of an open build writes no ledger change.
If the plan write fails after the ledger write, a rerun finds the chunk current and writes the plan.

### Override

`bench commitment override --assignment <id> --reason <text>` is a new commitment form.
The reason follows the one rule that `bench commitment block` applies to its reason.
Each refusal renders through the commitment form's refusal renderer (`internal/commitment/commitcmd/command.go:362`), so it prints `error: bench commitment override refused — <reason>`.
The verb refuses in this order, each with exit 1 and with the reasons below:

1. A root that is not the primary checkout: `error: commitment override runs from the primary checkout`.
2. A stdin that `terminal.IsTerminal` rejects: `error: commitment override needs the reviewer at a terminal; run it in your own terminal`.
3. An assignment that is not active: `error: commitment override assignment <id> is not active`.
4. A confirmation line that differs from `<id>`, or EOF: `error: commitment override not confirmed`.

Before the confirmation read, the verb prints one `override` row with the assignment, its label, the refusing open builds, and the reason.
It then prints the prompt `type the assignment id to confirm:`.
On success it writes the override record in one ledger decision and prints `overridden{assignment=<id>}`.

`PublishAdmitted` removes the override record of the publication assignment in the publish decision.
A landing that refuses before publish leaves the record.

The follow-on guard (`internal/benchguard`) is wired on Claude Code and Codex.
It denies a Bash tool call whose Bench segment runs `commitment override`, also as the child of `bench worktree exec`.
The guard manifest line in `.bench/hooks/block-bench-follow-on.sh` names the reviewer verb.

### Trust chain for the override

The reviewer's own terminal shell runs the `bench` wrapper in the primary checkout.
The wrapper resolves the primary checkout's `dist/bench`, which the landing broker refresh publishes from landed sources (`projects/benchkit.md:638-659`).
That executable checks the checkout, the terminal, the assignment, and the typed confirmation before its one ledger write.

An agent tool call meets two layers first.
The guard denies the verb on Claude Code and Codex, and it fails open when the core is missing.
The Claude Code Bash tool has no terminal on stdin, as this session observed; the Codex shell's terminal state is unknown.

The trust assumption is reviewer-visible.
A fabricated pseudo-terminal, a direct ledger edit, or a hand-built executable defeats these layers.
The kit's hooks declare an honest-mistake threat model (`.bench/hooks/check-agent-line.sh:12`), and this spec keeps that model.

### Fold record

`bench record fold <slug> --source <commit>` is a new record form in `internal/reviewrecord/recordcmd`.
`Record` gains `folds`, a list of `Fold{ID, After, From, Commit, Incoming, SourceDigest, Resolution, ResolutionDigest, TestPackages, Reviews}`.
`ID` counts `F1`, `F2`, and so on.
`After` names the chunk that the fold follows.

The form refuses, each before any write:

- A source that is not HEAD: `fold source <commit> is not HEAD; record the fold before the next commit`.
- A commit without exactly two parents: `fold <commit> is not a two-parent merge`.
- A second parent that is not an ancestor of the default branch tip: `fold <commit> does not merge a default-branch commit`.
- A first parent whose source digest is neither the covered digest nor a comment-only gap from it: `fold <commit> does not follow the reviewed source`.

The covered digest is the source digest of the last chunk entry, or of the last fold after it.

The automatic tree is the landing composition of the two parents (`internal/landing/composition.go:54-90`), which `bench worktree merge` publishes.
For a conflicted composition, the automatic tree is the tree that `git merge-tree --write-tree` writes with its conflict markers, through the same owner.
`Resolution` lists, in byte order, each path where the fold commit's tree differs from the automatic tree, without the review record path.
`ResolutionDigest` is the digest of that diff.
A fold that `bench worktree merge` published has an empty resolution.

`TestPackages` is the sorted package list that `testreport.WritesFloor` derives from the `_test.go` paths that differ between `From` and `Commit`.
That function is the one path-to-package owner (`internal/testreport/floor.go:24`).
`internal/reviewrecord` cannot import `internal/testreport`, because `internal/testreport` already depends on it; the record command computes the list and stores it.

`bench record review` takes `--fold <id>` in place of `--chunk <id>`.
A fold review copies the fold's `From` and `Commit` as its base and tip, and the three axes follow the chunk rules of `checkReviews`.
A fold with an empty resolution needs no review.

### Review chain

`checkSource` (`internal/reviewrecord/coverage.go:19`) changes in three places.

1. A later chunk base passes when the last fold after the predecessor covers it. The base digest equals that fold's `SourceDigest`, the fold commit is an ancestor of the base, and the fold's `From` covers the predecessor's digest. So the record commit after the fold serves as the base.
2. A source past the last chunk tip passes when fold records cover it. The requested digest equals the last fold's digest, or a comment-only gap separates them.
3. A fold with a non-empty resolution needs three passing axes, as `checkReviews` grades a chunk.

The chain-gap refusal becomes a repair of form `fold`.
Its text becomes `chunk <id>: stale review chain gap: expected base <tip>, the chunk <prev> tip, got <base>; plan commits land before the ticket merge, a default-branch merge lands before the first chunk or as a recorded fold, and only record commits and comment-only corrections follow a chunk tip; record the fold or review the uncovered delta`.
The stale-source refusal becomes a repair of form `fold` when the uncovered range holds a fold of a default-branch commit.
Its text keeps `stale reviewed source: no chunk review covers`.

### Final checks

`checkCompletion` (`internal/reviewrecord/check.go:71-84`) adds one derived requirement for each package in any fold's `TestPackages`.
The requirement id is `merged-test:<package>`, and its command is `bench test --package <package>`.
`bench record verification --final` accepts these ids.
The completion requires a passing result for each derived requirement at the completion source, as it requires for the planned final list.

### Review preflight

`unauthorizedPaths` (`internal/preflight/paths_authorized.go:24`) also skips a path that a fold entry of the record at the source tip lists in `Resolution`.
The fold review grades that path, so the fence does not need to name it.
The skip holds only while the path at the source tip keeps its bytes and mode from the fold commit.
The fold review graded only those bytes.

### Fold red

The `merge-fold-red` face keeps reviewer authority.
Its route becomes four steps: the existing first instruction, the hand merge `<checkout-git> merge <from-commit> (repair the red in the merge commit)`, `bench commit`, and `bench record fold`.
The route no longer names FT342.

### Moved destination

A new land face `land-destination-moved` has agent authority and the sentence `landing destination moved past the reviewed source`.
Its route has three steps: `bench worktree merge --from <branch> <label>`, the instruction `for a spec landing, record the fold, its review, the merged-test verification, and the completion at the new tip`, and the caller's rerun at the repaired tip.

Three refusals take this face:

- The post-gate recheck at `internal/landing/landing.go:283-285`, when the destination ref no longer equals the destination base. The detail text stays.
- A `TipMovedError` from the compare-and-swap at publish (`internal/landing/gitexec.go:22-31`).
- An authorization result with the completion proof cause, when the destination base is not an ancestor of the source tip.

A post-gate checkout change with an unmoved ref keeps `land-handback`.
A completion proof cause over a source that holds the destination base routes to `land-handback`.

The gate carries the proof cause as a type.
`funnelFace` already maps `completionProofError` to its face (`internal/gate/checkpoint.go:152-153`).
Outside a verb context, the gate result records that face, and `authorization.Result` copies it into a typed cause that the landing reads.
No reason text decides the route.

### Guidance

The `Land` section of `.agents/commands/bench-implement-spec.md` changes its fold rule.
A build folds `main` before the first chunk and in a gap before the next chunk's first dispatch.
After the last freeze, `main` moves only through the override.
Each fold after the first chunk takes `bench record fold`.
The anchors at `internal/anchors/registry_chunk_chain.go:30-36` and `:56` move to the new sentences.

`.agents/commands/bench-review-implementation.md:99`, `:217`, and `:244` state the fold record, the resolution review, and the landing base.
`.agents/commands/bench-final-check.md:24` names the moved-destination face.

`.bench/BENCH-reference.md` Command Notes state the lease, the window, the gap, and the override.
They also name `bench commitment block` as the reviewer's stop for a stalled build that still holds the window.
They also state the back-out route: `git revert --no-commit <commit>` in a fresh worktree, `bench commit`, and a landing with the override.
`CONTEXT.md` defines **landing lease**, **open build**, **chunk gap**, **fold record**, and **landing override**.
`docs/known-issues.md:14` leaves, and `CHANGELOG.md` gains one entry.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| BW-A / 1-hold-the-landing-lease.md | One lease per destination from composition through publish, with liveness and takeover | BW1-BW16 | planned lease, land-lease, ledger, and system race tests; recovery and registry tests | no |
| BW-B / 2-route-a-moved-destination.md | One moved-destination face with a typed proof cause | BW17-BW24 | planned moved-destination and authorization tests; the system race test | no |
| BW-C / 3-record-a-fold.md, 4-review-a-fold-resolution.md | Fold record, chain acceptance, resolution review, merged-test final checks, fold red route | BW25-BW45, BW93, BW98 | planned fold, record-fold, preflight, and merge-route tests; the chain-gap test; help inventory | yes |
| BW-D / 5-hold-main-for-an-open-build.md, 6-open-and-close-chunk-gaps.md | Open-build window with its lifetime and its gaps | BW46-BW70, BW94-BW97, BW99, BW100 | planned window, land-window, chunk-gap, and dispatch-gap tests | yes |
| BW-E / 7-grant-a-reviewer-override.md | Terminal-bound reviewer override, consumption, and guard | BW71-BW84 | planned override, land-override, and guard tests; help inventory | yes |
| BW-F / 8-document-the-landing-window.md | Guidance, glossary, known issues, and changelog | BW85-BW92 | review-owned reads, anchors, prose | no |

BW-A comes first, because BW-D decides the window inside the lease's acquire decision.
BW-B follows, because it needs no window and it stands alone.
BW-C comes before BW-D, because the dispatch refusal of BW-D names `bench record fold`.
BW-E consumes the window of BW-D.
BW-F comes last, because its guidance names every form.

### Completion plan

Run line, by reviewer direction on 2026-10-11: opus / high for each ticket author, opus / high for each review axis, and fable / high for each consultation. The reviewer gave `--auto-approve`, so the orchestrator decides each fence, spec, and repair question and flags it.

This version 2 completion plan declares the delegated run and one verification inventory for each ticket.
Each probe field names the mutation kind, its site, and the test that turns red, and the author derives the exact mutated text.

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"build-landing-window-build-20261011","orchestrator_session":"claude:blw-orchestrator-20261011","author_limit":1,"assignments":{"1-hold-the-landing-lease.md":[{"session":"claude:blw_t1","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"4206c20488bce8b4711d7760c53202bebe392a86","native_ref":"claude-agent:blw_t1"},{"session":"claude:blw_t1_rep1","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"fd00c07f85fb466ca6560d7e7bfc81903e7c2890","native_ref":"claude-agent:blw_t1_rep1","predecessor":"claude:blw_t1","trigger":"user-directed","stopped":"claude:blw_t1 ended after its ticket commit ced1dedc and its BW-A verification entries in record commit 3401d028","preserved":"no uncommitted work; the ticket commit ced1dedc and the record commit 3401d028 on the FT342-FT253 integration branch"}],"2-route-a-moved-destination.md":[{"session":"claude:blw_t2","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"e002c6cab0f376a127f18afe4509535a51a090fc","native_ref":"claude-agent:blw_t2"},{"session":"claude:blw_t2_rep1","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"d7656bcd2a56664c47733ed5d17c098770e284a8","native_ref":"claude-agent:blw_t2_rep1","predecessor":"claude:blw_t2","trigger":"user-directed","stopped":"claude:blw_t2 ended after its ticket commit 6c90a3c6 and its BW-B verification entries in record commit 0ceb0a3d","preserved":"no uncommitted work; the ticket commit 6c90a3c6 and the record commit 0ceb0a3d on the FT342-FT253 integration branch"}],"3-record-a-fold.md":[{"session":"claude:blw_t3","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"64ed906be5a2dec8df48a850d6cd532c8b0f2150","native_ref":"claude-agent:blw_t3"},{"session":"claude:blw_t3_rep1","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"3032523a32f5c2367101fdd14753b2bd27b142b2","native_ref":"claude-agent:blw_t3_rep1","predecessor":"claude:blw_t3","trigger":"user-directed","stopped":"claude:blw_t3 ended after its ticket commit 12a2d5e1 and its BW-C verification entries in record commit 96f67c82","preserved":"no uncommitted work; the ticket commit 12a2d5e1 and the record commit 96f67c82 on the FT342-FT253 integration branch"},{"session":"claude:blw_t3_rep2","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"fa4a6ee7975ec272730bf55288829f0b51602762","native_ref":"claude-agent:blw_t3_rep2","predecessor":"claude:blw_t3_rep1","trigger":"user-directed","stopped":"claude:blw_t3_rep1 ended after its repair commit eae5349c and its BW-C verification entries in record commit 0bea3e83","preserved":"no uncommitted work; the repair commit eae5349c and the record commit 0bea3e83 on the FT342-FT253 integration branch"}],"4-review-a-fold-resolution.md":[{"session":"claude:blw_t4","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"12a2d5e16090abb7a5aa80febbd69548ea631662","native_ref":"claude-agent:blw_t4"},{"session":"claude:blw_t4_rep1","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"450e75d8523645ef8488cbddf7bdd3d5c36ef359","native_ref":"claude-agent:blw_t4_rep1","predecessor":"claude:blw_t4","trigger":"user-directed","stopped":"claude:blw_t4 ended after its ticket commit db1bdc5d and its BW-C freeze and verification entries in record commit 3dd6beda","preserved":"no uncommitted work; the ticket commit db1bdc5d and the record commit 3dd6beda on the FT342-FT253 integration branch"}],"5-hold-main-for-an-open-build.md":[{"session":"claude:blw_t5","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"c470114448c428d4d7c3a2cf918436d3b3c367b3","native_ref":"claude-agent:blw_t5"},{"session":"claude:blw_t5_rep1","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"d5b5c1ae81fc92162eb843b050c074a66be89855","native_ref":"claude-agent:blw_t5_rep1","predecessor":"claude:blw_t5","trigger":"user-directed","stopped":"claude:blw_t5 ended after its ticket commit bd34ab9d and its BW-D verification entries in record commit debd4735","preserved":"no uncommitted work; the ticket commit bd34ab9d and the record commit debd4735 on the FT342-FT253 integration branch"}],"6-open-and-close-chunk-gaps.md":[{"session":"claude:blw_t6","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"e50fb99d324226ed149b37429175b1fe21b1c5f7","native_ref":"claude-agent:blw_t6"},{"session":"claude:blw_t6_rep1","assignment":"5e861f162efb8534058c864f8d4b9594","model":"opus","effort":"high","source":"da8a1550ee061b517fd53f1618e3675e2dbec75f","native_ref":"claude-agent:blw_t6_rep1","predecessor":"claude:blw_t6","trigger":"user-directed","stopped":"claude:blw_t6 ended after its ticket commit b68ee13d and its BW-D freeze and verification entries in record commit e09c7f3f","preserved":"no uncommitted work; the ticket commit b68ee13d and the record commit e09c7f3f on the FT342-FT253 integration branch"}],"7-grant-a-reviewer-override.md":[],"8-document-the-landing-window.md":[]}},"chunks":[{"id":"BW-A","tickets":["1-hold-the-landing-lease.md"],"verification":[{"id":"t1-land-focused","command":"bench test --package ./internal/worktree --run 'TestLandHoldsTheDestinationLease|TestSecondLandingMeetsTheLease|TestLandReleasesTheLeaseOnRefusal|TestLandResumeTakesNoLease|TestLeaseLockUnavailable|TestLandingRefusalRegistryHasAProducingFixture'","probe":"Omit the lease acquire call that the landing makes after its admission proofs and before landReviewed, so that a second landing reaches composition and the gate lock. TestSecondLandingMeetsTheLease turns red.","ticket":"1-hold-the-landing-lease.md"},{"id":"t1-lease-focused","command":"bench test --package ./internal/commitment/repository --run 'TestPublishReleasesTheLeaseInItsDecision|TestLeaseTakesOverADeadOwner'","ticket":"1-hold-the-landing-lease.md"},{"id":"t1-records-focused","command":"bench test --package ./internal/intent --run 'TestValidateCommitmentRefusesDuplicateBuildRecords|TestTolerantReadKeepsTheBuildRecords|TestIntentReExportsEveryMovedLedgerName'","ticket":"1-hold-the-landing-lease.md"},{"id":"t1-inventory-focused","command":"bench test --package ./internal/conformance --run 'TestCommitmentRecordInventoryNamesTheBuildRecords'","ticket":"1-hold-the-landing-lease.md"},{"id":"t1-floor","command":"bench test --package ./internal/intent/... --package ./internal/commitment/... --package ./internal/worktree --package ./internal/refusalroute --package ./internal/gate/... --package ./internal/gocache --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance","ticket":"1-hold-the-landing-lease.md"},{"id":"t1-data-handling","command":"bench test --check data-handling-derivation","ticket":"1-hold-the-landing-lease.md"},{"id":"t1-system-lease","command":"bench test --check system --run TestTwoLandingProcessesShareOneLease","ticket":"1-hold-the-landing-lease.md"},{"id":"t1-root-conformance","command":"bench test --check root-conformance","ticket":"1-hold-the-landing-lease.md"}]},{"id":"BW-B","tickets":["2-route-a-moved-destination.md"],"verification":[{"id":"t2-land-focused","command":"bench test --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'","probe":"Swap the ancestry condition that maps a completion proof cause to land-destination-moved, so that a completion proof cause over a moved destination routes to land-handback. TestCompletionOverAMovedDestinationNamesTheFold turns red.","ticket":"2-route-a-moved-destination.md"},{"id":"t2-authorization-focused","command":"bench test --package ./internal/gate/authorization --run 'TestCompletionProofCauseTravelsAsAType'","ticket":"2-route-a-moved-destination.md"},{"id":"t2-floor","command":"bench test --package ./internal/gate/... --package ./internal/landing/... --package ./internal/worktree --package ./internal/refusalroute --package ./internal/commit --package ./cmd/bench --package ./internal/conformance","ticket":"2-route-a-moved-destination.md"},{"id":"t2-system-race","command":"bench test --check system --run TestWorktreeLandPublicRaceAndRerun","ticket":"2-route-a-moved-destination.md"},{"id":"t2-root-conformance","command":"bench test --check root-conformance","ticket":"2-route-a-moved-destination.md"}]},{"id":"BW-C","tickets":["3-record-a-fold.md","4-review-a-fold-resolution.md"],"verification":[{"id":"t3-record-focused","command":"bench test --package ./internal/reviewrecord/recordcmd --run 'TestRecordFold'","ticket":"3-record-a-fold.md"},{"id":"t3-chain-focused","command":"bench test --package ./internal/reviewrecord --run 'TestFoldChainAcceptsARecordedFoldBetweenChunks|TestFoldCoversACompletionSourcePastTheLastChunk'","probe":"Omit the fold-coverage branch for a later chunk base in checkSource, so that the record commit after a recorded fold reads as a chain gap. TestFoldChainAcceptsARecordedFoldBetweenChunks turns red.","ticket":"3-record-a-fold.md"},{"id":"t3-gate-focused","command":"bench test --package ./internal/gate --run 'TestReviewCheckpointChainGapNamesTheExpectedBase'","ticket":"3-record-a-fold.md"},{"id":"t3-floor","command":"bench test --package ./internal/reviewrecord/... --package ./internal/testreport --package ./internal/gate/... --package ./internal/landing/... --package ./internal/worktree --package ./cmd/bench --package ./internal/conformance --package ./internal/git","ticket":"3-record-a-fold.md"},{"id":"t3-root-conformance","command":"bench test --check root-conformance","ticket":"3-record-a-fold.md"},{"id":"t4-review-focused","command":"bench test --package ./internal/reviewrecord --run 'TestFoldResolutionNeedsReview|TestFoldDerivedFinalRequirements'","probe":"Omit the merged-test derived requirements that checkCompletion adds for each fold package, so that the complete checkpoint passes without a merged-test result. TestFoldDerivedFinalRequirements turns red.","ticket":"4-review-a-fold-resolution.md"},{"id":"t4-preflight-focused","command":"bench test --package ./internal/preflight --run 'TestPathsAuthorizedAdmitsARecordedResolutionPath'","ticket":"4-review-a-fold-resolution.md"},{"id":"t4-route-focused","command":"bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'","ticket":"4-review-a-fold-resolution.md"},{"id":"t4-help-focused","command":"bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'","ticket":"4-review-a-fold-resolution.md"},{"id":"t4-floor","command":"bench test --package ./internal/reviewrecord/... --package ./internal/preflight/... --package ./internal/gate/... --package ./internal/worktree --package ./internal/refusalroute --package ./cmd/bench --package ./internal/conformance --package ./internal/git","ticket":"4-review-a-fold-resolution.md"},{"id":"t4-root-conformance","command":"bench test --check root-conformance","ticket":"4-review-a-fold-resolution.md"}]},{"id":"BW-D","tickets":["5-hold-main-for-an-open-build.md","6-open-and-close-chunk-gaps.md"],"verification":[{"id":"t5-land-focused","command":"bench test --package ./internal/worktree --run 'TestLandRefusesAnUnrelatedLandingInAnOpenBuild|TestWindowAdmitsTheBuildsOwnLanding|TestWindowClosesAtCompletion|TestWindowClosesOnBlock|TestWindowClosesOnAbandon|TestWindowRecheckAtPublish|TestWindowWinsOverTheLease|TestLandingRefusalRegistryHasAProducingFixture'","probe":"Omit the window decision from the publish decision in PublishAdmitted, so that only the acquire decision decides the window. TestWindowRecheckAtPublish turns red.","ticket":"5-hold-main-for-an-open-build.md"},{"id":"t5-window-focused","command":"bench test --package ./internal/commitment/repository --run 'TestWindowNamesTheGapToWaitFor|TestOpenBuildsDeriveFromActiveOwners|TestWindowNeedsAGapOfEveryOtherBuild|TestStartDropsAStaleGap|TestClosersDropTheGapRecord|TestWindowFailsClosedOnAGitFault'","ticket":"5-hold-main-for-an-open-build.md"},{"id":"t5-floor","command":"bench test --package ./internal/intent/... --package ./internal/commitment/... --package ./internal/worktree --package ./internal/refusalroute --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance","ticket":"5-hold-main-for-an-open-build.md"},{"id":"t5-data-handling","command":"bench test --check data-handling-derivation","ticket":"5-hold-main-for-an-open-build.md"},{"id":"t5-system","command":"bench test --check system","ticket":"5-hold-main-for-an-open-build.md"},{"id":"t5-root-conformance","command":"bench test --check root-conformance","ticket":"5-hold-main-for-an-open-build.md"},{"id":"t6-gate-focused","command":"bench test --package ./internal/gate --run 'TestGreenChunkCheckpointOpensAGap|TestLastChunkOpensNoGap|TestRedCheckpointOpensNoGap|TestUnboundCheckpointOpensNoGap'","ticket":"6-open-and-close-chunk-gaps.md"},{"id":"t6-dispatch-focused","command":"bench test --package ./internal/reviewrecord/recordcmd --run 'TestFirstDispatchClosesTheGap|TestDispatchRefusesAnUnfoldedMain|TestDispatchInsideAChunkKeepsTheGap'","probe":"Omit the default-branch ancestry refusal from the first-dispatch ledger decision, so that a dispatch over an unfolded main closes the gap. TestDispatchRefusesAnUnfoldedMain turns red.","ticket":"6-open-and-close-chunk-gaps.md"},{"id":"t6-land-focused","command":"bench test --package ./internal/worktree --run 'TestOpenGapAdmitsAnUnrelatedLanding|TestWindowRecheckAtPublish'","ticket":"6-open-and-close-chunk-gaps.md"},{"id":"t6-floor","command":"bench test --package ./internal/commitment/... --package ./internal/gate/... --package ./internal/reviewrecord/... --package ./internal/worktree --package ./cmd/bench --package ./internal/conformance","ticket":"6-open-and-close-chunk-gaps.md"},{"id":"t6-system","command":"bench test --check system","ticket":"6-open-and-close-chunk-gaps.md"},{"id":"t6-root-conformance","command":"bench test --check root-conformance","ticket":"6-open-and-close-chunk-gaps.md"}]},{"id":"BW-E","tickets":["7-grant-a-reviewer-override.md"],"verification":[{"id":"t7-verb-focused","command":"bench test --package ./internal/commitment/commitcmd --run 'TestOverride'","probe":"Omit the terminal predicate check from the override verb, so that a pipe on stdin reaches the confirmation read. TestOverrideRefusesWithoutATerminal turns red.","ticket":"7-grant-a-reviewer-override.md"},{"id":"t7-land-focused","command":"bench test --package ./internal/worktree --run 'TestOverrideAdmitsOneLandingThroughEveryWindow|TestOverrideIsConsumedAtPublish|TestOverrideSurvivesARefusedLanding'","ticket":"7-grant-a-reviewer-override.md"},{"id":"t7-guard-focused","command":"bench test --package ./internal/benchguard --run 'TestGuardDeniesTheOverrideVerb'","ticket":"7-grant-a-reviewer-override.md"},{"id":"t7-help-focused","command":"bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'","ticket":"7-grant-a-reviewer-override.md"},{"id":"t7-floor","command":"bench test --package ./internal/commitment/... --package ./internal/intent/... --package ./internal/refusalroute --package ./internal/benchguard --package ./internal/guards --package ./internal/anchors --package ./internal/worktree --package ./cmd/bench --package ./internal/conformance","ticket":"7-grant-a-reviewer-override.md"},{"id":"t7-data-handling","command":"bench test --check data-handling-derivation","ticket":"7-grant-a-reviewer-override.md"},{"id":"t7-root-conformance","command":"bench test --check root-conformance","ticket":"7-grant-a-reviewer-override.md"}]},{"id":"BW-F","tickets":["8-document-the-landing-window.md"],"verification":[{"id":"t8-anchors-focused","command":"bench test --package ./internal/anchors --run 'TestChunkChainAnchors'","probe":"Omit the forbid row for the retired sentence that folds main at the start of the last chunk from internal/anchors/registry_chunk_chain.go, so that the old sentence passes the registry. TestChunkChainAnchors turns red.","ticket":"8-document-the-landing-window.md"},{"id":"t8-floor","command":"bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance","ticket":"8-document-the-landing-window.md"},{"id":"t8-guidance","command":"bench test --check docs-currency-workflow","ticket":"8-document-the-landing-window.md"},{"id":"t8-conformance-meta","command":"bench test --check conformance-meta","ticket":"8-document-the-landing-window.md"},{"id":"t8-prose","command":"bench test --check prose","ticket":"8-document-the-landing-window.md"},{"id":"t8-root-conformance","command":"bench test --check root-conformance","ticket":"8-document-the-landing-window.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/build-landing-window/spec.md"},{"id":"intent","command":"bench test --package ./internal/intent/..."},{"id":"commitment","command":"bench test --package ./internal/commitment/..."},{"id":"reviewrecord","command":"bench test --package ./internal/reviewrecord/..."},{"id":"gate","command":"bench test --package ./internal/gate/..."},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"command-package","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"}]}
```

## Testing decisions

A good landing test drives `worktree.LandCommand` through `runVerb` over a fixture repository with a bound delivery, as `TestCommitmentGateRace` does.
It observes the refusal record, its `next=` route, the destination ref, and the ledger.
A race row holds the gate at a FIFO barrier (`landingFixtureWithGateStep`) and changes the ledger or the destination inside that window.
A row that needs two processes runs in `internal/systemtest`, as `TestWorktreeLandPublicRaceAndRerun` does, because a record lock belongs to its process.

A record row drives `recordcmd` over `recordtest.Attach`, as `TestReviewCheckpointChainGapNamesTheExpectedBase` does.
A fold fixture makes the fold with the real `bench worktree merge` owner, or with a hand `git merge` for a resolution.
A gap row drives `gate.RunCommand` with `--checkpoint` and reads the ledger after it.

An override row injects the terminal predicate and the confirmation reader, as `bench setup` injects `isTerminal` (`internal/adopt/setup.go:25`).
Non-terminal rows drive a pipe, never the null device, because `/dev/null` reads as a character device.

Every expected TOON row derives through `toon.Table`, so a numeric-looking chunk id keeps the encoder's quotes.

### Seam diagram

    trigger: an operator runs bench worktree land
        │
        ▼
    request, ledger  ──▶  [ acquire: window, record lock, lease record ] ──▶ [ compose, gate ] ──▶ [ publish: window, lease and override removal, CAS ]  ──▶  landed or refused{next=...}
                      ◀ tests attach here: runVerb with a FIFO gate barrier, and two real processes in the system suite

    trigger: the build runs bench gate --checkpoint --chunk, then bench record assignment
        │
        ▼
    plan at HEAD, ledger  ──▶  [ gap open on green, gap close on the first dispatch with the fold check ]  ──▶  ledger gap record or dispatch refusal
                      ◀ tests attach here: gate.RunCommand and recordcmd over a bound fixture build

    trigger: the build author runs bench record fold after bench worktree merge
        │
        ▼
    fold commit, review record  ──▶  [ fold entry, resolution, merged-test packages ] ──▶ [ checkSource, checkCompletion ]  ──▶  checkpoint verdict
                      ◀ tests attach here: recordtest fixtures with a clean fold and a hand-resolved fold

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| BW1 | 1, 9 | A landing that reaches its gate holds a lease record for `refs/heads/main` with its own assignment id in the ledger | `internal/worktree/land_lease_test.go` (`TestLandHoldsTheDestinationLease`) | A lease taken after the gate leaves the composition unleased. |
| BW2 | 2, 8 | A second landing during the first landing's gate refuses with the detail `landing destination refs/heads/main is leased by assignment <first-id>` | `internal/worktree/land_lease_test.go` (`TestSecondLandingMeetsTheLease`) | The current second landing meets the gate lock and names `bench doctor`. |
| BW3 | 3 | That refusal's `next=` begins with `wait until the landing of assignment <first-id> ends` and ends with the caller's rerun | `internal/worktree/land_lease_test.go` (`TestSecondLandingMeetsTheLease`) | A route without the wait sends the operator to rerun into the same refusal. |
| BW4 | 2 | The second landing runs no composition and no gate, so the gate tally holds one run | `internal/worktree/land_lease_test.go` (`TestSecondLandingMeetsTheLease`) | A lease checked after composition spends a gate run before the refusal. |
| BW5 | 8 | The second landing's output holds no `bench doctor` and no `gate owner: pid` | `internal/worktree/land_lease_test.go` (`TestSecondLandingMeetsTheLease`) | A lease that loses the order to the gate lock keeps the old faces. |
| BW6 | 4 | After a published landing, the ledger holds no lease for the destination | `internal/worktree/land_lease_test.go` (`TestLandHoldsTheDestinationLease`) | A lease kept past publish blocks the next landing. |
| BW7 | 4 | A competing acquire that waits on the intent lock during the publish decision sees no lease | `internal/commitment/repository/lease_test.go` (`TestPublishReleasesTheLeaseInItsDecision`) | A release in a later transaction leaves a window where the published landing still holds the lease. |
| BW8 | 5 | After a landing refuses red at its gate, the next landing acquires the lease | `internal/worktree/land_lease_test.go` (`TestLandReleasesTheLeaseOnRefusal`) | A release only on publish blocks the destination after every red. |
| BW9 | 6 | A seeded lease record for another assignment, with no holder of the lock file, is replaced by the next landing's record | `internal/commitment/repository/lease_test.go` (`TestLeaseTakesOverADeadOwner`) | A presence-only lease refuses forever after a kill. |
| BW10 | 6, 1 | Two landing processes contend, and the loser refuses with the lease detail while the winner publishes | `internal/systemtest/owner_land_lease_test.go` (`TestTwoLandingProcessesShareOneLease`) | An in-process map alone cannot see a second process. |
| BW11 | 7 | A `--resume` of a published landing completes while another landing holds the lease | `internal/worktree/land_lease_test.go` (`TestLandResumeTakesNoLease`) | A resume that takes the lease waits on unrelated work. |
| BW12 | 10 | `ValidateCommitment` refuses a ledger with two leases on one destination | `internal/intent/build_records_test.go` (`TestValidateCommitmentRefusesDuplicateBuildRecords`) | A hand-edited second owner then authorizes two landings. |
| BW13 | 12 | A tolerant transaction over a ledger that holds a lease, a gap, and an override writes all three back | `internal/intent/build_records_test.go` (`TestTolerantReadKeepsTheBuildRecords`) | The second field list drops a new record on a tolerant write. |
| BW14 | 11 | The commitment record inventory names `commitment.leases.destination`, `commitment.gaps.chunk`, and `commitment.overrides.reason` | `internal/conformance/commitment_record_test.go` (`TestCommitmentRecordInventoryNamesTheBuildRecords`) | An undocumented field passes only if the inventory skips the new lists. |
| BW15 | 13 | A directory at the landing lock path refuses with `landing lease unavailable:` before composition | `internal/worktree/land_lease_test.go` (`TestLeaseLockUnavailable`) | A lease fault that proceeds runs the landing without a lease. |
| BW16 | 2 | The recovery matrix and the landing face fixtures hold `land-lease-held` | `internal/worktree/identity_component_test.go` (`TestLandingRefusalRegistryHasAProducingFixture`) | A face without a producing fixture reds the registry test. |
| BW17 | 59 | A completion landing over a destination with one new commit refuses with face `land-destination-moved` | `internal/worktree/land_moved_destination_test.go` (`TestCompletionOverAMovedDestinationNamesTheFold`) | The current route is `land-infrastructure`. |
| BW18 | 59 | That refusal's `next=` holds `bench worktree merge --from main` and ends with the caller's rerun | `internal/worktree/land_moved_destination_test.go` (`TestCompletionOverAMovedDestinationNamesTheFold`) | A route that keeps the family word without the merge command passes a token check. |
| BW19 | 59 | That refusal's output holds no `bench doctor` | `internal/worktree/land_moved_destination_test.go` (`TestCompletionOverAMovedDestinationNamesTheFold`) | The infrastructure face survives beside the new face. |
| BW20 | 60 | A raw commit on the destination during the loser's gate gives a refusal whose `next=` holds `bench worktree merge --from main` | `internal/systemtest/owner_land_race_test.go` (`TestWorktreeLandPublicRaceAndRerun`) | The post-gate recheck still routes to the reviewer handback. |
| BW21 | 61 | A destination that moves inside the publication gap gives a `TipMovedError` refusal with face `land-destination-moved` | `internal/worktree/land_moved_destination_test.go` (`TestCompareAndSwapRefusalNamesTheFold`) | The landing gives the compare-and-swap refusal no face of its own. |
| BW22 | 62 | An ignored-file change in the destination checkout during the gate, with an unmoved ref, keeps face `land-handback` | `internal/worktree/land_moved_destination_test.go` (`TestCheckoutChangeWithoutARefMoveKeepsTheHandback`) | A face picked by the recheck alone prescribes a fold for an ignored file. |
| BW23 | 63 | A completion proof fault over a source that holds the destination base gives face `land-handback` | `internal/worktree/land_moved_destination_test.go` (`TestProofFaultOverACurrentSourceHandsBack`) | An ancestry-blind mapping prescribes a fold where none helps. |
| BW24 | 64 | `authorization.Result` carries the completion proof cause as a typed field when the gate refuses a completion proof | `internal/gate/authorization/authorization_test.go` (`TestCompletionProofCauseTravelsAsAType`) | A reason-string match breaks when the proof text changes. |
| BW25 | 46 | `bench record fold` over a fold that `bench worktree merge` published writes entry `F1` with `After`, `From`, `Commit`, `Incoming`, and the source digest | `internal/reviewrecord/recordcmd/fold_test.go` (`TestRecordFoldWritesTheFoldEntry`) | A form that writes no parents cannot bind the fold to the chain. |
| BW26 | 47 | That clean fold's `Resolution` is empty | `internal/reviewrecord/recordcmd/fold_test.go` (`TestRecordFoldWritesTheFoldEntry`) | A resolution that lists every merged path demands review of a clean fold. |
| BW27 | 48 | A hand-resolved conflict fold lists exactly the conflicted path in `Resolution` | `internal/reviewrecord/recordcmd/fold_test.go` (`TestRecordFoldNamesTheResolutionPaths`) | A first-parent diff lists every merged path. |
| BW28 | 48 | A clean fold with one extra hand edit lists exactly the edited path in `Resolution` | `internal/reviewrecord/recordcmd/fold_test.go` (`TestRecordFoldNamesTheResolutionPaths`) | A conflict-only rule hides a hand edit in a clean fold. |
| BW29 | 48 | A resolution path that holds a space and a double quote is recorded with its exact bytes | `internal/reviewrecord/recordcmd/fold_test.go` (`TestRecordFoldNamesTheResolutionPaths`) | A quoted-header parse records `"a \"b"` instead of the path. |
| BW30 | 53 | A source that is not HEAD refuses with `is not HEAD` before any write | `internal/reviewrecord/recordcmd/fold_test.go` (`TestRecordFoldRefusesANonFold`) | A form that records an old commit binds the wrong fold. |
| BW31 | 53 | A one-parent commit refuses with `is not a two-parent merge` | `internal/reviewrecord/recordcmd/fold_test.go` (`TestRecordFoldRefusesANonFold`) | Any commit then passes as a fold. |
| BW32 | 53 | A merge of a sibling tip outside the default branch refuses with `does not merge a default-branch commit` | `internal/reviewrecord/recordcmd/fold_test.go` (`TestRecordFoldRefusesANonFold`) | A sibling fold skips the chunk review. |
| BW33 | 53 | A fold whose first parent holds an unreviewed code commit refuses with `does not follow the reviewed source` | `internal/reviewrecord/recordcmd/fold_test.go` (`TestRecordFoldRefusesANonFold`) | A fold record covers build work that no chunk reviewed. |
| BW34 | 50, 58 | A chunk-2 checkpoint passes when chunk 2's base is the record commit that follows a recorded clean fold over the chunk-1 tip | `internal/reviewrecord/fold_test.go` (`TestFoldChainAcceptsARecordedFoldBetweenChunks`) | The current chain refuses that base as a stale gap. |
| BW35 | 52 | The same chunk-2 checkpoint without the fold entry refuses as a repair of form `fold` | `internal/reviewrecord/fold_test.go` (`TestFoldChainAcceptsARecordedFoldBetweenChunks`) | A chain that accepts any merge base skips the fold record. |
| BW36 | 52 | The chain-gap refusal holds `a default-branch merge lands before the first chunk or as a recorded fold` | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointChainGapNamesTheExpectedBase`) | The old text sends the operator to a last-chunk fold. |
| BW37 | 51 | `bench record completion` passes at a source past the last chunk tip when a clean recorded fold covers it | `internal/reviewrecord/fold_test.go` (`TestFoldCoversACompletionSourcePastTheLastChunk`) | The current refusal demands a refreeze and a full re-record. |
| BW38 | 52 | A completion at that source without the fold entry refuses with `stale reviewed source` as a repair of form `fold` | `internal/reviewrecord/fold_test.go` (`TestFoldCoversACompletionSourcePastTheLastChunk`) | A refusal without the form leaves the operator on the refreeze route. |
| BW39 | 49 | A fold with a non-empty resolution and two passing axes refuses the checkpoint until the third axis passes | `internal/reviewrecord/fold_test.go` (`TestFoldResolutionNeedsReview`) | A rule that counts one review passes a partly reviewed repair. |
| BW40 | 47 | A fold with an empty resolution and no review passes the checkpoint | `internal/reviewrecord/fold_test.go` (`TestFoldResolutionNeedsReview`) | A rule that always demands review costs a round for a clean fold. |
| BW41 | 54, 55 | A fold that brings in `pkg/a/a_test.go` stores `TestPackages` equal to `testreport.WritesFloor` over that path | `internal/reviewrecord/recordcmd/fold_test.go` (`TestRecordFoldDerivesMergedTestPackages`) | A second path-to-package rule drifts from the `--changed` index. |
| BW42 | 54 | The complete checkpoint refuses until a passing final result exists for `merged-test:<package>` | `internal/reviewrecord/fold_test.go` (`TestFoldDerivedFinalRequirements`) | The ft336 final checks omitted the merged package and stayed green. |
| BW43 | 56 | `paths-authorized` passes for a resolution path that the record's fold entry lists and the fence omits | `internal/preflight/folded_main_test.go` (`TestPathsAuthorizedAdmitsARecordedResolutionPath`) | The FT337 landing needed a fence expansion for merged paths. |
| BW44 | 56 | `paths-authorized` stays red for the same path when no fold entry lists it | `internal/preflight/folded_main_test.go` (`TestPathsAuthorizedAdmitsARecordedResolutionPath`) | A rule that admits every merge path unfences the build's own writes. |
| BW45 | 57 | The `merge-fold-red` route names the hand merge, `bench commit`, and `bench record fold`, and names no `FT342` | `internal/worktree/merge_route_test.go` (`TestMergeFacesFollowTheirRoutes`) | The current route names the open FT342 question. |
| BW46 | 14, 17 | After `bench commitment start` binds a staged spec, an unbound landing refuses with face `land-build-window` | `internal/worktree/land_window_test.go` (`TestLandRefusesAnUnrelatedLandingInAnOpenBuild`) | The current landing admits any landing during a build. |
| BW47 | 15, 17 | That refusal names the outcome, the deliverable, and `its first chunk checkpoint` | `internal/worktree/land_window_test.go` (`TestLandRefusesAnUnrelatedLandingInAnOpenBuild`) | A refusal without the gap leaves the operator to poll. |
| BW48 | 15 | With a closed gap record for chunk `C2` that is not last, the refusal names `the checkpoint of chunk C2` | `internal/commitment/repository/window_test.go` (`TestWindowNamesTheGapToWaitFor`) | A refusal that names the recorded chunk's predecessor sends the operator to a past gap. |
| BW49 | 15, 27 | With a closed gap record whose `Last` is true, the refusal names `its completion landing` | `internal/commitment/repository/window_test.go` (`TestWindowNamesTheGapToWaitFor`) | A refusal that names a checkpoint after the last chunk waits for a gap that never opens. |
| BW50 | 16 | A landing of the bound spec from a sibling assignment that `--from` bound to the build passes the window | `internal/worktree/land_window_test.go` (`TestWindowAdmitsTheBuildsOwnLanding`) | A rule keyed to one assignment refuses the build's own worktrees. |
| BW51 | 18 | After the build's completion landing publishes, an unbound landing publishes | `internal/worktree/land_window_test.go` (`TestWindowClosesAtCompletion`) | A window kept past delivery freezes `main` for good. |
| BW52 | 19 | After `bench commitment block` on the outcome, an unbound landing publishes | `internal/worktree/land_window_test.go` (`TestWindowClosesOnBlock`) | A window keyed to bindings survives the block, which keeps the binding. |
| BW53 | 20 | After the release of every bound assignment, an unbound landing publishes | `internal/worktree/land_window_test.go` (`TestWindowClosesOnAbandon`) | A window keyed to the claim alone survives an abandoned build. |
| BW54 | 20 | With one of two bound assignments released, an unbound landing still refuses | `internal/commitment/repository/window_test.go` (`TestOpenBuildsDeriveFromActiveOwners`) | A rule that needs every bound assignment active closes the window early. |
| BW55 | 21 | A started outcome whose deliverable is a tickets-only folder refuses no landing | `internal/commitment/repository/window_test.go` (`TestOpenBuildsDeriveFromActiveOwners`) | A window for every binding holds `main` for a one-ticket outcome. |
| BW56 | 22, 33 | A build that starts during an unrelated landing's gate refuses that landing at publish, and `main` stays unchanged | `internal/worktree/land_window_test.go` (`TestWindowRecheckAtPublish`) | A check only at acquisition lets the landing publish under the new build. |
| BW57 | 23 | A light-path landing with a tickets-only `--spec` refuses with `land-build-window` while a build is outside a gap | `internal/worktree/land_window_test.go` (`TestLandRefusesAnUnrelatedLandingInAnOpenBuild`) | A rule that admits every closing landing lets the build's own light-path dispatch move `main`. |
| BW58 | 24 | With two open builds, a landing inside build A's gap and outside build B's gap refuses and names build B | `internal/commitment/repository/window_test.go` (`TestWindowNeedsAGapOfEveryOtherBuild`) | A rule that checks one build ignores the parallel grant. |
| BW59 | 25 | With an open build outside a gap and a live lease holder, the refusal is the window refusal | `internal/worktree/land_window_test.go` (`TestWindowWinsOverTheLease`) | The reversed order sends the operator to wait for a lease and then meet the window. |
| BW60 | 26 | A green `bench gate --checkpoint <spec> --chunk C1` in the bound integration worktree writes an open gap for chunk `C1` | `internal/gate/chunk_gap_test.go` (`TestGreenChunkCheckpointOpensAGap`) | No gap opens, so no landing can ever land during the build. |
| BW61 | 29 | During that open gap, an unbound landing publishes | `internal/worktree/land_window_test.go` (`TestOpenGapAdmitsAnUnrelatedLanding`) | A window that ignores the gap record refuses every landing. |
| BW62 | 27 | A green checkpoint of the plan's last chunk writes no gap | `internal/gate/chunk_gap_test.go` (`TestLastChunkOpensNoGap`) | A gap after the last chunk lets `main` move after the freeze. |
| BW63 | 28 | A red chunk checkpoint writes no gap | `internal/gate/chunk_gap_test.go` (`TestRedCheckpointOpensNoGap`) | A gap before acceptance lets `main` move under an unreviewed chunk. |
| BW64 | 34 | A green chunk checkpoint in an unbound worktree writes no gap | `internal/gate/chunk_gap_test.go` (`TestUnboundCheckpointOpensNoGap`) | Any checkpoint run then opens a gap for an unrelated build. |
| BW65 | 30 | The first `bench record assignment` of chunk `C2` writes a closed gap for `C2` | `internal/reviewrecord/recordcmd/dispatch_gap_test.go` (`TestFirstDispatchClosesTheGap`) | A gap that stays open lets `main` move during chunk `C2`. |
| BW66 | 31 | That dispatch refuses with `dispatch refused: the default branch tip` when `main` moved past HEAD | `internal/reviewrecord/recordcmd/dispatch_gap_test.go` (`TestDispatchRefusesAnUnfoldedMain`) | The charge then pins a base behind `main`. |
| BW67 | 31 | That refusal names `bench worktree merge --from main` and `bench record fold` and leaves the plan fence unchanged | `internal/reviewrecord/recordcmd/dispatch_gap_test.go` (`TestDispatchRefusesAnUnfoldedMain`) | A refusal that writes the plan first records a dispatch that never ran. |
| BW68 | 32 | A second dispatch inside chunk `C2` leaves the ledger bytes unchanged | `internal/reviewrecord/recordcmd/dispatch_gap_test.go` (`TestDispatchInsideAChunkKeepsTheGap`) | A dispatch that rewrites the record can reopen or move the gap. |
| BW69 | 17, 30 | The first dispatch of the first chunk writes a closed gap for that chunk | `internal/reviewrecord/recordcmd/dispatch_gap_test.go` (`TestFirstDispatchClosesTheGap`) | Without it, the refusal cannot name the chunk to wait for. |
| BW70 | 33 | A first dispatch that closes the gap during an unbound landing's gate refuses that landing at publish | `internal/worktree/land_window_test.go` (`TestWindowRecheckAtPublish`) | A landing that publishes after the close moves the dispatched chunk's base. |
| BW71 | 35 | An override for an active assignment writes one record with its reason and prints `overridden{assignment=<id>}` | planned internal/commitment/commitcmd/override_test.go (TestOverrideWritesOneRecord) | A verb that prints success without the record admits nothing. |
| BW72 | 36 | A landing whose assignment holds an override publishes while two builds are outside their gaps | planned internal/worktree/land_override_test.go (TestOverrideAdmitsOneLandingThroughEveryWindow) | An override scoped to one build leaves the parallel build's window shut. |
| BW73 | 36 | A landing of another assignment still refuses while the override exists | planned internal/worktree/land_override_test.go (TestOverrideAdmitsOneLandingThroughEveryWindow) | An override that lifts the window for all landings is a second open gap. |
| BW74 | 37 | After the overridden landing publishes, the ledger holds no override for its assignment | planned internal/worktree/land_override_test.go (TestOverrideIsConsumedAtPublish) | A kept override admits a later publication. |
| BW75 | 38 | After the overridden landing refuses red at its gate, the override record stays | planned internal/worktree/land_override_test.go (TestOverrideSurvivesARefusedLanding) | A consumption at acquisition forces a second grant for each rerun. |
| BW76 | 39 | With a pipe on stdin, the verb refuses with `needs the reviewer at a terminal` and writes nothing | planned internal/commitment/commitcmd/override_test.go (TestOverrideRefusesWithoutATerminal) | A verb without the predicate lets an agent tool call grant itself. |
| BW77 | 40 | With `/dev/null` on stdin, the verb refuses with `not confirmed` and writes nothing | planned internal/commitment/commitcmd/override_test.go (TestOverrideNeedsTheTypedAssignment) | The null device passes the character-device predicate, so only the confirmation stops it. |
| BW78 | 40 | With an injected terminal and a confirmation line that differs in one byte, the verb refuses with `not confirmed` | planned internal/commitment/commitcmd/override_test.go (TestOverrideNeedsTheTypedAssignment) | A prefix or trimmed compare accepts a wrong id. |
| BW79 | 41 | From a Bench worktree, the verb refuses with `runs from the primary checkout` before the terminal check | planned internal/commitment/commitcmd/override_test.go (TestOverrideRefusesOutsideThePrimaryCheckout) | A worktree's own executable writes the override. |
| BW80 | 43 | An override for an unknown assignment refuses with `is not active` | planned internal/commitment/commitcmd/override_test.go (TestOverrideRefusesAnInactiveAssignment) | A typo then grants a record that no landing uses. |
| BW81 | 44 | A reason with a newline, and an empty reason, each refuse with the block verb's reason refusal | planned internal/commitment/commitcmd/override_test.go (TestOverrideRefusesABadReason) | A second reason rule drifts from the block verb's rule. |
| BW82 | 42 | The follow-on guard denies `bench commitment override --assignment x --reason y` in a Bash envelope | planned internal/benchguard/benchguard_test.go (TestGuardDeniesTheOverrideVerb) | The guard classifies only follow-ons today, so the verb runs. |
| BW83 | 42 | The guard denies `bench commitment --in primary override` and `bench worktree exec t -- bench commitment override` | planned internal/benchguard/benchguard_test.go (TestGuardDeniesTheOverrideVerb) | A first-word match misses the flag form and the child form. |
| BW84 | 35 | The help inventory lists the `bench commitment override` form | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A hand-copied help line drifts from the commitment form table. |
| BW85 | 65 | The implement-spec `Land` section states the gap fold before the first dispatch of the next chunk and the fold record | review-owned: the reviewer reads the section; the anchor check reds a drift of the new needles | No check reads the meaning of the fold rule. |
| BW86 | 65 | The implement-spec guidance holds no sentence that folds `main` at the start of the last chunk | review-owned: the reviewer reads the section; a forbid anchor names the old sentence | A surviving old sentence routes a fold that the chain refuses. |
| BW87 | 66 | The review guidance states the fold record, the resolution review, and the landing base after a fold | review-owned: the reviewer reads step 1 and the landing step | No check reads the review route. |
| BW88 | 59, 66 | The final-check guidance names the moved-destination route and no `bench doctor` route for a moved destination | review-owned: the reviewer reads the landing step | No check reads that step. |
| BW89 | 67, 45 | The reference guide states the lease, the window, the gap, the override, the revert back-out route, and `bench commitment block` as the stop for a stalled build | review-owned: the reviewer reads the Command Notes | No check reads the Command Notes content. |
| BW90 | 68 | `CONTEXT.md` defines the five terms, each with an Avoid list | review-owned: the reviewer reads the entries | The prose check grades sentences, not terms. |
| BW91 | 58, 69 | `docs/known-issues.md` holds no `base-current` row | review-owned: the reviewer reads the table | No check reads the issue table. |
| BW92 | 69 | `CHANGELOG.md` names the lease, the window, the override, the fold record, and the moved-destination face | review-owned: the reviewer reads the entry | No check reads the changelog content. |
| BW93 | 46 | The help inventory lists the `bench record fold` form, and the review form names `--fold` | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A hand-copied help line drifts from the record form registry. |
| BW94 | 16, 23 | A spec-less landing from a sibling assignment that `--from` bound to the build refuses with face `land-build-window` while the build is outside a gap | `internal/worktree/land_window_test.go` (`TestWindowAdmitsTheBuildsOwnLanding`) | A rule keyed to the bound assignment alone lets a build worktree land spec-less under a chunk. |
| BW96 | 17, 19 | After `bench commitment block` during an open gap and a later `bench commitment unblock` and `bench commitment start`, an unbound landing refuses with face `land-build-window` | `internal/commitment/repository/window_test.go` (`TestStartDropsAStaleGap`) | A start that keeps the old open gap reopens the build with `main` already free to move. |
| BW99 | 27, 30 | A green checkpoint of chunk `C1` while the gap record holds chunk `C3` with `Last` true leaves the ledger bytes unchanged | `internal/gate/chunk_gap_test.go` (`TestEarlierChunkCheckpointKeepsTheGap`) | A repair checkpoint of an earlier chunk reopens the gap after the last freeze. |
| BW100 | 14 | An unresolved default branch with a claimed outcome refuses the landing with face `land-infrastructure` | `internal/commitment/repository/window_test.go` (`TestWindowFailsClosedOnAGitFault`) | A Git fault opens every window. |
| BW97 | 18, 19 | After `bench commitment block`, and after the build's completion landing publishes, the ledger holds no gap record for the outcome | `internal/commitment/repository/window_test.go` (`TestClosersDropTheGapRecord`) | A closer that keeps the gap record leaves a stale record for the next start to trust. |
| BW98 | 56 | `paths-authorized` stays red for a recorded resolution path that a later commit changes after the fold | `internal/preflight/folded_main_test.go` (`TestPathsAuthorizedAdmitsARecordedResolutionPath`) | A name-only skip lets a later build write to the path pass the fence. |
| BW95 | 32 | A repair dispatch to a chunk-`C1` ticket while the gap record holds chunk `C3` with `Last` true leaves the ledger bytes unchanged | `internal/reviewrecord/recordcmd/dispatch_gap_test.go` (`TestDispatchInsideAChunkKeepsTheGap`) | A chunk-inequality rule rewrites the record to `C1`, so the refusal names a past gap and the build loses its last-chunk state. |

Not covered: story 70 — a reviewed exclusion; the Out of scope section prices it.
Not covered: story 71 — a reviewed exclusion; the Edge inventory keeps the current gate-lock face.
Not covered: story 72 — a reviewed exclusion; the trust chain states the assumption.
Not covered: story 73 — a reviewed exclusion; the Out of scope section prices it.

### Edge inventory

The audience is this repository and every repository that links the kit.
A linked repository runs the same landing, record, and gate owners, so every row applies there.
An absent `commitment` block and an empty `gaps` list both read as no recorded gap (BW47).
An absent landing lock file is created; a directory at that path refuses (BW15).

The hostile-input walk covers these classes:

- An assignment id operand with spaces or glob characters fails the identity pattern before any read.
- A reason with ESC, BEL, a newline, or no text refuses under the block verb's reason rule (BW81).
- A confirmation line with a trailing CR or a NBSP differs in a byte and refuses (BW78).
- A non-terminal stdin drives a pipe; the null device reads as a character device and stops at the confirmation (BW76, BW77).
- A guard envelope can spell `&&` as `&&`; the existing benchguard separator tests cover it, and BW83 adds the flag and child forms.
- A resolution path with a double quote, a backslash, or a space reads through `-z` plumbing (BW29).
- A chunk id that reads as a number keeps the encoder's quotes, because each expected row derives through `toon.Table`.
- A dangling or live symbolic link at the landing lock path is opened without following.
- Two processes reload the ledger across a real boundary in BW10.
- A cwd deeper than the root resolves the root before the primary-checkout test of BW79.

**Won't handle** — a non-landing gate run that holds the gate lock — the landing keeps `land-infrastructure`, and the caller's rerun stays the surviving route.
**Won't handle** — a `--wait` queue on the lease — the refusal names the owner, and the rerun stays the surviving route.
**Won't handle** — an older Bench binary that reads a new ledger record — its strict decode refuses, and the broker refresh stays the surviving route.
**Won't handle** — two open builds that are both past their last freeze — each completion landing refuses, and the reviewer override stays the surviving route.
**Won't handle** — a forged override record written by a file tool — the honest-mistake model excludes it, and the override verb stays the surviving route.

**Won't handle** — a fold of a sibling tip — the chunk review covers a sibling fold, and `bench worktree merge` stays its surviving route.
**Won't handle** — a fold recorded after a later commit — the form refuses, and a fresh fold with an immediate record stays the surviving route.
**Won't handle** — a legacy continuation run — it holds no delivery binding, so it opens no window; its own landing stays the surviving caller.
**Won't handle** — a stale override for a released assignment — no landing can present it, and a new override stays the surviving route.

## Ownership fences

The ticket slice derives this fence, so it equals the union of the ticket `Writes:` lines, plus the review pickup and the spec folder.
It holds the owner and test files, the command registry closure, the anchor registry closure, and each canary unit that pins an edited file.

- `.agents/commands/bench-drain.md`
- `.agents/commands/bench-final-check.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-write-spec.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.bench/BENCH-reference.md`
- `.bench/commitment.json`
- `.bench/hooks/block-bench-follow-on.sh`
- `.bench/structure-accept`
- `.bench/structure.budgets`
- `CHANGELOG.md`
- `CONTEXT.md`
- `DATA_HANDLING.md`
- `README.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/commitment_command.go`
- `cmd/bench/commitment_test.go`
- `cmd/bench/guards.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/isolated_command_fixtures_test.go`
- `cmd/bench/main.go`
- `cmd/bench/record_operations_test.go`
- `cmd/bench/response_bound_exempt_test.go`
- `docs/adr/0014-main-receives-writes-only-through-landings.md`
- `docs/known-issues.md`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `internal/anchors/registry_charge_binding.go`
- `internal/anchors/registry_charge_binding_test.go`
- `internal/anchors/registry_check_floor_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_commitment.go`
- `internal/anchors/registry_craft_research.go`
- `internal/anchors/registry_craft_research_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `internal/benchguard/benchguard.go`
- `internal/benchguard/benchguard_test.go`
- `internal/benchguard/reviewer.go`
- `internal/commit/commit.go`
- `internal/commitment/admission.go`
- `internal/commitment/admission_test.go`
- `internal/commitment/commitcmd/command.go`
- `internal/commitment/commitcmd/override.go`
- `internal/commitment/commitcmd/override_route_test.go`
- `internal/commitment/commitcmd/override_test.go`
- `internal/commitment/commitcmd/refusal_route_test.go`
- `internal/commitment/commitmenttest/admission.go`
- `internal/commitment/concurrent_test.go`
- `internal/commitment/continuation_test.go`
- `internal/commitment/repository/admission.go`
- `internal/commitment/repository/closure.go`
- `internal/commitment/repository/gap.go`
- `internal/commitment/repository/lease.go`
- `internal/commitment/repository/lease_test.go`
- `internal/commitment/repository/light_path_test.go`
- `internal/commitment/repository/override.go`
- `internal/commitment/repository/publication.go`
- `internal/commitment/repository/publication_lock_test.go`
- `internal/commitment/repository/publication_test.go`
- `internal/commitment/repository/readiness.go`
- `internal/commitment/repository/sources.go`
- `internal/commitment/repository/window.go`
- `internal/commitment/repository/window_test.go`
- `internal/commitment/store_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/commitment_record_test.go`
- `internal/conformance/data_handling_test.go`
- `internal/conformance/ordinary_build_census_test.go`
- `internal/conformance/refusal_route_guard_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/gate/authorization/authorization.go`
- `internal/gate/authorization/authorization_test.go`
- `internal/gate/checkpoint.go`
- `internal/gate/chunk_gap.go`
- `internal/gate/chunk_gap_test.go`
- `internal/gate/complete_checkpoint.go`
- `internal/gate/complete_checkpoint_test.go`
- `internal/gate/gate.go`
- `internal/gate/refusal_route_test.go`
- `internal/gate/review_checkpoint_commits_test.go`
- `internal/gate/review_checkpoint_test.go`
- `internal/gate/run_transaction.go`
- `internal/gate/verdict.go`
- `internal/gocache/lock.go`
- `internal/git/admin_readers_test.go`
- `internal/git/default_branch.go`
- `internal/git/git.go`
- `internal/guards/guards_test.go`
- `internal/intent/build_records_test.go`
- `internal/intent/ledger/commitment.go`
- `internal/intent/ledger/validate.go`
- `internal/intent/ledger_aliases.go`
- `internal/intent/ledger_aliases_test.go`
- `internal/intent/transaction.go`
- `internal/landing/completion_evidence_test.go`
- `internal/landing/composition.go`
- `internal/landing/composition_fold_test.go`
- `internal/landing/composition_test.go`
- `internal/landing/gitexec.go`
- `internal/landing/landing.go`
- `internal/landing/landing_reviewed_test.go`
- `internal/landing/landing_test.go`
- `internal/landing/merge.go`
- `internal/landing/state_test.go`
- `internal/preflight/decision.go`
- `internal/preflight/default_branch.go`
- `internal/preflight/folded_main_test.go`
- `internal/preflight/gather.go`
- `internal/preflight/gather_inputs.go`
- `internal/preflight/check_floor.go`
- `internal/preflight/paths_authorized.go`
- `internal/preflight/preflighttest/fixture.go`
- `internal/preflight/review.go`
- `internal/refusalroute/command_test.go`
- `internal/refusalroute/faces_commitment.go`
- `internal/refusalroute/faces_land.go`
- `internal/refusalroute/faces_merge.go`
- `internal/refusalroute/registry.go`
- `internal/refusalroute/registry_test.go`
- `internal/reviewrecord/assignment.go`
- `internal/reviewrecord/check.go`
- `internal/reviewrecord/coverage.go`
- `internal/reviewrecord/delegated.go`
- `internal/reviewrecord/fold.go`
- `internal/reviewrecord/fold_test.go`
- `internal/reviewrecord/parse.go`
- `internal/reviewrecord/plan.go`
- `internal/reviewrecord/record.go`
- `internal/reviewrecord/recordcmd/assignment.go`
- `internal/reviewrecord/recordcmd/assignment_test.go`
- `internal/reviewrecord/recordcmd/command.go`
- `internal/reviewrecord/recordcmd/dispatch_gap_test.go`
- `internal/reviewrecord/recordcmd/fold.go`
- `internal/reviewrecord/recordcmd/fold_test.go`
- `internal/reviewrecord/recordcmd/ordering_test.go`
- `internal/reviewrecord/recordcmd/refusal.go`
- `internal/reviewrecord/recordcmd/review_test.go`
- `internal/reviewrecord/recordcmd/verification.go`
- `internal/reviewrecord/recordcmd/verification_test.go`
- `internal/reviewrecord/recordtest/fixture.go`
- `internal/reviewrecord/recordtest/fold.go`
- `internal/reviewrecord/repair.go`
- `internal/reviewrecord/repair_test.go`
- `internal/reviewrecord/source_test.go`
- `internal/reviewrecord/write.go`
- `internal/systemtest/adoption_test.go`
- `internal/systemtest/bench_follow_on_test.go`
- `internal/systemtest/otel_verbs_test.go`
- `internal/systemtest/owner_artifact_recovery_test.go`
- `internal/systemtest/owner_land_lease_test.go`
- `internal/systemtest/owner_land_race_test.go`
- `internal/systemtest/owner_landing_fixture_test.go`
- `internal/systemtest/owner_test.go`
- `internal/testreport/floor.go`
- `internal/worktree/commitment_landing_fixture_test.go`
- `internal/worktree/commitment_light_landing_test.go`
- `internal/worktree/completion_fixture_test.go`
- `internal/worktree/delegated_integration_test.go`
- `internal/worktree/identity_component_test.go`
- `internal/worktree/joins.go`
- `internal/worktree/land.go`
- `internal/worktree/land_bench_home_test.go`
- `internal/worktree/land_broker_notice_test.go`
- `internal/worktree/land_fixtures_test.go`
- `internal/worktree/land_identity_test.go`
- `internal/worktree/land_journey_test.go`
- `internal/worktree/land_lease.go`
- `internal/worktree/land_lease_test.go`
- `internal/worktree/land_local_capture_test.go`
- `internal/worktree/land_marker_fixture_test.go`
- `internal/worktree/land_moved_destination_test.go`
- `internal/worktree/land_override_test.go`
- `internal/worktree/land_prunes_landed_siblings_test.go`
- `internal/worktree/land_refusal.go`
- `internal/worktree/land_refusal_fixture_test.go`
- `internal/worktree/land_specless_test.go`
- `internal/worktree/land_surface_test.go`
- `internal/worktree/land_tickets_only_test.go`
- `internal/worktree/land_trace_test.go`
- `internal/worktree/land_window_test.go`
- `internal/worktree/merge.go`
- `internal/worktree/merge_fold_red_test.go`
- `internal/worktree/merge_refusal.go`
- `internal/worktree/merge_route_test.go`
- `internal/worktree/merge_target_grade_test.go`
- `internal/worktree/parallel_census_test.go`
- `internal/worktree/pool_root_test.go`
- `internal/worktree/refusal_route_follow_test.go`
- `internal/worktree/refusal_route_test.go`
- `internal/worktree/tree_target_test.go`
- `projects/benchkit.md`
- `reviews/build-landing-window.md`
- `specs/build-landing-window/`
- `tests/canary/claude-agent-definitions/agent-unnamed-in-skill`
- `tests/canary/claude-agent-definitions/model-declared`
- `tests/canary/claude-agent-definitions/name-mismatch`
- `tests/canary/claude-agent-definitions/shell-tool-absent`
- `tests/canary/claude-agent-definitions/skill-names-missing-agent`
- `tests/canary/claude-agent-definitions/spawning-tool`
- `tests/canary/claude-agent-definitions/tools-absent`
- `tests/canary/data-handling-derivation/undocumented-passlist-var`
- `tests/canary/docs-currency-token-diet`
- `tests/canary/guidance-prose-budgets`
- `tests/canary/line-routing`
- `tests/canary/package-core-guard/bounds-raw-wait-deadline`
- `tests/canary/package-core-guard/bounds-raw-wait-duration`
- `tests/canary/package-core-guard/bounds-reassigned-wait-duration`
- `tests/canary/package-core-guard/bounds-redeclared-wait-duration`
- `tests/canary/package-core-guard/bounds-worktree-window-unwrapped`
- `tests/canary/load-validity-metadata/readme-shared-rule-drift`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/row-next-grammar/token-table-lacks-kit-edit`
- `tests/canary/skill-description-budgets`
- `tests/canary/skills-index-command-adapters`
- `tests/canary/workflow-guidance-anchors`

## Out of scope

- A `--wait` queue on the lease: 6 edits, 3 gate runs. The refusal and the rerun remain the surviving caller.
- A pty-backed test of the override's real terminal path: 4 edits, 2 gate runs, plus a dependency decision under the dependency standard. The injected predicate remains the surviving caller.
- A landing that names a non-landing gate owner instead of `bench doctor`: 5 edits, 2 gate runs. The rerun remains the surviving caller.
- Signed reviewer authority for the override: 12 edits, 4 gate runs, and a key-management decision. The honest-mistake layers remain the surviving caller.
- An `isatty` predicate in place of the character-device test: 3 edits, 2 gate runs, with a `bench setup` posture change. The typed confirmation remains the surviving caller.

## Further notes

Flagged additions:

- The resolution paths of a recorded fold pass `paths-authorized` (R8).
- A completion proof fault over a current source routes to `land-handback` (R10).
- `bench commitment start` drops a stale gap record of the outcome that it claims.
- `ReconcileDelivered` and `SetBlocker` drop the gap record of the outcome whose claim they drop.
- The override verb runs only from the primary checkout.
- The guard also denies the `bench worktree exec` child form of the verb.
- The gap open writes from the gate verb, and a write failure exits 1 after a green verdict.

Posture changes that red existing tests:

- `TestReviewCheckpointChainGapNamesTheExpectedBase` expects the new chain rule text.
- `foldRedFixture`, which `TestMergeFacesFollowTheirRoutes` follows, stops expecting `FT342` in the fold-red route.
- `TestMergeFacesFollowTheirRoutes` follows the new fold-red route.
- `TestHelpInventoryIsComplete` expects the `override` and `fold` forms and the `--fold` flags.
- `TestIntentReExportsEveryMovedLedgerName` expects the three new ledger types.
- `TestLandingRefusalRegistryHasAProducingFixture` needs producing fixtures for `land-lease-held`, `land-build-window`, and `land-destination-moved`.
- The `data-handling-derivation` check reds until `DATA_HANDLING.md` names the new fields.
- `TestWorktreeLandPublicRaceAndRerun` and the `default-branch-policy` case of `TestCommitmentGateRace` keep their detail text; their route moves to `land-destination-moved`.
- `specLessLandingFixture` and `ticketsOnlyLandingFixture` (`internal/worktree/land_fixtures_test.go:110,124`) bind a staged spec, so their spec-less and tickets-only landings refuse under the window. Ticket 5 changes those fixtures, and their callers in `internal/worktree` and `internal/systemtest` join its `Writes:` line.
- `internal/worktree/parallel_census_test.go` pins the package's exact test count, so each ticket that adds a top-level test in `internal/worktree` raises that pin.

Boundaries:

- `TestLandCommandRefusesPostGateUnknownIgnoredMutation` keeps `land-handback` (BW22).
- `TestLandSuppliesTheDeclaredBenchHomeItself` keeps its `bench doctor` route, because its cause is infrastructure.
- `TestGateRunRefusesInFlightExecution` keeps its gate-lock text.

Source-sentence-to-row table:

- "`bench worktree land` holds one lease in the intent ledger, keyed by the destination branch, from composition through publish" -> BW1, BW6, BW7, BW10.
- "a second landing on the same destination refuses with the owner's assignment id and a route that waits" -> BW2, BW3, BW4.
- "Today the gate lock covers only the gate phase, and its refusal names a PID" -> BW5.
- "`bench worktree land` refuses a landing that is not the open build's own landing" -> BW46, BW50, BW57, BW94.
- "The open-build state lives in the shared intent ledger, under the intent lock" -> BW56, BW60.
- "The refusal names the build and the gap to wait for" -> BW47, BW48, BW49.
- "The window opens at `bench commitment start`" -> BW46, BW69.
- "It closes on completion, on `bench commitment block`, or on abandon" -> BW51, BW52, BW53, BW96, BW97.
- "A gap opens at a recorded chunk checkpoint, and the next chunk dispatch closes it" -> BW60, BW65, BW95.
- "The build folds `main` once, before the author dispatch, so the charge pins the new base" -> BW66, BW67.
- "The build's own light-path dispatches also wait for a gap" -> BW57.
- "An override is a recorded, reviewer-owned flag with a reason" -> BW71, BW81.
- "An agent cannot grant it to itself" -> BW76-BW79, BW82, BW83.
- "A move after the last checkpoint occurs only through the override" -> BW62, BW72.
- "A new record form accepts a fold-delta review" -> BW25, BW34, BW37.
- "A clean fold needs no review round, and only conflict-resolved hunks get review" -> BW26, BW27, BW39, BW40.
- "The final checks include the packages of the merged test files" -> BW41, BW42.
- "Back-out is a revert commit that lands from a fresh worktree with the override. `main` is never reset" -> BW89.
- "The `merge-fold-red` refusal face routes a fold that adds a red to this row" -> BW45.
- "A refusal for a moved destination names `bench worktree merge` and the landing rerun, not `bench doctor`" -> BW17-BW21.
- "Bare `bench preflight build` still grades `base-current` red between chunks and names a `main` merge that the review chain refuses" -> BW34, BW91.
- "After a fold, the final checks need not include the package of each merged test file" -> BW42.
- "`bench record completion` refuses the stale source" -> BW37, BW38.

The pre-review proof checklist follows.

- Cited symbols: each symbol below resolves at b94a60ec.
  - Ledger: `CommitmentState`, `ValidateCommitment`, `readPathTolerant`, and `Transact`.
  - Landing: `PublishAdmitted`, `AdmitPublication`, `commitmentAdmission`, `landAttributed`, `landReviewed`, `landingCompositionRoute`, `Compose`, and `TipMovedError`.
  - Gate: `funnelFace`, `completionProofError`, `AuthorizeWithWriters`, and `RunCommand`.
  - Records: `checkSource`, `checkCompletion`, `checkReviews`, `RecordAssignment`, and `unauthorizedPaths`.
  - Others: `WritesFloor`, `IsTerminal`, `RecordLock`, `ReconcileDelivered`, and `SetBlocker`.
- Import edges: `internal/gate` already imports `internal/commitment/repository` and `internal/intent`. `internal/reviewrecord/recordcmd` already depends on `internal/commitment/repository` and `internal/testreport` through `internal/probe`; its new direct edges make no cycle, because neither package depends on `recordcmd`. `internal/commitment/repository` adds `internal/gocache`, whose dependencies are `bounds`, `toon`, `sanitize`, and `usage`. `internal/reviewrecord` cannot import `internal/testreport`, because `go list -deps ./internal/testreport` names `internal/reviewrecord`.
- Source-row clauses and occurrences: the full bodies of `roadmap/FT342.md` and `roadmap/FT253.md` and their seventeen occurrence lines were read, with the retros in Git history.
- Promised field labels: each label below is exact.
  - Ledger: `leases`, `gaps`, `overrides`, `destination`, `assignment`, `milestone`, `outcome`, `chunk`, `open`, `last`, and `reason`.
  - Record: `folds`, `after`, `from`, `commit`, `incoming`, `resolution`, `resolution_digest`, `test_packages`, and `merged-test:<package>`.
  - Faces: `land-lease-held`, `land-build-window`, and `land-destination-moved`.
- Changed-function callers: `bench consumers --production` at b94a60ec names each caller below.
  - `PublishAdmitted`: `internal/worktree/land.go:308`.
  - `RecordAssignment`: `internal/reviewrecord/recordcmd/assignment.go:40`.
  - `AuthorizeWithWriters`: `internal/gate/authorization/authorization.go:145` and `internal/landing/landing.go:117`.
  - `checkSource`: `internal/reviewrecord/check.go:20`, `completion.go:37`, and `coverage.go:16`.
  - `checkCompletion`: `internal/reviewrecord/coverage.go:144`.
  - `RunCommand`: `cmd/bench/main.go:151` and `internal/gate/gate.go:225`.
  - `landingCompositionRoute`: `internal/worktree/land.go:210`.
  - `ReconcileDelivered`: `internal/worktree/joins.go:61`.
- Copy survival: BW45 reds a surviving FT342 question. BW84 reds a hand-copied help line. BW81 reds a second reason rule. BW41 reds a second path-to-package rule.
- Rendered-shape readers: the old text "a fold in the last chunk delta" sits in `internal/reviewrecord/coverage.go:86` and `internal/gate/review_checkpoint_commits_test.go:49`. "the open FT342 question" sits in `internal/refusalroute/faces_merge.go:26` and the `foldRedFixture` names at `internal/worktree/merge_route_test.go:236`. The slicing runs `rg` for each and for the guard manifest line.
- Pin operators: the rows compare exact detail substrings, exact `next=` tokens, exact ledger records, and exact record fields.
- Entry reads: the lease reads the Git common directory through the intent address owner. The override reads stdin through an injected predicate and reader.
- Derived expectations: expected rows derive through `toon.Table`, and BW41 derives its expected list through `testreport.WritesFloor`.
- Consolidated rules: the override reason reuses the block verb's reason rule (BW81). The moved-destination face replaces three routes:

  | site | old route | new route | row |
  | --- | --- | --- | --- |
  | completion proof over a moved destination | `land-infrastructure` | `land-destination-moved` | BW17 |
  | post-gate ref move | `land-handback` | `land-destination-moved` | BW20 |
  | compare-and-swap refusal | `land-handback` | `land-destination-moved` | BW21 |

- Quantified obligations: "every landing that is not the build's own" covers spec-less, tickets-only, and spec landings (BW46, BW57). "each open build" covers two parallel builds (BW58).
- Workflow-step writes: the gap open writes the ledger, which no gate digest reads. The gap close writes the ledger before the plan fence, which the plan digest reads, and BW67 proves that a refusal leaves the fence unchanged.

Sources not re-read: three learning texts, because Git ignores `capture/learnings.md`.
They are the ft336 `base-current` occurrence, the adoption golden-fixture `bench doctor` detail, and the 2026-08-31 liveness-sentinel wording.
Unknown: the Codex shell's terminal state.
`bench commit` does not grade a hand-resolved fold commit against the spec fence. The production readers of `unauthorizedPaths` are in `internal/preflight` and `internal/worktree/land_identity.go`, and the commit lane lives in `internal/gate/lane.go`. The landing grades the fence after the record holds the fold entry, so R8 reaches it through `internal/preflight/gather.go:187`.
