# Split the worktree eligibility tests

Blocked by: 06-split-worktree-land-files.md
Writes: internal/worktree/eligibility_test.go

## What to build

Split `internal/worktree/eligibility_test.go` (795 lines) into two files. The explicit-decision file takes `newUnsafeTargetOwnedWorktree` and `gatherExplicitFactsForTest`. The automatic-matrix file also takes the verdict-projection group: `TestEligibilityVerdictProjectsWithoutSecondDecision` and `assertVerdictMatchesPlan`. Pure moves only.

The helpers stay in package `worktree`, because they build the private `explicitFacts` type. The blocker exists only because this ticket shares the `internal/worktree/` fence with ticket 06. No value contract crosses the two tickets.

## Acceptance

- [ ] R06: `bench structure` no longer lists `internal/worktree/eligibility_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R08: `go test -list '.*' ./internal/worktree/` emits the same test-name set at base and tip.
- [ ] R11: each moved helper has exactly one definition in the package.
- [ ] R18: `bench gate` exits zero before the commit.
