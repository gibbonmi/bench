# Project capture occurrences

Blocked by: Project occurrence ledgers and migrate recurrence

## What to build

Roadmap context consumes source-owned idea, learning, and retrospective recommendation
units, interprets only their final occurrence token, normalizes owner/incident pairs,
and renders deterministic capture and discrepancy tables without mutating the tree.

## Acceptance

- [x] Schema 3 renders `sequence_trusted`, both roadmap occurrence columns,
  `capture_occurrences`, and `occurrence_discrepancies` in fixed block order.
- [x] Capture rows sort by owner, incident, source, capture unit, and state; duplicate
  sources preserve both locations but normalize to one pending pair.
- [x] Recorded pairs become `already-recorded`, add only a non-structural discrepancy,
  and never contribute another pending pair.
- [x] Malformed final tokens, unknown owners, multiple tokens, and malformed ledgers
  emit their named structural discrepancies and make the sequence untrusted.
- [x] Learning and retrospective tokens stay attached to their source-owned units;
  same incident keys for different owners remain independent pairs.
- [x] Repeated and nested-cwd context calls stay byte-identical and read-only across
  special files, dangling symlinks, bounded inputs, and TOON byte classes.
