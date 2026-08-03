# Tighten the agreement check's blind spots

Blocked by: land-the-example-agreement-check.md
Ownership fence: `internal/conformance/example_agreement_test.go`
Assumptions: the per-ID mutations assertion matches an ID anywhere after the Red-mutations heading rather than in the table's criterion cell; the extractor appends a newline before materializing, so the EOF case never exercises un-normalized bytes; the passlist token regex anchored to the table's first cell is the shape precedent. Re-derive from the tree at pickup.

## What to build

FT164 repair round, falsification findings: the example-agreement check's
per-ID mutations assertion accepts an acceptance ID mentioned anywhere in the
section — a deleted table row passes when another row's prose names the ID —
so the assertion anchors to the criterion cell of a table row, mirroring the
passlist first-cell regex. And the EOF case normalizes the boundary away: the
materialization appends a newline unconditionally, so the check now
materializes the extracted bytes exactly as found, appends nothing, and the
EOF bite fixture proves an unterminated final line parses identically through
the real path.

## Acceptance

- [ ] [TA1] the per-ID mutations assertion matches only a table row whose criterion cell is the ID, and a deleted row reds even when the ID appears in another row's text.
- [ ] [TA2] the materialized example carries the extracted bytes without an appended newline, and the EOF-without-newline bite proves identical parses through the un-normalized path.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TA1 | delete RC1's table row while RC2's operation text mentions RC1 | `TestExampleAgreementPerIDRowsAnchorToCriterionCell` | build the temp doc with the deleted row and the stray mention, run the check, expect the no-mutation-for-RC1 diagnostic |
| TA2 | re-add the unconditional newline append before materializing | `TestExampleAgreementEOFWithoutNewline` | strip the final newline in the fixture, run the check through the real materialization, expect the parses compared without normalization |
