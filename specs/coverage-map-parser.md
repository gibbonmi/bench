# Coverage-map parser

## Problem
The gate anchors the acceptance-coverage *workflow prose* (the phrases must
exist in the commands), but nothing reads an actual spec's map. A future spec
can ship a coverage table with a missing cell, an empty red signal, or a row
pointing at a story that doesn't exist, and the gate stays green — the veto
surface silently degrades.

## Solution
A gate check parses every `specs/*.md` that carries an
`### Acceptance coverage map` heading: the canonical five-column header must be
present, every data row must have five non-empty cells, and every row's story
reference must resolve — a number no higher than the spec's numbered stories,
or an `edge…` reference. Errors name the spec and the row. Specs from before
the convention (or deliberately exempt) opt out with a
`<!-- coverage-map: historical -->` marker, mirroring the command-currency
marker.

## User stories
1. As the reviewer, I want a malformed coverage row (wrong cell count, empty
   cell) to fail the gate naming the spec and row, so a degraded veto surface
   is caught at commit time instead of at review.
   Line: claude-fable-5 (inline, session model) / medium. Gate/conformance
   logic — oracle correctness over speed.
2. As the reviewer, I want a story reference that doesn't resolve (story 7 of
   5) to fail the gate, so coverage rows can't drift from the stories they
   claim to cover.
   Line: claude-fable-5 (inline, session model) / medium. Same check.
3. As the reviewer, I want pre-convention specs opt-out-able via a historical
   marker, so history doesn't have to be rewritten to keep the gate green.
   Line: claude-fable-5 (inline, session model) / low. One skip clause.
4. As the kit maintainer, I want a canary fixture proving the parser bites, so
   the check can't rot into an always-pass.
   Line: claude-fable-5 (inline, session model) / low. Fixture + EXPECT, the
   established pattern.

## Implementation decisions
- Node heredoc appended to the docs-contracts gate fragment (the gate's
  existing pattern for structured parsing); no new files, nothing ships.
- Scope: only files under `specs/` containing the map heading. No heading, no
  check — the write-spec phase is what obliges new feature specs to carry one.
- Canonical header: `| story | behavior | seam | red signal | why it catches
  the failure |`. A map heading without it is an error (that is the drift).
- Story references accepted: an integer within 1..N where N = count of
  `^<n>. ` items under `## User stories`, or a reference beginning with
  `edge`. Anything else is an unrecognized reference error.
- Existing specs that fail the parse get the historical marker rather than a
  rewrite — history lives in git; the parser is for future specs.

## Testing decisions
- Good test = the canary: plant a spec with a defect, assert the inner gate
  goes red with the attributed message.
- Seam: the gate/canary seam; self-proving on this repo's own specs (the gate
  must stay green over `specs/` as it exists).
- Gate command: `.bench/gate.sh`.

### Seam diagram

    trigger: gate run (kit repo or canary inner run)
        │
        ▼
    specs/*.md with map heading ──▶ [ coverage-map parser ] ──▶ attributed errors
    ## User stories numbering   ──▶ [  (node, in gate)    ] ──▶ exit code via fail
                                        ◀ tests attach here: canary fixture plants
                                          a defective spec, asserts red + message

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | row with an empty cell → red naming spec + row | canary | new fixture fails "canary did not bite" until parser lands | empty red-signal cells are the observed degradation |
| 2 | story number above the spec's count → red | canary (same fixture, second defect asserted via its message on this repo during dev) | parser red on planted spec | dangling story refs are undetectable by anchors |
| 3 | marker file skipped | self-proving: any pre-convention spec that fails gets the marker and the gate returns green | gate red on this repo until markers placed | proves opt-out actually skips |
| 4 | parser can't rot silently | canary vacuous-baseline + bite assertions | already covered by canary framework | the framework asserts both baseline and bite |
| edge of 1 | map heading but no table/header → red | canary fixture family (header-missing variant folded into fixture) | parser red | heading-without-table is the laziest drift |

### Edge inventory
- empty/absent input → no heading = out of scope by decision (write-spec obliges
  new specs); heading-without-table = coverage row above.
- malformed input → story 1 row.
- boundary values → story ref exactly N is valid, N+1 red (story 2 row).
- **Won't handle:** validating `seam`/`behavior` cell *semantics* — prose
  quality stays a review responsibility, per the anchors' own comment.
- **Won't handle:** maps outside `specs/` — the convention lives there.
- re-run idempotency → read-only check; nothing to assert.

## Out of scope
- Parsing edge-inventory Won't-handle lines for structure — a second grammar
  over looser prose; ~6 edits, ~3 gate runs if the map parser proves its worth.
