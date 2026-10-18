# Witness tip refusals and undecided manual review

Blocked by: none
Writes: internal/preflight/staleness/identity_test.go, internal/preflight/staleness/completeness_test.go, internal/preflight/staleness/decision_test.go
Covers: none

## What to build

This ticket adds tests only. It comes from drain d-359accf567c1, learnings L152 and L172.

L152: each oversized-object refusal witness of the staleness audit makes only the baseline object oversized.
Thus a change that drops the tip-side refusal of the test stage or of the ticket Covers read stays green.
Add a tip-side case beside each baseline case.
The command reads the tickets of its clean checkout at the tip before it collects.
Thus only `Collect` itself reaches a tip ticket object that does not read, and the ticket tip case calls `Collect`.

L172: in `Decide`, a check that is not green holds manual review required.
Each collector that makes a check diagnostic also adds a diagnostic, which holds manual review required by itself.
Thus a change of that branch to "only a red check" stays green.
Add a `Decide` case with an undecided test discrepancy, no diagnostic, and a ready caller.

The production code does not change.

## Acceptance

- [ ] An audit whose tip test object is over the control-record bound refuses at the cited test source stage and prints no staleness result.
- [ ] `Collect` over a pair whose tip ticket object is over the control-record bound refuses at the ticket source stage for that ticket.
- [ ] `Decide` over one undecided test discrepancy, no diagnostic, and a ready caller gives a complete audit with no drift, a diagnostic cited-tests check, and manual review required.
- [ ] A mutation that drops each tip-side refusal, and a mutation that holds manual review only for a red check, each turn the focused package check red.
