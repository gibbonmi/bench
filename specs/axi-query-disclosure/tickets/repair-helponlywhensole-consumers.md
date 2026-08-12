# Repair HelpOnlyWhenSole consumers

Blocked by: none
Writes: `internal/worktree/list.go`, `internal/usage/parse.go`

## What to build

Close both consequences of `HelpOnlyWhenSole`'s real effect (review findings
R3 and R4, one root cause): the option disables all flag recognition and the
`--` terminator for its grammar, not only mid-args help. Remove the dead
consumer — the unreachable `parsed.EndedFlags` branch in `ListCommand`
(`list.go:29`; `Result.EndedFlags` is assigned only inside the block the flag
skips) — and document the full effect in the `Grammar` doc comment. No
behavior change anywhere.

## Acceptance

- [ ] [HS1] (covers local) the `parsed.EndedFlags` branch is removed from
  `ListCommand`; the QD2 paired argv fixtures (including bare `--` and
  `-- x`) and every QD6 worktree pair stay byte-identical.
- [ ] [HS2] (covers local) the `Grammar` doc comment states that
  `HelpOnlyWhenSole` recognizes help only as the sole argument AND disables
  flag recognition and the `--` terminator for the grammar; comment-only edit
  in `internal/usage/parse.go`.
