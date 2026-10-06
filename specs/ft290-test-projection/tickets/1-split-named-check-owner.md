# 1. Split the named-check owner out of the command file

Blocked by: none
Writes: internal/testreport/
Covers: none

## What to build

Chunk: TP-C1a.

Move the named-check owner out of `internal/testreport/command.go` into a new file in the same package.
The owner is the check name set (`namedChecks`, `isNamedCheck`), the unknown-check refusal (`unknownCheck`, `namedCheckInventory`), the three named-check runners (`runNamedCheck`, `runProseCheck`, `runSystemCheck`), and the `proseCheckName` constant.
The conformance environment already lives in `environment.go`, so it does not move.

Change no behavior and no output byte.
`command.go` is under its line budget today, but tickets 3 and 6 to 9 add more lines than its room. The move makes that room.

## Acceptance

- [ ] `command.go` holds none of the moved symbols, and its line count drops by the moved lines.
- [ ] A named-check result prints the same bytes before and after the move.
- [ ] The unknown-check refusal prints the same bytes before and after the move.
