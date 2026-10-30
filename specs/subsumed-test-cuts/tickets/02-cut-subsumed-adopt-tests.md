# Cut the subsumed adoption unit tests

Blocked by: none
Writes: internal/adopt/compatibility_test.go, internal/adopt/setup_test.go, internal/adopt/decision_test.go, internal/adopt/lifecycle_decision_test.go, internal/adopt/link_hook_test.go, internal/conformance/ordinary_build_census_test.go
Covers: TC1, TC10, TC11, TC12, TC13, TC14, TC15, TC16, TC17, TC18, TC19

## What to build

The adoption package loses five tests that broader tests subsume. The spec rows T1-2, T1-6, T1-9, T1-10, and T1-11 own this ticket.

- Remove `TestCompatibilityLegacyDoctor` (T1-2), `TestSetupGateDeclaresConsumerOnlyValidityScope` (T1-6), and `TestPrePushHookTopicPushIsSilent` (T1-11). Their kept tests already bit at the spec stage.
- T1-9 is a move row. First add an exact comparison with `wantSeededGateInputs` to `TestSetupTwiceLeavesGateInputsIdentical`. Probe `tc12` on it and see `bit`, then remove `TestSetupSeedsGateInputs`.
- T1-10 is a move row. First add three rows over the legacy inventory path to `TestCompleteLifecycleDecisions`: a duplicated current asset refuses, a desired asset with no fingerprint refuses, and an empty input plans no operation and no refusal. Probe `tc13`, `tc14`, and `tc15` on it and see `bit`, then remove `TestPlanLifecycleFromImmutableInventories`.
- The T1-10 cut empties `internal/adopt/decision_test.go`. Delete the file, and remove its key from `directArchitectureTests` in `internal/conformance/ordinary_build_census_test.go`.
- Remove each helper that only a removed test used. `runLegacyFixture` and `doctorResult` stay, because the Upgrade tests use them.

The build changes no production file. Run each probe through `bench probe` with `--expect` on the kept test, and read the `restored` cell of each answer.

## Acceptance

- [ ] The five candidate tests and `internal/adopt/decision_test.go` no longer exist.
- [ ] The completion-plan probes `tc1`, `tc10`, `tc11`, `tc16`, `tc17`, and `tc19` read `bit` on their kept tests.
- [ ] The completion-plan probe `tc12` reads `bit` on `TestSetupTwiceLeavesGateInputsIdentical` after the move.
- [ ] The completion-plan probes `tc13`, `tc14`, and `tc15` read `bit` on `TestCompleteLifecycleDecisions` after the move.
- [ ] `bench test --check ordinary-build-census` passes, and `directArchitectureTests` names no absent file.
- [ ] `bench test --package ./internal/adopt` passes.
