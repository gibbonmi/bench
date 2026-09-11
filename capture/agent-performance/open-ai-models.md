# OpenAI model scorecard

Last incorporated landing: `ft311-recoverable-reset` (`5201c9fc21f0dd85563edb704f401294816edeb7`, 2026-09-11). The landed source was the frozen benchmark candidate that Sol/high coordinated and one retained Astra/medium author built, with Terra/medium review axes.

The candidate passed every per-ticket probe and its final gate. It carried three behavior defects that a later Fable/high round found. An ignored file under a tracked directory was overwritten, an ignore-rule change deleted bytes, and a staged-only change did not stale the plan.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Astra / medium | implementer, 5 tickets on 1 landing + reviewer, 3 axes + 2 re-reviews on 1 landing | On `ft311-recoverable-reset` one retained author built all five tickets with biting probes and one same-model context replacement after two auto-review refusals. Its fingerprint omitted the raw index, and its move order deleted ignored bytes under a changed rule, which the candidate's Terra/medium review did not find. | Retained author for a ticket sequence under an independent review at a higher tier; give the Coverage axis a writable worktree |
| Terra / medium, high | implementation, latest 10 bounded tickets and repairs | Terra/high implemented FT311 ticket 4 with live evidence open. Review found the missing triage input contract, and repair added its five inputs. | Medium for exact one-seam tickets under coordinator mutation. High when one fact crosses multiple policy consumers. |
| Terra / high, medium, low | semantic review and closeout, latest 10 independent axis passes | On `ft311-recoverable-reset` three Terra/medium axes over the Astra candidate returned 5 raw findings and 3 targets, and they named the undecided collision partition. They did not find the staged-index gap, the ignore-rule drift, or the below-path collision that a Fable/high round found on the same tip. | Standards, Spec, and Coverage review in separate contexts; pair a medium review of a same-provider build with one higher-tier round before a landing |
| Luna / max, medium | implementation, 7 bounded tickets/repairs + reviewer, 4 axis passes on `stable-owner-landing` | Implementation: 2/7 first-pass with all terminal gates green. Review: three initial axes returned 9 raw findings and 8 de-duplicated targets with the citation standard held, and the repair-scoped re-review verified all seven predicates. | Low-cost writer for narrow slices under mandatory inspection |
| Sol / low | implementation, bounded tickets and repairs on 1 landing | The delegates returned focused tests and mutation probes. The retirement repair reproduced the FT94 ledger red, changed one owner, proved that restoring the old value made the test red, and restored green. | Exact ticket seams and small repairs under coordinator verification |
| Sol / high | orchestrator, 1 benchmark arm + implementation and semantic review, latest comparable lifecycle charges | On `ft311-recoverable-reset` it coordinated the Astra author through five tickets and one context replacement, ran per-ticket probes, and froze a green candidate with one undecided partition named. Sol/high implemented FT311 tickets 1 and 2 earlier, with one repair on ticket 2. | Kit and security-seam implementation, review, and orchestration of a retained lower-tier author |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| FT311 recoverable-reset candidate | Astra/medium built five tickets green under Sol/high; a Fable/high round then found three behavior defects the Terra/medium review missed. | delegate and reviewer | A same-provider medium review of a medium build needs one higher-tier round before a landing. |
| FT311 ticket 2 repair | Sol/high unified charge and proposal orchestration and completed every DP10-DP14 partition. Exact-baseline tests and two restored mutations bit. | delegate | Use Sol/high when bounded repair crosses command retry, authority, and graph semantics. |
| FT311 ticket 4 | Terra/high produced build and triage guidance. Review found that active guidance omitted the required triage inputs. | delegate | Trace each spec input into final phase text before returning a guidance ticket. |
| FT311 linked dogfood | Terra used `gpt-5.6-terra`; the shift exited 0 with one commit, two iterations, and no recovery. | orchestrator | Use a linked repository when the kit repository cannot exercise its installed surface. |
| landing-refusal orchestration | Sol/high preserved exact review identity, proved four independent mutations, fixed the review-loop policy, and resumed a green published-but-incomplete landing; it paid avoidable pickup and full-gate churn before breaking the loop | orchestrator | Treat recursive review as a process red and build a tight repro before accepting another unrelated repair target |

## Current decisions

- Use Terra/medium for ticket-sized write charges under coordinator mutation probes.
- Use Astra/medium as a retained author for a ticket sequence when the reviewer names it, under an independent review at a higher tier.
- Pair a Terra/medium review of a same-provider build with one Fable/high or Sol/high round before the landing.
- Use Terra/high for final independent review when closeout changes kit guidance and conformance together.
- Use Sol/low for exact ticket seams and small repairs when a coordinator can run a distinct mutation probe.
- Use Sol/high for bounded repair crossing command retry, authority, and dependency-graph behavior, and for orchestration of a retained lower-tier author.
- Settle disputed findings against the frozen candidate with an exact reproduction.
- On a landing gate red without a green baseline, run the fresh baseline before diagnosis.
- Give a Codex Coverage axis a writable isolated worktree when it must run throwaway probes.
- Change routing only after two comparable runs or one controlled model comparison.
