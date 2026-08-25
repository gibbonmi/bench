# Guard the harness state for parallel siblings

Line: opus / medium.

Blocked by: none
Writes: internal/worktree/test_run_test.go, internal/worktree/journey_test.go, internal/worktree/journey_race_test.go (new), internal/racetests/racetests.go, internal/gittest/gittest.go

## What to build

Two parallel journeys share the harness and record every effect. The gate's
`race` phase reports no data race.

Give the run-binary selector's `journeys` log a mutex. Keep the `sync.Once`
selection, so two parallel journeys still select one executable. Keep
`descendant` and its per-test cleanup as they are.

Add one worktree race sentinel. The sentinel runs two parallel journeys against
two repositories and reads the harness effect log. Name the sentinel in the
race-test registry, so the gate's `race` phase runs it.

Give `gittest.StubGit` a form that returns the stub directory and calls no
`t.Setenv`. A journey then passes the stub path on the child's `Env` and stays
eligible.

## Acceptance

- [ ] WF06 pins two goroutines that call the run-binary selector, and the log holds two journey lines.
- [ ] WF07 pins one worktree sentinel in the race-test registry.
- [ ] The sentinel passes under `go test -race`.
