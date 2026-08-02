# Reach recomposition before the gates it discards

Blocked by: Name the refused operation in the working-subject refusals
Ownership fence: `internal/specbuild/assign.go` (`Promote` only),
`internal/specbuild/promotion_recompose_test.go`
Assumptions: `Promote` today runs, in order — the terminal short-circuit, the
promotion recovery fast path, the clean-review check, the released-and-
integrated assignment loop, `validatePromotionEvidence`, and only then
`preconditions` with its recomposition branch. `recomposePromotion` already
replays the candidate onto the new tip, re-bootstraps green evidence there,
clears the review, resets the promotion fields, and drops the stale promote
operation. Nothing inside recomposition changes here. Re-derive the ordering
from the tree at pickup — an earlier ticket has edited this file.

## What to build

This is the deadlock. A run whose review returned accepted findings must commit
repair tickets to assign them; committing moves the branch tip; a moved tip is
exactly what makes `assign`, `checkpoint`, `integrate`, and `review` refuse with
a message pointing at `promote`. `promote` then refuses too, because it grades
the review and the assignment releases *before* it reaches recomposition — and a
mid-repair run can satisfy neither by construction. An operation that exists to
resolve a state must never refuse on that state.

Move the precondition call and its recomposition branch ahead of the three gates
that recomposition is about to discard, so a moved tip recomposes and reports
that a fresh review is next.

**Three gates precede recomposition, not two.** The clean-review check and the
assignment loop are the visible pair. `validatePromotionEvidence` is the third,
and it independently refuses the mid-repair state on both an absent review and
an unreleased assignment — so a reorder that moves only the first two leaves the
deadlock exactly where it was. That third call is *also* made by the promotion
recovery fast path, which must keep it. The edit therefore duplicates it into
the two paths that need it rather than relocating one call.

**The recovery fast path stays on top.** It resolves an already-published
promotion commit and never consults the precondition layer. A reorder placed
above it would re-gate a terminal run.

**The checks are not weakened, only repositioned.** They still guard the
promotion path; they stop guarding the path to recomposition. Deleting them is
the degenerate fix, and half of that degenerate fix — dropping the evidence
validation from the promotion path while duplicating it into recovery — is
invisible unless each evidence fault is asserted on its own.

**One existing expectation inverts.** The promotion test asserting that an
unreleased run is refused before recomposition encodes today's contract and is
authored behavior change under this spec. Invert it deliberately and say so in
your checkpoint; do not delete it. The test pinning the recomposition refusal
wording keeps its expectation, because that message still points at `promote`.

## Acceptance

- [ ] PR1 — a reviewed run whose branch tip advanced returns no error from `Promote` and reports review as its next action.
- [ ] PR2 — a run with accepted review findings *and* unreleased assignments recomposes on `Promote`, with no refusal.
- [ ] PR3 — after recomposition the review is cleared, the recorded base equals the new tip, and the candidate ref names the replayed commit.
- [ ] PR4 — with the tip unmoved, `Promote` still refuses an absent review, a stale-candidate review, a review with accepted findings, and any unreleased or unintegrated assignment.
- [ ] PR5 — with the tip unmoved, `Promote` still refuses each retained-evidence fault: drifted candidate ref, drifted review binding, incomplete checkpoint fields, drifted checkpoint ref, and an integration outside candidate ancestry.
- [ ] PR6 — the promotion recovery fast path is unchanged for an already-published promotion commit, including its own evidence validation.
- [ ] PR7 — `Promote` called twice on a moved tip recomposes once and then reports review as next, without a second replay.
- [ ] PR8 — a tip advance that is not a recognized ancestor still refuses with the subject-mismatch message rather than recomposing.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PR1 | move the precondition call above the clean-review check only | `TestPromoteRecomposesOnMovedTip` | compose a reviewed run, commit on the working branch, call `Promote`, assert no error and next equals review |
| PR2 | move the clean-review check and the assignment loop but leave `validatePromotionEvidence` in place | `TestPromoteRecomposesMidRepairRun` | compose a run whose review carries an accepted finding and whose assignment is unreleased, move the tip, call `Promote`, assert no error |
| PR3 | return early on the recomposition branch without calling `recomposePromotion` | `TestRecompositionClearsReviewAndAdvancesCandidate` | as PR1, then reload the record and assert review nil, base equals the new tip, and the candidate ref equals the replayed commit |
| PR4 | delete the clean-review check and the assignment loop instead of moving them | `TestPromoteStillRefusesUnreadyRunOnUnmovedTip` | four subtests on an unmoved tip, each composing one unready state, expecting the matching refusal |
| PR5 | duplicate the evidence validation into the recovery path only, dropping it from the promotion path | `TestPromoteStillRefusesEachEvidenceFault` | five subtests on an unmoved tip, each corrupting one retained-evidence fact, expecting a refusal before any gate execution |
| PR6 | place the precondition call above the recovery fast path | `TestPublishedPromotionRecoveryIsUnchanged` | drive a promotion to a published commit, call `Promote` again, assert the recovery result and that no recomposition occurred |
| PR7 | recompose without clearing the stale promote operation | `TestPromoteTwiceOnMovedTipRecomposesOnce` | as PR1, call `Promote` a second time, assert no error, next equals review, and the candidate ref is unchanged from the first call |
| PR8 | treat any tip inequality as recomposition | `TestUnrecognizedHeadMoveStillRefusesPromote` | compose a run, amend the base so it is no longer an ancestor of the tip, call `Promote`, expect the subject-mismatch refusal |
