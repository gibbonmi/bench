# Refuse a non-contention gate lock error with its cause

Blocked by: none
Writes: internal/gate/execution_probe_test.go, internal/gate/run_transaction.go
Covers: none

## What to build

Make `executeSubjectReplacing` classify an execution lock error through `gocache.Contended`. When a holder holds the lock, the gate reports `gate execution already in progress`. When any other error blocks the lock request, the gate refuses with `gate lock unavailable` and the error cause.

## Acceptance

- [ ] A contended execution lock reports `gate execution already in progress`.
- [ ] A non-contention execution lock error refuses with its cause and does not report `in progress`.
