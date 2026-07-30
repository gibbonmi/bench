# Project the active-map model

Blocked by: Migrate the tracked decision-map corpus

## What to build

`bench maps`, its count mode, and ambient status switch atomically to the
migrated model, and the legacy Handoff parser leaves `internal/maps`.

## Acceptance

- [x] Default output is `maps[N]{map,title,type,state,blockers}` and shows every
  unresolved decision ticket as frontier, blocked, or deferred with unresolved
  blocker titles.
- [x] A shaping map with only honest fog emits one map-level shaping row.
- [x] Listing and count share one scan: a shaping map counts once, an invalid
  active candidate counts once, and a valid ready map counts zero.
- [x] Ambient status consumes the same distinct active-map count and reports
  unknown when the active scan fails.
- [x] Compiled candidates never enter active rows or count.
- [x] Default, count, template, and validator calls are read-only and stable
  across repeated invocations.
- [x] The terminal package contains no legacy Handoff acceptance path or second
  parser.
