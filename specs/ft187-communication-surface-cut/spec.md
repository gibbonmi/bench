# FT187 — cut the communication surface to the rules that fire

Status: staged

Decision source: reviewer-confirmed current conversation, 2026-08-03 — the reviewer reviewed six findings against `.bench/BENCH.md`, accepted the direction ("i like the changes"), directed the roadmap row and this spec, and ruled that FT156 does not gate the work.

## Problem

`.bench/BENCH.md` is always loaded, and its "How to talk to me" section is the single largest block in it. Three of its rules point in opposite directions, so a session resolves the conflict toward whichever clause it read most recently and the voice drifts turn to turn. Two rules elsewhere in the file are absolutes that no session can satisfy, which teaches that the file's absolutes are rhetorical — expensive standing beside four invariants whose absolutes are real.

The reviewer reads every session's output through this surface. A rule that cannot be followed costs attention on every turn and buys nothing.

## Solution

Cut and rebind the five contradictions, leaving the sections shorter and each surviving rule followable. Sessions get one governing voice rule instead of three competing ones, an observable trigger for the structured-phase clauses, and two absolutes rebound to what they actually mean. Nothing new is added to the always-loaded surface; every story either removes text or replaces text with shorter text.

## User stories

1. **The formatting rules name one governing rule.** "Format for scan: tables and lists — use them", the cohesion clause's "bullets or tables only for genuinely parallel facts", and "one-sentence paragraphs stacked in a row read as a formatting error" are reconciled under one rule (recommend the colleague-voice bullet governs), so a three-fact update has a sanctioned shape.
   `Line: fable / high.` Guidance prose compounds through every session that loads it, which is `craft-line`'s leverage override.

2. **The progress labels are demoted to updates that carry scannable state.** The `progress` clause drops "even if the entire update fits in one sentence", which today forces two bold labels onto a one-sentence update — the exact register the section's closing bullet forbids.
   `Line: fable / high.` Same leverage override; the wording is the deliverable.

3. **The structured-phase clauses get an observable trigger.** "An in-progress Bench phase update" is replaced by a visible condition: a phase is active when a `/bench-*` command was invoked this session, and not otherwise.
   `Line: fable / high.` The trigger's precision is the whole point of the story.

4. **"NEVER assume, always verify" is rebound to what it means.** The Roles absolute becomes: never assume the reviewer's decisions, and never assume a claim the gate could check instead. The shared-rule marker list moves with the sentence.
   `Line: fable / high.` Prose bearing; `projects/benchkit.md` caches doc authoring at top/high.

5. **Invariant 2 states which half of the declaration binds.** A session cannot switch the main loop's model — only the reviewer can — so the declaration binds delegates and headless runs and is ceremony for the main loop. The invariant says so.
   `Line: fable / high.` Editing an invariant is the highest-leverage prose in the kit.

6. **The lighter-path sentence stops contradicting its table.** "You must get an explicit OK *before* skipping canonical steps" folds in the exception the table directly beneath it already grants.
   `Line: fable / high.` One sentence, but it is always-loaded guidance prose.

7. **An anchored shared rule cannot be relocated out of its owning section.** `checkSharedRuleSingleSource` matches markers against the whole file, so a rule moved from Roles into any other section — or commented out — still satisfies its anchor. This spec restructures sections, which is exactly the edit that failure mode hides. The check becomes section-scoped and comment-stripped for the markers this spec touches, with a canary fixture proving the relocation reds.
   `Line: opus / medium.` `projects/benchkit.md` routes gate and conformance logic at mid effort.

## Implementation decisions

**The wording changes are reviewed, not gate-graded.** Stories 1, 2, 3, 5, and 6 change semantics the conformance suite cannot judge — no check can decide whether a voice rule is followable. Semantic correctness rests on reviewer sign-off and the falsification pass, and this spec does not claim otherwise. What the gate can prove is narrower and is story 7's job: that nothing anchored was lost, moved, or commented out while the sections were restructured.

**Both anchor mechanisms are weaker than they look, and story 7 closes the gap before the prose edits rely on them.** The 2026-08-03 falsification pass established two defects by reading the checkers. `checkSharedRuleSingleSource` matches its marker list against the raw whole file, so a marker that is relocated to another section, or wrapped in an HTML comment, still satisfies its anchor. `checkStructuredPhaseContract` derives clause names from whatever backticks survive in the declaration line rather than pinning the four, so deleting a clause's name *and* its body together passes. Every "already covered" claim in this spec's first draft rested on strength these checks do not have. Story 7 gives them that strength — section-scoped, comment-stripped marker matching for the markers this spec touches, and a fixed four-name set for the structured-phase contract — and lands first, so the prose stories are guarded when they arrive.

**This is a strengthening, not a weakening.** No check loses reach. The Progress clause keeps a live body under the demoted wording, so the structured-phase contract stays satisfied by the intended edit and newly reds on the deletion it previously allowed.

**Story 7 writes in today's hand-written anchor format.** FT156 may replace that format with a declarative registry across roughly a hundred anchors; writing these in the current shape costs nothing extra, because the migration rewrites every anchor regardless. This spec does not wait on that ruling — with story 7 landing the strength the prose stories need, nothing here depends on how FT156 rules.

**The acceptance test for the prose is a cold read by the mid tier, and it is mandatory.** This row exists because always-loaded rules that read fine to their author drift when a session actually tries to apply them — the whole defect class is misinterpretation, not absence. So the revised sections are validated the way they fail: a fresh session at the `opus` binding reads them cold, with no memory of this spec or conversation, and reports what each rule tells it to do. The check is not "does this read well" but "does an ordinary session act on this the way the rule intends", and a rule that reads correctly only to the tier that wrote it has failed.

The tier is load-bearing rather than incidental. `fable` authors the prose under the leverage override, but `opus` is this repo's mid binding and the tier that runs ordinary sessions, so it is the representative reader. Reviewing at the authoring tier would prove the least useful thing — that the writer understood the writing.

This composes `craft-synthesis`'s third quality loop rather than inventing a procedure: the dogfood loop already asks whether a kit change survives real use. What this spec adds is naming the tier and requiring the read before sign-off, not after. Two failure signals count as findings and route back to the wording: the session paraphrases a rule into something the rule does not say, and the session cannot tell which of two surviving rules governs a case the old contradiction covered. This catches the passing degenerate directly — a section gutted to one named clause plus its markers is mechanically green and produces no actionable reading.

**Section sizes are not a target.** The cut is judged by whether each surviving rule is followable. No story lands a numeric budget.

## Testing decisions

- **External behavior a good test exercises.** That an anchored shared rule cannot vanish, move sections, or be commented out unnoticed while `.bench/BENCH.md` is restructured, and that a structured-phase clause cannot be deleted by dropping its name alongside its body.
- **Seams receiving tests.** `checkSharedRuleSingleSource` and `checkStructuredPhaseContract`, both in `internal/conformance`. Prior art: `TestStructuredPhaseContractIgnoresInactiveGuidance` already proves the comment-out, quotation, negation, and wrong-section decoys for the clause bodies, and `structured-phase-progress-anchor` is the existing canary for a section-scoped `.bench/BENCH.md` mutation — story 7's fixtures join that family and must not regress either.
- **Gate seam.** `bench gate`'s conformance suite observes stories 4 and 7; the canary phase observes story 7's fixtures.
- **The semantic seam is a fresh `opus` session, not a check.** The revised sections are handed to a cold reader at the mid binding that reports what each rule instructs; paraphrase drift and unresolvable precedence are findings. This is the only acceptance the prose stories have, so it runs before sign-off and its result is reported with the build.

### Seam diagram

    trigger: any edit to .bench/BENCH.md
        │
        ▼
    file text  ──▶  [ checkSharedRuleSingleSource ]  ──▶  diagnostics
                    [ checkStructuredPhaseContract ]
                        ◀ tests attach here: pass mutated guide text to each
                          checker directly, and mutate the real file through a
                          canary fixture in workflow-guidance-anchors

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 7 | A shared-rule marker relocated out of its owning section reds | `checkSharedRuleSingleSource` | to be observed at build time — the relocation must be demonstrated red before the check is trusted | Today the whole-file `strings.Contains` passes on any placement; a section-scoped predicate is the only thing that distinguishes a marker present in its owning section from one parked elsewhere — placement and presence, not the rule's semantics |
| 7 | A shared-rule marker wrapped in an HTML comment reds | `checkSharedRuleSingleSource` | to be observed at build time | The checker reads raw file text with no comment stripping, so an inactive rule currently satisfies its own anchor |
| 7 | Dropping a structured-phase clause name together with its body reds | `checkStructuredPhaseContract` | to be observed at build time | Clause names are derived rather than pinned, so a simultaneous deletion leaves nothing to validate; a fixed four-name set makes the absence itself the failure |
| 4 | The rebound Roles sentence is present in `.bench/BENCH.md`, inside Roles, and absent from `AGENTS.md`/`README.md` | `checkSharedRuleSingleSource` | already covered **once story 7 lands** — before that, only whole-file presence is checked | Renaming the sentence without renaming the marker reds; after story 7 the anchor also holds the sentence in the section that owns it |
| 1 | `"Clear beats dense"` survives the formatting reconciliation, in its own section | `checkSharedRuleSingleSource` | already covered once story 7 lands | The marker sits inside the reconciled bullet; losing or relocating it during the cut reds |
| 6 | `"Right-size the process"` survives the lighter-path rewording, in its own section | `checkSharedRuleSingleSource` | already covered once story 7 lands | The marker opens the edited paragraph; a rewrite that drops or moves it reds |
| 1, 2, 3 | The four structured-phase clauses stay declared with live bodies | `checkStructuredPhaseContract` | already covered for an emptied or negated body; covered for a dropped name only once story 7 lands | The demotion in story 2 keeps a live Progress body and passes, which is the intended distinction from deleting the clause |
| 5 | Invariant 2's binding-scope sentence is accurate | none | **not gate-observable** — no check can grade whether a scoping claim is true | Reviewer sign-off and the falsification pass are the only judges; recorded here rather than given a decorative row |
| 1, 2, 3, 5, 6 | A cold `opus` session reads each revised rule and reports the action it prescribes, without paraphrase drift | fresh mid-tier session, no spec or conversation context | **not gate-observable** — the finding is a misreading, which only a reader produces | The passing degenerate is a section gutted to one named clause plus its markers: mechanically green, and it yields no actionable reading, so the cold read fails where every check passes |
| 1, 6 | Where the old contradiction covered a case, the cold reader can name which surviving rule governs it | same session | **not gate-observable** — reviewed | An unresolvable precedence is the exact defect this row exists to remove; if the reader still cannot choose, the reconciliation did not happen |

### Edge inventory

- **Empty or absent input** — covered: `checkStructuredPhaseContract` returns `structuredPhaseUnavailable` on empty shared rules, proven by `TestRunConformanceDistinguishesAbsentAndEmptyInputs`.
- **Whole-section deletion** — covered: `markdownH2Section` returns empty and the checker reds with the dropped-contract diagnostic.
- **Inactive-but-present clause bodies** (commented out, quoted, negated, wrong section) — covered by the four existing decoy cases; story 7 must not regress them.
- **Inactive-but-present shared-rule markers** — story 7 row above; this is the gap the falsification pass found.
- **Duplicate section heading** — covered: `markdownH2Section` takes the first matching section, and a duplicated-heading check already exists in the docs-workflow suite.
- **Marker substring collision** — **Won't handle**: the rebound sentences are long and distinctive, and a section-scoped predicate narrows the match window further.
- **Paraphrase evasion** (a sentence reworded to evade its marker) — **Won't handle**: FT156's second face by name; adopting a stronger mechanism here would fork the ruling that row exists to make. Story 7 fixes placement and activation only, not paraphrase.
- **Concurrent edit by FT107's batch** — **Won't handle**: prevented by ordering, not by a check — the Dependencies table makes FT107 literal-blocked on this row.

## Out of scope

- **Single-sourcing the `CLI Inventory` section against `bin/bench.sh`.** No conformance check enforces the declared sync, so the list can drift silently. This is FT89's inventories half; moving it here would put one fact under two owners. `3 edits, 1 gate run`.
- **The declarative anchor registry.** FT156, roughly a hundred anchors and a 864-versus-660 structure violation. Story 7 strengthens two checks in place; converting the format is a separate build. `40+ edits, 6 gate runs`.
- **Paraphrase-resistant anchoring.** FT156's second face — substring forbids die to rewording. Independent of placement and activation. `unpriced until FT156 rules`.
- **The named-lighter-path routing decision.** Story 6 reconciles two statements as they stand; whether the light path should *exist* in its current form is FT156's grill and FT180's spec-optional route. `0 edits here`.
- **The demonstrated-delta audit over the rest of the always-loaded surface.** FT100 retains it after this row takes the "How to talk to me" slice. `12 edits, 3 gate runs`.
