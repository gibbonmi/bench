# Test capture transaction recovery in each interrupted state

Blocked by: none
Writes: internal/capturetx/recovery_test.go (new), internal/capturetx/store_test.go
Covers: none

## What to build

The quality survey of 2026-09-29, card 07, found that `internal/capturetx` owns crash recovery for capture writes but has one test and 17% statement coverage. This ticket adds tests only. It changes no production code.

Each test starts from a real repository in a temporary directory. A test that needs a crash state first runs `Begin` and then rewrites the manifest to the state that a crash leaves. Each test then asserts the result through `Append`, `Begin`, `Current`, `Sources`, `Commit`, or `Abort`. The tests use the shared `gittest.Repo` fixture. `store_test.go` also changes to use that fixture instead of its own `git init` call.

The tests cover these states:

- A staged directory from a crash before publish: no generation is open, and `Begin` succeeds.
- `preparing` with the live documents unchanged: recovery clears each live document and seals the generation.
- `preparing` with one live document already clear: recovery seals the generation.
- `preparing` with a changed live document: recovery refuses with "changed during cutover" and keeps the document.
- `sealed`: `Current`, a repeated `Begin`, and `Sources` return the open generation.
- `committing`: recovery retires the generation and keeps the live document.
- `aborting` before the restore: recovery puts the sealed text before the live suffix and retires the generation.
- `aborting` after the restore: recovery keeps the restored document and retires the generation.
- `aborting` with a changed live document: recovery refuses with "changed during abort".
- An unknown state, a manifest that does not parse, and an unsupported schema: recovery refuses each one.
- A sealed blob with a bad digest: `Sources` refuses with "failed its digest".
- `Commit` and `Abort` with an identifier that is not open: each call refuses.

The tests also cover a full `Commit` and a full `Abort` with no crash.

## Acceptance

- [ ] `internal/capturetx/recovery_test.go` has one test case for each state in the list.
- [ ] Each case asserts through an exported function of `internal/capturetx`.
- [ ] No production file in `internal/capturetx` changes.
- [ ] `store_test.go` uses `gittest.Repo` and calls `git init` nowhere.
- [ ] A `bench probe` mutation of two recovery lines turns a test red, and the probe restores the file.
- [ ] `go test ./internal/capturetx/` passes, and the statement coverage is above 17%.
