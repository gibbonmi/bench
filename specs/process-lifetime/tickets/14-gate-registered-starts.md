# Gate every registered start

Blocked by: 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md
Writes: internal/subprocess/control.go, internal/subprocess/startgate.go (new), internal/subprocess/startgate_test.go (new), internal/subprocess/lifetime_test.go, internal/subprocess/protection/recovery.go, internal/subprocess/protection/recovery_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/injected_ports_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gate/lifetime_resources_test.go
Covers: PL25, PL128, PL129

## What to build

The reviewer decided on 2026-10-23 to build the start gate in this build.
A registered launch starts its program through a start gate in the program's new group.
The gate waits on an inherited pipe and executes the program only after the open signal.
The actor sends that signal only after Started succeeds in every registered store.
At end of file, the gate exits and executes nothing.

An exec failure behind the gate remains a Start error with its existing meaning.
Resolve the program path before the gate starts.
A launch with no registration uses no start gate and keeps its current path.
The gate must work for every actor that installs a protocol, including a Go test binary.

Recovery grades a pending user whose actor is currently absent as no-start.
A pending user whose actor answers a probe stays protected.
Recovery still sends no signal and executes no recovered content.
The user record stays schema 1.

Keep cancellation, escalation, single Wait, and stream ownership as PL2-PL14 and PL124-PL125 state.
Keep each caller's selected normal-exit policy.

This ticket belongs to PL-C3.

## Acceptance

- [ ] PL25: A pending user resolves as no-start when its actor is absent and stays protected while its actor answers.
- [ ] PL128: A gated program never runs when its actor dies before the open signal.
- [ ] PL129: An exec failure behind the start gate remains a Start error.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.

- `bench test --package ./internal/subprocess/...`
- `bench test --package ./cmd/bench`
- `bench test --package ./internal/anchors`
- `bench test --package ./internal/conformance`
- `bench test --check root-conformance`
- `bench test --check system`
- `bench test --package ./internal/gate`

Mutation witness: Execute the program at end of file, or grade a pending user with an absent actor as protected. The ActorDeathBeforeOpen or PendingActorAbsent witness must fail.
