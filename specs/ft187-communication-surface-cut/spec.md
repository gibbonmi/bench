# FT187 — cut the communication surface to the rules that fire

Status: implemented

Decision source: reviewer-confirmed current conversation, 2026-08-03 and 2026-08-06 — the reviewer reviewed six findings against `.bench/BENCH.md`, accepted the direction ("i like the changes"), directed the roadmap row and this spec, then required story 7 to compose with FT156's landed declarative anchor registry rather than predate it.

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

4. **"NEVER assume, always verify" is rebound to what it means.** The Roles absolute becomes: "Never assume the reviewer's decisions, and never assume a claim the gate could check instead." The shared-rule marker list moves with the sentence.
   `Line: fable / high.` Prose bearing; `projects/benchkit.md` caches doc authoring at top/high.

5. **Invariant 2 states which half of the declaration binds.** A session cannot switch the main loop's model — only the reviewer can — so the declaration binds delegates and headless runs and is ceremony for the main loop. The invariant says so.
   `Line: fable / high.` Editing an invariant is the highest-leverage prose in the kit.

6. **The lighter-path sentence stops contradicting its table.** "You must get an explicit OK *before* skipping canonical steps" folds in the exception the table directly beneath it already grants.
   `Line: fable / high.` One sentence, but it is always-loaded guidance prose.

7. **Every shared-rule marker this cut edits is registered in its owning section.** FT156 landed the section-scoped, comment-stripped matcher and declarative `internal/anchors` registry, but the three markers this spec edits still live only in `checkSharedRuleSingleSource`'s whole-file list. Register the Roles, How to talk to me, and Workflow markers as `RequireInSection` rows, single-source each marker value between the registry and the duplication check, and pin the structured-phase contract to its fixed four clause names in the bespoke parser FT156 deliberately left outside the registry.
   `Line: opus / medium.` `projects/benchkit.md` routes gate and conformance logic at mid effort.

## Implementation decisions

**The wording changes are reviewed, not gate-graded.** Stories 1, 2, 3, 5, and 6 change semantics the conformance suite cannot judge — no check can decide whether a voice rule is followable. Semantic correctness rests on reviewer sign-off and the falsification pass, and this spec does not claim otherwise. What the gate can prove is narrower and is story 7's job: that nothing anchored was lost, moved, or commented out while the sections were restructured.

**FT156 owns uniform anchors; the structured-phase parser owns its grammar.** FT156 landed `internal/anchors` as the one declarative table and shared matcher for uniform needle anchors. Registry evaluation already strips HTML comments before resolving whole-file or section-scoped requirements, and its existing relocation canary proves that a `RequireInSection` needle cannot survive by moving elsewhere in `.bench/BENCH.md`. Story 7 therefore adds registry rows instead of teaching `checkSharedRuleSingleSource` a second placement matcher. The bespoke check retains canonical-source duplication enforcement for `AGENTS.md` and `README.md` and consumes the same exported marker values as the registry rows. Separately, `checkStructuredPhaseContract` remains the bespoke grammar owner FT156 explicitly left behind; story 7 replaces its derived clause-name set with the fixed `progress`, `exit`, `omission`, and `cohesion` contract. Story 7 lands first, so the prose stories are guarded when they arrive.

**This is a strengthening, not a weakening.** No check loses reach. The Progress clause keeps a live body under the demoted wording, so the structured-phase contract stays satisfied by the intended edit and newly reds on the deletion it previously allowed.

**Story 7 adds no second anchor inventory.** The three shared marker values are exported once from `internal/anchors`, referenced by their registry rows and by `checkSharedRuleSingleSource`, following FT156's existing `FixDontParkMarker` and `SourceWarrantMarker` precedent. An independent test expectation enumerates the three new `(file, kind, section, needle)` tuples only because omitting a row would otherwise stay green; this is the working agreement's named omission-mutation exception, not a production copy of the registry.

| file | kind | section | needle |
|---|---|---|---|
| `.bench/BENCH.md` | `RequireInSection` | `Roles` | `Never assume the reviewer's decisions, and never assume a claim the gate could check instead` |
| `.bench/BENCH.md` | `RequireInSection` | `How to talk to me` | `Clear beats dense` |
| `.bench/BENCH.md` | `RequireInSection` | `Workflow` | `Right-size the process` |

**The acceptance test for the prose is a cold `opus / low` read, and it is mandatory.** This row exists because always-loaded rules that read fine to their author drift when a session actually tries to apply them — the whole defect class is misinterpretation, not absence. So the revised sections are validated the way they fail: a fresh cross-harness session at `opus / low` reads them cold, with no memory of this spec or conversation, and reports what each rule tells it to do. The check is not "does this read well" but "does an ordinary session act on this the way the rule intends", and a rule that reads correctly only to the tier that wrote it has failed.

The tier is load-bearing rather than incidental. `fable` authors the prose under the leverage override, but `opus` is this repo's mid binding and the tier that runs ordinary sessions, so it is the representative reader. Reviewing at the authoring tier would prove the least useful thing — that the writer understood the writing.

This composes `craft-synthesis`'s third quality loop rather than inventing a procedure: the dogfood loop already asks whether a kit change survives real use. What this spec adds is naming the tier and requiring the read before sign-off, not after. Two failure signals count as findings and route back to the wording: the session paraphrases a rule into something the rule does not say, and the session cannot tell which of two surviving rules governs a case the old contradiction covered. This catches the passing degenerate directly — a section gutted to one named clause plus its markers is mechanically green and produces no actionable reading.

The cold reader receives these cases in order and reports the action each revised rule prescribes; the expected action is the acceptance oracle, not an invitation to choose an easy example:

| case | expected action |
|---|---|
| A three-fact update, tried once with genuinely parallel facts and once with dependent facts | Use one scannable list or table for the parallel facts and cohesive prose for the dependent facts; the colleague-voice rule governs both, and no rule demands stacked one-sentence paragraphs. |
| A one-sentence routine acknowledgement with no meaningful intermediate state or continued work | Use plain prose without `Status:` or `Next:` labels; the progress clause fires only when the update carries both meaningful state and continued work. |
| The same status update before and after a `$bench-*` phase is invoked in the session | Use ordinary conversation before invocation and the structured phase clauses after invocation while that phase remains active. |
| A multi-cycle main-session stage versus launching a delegate or headless run | Do not pretend a declaration can switch the already-running main model; declare the binding line for the delegate or headless run it governs. |
| A change satisfying both observables in the lighter-path table versus one that does not | Take the standing lighter-path approval without asking again when both observables hold; otherwise ask before skipping the canonical phases. |

**The surface cut carries a mechanical shrink receipt.** For stories 1–6, the build records one story-scoped `.bench/BENCH.md` passage with its before/after byte counts; adjacent passages may share a unified diff hunk. Every passage must be a deletion or a shorter replacement, and the six passages may add no seventh guidance rule. Reviewer interpretation still grades whether the shorter text says the right thing; the receipt prevents additive prose from passing under the claim that this is a cut.

**Section sizes are not a target.** The cut is judged by whether each surviving rule is followable. No story lands a numeric budget.

## Testing decisions

- **External behavior a good test exercises.** That every marker this spec edits is a section-scoped registry member, that the shared matcher still rejects relocated or commented-out requirements, and that a structured-phase clause cannot be deleted by dropping its name alongside its body.
- **Seams receiving tests.** Registry enumeration in `internal/anchors` receives the three-row omission test; the existing `fix-dont-park-section-relocated` and `commented-required-anchor` canaries retain matcher bite; `checkStructuredPhaseContract` receives the fixed-name test. Prior art: `TestStructuredPhaseContractIgnoresInactiveGuidance` already proves the comment-out, quotation, negation, and wrong-section decoys for clause bodies, and `structured-phase-progress-anchor` exercises the real `.bench/BENCH.md` mutation.
- **Gate seam.** `bench gate`'s conformance suite observes the registry membership, shared-rule duplication, and fixed structured-phase contract; the canary phase retains FT156's relocation and comment-stripping bite.
- **The semantic seam is a fresh `opus / low` session, not a check.** The revised sections are handed to a cold cross-harness reader that reports what each rule instructs; paraphrase drift and unresolvable precedence are findings. This is the only acceptance the prose stories have, so it runs before implementation review sign-off and its result is reported with the build.

### Seam diagram

    trigger: any edit to .bench/BENCH.md
        │
        ▼
    file text  ──▶  [ internal/anchors registry + shared matcher ]  ──▶  placement diagnostics
               └──▶  [ checkSharedRuleSingleSource ]               ──▶  duplication diagnostics
               └──▶  [ checkStructuredPhaseContract ]             ──▶  clause diagnostics
                           ◀ tests attach here: enumerate the three registry
                             tuples, drive the bespoke clause parser, and retain
                             FT156's real-tree relocation/comment canaries

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 7 | The Roles, How to talk to me, and Workflow markers this spec edits are `RequireInSection` registry rows with their exact owning sections | `internal/anchors` registry enumeration | to be observed at build time — write the three-tuple expectation first and demonstrate that omitting each row reds it | FT156 already proves matcher behavior; the remaining cheapest wrong implementation is adding zero or only some of these rows while leaving the bespoke whole-file list green, so the independent omission expectation enumerates all three |
| 7 | A section-scoped required marker cannot survive by moving elsewhere in `.bench/BENCH.md` or by being wrapped in an HTML comment | shared `internal/anchors` evaluator through the real graded root | already covered by FT156's `fix-dont-park-section-relocated` and `commented-required-anchor` canaries plus the registry kind expectation above | The canaries prove the one matcher bites on relocation and inactive text; the tuple row proves each FT187 marker is routed through that matcher rather than a private path |
| 7 | The structured-phase declaration names exactly `progress`, `exit`, `omission`, and `cohesion`, once each; deleting, renaming, duplicating, or adding a name reds | `checkStructuredPhaseContract` | to be observed at build time — table-drive the four omissions plus duplicate and unknown-name cases before replacing the derived set | Clause names are currently derived from the document, so a self-consistent deletion or addition defines its own passing contract; an independent exact-set expectation catches every cheaper grammar |
| 4 | The rebound Roles sentence is present in `.bench/BENCH.md`, inside Roles, and absent from `AGENTS.md`/`README.md` | Roles registry row + `checkSharedRuleSingleSource` | already covered once story 7 lands | The shared exported marker feeds placement and duplicate-source enforcement without a second marker definition |
| 1 | `"Clear beats dense"` survives the formatting reconciliation inside How to talk to me | How to talk to me registry row + `checkSharedRuleSingleSource` | already covered once story 7 lands | Losing or relocating the one shared marker fails registry evaluation; duplicating it into project-owned guidance fails the bespoke source check |
| 6 | `"Right-size the process"` survives the lighter-path rewording inside Workflow | Workflow registry row + `checkSharedRuleSingleSource` | already covered once story 7 lands | Losing or relocating the one shared marker fails registry evaluation; duplicating it into project-owned guidance fails the bespoke source check |
| 1, 2, 3 | The four structured-phase clauses stay declared with live bodies | `checkStructuredPhaseContract` | already covered for an emptied or negated body; covered for a dropped name only once story 7 lands | The demotion in story 2 keeps a live Progress body and passes, which is the intended distinction from deleting the clause |
| 5 | Invariant 2's binding-scope sentence is accurate | none | **not gate-observable** — no check can grade whether a scoping claim is true | Reviewer sign-off and the falsification pass are the only judges; recorded here rather than given a decorative row |
| 1, 2, 3, 5, 6 | A cold `opus / low` session reads each revised rule and reports the action it prescribes, without paraphrase drift | fresh cross-harness session, no spec or conversation context | **not gate-observable** — the finding is a misreading, which only a reader produces | The passing degenerate is a section gutted to one named clause plus its markers: mechanically green, and it yields no actionable reading, so the cold read fails where every check passes |
| 1–6 | The cold reader returns the specified expected action for each of the five enumerated conflict cases | same session | **not gate-observable** — reviewed against the case table above | The fixed cases prevent the reader from selecting an easy example, and each expected action makes paraphrase drift or unresolved precedence a concrete finding |
| 1–6 | Each story-scoped `.bench/BENCH.md` passage is a deletion or shorter replacement, and the six passages add no seventh guidance rule | reviewer-inspected diff plus before/after byte receipt per story | **not gate-observable** — the build reports six measured passages and review rejects any non-shrinking or additive passage | An implementation that preserves the markers while adding more standing guidance violates the solution even if every mechanical anchor stays green |

### Edge inventory

- **Error path** — covered: missing registry rows fail the independent three-tuple expectation; missing, duplicate, or unknown structured-phase names fail with clause diagnostics rather than an empty green.
- **Empty or absent input** — covered: `checkStructuredPhaseContract` returns `structuredPhaseUnavailable` on empty shared rules, proven by `TestRunConformanceDistinguishesAbsentAndEmptyInputs`; the registry already distinguishes a missing anchored file from a present file missing its needle.
- **Boundary values** — covered: the exact-set row enumerates all four required names and exercises zero through four by omission, then five with an unknown name; duplicate-name behavior is explicit rather than folded into count alone.
- **Malformed input** — covered for whole-section deletion, duplicate H2 headings, unterminated HTML comments, and inactive clause decoys, which already fail closed; fenced headings remain ignored as headings by the recorded parser test. **Won't handle:** changing the existing unclosed-fence-to-end-of-file interpretation — story 7 adds registry members and a fixed clause set, not new markdown grammar.
- **Interrupted or partial state** — covered by the malformed and missing-input cases: these checks hold no mutable state, so the only observable partial state is partially written document text, which fails its marker or exact-clause contract.
- **Re-run idempotency and process-boundary lifecycle** — **Won't handle**: registry and parser evaluation are pure reads over tracked text and persist no state; a fresh process receives the same compiled registry and file bytes.
- **Hostile environment** — **Won't handle**: story 7 adds no path discovery, command input, renderer field, network, prompt, timing, or filesystem write; it registers kit-authored literals against the existing fixed `.bench/BENCH.md` path. Existing anchor-reader special-file posture is unchanged.
- **Inactive-but-present clause bodies** — covered by the existing comment, quotation, negation, and wrong-section cases; story 7 must not regress them.
- **Inactive-but-present shared-rule markers** — covered only for HTML comments and wrong-section relocation by FT156's `commented-required-anchor` and `fix-dont-park-section-relocated` canaries; story 7's registry-kind expectation proves the edited markers take that shared path. **Won't handle:** a marker quoted or negated in otherwise active prose, because the uniform anchor contract is deliberately substring-and-placement only and semantic-context matching belongs to FT156's stronger-matcher deferral.
- **Marker substring collision** — **Won't handle**: the rebound sentences are long and distinctive, and a section-scoped predicate narrows the match window further.
- **Paraphrase evasion** — **Won't handle**: FT156 explicitly deferred stronger-than-substring matching. Story 7 registers placement and activation guarantees only.
- **Concurrent FT107 prose edits** — **Won't handle**: `ROADMAP.md`'s literal dependency table blocks FT107 on FT187, so the ownership conflict is prevented by ordering rather than document matching.

## Out of scope

- **Single-sourcing the `CLI Inventory` section against `bin/bench.sh`.** No conformance check enforces the declared sync, so the list can drift silently. This is FT89's inventories half; moving it here would put one fact under two owners. `3 edits, 1 gate run`.
- **Paraphrase-resistant anchoring.** FT156's closed deferral: stronger-than-substring matching is independent of placement and activation and remains unpriced until real paraphrase-evasion evidence reopens it.
- **The named-lighter-path routing decision.** Story 6 reconciles two statements as they stand; FT144 owns the named-lighter-path ruling and FT180 owns the spec-optional route. This spec changes neither capability. `0 FT187 edits, 0 FT187 gate runs`.
- **The demonstrated-delta audit over the rest of the always-loaded surface.** FT100 retains it after this row takes the "How to talk to me" slice. `12 edits, 3 gate runs`.
