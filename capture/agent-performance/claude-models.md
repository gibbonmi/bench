# Claude model scorecard

Last incorporated landing: `build-directions` (`287f23d8935168f2acdcfb98421b07e82c17ed42`, 2026-10-10).

The `build-directions` build ran four chunks and six tickets on Claude Code. An Opus/high orchestrator dispatched six fresh Opus/high authors and four Opus/high repair sessions, and Opus/high ran every review axis. By reviewer direction, seven Opus shadow authors ran an effort benchmark on tickets 01, 04, and 06, and three Fable/high blind graders scored every variant. Provider usage and charges remain unknown.

The `adoption-lifecycle-consumption` build ran three chunks and four tickets on Claude Code in 11 hours. A Fable orchestrator ran it until a pause, and an Opus orchestrator finished it. It dispatched five Opus/high authors, six Opus/high repair sessions, and two verification-only sessions. By reviewer direction, Opus/medium ran 30 review axes in parallel. Sonnet/high ran one consultation, three staleness slices, and one fix review. An Opus/high debug delegate fixed a lock leak, and provider usage and charges remain unknown.

The `native-record-operations` build ran three chunks on Claude Code. A Fable/high orchestrator dispatched three fresh Opus/high authors, one Opus/high successor after a reviewer restart, and three Opus/high repair sessions. By reviewer direction, every review axis ran on Fable/medium, nine rounds in all, three axes in parallel with one probing axis. Haiku ran one read-only roadmap match. Provider usage and charges remain unknown; delegate-reported tokens were 190k to 370k per author and 20k to 110k per axis.


## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 40 landings + implementer, 10 charges + reviewer, 48 axes | On `adoption-lifecycle-consumption` Fable orchestrated the first half without a witness cap or a package time budget in its charges, so ticket 01 ran 63 probes and the adopt package doubled unseen. It declared no expected repair rounds. | Coordination of a delegated build and adversarial spec review; review axes in parallel by reviewer direction; it implements only when the reviewer names it | orchestrator 0.36 over 1 pair, 1 abstained; medium reviewer 0.256 over 20 pairs, 1 abstained |
| Fable / high | reviewer, 124 axes + 4 blind comparisons + spec reviewer, 4 passes + control, 1 full-retrieval run + decider, 17 consultations + debug, 2 charges | On `build-directions` three Fable/high blind graders scored ten effort variants under shuffled labels. One grader found a pasted fixture harness that the Opus/high chunk review missed. | Review axes and spec review by reviewer direction; consultations whenever the coordinator is uncertain; the RE2-style full control | reviewer: 0.168 over 32 pairs, 1 abstained |
| Opus / high, xhigh | implementer, latest 10 of 57 current attempts; reviewer, 197 axes; repair, latest 10 of 50 current attempts; orchestrator, 2 builds; spec author, slicer, and reviewer, latest 10 passes | On `build-directions` six high authors committed inside their fences, and every coordinator probe bit. Blind graders ranked high first on one ticket, second on one, and last of four on the code ticket; one xhigh shadow ranked first on the code ticket at 1.6 times the tokens of high. | High for first-round review axes, orchestration, and debug runs; xhigh for a code ticket where a repair round costs more than about 100k tokens, or that crosses owner seams; repairs at the ticket's effort. | implementation: 0.098 over 68 pairs, 2 abstained, 0.028 over 5 pairs, 0.097 over 4 pairs, and 0.01 over 6 new pairs; orchestrator 0.36 over 1 pair; reviewer: 0.162 over 74 pairs, 9 abstained; repair: 0.010 over 11 pairs |
| Opus / medium, low | implementer, latest 10 tickets + 5 benchmark shadows; repair, latest 10 sessions; orchestrator, 5 builds; reviewer, 30 axes | On `build-directions` low shadows put the same production bytes as high on two guidance tickets at about half the tokens, with thinner red evidence and more stale comments. A medium shadow matched the xhigh grade on the code ticket at 0.8 times the tokens of high. | Low for an exact-spec guidance ticket; medium as the default for a code ticket, for review axes by reviewer direction, and for a mechanical test migration; a repair keeps the ticket's declared effort. | 0.129 over 67 pairs; 18 abstained; reviewer 0.156 over 45 pairs |
| Sonnet / high | orchestrator, 3 landings + reviewer, 132 axes + spec reviewer, 5 passes + researcher, 3 charges + implementer, 9 comparison tickets | On `adoption-lifecycle-consumption` Sonnet/high ran three staleness slices, one consultation, and one fix review. Its slices found every stale citation, its consultation verified the facts it was given, and its fix review found a mixed-binary race and a silent mutant. | First-round and confirming review axes, consultations, and fix reviews by reviewer direction, and orchestration; not fresh ticket authorship | high 0.145 over 2 pairs; xhigh 0.40 over 2 pairs; reviewer 0.203 over 4 pairs; consultant 0.09 over 1 pair |
| Sonnet / high (research and polish) | researcher, 6 paired charges + polish, 1 session | On `stdlib-reimplementation-check` Sonnet made no wrong claim in six research slices and found every key fact under the paired-reader rules. Its high polish of a Haiku draft closed nine audit items in one pass and found a hole in the live-tree meta check. | Default research reader, paired with Haiku; polish of a cheap draft when the reviewer names it | research 6 of 6 slices clean; polish 1 pass |
| Haiku / default–xhigh | researcher, 12 charges + implementer, 2 comparison tickets + auditor, 3 charges | On `native-record-operations` Haiku matched 12 workflow bottlenecks against the roadmap, the known issues, and the learnings. It cited seven open rows with quoted lines, named its search terms for the five unmatched, and flagged one unknown instead of guessing. | Research reader, paired with Sonnet; staleness audits; not ticket authorship | research 10 of 12 slices clean; authorship 0 landings |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing | unknown |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `worktree-seam-reduction` comparison | Opus and Sonnet / high / implementer; Fable / high / reviewer | One blind run ranked Sonnet-sliced first, Opus-sliced second, and Opus-full third. The reviewer selected the second arm, and all three arm gates passed. |
| `markdown-block-reader` spec review | Sonnet / high, xhigh; Opus / xhigh; Fable / high, xhigh / spec reviewer | Five reviewers graded one commit with one charge. List costs were Sonnet high $2.14, Fable high $1.98, Opus xhigh $2.44, Fable xhigh $2.50, and Sonnet xhigh $6.31. No reviewer found every confirmed finding, and only Fable high reported a false one. |
| `build-directions` effort benchmark | Opus / low, medium, high, xhigh / implementer; Fable / high / blind grader | Ten first-pass commits on three tickets met every acceptance row. On two guidance tickets low matched high at half the tokens. On the code ticket xhigh ranked first, medium second, and high last of four. |
| `commitment-delivery-integrity` chunk reviews | Sonnet / high / reviewer | The Coverage axes found four untested retention edges and two untested comparator terms. Each gap came from a silent probe, and each repair probe then bit. |
| `review-evidence-file-pages` RE2 control | claude-fable-5-1 / high / reviewer | The narrow round found 8 findings in 399474 tokens and the full control found 6 in 227918 tokens. They agreed on 3, so neither shape covered the other. |

## Current decisions

- Read this scorecard before each line declaration, because its decisions carry the orchestration rules that a build repeats.
- Declare the expected repair rounds and a confidence in each line declaration, so the orchestrator claim gets a label.
- Route Opus ticket authors by ticket type, by reviewer decision on `build-directions`: low for an exact-spec guidance ticket, medium as the code-ticket default.
- Route a code ticket to Opus/xhigh when a repair round would cost more than about 100k tokens, or when the ticket crosses owner seams.
- Change other routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Keep the `worktree-seam-reduction` comparison and the `markdown-block-reader` five-reviewer comparison descriptive until a comparable repeat.
- Run review axes on the line that the reviewer directs, and let only the Coverage axis run `bench probe` on the shared tree with focused selections.
- Give each repair session every finding from every axis, so that one repair cycle closes a round.
- Run each anchor-row omission probe through the docs-currency-workflow check, because the anchor harness grades a synthetic tree.
- Name a package time budget in each author charge, and require the floor's elapsed time in the return.
- When a delegate nears 300k tokens, let it finish its ticket, and send later work to a fresh recorded session.
- Name the existing owners, seams, fixture helpers, and the recurring review lessons in each author charge, including the Forbid-mirror rule and the comment sentence bound.
- Before dispatch, put every symbol-move file and every package test-count pin on the ticket `Writes:` line.
- Scope a repair fence to every ticket whose files hold a finding.
- Plan a verification-only session for each sibling ticket before a repair commits in a multi-ticket chunk.
- Make every fix delegate follow `bench-debug` from its command file, and require per-phase evidence in its return.
- Ask a read-only consultation on the line the reviewer directs when the coordinator is uncertain, and use a different-family blind grader for a benchmark.
- Do not route fresh ticket authorship to Sonnet/high until a comparable repeat confirms the latest blind result.
- Charge each author with root conformance, every package it writes, and `cmd/bench` for a public response change, in one `bench test` call.
- Charge each author and repair with the duplicated-facts sweep, its comment sweep, and one recorded red per independent expectation.
- Run `bench worktree build` after each Go commit and after each `main` fold in a build worktree.
- Run `bench record amendment` after each plan commit, and record each assignment with the worktree's own assignment id.
- After a repair commit, refreeze the chunk at the repair tip before the repair session records its verification.
- Fold a review's spec-only fix into the repair commit, because a plan commit after the chunk tip forces a refreeze.
- Take the predecessor chunk tip as the next chunk base when a comment-only correction follows that tip.
- Commit every record before preparing a review charge, and write no record while review axes bind their evidence.
- Treat `--auto-approve` as approval for fence expansions, spec and ticket expansions, and repair rounds, and flag each judgment call.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims.
- Apply the repair allowance to blocking findings that change source; a record-only or comment-only correction consumes no repair cycle.
- Read the assignment census before landing releases the worktree.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
- Run the doctor rehearsal on the candidate build before a landing that changes a Bench build input.
