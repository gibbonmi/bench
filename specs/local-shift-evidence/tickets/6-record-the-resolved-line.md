# 6. Record the resolved line

Blocked by: 1-move-the-trace-handoff-into-the-record.md, 5-record-each-pass-under-the-shift.md
Writes: internal/shift/session.go, internal/shift/pass_test.go (new), cmd/bench/main.go, cmd/bench/guards.go, cmd/bench/otel_hook_seams_test.go, tests/canary/package-core-guard/unrouted-subcommand, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/otelrecord/attributes.go, internal/otelrecord/registry.go
Covers: LE45, LE46, LE47, LE48

## What to build

Chunk: LE-B2.

Add the trace handoff for the current pass span to the adapter environment. Use the record package's handoff function from ticket 1. Append the handoff after `BENCH_SHIFT=1`, so an inherited value never wins.

Make `bench resolve-model` record one `line.resolve` span. It joins the handed-off trace when the handoff is present, and it starts a new trace otherwise. The span carries `bench.line.harness`, `bench.line.tier` when `BENCH_MODEL` names a tier, and `bench.line.model` when the resolved model is a safe model token. Its outcome is green on exit 0 and red otherwise. Add the seam to the registry and declare the keys.

`cmd/bench/main.go` is over its line budget. Move the resolve-model block beside the other binding consumer, into `cmd/bench/guards.go`, so `main.go` shrinks. The command registry and its conformance tests join the Writes line through the binding closure only, and the build expects no edit there.

## Acceptance

- [ ] The adapter environment carries `BENCH_OTEL_ROOT` and a `BENCH_OTEL_TRACEPARENT` whose span id is the current pass span.
- [ ] With the handoff, a routed binding, and `BENCH_MODEL=mid`, `bench resolve-model --harness claude` writes a `line.resolve` span under the handoff span.
- [ ] That span carries harness `claude`, tier `mid`, and the bound model.
- [ ] `bench resolve-model` with no handoff writes a `line.resolve` span with no parent.
- [ ] A refused resolution in a routed repo with no `BENCH_MODEL` writes a `line.resolve` span with outcome `red` and no model key.
- [ ] The stdout of `bench resolve-model` is unchanged for each case.
