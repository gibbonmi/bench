# Refuse a Contracts crossing that anchors no fenced path

Blocked by: teach-fence-derivation-from-the-contracts-crossing.md
Ownership fence: `internal/specbuild/assign.go`, `internal/specbuild/parse_ticket_test.go`, `internal/specbuild/checkpoint_fixture_test.go`
Contracts: the parsed crossing's operands cross the ticket file→`internal/specbuild/assign.go`, asserted by CA3 against the real assign path
Assumptions: `ParseTicket` does not read `Contracts:` today, so the field must be parsed before it can be graded; the parse stays permissive because the conformance example-agreement check grades deliberately-mutated examples through it and a parse refusal would shadow its fence-mismatch diagnostic; the refusal fires only when a `Contracts:` line is present and unanchored, because 108 of 235 ticket blobs in history carry a fence with no `Contracts:` line and live specs still hold such tickets; the rule was validated against the only 10 tickets that carry the field, where it refuses 3 and passes 7 with no false refusal; `none crosses` stays exempt; claims re-derived from the tree at pickup

## What to build

`bench spec build assign` refuses a ticket whose `Contracts:` line names no
backticked path lying inside that ticket's own ownership fence, so a repair
scoped to a finding's cited lines is rejected at lease time instead of
discovered a review round later. A crossing may still name an unfenceable
surface on its other side (`bash`, every audited package); only the written side
must anchor. The judgment is assignment policy rather than parse validity: the
shared parser stays permissive so its other consumers keep grading the grammar
they own.

## Acceptance

- [ ] [CA1] `ParseTicket` reads the `Contracts:` line into the parsed ticket and refuses nothing new.
- [ ] [CA2] `ContractsAnchored` answers false only when every operand falls outside the fence, and stays true for an absent line, `none crosses`, and a crossing whose far side names no path.
- [ ] [CA3] `Assign` refuses an unanchored ticket and names the unanchored crossing in its error.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CA1 | drop the `Contracts:` prefix branch from the parse loop | the parse test | remove the branch, run `go test -count=1 ./internal/specbuild -run TestParseTicketReadsTheContractsLine`, expect the empty-contracts failure |
| CA2 | require *every* operand to anchor | the anchoring table test | tighten the check, run `go test -count=1 ./internal/specbuild -run TestTicketContractsAnchored`, expect the unfenceable-far-side failure |
| CA3 | delete the refusal from the assign path | the assign refusal test | remove the guard, run `go test -count=1 ./internal/specbuild -run TestAssignRefusesAContract`, expect the leased-unanchored-ticket failure |
