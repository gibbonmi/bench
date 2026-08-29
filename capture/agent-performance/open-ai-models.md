# OpenAI model scorecard

Last incorporated landing: `resolved-consumer-surface` (`41b329f9`, 2026-08-28) —
Sol/high as the cross-family reviewer for three initial axes, one repair-scoped
re-review, and one scoped re-check, through `codex exec`, under a Claude
orchestrator with Opus/low writers.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches
dollar break-even near 5x Terra's token use. Per-agent token telemetry is not
currently available.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Terra / medium, high | implementation, latest 10 bounded tickets/repairs | On `bench-shell-follow-on-guard`, medium returned clean bounded seams and biting probes but missed the routing registry, value-taking prefix options, and an ownership fence; high single-sourced projection/prefix policy across two guards and closed the expanded process matrix in one pass | Medium for exact one-seam tickets under coordinator mutation; high when one fact crosses multiple policy consumers |
| Luna / max, medium | implementation, 7 bounded tickets/repairs + reviewer, 4 axis passes on `stable-owner-landing` | Implementation: 2/7 first-pass with all terminal gates green. Review: three initial axes returned 9 raw findings and 8 de-duplicated targets with the citation standard held (one axis refuted four of its own leads by enumeration), and the repair-scoped re-review verified all seven predicates and stayed inside its blocking scope | Low-cost writer for narrow slices under mandatory inspection; standing tier for the three review axes |
| Terra / high | implementation, 3 cross-fence lifecycle charges | 0/3 first-pass on review. On `prospective-artifact-recovery` ticket 01 landed with green rows, probes, and lane, but one branch left an authored binary outside the bundle, an unresolved temporary root leaked the Git registration, and its tests violated the branch-native census | Use when one behavior crosses classifier, planner, and lifecycle consumers; require consumer inventory in the charge |
| Terra / high | semantic review, latest 10 independent axis passes across 4 landings | Found concrete wrapper, descriptor, prefix-value, query-seam, grammar, and one-source gaps with executable probes; final composition axes were clean, while one repair-scoped Standards pass reopened an untouched wrapper-grammar observation that the coordinator correctly kept advisory | Standards, Spec, and Coverage review in separate contexts; pin accepted predicates and prior tip for correction validation |
| Sol / high | implementation and semantic review, 4 lifecycle charges + latest 10 axis passes | On `resolved-consumer-surface` three initial axes returned 28 raw findings that collapsed to 10 repair targets, two of them real behavior gaps in a new blast derivation, and each axis refuted its own leads by enumeration. The repair-scoped re-review found one more real gap the repair had missed (git C-quotes control bytes in patch headers) and stayed inside its scope; the final re-check found one unpinned decoder arm. | Kit-level and security-seam implementation and review; the cross-family reviewer on a large new-package landing |
| Sol / high | orchestration, 5 full spec lifecycles + 1 partial | On `prospective-artifact-recovery` it accepted ticket 01 after a repository-identity probe, then hit the agent thread limit and wrote a handoff that a later session resumed cleanly; its probe did not reach the run-binary branch or the symlinked root that review found later | Complex lifecycle coordination; separate review and composition bases, preflight broker toolchains, and treat merge continuation as its own state transition |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| shell-follow-on guard lifecycle | Terra/medium+high implemented eight ticket slices; Terra/high reviews drove two scoped repair rounds; Sol recovered a destination merge and broker `PATH` refusal before green landing and retirement | delegate, reviewer, and orchestrator | Value-taking option partitions and ownership seams belong in the initial charge; composition bases and broker toolchains need preflight before landing |
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
- A cross-family Sol/high review through `codex exec` is a valid reviewer for a
  Claude-built landing. The recipe closes stdin, sets `-C` to the source, and writes
  the final message with `-o`. A scoped re-review at that tier found a gap the first
  repair missed, so a repair pass gets its own scoped check before it lands.
- Change routing only after two comparable runs or one controlled model
  comparison.
