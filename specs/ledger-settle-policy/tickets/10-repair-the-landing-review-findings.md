# Repair the landing review findings

Blocked by: 08-surface-the-settle-refusal-reason.md
Writes: internal/landing/settlepolicy/settlepolicy.go, internal/landing/settlepolicy/settlepolicy_test.go, internal/landing/composition.go, internal/landing/composition_test.go, internal/landing/landing_helpers_test.go, internal/landing/merge_test.go
Covers: LS24, LS41

## What to build

Verify the premise first. Read `Settle`, `CaptureSide`, and the mode
constants in internal/landing/settlepolicy/settlepolicy.go. Read
`resolveCaptureConflict` and `unionStages` in
internal/landing/composition.go. Read `TestComposeSettlesPhaseOwnedConflictsByRule`
in internal/landing/composition_test.go and the conflict fixtures in
internal/landing/landing_helpers_test.go and internal/landing/merge_test.go.
Read the two deleted rows `union-deleted-on-one-side` and
`union-added-on-both-sides` at the base commit 897f52be with
`git show 897f52be:internal/landing/composition_test.go`.

In `Settle`, answer a removal verdict for a union path whose stage 2 and
stage 3 are both absent. The adapter then never hands a side-less union to its
text merge, so it never names `union content not text` for it. Add the table
case to `TestSettlePolicyAnswersTheCaptureRule`.

Restore the two union journeys to the composition suite as rows of
`TestComposeSettlesPhaseOwnedConflictsByRule`, with their test logic as it was
at the base. They exercise the adapter's one-sided union arms, which no policy
table reaches. Place shared fixtures in landing_helpers_test.go or
merge_test.go so that every file stays under 400 lines, and write no new file
inside internal/landing/.

Correct the comment above the mode constants in settlepolicy.go. It says "The
two regular file modes" above four constants; name what each group is.

## Acceptance

- [ ] The settle policy answers a removal for a union path whose two sides are both absent.
- [ ] `Compose` settles a union path deleted on one side by publishing the present side.
- [ ] `Compose` settles a union path added on both sides by their union.
- [ ] The mode-constant comment describes the four constants it heads.
- [ ] Every `internal/landing` file stays under 400 lines, and the directory holds 16 source files.
- [ ] The pre-existing `internal/landing` suite passes with its test logic unchanged.
- [ ] Self-probe: hand the side-less union to the text merge again, and report the new policy case red.
