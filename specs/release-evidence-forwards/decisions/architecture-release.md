# Release forwards and capability scope (C11; FT364)

Status: ready

## Destination

Remove redundant internal forwards only after their consumers move to the evidence owner.
Retain tested publication behavior and keep unsupported npm staging explicit.

## Notes

The reviewer delegated architecture recommendations on 2026-10-06.
A forwarding export adds no independent policy.
A publication state machine owns ordering, authorization, and recovery.
Avoid: treating both as unused release machinery.
Use craft-spec and craft-seams for the narrow cut spec.

## Decisions so far

- [Which release machinery should remain?](architecture-release/tickets/1.md): keep the state machine and the tested capability set.
- [Which forwards are suitable for removal?](architecture-release/tickets/2.md): move internal callers to releaseevidence before cuts.
- [How should surfaces with unknown external consumers be treated?](architecture-release/tickets/3.md): retain public surfaces without a demonstrated removal case.

## Not yet specified

## Spec-writer discretion

- Mechanical import and private-name changes with unchanged record schemas.
- Caller migration order after the complete consumer census.

## Out of scope

- Deletion of the publication state machine or fixture staging.
- A claim that fixture tests qualify live npm behavior.
- A cut to a public CLI solely because no workflow calls it.
- The remaining non-release FT368 surfaces.
- A new publication, release, or external registry operation.

## Sources

- Path: `roadmap/FT364.md`
  Supports: the release cut proposal and prerequisite review.
  Drift: a change to its intended scope.
- Path: `roadmap/FT368.md`
  Supports: the broader public-surface decisions that remain separate.
  Drift: a change to those surface contracts.
- Path: `roadmap/FT142.md`
  Supports: the retained release residuals.
  Drift: resolution of the relevant residuals.
- Path: `internal/releasepreflight/types.go`
  Supports: the pass-through evidence exports.
  Drift: a change to their consumers or signatures.
- Path: `internal/publication/command.go`
  Supports: the actual adapter and staged-path behavior.
  Drift: a change to adapter selection or submission.
- Path: `internal/publication/command_adapter_test.go`
  Supports: staged refusal and real command journeys.
  Drift: a change to those tests.

