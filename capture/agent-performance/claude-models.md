# Claude model scorecard

Last incorporated landing: `ticket-checkpoint-floor` (`d0b62266587865469dff55076c9e3d6432046652`, 2026-10-09) and `known-issue-table` (`9b99dc0c8b44eb0640aa31faa53cef63d7a3ec5a`, 2026-10-08).

The `ticket-checkpoint-floor` build ran five chunks on Claude Code. An Opus/high orchestrator dispatched five fresh Opus/high authors and seven repair sessions. By reviewer direction, first rounds ran on Fable/high, then Opus/high, then Fable/medium, and confirming rounds ran on Sonnet/high and then Fable/medium. Haiku and Sonnet ran read-only research, and one Opus/medium author landed the `known-issue-table` light path. Provider usage and charges remain unknown.

The `refusal-route-registry` build ran seven chunks and thirteen tickets on Claude Code. An Opus/high orchestrator dispatched 13 fresh Opus/high authors, one successor, 18 repair sessions, and 5 verification-only transfers. By reviewer direction, Sonnet/high ran 21 first-round axes and 34 confirming axes, and Fable/high ran one consultation. Provider usage and charges remain unknown.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 38 landings + implementer, 10 charges + reviewer, 21 axes | On `ticket-checkpoint-floor` Fable/medium ran first and confirming axes for three chunks. Each axis took 1 to 4.5 minutes, and its Coverage axes found the two false passes that became rows TF73 and TF74. | Coordination of a delegated build and adversarial spec review; it implements only when the reviewer names it | orchestrator 0.36 over 1 pair; medium reviewer 0.188 over 17 pairs; 0 abstained |
| Fable / high | reviewer, 124 axes + 1 blind comparison + spec reviewer, 4 passes + control, 1 full-retrieval run + decider, 17 consultations + debug, 2 charges | On `ticket-checkpoint-floor` Fable/high ran the first TF-C1 round. Its Coverage axis proved an untested kit clause with a silent probe, and its Standards axis found three duplicated facts. | Review axes and spec review by reviewer direction; consultations whenever the coordinator is uncertain; the RE2-style full control | reviewer: 0.168 over 32 pairs, 1 abstained |
| Opus / high, xhigh | implementer, latest 10 of 49 current attempts; reviewer, 161 axes; repair, latest 10 of 46 current attempts; orchestrator, 1 build; spec author, slicer, and reviewer, latest 10 passes | On `ticket-checkpoint-floor` five high authors committed inside their fences or stopped at them, and every named probe bit. Two authors stopped for fences that the spec omitted. Charges that named the existing owners removed the blocking duplication from later first rounds. | High for first-round review axes, debug runs, guidance tickets, and foundational Go seams such as an oracle route; xhigh for a spec or a ticket that crosses owner seams; repairs at the ticket's effort. | implementation: 0.098 over 68 pairs, 2 abstained, 0.010 over 1 pair, and 0.028 over 5 new pairs; reviewer: 0.162 over 74 pairs, 9 abstained; repair: 0.010 over 11 pairs |
| Opus / medium, low | implementer, latest 10 tickets; repair, latest 10 sessions; orchestrator, 5 builds; earlier reviewer and repair sessions | On `known-issue-table` one medium author landed a light-path ticket with a four-row table and a one-line charge rule. It reported one stale reference outside its fence before it edited. | Medium for a CLI form ticket with named per-member probes, for a mechanical test migration with a named regression probe, and for orchestration. Low for exact tickets; a repair keeps the ticket's declared effort. | 0.129 over 67 pairs; 18 abstained |
| Sonnet / high | orchestrator, 3 landings + reviewer, 132 axes + spec reviewer, 5 passes + researcher, 3 charges + implementer, 9 comparison tickets | On `ticket-checkpoint-floor` Sonnet/high ran the TF-C1 confirming rounds and four staleness slices. It found one fixture duplication, and its staleness verdicts were all correct. | First-round and confirming review axes by reviewer direction, and orchestration; not fresh ticket authorship | high 0.145 over 2 pairs; xhigh 0.40 over 2 pairs |
| Sonnet / high (research and polish) | researcher, 6 paired charges + polish, 1 session | On `stdlib-reimplementation-check` Sonnet made no wrong claim in six research slices and found every key fact under the paired-reader rules. Its high polish of a Haiku draft closed nine audit items in one pass and found a hole in the live-tree meta check. | Default research reader, paired with Haiku; polish of a cheap draft when the reviewer names it | research 6 of 6 slices clean; polish 1 pass |
| Haiku / default–xhigh | researcher, 12 charges + implementer, 2 comparison tickets + auditor, 3 charges | On `ticket-checkpoint-floor` Haiku ranked 14 roadmap rows for time savings. Its top two citations held, and it named what it did not read. | Research reader, paired with Sonnet; staleness audits; not ticket authorship | research 10 of 12 slices clean; authorship 0 landings |
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

- Read this scorecard before each line declaration, because its decisions carry the orchestration rules that a build repeats.
- Change routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Keep the `worktree-seam-reduction` comparison and the `markdown-block-reader` five-reviewer comparison descriptive until a comparable repeat.
- Run review axes on the line that the reviewer directs. Fable/medium served the later rounds of `ticket-checkpoint-floor` by reviewer direction, and it was faster than Opus/high at equal finding quality.
- When a delegate nears 300k tokens, let it finish its ticket, and send later work to a fresh recorded session.
- Name the existing owners, seams, and fixture helpers in each author charge, together with the recurring review lessons.
- Before a ticket adds a printed route to a shared funnel or a shared face, name every caller and every printing verb in its charge.
- Give each ticket a fresh Opus author at the spec's declared effort, and give each post-review repair a fresh Opus session at the ticket's effort.
- Make every fix delegate follow `bench-debug`, and accept a debug result that disproves the reported cause.
- Ask a Fable/high read-only consultation when the coordinator is uncertain, and check its claim against the tree.
- Do not route fresh ticket authorship to Sonnet/high until a comparable repeat confirms the latest blind result.
- Charge each author with root conformance, every package it writes, and `cmd/bench` for a public response or embed change. The `check-floor` row grades that floor.
- Charge each author and repair with the duplicated-facts sweep, its comment sweep, and one recorded red per independent expectation.
- When an active author reports a copy outside its fence, widen the fence before the freeze. The same author then folds the copy before review.
- Run `bench worktree build` after each Go commit in a build worktree, because a stale binary turns `binary-seal` red.
- Write each probe so that its mutant compiles and changes behavior; a mutant with an unused import or an equal commit name proves nothing.
- Run `bench record amendment` after each plan commit and after each ticket commit that rewrites seam cells, before `bench record chunk`.
- Apply a plan-text correction before the chunk freeze, or freeze the chunk again before any record at the old source.
- Record the mutated-run probe exit as 1 from the failing test run; `docs/known-issues.md` holds the workaround.
- With reviewer approval, let the ticket author record the amendment, freeze its chunk, and record its verification in one pass before it returns.
- Treat `--auto-approve` as approval for fence expansions, spec and ticket expansions, and repair rounds, and fix a found defect at once instead of deferring it.
- When `main` moves during a build, fold `main` into the source before the completion landing. Fold it earlier when a later ticket writes a path that `main` changed. The fold joins the last chunk's review delta.
- Commit every record before the complete checkpoint, because the checkpoint refuses a dirty checkout and grades the published tree.
- Run tests and probes serially when review axes or authors share a tree.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims.
- Apply the repair allowance to blocking findings that change source. A record-only or plan-text change consumes no repair cycle.
- Run the whole conformance package for a guidance ticket, because a Markdown lane misses a test that reads guidance.
- Read the assignment census before landing releases the worktree.
- Bind a deliverable to its outcome through a landed plan before `bench commitment start`, and run the start through `bench worktree exec` in the owning worktree.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
- Run the doctor rehearsal on the candidate build before a landing that changes a Bench build input.
