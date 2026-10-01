# Recover a capture transaction whose removal a crash cut short

Blocked by: none
Writes: internal/capturetx/transaction.go, internal/capturetx/public.go, internal/capturetx/recovery_test.go
Covers: none

## What to build

`Commit` and `Abort` in `internal/capturetx` end with the removal of the whole transaction directory. A crash can stop that removal after it deletes the manifest and before it deletes the blobs. Recovery then finds no manifest and does nothing. `Begin` then fails to publish its staged directory, because the old directory is not empty. Each later `Begin` fails the same way until a person removes the directory by hand.

The reviewer decision is to recover forward. `Begin` writes the manifest into the staged directory before it publishes that directory with one rename. Each later manifest write replaces the manifest with one rename, and it occurs only after a manifest read succeeds under the lock. Thus, under the lock, a transaction directory without a manifest is always a removal that `Commit` or `Abort` started. Recovery completes that removal. Recovery never removes a transaction directory that has a manifest.

A crash between the creation of the staged directory and its publish also leaves an inert staged directory. Nothing removes it. Recovery also removes each staged directory whose name has the exact prefix that `Begin` gives to its staged directory. Recovery runs under the lock, so no `Begin` can be in progress at that time. One constant holds that prefix for both the creation and the removal.

## Acceptance

- [ ] A transaction directory that holds only a blob and no manifest does not stop the next `Begin`, and `Begin` opens a fresh generation.
- [ ] Recovery removes a stale staged directory from a crash before publish.
- [ ] Recovery keeps a transaction directory that has a manifest. The existing tests `TestSealedGenerationStaysOpen` and `TestPreparingRecoversForward` cover this row.
- [ ] Each new test fails before the production change and passes after it.
- [ ] A `bench probe` mutation that removes the completion of the cut-short removal turns the blob-only test red, and the probe restores the file.
- [ ] `go test ./internal/capturetx/ ./internal/roadmap/` passes.
