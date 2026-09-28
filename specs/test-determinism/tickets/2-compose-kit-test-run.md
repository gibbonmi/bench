# 2. Compose the kit test run into bench test and release preflight

Blocked by: 1-open-kit-test-run.md
Writes: internal/gittest/gittest.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/probe/probe_test.go, internal/testreport/environment.go, internal/testreport/environment_test.go, internal/releasepreflight/command.go, internal/releasepreflight/external_test.go, internal/env/, internal/gate/phases.go, internal/testreport/command.go, internal/testreport/check_test.go, internal/testreport/testreport_test.go, internal/testreport/outcome_test.go, internal/testreport/selection_test.go, internal/testreport/cancel_test.go, internal/gate/phases_test.go
Covers: TD7, TD16, TD17, TD18, TD19

## What to build

Chunk: TD-C1b.

The five registry files join the fence through the binding registry closure only.
This ticket expects no edits to them.

`bench test` opens a kit test run when its Go child runs in the kit, and it closes the run after the child exits. The runner derives the Bench build cache from the operator's `HOME` before it merges the entries, so the child keeps a warm cache. Each release preflight external phase opens a run, because release preflight grades the kit alone.

After both runners use the run, remove the git test configuration function and any gate helper that only forwarded it. The git maintenance policy then has one source, the kit test run. The owner's documentation states that a plain `go test` opens no run.

The command owner keeps the run alive until its child exits.
The shared probe checks that Go telemetry is off and its configuration belongs to the run.
The repair preserves the empty HOME and TMPDIR guarantees.
The affected Go fixtures forward environment queries to Go and retain their test output.
The linked-root phase test replaces the retired forwarding-helper test.

## Acceptance

- [ ] The `bench test` child for the kit runs the kit-run probe with exit 0.
- [ ] The `bench test` child carries `GOCACHE` equal to the Bench cache that the operator's `HOME` derives.
- [ ] A release preflight external phase runs the kit-run probe with exit 0.
- [ ] No non-test Go file outside the owner's file names `maintenance.auto`.
- [ ] The owner's documentation states that a plain `go test` is outside the policy.
