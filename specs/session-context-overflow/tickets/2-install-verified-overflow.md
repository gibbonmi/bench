# 2. Install verified overflow projection

Blocked by: 1-prove-runtime-paths.md
Writes: internal/harnessoverflow (new), .bench/hooks/result-overflow.sh (new), .codex/hooks.json, .claude/settings.json, internal/harnesses/harnesses.go, internal/harnesses/harnesses_test.go, internal/conformance/harness_record_test.go, cmd/bench/main.go, cmd/bench/otel_hook_seams_test.go, internal/otelrecord/registry.go, internal/systemtest/harness_overflow_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: OV2, OV3, OV4, OV5, OV6, OV7, OV8, OV9, OV10, OV11, OV12, OV13, OV14, OV15, OV16, OV17, OV18, OV19

## What to build

Entry requires runtime proof, approved numeric budgets, and a completed spec-authoring checkpoint.
Stop before product writes when any prerequisite is absent.
Implement the verified adapter, artifact owner, bounded projection, configuration, and end-to-end tests as one vertical slice.
Do not install an unsupported path or select a budget during the build.
The final fence must include the actual adoption readers before this ticket receives approval.

## Acceptance

- [ ] A complete artifact exists before any oversized replacement reaches the model.
- [ ] The result retains producer status, approved diagnostics, true bytes, and an exact retrieval route.
- [ ] Unsupported inputs and preservation failures retain the original tool result with a limitation.
- [ ] The producer sentinel proves exactly one invocation across duplicate callbacks and failed spills.
- [ ] Runtime drift cannot silently retain verified status.
- [ ] Installed-path tests preserve existing guards and private artifact ownership.
