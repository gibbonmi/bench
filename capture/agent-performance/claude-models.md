# Claude model scorecard

Last incorporated landing: `session-context-measurement` (`e577ca53cc9ed20dca4860819039b0b9330063a7`, 2026-09-12). Fable coordinated a two-chunk spec build with one Opus/high retained author, three Opus/medium review axes in separate venues, and one Sonnet/xhigh debug delegate. Two fast-track fixes landed on main during the build.

The first review round returned 12 findings on chunk 1 and 10 on chunk 2, with seven and eight repair targets. Every reaffirmation converged in two rounds. The coordinator's independent probes bit at two sites the author never probed.

Seventy-five completed landings are recorded. Routing follows the harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 35 landings + implementer, 10 charges + reviewer, 6 axes | On `session-context-measurement` it added the missing completion plan before the first checkpoint, probed every return at a fresh site, and disposed two review findings on the rule text. It paid one landing refusal on the base after a mid-build main merge, and it scripted eight raw writes into the pool path. | Coordination of a parallel build and adversarial spec review; it implements only when the reviewer names it |
| Fable / high | reviewer, 3 axes on 1 landing | On `ft311-recoverable-reset` the three axes found the below-path collision, the ignore-rule drift that deleted bytes, the hidden index flags, and the primary-side checkpoint resolution, each with an executed probe. | Review axes over a candidate another provider built, when the reviewer names the tier |
| Opus / high | implementer, latest 10 guidance and Go-seam charges | On `session-context-measurement` one retained author landed 21 rows red-first with five named probes, a 20-case evidence report, and three repair rounds. It claimed per-dimension absence tests that did not exist and labelled one capture with the wrong command. | High for the retained author of a spec build, process lifecycle, anchored guidance, and foundational Go seams |
| Opus / medium, low | implementer, orchestrator, and reviewer combined | On `session-context-measurement` three medium axes found an untested producer shape through a silent probe, an undecided oversized-line partition, and a missing-key zero, then re-measured every disputed byte count themselves. One axis withdrew a fence finding after it read the rule the coordinator cited. | Medium for the three review axes, reaffirmations, gates, conformance, repair, and triage. Low for exact tickets |
| Sonnet / high | orchestrator, 3 landings | On `roadmap-light-path-fixes-2` it ran ten ticket charges and two review rounds across two shared worktrees, caught a dirty integration worktree before the next commit, and routed two acceptance shortfalls to the reviewer. | Continues to hold at high effort; compare again after a fourth orchestrated build |
| Sonnet / low–xhigh | implementer, latest 10 of 80 charges + debug, 1 landing | On `learnings-inbox-shape` one xhigh debug charge built the red loop first, falsified the two-parser hypothesis, fixed the writer with a round-trip test, then added the prose gate that refused the coordinator's next entry. It left a reproduction copy of the ignored inbox in its worktree. | Xhigh for a bounded debug path with a named symptom; low for an exact prose or spec repair when the reviewer names it |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `session-context-measurement` Coverage axis | Opus / medium / reviewer | The axis probed the string-form result path with `bench probe`, saw 48 tests stay silent, and named the fixture that turned the mutation red. |
| `session-context-measurement` Spec axis | Opus / medium / reviewer | The axis re-ran four measured cases through the real command and separated tree drift from a misreported byte count. |
| `learnings-inbox-shape` debug | Sonnet / xhigh / implementer | The delegate reproduced the symptom with a constructed entry, showed both readers share one parser, and fixed the one writer. |
| `ft311-recoverable-reset` Coverage axis | Fable / high / reviewer | The axis built the tip, added an ignore rule after the checkpoint, ran the apply, and observed the ignored bytes deleted with no envelope. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The delegate found that one environment change crossed three packages and stopped at the fence. |

## Current decisions

- Route changes only after two comparable runs or one controlled model comparison.
- The top tier implements only when the reviewer names it for the run.
- Use Opus/high as the retained author of a spec build when the reviewer names delegated authorship.
- Use Opus/medium for the three independent review axes and for every reaffirmation.
- Use Fable/high for the review axes over a candidate another provider built, when the reviewer names the tier.
- Use Sonnet/xhigh for a bounded debug path with a named symptom and a fenced owner.
- Keep Standards, Spec, and Coverage in separate native contexts and retain each return independently.
- Write and commit the review pickup before the repair charge goes out.
- Merge the latest main tip into the source before the last chunk's review, and pass that tip as the landing base.
- Probe every done-claim at a distinct site and mutation kind, and prove every new test red under a named mutation before its commit.
- Route undecided behavior partitions to the reviewer in one question with a recommendation each.
- Read the assignment census before the landing releases a worktree, and script no write into the pool path.
- Finish a hand-resolved merge with `git merge --continue`, never with `bench commit`.
- Create no sibling worktree while another session is inside a landing tail.
