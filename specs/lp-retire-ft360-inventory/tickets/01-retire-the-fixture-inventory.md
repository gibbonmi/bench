# Retire the shared-test-fixtures inventory with its spec

Blocked by: none
Writes: specs/shared-test-fixtures/inventory.json
Covers: none

## What to build

The `shared-test-fixtures` spec landed, and `bench spec retire` removes its whole folder. The folder holds `inventory.json`, the caller inventory of the spec stage. The commit refuses that deletion as a production path, because the landed outcome holds no current delivery binding.

This ticket covers that one deletion, so the retirement commit lands through the light path. No code reads the inventory after the landing, because ticket 13's registered checks now enforce each owner.

## Acceptance

- [ ] The `specs/shared-test-fixtures` folder, with its inventory, is gone from the tree.
- [ ] The whole-project gate passes on the landing.
