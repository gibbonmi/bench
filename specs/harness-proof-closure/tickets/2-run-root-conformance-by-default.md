# Run root conformance by default

Blocked by: none
Writes: internal/conformance/gate_entry_test.go, internal/conformance/harness_test.go, internal/conformance/registry/registry.go
Covers: HP4, HP5

## What to build

Make the live-tree conformance entry use its harness repository root when
`BENCH_CONFORMANCE_ROOT` is unset. Keep an explicit environment value as the
higher-precedence graded root. Remove the environment-class skip for the unset
case and update the public registry contract to describe the new default.

## Acceptance

- [ ] A package run with no conformance-root environment evaluates the current repository root.
- [ ] An explicit conformance-root environment evaluates that root instead of the current repository.
- [ ] The ordinary gate's existing explicit root and tier route keeps its current behavior.
