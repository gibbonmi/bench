# Close the review findings

Line: opus / medium.

Blocked by: 16-pin-the-serial-ceiling-and-close.md
Writes: internal/racetests/racetests.go, internal/worktree/parallel_census_test.go, cmd/bench/command_registry_test.go, reviews/worktree-test-floor.md (delete)

## What to build

Three accepted findings from `reviews/worktree-test-floor.md`.

`TestParallelJourneysRecordEverySelection` joins `racetests.Tests`, so the
gate's `race` phase runs WF06's test under `-race`.

The census gains one serial edge. An assignment whose left-hand side is a
selector on an imported package (`pkg.Ident = …`) marks the test serial, in
the class of `bindEnv`. A synthetic test plants an import and the assignment
and expects the serial verdict. A second expects the pair report when the
test also calls `t.Parallel()`. The live-tree census and the serial ceiling
stay green; the ceiling moves only if the live count changes, and the ticket
states the new count.

`TestWorktreeRoutesKeepTheirBytesFromASubdirectory` runs its two invocations
against a fixture repository from `gittest`, not the live checkout, so a
worktree change during the gate cannot flip it. The bytes it pins stay the
same in kind: the root run and the subdirectory run agree.

The green fix commit deletes `reviews/worktree-test-floor.md`.

## Acceptance

- [ ] `racetests.Tests` names `TestParallelJourneysRecordEverySelection`, and the race runner's expected-runs check counts it.
- [ ] A synthetic test that assigns an imported package's variable is serial under the census, and the pair is reported with `t.Parallel()`.
- [ ] `TestWorktreeRoutesKeepTheirBytesFromASubdirectory` passes while a worktree is created between its two runs.
- [ ] WF01 and WF18 stay green; `go test -race -count=1 -parallel 2 ./internal/worktree` passes.
