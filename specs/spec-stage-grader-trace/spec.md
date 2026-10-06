# The spec stage traces each pinned check, entry read, and derived expectation to its grader

Status: staged

Roadmap: FT376

Decision source: `roadmap/FT376.md`, a named reviewed artifact from drain `d-007c40a25f47`.

Verification log: 2 iteration(s) to accept — the fable-high spec review accepted with fixes in iteration 1. Its fixes were an unanchored gloss, two unbounded two-clause needles, an exact canary fence, and three wording nits. The sonnet-high ticket review accepted with fixes in iteration 1, for missing N1 and N6 probes and the fence-union sentence. It accepted in iteration 2, after the delegated consultant merged five tickets into one.

## Problem

A spec author writes coverage rows before the first review charge. `craft-spec` asks for the cheapest wrong answer for each row. No pre-review check traces a row through the code that grades it, so review finds the gap instead.

Four occurrences of this gap fell between 2026-10-02 and 2026-10-03:

- The worktree-seam-reduction spec stage had four untraced rows. One row pinned a count, and its grader read the count as a floor.
- In the same build, ticket WS51 and ticket 8 had a named probe that stayed silent until the test graded the refusal line.
- The FT358 spec review took two iterations. It found unverified claims, per-site behavior differences, and obligations that some tickets did not carry.
- The FT370 spec review found a missed `RecordCompletion` caller and a repair assignment that changed the plan digest.

The current pre-review proof checklist in `.agents/skills/bench-craft-spec/references/map-discipline.md` has seven classes. No class asks for an operator, an entry read, an expectation source, a consolidation table, a quantified obligation, or a digest trace. The `Changed-function callers` class names no method, so an author can run a text search and miss an unexported caller.

## Solution

The pre-review proof checklist in `map-discipline.md` gains six classes, and the `Before the map locks` section gains one caller-sweep rule. Each new rule sentence is an anchor needle. The one unanchored new sentence is the gloss that defines an internal form. The gate reds when a needle leaves its section, and seven canaries show that the gate reds a weakened needle.

The new classes make the author record, before the first review charge:

- the comparison operator in the grader of each pin row
- each unexported read below an entry that has no internal form, with its grader
- the grader of each derived expectation, and the rule that no expectation comes from the code under test
- a consolidation table of each site's old rule and new rule, with a disposition for each changed cell
- a check of every quantified obligation at each affected site and across all tickets
- a trace of every write of a new workflow step through the digests that later checkpoints compare

The caller-sweep rule makes the author run `bench consumers` for each changed function, unexported functions included. The verb already resolves an unexported symbol, so this spec changes no CLI.

## User stories

Line: opus / high.
Implementation-line reason: GT-C1 is the only chunk. It adds guidance prose that every later spec author reads, so the leverage override in `craft-line` sets mid tier at high effort. The spec quotes each needle, each diagnostic, and each canary mutation. The seam is the existing anchor harness, and its unit tests red each row cheaply.
Harder chunks: none.

### The checklist traces each pinned check, entry read, and derived expectation to its grader

1. As a spec author, I want each pin row's grader operator quoted, so that a count pin read as a floor shows early.
2. As a spec author, I want each unexported read below an entry without an internal form listed, so that the read gets a grader.
3. As a spec author, I want the grader of each derived expectation named, so that a silent probe shows before review.
4. As a spec reviewer, I want no expectation from the code under test, so that no row passes by construction.

### A consolidation maps each changed site rule

5. As a spec author who consolidates rules, I want each site's old and new rule in one table, so that site differences show.
6. As a spec reviewer, I want each changed table cell to take a disposition, so that no site change ships silently.

### Quantified obligations hold at each site and in each ticket

7. As a spec author, I want each quantified obligation checked at every affected site, so that no site misses it.
8. As a ticket slicer, I want each quantified obligation checked across all tickets, so that no ticket drops it.

### The caller sweep uses the consumers query

9. As a spec author, I want the caller sweep to run `bench consumers`, so that it finds callers that a text search misses.
10. As a spec author, I want the caller sweep to include unexported functions, so that no unexported caller reaches review unseen.

### A new workflow step traces its writes

11. As a spec author, I want each new workflow step's writes traced through compared digests, so that a digest change shows early.

### The gate holds the new rules

12. As a kit maintainer, I want a conformant tree to raise no new diagnostic, so that the anchors cannot red good guidance.
13. As a kit maintainer, I want the live kit to satisfy each new anchor, so that the gate grades the shipped text.
14. As a kit maintainer, I want a canary for each needle that can lose a clause, so that the gate reds a weakened sentence.
15. As a linked-repository user, I want a changelog entry for the new rules, so that the release notes show the change.

## Implementation decisions

### Guidance text

The `Before the map locks` section of `map-discipline.md` keeps every current sentence byte for byte. The fixed pre-review proof checklist gains six sub-bullets after `Rendered-shape readers`. Each sub-bullet stays on one physical line. The exact sentences are:

- N1: "`Pin operators` quotes the comparison operator that the grader of each pin row applies."
- N2: "`Entry reads` lists each unexported read below an entry without an internal form, and names its grader." The same line adds the unanchored gloss "An internal form takes injected values in place of ambient reads." The gloss defines a term and states no rule, so it is the one new sentence without an anchor.
- N3 and N4 share one line: "`Derived expectations` names the grader of each derived expectation." and "No expectation comes from the code under test."
- N5 and N6 share one line: "`Consolidated rules`, when a spec consolidates repeated rules, gives a consolidation table of each site's old rule and new rule." and "Each changed cell of the consolidation table maps to an acceptance row, a flagged addition, or a Won't handle line."
- N7: "`Quantified obligations` checks every quantified obligation at each affected site and across all tickets."
- N9: "`Workflow-step writes`, for each new workflow step, traces every write through the digests that later checkpoints compare."

A new top-level bullet follows the checklist, with N8: "The changed-function caller sweep runs `bench consumers` for each changed function, unexported functions included."

`craft-spec/SKILL.md` does not change. It holds 155 lines against its budget of 155 in `projects/benchkit.md`, and `map-discipline.md` is the reference that `craft-spec` charges for the checklist.

### Anchor registry

A new themed registry pair holds the nine anchors: `internal/anchors/registry_spec_trace.go` declares `specTraceAnchors`, and `registry_spec_trace_test.go` tests them. Each anchor has group `AfterImplementSpec`, kind `RequireInSection`, file `mapDiscipline`, and section `Before the map locks`. The file reuses the `mapDiscipline` constant from `registry_ticket_passes.go`.

The diagnostics are:

| needle | diagnostic |
| --- | --- |
| N1 | `map discipline: the pre-review checklist quotes each pin row's grader operator` |
| N2 | `map discipline: the pre-review checklist traces each unexported entry read to its grader` |
| N3 | `map discipline: the pre-review checklist names the grader of each derived expectation` |
| N4 | `map discipline: no expectation comes from the code under test` |
| N5 | `map discipline: a rule consolidation gives each site's old and new rule` |
| N6 | `map discipline: each changed consolidation cell takes a row, a flagged addition, or a Won't handle line` |
| N7 | `map discipline: quantified obligations hold at each site and across all tickets` |
| N8 | `map discipline: the caller sweep runs bench consumers on unexported functions too` |
| N9 | `map discipline: a new workflow step traces its writes through checkpoint digests` |

The `registry` declaration in `registry_data.go` appends `specTraceAnchors` at the end of its chain. That edit changes one line and adds none. `registry_data.go` holds 479 lines and `registry_data_test.go` holds 1286, both above the 400-line limit, so `bench structure --growth` reds any line that either file gains.

The test `TestSpecGraderTraceAnchors` writes its needles, sections, and diagnostics independently of the registry, after the precedent of `TestReviewRuleAnchors`. The `anchorHarness` check then proves both directions. The conformant tree raises none of the nine diagnostics. Each broken tree raises its own diagnostic and no other.

### Canaries

Seven new fixtures under `tests/canary/workflow-guidance-anchors/` each weaken one needle. N3 and N4 take no canary, because each one states a single clause. Each fixture has `BASE` with the line `.agents/skills/bench-craft-spec/references/map-discipline.md`, a `MUTATE.json` with one replacement, and an `EXPECT` that holds the diagnostic of the mutated needle.

| fixture | old | new | EXPECT |
| --- | --- | --- | --- |
| `map-discipline-pin-operator-trace` | N1 | "`Pin operators` lists each pin row." | N1 diagnostic |
| `map-discipline-entry-read-grader` | N2 | "`Entry reads` lists each unexported read below an entry without an internal form." | N2 diagnostic |
| `map-discipline-consolidation-both-rules` | N5 | "`Consolidated rules`, when a spec consolidates repeated rules, gives a consolidation table of each site's new rule." | N5 diagnostic |
| `map-discipline-consolidation-cells` | N6 | "Each changed cell of the consolidation table maps to an acceptance row." | N6 diagnostic |
| `map-discipline-quantified-tickets` | N7 | "`Quantified obligations` checks every quantified obligation at each affected site." | N7 diagnostic |
| `map-discipline-unexported-callers` | N8 | "The changed-function caller sweep runs `bench consumers` for each changed function." | N8 diagnostic |
| `map-discipline-workflow-step-digests` | N9 | "`Workflow-step writes`, for each new workflow step, lists every write." | N9 diagnostic |

`canary.Fixtures` discovers each fixture, and no file pins the fixture count.

### Changelog

`CHANGELOG.md` gains one `### Spec grader trace` entry under `## [Unreleased]`. It names the six new checklist classes and the caller-sweep rule.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| GT-C1 / `1-trace-spec-rows-to-graders.md` | The pre-review checklist and the caller sweep carry the FT376 rules, and the gate holds each rule | GT1, GT2, GT3, GT4, GT5, GT6, GT7, GT8, GT9, GT10, GT11, GT12, GT13, GT14, GT15, GT16, GT17, GT18, GT19 | `bench test --package ./internal/anchors`, `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner`, `bench test --check canary-fixture-compliance`, `bench test --check docs-currency-workflow` | no |

One ticket delivers the chunk and is its one serial green checkpoint. Each needle adds one guidance line, one anchor row, one test rule, and at most one canary to the seam that the ticket opens. A split would give no earlier review checkpoint, because the chunk has one review. The commitment criteria in the source trace keep each FT376 criterion traceable to its rows.

```bench-completion-plan
{"version":1,"chunks":[{"id":"GT-C1","tickets":["1-trace-spec-rows-to-graders.md"],"verification":[{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"fixture-bites","command":"bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner"},{"id":"canary-fixture-compliance","command":"bench test --check canary-fixture-compliance"},{"id":"docs-currency-workflow","command":"bench test --check docs-currency-workflow"}]}],"final_verification":[{"id":"coverage-check","command":"bench coverage --check specs/spec-stage-grader-trace/spec.md"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"docs-currency-workflow","command":"bench test --check docs-currency-workflow"}]}
```

## Testing decisions

- A good test writes a minimal tree with the needles, removes one needle, and observes the diagnostic that `EvaluateGroup` returns. It does not read the live guidance, because the live read belongs to the gate.
- The anchor rows attach at `anchorHarness` in `internal/anchors`, after the precedent of `TestReviewRuleAnchors` and `TestMapDisciplineAnchorsRedOnRemoval`.
- The canary rows attach at `TestEveryRetainedFixtureBitesThroughRegisteredOwner`, which applies each `MUTATE.json` to a copy of the kit and runs the registered owner `docs-currency-workflow`.
- The live-tree row attaches at the `docs-currency-workflow` conformance check in `bench gate`. Its executed root is `checkDocsCurrencyAndWorkflow`, which calls `checkWorkflowAnchors`, which calls `anchors.EvaluateGroup(root, anchors.AfterImplementSpec)`.

### Seam diagram

    trigger: bench gate, the anchors package tests, and the fixture-bite test
        │
        ▼
    map-discipline.md  ──▶  [ anchors.EvaluateGroup over specTraceAnchors ]  ──▶  diagnostics
                      ◀ tests attach here: anchorHarness writes a tree without one needle;
                        a canary mutates one needle in a kit copy; the gate reads the live kit

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| GT1 | 1 | A `Before the map locks` section without N1 raises the N1 diagnostic | planned TestSpecGraderTraceAnchors in internal/anchors/registry_spec_trace_test.go | A registry without the N1 anchor raises nothing when the sentence leaves. |
| GT2 | 2 | A `Before the map locks` section without N2 raises the N2 diagnostic | planned TestSpecGraderTraceAnchors in internal/anchors/registry_spec_trace_test.go | A registry without the N2 anchor raises nothing when the sentence leaves. |
| GT3 | 3 | A `Before the map locks` section without N3 raises the N3 diagnostic | planned TestSpecGraderTraceAnchors in internal/anchors/registry_spec_trace_test.go | A registry without the N3 anchor raises nothing when the sentence leaves. |
| GT4 | 4 | A `Before the map locks` section without N4 raises the N4 diagnostic | planned TestSpecGraderTraceAnchors in internal/anchors/registry_spec_trace_test.go | An anchor that joins N3 and N4 in one needle cannot name which rule left. |
| GT5 | 5 | A `Before the map locks` section without N5 raises the N5 diagnostic | planned TestSpecGraderTraceAnchors in internal/anchors/registry_spec_trace_test.go | A registry without the N5 anchor raises nothing when the sentence leaves. |
| GT6 | 6 | A `Before the map locks` section without N6 raises the N6 diagnostic | planned TestSpecGraderTraceAnchors in internal/anchors/registry_spec_trace_test.go | A registry without the N6 anchor raises nothing when the sentence leaves. |
| GT7 | 7 | A `Before the map locks` section without N7 raises the N7 diagnostic | planned TestSpecGraderTraceAnchors in internal/anchors/registry_spec_trace_test.go | A registry without the N7 anchor raises nothing when the sentence leaves. |
| GT8 | 8 | The `map-discipline-quantified-tickets` mutation raises the N7 diagnostic | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An anchor on a shorter needle without `across all tickets` stays green when the ticket clause leaves. |
| GT9 | 9 | A `Before the map locks` section without N8 raises the N8 diagnostic | planned TestSpecGraderTraceAnchors in internal/anchors/registry_spec_trace_test.go | A registry without the N8 anchor raises nothing when the sentence leaves. |
| GT10 | 10 | The `map-discipline-unexported-callers` mutation raises the N8 diagnostic | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An anchor on a shorter needle without `unexported functions included` stays green when that clause leaves. |
| GT11 | 11 | A `Before the map locks` section without N9 raises the N9 diagnostic | planned TestSpecGraderTraceAnchors in internal/anchors/registry_spec_trace_test.go | A registry without the N9 anchor raises nothing when the sentence leaves. |
| GT12 | 12 | A tree that holds N1 to N9 in `Before the map locks` raises none of the nine diagnostics | planned TestSpecGraderTraceAnchors in internal/anchors/registry_spec_trace_test.go | An anchor with the wrong section or a typo in its needle reds a conformant tree. |
| GT13 | 13 | The `docs-currency-workflow` check over the kit root raises none of the nine diagnostics | the `docs-currency-workflow` conformance check in `bench gate` | A needle that the live guidance does not hold verbatim reds the gate. |
| GT14 | 14 | The `map-discipline-pin-operator-trace` mutation raises the N1 diagnostic | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An anchor on the class label alone stays green when the operator clause leaves. |
| GT15 | 14 | The `map-discipline-consolidation-cells` mutation raises the N6 diagnostic | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An anchor that stops before the Won't handle arm stays green when that arm leaves. |
| GT16 | 14 | The `map-discipline-workflow-step-digests` mutation raises the N9 diagnostic | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An anchor on the class label alone stays green when the digest clause leaves. |
| GT17 | 15 | The `## [Unreleased]` section of `CHANGELOG.md` holds a `### Spec grader trace` entry that names the six classes and the caller-sweep rule | review-owned | No check grades changelog content, so review reads the entry. |
| GT18 | 14 | The `map-discipline-entry-read-grader` mutation raises the N2 diagnostic | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An anchor that stops before `and names its grader` stays green when the grader clause leaves. |
| GT19 | 14 | The `map-discipline-consolidation-both-rules` mutation raises the N5 diagnostic | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An anchor on the table clause alone stays green when the old-rule column leaves. |

### Edge inventory

The in-scope edges are a removed needle (GT1 to GT11) and a weakened needle (GT8, GT10, GT14 to GT16, GT18, GT19). The others are a conformant tree (GT12) and the live kit (GT13). A needle that moves to another section is a removed needle for its own section, so the section-scoped rows cover it.

The audience of each behavior is every repository that links the kit, because `map-discipline.md` ships to each one. The anchors and canaries serve this repository only. `map-discipline.md` is one file, so the absent-versus-empty directory pair does not apply. An absent file raises `section-scoped anchor file missing`, and the existing map-discipline anchors already cover that state.

The hostile-input classes of `projects/benchkit.md` do not reach this surface:

- No path, branch name, or commit subject enters the needles or the diagnostics.
- No `git diff` header reaches the anchor evaluator.
- No TOON cell carries the new text, because `bench anchors` already renders every needle through its existing path.

**Won't handle** — a retrofit of the new classes into a staged spec — each staged spec in `specs/` is past its first review charge.

**Won't handle** — a gate check that parses Further notes for the new labels — no check parses the seven current classes either.

**Won't handle** — a change to `craft-spec/SKILL.md` — the file sits at its 155-line budget, and it already points to `map-discipline.md`.

## Ownership fences

- `.agents/skills/bench-craft-spec/references/map-discipline.md`
- `internal/anchors/registry_spec_trace.go`
- `internal/anchors/registry_spec_trace_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `tests/canary/workflow-guidance-anchors/map-discipline-consolidation-both-rules`
- `tests/canary/workflow-guidance-anchors/map-discipline-consolidation-cells`
- `tests/canary/workflow-guidance-anchors/map-discipline-entry-read-grader`
- `tests/canary/workflow-guidance-anchors/map-discipline-pin-operator-trace`
- `tests/canary/workflow-guidance-anchors/map-discipline-quantified-tickets`
- `tests/canary/workflow-guidance-anchors/map-discipline-unexported-callers`
- `tests/canary/workflow-guidance-anchors/map-discipline-workflow-step-digests`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/craft-spec-transaction-failure-rows`
- `tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory`
- `tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition`
- `tests/canary/workflow-guidance-anchors/map-discipline-copy-survival-proof`
- `tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows`
- `tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller`
- `tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace`
- `tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state`
- `tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions`
- `tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep`
- `tests/canary/workflow-guidance-anchors/map-discipline-pre-review-proof`
- `tests/canary/workflow-guidance-anchors/map-discipline-pre-review-source-consumers`
- `tests/canary/workflow-guidance-anchors/map-discipline-promise-rows`
- `tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands`
- `tests/canary/workflow-guidance-anchors/map-discipline-shipped-surface-claim`
- `tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table`
- `tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound`
- `tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers`
- `tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers`
- `tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence`
- `tests/canary/workflow-guidance-anchors/reader-sweep-term`
- `CHANGELOG.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `reviews/spec-stage-grader-trace.md`

Edits go only to the guidance file, the new registry pair, and one registry line in the anchor data file. The seven new canary directories and the changelog also take edits.

The other paths are closure entries that build preflight requires. Each anchor registry file that names the guidance file is a closure entry. Each fixture directory that pins an edited path is a closure entry, and the fixture-closure row of build preflight enumerates those directories. The binding registry binds the anchors package to the five command-registry files. This list equals the union of the ticket write lines, plus the review pickup.

## Out of scope

- A `bench coverage --check` rule that refuses a spec whose Further notes lack a checklist class. Estimate: 8 edits, 2 gate runs.

## Further notes

### Source trace

| source sentence (`roadmap/FT376.md`) | rows |
| --- | --- |
| "A pin row quotes its comparison operator." | GT1, GT14 |
| "The reader sweep lists each unexported read below an entry without an internal form." | GT2, GT18 |
| "No expectation comes from the code under test." | GT3, GT4 |
| "A spec that consolidates repeated rules includes a table of each site's old rule and new rule." | GT5, GT19 |
| "Each changed cell maps to an acceptance row, a flagged addition, or a Won't handle line." | GT6, GT15 |
| "The author checks every quantified obligation at each affected site and across all tickets." | GT7, GT8 |
| "The caller sweep uses `bench consumers` for each changed function, including unexported functions." | GT9, GT10 |
| "For each new workflow step, the author traces every write through the digests that later checkpoints compare." | GT11, GT16 |
| Title: "traces each pinned check, entry read, and derived expectation to its grader before the first review" | GT1, GT2, GT3 |

The commitment criteria in `.bench/commitment.json` restate these sentences. `FT376.trace` maps to GT1 to GT4, `FT376.consolidation` to GT5 and GT6, `FT376.quantified` to GT7 and GT8, `FT376.callers` to GT9 and GT10, and `FT376.digests` to GT11.

### Flagged calls for reviewer veto

- **Checklist classes, not free rules.** The source says the author "makes three checks" before the first review charge. This spec puts each check in the fixed pre-review proof checklist, so Further notes records each result or `none`. A free rule would leave no record that review can read.
- **"Internal form" keeps the FT356 meaning.** The FT356 spec defines an internal form as the entry variant that takes injected kit, clock, or joins values. N2 adds a one-sentence gloss, because `map-discipline.md` has no glossary of that term.
- **"Code under test" excludes a shared lower seam.** The hostile-input checklist tells an author to derive a TOON expectation through the same `toon.Table` call that the producer makes. That call is a shared encoder, not the code under test, so N4 does not contradict it. The working agreement's one-source rule for test expectations governs tests; N4 governs the source of a spec row's expected value.
- **The caller-sweep rule is a separate bullet.** N8 adds a method beside the anchored `Changed-function callers` class. That choice keeps the existing needle and its canary byte for byte.
- **A new registry pair.** The growth ratchet reds any line that `registry_data.go` or `registry_data_test.go` gains, so the anchors go in a themed pair after the `registry_review_rules.go` precedent.
- **FT378 is not this spec's work.** An FT378 occurrence names FT376 as the owner of a clause-to-row check. The existing `Source-row clauses and occurrences` class owns that check, so this spec adds no clause-to-row rule.

### Flagged additions

- GT12 pins the conformant direction. The source does not name it, but the anchor harness proves both directions for each rule set.
- GT14, GT15, GT16, GT18, and GT19 add canaries for needles whose rows already have a removal test. They show that the gate reds a weakened sentence. The working agreement admits an independent expectation only when its red is demonstrated, and each canary is that demonstration.
- GT17 adds a changelog entry under the `craft-synthesis` rule for user-visible behavior.

### Reader sweep

- Readers of the pre-review checklist: `.agents/skills/bench-craft-spec/SKILL.md` points to `map-discipline.md`. `.agents/skills/bench-craft-tickets/SKILL.md` links the `Before the map locks` proof rules. The `/bench-write-spec` command charges `craft-spec`. No other surface restates the checklist.
- Go readers of the checklist labels: only `internal/anchors/registry_data.go` and `registry_data_test.go`. No Go code parses a spec's Further notes; `preflighttest/reviewcases.go` and `evidence_file_pages_test.go` only write a bare `## Further notes` heading into fixtures.
- Readers of `map-discipline.md` in the anchors package: `registry_data.go`, `registry_data_test.go`, `registry_ticket_passes.go`, and `registry_ticket_passes_test.go`. The fence names each one.
- The executed root of the anchor rows: `checkDocsCurrencyAndWorkflow` calls `checkWorkflowAnchors`, which calls `anchors.EvaluateGroup` for `AfterImplementSpec`. The registry owner `docs-currency-workflow` owns the `workflow-guidance-anchors` canary family in `internal/conformance/registry/registry.go`.
- Shipped-surface claim words: none. The new sentences name no repo-only path.

### Pre-review proof checklist

- `Cited symbols`: `EvaluateGroup`, `Entries`, `registry`, `mapDiscipline`, `anchorHarness`, `TestReviewRuleAnchors`, `TestMapDisciplineAnchorsRedOnRemoval`, `TestEveryRetainedFixtureBitesThroughRegisteredOwner`, `checkDocsCurrencyAndWorkflow`, `checkWorkflowAnchors`, `canary.Fixtures`, and `BoundFiles` resolve at `44882814`.
- `Import edges`: none. The new registry pair imports only `testing`.
- `Source-row clauses and occurrences`: the source trace quotes each clause. The four occurrences are in the Problem.
- `Promised field labels`: the six class labels and the nine diagnostics above.
- `Changed-function callers`: none. The spec changes no function; `registry` is a package variable, and `Entries`, `FilesBelow`, and `evaluate` read it unchanged.
- `Copy survival`: none. No new owner replaces copies.
- `Rendered-shape readers`: none. No rendered output changes.
- `Pin operators`: none. No row pins a count or a value.
- `Entry reads`: none. The spec adds no entry and no ambient read.
- `Derived expectations`: each expected diagnostic comes from the table above, and `TestSpecGraderTraceAnchors` restates it independently of the registry.
- `Consolidated rules`: none. The spec consolidates no repeated rule.
- `Quantified obligations`: "each new rule sentence is an anchor needle" holds for N1 to N9 in GT1 to GT11. The internal-form gloss is the one excluded sentence. "A canary for each needle that can lose a clause" holds for N1, N2, and N5 to N9. Rows GT8, GT10, GT14 to GT16, GT18, and GT19 hold those canaries.
- `Workflow-step writes`: none. The spec adds no workflow step.

`bench anchors .agents/skills/bench-craft-spec/references/map-discipline.md` reported 31 anchors at `44882814`. No anchored sentence changes, so no anchor claim of this spec conflicts with that output.

### Bootstrap authority

None. The change adds guidance text and anchors, and it launches no executable.

### Source disclosure

The roadmap row cites `capture/learnings.md` entries from drain `d-007c40a25f47`. That drain removed them from the journal, so this session could not re-read them. This session read `roadmap/FT376.md`, the FT376 criteria in `.bench/commitment.json`, `roadmap/FT378.md`, and the FT356 spec text at commit `923b0f00` for the term "internal form".
