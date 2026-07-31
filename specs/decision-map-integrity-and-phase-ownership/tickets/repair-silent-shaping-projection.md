# Repair silent shaping projection

Blocked by: Repair review map findings

## What to build

Make the detailed active-map projection agree with the distinct-map count for
a valid shaping map whose tickets are resolved and whose fog body is empty.

## Acceptance

- [x] Every shaping map with zero unresolved decision tickets emits exactly one
  `Not yet specified` / `fog` / `shaping` row, whether its fog body is empty or
  non-empty.
- [x] The synthetic row contributes no extra count; list and count remain two
  projections of the same active scan.
- [x] Ready maps remain omitted, invalid candidates remain visible once, and
  compiled maps remain excluded.
- [x] Focused model and AXI contracts cover the empty-fog shaping case and keep
  the non-empty-fog case green.
