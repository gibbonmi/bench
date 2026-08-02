# Exempt abandon from the recomposition refusal

Blocked by: Name the refused operation in the working-subject refusals
Ownership fence: `internal/specbuild/precondition.go` (the recomposition return
site only), `internal/specbuild/abandon_test.go`
Assumptions: `preconditions` computes `recompose` early, refuses an
unrecognized advance immediately with the subject-mismatch message, and returns
the `errRecompose` wrapper as its *last* act, after branch identity, spec
identity, candidate currency, assignment ownership, and operation evidence have
all passed. Only that last return is this ticket's. Re-derive the site from the
tree at pickup: the preceding ticket has already edited this file.

## What to build

An escape hatch that requires the state it exists to escape is not an escape
hatch. `bench spec build abandon --apply` sits behind the shared recomposition
refusal, so a run whose branch tip moved — which is every mid-repair run, since
committing repair tickets is what moves it — cannot be retired at all. The run
that motivated this spec is still `active` on `main` for exactly this reason.

Exempt the abandon mutation at the single site that returns `errRecompose`.
Nothing else moves.

**The exemption is from recomposition, not from identity.** A branch that was
rebased or amended fails the recognized-advance check *earlier* in the same
function, with the subject-mismatch message, and must keep failing it. That is
identity drift, not recomposition: cleaning up against a history that no longer
contains the run's own base acts on the wrong tree. Exempting the whole
precondition call rather than this one branch is the cheap wrong fix, and the
criteria below are sized to catch it.

The fingerprint drift check that makes `--apply` safe lives in `ApplyAbandon`,
not here, and is untouched.

## Acceptance

- [ ] AB1 — `ApplyAbandon` with a valid fingerprint succeeds when the branch tip has advanced by a recognized ancestor commit.
- [ ] AB2 — `ApplyAbandon` on a moved tip still refuses a drifted fingerprint.
- [ ] AB3 — `ApplyAbandon` on a moved tip still refuses on branch identity drift, on spec identity drift, and on candidate-ref drift.
- [ ] AB4 — `ApplyAbandon` on a rebased or amended branch still refuses with the subject-mismatch message, not with recomposition.
- [ ] AB5 — `start`, `assign`, `checkpoint`, `integrate`, and `review` each still receive the recomposition refusal on the same moved tip; `promote` is the one exclusion, because on a moved tip it recomposes instead of refusing.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AB1 | leave the refusal in place and exempt abandon in `ApplyAbandon`'s caller instead | `TestApplyAbandonSucceedsOnMovedTip` | compose a run with assignments, commit a new commit on the working branch, plan, apply with the plan's fingerprint, assert no error and a terminal run |
| AB2 | exempt abandon from the fingerprint comparison alongside the recomposition refusal | `TestApplyAbandonRefusesDriftedFingerprintOnMovedTip` | as AB1, but apply with a fingerprint from a stale plan, expect the drift refusal |
| AB3 | return early for abandon before the identity checks rather than at the recomposition return | `TestApplyAbandonStillRefusesIdentityDriftOnMovedTip` | three subtests, each mutating one of branch, spec tip, and candidate ref, then applying on a moved tip, expecting the matching refusal |
| AB4 | treat any tip inequality as recomposition for abandon | `TestApplyAbandonRefusesUnrecognizedHeadMove` | compose a run, amend the base commit so the recorded base is no longer an ancestor, apply, expect the subject-mismatch refusal |
| AB5 | key the exemption off the presence of a moved tip rather than off the mutation | `TestNonAbandonMutationsStillRecomposeOnMovedTip` | enumerate the `mutation` constants at test time, exclude `abandon` and `promote` by name with the reason recorded, invoke each remaining one on a moved tip, and assert each returns the recomposition refusal naming `promote` |
