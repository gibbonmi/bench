# Correct the stale diagnostics and pin the prose

Blocked by: none
Writes: internal/structure, .bench/structure-accept, .agents/commands/bench-review-implementation.md, .bench/BENCH-reference.md, internal/anchors/registry_data.go

## What to build

Four accept rows in `.bench/structure-accept` name paths that are no longer
scanned source files, so `bench structure` opens with bookkeeping instead of
findings. Those rows go, and every surviving grant keeps its exact reason. The
review-phase and reference passages that describe the gate's conformance shape
are corrected to match the phase table the gate actually runs, and an anchor
pins each corrected passage so it reds when the wording drifts back.

## Acceptance

- [ ] `bench structure` prints no stale accept row.
- [ ] every surviving accept grant resolves to a scanned source file and keeps its reason text.
- [ ] the review-phase and reference passages describe the conformance shape the phase table runs.
- [ ] removing a corrected passage turns the anchors check red.
