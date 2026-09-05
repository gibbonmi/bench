# Migrate the worktree file sites to the reader

Blocked by: migrate-the-worktree-directory-sites.md
Writes: internal/worktree/subshell.go, internal/worktree/clean.go, internal/worktree/worktree.go, internal/worktree/eligibility_test.go, internal/worktree/lifecycle_acquire_test.go, internal/worktree/reauthorize_test.go, internal/worktree/admin_readers_test.go (new), internal/worktree/parallel_census_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: GR18, GR24, GR25

## What to build

Verify the premise first. Read `AdminDir` and `AdminPath` in
internal/git/worktree_admin.go, which the expand ticket adds. Read `LeaseFile` at
internal/worktree/worktree.go line 33, which joins a relative answer onto the
worktree path. Read the index read inside `planExplicitWith` at
internal/worktree/subshell.go line 261. Read `realIndexTree` at
internal/worktree/clean.go line 206. Read `BenchLeaseFilename` in
internal/git/worktree_admin.go.

Read the four test sites that spell a Git flag. They sit at
internal/worktree/eligibility_test.go lines 83 and 113, at
internal/worktree/lifecycle_acquire_test.go line 188, and at
internal/worktree/reauthorize_test.go line 242. The tests are
`TestEligibilityVerdictProjectsWithoutSecondDecision`,
`TestRecoveryPreservesEveryGitVisibleLayerWithoutMovingBranchOrIndex`, and
`TestReauthorizeCommandRollsBackLockRefreshAndCASLoss`.

Change `LeaseFile` to return the answer of `AdminPath` or its typed failure.
Remove the join onto the worktree path. Change the two index reads to call
`AdminPath` with the name `index`. Change the four test sites to call the reader,
and leave every assertion of those tests unchanged.

Three of the files this ticket writes sit over the structure budget. The
subshell, the cleanup, and the worktree file must each end at or under the base
line count. The fast lane reds the structure growth check on a worktree commit
that adds a line to one of them.

The five command registries in the `Writes:` line are closure headroom. They
close the worktree binding, and this ticket edits none of them.

Add the two new tests to internal/worktree/admin_readers_test.go, which the
blocker ticket creates. Name them `TestLeaseFileMatchesIndependentRevParse` and
`TestLeaseFileRefusesRelativeAnswer`. Drive the first over a linked worktree
against an independent `rev-parse` run. Drive the second with the
`relative-git-path` stub on `PATH`.

## Acceptance

- [ ] `LeaseFile` answers the independent `rev-parse --git-path bench-lease` path over a linked worktree.
- [ ] `LeaseFile` returns a typed resolution failure under the `relative-git-path` stub.
- [ ] The two index reads call the file reader and spell no Git flag.
- [ ] The four migrated test sites call the file reader or the directory reader.
- [ ] The three named tests pass with their assertions unchanged.
- [ ] `bench structure --growth <base>` reports no over-budget file that gained a line.
- [ ] Self-probe: restore the join onto the worktree path in `LeaseFile`, and report the relative-answer test red.
