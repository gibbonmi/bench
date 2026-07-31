# Enforce the binding matrix and command tokens

Blocked by: Migrate line binding through the harness matrix

## What to build

Root conformance validates the complete declared harness matrix against the
profile and rejects unbound model tokens in every discovered portable command.

## Acceptance

- [ ] Missing, malformed, incomplete, and prose-drifted matrix cells emit targeted diagnostics.
- [ ] Every discovered regular command file is scanned from the directory-owned inventory.
- [ ] One unbound literal emits its file-and-token diagnostic from an otherwise clean fixture.
- [ ] Non-regular command entries are rejected before any read.
