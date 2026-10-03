# Claude model scorecard

Last incorporated landing: `worktree-seam-reduction` (`eaa2201cea5fd35cbedb79bf91a5b41a9d5edad1`, 2026-10-02).

The `worktree-seam-reduction` build ran in Claude Code with a Fable/high orchestrator. Opus/high sessions implemented, repaired, and reviewed 15 tickets in seven chunks. Fable/high reviewed the three comparison arms, and Sonnet/high implemented the six-ticket arm C. Provider usage and charges remain unknown.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 38 landings + implementer, 10 charges + reviewer, 21 axes | On `worktree-seam-reduction` the high orchestrator ran seven chunks, preserved the three-arm blind, and probed each accepted repair. All chunk checkpoints and the original landing gate passed. | Coordination of a delegated build and adversarial spec review; it implements only when the reviewer names it | orchestrator 0.36 over 1 pair; medium reviewer 0.188 over 17 pairs; 0 abstained |
| Fable / high | reviewer, 99 axes + 1 blind comparison + spec reviewer, 2 passes + control, 1 full-retrieval run + decider, 15 consultations + debug, 2 charges | Ten comparison reviews accepted all slices. The blind review ranked Sonnet-sliced first, Opus-sliced second, and Opus-full third, but it used nine tests against its cap of six. | Review axes and spec review by reviewer direction; consultations whenever the coordinator is uncertain; the RE2-style full control | reviewer: 0.171 over 12 pairs, 1 abstained |
| Opus / high, xhigh | implementer, latest 10 of 27 current attempts; reviewer, 155 axes; repair, latest 10 of 22 current attempts; spec author and reviewer, latest 6 passes | Fifteen Opus/high ticket assignments completed within their caps, and all seven chunks closed green. Opus/high review required repairs for ten tickets and confirmed every repair. | High for first-round review axes, debug runs, guidance tickets, and foundational Go seams; xhigh for a spec or a ticket that crosses owner seams; repairs at the ticket's effort. | implementation: 0.099 over 67 pairs, 2 abstained; reviewer: 0.161 over 65 pairs, 8 abstained; repair: 0.010 over 10 pairs |
| Opus / medium, low | implementer, latest 10 tickets; repair, latest 10 sessions; orchestrator, 2 builds; earlier reviewer and repair sessions | On `record-evidence` all five medium authors finished in one attempt, and every coordinator probe bit. Reviews still found one-source and per-member test gaps in four of five tickets, and one author edited the pool path with a script. | Medium for a CLI form ticket with named per-member probes, for a mechanical test migration with a named regression probe, and for orchestration. Low for exact tickets; a repair keeps the ticket's declared effort. | 0.129 over 53 pairs; 15 abstained |
| Sonnet / high | orchestrator, 3 landings + reviewer, 14 axes + spec reviewer, 2 passes + implementer, 9 comparison tickets | The six-ticket arm C ranked first in one blind review. Its API-equivalent estimate was $12.90 over 71.2 minutes, versus $12.10 and 69.0 minutes for the selected Opus-sliced arm. | Later review passes after a repair, and orchestration; not fresh ticket authorship | high 0 labeled pairs; xhigh 0.40 over 2 pairs |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing | unknown |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `worktree-seam-reduction` comparison | Opus and Sonnet / high / implementer; Fable / high / reviewer | One blind run ranked Sonnet-sliced first, Opus-sliced second, and Opus-full third. The reviewer selected the second arm, and all three arm gates passed. |
| `tree-targets` TT-C4 round 2 | Fable / high / reviewer | The Spec axis found that a repair through the shared refusal printer doubled a backslash, so the printed repair command named a different label. No row pinned that label. |
| `markdown-block-reader` spec review | Sonnet / high, xhigh; Opus / xhigh; Fable / high, xhigh / spec reviewer | Five reviewers graded one commit with one charge. List costs were Sonnet high $2.14, Fable high $1.98, Opus xhigh $2.44, Fable xhigh $2.50, and Sonnet xhigh $6.31. No reviewer found every confirmed finding, and only Fable high reported a false one. |
| `ft311-recoverable-reset` Coverage axis | Fable / high / reviewer | The axis built the tip, added an ignore rule after the checkpoint, ran the apply, and observed the ignored bytes deleted with no envelope, which no row had decided. |
| `review-evidence-file-pages` RE2 control | claude-fable-5-1 / high / reviewer | The narrow round found 8 findings in 399474 tokens and the full control found 6 in 227918 tokens. They agreed on 3, so neither shape covered the other. |

## Current decisions

- Change routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Keep the `worktree-seam-reduction` comparison descriptive until a comparable repeat supports a routing change.
- Keep the `markdown-block-reader` five-reviewer spec comparison descriptive. Repeat it on later specs before a spec-review routing change.
- Run a first review round on the harness mid binding, unless the reviewer directs Fable for a chunk or a spec. Escalate review to Fable only from round two on, by reviewer rule.
- Give each ticket a fresh Opus author at the spec's declared effort, and give each post-review repair a fresh Opus session at the ticket's effort.
- Use Opus/medium for a mechanical test migration only with a named regression probe that the author runs before and after.
- Make every fix delegate follow `bench-debug`, and accept a debug result that disproves the reported cause.
- Ask a Fable/high read-only consultation when the coordinator is uncertain, and check its claim against the tree.
- Do not route fresh ticket authorship to Sonnet/high until a comparable repeat confirms the latest blind result.
- Charge each author with root conformance, every package it writes, and `cmd/bench` for a public response or embed change.
- Charge each author and repair with the duplicated-facts sweep, its comment sweep, and one recorded red per independent expectation.
- Read a helper before a repair charge names it as the fold target.
- Run the propose-writes preflight before dispatch, and fold fixture paths into the Writes lines; fence expansions are pre-approved.
- Keep the narrow review shape provisional. Pair narrow axes with one full reader in a later trial.
- Probe author claims independently, and rerun every ticket's verification at the final chunk source before the chunk review.
- Record ticket assignments before each dispatch, and record each plan amendment as an identity chunk mapping when chunks do not change.
- Add a repair test's coverage row in the repair plan commit, before the repair dispatch and the confirming round.
- Run tests and probes serially when review axes share a tree.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims.
- Apply the repair allowance to blocking findings, and return exhaustion to the reviewer.
- Do not land an unrelated change on `main` while a chunked build is open; on `worktree-verb-runner` that landing forced a revert, land, and re-land.
- Run the whole conformance package for a guidance ticket, because a Markdown lane misses a test that reads guidance.
- Write chunk, verification, review, and amendment entries with `bench record`, and run the amendment form before the chunk form at each freeze after the first chunk.
- Charge each author of a CLI form with one probe for each member of each refusal step that the form joins.
- Keep a built executable out of a worktree before its checkpoint, because a suite test rewrites its broker manifest. Build a scratch binary with `--manifest-dir` outside the worktree.
- Charge each author to edit pool-path files only with the file tools or through `bench worktree exec`, never with a script on the pool path.
- Stop a consultation that produces no output past its expected window, and reroute it.
- Read the assignment census before landing releases the worktree.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
- Run the doctor rehearsal on the candidate build before a landing that changes a Bench build input.
