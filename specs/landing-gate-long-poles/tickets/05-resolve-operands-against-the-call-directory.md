# Resolve operands against the call directory

Blocked by: 04-read-worktree-locations-through-the-run-value.md
Writes: internal/canonicalpath/canonicalpath.go, internal/canonicalpath/canonicalpath_test.go, internal/worktree/effects.go, internal/worktree/subshell.go, internal/worktree/path.go, internal/worktree/land.go, internal/worktree/land_resume.go, internal/worktree/reauthorize.go, internal/worktree/worktree.go, internal/worktree/verb_runner_test.go, internal/worktree/verb_runner_check_test.go, internal/worktree/single_read_census_cases_test.go, internal/worktree/serial_ceiling_test.go, internal/worktree/parallel_census_test.go, internal/worktree/identifier_operand_test.go, internal/worktree/resume_reconcile_test.go, internal/worktree/request_token_test.go, internal/worktree/identity_component_test.go, internal/worktree/land_resume_test.go, internal/worktree/land_reauthorize_operand_test.go, internal/worktree/path_identifier_test.go, internal/worktree/lifecycle_policy_test.go
Covers: LGP18, LGP20, LGP21, LGP22, LGP23, LGP34

## What to build

A verb call carries its working directory, so the operand tests no longer change the process directory.

`internal/canonicalpath` gains a form of `Resolve` that takes the directory as an argument. It keeps exactly the semantics of `Resolve` (`internal/canonicalpath/canonicalpath.go:33-58`): it joins the parts as text and resolves symlinks before it cleans. `canonicalpath.Operand` is not that form.

`effects.go` reads the process directory once at the verb boundary. A failed read is carried forward, and it surfaces only when a relative operand needs it. So an absolute operand still works in a removed directory.

Route the directory to operand resolution: `resolveOperand` and `canonicalPath` (`internal/worktree/subshell.go:117-133`), `targetPath` (`path.go:293-307`), and the `resolveVerbOperand` callers in `land.go`, `reauthorize.go`, and `worktree.go`. Add a form beside `canonicalPath` rather than change its 22 call sites.

`verbCall` (`internal/worktree/verb_runner_test.go:47-55`) gains a working-directory field. Extend `checkVerbCall` and its refusal tests as the kit and clock fields do.

Move the 16 movable `chdir` tests in these files to `t.Parallel`: `identifier_operand_test.go`, `resume_reconcile_test.go`, `request_token_test.go`, `identity_component_test.go`, `land_resume_test.go`, `land_reauthorize_operand_test.go`, `path_identifier_test.go`, and `lifecycle_policy_test.go`. Record each test that moves and each test that stays.

`TestLandingHandsBackAPathWithNoCanonicalForm` (`land_resume_refusal_test.go:343`) stays serial. Put the LGP23 test in `serial_ceiling_test.go`, because `parallel_census_test.go` has 1096 lines and no grant.

Set both pins by the recorded arithmetic: the ceiling is the current value minus the moved tests plus the new serial tests. Update the four-class comment at `parallel_census_test.go:1033-1040` with no net line gain.

## Acceptance

- [ ] A relative operand resolves against the verb call's directory, not the process directory (LGP18).
- [ ] `TestParallelCensusOnTheLiveTree` passes with every moved test parallel (LGP20).
- [ ] The serial set equals the recorded ceiling, and the test count equals the recorded sum (LGP21, LGP22).
- [ ] The removed-directory test stays serial and keeps its real `filepath.Abs` failure (LGP23).
- [ ] The probe "leave the ceiling at its ticket-04 value" bites.
- [ ] The verification record holds the package time before and after the change (LGP34).
