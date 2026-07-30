# Project the active-map model

Blocked by: Validate the decision graph and readiness

## What to build

Every tracked active and compiled map migrates without decision loss, then
`bench maps`, its count mode, and ambient status switch atomically to the new
model and the legacy Handoff parser leaves `internal/maps`.

## Acceptance

- [ ] Every tracked active map is honestly shaping or ready, every compiled map
  is ready, and no migrated map loses an answer, exclusion, or source fact.
- [ ] Default output is `maps[N]{map,title,type,state,blockers}` and shows every
  unresolved decision ticket as frontier, blocked, or deferred with unresolved
  blocker titles.
- [ ] A shaping map with only honest fog emits one map-level shaping row.
- [ ] Listing and count share one scan: a shaping map counts once, an invalid
  active candidate counts once, and a valid ready map counts zero.
- [ ] Ambient status consumes the same distinct active-map count and reports
  unknown when the active scan fails.
- [ ] Compiled candidates never enter active rows or count.
- [ ] Default, count, template, and validator calls are read-only and stable
  across repeated invocations.
- [ ] The terminal package contains no legacy Handoff acceptance path or second
  parser.
