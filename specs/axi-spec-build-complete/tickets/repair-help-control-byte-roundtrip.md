# Preserve fixed control bytes through rendered help

Blocked by: none
Ownership fence: `internal/axi/action.go`, `internal/axi/action_test.go`
Integration surfaces: fixed argv serialization and help rendering→`internal/axi/action.go`; existing shell quoting→`internal/sanitize/sanitize.go` exercised unchanged by HR1 and HR2; existing TOON cell encoding→`internal/toon/toon.go` exercised unchanged by HR1 and HR2
Contracts: a fixed argument's exact bytes cross `internal/axi/action.go` command construction→the existing TOON table encoder, asserted by HR1 and HR2 through the rendered cell and an executing shell; FR6's permitted-byte decision crosses construction→rendering inside `internal/axi/action.go`, asserted by HR3
Closure: HR1/tab-byte-roundtrip, HR2/carriage-return-byte-roundtrip, HR3/tab-construction-preserved, HR3/carriage-return-construction-preserved, HR4/newline-refusal, HR5/esc-refusal

## What to build

Close the accepted round-8 Spec and Coverage findings P1/C1 without reopening
FR6. Fixed tab and carriage-return bytes remain valid action arguments, and a
command copied from `RenderHelp` reconstructs those exact bytes when the shell
executes it. TOON's visible escape sequences must not silently change either
byte into a backslash-plus-letter argument.

## Acceptance

- [ ] [HR1] (covers local) (P1, C1) a fixed argument containing a tab survives `RenderHelp` and shell execution as the original tab byte.
- [ ] [HR2] (covers local) (P1, C1) a fixed argument containing a carriage return survives `RenderHelp` and shell execution as the original carriage-return byte.
- [ ] [HR3] (covers local) (P1, C1) action construction continues accepting fixed tab and carriage-return values.
- [ ] [HR4] (covers local) (P1, C1) action construction continues refusing a fixed newline value.
- [ ] [HR5] (covers local) (P1, C1) action construction continues refusing a fixed ESC value.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| HR1/tab-byte-roundtrip | return the literal TOON-visible `\\t` spelling inside the copied shell argument | the rendered-help shell round-trip test | render an action carrying a tab, execute the rendered command through the shell with an argv-byte probe, and require the original tab byte rather than backslash plus `t` |
| HR2/carriage-return-byte-roundtrip | return the literal TOON-visible `\\r` spelling inside the copied shell argument | the rendered-help shell round-trip test | render an action carrying a carriage return, execute the rendered command through the shell with an argv-byte probe, and require the original carriage-return byte rather than backslash plus `r` |
| HR3/tab-construction-preserved | reject the FR6-permitted fixed tab value during action construction | the pure action construction test | construct a fixed tab value and require construction to succeed |
| HR3/carriage-return-construction-preserved | reject the FR6-permitted fixed carriage-return value during action construction | the pure action construction test | construct a fixed carriage-return value and require construction to succeed |
| HR4/newline-refusal | admit a fixed newline value during action construction | the pure action construction test | construct the fixed newline value and require refusal |
| HR5/esc-refusal | admit a fixed ESC value during action construction | the pure action construction test | construct the fixed ESC value and require refusal |
