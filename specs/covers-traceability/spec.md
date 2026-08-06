# Spec→ticket coverage traceability

Status: staged

Decision source: reviewer-confirmed conversation of 2026-08-05 (grill: mechanism, row IDs, enforcement loci, rollout, local rows, grammar)

## Problem

A spec's acceptance coverage map enumerates behaviors per subclause — hostile
paths, special-file targets, SIGINT — but nothing binds ticket acceptance rows
to those map rows. A ticket author can compress several mapped behaviors into
one compound row ("validation and interruption failures preserve prior
output"), a delegate proves one example, the per-row receipt marks the whole
row, and the lost subclauses reach semantic review as untested gaps. In a real
build this shipped two bugs and two missing test classes past checkpoint.

## Solution

Give coverage-map rows stable spec-local IDs, and give every ticket acceptance
row an inline `covers` annotation naming exactly one map row (or declaring
itself `local`). `assign` refuses a malformed or dangling mapping per ticket;
`promote` refuses a composition that leaves any map row uncovered. Because
mapping is 1:1 at the row level, the existing per-row checkpoint receipts
become per-subclause evidence with no receipt schema change. Enforcement keys
off the spec: a map without row IDs takes today's path unchanged, so no
in-flight build breaks.

## User stories

1. **A spec author opts a coverage map into row IDs and `bench coverage
   --check` validates them.** A map whose canonical header leads with a `row`
   column is opted-in; `--check` validates ID grammar, spec-local uniqueness,
   and no empty ID cells, and reports a map that mixes ID and non-ID rows.
   `Line: opus / medium.` Validation logic at a known parser seam, fully
   gate-observable.
2. **`ParseTicket` reads a per-row `covers` annotation.** Each acceptance row
   may carry `(covers <ID>)` or `(covers local)` after its row-ID bracket; the
   parse captures the mapping and stays permissive, so legacy tickets and the
   conformance-mutated examples parse exactly as today.
   `Line: opus / medium.` Grammar extension beside existing parse tests.
3. **`assign` refuses a covers-policy violation for an opted-in spec.** Under
   an opted-in spec, every row must carry a `covers` annotation resolving to
   exactly one existing map ID or `local`; a missing, malformed, or dangling
   annotation is refused with the offending row named. The annotation grammar
   attaches only to a single-ID row, so an R-range line's expanded rows are
   unannotated and refuse under the same missing-covers rule. Tickets of a
   non-opted-in spec are untouched.
   `Line: opus / medium.` Assignment policy following the `ContractsAnchored`
   precedent.
4. **`promote` refuses an uncovered map row.** Before the prospective gate
   owner executes, the union of non-`local` covers across the run's integrated
   assignments' tickets must include every map ID; the refusal names the
   uncovered IDs and spends no gate run.
   `Line: opus / medium.` Lifecycle logic with existing promote-refusal prior
   art.
5. **The skills teach the schema and the taught example stays parser-agreed.**
   `craft-spec` gains the optional `row` column and its opt-in meaning;
   `craft-tickets` gains the covers grammar, the `local` rule, and an updated
   taught example. The conformance example-agreement check's expected literals
   grow a covers expectation, so stripping the annotations from the taught
   example turns that check red rather than merely thinning the prose.
   `Line: fable / high.` Kit guidance prose takes the leverage override.

## Implementation decisions

**One parser per artifact, consumed cross-package.** The coverage package
already owns spec-map parsing; it gains the optional leading `row` column
(6-cell canonical header) and exports the parsed IDs. The lifecycle consumes
that parse for opt-in detection, assign resolution, and promote totality — it
never re-derives map structure. Likewise `ParseTicket` stays the one ticket
grammar; assign and promote read its parsed covers rather than rescanning
files.

**ID grammar is the existing ticket-ID grammar, spec-local.** An uppercase tag
plus a number, unique within the spec's map. Ticket row IDs and map row IDs are
separate namespaces; the covers annotation is the only bridge.

**Opt-in is a property of the map, not a version flag.** A 6-cell map with a
leading `row` column is opted-in; a 5-cell map is legacy and grades exactly as
today. A 6-cell map with invalid IDs fails `--check`, and assign under such a
map refuses rather than guessing — fail closed once the author has opted in.

**Parse stays permissive; policy refuses.** Following the `ContractsAnchored`
precedent, `ParseTicket` never refuses on covers content — the conformance
example-agreement check grades deliberately mutated examples through it, and a
parse refusal would shadow the policy diagnostic. Assign owns every refusal.

**`local` is accepted machinery, graded honesty.** `(covers local)` marks a
ticket-time discovery or repair row; assign accepts it, promote's totality
ignores it, and review grades whether the marker is honest. A silent missing
annotation is refused, so omission cannot dodge the mapping.

**Promote grades the integrated composition, not the tickets directory.**
Totality is computed over exactly the tickets bound to the run's integrated
assignments: each assignment records its ticket path and content digest, and
promote re-parses those files and refuses on a digest mismatch, so a
post-integration covers edit and a decoy ticket file that was never assigned
both refuse rather than count. The totality check runs before the prospective
gate owner executes, so an uncovered map refuses without spending a gate run.
Over-coverage (two rows covering one map ID, across tickets or within one) is
legal; only a zero-cover map ID refuses.

**The covers grammar attaches to single-ID rows only.** `ParseTicket` expands
an R-range before policy runs and discards range provenance, so no annotation
attaches to a range line; its expanded rows are simply unannotated, and the
opted-in missing-covers refusal covers them with no special case.

**No state or receipt schema change.** The assignment record and checkpoint
receipts are untouched: the ticket digest already pins covers edits after
assign, and 1:1 row mapping makes the existing per-row receipt the
per-subclause evidence.

**This spec's own map stays legacy.** The 6-cell checker does not exist until
this build lands, and the staged-spec `--check` sweep is fail-closed, so the
map below is 5-cell by necessity, not oversight.

## Testing decisions

- **External behavior.** Tests drive the real commands and parsers: `--check`
  over authored map fixtures, `ParseTicket` over ticket files on disk, assign
  and promote through the lifecycle service against fixture repos, never
  through re-implemented grammar.
- **Coverage-checker seam.** The coverage package's `Check` tests are the prior
  art; new cases cover the 6-cell header, ID grammar, duplicates, empty ID
  cells, and mixed maps.
- **Ticket-parse seam.** The existing `ParseTicket` tests are the prior art;
  new cases cover annotated rows, unannotated legacy rows, malformed
  annotations parsing as unannotated, and digest sensitivity to a covers edit.
- **Assign-policy seam.** The lifecycle refusal tests are the prior art; new
  cases cover each refusal class and the legacy-spec pass-through.
- **Promote seam.** The lifecycle promote tests are the prior art; new cases
  cover the uncovered-ID refusal, its naming, the all-covered pass, and
  `local` contributing nothing.
- **Conformance seam.** The example-agreement check grades `craft-tickets`'
  updated taught example through the production `ParseTicket`; the docs needle
  checks pin the new schema clauses in both skills.
- **Gate seam.** The dev-tier `bench gate` core-package and conformance phases
  observe every row above.

### Seam diagram

    trigger: bench coverage --check / spec build assign / spec build promote
        │
        ▼
    spec.md map ──▶ [ coverage parse: IDs, opt-in ] ──▶ map IDs + violations
        │                 ◀ tests attach: Check over authored map fixtures
        ▼
    ticket file ──▶ [ ParseTicket: rows + covers ] ──▶ per-row mapping
        │                 ◀ tests attach: parse over ticket files on disk
        ▼
    [ assign policy: resolve each row ]──refusal──▶ named offending row
        │                 ◀ tests attach: lifecycle refusal fixtures
        ▼
    [ promote totality: union of covers ]──refusal──▶ named uncovered IDs
                          ◀ tests attach: lifecycle promote fixtures

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A 6-cell map with a leading `row` column parses and exports its IDs | coverage parse tests | current `Check` refuses at the header comparison with "coverage map missing the canonical header"; export test to be observed red at build time | A checker that still refuses the 6-cell header makes opt-in impossible |
| 1 | Duplicate IDs, empty ID cells, bad-grammar IDs, and ID/non-ID mixes are `--check` violations | coverage `Check` tests | to be observed red at build time; no ID validation exists | An unvalidated ID column lets a dangling covers reference grade as author error later instead of at authoring |
| 1 | A 5-cell map still checks exactly as today | existing coverage `Check` tests | already covered; run unchanged | The legacy path is the rollout guarantee |
| 2 | `ParseTicket` captures `(covers <ID>)` and `(covers local)` per row | parse tests over ticket files | to be observed red at build time; `Ticket` has no covers field | Assign and promote cannot enforce a mapping the parse never captured |
| 2 | Unannotated and malformed-annotation rows parse as unannotated, and legacy tickets parse byte-for-byte as today | existing plus new parse tests | already covered for legacy tickets; malformed cases to be observed at build time | A parse refusal would shadow the conformance example-agreement diagnostic and break staged legacy tickets |
| 2 | Editing a covers annotation changes the ticket digest | parse digest test | already covered: digest hashes full content; regression control | Digest pinning is what keeps checkpoint honest about post-assign covers edits |
| 3 | Under an opted-in spec, a missing, malformed, or dangling covers is refused with the row named, and an R-range line's expanded rows refuse as unannotated | lifecycle assign fixtures | to be observed red at build time per refusal class | Any silent acceptance reopens the compound-row failure this spec exists to close |
| 3 | `(covers local)` and valid mappings assign; a non-opted-in spec's tickets assign exactly as today | lifecycle assign fixtures | legacy pass-through already covered by existing assign tests; local acceptance to be observed at build time | Over-refusal would force spec edits on every repair ticket and break in-flight builds |
| 3 | An opted-in spec whose map fails ID validation refuses assign | lifecycle assign fixture | to be observed red at build time | Failing open on an invalid opted-in map lets a broken map disable enforcement silently |
| 4 | Promote refuses when any map ID has zero non-`local` covers across the integrated assignments' tickets, naming the uncovered IDs, before the gate owner executes and with the fixture gate recording zero executions | lifecycle promote fixtures with a recording gate owner | to be observed red at build time | The compound row reached review because no composition point ever computed totality, and a post-gate check would spend a full gate run to say no |
| 4 | A decoy ticket file in the tickets directory that no assignment binds contributes nothing to totality | lifecycle promote fixture with an unassigned covering ticket | to be observed red at build time | Directory-scoped totality re-opens the failure one level up: an unassigned, unverified file could satisfy every map ID |
| 4 | A post-integration covers edit to an assigned ticket refuses at promote on the digest mismatch | lifecycle promote fixture editing the ticket after integrate | to be observed red at build time | Without the digest check, coverage claims could be rewritten after the evidence that graded them |
| 4 | Full coverage promotes; over-coverage and `local` rows do not refuse | lifecycle promote fixtures | to be observed at build time | Refusing over-coverage would forbid legitimate defense in depth |
| 5 | `craft-tickets`' taught example carries covers annotations, and the example-agreement check's grown expected literals turn red when the annotations are stripped | conformance example-agreement check with a covers expectation | current check grades only row IDs, fence, blockers, and mutations against hand-authored literals, so stripping covers stays green today | The taught example is the grammar most ticket authors copy, and without a covers literal the check cannot see the annotations at all |
| 5 | Both skills' new schema clauses are pinned so deleting one turns the conformance docs phase red | docs needle checks | to be observed red at build time per clause | Unpinned guidance prose drifts from the parser it describes |

The cheapest degenerate is a promote that computes totality over the tickets
*directory* instead of the integrated assignments — a decoy file nobody
assigned satisfies every map ID; the decoy row goes red on it. The next
cheapest is an assign that checks only annotation presence, not resolution —
every row says `covers` something and nothing verifies the target exists —
combined with a promote that never reads the map; the dangling-refusal row and
the uncovered-ID promote row each go red on that.
The composition degenerate keeps both package rows green on fixtures while the
real commands never wire coverage's parse into the lifecycle; the assign and
promote rows are therefore driven through the lifecycle service against real
spec and ticket files, not against a fixture map handed to the policy function.

### Edge inventory

- Error path — resolved by the refusal rows (assign classes, promote
  uncovered-ID, invalid-map fail-closed).
- Empty or absent input — an empty map already fails `Check` ("no data rows");
  an opted-in spec with zero assignments passes the vacuously-true
  integrated-and-released loop, and the totality check is what refuses it:
  every map ID is uncovered. Resolved by the promote rows.
- Boundary values — one-row map fully covered by one ticket row; a map ID
  covered only by `local` rows refuses. Resolved by the promote rows.
- Malformed input — malformed annotations parse as unannotated (story 2 row)
  and then refuse at assign under opt-in (story 3 row); IDs are constrained to
  the uppercase-tag grammar, so bracket or backtick payloads never match.
- Interrupted or partial state — assign and promote keep their existing
  operation journal and re-entry semantics; no new filesystem state is
  introduced. Resolved by existing lifecycle controls as regression.
- Re-run idempotency — resolved by existing assign/promote re-entry tests
  running unchanged.
- Process-boundary lifecycle — a covers edit after assign changes the ticket
  digest and trips the existing checkpoint digest refusal; regression control
  (story 2 digest row).
- Hostile environment — **Won't handle:** hostile filesystem paths for spec
  and ticket files; both are repo-committed paths already constrained by
  `ParseTicket`'s traversal refusals, which run unchanged.

## Out of scope

- Per-covered-row receipt schema extension — rejected by the mechanism
  decision, not deferred: 1:1 rows make it redundant.
- Retrofit tooling that adds IDs and covers to legacy specs and tickets — a
  separate migration capability, about 6 edits and 1 gate run.
- A reporting command rendering build-wide coverage totality outside promote —
  a separate observability capability, about 4 edits and 1 gate run.
- Cross-spec coverage (one spec's ticket covering another spec's row) — no
  decided semantics; a future decision if the need appears.
