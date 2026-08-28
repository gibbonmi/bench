# Repair the ungraded gate anchors

Blocked by: none
Writes: internal/conformance/workflow_checks_test.go, internal/conformance/release_probe_fixture_test.go, internal/releaseevidence/release_plan_test.go

## What to build

The gate grades the matrix producers, not only their consumers. Review found
that changing one command in the `preflight` job restarts both macOS runners
with the gate still green.

`checkNativeRuntimeWorkflow` asserts that the `preflight` job declares the
`proven` output from its proof-matrix step, and that the step's `run` line names
`proof-matrix-json`. It makes the same assertion for the shipped `matrix` output
and `matrix-json`. Each assertion carries a mutation that turns it red.

The release evidence probe removes one proven proof file and requires the
`native target proof is incomplete` diagnostic. That diagnostic occurs once in
the tree today, at its definition, so the count expression has no red.

The plan decode test drops its four-row expectation table. The decode assertion
alone serves the test's stated purpose, and the table is a second copy of the
plan.

## Acceptance

- [ ] The check names a diagnostic when the proof-matrix step runs `matrix-json` (row B4).
- [ ] The check names a diagnostic when the `proven` output is absent (row B4).
- [ ] The check names a diagnostic when the shipped matrix step runs `proof-matrix-json` (row B5).
- [ ] Every added assertion carries a mutation that turns it red (row C2).
- [ ] The probe fails with `native target proof is incomplete` when one proven proof is removed (row B10).
- [ ] The plan decode test carries no enumerated target table.
