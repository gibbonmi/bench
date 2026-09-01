# Spec-authoring discipline

Status: staged

Roadmap: FT220, FT257, FT278

Decision source: reviewer-confirmed conversation, 2026-09-01 (fifteen closed predicates)

Verification log: 2 iteration(s) to accept — one round returned revise; the author folded three blocking findings and eleven acceptance findings.

## Problem

Three separate failures cost repair cycles in the spec-authoring phases.

A decision-map author cannot see the grammar the parser enforces. `bench maps --template`
renders an empty `## Sources` section, so the author guesses the record shape. The parser
then refuses a wrapped field with a message that names the field and not the rule. The
template also hides the rule that a resolved decision ticket cannot stay blocked by an
unresolved one.

A spec author slices tickets before the author knows who reads the decision fact. On
the FT216 write-spec the two obvious planners were named and the shared reader was
missed. The first ticket slice bundled an independently shippable consumer with another
one, and that cost a second review round.

A spec author ships a promise without a coverage row. On the path-aware-lane spec review
three source sentences, one destination-side case, and five additions beyond the map had
no row. Round 1 returned revise with seven blocking findings.

## Solution

The kit closes each failure at the place the author already reads.

`bench maps --template` shows one complete `## Sources` entry and states the
resolved-blocked rule. The wrapped-field refusal names the one-physical-line rule.
`/bench-shape-idea` owns the decision-map authoring steps: read one ready map's Sources
block first, then run `bench maps` and `bench gate-prose` on the first skeleton.

`craft-spec` extends its existing sweep sentence into a named **reader sweep**. The sweep
lists each named consumer of the decision fact and each helper a named consumer calls
directly. `/bench-write-spec` sequences the sweep before the `craft-tickets` charge, and
the ship-test question moves down beside it.

`references/map-discipline.md` gains the coverage-row rules. Each in-scope promise takes
one red-capable row. An either-side predicate takes two rows. A flagged-additions list
and a source-sentence-to-row table sit under Further notes.

## User stories

### Group A — the decision-map template and its diagnostics

Line: opus / medium. The group edits the maps parser, the template renderer, and the
refusal text, so it needs the tier that carries Go-seam work.

1. As a decision-map author, I want `bench maps --template` to show one complete Sources
   entry, so that I copy the record shape instead of guessing it.
2. As a decision-map author, I want a `URL:` locator in the template's Sources entry, so
   that it validates without a local file.
3. As a decision-map author, I want Supports and Drift on separate physical lines, so that
   the example teaches me the one-physical-line rule.
4. As a decision-map author, I want the rendered template to validate clean, so that my
   first skeleton is green before I write an answer.
5. As a decision-map author, I want the wrapped-field refusal to name the one-physical-line
   rule, so that I repair the record without the parser.
6. As a decision-map author, I want a second locator line to keep its current refusal, so
   that the recorded fixture stays honest.
7. As a kit maintainer, I want the template to keep one decision ticket, so that every
   single-count answer replacement stays green.
8. As a decision-map author, I want the template to state the resolved-blocked rule, so
   that I do not learn it from a refusal.
9. As a kit maintainer, I want the resolved-blocked enforcement in the graph walk, so that
   the template duplicates no check.

### Group B — coverage-row discipline

Line: opus / medium. The group writes anchored normative prose with canary fixtures, and
medium is the standing implementer effort for guidance prose.

10. As a spec author, I want one red-capable row for every in-scope edge-inventory,
    source-promise, and fence-closure promise, so that no promise ships unobserved.
11. As a spec author, I want a Won't handle line with a surviving caller for each excluded
    edge, so that the cut stays recorded.
12. As a spec author, I want the flagged-additions list under Further notes before the
    first review charge, so that the round can remove additions.
13. As a spec author, I want the source-sentence-to-row table under Further notes before
    the first review charge, so that the round checks each sentence.
14. As a reviewer, I want a row for each listed addition, and no unlisted addition, so that
    promise inflation stops.
15. As a spec author, I want an either-side predicate to produce two rows, one side per
    row, so that neither side hides.
16. As a kit maintainer, I want `bench coverage --check` to stay unchanged, so that the new
    rules stay reviewer-graded prose.
17. As a spec author, I want each canary or conformance row traced to the executed root, so
    that no row is unreachable.
18. As a spec author, I want each named diagnostic state addable or mutable in a fixture,
    so that the canary class can red.
19. As a coverage-map author, I want to quote each pasted operand in the charge, so that the
    delegate reads the operand I read.
20. As a kit maintainer, I want the SKILL.md template heading Further notes to stay bare, so
    that the reference file holds the contract alone.
21. As a kit maintainer, I want the two existing pre-charge rules to stay single-sourced, so
    that this build adds no second copy.

### Group C — the reader sweep and the ship test

Line: opus / medium. The group moves anchored sentences inside a file at its line budget,
so it needs the tier that reads the anchor registry before it edits.

22. As a spec author, I want the craft-spec sweep sentence to name the reader sweep, so
    that I run it on the decision fact.
23. As a spec author, I want the reader sweep to list each named consumer of the decision
    fact, so that no shared reader escapes.
24. As a spec author, I want the sweep to list each helper a named consumer calls directly,
    so that no indirect reader stays unfenced.
25. As a spec author, I want a deeper callee to join only when the callee reads the fact
    itself, so that the sweep stays bounded.
26. As a spec author, I want each shared reader under an exact ownership fence, so that a
    writer cannot edit an unassigned reader.
27. As a kit reader, I want the canonical term "reader sweep", so that "census" still means
    the raw-call census or the architecture census.
28. As a spec author, I want the phase file to sequence the sweep before the craft-tickets
    charge, so that the fences exist before the slice.
29. As a kit maintainer, I want the phase file to sequence the sweep and to restate no rule,
    so that one source holds it.
30. As a spec author, I want the ship-test question immediately before the ticket slice, so
    that I split the graph where I slice.
31. As a reviewer, I want the review-round paragraph to drop the ship-test question, so
    that the move leaves no second copy.
32. As a spec author, I want the ticket graph to split where a consumer branch lands green
    alone, so that the branch ships alone.
33. As a kit maintainer, I want the move to add no net line to the phase file, so that the
    73-line budget holds.

### Group D — decision-map authoring steps and the asset convention

Line: opus / medium. The group edits two phase files and the template together, and the
three copies of the asset path must agree exactly.

34. As a decision-map author, I want to read one ready decision map's Sources block before
    my first write, so that I copy a live record.
35. As a decision-map author, I want to run `bench maps` and `bench gate-prose` on my first
    skeleton, so that I repair the grammar early.
36. As a kit reader, I want no "prose preflight" term in the kit, so that the named verb
    stays the one handle.
37. As a decision-map author, I want map-owned assets under `decisions/assets/`, so that the
    candidate scanner never reads an asset as a map.
38. As a kit maintainer, I want the scanner to keep its current skip behavior, so that the
    convention needs no Go change.

### Group E — the glossary terms

Line: opus / medium. Story 42 audits every sentence this build adds or moves, so the group
needs the tier that reads the whole diff.

39. As a kit reader, I want CONTEXT.md to define "coverage map" with an Avoid list, so that
    nobody writes a bare "map" for it.
40. As a kit reader, I want CONTEXT.md to keep "decision map" with its Avoid list, so that the
    two map terms stay distinct.
41. As a kit reader, I want CONTEXT.md to define "reader sweep" with an Avoid list, so that
    "census" stays reserved.
42. As a kit author, I want new prose to write the two-word terms, so that a bare "map" never
    enters the kit.

## Implementation decisions

**The template gains a Sources example rendered from the schema.** `DecisionMapTemplate`
already renders each terminal section from `canonicalDecisionMapSchema`. The Sources
example joins that render as body text under the Sources terminal heading. The example
uses a `URL:` locator, because `validateSourcePath` resolves a `Path:` locator against
the repository root and no placeholder path exists there. A `URL:` locator routes to
`validSourceURL` instead, so the rendered template validates clean in a temporary root.

**The wrapped-field refusal splits from the unknown-field refusal.** One code path
currently produces both messages. A record line is a wrapped continuation when it holds
no field separator. That case gets the new message, and the message names the
one-physical-line rule.

Every other line keeps the current message word for word. The refusal states the rule for
a Sources record only. Wrapping stays legal in a terminal bullet list. The message must
not claim that wrapping is illegal everywhere.

**The resolved-blocked rule enters the template as a stated line, not as a second check.**
The graph walk stays the one enforcement. The template line describes the rule the walk
already applies.

**The template keeps one decision ticket.** Every consumer resolves the answer placeholder
with a single-count replacement. A second ticket would leave the second placeholder
unresolved and would red a ready-map consumer.

**Each new normative sentence ships as a rule, an anchor, and a canary fixture.** The build
mirrors the FT144 landing shape. The sentence lands in its guidance file. One anchor tuple
in `internal/anchors/registry_data.go` requires or forbids the exact needle. One registry
test proves the anchor reds on removal. One fixture under
`tests/canary/workflow-guidance-anchors/` carries `BASE`, `EXPECT`, and `MUTATE.json`.

**New coverage-row rules land in `references/map-discipline.md`.** The reference carries no
line budget. `craft-spec/SKILL.md` sits at 150 of 152 lines, so only the reader-sweep
extension of the existing sweep sentence lands there. No budget row moves in this build.

**The ship-test question moves as an exact byte sequence, and it carries the split arm.**
An anchor requires the needle "could a narrower capability ship on its own gate" in the
phase file with no section scope. The move keeps that needle byte for byte. The grill
closed the relocated question and the branch-split instruction as one fact. The split arm
therefore joins the same moved sentence rather than a separate line.

**The net-line accounting per budgeted file is exact.** In
`.agents/commands/bench-write-spec.md` the removal frees two lines, and the relocated
question with its split arm takes two lines. That is a net of zero. The sweep sequencing
clause lands as an in-place edit of an existing step sentence,
so it takes no line. The file stays at 73 of 73.

In `.agents/skills/bench-craft-spec/SKILL.md` the reader-sweep extension adds one line,
so the file moves from 150 to 151 of 152. `projects/benchkit.md` stays outside the
fences, so no budget row moves.

**Anchor placement inside one H2 section is not machine-checkable.** `internal/anchors`
scopes a sectioned anchor to an H2 section. The phase file holds the slicing step and the
review-round paragraph under one `## Process` heading. No section-scoped anchor can
express the placement, and a section-scoped forbid would also forbid the required
placement. The Standards axis therefore owns the placement and the no-second-copy
property, and the whole-file require anchor holds the byte-exactness.

**The asset convention stays prose in three files.** `/bench-write-spec`,
`/bench-shape-idea`, and the template each name `decisions/assets/`. The scanner already
lists only the direct children of a decisions directory, so it already skips the nested
asset directory. No Go change follows.

**Two pre-charge rules already exist and stay untouched.** `map-discipline.md` already
states the authoritative-inventory rule for a universal claim and the named-test-function
rule for a row whose seam is the existing tests. This build adds only the pasted-operand
rule.

## Testing decisions

- A good test drives one published behavior. The published behaviors are the rendered
  template text, the validator's diagnostic list, the anchor registry's verdict, and the
  canary fixture's bite.
- The Go seams are the existing `internal/maps` tests. `TestParseDecisionMapSchemaAndTemplate`
  already validates the rendered template. `TestMapSourcesRequireExactRecordShape` already
  pins each record-shape message. `TestMapTerminalContinuationAndEmptyAnswer` already proves
  that a wrapped bullet stays legal in a terminal list.
- The prose seams are `internal/anchors/registry_data_test.go` for the red-on-removal proof
  and `internal/conformance/fixture_bite_test.go` for the fixture bite.
- The gate observes the feature through the `guidance-anchors` check, the
  `guidance-prose-budgets` check, the `decision-map-integrity` canary family, and the
  `workflow-guidance-anchors` canary family.

### Seam diagram

    trigger: bench maps --template, or bench maps over a tree
        │
        ▼
    canonicalDecisionMapSchema ──▶ [ DecisionMapTemplate ] ──▶ template text
                                          │
                                          ▼
                              [ ValidateDecisionMap ] ──▶ ordered diagnostics
                      ◀ tests attach here: internal/maps unit tests render the
                        template and validate it, and canary fixtures bite the
                        diagnostics through bench gate

    trigger: bench gate guidance-anchors phase
        │
        ▼
    guidance file text ──▶ [ anchors registry ] ──▶ diagnostic or silence
                      ◀ tests attach here: registry_data_test.go mutates one
                        needle, and a canary fixture pins the exact diagnostic

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| SAD1 | 1 | The rendered template holds one Sources bullet under the Sources heading | `internal/maps/maps_parse_test.go` (`TestParseDecisionMapSchemaAndTemplate`) | An empty Sources section renders no bullet, so the token assertion fails |
| SAD2 | 2 | The template's Sources bullet starts with a `URL:` locator | `internal/maps/maps_parse_test.go` (`TestParseDecisionMapSchemaAndTemplate`) | A `Path:` locator sends the value to validateSourcePath, which finds no such file and reds |
| SAD3 | 3 | The template's Supports line and Drift line are two separate physical lines | `internal/maps/maps_parse_test.go` (`TestParseDecisionMapSchemaAndTemplate`) | One joined line parses as an unexpected field, so the template diagnostics stop being empty |
| SAD4 | 4 | The rendered template validates with zero diagnostics as shaping and as ready | `internal/conformance/decision_map_integrity_test.go` (`TestDecisionMapIntegrityCheckValidatesEveryCandidate`) | A malformed example reaches every consumer that resolves the answer placeholder and validates the result |
| SAD5 | 5 | A Sources record line with no field separator reds with a message naming the one-physical-line rule | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/decision-map-integrity/source-wrapped-field` | A generic unknown-field message does not contain the rule text, so the fixture EXPECT does not match |
| SAD6 | 6 | A second `URL:` locator line still reds with the exact current message | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/decision-map-integrity/source-second-locator` | A blanket rewrite of the shared message changes these bytes and the pinned EXPECT fails |
| SAD7 | 7 | The rendered template holds exactly one decision-ticket heading | `internal/maps/maps_parse_test.go` (`TestParseDecisionMapSchemaAndTemplate`) | A second ticket leaves a second unresolved answer, so a ready-map consumer reds on the unresolved ticket |
| SAD8 | 8 | The rendered template states that a resolved ticket cannot stay blocked by an unresolved ticket | `internal/maps/maps_parse_test.go` (`TestParseDecisionMapSchemaAndTemplate`) | A template without the sentence fails the token assertion for that line |
| SAD9 | 9 | A resolved ticket blocked by an unresolved ticket still reds through the graph walk | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/decision-map-integrity/graph-resolved-on-unresolved` | Moving the rule into the template and dropping the walk edge silences the fixture's diagnostic |
| SAD10 | 10 | map-discipline.md requires one red-capable row for every in-scope edge-inventory, source-promise, and fence-closure promise | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-promise-rows` | Weakening the sentence to a recommendation removes the needle bytes and reds the anchor |
| SAD11 | 11 | map-discipline.md requires a Won't handle line with a surviving in-scope caller for an excluded edge | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller` | Dropping the surviving-caller clause leaves an exclusion nobody can trace, and the anchor reds |
| SAD12 | 12 | map-discipline.md requires the flagged-additions list under Further notes before the first review charge | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions` | Removing the before-the-charge clause lets the list arrive after review, and the anchor reds |
| SAD13 | 13 | map-discipline.md requires the source-sentence-to-row table under Further notes before the first review charge | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table` | Removing the table sentence leaves review with no per-sentence check, and the anchor reds |
| SAD14 | 14 | map-discipline.md states that the round demands a row for each listed addition and removes each unlisted addition | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition` | Keeping only the demand arm lets an unlisted addition survive, and the mutated needle reds |
| SAD15 | 15 | map-discipline.md requires two rows for an either-side predicate, one side per row | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows` | A single row that names both sides passes a one-row reading, so the two-row needle must be exact |
| SAD16 | 16 | No file under `internal/coverage/` changes | review-owned: the Standards axis reads the diff against the ownership fences | An edit to the checker falls outside every fence, and only a reader compares the diff to the fence list |
| SAD17 | 17 | map-discipline.md requires an executed-root trace for each canary or conformance row before the map locks | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace` | Dropping the trace lets an unreachable row look covered, and the anchor reds |
| SAD18 | 18 | map-discipline.md requires a named diagnostic state to be addable or mutable in a fixture | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state` | A class nobody can create in a fixture cannot bite, and removing the sentence reds the anchor |
| SAD19 | 19 | map-discipline.md requires the author to quote each pasted operand in the charge | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands` | A paraphrased operand reaches the delegate as a different value, and the anchor reds on removal |
| SAD20 | 20 | The Template section of craft-spec SKILL.md holds no contract text under the Further notes heading | `internal/anchors/registry_data_test.go` (`TestCraftSpecMapDisciplineAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/craft-spec-further-notes-bare` | A forbid-in-section tuple over the Template section reds the moment contract text lands under that heading |
| SAD21 | 21 | The authoritative-inventory rule and the named-test-function rule appear once each in map-discipline.md | review-owned: the Standards axis reads the reference file | A second copy of either rule drifts from the first, and only a reader can see the duplication |
| SAD22 | 22 | The craft-spec sweep sentence names the reader sweep | `internal/anchors/registry_data_test.go` (`TestCraftSpecMapDisciplineAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/craft-spec-reader-sweep-name` | An unnamed sweep gives the phase file nothing to sequence, and the anchor reds |
| SAD23 | 23 | map-discipline.md requires each named consumer of the decision fact in the sweep | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers` | The FT216 miss was a named shared consumer, so removing this clause restores that miss |
| SAD24 | 24 | map-discipline.md requires each helper a named consumer calls directly | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers` | A sweep of named consumers alone leaves the shared helper unfenced, and the anchor reds |
| SAD25 | 25 | map-discipline.md admits a deeper callee only when the callee reads the fact itself | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound` | Without the bound the sweep walks the whole call graph, and removing the clause reds the anchor |
| SAD26 | 26 | map-discipline.md requires an exact ownership fence for each shared reader | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence` | A sweep with no fence rule produces an inventory nobody assigns, and the anchor reds |
| SAD27 | 27 | The kit writes "reader sweep" and never "reader census" for this sweep | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/reader-sweep-term` | A "census" spelling collides with the raw-call census, and the forbid needle reds on it |
| SAD28 | 28 | The phase file sequences the reader sweep before the craft-tickets charge | `internal/anchors/registry_data_test.go` (`TestCraftDelegateDisciplineAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/write-spec-reader-sweep-sequence` | A sweep sequenced after slicing arrives too late to shape the fences, and the anchor reds |
| SAD29 | 29 | The phase file names the sweep and states no sweep rule of its own | review-owned: the Standards axis reads the phase-file diff | Only a reader can see that the phase file restated a rule the skill owns |
| SAD30 | 30 | The ship-test question sits immediately before the ticket-slicing step | review-owned: the Standards axis reads the placement | The slicing step and the review-round paragraph share one H2 section, so no section-scoped anchor can express the placement |
| SAD31 | 31 | The review-round paragraph holds no copy of the ship-test question | review-owned: the Standards axis reads the placement | A section-scoped forbid over Process would also forbid the required placement, so only a reader can see a second copy |
| SAD32 | 32 | The moved ship-test sentence carries the split arm, which states that the graph splits where a consumer branch lands green alone | `internal/anchors/registry_data_test.go` (`TestCraftDelegateDisciplineAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/write-spec-branch-split` | The whole-file require needle covers the moved sentence, so dropping the split arm reds the anchor |
| SAD33 | 33 | The phase file stays at 73 lines or fewer | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsHoldOnTheLiveTree`) | An added sweep sentence without the compensating move pushes the file past its budget |
| SAD34 | 34 | bench-shape-idea.md tells the author to read one ready decision map's Sources block before the first write | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/shape-idea-read-ready-sources` | An author who reads no live record guesses the shape again, and the anchor reds on removal |
| SAD35 | 35 | bench-shape-idea.md tells the author to run bench maps and bench gate-prose after the first skeleton | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/shape-idea-skeleton-checks` | Naming only one of the two verbs leaves one lane unrun, and the exact needle reds |
| SAD36 | 36 | No kit file contains the term "prose preflight" | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/prose-preflight-term` | The roadmap wording invites the term, so only a forbid needle keeps it out |
| SAD37 | 37 | bench-write-spec.md, bench-shape-idea.md, and the template each spell `decisions/assets/` | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/decision-map-asset-path` | A file that spells a different path sends an asset where the scanner reads it, and the anchor reds |
| SAD38 | 38 | The candidate scanner still skips a map file nested under a decisions asset directory | `internal/maps/maps_parse_test.go` (`TestDiscoverDecisionMapCandidatesDirectChildren`) | A recursive walk would return the nested file, and the golden candidate list fails |
| SAD39 | 39 | CONTEXT.md defines "coverage map" and lists the terms to avoid | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/context-coverage-map-term` | Without the entry an author writes a bare "map", and the anchor reds on the missing needle |
| SAD40 | 40 | CONTEXT.md keeps the "decision map" entry with its Avoid list | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/context-decision-map-term` | A merged map entry loses the distinction between the two artifacts, and the anchor reds |
| SAD41 | 41 | CONTEXT.md defines "reader sweep" and reserves "census" for the raw-call census and the architecture census | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/context-reader-sweep-term` | Without the reservation "census" spreads to a third meaning, and the anchor reds |
| SAD42 | 42 | Every sentence this build adds or moves in craft-spec SKILL.md, map-discipline.md, bench-write-spec.md, bench-shape-idea.md, CONTEXT.md, and the rendered template writes a two-word map term | review-owned: the Standards axis reads the whole diff | A bare "map" reads as correct in isolation, so only a reader across the diff can see it |

### Edge inventory

- **The template-parses-clean edge.** The rendered template validates in a temporary root
  with no repository files. A `Path:` locator would red there. Row SAD2 holds this edge.
- **The ready-status edge.** Two consumers convert the rendered template to `Status: ready`
  before they validate it. A ready map runs the readiness diagnostics, which include the
  Sources diagnostics. Row SAD4 holds this edge.
- **The wrapped-continuation edge.** A wrapped bullet stays legal in a terminal list. The
  new refusal speaks about a Sources record alone. Row SAD5 holds the new message, and the
  existing terminal-list test holds the legality.
- **The shared-message edge.** The second-locator fixture pins bytes the same code path
  produces. Row SAD6 holds this edge.
- **The budget-ceiling edges.** craft-spec SKILL.md takes one line and ends at 151 of 152.
  craft-tickets SKILL.md takes none. The phase file nets zero and stays at 73. Row SAD33
  holds the phase file, and row SAD20 holds the bare Further notes heading.
- **The anchor-needle-reflow edge.** Existing fixtures pin exact bytes in the phase file
  near the slicing step. A reflow of those sentences reds them. Row SAD32 holds the moved
  needle, and the fixture family holds the untouched ones.
- **The single-H2-section edge.** The phase file holds the slicing step and the
  review-round paragraph under one `## Process` heading. Rows SAD30 and SAD31 are
  review-owned for that reason.
- **The forbid-needle edges.** Two new terms need a forbid tuple rather than a require
  tuple, because absence is the property. Rows SAD27 and SAD36 hold these edges.

**Won't handle** — a `Path:` locator example in the template — a `URL:` example teaches the
same record shape. `/bench-shape-idea` still sends the author to a live map's Sources
block for a path example.

**Won't handle** — a Go check on a map-owned asset outside `decisions/assets/` — the
scanner already skips the nested directory. `/bench-write-spec` still moves assets at
compile time.

**Won't handle** — a second decision ticket in the template — eleven consumers resolve one
answer placeholder, and `/bench-shape-idea` still tells the author to add tickets.

**Won't handle** — a machine check for the two required lists under Further notes —
`bench coverage --check` stays unchanged by decision. The review round still demands both
lists.

**Won't handle** — a forbid needle for a bare "map" in new prose — the word is legal
inside "decision map" and "coverage map". The Standards axis still reads the diff.

## Ownership fences

- `internal/maps/schema.go`
- `internal/maps/validation.go`
- `internal/maps/maps_parse_test.go`
- `internal/maps/maps_graph_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/conformance/decision_map_integrity_test.go`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.agents/skills/bench-craft-spec/references/map-discipline.md`
- `.agents/commands/bench-write-spec.md`
- `.agents/commands/bench-shape-idea.md`
- `CONTEXT.md`
- `tests/canary/decision-map-integrity`
- `tests/canary/workflow-guidance-anchors`
- `specs/spec-authoring-discipline`
- `reviews/spec-authoring-discipline.md`

## Out of scope

- A budget-row raise for `.agents/commands/bench-write-spec.md` or for
  `.agents/skills/bench-craft-spec/SKILL.md`. The phase file nets zero lines and stays at
  73 of 73. The skill file takes one line and ends at 151 of 152. Estimate: 2 edits, 1
  gate run.
- A structured Further-notes grammar that `bench coverage --check` parses, with rows for the
  flagged-additions list and the source-sentence-to-row table. Estimate: 6 edits, 3 gate runs.
- A `bench maps --explain <diagnostic>` verb that prints the rule behind each refusal.
  Estimate: 8 edits, 3 gate runs.
- A scanner refusal that reports a decision map found outside a scanned directory.
  Estimate: 5 edits, 2 gate runs.

## Further notes

### Enforcement reads

The author opened each file below in this session.

- `internal/maps/schema.go` — `canonicalDecisionMapSchema` at line 116 and
  `DecisionMapTemplate` at line 427. The renderer writes each terminal heading with no body,
  so the Sources example is new body text.
- `internal/maps/validation.go` — `sourceDiagnostics` at lines 134 to 202. The unknown-field
  message at lines 154 to 156 serves both a wrapped line and a second locator line. The
  resolved-blocked edge rule sits at lines 57 to 59. `validateSourcePath` runs for a `Path:`
  locator alone, at lines 183 to 186.
- `internal/maps/maps_parse_test.go` — `TestParseDecisionMapSchemaAndTemplate` at line 14,
  with the token assertion at line 57 and the template-diagnostics assertions at lines 65 and
  68. `TestDiscoverDecisionMapCandidatesDirectChildren` at line 205 already writes
  `decisions/assets/nested.md` and expects the scanner to skip it.
  `TestDecisionMapDiagnosticsGolden` at line 335 pins the ordered diagnostic slice.
- `internal/maps/maps_graph_test.go` — `TestMapSourcesRequireExactRecordShape` at line 264,
  whose message pins run from line 270 to line 276.
  `TestMapTerminalContinuationAndEmptyAnswer` at lines 227 to 232 proves a wrapped bullet is
  legal in a terminal list.
- `internal/anchors/registry_data.go` — the phase-file tuples at lines 25 to 47, the
  narrower-capability needle at line 132, and the ticket-slicing needles at lines 181 to 184.
  Each tuple has an empty section, so an insert elsewhere in the file keeps it green.
  `internal/anchors/registry_data_test.go` holds
  `TestCraftSpecMapDisciplineAnchorsRedOnRemoval` at line 541 and
  `TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval` at line 1148.
- `internal/conformance/fixture_bite_test.go` — `TestEveryRetainedFixtureBitesThroughRegisteredOwner`
  at line 21 runs every fixture under `tests/canary`.
- `internal/conformance/prose_budget_test.go` — lines 23 to 38 name
  `projects/benchkit.md` and the `Guidance prose budgets` section as the one budget source.
- `projects/benchkit.md` — the budget table at lines 476 to 486, and the sentence at line 489
  that keeps every other `.agents/commands/*.md` file outside the reviewed universe.
  `.agents/commands/bench-shape-idea.md` and the reference files therefore carry no budget.

The author did not read `internal/coverage`, because row SAD16 rests on the fence and not on
that package's behavior.

### Precedent read

The FT144 landing at commit `45deacc9` added two rules to
`.agents/skills/bench-craft-spec/references/map-discipline.md`. Each rule received an anchor
tuple, a red-on-removal registry test, and a canary fixture under
`tests/canary/workflow-guidance-anchors/`. The author read
`tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory`, whose fixture
holds `BASE`, `EXPECT`, and `MUTATE.json`. Every new normative sentence in this spec ships
in that shape.

### Reader sweep

- Eleven test call sites resolve the template's `<answer>` with a single-count
  replacement. The author found each site with one repository-wide `rg` run. The
  one-ticket decision keeps each site green. The sites are:
  - `internal/maps/maps_command_test.go` at lines 95, 309, and 325
  - `internal/maps/maps_parse_test.go` at line 253
  - `internal/conformance/decision_map_integrity_test.go` at line 90
  - `internal/status/status_signals_test.go` at lines 147, 165, and 335
  - `internal/status/status_producible_test.go` at line 92
  - `internal/status/status_command_test.go` at lines 126 and 160
- `internal/status/status_producible_test.go` at line 298 writes the raw template as a
  shaping map. The Sources example must also satisfy that call site.
- The refusal message has two readers:
  `tests/canary/decision-map-integrity/source-second-locator/EXPECT`, which holds
  `Sources https://example.invalid/second unexpected field URL`, and
  `internal/maps/maps_graph_test.go` at lines 270 to 276.
- The phase-file text near lines 57 to 64 has one reader:
  `internal/anchors/registry_data.go`. The author checked the needles at lines 132 and 181
  to 184. The move keeps each sentence byte for byte, so those needles survive.

### Executed-root trace

This spec applies its own SAD17 and SAD18 rules to its canary and conformance rows.

- Every canary row reaches the executed root through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`
  in `internal/conformance/fixture_bite_test.go`. That test enumerates `tests/canary` from
  the `canary.Fixtures` producer, so a new fixture joins the run with no registry edit.
- `internal/conformance/registry/registry.go` binds the family at line 171 for
  `workflow-guidance-anchors` and at line 174 for `decision-map-integrity`. Both new
  fixtures land inside an already-bound family.
- The gate executes `internal/conformance` and `internal/maps` as ordinary Go test
  packages, so the cited functions sit inside an executed phase.
- The `source-wrapped-field` diagnostic state is addable. A `BASE` map with a wrapped
  Sources continuation line creates it, and no other tree state is needed.
- Each `workflow-guidance-anchors` diagnostic state is mutable. `MUTATE.json` rewrites one
  needle in a real guidance file, which is the FT144 fixture shape.
- Row SAD4 reaches the executed root through
  `TestDecisionMapIntegrityCheckValidatesEveryCandidate`, which drives the live validator
  over the rendered template.

### Flagged additions

Each item below goes beyond the fifteen predicates and the three roadmap bodies. Review
demands a row for each item, or removes the item.

1. The ownership fence covers `internal/conformance/decision_map_integrity_test.go`, because
   row SAD4 drives the rendered template through the live validator there.
2. Row SAD16 rests on the ownership fence rather than on a test. Predicate 9 forbids a
   coverage-checker change, and no test can observe an absent edit.
3. Rows SAD21, SAD29, and SAD42 are review-owned. No check can see a duplicated rule or a bare term
   inside otherwise correct prose.
4. Rows SAD30 and SAD31 are review-owned. The anchor mechanism scopes a sectioned anchor
   to an H2 section, and the phase file holds both positions under one `## Process`
   heading.
5. Row SAD20 uses a forbid-in-section tuple rather than the line budget. A one-line
   addition under the Further notes heading does not cross the 152-line budget, so the
   budget cannot grade that heading.

### Source sentences and their rows

| source sentence | rows |
|---|---|
| Predicate 1 — one spec covers FT220, FT257, and FT278 | the whole map |
| Predicate 2 — a new rule lands in an unbudgeted reference file | SAD20, SAD33 |
| Predicate 3 — the sweep sentence gains the extension and the phase file sequences it | SAD22, SAD28, SAD29 |
| Predicate 4 — named consumers plus direct helpers, deeper callees only when they read the fact | SAD23, SAD24, SAD25 |
| Predicate 5 — the canonical term is reader sweep | SAD27, SAD41 |
| Predicate 6 — the ship-test question moves and the review section drops it | SAD30, SAD31, SAD32 |
| Predicate 7 — one red-capable row per in-scope promise, and a Won't handle for each cut edge | SAD10, SAD11 |
| Predicate 8 — the flagged-additions list and the source table under a bare Further notes heading | SAD12, SAD13, SAD14, SAD20 |
| Predicate 9 — two rows for an either-side predicate, and no coverage-checker change | SAD15, SAD16 |
| Predicate 10 — a URL Sources example that parses clean | SAD1, SAD2, SAD3, SAD4 |
| Predicate 11 — one ticket, a stated resolved-blocked rule, and the existing enforcement | SAD7, SAD8, SAD9 |
| Predicate 12 — the shape-idea authoring steps and no prose-preflight term | SAD34, SAD35, SAD36 |
| Predicate 13 — the two-word terms and the CONTEXT.md entries | SAD39, SAD40, SAD41, SAD42 |
| Predicate 14 — the decisions asset path in three files, with no Go change | SAD37, SAD38 |
| Predicate 15 — the wrapped-field refusal, the pre-charge checks, and fixture reachability | SAD5, SAD6, SAD17, SAD18 |
| Predicate 15 — the pasted-operand check and the two rules that stay single-sourced | SAD19, SAD21 |
| FT220 body — each shared reader under an exact fence | SAD26 |
