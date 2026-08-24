# Split the lines tests

Blocked by: none
Writes: internal/lines/lines_test.go

Rewritten 2026-08-24. The original target, `internal/worktree/eligibility_test.go`, has held 209 lines since before the spec's base. The reviewer widened the scope by one file so row R13 can hold at 55.

## What to build

Split `internal/lines/lines_test.go` (612 lines) into three files plus one fixture file. Pure moves only.

- The parse file takes `TestTierValue`, `TestTierValueEdgeCasesReachTheMatrix`, `TestTierValueEdgeCases`, `TestModelFromEnvelope`, `TestParseBinding`, `TestParseBindingRejectsForeignHarnessKeys`, `TestRetiredSchemaBindsNothing`, and `TestUnreadableSourceIsNotAbsent`.
- The verdict-resolution file takes `TestResolveModelVerdict`, `TestResolveModelVerdictResolvesEveryCell`, `TestResolveModelVerdictNamesTheBindingPath`, `TestResolveModelVerdictRejectsMalformedColumn`, and `TestCellFault`.
- The agent-line file takes the `TestAgentLineVerdict*` group and `TestSubagentTypeNeverImpersonatesAFork`.
- The fixture file takes `bound`, `envelope`, `agentEnvelope`, `forkEnvelope`, and `contains`.

The helpers stay in package `lines`, because they build package-private values. The delegate may regroup within the package when a measured count differs, under the coverage rows.

## Acceptance

- [ ] R20: `bench structure` no longer lists `internal/lines/lines_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R08: `go test -list '.*' ./internal/lines/` emits the same test-name set at base and tip.
- [ ] R11: each moved helper has exactly one definition in the package.
- [ ] R18: `bench gate` exits zero before the commit.
- [ ] R13: the final `bench structure` total is at most 55 issues once every ticket has landed.
