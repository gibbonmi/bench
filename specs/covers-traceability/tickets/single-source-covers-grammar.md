# Single-source the covers ID grammar and pin the reviewed edge cases

Blocked by: export-coverage-map-row-ids.md, parse-ticket-covers-annotations.md, refuse-covers-violations-at-assign.md, refuse-uncovered-map-rows-at-promote.md
Ownership fence: `internal/coverage`, `internal/specbuild`
Integration surfaces: exported grammar consumer→internal/specbuild (this fence); example-agreement taught example→existing internal/conformance/example_agreement_test.go path, unchanged (single valid annotations only)
Contracts: the row-ID grammar crossing `internal/coverage`→`internal/specbuild`, asserted by RG1 against the real exported pattern

## What to build

The review found two derivations of the row-ID grammar that disagree
(`rowIDRe = ^[A-Z]+[0-9]+$` in coverage, `coversValue = ^(local|[A-Z][A-Z0-9]*)$`
in specbuild — `AB` and `A1B2` pass one and fail the other). Export the grammar
once from `internal/coverage` and have the covers operand compose it, so an
operand failing the one grammar parses as unannotated and refuses at assign
under an opted-in spec. A row carrying more than one covers annotation is the
compound-row failure one level down: it parses as unannotated (malformed
class) rather than first-wins, and so refuses at assign the same way — the
parser stays permissive, policy still owns every refusal. Pin the reviewed
edge cases: backtick and bracket payloads in the malformed parse table, a
zero-assignment promote under an opted-in map refusing with every map ID
named, and the existing invalid-map refusal branch in `requireCoversTotality`.

Also correct two comments to state what the code does: the `rowIDRe` comment
claims ParseTicket's row IDs share the grammar (they are unconstrained
`[^]]+`), and the totality call-site comment says the check runs before the
gate owner is *resolved* (it runs before the owner *executes*).

## Acceptance

- [ ] [RG1] the covers operand composes coverage's one exported row-ID grammar, so `(covers AB)` and `(covers A1B2)` parse as unannotated and refuse at assign under an opted-in spec
- [ ] [RG2] a row carrying more than one covers annotation parses as unannotated and refuses at assign under an opted-in spec
- [ ] [RG3] backticked and bracketed covers payloads parse as unannotated
- [ ] [RG4] promote under an opted-in map with zero assignments refuses naming every map ID
- [ ] [RG5] promote under an opted-in map that fails ID validation refuses

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RG1 | restore a locally-defined `^(local|[A-Z][A-Z0-9]*)$` operand | the grammar-agreement parse test | apply, run `go test ./internal/specbuild -run TestParseTicket`, expect the AB-annotated failure |
| RG2 | keep first-wins FindStringSubmatch over the row's trailing text | the multiple-annotation parse test | apply, run `go test ./internal/specbuild -run TestParseTicket`, expect the first-wins failure |
| RG3 | strip backticks and brackets from the operand before matching | the hostile-payload parse rows | apply, run `go test ./internal/specbuild -run TestParseTicket`, expect the annotated-payload failure |
| RG4 | return nil from totality when the run has no assignments | the zero-assignment promote test | apply, run `go test ./internal/specbuild -run TestPromote`, expect the accepted-empty-run failure |
| RG5 | treat a map with violations as not opted in inside totality | the invalid-map promote test | apply, run `go test ./internal/specbuild -run TestPromote`, expect the failed-open failure |
