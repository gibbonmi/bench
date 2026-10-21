# Prove the guard migration differential

Blocked by: 1-share-git-command-projections.md, 2-add-envelope-owner-for-git.md, 3-migrate-bench-and-write-readers.md, 4-migrate-bench-projection-consumers.md
Writes: internal/systemtest/guard_projection_differential_test.go (new), internal/systemtest/testdata/guard-projection-family-d.golden (new)
Covers: GP46

## What to build

Prove that the complete migration changed no supported caller result, by comparison with recorded pre-migration behavior.
Tickets 1–4 supply the migrated consumers. The coordinator supplies a pre-migration executable, built before ticket 1 changed source, and records its build commit in the build plan. If the executable is lost, rebuild it from that commit.

Write the GP46 driver in `guard_projection_differential_test.go`. It runs family D and the full decoder preservation matrix through the actual core and hook entries, using the existing system process owner under `BENCH_KIT`.
Run the driver once against the pre-migration executable, and store its exact exit, stdout, and stderr tuples as `internal/systemtest/testdata/guard-projection-family-d.golden`.
Before you record or compare a tuple, replace each temporary fixture root in it with one fixed token, so the golden matches across runs.
Review the golden apart from the candidate code. The landed test compares the candidate with the golden, so it needs no retained executable.
Feed the ordered raw envelopes without JSON remarshal or key normalization. Compare no elapsed time.
This ticket authorizes no intentional difference.

Give the new test file the `//go:build system` tag. Write no literal wait deadline (`wait_deadline_literal_test.go:33-39`) and no `Skip` call (`cmd/bench/skip_ownership_test.go:30-64`).
Each new file stays within 400 lines. Add no grant.

## Acceptance

- [ ] Candidate tuples over family D and the decoder matrix equal `guard-projection-family-d.golden` (GP46).
- [ ] Swapping the Git caller's first-child selection for an every-child search makes `TestGuardProjectionDifferentialFamilies` red on the W4 member, and restoration returns the same green.
- [ ] The golden holds no temporary path and no elapsed time.
