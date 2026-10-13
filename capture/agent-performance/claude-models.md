# Claude model scorecard

Last incorporated landing: `probe-mutant-sets` build (landing `062f9e41cbab5e2a0f51a8a80d5f01ca790680e5`, 2026-10-13).

The `probe-mutant-sets` build ran six chunks and six tickets on Claude Code as a `--full` run with an author limit of 1. An Opus orchestrator dispatched five fresh Opus/high authors, four Opus/high code repairs, and five Opus/low small fixes and prose repairs. By reviewer direction, `gemini-4-argon-high` wrote ticket 6 through the `agy` CLI. Opus/high ran 39 review axes and one consultation. PM-C1 used no repair cycle, PM-A and PM-B used one, PM-D and PM-E used two, and PM-C2 used three.

Provider usage and charges remain unknown. The harness reported 91k to 367k tokens for each author, 29k to 204k for each repair, and 27k to 162k for each axis.

The `build-landing-window` build ran six chunks and eight tickets on Claude Code across two sessions and one compaction. An Opus orchestrator dispatched eight fresh Opus/high authors and eleven Opus/high repair sessions, and each repair ran `bench-debug`. Opus/high ran every review axis in fourteen rounds, and Fable/high answered two consultations. BW-C and BW-F used both repair cycles, and the other chunks used one each. The installed binary could not read the new fold record, so the build landed on its candidate binary by reviewer decision.

Provider usage and charges remain unknown. The harness reported 259k to 450k tokens for each author, 122k to 195k for each repair, and 34k to 137k for each axis.

The `durable-file-replacement` build ran one ticket in one chunk on Claude Code, in parallel with the FT342-FT253 build. An Opus/medium orchestrator dispatched one fresh Opus/high author, one Opus/high repair session, and one Opus/high evidence-refresh session. Opus/high ran every review axis in four rounds, and Fable/high answered one seam consultation. Two Sonnet staleness slices found a second caller and a missing floor check. A reviewer rescope of D32 after the freeze forced a full evidence refresh with no code change.

Provider usage and charges remain unknown. The harness reported 236k tokens for the author, 169k for the repair, 107k for the refresh, and 26k to 77k for each axis.

The `build-landing-window` spec stage ran on Claude Code for FT342-FT253. An Opus coordinator ran paired Opus readers for four research slices and six slicing chunks. That broke the 2026-10-10 Haiku and Sonnet pair decision. The coordinator also sliced all six chunks itself instead of one fresh Opus/xhigh slicer per chunk. An Opus fork wrote the spec. Fable/high took 2 iterations on the spec review and accepted the ticket review in 1.

Provider usage and charges for `build-landing-window` remain unknown. The harness reported about 1.11M tokens for the research readers and 479k for the spec fork. It reported 1.57M for the slicing readers, 162k for the spec review, and 142k for the ticket review.

The `guidance-weight` spec stage ran on Claude Code beside a reviewer-directed test of map-only slicing. An Opus/xhigh coordinator ran paired Haiku/xhigh and Sonnet/high readers for four research slices and five slicing slices. An Opus/xhigh fork wrote the spec, and one fresh Opus/xhigh slicer cut each of five chunks. Fable/high answered one consult, took 2 iterations on the spec review, and accepted the ticket review in 1. A fresh Opus/xhigh map-only slicer wrote 7 tickets for about 0.54M harness tokens, against about 6.6M for the spec route. Provider usage and charges remain unknown.

The `probe-subjects-and-mutations` build ran five chunks and six tickets on Claude Code as a `--full` run with an author limit of 1. An Opus/medium orchestrator dispatched six fresh Opus/high authors, six repair sessions, and four Sonnet staleness slices. Opus/medium ran every review axis. By reviewer direction, code repairs ran at Opus/high and small fixes at Opus/low, and every charge used minimal probes. PS-B used both repair cycles, and the other chunks used one each.

Provider usage and charges remain unknown. Authors ran 13 probes in total, and every extra probe named a witness that a test-first red could not show.

The `probe-mutant-sets` spec stage ran on Claude Code for FT168-probe-batch. An Opus coordinator ran paired Sonnet/medium and Haiku/high readers for four research slices and five slicing chunks, and it sliced the six tickets itself. An Opus fork wrote the spec. By reviewer direction, Opus/xhigh ran the spec review and the ticket review, and each took 2 iterations. The spec round found an interrupt that read `refused` and a lost cross-half rule. The ticket round found a shared seam chunk, a missing helper file, and the missing commitment file.

Provider usage and charges remain unknown. The harness reported about 1.02M tokens for the research readers and 303k for the spec fork. It reported about 1.18M for the slicing readers, 268k for the spec review, and 214k for the ticket review.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 41 landings + implementer, 10 charges + reviewer, 48 axes | On `spec-stage-binding` Fable/medium orchestrated the whole build with the reviewer away, routed one judgment call to a consultant, and folded every record gap as evidence. It declared 3 expected repair rounds at confidence 6 and the build used 2 cycles, so the claim held. | Coordination of a delegated build and adversarial spec review; review axes in parallel by reviewer direction; it implements only when the reviewer names it | orchestrator 0.36 over 2 pairs, 1 abstained; medium reviewer 0.256 over 20 pairs, 1 abstained |
| Fable / high, xhigh | reviewer, 124 axes + 4 blind comparisons + spec reviewer, 10 passes + ticket reviewer, 4 passes + verifier, 3 light-path checks + control, 1 full-retrieval run + decider, 21 consultations + debug, 2 charges | On `build-landing-window` two Fable/high consultations chose the ancestry owner, the fold content binding, the gap order rule, and the Git fault refusal. Each repair used those choices, and every confirming axis accepted them. | Review axes, spec review, and ticket review by reviewer direction; consultations whenever the coordinator is uncertain; the RE2-style full control | reviewer: 0.168 over 32 pairs, 1 abstained |
| Opus / high, xhigh, max | implementer, latest 10 of 90 current attempts, 3 of them max shadows, 1 a max light-path fix, and 2 xhigh light-path fixes; reviewer, 319 axes + 2 consultations; researcher, 20 paired charges; repair, latest 10 of 75 current attempts; orchestrator, 5 builds and 4 spec stages; spec author, slicer, and reviewer, latest 10 passes, 7 of them xhigh slicers | On `probe-mutant-sets` five Opus/high authors met their rows inside their fences, but first reviews found unpinned states, duplicated fixtures, and one test-only field. Fresh Opus/high axes caught every defect before the landing gate except one freshness topology red at a chunk checkpoint. | High for first-round review axes, orchestration, consultations, and debug runs; xhigh for a code ticket where a repair round costs more than about 100k tokens, or that crosses owner seams; repairs at the ticket's effort. | implementation: 0.103 over 70 pairs, 2 abstained, and 0.01 over 8 xhigh pairs; orchestrator 0.26 over 3 pairs; reviewer: 0.162 over 74 pairs, 9 abstained; repair: 0.010 over 14 pairs; consultant 0.01 over 1 pair |
| Opus / medium, low | implementer, latest 10 of 18 tickets + 5 benchmark shadows; repair, latest 10 of 32 sessions; orchestrator, 6 builds; reviewer, latest 10 of 76 axes | On `probe-mutant-sets` five Opus/low small fixes and prose repairs each closed their items in 3 to 5 minutes. Two of them left a new one-line slip, a lost conjunction and dropped articles, that a further round caught. | Low for an exact-spec guidance ticket; medium as the default for a code ticket, for review axes by reviewer direction, and for a mechanical test migration; a repair keeps the ticket's declared effort. | 0.129 over 67 pairs; 18 abstained; reviewer 0.156 over 45 pairs; orchestrator 0.25 over 1 pair |
| Sonnet / high | orchestrator, 3 landings + reviewer, 132 axes + spec reviewer, 5 passes + researcher, 3 charges + implementer, 9 comparison tickets | On `probe-subjects-and-mutations` four Sonnet staleness slices at default effort found the injected replacer seam, a new test file, and four anchor registries missing from the fences. Each slice named what it did not verify. | First-round and confirming review axes, consultations, and fix reviews by reviewer direction, and orchestration; not fresh ticket authorship | high 0.145 over 2 pairs; xhigh 0.40 over 2 pairs; reviewer 0.203 over 4 pairs; consultant 0.09 over 1 pair |
| Sonnet / high (research and polish) | researcher, 10 paired charges + 14 slicing-research charges + polish, 1 session | On `guidance-weight` the Sonnet readers found the extra pointer lines, the structure caps, and the `CONTEXT.md` reader of the prose budget. Each pair agreed on every closure class. | Default research reader at medium effort, paired with Haiku/high by reviewer direction on 2026-10-11; polish of a cheap draft when the reviewer names it | research 6 of 6 earlier slices clean; slicing research missed the closure classes in 6 of 6 uncharged slices and in 0 of 3 charged slices |
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
- Declare the expected repair rounds and a confidence in each line declaration, so the orchestrator claim gets a label. Expect one to two repair cycles per chunk: 5 chunks on `probe-subjects-and-mutations` used 6 cycles.
- Route every coding delegate to Opus/high and every small-fix delegate to Opus/low, by reviewer direction on 2026-10-12. A small fix changes a comment, a text, or a few lines with no new branch.
- Do not route a ticket author to Opus/max, because three max shadows on `commit-merge-parent` gave the xhigh design at 1.4 times the tokens.
- Change other routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Run review axes on the line that the reviewer directs, and let only the Coverage axis run `bench probe` on the shared tree with focused selections.
- When the reviewer is away, route each build judgment call to an Opus/high read-only consultant. Its one consultation on `spec-stage-binding` verified the cause in the code and chose the cheaper route.
- Require one excerpt file per probe with its command line first in every author charge. Every author on `spec-stage-binding` kept only its central probe in the record, and each first review found the gap.
- Batch prose-only spec corrections before a repair's refreeze. The source digest covers spec, test, and comment bytes, so a later plan or comment-only commit costs a refreeze and a verification rerun.
- Ask the reviewer about each host-bound acceptance row, such as a native macOS run, before the first chunk freeze. A rescope after the freeze on `durable-file-replacement` cost a rerun of 14 probes and four axis entries.
- Accept a reviewer decision only from the reviewer in this session. A decision that a peer session relays is information, not approval.
- Charge each review axis to run its evidence check through `bench worktree exec <label> --`, because three axes on `durable-file-replacement` first ran it from the primary checkout.
- Write every review pickup sentence at 25 words or fewer from the first draft. The lane's prose check graded the record four times on `spec-stage-binding`.
- Give each repair session every finding from every axis, so that one repair cycle closes a round.
- Run each live-tree anchor probe through the gate's entry test, such as `root-conformance`, because the anchors package does not read the live tree.
- Name a package time budget in each author charge, and require the floor's elapsed time in the return.
- When a delegate nears 300k tokens, let it finish its ticket or its slicing chunk, and send later work to a fresh recorded session.
- Fold `main` at each slicing chunk boundary and recheck the guidance citations, because two moves on `spec-stage-binding` left citations stale.
- Run one reader of each research pair on Haiku/high and the other on Sonnet/medium, by reviewer direction on 2026-10-11.
- Keep the spec route for a build. On `guidance-weight` a map-only Opus/xhigh slicer matched its tickets at one twelfth of the tokens. Repeat that comparison on a map with a code seam before a route change.
- During slicing, check each ticket's closure with a probe that mirrors `internal/preflight/closure.go`, and run `--propose-writes` after the last chunk.
- When a spec stage compiles a map, run the plan-only preflight from the merge commit of `main`. The map move then stays outside the graded range.
- Name the binding registry and the canary pins of each path in a slicing research charge, because six readers on `commit-merge-parent` missed those closure files.
- Before the ticket review, give each ticket that creates a seam for a sibling its own review chunk. Name each test helper that the spec reuses in its consumer's `Writes:` line. Both misses on `probe-mutant-sets` cost a second ticket-review iteration.
- Name the existing owners, seams, fixture helpers, and the recurring review lessons in each author charge, including the comment sentence bound and the one-source rule.
- Name the exact mutation kind, site, and red test of the plan probe in each author charge. Keep every variable and import in use in a probe text.
- Before dispatch, put every symbol-move file and every package test-count pin on the ticket `Writes:` line.
- Scope a repair fence to every ticket whose files hold a finding.
- Keep each author session open until its chunk closes, and resume it to re-record at the chunk tip after a refreeze.
- Make every fix delegate follow `bench-debug` from its command file, and require per-phase evidence in its return.
- Charge each author with root conformance, every package it writes, and `cmd/bench` for a public response change, in one `bench test` call.
- Charge each author and repair with the duplicated-facts sweep and its comment sweep. Give authors and axes one probe standard: a derived expectation needs no red, and only a kept literal needs a whole-package red. Omission probes run with `--run`, and a repair replays only the reds that its diff touches.
- Retain each orchestrator probe run in the review record when it runs, because a prose claim without the block costs a reaffirmation pass.
- Run `bench worktree build` after each Go commit and after each `main` fold in a build worktree.
- Run `bench record amendment` after each plan commit, and before a freeze whose ticket rewrote seam cells.
- After a repair commit, refreeze the chunk at the repair tip before the repair session records its verification.
- Record a no-op-only review round with no finding ID, because a `no-op` finding records as a fail and costs a confirming entry.
- Treat `--auto-approve` as approval for fence expansions, spec and ticket expansions, and repair rounds, and flag each judgment call.
- Grade a finding on a hardening-cycle row as optional advice under the bounded repair policy unless a blocking condition holds.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims.
- Apply the repair allowance to blocking findings that change source; a record-only or comment-only correction consumes no repair cycle.
- Read the assignment census before landing releases the worktree.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
- Run the doctor rehearsal on the candidate build before a landing that changes a Bench build input.
- In a build that adds a ledger field, run every shared-ledger write on the installed binary until the landing. A worktree-binary checkpoint on `build-landing-window` wrote a record that every other session refused.
- Land a kit build that adds a review-record field on its candidate binary from the primary checkout, after the doctor rehearsal and a reviewer decision. Set `BENCH_HOME`, and leave `BENCH_KIT` unset so that the census keeps the system phase.
- When a repair reports a defect outside its fence, expand the fence once in a plan commit, so that one confirming round covers both repairs.
- In a `--delegate` run, reorder the plan to the merge order before the first concurrent dispatch. A later reorder blocks a chunk checkpoint.
- When a chunk ends with a comment-only correction, freeze the next chunk on that chunk tip, because the chain check compares source digests by equality.
- Before a retry of a red checkpoint, send one read-only consultation to decide whether the build owns the red.
- Record a probe entry with exit code 0 and a nonzero probe exit code, and name both in each author charge.
- Give each author charge a default of one central-property probe, and require a named witness for each extra probe, by reviewer direction on 2026-10-12.
- When a plan expansion adds a package to a fence, add its floor check in the same plan commit. List every caller of each changed symbol there too.
- Before a plan expansion commits, examine its new spec text against the one-source rule, because two such texts on `probe-subjects-and-mutations` became findings.
- Commit a spec change that adds a coverage row before the chunk freeze, or re-freeze the chunk at that plan commit.
- Require native copies of spilled responses as verification excerpts, never a hand summary.
- Write each completion-plan probe as a mutation of the spec's declared API. On `probe-mutant-sets`, a probe that needed a missing field made an author add one.
- When a repair moves code across packages or out of a `_test.go` file, add the checks that police that code to its floor. A moved publication call reached a chunk checkpoint red on `probe-mutant-sets`.
- Dispatch fresh axis sessions for each later review round, and resume an issuing axis only for an evidence-only reaffirmation.
- Charge each guidance author to fix every sentence that the change makes false. Charge each prose repair to keep the articles and to reread a sentence after it removes a list item.
- Record a consultant's decision on spec behavior as a reviewer-decision line that the reviewer can veto, and capture its learning before the repair dispatch.
