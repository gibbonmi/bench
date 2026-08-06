# Refuse covers-policy violations at assign under an opted-in spec

Blocked by: export-coverage-map-row-ids.md, parse-ticket-covers-annotations.md
Ownership fence: `internal/specbuild`
Integration surfaces: map-ID export producer→export-coverage-map-row-ids.md; covers parse producer→parse-ticket-covers-annotations.md; refusal operation naming→existing internal/specbuild/refusal_operation_test.go exercised by the AS1/AS2/AS4 refusal strings naming assign
Contracts: the ordered map row IDs with the opt-in verdict crossing internal/coverage→`internal/specbuild`, asserted by AS1 and AS4 against the real coverage export over real spec files on disk

## What to build

Under a spec whose coverage map is opted into row IDs, `bench spec build
assign` refuses any ticket acceptance row whose covers annotation is missing,
malformed, or names no existing map ID, naming the offending row; tickets of a
non-opted-in spec assign exactly as today.

The policy sits in `Service.Assign` between `ParseTicket`/`requireCommittedTicket`
and `beginOperation`, following the `ContractsAnchored` refusal precedent
(single-sentence lowercase errors naming the subject and the word assign, per
refusal_operation_test.go). Opt-in and ID resolution come from the coverage
package's new exported entry over the run's real spec file — never a re-derived
map parse. An R-range line's expanded rows are unannotated by construction and
refuse under the same missing-covers rule, no special case. An opted-in map
that fails ID validation refuses assign (fail closed). `(covers local)` is
accepted. Fixtures follow `repo(t)`/`newCheckpointFixture` prior art: real git
repos, real spec and ticket files, driving `service.Assign`; the base fixture
spec has no coverage map, so opted-in cases write a 6-cell map into spec.md.

## Acceptance

- [ ] [AS1] under an opted-in spec, a row missing its covers annotation is refused naming the row, and an R-range line's expanded rows refuse under the same rule
- [ ] [AS2] a malformed or dangling covers annotation is refused naming the offending row
- [ ] [AS3] `(covers local)` and valid mappings assign, and a non-opted-in spec's tickets assign exactly as today
- [ ] [AS4] an opted-in spec whose map fails ID validation refuses assign

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AS1 | skip the covers check for rows expanded from a range | the range-refusal assign test | apply, run `go test ./internal/specbuild -run TestAssign`, expect the accepted-range failure |
| AS2 | check annotation presence but never resolve the target against the map IDs | the dangling-refusal assign test | apply, run `go test ./internal/specbuild -run TestAssign`, expect the accepted-dangling failure |
| AS3 | refuse `(covers local)` like a dangling ID | the local-acceptance assign test | apply, run `go test ./internal/specbuild -run TestAssign`, expect the refused-local failure |
| AS4 | treat a map with ID violations as not opted in | the invalid-map refusal test | apply, run `go test ./internal/specbuild -run TestAssign`, expect the failed-open failure |
