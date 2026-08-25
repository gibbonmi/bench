# The shared fixtures pass a private home

Line: opus / medium.

Blocked by: 12-take-explicit-home-in-every-verb.md
Writes: internal/worktree/land_fixtures_test.go, internal/worktree/resume_test.go, internal/worktree/pool_reclaim_test.go, internal/worktree/reauthorize_test.go, internal/worktree/live_binary_test.go, internal/worktree/land_freshness_test.go, and every `internal/worktree/*_test.go` file the census names once the fixtures stop binding (marks only)

## What to build

`landingFixtureAtHome` binds nothing. It writes the tally path into the two
gate scripts as a literal. It declares an empty `environment` in
`gate-inputs.json`. It passes its home to `mustCreate` and to every verb it
runs. `newOwnedAssignment`, `newPendingAssignment`, `mustCreate`,
`newReclaimPool`, `reauthorizeFixture`, `newResidueGuardFixture`, and
`redProspectiveGateLanding` take or own a private home under `t.TempDir()` and
pass it down. A test that grades the boundary default keeps its bind.

Once the fixtures stop binding, the live-tree census reports every test they
serve as eligible. Add `t.Parallel()` to each one in the same ticket, so the
census stays green at the commit. Run the whole package under `-race`.

## Acceptance

- [ ] WF17 pins the landing fixture on no environment bind, and a test it serves calls `t.Parallel()` under the census.
- [ ] WF01 stays green on the live tree.
- [ ] The package passes `go test -race`, and the home-residue guard reports no write below the operator's home.
