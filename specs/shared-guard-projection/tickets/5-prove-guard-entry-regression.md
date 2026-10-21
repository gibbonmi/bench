# Prove guard entry regression

Blocked by: 1-share-git-command-projections.md, 2-migrate-guard-envelope-consumers.md, 3-migrate-bench-projection-consumers.md
Writes: internal/systemtest/guard_projection_test.go (new), internal/systemtest/testdata/guard-projection-family-d.golden (new), cmd/bench/guard_projection_test.go (new)
Covers: GP34, GP35, GP36, GP37, GP38, GP39, GP40, GP41, GP44, GP45, GP46, GP71

## What to build

Prove through the actual hook and core entries that the migration changed no supported behavior.
Tickets 1–3 supply the migrated Git, Bench, reviewer, and write consumers; this ticket adds only regression witnesses and their reference fixture.

Run every actual hook and selected-core witness through the existing system process owner under `BENCH_KIT`, as the spec's Testing decisions state. Add no second subprocess harness.
GP46 compares the candidate with the pre-migration behavior over family D and the full decoder matrix, without elapsed-time equality.
The coordinator retains a pre-migration executable built before ticket 1 changes source.
Run the GP46 driver once against that executable and store its exact exit, stdout, and stderr tuples as `internal/systemtest/testdata/guard-projection-family-d.golden`.
Review that golden apart from the candidate code. The landed test then compares the candidate with the golden, so it needs no retained executable.
This regression repeats earlier cases; it does not supply a missing earlier test or repair.
GP44 checks census exclusion through the core entry in `cmd/bench/guard_projection_test.go`, which ticket 2 created.

Give the new system test file the `//go:build system` tag. Write no literal wait deadline in a new test (`wait_deadline_literal_test.go:33-39`).
Do not call `Skip`, `Skipf`, or `SkipNow` in a new test (`cmd/bench/skip_ownership_test.go:30-64`).
Each new file stays within 400 lines. Add no grant.

## Acceptance

- [ ] The landed degraded-guard witnesses for family A, parse advice, missing library, missing core, core error, and recovery stay unchanged and pass (GP34–GP37, GP39, GP40).
- [ ] The file-write hook keeps its warning-allow tuple without its library or core, and the follow-on hook keeps its warning-allow tuple without its resolver, wrapper, or core (GP38, GP71).
- [ ] Stale-core Bench and non-Bench calls keep their different refusal tuples (GP41).
- [ ] A supported Bench invocation makes no raw census record (GP44).
- [ ] An expansion-shaped command leaves the sentinel absent through the actual entry (GP45).
- [ ] Candidate tuples over family D and the decoder matrix equal `guard-projection-family-d.golden` (GP46).
- [ ] Changing one golden tuple makes `TestGuardProjectionDifferentialFamilies` red, and restoration returns the same green.
