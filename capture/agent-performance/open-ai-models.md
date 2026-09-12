# OpenAI model scorecard

Last incorporated landing: `retained-implementation-workflow` (`d11d9b7e2dbe12803f05de1f5e9a6a40c4766388`, 2026-09-11).
Sol/high implemented three chunks and repaired all accepted findings in one retained session.
Nine Astra/high axes reviewed the three chunks once each.
The closeout retired the implemented spec at `9d1c4a81587f61650a432ac4df02a51a7f51ff20`.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use.
The retained-workflow session used 95,333,871 input tokens, including 93,844,992 cached tokens, and 265,635 output tokens across Sol and Astra.
At the 2026-09-11 standard API prices, these tokens cost $56.36.
An internal auto-review model has no public API price, so this total excludes its tokens.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Astra / high, medium | reviewer, 9 high axes on 1 landing + implementer, 5 medium tickets on 1 landing | On `retained-implementation-workflow`, nine Astra/high axes found eight accepted targets across three chunks. The reviews found a permanent-gate defect, stale handoff prose, a retro parser mismatch, and missing chunk identities. | High for three independent review axes over a retained Sol build; medium as a retained author under an independent higher-tier review |
| Terra / medium, high | implementation, latest 10 bounded tickets and repairs | Terra/high implemented FT311 ticket 4 with live evidence open. Review found the missing triage input contract, and repair added its five inputs. | Medium for exact one-seam tickets under coordinator mutation. High when one fact crosses multiple policy consumers. |
| Terra / high, medium, low | research, 2 charges; semantic review, latest 10 independent axis passes | Both Terra/high research returns needed coordinator corrections. One missed Codex contracts; the other missed Claude telemetry and conflated agent ancestry with nested tool calls. | Standards, Spec, and Coverage review in separate contexts; verify research claims against primary sources |
| Luna / max, medium | implementation, 7 bounded tickets/repairs + reviewer, 4 axis passes on `stable-owner-landing` | Implementation: 2/7 first-pass with all terminal gates green. Review: three initial axes returned 9 raw findings and 8 de-duplicated targets with the citation standard held, and the repair-scoped re-review verified all seven predicates. | Low-cost writer for narrow slices under mandatory inspection |
| Sol / low | implementation, bounded tickets and repairs on 1 landing | The delegates returned focused tests and mutation probes. The retirement repair reproduced the FT94 ledger red, changed one owner, proved that restoring the old value made the test red, and restored green. | Exact ticket seams and small repairs under coordinator verification |
| Sol / high | orchestrator, 1 benchmark arm + retained implementer, 3 tickets on 1 landing + semantic review | On `retained-implementation-workflow`, Sol/high implemented three chunks, repaired eight accepted review targets, and passed the final gate. Final reconciliation caught three missed consumers and one late fence addition. | Retained implementation of connected kit-policy chunks, review, and orchestration of a retained lower-tier author |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| FT311 recoverable-reset candidate | Astra/medium built five tickets green under Sol/high; a Fable/high round then found three behavior defects the Terra/medium review missed. | delegate and reviewer | A same-provider medium review of a medium build needs one higher-tier round before a landing. |
| FT311 ticket 2 repair | Sol/high unified charge and proposal orchestration and completed every DP10-DP14 partition. Exact-baseline tests and two restored mutations bit. | delegate | Use Sol/high when bounded repair crosses command retry, authority, and graph semantics. |
| FT311 ticket 4 | Terra/high produced build and triage guidance. Review found that active guidance omitted the required triage inputs. | delegate | Trace each spec input into final phase text before returning a guidance ticket. |
| FT311 linked dogfood | Terra used `gpt-5.6-terra`; the shift exited 0 with one commit, two iterations, and no recovery. | orchestrator | Use a linked repository when the kit repository cannot exercise its installed surface. |
| Retained implementation workflow | The Sol/high author completed three chunks, and nine Astra/high axes found eight accepted targets. | orchestrator and reviewer | Use Astra/high after each Sol/high chunk, and keep the Sol author for all accepted repairs. |

## Current decisions

- Use Terra/medium for a ticket-sized write charge under a coordinator mutation probe.
- Use Astra/medium as a retained author when the reviewer names it and a higher tier reviews the result.
- Use Astra/high for the three independent review axes after each Sol/high chunk.
- Pair a Terra/medium review of a Terra/medium build with one Fable/high or Sol/high round before landing.
- Use Terra/high for final independent review when closeout changes kit guidance and conformance together.
- Use Sol/low for an exact ticket seam or a small repair when the coordinator runs a distinct mutation probe.
- Retain Sol/high across connected kit-policy chunks and keep it for all accepted same-session repairs.
- Settle a disputed finding against the frozen candidate with an exact reproduction.
- If a landing gate is red without a green baseline, run the fresh baseline before diagnosis.
- Give a Codex Coverage axis a writable worktree when it must run a temporary probe.
- Change routing only after two comparable runs or one controlled model comparison.
