# Review pickup: worktree-native-forms

Frozen base `a8724eb4341db33993db93ffbe299d4004f5d345`, reviewed tip
`b768fcc99322f6cbce92571a65387201fd396aea`. Three Opus axes at medium effort.

## Standards

Count: 8 raw, 6 repair targets before collapse. Worst: S2.

- S1, `auto-fix`. Rule: one source per fact (`AGENTS.md`, Code standard). The
  `--from contains control characters` guard is spelled twice, at
  `internal/worktree/merge.go:183` and `internal/worktree/worktree.go:625`.
  Ticket `repair-worktree-single-resolution-and-replay.md`.
- S2, `auto-fix`. Same rule. `resolveBuildTarget` in `internal/worktree/build.go`
  reads the ledger and selects the assignment a second time after
  `resolveWorktree` did the same and kept the path alone. No other verb pairs
  the two. Same ticket.
- J2, `auto-fix`. `craft-comments`, Aging. `internal/worktree/path.go:143` says
  "both target-taking verbs"; four verbs share the printer. Same ticket.
- S3, `ask-user`. Same rule. `internal/preflight/gather.go` `canonicalRoot` is a
  third copy of the canonical-path form beside `internal/worktree/subshell.go`
  `canonicalPath` and `internal/canary/mutation.go` `resolvePath`.
  `internal/worktree/land.go` imports `internal/preflight`, so a direct call
  cycles. The collapse needs a below-both leaf package, which is a new seam
  outside this spec's fence. The reviewer decides the home.
- J3, `no-op`. WF42's kept-grammar check reads the same const the help joins.
  The house idiom already holds for `usage.WorktreeShow`, and WF31 pins the
  literal grammar independently.
- J1, J4, J5, `no-op`.

## Spec

Count: 2 raw, 0 repair targets. All 44 rows closed with named tests.

- Note 1, `no-op`. The spec says "one system-phase producer"; the gate exports
  four names, each single-sourced with two call sites.
- Note 2, `no-op`. The `bench worktree --help` row description gained `build`,
  because the row enumerates the grammars the help prints.

## Coverage

Count: 3 raw, 2 repair targets after refutation. Worst: C3.

- C1, `no-op`. Claim: with `BENCH_KIT` unset the kit-only guard is vacuous and
  the child loses `BENCH_KIT`. Refutation: `selectedRunEnvironment` appends
  `BENCH_KIT=<selection.SourceRoot>` at `internal/testreport/command.go:193`.
  `gate.kitRoot` falls back to the root exactly as `testBenchSource` does, so
  the check mirrors the gate's own kit derivation.
- C2, `auto-fix`. A future conformance registry check named `system` is
  shadowed silently at `internal/testreport/command.go:93`. One reservation
  test in `internal/testreport/check_test.go` asserts the registry never
  names `gate.SystemPhaseName`. Ticket `repair-system-check-name-reservation.md`.
- C3, `auto-fix`. `CreateCommand` resolves `--from` before `createAt` runs the
  replay lookup. So a replay with the same `--request` after the sibling
  dirties refuses instead of returning the record. The Edge inventory promises
  the replay returns the existing record. New row WF45. Ticket
  `repair-worktree-single-resolution-and-replay.md`.
