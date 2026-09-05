# Migrate the gate and dashboard sites to the reader

Blocked by: add-the-named-git-admin-readers.md
Writes: internal/gate/lane.go, internal/gate/verdict.go, internal/gate/run_transaction.go, internal/gate/composed_green.go, internal/gate/phases.go, internal/gate/admin_readers_test.go (new), internal/dashboard/dashboard.go, internal/dashboard/dashboard_test.go, tests/canary/package-core-guard/git-flag-retyped
Covers: GR13, GR14, GR15, GR16, GR17

## What to build

Verify the premise first. Read `AdminDir` in internal/git/worktree_admin.go, which
the blocker ticket adds. Read `runLane` at internal/gate/lane.go line 100. Read
`inspectSubjectAt` at internal/gate/verdict.go line 194. Read
`executeSubjectWithRunBinary` at internal/gate/run_transaction.go line 42. Read
`composedGreenAtKit` at internal/gate/composed_green.go line 15.

Read `pinPath` at internal/gate/phases.go line 412. Read `Command` at
internal/dashboard/dashboard.go line 64.

Each of the six sites spells `rev-parse --absolute-git-dir` today. `pinPath` is
the one site the staged pin-removal spec deletes. Migrate `pinPath` only when it
still exists at the build's base. Report the skip when it does not.

Replace the Git call at each site with `git.AdminDir`. Keep the exact refusal
outcome of every site. The fast lane keeps the error text `gate: git directory
unavailable`. The verdict inspection keeps the state `Unavailable` with the reason
`git directory unavailable`. The run transaction keeps its operational result
with `git directory unavailable`.

The composed-green probe keeps its false answer. The dashboard keeps its exit 1
with the text `cannot resolve git dir`. `pinPath` keeps its fallback path under
`.git`. Drop the hand-written empty-answer guards the reader now owns.

The fixture directory in the `Writes:` line is closure headroom for the pin the
contract ticket's fixture puts on the dashboard file. This ticket edits nothing there.

The fast lane file and the verdict file sit over the structure budget. Each must
end at or under its base line count. The fast lane reds the structure growth check
on a worktree commit that adds a line to one of them.

Write the four new gate tests in internal/gate/admin_readers_test.go. Name them
`TestComposedGreenIsFalseWhenTheReaderFails`,
`TestRunLaneRefusesWhenTheReaderFails`,
`TestInspectSubjectReportsUnavailableWhenTheReaderFails`, and
`TestRunTransactionRefusesWhenTheReaderFails`. Add
`TestDashboardWriteRefusesWhenTheReaderFails` to
internal/dashboard/dashboard_test.go. Drive each over a real repository with the
`fail-git-dir` stub on `PATH`. A non-repository root stops before the reader in
the composed-green probe and the dashboard.

Read `TestGateTransactionRefusesAnUnheldCache` and `plannedEvaluation` in
internal/gate/cache_env_test.go. Drive the transaction test through
`executeSubjectWithRunBinary` with a planned closed subject and `forceRun`.
Assert `ActionExit` 1, `GateExit` 0, and `git directory unavailable` on stderr.

## Acceptance

- [ ] The migrated gate sites and the dashboard site call `git.AdminDir` and spell no Git flag.
- [ ] `ComposedGreen` answers false over a real repository under the `fail-git-dir` stub.
- [ ] The fast lane returns the error `gate: git directory unavailable` under the stub.
- [ ] The verdict inspection reports `Unavailable` with the reason `git directory unavailable` under the stub.
- [ ] The dashboard file write exits 1 and prints `cannot resolve git dir` under the stub.
- [ ] The run transaction returns `ActionExit` 1 with `GateExit` 0 and prints `git directory unavailable` under the stub.
- [ ] The pre-existing `internal/gate` and `internal/dashboard` suites pass with their test logic unchanged.
- [ ] Self-probe: return the typed reader text from the fast lane, and report the lane refusal test red.
