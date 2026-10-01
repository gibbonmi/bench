# Strip the wrapper routing from the environment of a focused run

Blocked by: none
Writes: internal/testreport/environment.go, internal/testreport/environment_test.go, CHANGELOG.md
Covers: none

## What to build

Roadmap item FT327 reports that one adoption system test fails in a clean worktree.
`bench worktree exec` sets `BENCH_WRAPPER` to the `bin/bench.sh` of the worktree.
`bench test` gives that value to its Go child, because the focused run environment
does not remove the wrapper routing. The `bench setup` call in the system suite then
grades the old broker manifest of the worktree and exits 3.

The gate already removes the wrapper routing from its children. The focused run
environment now uses the same owner, `env.WithoutWrapperRouting`, so the Go child
carries no `BENCH_WRAPPER` and no `BENCH_KIT`. That owner holds the `BENCH_KIT` name,
so the conformance strip list does not repeat it. A conformance run still adds
`BENCH_KIT` again with the source root of the selected executable.

## Acceptance

- [x] The focused run environment carries no wrapper routing entry when the caller environment has `BENCH_WRAPPER` and `BENCH_KIT`.
- [x] The conformance strip list does not name `BENCH_KIT`.
- [x] Removal of the wrapper routing strip turns `TestTestEnvironmentStripsWrapperRouting` red.

## Verification

`TestTestEnvironmentStripsWrapperRouting` was red before the production edit, because
the child carried `BENCH_WRAPPER`. The `internal/testreport` packages pass after the
edit. A `bench probe` run removed the wrapper routing strip, and the new test was red.
The probe restored the file.
