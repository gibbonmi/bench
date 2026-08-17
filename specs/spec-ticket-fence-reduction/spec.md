# Spec and ticket fence reduction

Status: staged

Decision source: `specs/spec-ticket-fence-reduction/decisions/spec-ticket-fence-reduction.md` (ready compiled map, all eighteen decision tickets reviewer-closed 2026-08-16; #15 carries a same-day superseding note fixing this build's review tier at `fable`/high)

Verification log: spec 1 + tickets 1 iteration(s) to accept — one round over the pair with `fable`/high through the native agent surface, per map #3 applied to this build's own process. BLOCK on 7 findings, all folded: three false or non-discriminating red-signal cells (SR3, SR6, and the `wantTable` justification), a missing ticket blocker, three fence/`Writes:` mismatches, an incomplete anchor inventory, and a decision-source divergence on `--reviewer` the reviewer closed as tier-only. The round also caught that commit 24cad87d already landed the `craft-tickets` prose remake, which cut stories 28–31 down to enforcement-only work. At approval the reviewer split the template/loop ticket in two along the file-versus-process seam, taking the breakdown from 7 slices to 8.

## Problem

The spec artifact predicts test behavior before any test runs. Its `red signal`
cell asserts something about current code — "already covered: X", "today Y does
Z" — at authoring time, so the cell can be false, and catching a false cell is
what the pre-implementation review loop spends its rounds on. On FT210 five of
ten blocking findings were exactly that: false current-code claims in coverage
cells, one of which reopened a closed decision built on the wrong premise. Four
delegated rounds (spec ×2, tickets ×2) then terminated on the two-round cap with
partials folded rather than on convergence, and the resulting spec ran 365 lines
over a 263-line map for a feature that sliced to five tickets.

## Solution

Cut the cell that makes the prediction and the loop that verifies it. The
acceptance coverage map keeps `row`, `story`, `behavior`, `seam`, and `why it
catches the failure` and loses `red signal`; the four-form red-signal grammar is
deleted from `craft-spec`, leaving `craft-tdd` — which classifies a row at the
moment it runs — as its one owner. `internal/coverage` accepts the reduced header
beside the existing one so nothing in flight breaks, and projects one uniform
`rows[N]{story,behavior,seam}` shape whichever schema a spec uses. The two
review loops become one round over the spec-and-tickets pair. The edge-class walk
moves to the TDD loop where the seam is visible, leaving `**Won't handle**` for
reviewed exclusions. `craft-tickets` takes `to-tickets`' sizing rule, and the
spec template moves into `craft-spec` so `bench-write-spec.md` fits a new 73-line
budget (reviewer-accepted 2026-08-17 in place of the map's 60; see the prose-landing bullet).

## User stories

**Reading a spec under the reduced schema** — Line: `opus` / high. Gate and
conformance logic: a parser that mis-accepts a header decides what every later
spec is graded against, and correctness of the oracle outranks speed.

1. As a spec author, I want `| row | story | behavior | seam | why it catches the failure |` accepted as a canonical header, so that a spec can omit the red-signal cell.
2. As a spec author, I want `| story | behavior | seam | why it catches the failure |` accepted, so that the row-ID opt-out still exists under the reduced schema.
3. As a spec author, I want both existing headers still accepted unchanged, so that a staged spec keeps validating while the six-column form drains.
4. As a spec author, I want the parser to choose a schema by the header's cell *names* and never by its cell count, so that a five-cell reduced header and a five-cell legacy header cannot be confused.
5. As a spec author writing a reduced map, I want a wrong-width data row reported against the reduced width, so that the error names the schema I am actually using.
6. As a spec author, I want the empty-cell violation to name the reduced schema's field, so that the message points at a column that exists.
7. As a spec author, I want the one-predicate-per-behavior check to read the behavior cell at its reduced offset, so that a `;` in a behavior is still refused.
8. As a spec author, I want the story-reference, fan-out-bound, duplicate-row-ID, malformed-row-ID, and orphan-story checks to behave identically under both schemas, so that cutting a column does not silently cut a check.
9. As an agent running `bench coverage <spec>`, I want the projection to be `rows[N]{story,behavior,seam}` for every accepted header, so that one shape seeds implementation tasks regardless of the spec's schema.
10. As an agent, I want the AXI action list and its per-row remedy unchanged, so that repairing a broken map is still one named invocation.
11. As a maintainer, I want `ParseSpec`'s opt-in verdict, ordered row IDs, and violations unchanged for a six-column spec, so that `internal/preflight`'s fence verdict does not move with the column cut.
12. As a spec author, I want a header that renames only some cells refused as "missing the canonical header", so that a half-migrated table is a loud failure rather than a silent partial parse.
13. As a spec author, I want a reduced map with no data rows to report "coverage map has no data rows", so that an empty table is not mistaken for a valid one.
14. As a spec author, I want a reduced spec with no map and no historical marker still reported as a violation, so that the reduced schema is not a route to an ungraded spec.
15. As a spec author, I want the historical marker to skip validation under either schema, so that the opt-out keeps working.
16. As a maintainer, I want every in-tree test fixture that embeds a coverage map migrated to the reduced header, so that the reduced schema is exercised through its real consumers and the eventual contract is a parser-only edit.

**Authoring under the reduced fences** — Line: `fable` / high. Guidance prose
compounds through every session that loads it while the edit costs few tokens —
the profile's doc-authoring leverage override.

17. As a spec author, I want the spec template to live in `craft-spec`, so that the authoring discipline and the artifact shape have one owner.
18. As a cold session, I want `bench-write-spec.md` to hold only the entry contract, ownership, and exit handoff, so that the phase file is read rather than skimmed.
19. As a maintainer, I want `.agents/commands/bench-write-spec.md` budgeted at 73 lines in the profile's table, so that the shrink is enforced rather than a one-time diff that re-accretes.
20. As a spec author, I want one review round over the spec-and-tickets pair instead of two rounds per artifact, so that the loop converges on a verdict rather than on a cap.
21. As a spec author, I want `craft-tickets`' granularity / blocking-edges / merge-split quiz to be that round's approval step, so that ticket approval is not a separate loop.
22. As a spec author, I want `--reviewer <tier> [effort]` to resolve same-family through the invoking harness's own `.bench/lines.env` column only, so that one resolution path replaces both the family fork and the bound-model-id lookup.
23. As a spec author, I want the four-form red-signal grammar deleted from `craft-spec`, so that no cell asks me to assert current-code behavior before a test runs.
24. As a builder, I want `craft-tdd` to be the one owner of `already covered` and `not TDD-able` classification, so that the classification happens where it is observed.
25. As a spec author, I want the canonical edge-class walk to live in `craft-tdd` at the seam, so that edges are enumerated where they are visible.
26. As a reviewer, I want `**Won't handle**` kept for deliberately excluded edges, so that exclusions still have a veto surface.
27. As a spec author, I want the requirement that every walked edge class produce a row or a `**Won't handle**` line retired, so that an exhaustive disposition list stops padding the artifact.
28. As a maintainer, I want `craft-tickets`' acceptance-row clause anchored with a canary that fails for its own planted reason, so that the sizing rule cannot silently regress.
29. As a spec author, I want this spec to record that the vertical-slice rule, prefactoring-first, and the retirement of "smallest independently-green" already landed in commit 24cad87d, so that no ticket re-does landed work.
30. As a maintainer, I want "independently-green implementation tickets" kept as README, CHANGELOG, and CONTEXT vocabulary, so that implementation tickets stay distinguishable from shaping decision tickets.
31. As a maintainer, I want `craft-tdd`'s anchored sentence naming `craft-spec` as the one source of the red-signal definition reworded when that ownership moves, so that an anchored sentence does not survive as a false claim.
32. As a maintainer, I want every anchor whose needle moved to a new file or section retargeted rather than deleted, so that the enforcement surface still bites after the move.
33. As a maintainer, I want each retargeted anchor's canary fixture updated to keep failing for its own planted reason, so that a green canary suite still proves the anchors bite.
34. As a maintainer, I want `craft-delegate`'s write-delegation charge to carry behavior and seam rather than a red signal, so that the charge matches the rows a spec now holds.
35. As a maintainer, I want `craft-review`'s Coverage axis and `/bench-implement-spec`'s task-seeding prose to name the reduced projection, so that the consumers describe the shape they receive.
36. As a maintainer, I want `CONTEXT.md` to define **coverage row** and **acceptance row** in the same green change as the behavior, so that the glossary never describes a tree that does not exist yet.
37. As a maintainer, I want the shipped docs that describe the coverage map to state the reduced schema, so that the field guide does not teach a retired column.
38. As a maintainer, I want a CHANGELOG entry for the reduced schema and the single review round, so that the change is discoverable.
39. As a maintainer, I want `bench skills-index --check` green after the skill descriptions change, so that the index and the skills cannot disagree.

**Explicitly not wanted** (see Out of scope)

40. As a maintainer, I do not want the six-column parser branch deleted in this change, so that specs staged under it keep validating.
41. As a maintainer, I do not want the light-path threshold moved, so that whether the cheaper heavy path gets used stays measurable.
42. As a maintainer, I do not want `/bench-implement-spec`'s worktree, preflight, or landing choreography changed, so that a real capability is not cut with the ceremony.
43. As a maintainer, I do not want `references/cross-harness-reviewers.md` deleted, so that cross-family delegation stays available to other phases.

Story partition: 1–16 land in `internal/coverage` and its fixture consumers;
17–39 land in guidance prose and its enforcement surface. The two share the
column list — the anchors and canaries pin it to the prose — so they ship as one
bundle (map #9).

## Implementation decisions

- **Header choice is by cell names, never cell count.** The parser already
  switches on the joined lowercase header names. Two more cases join it:
  `row|story|behavior|seam|why it catches the failure` (reduced, opted into row
  IDs) and `story|behavior|seam|why it catches the failure` (reduced, legacy row
  shape). A five-cell header is therefore unambiguous — the legacy five-cell form
  names `red signal`, the reduced one does not.
- **One schema descriptor, three derived values.** `width`, `fields`, and
  `storyOffset` currently derive from one opt-in flag. They derive instead from a
  single schema value carrying the field-name list; every offset-taking check
  reads its cell through that descriptor rather than a literal index, so a check
  cannot silently read the wrong column under one schema. Adding a schema is
  adding one descriptor, not editing five call sites.
- **The projection becomes `story, behavior, seam`.** `Rows` today returns
  story/seam/red-signal, and `bench coverage <spec>` renders it as
  `rows[N]{story,seam,red_signal}`. Behavior is the cell that names what to
  build, and it exists under every accepted header — so the projection is uniform
  across schemas instead of returning a column half the specs will not have. This
  is an observable output change to an AXI-approved command; the AXI action list,
  its per-row remedy, and the `spec:`/`state:` lines are unchanged.
- **The checked-in goldens are the independent expectation.**
  `internal/coverage/testdata/` holds five `.stdout` files that
  `TestCommandPreservesCheckedInPreDisclosureResponses` compares `Command` output
  against; `pre-disclosure-mapped.stdout` and `pre-disclosure-malformed.stdout`
  carry `rows[1]{story,seam,red_signal}` on line 3, so the column-list mutation
  already goes red today. Those five files are hand-edited to the new header as
  part of the projection change — never regenerated by running the code, which
  would make them agree with any implementation. Separately,
  `coverage_test.go`'s `wantTable` helper builds its expectation by calling
  `toon.Table` with the implementation's own column list; it follows a projection
  change silently, so `TestCommand`'s projection assertions move to a literal
  expected block. The goldens are the red; fixing `wantTable` stops a second,
  tautological assertion from masking a future change.

- **`ParseSpec` is untouched in shape.** It keeps returning opt-in, ordered row
  IDs, and violations; `internal/preflight` is not edited.
- **Expand only; contract is a later spec.** Both existing headers stay accepted.
  The in-tree test fixtures that embed a coverage map migrate to the reduced
  header in this change, so the contract that eventually deletes the six-column
  branch is a parser-and-two-specs edit rather than a fixture sweep. This spec's
  own map is authored six-column because that is what the tree enforces today.
- **The template moves, so its anchors move with it.** Every anchor needle that
  lives inside the template block — the acceptance-coverage vocabulary, the seam
  diagram, `tests attach here`, `Won't handle`, `Status: staged`, `why it catches
  the failure`, the header line, and the approval paragraph — retargets from
  `.agents/commands/bench-write-spec.md` to
  `.agents/skills/bench-craft-spec/SKILL.md`, with the sections named. Anchors
  whose needle is retired outright (the `red signal` requirement, the two-round
  sentence, the edge-class disposition requirement) are deleted with their canary
  fixtures; anchors whose needle is reworded (the `--reviewer` grammar, the
  falsification-question line, `craft-tdd`'s red-signal-ownership sentence, and
  the profile's loop-1/loop-2 routing) are rewritten in place.
  `.agents/commands/bench-write-spec.md` carries **47** anchored registry rows
  today (34 `Require`, 8 `RequireInSection`, 4 `Forbid`, 1 `ForbidInSection`), all
  live and green — the "twenty" this spec carried at staging was a miscount, not a
  stale inventory. A `BASE`/`MUTATE.json` fixture that mutates a moved needle
  follows it; payload-only fixtures are hand-written stand-ins, not copies of the file.
  The shrink dispositions **every one** — retained in the budgeted file,
  retargeted to `craft-spec` with its section named, or retired with its canary —
  and the ticket's acceptance is the full enumeration, not a sample. `bench
  anchors` green over the tree is the check that no needle was left behind.
- **Where the surviving prose lands (reviewer-closed 2026-08-16).** Only seven
  needles live exclusively inside the moved template block; 31 more live in prose
  that stays behind, totalling ~3,300 characters against a line budget that
  counts physical newlines. They do not all fit, and the resolution is a split by
  ownership rather than deletion — nothing is retired for being surplus:
  - The command keeps the entry contract, ownership, the exit handoff, the
    spec-retire lifecycle, and the stale-command sweep. That is ~3,300 characters
    of verbatim needle once the retained anchors are counted — the map's 60-line
    estimate rested on a ~1,100 miscount — and it lands at 73 lines, the budget the
    reviewer accepted at build time (2026-08-17).
  - The **review rubric** — the materiality exit, the cheapest-plausible degenerate
    standard, and the falsification questions — moves to `craft-spec`, which already
    owns the process. It is rubric, not phase choreography.
  - Three sentences are **two-loop residue** and are reworded to one round rather
    than relocated: the slicing step ("After loop 1 accepts … then run loop 2"), the
    learnings hook ("When *either loop* takes more than one iteration"), and the
    `Verification log: spec <n> + tickets <m>` two-count schema. Left alone they
    survive as anchored false claims — the defect story 31 exists to prevent.
  - `craft-spec` therefore overruns the 120-line `*/SKILL.md` glob and takes an
    exact budget row of its own, resolving the headroom question the Further notes
    left open.

- **Confirmed rather than changed.** `projects/benchkit.md` already routes both
  write-spec loops to **mid model, high effort, read-only and same-family through
  the harness's native agent surface**. Map #12 and #15 therefore confirm the
  existing default; the delta is the loop *count* and the removal of
  `--reviewer`'s cross-family branch, not the tier or the venue. FT210's
  cross-family Sol rounds came from an explicit override.
- **No bootstrap-authority claim.** `craft-spec`'s
  `Bootstrap authority before execution` rule applies to a trusted-execution or
  refusal-before-execution claim; this change makes none — the parser reads
  markdown and the prose edits carry no executable hop — so the rule is
  discharged as not applicable rather than skipped.

## Testing decisions

- A good test drives the public interface: `Check`, `Rows`, `ParseSpec`, and
  `Command` over spec bytes, asserting the violation strings and the rendered
  TOON a caller actually receives — never the unexported schema descriptor.
- **Engineering seam: `internal/coverage`.** Prior art is `coverage_test.go`,
  whose `TestCheck` pins every validation phrasing because they are matched by
  substring downstream, and whose `TestCommand` asserts the rendered output.
  Reduced-schema cases join those tables rather than opening a parallel file.
- **Engineering seam: the gate's docs layer.** The prose half is observed by
  `bench anchors` over the real tree and by the `workflow-guidance-anchors` and
  `coverage-map-validation` canary families, each fixture failing for its own
  planted reason. Prior art is every existing fixture under those directories.
- **Gate seam.** `bench gate` observes both: the `test` phase runs
  `internal/coverage`, and the `docs` layer runs the anchors, the canary
  inventory, `bench skills-index --check`, and the guidance-prose-budget check
  that reads the profile's table.

### Seam diagram

    trigger: bench coverage [--check] <spec> · bench gate · internal/preflight
        │
        ▼
    spec.md bytes ──▶ [ internal/coverage: parse → schema → Check / Rows ] ──▶ violations, rows[N]{story,behavior,seam}
                          ◀ tests attach here: build spec bytes in-test, call Check/Rows/Command, assert strings

    trigger: a session loading guidance
        │
        ▼
    .agents prose ──▶ [ anchors registry + canary fixtures + budget table ] ──▶ green or a named diagnostic
                          ◀ tests attach here: bench anchors over the tree; each canary mutates one needle

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| SR1 | 1 | `row\|story\|behavior\|seam\|why it catches the failure` parses as mapped with row IDs opted in | `internal/coverage` | not observed: `TestReducedSchemaHeaders` will assert `State` is `mapped` and `ParseSpec` returns opt-in true for this header | a parser that ignores the reduced header returns `no-map`, so the spec is reported unmapped |
| SR2 | 2 | `story\|behavior\|seam\|why it catches the failure` parses as mapped with row IDs absent | `internal/coverage` | not observed: `TestReducedSchemaHeaders` will assert `State` is `mapped` and `rowIDs` is nil for this header | a four-cell reduced header left unrecognized makes the row-ID opt-out unusable |
| SR3 | 3 | both legacy headers still parse as mapped with their existing opt-in verdicts | `internal/coverage` | already covered: `TestCheck`'s `valid` fixture asserts zero violations for the legacy five-column header (dropping that case from the switch makes it report `missing the canonical header`), and `TestParseSpecOptIn` asserts opt-in true for the six-column header | a switch rewritten rather than extended silently unmaps every staged spec |
| SR4 | 4 | a five-cell header naming `red signal` selects the legacy schema and one omitting it selects the reduced schema | `internal/coverage` | not observed: `TestReducedSchemaHeaders` will assert the two five-cell headers produce different `Rows` outputs from identical row bytes | a count-based choice reads the wrong column for one of the two five-cell forms |
| SR5 | 5 | a reduced-map row with the wrong cell count reports `has N cells (want 5)` | `internal/coverage` | not observed: `TestCheck` gains a reduced-schema case asserting the `(want 5)` clause | a width taken from the legacy schema tells the author to add a column that no longer exists |
| SR6 | 6 | an empty cell in a reduced map names a reduced field | `internal/coverage` | not observed: `TestCheck` gains a case asserting `empty 'seam' cell` for a reduced map's fourth cell — index 3, the only offset where the reduced and legacy five-cell field lists disagree | a field list taken from the legacy schema names `red signal` at index 3, a column the author never wrote; asserting the last cell would not catch it, since index 4 is `why it catches the failure` under both lists |
| SR7 | 7 | a `;` in a reduced map's behavior cell is refused | `internal/coverage` | not observed: `TestCheck` gains a case asserting the `states more than one predicate` phrasing on a reduced map | a hard-coded behavior offset reads the seam cell instead, so a two-predicate behavior passes |
| SR8 | 8 | story references, the four-story fan-out bound, duplicate and malformed row IDs, and orphan stories are refused identically under both schemas | `internal/coverage` | not observed: `TestCheck`'s existing cases run against both schemas from one table, asserting identical violation strings | a check reading a literal index passes on one schema and not the other, so cutting a column cuts a check |
| SR9 | 9 | `bench coverage <spec>` renders `rows[N]{story,behavior,seam}` for a legacy spec | `internal/coverage` | observed red: `internal/coverage/testdata/pre-disclosure-mapped.stdout` line 3 is `rows[1]{story,seam,red_signal}:`, so `TestCommandPreservesCheckedInPreDisclosureResponses` fails until the golden is hand-edited; `TestCommand`'s own assertion moves to a literal block | a projection left on story/seam/red-signal feeds task seeding a column half the specs lack |
| SR10 | 9 | `bench coverage <spec>` renders the same `rows[N]{story,behavior,seam}` shape for a reduced spec | `internal/coverage` | not observed: `TestCommand` gains a reduced-spec case asserting the same literal TOON header | two projections by schema make a caller branch on the spec's schema |
| SR11 | 10 | a reduced spec with violations still renders one `retry after repairing coverage map` action, and a clean one still renders a per-row action | `internal/coverage` | not observed: `TestCommand`'s reduced case asserts both action forms | an action list rebuilt with the projection loses the remedy an agent reads |
| SR12 | 11 | `ParseSpec` returns unchanged opt-in, ids, and violations for a six-column spec | `internal/coverage` | already covered: `TestParseSpecOptIn` asserts opt-in true and the exact ordered id slice for the six-column fixture | a signature or ordering change breaks `internal/preflight`'s fence verdict without touching that package |
| SR13 | 12 | a header renaming only some cells reports `coverage map missing the canonical header` | `internal/coverage` | not observed: `TestCheck` gains a half-renamed header case asserting that phrasing | a prefix or fuzzy match parses a half-migrated table under the wrong schema |
| SR14 | 13 | a reduced map with a header and no data rows reports `coverage map has no data rows` | `internal/coverage` | not observed: `TestCheck` gains a reduced empty-table case | a schema branch that returns early on an unknown width reports a valid map |
| SR15 | 14 | a spec with no map and no historical marker is still a violation | `internal/coverage` | already covered: the `coverage-map-validation/no-map-not-historical` canary reproduces `coverage map missing and spec is not marked historical` | a new schema path that defaults to "mapped" makes an unmapped spec pass |
| SR16 | 15 | the historical marker skips validation under a reduced header | `internal/coverage` | not observed: `TestStateAndRows` gains a reduced historical case asserting `historical` | a marker check bound to the legacy header grades an opted-out spec |
| SR17 | 16 | the six Go test files outside `internal/coverage` that embed a coverage map use the reduced header, and their suites stay green | `internal/preflight`, `internal/worktree`, `internal/systemtest`, `cmd/bench` | not observed: those suites run unchanged against migrated fixtures | a fixture left six-column means the reduced header is never exercised through a real consumer |
| SR18 | 17, 18, 19 | `bench-write-spec.md` is at most 73 lines and the profile's budget table carries its row | gate docs layer | not observed: the `guidance-prose-budgets` check will fail on the 190-line file the moment its row is added | without the row the shrink is unenforced and re-accretes |
| SR19 | 20, 21 | `bench-write-spec.md` states one review round over the spec-and-tickets pair with the ticket quiz as its approval step | gate docs layer | not observed: the two-round anchor is replaced by a one-round needle, and its canary fixture mutates the new sentence | a reworded round paragraph that keeps two loops leaves the cost the change exists to cut |
| SR20 | 22 | the `--reviewer` anchor names a tier-only grammar resolving same-family through the harness's own column, with no cross-family route and no bound-model-id form | gate docs layer | not observed: the reworded anchor's canary fixture drops the tier-only clause and expects the diagnostic | a needle left naming the cross-family recipe or the model-id form keeps a removed branch documented |
| SR21 | 23, 24 | `craft-spec` carries no red-signal grammar and `craft-tdd` names `already covered` and `not TDD-able` | gate docs layer | not observed: a new `ForbidInSection` anchor on `craft-spec` reports a diagnostic while the four-form grammar remains; the existing `not TDD-able` requirement on `craft-tdd` holds the other half and no canary plants it today | deleting the grammar from both files at once leaves no owner for the classification |
| SR22 | 25, 26, 27 | the canonical edge-class run lives in `craft-tdd` and `**Won't handle**` survives for reviewed exclusions | gate docs layer | not observed: `bench anchors` reports both craft-spec edge-class needles — the plain `re-run idempotency` Require and the `RequireInSection("The edge inventory")` class run — until both retarget to `craft-tdd`; a new canary fixture is created for the moved section, since none plants it today | moving the walk without moving its anchor drops the class list with no diagnostic |
| SR23 | 28 | the acceptance-row clause in `craft-tickets` is anchored and its canary reproduces its own diagnostic | gate docs layer | not observed: no anchor names that clause today, so deleting it is silent; the new anchor's canary plants the deletion | an unanchored sizing clause regresses to the 17-versus-5 split with no diagnostic |
| SR31 | 31 | `craft-tdd`'s `row schema and the red-signal definition are` needle names the definition's new owner and still bites | gate docs layer | not observed: the reworded needle's canary mutates the new sentence; leaving the needle unchanged keeps an anchored sentence that points at deleted prose | an anchor guarding a false sentence enforces the wrong thing while reading green |
| SR24 | 30 | `independently-green implementation tickets` still appears in README and CHANGELOG | gate docs layer | already covered: the `changelog-ticket-vocabulary` canary reproduces `CHANGELOG.md dropped the decision-ticket phase-ownership change`; README's needle is anchor-only, so `bench anchors` reports its diagnostic on removal and no canary plants it | a vocabulary sweep that removes the phrase erases the decision-ticket distinction |
| SR25 | 32, 33 | every retargeted anchor names a file that contains its needle, and `bench anchors` is green over the tree | gate docs layer | not observed: `bench anchors` reports each moved needle's diagnostic until its row is retargeted | an anchor left pointing at the old file passes vacuously or fails for the wrong reason |
| SR26 | 33 | each touched canary fixture still fails for its own planted reason | gate docs layer | not observed: the canary suite fails if a fixture stops reproducing its `EXPECT` line | a fixture updated to match new prose without re-planting its reason proves nothing |
| SR27 | 34, 35 | `craft-delegate`'s charge names behavior and seam, and `craft-review` and `/bench-implement-spec` name the reduced projection | gate docs layer | not observed: `bench anchors` and the docs-currency sweep run over the reworded files | a consumer describing a column the producer no longer emits sends a delegate hunting a missing cell |
| SR28 | 36 | `CONTEXT.md` defines **coverage row** and **acceptance row** with Avoid lists | gate docs layer | not observed: the glossary-only structure check runs over the added terms | a term resolved in the spec but absent from the glossary is re-invented per session |
| SR29 | 37, 38 | the field guide and the distillation doc describe the reduced schema, and CHANGELOG carries an entry | gate docs layer | not observed: the docs-currency canary family runs over the edited docs | shipped docs teaching a retired column are worse than no docs |
| SR30 | 39 | `bench skills-index --check` is green after the skill descriptions change | gate docs layer | not observed: the check fails while an edited description and its index row disagree | an index left stale sends a session to a skill whose description no longer matches |

Not covered: story 29 — an accounting statement about landed work, recorded in the implementation decisions; no behavior to observe.
Not covered: story 40 — the six-column branch is deliberately retained; its deletion is priced in Out of scope.
Not covered: story 41 — the right-size table is unedited, so no row can observe a change that does not happen.
Not covered: story 42 — `/bench-implement-spec`'s worktree choreography is unedited beyond the projection sentence SR27 covers.
Not covered: story 43 — `references/cross-harness-reviewers.md` is unedited; its existing canaries already observe it.

### Edge inventory

The generic classes and the profile's hostile-input checklist, walked against the
parser (the prose half has no input surface of its own — it is graded, not
parsed):

- **Error path** — spec unreadable or absent: SR15's sibling paths are unchanged; `Command` already returns `spec not readable` / `spec not found`. **Won't handle** as a new row — no reduced-schema branch touches file resolution.
- **Empty / absent input** — header with no data rows: SR14. Spec with no map: SR15.
- **Boundary values** — the four-story fan-out bound at exactly `bounds.CoverageRowStories` and one past it: SR8.
- **Malformed input** — half-renamed header: SR13. Wrong cell count: SR5. Malformed row ID: SR8.
- **Interrupted / partial state** — not applicable: the parser is a pure read over a byte slice with no partial write.
- **Re-run idempotency** — `bench coverage --check` is read-only and reports the same verdict on repeat. **Won't handle** — no row: the command performs no write, so the profile's self-falsifying-write class cannot arise.
- **Process-boundary lifecycle** — `Command` returns a string and a code to one caller; unchanged. **Won't handle** — no process boundary moves in this change.
- **Hostile environment** — control bytes in a spec *path*: already covered by `TestCommandControlBearingSpecPathPreservesPrimaryAndHonestFallback` and `TestCommandAngleBracketSpecPathPreservesPrimaryAndHonestFallback`, both unchanged by the schema. Control bytes in a *cell* reaching the TOON sink: the `behavior` cell newly reaches `toon.Table` through the projection, where `red_signal` used to — SR9 and SR10 assert the rendered table, and `toon.Table` refuses unrepresentable bytes rather than rendering them.
- **Hostile-input checklist — no trailing newline**: a spec whose last line is the final map row. **Won't handle** — no row: `parse` splits the whole file on `\n`, so a final unterminated line is yielded as a line under either schema; the diff changes no part of that path.
- **Hostile-input checklist — absent vs present-but-empty file**: unchanged from today's resolution path, which this change does not touch.
- **Hostile-input checklist — special files and dangling symlinks in the spec sweep**: unchanged; `specref.Resolve` owns that classification and is not edited.
- **Hostile-input checklist — escaped pipes in a cell** (`\|`): the parser's escape sentinel must survive the schema change, since a behavior cell legitimately contains `|`. Covered by SR8's shared table, whose fixtures carry `does x \| y`.

## Ownership fences

- `specs/spec-ticket-fence-reduction/spec.md`
- `specs/spec-ticket-fence-reduction/tickets/`
- `reviews/spec-ticket-fence-reduction.md` (transient review pickup; deleted with its repairs)
- `capture/session-handoff.md`
- `internal/coverage/coverage.go`
- `internal/coverage/coverage_test.go`
- `internal/coverage/testdata/`
- `internal/preflight/command_bootstrap_test.go`
- `internal/preflight/command_review_test.go`
- `internal/preflight/command_harness_test.go`
- `internal/worktree/land_test.go`
- `internal/systemtest/owner_test.go`
- `cmd/bench/command_registry_test.go`
- `internal/anchors/registry_data.go`
- `internal/conformance/fixture_bite_test.go` (the four `TestSpecTicketHandoffWorkflowFixturesAreComplete` diagnostic literals that mirror retargeted anchors, and nothing else)
- `.agents/commands/bench-write-spec.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `.agents/skills/bench-craft-tdd/SKILL.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-review/SKILL.md`
- `tests/canary/workflow-guidance-anchors/`
- `tests/canary/coverage-map-validation/`
- `projects/benchkit.md`
- `.bench/BENCH-reference.md`
- `CONTEXT.md`
- `CHANGELOG.md`
- `README.md`
- `docs/field-guide.html`
- `docs/reporesident-distillation.md`

## Out of scope

- **Contract the six-column branch** — delete the two legacy header cases, their
  field descriptors, and the six-column canary fixture. `worktree-landed-retirement`
  retired at `64bca130`, so this spec's own map is now the only staged
  six-column map in the tree: the contract unblocks as soon as this spec reaches
  `Status: implemented` and its map is migrated or retired with it.
  4 edits, 2 gate runs.
- **Move the right-size table's light-path threshold** (map #16). 2 edits, 1 gate
  run.
- **Change `/bench-implement-spec`'s worktree, preflight, or landing
  choreography** (map #1). Unpriced — a separate capability with its own spec.
- **Delete `references/cross-harness-reviewers.md`** (map #13). 4 edits, 1 gate
  run.

## Further notes

Three items for reviewer veto rather than assumption:

- **Story 2 extends map #5.** The map fixes the reduced header as exactly
  `| row | story | behavior | seam | why it catches the failure |`. Story 2 also
  accepts the four-cell form without `row`, extending the existing row-ID opt-out
  convention to the reduced schema. Flagged as a non-behavioral extension of the
  source; veto it and story 2, SR2, and one accept-ticket bullet come out.
- **`craft-spec`'s budget headroom.** It is 71 lines under the 120-line
  `*/SKILL.md` glob. Taking the template while shedding the grammar and the edge
  walk lands it near 110 — inside the glob but thin. Adding an exact budget row is
  a "what the gate proves" call, so the build will surface the landed line count
  rather than add a row on its own.
- **Map #18** (closing the two `capture/learnings.md` entries against this map)
  carries no coverage row because it is a `/bench-what-next` action outside this
  diff.

`bench maps` reports `decisions/spec-build-review-gate-cadence.md` invalid — its
`Sources` names `internal/specbuild/checkpoint.go`, which no longer exists. It is
unrelated to this change and repointing it needs a reviewer call on what replaced
that seam, so it stays parked rather than folded in here.
