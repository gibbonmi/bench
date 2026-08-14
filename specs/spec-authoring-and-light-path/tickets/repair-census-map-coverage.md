# Cover the repair fixture census

Blocked by: repair-ticket-self-containment.md
Writes: specs/spec-authoring-and-light-path/spec.md, specs/spec-authoring-and-light-path/tickets/repair-profile-loop-routing.md

## What to build

Add one acceptance-map row for the fixture-producing profile-loop repair ticket's independently authored canary count update. Tag that ticket's loop-routing acceptance to WF10 and its count acceptance to the new row, keeping one owner per row.

## Acceptance

- [ ] the profile-loop repair acceptance explicitly covers WF10
- [ ] a new one-owner row covers its 202-to-204 canary census update
- [ ] `bench coverage --check` passes
