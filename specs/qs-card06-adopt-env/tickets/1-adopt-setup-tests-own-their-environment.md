# Make the interactive setup tests own their environment

Blocked by: none
Writes: internal/adopt/setup_prompt_test.go
Covers: none

## What to build

The quality survey of 2026-09-29, card 06, found that a plain `go test ./internal/adopt/` fails `TestSetupInteractiveSingleConfirm` and `TestSetupInteractiveResolvesAmbiguityOneAtATime`. The failure has an environment cause, not a code cause. The two tests pass under `bench test`.

Setup seeds `.bench/gate-inputs.json`, and the seed declares `BENCH_HOME`, `BENCH_KIT`, `BENCH_RUN_BINARY`, and `HOME`. Setup then runs the doctor rows. The gate inputs row is red when the process environment does not supply a declared name. A red row makes setup exit 3. The fixture sets only `BENCH_KIT`, so the result depends on the shell that starts `go test`.

The tests print only stderr on a failure, but the doctor writes the red row to stdout. Thus the failure message is empty.

Change the test file only:

- Add one helper that runs the interactive setup for both tests.
- Make the helper supply each name that the seeded manifest declares and that the shell does not supply. Read the names from `scaffoldGateInputs`, so the test keeps no second copy of the list.
- Make the helper print the setup stdout and stderr when the exit code is not zero.
- Replace the `slicesEqual` helper with `slices.Equal`.

## Acceptance

- [ ] A plain `go test -count=1 -run TestSetupInteractive ./internal/adopt/` passes in a shell that does not export `BENCH_RUN_BINARY`.
- [ ] The test file does not spell out the seeded environment names a second time.
- [ ] A probe that makes setup exit 3 reds both tests, and each failure prints the setup stdout and stderr.
- [ ] `internal/adopt/setup_prompt_test.go` has no `slicesEqual` helper.
- [ ] `go vet ./...` and the `internal/adopt` package tests pass.
