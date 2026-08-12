# Repair worktree review closure

Blocked by: repair-disclosure-value-proof-seam.md
Writes: `internal/usage/parse.go`, `internal/usage/parse_test.go`, `internal/worktree/list.go`, `internal/worktree/list_actions_test.go`, `internal/worktree/testdata`

## What to build

Close the three worktree review gaps without reopening the accepted grammar decision: preserve the real pre-change empty-token refusal bytes through a grammar-scoped option, consume producer order directly in the public oracle, and add a public non-actionable-terminal compatibility pair.

## Acceptance

- [ ] [WR1] (covers QD2) the checked old empty-token response is corrected to `usage: bench worktree list (unknown argument: )\n`, and current `ListCommand([]string{""})` reproduces that stdout, empty stderr, and exit 2 through `usage.Parse`. Only sole `--help`, `-h`, and `help` differ between old and new rows.
- [ ] [WR2] (covers QD2) the empty-token rendering exception is declared on `worktreeListGrammar`; every existing grammar retains the shared parser's quoted-empty refusal. Mutating the worktree option or enabling it on an unrelated grammar is red.
- [ ] [WR3] (covers QD1) the active/orphan public oracle consumes the assignment producer's creation and serialized order directly and never compares or sorts IDs independently. Re-sorting production rows makes the exact public response test red.
- [ ] [WR4] (covers QD6) a checked-in public fixture for a non-actionable terminal assignment or present foreign row pins the pre-disclosure primary bytes, empty stderr, and exit 0; the candidate adds exactly `help[0]{cmd,why}:\n`.
- [ ] [WR5] (covers QD6) changing terminal primary bytes, exit, or help cardinality makes the terminal fixture test red without relying on the direct `actionsForRows` test.
