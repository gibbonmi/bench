# Remove the cmd test home

Blocked by: none
Writes: cmd/bench/spill_support_test.go, cmd/bench/gate_route_test.go
Covers: none

## What to build

A `cmd/bench` test run leaves no `bench-cmd-home-*` directory in the temporary directory. Today each run leaves ten of them.

The cause is the gate route owner process. The `gateRoute` helper runs the public shell wrapper, and the wrapper starts the test binary again as the owner process. That child runs `TestMain`, which makes a private home. Then `TestGateRouteOwnerProcess` calls `os.Exit` inside `m.Run`, so the removal in `TestMain` never runs. Each of the ten `gateRoute` calls leaves one home. The child also ignores the `BENCH_HOME` that `gateRoute` sets, because `TestMain` replaces it.

`TestMain` runs the owner process before it makes the private home. The owner process uses the home that `gateRoute` gives it, and that home is inside the test's own temporary directory. `TestGateRouteOwnerProcess` goes away. If `TestMain` cannot remove the private home, it writes the error and the run fails.

`gateRoute` gives the child a private temporary directory and fails when that directory is not empty after the child exits. Thus a child that leaves a temporary entry turns the gate route tests red.

## Acceptance

- [ ] One `go test -count=1 ./cmd/bench/` run with a private `TMPDIR` leaves the same number of `bench-cmd-home-*` entries as before the run. The author reports both counts and the wall time.
- [ ] If `TestMain` cannot remove the private home, it writes the error to stderr and the run exits nonzero.
- [ ] The gate route tests fail when the owner process leaves an entry in its temporary directory. A mutation probe that makes the private home before the owner process starts turns them red.
- [ ] `bench test --package ./cmd/bench/` and `bench test --check skip-ownership` pass.
