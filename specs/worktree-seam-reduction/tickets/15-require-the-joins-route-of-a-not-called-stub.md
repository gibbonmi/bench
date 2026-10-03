# 15. Require the joins route of a not-called stub

Blocked by: 14-name-each-table-in-production.md
Writes: internal/worktree/verb_runner_test.go, internal/worktree/land_identity_test.go, internal/worktree/verb_result_route_test.go (new), internal/worktree/reset_apply_test.go, internal/worktree/land_effects_test.go, internal/worktree/land_flags_test.go, internal/worktree/parallel_census_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS75, WS76, WS77, WS78, WS79, WS80

## What to build

Chunk: SR-C7.

Add a `viaJoins` field to the verb result. The verb runner sets it when the run took the
internal form with the call's joins value. Add a `mustViaJoins` form that fails the
test when the field is false. The name stays apart from the `joined` field of `verbForm`. Put the new runner test in a new test file, because
`verb_runner_check_test.go` is near its line budget.

Make each test that asserts that a joins stub was not called also call `mustViaJoins`
on that run. The tests are `TestResetApplyTakesTheCleanupLock`,
`TestLandSkipsTheRefreshWithoutBuildInputs`, `TestLandSkipsAFreshBroker`,
`TestResumeReadsEffectStateFromTheTree`, `TestLandCommandHostileSourceInputsRefuseBoundedly`,
`TestLandCommandRefusesDestinationAndSourceStateBeforeGate`, and the six identity-refusal tests in
`land_identity_test.go`.

For each named test, run `bench probe` that drops the joins value from the not-called run.
Record each probe command and its red in the verification note. Raise
`worktreeTestCount` by the number of new top-level tests in this commit.

## Acceptance

- [ ] A joins-form run sets `viaJoins`, and a public-entry run leaves it false.
- [ ] Each named not-called test turns red when its not-called run drops the joins value.
- [ ] The package declares exactly `worktreeTestCount` top-level tests at this commit.
