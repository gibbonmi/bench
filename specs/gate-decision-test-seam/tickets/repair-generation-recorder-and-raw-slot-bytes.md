# Repair generation recording and raw slot preservation

Blocked by: project-public-document-decisions-directly
Ownership fence: `internal/gate/check_slots_test.go`
Integration surfaces: decision-seam generation capture→`internal/gate/check_slots_test.go` + OP1; seeded check-slot store→`internal/gate/check_slots_test.go` + RS1
Contracts: every exhaustive changed-state generation capture crosses the matrix recorder→OP1 assertion; raw seeded check-slot paths and bytes cross fixture seed→every changed/restored state comparison
Closure: OP1/generation-recorder, OP1/second-capture-red, RS1/raw-slot-bytes

## What to build

Repair the integrated matrix so its generation-operation bound observes captures at the actual test seam rather than incrementing a counter beside one expected call, and so its restoration control compares the raw persisted slot-store bytes rather than normalized JSON. Keep the production oracle unchanged and retain the literal matrix, one seed, direct decision calls, existing operation markers, and representative controls.

## Acceptance

- [ ] [OP1] The matrix records every tree-generation capture used by each changed or restored state at the capture seam, asserts exactly one capture per state, and an added second capture makes the focused test red.
- [ ] [RS1] The matrix snapshots and compares the raw seeded check-slot store paths and bytes after every changed and restored state without JSON normalization.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OP1/generation-recorder | add a second generation capture inside one matrix state | the seam-owned operation recorder | temporarily add the second capture, run the focused matrix test, require an exact capture-count failure, then restore the test |
| OP1/second-capture-red | bypass the expected single-capture call site with another real capture | the recorder attached to the capture seam | temporarily invoke the capture seam twice for one state, run the focused matrix test, require red, then restore the test |
| RS1/raw-slot-bytes | rewrite one seeded slot file with semantically equivalent differently formatted JSON | the raw slot-store comparison | temporarily rewrite one seeded slot file after capture, run the focused matrix test, require a raw-byte mismatch, then restore the test |
