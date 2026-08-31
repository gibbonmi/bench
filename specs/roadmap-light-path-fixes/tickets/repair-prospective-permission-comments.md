# Repair prospective permission comments

Blocked by: add-negative-read-published-caller.md
Writes: internal/gate/prospectiveartifact/prospectiveartifact_test.go
Covers: LF24

## What to build

Replace retired PAR14 citations on the 0644 and 0400 record tests. Each edited
comment states the permission edge that its test proves.

## Acceptance

- [ ] The 0644 test comment states its live permission edge.
- [ ] The 0400 test comment states its live permission edge.
- [ ] Neither edited comment cites the retired PAR14 row.

