# Adoption lifecycle consumption (FT217; C06)

Status: ready

## Destination

Make adoption preview and execution consume the existing lifecycle decision.
Use FT217's closed preserving-refactor contract.

## Notes

The reviewer delegated architecture recommendations on 2026-10-06.
FT217 already owns the architecture decision.
A lifecycle decision classifies immutable inventory facts.
Avoid: a second filesystem transaction.
Use craft-spec and craft-seams for its spec.

## Decisions so far

- [Should the existing lifecycle decision be replaced?](architecture-adoption/tickets/1.md): deepen PlanLifecycle through real consumers.
- [What prevents a refactor from changing adoption behavior?](architecture-adoption/tickets/2.md): keep FT217's unchanged-assertion exit proof.

## Not yet specified

## Spec-writer discretion

- Private names and file placement within the decided modules.
- Ticket slices that preserve the existing command behavior.

## Out of scope

- A new dry-run command surface.
- A second lifecycle planner.
- Changes to pre-existing assertions or expected outcomes during the preserving refactor.
- Generic fixture migration already owned by FT360 or landing-test-efficiency.

## Sources

- Path: `roadmap/FT217.md`
  Supports: the closed behavior, route, and exit proof.
  Drift: an approved change to that decision.
- Path: `internal/adopt/decision.go`
  Supports: the existing immutable inventory decision.
  Drift: a change to PlanLifecycle.
- Path: `internal/adopt/upgrade.go`
  Supports: the current decision consumer in upgrade counts.
  Drift: a change to upgradePlanCounts.
- Path: `internal/adopt/link_transaction.go`
  Supports: execution classification outside the decision.
  Drift: a change to managed-asset classification.
- Path: `internal/adopt/unlink.go`
  Supports: unlink preservation and removal classification.
  Drift: a change to unlink plan construction.

