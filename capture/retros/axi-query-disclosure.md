# Retro — axi-query-disclosure

Closed 2026-08-12 at `2747e38b` (Status: implemented), after a full three-axis
review plus a rerun Fable cross-harness pass and a nine-ticket repair frontier.

## What the failures actually were

The amended spec already named every missed predicate; the defects were
evidence/oracle failures, not a bad spec. Concretely: assumed base bytes,
incomplete enumeration of QD6 states (unreadable proved only via the oversized
route; the completed-assignment row never driven publicly), weak marker
assertions standing in for structured TOON decoding against QD3, and review
passes that accepted ticket intent without falsifying the composed
implementation. The repair pass therefore strengthened oracles almost
exclusively: of nine tickets, seven changed only tests or comments, one
reworded a disclosure cell, one collapsed duplicated test scaffolding.

## What worked

- The fence discipline caught the one genuinely reviewer-owned question (the
  scaffold's home) at ticket review, before implementation.
- Disjoint `Writes:` sets let seven write delegates land serially with blind
  patch application and zero conflicts.
- Demanded red demonstrations (MB2, GS2, SE2) produced real counterexamples —
  including one (schema-wrong help table) the old assertions provably passed.

## Residuals, all routed

- R11 (active assignment, deleted tree) — parked idea, needs a spec amendment.
- Census scope for process-backed fixtures — parked idea, coupled with the
  gittest visibility note.
- Standing test-support commons fence — parked idea.
- Pre-existing `TestListCommandPublicRowsAndDisclosure` flake — parked idea,
  reproduced on the pre-repair baseline.
- `leasedRepo` left outside the gittest collapse (would need a third API
  shape) — contestable residue, flagged in the T1 ticket evidence.
