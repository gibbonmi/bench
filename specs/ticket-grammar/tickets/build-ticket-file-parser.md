# Build the ticket-file parser

Blocked by: lift-shared-field-grammar.md
Writes: internal/tickets/tickets.go (new), internal/tickets/tickets_test.go (new)
Covers: TG1, TG2, TG3, TG7, TG8, TG9, TG10, TG11, TG20, TG23, TG24, TG32, TG33, TG34

## What to build

A new package `internal/tickets` owns the ticket-file schema. The parser takes
immutable bytes, the sibling basename set, and the spec tag. It returns the
parsed ticket and its diagnostics. It reads no file and writes no file.

The parser requires the title line, three field lines, and two headings. The
field lines are `Blocked by:`, `Writes:`, and `Covers:`.
The two headings are `## What to build` and `## Acceptance`.
The parser tolerates an optional `## Delegate charge` heading, and it grades no
charge content.

The dependency grammar accepts `none` or the sibling ticket basenames. The
parser names a dangling blocker, a duplicate blocker, and a self-edge. It names
one edge of a blocker cycle.

The mutation grammar accepts `none` or a comma-separated row-ID list. The
parser names a repeated token. It also names a token whose tag differs from the
supplied spec tag.

Three hostile inputs stay fail-closed. A field prefix with a non-breaking space
yields the missing-field diagnostic. An unterminated fence yields its own
diagnostic, and the parser grades nothing after the opening marker. A required
field on the last line without a trailing newline still parses.

## Acceptance

- [ ] TG1 — the parser names each absent required field in one diagnostic.
- [ ] TG2 — a second `Blocked by:` line yields a duplicate-field diagnostic.
- [ ] TG3 — a field line inside a fenced code block parses as no field.
- [ ] TG7 — `Blocked by: none` parses to zero dependency edges.
- [ ] TG8 — a dangling blocker yields a diagnostic that names both basenames.
- [ ] TG9 — a repeated blocker basename yields a duplicate-blocker diagnostic.
- [ ] TG10 — a self-naming blocker yields a self-edge diagnostic.
- [ ] TG11 — a blocker cycle yields a cycle diagnostic that names one edge.
- [ ] TG20 — a row ID in body prose lands in no parsed `Covers:` set.
- [ ] TG23 — a `Covers:` token with a foreign tag yields a diagnostic.
- [ ] TG24 — a repeated `Covers:` token yields a duplicate diagnostic.
- [ ] TG32 — a non-breaking space in a field prefix yields the missing-field diagnostic.
- [ ] TG33 — an unterminated fence yields an unterminated-fence diagnostic.
- [ ] TG34 — a required field on the last line without a newline parses.

## Delegate charge

You work in the Bench repo on the `ticket-grammar` spec. Read
`specs/ticket-grammar/spec.md` first. Then read `internal/maps/fields.go`,
which a sibling ticket landed, and `internal/maps/schema.go` in full. Read
`internal/maps/maps_parse_test.go` for the table-test shape.

Create `internal/tickets`. Drive the field scan and the graph walk through the
lifted symbols in `internal/maps`. Do not write a second scan and do not write
a second walk.

Take immutable bytes, the sibling basename set, and the spec tag. Return the
parsed ticket and an ordered diagnostic list. Keep every file read out of this
package.

Require the title, `Blocked by:`, `Writes:`, `Covers:`, `## What to build`, and
`## Acceptance`. Tolerate `## Delegate charge`. Name each absent required
field.

Accept `none` or sibling basenames in `Blocked by:`. Name a dangling blocker
with both basenames. Name a duplicate blocker, a self-edge, and one cycle
edge.

Accept `none` or a comma-separated row-ID list in `Covers:`. Name a repeated
token. Name a token whose tag differs from the spec tag. Assert that a row ID
in body prose lands in no parsed `Covers:` set.

Keep the three hostile inputs fail-closed. Report the missing-field diagnostic
for a non-breaking space in a field prefix. Report an unterminated fence and
grade nothing after the opening marker. Parse a required field on an unterminated
last line.

Add `TestParseTicketRequiredFields`, `TestParseTicketFieldMechanics`,
`TestParseTicketBlockerGraph`, `TestParseTicketCoversGrammar`, and
`TestParseTicketHostileInput` in `internal/tickets`. Write each one as a table
test. Assert the exact diagnostic text.

Run only `bench worktree exec ft174-ticket-grammar -- go test ./internal/tickets/`.
Do not commit. Do not edit the spec.
