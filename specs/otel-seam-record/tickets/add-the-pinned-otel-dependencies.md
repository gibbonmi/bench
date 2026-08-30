# 1. Add the pinned OpenTelemetry dependencies

Blocked by: none
Line: opus / medium
Rows: OT13
Writes: go.mod, go.sum, internal/conformance/otel_footprint_test.go (new)

## What to build

The tree gains three dependencies: `go.opentelemetry.io/otel`,
`go.opentelemetry.io/otel/trace`, and `go.opentelemetry.io/otel/sdk`. Each one
is pinned at v1.46.0 under Apache-2.0. The transitive set is logr, stdr,
xxhash, uuid, x/sys, and the passive `auto/sdk` shim.

The module `sdk/metric` stays out, because this spec ships spans alone. No OTLP
exporter enters the tree, because every official exporter drags
`google.golang.org/protobuf` in.

A new footprint test in `internal/conformance` reads `go.mod`. The test reds
when `go.mod` names a protobuf module or a gRPC module. The dependency standard
owns this shape: an official-org source, an MIT-compatible license, and a
build-time-only footprint.

This ticket adds no span and no package. It lands the dependency ahead of the
encoder, so the encoder ticket reads a green tree.

## Acceptance

- [ ] OT13: the footprint test reds when `go.mod` names a protobuf module or a gRPC module.
- [ ] the footprint test passes over the live `go.mod` with the three pinned modules present.
- [ ] `go build ./...` and the gate `test` phase stay green.
