# Migrate spec lifecycle and command contracts

Blocked by: Expand spec resolution with folder form

## What to build

Move the spec lifecycle and direct command contracts to folder specs: coverage,
commit status flips, implemented and retire behavior, history, and gate-action
proofs. Preserve retired flat-path history while making folder deletion and
interrupted-retire behavior observable.

## Acceptance

- [x] Coverage, commit, implemented, retire, and history contracts use folder specs.
- [x] Retire removes the pickup and complete folder, including tickets.
- [x] Recoverable interrupted retire states rerun cleanly.
- [x] Terminal folder residue refuses with the named hand-clean instruction.
- [x] Both historic flat and new folder deletions appear in history.
