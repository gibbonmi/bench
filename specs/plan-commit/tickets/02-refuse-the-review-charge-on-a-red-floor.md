# Show that a red check floor refuses the review charge

Blocked by: 01-grade-the-chunk-table-against-the-map.md
Writes: internal/preflight/check_floor_command_test.go
Covers: PC19, PC20, PC21

## What to build

The review charge already refuses on any red verdict row (`internal/preflight/review.go:207-212`), but no test shows that a red `check-floor` row refuses it. Add `TestCheckFloorRefusesTheReviewCharge` beside `TestCheckFloorRefusesTheBuildCharge` in `internal/preflight/check_floor_command_test.go`. Seed the repository with `seedFloorRepo`, activate the assignment with `preflighttest.ActiveAssignment`, and run the review charge with `preflighttest.ReviewArgs`. `ReviewArgs` activates nothing, so the activation is part of the setup.

The review refusal order is the checkout, then the red verdict, then the sources (`internal/preflight/review.go:208-225`), so the seeded repository needs no review sources. The test asserts exit 1 and the refusal text `preflight required: check-floor`, the same text that the build-charge test asserts.

This ticket changes no production code. Before the first edit, prove that the new test is red-capable: run it with `bench probe` against a mutation that removes the `verdict.Red` refusal from the review charge, and record the verdict.

PC20 and PC21 cite existing tests. Read `TestCheckFloorRefusesTheBuildCharge` and `TestCommitChainReportsRedPreflight` and confirm that each still asserts its row at the ticket tip.

## Acceptance

- [ ] `TestCheckFloorRefusesTheReviewCharge` exits 1 with `preflight required: check-floor` for a plan that omits a floor check.
- [ ] A mutation that removes the review charge's `verdict.Red` refusal turns the new test red.
- [ ] `TestCheckFloorRefusesTheBuildCharge` and `TestCommitChainReportsRedPreflight` stay green and unchanged.
