# Claude model scorecard

Last incorporated landing: `anchor-admission-prune` build (landing `3854648849e2193c95c9933b9639c6f1d858a75f`, 2026-10-18).

The `anchor-admission-prune` build ran ten tickets in seven chunks on Claude Code as a `--full --auto-approve` run with an author limit of 1. An Opus orchestrator dispatched fresh Opus/medium authors and repairs, with Opus/high for ticket 8 and its repairs. Opus/medium ran every review axis, and one Fable/high and one Opus/high consultation judged the per-ticket verification rule. The chunks used 0, 1, 0, 1, 1, 2, and 2 repair cycles, and two in-build preflight fixes added a ticket and a chunk.

Provider usage and charges remain unknown. The landing census counted 34 raw calls: rg 20, cat 5, sed 4.

The `preflight-ownership-closure` build ran eight tickets in four chunks on Claude Code as a `--full --auto-approve` run across several sessions. Its author limit was 1. Opus authors wrote tickets 01 to 05, and `gemini-4-argon-medium` wrote tickets 06 and 07 through `agy`. A fresh Opus/medium author wrote ticket 08 by reviewer direction. Opus/medium ran the later axes and repairs, and Opus/high answered the consultations. The chunks used 1, 6, 3, and 3 repair cycles, and OC-C2 and OC-C4 went past the allowance under `--auto-approve`.

Provider usage and charges remain unknown. The OC-C4 census counted 75 raw calls: rg 43, python3 8, and sed 8. The harness reported about 370k tokens for the ticket 08 author and 105k to 195k for each repair.

The `anchor-admission-prune` spec stage ran on Claude Code from a reviewer-confirmed conversation. An Opus coordinator ran four Opus/high research readers and twelve paired Sonnet/high slicing readers, and it sliced the eight tickets itself. An Opus fork wrote the spec. Opus/medium ran the spec review in 2 iterations and the ticket review in 2. The reviewer then directed that Sonnet research delegates run at medium effort.

Provider usage and charges remain unknown. The harness reported about 140k to 250k tokens for each slicing reader.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 42 landings + implementer, 10 charges + reviewer, 48 axes | On `repair-confirmation` Fable/high orchestrated four chunks under `--auto-approve` and added four spec guarantees, each flagged for veto. It retried a stale-digest refusal three times and a `check-floor` red three times before each plan commit closed it. | Coordination of a delegated build and adversarial spec review; review axes in parallel by reviewer direction; it implements only when the reviewer names it | orchestrator 0.36 over 3 pairs, 1 abstained; medium reviewer 0.256 over 20 pairs, 1 abstained |
| Fable / high, xhigh | reviewer, 124 axes + 4 blind comparisons + spec reviewer, 10 passes + ticket reviewer, 4 passes + verifier, 3 light-path checks + control, 1 full-retrieval run + decider, 21 consultations + debug, 2 charges | On `anchor-admission-prune` one Fable/high consultation judged that the per-ticket verification split adds value and that author-performed entries are ceremony. It agreed with the Opus/high consultation, and it cited each claim. | Review axes, spec review, and ticket review by reviewer direction; consultations whenever the coordinator is uncertain; the RE2-style full control | reviewer: 0.168 over 32 pairs, 1 abstained |
| Opus / high, xhigh, max | implementer, latest 10 of 97 current attempts, 3 of them max shadows, 1 a max light-path fix, and 2 xhigh light-path fixes; reviewer, 319 earlier axes + the `record-ticket` axes + 3 consultations; researcher, 20 paired charges; repair, latest 10 of 79 current attempts; orchestrator, 5 builds and 4 spec stages; spec author, slicer, and reviewer, latest 10 passes, 7 of them xhigh slicers | On `anchor-admission-prune` the Opus/high ticket 8 author met all seven guidance rows with no needle change, and two Opus/high repairs closed a self-contradiction and its leftover. The Opus/high consultant agreed with Fable/high that per-ticket verification by author is ceremony. | High for first-round review axes, orchestration, consultations, and debug runs; xhigh for a code ticket where a repair round costs more than about 100k tokens, or that crosses owner seams; repairs at the ticket's effort. | implementation: 0.097 over 82 pairs, 3 abstained, and 0.01 over 8 xhigh pairs; orchestrator 0.26 over 3 pairs; reviewer: 0.162 over 74 pairs, 9 abstained; repair: 0.029 over 29 pairs; consultant 0.01 over 1 pair |
| Opus / medium, low | implementer, latest 10 of 23 tickets + 5 benchmark shadows; repair, latest 10 of 35 sessions; orchestrator, 8 builds; reviewer, latest 10 of 88 axes | On `anchor-admission-prune` nine Opus/medium authors held their fences, and three stopped correctly at a preflight or structure refusal. Seven Opus/medium repairs closed every accepted finding, but one repair unpinned a test bound that a coordinator probe then caught. | Low for an exact-spec guidance ticket; medium as the default for a code ticket, for review axes by reviewer direction, and for a mechanical test migration; a repair keeps the ticket's declared effort. | 0.087 over 116 pairs; 28 abstained; reviewer 0.156 over 45 pairs; orchestrator 0.287 over 3 pairs |
| Sonnet / high | orchestrator, 3 landings + reviewer, 132 axes + spec reviewer, 5 passes + researcher, 3 charges + implementer, 9 comparison tickets | On `probe-subjects-and-mutations` four Sonnet staleness slices at default effort found the injected replacer seam, a new test file, and four anchor registries missing from the fences. Each slice named what it did not verify. | First-round and confirming review axes, consultations, and fix reviews by reviewer direction, and orchestration; not fresh ticket authorship | high 0.145 over 2 pairs; xhigh 0.40 over 2 pairs; reviewer 0.203 over 4 pairs; consultant 0.09 over 1 pair |
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

- Read this scorecard before each line declaration, because its decisions carry the orchestration rules that a build repeats.
- Declare the expected repair rounds for each chunk with a confidence, so the orchestrator claim gets a label. Expect zero to two repair cycles per chunk; the seven chunks of `anchor-admission-prune` used 0 to 2.
- Route every coding delegate to Opus/high and every small-fix delegate to Opus/low, by reviewer direction on 2026-10-12. The reviewer can name another line for one run. On `preflight-ownership-closure` the reviewer named Opus/medium for the later authors and repairs.
- Do not route a ticket author to Opus/max, because three max shadows on `commit-merge-parent` gave the xhigh design at 1.4 times the tokens.
- Change other routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Run review axes on the line that the reviewer directs, and let only the Coverage axis run `bench probe` on the shared tree with focused selections.
- When the reviewer is away, route each build judgment call to an Opus/high read-only consultant. On `preflight-ownership-closure` consultants decided two fence questions and traced a cross-chunk charge refusal.
- Quote the advice rule and the hardening rule of the bounded repair policy in each review charge. Promote a judgment finding to a repair only when a binding requirement holds.
- Grade a finding on a hardening-cycle row as optional advice unless a blocking condition holds. A record-only or comment-only correction consumes no repair cycle.
- Grade a Coverage finding only against the lines that the chunk delta adds. Send a gap in an accepted chunk's code to the drain, unless it invalidates that chunk's evidence.
- Require one raw native excerpt per record with its command line first, never a hand summary. Copy a spilled response natively into its excerpt.
- Batch prose-only spec corrections before a repair's refreeze, because the source digest covers spec, test, and comment bytes.
- Ask the reviewer about each host-bound acceptance row before the first chunk freeze.
- Accept a reviewer decision only from the reviewer in this session. A decision that a peer session relays is information, not approval.
- Charge each review axis to run its evidence check through `bench worktree exec <label> --`, because three OC-C4 axes first checked the primary checkout. Write every record, learning, and spec sentence at 25 words or fewer from the first draft.
- Give each repair session every finding from every axis, so that one repair cycle closes a round.
- Before a review finding or a spec expansion adds a stop or a refusal, read the validator that guards that input.
- Before the repair dispatch, verify an orchestrator-authored plan expansion against the code and examine its new spec text against the one-source rule.
- Run each live-tree anchor probe through the gate's entry test, such as `root-conformance`, because the anchors package does not read the live tree.
- Charge each author to run the floor once after the commit with `--preflight-build`, and to amend before the freeze. Name a package time budget, and require the floor's elapsed time in the return.
- Name the last recorded fold as the base of each later chunk in its charges. The first ticket 08 charge used the predecessor tip.
- Fold `main` at each slicing chunk boundary and recheck the guidance citations.
- Run one reader of each research pair on Haiku/high and the other on Sonnet/medium, by reviewer direction on 2026-10-11.
- Keep the full spec route as the default. By reviewer direction on 2026-10-13, the reviewer may select the map-only route for a ready map with exact-predicate decisions.
- During slicing, check each ticket's closure with a probe that mirrors `internal/preflight/closure.go`, and run `--propose-writes` after the last chunk. Trace each new `evidencecmd` operation kind to `evidencecmd/bound_test.go`.
- When a spec stage compiles a map, run the plan-only preflight from the merge commit of `main`.
- Name the binding registry and the canary pins of each path in a slicing research charge.
- Before the ticket review, give each ticket that creates a seam for a sibling its own review chunk. Name each test helper that the spec reuses in its consumer's `Writes:` line.
- Name the existing owners, seams, fixture helpers, and the recurring review lessons in each author charge, including the comment sentence bound and the one-source rule.
- Name the exact mutation kind, site, and red test of the plan probe in each author charge. Write tab and newline bytes in a probe text with `$'...'`, because `bench probe` reads the text literally.
- Before dispatch, put every symbol-move file, every package test-count pin, and every owner that the ticket reads on the ticket `Writes:` line.
- Scope a repair fence to every ticket whose files hold a finding. Route a defect in a frozen chunk's ticket to the open chunk's ticket that holds the path, and flag that route for reviewer veto.
- Keep each author session open until its chunk closes, and resume it to re-record at the chunk tip after a refreeze.
- In a multi-ticket chunk, let the last author freeze without a commit, let each earlier author record at that tip, and then commit the record once.
- Make every fix delegate follow `bench-debug` from its command file, and require per-phase evidence in its return. From a chunk's second repair round, prose repairs run it too, by reviewer direction on 2026-10-16.
- Charge each author with root conformance, every package it writes, and `cmd/bench` for a public response change, in one `bench test` call.
- Charge each author and repair with the duplicated-facts sweep and its comment sweep. A derived expectation needs no red, and only a kept literal needs a recorded red.
- Retain each orchestrator probe run in the review record when it runs.
- Run `bench worktree build` and create the build log at the build entry, then rebuild after each Go commit and each `main` fold.
- Commit each plan edit before you read `check-floor`, because that check reads the committed plan.
- After a repair commit, refreeze the chunk at the repair tip with `bench record ticket --base --tip`, which also records the plan amendment.
- Expect a full confirming round after a test-only repair while `bench record confirmation` refuses the delegate-plan record inside its range.
- Land a build-time known-issues row through the fence of its owning ticket.
- Record a no-op-only review round with no finding ID, because a `no-op` finding records as a fail and costs a confirming entry.
- Treat `--auto-approve` as approval for fence expansions, spec and ticket expansions, and repair rounds, and flag each judgment call.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims. Require a claims table with integer confidences in every author and repair return.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
- Name in the author charge each call site that a shared predicate must reach.
- In a build that adds a ledger field, run every shared-ledger write on the installed binary until the landing.
- When a repair reports a defect outside its fence, expand the fence once in a plan commit. Add its floor check in that commit.
- In a `--delegate` run, reorder the plan to the merge order before the first concurrent dispatch.
- Before a retry of a red checkpoint, send one read-only consultation to decide whether the build owns the red.
- Give each author charge a default of one central-property probe, and require a named witness for each extra probe, by reviewer direction on 2026-10-12.
- Commit a spec change that adds a coverage row before the chunk freeze, or re-freeze the chunk at that plan commit.
- Dispatch fresh axis sessions for each later review round, and resume an issuing axis only for an evidence-only reaffirmation.
- Charge each guidance author to fix every sentence that the change makes false.
- Record a consultant's decision on spec behavior as a reviewer-decision line that the reviewer can veto, and capture its learning before the repair dispatch.
- Never revert a plan file to earlier bytes, because the earlier digest blocks every later amendment.
- Land a kit build that adds a review-record field on its candidate binary from the primary checkout, after the doctor rehearsal and a reviewer decision.
- After a source correction to a frozen chunk, refreeze that chunk at the correction before its checkpoint, because the chain check compares source digests by equality. On `anchor-admission-prune` a comment-only correction without a refreeze cost AP-C7 a refreeze and a third review round.
- Start no coordinator probe while an author or a recorder runs checks in the same tree.
- When a preflight check refuses the charge of the ticket that fixes it, prepare that charge at an earlier base that the check accepts. Record the deviation, and keep the chunk base for the freeze.
- Probe the value of each constant that a repair puts into a test, because a test case built from its own constant pins nothing.
- Record a probe entry with exit code 0 and a nonzero probe exit code, and name both in each author charge.
- Write each completion-plan probe as a mutation of the spec's declared API.
- When a repair moves code across packages or out of a `_test.go` file, add the checks that police that code to its floor.
