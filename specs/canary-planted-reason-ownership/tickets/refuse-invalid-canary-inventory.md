# Refuse invalid canary inventory

Blocked by: none
Writes: `internal/canary`

## What to build

Keep every producer-derived binding partition fail-closed and replace the empty-tree execution-proof diagnostic with inventory-only truth. Exercise the selection decision directly so the ticket lands without depending on command vocabulary, dispatch removal, or fixture-file metadata repair.

## Acceptance

- [ ] (covers CI3) Empty, duplicate, unbound, unsafe, and control-bearing inventories refuse deterministically, and the empty refusal claims only that no accepted binding exists.
