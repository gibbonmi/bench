# Guard Bench commands from shell follow-ons

Frozen base: `f05c5e7a4a4e03b7539455822d045aa2cd0a9e92`

Reviewed tip: `dc8c0880adee49315506a8c884879222d5470446`

## Standards

Finding count: 0. Worst issue: none.

The review found no actionable Standards issue.

## Spec

Finding count: 1. Worst issue: P1 wrapper outer-follow-on bypass.

- P1, `auto-fix`: The spec requires one wrapper depth and refusal of outer shell
  syntax. The classifier allows `bash -lc 'bench help' | touch <marker>` because
  wrapper recursion checks only the quoted child. Repair the wrapper result so the
  outer stream still decides refusal. Citation: spec FOG18 and
  `internal/benchguard/benchguard.go` wrapper scan.

## Coverage

Finding count: 4. Worst issue: P1 Bench detection gaps.

- P1, `auto-fix`: A leading file-descriptor redirection leaves its digit in command
  position. `2>/dev/null bench help | touch <marker>` therefore hides Bench.
  Citation: spec redirection inventory and `internal/benchguard/benchguard.go`
  command-word projection.
- P1, `auto-fix`: Routine-prefix option parsing can skip the Bench executable.
  `xargs bench help | touch <marker>` is the direct case. Cover the supported
  `env`, `command`, `nohup`, `timeout`, and `xargs` forms. Citation: spec FOG27
  and the routine-prefix decision.
- P2, `auto-fix`: FOG23 requires the rendered `bench guards` query seam. Existing
  coverage calls `Rows` directly, so a rendering or routing regression can pass.
  Citation: spec FOG23 and `internal/guards/guards_test.go`.
- P2, `auto-fix`: FOG25 requires `bench gate --brief` to remain a usage error.
  Existing generic argument coverage does not name this flag. Citation: spec FOG25
  and the gate command tests.
