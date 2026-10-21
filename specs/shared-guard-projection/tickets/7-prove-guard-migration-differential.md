# Prove the guard migration differential

Blocked by: 1-share-git-command-projections.md, 2-add-envelope-owner-for-git.md, 3-migrate-bench-and-write-readers.md, 4-migrate-bench-projection-consumers.md
Writes: internal/systemtest/guard_projection_differential_test.go (new), internal/systemtest/testdata/guard-projection-family-d.golden (new)
Covers: GP46

## What to build

Prove that the complete migration changed no supported caller result, by comparison with recorded pre-migration behavior.
Tickets 1–4 supply the migrated consumers. The coordinator supplies a Bench worktree at the pre-migration build-entry commit, built before ticket 1 changed source, and records that commit in the build plan. The spec's Graph ordering section names the commit and the worktree label. If the worktree is lost, recreate and rebuild it from that commit.

Write the GP46 driver in `guard_projection_differential_test.go`. It runs family D and the full decoder preservation matrix through the actual core and hook entries, using the existing system process owner under `BENCH_KIT`.
The system owner verifies the run binary's seal against `BENCH_KIT` (`internal/systemtest/owner_test.go:86-91`), so record with `BENCH_RUN_BINARY` set to the pre-migration worktree's executable and `BENCH_KIT` set to that worktree.
If `bench test --check system` cannot select that pair, run the driver once through `go test -tags system` with those two variables, and report that raw call.
Store the recorded exit, stdout, and stderr tuples as `internal/systemtest/testdata/guard-projection-family-d.golden`. The landed test only compares; add no update flag to it.
Before you record or compare a tuple, replace each temporary fixture root in it with one fixed token, so the golden matches across runs.
Review the golden apart from the candidate code. The landed test compares the candidate with the golden, so it needs no retained executable.
Feed the ordered raw envelopes without JSON remarshal or key normalization. Compare no elapsed time.
This ticket authorizes no intentional difference.

Prefix each new package-level helper name with `differential`, so it cannot collide with a helper that ticket 6 adds to the same package.
Give the new test file the `//go:build system` tag. Write no literal wait deadline (`internal/conformance/wait_deadline_literal_test.go:33-39`) and no `Skip` call (`internal/conformance/skip_ownership_test.go:30-64`).
Each new file stays within 400 lines. Add no grant.

## Acceptance

- [ ] Candidate tuples over family D and the decoder matrix equal `guard-projection-family-d.golden` (GP46).
- [ ] Swapping the Git caller's first-child selection for an every-child search makes `TestGuardProjectionDifferentialFamilies` red on the W4 member, and restoration returns the same green.
- [ ] The golden holds no temporary path and no elapsed time.
