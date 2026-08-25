# Mark the harness tests as parallel and close the census

Line: opus / medium.

The live-tree census can red across 53 test files, and that judgment needs the
mid tier.

Blocked by: 03-build-parallel-census-and-pins.md, 04-mark-five-slowest-tests-parallel.md, 05-mark-landing-surface-tests-parallel.md, 06-mark-landing-resume-tests-parallel.md, 07-mark-clean-and-landed-tests-parallel.md, 08-mark-lifecycle-and-orphan-tests-parallel.md, 09-mark-reclaim-and-list-tests-parallel.md, 10-mark-core-unit-tests-parallel.md
Writes: internal/worktree/reauthorize_test.go, internal/worktree/live_binary_test.go, internal/worktree/journey_test.go, internal/worktree/main_test.go, internal/worktree/test_run_test.go, internal/worktree/effect_census_test.go, internal/worktree/parallel_census_test.go, CHANGELOG.md

## What to build

The last six files run their eligible tests in parallel. Every eligible test in
the package then calls `t.Parallel()`, so the live-tree census can turn green.

Add `t.Parallel()` to every eligible test in this file set. A test that binds
environment or changes the directory stays serial. `TestMain` is neither
eligible nor serial, and it keeps its shape.

Add the live-tree census test. It runs the census of ticket 03 over the
package's own `_test.go` files and reports each mismatch by file and line.

Record the user-facing change in `CHANGELOG.md` under Unreleased.

## Acceptance

- [ ] WF01 pins the parallel census on the live tree.
- [ ] Every eligible test in the six files calls `t.Parallel()`.
- [ ] The package passes `go test -race`.
- [ ] `CHANGELOG.md` names the parallel `internal/worktree` suite under Unreleased.
