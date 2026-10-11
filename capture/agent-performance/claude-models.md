# Claude model scorecard

Last incorporated landing: `spec-stage-binding` build (landing commit `88de8d9d16ea1ab595687515504cbf006e66a5f2`, 2026-10-11).

The `spec-stage-binding` build ran three chunks and three tickets on Claude Code in about 6 hours with the reviewer away. A Fable/medium orchestrator dispatched three fresh Opus/xhigh authors and two Opus/xhigh repair sessions, and Opus/high ran every review axis in six rounds. One Opus/high consultation decided the stale-base fold. Provider usage and charges remain unknown. The harness reported 396k, 280k, and 266k tokens for the authors, 245k and 180k for the repairs, and 34k to 147k for each axis.

The `spec-stage-binding` spec stage ran on Claude Code, with one compaction before the slicing research. An Opus orchestrator ran paired Haiku/xhigh and Sonnet/high readers for each research slice, and an Opus spec-author fork. Fresh Opus slicers ran at high for two chunks and at xhigh for the third, by reviewer direction. Fable/xhigh accepted the spec review and the ticket review in one iteration each, and an Opus/max delegate shipped one light-path fix by reviewer direction. Provider usage and charges remain unknown. The harness reported 148k, 131k, and 275k tokens for the slicers, 296k for the ticket review, and 238k for the light-path author.

The `commit-merge-parent` build ran three chunks and five tickets on Claude Code in about 8 hours. A Fable orchestrator dispatched five fresh Opus/xhigh authors and four Opus/xhigh repair sessions. Opus/high ran every review axis in eight rounds. By reviewer direction, three Opus/max shadow authors ran an effort assessment on tickets 1 to 3, and no consultation was needed. Provider usage and charges remain unknown. The harness reported 210k, 283k, and 257k tokens for the xhigh authors, and 301k, 399k, and 369k for the max shadows.

The `commit-merge-parent` spec stage ran on Claude Code. An Opus orchestrator ran paired Sonnet/high research readers, an Opus spec-author fork, and three fresh Opus/max slicers, one for each chunk. Fable/high accepted the spec review and the ticket review in one iteration each. Provider usage and charges remain unknown; the harness reported 258k, 383k, and 278k tokens for the three slicers and 139k for the ticket review.


## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 41 landings + implementer, 10 charges + reviewer, 48 axes | On `spec-stage-binding` Fable/medium orchestrated the whole build with the reviewer away, routed one judgment call to a consultant, and folded every record gap as evidence. It declared 3 expected repair rounds at confidence 6 and the build used 2 cycles, so the claim held. | Coordination of a delegated build and adversarial spec review; review axes in parallel by reviewer direction; it implements only when the reviewer names it | orchestrator 0.36 over 2 pairs, 1 abstained; medium reviewer 0.256 over 20 pairs, 1 abstained |
| Fable / high, xhigh | reviewer, 124 axes + 4 blind comparisons + spec reviewer, 6 passes + ticket reviewer, 2 passes + control, 1 full-retrieval run + decider, 17 consultations + debug, 2 charges | On `spec-stage-binding` Fable/xhigh accepted the spec and the ticket set in one iteration each, with eight and nine non-blocking findings. Its ticket review re-derived each ticket contract from the tree and found the stale citations that two main moves left. | Review axes, spec review, and ticket review by reviewer direction; consultations whenever the coordinator is uncertain; the RE2-style full control | reviewer: 0.168 over 32 pairs, 1 abstained |
| Opus / high, xhigh, max | implementer, latest 10 of 69 current attempts, 3 of them max shadows and 1 a max light-path fix; reviewer, 240 axes + 1 consultation; repair, latest 10 of 56 current attempts; orchestrator, 2 builds and 2 spec stages; spec author, slicer, and reviewer, latest 10 passes | On `spec-stage-binding` three xhigh authors each landed in one attempt with every row red first, and two xhigh repairs each closed every finding in one attempt. Nineteen high axes found one real defect, the stale-base refusal, and otherwise record gaps and comment form; each author kept only its central probe in the record. | High for first-round review axes, orchestration, consultations, and debug runs; xhigh for a code ticket where a repair round costs more than about 100k tokens, or that crosses owner seams; repairs at the ticket's effort. | implementation: 0.098 over 68 pairs, 2 abstained, and 0.01 over 8 xhigh pairs; orchestrator 0.36 over 1 pair; reviewer: 0.162 over 74 pairs, 9 abstained; repair: 0.010 over 13 pairs; consultant 0.01 over 1 pair |
| Opus / medium, low | implementer, latest 10 tickets + 5 benchmark shadows; repair, latest 10 sessions; orchestrator, 5 builds; reviewer, 30 axes | On `build-directions` low shadows put the same production bytes as high on two guidance tickets at about half the tokens, with thinner red evidence and more stale comments. A medium shadow matched the xhigh grade on the code ticket at 0.8 times the tokens of high. | Low for an exact-spec guidance ticket; medium as the default for a code ticket, for review axes by reviewer direction, and for a mechanical test migration; a repair keeps the ticket's declared effort. | 0.129 over 67 pairs; 18 abstained; reviewer 0.156 over 45 pairs |
| Sonnet / high | orchestrator, 3 landings + reviewer, 132 axes + spec reviewer, 5 passes + researcher, 3 charges + implementer, 9 comparison tickets | On `adoption-lifecycle-consumption` Sonnet/high ran three staleness slices, one consultation, and one fix review. Its slices found every stale citation, its consultation verified the facts it was given, and its fix review found a mixed-binary race and a silent mutant. | First-round and confirming review axes, consultations, and fix reviews by reviewer direction, and orchestration; not fresh ticket authorship | high 0.145 over 2 pairs; xhigh 0.40 over 2 pairs; reviewer 0.203 over 4 pairs; consultant 0.09 over 1 pair |
| Sonnet / high (research and polish) | researcher, 6 paired charges + 9 slicing-research charges + polish, 1 session | On `spec-stage-binding` the charge named the binding registry and the canary pins, and the three Sonnet slicing readers found every closure class. One reader also found a planned line join that broke three byte-exact canary needles. | Default research reader, paired with Haiku; polish of a cheap draft when the reviewer names it | research 6 of 6 earlier slices clean; slicing research missed the closure classes in 6 of 6 uncharged slices and in 0 of 3 charged slices |
| Haiku / default–xhigh | researcher, 12 charges + 3 slicing-research charges + implementer, 2 comparison tickets + auditor, 3 charges | On `spec-stage-binding` the three Haiku/xhigh slicing readers matched the Sonnet closure sets in each chunk and also found the line-join hazard. Only the Sonnet readers named the shared `CHANGELOG.md` closure for the first two chunks. | Research reader, paired with Sonnet; staleness audits; not ticket authorship | research 10 of 12 slices clean; slicing research agreed with Sonnet in 3 of 3 chunks; authorship 0 landings |
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
- Declare the expected repair rounds and a confidence in each line declaration, so the orchestrator claim gets a label. Expect about one repair cycle per chunk on an owner-seam build.
- Route Opus ticket authors by ticket type: low for an exact-spec guidance ticket, and medium as the code-ticket default. Use xhigh for a code ticket that crosses owner seams or whose repair round costs more than about 100k tokens.
- Do not route a ticket author to Opus/max, because three max shadows on `commit-merge-parent` gave the xhigh design at 1.4 times the tokens.
- Change other routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Run review axes on the line that the reviewer directs, and let only the Coverage axis run `bench probe` on the shared tree with focused selections.
- When the reviewer is away, route each build judgment call to an Opus/high read-only consultant. Its one consultation on `spec-stage-binding` verified the cause in the code and chose the cheaper route.
- Require one excerpt file per probe with its command line first in every author charge. Every author on `spec-stage-binding` kept only its central probe in the record, and each first review found the gap.
- Batch prose-only spec corrections before a repair's refreeze. The source digest covers spec, test, and comment bytes, so a later plan or comment-only commit costs a refreeze and a verification rerun.
- Write every review pickup sentence at 25 words or fewer from the first draft. The lane's prose check graded the record four times on `spec-stage-binding`.
- Give each repair session every finding from every axis, so that one repair cycle closes a round.
- Run each live-tree anchor probe through the docs-currency-workflow check, and keep the ticket's mirror route for the anchor-row omission.
- Name a package time budget in each author charge, and require the floor's elapsed time in the return.
- When a delegate nears 300k tokens, let it finish its ticket or its slicing chunk, and send later work to a fresh recorded session.
- Fold `main` at each slicing chunk boundary and recheck the guidance citations, because two moves on `spec-stage-binding` left citations stale.
- Run one reader of each research pair on Haiku/xhigh and the other on Sonnet/high, by reviewer direction on 2026-10-10.
- Name the binding registry and the canary pins of each path in a slicing research charge, because six readers on `commit-merge-parent` missed those closure files.
- Name the existing owners, seams, fixture helpers, and the recurring review lessons in each author charge, including the comment sentence bound and the one-source rule.
- Name the exact mutation kind, site, and red test of the plan probe in each author charge. Keep every variable and import in use in a probe text.
- Before dispatch, put every symbol-move file and every package test-count pin on the ticket `Writes:` line.
- Scope a repair fence to every ticket whose files hold a finding.
- Keep each author session open until its chunk closes, and resume it to re-record at the chunk tip after a refreeze.
- Make every fix delegate follow `bench-debug` from its command file, and require per-phase evidence in its return.
- Charge each author with root conformance, every package it writes, and `cmd/bench` for a public response change, in one `bench test` call.
- Charge each author and repair with the duplicated-facts sweep, its comment sweep, and one recorded red per independent expectation.
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
