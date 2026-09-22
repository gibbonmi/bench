# Refresh a ready map before the survey handoff

Blocked by: none
Writes: .agents/commands/bench-deepen.md, specs/ft219-ready-map-refresh/tickets/refresh-ready-map.md (new), reviews/ft219-ready-map-refresh.md (new)
Covers: none

## What to build

Add the approved FT219 branch after source verification and before candidate selection.
A survey that revisits a ready map updates that map to the current remaining work.
The phase preserves closed decisions and exits through the normal handoff.
An unresolved frontier keeps the existing report and grill route.

## Acceptance

- [ ] The phase verifies an existing map's frontier against current source before candidate selection.
- [ ] An empty frontier refreshes the existing map around remaining work and preserves closed decisions.
- [ ] Duplicate roadmap decision claims become pointers to the existing map.
- [ ] The phase runs `bench maps` and rewrites the session handoff after the refresh.
- [ ] The ready branch creates no second map and does not grill settled predicates again.
- [ ] A frontier with open questions continues through the existing candidate report and grill.

## Verification

Run the prose lane and the workflow guidance check.
Read every steered surface and obtain fresh-session adoption against the committed owner bytes.
Run separate native Standards, Spec, and Coverage reviews before the landing gate.

## Line and allowance

The retained coordinator uses Astra at ultra effort, with three implementation attempts and two post-review repair cycles.
No repair cycle is consumed at entry. The expected repair count is zero, with confidence eight.
