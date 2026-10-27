# Share the plan-command parser with the test report

Blocked by: none
Writes: internal/testreport/selection_facts.go, internal/testreport/selection_facts_test.go, internal/preflight/check_floor.go
Covers: CF1, CF2, CF3

## What to build

The check floor of the build preflight parses each plan command with `planSelection`, and it compares package expressions with `packageCovers`. Move both rules into `internal/testreport`, beside the selection forms that they read, and export them. The check floor then calls the exported rules, and the build preflight keeps no copy.

The exported parser takes the root and one plan command. It answers the prepared request and whether the command is one simple `bench test` call that the test grammar accepts. The exported cover predicate answers whether one package expression equals an element or covers it at a segment boundary. Both keep their current behavior exactly. The carry decision of `03-decide-the-carry.md` is their second caller.

## Acceptance

- [ ] The check floor refuses a compound plan command, as `TestCheckFloorRejectsACompoundCommand` shows.
- [ ] The check floor keeps the segment boundary of a package pattern, as `TestCheckFloorKeepsTheSegmentBoundary` shows.
- [ ] Each existing check floor test in `internal/preflight` passes with no edit.
- [ ] `rg -n "func planSelection|func packageCovers" internal/preflight` finds no declaration.
