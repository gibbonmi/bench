# Prove guard entry regression

Blocked by: 1-share-git-command-projections.md, 2-add-envelope-owner-for-git.md, 3-migrate-bench-and-write-readers.md, 4-migrate-bench-projection-consumers.md
Writes: internal/systemtest/guard_projection_test.go (new), cmd/bench/guard_projection_test.go (new)
Covers: GP34, GP35, GP36, GP37, GP38, GP39, GP40, GP41, GP44, GP45, GP71

## What to build

Prove through the actual hook and core entries that the migrated guards keep each degraded and recovery behavior.
Tickets 1–4 supply the migrated Git, Bench, reviewer, and write consumers; this ticket adds only regression witnesses.

Run every actual hook and selected-core witness through the existing system process owner under `BENCH_KIT`, as the spec's Testing decisions state. Add no second subprocess harness.
Most rows cite the landed degraded-guard witnesses; new tests cover only GP38, GP44, and GP45.
GP44 checks census exclusion through the core entry in `cmd/bench/guard_projection_test.go`, which ticket 2 created.
This regression repeats earlier cases; it does not supply a missing earlier test or repair.

Prefix each new package-level helper name in `internal/systemtest` with `projection`, so it cannot collide with a helper that ticket 7 adds to the same package.
Give the new system test file the `//go:build system` tag. Write no literal wait deadline in a new test (`internal/conformance/wait_deadline_literal_test.go:33-39`).
Do not call `Skip`, `Skipf`, or `SkipNow` in a new test (`internal/conformance/skip_ownership_test.go:30-64`).
Each new file stays within 400 lines. Add no grant.

## Acceptance

- [ ] The landed degraded-guard witnesses for family A, parse advice, missing library, missing core, core error, and recovery stay unchanged and pass (GP34–GP37, GP39, GP40).
- [ ] The file-write hook keeps its warning-allow tuple without its library or core, and the follow-on hook keeps its warning-allow tuple without its resolver, wrapper, or core (GP38, GP71).
- [ ] Stale-core Bench and non-Bench calls keep their different refusal tuples (GP41).
- [ ] A supported Bench invocation makes no raw census record (GP44).
- [ ] An expansion-shaped command leaves the sentinel absent through the actual entry (GP45).
