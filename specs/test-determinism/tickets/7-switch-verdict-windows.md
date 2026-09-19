# 7. Turn off verdict windows in a kit test run

Blocked by: 2-compose-kit-test-run.md
Writes: internal/bounds/bounds.go, internal/bounds/bounds_test.go, internal/env/, internal/git/git.go, internal/git/worktree_admin_enum_test.go, internal/gate/gate.go, internal/models/models.go, internal/sessioninspect/sessioninspect.go, internal/sessioninspect/sessioninspect_test.go, internal/guards/guards.go, internal/coverage/citation_execution.go, internal/refresh/refresh.go, internal/systemtest/owner_test.go, internal/systemtest/owner_land_race_test.go, internal/systemtest/owner_artifact_recovery_test.go, internal/conformance/bounds_policy_test.go, internal/conformance/fixture_bite_test.go, tests/canary/package-core-guard/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: TD33, TD34, TD35, TD36, TD37, TD38, TD39, TD42, TD43, TD44, TD45, TD46

## What to build

Chunk: TD-C3.

The bounds package owns the switch `BENCH_TEST_UNBOUNDED_WAITS`, whose only accepted value is exactly `1`. It owns the unbounded value and two window accessors. The first accessor returns the unbounded value for a policy window when the switch is on, and the policy window otherwise. The second accessor returns its argument unchanged, for a window that the switch must not change. Every bounds wait function treats the unbounded value as no deadline. The bounds package documentation states that the switch removes every verdict window, the 45-minute gate timeout included.

Each package variable that holds one of the seven verdict windows initializes through the first accessor. Session inspection gains a variable for its provider window and its discovery window, both through the first accessor, and a setter for tests. Each setter for tests keeps its raw assignment. The bounds-policy check's required list and owner table name all seven windows, the discovery window included. The owner table asks for the first accessor around each one.

Two new canaries in the `package-core-guard` family prove that rule: one for the worktree list window and one for the discovery window. The kit test run sets the switch in its entries.

The system suite removes the switch from each Bench child through one base-environment helper. `childEnvironment`, `systemStartSelected`, and `startArtifactLand` call that helper. The helper names the switch through the bounds package. The system suite files are system-tagged, so run them with `BENCH_KIT` set, through `bench test --check system`.

The switch census in the spec lists every test that the switch reaches. `TestEnvironmentPhaseTE15StopsAtDiscoveryBound` sets its own discovery window through the new setter. `TestWorktreeListTimeoutDefaultUsesPolicy` expects the first accessor's result.

## Acceptance

- [ ] With the switch at `1`, the first accessor returns the unbounded value for each of the seven verdict windows.
- [ ] With the switch absent, or at `0`, `true`, ` 1`, or `1` plus a newline, the first accessor returns the policy window.
- [ ] `bounds.Run` with the unbounded value completes a 200 ms child, and `bounds.Context` with it has no deadline.
- [ ] A test-set worktree list window and a test-set gate timeout still expire with the switch at `1`.
- [ ] With the switch at `1`, the second accessor returns the gate's cancel grace unchanged.
- [ ] The kit test run's entries carry `BENCH_TEST_UNBOUNDED_WAITS=1`.
- [ ] The bounds-policy check reds a verdict window variable that skips the first accessor. It also reds a session inspection file that reads the discovery constant directly.
- [ ] `TestSessionStartTE15BoundsDiscoveryAndContinues` passes with the switch at `1` in the test process.
- [ ] The whole gate is green with the switch on.
