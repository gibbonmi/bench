# OpenAI model scorecard

Last incorporated landing: `worktree-landed-retirement` (`bdc7d978`, 2026-08-16) —
Luna/max and Terra/high as implementers, Terra/high as reviewer, and Sol as orchestrator.

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
| Terra / high | semantic review, latest 10 independent axis passes across 2 landings | Actionable stale-output, fresh-process, control-byte, derived-count, and hostile-consumer findings; the final exact-tip cycle was clean with no unsupported finding surviving verification | Standards, Spec, and Coverage review; use exact frozen base/tip and separate contexts |
| Sol / high | semantic review, 4 axis passes (3 full + 1 targeted) | 23 raw findings, 17 repair targets; every finding carried a citation and an enumeration; 1 finding refuted by a project rule the axis had not read; 3 axes independently converged on the same worst issue | Cross-family review of a refactor's diff; charge must name closed reviewer decisions or they are re-litigated |
| Sol / effort not exposed | orchestration, 2 full spec lifecycles | Preserved the frozen source through four reviews, verified delegate mutations, recovered two moving-destination landing refusals, and completed release/retirement; accepted too many local-seam proofs before forcing the real Git boundary and paid repeated full-gate churn | Complex lifecycle coordination; require consumer inventory and external-boundary proof before freezing review |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| verdict-class registry | Luna/xhigh produced the registry and per-validator mutations; omitted retained public cases until challenged | delegate | Cheap tier handles enumerated data shapes, but charges must name retained behavior, not only probes |
| three-axis refactor review | Sol/high found a mutating invocation advertised as invalid, two surviving duplicate derivations inside the very module built to remove duplication, and correctly separated inherited from diff-owned behavior once asked | reviewer | Cross-family review earns its cost on one-source judgment; it misses project-profile rules unless the charge names them |
| landed-set tickets | Luna/max delivered four vertical slices with green mutations; three needed coordinator correction at an established renderer, dispatch, or terminal-lifecycle seam | spec/ticket and delegate | Five vertical tickets beat 17 row leaves, but advisory fences must include the existing owner seam before delegation |
| hostile-metadata repair | Terra/high first instrumented only its new planner seam; a real-process `PATH` Git wrapper then exposed older list/resume lease routes and drove the shared guard | delegate | Universal safety claims require enumerated consumers and an oracle above every caller |
| exact semantic review | Terra/high found stale-output, fresh-process, control-byte, derived-count, and cross-consumer holes over four cycles, then returned three clean exact-tip verdicts | reviewer | Higher-cost review paid for itself on boundary reachability and one-source judgment |
| full orchestration | Coordinator preserved exact review identity and completed a moving-main landing, but repeated candidate repair and gate cycles before demanding the highest observable boundary | orchestrator | Freeze only after a universal claim's consumer inventory and real-boundary mutation are present |

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
- Give a review charge the closed decisions explicitly. Sol re-derived facts well but
  cited a general skill over a project-profile rule that contradicted it.
- Change routing only after two comparable runs or one controlled model
  comparison.
