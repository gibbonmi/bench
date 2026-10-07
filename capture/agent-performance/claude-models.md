# Claude model scorecard

Last incorporated landing: `ft290-test-projection` (`e558624eb20b848ed4c2b6680d04fa8cc0f17557`, 2026-10-07).

The `ft290-test-projection` build ran five chunks and ten tickets on Claude Code. An Opus/medium orchestrator dispatched ten fresh Opus/medium authors and fourteen fresh Opus/medium repair sessions. By reviewer direction, Sonnet/high ran all review axes and Fable/high ran two consultations. Provider usage and charges remain unknown.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 38 landings + implementer, 10 charges + reviewer, 21 axes | On `worktree-seam-reduction` the high orchestrator ran seven chunks, preserved the three-arm blind, and probed each accepted repair. All chunk checkpoints and the original landing gate passed. | Coordination of a delegated build and adversarial spec review; it implements only when the reviewer names it | orchestrator 0.36 over 1 pair; medium reviewer 0.188 over 17 pairs; 0 abstained |
| Fable / high | reviewer, 124 axes + 1 blind comparison + spec reviewer, 4 passes + control, 1 full-retrieval run + decider, 17 consultations + debug, 2 charges | On `ft290-test-projection` Fable/high ran two read-only consultations. One found a record-verb route that avoided a hand edit, and one proved a kit conflict between the completion gate and the review preflight. | Review axes and spec review by reviewer direction; consultations whenever the coordinator is uncertain; the RE2-style full control | reviewer: 0.168 over 32 pairs, 1 abstained |
| Opus / high, xhigh | implementer, latest 10 of 33 current attempts; reviewer, 155 axes; repair, latest 10 of 27 current attempts; spec author, slicer, and reviewer, latest 10 passes | On `ft290-test-projection` one Opus/high light-path author fixed the review preflight for a folded main. It stopped correctly at its fence twice, and a coordinator omission probe bit; its first landing failed a prose budget that its lane did not run. | High for first-round review axes, debug runs, guidance tickets, and foundational Go seams; xhigh for a spec or a ticket that crosses owner seams; repairs at the ticket's effort. | implementation: 0.098 over 68 pairs, 2 abstained, and 0.010 over 1 new pair; reviewer: 0.163 over 70 pairs, 9 abstained; repair: 0.010 over 11 pairs |
| Opus / medium, low | implementer, latest 10 tickets; repair, latest 10 sessions; orchestrator, 4 builds; earlier reviewer and repair sessions | On `ft290-test-projection` ten medium authors each committed first-pass, and every named probe bit; reviews still drove repairs on all chunks, mostly for untested edges. The medium orchestrator repeated two recorded scorecard lessons and charged one probe that could not run. | Medium for a CLI form ticket with named per-member probes, for a mechanical test migration with a named regression probe, and for orchestration. Low for exact tickets; a repair keeps the ticket's declared effort. | 0.129 over 66 pairs; 18 abstained |
| Sonnet / high | orchestrator, 3 landings + reviewer, 77 axes + spec reviewer, 5 passes + researcher, 3 charges + implementer, 9 comparison tickets | On `ft290-test-projection` Sonnet/high ran every axis of five chunk reviews and their confirming rounds. Its Coverage axes found real gaps with silent probes, including a self-import cause and an inherited selector defect, and no finding was false. | First-round and confirming review axes by reviewer direction, and orchestration; not fresh ticket authorship | high 0.145 over 2 pairs; xhigh 0.40 over 2 pairs |
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
- Keep the `worktree-seam-reduction` comparison and the `markdown-block-reader` five-reviewer comparison descriptive until a comparable repeat.
- Run a first review round on the harness mid binding, unless the reviewer directs another line. Sonnet/high served every axis of `ft290-test-projection` by reviewer direction.
- Give each ticket a fresh Opus author at the spec's declared effort, and give each post-review repair a fresh Opus session at the ticket's effort.
- Make every fix delegate follow `bench-debug`, and accept a debug result that disproves the reported cause.
- Ask a Fable/high read-only consultation when the coordinator is uncertain, and check its claim against the tree.
- Prove a review claim about Git behavior in a scratch tree before a repair dispatch.
- Do not route fresh ticket authorship to Sonnet/high until a comparable repeat confirms the latest blind result.
- Charge each author with root conformance, every package it writes, and `cmd/bench` for a public response or embed change.
- Charge each author and repair with the duplicated-facts sweep, its comment sweep, and one recorded red per independent expectation.
- Charge each repair of a shared helper with every production caller package, and require its mutation replay in each.
- Run each named-check row through `bench worktree build` and then `./dist/bench test --check`, because the installed binary grades its own source.
- Write each probe in a charge in a form that `bench probe` can express; a file omission empties a file and does not delete it.
- Put these orchestration rules into the shared charge file before the first dispatch: amendment before chunk record, mutated-run probe exit, and the named-check build.
- Read a helper before a repair charge names it as the fold target.
- Run the propose-writes preflight before dispatch, and fold fixture paths into the Writes lines.
- Treat `--auto-approve` as approval for fence expansions, spec and ticket expansions, and repair rounds, and fix a found defect at once instead of deferring it.
- When `main` moves during a build, fold `main` into the source before the completion landing. The fold joins the last chunk's review delta, and the landing base stays the pre-chunk `main` tip.
- Before a completion landing, sweep for references to each roadmap row that the closure deletes. The completion checkpoint grades the tree before the closure.
- Keep the narrow review shape provisional. Pair narrow axes with one full reader in a later trial.
- Probe author claims independently, and rerun every ticket's verification at the final chunk source before the chunk review.
- Record ticket assignments before each dispatch, and give each verification of a version 2 chunk an ID that is unique in that chunk.
- Add a repair test's coverage row in the repair plan commit, before the repair dispatch and the confirming round.
- Run tests and probes serially when review axes or authors share a tree.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims.
- Apply the repair allowance to blocking findings that change source. A record-only or comment-only change consumes no repair cycle.
- Correct spec prose inside the chunk that needs it, because the checkpoint needs a passing Spec result.
- Name every anchor needle in the plan probe list with its expected verdict.
- Write each fact in a delegate charge from git output, not from memory.
- Run the whole conformance package for a guidance ticket, because a Markdown lane misses a test that reads guidance, such as a prose budget.
- Charge each author of a CLI form with one probe for each member of each refusal step that the form joins.
- Charge each author to edit pool-path files only with the file tools or through `bench worktree exec`.
- Stop a consultation that produces no output past its expected window, and reroute it.
- Read the assignment census before landing releases the worktree.
- Bind a deliverable to its outcome through a landed plan before `bench commitment start`.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
- Run the doctor rehearsal on the candidate build before a landing that changes a Bench build input.
