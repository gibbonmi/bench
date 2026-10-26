# Claude model scorecard

Last incorporated landing: `gate-reuse` build (landing `8e1a8828340e58e9c164aee23d385437d7e1eb83`, 2026-10-26).

The `chunk-repair-cycle` build (FT445) ran 12 tickets in 10 chunks on Claude Code as a `--full --delegate --auto-approve` run with an author limit of 3. An Opus/high orchestrator dispatched Opus/high and Opus/xhigh authors and repairs, and a Sonnet/high author ran ticket 12 as a cheap-tier trial. Sonnet/high ran every review axis. Opus/medium and Sonnet/high consulted, and Fable/medium broke their splits. Each chunk used one repair cycle or none.

Provider usage and charges remain unknown. The landing census counted 15 raw calls: python3 4, cat 3, rg 3, cp 2, ls 2, sed 1. The ticket 11 author session was lost after its commit. One test-only repair cycle took 40 minutes, and the reviewer dropped ticket 13 after its measurement.

The `no-go-schedule` build (FT449) ran nine tickets in six chunks on Claude Code as a `--full --delegate --auto-approve` run, and the reviewer raised the author limit to 3 after NG-C3. An Opus/high orchestrator dispatched Opus/high authors and Opus/high and Opus/medium repairs. Sonnet/high ran every review axis. Opus/high and Sonnet/high consulted, and Fable/medium broke their two splits. NG-C3 used both repair cycles, and the other chunks used one or none.

Provider usage and charges remain unknown. The landing census counted 113 raw calls: rg 74, python3 25, cat 5, sed 5, cd 3, wc 1. Three concurrent authors collided on one symbol, the serial review record refused the concurrent chunks, and the landing gate found a decision-map source that the landing deletes.

The `no-go-schedule` shaping and spec stages ran on Claude Code for FT449. By reviewer direction, one Opus coordinator did the research, the spec, and the slicing of nine tickets inline, with no reader pair and no fork. Its own traced test pass and permission-denial probe replaced the research readers. A Fable/medium map review found two soundness holes and a ticket contradiction before readiness.

Opus/medium accepted the spec at iteration 2, after iteration 1 found misplaced chunk rows, a dropped refusal, and unquoted `-exec` operands. Opus/medium accepted the tickets at iteration 2, after iteration 1 found two tests that pin the old argv and a record query with no producer. Provider usage and charges remain unknown.

The `gate-reuse` build (FT314) ran six tickets in four chunks on Claude Code as a `--full --delegate --auto-approve` run. An Opus/medium orchestrator dispatched Opus/medium authors, continuations, and repairs, and Sonnet/high ran every review axis. Opus/high and Sonnet/high consulted once, and Fable/high broke their split; after GR-C1 the reviewer limited consultations to blocking findings. Each chunk used one repair cycle.

Provider usage and charges remain unknown. The landing census counted 16 raw calls: cat 7, rg 7, wc 2. Four author contexts ran past the 150k-token target, and two planned probes could not witness their rows until the orchestrator re-cut them.

The `gate-reuse` stages ran on Claude Code for FT314. A first spec stage reached sign-off, and the reviewer reopened it because it missed the outcome's commitment criteria.
One Opus coordinator then shaped a 14-ticket decision map and sliced six tickets on the map-only route in the same session, by reviewer direction.
An Opus/high delegate researched the executable roles, and the coordinator measured the saving from git, the gate run logs, and the session records.

Fable/medium accepted the tickets at iteration 2. Iteration 1 found a row whose producer lands in a later ticket, a seam outside its test package, and absent helper inputs.
The reopened draft used Sonnet/high, at two iterations each for its spec and its tickets. Provider usage and charges remain unknown.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 42 landings + implementer, 10 charges + reviewer, 48 axes | On `repair-confirmation` Fable/high orchestrated four chunks under `--auto-approve` and added four spec guarantees, each flagged for veto. It retried a stale-digest refusal three times and a `check-floor` red three times before each plan commit closed it. | Coordination of a delegated build and adversarial spec review; review axes in parallel by reviewer direction; it implements only when the reviewer names it | orchestrator 0.36 over 3 pairs, 1 abstained; medium reviewer 0.256 over 20 pairs, 1 abstained |
| Fable / high, xhigh, medium | reviewer, 124 axes + 4 blind comparisons + spec reviewer, 10 passes + ticket reviewer, 6 passes + verifier, 3 light-path checks + control, 1 full-retrieval run + decider, 28 consultations + debug, 2 charges | On `gate-reuse` Fable/high broke a three-question split: it took the fail-closed executable rule and the relative link spelling, and it sent the run-binary source question to a learning. On `no-go-schedule` Fable/medium broke two consultant splits. | Review axes, spec review, and ticket review by reviewer direction; consultations whenever the coordinator is uncertain; the RE2-style full control | reviewer: 0.168 over 32 pairs, 1 abstained; decider 0.04 over 1 pair |
| Opus / high, xhigh, max | implementer, latest 10 of 104 current attempts, 3 of them max shadows, 1 a max light-path fix, and 2 xhigh light-path fixes; reviewer, 319 earlier axes + 42 `closure-removal` axes + 4 consultations; researcher, 20 paired charges + 10 off-route `test-cost-signal` charges; repair, latest 10 of 92 current attempts; orchestrator, 6 builds and 5 spec stages; spec author, slicer, and reviewer, latest 10 passes, 7 of them xhigh slicers | On `no-go-schedule` the Opus/high authors landed nine tickets, three of them at once under `--delegate 3`. On `gate-reuse` Opus/high consulted once and split with Sonnet/high on three of four questions. | High for first-round review axes, orchestration, consultations, and debug runs; xhigh for a code ticket where a repair round costs more than about 100k tokens, or that crosses owner seams; repairs at the ticket's effort. | implementation: 0.099 over 128 pairs, 3 abstained, and 0.007 over 18 xhigh pairs; orchestrator 0.26 over 3 pairs; reviewer: 0.162 over 74 pairs, 9 abstained; repair: 0.035 over 44 pairs; consultant 0.05 over 2 pairs |
| Opus / medium, low | implementer, latest 10 of 27 tickets + 5 benchmark shadows; repair, latest 10 of 41 sessions; orchestrator, 10 builds; reviewer, latest 10 of 100 axes + spec and ticket reviewer, 10 passes | On `gate-reuse` Opus/medium authored all six tickets and every repair, and no author claim failed its row. Four tickets stopped correctly on a test outside their fence, and four author contexts ran past the 150k target. | Low for an exact-spec guidance ticket; medium as the default for a code ticket, for review axes by reviewer direction, and for a mechanical test migration; a repair keeps the ticket's declared effort. | 0.082 over 135 pairs; 34 abstained; reviewer 0.176 over 76 pairs; orchestrator 0.299 over 9 pairs |
| Sonnet / high | orchestrator, 3 landings + reviewer, 183 axes + spec reviewer, 7 passes + ticket reviewer, 2 passes + researcher, 3 charges + implementer, 9 comparison tickets and 1 trial ticket + staleness, 8 slices + consultant, 2 questions and the FT445 consultations | On `gate-reuse` Sonnet/high ran every axis, three staleness slices, and one consultation. Its Coverage axes found that two planned probes could not witness their rows and a broker re-exec loop; two axes ran their evidence check from the primary checkout. | First-round and confirming review axes, consultations, staleness slices, and fix reviews by reviewer direction, and orchestration; a measurement or prose ticket when the reviewer names it, and not yet a code ticket | high 0.145 over 2 pairs; xhigh 0.40 over 2 pairs; reviewer 0.131 over 19 pairs; consultant 0.225 over 2 pairs |
| Sonnet / high (research and polish) | researcher, 15 paired charges + 5 single charges + 19 slicing-research charges + polish, 1 session | On the `landing-gate-long-poles` spec stage Sonnet/medium readers found the over-budget test files and a hand-mirrored census copy. One Sonnet reader undercounted the owner-marker starts, and its Haiku pair got the count right. | Default research reader at medium effort, paired with Haiku/high by reviewer direction on 2026-10-11; polish of a cheap draft when the reviewer names it | research 6 of 6 earlier slices clean; slicing research missed the closure classes in 6 of 6 uncharged slices and in 0 of 3 charged slices |
| Haiku / default–xhigh | researcher, 21 charges + 12 slicing-research charges + implementer, 2 comparison tickets + auditor, 3 charges | On the `landing-gate-long-poles` spec stage Haiku/high readers found the gate's run-binary binding, the `poolkey` common-directory start, and the TMPDIR children. Each return cost about 150k to 350k tokens, against about 55k to 100k for a Sonnet/medium pair. | Research reader at high effort by default, paired with Sonnet/medium; staleness audits; not ticket authorship | research 10 of 12 slices clean; slicing research agreed with Sonnet in 3 of 3 chunks; authorship 0 landings |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges; reviewer, 42 medium axes | On `test-cost-signal` Sonnet/medium ran every axis by reviewer direction. Its Coverage probes found the max-wins and min-wins survivors, but no axis found the spill regression, and two TC-C1 axes skipped their required reads. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; medium review axes only by reviewer direction, with a coordinator probe on each chunk | implementer unknown; reviewer 0.199 over 12 pairs |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `worktree-seam-reduction` comparison | Opus and Sonnet / high / implementer; Fable / high / reviewer | One blind run ranked Sonnet-sliced first, Opus-sliced second, and Opus-full third. The reviewer selected the second arm, and all three arm gates passed. |
| `markdown-block-reader` spec review | Sonnet / high, xhigh; Opus / xhigh; Fable / high, xhigh / spec reviewer | Five reviewers graded one commit with one charge. List costs were Sonnet high $2.14, Fable high $1.98, Opus xhigh $2.44, Fable xhigh $2.50, and Sonnet xhigh $6.31. No reviewer found every confirmed finding, and only Fable high reported a false one. |
| `build-directions` effort benchmark | Opus / low, medium, high, xhigh / implementer; Fable / high / blind grader | Ten first-pass commits on three tickets met every acceptance row. On two guidance tickets low matched high at half the tokens. On the code ticket xhigh ranked first, medium second, and high last of four. |
| `commit-merge-parent` effort assessment | Opus / xhigh, max / implementer | Three max shadows matched the xhigh design on every code ticket at 1.41 to 1.44 times the tokens and up to 2.3 times the time. Eight review rounds found 17 findings on the xhigh trees, 10 of them repaired. |
| `review-evidence-file-pages` RE2 control | claude-fable-5-1 / high / reviewer | The narrow round found 8 findings in 399474 tokens and the full control found 6 in 227918 tokens. They agreed on 3, so neither shape covered the other. |

## Current decisions

- Read this scorecard before each line declaration and again after each context compaction, because its decisions carry the orchestration rules that a build repeats. `landing-gate-long-poles` broke two of them after a compaction.
- Route each read-only research and slicing-research delegate to Haiku/high, by reviewer direction on 2026-10-21. Never copy the reviewer line to a research delegate; `test-cost-signal` ran ten Opus/high research charges that way.
- Before a spec or ticket commit adds or changes a quantified predicate, sweep its set against the tree and cite the count. `landing-gate-long-poles` took six review iterations, and four of its blockers were unswept sets.
- Size each ticket to about 150k tokens of author context, by reviewer direction on 2026-10-21. Count the read set outside `Writes:`, about 5k tokens for each witness probe, and one lane retry; `gate-reuse` ran four authors past the target.
- When a spec `Line:` names a model that this harness cannot run, ask the reviewer for the author line before the first dispatch, and record the direction in the assignment commit.
- Declare the expected repair rounds for each chunk with a confidence, so the orchestrator claim gets a label. Expect one repair cycle per chunk for a small seam and two to three when a chunk deletes a shared seam.
- Route every coding delegate to the line that the reviewer directs for the run. Move a ticket that the spec rates harder, with an uncertain seam, one effort step up, and record the reason in its assignment commit.
- Do not route a ticket author to Opus/max, because three max shadows on `commit-merge-parent` gave the xhigh design at 1.4 times the tokens.
- Change other routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Run review axes on the line that the reviewer directs, and let only the Coverage axis run `bench probe` on the shared tree with focused selections.
- During a review round, read source through `git show` at the frozen tip, because a Coverage probe can hold a mutant in the shared working file.
- Send a finding to the consultant pair only when it is blocking and would otherwise stop for the reviewer, by reviewer direction on 2026-10-26; decide every other finding in the orchestrator. Add the tie-breaker only when the pair splits.
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
- Widen a ticket with `bench spec widen`, and record each plan expansion with `bench learning` before dispatch. When a plan commit drops a ticket, sweep the spec for its rows, its name, and its seam block.
- Name the first chunk base as the `main` tip, and each later chunk base as the recorded chunk tip of its predecessor, because the review chain refuses any source commit between them. Pass that base to the ticket charge, the post-ticket preflight, and `bench record ticket`; `test-cost-signal` used an assignment commit and needed a base move.
- Commit a comment-only correction with no refreeze, because a refreeze orphans every recorded verification and the checkpoint proves the comment gap.
- Name the existing owners, seams, fixture helpers, census tests, and the recurring review lessons in each author charge, including the comment sentence bound and the one-source rule.
- Name the exact mutation kind, site, and red test of the plan probe in each author charge.
- Before dispatch, list on the ticket `Writes:` line every caller of a changed helper and every test that counts or pins a removed row.
- At build entry, rewrite every decision-map source that the landing deletes, such as the outcome's own roadmap row. The FT449 landing ran red on that source.
- Record each ticket at its own commit, and restage no ticket, because FT445 binds each verification entry to the commit that it graded.
- Charge each author to run each plan probe with its exact planned selection immediately before its `bench record excerpt` call, because the excerpt reads only the last call.
- Run the `bench worktree build` and `doctor` rehearsal before the landing of a spec that changes Go build inputs.
- Let each ticket author rewrite its own planned coverage cells in the ticket commit, so that each citation lands with its cited test.
- In a prose repair charge that merges sentences under a word bound, name each row phrase that must survive.
- Make every fix delegate follow `bench-debug` from its command file, and require per-phase evidence in its return.
- Charge each author and repair with the duplicated-facts sweep and its comment sweep.
- Commit each plan edit before you read `check-floor`, because that check reads the committed plan. A widening to a new package raises the floor, so add its floor check in the same plan commit.
- Send a fresh narrow session for a ticket continuation, not a resume of a full author context, by reviewer direction on 2026-10-26.
- Before a ticket relies on a planned probe, trace the probe to the assertion it claims and compile its mutant. Two `gate-reuse` probes failed this.
- Start a prepared checkpoint gate only when no plan commit or author edit can follow before it ends; a `gate-reuse` plan commit made one refuse.
- After a repair commit, refreeze the chunk at the repair tip, never at a later record commit.
- Charge each repair session to record its own ticket verifications under its session identity. After a lost author session, record a `session-lost` replacement that verifies the preserved commit.
- Close a test-only repair with one `bench record confirmation` that lists every open finding of the round, but only when no sibling chunk's commit sits between the reviewed tip and the repair.
- Record a no-op-only review round with no finding ID, because a `no-op` finding records as a fail and costs a confirming entry.
- Treat `--auto-approve` as approval for fence expansions, spec and ticket expansions, and repair rounds, and flag each judgment call.
- Under `--delegate`, freeze concurrent chunks in plan order and chain each base to the previous chunk tip, because the review record walks its chunks serially.
- When two concurrent tickets write one Go package, name in each charge the owner of each new shared helper. FT449 tickets 07 and 08 collided on one symbol.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims. Require a claims table with integer confidences in every author and repair return; six `test-cost-signal` repairs returned `high` and abstained.
- Charge a repair that pins a selection rule with the whole mutant family, such as first-wins, max-wins, and min-wins. Three `test-cost-signal` confirming rounds each found the next survivor of one test.
- Run one coordinator probe on each returned tree, at a site that no author probed, and never while an author or a recorder runs checks there. On `chunk-repair-cycle` it found a silent reuse guard before any axis.
- Charge each author and repair session with the plan's probes only, by reviewer direction on minimal probing. The coordinator probe carries the finding mutants.
- Before a spec adds lines to a bounded response, count each reader's shortest response against `bounds.ResponseLines`. On `test-cost-signal` the checkpoint gate, not review, found the spill.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
- Read the assignment census record before `bench worktree land`, because the landing removes it.
- Run the complete checkpoint before the final verification rows, because a published-tree red forces a repair and a rerun of every final row.
- Land with the request token of the outcome's `bench commitment start`, and copy it into the handoff state, because a new token or a reauthorize to it loses the delivery binding.
- When a Go file nears its line cap, widen the ticket with a new sibling file before the author edit, not after a stop.
- Expect the Opus/medium axes to find real behavior defects on a multi-actor lifecycle spec, so keep three axes on each confirming round.
- In a confirming round, send a new finding below severity 5 to Advice and open no repair cycle, by reviewer direction on 2026-10-23.
- Put the system suite in the checks of a ticket that changes a CLI verb's exit or output. On `process-lifetime` only the checkpoint gate caught the BO69 regression.
- At build entry, flag every native-platform row that the build host cannot run, and ask the reviewer before the last chunk. PL90 stopped `process-lifetime` at its freeze.
- After `bench commitment block` and `unblock`, rerun `bench commitment start` with the same token and deliverable, because the unblock drops the delivery binding.
