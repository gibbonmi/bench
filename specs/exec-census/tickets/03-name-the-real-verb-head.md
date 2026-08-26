# Name the real verb head in a census record

Blocked by: 02-record-one-raw-call-per-bash-call.md
Writes: internal/census/census.go, internal/census/census_test.go

## What to build

The head must name the command the agent ran, because the drain reads the head
and proposes the missing Bench verb. This ticket refines the head that
`census.Record` writes. The signature and the match do not change.

`shellcommand.Parse` removes a heredoc body before it returns the tokens. A
pool path inside that body therefore never reaches a command's words. When no
simple command's words name a pool path, but the raw text does, the head is the
first simple command's head. A scripted edit then counts with a real head.

`shellcommand.ResolveRoutinePrefix` steps over `env`, `timeout`, `xargs`, and
leading assignments. The head is the word that prefix resolves to. When that
word is `git`, the head also carries the first subcommand word, for example
`git add`. A flag before the subcommand does not become the subcommand.

Two edges hold. A text `cd <pool>/<owner>-<id> && make` records the head `cd`.
An unterminated quote takes the tokenizer's fallback split, and the text still
matches.

## Acceptance

- [ ] A heredoc text whose body names a pool path records the head `python3`. (EC05)
- [ ] A text `git add <pool>/<owner>-<id>/x` records the head `git add`. (EC07)
- [ ] A text `FOO=1 timeout 5 sed -i x <pool>/<owner>-<id>/y` records the head `sed`. (EC08)
- [ ] A text `cd <pool>/<owner>-<id> && make` records the head `cd`.
