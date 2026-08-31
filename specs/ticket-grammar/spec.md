# Ticket-file grammar

Status: implemented

Roadmap: FT174

Decision source: roadmap/FT174.md (named reviewed artifact); the reviewer confirmed the full-mechanization enforcement shape in the 2026-08-31 spec session.

Verification log: 2 iteration(s) to accept — eight blocking findings and one merge folded after round one; the scoped re-review accepted with three prose fixes, folded.

## Problem

No parser reads ticket files. Preflight scrapes bare `[A-Z]+[0-9]+` tokens from
ticket prose, so a row ID in a sentence counts as ownership. `Blocked by:` and
`Writes:` are unparsed text. A dangling blocker, a duplicate field, a typo
path, or an omitted fixture registry reaches a delegate charge unnamed. Four
recorded builds paid repair rounds for exactly these faults.

## Solution

One parser owns the ticket-file schema. It enforces three grammars: dependency
(`Blocked by:`), ownership (`Writes:` with three closures), and mutation
(`Covers:` row citations). Preflight replaces the token scrape with parsed
rows. A new conformance check sweeps every staged spec's tickets inside the
gate, and a canary family proves each diagnostic bites. The parser composes
the field grammar `internal/maps` already owns; it writes no second scan.

## User stories

Group 1 — one ticket parser. Line: opus / medium. The parser is oracle-adjacent conformance logic, which routes mid effort.

1. As a build coordinator, I want one parser to read every ticket file, so that no phase scrapes bare tokens.
2. As a spec author, I want the required fields enforced, so that a half-written ticket is named
  before the charge.
3. As a reviewer, I want a duplicate field line rejected, so that two `Blocked by:` values cannot disagree silently.
4. As a reviewer, I want fenced code skipped, so that a quoted example cannot parse as live grammar.
5. As a maintainer, I want the decision-map parser and the ticket parser to share one field grammar, so that the two grammars cannot drift.
6. As a maintainer, I want the decision-map diagnostics unchanged by the lift, so that existing maps keep their verdicts.
7. As an operator, I want a special file or an unreadable ticket refused by name, so that a
  broken entry never reads as green.

Group 2 — dependency grammar. Line: opus / medium. The graph rules are small compositions of lifted logic.

8. As a coordinator, I want `Blocked by:` to hold `none` or sibling ticket basenames, so that the frontier derives from files alone.
9. As a coordinator, I want a dangling blocker named, so that a retitled sibling cannot orphan an edge.
10. As a coordinator, I want a duplicate blocker and a self-edge named, so that a careless list cannot pass.
11. As a coordinator, I want a blocker cycle named with one edge, so that a build cannot start against an empty frontier.

Group 3 — ownership grammar. Line: opus / medium. The closures compose existing enumerations and stay mid effort.

12. As a delegate, I want every `Writes:` entry to exist or carry `(new)`, so that a typo path is
  caught before the charge.
13. As a reviewer, I want a ticket that writes a fixture-pinned path to name that fixture, so that the red-capable fixture rides in the ticket.
14. As a reviewer, I want a ticket that writes a bound command package to name its fixture registries, so that no registry is found mid-build.
15. As a maintainer, I want the command-to-registry binding gate-checked for existence and completeness, so that the binding cannot rot.
16. As a delegate, I want a ticket that writes a system-tag test to pin `BENCH_KIT`, so that an ambient kit cannot flip the fixture.
29. As a reviewer, I want the dispatcher, renderer, and terminal-lifecycle owners bound in the same table, so that nobody names them by hand.

Group 4 — mutation grammar. Line: opus / medium. The citation reads move from scrape to structure at a known seam.

17. As a preflight consumer, I want row ownership read from a `Covers:` field, so that a row ID in prose is not evidence.
18. As a preflight consumer, I want every declared row covered by some ticket's `Covers:`, so that an unowned row stops the phase.
19. As a preflight consumer, I want every `Covers:` token to name a declared row, so that a phantom or
  foreign token stops the phase.
20. As a reviewer, I want a repeated `Covers:` token in one ticket named, so that a padded list cannot inflate coverage.
21. As a light-path author, I want the `Covers:` checks skipped in a tickets-only folder, so that a spec-less ticket needs no phantom rows.

Group 5 — enforcement venues. Line: opus / medium. Gate and conformance logic routes mid effort.

22. As a phase runner, I want the grammar reds as rows in `bench preflight build` and `bench preflight review`, so that a red stops the phase at entry.
23. As a lander, I want the gate to sweep every staged spec's tickets, so that a malformed ticket cannot land green.
24. As a lander, I want the sweep to tolerate an absent tickets directory, so that a staged spec before slicing stays green.
25. As a maintainer, I want canary fixtures that red a real gate run for the new diagnostic classes, so that the sweep provably bites.

Group 6 — guidance. Line: opus / high. Guidance prose compounds through every session that loads it.

26. As a ticket author, I want the craft-tickets template to advertise the enforced grammar, so that a cold author writes a passing ticket.
27. As a ticket author, I want the template's example to parse clean under the new grammar, so that a copied example cannot red.
28. As a reviewer, I want the guidance to stay within its prose budget, so that the skill stays a bounded read.

## Implementation decisions

- A new package `internal/tickets` owns the ticket-file schema: title,
  `Blocked by:`, `Writes:`, `Covers:`, `## What to build`, `## Acceptance`,
  and an optional `## Delegate charge`. It is a pure decision domain in the
  preflight style: immutable bytes in, diagnostics out, no I/O in the core.
- The generic field scan lifts out of `internal/maps` into one shared exported
  form. The scan covers the fence skip, the prefix fields, the duplicate-field
  diagnostics, the blocker-list split, and the duplicate/dangling/self/cycle
  graph walk. Both
  schemas drive it. The maps package keeps its diagnostics byte-identical.
- `Covers:` is the one citation source. The value is `none` or a
  comma-separated row-ID list. In a folder with `spec.md` the field is
  required; in a tickets-only folder the mutation checks are skipped. The
  bare-token scrape (`tokenRe`, `Facts.TicketTokens`) deletes. The spec tag
  still derives from the declared row IDs.
- Preflight gains one verdict row per grammar: `tickets-parse`,
  `blockers-resolve`, `writes-resolve`, `fixture-closure`, `registry-closure`,
  and `kit-pin`, ordered after `paths-authorized` and before `rows-owned`.
  `rows-owned` and `rows-membership` read parsed `Covers:` tokens. Decide stays
  pure: the gatherer collects the existence bits, the fixture-pin map, and the
  constraint tags into Facts. Build mode
  without a tickets directory keeps today's not-applicable posture. An
  unreadable directory or special file stays a bootstrap failure; a grammar
  fault is a named verdict row.
- The ticket enumeration stays recursive with the lstat-first refusal at
  every depth. Only `.md` files parse as tickets, and a non-`.md` file is
  ignored. A basename that appears at two depths yields a duplicate-identity
  diagnostic, because blocker edges resolve by basename.
- The ownership closures derive from producers, never from copied lists.
  The fixture closure enumerates the live canary inventory. A fixture pins a
  path when its `BASE` list, its `files/` overlay, or its `MUTATE.json` names
  that path. The `BASE` walk follows `@` includes. `internal/canary` exports
  that enumeration.
- The kit-pin rule parses the written test file's build constraints, as the
  coverage citations already do. When the system tag is present, the rule
  requires the literal `BENCH_KIT` in the ticket body.
- The command-to-registry binding is declared data in `internal/tickets`,
  in the shape `internal/anchors/registry_data.go` set. Each row binds a package prefix to
  the files a ticket must co-name. The command rows include the help
  projection and the envelope cases. Seed rows also bind the dispatcher
  (`cmd/bench`), the renderer (`internal/toon`), and the terminal-lifecycle
  owner (`internal/terminal`) to their pinned test files. The
  new conformance check proves every bound file exists and that every AXI
  query command's package has a row. The check also proves the three owner
  rows are present.
- A new registered conformance check `ticket-grammar` (dev tier) sweeps every
  `specs/*/tickets/` directory with the same parser, fail-closed. Its canary
  family is `tests/canary/ticket-grammar`. The check joins the profile's
  input-binding table. `bench test --check ticket-grammar` reaches it through
  the existing check surface; no new command ships.
- The landing keeps its narrower final authorization: `bench worktree land`
  still consumes `paths-authorized` alone, and the destination's gate sweep
  grades ticket grammar.
- No staged spec exists in the tree today, so no landed ticket needs a
  migration.

## Testing decisions

- A good test drives observable behavior. It feeds a ticket-file fixture and
  observes a named diagnostic, a verdict row, or a gate red.
- Seams and prior art: `internal/tickets` table tests follow
  `maps_parse_test.go`; preflight rows follow `decision_test.go` table tests;
  the sweep follows the docs-currency check tests; the canary family follows
  the fixture-bite proof architecture.
- The gate observes the feature through the registered `ticket-grammar`
  conformance check inside the test phase. The canary family proves each
  diagnostic class through a real gate run.

### Seam diagram

    bench preflight build|review <slug>          bench gate
        │                                            │
        ▼                                            ▼
    specs/<slug>/tickets/*.md ──▶ [ internal/tickets parser ] ──▶ diagnostics
        │                                 │
        ▼                                 ▼
    [ preflight Decide rows ] ──▶ TOON verdict table
                                      ◀ tests attach here: table tests feed
                                        immutable facts; the conformance check
                                        sweeps staged specs; canary fixtures
                                        red a real gate run

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| TG1 | 1, 2 | the parser reports one named diagnostic for each absent required field | tickets parser table test | a token scraper reports nothing on a field-less ticket |
| TG2 | 3 | a second `Blocked by:` line in one ticket yields a duplicate-field diagnostic | tickets parser table test | last-wins parsing lets two values disagree silently |
| TG3 | 4 | a field line inside a fenced code block parses as no field | tickets parser table test | a prefix grep reds the template's own quoted example |
| TG4 | 5 | the maps package and the tickets package drive the field scan through one lifted shared symbol | review-owned (diff inspection) | a copied scan drifts on the next grammar change |
| TG5 | 6 | one golden test asserts the exact ordered diagnostic slice over the maps test inputs across the lift | `TestDecisionMapDiagnosticsGolden` in the maps package, added green before the lift | a presence-only assertion lets the lift reorder or drop a verdict |
| TG6 | 7 | a FIFO, a device, or a dangling symlink under tickets/ yields a named refusal | preflight gather table test | a plain read blocks on the file or reports a broken link as empty |
| TG7 | 8 | `Blocked by: none` parses to zero dependency edges | tickets parser table test | a literal-minded parser reports `none` as a dangling basename |
| TG8 | 9 | a blocker that names no sibling file yields a dangling-blocker diagnostic naming both basenames | tickets parser table test | a retitled sibling orphans the edge with no signal |
| TG9 | 10 | a repeated blocker basename yields a duplicate-blocker diagnostic | tickets parser table test | silent deduplication hides the authoring fault |
| TG10 | 10 | a ticket that names itself as a blocker yields a self-edge diagnostic | tickets parser table test | a self-edge removes the ticket from every frontier with no signal |
| TG11 | 11 | a blocker cycle yields a cycle diagnostic that names one edge | tickets parser table test | a cycle leaves an empty frontier the coordinator debugs by hand |
| TG12 | 12 | a `Writes:` entry absent from the tree without the `(new)` marker yields a diagnostic naming the path | preflight decision table test | a typo path charges a delegate against nothing |
| TG13 | 12 | a `Writes:` entry marked `(new)` and absent from the tree stays green | preflight decision table test | an unconditional exists check reds every new-file ticket |
| TG14 | 13 | a ticket that writes a fixture-named path reds unless `Writes:` also names that fixture directory | preflight decision table test with a planted fixture | the red-capable fixture stays outside the ticket and the bite breaks unnoticed |
| TG15 | 13 | the fixture-pin enumeration derives from the live canary loader | `TestFixturePinsEnumeratesLiveInventory` plants a synthetic fixture and observes the new pin | a copied fixture list goes stale the day a fixture moves |
| TG16 | 14 | a ticket that writes a bound package reds unless `Writes:` names every file bound to that package | preflight decision table test | an AXI ticket names five of eight registries and pays repair rounds |
| TG17 | 15 | a binding row that names a nonexistent registry file reds the gate check | conformance check test over a synthetic root | the binding rots and the closure enforces a fiction |
| TG18 | 15 | an AXI query command package with no binding row reds the gate check | conformance check test | a new verb ships with no registry closure at all |
| TG19 | 16 | a written Go test file whose build constraint carries the system tag reds unless the ticket body states `BENCH_KIT` | preflight decision table test | an ambient kit flips the fixture verdict under composition |
| TG20 | 17 | a spec-tag row ID in ticket prose lands in no parsed `Covers:` set | tickets parser table test | the bare-token scrape returns and prose mentions read as evidence |
| TG21 | 18 | a declared row in no ticket's `Covers:` reds rows-owned naming the row | preflight decision table test | an unowned row builds nothing and nobody notices |
| TG22 | 19 | a `Covers:` token that names no declared row reds rows-membership naming the token | preflight decision table test | a phantom row ID claims coverage that does not exist |
| TG23 | 19 | a `Covers:` token with a foreign tag yields a diagnostic | tickets parser table test | today's tag scoping silently ignores the stray token |
| TG24 | 20 | a repeated token in one ticket's `Covers:` yields a duplicate diagnostic | tickets parser table test | a padded list inflates apparent coverage |
| TG25 | 21 | a tickets folder with no spec.md is graded with the `Covers:` checks skipped | conformance check test | requiring rows without a map reds every light-path landing |
| TG26 | 22 | `bench preflight build` renders the six rows `tickets-parse`, `blockers-resolve`, `writes-resolve`, `fixture-closure`, `registry-closure`, and `kit-pin`, each red with its own detail | preflight command test | a red folded into one shared row loses failure attribution |
| TG27 | 23 | the gate sweep reds a staged spec whose ticket carries a dangling blocker | conformance check test | a malformed ticket lands green and the next session inherits it |
| TG28 | 24 | a staged spec with no tickets directory stays green in the sweep | conformance check test | the sweep reds the spec phase before slicing can run |
| TG29 | 25 | every fixture in the ticket-grammar canary family reds a real gate run with its EXPECT diagnostic | fixture-bite proof over the family | the diagnostics exist in tests and never reach the oracle |
| TG30 | 26, 27 | the craft-tickets example between the ticket-example markers parses clean through the live parser | ordinary test feeding the marked example to the parser | the copied starting point reds the first ticket an author writes |
| TG31 | 28 | the craft-tickets skill stays within its guidance-prose budget row | guidance-prose-budgets check | the grammar advertisement bloats the skill past a cold read |
| TG32 | Edge | a field-prefix line whose leading bytes carry an NBSP yields the missing-required-field diagnostic | tickets parser table test | a silent non-field read hides the unanchored token |
| TG33 | Edge | a fence opened and never closed yields an unterminated-fence diagnostic | tickets parser table test | a swallowing parser grades nothing after the opening marker |
| TG34 | Edge | a required field on the last line without a trailing newline parses | tickets parser table test | a hand-edited ending reds a well-formed ticket |
| TG35 | Edge | a control byte in a diagnostic path is refused before the verdict table renders | preflight command test | TOON cannot represent the cell and the verdict dies mid-render |
| TG36 | 29 | the binding table carries rows for the dispatcher, the renderer, and the terminal-lifecycle owners | conformance check test | the advisory fence omits established owners and the coordinator reconstructs them by hand |
| TG37 | Edge | a non-`.md` file under tickets/ is ignored by the grammar | preflight gather table test | a stray asset reds as a malformed ticket |
| TG38 | Edge | two ticket files with one basename at different depths yield a duplicate-identity diagnostic | preflight gather table test | blocker edges resolve against an ambiguous name silently |
| TG39 | Edge | in build mode with no tickets directory the six grammar rows render not-applicable | preflight decision table test | a silent green and a graded green become indistinguishable |

### Edge inventory

The hostile-input walk for this shell-CLI domain, and the deliberate
exclusions:

- Special files, dangling symlinks, and unreadable entries: TG6 keeps the
  lstat-first refusal at every depth.
- Unterminated fence: TG33; the parser reports the region and grades nothing
  after it, fail-closed.
- Non-ASCII whitespace in a field prefix: TG32; the diagnostic stays
  fail-closed rather than false-positive.
- Missing trailing newline: TG34.
- Absent versus empty tickets directory: the existing preflight semantics stay.
  Absent is not-applicable in build mode, and a present empty directory runs
  the row checks.
- Control bytes in paths: TG35 through the existing representability refusal.
- A `Writes:` glob or a comma inside a path: the entry fails the exists check
  and reds through TG12 with the entry named.
- **Won't handle:** a `(new)` entry whose path already exists stays green, because a blocker
  ticket may land the file first. The exists check survives for typo paths.
- **Won't handle:** `Writes:` versus the actual diff. The list stays advisory for what a delegate
  edits, and `paths-authorized` still grades the real diff against the spec
  fence.
- **Won't handle:** cross-folder blocker edges — a ticket blocks only on
  siblings; a build inside one spec survives as the in-scope caller.
- **Won't handle:** `## Delegate charge` content grading — `craft-delegate`
  and review own charge quality; the parser only tolerates the heading.
- **Won't handle:** needle tables pinned inside Go test files — the fixture
  closure reads fixture bytes only. A guidance ticket names
  `internal/conformance/fixture_bite_test.go` by hand, and review survives as
  the in-scope caller.
- **Won't handle:** semantic prose of ticket bodies — the prose-mechanics
  check and review keep grading sentences; field lines stay exempt.
- **Won't handle:** two tickets covering the same row — shared coverage is
  legal today and rows-owned needs one owner, not exactly one.

## Ownership fences

- `internal/tickets` (new)
- `internal/maps`
- `internal/preflight`
- `internal/conformance`
- `internal/canary`
- `tests/canary/ticket-grammar` (new)
- `reviews/ticket-grammar.md` (new)
- `tests/canary/workflow-guidance-anchors`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `projects/benchkit.md`
- `internal/systemtest/owner_land_race_test.go`

## Out of scope

- A `bench tickets <slug>` frontier projection verb — 6 edits, 2 gate runs.
- The FT164 repair-charge template and done-claim owner resolution — 4 edits,
  1 gate run.
- Full unification of the maps and tickets schemas into one schema engine
  beyond the lifted field scan — 8 edits, 2 gate runs.
- Cross-spec blocker edges — 3 edits, 1 gate run.

## Further notes

- The guidance edit must keep every anchor needle the anchors registry pins in
  `bench-craft-tickets/SKILL.md`. The workflow-guidance canary fixtures that
  anchor its clauses ride in the guidance ticket's `Writes:`.
- The craft-review and craft-delegate mentions of an advisory `Writes:` stay
  true: disjointness and diff conformance remain advisory; existence and the
  three closures become enforced.
- The skill file sits at exactly its 100-line budget today, with zero
  headroom. The guidance ticket must cut as much as it adds, or the reviewer
  raises the budget row in the profile.
- The `Covers:` field formalizes the live convention the citation-execution-
  proof tickets already used; preflight reads ids, not ranges, unchanged.
