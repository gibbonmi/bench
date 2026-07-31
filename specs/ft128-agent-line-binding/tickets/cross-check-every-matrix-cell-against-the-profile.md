# Cross-check every matrix cell against the profile

Blocked by: Resolve the line through the harness matrix

## What to build

`checkLineBinding` grades the reviewer-owned binding as a matrix rather than a
sampled pair. It cross-checks every declared cell against the profile's
rendering, anchored to that profile's `Lines` section so a cell named in an
unrelated paragraph no longer satisfies the prose check, and it grades
completeness per *declared* harness so a partial column is a diagnostic while an
absent one is not.

Covers story 6.

## Acceptance

- [ ] Each of the six declared cells is cross-checked against the profile, one
      mutation per cell, with no cell satisfied by a sample of another.
- [ ] A matrix cell named only outside the profile's `Lines` section does not
      satisfy the prose check.
- [ ] A declared Claude harness missing its mid cell emits a diagnostic.
- [ ] A malformed matrix token and a missing `lines.env` each emit their own
      named diagnostic under the matrix keys.
