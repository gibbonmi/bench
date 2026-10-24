# Keep the shift response on a checkpoint exit

Blocked by: none

Writes: internal/shift/session.go, internal/shift/loop.go, internal/shift/fault_test.go, internal/shift/lifetime_test.go, internal/shift/shift_test.go, internal/systemtest/process_lifetime_callers_test.go, internal/systemtest/process_lifetime_shift_test.go, tests/canary/process-lifetime-closure/private-shift-loop/EXPECT, CHANGELOG.md

Covers: none

## What to build

The `bench shift` verb runs under the bounded response. The bounded owner keeps stdout and stderr, and it prints them after the verb returns. A shift checkpoint stops the loop when the wall deadline or a signal has fired. The checkpoint must not exit the process. It must preserve the work, resolve the Result, and return that Result to the loop. The loop then returns through `finish`, so the bounded owner prints the response.

The refactor phase has checkpoints too. It returns its stop Result to the loop in the same way. A stage, commit, or rollback failure in the refactor phase also returns its evidence Result.

## Acceptance

- [ ] A shift that passes its wall deadline exits 3, and its stdout holds the `shift_result` block with the `incomplete` outcome. The system test drives the sealed `bench shift` executable.
- [ ] A shift that receives SIGINT exits 130, and its stdout holds the `shift_result` block with the `interrupted` outcome. The system test drives the sealed `bench shift` executable.
- [ ] A checkpoint stop whose teardown fails makes Loop return 1, and its stdout holds the `shift_result` block with the `failed` outcome and the teardown detail.
- [ ] The adapter group is absent after the deadline exit and after the interrupt exit.
- [ ] No checkpoint path calls `os.Exit`.
