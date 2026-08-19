# Run the worktree test binary under a private BENCH_HOME that reds on residue

Blocked by: 01-bind-the-reauthorize-fixture-home.md
Writes: internal/worktree/main_test.go, internal/worktree/reauthorize_test.go (probe (a) only, reverted)

## What to build

A `TestMain` in `internal/worktree` creates an empty process-private directory
with `os.MkdirTemp`, records the inherited `BENCH_HOME`, exports the private
path as `BENCH_HOME`, runs the tests, computes residue under the private home
through one unexported predicate, removes the private home, and exits non-zero
if `m.Run` failed, residue exists, or removal failed. Creation failure exits
non-zero naming the error before any test runs. The residue report on stderr is
one heading line (count and private home path) then one line per residue path;
for each residue worktree holding a `.git` pointer file, its `gitdir:` line is
appended so the report names the originating test's temporary root. The
predicate is unit-tested on a synthetic home; `TestPackageHomeIsPrivate` asserts
the exported home at test start. Spec rows IS1, IS2, IS3, DT1, DT2, DT3, DT4, EX1.

The contract with ticket 01: this ticket is green only because ticket 01 landed
— with the fixture binding removed, this `TestMain` must turn the package red
naming a `TestReauthorize…` root, and that red is the build's recorded mutation
probe.

## Acceptance

- [x] `TestPackageHomeIsPrivate`: at test start `BENCH_HOME` exists and is
      empty, differs from the recorded inherited value, is not under it, is not
      under `$HOME/.bench`, and `Pool` of a fresh root resolves under it.
- [x] `TestHomeResidueListsLeakedWorktreesWithOrigin`: a synthetic home with a
      planted `worktrees/001-1/<id>/.git` pointer and a stray top-level file
      reports exactly two top-level entries plus the planted `gitdir:` target;
      a home holding only an empty `worktrees/` directory reports exactly one
      entry; an empty home reports nothing.
- [x] Probe (a), recorded: with the ticket 01 binding temporarily removed,
      `go test -count=1 ./internal/worktree -run Reauthorize` exits non-zero
      and the report names a `TestReauthorize…` temporary root; restored, it
      exits zero.
- [x] Probe (b), recorded: `TMPDIR=/nonexistent go test -count=1
      ./internal/worktree -run TestPool` exits non-zero with the `TestMain`
      creation error and runs no test.
- [x] Probe (c), recorded: a temporary `t.Fatal` in one test with a clean home
      makes the package FAIL; removed, the package passes.
- [x] After `go test -count=1 ./internal/worktree`, the private home path the
      run printed (or a `find "$TMPDIR" -maxdepth 1 -name '<prefix>*'`) shows no
      leftover directory.
- [x] `TestPoolDefaultBenchHome` is unchanged and passes.
- [x] `go test -count=1 ./internal/worktree -run TestActualSIGINT` passes under
      the new `TestMain` (the re-exec'd helper child gets its own private home).
- [x] `go test -race -count=1 -v ./internal/worktree -run
      TestConcurrentCleanupRecordsOneTransaction` (the race phase's worktree
      entry) passes under the new `TestMain`.
