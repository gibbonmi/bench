# OpenAI model scorecard

Last incorporated landing: `git-admin-readers` (`5589e73a`, 2026-09-05) — Fable
(Claude Code) orchestrated. Astra/medium, a reviewer-named model outside the
tier binding, ran the three review axes and two re-reviews through `codex exec`
inside `bench worktree exec` with an empty quoted heredoc.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches
dollar break-even near 5x Terra's token use. Per-agent token telemetry is not
currently available.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Terra / medium, high | implementation, latest 10 bounded tickets/repairs | On `landing-refusal-standard`, medium landed 11 tickets and one 5-target repair pass with biting probes and honest fence-limit reports; 7 of 12 charges were first-pass accepted, and the misses were sampled assertions and unobserved helpers. | Medium for exact one-seam tickets under coordinator mutation; high when one fact crosses multiple policy consumers |
| Luna / max, medium | implementation, 7 bounded tickets/repairs + reviewer, 4 axis passes on `stable-owner-landing` | Implementation: 2/7 first-pass with all terminal gates green. Review: three initial axes returned 9 raw findings and 8 de-duplicated targets with the citation standard held (one axis refuted four of its own leads by enumeration), and the repair-scoped re-review verified all seven predicates and stayed inside its blocking scope | Low-cost writer for narrow slices under mandatory inspection; standing tier for the three review axes |
| Terra / high, low | semantic review, latest 10 independent axis passes across 7 landings | On `kit-guidance-fold`, one high-effort falsification pass over eight guidance files returned three cited findings; two were accepted, and one of them was the term collision on "disposition" that the three Claude axes missed. It grepped the registry before it claimed an uncovered sentence. | Standards, Spec, and Coverage review in separate contexts, and the standing cross-harness falsification pass on a kit-guidance diff |
| Sol / low | implementation, bounded tickets and repairs on 1 landing | The delegates returned focused tests and mutation probes. The retirement repair reproduced the FT94 ledger red, changed one owner, proved that restoring the old value made the test red, and restored green. | Exact ticket seams and small repairs under coordinator verification |
| Sol / high | implementation and semantic review, 4 lifecycle charges + latest 10 axis passes | On `resolved-consumer-surface`, three initial axes collapsed 28 raw findings to 10 repair targets. A repair-scoped re-review found one missed control-byte gap, and the final re-check found one unpinned decoder arm. | Kit-level and security-seam implementation and review; the cross-family reviewer on a large new-package landing |
| Astra / medium | semantic review, 3 axes + 2 re-reviews on 1 landing | On `git-admin-readers` the three axes returned 13 raw findings and 6 repair targets; the Coverage axis probed a symlink-plus-`..` root and a symlinked temp parent and observed both breaks, and the Spec axis found the check walked only two directories. Both re-reviews closed every predicate and stayed inside the blocking scope. A read-only sandbox refused the Go cache, so the blast table and `bench test --check` ran outside it. | Reviewer-named implementation review when the tier binding is not the reviewer's choice; give the Coverage axis a writable worktree |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| roadmap light-path lifecycle | Sol/low implemented bounded tickets and the retirement repair; Terra/high found the LF2 and retrospective safety gaps; Codex held the frozen source and rejected a false LF2 claim | spec/ticket, delegate, reviewer, and orchestrator | Cross-consumer inputs and race-test reachability belong in the initial charge; exact frozen repros settle disputed findings |
| landing-refusal build | Terra/medium landed a 24-file, 11-ticket registry build with one repair commit; two delegates reported fence limits instead of routing around them, and the coordinator ran one blocked probe itself | delegate and spec/ticket | A charge that names the shared constructor lets eight serialized tickets compose without a merge conflict |
| hostile-metadata repair | Terra/high first instrumented only its new planner seam; a real-process `PATH` Git wrapper then exposed older list/resume lease routes and drove the shared guard | delegate | Universal safety claims require enumerated consumers and an oracle above every caller |
| repair-scoped semantic review | Sol/high found real ancestry, resume, control-byte, and 64-hex-token gaps, then kept reopening the original range until a late changelog observation exposed the missing convergence rule | reviewer | Full discovery and correction validation are different modes; only accepted predicates and repair-induced changes block re-review |
| landing-refusal orchestration | Sol/high preserved exact review identity, proved four independent mutations, fixed the review-loop policy, and resumed a green published-but-incomplete landing; it paid avoidable pickup and full-gate churn before breaking the loop | orchestrator | Treat recursive review as a process red and build a tight repro before accepting another unrelated repair target |

## Current decisions

- Use Terra/medium for ticket-sized write charges under coordinator mutation
  probes; it carried an 11-ticket build with 7 of 12 first-pass acceptances.
- Use Terra/low for the three review axes; its findings held the citation
  standard, and the coordinator settles a predicted red against the oracle.
- Use Terra/high for the standing cross-harness falsification pass on a kit-guidance
  diff; run it through `bench worktree exec` with an empty quoted heredoc.
- Use Sol/low for exact ticket seams and small repairs when a coordinator can
  rerun a distinct mutation probe.
- Settle a disputed finding against the frozen candidate with an exact repro.
  Do not repair a failure that the candidate does not produce.
- On a landing gate red with no green baseline, run the fresh baseline before
  diagnosis; one transient red cost one avoidable debugging pass here.
- Astra/medium is a reviewer-named review model outside the tier binding; it held
  the citation standard and probed its Coverage claims on `git-admin-readers`.
  Compare again before it joins a binding.
- Give a `codex exec` Coverage axis a writable sandbox in its own worktree; a
  read-only sandbox refuses the Go build cache and the blast table.
- Change routing only after two comparable runs or one controlled model
  comparison.
