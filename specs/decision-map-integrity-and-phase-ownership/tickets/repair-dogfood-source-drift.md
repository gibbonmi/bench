# Repair dogfood source drift

Blocked by: Separate shaping from spec authoring

## What to build

Repair the compiled decision map facts that a fresh map-to-spec session found
stale after the phase-ownership migration.

## Acceptance

- [x] The map no longer claims that the current workflow requires or retains a
  transitional Handoff.
- [x] The command-source descriptions match the situational shaping and
  three-source spec-authoring contracts now in the tree.
- [x] The `internal/maps` Sources identify the current schema, validation,
  tree-validation, and query owners without dropping the original decision
  provenance.
- [x] The ready compiled map passes focused tree validation, and a fresh reader
  can re-verify every repository Path source.
