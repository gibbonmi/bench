# Export the gate execution probe

Blocked by: none
Writes: internal/gate/run_transaction.go, internal/gate/execution_probe_test.go (new)
Covers: PB42

## What to build

Verify the premise first. Read `lockHeld` and `acquireExecutionLock` in
internal/gate/run_transaction.go. Read the one caller of `lockHeld` in
internal/gate/verdict.go. Read `startGateLockHolder` and `runGateLockHolder` in
internal/gate/run_failure_outcomes_test.go, which hold the lock from a child process.

Add `ExecutionInProgress(root string) (bool, error)` beside `lockHeld`. It resolves the
checkout administration directory through `benchgit.AdminDir` and answers `lockHeld`.
A root that is not a repository returns the adapter's error. Keep `lockHeld` and its
inspection caller unchanged.

Write `TestExecutionInProgressReadsTheLock` in the new file. It answers false over a
fresh repository, and true while `startGateLockHolder` holds the lock. It answers
false after the release, and an error over a temporary directory that is not a
repository.

Self-probe: make the reader stat the lock file instead of reading the lock and show
the after-release case red.

## Acceptance

- [ ] `gate.ExecutionInProgress` answers true under a child lock holder and false after the release.
- [ ] `gate.ExecutionInProgress` returns an error over a directory that is not a repository.
- [ ] `go test ./internal/gate -run 'ExecutionInProgress|Inspect'` passes.
