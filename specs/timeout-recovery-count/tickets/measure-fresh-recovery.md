# Measure fresh recovery after timeout

Blocked by: none
Writes: internal/gate/run_failure_outcomes_test.go, internal/gate/timeout_recovery_count_test.go (new)
Covers: none

## What to build

Make the timeout recovery test independent of whether the timed-out child starts before its deadline.
Compare the recovery run count with the observed count after timeout.
Preserve the timeout exit, timeout evidence, non-reusability, and successful fresh recovery assertions.
Keep both early-timeout and started-child consequences observable at the existing test seam.
Move the affected test if its current file needs headroom.
Do not change production behavior or timeout policy.

## Acceptance

- [ ] A timeout before the child counter increments still requires exactly one fresh recovery execution.
- [ ] A timeout after child startup still requires exactly one fresh recovery execution.
- [ ] The timeout result remains exit 124 with timeout evidence that cannot be reused.
- [ ] Omitting recovery execution or reusing old evidence makes the test fail.
