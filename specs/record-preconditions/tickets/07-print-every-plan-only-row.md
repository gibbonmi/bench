# Print every check row for a plan-only preflight

Blocked by: none
Writes: internal/preflight/command.go, internal/preflight/evidencecmd/operations.go, internal/preflight/verdict_summary_test.go, internal/preflight/map_sources_test.go, cmd/bench/help_inventory_test.go, CHANGELOG.md
Covers: RP42, RP43, RP44

## What to build

Review chunk: RP-C5.

Make `renderChecks` (`internal/preflight/command.go:220-250`) take the operation kind. The plan-only kind prints every check row in check order, and every other kind keeps the red-only table of `redRows`. The count line and the exit rule stay the same. Change the plan-only help description (`internal/preflight/evidencecmd/operations.go:107-108`) to state that it prints every check row, and move the pinned row at `cmd/bench/help_inventory_test.go:112`. Update the comment at `command.go:216-217` and the comment at `map_sources_test.go:54-55` to the new state.

A full table holds more lines than the public response bound, so the CLI response spills, and the spill file holds the complete table. Do not change the response bound. Add one `CHANGELOG.md` entry under `## [Unreleased]` for RP-C5.

## Acceptance

- [ ] A green plan-only verdict prints a `checks` table with one row per check, green and not-applicable rows included (RP42).
- [ ] A green ordinary build verdict still prints only the count line (RP43).
- [ ] A plan-only verdict still exits 0 when green and 1 when a row is red (RP44).

## Verification

Run `bench test --package ./internal/preflight/...` and `bench test --package ./cmd/bench` whole. Observe RP42 red before the production edit, and record the route of the red. Self-probe: keep the red filter for the plan-only kind, and report that RP42 goes red.
