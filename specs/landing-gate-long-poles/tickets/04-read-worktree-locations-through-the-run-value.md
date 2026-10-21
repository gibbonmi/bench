# Read worktree locations through the run value

Blocked by: 03-add-the-git-run-value.md
Writes: internal/worktree/effects.go, internal/worktree/joins.go, internal/worktree/subshell.go, internal/worktree/classifier.go, internal/worktree/worktree.go, internal/worktree/identity_component.go, internal/worktree/ownership.go, internal/worktree/reset.go, internal/worktree/reauthorize.go, internal/worktree/journey_test.go, internal/worktree/parallel_census_test.go, internal/worktree/admin_readers_test.go, internal/worktree/eligibility_test.go, internal/worktree/single_read_census_cases_test.go
Covers: LGP16, LGP17, LGP21, LGP22, LGP34

## What to build

Worktree commands read the common directory once per run through the git run value of ticket 03. The run value is built once at the verb boundary and travels with the per-run seam.

`planExplicitWith` (`internal/worktree/subshell.go:160`) already takes `joins`, so the run value can travel there with no signature change. Its seven callers then stay unchanged.

Give `validateOwnerMarker` (`internal/worktree/classifier.go:396`) and `canonicalRoot` (`internal/worktree/worktree.go:100`) run-value forms beside the current ones. Move the exec-target path (`identity_component.go:69-73`) to the run-value form, so each exec target costs fewer starts.

Count only the `--git-common-dir` lines of the `journeyStubGit` log (`internal/worktree/journey_test.go:64-68`). Write each new test first, and record its red count before the change. An owned explicit plan reads 8 today, and one `validateOwnerMarker` call reads 3. If a red count differs, stop and report a wrong spec.

Call `planExplicitWith` directly, for an owned target under a canonical primary root. The two ledger reads through `intent.Address` stay live.

`internal/worktree/classifier.go` has 463 lines and no grant, and `clean_landed.go` has exactly 400. Neither file may gain net lines.

The two new tests use `journeyStubGit`, which binds PATH, so both are serial. Set `worktreeTestCount` to 807 and `worktreeSerialCeiling` to 47, and record that arithmetic. Change only the two literals, because `parallel_census_test.go` has 1096 lines and no grant.

Record `bench test --full --package ./internal/worktree` before and after the change, for LGP34.

## Acceptance

- [ ] A direct `planExplicitWith` call for an owned target under a canonical primary root starts the common-directory query exactly three times (LGP16).
- [ ] One `validateOwnerMarker` call under a canonical primary root starts the common-directory query exactly once (LGP17).
- [ ] The serial set equals the recorded ceiling, and the test count equals the recorded sum (LGP21, LGP22).
- [ ] The probe "skip the run value in validateOwnerMarker" bites.
- [ ] The verification record holds the package time before and after the change (LGP34).
