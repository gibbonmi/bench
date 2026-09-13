# FT232 repair-loop tripwire and bounded repair policy

Status: shaping

## Destination

Shape two separate deliverables: a bounded repair policy and an experimental repair-loop advisory.
Determine whether existing records can distinguish a repair loop from ordinary iteration.

## Notes

Origin: learnings; roadmap/FT232.md records the reviewed candidate.

Consult craft-domain, craft-grill, craft-line, and craft-synthesis.
Use craft-research for the evidence question before selecting a signal.

The lane span and the lane record are distinct records.
The lane span carries the first failing check and diagnostic.
The lane record omits those fields.
This observation does not prove that a repeated failure identifies a repair loop.

## Decisions so far

- [Separate deliverables](ft232-repair-loop/tickets/1.md): The policy and advisory receive separate specifications.
- [Blocking findings](ft232-repair-loop/tickets/2.md): Concrete defects and unmet requirements block completion; optional improvements remain advisory.
- [Repair allowance](ft232-repair-loop/tickets/3.md): Two repair cycles follow each chunk's initial review; a fresh review cannot reset the count.
- [Applicable runs](ft232-repair-loop/tickets/4.md): The allowance applies to all implementation runs, with an explicit reviewer override.
- [Evidence standard](ft232-repair-loop/tickets/5.md): Real examples must distinguish stalled repairs from productive iteration before detector specification.
- [Pre-review work](ft232-repair-loop/tickets/7.md): The cap applies after initial review; existing rules govern earlier work.
- [Reviewer confirmation](ft232-repair-loop/tickets/8.md): This pass is confirmed, with advisory research as the next frontier.

## Not yet specified

## Spec-writer discretion

## Out of scope

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
