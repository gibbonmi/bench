# Extend the bounded operation journal

Blocked by: allow-verified-clean-provisional-release.md
Ownership fence: `internal/specbuild/integrate.go`, `internal/specbuild/operation_journal_test.go`
Integration surfaces: operation admission and shared finite bound→`internal/specbuild/integrate.go`; decoded operation-count validation→existing `internal/specbuild/state.go` plus OJ2; 65th-entry and full-journal proofs→`internal/specbuild/operation_journal_test.go` plus OJ1-OJ2
Contracts: operation records and their finite count cross `internal/specbuild/integrate.go` admission→existing `internal/specbuild/state.go` decode validation, asserted by OJ1-OJ2 against durable save and reload in `internal/specbuild/operation_journal_test.go`

## What to build

Raise the durable operation-journal bound from 64 to 128 so an unusually long
but valid lifecycle run can finish checkpoint, review, and promotion. Keep one
shared constant as the admission and decode-validity source. A request at the new
bound must still fail closed without changing the journal, while an idempotent
replay of an existing operation must remain available even when the journal is
full. Do not add compaction or discard prior operation evidence.

## Acceptance

- [ ] [OJ1] A valid run with 64 completed operations admits and durably records one more distinct prepared operation.
- [ ] [OJ2] A run at 128 operations refuses a new request without mutation while an exact replay of an existing operation remains idempotently readable.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OJ1 | restore the 64-entry operation limit | `TestOperationJournalAdmitsSixtyFifthEntry` | run `go test ./internal/specbuild -run '^TestOperationJournalAdmitsSixtyFifthEntry$' -count=1`; expect the 65th operation to be refused |
| OJ2 | remove or weaken the finite-bound refusal | `TestOperationJournalRetainsFiniteBoundAndExistingReplay` | run `go test ./internal/specbuild -run '^TestOperationJournalRetainsFiniteBoundAndExistingReplay$' -count=1`; expect the full-journal refusal or idempotent replay assertion to fail |
