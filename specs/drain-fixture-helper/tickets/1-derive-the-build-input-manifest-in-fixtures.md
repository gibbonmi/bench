# Derive the build-input manifest path and line in the kit test fixtures

Blocked by: none
Writes: cmd/bench/commands_brief_test.go, cmd/bench/command_registry_test.go, cmd/bench/build_subject_mode_test.go, internal/worktree/land_effects_test.go, internal/worktree/land_fixtures_test.go, internal/gate/prospective_owner_test.go, internal/preflight/binary_seal_test.go, internal/systemtest/owner_landing_fixture_test.go, internal/adopt/broker_test.go, internal/freshness/freshness_digest_test.go, internal/probe/probe_test.go, internal/runbinary/runbinary_test.go, internal/conformance/gate_entry_test.go
Covers: none

## What to build

The package `internal/freshness` owns the path of the build-input manifest and the grammar of one manifest line. `freshness.BuildInputsManifest` gives the path, and `freshness.BuildInputLine` gives one line. About 13 test files spell the path `scripts/go-build.inputs` or a manifest line as literal text. So a change to the owner leaves each fixture with a stale copy.

Make each fixture that writes or reads the manifest use the two owner names. Keep a literal only where it is a deliberate independent expectation that a named mutation of the owner turns red. Record each kept literal and its mutation in the commit message.

## Acceptance

- [ ] No fixture in the Writes list spells the manifest path or a manifest line as literal text, except a recorded independent expectation.
- [ ] Each kept independent expectation turns red under a named `bench probe` mutation of `internal/freshness`.
- [ ] Each affected package test passes.
