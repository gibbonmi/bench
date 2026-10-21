# Pass the child environment through the call

Blocked by: 05-resolve-operands-against-the-call-directory.md
Writes: internal/worktree/effects.go, internal/worktree/exec.go, internal/worktree/clean.go, internal/worktree/build.go, internal/worktree/build_restore.go, internal/worktree/layers.go, internal/worktree/subshell.go, cmd/bench/worktree_leaves.go, internal/worktree/verb_runner_test.go, internal/worktree/verb_runner_check_test.go, internal/worktree/exec_test.go, internal/worktree/exec_call_env_test.go (new), internal/worktree/exec_pwd_test.go, internal/worktree/reset_refusal_test.go, internal/worktree/subshell_test.go, internal/worktree/effect_census_test.go, internal/worktree/single_read_census_cases_test.go, internal/worktree/parallel_census_test.go
Covers: LGP19, LGP20, LGP21, LGP22, LGP24, LGP34

## What to build

A verb call carries its child base environment, so the exec tests no longer bind the process environment.

`effects.go` reads `os.Environ()` once at the verb boundary into `ambient`. Four production sites read it today, and each one takes the ambient value instead:

- `execEnv` (`internal/worktree/exec.go:201`);
- `gitInput` (`internal/worktree/clean.go:218`);
- the doctor child in `repairPrimary` (`internal/worktree/build.go:331`);
- the subshell entry (`internal/worktree/subshell.go:40`).

`TestEffectBoundaryCensus` gains `os.Environ` in `ambientReadPattern` (`internal/worktree/effect_census_test.go:14`). No other `os.Environ` read stays in the package's production files. Update the synthetic copy of `effects.go` in `single_read_census_cases_test.go:19-27`.

`verbCall` gains a child base environment field. The exec verb has no joins form today (`verb_runner_test.go:111-113`), so give the exec path an entry that takes the field.

Move the 10 exec child-environment tests in `exec_test.go` to `t.Parallel`. Nine reach `bindEnv` through `childOutputAtHome` (`exec_test.go:26-31`), and `TestExecChildTakesTheExplicitHome` calls it at `:106`. Update the stale comments at `exec_test.go:24-26` and `:99-103`.

`exec_test.go` has 421 lines and no grant, so it may not gain net lines. Put the LGP19 test in the new `internal/worktree/exec_call_env_test.go`.

Set both pins by the recorded arithmetic. The Home-default graders and the PATH-stub tests stay serial.

## Acceptance

- [ ] An exec child sees a value from the call's base environment and not a process-only value (LGP19).
- [ ] `TestParallelCensusOnTheLiveTree` passes with every moved exec test parallel (LGP20).
- [ ] The serial set equals the recorded ceiling, and the test count equals the recorded sum (LGP21, LGP22).
- [ ] An `os.Environ` read outside `effects.go` reds the effect census (LGP24).
- [ ] The probe "restore the os.Environ read in execEnv" bites.
- [ ] The verification record holds the package time before and after the change (LGP34).
