# Carry the child environment

Blocked by: 07-move-the-working-directory-tests.md
Writes: internal/worktree/effects.go, internal/worktree/exec.go, internal/worktree/clean.go, internal/worktree/build.go, internal/worktree/build_restore.go, internal/worktree/layers.go, internal/worktree/subshell.go, cmd/bench/worktree_leaves.go, internal/worktree/exec_test.go, internal/worktree/exec_pwd_test.go, internal/worktree/reset_refusal_test.go, internal/worktree/subshell_test.go, internal/worktree/effect_census_test.go, internal/worktree/single_read_census_cases_test.go, internal/worktree/verb_runner_test.go, internal/worktree/verb_runner_check_test.go, internal/worktree/exec_call_env_test.go (new), internal/worktree/parallel_census_test.go, internal/worktree/lifecycle.go, internal/worktree/reset.go, internal/worktree/reset_apply.go, internal/worktree/reset_envelope.go, internal/worktree/reset_restore.go, internal/worktree/joins.go, internal/worktree/reset_apply_test.go, internal/worktree/reset_restore_test.go, internal/worktree/reset_restore_refusal_test.go, internal/worktree/resume_reconcile_test.go
Covers: LGP19, LGP22, LGP24, LGP34

## What to build

A verb run carries its child base environment, and `effects.go` becomes the one owner of the `os.Environ` read.

`effects.go` reads `os.Environ()` once at the verb boundary into `ambient`. Four production sites read it today, and each one takes the ambient value instead:

- `execEnv` (`internal/worktree/exec.go:201`), under `runWorktreeChild` (`exec.go:111`);
- `gitInput` (`internal/worktree/clean.go:218`), under `commitTree` and `captureLayers` (`layers.go:149`, `:153`);
- the doctor child of `repairPrimary` (`internal/worktree/build.go:331`), under `installAttributed` and `restoreAttributed` (`build_restore.go:50`);
- the subshell entry (`internal/worktree/subshell.go:40`), called from `cmd/bench/worktree_leaves.go:25`.

An exported chain that holds no `ambient` takes the base environment from one `effects.go` function, as `currentTime` serves the clock (`internal/worktree/effects.go:23`). One chain holds none: `RunTreeChild` (`exec.go:119`) through `runChild` to `execEnv` (`exec.go:130`), whose one caller is `internal/treetarget/run.go:66`. Keep its signature. The doc comment of the `effects.go` function names that chain.

The single-read census refuses a reference to an `effects.go` read from an unexported declaration (`internal/worktree/single_read_census_test.go:283-284`). So the reset recovery chain cannot use that function. `gitInput`, `commitTree`, and `captureLayers` each take an `environ []string` parameter, and the value comes from the `ambient` that each chain's verb boundary holds:

- the cleanup chain: `retireCheckout` (`lifecycle.go:492`) passes `a.environ` through `recoverAssignmentWithFault` (`clean.go:33`);
- the reset chain: `resetWith` (`reset.go:71`) passes `a.environ` as a call parameter to `applyReset` (`reset_apply.go:15`), which passes it to `resetRestoredLayersMatch` (`reset_restore.go:48`) and to the `resetEnvelope` join. The join's function type (`joins.go:62`) and `writeResetEnvelope` (`reset_envelope.go:10`) gain the parameter. `resetPlan` carries no environment, because it is a fingerprinted decision value and `applyReset` re-plans.

Test callers of the reset recovery chain pass `os.Environ()`, because the census scans production files only.

The exec entry `ExecCommandResolving` (`exec.go:52-108`), called from `cmd/bench/worktree_leaves.go:20`, is a verb boundary, so it builds its own `ambient` value. `installAttributed` (`build.go:158`, called at `:71`) has no `ambient` parameter today; pass the value down from its verb boundary. `subshellAt` (`subshell.go:43`) keeps its `environ` parameter for its test callers (`subshell_test.go:27`, `:46`, `:82`); only its entry stops its own read.

The exec tests call `runWorktreeChild` directly (`exec_test.go:33`), not `runVerb`. So the base environment enters through `runWorktreeChild` and `execEnv`. `exec_pwd_test.go:84` also calls `execEnv`.

A signature change of `runWorktreeChild` reaches its calls at `exec_test.go:33`, `:47`, `:272`, and `:343`. `exec_test.go` has 421 lines and no grant, so change only those call lines, with no net line gain.

`TestEffectBoundaryCensus` gains `os.Environ` in `ambientReadPattern` (`internal/worktree/effect_census_test.go:14`). No other `os.Environ` read stays in the package's production files. Update the synthetic copy of `effects.go` in `single_read_census_cases_test.go:15-29`.

`verbCall` gains a child base environment field. Extend `checkVerbCall` and its refusal tests in `verb_runner_check_test.go`, which ticket 06 also extended. Keep that file at or under 400 lines. The field applies to the verbs with a joins form. Exec has none (`verb_runner_check_test.go:168-178`), so the LGP19 test passes the base environment to `runWorktreeChild` directly.

Put the LGP19 test in the new `internal/worktree/exec_call_env_test.go`. It calls `t.Parallel`, and so does the `verbCall` environment refusal test, so the count rises by two. Record that arithmetic.

Record `bench test --full --package ./internal/worktree` before and after the change, for LGP34.

## Acceptance

- [ ] An exec child sees a value from the call's base environment and not a process-only value (LGP19).
- [ ] An `os.Environ` read outside `effects.go` reds the effect census (LGP24).
- [ ] The test count equals the recorded sum (LGP22).
- [ ] The probe "restore the os.Environ read in execEnv" bites.
- [ ] The verification record holds the package time before and after the change (LGP34).
