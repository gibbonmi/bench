# Finish guard scan workers before timeout return

Blocked by: none

## What to build

`guards.Scan` preserves its timeout result while joining the enumeration and
inspection work it started, and the dev gate's existing race phase permanently
checks the two guards timeout paths alongside the worktree cleanup race.

## Acceptance

- [x] An enumeration or inspection timeout returns the existing incomplete
  result only after the worker started by `Scan` has finished; no worker leaks
  into the next test.
- [x] `go test -race -count=1 ./internal/guards` passes, including the FIFO
  non-blocking proof under race instrumentation.
- [x] The dev gate race phase runs both guards timeout tests and the existing
  worktree cleanup race, fails closed if a named test does not execute, and is
  observed red against the pre-fix `guards.Scan` race before going green.
