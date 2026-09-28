# Align filter help

Blocked by: none
Writes: internal/consumers/command.go, internal/consumers/query_test.go, cmd/bench/help_inventory_test.go
Covers: none

## What to build

Use the same compact filter grammar in the consumers and outline help.
Keep each command's suffix as the source for its usage and top-level help.

## Acceptance

- [ ] Both commands show `[--production|--test]` in their usage lines.
- [ ] The top-level help shows that spelling for both commands.
