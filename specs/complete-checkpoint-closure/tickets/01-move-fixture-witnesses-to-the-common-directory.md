# Move the gate fixture witnesses to the Git common directory

Blocked by: none
Writes: internal/gate/run_outcomes_test.go, internal/gate/review_checkpoint_test.go, internal/gate/review_checkpoint_commits_test.go, internal/gate/run_failure_outcomes_test.go
Covers: none

## What to build

This ticket is test-only prefactoring for ticket 2.
A prospective run executes the project gate in a private linked checkout, and the run removes that checkout.
A witness file that the fixture gate writes in the checkout therefore disappears after a prospective run.
This ticket moves each witness to the Git common directory, which the primary checkout and each linked checkout share.

The `outcomeFixture` gate script resolves `git rev-parse --path-format=absolute --git-common-dir` once.
It writes `.gate-run-count` and `.gate-record-during` in that directory.
It copies `bench-last-gate` from that directory.
The other fixture markers, such as `.gate-red` and `.gate-wait`, stay in the checkout.

One test helper owns the witness path for every reader.
These readers move to it:

- `outcomeRuns` in `internal/gate/run_outcomes_test.go`, which keeps its signature so that `internal/gate/timeout_recovery_count_test.go` needs no edit
- the `.gate-record-during` read in `internal/gate/run_outcomes_test.go`
- the no-run assertion in `internal/gate/review_checkpoint_test.go`
- the two no-run assertions in `internal/gate/review_checkpoint_commits_test.go`
- the no-run and no-record assertions and the record wait in `internal/gate/run_failure_outcomes_test.go`

No assertion changes its meaning.
Today every fixture run is an ordinary run in the primary checkout, where the Git directory and the Git common directory are the same.
So the moved assertions stay green on the current gate.
CC25 and CC32 prove the move against the prospective route in ticket 2.

## Acceptance

- [ ] After a green fixture run, `.gate-run-count` and `.gate-record-during` exist in the Git common directory and do not exist in the checkout.
- [ ] `rg` over `internal/gate` finds `.gate-run-count` and `.gate-record-during` only in the fixture script and in the one witness helper.
- [ ] Each reader in the list above resolves its witness path through the one helper.
- [ ] The `internal/gate` package tests pass with no assertion removed or weakened.
