# Route the tickets-only close to the landing verb

Blocked by: none
Writes: internal/status/status.go, internal/status/status_counters_test.go, internal/status/status_producible_test.go, tests/canary/docs-currency-token-diet/signal-vocabulary-drift
Covers: LP6, LP7, LP8

## What to build

Verify the premise first: `closeTicketsAction` in internal/status/status.go
prints `bench commit --spec`, and `bench commit --help` accepts no `--spec`
flag. Then change the definition's command to `bench worktree land --spec` and
keep the one-word argument shape, so the row prints
`bench worktree land --spec <slug>`. Update the comment beside
`appendTicketsOnly` to name the landing's `--spec` close. Update the two pinned
test literals to the new string.

## Acceptance

- [ ] The tickets-only case in `TestAllProducibleBoardActionsAreInvocableOrEmpty` expects and receives `bench worktree land --spec <slug>`.
- [ ] `TestTicketsOnlyResidueRowCountsAndRanksBelowItsBand` expects and receives the new route for two folders.
- [ ] `TestActionDefinitionsRenderAndParseTheSameCommand` passes with the new definition.
- [ ] Self-probe: drop the argument to a bare command, and report both literal assertions red as the observed result.
