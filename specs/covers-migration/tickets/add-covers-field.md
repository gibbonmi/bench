# Add the Covers field to the citation ticket

Blocked by: none
Writes: specs/citation-phase-package-scope/tickets/bind-citations-to-phase-package-scope.md
Covers: none

## What to build

The ticket-grammar landing enforces a `Covers:` field on every ticket file.
The staged citation ticket predates the rule. Add the field with the twelve
PS rows its acceptance list already names.

## Acceptance

- [ ] The citation ticket carries `Covers:` after `Writes:`.
- [ ] `bench gate-prose` passes on both edited files.
