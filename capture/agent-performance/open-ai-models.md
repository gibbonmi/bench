# OpenAI model scorecard

Last incorporated landing: `landing-refusal-diagnostics` (`f816e068`, 2026-08-21) —
Sol/high as implementer, reviewer, and orchestrator; Terra/high as reviewer.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches
dollar break-even near 5x Terra's token use. Per-agent token telemetry is not
currently available.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Luna / xhigh | implementation, 1 bounded registry ticket | 0/1 first-pass; 1 completeness catch; mutation and gate evidence strong after repair | Mechanical tickets with explicit retained-test inventory and exact write fence |
| Luna / max | implementation, 7 bounded tickets/repairs across 2 landings | 2/7 first-pass; 5 fence, lifecycle-boundary, or scope catches; all terminal gates green and required mutations bit | Preferred low-cost writer for narrow vertical slices when coordinator inspection is already mandatory |
| Terra / high | implementation, 2 cross-fence lifecycle charges | 0/2 first-pass; found the exact landed predicate but initially preserved a conflicting fixture assumption, and the review-repair charge needed a real-process boundary oracle after a self-referential seam stayed green | Use when one behavior crosses classifier, planner, and lifecycle consumers; require consumer inventory in the charge |
| Terra / high | semantic review, latest 10 independent axis passes across 3 landings | Actionable stale-output, fresh-process, control-byte, derived-count, hostile-consumer, ancestry, resume-recovery, and hostile-matrix findings; correction checks were strongest when charged with exact accepted predicates instead of reopening full discovery | Standards, Spec, and Coverage review; use separate contexts and distinguish initial full review from repair-scoped validation |
| Sol / high | implementation and semantic review, 1 lifecycle + latest 10 axis passes | Repair delegates supplied real reds and restored production exactly; review found genuine ancestry, resume, hostile-matrix, and mutation-strength gaps, but repeated full-range discovery surfaced an unrelated changelog miss after the accepted corrections were already clean | Kit-level repair and review; initial discovery stays full-range, while repair validation must use accepted predicates and the repair delta |
| Sol / high | orchestration, 3 full spec lifecycles | Preserved the frozen source, independently killed resume/token/path mutations, diagnosed native-Go PATH and trusted-binary refusals, and completed a published-incomplete resume plus retirement; failed to stop recursive full re-review until the user challenged it, then reproduced and fixed the canonical workflow defect | Complex lifecycle coordination; treat non-convergence as a workflow bug immediately instead of paying another review/gate cycle |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| three-axis refactor review | Sol/high found a mutating invocation advertised as invalid, two surviving duplicate derivations inside the very module built to remove duplication, and correctly separated inherited from diff-owned behavior once asked | reviewer | Cross-family review earns its cost on one-source judgment; it misses project-profile rules unless the charge names them |
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
