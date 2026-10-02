# 6. Land the stubbed landing tests for real

Blocked by: 5-fault-the-landing-follow-on-steps-with-real-fixtures.md
Writes: internal/worktree/land_flags_test.go, internal/worktree/land_reauthorization_test.go, internal/worktree/land_fixtures_test.go, internal/worktree/joins.go, internal/worktree/land.go, internal/worktree/land_resume.go, internal/worktree/land_identity.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS39, WS40, WS41, WS42, WS43, WS44

## What to build

Chunk: SR-C4.

Convert each user of `stubLandJoins` to a real landing with a shell gate, then delete
`stubLandJoins`. The six users are in `land_flags_test.go` and
`land_reauthorization_test.go`. `TestLandCommandPostCASTerminalTable` faults the marker
with the fixture of ticket 4 and the reconcile with the fixture of ticket 5.
`TestLandCommandProjectGreenOrderTable` asserts the real marker ref after the landing, in
place of the recorded arguments. If ticket 4 recorded a failed marker probe, each test in
this ticket keeps its `advanceLandingMarker` stub and converts its other stubs.

Remove `advanceLandingMarker`, `reconcileLanding`, `authorizeLandingSource`, and
`pruneLandedBranches` from the joins value. Keep each field whose probe failed in ticket 4
or ticket 5. The landing then calls the real functions directly.

Each converted test keeps its name.

For a field without a probe, a fixture that cannot make its test pass stops the build.
The build names that test, per decision 3.

## Acceptance

- [ ] The digest-shaped token, the abbreviated tip, and the abbreviated base tests pass through a real landing.
- [ ] The post-swap terminal table prints a row for each real fault.
- [ ] A release diagnostic in a real landing cannot forge a terminal line.
- [ ] A real landing advances the green marker to the published commit.
- [ ] `stubLandJoins` is gone, and the joins value declares no landing field whose probe turned red.
