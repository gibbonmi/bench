# Export coverage-map row IDs behind the 6-cell opt-in header

Blocked by: none
Ownership fence: `internal/coverage`
Integration surfaces: exported row-ID parse consumer→refuse-covers-violations-at-assign.md; staged-spec `--check` sweep→existing internal/conformance/docs_workflow_checks_test.go checkCoverageMaps exercised by CV3; `bench coverage` TOON schema→internal/coverage (Rows stays the 3-column table TestCommand pins)
Contracts: the ordered map row IDs with the opt-in verdict and Check violations crossing `internal/coverage`→internal/specbuild, asserted by AS1 and AS4 on refuse-covers-violations-at-assign.md against this real export

## What to build

A spec author opts a coverage map into row IDs by leading the canonical header
with a `row` column, `bench coverage --check` validates the IDs, and an exported
entry point hands consumers outside the package the parsed IDs.

The canonical 6-cell header is `row|story|behavior|seam|red signal|why it
catches the failure`; the 5-cell header stays valid and grades exactly as
today. ID grammar is the ticket-ID shape: an uppercase tag plus a number
(`^[A-Z]+[0-9]+$`), unique within the map, no empty ID cells, and a map mixing
ID and non-ID rows is a violation. `Check`, `State`, and `Rows` take the
unexported `parsed` type, so add one exported function (parse a spec path,
return the opt-in verdict, the ordered row IDs, and the Check violations) as
the lifecycle's sole entry; do not export `parsed` itself. `dataRow` is
hard-capped at 5 cells (`all [5]string`), so the widening touches the struct.
Violation strings follow the existing lowercase substring-matched style.

## Acceptance

- [ ] [CV1] a 6-cell map with a leading `row` column parses, and the exported entry point returns its ordered row IDs and opt-in verdict
- [ ] [CV2] duplicate IDs, empty ID cells, bad-grammar IDs, and an ID/non-ID mix are `--check` violations
- [ ] [CV3] a 5-cell map checks exactly as today and reports not opted in

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CV1 | restore the whole-row 5-cell header equality so a 6-cell header is refused | the exported-parse test | apply, run `go test ./internal/coverage`, expect the missing-header failure on the 6-cell fixture |
| CV2 | drop the ID-grammar validation branch | the ID-violation Check tests | apply, run `go test ./internal/coverage`, expect the bad-grammar and duplicate cases to miss their violations |
| CV3 | require the `row` column unconditionally | the existing 5-cell Check tests | apply, run `go test ./internal/coverage`, expect the legacy-header cases to fail |
