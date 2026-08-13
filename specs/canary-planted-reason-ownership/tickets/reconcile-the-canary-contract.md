# Reconcile the complete canary contract

Blocked by: report-truthful-canary-inventory.md, remove-production-canary-dispatch.md, refuse-invalid-canary-inventory.md, reject-special-fixture-controls.md, reconcile-canary-inventory-consumers.md, advertise-canary-inventory-help.md, prove-every-retained-fixture.md, stop-seeding-linked-proof.md
Writes: `.agents/skills/bench-craft-gate/SKILL.md`, `.bench/BENCH-reference.md`, `projects/benchkit.md`, `docs/adr/0001-working-tree-gate-tripwire.md`, `docs/adr/0003-canary-latency-budget.md`, `docs/adr/0009-canary-concurrency-budget.md`, `CHANGELOG.md`

## What to build

Reconcile every current-state guidance and release-note surface with the landed contract: the kit's ordinary tests directly prove retained fixtures, the public and ship commands validate inventory only, and linked repos own native proof. Deprecate the two nested-sweep performance decisions and remove stale dispatch, sweep, inner-gate, worker, and fixture-count claims without reopening the retired architecture or FT168.

## Acceptance

- [ ] (covers DOC1) ADR 0001, `craft-gate`, and the project profile assign direct planted-reason proof to kit ordinary tests and assign linked repos only inventory validation plus project-native proof ownership.
- [ ] (covers DOC2) ADRs 0003 and 0009 are current-state records with `Status: deprecated` and no live nested-sweep latency, worker, inner-gate, or `GOMAXPROCS` policy.
- [ ] (covers DOC3) Production comments and fields, conformance rationale, system expectations, wrapper help, compatibility prose, shipped guidance, platform reference, profile, ADRs, and the Unreleased changelog use inventory/proof terms consistently and store no fixture count outside the omission-grading test expectation.
- [ ] (covers CC1) The exact composed candidate simultaneously reports inventory only, directly proves all 182 retained kit fixtures, preserves the two migrated release tripwires in owning-package tests, scaffolds no false linked-repo proof, and documents exactly those guarantees.
