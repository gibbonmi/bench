# mixed tag row ids (canary fixture)

This fixture uses the opt-in header, which carries a `row` cell, and proves the
mixed-tag diagnostic reaches the gate: row IDs that carry two tags fail a real
gate run.

## User stories
1. As a reviewer, I want a map whose row IDs carry two tags to fail the gate.

### Acceptance coverage map
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| MT1 | 1 | a behavior | a seam | why it catches |
| XT2 | 1 | another behavior | another seam | why it catches |
