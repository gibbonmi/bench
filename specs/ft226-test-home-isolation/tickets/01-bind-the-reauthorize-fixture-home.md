# Bind the reauthorize fixture's BENCH_HOME under its test root

Blocked by: none
Writes: internal/worktree/reauthorize_test.go

## What to build

`reauthorizeFixture` binds `BENCH_HOME` under the fixture's own repository root
for the test's duration, the same per-test binding its sibling fixture
`newOwnedAssignment` in `resume_test.go` already uses, before the first
`Create`. After this lands, `go test -count=1 ./internal/worktree` run against
an empty sentinel `BENCH_HOME` leaves the sentinel empty — the ten
`worktrees/001-*` entries the package writes today are gone. Spec rows IS1 (the
binding half), EV1.

## Acceptance

- [x] `BENCH_HOME=<empty sentinel dir> go test -count=1 ./internal/worktree`
      passes and the sentinel holds no entries afterwards.
- [x] Every `TestReauthorize*` creation path lies under `Pool(root)` for the
      fixture's own root, which resolves under that root.
- [x] `BENCH_HOME=<empty sentinel dir> go test -count=1 ./...` passes and the
      sentinel holds no entries afterwards; the before/after `find` output is
      recorded in the spec's verification log (EV1).
