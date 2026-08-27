# Hold the cache lock and add bench cache clean

Blocked by: 04-add-the-bench-cache-verb.md

Writes: internal/gocache/ (new), internal/gate/run_transaction.go, internal/gate/lane.go, internal/testreport/testreport.go, cmd/bench/main.go, cmd/bench/command_registry_test.go, internal/conformance/subcommand_routing_test.go, projects/benchkit.md, CHANGELOG.md

## What to build

An operator reclaims the disk on demand, and no running compile loses an entry.
Add the cache lock and the `clean` child of `bench cache`.

The cache lock is a POSIX record lock on `bench.lock` inside the Bench build
cache. A holder takes a read lock and keeps its descriptor open for the run's
span. Extend `internal/gocache/` with the lock, and mirror the gate execution
lock in `internal/gate/run_transaction.go`. That file already holds
`recordLock`, `acquireExecutionLock`, `unlockExecutionLock`, and `lockHeld`.

Three holders take the shared lock. The gate run takes it in
`internal/gate/run_transaction.go`, and its span opens before the run-binary
build and closes after teardown. The lane run takes it in
`internal/gate/lane.go` for its span. The focused run takes it in
`internal/testreport/testreport.go` for its span. Two holders acquire the
shared lock at the same time, so the bound never serializes gates.

`bench cache clean` requests a write lock with a no-wait set, and `EAGAIN` is
the refusal. On a refusal it exits 1, removes no file, and names the blocking
pid that `F_GETLK` reports, in the form `pid <n>`. With no holder it measures
the footprint before the removal. Then it runs `go clean -cache` against the
directory and reports the bytes and the files that measurement names. Go's own clean
removes the two-hex subdirectories only, so `bench.lock`, `trim.txt`, and
`README` survive. An absent directory reports zero removed and exits 0, and a
missing `go` on `PATH` is a refusal that names `go`.

Add `clean` to the command registry in `cmd/bench/main.go` as a mutation child
of the `cache` verb that `04-add-the-bench-cache-verb.md` registered. The
`setup` row is the mutating sibling to mirror. Update the routing map and the
`projects/benchkit.md` note beside it.

## Acceptance

- [ ] L01 — While a gate run holds the cache lock, `bench cache clean` exits 1.
- [ ] L02 — While a `bench test` run holds the cache lock, `bench cache clean` exits 1.
- [ ] L03 — While a lane run holds the cache lock, `bench cache clean` exits 1.
- [ ] L04 — A refused `bench cache clean` removes no file from the directory.
- [ ] L05 — Two holders acquire the shared lock at the same time.
- [ ] L06 — With no holder, `bench cache clean` removes every two-hex subdirectory.
- [ ] L07 — `bench cache clean` leaves `bench.lock`, `trim.txt`, and `README` in place.
- [ ] L08 — `bench cache clean` on an absent directory reports zero removed and exits 0.
- [ ] L10 — `bench cache clean` with no `go` on `PATH` refuses and names `go`.
- [ ] L11 — The clean refusal names the holder's pid as `pid <n>`.
- [ ] L12 — With no holder, `bench cache clean` reports the bytes and files removed as measured before the removal.

Delivered outcome: an operator empties the Bench build cache on demand, and a
concurrent gate keeps every archive it is reading.
