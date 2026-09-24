# 2. Redact and version every record line

Blocked by: none
Writes: internal/otelrecord/encode.go, internal/otelrecord/encode_test.go, internal/otelrecord/attributes.go, internal/otelrecord/reader.go, internal/otelrecord/reader_test.go, internal/otelrecord/processor_test.go, cmd/bench/process_env.go, cmd/bench/main.go, tests/canary/package-core-guard/unrouted-subcommand, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/systemtest/otel_gate_test.go
Covers: LE1, LE2, LE3, LE4, LE5, LE96, LE6, LE7, LE8, LE9, LE10

## What to build

Chunk: LE-A.

Make the encoder write its own resource block for each line. After the version is set, the block holds exactly `service.name` with the value `bench`, `service.version`, and `bench.record.schema` with the value `1`. The encoder ignores the SDK resource, so `OTEL_RESOURCE_ATTRIBUTES` and `OTEL_SERVICE_NAME` never reach a line.

Make the encoder drop each span attribute whose key is not in the declared set. Update the declared-set comment, because the set is now enforced.

Add a setter for the Bench version in the record package. The command layer calls it once at process start with its stamped version. `cmd/bench/main.go` is over its line budget, so the call goes into the process preparation in `cmd/bench/process_env.go`, and `main.go` does not grow. A process that sets no version writes no `service.version` key, so its block holds exactly two keys.

Make both readers treat a line whose schema attribute is not `1` as a malformed line. A line with no schema attribute is a legacy line and reads as before.

The version row lives in the system suite, whose files are system-tagged. Run it with `BENCH_KIT` set, through `bench test --check system`. The row compares with the second field of the `bench version` line, because that line also names the platform. The command registry and its conformance tests join the Writes line through the binding closure only, and the build expects no edit there.

## Acceptance

- [ ] An encoded line carries `bench.record.schema` with the value `1`.
- [ ] After the version is set to `9.9.9-test`, an encoded line carries `service.version` with that value.
- [ ] A `bench gate` run through the built binary writes lines whose `service.version` equals the second field of the `bench version` line, such as `0.2.0` in `bench 0.2.0 (linux/amd64)`.
- [ ] With `OTEL_RESOURCE_ATTRIBUTES=bench.leak=PROBEMARK` and `OTEL_SERVICE_NAME=SVCMARK` set, no line that `Begin` writes holds either marker.
- [ ] After the version is set, the resource block of a line holds exactly `service.name`, `service.version`, and `bench.record.schema`.
- [ ] With no version set, the resource block holds exactly `service.name` and `bench.record.schema`.
- [ ] A span with the undeclared attribute `bench.payload` encodes to a line with no such key, and a span with every declared attribute keeps each one.
- [ ] `ReadSelected` reports a schema `2` line as a malformed-line problem, and `ReadSpans` returns no span from it.
- [ ] `ReadSpans` returns the finished span of a legacy line with no schema attribute.
