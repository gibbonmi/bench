# Add the background width input

Blocked by: 01-measure-the-background-census.md
Writes: internal/gate/phases.go, internal/gate/gate_go.go, internal/gate/gate.go, internal/gate/width.go (new), internal/gate/width_test.go (new), internal/gate/branch_native_phases_test.go, CHANGELOG.md
Covers: BG19, BG20, BG21, BG22

## What to build

The gate gives its children a closed environment, so an ambient `GOFLAGS` never reaches a gate's `go test` (`internal/gate/subject.go:69-71`, `internal/gate/subject.go:101-111`).
No Bench setting selects a gate width today.

Add one Bench-owned width input. When it is set, each Go test phase of the gate runs with `-p=2` and `-parallel=2`.
The Go test phases build their argv through `BaseTestArgv` and the race and system phase builders (`internal/gate/phases.go:63`, `internal/gate/gate_go.go:126`, `internal/gate/gate_go.go:212`, `internal/gate/gate_go.go:232`).
The input is not a member of the grade key.
`gateEnv` removes it from each phase child, as it removes the delta schedule entry (`internal/gate/gate.go:158-173`).
Follow the precedent of `BENCH_GATE_DELTA` (`internal/gate/delta_checkpoint.go:31`, `internal/gate/delta_checkpoint.go:84`).

A run without the input keeps its argv byte for byte.

## Acceptance

- [ ] With the width input set, each Go test phase argv holds `-p=2` and `-parallel=2`.
- [ ] The width input leaves the grade key of one subject unchanged.
- [ ] Without the width input, no Go test phase argv holds `-p` or `-parallel`.
- [ ] A phase child does not inherit the width input.
