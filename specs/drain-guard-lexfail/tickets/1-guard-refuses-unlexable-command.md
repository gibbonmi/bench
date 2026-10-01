# Refuse a guarded command the lexer cannot parse

Blocked by: none
Writes: internal/shellcommand/shellcommand.go, internal/shellcommand/unlexed.go (new), internal/shellcommand/shellcommand_test.go, internal/shellcommand/unlexed_test.go (new), internal/gitguard/gitguard.go, internal/gitguard/scan.go, internal/gitguard/unlexed_test.go (new), internal/benchguard/benchguard.go, internal/benchguard/unlexed_test.go (new)
Covers: none

## What to build

The shell lexer in `internal/shellcommand` stops at an unterminated quote or a trailing backslash. `Parse` then uses a fallback split, and no caller knows that the split is a guess. The fallback split keeps a shell operator inside a word, so `git status;git push --force` stays one word and the guard does not see the `push` subcommand. A command with a second line that opens a quote therefore passes the destructive-git guard, and the shell still runs the first line.

The reviewer decided that a guard fails closed here.

- The git guard refuses a command that the lexer cannot parse and that holds a `git` word.
- The Bench follow-on guard refuses a command that the lexer cannot parse and that holds a Bench word.

The shell rim of the destructive-git hook takes the same posture for a command that it cannot read.

Make these changes:

- `Parse` records the lexer failure on the returned `Stream`.
- `shellcommand` gives one query that names the first simple command of a failed parse that holds a guarded word. The fallback keeps quote and escape characters, so the query reads each word without them.
- `shellcommand` owns the one repair sentence for a command the lexer cannot parse.
- The fallback split puts each run of shell operator characters in its own token, so the tokens show each separator.
- The git guard refuses with its own deny label and that repair sentence.
- The Bench follow-on guard refuses with that repair sentence and the segment that holds the Bench word.
- A command that the lexer parses keeps its current verdict in both guards.

The fallback code moves to a new file, because `shellcommand.go` is over its structure budget.

## Acceptance

- [ ] The git guard refuses `git status;git push --force` followed by a line that opens a quote, and the hook script exits 2 for it.
- [ ] The git guard refuses a command that the lexer cannot parse and that holds `git` with no destructive subcommand.
- [ ] The git guard allows a command that the lexer cannot parse and that holds no `git` word.
- [ ] The Bench follow-on guard refuses a command that the lexer cannot parse and that holds a Bench word.
- [ ] The Bench follow-on guard allows a command that the lexer cannot parse and that holds no Bench word.
- [ ] `Parse` marks the lexer failure for an unterminated quote and does not mark it for a clean command.
- [ ] Each new test is red before the production edit, and a `bench probe` omission of the git guard refusal turns a test red.
- [ ] `go test ./internal/shellcommand/ ./internal/gitguard/ ./internal/benchguard/ ./internal/census/` passes.
