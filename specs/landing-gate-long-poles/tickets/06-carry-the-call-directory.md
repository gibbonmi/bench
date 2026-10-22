# Carry the call directory

Blocked by: 05-read-the-explicit-plan-locations-once.md
Writes: internal/canonicalpath/canonicalpath.go, internal/canonicalpath/canonicalpath_test.go, internal/worktree/effects.go, internal/worktree/subshell.go, internal/worktree/resume.go, internal/worktree/ownership.go, internal/worktree/land.go, internal/worktree/land_resume.go, internal/worktree/reauthorize.go, internal/worktree/worktree.go, internal/worktree/path.go, internal/worktree/tree_target.go, internal/worktree/verb_runner_test.go, internal/worktree/verb_runner_check_test.go, internal/worktree/single_read_census_cases_test.go, internal/worktree/identifier_operand_test.go, internal/worktree/parallel_census_test.go
Covers: LGP18, LGP22, LGP34

## What to build

A verb run carries its working directory, and a relative operand resolves against it.

`internal/canonicalpath` gains a form of `Resolve` that takes the directory as an argument. It keeps exactly the semantics of `Resolve` (`internal/canonicalpath/canonicalpath.go:46-59`): it joins the parts as text and resolves symlinks before it cleans. `canonicalpath.Operand` is not that form, because it cleans lexically.

`effects.go` reads the process directory once at the verb boundary into `ambient`. A failed read is carried forward, and it surfaces only when a relative operand needs it. So an absolute operand still works in a removed directory.

Route the directory to each relative-operand site:

- `canonicalPath` and `resolveOperand` (`internal/worktree/subshell.go:117-138`);
- `canonicalPath` in `releaseAssignment` (`ownership.go:217`);
- `canonicalPath` in `resumeLandWith` (`land_resume.go:36`);
- `canonicalPath(target)` in the `cleanCommandWith` error path (`worktree.go:354`);
- the `resolveVerbOperand` callers in `land.go`, `reauthorize.go`, and `worktree.go`.

`planExplicitWith` (`subshell.go:160`) and `cleanupIdentity` (`resume.go:22-23`) call `resolveOperand` and hold no `ambient`. Do not add a directory parameter to them. The verb boundary that holds `ambient` makes a relative path operand absolute against the call directory before the call, for example at `worktree.go:275`. `resolveOperand` then sees an absolute path, which `canonicalpath.Resolve` leaves unchanged.

Leave an operand that starts with `~` unchanged at the boundary, because `expandHomeTarget` (`path.go:276`) owns its refusal text. A ledger-fed caller passes a recorded absolute path, and the public entries `PlanExplicit`, `PlanAutomatic`, and `ApplyExplicit` keep the process directory. A carried directory-read failure surfaces with the same refusal text as a failed process-directory read today.

Add a form beside `canonicalPath` rather than change its 22 call sites. `targetPath` (`path.go:293-309`) already refuses a relative path, so it needs no directory. If its signature changes anyway, update `internal/worktree/tree_target.go:64`.

The new process-directory read in `effects.go` joins the read set of the single-read census. Update the synthetic copy of `effects.go` in `single_read_census_cases_test.go:15-29` to match.

`land.go` has 397 lines, so only 3 lines of headroom remain. `verb_runner_check_test.go` has 378 lines, and ticket 08 also extends it.

`verbCall` (`verb_runner_test.go:47-55`) gains a working-directory field. Extend `checkVerbCall` and its refusal tests as the kit and clock fields do. A verb with no joins form (create, reclaim, list, path, show, exec, pool, lease-file) refuses the field.

Put the LGP18 test in `identifier_operand_test.go`. It calls `t.Parallel`, so `worktreeTestCount` rises by one and the ceiling stays unchanged. Record that arithmetic.

## Acceptance

- [ ] A relative operand resolves against the verb call's directory, not the process directory (LGP18).
- [ ] An absolute operand still resolves when the process directory is gone.
- [ ] The test count equals the recorded sum (LGP22).
- [ ] The probe "resolve a relative operand against the process directory" bites.
- [ ] The verification record holds the package time before and after the change (LGP34).
