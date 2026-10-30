# Cut the census, effects, and freshness tests

Blocked by: none
Writes: internal/worktree/land_census_test.go, internal/worktree/land_census_output_test.go, internal/worktree/land_effects_test.go, internal/worktree/land_freshness_test.go, internal/worktree/parallel_census_test.go
Covers: TC25, TC26, TC29, TC30, TC31, TC32, TC33

## What to build

The worktree package loses five top-level landing tests that broader tests subsume or that a merge replaces. The spec rows WT-2, WT-3, WT-6, and WT-8 own this ticket.

- Remove `TestLandCommandStatesZeroForAnAssignmentWithNoRecords` (WT-2). Its kept test already bit at the spec stage.
- WT-3 is a move row. First make `TestLandSkipsTheRefreshWithoutBuildInputs` assert that the effects table is the tail before the landed record, and that the landed record is the last stdout line. Probe `tc26` on it and see `bit`, then remove `TestLandEffectsRowPrecedesTheLandedRecord`. `land_effects_test.go` is at its 400-line budget, so the file must not grow.
- WT-6 is a move row. First make `TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion` assert that the released checkout is absent after the landing. Probe `tc29` on it and see `bit`, then remove `TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease`. Its dead process check goes with it.
- WT-8 is a merge. `TestLandCommandStatesTheCensusCountAndDropsTheRecords` seeds the head records and the output records before one landing. It asserts the count suffix, the dropped record file, the surviving handoff sections, the head breakdown, and the output breakdown. Remove `TestLandCommandPrintsTheCensusHeadBreakdown` and `TestLandingPrintsOutputBreakdown`, and delete the then-empty `land_census_output_test.go`.
- Lower `worktreeTestCount` by 5 from its value at this ticket's base. Rewrite the doc sentence of `TestPackageTestCountPin` that says the pin proves no test is removed or merged for wall-clock, so that it states the exact-count purpose. Add no line to `parallel_census_test.go`, because it is over its line budget.
- Remove each helper that only a removed test used. Slicing research found none.
- Keep moved assertions inside the package censuses: use the existing `wantEffects` helper, and write no literal that opens with a production table name and `[`. Add no `exec.Command`, `Setenv`, or `Chdir` call outside the harness files.

The build changes no production file. Run each probe through `bench probe` with `--expect` on the kept test, and read the `restored` cell of each answer.

## Acceptance

- [ ] The five candidate tests and `land_census_output_test.go` no longer exist.
- [ ] The completion-plan probe `tc25` reads `bit` on `TestLandCommandPrintsNoHeadsLineWithoutRecords`.
- [ ] The completion-plan probe `tc26` reads `bit` on `TestLandSkipsTheRefreshWithoutBuildInputs` after the move.
- [ ] The completion-plan probe `tc29` reads `bit` on `TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion` after the move.
- [ ] The completion-plan probes `tc30`, `tc31`, `tc32`, and `tc33` read `bit` on the merged census test.
- [ ] `TestPackageTestCountPin` passes, and no line count of `parallel_census_test.go` or `land_effects_test.go` grows.
- [ ] `bench test --package ./internal/worktree` passes.
