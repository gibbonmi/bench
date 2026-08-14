# Order the profile repair after its census row

Blocked by: none
Writes: specs/spec-authoring-and-light-path/tickets/repair-profile-loop-routing.md

## What to build

Set the profile-loop repair's blocker to `repair-census-map-coverage.md`, which authors WF27 before the profile repair claims that row.

## Acceptance

- [ ] a fresh ticket graph orders the WF27 author before its consumer
