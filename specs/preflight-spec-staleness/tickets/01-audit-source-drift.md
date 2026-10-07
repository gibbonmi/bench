# Audit exact source drift through the real command

Blocked by: none
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/preflight_version_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/preflight/binary_seal_test.go, internal/preflight/charge_pack.go, internal/preflight/charge_test.go, internal/preflight/command.go, internal/preflight/command_build_test.go, internal/preflight/completion_plan_test.go, internal/preflight/decision.go, internal/preflight/decision_test.go, internal/preflight/evidencecmd/evidence_grammar_test.go, internal/preflight/evidencecmd/operations.go, internal/preflight/evidencecmd/operations_test.go, internal/preflight/explicit_base_test.go, internal/preflight/gather.go, internal/preflight/gather_inputs.go, internal/preflight/gather_test.go, internal/preflight/preflighttest/fixture.go, internal/preflight/source_tip_test.go, internal/preflight/staleness (new), internal/preflight/verdict_summary_test.go
Covers: SS1, SS2, SS3, SS4, SS5, SS6, SS7, SS8, SS9, SS10, SS11, SS12, SS35, SS36, SS37, SS38, SS39, SS40, SS51, SS52, SS56

## What to build

The actual --audit-spec command reports exact source drift over declared scope.
Before any implementation, require complete independently accepted and landed preflight-ownership-closure delivery.
A draft or provider chunk is insufficient.
Reconcile graph pins and claims against that landed source through the existing enabling-plan procedure.
No FT318, FT317, or FT125 delivery is required.

Expose staleness.Collect(source ownership.Source, manifest ownership.Manifest, sourceFacts ownership.Snapshot) (AuditFacts, error) and staleness.Decide(facts AuditFacts) Audit.
Typed staleness.Error preserves Stage, Subject, Cause, and Unwrap.
AuditFacts carries baseline, tip, spec identity, path, symbol, test, coverage states, and diagnostics.
Audit carries three mechanical checks, drift, diagnostics, completeness, and ManualRequired.
The child imports neither its parent nor another workflow state owner.
Consume the complete landed metadata, citation, and authorization owners.

Gather supplies one exact committed pair and typed source bytes.
Enumerate both fence-prefix sets plus claims, metadata owners, coverage, active spec, and tickets.
Retain additions, deletions, renames, mode changes, dot paths, and both exact object IDs.
Union before and after members; use explicit absent cells and segment boundaries.
Missing before spec or unreadable objects refuse; checkout bytes cannot replace committed objects.
Spec prose has no semantic-digest exemption.

Add the selector only through the existing operation registry, requiring both source identities.
Reject every incompatible operation selector.
Render source, checks, spec_checks, drift, audit_diagnostics, and staleness with existing complete spills.
This checkpoint supplies spec-drift only.
Missing successor collectors remain diagnostic and manual_required true; the documented manual workflow stays unchanged.

Use types.go, collect.go, drift.go, decision.go, and drift_test.go in the new child.
Drive actual CommandWithVersion through committed fixtures and the existing seal, current-action, source-precedence, and movement seams.
Pin BENCH_KIT for system-tagged fixtures.
The complete prerequisite supplies parent headroom; add no root sibling or grant.
Remeasure actual composed Growth and crowding against that landed source.

This ticket belongs to SS-C1.
The explicit predecessor edge orders every overlapping write.
No unimplemented successor is required for this ticket's owned acceptance rows.

## Acceptance

- [ ] SS1: A modified fenced path appears in drift.
- [ ] SS2: A deleted fenced path appears with its before identity.
- [ ] SS3: A new fenced path appears with its after identity.
- [ ] SS4: A renamed fenced path retains both names.
- [ ] SS5: A file-mode change appears in drift.
- [ ] SS6: A fence prefix includes each tracked dot-path member.
- [ ] SS7: The audit prints exact resolved baseline and tip identities.
- [ ] SS8: An absent before spec refuses audit collection.
- [ ] SS9: Changed active spec bytes appear as drift.
- [ ] SS10: An unchanged exact member set prints a complete empty drift table.
- [ ] SS11: A segment-prefix collision stays outside the fenced scope.
- [ ] SS12: A claim source outside the write fence joins audit scope.
- [ ] SS35: Dirty source refuses before a staleness decision.
- [ ] SS36: Persistent snapshot movement publishes no skip decision.
- [ ] SS37: A source-tip mismatch retains its current refusal.
- [ ] SS38: A required executable-seal refusal remains binding.
- [ ] SS39: Audit-spec rejects combination with another selector.
- [ ] SS40: Audit-spec requires both explicit source identities.
- [ ] SS51: A special before/after source refuses without following it.
- [ ] SS52: An audit response retains all drift rows through existing spill.
- [ ] SS56: Existing required evidence bytes retain current-action binding.

## Verification

Run these future checks through their actual production entries.
Use a few sufficient multi-case witnesses; do not create a redundant test per row.
Retain the named omission mutation, its behavioral failure, restoration, and passing command result.
Compilation failure supplies no behavioral evidence.

- `bench test --package ./internal/preflight/...`
- `bench test --package ./cmd/bench`
- `bench test --check ticket-grammar`

Focused producer command: `bench test --package ./internal/preflight/... --run 'TestDriftObjects|TestDriftIdentity|TestAuditMovement'`.
Named mutation identity: `ss-01-omission`.
This obligation uses the named production seam and observable omission in this ticket, not a compile-failure substitute.

Mutation obligation: Omit a deleted or dot member, collapse mode or rename identity, or accept an absent before spec. The exact-object or source-refusal witness must fail.

The completion plan names this ticket's focused command and mutation identity.
Execute bench probe against the specified production seam after its code exists, then restore and rerun that exact focused command.
If the planned mutation cannot reach that producer, amend the enabling plan before verification; no substitute earns completion.

Read the canonical structure report before dispatch.
Remeasure Growth from the actual preceding tip and whole-tree crowding.
Keep cohesive extraction and necessary pin relocation in this same ticket.
