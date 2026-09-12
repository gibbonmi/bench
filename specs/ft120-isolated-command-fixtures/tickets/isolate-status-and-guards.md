# Isolate status and guards command fixtures

Blocked by: none
Writes: cmd/bench, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: none

## What to build

Give the two command tests private repositories and private BENCH_HOME directories.
Resolve the kit asset source before changing the test directory.
Copy the shipped follow-on hook and its shipped harness configurations into the guards fixture.
Move the two tests out of the oversized main test file.

The coordinator reconciles the roadmap after this fix lands.

## Acceptance

- [ ] The status route test reads only its private repository and pool.
- [ ] The guards test reads only its private repository and pool.
- [ ] The guards test still verifies the shipped follow-on guard and both harness wires.
- [ ] Both tests work when the caller starts outside a Git repository.
