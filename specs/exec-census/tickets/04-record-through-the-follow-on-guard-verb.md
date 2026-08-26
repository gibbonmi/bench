# Record a raw call through the follow-on guard verb

Blocked by: 02-record-one-raw-call-per-bash-call.md
Writes: cmd/bench/main.go, cmd/bench/main_test.go

## What to build

`bench guard-bench-follow-on` reads every Bash tool call's envelope today. It
gains one duty. Before it returns the verdict, it calls `census.Record` with
the decoded command text, the repository root, the resolved Bench home, and
the current time.

The verb resolves the home at the command boundary. It tests the text for the
`<home>/worktrees/` substring before it resolves any root, so an ordinary
call spawns no git process. When the substring is present, the verb resolves
the root with `git.Root()` and passes it down. `poolkey.Key` canonicalizes
it, so a call from inside a linked worktree keys the primary repository.

The record never changes the verdict or the exit code. Its own error is
silent: the verb writes no diagnostic and still exits 0. The record happens
before the verdict, so the verb's own control flow cannot skip it.

Two existing paths stay unchanged. An envelope with no command field takes the
warning path and records nothing. A command with a NUL byte takes the refusal
path and records nothing.

This ticket adds no hook script and no new verb. Both harness wirings, the
guards manifest, the conformance registry, and the adopt payload stay as they
are.

## Acceptance

- [ ] A raw-call envelope with a plain command makes the verb exit 0 and write one record. (EC09)
- [ ] An unwritable home makes the verb exit 0 with empty stderr. (EC10)
- [ ] An envelope with no command field records nothing and keeps the warning path.
- [ ] An envelope whose text lacks the `<home>/worktrees/` substring records nothing and resolves no root.
