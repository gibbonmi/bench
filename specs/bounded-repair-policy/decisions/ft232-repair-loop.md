# FT232 repair-loop tripwire and bounded repair policy

Status: ready

## Destination

Specify two separate deliverables: a bounded repair policy and a bounded collection pilot with an evidence report.
Ticket #14 confirms these outcomes.
The pilot supports a later decision about a repair-loop detector.

## Notes

Origin: learnings; roadmap/FT232.md records the reviewed candidate.

Consult craft-domain, craft-grill, craft-line, and craft-synthesis.
Consult craft-research when evaluating the pilot results.
The record signal report holds the existing evidence and its limits.

## Decisions so far

- [Separate deliverables](ft232-repair-loop/tickets/1.md): The policy and advisory work remain separate deliverables.
- [Blocking findings](ft232-repair-loop/tickets/2.md): Concrete defects and unmet requirements block completion; optional improvements remain advisory.
- [Repair allowance](ft232-repair-loop/tickets/3.md): Two repair cycles follow each chunk's initial review; a fresh review cannot reset the count.
- [Applicable runs](ft232-repair-loop/tickets/4.md): The allowance applies to all implementation runs, with an explicit reviewer override.
- [Evidence standard](ft232-repair-loop/tickets/5.md): Real examples must distinguish stalled repairs from productive iteration before detector specification.
- [Pre-review work](ft232-repair-loop/tickets/7.md): The cap applies after initial review; existing rules govern earlier work.
- [Reviewer confirmation](ft232-repair-loop/tickets/8.md): The reviewer confirms the policy decisions; ticket #14 holds the final scope.
- [Record evidence](ft232-repair-loop/tickets/6.md): Existing records do not demonstrate the required distinction between stalled repairs and productive iteration.
- [Collection experiment](ft232-repair-loop/tickets/9.md): The reviewer authorizes bounded evidence collection before a later detector decision.
- [Pilot scope](ft232-repair-loop/tickets/10.md): Explicit activation limits the pilot to ordinary work in the Bench kit.
- [Progress evidence](ft232-repair-loop/tickets/11.md): Verification supports each progress label; unsupported labels remain unknown.
- [Pilot bounds](ft232-repair-loop/tickets/12.md): Ten completed sequences or fourteen days stop the pilot; missing example classes make the result inconclusive.
- [Repair sequence](ft232-repair-loop/tickets/13.md): One chunk's sequence runs from its first blocker to verified closure or reviewer handoff.
- [Specification scope](ft232-repair-loop/tickets/14.md): Separate policy and collection specs exclude a detector and warnings pending a later reviewer decision.

## Not yet specified

## Spec-writer discretion

## Out of scope

- A detector or warning before the reviewer evaluates the pilot results.
- Automatic model changes or gate overrides from the advisory.
- Weaker required checks or behavioral guarantees.
- A fixed repair cap before initial review.

## Sources

- Path: `roadmap/FT232.md`
  Supports: The candidate outcomes and existing constraints.
  Drift: The reviewer changes the candidate scope or constraints.
- Path: `internal/gate/lane.go`
  Supports: The distinction between the lane span and the lane record.
  Drift: The lane changes the fields it retains or the record it writes.

- Path: `specs/bounded-repair-policy/decisions/ft232-repair-loop/assets/record-signal-evidence.md`
  Supports: Tickets #6 and #9; existing record limits and the collection decision.
  Drift: Record producers change or new labeled repair sequences become available.
