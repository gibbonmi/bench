# Extend the bounded operation journal

Blocked by: allow-verified-clean-provisional-release.md
Ownership fence: `internal/specbuild/integrate.go`, `internal/specbuild/operation_journal_test.go`
Integration surfaces: lifecycle transition -> durable operation journal; long repair run -> checkpoint, review, and promotion transitions; decoded run state -> bounded validity check
Contracts: `internal/specbuild/integrate.go` admits the 65th distinct durable transition required by a long but valid repair run while retaining a finite 128-entry bound shared by mutation admission and state validation; `internal/specbuild/operation_journal_test.go` proves the former 64-entry boundary is live and the new upper boundary still refuses without mutation

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

- [ ] [MOJ1] Restoring the 64-entry limit makes the 65th-operation acceptance test red.
- [ ] [MOJ2] Removing or weakening the finite-bound refusal makes the full-journal test red.
