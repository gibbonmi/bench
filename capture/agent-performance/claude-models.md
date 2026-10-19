# Claude model scorecard

Last incorporated landing: `degraded-guard-refusal` build (landing `bd7beca66989d7e56a4b515775f67b8682c44cb7`, 2026-10-19).

The `degraded-guard-refusal` build (FT366) ran one ticket in one chunk on Claude Code as a `--full --auto-approve` run with an author limit of 1. The spec declared a Codex author line, so the reviewer directed Opus/high for the author and the repair. An Opus orchestrator ran two Sonnet staleness slices, and Opus/medium ran every review axis by reviewer direction. No consultation question arose, and the chunk used 2 repair cycles.

Provider usage and charges remain unknown. The landing census counted 18 raw calls: rg 16, cp 1, python3 1.

The `closure-removal` build (FT293) ran six tickets in four chunks on Claude Code as a `--full --auto-approve` run with an author limit of 1. An Opus orchestrator dispatched fresh Opus/high authors and eleven Opus/high repair sessions, and Opus/high ran every review axis by reviewer direction. Four Sonnet staleness slices opened the build, and one Fable/medium and Sonnet/high consultation pair agreed on the only judgment question. The chunks used 3, 2, 2, and 3 repair cycles, and `--auto-approve` covered the third cycles.

Provider usage and charges remain unknown. The landing census counted 24 raw calls: rg 16, cp 4, cat 2, cd 1, sed 1.

The `closure-removal` spec stage ran on Claude Code for FT293. An Opus coordinator ran eight paired Haiku/high research readers and eight paired Haiku/high slicing readers by reviewer direction, and it sliced the six tickets itself. An Opus fork wrote the spec, and the coordinator dropped the fork's staged-spec sweep before review. Fable/high ran the spec review in 1 iteration and the ticket review in 2, by reviewer direction. The ticket review caught a deleted file whose in-package callers the slicer had not searched.

Provider usage and charges remain unknown. The harness reported about 150k to 270k tokens for each research and slicing reader.

The `plan-commit` spec stage ran on Claude Code for FT380. An Opus coordinator ran four Sonnet/medium and Haiku/high research pairs and two slicing pairs, and it sliced the three tickets itself. An Opus fork wrote the spec. Opus/high and Sonnet/high reviewed both rounds by reviewer direction. In the spec round, Opus/high rejected iteration 1 with three blocking gaps, and Sonnet/high accepted it. Both accepted iteration 2 and the ticket round in 1.

The research pairs disagreed twice, and a targeted read closed each disagreement. A Sonnet reader denied that the Markdown lane runs the coverage parser, and the code showed that it does. A Haiku reader found the `./cmd/bench` floor element that its Sonnet pair missed. Provider usage and charges remain unknown. The harness reported about 110k to 290k tokens for each reader.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 42 landings + implementer, 10 charges + reviewer, 48 axes | On `repair-confirmation` Fable/high orchestrated four chunks under `--auto-approve` and added four spec guarantees, each flagged for veto. It retried a stale-digest refusal three times and a `check-floor` red three times before each plan commit closed it. | Coordination of a delegated build and adversarial spec review; review axes in parallel by reviewer direction; it implements only when the reviewer names it | orchestrator 0.36 over 3 pairs, 1 abstained; medium reviewer 0.256 over 20 pairs, 1 abstained |
| Fable / high, xhigh, medium | reviewer, 124 axes + 4 blind comparisons + spec reviewer, 10 passes + ticket reviewer, 4 passes + verifier, 3 light-path checks + control, 1 full-retrieval run + decider, 26 consultations + debug, 2 charges | On `closure-removal` one Fable/medium consultation agreed with Sonnet/high to accept the refusal of a root headroom premise, so no tie-breaker ran. It cited the structure engine's root spelling as the reason. | Review axes, spec review, and ticket review by reviewer direction; consultations whenever the coordinator is uncertain; the RE2-style full control | reviewer: 0.168 over 32 pairs, 1 abstained |
| Opus / high, xhigh, max | implementer, latest 10 of 103 current attempts, 3 of them max shadows, 1 a max light-path fix, and 2 xhigh light-path fixes; reviewer, 319 earlier axes + 42 `closure-removal` axes + 3 consultations; researcher, 20 paired charges; repair, latest 10 of 91 current attempts; orchestrator, 6 builds and 4 spec stages; spec author, slicer, and reviewer, latest 10 passes, 7 of them xhigh slicers | On `degraded-guard-refusal` the Opus/high author met all 41 original rows in one pass with biting probes, but the review found seven repair targets, four of them one-source defects behind a narrow fence. One Opus/high repair closed all seven in one cycle, and every claim held. | High for first-round review axes, orchestration, consultations, and debug runs; xhigh for a code ticket where a repair round costs more than about 100k tokens, or that crosses owner seams; repairs at the ticket's effort. | implementation: 0.123 over 99 pairs, 3 abstained, and 0.01 over 8 xhigh pairs; orchestrator 0.26 over 3 pairs; reviewer: 0.162 over 74 pairs, 9 abstained; repair: 0.029 over 31 pairs; consultant 0.01 over 1 pair |
| Opus / medium, low | implementer, latest 10 of 23 tickets + 5 benchmark shadows; repair, latest 10 of 35 sessions; orchestrator, 8 builds; reviewer, latest 10 of 88 axes | On `degraded-guard-refusal` Opus/medium axes found a silent wrapper-continuation mutation and three duplicated test helpers, and each confirming round closed its findings. The Opus orchestrator caused one extra review round and one wasted re-verification through a spec correction after the chunk freeze. | Low for an exact-spec guidance ticket; medium as the default for a code ticket, for review axes by reviewer direction, and for a mechanical test migration; a repair keeps the ticket's declared effort. | 0.087 over 116 pairs; 28 abstained; reviewer 0.156 over 45 pairs; orchestrator 0.287 over 3 pairs |
| Sonnet / high | orchestrator, 3 landings + reviewer, 132 axes + spec reviewer, 5 passes + researcher, 3 charges + implementer, 9 comparison tickets + staleness, 4 slices | On `degraded-guard-refusal` two Sonnet staleness slices agreed that ten fence paths were obsolete and found five drifted citations, and each named what it did not verify. The orchestrator accepted every finding without a tie-breaker. | First-round and confirming review axes, consultations, staleness slices, and fix reviews by reviewer direction, and orchestration; not fresh ticket authorship | high 0.145 over 2 pairs; xhigh 0.40 over 2 pairs; reviewer 0.203 over 4 pairs; consultant 0.09 over 1 pair |
| Sonnet / high (research and polish) | researcher, 10 paired charges + 5 single charges + 15 slicing-research charges + polish, 1 session | On `repair-confirmation` Sonnet delegates ran the build research by reviewer direction. The build record keeps no result for each slice, so this landing adds no quality evidence. | Default research reader at medium effort, paired with Haiku/high by reviewer direction on 2026-10-11; polish of a cheap draft when the reviewer names it | research 6 of 6 earlier slices clean; slicing research missed the closure classes in 6 of 6 uncharged slices and in 0 of 3 charged slices |
| Haiku / default–xhigh | researcher, 16 charges + 8 slicing-research charges + implementer, 2 comparison tickets + auditor, 3 charges | On `guidance-weight` the Haiku/xhigh readers matched the Sonnet closure sets in each slice. They also found the B42 drain duplicate and the renamed canary EXPECT of an unmoved row. | Research reader at high effort by default, paired with Sonnet/medium; staleness audits; not ticket authorship | research 10 of 12 slices clean; slicing research agreed with Sonnet in 3 of 3 chunks; authorship 0 landings |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing | unknown |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `worktree-seam-reduction` comparison | Opus and Sonnet / high / implementer; Fable / high / reviewer | One blind run ranked Sonnet-sliced first, Opus-sliced second, and Opus-full third. The reviewer selected the second arm, and all three arm gates passed. |
| `markdown-block-reader` spec review | Sonnet / high, xhigh; Opus / xhigh; Fable / high, xhigh / spec reviewer | Five reviewers graded one commit with one charge. List costs were Sonnet high $2.14, Fable high $1.98, Opus xhigh $2.44, Fable xhigh $2.50, and Sonnet xhigh $6.31. No reviewer found every confirmed finding, and only Fable high reported a false one. |
| `build-directions` effort benchmark | Opus / low, medium, high, xhigh / implementer; Fable / high / blind grader | Ten first-pass commits on three tickets met every acceptance row. On two guidance tickets low matched high at half the tokens. On the code ticket xhigh ranked first, medium second, and high last of four. |
| `commit-merge-parent` effort assessment | Opus / xhigh, max / implementer | Three max shadows matched the xhigh design on every code ticket at 1.41 to 1.44 times the tokens and up to 2.3 times the time. Eight review rounds found 17 findings on the xhigh trees, 10 of them repaired. |
| `review-evidence-file-pages` RE2 control | claude-fable-5-1 / high / reviewer | The narrow round found 8 findings in 399474 tokens and the full control found 6 in 227918 tokens. They agreed on 3, so neither shape covered the other. |

## Current decisions

- Read this scorecard before each line declaration, because its decisions carry the orchestration rules that a build repeats. `degraded-guard-refusal` skipped it and repeated three lessons below.
- When a spec `Line:` names a model that this harness cannot run, ask the reviewer for the author line before the first dispatch, and record the direction in the assignment commit.
- Declare the expected repair rounds for each chunk with a confidence, so the orchestrator claim gets a label. Expect two to three repair cycles per chunk when a chunk deletes a shared seam; `closure-removal` used 3, 2, 2, and 3.
- Route every coding delegate to Opus/high and every small-fix delegate to Opus/low, by reviewer direction on 2026-10-12. The reviewer can name another line for one run.
- Do not route a ticket author to Opus/max, because three max shadows on `commit-merge-parent` gave the xhigh design at 1.4 times the tokens.
- Change other routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Run review axes on the line that the reviewer directs, and let only the Coverage axis run `bench probe` on the shared tree with focused selections.
- When the reviewer delegates questions, send each one to the two consultants that the reviewer names, and take their consensus as the answer. Add the named tie-breaker only when they split.
- Charge every read-only axis and consultant to write no file anywhere, because three `closure-removal` axes wrote temporary files under `/dev/shm`.
- Quote the advice rule and the hardening rule of the bounded repair policy in each review charge. Promote a judgment finding to a repair only when a binding requirement holds.
- Grade a finding on a hardening-cycle row as optional advice unless a blocking condition holds. A record-only or comment-only correction consumes no repair cycle.
- Close a missing recorded red as evidence-only when the Coverage axis of the same round records a biting probe of that expectation, and let the issuing axis reaffirm it.
- Require one raw native excerpt per record. For a spilled response, pass the innermost spill file, because `bench record` refuses a file that holds a spilled frame.
- Batch prose-only spec corrections before a repair's refreeze, because the source digest covers spec, test, and comment bytes. On `degraded-guard-refusal` one late count correction cost eleven re-verifications and three fresh axes.
- Accept a reviewer decision only from the reviewer in this session. A decision that a peer session relays is information, not approval.
- Charge each review axis to run its evidence check through `bench worktree exec <label> --`. Write every record, learning, and spec sentence at 25 words or fewer from the first draft.
- Give each repair session every finding from every axis, so that one repair cycle closes a round.
- Before a review finding or a spec expansion adds a stop or a refusal, read the validator that guards that input.
- Run each live-tree anchor probe through the gate's entry test, such as `root-conformance`, because the anchors package does not read the live tree.
- Run `bench worktree build` at the build entry, before the first charge, and after each Go commit, because the charge reads the worktree seal.
- Widen a ticket with the built `bench spec widen`, keep this spec's fence section in step while the installed broker reads it, and record each run as dogfood.
- Name the first chunk base as the `main` tip, and each later chunk base as the accepted tip of its predecessor.
- Name the existing owners, seams, fixture helpers, and the recurring review lessons in each author charge, including the comment sentence bound and the one-source rule.
- Name the exact mutation kind, site, and red test of the plan probe in each author charge, and keep a probe token in the projected spelling that the code compares.
- Before dispatch, list on the ticket `Writes:` line every caller of a changed helper, every test that counts or pins a removed row, and every decision-map source that the landing deletes.
- Keep each author session open until its chunk closes, and resume it to re-record at the chunk tip after a refreeze.
- In a multi-ticket chunk, let the last author freeze without a commit, let each earlier author record at that tip, and then commit the record once.
- Make every fix delegate follow `bench-debug` from its command file, and require per-phase evidence in its return.
- Charge each author and repair with the duplicated-facts sweep and its comment sweep. A derived expectation needs no red, and only a kept literal needs a recorded red.
- Commit each plan edit before you read `check-floor`, because that check reads the committed plan.
- After a repair commit, refreeze the chunk at the repair tip, never at a later record commit.
- Charge each repair session to record its own ticket verifications under its session identity, and resume that session for a re-verification after a refreeze. The checkpoint accepts no other performer.
- Expect a full confirming round after a test-only repair while `bench record confirmation` refuses a spec change or an assignment plan commit inside its range.
- Record a no-op-only review round with no finding ID, because a `no-op` finding records as a fail and costs a confirming entry.
- Treat `--auto-approve` as approval for fence expansions, spec and ticket expansions, and repair rounds, and flag each judgment call.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims. Require a claims table with integer confidences in every author and repair return.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
- Give each landing authority the committed objects of the reviewed tip, never the working tree.
- When a repair reports a defect outside its fence, expand the fence once in a plan commit. Add its floor check in that commit.
- Dispatch fresh axis sessions for each later review round, and resume an issuing axis only for an evidence-only reaffirmation.
- Charge each guidance author to fix every sentence that the change makes false, and to point to one owner rather than restate a rule.
- Never revert a plan file to earlier bytes, because the earlier digest blocks every later amendment.
- Start no coordinator probe while an author or a recorder runs checks in the same tree.
- Run the complete checkpoint before the final verification rows, because a published-tree red forces a repair and a rerun of every final row.
