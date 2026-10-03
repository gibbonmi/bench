# 4. Interrupt the landing marker with a gate script

Blocked by: 3-pass-the-kit-value-to-merge-and-land.md
Writes: internal/worktree/land_marker_fixture_test.go (new), internal/worktree/land_specless_test.go, internal/worktree/land_resume_test.go, internal/worktree/land_local_capture_test.go, internal/worktree/land_tickets_only_test.go, internal/worktree/identity_component_test.go, internal/worktree/land_reauthorization_test.go, internal/worktree/land_fixtures_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS26, WS27, WS28, WS29, WS30, WS31, WS32, WS33

## What to build

Chunk: SR-C4.

Add a landing fixture whose gate script deletes `refs/bench/green/<branch>` in the
destination repository during the gate run. The real landing then publishes, and the
real marker swap fails, because the expected prior tip is gone. Put the fixture in a new
test file.

First run the probe. Convert `TestResumeLandCommandCompletesAnInterruptedMarker`, then run
`bench probe` on `land.go` with the marker error branch omitted, scoped to that test. If
the test turns red, convert the other marker tests in this ticket. Record the probe
command and its red in the verification note.

A probe fails when it stays green, or when the fixture cannot make the converted test pass. After a failed
probe, keep `advanceLandingMarker` and its tests unchanged. Record one `bench learning`
entry that names the field and the probe, and end this ticket's conversion without a build
stop.

Each converted test keeps its name. `identity_component_test.go` is over its line budget,
so the conversion does not grow it. The field itself leaves in ticket 6, because
`stubLandJoins` still sets it.

## Acceptance

- [ ] The marker probe turns `TestResumeLandCommandCompletesAnInterruptedMarker` red, and the verification note records it.
- [ ] Each of the eight marker tests interrupts a real landing with the gate-script fixture and passes.
- [ ] No marker test in this ticket sets `advanceLandingMarker`.
- [ ] If the probe failed, the learning entry exists and the eight tests are unchanged.
