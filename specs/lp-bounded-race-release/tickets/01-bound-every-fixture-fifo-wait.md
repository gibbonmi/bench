# Bound every fixture FIFO wait

Blocked by: none
Writes: internal/systemtest/owner_landing_fixture_test.go, internal/systemtest/owner_landing_barrier_test.go (new)
Covers: none

## What to build

The landing race fixture gate signals the ready FIFO and then reads the release FIFO. Each FIFO open blocks until a partner opens the other end. Before this ticket, the gate made both opens with no bound. The gate runs in its own process group. When the test kills the landing, or the test binary stops, the gate stays blocked with no end.

Put the barrier in one helper, `landingRaceBarrier`. The helper prefixes the barrier with `boundedWait`, so the ambient timeout ends a barrier that loses its partner. The landing race gate and the artifact phase both call the helper with one window, `landingBarrierWindow`.

Add a system test that runs the barrier with a short window. One case has no ready reader. One case has no release writer. Each case must end with the timeout exit inside the deadline that `bounds.TestDeadline` derives from the window.

## Acceptance

- [ ] The landing race gate and the artifact phase get their barrier from `landingRaceBarrier`, and no other fixture script in `internal/systemtest` makes a FIFO open with no bound.
- [ ] `TestLandingRaceBarrierEndsWithoutItsPartner` passes. A probe that removes the `boundedWait` prefix turns it red with a timeout verdict.
- [ ] `bench test --check system` passes.
