# Preserve focused system verdicts

Blocked by: none
Writes: internal/systemtest/owner_test.go, internal/systemtest/owner_selection_test.go
Covers: HP11, HP12, HP13, HP19

## What to build

Make the system-suite owner distinguish an unfiltered suite from a run selected
with Go's parsed `test.run` flag. A focused run returns the selected tests'
verdict after the normal owner cleanup and does not apply the full-suite ledger
verification. An unfiltered run keeps the current ledger verification and its
failure diagnostics.
Bind `BENCH_KIT` to the assignment worktree for every focused and unfiltered
system-suite invocation.

## Acceptance

- [ ] A passing `-run` selection exits 0 without a missing-ledger trailer.
- [ ] A failing `-run` selection keeps the selected test's nonzero verdict without replacing it with a ledger diagnostic.
- [ ] An unfiltered system suite still verifies repository, executable, and terminal-outcome observations.
- [ ] Owner cleanup runs for focused and unfiltered invocations.
