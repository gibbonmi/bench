# Grant the publisher override between two others

Blocked by: none
Writes: internal/worktree/land_override_test.go
Covers: none

## What to build

`TestOverrideIsConsumedAtPublish` grants three overrides with the publisher's override in the middle. The publication consumes the publisher's override and keeps the overrides before and after it.

## Acceptance

- [ ] `TestOverrideIsConsumedAtPublish` grants the publisher's override between two other overrides and asserts that the publication keeps the other two overrides in order.
