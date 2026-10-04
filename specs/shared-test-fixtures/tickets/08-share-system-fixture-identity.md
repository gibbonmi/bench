# Share the system fixture identity

Blocked by: 07-share-worktree-fixture-identity.md
Writes: internal/systemtest/bench_follow_on_test.go, internal/systemtest/charge_evidence_test.go, internal/systemtest/land_route_test.go, internal/systemtest/otel_gate_test.go, internal/systemtest/otel_verbs_test.go, internal/systemtest/owner_artifact_recovery_test.go, internal/systemtest/owner_land_race_test.go, internal/systemtest/owner_landing_fixture_test.go, internal/systemtest/owner_selection_test.go, internal/systemtest/owner_test.go, internal/systemtest/status_route_converge_test.go
Covers: GF18, GF29, GF30

## What to build

Use the GF-C1 identity values in system fixtures and generated scripts. Keep systemGit, systemGitOutput, and owner.runAt as the process and repository owners. Keep the selected executable ledger and teardown behavior.

Create an owned system repository and a landing fixture. Both use the canonical default identity through owner.runAt, within the existing repository budget.

The tagged system suite receives BENCH_KIT from `bench test --check system`. Do not start an unsealed system binary.

## Acceptance

- [ ] System Git children retain owner.runAt and its repository budget (GF18).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --check system`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
