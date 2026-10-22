# Carry the run value and read owner markers once

Blocked by: 03-add-the-git-run-value.md
Writes: internal/worktree/effects.go, internal/worktree/owner_marker.go (new), internal/worktree/classifier.go, internal/worktree/worktree.go, internal/worktree/identity_component.go, internal/worktree/ownership.go, internal/worktree/reset.go, internal/worktree/reauthorize.go, internal/worktree/verb_runner_test.go, internal/worktree/journey_test.go, internal/worktree/parallel_census_test.go, internal/worktree/admin_readers_test.go, internal/worktree/eligibility_test.go, internal/worktree/single_read_census_cases_test.go
Covers: LGP17, LGP21, LGP22, LGP34

## What to build

Each worktree verb run carries one git run value, and owner-marker validation reads the common directory once through it.

`ambient` carries the run value (`internal/worktree/effects.go:22-36`). `newAmbient` builds a fresh run value for each verb run, and so does `callAmbient` (`verb_runner_test.go:154-163`). Never hold the run value in `joins`, because a test reuses one `joins` value across verb runs (`verb_runner_test.go:142-146`).

Move `validateOwnerMarker` (`internal/worktree/classifier.go:396`, to the end of the function) to the new `internal/worktree/owner_marker.go`, and give it a run-value form. `classifier.go` has 463 lines and no grant, so the move must leave it shorter. Keep the current signature as a thin form for the test callers `admin_readers_test.go:102` and `eligibility_test.go:137`.

Give `canonicalRoot` (`internal/worktree/worktree.go:100`) a run-value form beside the current one. It reads `poolkey`'s run-value form from ticket 03.

Move the exec-target path onto the run value. `ownerMarkerRefusal` (`identity_component.go:69-73`) has three callers: `identity_component.go:99`, `ownership.go:204`, and `reset.go:158`. `reauthorize.go:83` also calls `validateOwnerMarker`. The functions around the three callers (`identityBundleRefusal`, `validateCreationBundle`, and `planReset`) hold no `ambient`, so each one builds a fresh run value on its call line. One validation then still starts the common-directory query once.

`reauthorize.go:83` sits in `reauthorizeParsed` (`reauthorize.go:64`), which holds no `ambient`. Its caller `reauthorizeWith` holds `a` (`reauthorize.go:50`), so pass the run value of `a` down to it.

A package-qualified constructor call in `effects.go`, such as `git.NewRun()`, joins the read set of `TestSingleReadCensusOnTheLiveTree` (`internal/worktree/single_read_census_test.go:109-155`). A later call of the same constructor in an unexported function then reds that census as a read below a census entry. Pick a constructor shape that the census does not count, such as a composite literal or a bare-name helper outside `effects.go`, and prove that the census stays green.

Write the LGP17 test first in `journey_test.go`, and record its red count before the change. Count only the `--git-common-dir` lines of the `journeyStubGit` log. One call under a canonical primary root reads 3 today. If the red shows another count, stop and report a wrong spec.

The stub mode decides the count. Use `fail-git-path` or `relative-git-path`, which pass `--git-common-dir` through to the real git. The mode `fail-git-dir` stops before the third read, and the default modes answer `worktree` with empty output, so both give 2. Build the fixture with `newOwnedAssignment` (`internal/worktree/resume_test.go:527`) before you bind the stub.

Assert the second acceptance item ("two verb runs") inside the LGP17 test or an existing test. A new top-level test changes the recorded pin arithmetic.

The new test uses `journeyStubGit`, which binds PATH, so it is serial. Set `worktreeTestCount` to 806 and `worktreeSerialCeiling` to 46, and record that arithmetic. Change only the two literals, because `parallel_census_test.go` has 1096 lines and no grant.

Update the synthetic copy of `effects.go` in `single_read_census_cases_test.go:15-29` if the `ambient` shape changes.

Record `bench test --full --package ./internal/worktree` before and after the change, for LGP34.

## Acceptance

- [ ] One `validateOwnerMarker` call under a canonical primary root starts the common-directory query exactly once (LGP17).
- [ ] Two verb runs in one test each build their own run value.
- [ ] The serial set equals 46, and the test count equals 806 (LGP21, LGP22).
- [ ] The probe "skip the run value in validateOwnerMarker" bites.
- [ ] The verification record holds the package time before and after the change (LGP34).
