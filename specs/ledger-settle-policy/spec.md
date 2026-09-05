# Pure policy behind the intent ledger and the composition

Status: staged

Roadmap: FT302

Decision source: the named reviewed artifact `roadmap/FT302.md` (candidates 3 and 4 paragraph, restored 2026-09-05 at 30b00d59), with the survey report `capture/architecture-review-20260905T112705.html` as its source evidence

Verification log: 1 iteration(s) to accept — the round on opus at medium effort returned seven findings and two ticket-graph blockers inside the cap of one. The author folded four spec findings: the LS6, LS22, and LS33 citations, and the LS35 flagged addition. The slicer folded three ticket findings: the 02→01 edge, the 07 split, and the 07 Writes line. The reviewer's sign-off is the acceptance.

## Problem

Two owners bury a decision inside an effect.

Seven intent mutators re-type one envelope. `Upsert`, `Compact`,
`PutCleanupReceipt`, `PutAssignment`, `ReauthorizeAssignment`,
`PurgeAssignments`, and `DeleteAssignment` each resolve the ledger address,
acquire the lock, read the file, decide, and replace the file. So no admission
rule is testable without a repository and a lock. The seven copies also carry
the lock order, the stale-lock rule, and the atomic replacement seven times.
`internal/intent/assignment.go` is 584 lines against a 400-line budget.

The composition has the same shape. `resolveCaptureConflict` holds the settle
precedence, the mode agreement, and the Git tree edit in one function. Seven
tests build real repositories to reach a classification a mode list decides.
The refusal reason is a `false` return inside that function, so no caller can
name why a settle refused. `internal/landing/composition_test.go` is 585 lines
against the same budget.

## Solution

One ledger transaction applies a pure ledger-to-ledger decision under the lock.
It carries a read mode, so the tolerant purge runs inside it. It carries an
optional compensation step, so the reauthorization keeps its external rollback.
Each mutator becomes its decision.

The ledger types, the validators, and `RequestDigest` move to a leaf package.
The intent package re-exports them as aliases, so no consumer file moves. A
policy child holds the admission rules, and the compaction takes typed liveness
facts the adapter translates.

A second policy child under `internal/landing` answers the settle verdict per
path from stage records. The adapter owns Git alone. The refusal reason
surfaces on `Conflict`, so a conflicted landing names why the settle refused.
`CONTEXT.md` gains the terms **ledger transaction** and **settle verdict**.

## User stories

### Group A — the ledger transaction

Line: opus / medium. The transaction owns the lock order and the atomic
replacement for seven mutators. The scorecard routes correctness-critical owner
logic to the mid tier at medium effort.

1. As a mutator author, I want one transaction that locks, reads, runs my closure, and writes, so that I never re-type the envelope.
2. As a mutator author, I want a write only on a reported change, so that an identical write preserves the ledger bytes.
3. As a mutator author, I want the transaction to release the lock on every exit, so that a refusing closure never strands the lock.
4. As a mutator author, I want a closure error to leave the ledger bytes untouched, so that a refusal persists nothing.
5. As the purge, I want a tolerant read mode that drops each unreadable record alone, so that an older ledger still clears.
6. As the strict readers, I want the strict mode to keep refusing a ledger it cannot account for, so that no unaccounted record authorizes anything.
7. As the reauthorization, I want an optional compensation step that runs when the write fails, so that my external transition rolls back.
8. As every other mutator, I want no compensation step, so that only the reauthorization pays for one.
9. As a concurrent writer, I want the transaction to reclaim a stale lock and keep every entry, so that two writers do not lose work.
10. As the purge, I want an absent ledger to answer before the lock, so that a converged re-run stays a no-op.
11. As a reviewer, I want each of the seven mutators to keep its current external behavior, so that the deepening preserves behavior.

### Group B — the leaf ledger package

Line: opus / low. Each move is a mechanical relocation at a known seam under
the covering suite, which is the scorecard's low-effort shape.

12. As a maintainer, I want the types, the validators, and `RequestDigest` in a leaf package, so that the schema has one pure owner.
13. As a consumer, I want the intent package to re-export every moved name as an alias, so that no consumer file changes an import.
14. As a maintainer, I want the leaf package to import nothing under `internal/`, so that no importer of it can close a cycle.
15. As a file owner, I want `internal/intent/assignment.go` to end under its 400-line budget, so that the structure check turns green.

### Group C — the admission policy

Line: opus / medium. The admission rules decide what the ledger authorizes, and
the scorecard routes owner-decision logic to the mid tier at medium effort.

16. As a maintainer, I want the admission rules of the seven mutators in a policy child, so that each rule is testable without a repository.
17. As the compaction, I want typed liveness facts for its three probes, so that the adapter owns every Git and filesystem read.
18. As a test author, I want each admission partition as a policy table case, so that a partition costs no repository and no lock.
19. As a reviewer, I want one real-filesystem journey per mutator to stay, so that the adapter's wiring keeps an end-to-end proof.
20. As a maintainer, I want the policy child to import neither `internal/git` nor `internal/intent`, so that the effect boundary holds.

### Group D — the settle policy

Line: opus / medium. The settle precedence decides what a landing publishes,
and the scorecard routes owner-decision logic to the mid tier at medium effort.

21. As a maintainer, I want a policy child under `internal/landing` with its own stage-record type, so that the settle decision reads no Git.
22. As a maintainer, I want the merge-tree parser to stay adapter translation, so that Git's output format has one reader.
23. As a maintainer, I want the mode-list classifier in the policy, so that a conflict kind is decided from a mode list alone.
24. As a maintainer, I want the capture rule table in the policy, so that the side, the union, and the removal come from one source.
25. As a test author, I want each settle partition as a policy table case, so that a partition costs no repository.
26. As a reviewer, I want the seven real-Git classification journeys to stay, so that the parser keeps its representative proof.
27. As a reviewer, I want one settle journey per verb and one refusal journey, so that the tree edit keeps an end-to-end proof.
28. As a file owner, I want `internal/landing/composition_test.go` to end under its 400-line budget, so that the structure check turns green.
29. As a file owner, I want no new file inside `internal/landing/` itself, so that the crowded-directory count does not grow.

### Group E — the surfaced refusal reason

Line: opus / medium. The refusal text is observable output every landing
consumer reads, and the scorecard routes an output change to medium effort.

30. As a landing operator, I want a refused settle to name the reason and the paths, so that repair starts from the cause.
31. As a landing operator, I want a path outside the capture table named as its own reason, so that I repair the path by hand.
32. As a landing operator, I want a non-regular mode named as its own reason, so that a symlink or a gitlink under `capture/` reads plainly.
33. As a landing operator, I want a mode disagreement named as its own reason, so that a permission change reads plainly.
34. As a landing operator, I want union content Git cannot merge as text named as its own reason, so that a binary journal reads plainly.
35. As a landing operator, I want one reason when two compete, so that the message stays one predicate.
36. As a consumer of the refusal, I want the kind to stay at the front of the text, so that every existing match holds.
37. As a reviewer, I want the old bare-kind text excluded when a reason exists, so that a build cannot ship the reason unread.

### Group F — the purity censuses and the glossary

Line: opus / medium. Guidance prose carries the leverage override, and the
reviewer's 2026-08-26 rule caps a subagent at medium effort.

38. As a maintainer, I want each new package to carry its own source census wrapper, so that the purity boundary is executable.
39. As a maintainer, I want the policy census to forbid the two new parent adapters, so that a policy child cannot import its parent.
40. As a cold session, I want `CONTEXT.md` to define **ledger transaction** with its Avoid list, so that the vocabulary does not drift.
41. As a cold session, I want `CONTEXT.md` to define **settle verdict** with its Avoid list, so that the vocabulary does not drift.

## Implementation decisions

- The intent package gains one exported transaction. It resolves the ledger address and acquires the lock file. It then reads the ledger, calls the caller's closure with it, and writes on a reported change. The closure returns the next ledger, a changed report, and an error. The transaction releases the lock on every exit. It performs the same atomic replacement `writePath` performs today.
- The transaction lives in its own file inside the intent package. `internal/intent/intent.go` is 395 lines against a 400-line budget, so the transaction cannot land there.
- The transaction takes a read mode. The strict mode is today's `readPath`, which refuses a ledger it cannot account for. The tolerant mode is today's purge read. It decodes each assignment as a raw value and drops the records it cannot decode or validate. It reads the entries and the cleanup receipts separately. `PurgeAssignments` is the one tolerant caller, and it keeps its pre-lock absent-file answer.
- The transaction takes an optional compensation step. The step runs when the write fails. `ReauthorizeAssignment` is the one caller that supplies one, and it supplies the rollback its `transition` returns. Every other mutator supplies none.
- The seven mutators become decisions. `Upsert`, `Compact`, `PutCleanupReceipt`, `PutAssignment`, `ReauthorizeAssignment`, `PurgeAssignments`, and `DeleteAssignment` each call the transaction with their closure. Their exported signatures and their external behavior do not change.
- A new leaf package holds the ledger schema. It holds `Ledger`, `Entry`, `Kind`, `Assignment`, `AssignmentState`, `Recovery`, and `CleanupReceipt`, with their constants. It also holds the schema constants and `RequestDigest`. It holds `ValidateAssignment`, `ValidIdentity`, the entry validator, and the cleanup-receipt validator. It holds the identity patterns and the branch-ref helpers. The package imports nothing under `internal/`.
- The intent package re-exports every moved name. A type becomes a type alias, and a function and a constant become a declared value. So every importing file keeps its `intent.` spelling and its import line.
- A new policy child holds the admission rules of the seven mutators. Each rule answers from the supplied ledger and the mutator's operand alone. The compaction rule takes typed liveness facts: one candidate-probe boolean, one worktree-existence answer per entry, and one landed answer per branch. The intent package translates `git.Worktrees`, `git.LocalBranches`, `os.Stat`, `git.ResolvedDefault`, and `git.LandedInDefault` into those facts at its boundary.
- A new policy child under `internal/landing` holds the settle decision. It owns its own stage-record type, the capture rule table, the mode-agreement rule, the mode-list conflict classifier, and the settle verdict per path. It answers a verdict of side, union, remove, or refuse, and a refusal carries its reason.
- The merge-tree parser stays in the adapter. The adapter parses `merge-tree --write-tree -z` into the policy's stage-record type, calls the policy, and applies the verdict with `update-index`. The union verdict's text merge stays in the adapter, because it shells to `git merge-file`.
- `Conflict` gains a settle-refusal reason field. The four reasons are `path outside the capture table`, `non-regular mode`, `mode disagreement`, and `union content not text`. `ConflictError.Error` reads exactly `composition conflict: <kind>; settle refused: <reason> (<paths>)` when a reason exists, and it keeps today's `composition conflict: <kind>` when none exists. The paths are the refusing paths, comma-separated in merge-tree order.
- One reason wins. The policy scans the stage records in merge-tree order and answers the first refusal it finds. Within one record the table check precedes the mode check, so `path outside the capture table` beats `non-regular mode`. The record scan precedes the settle scan, so both record reasons beat `mode disagreement` and `union content not text`.
- `internal/puritycensus` gains `internal/intent` and `internal/landing` to the policy package's forbidden-import pattern. The pattern is the one source for the policy boundary, so the two new children and the three existing children grade against one rule.
- Each new package carries a `purity_census_test.go` wrapper. The two policy children run `puritycensus.PolicyPackage()`, and the leaf package runs `puritycensus.LeafPackage()`. Each wrapper names its own package source through `MustHold`.
- The two policy children and the leaf package are new directories. No new file lands in `internal/landing/` itself, because that directory already holds 16 source files against a 12-file budget.
- `CONTEXT.md` gains **ledger transaction**. The term names the one envelope over the intent ledger. It locks, reads in a named mode, runs one decision closure, and writes a reported change. Its Avoid list holds `lock helper` and `mutator envelope`.
- `CONTEXT.md` gains **settle verdict**. The term names the policy's answer for one conflicted path: a side, a union, a removal, or a refusal with a reason. Its Avoid list holds `resolution` and `merge result`.
- The exit proof is a differential run: the pre-existing suite passes with test logic unchanged, except the named moved partitions and the surfaced-reason row. The named exceptions are the admission partitions that become policy table cases, the settle partitions that become policy table cases, and the new reason assertions.

## Testing decisions

- A good transaction test drives a real ledger file under a real lock and observes the file bytes, the file time, and the lock file. The prior art is `TestConcurrentWritersKeepEveryEntryAndStaleLockReclaims` and `TestIdenticalStampedAssignmentWritePreservesBytes`.
- A good admission test drives the policy with a literal ledger value and observes the returned ledger and the changed report. It builds no repository and takes no lock. The prior art is `internal/worktree/landingpolicy/landingpolicy_test.go`.
- A good settle test drives the policy with a literal stage-record slice and observes the verdict per path. The prior art is the same policy test file.
- Each mutator keeps one real-filesystem journey, and the composition keeps its seven classification journeys plus one settle journey per verb and one refusal journey.
- The gate's structure phase observes the two budget stories. The conformance phase observes nothing new, because no standing check enumerates census packages.

### Seam diagram

    intent mutator (Upsert, Compact, PutCleanupReceipt, PutAssignment,
                    ReauthorizeAssignment, PurgeAssignments, DeleteAssignment)
        │
        ▼
    root, read mode, decision closure, optional compensation
        │
        ▼
    [ ledger transaction: address, lock, stale lock, read, decide, atomic replace ]
        │                     ◀ tests attach here: a real ledger file and lock
        ▼
    ledger  ──▶  [ admission policy (pure) ]  ──▶  next ledger, changed
                     ◀ tests attach here: literal ledger values in a table

    bench worktree land | bench worktree merge
        │
        ▼
    merge-tree -z output  ──▶  [ composition adapter: parse, apply, union merge ]
                                   │
                                   ▼
                          stage records  ──▶  [ settle policy (pure) ]  ──▶  verdict: side | union | remove | refuse(reason)
                                                  ◀ tests attach here: literal stage records in a table
                                   │
                                   ▼
                          tree | `composition conflict: <kind>; settle refused: <reason> (<paths>)`
                                   ◀ tests attach here: real-Git journeys and the refusal text

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LS1 | 1, 11 | each of the seven mutators produces the same ledger file bytes over its own happy-path operand as the pre-transaction implementation produces | a new test `TestLedgerTransactionPreservesEveryMutatorOutcome` in `internal/intent`, which runs the seven mutators over one fixture ledger and compares each result with the recorded pre-change bytes | a transaction that drops one mutator's decision writes different bytes |
| LS2 | 2 | an `Upsert` of an identical entry and a `PutAssignment` of an identical assignment each leave the ledger bytes and the file modification time unchanged | `internal/intent/intent_test.go` (`TestIdenticalStampedAssignmentWritePreservesBytes`) plus a new `Upsert` case in that test | a transaction that always writes replaces the file and moves the modification time |
| LS3 | 3, 4 | with a closure that returns an error, the transaction returns that error, the ledger file bytes are unchanged, and the lock file is absent | a new test `TestLedgerTransactionRefusalPersistsNothingAndReleasesTheLock` in `internal/intent` | a transaction that writes before the closure returns changes the bytes, and a transaction that returns before the release leaves the lock file |
| LS4 | 4 | with a closure that reports a change and a ledger directory made read-only, the transaction returns a write error and the ledger file bytes are unchanged | a new test `TestLedgerTransactionTerminalWriteFailureKeepsThePreviousBytes` in `internal/intent` | a transaction that writes in place instead of replacing atomically truncates the file |
| LS5 | 5 | with a ledger holding one valid and one undecodable assignment record, `PurgeAssignments` drops the undecodable record and keeps the valid one | `internal/intent/purge_test.go` (`TestPurgeAssignmentsDropsRecordsReadRefuses`) | a purge moved to the strict mode refuses the whole file and drops nothing |
| LS6 | 6 | with a ledger holding an unsupported schema number, `Read` returns an error and no ledger | `internal/intent/intent_test.go` (`TestReadEvidenceStates`) | a strict reader moved to the tolerant mode returns an empty ledger |
| LS7 | 7 | with a transition that returns a rollback and a ledger directory made read-only, `ReauthorizeAssignment` returns an error and the rollback ran exactly once | a new test `TestReauthorizeCompensatesOnTheTerminalWriteFailure` in `internal/intent` | a transaction that drops the compensation step leaves the external transition applied |
| LS8 | 8 | with a ledger directory made read-only, `PutAssignment` returns an error, and the transaction runs no compensation step | the same new test, second case | a transaction that requires a compensation step panics on a nil step |
| LS9 | 9 | two concurrent writers each keep their entry, and a lock file left by a dead process is reclaimed | `internal/intent/intent_test.go` (`TestConcurrentWritersKeepEveryEntryAndStaleLockReclaims`) | a transaction that drops the stale-lock reclaim times out |
| LS10 | 10 | with no ledger file present, `PurgeAssignments` returns zero, and no lock file is created | `internal/intent/purge_test.go` (`TestPurgeAssignmentsIsANoOpWithoutALedger`) | a purge that acquires the lock first creates a lock file beside an absent ledger |
| LS11 | 12, 13 | every moved name resolves through its `intent.` spelling and equals the leaf package's declaration | a new test `TestIntentReExportsEveryMovedLedgerName` in `internal/intent`, which assigns each alias to the leaf value | a move that renames a type reds compilation, and a move that drops one alias reds this test |
| LS12 | 13 | the whole module builds with no import of the leaf package outside `internal/intent`, `internal/intent/admissionpolicy`, and their tests | a new test `TestOnlyTheIntentOwnerImportsTheLeafLedger` in `internal/intent` | a move that rewrites consumer imports adds an import edge this test names |
| LS13 | 14 | the leaf package's own directory holds no import under `internal/` other than `internal/puritycensus`, no ambient effect, and no `t.Parallel` call | a new census wrapper test `TestPurePackageSourceCensus` in `internal/intent/ledger`, which runs `puritycensus.LeafPackage()` | a leaf that imports `internal/git` for the address reds the forbidden-import diagnostic |
| LS14 | 15 | `bench structure` reports `internal/intent/assignment.go` under its 400-line budget | the gate's structure phase | a move that leaves the schema in place keeps the file at 584 lines |
| LS15 | 16 | for each of the seven mutators, the admission policy answers the next ledger and the changed report from a literal ledger value with no repository present | a new table test `TestAdmissionPolicyDecidesEveryMutator` in `internal/intent/admissionpolicy` | a rule left in the adapter cannot be reached without a repository, so this test cannot compile |
| LS16 | 17 | the compaction rule keeps an entry whose worktree exists, drops an entry whose worktree is absent, drops an entry whose branch landed, and keeps an uncorrelated Claude entry when the candidate probe is true | a new table test `TestCompactionRuleReadsTypedLivenessFacts` in `internal/intent/admissionpolicy` | a rule that calls `os.Stat` reds the census, and a rule that ignores the candidate probe keeps the wrong entry |
| LS17 | 17 | `Snapshot` and `Compact` over a real repository return the same live entries as before the change for a worktree entry, a landed-branch entry, and an uncorrelated Claude entry | `internal/intent/intent_test.go` (`TestSnapshotProofLifecycle` and `TestUncorrelatedEntriesUseCandidateSet`) | an adapter that translates the landed fact for the wrong branch keeps a landed entry |
| LS18 | 18, 19 | each mutator keeps exactly one real-filesystem journey, and every other admission partition is a policy table case | review-owned: the reviewer reads the seven journeys and the two policy table tests | the suite passes either way, so no check counts journeys |
| LS19 | 20 | the admission policy's own directory holds no import of `internal/git`, `internal/intent`, `internal/worktree`, `internal/bounds`, `os/exec`, or `syscall`, no ambient effect, and no `t.Parallel` call | a new census wrapper test `TestPurePackageSourceCensus` in `internal/intent/admissionpolicy`, which runs `puritycensus.PolicyPackage()` | a policy that imports its parent for the ledger type reds the forbidden-import diagnostic |
| LS20 | 21, 24 | the settle policy answers `source` for `capture/session-handoff.md`, `union` for `capture/learnings.md` and `capture/IDEAS.md`, `destination` for another `capture/` path, and a refusal for `capture.md` | a new table test `TestSettlePolicyAnswersTheCaptureRule` in `internal/landing/settlepolicy` | a table that keys on the bare word `capture` settles `capture.md` |
| LS21 | 21 | the settle policy answers a removal for a path whose winning stage is absent | the same new table test, one case | a policy that answers the other side republishes a deleted file |
| LS22 | 22 | `parseConflict` turns a `merge-tree -z` output into the policy's stage records, and it returns an error for empty output and for a malformed record | `internal/landing/composition_test.go` (`TestComposeClassifiesRealGitConflictsWithoutMutation` and `TestConflictKindRejectsEmptyMergeTreeOutput`) plus a new malformed-record case; the first test reaches the parser through the composition owner for seven real outputs | a parser moved into the policy makes the policy read Git's format, which reds LS26 |
| LS23 | 23 | the conflict classifier answers `gitlink` for a mode list holding `160000`, `symlink` for one holding `120000`, `mode` for two unequal ordinary modes, and `textual` for one repeated mode | a new table test `TestConflictKindClassifiesModeLists` in `internal/landing/settlepolicy` | a classifier that checks the ordinary modes before the special ones answers `mode` for a gitlink list |
| LS24 | 25, 27 | the composition keeps one settle journey per verb and one refusal journey, and every other settle partition is a policy table case | review-owned: the reviewer reads `TestComposeSettlesPhaseOwnedConflictsByRule` and `TestComposeRefusesConflictsTheRuleTableCannotSettle` against the policy tables | the suite passes either way, so no check counts journeys |
| LS25 | 26 | over real Git conflicts, `Compose` answers the kinds `textual`, `modify/delete`, `rename/rename`, `file/directory`, `mode`, `symlink`, and `gitlink` without mutating the repository | `internal/landing/composition_test.go` (`TestComposeClassifiesRealGitConflictsWithoutMutation`) | a move that drops one kind from the parser reds that case |
| LS26 | 21 | the settle policy's own directory holds no import of `internal/git`, `internal/landing`, `internal/worktree`, `internal/bounds`, `os/exec`, or `syscall`, no ambient effect, and no `t.Parallel` call | a new census wrapper test `TestPurePackageSourceCensus` in `internal/landing/settlepolicy`, which runs `puritycensus.PolicyPackage()` | a policy that shells to `git merge-file` for the union reds the forbidden-import diagnostic |
| LS27 | 29 | `bench structure` reports no crowded-directory growth for `internal/landing/` | the gate's structure phase | a policy placed in a new file inside `internal/landing/` raises the count to 17 |
| LS28 | 30, 31 | with a conflict on `capture/learnings.md` beside a conflict on `named`, `ConflictError.Error` reads `composition conflict: textual; settle refused: path outside the capture table (named)` | a new test `TestConflictErrorNamesTheSettleRefusalReason` in `internal/landing` | a build that stores the reason without rendering it returns the bare kind |
| LS29 | 32 | with a symlink at `capture/learnings.md` on one side, `ConflictError.Error` holds `settle refused: non-regular mode` and names `capture/learnings.md` | the same new test, second case | a policy that checks the table alone answers the table reason for a symlink |
| LS30 | 33 | with a permission change on `capture/learnings.md` on one side and a content change on the other, `ConflictError.Error` holds `settle refused: mode disagreement` | the same new test, third case | a policy that publishes one side's mode loses the other side's mode |
| LS31 | 34 | with binary content on both sides of `capture/learnings.md`, `ConflictError.Error` holds `settle refused: union content not text` | the same new test, fourth case | a policy that answers the destination side publishes a half-merged binary |
| LS32 | 35 | with `capture.md` conflicting beside a symlinked `capture/learnings.md`, `ConflictError.Error` names exactly one reason, and that reason is `path outside the capture table` | the same new test, fifth case | a policy that answers the last refusal it finds names `non-regular mode` |
| LS33 | 36 | the landing refusal surface still holds the text `composition conflict: <kind>` at the front and the conflicted path, for the kinds `textual` and `mode` | `internal/worktree/land_surface_test.go` (`TestLandCommandConflictRefusalNamesThePath`), which asserts the `textual` kind, and `internal/worktree/merge_test.go` (`TestMergeRefusesACaptureAddAddWithDisagreeingModes`), which asserts the `mode` kind | a reason placed before the kind breaks every existing prefix match |
| LS34 | 37 | with a refusing settle, `ConflictError.Error` does not equal `composition conflict: <kind>`, and it holds the separator `; settle refused: ` | the same new test, an assertion on every case | a build that renders the reason only in the typed field leaves the old text whole |
| LS35 | 38 | each of the three new packages holds a `purity_census_test.go` whose census names that package's own source | review-owned: the reviewer reads the three wrapper files | no standing check enumerates census packages, so an omitted wrapper passes the gate |
| LS36 | 39 | the policy census reports a forbidden import for a source that imports `internal/intent` and for a source that imports `internal/landing` | `internal/puritycensus/census_test.go` (`TestCensusRefusesTheIntentParentAdapter` and `TestCensusRefusesTheLandingParentAdapter`), two new cases through `diagnose` with `PolicyPackage()` | a pattern left unchanged lets a policy child import its parent adapter |
| LS37 | 40 | `CONTEXT.md` holds the term `ledger transaction` with its Avoid list naming `lock helper` and `mutator envelope` | review-owned: the reviewer reads the entry | the prose mechanics check grades sentences, not terms |
| LS38 | 41 | `CONTEXT.md` holds the term `settle verdict` with its Avoid list naming `resolution` and `merge result` | review-owned: the reviewer reads the entry | the prose mechanics check grades sentences, not terms |
| LS39 | 28 | `bench structure` reports `internal/landing/composition_test.go` under its 400-line budget | the gate's structure phase | a move that leaves the settle tables in place keeps the file at 585 lines |

### Edge inventory

- Error paths: a closure error persists nothing (LS3). A terminal write failure keeps the previous bytes (LS4) and runs the compensation when one exists (LS7). A lock timeout returns today's `lock intent ledger: timed out` text, which no row changes.
- Fail-closed cleanup: the transaction removes the lock file on every exit, including the error exits (LS3). The compensation runs once, and a second failure inside it does not run it again (LS7).
- Empty input: an absent ledger file reads as an empty ledger under the strict mode, and it answers no work under the purge (LS10). An empty conflict record set keeps today's bare kind. LS34 asserts the separator on a refusing settle alone.
- Boundary values: `capture.md` is the prefix boundary of the capture table (LS20). The 256-receipt window is unchanged, and `TestCleanupReceiptWindowKeepsExactlyLast256Completions` observes it.
- Hostile paths: a capture path holding a space settles by the table (LS20 keeps the existing `path-with-a-space` case). A control byte in a conflicted path renders through the existing landing refusal surface, which `internal/worktree/merge_test.go` observes at its `boardfile.md` case.
- Currency: the transaction reads the ledger inside the lock, so no read is stale.
- Re-run idempotency: an identical mutation writes nothing (LS2), and a purge over converged state writes nothing (LS10).
- Partial implementation: a build that adds the transaction and leaves one mutator on the old envelope reds LS1. A build that moves the schema and drops one alias reds LS11. A build that adds the reason field and never renders it reds LS34.
- Audience: every behavior here serves this repository. The two verbs `bench worktree land` and `bench worktree merge` ship to every repository that links the kit, so LS28 to LS34 name that audience.
- Absent versus empty for each directory read: the composition reads no directory. The censuses read their own package directory, and an empty directory fails the census with `census found no Go sources`.
- Package-variable swaps: neither package declares a swappable package variable. The landing tests swap three `Owner` struct fields in the test process. `Compose` runs in that same process, so the swap reaches it.
- Transaction verification failures: persistence before the oracle runs (LS3), interruption inside the oracle (LS3), and persistence at the terminal step (LS4, LS7, and LS8).

**Won't handle** — a standing check that requires a census wrapper in every pure package — the three existing census packages carry the same review-owned obligation today. `internal/worktree/landingpolicy` stays its surviving caller.

**Won't handle** — the dead `changed` variable in `Upsert` — the transaction's changed report replaces it, and the identical-entry early return stays its surviving caller.

**Won't handle** — a union merge inside the settle policy — the census forbids the `git merge-file` effect, so the adapter keeps the text merge. `unionStages` stays its surviving caller.

**Won't handle** — splitting `internal/landing/` below its 12-file directory budget — the directory is already crowded, and this spec only promises no growth. `internal/landing/landing.go` stays its surviving caller.

**Won't handle** — a machine-readable reason field on the rendered refusal — the refusal renders `conflict.Error()` through one detail field, and `landingConflictRefusal` stays its surviving caller.

## Ownership fences

This section is derived from the reader sweep; the ticket slice confirms it.

- `specs/ledger-settle-policy/`
- `reviews/ledger-settle-policy.md`
- `internal/intent/intent.go`
- `internal/intent/transaction.go` (new)
- `internal/intent/transaction_test.go` (new)
- `internal/intent/assignment.go`
- `internal/intent/validate.go`
- `internal/intent/branches.go`
- `internal/intent/worktree_owner.go`
- `internal/intent/intent_test.go`
- `internal/intent/purge_test.go`
- `internal/intent/assignment_lookup_test.go`
- `internal/intent/worktree_owner_test.go`
- `internal/intent/ledger/` (new)
- `internal/intent/ledger_aliases.go` (new)
- `internal/intent/ledger_aliases_test.go` (new)
- `internal/intent/admissionpolicy/` (new)
- `internal/worktree/lifecyclepolicy/lifecyclepolicy.go`
- `internal/worktree/lifecyclepolicy/lifecyclepolicy_test.go`
- `internal/landing/composition.go`
- `internal/landing/merge.go`
- `internal/landing/landing.go`
- `internal/landing/composition_test.go`
- `internal/landing/merge_test.go`
- `internal/landing/landing_helpers_test.go`
- `internal/landing/settlepolicy/` (new)
- `internal/puritycensus/census.go`
- `internal/puritycensus/census_test.go`
- `CONTEXT.md`
- `internal/worktree/land_surface_test.go` — closure headroom only
- `internal/worktree/merge_test.go` — closure headroom only
- `internal/worktree/land_journey_test.go` — closure headroom only
- `internal/worktree/identity_component_test.go` — closure headroom only
- `internal/worktree/land_refusal.go` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-decision-map-term` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary` — closure headroom only
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift` — closure headroom only

A closure headroom entry creates no blocker edge and takes no edit. The five
worktree files hold the prefix matches LS33 keeps. The eight fixture
directories close the pins on `CONTEXT.md`. The eighth is the token-diet
vocabulary fixture. `bench anchors CONTEXT.md` names seven required needles,
and none of them sits in the two new entries.

## Out of scope

- A standing conformance check that enumerates every pure package and requires its census wrapper and its registry row. 5 edits, 1 gate run.
- Splitting `internal/landing/` into modules under its 12-file directory budget. 12 edits, 2 gate runs.
- A typed settle-reason field on the rendered landing refusal, which the disclosure and the TOON surface would read. 6 edits, 2 gate runs.
- Named readers for the peeled-commit and porcelain-status Git facts. FT302 holds them, and each needs a fresh census. 6 edits, 1 gate run.

## Further notes

Flagged additions beyond the decision source:

- The `internal/puritycensus` pattern change. The source puts both children under the policy census. Today's pattern names `internal/worktree` alone, so a new child could import its parent.
- The no-new-file rule for `internal/landing/` itself. `bench structure` reports that directory crowded at 16 files against a 12-file budget, so a new file there worsens a standing red.
- The winning-reason order. The source names four reasons and one message shape, and it does not say which wins when two compete. `craft-spec`'s map discipline requires an ordering row.
- The census wrapper in each new pure package. The three existing pure packages carry the wrapper, so a new pure package without one breaks the tree convention.
- The two structure-budget stories. The source's exit proof does not name them, and `bench structure` reports both files over the 400-line budget today.

Build decisions recorded for reviewer veto:

- The package names are `internal/intent/ledger`, `internal/intent/admissionpolicy`, and `internal/landing/settlepolicy`. The tickets fix one spelling so that sibling tickets compile together.
- The compensation step runs on the write failure alone. Today `ReauthorizeAssignment` also runs its rollback when the compare-and-swap fails. Under the transaction the compare-and-swap sits inside the closure, so the closure's own error path must run the rollback to keep today's behavior. LS7 pins the write-failure arm, and `TestCompareAndSwapRequestDigestRefusesConcurrentMovementAndPreservesOtherFields` pins the other arm.
- The seven mutators keep their exported signatures. A signature change would move consumer files, which the source forbids.

Build decisions recorded during the build (2026-09-05, `--full` run) for reviewer veto:

- The `--full` invocation from the handoff's Next command is read as the sign-off, with the four recommended answers to the reviewer questions above.
- `internal/worktree/lifecyclepolicy` imports `internal/intent` for schema types alone, so the widened policy census reds its wrapper. It moves onto `internal/intent/ledger`, and the LS12 allowlist gains that package. The fence gains its two Go files.
- The LS36 cases are two named tests, `TestCensusRefusesTheIntentParentAdapter` and `TestCensusRefusesTheLandingParentAdapter`, because the cited test counts three diagnostics over one source.
- The LS11 and LS12 tests live in a new file `internal/intent/ledger_aliases_test.go`, because `intent_test.go` sits at 397 lines. The aliases live in `internal/intent/ledger_aliases.go`.
- A read-only ledger directory set before the call fails at the lock, not at the write. LS4, LS7, and LS8 set the directory read-only inside the closure, after the lock exists, and restore it in the compensation or the test cleanup.
- The tolerant read decodes the entries and the cleanup receipts on every read. The old purge decoded them only before a write. So a file whose assignments parse and whose entries do not now returns an error, where it returned zero with nothing to drop.
- The composition keeps two refusal journeys, not one. The policy refusal is `code-path-beside-a-capture-path`. The adapter's own union text-merge refusal is `binary-union-path`, which no policy table can cover. It gains a fourth settle journey, `removal`, because the coordinator's probe showed the removal verdict had no journey.
- A settle refusal exists only when the policy is engaged, that is, when at least one stage record's path has a capture rule. A conflict whose paths all sit outside the table keeps today's bare `composition conflict: <kind>`. `TestLandCommandPublicConflictRepairRequiresNewReviewedTip` pins `detail=composition conflict: textual,next=` for a plain `owned.txt` conflict, and the fence names that file as closure headroom. Ticket 08's fence gains `internal/landing/settlepolicy/` for the engagement rule and its table case.
- `compareAndSwapRequestDigest` moved to the admission policy as `CompareAndSwapRequestDigest`. Its intent test moved with it into two policy table cases, so the cited test name no longer exists in `internal/intent`. Ticket 05's fence gains `internal/intent/admissionpolicy/admissionpolicy_test.go` for that move.
- `PutCleanupReceipt` validates the receipt inside its rule, under the lock, where it validated before the lock. The error text does not change.
- `PutAssignment` exposes no hook that runs under the lock. So LS8 observes the nil compensation through a transaction over the `PutAssignment` rule on a read-only directory. It then drives `PutAssignment` itself over the same directory for the error.
- `DeleteAssignment` had no real-filesystem journey before. The LS1 test supplies its journey.

Source-sentence-to-row table:

| source sentence | rows |
|---|---|
| one ledger transaction runs a closure under the lock, hands it the read ledger, and writes only on a reported change | LS1, LS2, LS3, LS9 |
| it carries a read mode (strict or tolerant) so the purge runs inside it | LS5, LS6, LS10 |
| it carries an optional compensation step for a failed write, which the reauthorization alone uses | LS4, LS7, LS8 |
| the ledger types, the validators, and `RequestDigest` move to a leaf package under the leaf census | LS11, LS13, LS14 |
| the intent package re-exports them as aliases, so no consumer file moves | LS11, LS12 |
| a policy child under the policy census holds the admission rules | LS15, LS19 |
| the compaction takes typed liveness facts translated at the adapter | LS16, LS17 |
| admission partitions move to policy tables, and one real-filesystem journey per mutator stays | LS18 |
| a policy child under `internal/landing` with its own stage-record type | LS20, LS21, LS26, LS27 |
| the merge-tree parser stays adapter translation, and the mode-list classifier moves to the policy | LS22, LS23 |
| the settle refusal reason surfaces on `Conflict` in the named format, with a coverage row | LS28, LS29, LS30, LS31, LS32, LS33, LS34 |
| the seven classification journeys stay | LS25 |
| one settle journey per verb and one refusal journey stay | LS24 |
| the exit proof: the suite passes with test logic unchanged, except the named moved partitions and the surfaced-reason row | LS1, LS2, LS5, LS6, LS9, LS17, LS22, LS25, LS33, LS39 |
| `CONTEXT.md` gains **ledger transaction** and **settle verdict** inside the spec's fence | LS37, LS38 |
| flagged addition: census wrapper per new pure package | LS13, LS19, LS26, LS35 |

Pre-review proof checklist:

- Cited symbols: each symbol below resolves in the tree at the spec commit. The tests named `new` are new.
- Cited intent envelope symbols: `Address`, `Read`, `readPath`, `writePath`, `acquire`, `staleLock`, `Snapshot`, `filterLive`, `claudeCandidates`, `NewEntry`.
- Cited intent mutator symbols: `Upsert`, `Compact`, `PutCleanupReceipt`, `PutAssignment`, `ReauthorizeAssignment`, `PurgeAssignments`, `DeleteAssignment`, `compareAndSwapRequestDigest`.
- Cited intent schema symbols: `Ledger`, `Entry`, `Kind`, `Assignment`, `AssignmentState`, `Recovery`, `CleanupReceipt`, `RequestDigest`.
- Cited intent validator symbols: `validEntry`, `ValidateAssignment`, `validateCleanupReceipts`, `ValidIdentity`, `AssignmentBranchRef`, `RecoveryRefPrefix`.
- Cited landing symbols: `Compose`, `parseConflict`, `resolveCaptureConflict`, `unionStages`, `contentConflictKind`, `CaptureSide`, `stageRecord`, `unionStage`, `Conflict`, `ConflictError`, `CompositionResult`, `mergeTree`, `mergeTreeResult`, `editTree`, `indexRun`, `Owner`, `New`.
- Cited census symbols: `puritycensus.Scan`, `puritycensus.PolicyPackage`, `puritycensus.LeafPackage`, `puritycensus.Sources.MustHold`, `policyImport`, `leafImport`, `diagnose`.
- Cited existing tests: `TestConcurrentWritersKeepEveryEntryAndStaleLockReclaims`, `TestIdenticalStampedAssignmentWritePreservesBytes`, `TestReadEvidenceStates`, `TestSnapshotProofLifecycle`, `TestUncorrelatedEntriesUseCandidateSet`, `TestCleanupReceiptWindowKeepsExactlyLast256Completions`, `TestCompareAndSwapRequestDigestRefusesConcurrentMovementAndPreservesOtherFields`, `TestPurgeAssignmentsDropsRecordsReadRefuses`, `TestPurgeAssignmentsIsANoOpWithoutALedger`.
- Cited existing tests, continued: `TestComposeClassifiesRealGitConflictsWithoutMutation`, `TestComposeSettlesPhaseOwnedConflictsByRule`, `TestComposeRefusesConflictsTheRuleTableCannotSettle`, `TestConflictKindRejectsEmptyMergeTreeOutput`, `TestMergeRefusesAConflictOutsideTheCaptureTable`, `TestMergeResolvesAConflictedLearningsJournalAsAUnion`.
- Import edges: `internal/intent` gains `internal/intent/ledger` and `internal/intent/admissionpolicy`. `internal/intent/admissionpolicy` gains `internal/intent/ledger`. `internal/landing` gains `internal/landing/settlepolicy`. The three new packages gain `internal/puritycensus` in their test files only. No other package gains an import edge.
- Source-row clauses and occurrences: the source is the candidates 3 and 4 paragraph of `roadmap/FT302.md`. The table above lists each clause once, and each clause occurs once in that paragraph.
- Promised field labels: the settle-refusal reason field on `Conflict`. The four reason values `path outside the capture table`, `non-regular mode`, `mode disagreement`, and `union content not text`. The message `composition conflict: <kind>; settle refused: <reason> (<paths>)`.
- Changed-function callers: the seven intent mutators keep their signatures, so every caller keeps its call. `resolveCaptureConflict` has one caller, `Compose`. `CaptureSide` has two callers, both inside `resolveCaptureConflict`. `contentConflictKind` has one caller, `parseConflict`. `ConflictError` has two producers, `Merge` and `LandReviewed`, and three consumers, `land.go`, `merge.go`, and `merge_test.go`.
- Copy survival: LS12 reds a surviving leaf import outside the intent owner. LS14 reds a schema copy left in `internal/intent/assignment.go`, and LS39 reds a settle-table copy left in `internal/landing/composition_test.go`.

Reader sweep of the two decision facts. The readers are:

- the 68 Go files that import `internal/intent`, of which 30 are not test files. The enumeration command is `rg -l 'bench/internal/intent"' --type go`. No file moves, because the aliases keep every spelling, and LS11 and LS12 red a move.
- the five worktree files that match the composition-conflict text, which the closure headroom names and LS33 covers
- `internal/landing/merge.go` and `internal/landing/landing.go`, the two `ConflictError` producers, which take LS28
- `internal/worktree/land_refusal.go`, which renders `conflict.Error()` and `conflict.Paths` through one refusal face, and which takes LS33
- `internal/puritycensus/census.go`, the one source of both censuses, which takes LS36
- `CONTEXT.md`, which takes LS37 and LS38
- the seven `context-*` canary fixtures, whose pinned needles the two new terms do not touch

The shipped-surface claim words: this spec adds no claim word beside a repo-only
path. The refusal message names a reason and a path list, and it names no
package path, so `package-core-guard` reads nothing new.

The trust chain is unchanged. The transaction runs the same lock file, the same
stale-lock rule, and the same atomic replacement the seven mutators run today.

Reviewer questions:

1. The source says "none of the 41 consumer files moves". The tree holds 68 Go files that import `internal/intent`, of which 30 are not test files. Neither number is 41. The spec uses 68 and 30, and it names `rg -l 'bench/internal/intent"' --type go` as the enumeration command. My recommendation: accept 68 and 30, and let the drain correct the roadmap row.
2. `internal/puritycensus`'s policy pattern names `internal/worktree` alone as the parent adapter. Two new children sit under other parents. My recommendation: add `internal/intent` and `internal/landing` to that one pattern, which LS36 pins.
3. No standing check requires a census wrapper in a pure package, so LS35 is review-owned. My recommendation: leave it review-owned and keep the enumerating check out of scope, because the three existing census packages carry the same obligation today.
4. The source does not say which settle reason wins when two compete. My recommendation: the first refusing record in merge-tree order wins. The table check beats the mode check inside one record. Both record reasons beat the two settle reasons, which LS32 pins.
