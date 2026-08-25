# Convert the direct home binds in the land and clean files

Line: sonnet / low.

Every verb takes the home, so a direct bind becomes an argument. The census
and `go test -race` grade the result.

Blocked by: 13-fixtures-pass-a-private-home.md
Writes: internal/worktree/land_flags_test.go, internal/worktree/land_release_refusal_test.go, internal/worktree/land_reauthorization_test.go, internal/worktree/land_facts_test.go, internal/worktree/land_resume_test.go, internal/worktree/land_journey_test.go, internal/worktree/clean_landed_test.go, internal/worktree/clean_landed_hostile_test.go, internal/worktree/clean_landed_apply_test.go, internal/worktree/landed_test.go, internal/worktree/squash_landed_test.go

## What to build

Each test in these files that calls `bindEnv(t, "BENCH_HOME", …)` for an
in-process verb passes that home to the verb instead and drops the bind. The
test then calls `t.Parallel()`. A test that binds another key for an
in-process read, or that grades the boundary default, stays serial. Remove no
test, and change no assertion.

## Acceptance

- [ ] No `bindEnv(t, "BENCH_HOME", …)` remains in these files except one that grades the boundary default.
- [ ] WF01 stays green on the live tree.
- [ ] The package passes `go test -race`.
