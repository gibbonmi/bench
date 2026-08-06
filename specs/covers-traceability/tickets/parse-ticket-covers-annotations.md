# Parse per-row covers annotations in tickets

Blocked by: none
Ownership fence: `internal/specbuild/assign.go`, `internal/specbuild/parse_ticket_test.go`
Integration surfaces: Ticket covers consumer→refuse-covers-violations-at-assign.md; taught-example grader→teach-covers-schema-in-skills.md; legacy grammar guarantee→existing internal/specbuild/legacy_ticket_test.go exercised by PT2
Contracts: the per-row covers mapping (aligned with Rows; empty for unannotated, `local`, or a map-row ID) crossing the `internal/specbuild/assign.go` parse→assign and promote policy, asserted by AS1–AS3 on refuse-covers-violations-at-assign.md against this real parser

## What to build

`ParseTicket` captures an optional `(covers <ID>)` or `(covers local)` written
after an acceptance row's ID bracket and exposes the mapping per row on
`Ticket`, staying permissive so every ticket that parses today parses
byte-for-byte identically.

The `ticketRow` regex captures only the bracket ID and discards the rest of
the line, so the annotation needs a second capture or follow-on scan. The
mapping aligns index-for-index with `Rows` (an ordered slice; empty string =
unannotated). A malformed annotation parses as unannotated — never a parse
error, per the ContractsAnchored precedent: policy refuses, the parser does
not. R-range expansion (`expandRows`) is untouched, so expanded rows carry no
annotation. `Digest` already hashes the full file bytes; add the regression
test that a covers-only edit moves it. Tests live in the external
`specbuild_test` package through the existing `ticketFixture` helper.

## Acceptance

- [ ] [PT1] ParseTicket captures `(covers <ID>)` and `(covers local)` after a row's ID bracket, exposed per row on Ticket
- [ ] [PT2] unannotated rows and malformed annotations parse as unannotated, and legacy tickets and R-range expansions parse exactly as today
- [ ] [PT3] editing only a covers annotation changes Ticket.Digest

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PT1 | discard the covers capture so every row parses as unannotated | the covers parse test | apply, run `go test ./internal/specbuild -run TestParseTicket`, expect the captured-mapping failure |
| PT2 | return a parse error on a malformed annotation | the malformed-as-unannotated test | apply, run `go test ./internal/specbuild -run TestParseTicket`, expect the permissiveness failure |
| PT3 | compute Digest over the file with annotations stripped | the digest sensitivity test | apply, run `go test ./internal/specbuild -run TestParseTicket`, expect the unchanged-digest failure |
