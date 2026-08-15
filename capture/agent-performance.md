# Agent implementation and orchestration scorecard

This working document guides model and effort routing from observed delivery
quality, cost, and coordination load. It is a current scorecard, not a run log.
Unknown measurements stay unknown; qualitative impressions never become invented
token counts.

## Update contract

- Rewrite rows in place after a completed landing or controlled comparison.
- Aggregate at most the latest 10 comparable assignments per model/effort/role.
- Keep at most 6 routing rows and 5 representative-evidence rows.
- Replace the least decision-relevant evidence when a stronger example arrives.
- Keep the file at or below 120 lines; shorten or remove evidence before adding a
  section.
- Attribute rework to one origin: delegate, spec/ticket, tree/tooling, reviewer,
  or orchestrator. Do not charge upstream omissions to the implementer.
- Compare dollars separately from tokens and latency. The current planning input
  prices Luna tokens at 0.2× Terra, so Luna reaches dollar break-even near 5×
  Terra's token use; per-agent token telemetry is not currently available.

## Measures

| measure | meaning |
| --- | --- |
| first-pass accepted | Done-claim needed no coordinator-authored correction before commit |
| coordinator catches | Material completeness, fence, correctness, or evidence misses |
| repair rounds | Returned write pass after an accepted finding; in-pass feedback is still a catch |
| mutation quality | Required mutations bit and production was restored exactly |
| terminal quality | Focused checks, full gate, and exact-tip review outcome |
| efficiency | Token evidence when available, relative dollar input, and wall-clock churn |

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
  below 5× Terra and wall-clock latency is acceptable.
- Use Terra/high for semantic review and for implementation requiring judgment
  across ownership fences, spec corrections, or several interacting lifecycle
  branches.
- Judge the orchestrator on both terminal correctness and avoidable lifecycle
  churn. A green landing does not erase unnecessary gates, candidate moves, or
  discovery refusals.
- Do not change routing from one run alone. Promote a signal to a general rule only
  after two comparable runs or one controlled model comparison.
