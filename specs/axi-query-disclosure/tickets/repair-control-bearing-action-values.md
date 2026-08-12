# Repair control-bearing action values

Blocked by: repair-worktree-compatibility-closure.md
Writes: `internal/axi/action.go`, `internal/axi/action_test.go`, `internal/maps/maps_test.go`, `internal/worktree/list_actions_test.go`, `internal/coverage/coverage_test.go`

## What to build

Let known action arguments carrying exactly tab U+0009, line feed U+000A, or carriage return U+000D survive through the existing shell-quoting and TOON escaping pipeline. Continue to refuse every other `unicode.IsControl` rune: the remaining C0 controls, DEL U+007F, and C1 U+0080–U+009F. When a surface-owned disclosure value is unsupported, preserve the already-computed primary response and exit and append the honest empty help block instead of replacing the response with a renderer error.

## Acceptance

- [ ] [CV1] (covers QD1, QD6) `KnownArgument` accepts exactly U+0009, U+000A, and U+000D from the control-rune class and refuses the remaining C0 controls, DEL U+007F, and C1 U+0080–U+009F. The oracle decodes the rendered `cmd` cell with the official `toon-go` decoder, invokes the decoded command through a POSIX shell argv probe, and byte-compares the recovered argument with the source. This acceptance explicitly rewrites RP1's original blanket control-bearing counterexample.
- [ ] [CV2] (covers QD1, QD6) public coverage and orphaned-worktree queries with each supported control-bearing path preserve their primary response and exit while appending the exact round-tripped action. An unsupported known argv value preserves the same primary response and exit and appends honest `help[0]{cmd,why}:` rather than replacing the response with an error.
- [ ] [CV3] (covers QD1, QD6) an invalid-map diagnostic whose `why` cell contains an unsupported control proves the same fallback posture: its primary diagnostic and exit are unchanged and its appended help block is honestly empty.
