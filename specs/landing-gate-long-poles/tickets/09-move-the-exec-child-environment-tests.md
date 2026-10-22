# Move the exec child-environment tests

Blocked by: 08-carry-the-child-environment.md
Writes: internal/worktree/exec_test.go, internal/worktree/exec_call_env_test.go (new), internal/worktree/parallel_census_test.go
Covers: LGP20, LGP21, LGP22, LGP34

## What to build

The ten exec child-environment tests pass the call's base environment instead of binding the process environment, and they call `t.Parallel`.

Nine tests reach `bindEnv` through `childOutputAtHome` (`internal/worktree/exec_test.go:26-31`), by way of `childOutput`, `childEnvironment`, or `childWrapper`. `TestExecChildTakesTheExplicitHome` also calls `bindEnv` at `:106`. The ten tests sit at `exec_test.go:103-240`.

`childOutputAtHome` passes the base environment through `runWorktreeChild` from ticket 08 instead of `bindEnv`. Update the stale comments at `exec_test.go:17-18`, `:24-25`, and `:98-102`. A test that still reaches the `effects.go` environment function of ticket 08 reads the process environment, so it stays serial; record it.

The exec class comment at `parallel_census_test.go:1033-1040` says that the exec child tests bind the caller environment. Update it with no net line gain.

`exec_test.go` has 421 lines and no grant, so it may not gain net lines. Move the helpers to `exec_call_env_test.go` if the ten `t.Parallel` lines need the room.

Set both pins by the recorded arithmetic. The Home-default graders and the PATH-stub tests stay serial.

Record `bench test --full --package ./internal/worktree` before and after the change, for LGP34.

## Acceptance

- [ ] `TestParallelCensusOnTheLiveTree` passes with the ten exec tests parallel (LGP20).
- [ ] The serial set equals the recorded ceiling, and the test count equals the recorded sum (LGP21, LGP22).
- [ ] The probe "bind the caller environment again in childOutputAtHome" bites.
- [ ] The verification record holds the package time before and after the change (LGP34).
