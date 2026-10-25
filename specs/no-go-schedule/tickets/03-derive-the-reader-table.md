# Derive and persist the reader table after a green whole kit run

Blocked by: 01-compute-the-oracle-identity.md, 02-record-a-test-log-per-package.md
Writes: internal/gate/reader_families.go (new), internal/gate/reader_table.go (new), internal/gate/reader_table_test.go (new), internal/gate/reader_store.go (new), internal/gate/reader_store_test.go (new), internal/gate/run_transaction.go, internal/git/git.go, internal/git/tree.go
Covers: NG13, NG14, NG15, NG16, NG17, NG18, NG19, NG20, NG21, NG22, NG23, NG24, NG25, NG26, NG27, NG28, NG71

## What to build

`reader_families.go` declares the eleven covered families once: `specs/`, `roadmap/`, `ROADMAP.md`, `capture/`, `reviews/`, `decisions/`, `CHANGELOG.md`, `.bench/commitment.json`, `.agents/`, `.bench/BENCH.md`, and `docs/`. A family that ends with `/` claims every path below it. It also declares the child-reader list, one row today: `internal/releaseevidence` with `CHANGELOG.md`.

The derivation is one pure function in `reader_table.go`. Its inputs are the package logs, the list of test packages, the graded root, the tracked symbolic links of the graded tree, the families, and the child rows. It applies the rules of the spec's "The reader table" section in order and returns a table or a refusal reason.

It refuses the table when a listed package has no log, when a log is empty or malformed, or when an operation word is unknown. It also refuses when `internal/conformance` does not read every covered family.

`internal/git` gains a reader of the tracked symbolic links of a tree, built on `ListTree` (`internal/git/tree.go:151`). It also gains the store's directory name beside `DeltaCacheFile` (`internal/git/git.go:177-179`).

`reader_store.go` owns the run's log directory and the table store. For a whole run of the kit's built-in table, the owner creates the log directory and sets `BENCH_GATE_TESTLOG` in the phase environment. After the run, it lists the test packages with `go list`, using the build flags and tags of the traced `test` call.

After a green verdict persists, the owner derives the table and persists it in the common Git directory, keyed by the graded tree, with the run's identity digest. A red run, an interrupted run, and a failed derivation persist no table, and a failed derivation leaves the green unchanged. The owner removes the log directory in every case. The store keeps at most 32 tables and never removes the table of the green marker's tree.

`internal/gate/run_transaction.go` holds 397 lines. Add only call sites there, and keep the file at 400 lines or fewer.

## Acceptance

- [ ] A log that opens `roadmap/FT1.md` makes its package a reader of `roadmap/`, and that table names the package for a change of `roadmap/FT2.md`.
- [ ] A read of `.claude/commands/x.md` under a tracked link to `../.agents/commands` makes the package a reader of `.agents/`.
- [ ] A read inside the package's own directory, a path outside the graded root, a `chdir` line, and a `getenv` line each add no family.
- [ ] A relative path after a `chdir` line resolves against the new directory.
- [ ] A malformed first line, an unknown operation, an empty log, and a missing log each refuse the table, and the refusal names the package.
- [ ] A table in which `internal/conformance` reads no `docs/` path is refused.
- [ ] The child row makes `internal/releaseevidence` a reader of `CHANGELOG.md` with no logged read.
- [ ] A green whole kit run with fixture logs persists one table keyed by its tree, with the run's identity digest.
- [ ] A red whole kit run persists no table, and a failed derivation keeps the green and persists no table.
- [ ] The store holds at most 32 tables and keeps the table of the green marker's tree.
- [ ] A whole run leaves no log directory behind.
