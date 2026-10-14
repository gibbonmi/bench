# Report a dirty record in the current binding

Blocked by: none
Writes: internal/preflight/preparation.go, internal/preflight/charge_pack.go, internal/preflight/charge.go, internal/preflight/review.go, internal/preflight/proposal.go, internal/preflight/command.go, internal/preflight/charge_evidence_test.go, internal/preflight/charge_legacy_pack_test.go, internal/preflight/proposal_command_test.go, internal/preflight/review_charge_test.go, internal/preflight/testdata/legacy-charge/, internal/preflight/preflighttest/fixture.go, internal/preflight/evidencecmd/evidence_modes_test.go, internal/preflight/evidencecmd/evidence_grammar_test.go, internal/preflight/evidencecmd/evidence_record_order_test.go, internal/preflight/evidencecmd/evidence_budget_test.go, internal/preflight/evidencecmd/operations.go, internal/chargeevidence/read.go, internal/chargeevidence/schema.go, internal/chargeevidence/reference.go, internal/chargeevidence/format_test.go, .agents/skills/bench-craft-delegate/references/charge-evidence-format.md, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, cmd/bench/commitment_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RT40, RT41, RT42, RT43, RT44, RT45, RT46

## What to build

`bench preflight evidence <id> --check-current` binds beside a dirty review record of its spec and reports it.
The current binding already passes the `check` action to `preparationCheckoutRefusal`.
For the check action only, the refusal skips one dirty path: `reviewrecord.RecordPath` of the selected spec.
The `current` response block gains a `dirty_record` flag field after `current`, true when that path is dirty.
The generated evidence reference lists the new field.

Any other dirty path still refuses the binding, and a dirty record of another spec refuses too.
The build charge, the review charge, and the proposal keep the strict refusal.

`internal/chargeevidence/read.go` has 11 lines of headroom, and `internal/preflight/review_charge_test.go` has 3.
The new binding cases go in `evidence_modes_test.go` beside CE123, and the RT45 charge case goes in `charge_evidence_test.go` or `charge_legacy_pack_test.go`.
The preflight directories are crowded, so this ticket adds no new test file there.
The fixture fence already lists `reviews/example.md`, so a test can dirty that path with no fixture change.

## Acceptance

- [ ] `--check-current` with a dirty record of its spec binds and reports `dirty_record` true.
- [ ] `--check-current` on a clean checkout reports `dirty_record` false.
- [ ] `--check-current` with a dirty untracked Go file still refuses with the dirty checkout refusal.
- [ ] `--check-current` with a dirty record and a dirty Go file refuses.
- [ ] `--check-current` with a dirty record of another spec refuses.
- [ ] The build charge with a dirty record of its spec still refuses with the dirty checkout refusal.
- [ ] The generated evidence reference lists `dirty_record` in the `current` block, and `TestEvidenceSchemaRegistry` passes.
