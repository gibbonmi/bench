# Mark each live command substitution in the shared projection

Blocked by: none
Writes: internal/shellcommand, tests/canary/cancel-signal-registrations/own-signal-set
Covers: R01, R02, R03, R04, R05, R06, R07, R55

## What to build

The shared `shellcommand` projection marks each live command substitution once, so that a guard can see a substitution that quote folding hides.

`Token` gains one boolean field, `Substitution`. The lexer records a mark beside each token, and `Parse` copies the mark into the field.
The lexer marks a Word token when the token holds a `$(` or a backtick outside single quotes and not after a backslash.
Inside double quotes, `\$` and a backslash before a backtick are escapes for the mark only. The folded text stays as it is now.
The lexer marks a ControlOperator token whose text starts with `(` when an unquoted, unescaped `$` comes directly before it with no space.
A `$((` arithmetic opener takes the same mark as `$(`.
The fallback split of an unlexed command gives every token an unmarked field, so the `Unlexed` path keeps its current facts.

The mark adds information only. Token text, token kinds, spans, and the `Unlexed` flag stay exactly as they are now, so no Git, Bench, census, or worktree verdict changes.

`internal/shellcommand/shellcommand.go` has 397 lines against a 400-line budget. Move the lexer (`lex` and its helpers) into a new file in the same package, such as `internal/shellcommand/lex.go`, and put the mark logic there.
Keep the import block and the `TokenKind` declaration at the top of `shellcommand.go` unchanged. The canary `tests/canary/cancel-signal-registrations/own-signal-set` anchors on those lines and pins `shellcommand.go:18`. If the move changes those lines, update the canary in this ticket.

Add field names to the unkeyed `Token` literals in `internal/shellcommand/shellcommand_test.go` (lines 11, 27, 49, 61, 70) and `internal/shellcommand/unlexed_test.go` (line 32). Change no expected text, kind, or span, so `TestParsePublishesTokenClassesAndSimpleCommands` keeps its witness for R07.

Put the new mark table in `internal/shellcommand/substitution_test.go` as `TestParseMarksLiveSubstitution`, with one case for each of R01 to R06. The expectations come from POSIX shell quoting, not from the code under test, so record one omission red under the `AGENTS.md` exception: drop the double-quote branch of the mark and record the R01 red, the restore, and the green.

## Acceptance

- [ ] `Parse("rg \"$(x)\" .")` yields the Word `$(x)` with `Substitution` true (R01).
- [ ] ``Parse("rg \"`x`\" .")`` and ``Parse("rg `x` .")`` each yield a marked Word (R02).
- [ ] `Parse("rg $(x) .")` yields a marked `(` ControlOperator directly after the Word `$` (R03).
- [ ] `Parse("rg 'lit$(x)' .")` and ``Parse("rg 'a`b' .")`` yield only unmarked tokens (R04).
- [ ] `Parse("rg \\$\\(x\\) .")`, ``Parse("rg \\`x\\` .")``, and `Parse("rg \"\\$(x)\" .")` yield only unmarked tokens (R05).
- [ ] `Parse("(rg x .)")` and `Parse("a $ (b)")` yield an unmarked `(` operator (R06).
- [ ] `TestParsePublishesTokenClassesAndSimpleCommands` passes with its existing expected text, kinds, and spans (R07).
- [ ] `bench test --check system --run TestGuardProjectionDifferentialFamilies` passes with the golden byte-equal (R55).
- [ ] `bench test --package ./internal/gitguard --package ./internal/benchguard --package ./internal/census` passes unchanged.
- [ ] `shellcommand.go` stays within its 400-line budget, and the cancel-signal canary stays green.
