# Propose helper callers and exact before-state ownership

Blocked by: 02-classify-grammar-effects.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/commitment/repository/light_path.go, internal/commitment/repository/light_path_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/consumers/blast.go, internal/consumers/blast_test.go, internal/consumers/loader.go, internal/consumers/loader_test.go, internal/consumers/resolve.go, internal/consumers/rows.go, internal/preflight/closure.go, internal/preflight/decision.go, internal/preflight/deletion_preflight_test.go, internal/preflight/fence_writes.go, internal/preflight/fence_writes_test.go, internal/preflight/gather.go, internal/preflight/ownership (new), internal/preflight/preflighttest/fixture.go, internal/preflight/proposal.go, internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/recordcmd/command_test.go, internal/tickets/tickets.go, internal/tickets/writes.go, internal/tickets/writes_test.go
Covers: OC20, OC21, OC22, OC23, OC24, OC25, OC26, OC27, OC28, OC29, OC30, OC31, OC82

## What to build

The real proposal lists all static production and test callers of a changed unexported helper, retaining before-state users after relocation or deletion.
Consume the established ownership and effect contracts.
Expose consumers.Load(root string) ([]*Package, error) in existing loader.go, backed by loadPackages, Resolve, and Rows.
Return complete typed packages or an error in the default clean committed context.

Ownership retains immutable resolved identities, not mutable Package objects.
Reuse the existing before-state AST machinery.
Reflection, opaque references, and non-default contexts remain diagnostic; add no historical arbitrary-context loader or shell-output parser.

Expose tickets.ParseWrite for Existing, New, and Deleted; keep WritesPath as its compatibility projection.
The exact path (deleted) marker needs one tracked before-state file at Source.Base.
Accept it before implementation and after removal.
Refuse nonexistent before paths, deleted directory shorthand, and combined markers.
Preserve the legacy exact committed D-status case.
Authorization, union, overlap, closure, commitment, and record consumers share this marker rule.

Add exact reused-owner sources and complete retirement, posture-caller, and printed-field producer censuses.
Omitted consumers or producers remain diagnostic.
A planned new helper has no invented existing callers.
Use helpers.go, deletion.go, and TestHelperCallers, TestDeletedEntry, and TestOwnerPremises in effects_test.go.
Drive actual before/after commits, the third helper caller, and both present-before and absent-tip deletion shapes.
The legacy D-status differential must remain green.

The loader wrapper fits its existing 60-line owner without a new root sibling.
Mandatory premise charge policy remains inactive until ticket 7.
No successor collector is required for these proposal predicates.

This ticket belongs to OC-C2.
The explicit predecessor edge orders every overlapping write.
No unimplemented successor is required for this ticket's owned acceptance rows.

## Acceptance

- [ ] OC20: A changed unexported signature proposes each static test caller.
- [ ] OC21: A fixture-helper change proposes its external call-site file.
- [ ] OC22: A relocation retains the before-state helper's caller proof.
- [ ] OC23: A deleted helper retains its before-state caller proof.
- [ ] OC24: A deleted marker resolves an exact tracked before path.
- [ ] OC25: A nonexistent before path refuses the deleted marker.
- [ ] OC26: A planned deleted marker accepts a still-present before path.
- [ ] OC27: Combined new and deleted markers refuse parsing.
- [ ] OC28: A reused owner adds its exact rule source to closure.
- [ ] OC29: A liveness premise reports an omitted retirement consumer.
- [ ] OC30: A posture premise reports an omitted static writer caller.
- [ ] OC31: A printed-field premise reports an omitted producer.
- [ ] OC82: Existing committed D-status ownership remains accepted.

## Verification

Run these future checks through their actual production entries.
Use a few sufficient multi-case witnesses; do not create a redundant test per row.
Retain the named omission mutation, its behavioral failure, restoration, and passing command result.
Compilation failure supplies no behavioral evidence.

- `bench test --package ./internal/preflight/...`
- `bench test --package ./internal/consumers`
- `bench test --package ./internal/tickets`
- `bench test --package ./internal/commitment/repository`
- `bench test --package ./internal/reviewrecord/recordcmd`
- `bench test --package ./cmd/bench`
- `bench test --check ticket-grammar`

Focused producer command: `bench test --package ./internal/preflight/... --run 'TestHelperCallers|TestDeletedEntry|TestOwnerPremises'`.
Named mutation identity: `oc-03-omission`.
This obligation uses the named production seam and observable omission in this ticket, not a compile-failure substitute.

Mutation obligation: Sample only two callers, erase base users after relocation, or accept a never-tracked deletion. The third-file or before-state sentinel must fail.

The completion plan names this ticket's focused command and mutation identity.
Execute bench probe against the specified production seam after its code exists, then restore and rerun that exact focused command.
If the planned mutation cannot reach that producer, amend the enabling plan before verification; no substitute earns completion.

Read the canonical structure report before dispatch.
Remeasure Growth from the actual preceding tip and whole-tree crowding.
Keep cohesive extraction and necessary pin relocation in this same ticket.
