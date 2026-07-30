# Recurrence tallying

Status: ready

## Destination

Give roadmap prioritization a durable count of independently observed
occurrences that parsers can consume, without creating duplicate FT rows or
inventing a second roadmap parser.

## #1: Does recurrence determine priority or inform it?

Blocked by: none
Type: Grill

### Question

Decide whether an occurrence count mechanically changes row ordering or is a
visible evidence input that `/bench-what-next` weighs alongside severity, cost,
dependencies, and reviewer pricing.

### Answer

Occurrence count is the deterministic first tie-breaker among rows with the
same severity and actionability. Higher severity, blocked versus actionable
state, literal dependencies, and explicit reviewer pricing remain stronger
ordering inputs, so recurrence cannot promote a minor row across those classes.

## #2: What is one occurrence?

Blocked by: #1
Type: Grill

### Question

Define the counted unit and duplicate posture when one incident appears in more
than one capture source or produces multiple recommendations in one retro.

### Answer

One occurrence is one unique `(primary FT owner, incident key)` pair. Multiple
capture artifacts or recommendations describing the same incident reuse the
same key and increment that FT only once. The same incident may increment
different FTs once each when it independently evidences each row.

## #3: How does capture cite a primary roadmap owner?

Blocked by: #2
Type: Grill

### Question

Choose one explicit owner-citation grammar across ideas, learning entries, and
retro recommendations, including unknown FT, malformed citation, and multiple
owner behavior. Inference from prose would make the count non-deterministic.

### Answer

Every capture unit ends with one visible token:
`[occurrence:FT98/2026-07-30-scoped-roadmap-commit]`. The owner is exactly one
existing roadmap ID. The incident key is lowercase ASCII letters, digits, and
hyphens, begins and ends alphanumeric, and is at most 64 bytes.

`bench idea --owner FT98 --incident <key> "<text>"` writes the token; both
flags are required together and the existing no-flag form remains valid for
capture that does not yet claim recurrence. A learning entry is one capture
unit. A retro recommendation paragraph or list item is one capture unit, so
one retro may carry several tokens on separate recommendations. More than one
token on a unit, a malformed token, or an owner absent from the current roadmap
is a visible context discrepancy that prevents trusting the recommended
sequence; it is never silently ignored or inferred from prose.

## #4: Where is the current tally stored and enforced?

Blocked by: #1, #2, #3
Type: Grill

### Question

Choose the durable ROADMAP representation and the `bench roadmap --context`
contract that checks citations against rows, exposes counts to the drain, and
fails visibly when the recorded tally and capture evidence disagree.

### Answer

Each roadmap row carries one exact `Occurrences:` line containing its sorted,
unique incident keys. The row ID supplies the owner, so `(row ID, incident key)`
is the stored identity and the count is always derived from the keys rather
than repeated in the heading or another ledger. Git history owns when each key
arrived.

The roadmap parser owns this line and rejects malformed or duplicate keys.
`bench roadmap --context` advances its schema and exposes each row's derived
`occurrence_count` and incident keys plus one capture-occurrence table spanning
ideas, learning entries, and retro recommendation units. Its discrepancy table
reports malformed tokens, unknown owners, multiple tokens on one unit, and
already-recorded pairs. An already-recorded pair is valid evidence for dropping
the duplicate capture; structural discrepancies keep the sequence untrusted.

`/bench-what-next` remains the reviewed mutator: it adds a new incident key to
the owning row before removing the capture source, then applies recurrence as
the tie-breaker from #1. No command silently reorders or rewrites ROADMAP.md.
The implementation migrates every header that currently claims `evidence
supplied` to incident keys recoverable from its row. Where an old count cannot
be mapped to distinct named incidents, deterministic baseline keys pinned to
the migration commit preserve the count without fabricating provenance.

## Not yet specified

## Sources

- Path: `ROADMAP.md`
  Supports: FT126 owns the roadmap parser and context seam and records the first FT98 recurrence.
  Drift: Update when the roadmap ownership or recorded recurrence changes.
- Path: `internal/roadmap/context_parse.go`
  Supports: One parser owns roadmap rows and all three capture-source facts in the schema-2 snapshot.
  Drift: Update when the parser's capture-source schema changes.
- Path: `internal/roadmap/context_types.go`
  Supports: Current snapshot fields carry any machine-readable tally.
  Drift: Update when tally fields move or change shape.
- Path: `internal/roadmap/context_render.go`
  Supports: Current snapshot fields render any machine-readable tally.
  Drift: Update when tally rendering moves or changes shape.
- Path: `internal/roadmap/roadmap.go`
  Supports: Ideas are dated free-text lines today and the roadmap command displays capture state without performing the drain.
  Drift: Update when idea capture or roadmap command ownership changes.
- Path: `.agents/commands/bench-what-next.md`
  Supports: Prioritization and reviewed mutation remain judgment-owned by the maintenance phase.
  Drift: Update when the maintenance phase changes ownership.

## Spec-writer discretion

## Out of scope
