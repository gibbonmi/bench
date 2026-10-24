# Record the lifetime owner decision

Blocked by: none
Writes: docs/adr/0030-one-lifetime-owner-and-process-user-protection.md (new), specs/lp-adr-lifetime-owner/tickets/01-record-the-lifetime-owner-decision.md (new)
Covers: none

## What to build

Add ADR 30, "One lifetime owner and process-user protection", to the ADR folder.
The ADR records the decided state of the subprocess lifetime owner and of durable process-user protection.
It uses the filename, heading, and Status line format of ADR 29.
It states the decided platform scope, and it does not state the status of a platform witness run.
It contains no history, no file paths, and no code snippets.

## Acceptance

- [ ] The ADR folder holds ADR 30 with a heading and a Status line in the format of ADR 29.
- [ ] ADR 30 records one lifetime owner for every production child process, its bounded shutdown, and its caller census.
- [ ] ADR 30 records durable process-user protection, the start gate, and the safe removal protocol.
- [ ] ADR 30 states the platform scope for Linux, Darwin, and Windows, and it does not state the status of a witness run.
- [ ] The commit lane passes, and its prose grade passes for both files.
