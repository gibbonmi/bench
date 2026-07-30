# Recurrence tallying

## Destination

Give roadmap prioritization a durable count of independently observed
occurrences that parsers can consume, without creating duplicate FT rows or
inventing a second roadmap parser.

## #1: Does recurrence determine priority or inform it?

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

n/a — the capture unit, migration, and multi-recommendation retro posture are
settled above.

## Out of scope

- A command that automatically mutates or globally sorts ROADMAP.md; the reviewed
  maintenance phase remains the prioritization owner.
- A second occurrence-history file or event database; the current key set lives
  on the row and Git is the historical record.
- Multiple primary owners on one capture unit; repeat the same incident key on
  separate units when one incident independently evidences several FTs.

## Sources

- `ROADMAP.md` — FT126 owns the roadmap parser/context seam and records the first
  FT98 recurrence.
- `internal/roadmap/context_parse.go` — one parser owns roadmap rows and all three
  capture-source facts in the schema-2 snapshot.
- `internal/roadmap/context_types.go` and `internal/roadmap/context_render.go` —
  drift-prone current snapshot fields that will carry any machine-readable tally.
- `internal/roadmap/roadmap.go` — ideas are dated free-text lines today; the
  roadmap command displays capture state but does not perform the drain.
- `.agents/commands/bench-what-next.md` — prioritization and the reviewed mutation
  remain judgment owned by the maintenance phase.

## Handoff

1. **Module boundaries.**
   - `internal/roadmap` is the deep owner of occurrence-token and row-ledger
     parsing, deduplication, owner validation, derived counts, context
     discrepancies, and the unified projection across capture sources.
   - `bench idea` argument handling is a thin producer of the same token grammar;
     its existing free-text capture behavior remains intact.
   - `internal/learnings` and `internal/retros` continue to own their document
     shapes. They expose capture units and bodies; recurrence meaning stays in
     `internal/roadmap` rather than being reimplemented per source.
   - `.agents/commands/bench-what-next.md` owns the reviewed add-key/remove-source
     and tie-breaker procedure; it does not carry another parser or count.
2. **Contracts.**
   - `bench idea --owner <FT> --incident <key> <text>` requires both metadata
     flags together, validates the token grammar, appends exactly one visible
     occurrence token, and preserves the existing exit and append contract.
   - A row has zero or one `Occurrences:` line. Its comma-separated incident
     keys are unique and canonical; the parser derives the count.
   - Context exposes owner, incident, source, capture unit, and state for pending
     occurrences, plus row count/keys and structured discrepancies. Structural
     discrepancy means the recommended sequence is untrusted.
   - Repeated `(owner, incident)` capture against a recorded row key is reported
     as already recorded, not counted again.
3. **Deep vs thin.** `internal/roadmap` hides grammar, cross-source normalization,
   deduplication, and discrepancy policy behind the existing document/context
   interfaces. CLI dispatch and source-specific readers pass typed text through
   and gain no recurrence policy of their own.
4. **Black-box assertables.**
   - Idea capture with valid metadata appends one canonical token; either flag
     alone, an invalid FT, or an invalid incident exits with usage/error and
     leaves the file unchanged.
   - Context reports two capture artifacts with the same owner/incident as one
     pending pair with both sources, then reports it as already recorded once
     the row contains the key.
   - A malformed token, absent owner row, duplicate row key, or multiple tokens
     on one capture unit appears as a structured discrepancy and prevents a
     trusted next-action projection.
   - Row count equals the number of unique stored keys; no heading count exists
     to drift.
   - Within equal severity and actionability, the reviewed sequence orders the
     higher count first; dependencies and explicit reviewer pricing still win.
5. **Gate attachment.** `internal/roadmap` unit tests own grammar, ledger,
   deduplication, and context projection. Runtime contracts own `bench idea`
   append/refusal behavior and `bench roadmap --context` AXI output. Conformance
   owns the command prose anchors and the absence of legacy header counts after
   migration. The gate cannot mechanically grade the maintenance phase's
   judgment that two rows share severity/actionability; the phase anchor and
   a fixture with equal-class rows grade the deterministic tie-break rule.
6. **Hostile-input owners.**
   - Spaces, glob characters, leading dashes, `--`, missing values, and
     multi-word idea text belong to the existing CLI grammar plus the new
     owner/incident arguments.
   - Control bytes, non-ASCII, separators, overlong keys, duplicate keys, CRLF,
     and a missing final newline belong to the occurrence and row parsers.
   - Absent versus empty capture files, special files, dangling symlinks, and
     bounded reads remain with the existing capture classifiers.
   - A retro containing several recommendation units and a learning body with
     token-like prose exercise capture-unit association rather than global text
     matching.
   - Re-run idempotency is the `(owner, incident)` deduplication contract.
7. **Uncertainty flags.** None. The reviewer approved the recommended posture
   for every ticket, including the syntax, tie-break rule, migration, and
   fail-visible behavior.
8. **Rejected alternatives.** Inferring ownership from prose; storing a mutable
   numeric count; a second ledger; silent malformed-token drops; multiple primary
   owners on one unit; counting every artifact independently; automatic global
   sorting; recurrence overriding severity, dependencies, or reviewer pricing.
9. **Domain watch-outs.**
   - Capture files disappear during a drain, so durable incident identity must
     land on the row before source removal.
   - The roadmap heading is already status-rich prose; putting the count there
     would duplicate the parser-owned key set and recreate drift.
   - A retired FT is not a valid current owner. Its history remains in Git rather
     than accepting new recurrence against a row that no longer exists.
   - Context schema consumers require an explicit schema advance when fields or
     tables change.

Dependency order: n/a — single spec.
