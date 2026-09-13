# Own gate fixture inputs

Blocked by: none
Writes: internal/testrepo/gate_fixture.go (new), internal/testrepo/gate_fixture_test.go (new), internal/commit/format_test.go, internal/gate/run_outcomes_test.go, internal/landing/completion_evidence_test.go, internal/shift/fault_test.go, internal/shift/shift_test.go, internal/status/status_producible_test.go, internal/systemtest/owner_landing_fixture_test.go, internal/worktree/delegated_integration_test.go, internal/worktree/land_fixtures_test.go, internal/worktree/land_journey_test.go
Covers: HP1, HP2, HP3

## What to build

Add one test-repository owner that writes a gate fixture's ordinary script,
prospective script, and input manifest from one declaration. A caller obtains an
ambient command token from this owner, and that same request adds the command to
the manifest. The owner can reuse one body for both script paths.

Move the gate fixtures that invoke ambient commands through this owner.
Give each fixture a private command path made only from its declarations.
An omitted declaration then fails where the fixture runs.
Keep absolute test executables and shell built-ins outside the ambient-command list.

## Acceptance

- [ ] One fixture declaration writes both requested gate scripts and one canonical gate-input manifest.
- [ ] Each command token used in a migrated script appears once in the manifest's sorted tools list.
- [ ] A migrated fixture cannot run an ambient command that it did not obtain from the owner.
- [ ] The affected landing, gate, shift, status, and commit fixtures pass with only their declared commands on the fixture path.
