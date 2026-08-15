# OpenAI model scorecard

Last incorporated landing: `gate-run-transaction` (`379035df`).

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches
dollar break-even near 5x Terra's token use. Per-agent token telemetry is not
currently available.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Luna / xhigh | implementation, 1 bounded registry ticket | 0/1 first-pass; 1 completeness catch; mutation and gate evidence strong after repair | Mechanical tickets with explicit retained-test inventory and exact write fence |
| Luna / max | implementation, 3 bounded tickets/repairs | 1/3 first-pass; 2 scope/over-deletion catches; all terminal gates green; mutation evidence strong | Preferred low-cost writer for narrow slices when coordinator inspection is already mandatory |
| Terra / high | semantic review, 12 independent axis passes | 6 actionable findings in the first two cycles; final two cycles clean; no unsupported finding survived verification | Standards, Spec, and Coverage review; use exact frozen base/tip and separate contexts |
| Sol / effort not exposed | orchestration, 1 full spec lifecycle | Strong fence preservation, mutation verification, repair routing, landing recovery, and cleanup; avoidable CLI-discovery and handoff churn | Complex lifecycle coordination; require command-source grammar checks before operational calls |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| verdict-class registry | Luna/xhigh produced the registry and per-validator mutations; omitted retained public cases until challenged | delegate | Cheap tier handles enumerated data shapes, but charges must name retained behavior, not only probes |
| narrow verdict reason | Luna/max reached the right registry-derived behavior; coordinator narrowed an intermediate cross-fence draft | delegate | Max improves difficult local reasoning but does not replace path-fence inspection |
| semantic-review repairs | Luna/max added cancellation, pending-write, and owner-write proofs; one pass over-deleted base prose, the next was clean | delegate and spec/ticket | Batch related public-seam tests, then compare every documentation deletion to the slice base |
| exact semantic review | Terra/high found lifecycle gaps and fixture duplication across two cycles, then returned clean exact-tip verdicts | spec/ticket and reviewer | Higher-cost review paid for itself on edge enumeration and one-source judgment |
| full orchestration | Coordinator caught unreachable mutations, missing retained tests, staged-spec drift, and cleanup state; also tried one unsupported `--help` form | orchestrator | Keep exact-candidate discipline; improve grammar discovery and phase-boundary economy |

## Current decisions

- Use Luna/max for bounded implementation when its expected token use is plausibly
  below 5x Terra and wall-clock latency is acceptable.
- Use Terra/high for semantic review and for implementation requiring judgment
  across ownership fences, spec corrections, or several interacting lifecycle
  branches.
- Judge the orchestrator on both terminal correctness and avoidable lifecycle
  churn. A green landing does not erase unnecessary gates, candidate moves, or
  discovery refusals.
- Change routing only after two comparable runs or one controlled model
  comparison.
