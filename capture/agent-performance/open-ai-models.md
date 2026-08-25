# OpenAI model scorecard

Last incorporated landing: `stable-owner-landing` (`0aecd5fc`, 2026-08-25) —
Sol/high and Terra/medium+high as implementers; Luna/medium as the three review
axes and the repair-scoped re-review, under a Fable orchestrator.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches
dollar break-even near 5x Terra's token use. Per-agent token telemetry is not
currently available.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Terra / medium, high | implementation, 3 bounded tickets + 6 semantic-review axes | Implementer mostly first-pass with biting probes; on `stable-owner-landing` its gate-and-publish ticket needed one coordinator-caught coverage repair (unit tests self-set the env the transport should carry) but caught its own wrong first SOL13 probe; its high-effort repair ticket closed six of seven targets first-pass and correctly stopped the seventh — an accepted predicate that 28 journey tests proved unimplementable — with the counter-proposal the reviewer adopted | Mid-tier default for exact bounded tickets and repairs; stop-and-report on a wrong predicate is holding |
| Luna / max, medium | implementation, 7 bounded tickets/repairs + reviewer, 4 axis passes on `stable-owner-landing` | Implementation: 2/7 first-pass with all terminal gates green. Review: three initial axes returned 9 raw findings and 8 de-duplicated targets with the citation standard held (one axis refuted four of its own leads by enumeration), and the repair-scoped re-review verified all seven predicates and stayed inside its blocking scope | Low-cost writer for narrow slices under mandatory inspection; standing tier for the three review axes |
| Terra / high | implementation, 2 cross-fence lifecycle charges | 0/2 first-pass; found the exact landed predicate but initially preserved a conflicting fixture assumption, and the review-repair charge needed a real-process boundary oracle after a self-referential seam stayed green | Use when one behavior crosses classifier, planner, and lifecycle consumers; require consumer inventory in the charge |
| Terra / high | semantic review, latest 10 independent axis passes across 3 landings | Actionable stale-output, fresh-process, control-byte, derived-count, hostile-consumer, ancestry, resume-recovery, and hostile-matrix findings; correction checks were strongest when charged with exact accepted predicates instead of reopening full discovery | Standards, Spec, and Coverage review; use separate contexts and distinguish initial full review from repair-scoped validation |
| Sol / high | implementation and semantic review, 4 lifecycle charges + latest 10 axis passes | On `stable-owner-landing` the wrapper-broker ticket landed in one build pass plus one fence-extended continuation: it refused to edit an out-of-fence stale test and reported the exact red instead, ran its own biting route-omission probe, and one env-override case exposed that the old route executed a repository decoy; earlier lifecycle and repair charges keep their attributed probes | Kit-level and security-seam implementation and review; the out-of-fence stop discipline is holding |
| Sol / high | orchestration, 4 full spec lifecycles | Preserved the frozen source and unrelated dirty destination bytes, independently killed a distinct discovery input, converged repair-scoped review, and completed landing plus retirement; paid avoidable stale-binary and orphan-roadmap-detail refusals | Complex lifecycle coordination; preflight trusted binaries and retirement companions before paying a full oracle run |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| inherited-toolchain full lifecycle | Terra/medium found five actionable initial gaps and cleared the exact repair predicates; Sol/high preserved unrelated dirty bytes by digest, proved a distinct discovery mutation, and landed the immutable pair | orchestrator | Medium review can be effective with an exact scope; full orchestration should preflight binary trust and retirement companions before expensive gates |
| landed-set tickets | Luna/max delivered four vertical slices with green mutations; three needed coordinator correction at an established renderer, dispatch, or terminal-lifecycle seam | spec/ticket and delegate | Five vertical tickets beat 17 row leaves, but advisory fences must include the existing owner seam before delegation |
| hostile-metadata repair | Terra/high first instrumented only its new planner seam; a real-process `PATH` Git wrapper then exposed older list/resume lease routes and drove the shared guard | delegate | Universal safety claims require enumerated consumers and an oracle above every caller |
| repair-scoped semantic review | Sol/high found real ancestry, resume, control-byte, and 64-hex-token gaps, then kept reopening the original range until a late changelog observation exposed the missing convergence rule | reviewer | Full discovery and correction validation are different modes; only accepted predicates and repair-induced changes block re-review |
| landing-refusal orchestration | Sol/high preserved exact review identity, proved four independent mutations, fixed the review-loop policy, and resumed a green published-but-incomplete landing; it paid avoidable pickup and full-gate churn before breaking the loop | orchestrator | Treat recursive review as a process red and build a tight repro before accepting another unrelated repair target |

## Current decisions

- Use Luna/max for bounded implementation when its expected token use is plausibly
  below 5x Terra and wall-clock latency is acceptable.
- Use Terra/high for semantic review and for implementation requiring judgment
  across ownership fences, spec corrections, or several interacting lifecycle
  branches.
- A universal safety claim names every executable consumer and proves the shared
  external boundary; a repair-owned seam is insufficient evidence.
- Judge the orchestrator on both terminal correctness and avoidable lifecycle
  churn. A green landing does not erase unnecessary gates, candidate moves, or
  discovery refusals.
- Give initial review the full frozen range and closed decisions explicitly. Give
  repair re-review the accepted predicates and prior reviewed tip; unrelated older
  observations are follow-ons, not blockers.
- Change routing only after two comparable runs or one controlled model
  comparison.
