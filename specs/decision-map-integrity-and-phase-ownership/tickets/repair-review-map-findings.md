# Repair review map findings

Blocked by: Repair dogfood source drift

## What to build

Close the two concrete map-validator defects found by the fresh three-axis
review.

## Acceptance

- [x] A source record rejects a second `Path`, second `URL`, mixed second
  locator, or any other field beyond its one locator, `Supports`, and `Drift`.
- [x] Source-field order and cardinality remain exact: one locator first,
  followed by one non-empty `Supports` and one non-empty `Drift`.
- [x] Cycle diagnostics identify both endpoint ticket titles, both graph
  handles, the map path supplied by tree validation, and the bad edge.
- [x] Focused tests cover the malformed second-locator cases and cycle message;
  the decision-map canary family bites the new source arm and retained cycle
  diagnostic.
