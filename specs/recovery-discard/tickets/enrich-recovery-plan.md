# Split orphaned from absent and report what a discard would drop

Blocked by: none
Ownership fence: `internal/worktree`, `internal/contract/runtime/runtime_worktree_test.go`
Contracts: the orphaned/absent verdict and the change-summary value cross `internal/worktree/classifier.go`→`internal/worktree/worktree.go` renderer and on into the `recovery_cleanup` receipt, asserted by EP1, EP2 and EP4 against real `PlanRecovery` output over refs produced through the package's own owned-assignment helpers
Assumptions: `PlanRecovery` stays the one classifier and no second planner appears; the ref-existence probe is a ref read, not a record read; the change summary is derived at plan time and persisted nowhere; claims re-derived from the tree at pickup

## What to build

`PlanRecovery` today collapses two different states onto one `retain` verdict
with the detail `recovery ref has no recovered assignment`: a recovery ref that
exists but whose owning assignment row is gone, and a ref that does not exist at
all. One is actionable and the other is already done. Split them, and give every
plan a derived summary of what its payload changes against its recorded base, so
an operator reading the plan can tell a one-file leftover from an abandoned
build before anything is dropped.

The classifier gains a ref-existence probe it does not have today. The change
summary is computed against the recorded base at plan time and degrades to a
definite unknown — never an error — when the base or the payload no longer
resolves, because an error there would make an otherwise-actionable ref
unplannable.

This ticket does not add the discard verb; it makes the plan able to say what
discard will later act on. The emitted receipt gains the summary column for
every action, including `retain` and `retire`.

## Acceptance

- [ ] [EP1] planning a recovery ref that exists with no owning assignment row reports an orphaned verdict distinguishable in the receipt from both `retain` and the absent verdict.
- [ ] [EP2] planning a ref name that resolves to no ref reports an absent verdict distinguishable in the receipt from the orphaned verdict.
- [ ] [EP3] the pre-existing verdicts are unchanged: a recovered assignment whose payloads all pass the landedness proof still plans `retire`, and one whose payloads do not still plans `retain`.
- [ ] [EP4] every plan carries a changed-path count for its payload against its recorded base, present for `retain`, `retire`, orphaned and absent plans alike.
- [ ] [EP5] a plan whose recorded base or payload no longer resolves carries a definite unknown in the change-summary field and returns no error.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| EP1 | make the row-less branch return the plain `retain` verdict the tree returns today | the orphaned-verdict plan test | drop the orphan arm in `PlanRecovery`, run `go test ./internal/worktree -run Recovery -timeout 120s`, expect the orphaned-verdict assertion to fail |
| EP2 | make the ref-existence probe report every name as present | the absent-verdict plan test | return true unconditionally from the probe, run `go test ./internal/worktree -run Recovery -timeout 120s`, expect the absent-verdict assertion to fail |
| EP3 | make the landedness aggregation accept a single landed payload instead of all | the retire/retain regression test | flip the `all` fold to `any` in `PlanRecovery`, run `go test ./internal/worktree -run Recovery -timeout 120s`, expect the retain-verdict assertion to fail |
| EP4 | return a fixed zero for the changed-path count | the change-summary plan test | hardcode the count, run `go test ./internal/worktree -run Recovery -timeout 120s`, expect the count assertion to fail |
| EP5 | return an error from the change summary when the base does not resolve | the honest-unknown plan test | propagate the resolve error out of `PlanRecovery`, run `go test ./internal/worktree -run Recovery -timeout 120s`, expect the no-error assertion to fail |
