# 1. Move the trace handoff into the record package

Blocked by: none
Writes: internal/otelrecord/provider.go, internal/otelrecord/processor_test.go, internal/gate/telemetry.go, internal/gate/runner.go, internal/gate/gate.go, internal/gate/otel_env_test.go
Covers: LE11, LE12

## What to build

Chunk: LE-A.

Move the two handoff variables, `BENCH_OTEL_ROOT` and `BENCH_OTEL_TRACEPARENT`, from the gate package to the record package. Move their inject and extract steps with them. The record package then offers one function that adds the handoff for a context to a child environment. It also offers one function that attaches a child process to the handed-off trace and record.

The gate calls these functions to compose its phase child and to attach the phase process. The gate still strips each inherited handoff value before it composes that child. This ticket changes no record line. `internal/gate/runner.go` is over its line budget, so the move must leave that file shorter.

Add a source test to `internal/gate/otel_env_test.go` that reads each non-test Go file of the gate package and fails on the literal `BENCH_OTEL_`. Build the searched literal from two parts, so the test file does not match itself.

Ticket 6 consumes this seam for the shift adapter and `bench resolve-model`.

## Acceptance

- [ ] A child span started from the handoff of a parent span carries the parent's trace id and the parent's span id as its parent.
- [ ] `TestGateEnvStripsTheRecordVariables` still removes both inherited values from the phase environment.
- [ ] The gate phase spans still sit under the root gate span in the existing gate span tests.
- [ ] No non-test Go file in `internal/gate` holds the literal `BENCH_OTEL_`, and the source test fails when a copy of either name returns.
