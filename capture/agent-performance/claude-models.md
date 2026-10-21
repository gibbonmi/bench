# Claude model scorecard

Last incorporated landing: `candidate-binary-install` build (landing `3b970d1c2d1686f05436dea57e5c845037723cb0`, 2026-10-20).

The `candidate-binary-install` build (FT419) ran six tickets in five chunks on Claude Code as a `--full --delegate --auto-approve` run with an author limit of 3. An Opus/high orchestrator dispatched fresh Opus/high authors and ten Opus/high repair sessions. Opus/medium ran every review axis by reviewer direction. One Opus/high and Sonnet/high consultant pair decided two behavioral questions by consensus, and a Haiku pass ran the staleness audit. The chunks used 2, 1, 3, 2, and 1 repair cycles.

Provider usage and charges remain unknown. The landing census counted 80 raw calls: rg 60, python3 6, cat 4, cd 3, wc 3, ls 2, awk 1, sed 1. Review found five real defects that the authors' tests passed through.

The `merge-target-lane` build (FT335) ran two tickets in two chunks on Claude Code as a `--full --auto-approve` run with an author limit of 1. An Opus/medium orchestrator ran two Sonnet staleness slices, an Opus/medium author for ticket 01, and an Opus/high author for the harder ticket 02. Opus/medium ran every review axis by reviewer direction, and no consultation question arose. Each chunk used 1 repair cycle.

Provider usage and charges remain unknown. The landing census counted 41 raw calls, and the landing removed the record before the close read its verbs.

The `plan-commit` build (FT380) ran three tickets in two chunks on Claude Code as a `--full --auto-approve` run with an author limit of 1. An Opus/medium orchestrator ran three Sonnet staleness slices, fresh Opus/high authors, and three Opus/high repair sessions. Opus/medium ran every review axis by reviewer direction, and no consultation question arose. PC-C1 used 1 repair cycle and PC-C2 used 2.

Provider usage and charges remain unknown. The landing census counted 8 raw calls: rg 7, git diff 1.

The `closure-removal` build (FT293) ran six tickets in four chunks on Claude Code as a `--full --auto-approve` run with an author limit of 1. An Opus orchestrator dispatched fresh Opus/high authors and eleven Opus/high repair sessions, and Opus/high ran every review axis by reviewer direction. Four Sonnet staleness slices opened the build, and one Fable/medium and Sonnet/high consultation pair agreed on the only judgment question. The chunks used 3, 2, 2, and 3 repair cycles, and `--auto-approve` covered the third cycles.

Provider usage and charges remain unknown. The landing census counted 24 raw calls: rg 16, cp 4, cat 2, cd 1, sed 1.

The `closure-removal` spec stage ran on Claude Code for FT293. An Opus coordinator ran eight paired Haiku/high research readers and eight paired Haiku/high slicing readers by reviewer direction, and it sliced the six tickets itself. An Opus fork wrote the spec, and the coordinator dropped the fork's staged-spec sweep before review. Fable/high ran the spec review in 1 iteration and the ticket review in 2, by reviewer direction. The ticket review caught a deleted file whose in-package callers the slicer had not searched.

Provider usage and charges remain unknown. The harness reported about 150k to 270k tokens for each research and slicing reader.

The `plan-commit` spec stage ran on Claude Code for FT380. An Opus coordinator ran four Sonnet/medium and Haiku/high research pairs and two slicing pairs, and it sliced the three tickets itself. An Opus fork wrote the spec. Opus/high and Sonnet/high reviewed both rounds by reviewer direction. In the spec round, Opus/high rejected iteration 1 with three blocking gaps, and Sonnet/high accepted it. Both accepted iteration 2 and the ticket round in 1.

The research pairs disagreed twice, and a targeted read closed each disagreement. A Sonnet reader denied that the Markdown lane runs the coverage parser, and the code showed that it does. A Haiku reader found the `./cmd/bench` floor element that its Sonnet pair missed. Provider usage and charges remain unknown. The harness reported about 110k to 290k tokens for each reader.

The `merge-target-lane` spec stage ran on Claude Code for FT335. An Opus coordinator ran four Opus/medium research pairs and one Opus/medium slicing pair. A fresh Opus/high writer sliced the two tickets. An Opus fork wrote the spec. Opus/medium reviewed both rounds by reviewer direction, and each round took 2 iterations. The spec round rejected a second production owner of the wrapper's kit rule, and the reviewer recorded two decisions. The ticket round rejected one chunk that held a seam ticket and its consumer.

The research pairs disagreed once, about the root of the whole-gate run log, and a targeted read closed it. The slicing pair found a structure budget and a serial-ceiling pin that the spec had missed. Provider usage and charges remain unknown. The harness reported about 40k to 130k tokens for each reader, 250k for the fork, and 170k to 185k for the spec reviewer.

The `candidate-binary-install` spec stage ran on Claude Code for FT419. An Opus coordinator ran four Opus/medium research pairs and five Haiku/high slicing pairs, and it sliced seven tickets itself, six after the review merged one. An Opus/high fork wrote the spec. Opus/medium reviewed the spec round, which took 3 iterations. Two consultations between iterations 2 and 3 (Fable/medium and Sonnet/high) both found the same lifecycle deadlock, an unrelated landing whose refresh uninstalled the candidate. The ticket round took 5 iterations. Opus/medium graded the first two, and Fable/high and Opus/high graded the last three under the reviewer's delegated sign-off. In those three, the two graders found different blockers twice and agreed to approve at the fifth.

Two of the four rejected re-reviews found a blocker that the coordinator's own previous repair had introduced. A consensus diagnosis by Fable/high and Opus/high attributed most misses to three causes: a fork that wrote claims without re-reading the code, no lifecycle walk across the actors, and re-reviews that graded only the repair. The Haiku slicing readers found the directory caps, the census pin, and the rule-sentence anchor conflict that the accepted spec had missed. Provider usage and charges remain unknown. The harness reported about 40k to 280k tokens for each reader, 260k for the fork, and 75k to 300k for each grader.

The FT366-guard-residual spec stage refreshed two staged specs on Claude Code. An Opus coordinator ran five Haiku/high research pairs and four Haiku/high slicing pairs by reviewer direction, and it sliced the tickets itself. An Opus fork refreshed both specs, and Opus/high graded every round by reviewer direction. Each spec round took 2 iterations, and the first ticket round took 2.

The reviewer then directed an experiment: each ticket fits about 200k tokens of author context for an Opus/medium or Opus/high author. The re-slice gave 9 tickets, and its fresh ticket round accepted in 1 iteration. The Haiku readers found a new consumer of the retired helpers and the headroom and canary facts that the old specs missed. Provider usage and charges remain unknown. The harness reported about 40k to 230k tokens for each reader, 300k for the fork, and 115k to 175k for each grader.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 42 landings + implementer, 10 charges + reviewer, 48 axes | On `repair-confirmation` Fable/high orchestrated four chunks under `--auto-approve` and added four spec guarantees, each flagged for veto. It retried a stale-digest refusal three times and a `check-floor` red three times before each plan commit closed it. | Coordination of a delegated build and adversarial spec review; review axes in parallel by reviewer direction; it implements only when the reviewer names it | orchestrator 0.36 over 3 pairs, 1 abstained; medium reviewer 0.256 over 20 pairs, 1 abstained |
| Fable / high, xhigh, medium | reviewer, 124 axes + 4 blind comparisons + spec reviewer, 10 passes + ticket reviewer, 4 passes + verifier, 3 light-path checks + control, 1 full-retrieval run + decider, 26 consultations + debug, 2 charges | On `closure-removal` one Fable/medium consultation agreed with Sonnet/high to accept the refusal of a root headroom premise, so no tie-breaker ran. It cited the structure engine's root spelling as the reason. | Review axes, spec review, and ticket review by reviewer direction; consultations whenever the coordinator is uncertain; the RE2-style full control | reviewer: 0.168 over 32 pairs, 1 abstained |
| Opus / high, xhigh, max | implementer, latest 10 of 104 current attempts, 3 of them max shadows, 1 a max light-path fix, and 2 xhigh light-path fixes; reviewer, 319 earlier axes + 42 `closure-removal` axes + 3 consultations; researcher, 20 paired charges; repair, latest 10 of 92 current attempts; orchestrator, 6 builds and 4 spec stages; spec author, slicer, and reviewer, latest 10 passes, 7 of them xhigh slicers | On `candidate-binary-install` six Opus/high authors landed every row, but review found five real defects that their tests passed through. Ten Opus/high repairs closed each round, and the Opus/high orchestrator caught two line-budget stops and one wrong chunk base. | High for first-round review axes, orchestration, consultations, and debug runs; xhigh for a code ticket where a repair round costs more than about 100k tokens, or that crosses owner seams; repairs at the ticket's effort. | implementation: 0.110 over 114 pairs, 3 abstained, and 0.01 over 8 xhigh pairs; orchestrator 0.26 over 3 pairs; reviewer: 0.162 over 74 pairs, 9 abstained; repair: 0.035 over 44 pairs; consultant 0.01 over 1 pair |
| Opus / medium, low | implementer, latest 10 of 24 tickets + 5 benchmark shadows; repair, latest 10 of 36 sessions; orchestrator, 9 builds; reviewer, latest 10 of 100 axes | On `candidate-binary-install` the Opus/medium axes ran 42 results over five chunks and found five real behavior defects, including two open declared builds and an absent-ledger restore. The CB-C5 confirming axes each found a different advice item. | Low for an exact-spec guidance ticket; medium as the default for a code ticket, for review axes by reviewer direction, and for a mechanical test migration; a repair keeps the ticket's declared effort. | 0.083 over 123 pairs; 28 abstained; reviewer 0.156 over 45 pairs; orchestrator 0.305 over 4 pairs |
| Sonnet / high | orchestrator, 3 landings + reviewer, 132 axes + spec reviewer, 5 passes + researcher, 3 charges + implementer, 9 comparison tickets + staleness, 6 slices | On `candidate-binary-install` Sonnet/high and Opus/high consultants agreed on both behavioral questions: two parser refusals and a new overlap refusal face. No tie-breaker ran. | First-round and confirming review axes, consultations, staleness slices, and fix reviews by reviewer direction, and orchestration; not fresh ticket authorship | high 0.145 over 2 pairs; xhigh 0.40 over 2 pairs; reviewer 0.203 over 4 pairs; consultant 0.09 over 1 pair |
| Sonnet / high (research and polish) | researcher, 10 paired charges + 5 single charges + 15 slicing-research charges + polish, 1 session | On `repair-confirmation` Sonnet delegates ran the build research by reviewer direction. The build record keeps no result for each slice, so this landing adds no quality evidence. | Default research reader at medium effort, paired with Haiku/high by reviewer direction on 2026-10-11; polish of a cheap draft when the reviewer names it | research 6 of 6 earlier slices clean; slicing research missed the closure classes in 6 of 6 uncharged slices and in 0 of 3 charged slices |
| Haiku / default–xhigh | researcher, 16 charges + 8 slicing-research charges + implementer, 2 comparison tickets + auditor, 3 charges | On `candidate-binary-install` a Haiku staleness pass ran by reviewer direction, and its corrections went into the first plan amendment. A Haiku reader also found that a roadmap row already held a process caveat. | Research reader at high effort by default, paired with Sonnet/medium; staleness audits; not ticket authorship | research 10 of 12 slices clean; slicing research agreed with Sonnet in 3 of 3 chunks; authorship 0 landings |
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

- Read this scorecard before each line declaration, because its decisions carry the orchestration rules that a build repeats. `candidate-binary-install` again prepared a charge before its assignment commit.
- When a spec `Line:` names a model that this harness cannot run, ask the reviewer for the author line before the first dispatch, and record the direction in the assignment commit.
- Declare the expected repair rounds for each chunk with a confidence, so the orchestrator claim gets a label. Expect one repair cycle per chunk for a small seam and two to three when a chunk deletes a shared seam.
- Route every coding delegate to the line that the reviewer directs for the run. Move a ticket that the spec rates harder, with an uncertain seam, one effort step up, and record the reason in its assignment commit.
- Do not route a ticket author to Opus/max, because three max shadows on `commit-merge-parent` gave the xhigh design at 1.4 times the tokens.
- Change other routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Run review axes on the line that the reviewer directs, and let only the Coverage axis run `bench probe` on the shared tree with focused selections.
- During a review round, read source through `git show` at the frozen tip, because a Coverage probe can hold a mutant in the shared working file.
- When the reviewer delegates questions, send each one to the two consultants that the reviewer names, and take their consensus as the answer. Add the named tie-breaker only when they split.
- Charge every read-only axis and consultant to write no file anywhere.
- Quote the advice rule and the hardening rule of the bounded repair policy in each review charge. Promote a judgment finding to a repair only when a binding requirement holds.
- Grade a finding on a hardening-cycle row as optional advice unless a blocking condition holds. A record-only or comment-only correction consumes no repair cycle.
- Close a missing recorded red as evidence-only when the Coverage axis of the same round records a biting probe of that expectation, and let the issuing axis reaffirm it.
- Close a test copy of a production fact under the test-expectation exception only with a recorded probe at the production site that reds the copy alone.
- Require one raw native excerpt per record. For a spilled response, pass the innermost spill file, because `bench record` refuses a file that holds a spilled frame.
- Batch prose-only spec corrections before a repair's refreeze, because the source digest covers spec, test, and comment bytes.
- Accept a reviewer decision only from the reviewer in this session. A decision that a peer session relays is information, not approval.
- Charge each review axis to run its evidence check through `bench worktree exec <label> --`. Write every record, learning, and spec sentence at 25 words or fewer from the first draft.
- Give each repair session every finding from every axis, so that one repair cycle closes a round.
- Commit each assignment record before you prepare its author charge, because the charge evidence binds the source tip.
- Pass `--base` to the first `bench record ticket` of the first chunk, because that chunk has no entry and no predecessor.
- Run each live-tree anchor probe through the gate's entry test, such as `root-conformance`, because the anchors package does not read the live tree.
- Widen a ticket with `bench spec widen`, and record each plan expansion with `bench learning` before dispatch.
- Name the first chunk base as the `main` tip, and each later chunk base as the accepted tip of its predecessor.
- Name the existing owners, seams, fixture helpers, census tests, and the recurring review lessons in each author charge, including the comment sentence bound and the one-source rule.
- Name the exact mutation kind, site, and red test of the plan probe in each author charge.
- Before dispatch, list on the ticket `Writes:` line every caller of a changed helper, every test that counts or pins a removed row, and every decision-map source that the landing deletes.
- In a multi-ticket chunk, let every author stage its runs at the chunk tip before the first record commit.
- In a prose repair charge that merges sentences under a word bound, name each row phrase that must survive.
- Make every fix delegate follow `bench-debug` from its command file, and require per-phase evidence in its return.
- Charge each author and repair with the duplicated-facts sweep and its comment sweep.
- Commit each plan edit before you read `check-floor`, because that check reads the committed plan.
- After a repair commit, refreeze the chunk at the repair tip, never at a later record commit.
- Charge each repair session to record its own ticket verifications under its session identity.
- Expect a full confirming round after a test-only repair while `bench record confirmation` refuses an assignment plan commit inside its range.
- Record a no-op-only review round with no finding ID, because a `no-op` finding records as a fail and costs a confirming entry.
- Treat `--auto-approve` as approval for fence expansions, spec and ticket expansions, and repair rounds, and flag each judgment call.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims. Require a claims table with integer confidences in every author and repair return.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
- Read the assignment census record before `bench worktree land`, because the landing removes it.
- Start no coordinator probe while an author or a recorder runs checks in the same tree.
- Run the complete checkpoint before the final verification rows, because a published-tree red forces a repair and a rerun of every final row.
- Land with the request token of the outcome's `bench commitment start`, because a new token loses the delivery binding.
- When a Go file nears its line cap, widen the ticket with a new sibling file before the author edit, not after a stop.
- Expect the Opus/medium axes to find real behavior defects on a multi-actor lifecycle spec, so keep three axes on each confirming round.
