# Contract live specs to folder only

Blocked by: Migrate spec lifecycle and command contracts; Migrate ambient and AXI spec consumers

## What to build

Remove transitional live flat-form support. Refuse flat-only, flat/folder
collisions, and folder residue without `spec.md` with the specified precedence
and migration diagnostics. Make conformance validate folder coverage maps and
reject stray flat specs while leaving retired flat history intact.

## Acceptance

- [ ] Every specified flat/collision/missing-file refusal is green.
- [ ] The combined collision state follows the specified precedence.
- [ ] Conformance validates `specs/*/spec.md`.
- [ ] Conformance rejects a stray live `specs/*.md`.
- [ ] Retired flat history remains green.
- [ ] The project gate is green.
