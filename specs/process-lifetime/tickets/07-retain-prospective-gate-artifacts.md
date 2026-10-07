# Retain prospective artifacts through gate and lane cleanup

Blocked by: 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/process_resources.go (new), docs/adr/0019-one-owner-holds-the-prospective-artifact-bundle.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/axi_query_registry_test.go, internal/conformance/injected_ports_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gate/engine.go, internal/gate/lane.go, internal/gate/lifetime_resources_test.go (new), internal/gate/phases.go, internal/gate/prospectiveartifact, internal/gate/run_transaction.go, internal/gate/runner.go, internal/systemtest/owner_artifact_recovery_test.go, internal/systemtest/process_lifetime_artifacts_test.go (new), internal/systemtest/process_lifetime_test.go (new), projects/benchkit.md, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/skill-description-budgets/budget-table-missing, tests/canary/skill-description-budgets/description-folded, tests/canary/skill-description-budgets/description-missing, tests/canary/skill-description-budgets/over-budget-command, tests/canary/skill-description-budgets/over-budget-description, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership, tests/canary/workflow-guidance-anchors/benchkit-system-suite-route
Covers: PL32, PL33, PL58, PL59, PL60, PL75, PL121, PL123

## What to build

A direct Close and a second-process sweep retain a prospective checkout with an unresolved registered group.
Open and publish protection before checkout use.
Keep ReadPublished, canonical repository binding, and Publish as prospective ownership authority.
The next owner schema references the protection generation and retains ambiguous legacy schema-1 records.

Install the prospective identity adapter in the established CLI table.
Sweep additionally requires current absence of the initiating bundle owner.
A live unused owner retains its checkout; a fresh unused closed bundle remains removable.
Direct Close uses the live object and observes the same durable user set.
No recovered PID grants signal or deletion authority.

Gate and lane consume the complete outcome and every required resource Close result.
A nested separate group prevents retained green evidence and lane pass even when its CLI exits zero.
Consume the lane selection-release error from ticket 4 and private-run cleanup from ticket 5.
Retain actual evidence paths and diagnostics for every unresolved obligation.

Reconcile ADR0019's retained-descendant exception with its recovery rule.
Retain the profile correction from ticket 1 without restating its ownership fact.
Preserve all anchor, canary, identity, phase-order, and evidence-reuse assertions.
Fixtures that formerly relied only on owner death need complete user proof or explicit retained disposition.

Keep over-budget source files from gaining lines at this checkpoint.
Move only cohesive outcome or phase composition into the co-owned runner or phases source when needed.
Do not defer headroom debt to a successor.

Before implementation starts, confirm complete accepted and landed durable-file-replacement delivery.
Its leaf, review-record caller, and native qualification must be complete.
Durable-caller-migration is not a prerequisite.
Use one retained integration source and a fresh author for this ticket.

Keep each caller's selected normal-exit policy.
A raw child status cannot certify required cleanup or publication.
Retain uncertain resources and every unresolved obligation.
Never signal from recovered records or add a blanket normal-exit kill.

This ticket belongs to PL-C3.
Finish the preceding chunk review before this chunk starts.
Its blockers include every prior overlapping Writes owner.
No successor caller must exist for this ticket's owned predicates.

BENCH_KIT is supplied by bench test --check system through the existing sealed system owner.
Do not introduce another subprocess fixture runner or private binary publisher.

## Acceptance

- [ ] PL32: An unresolved nested phase prevents retained green evidence.
- [ ] PL33: Incomplete cleanup prevents a lane pass record.
- [ ] PL58: Prospective Close retains an unresolved group user.
- [ ] PL59: A second-process sweep retains surviving prospective users.
- [ ] PL60: A legacy prospective record retains ambiguous cleanup proof.
- [ ] PL75: A fresh no-user prospective bundle remains removable.
- [ ] PL121: Sweep retains a live prospective owner with no group users.
- [ ] PL123: The lane consumes its selection-release failure.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/gate/prospectiveartifact`
- `bench test --package ./internal/gate`
- `bench test --check system`
- `bench test --check docs-currency-workflow`

Mutation witness: Restore owner-death-only sweep, ignore a nested registered group, or discard lane selection release failure. The real Close, second-process sweep, or green-record witness must fail.
