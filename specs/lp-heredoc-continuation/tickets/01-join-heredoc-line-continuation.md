# Join a line continuation before a heredoc body starts

Blocked by: none
Writes: internal/shellcommand/shellcommand.go, internal/shellcommand/lex.go, internal/shellcommand/shellcommand_test.go, specs/lp-heredoc-continuation/tickets/01-join-heredoc-line-continuation.md
Covers: none

## What to build

A backslash-newline after a heredoc operator hides the next command line from every Bash guard.
For example, the shell runs `git reset --hard` in the command `cat <<EOF && \`, `git reset --hard`, `hello`, `EOF`.
`block-dangerous-git.sh` allows that command, and it also allows the shape `cat <<EOF \`, `&& git reset --hard`, `hello`, `EOF`.

`Parse` removes heredoc bodies before it lexes the command.
The heredoc scanner starts the body after the physical line that holds the operator, so it removes the continued command line as body text.
The lexer also keeps the newline after a backslash as a word character, so the continued command word becomes a newline and `git`, and no guard matches it.

Make the heredoc scanner treat a line that ends in a backslash outside quotes as one logical line with the next line.
The body then starts after the logical operator line, as in the shell.
Make the lexer remove a backslash-newline outside quotes, as the shell does, so the continued words join their logical line.
A backslash-newline inside single quotes stays text.
`Parse` keeps its contract, and it still removes a real heredoc body.
The test seam is `Parse` and `ProjectCommandWords`.

## Acceptance

- [ ] A `Parse` test with each continued-operator shape that hides `git reset --hard` is red before the fix and green after it, and `block-dangerous-git.sh` blocks both shapes at exit 2.
- [ ] The existing `internal/shellcommand` tests and the guard tests that consume `Parse` stay green.
- [ ] The `bench commit` lane passes.
