# Worktree orphan reconcile — release and resume clear tree-gone records safely

Status: implemented

## Problem

FT93(a) made `bench worktree release` surface the retained verdict, but two orphan
paths remain and were the source of the 10 open assignment records at session start.

- **(b) Release dead-ends on an out-of-band removal.** When a request-less
  `bench worktree clean --discard-ignored --apply` removes an owned worktree, the
  assignment record survives. A later `bench worktree release --request <id> <path>`
  finds the tree gone, finds the explicit-clean cleanup receipt (which does not match
  the automatic-registration reconcile shape), and returns
  `cleanup receipt does not authorize release reconciliation` — the record is
  stranded with no next step.
- **(c) Resume never visits the orphans.** `ConservativeCleanup` iterates only
  *registered* git worktrees, so an assignment record whose tree is gone and is no
  longer a registered worktree is never reconciled; it accumulates on the ledger and
  inflates the open count forever.

Neither path distinguishes **residue** (a record whose tree is gone and which
preserves no work — safe to compact) from **preserved work** (a record holding a
recovery pointer that must never be silently discarded). The fix must honor that
distinction: this session found all 8 tree-gone records held live recovery refs, and
auto-deleting them would have severed the only ledger pointer to preserved work.

## Solution

One safety predicate, single-sourced: **a tree-gone assignment record is residue iff
it carries no recovery pointer (`len(Recovery) == 0`); otherwise it holds preserved
work.** The record's own `Recovery` field is the source — no on-disk ref probing.
Both paths consult that one predicate.

- **(b) Release** at the out-of-band dead-end reconciles residue — writes the terminal
  release receipt so `finishReleaseReceipt` compacts the record — and, for a record
  that preserves work, returns an actionable verdict naming
  `bench worktree recovery <ref>` (recover or retire), leaving the record and ref
  intact.
- **(c) Resume**, after its registered-worktree pass, sweeps assignment records whose
  tree is gone and which are not registered worktrees: it compacts residue (deletes the
  record, counted in the summary) and reports preserved-work orphans with their exact
  recovery command — never deleting them.

The recovery-retire guard is unchanged: un-landed preserved work still cannot be
auto-discarded; the reviewer retires it deliberately (as done this session). This
mechanism only clears records that hold *no* work, and *names* the deliberate path for
records that do.

## User stories

1. As a session, I want `release --request` on a tree removed out of band with no
   preserved work to reconcile and compact the record (exit 0, the record gone),
   instead of dead-ending on "cleanup receipt does not authorize release
   reconciliation", so an out-of-band clean does not strand the ledger. Line:
   in-session / high. This is the (b) reconcile seam and the load-bearing safety split.
2. As a session, I want `release --request` on a tree removed out of band that still
   holds preserved work to return a verdict naming `bench worktree recovery <ref>` and
   to leave the record and recovery ref intact, so release never silently discards
   preserved work. Line: in-session / high. The preserved branch must never compact.
3. As a session, I want `bench resume` to sweep tree-gone, unregistered residue records
   — compacting them and counting them — so orphaned records stop accumulating across
   sessions. Line: in-session / high. This is the (c) accumulation fix.
4. As a session, I want `bench resume` to report, never delete, tree-gone records that
   hold preserved work, naming each `bench worktree recovery <ref>` command, so
   preserved work is surfaced for a deliberate recover-or-retire. Line: in-session /
   high. The report-not-delete rule is the safety property.
5. As a reviewer, I want active-state records and any record whose tree still exists to
   be left untouched by the sweep, so a live session's registration is never reconciled
   out from under it. Line: in-session / medium. The conservative exclusion.

## Implementation decisions

- **One residue predicate.** `func residualAssignment(a intent.Assignment) bool { return len(a.Recovery) == 0 }`
  in `internal/worktree`, consumed by both the release reconcile and the resume sweep.
  The reconcile-vs-preserve decision lives in exactly one place (code standard: one
  source per fact).
- **(b) lives in `releaseAssignment` (`internal/worktree/ownership.go`).** At the
  missing-tree branch, the current `else if found { return error }` for a found-but-
  non-matching cleanup receipt is replaced: load the orphaned assignment (by the
  receipt's `Assignment`/`Request`); if residue, set it `StateComplete`, write
  `receiptFromRelease(...ActionRemoved)` and return it (the command layer's
  `finishReleaseReceipt` deletes the record); if it preserves work, return
  `errRecoveryPending` naming the first recovery ref and the `bench worktree recovery`
  command. Idempotent: a second call short-circuits on the completed release receipt
  already written at the top of `releaseAssignment`.
- **(c) lives in `ConservativeCleanup` (`internal/worktree/resume.go`).** After the
  registered-worktree loop, a second loop over `intent.Assignments(root)` visits each
  record whose `Worktree` is neither an existing path nor a registered worktree path
  and whose state is not `StateActive`. Residue records are deleted via
  `intent.DeleteAssignment` and counted in a new `ResumeResult.Reconciled`. Preserved-
  work records are collected into `ResumeResult.Preserved` (one entry per record: id +
  first recovery ref) and left intact.
- **Summary rendering (`renderResumeSummary`).** The `bench resume` line gains
  `reconciled <n>`; when preserved orphans exist, a follow-up block names each with its
  `bench worktree recovery <ref>` command, so a session at start sees exactly what work
  is parked and how to act on it.
- **Registered-path set.** Built once from `ClassifyRegisteredWorktrees(root)` (the same
  source the first loop uses) and compared with the existing `samePath` helper, so a
  registered-but-dir-missing worktree (prunable) is left to the git-worktree path, not
  swept as an orphan.

## Testing decisions

- **Assert observable outcome, not internals:** the release exit code and returned
  receipt, the presence/absence of the assignment record and recovery ref on the
  ledger, and the `ResumeResult` counts — never the private control flow. This keeps
  each expectation independent of the implementation so a named omission turns the gate
  red.
- **One in-process seam, `internal/worktree`.** All rows attach here, driving
  `ReleaseCommand`/`releaseAssignment` and `ConservativeCleanup` against fixture repos.
  Prior art to compose: `newOwnedAssignment`/`newPendingAssignment`/`markPending`,
  `commitInWorktree`, `requireTest` (`resume_test.go`), and the explicit-clean apply
  path (`CleanCommand` / `ApplyExplicitWithOptions`) to construct the out-of-band
  removal faithfully as a fixture. Recovery fixtures set `assignment.Recovery` directly.
- **Gate command:** the project gate, `.bench/gate.sh`; the worktree unit tests run
  under the compiled-core checks.

## Acceptance coverage map

| # | Story | Seam | Red signal (fails without the fix) | Asserts |
|---|---|---|---|---|
| 1 | 1 | worktree | Out-of-band-cleaned residue release returns "cleanup receipt does not authorize release reconciliation" | `ReleaseCommand` exits 0, returned receipt action `removed`, assignment record gone |
| 2 | 1 | worktree | — (idempotency) | A second `ReleaseCommand` with the same request exits 0 and returns the same receipt |
| 3 | 2 | worktree | Preserved-work out-of-band release dead-ends with the cryptic message | `ReleaseCommand` exits non-zero, stderr names `bench worktree recovery <ref>`, record and recovery ref still present |
| 4 | 3 | worktree | `ConservativeCleanup` leaves a tree-gone unregistered residue record on the ledger | After `ConservativeCleanup`, the residue record is deleted and `Reconciled == 1` |
| 5 | 4 | worktree | A tree-gone record holding a recovery pointer is deleted (or ignored) by the sweep | The preserved record survives, appears in `Preserved`, and its recovery ref is untouched |
| 6 | 5 | worktree | An over-broad sweep deletes an active or tree-present record | An `active` tree-gone record and a tree-present record both survive `ConservativeCleanup` |

## Out of scope

- **No new CLI subcommand** — reconcile rides `release` and `resume`.
- **The recovery-retire landedness guard is unchanged.** Un-landed preserved work is
  never auto-discarded; the reviewer retires it deliberately via
  `bench worktree recovery <ref> --apply <fingerprint>` (or an explicit override).
- **A one-shot sweep of pre-existing orphans is not part of this** — those were cleared
  by hand this session; this mechanism prevents recurrence.
