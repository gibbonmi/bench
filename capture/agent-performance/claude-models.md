# Claude model scorecard

Last incorporated landing: `spec-stage-grader-trace` (`5a2dc13e6cc6c27caeaadbb01464283a0a1706aa`, 2026-10-06).

The `spec-stage-grader-trace` build ran one chunk and one ticket on Claude Code. An Opus/medium orchestrator dispatched a fresh Opus/high author and a fresh Opus/high repair session. By reviewer direction, Fable/high ran all six review axes. Provider usage and charges remain unknown.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 38 landings + implementer, 10 charges + reviewer, 21 axes | On `worktree-seam-reduction` the high orchestrator ran seven chunks, preserved the three-arm blind, and probed each accepted repair. All chunk checkpoints and the original landing gate passed. | Coordination of a delegated build and adversarial spec review; it implements only when the reviewer names it | orchestrator 0.36 over 1 pair; medium reviewer 0.188 over 17 pairs; 0 abstained |
| Fable / high | reviewer, 124 axes + 1 blind comparison + spec reviewer, 4 passes + control, 1 full-retrieval run + decider, 15 consultations + debug, 2 charges | On `spec-stage-grader-trace` Fable/high ran one chunk review and its confirming round, and 4 of its 7 labeled findings held. Its Coverage axis found two untested needles with a silent probe; three low-confidence findings were refuted against the author returns. | Review axes and spec review by reviewer direction; consultations whenever the coordinator is uncertain; the RE2-style full control | reviewer: 0.168 over 32 pairs, 1 abstained |
| Opus / high, xhigh | implementer, latest 10 of 32 current attempts; reviewer, 155 axes; repair, latest 10 of 27 current attempts; spec author, slicer, and reviewer, latest 10 passes | On `spec-stage-grader-trace` one Opus/high author committed the ticket in one attempt, and all seven named probes bit; it first recorded the probe exit code as 0. One Opus/high repair session closed three findings first-pass, and all nine named probes bit. | High for first-round review axes, debug runs, guidance tickets, and foundational Go seams; xhigh for a spec or a ticket that crosses owner seams; repairs at the ticket's effort. | implementation: 0.098 over 68 pairs, 2 abstained, and 0.010 over 1 new pair; reviewer: 0.163 over 70 pairs, 9 abstained; repair: 0.010 over 11 pairs |
| Opus / medium, low | implementer, latest 10 tickets; repair, latest 10 sessions; orchestrator, 3 builds; earlier reviewer and repair sessions | On `spec-stage-grader-trace` the medium orchestrator recovered a missing deliverable binding and a mid-build commitment change on `main`. It also omitted one plan row from the chunk table, which the Spec axis caught. | Medium for a CLI form ticket with named per-member probes, for a mechanical test migration with a named regression probe, and for orchestration. Low for exact tickets; a repair keeps the ticket's declared effort. | 0.129 over 53 pairs; 15 abstained |
| Sonnet / high | orchestrator, 3 landings + reviewer, 32 axes + spec reviewer, 5 passes + researcher, 3 charges + implementer, 9 comparison tickets | On `commitment-delivery-integrity` Sonnet/high ran every axis of two chunk reviews and two confirming rounds. Its Coverage axes found six real test gaps with silent probes, and its confirming rounds raised no false finding. | First-round and confirming review axes by reviewer direction, and orchestration; not fresh ticket authorship | high 0.145 over 2 pairs; xhigh 0.40 over 2 pairs |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing | unknown |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `worktree-seam-reduction` comparison | Opus and Sonnet / high / implementer; Fable / high / reviewer | One blind run ranked Sonnet-sliced first, Opus-sliced second, and Opus-full third. The reviewer selected the second arm, and all three arm gates passed. |
| `markdown-block-reader` spec review | Sonnet / high, xhigh; Opus / xhigh; Fable / high, xhigh / spec reviewer | Five reviewers graded one commit with one charge. List costs were Sonnet high $2.14, Fable high $1.98, Opus xhigh $2.44, Fable xhigh $2.50, and Sonnet xhigh $6.31. No reviewer found every confirmed finding, and only Fable high reported a false one. |
| `light-path-commitment-exemption` chunk reviews | Fable / high / reviewer | The first rounds found a publication-scope spec deviation, a duplicated fixture template, and three untested branches. Two later findings were false: a Git pathspec claim and new words for a spec-pinned diagnostic. |
| `commitment-delivery-integrity` chunk reviews | Sonnet / high / reviewer | The Coverage axes found four untested retention edges and two untested comparator terms. Each gap came from a silent probe, and each repair probe then bit. |
| `review-evidence-file-pages` RE2 control | claude-fable-5-1 / high / reviewer | The narrow round found 8 findings in 399474 tokens and the full control found 6 in 227918 tokens. They agreed on 3, so neither shape covered the other. |

## Current decisions

- Change routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Keep the `worktree-seam-reduction` comparison descriptive until a comparable repeat supports a routing change.
- Keep the `markdown-block-reader` five-reviewer spec comparison descriptive. Repeat it on later specs before a spec-review routing change.
- Run a first review round on the harness mid binding, unless the reviewer directs another line for a chunk or a spec. Escalate review to Fable only from round two on, by reviewer rule.
- Give each ticket a fresh Opus author at the spec's declared effort, and give each post-review repair a fresh Opus session at the ticket's effort.
- When an author's session is lost, give its ticket a fresh session with the `session-lost` trigger before that ticket's verification is recorded again.
- Use Opus/medium for a mechanical test migration only with a named regression probe that the author runs before and after.
- Make every fix delegate follow `bench-debug`, and accept a debug result that disproves the reported cause.
- Ask a Fable/high read-only consultation when the coordinator is uncertain, and check its claim against the tree.
- Prove a review claim about Git behavior in a scratch tree before a repair dispatch.
- Check proposed new words against the spec acceptance rows before a repair dispatch, because a row can pin the exact text.
- Do not route fresh ticket authorship to Sonnet/high until a comparable repeat confirms the latest blind result.
- Charge each author with root conformance, every package it writes, and `cmd/bench` for a public response or embed change.
- Charge each author and repair with the duplicated-facts sweep, its comment sweep, and one recorded red per independent expectation.
- Read a helper before a repair charge names it as the fold target.
- Run the propose-writes preflight before dispatch, and fold fixture paths into the Writes lines; fence expansions are pre-approved.
- When a plan commit expands a fence into a package, update each spec sentence about that package in the same commit.
- List the count-pin file in the Writes line of a ticket that adds a top-level test to a pinned package.
- Keep the narrow review shape provisional. Pair narrow axes with one full reader in a later trial.
- Probe author claims independently, and rerun every ticket's verification at the final chunk source before the chunk review.
- Record ticket assignments before each dispatch, and record each plan amendment as an identity chunk mapping when chunks do not change.
- Give each verification of a version 2 chunk an ID that is unique in that chunk.
- Add a repair test's coverage row in the repair plan commit, before the repair dispatch and the confirming round.
- Run tests and probes serially when review axes or authors share a tree.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims.
- Apply the repair allowance to blocking findings that change source, and return exhaustion to the reviewer. A record-only or comment-only change consumes no repair cycle.
- Correct spec prose inside the chunk that needs it, because the checkpoint needs a passing Spec result and a spec edit moves the frozen source.
- Limit a confirming round on record-only corrections to the findings it confirms, and treat a new preference in that prose as advice.
- Do not land an unrelated change on `main` while a chunked build is open. If main moves anyway, fold it only with reviewer approval, and record an empty diff of each fold-only fence entry before the landing.
- If a commitment change on `main` refuses the build's commits before the first review, merge `main` and move the chunk base. Flag that fold for reviewer veto.
- Name every anchor needle in the plan probe list with its expected verdict, so a needle without a canary shows before dispatch.
- Write each fact in a delegate charge from git output, not from memory.
- Run the whole conformance package for a guidance ticket, because a Markdown lane misses a test that reads guidance.
- Write chunk, verification, review, and amendment entries with `bench record`. At each freeze after a plan change, run the amendment form before the chunk form.
- Record `--probe-exit-code` as the exit of the mutated test run, which is 1 when the probe bites.
- Charge each author of a CLI form with one probe for each member of each refusal step that the form joins.
- Keep a built executable out of a worktree before its checkpoint, because a suite test rewrites its broker manifest. Build a scratch binary with `--manifest-dir` outside the worktree.
- Charge each author to edit pool-path files only with the file tools or through `bench worktree exec`, never with a script on the pool path.
- Stop a consultation that produces no output past its expected window, and reroute it.
- Read the assignment census before landing releases the worktree.
- Bind a deliverable to its outcome through a landed plan before `bench commitment start`, because start reads the policy at `main`. Land the staged spec first, because the plan reads the deliverable at `main`.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
- Run the doctor rehearsal on the candidate build before a landing that changes a Bench build input, and run `bin/bench.sh doctor --fix` after that landing.
