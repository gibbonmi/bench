# Project commitment state through readers and guidance

Blocked by: 08-verify-milestone-outcomes.md
Writes: internal/commitment (new), internal/intent, internal/roadmap, internal/status, internal/dashboard, cmd/bench, internal/usage, internal/conformance, internal/anchors, tests/canary/workflow-guidance-anchors, .bench/BENCH.md, .bench/BENCH-reference.md, .agents/commands/bench.md, .agents/commands/bench-what-next.md, .agents/commands/bench-drain.md, .agents/commands/bench-write-spec.md, .agents/commands/bench-implement-spec.md, .agents/commands/bench-final-check.md, .agents/commands/bench-setup-repo.md, .agents/commands/bench-debug.md, README.md, DATA_HANDLING.md, CHANGELOG.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/data-handling-derivation/undocumented-passlist-var, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/docs-currency-token-diet/dogfood-referent-shipped, tests/canary/docs-currency-token-diet/introduces-undeclared-command, tests/canary/docs-currency-token-diet/missing-cli-inventory, tests/canary/docs-currency-token-diet/readme-command-first, tests/canary/docs-currency-token-diet/stale-cli-doc-reference, tests/canary/docs-currency-token-diet/stale-command-reference, tests/canary/load-validity-metadata/readme-shared-rule-drift, tests/canary/load-validity-metadata/shared-rule-drift, tests/canary/row-next-grammar/token-table-lacks-kit-edit, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/unindexed-skill
Covers: DC47, DC54, DC64, DC65

## What to build

Render one shared commitment projection in roadmap, status, and dashboard. The eligible outcome, blocker, missing approval, and adoption remedy must agree.
Ticket 08 supplies the completed command family. Do not introduce a second eligibility calculation in a reader.
Keep diagnostic detail bounded through existing output mechanisms.

Put the canonical commitment rule in BENCH.md. Phase guidance invokes its commands and preserves one source per rule.
Remove each unconditional implement-now or sequence-replacement grant. Give each removed granting sentence an independent forbid expectation and demonstrate its red.
Explain purpose-based defect and refactor priority without giving labels automatic admission.

Update help, command-route assertions, anchor registries, data inventory, and adoption guidance in the same checkpoint.
The route assertion covers every required admission consumer. A consumer omission must fail it, independently of the installed-owner system proof.
Document proposal references, runtime claims, blockers, and evidence retention without collecting transcripts.

Read only the projection consumers, phase clauses listed in the spec sweep, registry rows, and their pinned fixtures. Search each changed output token across tests, help, references, and anchors. Extract status projection work into a focused sibling before growing the existing oversized file.

## Acceptance

- [ ] Roadmap, status, and dashboard select B when A is blocked and B is independent (DC47).
- [ ] The canonical rule and phase routes preserve defect/refactor priority and explicit displacement (DC54).
- [ ] Removing one admission consumer fails the route test; record its red (DC64).
- [ ] The data inventory states record contents and local retention (DC65).
- [ ] Each retired granting sentence has a separate omission-sensitive prohibition. Existing anchors and generated help agree with the new routes.

## Checkpoint verification

Run `bench test --package ./internal/roadmap`, `bench test --package ./internal/status`, `bench test --package ./internal/dashboard`, `bench test --package ./cmd/bench`, `bench test --package ./internal/anchors`, and `bench test --package ./internal/conformance`. Run prose checks on each changed document.
