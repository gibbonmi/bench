# Align the upgrade hook plan

Blocked by: restrict-upgrade-hook-absence.md
Ownership fence: `internal/adopt/link_hook.go`, `internal/adopt/link_stage.go`, `internal/adopt/upgrade.go`, `internal/contract/surface/upgrade_test.go`
Integration surfaces: prospective rendered hook bytes→`internal/adopt/link_stage.go` + AP1; upgrade count→`internal/adopt/upgrade.go` + AP1; consumer upgrade contract→`internal/contract/surface/upgrade_test.go` + AP1
Contracts: the transaction's prospective rendered bytes cross `internal/adopt/link_stage.go`→`internal/adopt/upgrade.go`, asserted by AP1 against the real staged producer
Closure: AP1/rendered-byte-plan-parity

## What to build

Make unequal-version `bench upgrade --check` count the pre-push write exactly when the same accepted transaction would stage bytes different from the effective installed hook. The predicate is already narrowed by `restrict-upgrade-hook-absence.md`; a rendered-bytes-only cut must keep the producer, count consumer, and real upgrade fixture together because the live-branch-change contract otherwise has no consumer that can compare the transaction's bytes with the reported plan.

## Acceptance

- [ ] [AP1] On an unequal-version upgrade, `--check` adds one hook change exactly when the accepted transaction's rendered pre-push bytes differ from the effective installed hook, including a live branch that changed after an otherwise-current hook was installed.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AP1/rendered-byte-plan-parity | classify a byte-current hook as no refresh without comparing the transaction's new rendered bytes | upgrade surface contract | install a current hook for one baked branch, resolve a different live branch, run unequal-version `upgrade --check`, expect the missing change count red |
