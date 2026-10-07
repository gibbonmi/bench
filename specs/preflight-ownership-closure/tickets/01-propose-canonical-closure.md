# Propose canonical closure with source reasons

Blocked by: none
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/references.go, internal/anchors/references_test.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/bounds/bounds.go, internal/bounds/bounds_test.go, internal/bounds/classify.go, internal/bounds/classify_nofollow_test.go, internal/canary/inventory.go, internal/canary/inventory_test.go, internal/canary/mutation.go, internal/canary/mutation_test.go, internal/canary/pins.go (new), internal/canary/pins_test.go (new), internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/coverage/citation_form_test.go, internal/coverage/citations.go, internal/coverage/citations_test.go, internal/coverage/coverage.go, internal/coverage/coverage_command_test.go, internal/coverage/coverage_schema_test.go, internal/coverage/coverage_test.go, internal/coverage/projection.go (new), internal/coverage/rows.go (new), internal/maps/fields.go, internal/maps/maps_parse_test.go, internal/maps/schema.go, internal/maps/tickets.go, internal/maps/tickets_test.go, internal/preflight/anchor_closure_test.go, internal/preflight/binary_seal_test.go, internal/preflight/charge_pack.go, internal/preflight/charge_test.go, internal/preflight/closure.go, internal/preflight/command.go, internal/preflight/command_bootstrap_test.go, internal/preflight/command_build_test.go, internal/preflight/command_review_test.go, internal/preflight/completion_plan_test.go, internal/preflight/decision.go, internal/preflight/decision_test.go, internal/preflight/fence_writes.go, internal/preflight/gather.go, internal/preflight/gather_inputs.go, internal/preflight/gather_test.go, internal/preflight/kit_pin_new_test.go, internal/preflight/ownership (new), internal/preflight/preflighttest/fixture.go, internal/preflight/preparation.go, internal/preflight/proposal.go, internal/preflight/proposal_command_test.go, internal/preflight/proposal_edges_test.go, internal/preflight/proposal_readonly_test.go, internal/preflight/proposal_test.go, internal/preflight/system_tag.go, internal/preflight/verdict_summary_test.go, internal/spec/fences.go, internal/spec/fences_test.go, internal/spec/resolve.go, internal/spec/resolve_nofollow_test.go, internal/tickets/registry_data.go, internal/tickets/registry_data_test.go, tests/canary/package-core-guard/bounds-classify-limit-restated, tests/canary/package-core-guard/bounds-discovery-window-unwrapped, tests/canary/package-core-guard/bounds-dot-import-package-alias, tests/canary/package-core-guard/bounds-dot-import-wait, tests/canary/package-core-guard/bounds-duplicate-owner, tests/canary/package-core-guard/bounds-intent-window-fixed, tests/canary/package-core-guard/bounds-multiple-dot-import-wait, tests/canary/package-core-guard/bounds-parenthesized-wait, tests/canary/package-core-guard/bounds-raw-elapsed-wait, tests/canary/package-core-guard/bounds-raw-injected-wait, tests/canary/package-core-guard/bounds-raw-wait-deadline, tests/canary/package-core-guard/bounds-raw-wait-duration, tests/canary/package-core-guard/bounds-read-limit-restated, tests/canary/package-core-guard/bounds-reassigned-wait-duration, tests/canary/package-core-guard/bounds-redeclared-wait-duration, tests/canary/package-core-guard/bounds-worktree-window-unwrapped, tests/canary/workflow-guidance-anchors/decision-map-asset-path
Covers: OC1, OC2, OC3, OC4, OC5, OC6, OC7, OC8, OC9, OC10, OC11, OC12, OC13, OC55, OC56, OC57, OC58, OC59, OC60, OC71, OC72, OC73, OC74, OC75, OC76, OC77, OC78, OC79, OC80, OC81, OC83, OC84, OC85, OC90, OC91, OC92

## What to build

The real --propose-writes entry reports transitive canonical requirements and every source reason through Gather, immutable Facts, Decide, and the existing closure-family ordering.
Expose ownership.Collect(root string, source Source, manifest Manifest, tickets []tickets.Ticket) (Snapshot, error).
The child imports no parent.
Source carries exact base, tip, selected SpecPath, and committed path/status projection.

Snapshot carries Requirements, Premises, Diagnostics, Complete, and ChargeReadiness.
ChargeReadiness carries InventoryComplete, CurrentComplete, RequiredCurrent, and PendingNative.
This first checkpoint exposes the value shape without activating successor premise grading.

Requirement carries ticket, originating entry, required path, kind, and all Provenance records.
Provenance carries canonical owner path, line, rule, reason, base, tip, source identity, and tree-or-kit scope.
Typed ownership.Error preserves Stage, Subject, Cause, and Unwrap.

Add typed provenance at tickets, anchors, and canary, with unchanged compatibility projections.
The registry owner and its tests join this first-use Writes; ticket 2 cannot supply ticket 1's accessor.
Canonical facts bind required path, owner, rule, selector, and actual source location through the established verified source context.
Keep registry rows, Go literal scanning, BASE includes, and fixture materialization single-sourced.

Repeated fixed-point traversal preserves independent reasons and family order and terminates cycles.
A compiled registry location needs the caller's independently verified kit source; an unavailable identity remains diagnostic.
Preserve genuinely absent inventory behavior while refusing present unreadable or special canary inventory at its canonical pin owner.
Unreadable anchors cannot supply partial green.

Parse the declared metadata through FieldScan's new FenceRole and FenceLabel, preserving every legacy FieldLine value.
Reject duplicate fences, keys, IDs, unknown fields, versions, enums, unclosed fences, and false openers inside another fence.
Use the existing classified bounded reader and coverage's typed projection.
An authored complete inventory grants no approval or semantic guarantee.
This checkpoint implements metadata syntax and canonical closure; successor premise certification does not activate charge refusal yet.

Render writes_proposal, ordering, closure, closure_evidence, source, and premises through the actual proposal command.
An absent legacy manifest stays incomplete and diagnostic in the authoring view.
Keep the response read-only and retain checkout-first, source, seal, evidence, and movement refusals.
Use actual committed source fixtures and independent second-hop and reason sentinels.
Use proposal_test.go, source_test.go, and premises_test.go in the ownership child.

Pay parent headroom in this ticket.
Move ticket probing into existing closure.go and system_tag.go, and pure write authorization into existing fence_writes.go.
Move the pure coverage parser/model into rows.go and its projection into projection.go.
Place any citations.go projection growth in that same extraction.

Move FixturePins enumeration into pins.go, with pins_test.go; reuse inventory and mutation readers.
Coverage's two new siblings take nine source files to eleven; canary's two take ten to twelve.
Add no file to the 44-file preflight root and no budget or accept-list grant.

Preserve all co-owned assertions and fixture purposes; holders needing no update stay byte unchanged.
Pin BENCH_KIT for every system-tagged closure fixture.
No successor grammar, helper, premise, or breakdown implementation is needed for this checkpoint.

Parse both explicit manual applicability forms into ManualClaim and NativeWitness.
Keep every manual item diagnostic and Complete false.
Current-premise rejects native-only fields; native-future requires its exact descriptor and observation unobserved.
Canonical witness ownership and charge applicability activate in ticket 7, after all required collectors exist.

This ticket belongs to OC-C1.
The explicit predecessor edge orders every overlapping write.
No unimplemented successor is required for this ticket's owned acceptance rows.

## Acceptance

- [ ] OC1: A missing fixture pin appears as a required path.
- [ ] OC2: A proposal leaves spec and ticket bytes unchanged.
- [ ] OC3: A required registry's own fixture enters the fixed point.
- [ ] OC4: A closure cycle terminates without dropping its required path.
- [ ] OC5: Two independent reasons for one path survive in evidence.
- [ ] OC6: An anchor reason names its actual holder and source line.
- [ ] OC7: A fixture include reason names the actual included BASE source.
- [ ] OC8: A registry reason names the canonical rule owner.
- [ ] OC9: A proposal prints the exact resolved source tip.
- [ ] OC10: An absent linked-kit anchor directory adds no invented holder.
- [ ] OC11: An absent linked-kit fixture inventory adds no invented fixture.
- [ ] OC12: A complete empty proposal prints typed zero-row tables.
- [ ] OC13: An incomplete collector never reports complete proof.
- [ ] OC55: A duplicate metadata fence refuses complete collection.
- [ ] OC56: Duplicate JSON keys refuse complete collection.
- [ ] OC57: Unknown metadata fields refuse complete collection.
- [ ] OC58: Duplicate item IDs refuse complete collection.
- [ ] OC59: An unsupported metadata version remains diagnostic.
- [ ] OC60: An invalid metadata enum remains diagnostic.
- [ ] OC71: A symlink metadata source refuses before opening.
- [ ] OC72: A FIFO metadata source refuses without blocking.
- [ ] OC73: A directory at the metadata file refuses collection.
- [ ] OC74: An oversized metadata source refuses collection.
- [ ] OC75: A control-byte ownership path retains the shared refusal.
- [ ] OC76: A repository-escaping path refuses authority.
- [ ] OC77: Segment-prefix collisions remain outside authorization.
- [ ] OC78: Closure output spills without losing evidence reasons.
- [ ] OC79: Dirty checkout wins over a missing selected ticket.
- [ ] OC80: Persistent source movement publishes no proposal.
- [ ] OC81: The current source-tip mismatch remains a refusal.
- [ ] OC83: An unreadable present fixture inventory remains diagnostic.
- [ ] OC84: An unreadable present anchor registry remains diagnostic.
- [ ] OC85: A newly required system test retains its BENCH_KIT obligation.
- [ ] OC90: Canonical evidence-source mismatch retains charge refusal.
- [ ] OC91: Metadata-looking text inside another fence is not an opener.
- [ ] OC92: Added fence roles preserve every legacy field projection.

## Verification

Run these future checks through their actual production entries.
Use a few sufficient multi-case witnesses; do not create a redundant test per row.
Retain the named omission mutation, its behavioral failure, restoration, and passing command result.
Compilation failure supplies no behavioral evidence.

- `bench test --package ./internal/preflight/...`
- `bench test --package ./internal/anchors`
- `bench test --package ./internal/canary`
- `bench test --package ./internal/coverage`
- `bench test --package ./internal/maps`
- `bench test --package ./internal/bounds`
- `bench test --package ./internal/spec`
- `bench test --package ./internal/tickets`
- `bench test --package ./cmd/bench`
- `bench test --check ticket-grammar`

Focused producer command: `bench test --package ./internal/preflight/... --run 'TestProposalFixedPoint|TestProposalProvenance|TestMetadataGrammar'`.
Named mutation identity: `oc-01-omission`.
This obligation uses the named production seam and observable omission in this ticket, not a compile-failure substitute.

Mutation obligation: Omit the second-hop pin, one independent reason, or the FieldScan fence role. The actual proposal fixture or legacy parity witness must fail.

The completion plan names this ticket's focused command and mutation identity.
Execute bench probe against the specified production seam after its code exists, then restore and rerun that exact focused command.
If the planned mutation cannot reach that producer, amend the enabling plan before verification; no substitute earns completion.

Read the canonical structure report before dispatch.
Remeasure Growth from the actual preceding tip and whole-tree crowding.
Keep cohesive extraction and necessary pin relocation in this same ticket.
