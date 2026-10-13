# Document the mutant set and close the probe-cost issue

Blocked by: 5-import-one-mutant-row.md
Writes: .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .bench/BENCH-reference.md, CHANGELOG.md, CONTEXT.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, docs/known-issues.md, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, internal/anchors/registry_check_floor_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, projects/benchkit.md, roadmap/FT168.md, tests/canary/claude-agent-definitions/, tests/canary/docs-currency-token-diet/, tests/canary/guidance-prose-budgets/, tests/canary/line-routing/, tests/canary/skill-description-budgets/, tests/canary/skills-index-command-adapters/, tests/canary/workflow-guidance-anchors/
Covers: PM61, PM62, PM63, PM64, PM65, PM66, PM68, PM73, PM75

## What to build

The guidance names the mutant set, the repeated `--package`, and the `--mutant` import as current behavior.

- `CONTEXT.md` gains the **mutant set** entry with the Avoid list "probe batch", "mutant batch", and "multi-probe". The **combined mutant** entry stays byte-identical.
- `.bench/BENCH-reference.md` states the `--mutants` form, the repeated `--package`, the rows, the exits, the build rule, and the stop rules. The sentence "The focused run is one package or one check" changes.
- `.agents/skills/bench-craft-delegate/SKILL.md` names the mutant set beside the mutation forms.
- `references/delegation-discipline.md` adds one sentence after its first `restored` bullet: a coordinator confirms the `restored` cell of each mutant row.
- `.agents/commands/bench-implement-spec.md` adds one sentence beside the witness-cap sentences: an author probes several witnesses of one subject in one `bench probe --mutants` call.
- `docs/known-issues.md` loses the probe-cost row, and the `Require` anchor on its workaround sentence goes with its test row.
- `CHANGELOG.md` gains one Unreleased entry for the mutant set, the repeated `--package`, and `--mutant`.
- `roadmap/FT168.md` drops the delivered faces and the stale staged-spec reference, and it keeps the canary and failing-name faces.

The two anchored `restored` sentences of `delegation-discipline.md` stay byte-identical.
The anchored witness-cap sentences of `bench-implement-spec.md` stay byte-identical.
The `Forbid` anchor on "Cap the witnesses in a charge" stays.

`SKILL.md` and `bench-implement-spec.md` sit at their line budgets in `projects/benchkit.md`.
Fit each addition into an existing line, or raise that file's budget row in `projects/benchkit.md` with a one-clause reason.
Each anchor registry file and canary directory that pins a touched guidance path joins this `Writes:` line from the build preflight's closure list.

## Acceptance

- [ ] `CONTEXT.md` defines **mutant set**, and the **combined mutant** entry is byte-identical.
- [ ] The reference guide's probe notes state the `--mutants` form, the rows, the exits, the build rule, and the stop rules.
- [ ] The delegate skill's probe sentence names the mutant set.
- [ ] The probes guidance holds the per-row `restored` sentence, and the anchors check and the `delegate-probe-mutated-bytes` canary stay green.
- [ ] `/bench-implement-spec` holds the mutant-set witness sentence, and `TestWriteSpecChunkedSlicing` stays green.
- [ ] The probe-cost row and its `Require` anchor are gone, and the live anchors check stays green.
- [ ] `CHANGELOG.md` holds the Unreleased entry.
- [ ] `roadmap/FT168.md` names only its open faces.
