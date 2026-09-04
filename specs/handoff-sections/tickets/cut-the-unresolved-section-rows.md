# Cut the unresolved section rows

Blocked by: date-each-section-in-bench-status.md
Writes: internal/status/handoff.go, internal/status/handoff_test.go
Covers: HS21

## What to build

Repair for review finding Sp1. Verify the premise first. `bench status`
emits an unresolved-section row when a section's key names no active
assignment. It emits one too when the tip is empty or does not resolve.
Then remove those rows and their reasons.

A section the join cannot resolve reports nothing, in the default view and
under `--all`. Keep the behind rows and their ordering.

## Acceptance

- [ ] A document with one behind section and one section whose key names no assignment yields exactly one row, and it names the behind section.
- [ ] `TestAppendHandoffSectionNamesTheBehindSection` and `TestAppendHandoffSectionRewriteLeavesTheDistance` pass.
- [ ] Self-probe: keep one unresolved row, and report the one-row test red.
