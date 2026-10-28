# Refuse a plain chunk requirement

Blocked by: 01-accept-probe-only-plans-and-retire-check-floor.md
Writes: internal/reviewrecord/plan.go, internal/reviewrecord/plan_probe_test.go (new), internal/reviewrecord/recordtest/fixture.go, internal/reviewrecord/source_test.go, internal/reviewrecord/recordcmd/verification_test.go, internal/reviewrecord/recordcmd/ticket_test.go, internal/reviewrecord/recordcmd/ticket_carry_test.go, internal/reviewrecord/recordcmd/ordering_test.go, internal/reviewrecord/recordcmd/completion_test.go, internal/reviewrecord/recordcmd/confirmation_test.go, internal/reviewrecord/recordcmd/assignment_test.go, internal/reviewrecord/recordcmd/chunk_test.go, internal/reviewrecord/recordcmd/amendment_test.go, internal/gate/review_checkpoint_test.go, internal/preflight/preflighttest/fixture.go, internal/preflight/preflighttest/testdata/review-documents.toon, internal/preflight/completion_plan_test.go, internal/preflight/charge_test.go, internal/preflight/staleness/completeness_test.go, internal/preflight/ownership/source_test.go, internal/preflight/ownership/effects_test.go, internal/preflight/ownership/premises_test.go, internal/preflight/ownership/charge_fixture_test.go, internal/preflight/ownership/premises_fixture_test.go, internal/assessment/comparison_test.go, internal/assessment/comparison_report_test.go, internal/worktree/delegated_integration_test.go, cmd/bench/record_operations_test.go
Covers: BC23, BC33

## What to build

The plan grade refuses a chunk requirement without a `probe` key. The refusal names the chunk ID and the requirement ID. In `requirementsValid`, the refusal starts with the line `if owned != nil && item.Probe == "" {`, because the probe-02 command of the plan swaps that exact line.
The final list keeps plain requirements.

The shared plan fixtures hold plain chunk requirements now: `additional` in `internal/reviewrecord/recordtest/fixture.go` and `tests` in `internal/preflight/preflighttest/fixture.go`. Move each fixture and each test that reads it to probe requirements, or drop the plain requirement where a test does not need it.
Each moved test keeps the behavior that it asserts. A test whose subject is a plain chunk requirement now asserts the refusal.

`bench preflight build` reds the `completion-plan` row for a plan with a plain chunk requirement, and the row names the requirement. This red is the route by which each staged spec changes at its build entry.

## Acceptance

- [ ] The plan grade refuses a chunk requirement without `probe`, and it names the chunk ID and the requirement ID (BC23).
- [ ] A staged-spec fixture with one plain chunk requirement reds the `completion-plan` row of `bench preflight build` and names the requirement (BC33).
- [ ] No shared fixture builds a plain chunk requirement.
