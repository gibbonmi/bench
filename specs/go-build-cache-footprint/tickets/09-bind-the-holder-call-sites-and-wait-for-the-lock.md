# Bind the holder call sites and wait for the lock

Blocked by: none

Writes: internal/gocache/lock.go, internal/gocache/lock_test.go, internal/gocache/clean_test.go, internal/gate/lane.go, internal/gate/lane_test.go, internal/gate/run_transaction.go, internal/gate/run_failure_outcomes_test.go, internal/testreport/testreport.go, internal/testreport/testreport_test.go, internal/bounds/bounds.go, internal/conformance/bounds_policy_test.go, specs/go-build-cache-footprint/spec.md

## What to repair

Review finding (Spec 1, Coverage 1): no test binds the three production
`gocache.Hold` call sites. `TestCleanRefusesWhileEachHolderRuns` calls `Hold`
directly under the labels gate, focused, and lane. Delete any one of the three
call sites and the suite stays green.

Bind each site. The production entries are the gate run in
`run_transaction.go`, `RunLane` in `lane.go`, and `testreport.Command`. For
each entry, a test drives it and observes the lock while the run is in
flight. The observer is a second process or the lock file's state. If the
production path cannot pause, an injected holder seam is acceptable. The test
must then also prove that the production path calls that seam.

Review finding (Coverage 2), reviewer decision 2026-08-27: a holder waits
unbounded. Replace the `F_SETLK` retry loop in `acquireShared` with one
`F_SETLKW` request. Remove `CacheHoldWait` from `internal/bounds/bounds.go`
and its owners row from `internal/conformance/bounds_policy_test.go`. The
clean stays no-wait.

## Acceptance

- [ ] L01 — While a gate run holds the cache lock, `bench cache clean` exits 1.
- [ ] L02 — While a `bench test` run holds the cache lock, `bench cache clean` exits 1.
- [ ] L03 — While a lane run holds the cache lock, `bench cache clean` exits 1.

Delivered outcome: each holder's production path is proven to take the lock,
and a holder never runs unlocked beside a clean.
