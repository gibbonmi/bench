# Spec-row repair rounds: what the largest repair cause actually is

Recommendation: define `spec-row` and split it, because today the label covers three different failures with three different fixes. Then decide one policy question: do edge cases belong in the spec stage or in the build's mutation probes? The current rubric sends them to the build, so about a third of the `spec-row` rounds are the designed path, not a spec defect.

`spec-row` is the largest repair cause since August, with about 178 rounds, but no guidance file defines it. In the rounds that two readers classified, about 40 percent were edge cases that a build review found and no row decided. About 30 percent were not spec defects: the row was clear, and the code or its test missed it. The remaining 20 to 35 percent were defects in the spec text, and rules that were in force at the time already required those catches.

Scope: the `spec-row` rounds in implementation retros from 2026-09-15 to 2026-10-19, with monthly rates from August
Evidence status: 66 and 45 rounds of 78 classified by two independent readers against the review records and repair commits; the class of each round is a judgment
Source repository: /home/mgibs/workspace/bench
Source commit: b29c4df2d
Retrieved: 2026-10-19
Consumed by: the reviewer's decision on the repair-attribution vocabulary in `/bench-final-check` and on the edge rule in `craft-spec`
Drift: refresh when the repair-attribution vocabulary, the `craft-spec` review rubric, or the hostile-input checklist changes
Retire when: an approved vocabulary change lands and one month of retros uses it

This report changes no command, skill, or check. It proposes changes for the reviewer to decide.

## Question graph

1. Where is `spec-row` defined, and how do retros apply it?
2. What defect does each `spec-row` round actually repair?
3. Which review axis or check found each defect, and in which round?
4. Could the spec stage have caught each class, and did a rule already require the catch?
5. Did the rate change over time, and did any rule change move it?
6. What should Bench change?

Questions 1 to 5 hold the facts. Question 6 holds the proposals.

## Method

The coordinator extracted the repair-attribution table of every retro in git history. The table names each ticket, its repair rounds, and one cause for each round. 28 retros from 2026-09-15 to 2026-10-19 hold 78 `spec-row` rounds. Two read-only readers on opus at high effort received the same charge and the same dataset. Neither reader saw the other's return.

For each round, each reader found the review record at the parent of its retirement commit, the review finding, and the repair or plan commit that answered it. Reader A classified 66 rounds in 22 retros. Reader B classified 45 rounds in 17 retros from the full finding text, and 33 more from headlines only. The coordinator reopened the cited guidance lines and closed the edge-inventory question with a targeted read.

## 1. Definition

### Facts

- **No definition:** `.agents/commands/bench-final-check.md:132` lists `spec-row` in the cause vocabulary. Lines 134 to 137 define only `one-source` and `check-gap`, which `5013e740d` added on 2026-09-27.
- **Origin:** the term arrived in `d12728916` on 2026-08-07, from the one-shot-feedback spec, which also lists it with no definition. No glossary, reference, or `docs/` file defines it.
- **Working meaning:** retro prose uses it as "the coverage map did not ask for the behavior the repair added" (`c9c59cda8`, ft229-hygiene-batch) or "incomplete input-family coverage" (`7d6335f88`, test-determinism).
- **Loose use:** retros also charge `spec-row` to rounds where the Spec axis passed and only a test was weak. Examples are commitment-delivery-integrity FD-C1 and FD-C2, and durable-caller-migration D-C1. One retro charges it with no spec at all (`61c885586`, light-path-wrapper-home-guard).
- **Roadmap:** no roadmap row mentions `spec-row` or the repair-attribution vocabulary.

## 2. Defect class

### Facts

| Class | Reader A (66 rounds) | Reader B (45 rounds, full text) |
|---|---|---|
| Missing edge: a review found an input, state, or refusal that no row decided | 24 (36%) | 20 (44%) |
| Not a spec defect: the row was clear, and the test stayed silent or the code missed the row | 18 (27%) | 14 (31%) |
| Spec text defect: a wrong, ambiguous, or unowned row, a stale premise, an incomplete sweep, or a drifted table cell | 23 (35%) | 8 (18%) |
| Review-added scope beyond the approved spec | 1 (2%) | 2 (4%) |
| Unknown | 0 | 1 (2%) |

The readers agree on the order of the first two classes. They disagree on several single rounds, because a round that closed several findings fits more than one class. For example, reader B classed closure-removal ticket 02 round 1 as a missing edge, and reader A classed it as a silent test. Reader A split the spec text defects into finer classes and found more of them.

Examples:

- **Missing edge:** ref-inventory symbolic refs (rows RI83, RI93, RI102, RI103; `8b7927733`, `dd3c188a7`). preflight-ownership-closure ticket 03 took four rounds for Go write forms such as elided pointer literals, generic aliases, and conversions (`34548945f`, `a88a5a653`). commit-merge-parent needed a glob-shaped merge name (`cca1013b0`).
- **Not a spec defect:** commitment-delivery-integrity FD-C1, where the code met the row and a test stayed silent under a probe. roadmap-delivery-commitment DC-C1, where the code broke the explicit rules at spec lines 158, 164, and 464.
- **Spec text defect:** build-landing-window BW-D, where spec line 266 contradicts stories 27 and 30 (`814a14cce`). probe-mutant-sets PM-C2, where the spec states no precedence. roadmap-delivery-commitment DC-C9, which assumed a fresh link has a `ROADMAP.md`. build-directions BD62, where a live sentence still stated the retired rule (`4633b5c57`).

Missing-edge rounds cluster in two families: ref-inventory symbolic refs (6 rounds) and preflight-ownership-closure Go write forms (4 rounds). Both families come from rows that promise completeness over an open set of inputs.

## 3. Finder

### Facts

- **Spec axis:** found about 24 to 30 rounds.
- **Coverage axis:** found about 17 to 19 rounds, often through a silent `bench probe` mutation.
- **Standards axis:** found 2 to 3 rounds.
- **Gate and author:** found none of the classified rounds.
- **First round:** about 75 to 80 percent surfaced in the first review round.

The confirming-round findings cluster in the open families, where each repair exposed the next edge. Examples are preflight-ownership-closure rounds 2 to 6, ref-inventory C2b, and roadmap-delivery-commitment ticket 09.

## 4. Prevention

### Facts

- **Spec review ran every time:** reader B read the verification log of all 28 specs, and each spec passed a spec review before its build. The review lines ranged from sonnet at high effort to fable or GPT at extra-high effort.
- **The rubric defers edges by design:** `.agents/skills/bench-craft-spec/SKILL.md:85` states that "a degenerate that needs deliberate contrivance is the build's mutation-probe target, never a new spec row". It also states that a revision "may not add a promise beyond the decision source". `1f18922e6` carried this text on 2026-08-29.
- **Many missed classes were already on the checklist:** the edge inventory attaches the profile's hostile-input checklist (`SKILL.md:54`). Glob characters, control bytes, and symlinks were on that checklist before the builds that missed them (`projects/benchkit.md:182-231`).
- **Symbolic refs were not on the checklist yet:** the ref-inventory edge inventory at its first build commit `2d3535b6b` has no symbolic-ref line. `f19afe1ac` added symbolic refs to the checklist about two and a half hours after that commit. The checklist grows from build findings.
- **Rules for spec text defects were in force:** the rubric already required each catch for these classes, at the time of each miss:
  - "every source behavior becomes a row or an exception" (`SKILL.md:48`);
  - "Does every source behavior have a red-capable row?" (`SKILL.md:79`);
  - the exact-predicate question (`SKILL.md:81`);
  - claim verification against current code (`SKILL.md:18`, since `de2e743d4`);
  - the reader sweep (`SKILL.md:16`).

### Inference

| Class | Can the spec stage catch it? | What would catch it |
|---|---|---|
| Missing edge in a checklist class | yes, and the rules require it | a recorded walk that maps each checklist class to a row or a Won't handle for each surface, graded by spec review |
| Missing edge in an open family | partly | a rule that a row promising completeness over an open family states a fail-closed bound or a Won't handle |
| Missing edge that needs contrivance | no, by design | the build's mutation probe, as `SKILL.md:85` states |
| Spec text defect | yes, and the rules require it | mechanical support: a trace from each prose clause to a row, and derived seam and chunk cells |
| Clear row missed by the code or the test | no | this is `delegate-error` |
| Review-added scope | no | this is plan expansion, which `.bench/BENCH.md` permits |

## 5. Trend

### Facts

| Month | `spec-row` rounds per ticket | All rounds per ticket | `spec-row` share of rounds |
|---|---|---|---|
| August | 0.20 to 0.21 | 0.57 | 36% |
| September | 0.22 to 0.27 | 0.80 | 33 to 34% |
| October | 0.26 to 0.27 | 1.17 | 23% |

The three parses of the retro set (the coordinator's and each reader's) differ in retro selection: 124 to 139 retros, 567 to 592 rounds, and 177 to 179 `spec-row` rounds. The `spec-row` count is stable across the parses.

### Inference

The `spec-row` rate per ticket did not fall. Its share fell because all rounds per ticket doubled, and because `one-source` and `check-gap` began on 2026-09-27. No spec-stage rule change coincides with a lower rate. The changes that the readers checked include:

- the reader sweep (`a9e6cbed0`, 2026-09-01);
- claim verification and symbolic refs (2026-09-27);
- premise reads (`5fef6ad5d`, 2026-09-28);
- the grader trace (`3f6f99d6c`, 2026-10-06);
- the spec-alone review (`f191167a7`, 2026-10-10).

The doubling of all rounds per ticket from August to October is the larger signal. In October, `one-source` caused 85 of 205 rounds.

## 6. Proposals

- **S1. Define and split the label.** In the repair-attribution paragraph of `/bench-final-check`, define `spec-row` as a defect in the spec text: a wrong, ambiguous, or unowned row, a stale premise, or an incomplete sweep. Add one term for an edge that a build review found and no row decided, for example `edge-found`. Send a missed clear row to `delegate-error`, and include a test that stays silent under a probe, because `check-gap` covers only a ticket check that omitted a package or a check. The paragraph is the single source of the vocabulary, and its anchors-registry needle changes with it. The table then shows which of the three fixes each round needs.
- **S2. Decide where edges belong.** The reviewer decides between two options. The first option keeps `SKILL.md:85` and accepts edge rounds as the designed cost: about one round for every eight tickets. The second option makes the spec review grade a recorded checklist walk for each surface, so that checklist-class edges become rows or Won't-handle lines before the build. The first option costs nothing now. The second option adds spec-stage work but targets the misses that the checklist already names.
- **S3. Bound open families at spec time.** Add one rubric question: a row that promises completeness over an open family of inputs, such as syntax forms or ref shapes, states a fail-closed rule or a Won't-handle bound. That bound stops the chain of confirming rounds that preflight-ownership-closure and ref-inventory paid.
- **S4. Give mechanical support to spec text checks.** The rules for spec text defects exist and were in force, so another rule is unlikely to help. Candidates are a trace from each quantified prose clause to a row, and seam and chunk cells derived from one source. Each candidate needs a spec.
- **S5. Assess `one-source` next.** It is the largest cause in October and the main driver of the doubled round count. It is a separate assessment.

## Threats to validity

- **Judgment classes:** each round's class is a reader's judgment, and the two readers disagree on several single rounds.
- **Mixed rounds:** a round that closed several findings has one cause term, so its class follows the finding that drove the amendment.
- **Missing rounds:** 12 rounds were not classified by reader A, and 33 rounds were classified only from headlines by reader B.
- **Retro selection:** the three parses select retros differently, so the rates are ranges.
- **Rate and rules move together:** the rules and the review lines changed during the period, so no causal claim about a rule holds.

## What was read and what was not

The readers read these sources:

- the review records at the parent of each retirement commit;
- the repair and plan commit subjects and diffs;
- the verification log of each spec;
- the history of `bench-final-check.md`, `craft-spec`, `bench-write-spec.md`, and the profile checklist.

The coordinator reopened these sources:

- `SKILL.md:50-56` and `SKILL.md:76-86`;
- `bench-final-check.md:129-142`;
- the checklist lines in `projects/benchkit.md`;
- the ref-inventory edge inventory at `2d3535b6b`.

Nobody read these sources:

- the spec-review transcripts;
- the edge inventory of each spec, except ref-inventory;
- the repair commits of the headline-only rounds.

## Verification record

- [x] The output opens with the recommendation, the scope, and the evidence status.
- [x] It keeps facts, inferences, tested results, and proposals in separate sections.
- [x] It holds one section for each question.
- [x] It uses tables where the comparison matters.
- [x] It keeps each contradiction and each residual unknown.
- [x] It states what could not be verified.
- [x] It ends with a validation plan.
- [x] Each load-bearing claim names a commit, a path, or a recorded measurement.

## Validation plan

1. **Reclassify with the new terms.** After S1 lands, ask two readers to reclassify these 78 rounds with the new terms. Record their agreement, because the new terms must agree better than the current classes do.
2. **Test the checklist walk.** For S2, replay the spec review of commit-merge-parent and durable-caller-migration with a required checklist walk. Record whether the walk finds the glob-name and control-byte edges that their builds found.
3. **Test the open-family rule.** For S3, apply the open-family question to the next spec that promises completeness over an input family. Count its confirming rounds.
4. **Watch one month of retros.** After S1, compare the `spec-row` and `edge-found` rates per ticket with this report's October baseline.
