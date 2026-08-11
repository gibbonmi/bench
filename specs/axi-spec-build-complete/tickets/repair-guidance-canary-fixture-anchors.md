# Re-anchor workflow canaries on retained evidence routes

Blocked by: none
Ownership fence: `tests/canary/workflow-guidance-anchors/ticket-observed-red-evidence/MUTATE.json`, `tests/canary/workflow-guidance-anchors/ticket-observed-red-evidence/EXPECT`, `tests/canary/workflow-guidance-anchors/ticket-already-covered-evidence/MUTATE.json`, `tests/canary/workflow-guidance-anchors/ticket-already-covered-evidence/EXPECT`, `tests/canary/workflow-guidance-anchors/ticket-not-tdd-able-evidence/MUTATE.json`, `tests/canary/workflow-guidance-anchors/ticket-not-tdd-able-evidence/EXPECT`
Integration surfaces: retained evidence-routing sentences in `bench-craft-tickets`→the three auto-discovered workflow-guidance canary mutations; mutation materialization→each fixture's exact expected bite
Contracts: the retained evidence-route clause crosses each fenced `tests/canary/workflow-guidance-anchors/ticket-observed-red-evidence/MUTATE.json`, `tests/canary/workflow-guidance-anchors/ticket-already-covered-evidence/MUTATE.json`, and `tests/canary/workflow-guidance-anchors/ticket-not-tdd-able-evidence/MUTATE.json` into its sibling EXPECT bite, asserted by CF1-CF3 before and after the guidance rewrite

## What to build

Re-anchor the three workflow guidance canaries that currently include the
mandatory subject-mutation wording being retired. Narrow each `old` value to the
stable evidence-route clause that remains policy: observed-red rows carry their
failing public operation into ticket acceptance; already-covered rows retain
their named control; and not-TDD-able rows map to the first ticket where their
seam exists. Give each mutation a replacement that drops that retained route,
and update its `EXPECT` text to name the same omission. Preserve the existing
fixture directories, auto-discovery shape, and one exact bite per fixture.

## Acceptance

- [ ] [CF1] (covers local) the observed-red fixture materializes against current guidance by removing the retained public-operation-to-ticket route, produces its exact omission bite, and remains applicable after mandatory mutation prose is removed.
- [ ] [CF2] (covers local) the already-covered fixture materializes against current guidance by removing the retained named-control route, produces its exact omission bite, and remains applicable after mandatory mutation prose is removed.
- [ ] [CF3] (covers local) the not-TDD-able fixture materializes against current guidance by removing the retained blocker-to-first-seam-ticket route, produces its exact omission bite, and remains applicable after mandatory mutation prose is removed.
