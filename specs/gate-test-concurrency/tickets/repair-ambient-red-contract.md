# Repair the ambient ticket's red contract

Blocked by: none
Ownership fence: `specs/gate-test-concurrency/tickets/repair-ambient-kit-migrations.md`
Integration surfaces: `internal/gate/manifest.go` manifest-present branch→`specs/gate-test-concurrency/tickets/repair-ambient-kit-migrations.md` + RAT1; `internal/gate/component_decision.go` no-module scoping guard→`specs/gate-test-concurrency/tickets/repair-ambient-kit-migrations.md` + RAT1; `internal/gate/evaluation.go` prospective empty-scoping branch→`specs/gate-test-concurrency/tickets/repair-ambient-kit-migrations.md` + RAT1; refreshed repair-loop assignment→repair-ambient-kit-migrations.md
Contracts: the mutation claims in `specs/gate-test-concurrency/tickets/repair-ambient-kit-migrations.md` distinguish behaviorally kit-sensitive entries from present-manifest, no-module, and prospective paths, asserted by RAT1 against the real producer branches
Closure: RAT1/present-manifest-audit, RAT1/no-module-exemption, RAT1/prospective-exemption, RAT1/refresh-metadata

## What to build

Correct the blocked ticket after its two prescribed hostile-environment
mutations stayed green. Its exported-call audit must retain calls whose
throwaway roots have no `go.mod` and prospective calls whose evaluation returns
empty scoping. Its absent-tool change is a single-source cleanup: the fixture
already owns its kit identity, but its present manifest means `phaseTable`
cannot observe the kit argument.

Replace the false behavioral mutation claims with attributable checks. A
kit-sensitive exported-entry mutation must first make the audited fixture reach
a real kit consumer. Reintroducing the redundant `kitRoot` read in the
absent-tool helper is instead owned by an exact source audit; do not claim its
manifest-backed focused test becomes behaviorally red.

## Acceptance

- [ ] [RAT1] (covers local) every red claim in the ambient migration ticket reaches a present consumer, the manifest-backed cleanup is labeled as a source audit, and the blocked assignment can refresh without changing its ownership fence.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RAT1/present-manifest-audit | restore the claim that the manifest-backed absent-tool helper becomes behaviorally red from ambient kit alone | the manifest-source and hostile-control audit | apply, compare the claim to `manifest.go`, run the focused hostile control green, and require the audit to name the claim as unattributable |
| RAT1/no-module-exemption | claim that one audited no-module fixture consumes kit identity | the component-scoping guard audit | apply, compare the fixture to `component_decision.go`, and require the audit to name the missing `go.mod` guard |
| RAT1/prospective-exemption | claim that one audited prospective call consumes kit identity | the prospective-scoping guard audit | apply, compare the claim to `evaluation.go`, and require the audit to name the prospective early return |
| RAT1/refresh-metadata | omit the target ticket's reciprocal blocker or change its ownership fence | the lifecycle refresh preconditions | apply, run the same `assign --refresh` request, and expect the metadata or fence refusal before any worktree mutation |
