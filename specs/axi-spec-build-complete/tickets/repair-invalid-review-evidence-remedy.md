# Repair the invalid review-evidence remedy

Blocked by: none
Ownership fence: `internal/specbuild/lifecycle.go`, `internal/specbuild/lifecycle_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: RV1/invalid-review-evidence-advertises-review

## What to build

Close the accepted Spec and Coverage findings (P1, C1) from the Terra/xhigh
review of candidate `72cab0bd3e1947b724d6c417cb412fa3220d22f8`:
`Service.Review` (`internal/specbuild/lifecycle.go` around line 334) refuses
an invalid or mismatched review receipt with
`statusRefusal(RefusalInvalidEvidence, slug, errInvalidReviewReceipt)`.
`statusRefusal` (`internal/specbuild/disclosure.go`) hardcodes its operation
as `"status"`, so `RefusalForClass` routes `RefusalInvalidEvidence` through
`operationAction("status", slug)` — the catalog's `default` case, which
returns `bench spec build status <slug>`, not a command that lets the caller
retry with corrected review evidence. A caller whose review receipt was
rejected is told to check status, not to resubmit
`bench spec build review <slug> --evidence <receipt>` — the FT164 wrong-remedy
class story 2 and SB3 (`specs/axi-spec-build-complete/spec.md` around line 58)
exist to close. No existing test drives this exact path: the checkpoint-side
hostile-evidence tests (`internal/specbuild/checkpoint_hostile_test.go`) cover
`Checkpoint`'s own invalid-evidence remedy, not `Review`'s, and the disclosure
matrix's `review/state` fault point (`internal/specbuild/disclosure_observation.go`
around line 158) injects its fault after receipt validation, so it never
reaches this exact refusal either.

Fix: in `Service.Review`, replace
`statusRefusal(RefusalInvalidEvidence, slug, errInvalidReviewReceipt)` with a
call that threads `"review"` as the operation instead of `"status"` — for
example `operationRefusal(RefusalInvalidEvidence, "review", slug,
errInvalidReviewReceipt)` (see `operationRefusal` in
`internal/specbuild/disclosure.go`, which is exactly `RefusalForClass(class,
operation, slug, cause)`). With `operation == "review"`,
`operationAction("review", slug)` already returns the correct
`bench spec build review <slug> --evidence <receipt>` template — no change is
needed to `RefusalForClass`, `operationAction`, or any other refusal class's
remedy. Do not touch any other refusal construction in this file or in
`disclosure.go`.

New coverage in `internal/specbuild/lifecycle_test.go`: drive the real
`Service.Review` with an evidence file that fails `readReviewReceipt`'s
validation (or whose `Run`/`Candidate` fields mismatch the retained run —
whichever is simplest against the file's existing fixtures), read the
resulting typed refusal through `RefusalFacts`, and require its one action's
`Command()` is `bench spec build review '<slug>' --evidence <receipt>` — not
a `status` command.

## Acceptance

- [ ] [RV1] (covers local) (P1, C1) `Service.Review` given invalid or
  mismatched review evidence returns a `RefusalInvalidEvidence` refusal whose
  remedy is `bench spec build review '<slug>' --evidence <receipt>`, driven
  through the real public `Review` call (not a hand-constructed refusal);
  every other refusal class's remedy is unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RV1/invalid-review-evidence-advertises-review | restore `statusRefusal(RefusalInvalidEvidence, slug, ...)` in `Service.Review` | focused real-`Review` refusal test | drive `Service.Review` with invalid or mismatched evidence, read the remedy via `RefusalFacts`, and require the `review --evidence <receipt>` command rather than `status` |
