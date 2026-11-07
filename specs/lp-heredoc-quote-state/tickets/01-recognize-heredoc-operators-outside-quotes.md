# Recognize heredoc operators only outside quotes

Blocked by: none
Writes: internal/shellcommand/shellcommand.go, internal/shellcommand/shellcommand_test.go, internal/shellcommand/lex.go
Covers: none

## What to build

The shell parser removes heredoc bodies before the lexer reads the command. Make it recognize a `<<` operator only outside single quotes, outside double quotes, and outside a backslash escape. Track the quote state across lines, so a quote that opens on one line and closes on a later line keeps its lines in the command. A heredoc body starts only after the line on which every quote closes.

Use the lexer's quote rule for this scan, so the parser keeps one source for where a quote ends.

This change makes each Bash guard stricter: the parser removes fewer lines, so every command that a guard refuses now stays refused.

## Acceptance

- [ ] A `"<<X"` or `'<<X'` argument followed by a `git reset --hard` line parses with `git reset --hard` as a simple command.
- [ ] A `<<X` inside a double-quoted string that spans lines does not remove the lines after the string.
- [ ] A real `<<X` operator on a line that ends inside an open quote does not remove the next lines until that quote closes.
- [ ] A real heredoc outside quotes still has its body removed, and its operator stays in the token stream.
