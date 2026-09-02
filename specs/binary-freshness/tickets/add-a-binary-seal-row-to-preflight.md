# Add a binary-seal row to preflight

Blocked by: none
Writes: internal/preflight/decision.go, internal/preflight/gather.go, internal/preflight/decision_test.go
Covers: BF26, BF27

## What to build

Verify the premise first: `Decide` in internal/preflight/decision.go runs the
check table, and no row reads a seal. Then gather the destination
`dist/bench` verdict in `Gather` through `freshness.Verify` on the root, and
add a `binary-seal` row after `kit-pin`. A mismatch reds with the
`RebuildAction` sentence in the detail cell. An absent `dist/bench` reports
not applicable. Keep the row inside the existing `CheckResult` shape.

## Acceptance

- [ ] A root with a mismatched seal reds `binary-seal` with the sentence.
- [ ] A root with no `dist/bench` reports `binary-seal` not applicable and the other rows unchanged.
- [ ] Self-probe: report not applicable on a mismatch, and report the new decision test red.
