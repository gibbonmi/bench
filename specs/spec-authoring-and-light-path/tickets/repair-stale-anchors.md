# Repair the four stale anchor registry rows

Blocked by: none
Writes: internal/anchors/registry_data.go, internal/canary/inventory_test.go, tests/canary/workflow-guidance-anchors/

## What to build

Four registry rows no longer match the live tree and nothing reds because only
fixture-backed rows are enforced. Repair: the `projects/benchkit.md` needle
becomes `Hostile-input checklist` (matching the real casing) and the
`bench-review-implementation.md` needle becomes the current entry sentence
(`bench preflight review` in explicit-base mode); both gain new fixtures so
they are enforced from now on. Retire: the `shared-build-cache opt-in` row
(its subject was deliberately removed by the branch-native architecture change)
and the `new session on the mid tier` row (superseded by the authorship
successor authored in its own ticket — no file conflict: that ticket replaces
the row, this one does not touch it). The live-tree observable is the root-conformance sweep — `TestRootConformance`
with `BENCH_CONFORMANCE_ROOT` set (prep-release's run) — which reports all
four stale diagnostics today; record its before/after output in the ticket
evidence, with the ~8 unrelated pre-existing reds (including implement-spec's
missing Entry orientation / Exit handoff sections) named as inherited
baseline, never repaired here. Shares `registry_data.go` and the fixtures
tree with every sibling — those paths land serially across the whole spec.

## Acceptance

- [ ] the two repaired needles hold on the live tree and their new fixtures
      bite both halves (covers WF11)
- [ ] the `shared-build-cache opt-in` row is gone and the root-conformance
      sweep output shows zero stale diagnostics for the rows this ticket owns,
      inherited reds recorded (covers WF12)
- [ ] the ticket's fixture additions update the canary binding count in this
      ticket's green commit (covers WF23)
