# 1. Split the named-check owner out of the command file

Blocked by: none
Writes: internal/testreport/
Covers: none

## What to build

Chunk: TP-C1.

Move the named-check owner out of `internal/testreport/command.go` into a new file in the same package.
The owner is the check name set, the unknown-check refusal, the conformance environment, and the three named-check runners.

Change no behavior and no output byte.
After the move, `command.go` is under its line budget, so the later tickets can edit the request parser.

## Acceptance

- [ ] `command.go` holds fewer lines than its structure budget.
- [ ] `bench test --package ./internal/testreport` is green with no test expectation changed.
- [ ] `bench structure --growth bcc5543f` names no file of this ticket.
