# OpenAI model scorecard

Last incorporated landing: `ft311-delegate-preparation` (`4d0354d1784e8f5959759dc074b6ace1aeb98462`, 2026-09-08). Sol/high implemented tickets 1, 2, and the ticket 2 repair.

Terra/high implemented ticket 4 and completed closeout repairs through `676484a`. Three fresh Terra/high axes returned zero findings.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. FT311 telemetry is unknown.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Terra / medium, high | implementation, latest 10 bounded tickets and repairs | Terra/high implemented FT311 ticket 4 with live evidence open. Review found the missing triage input contract, and repair added its five inputs. | Medium for exact one-seam tickets under coordinator mutation. High when one fact crosses multiple policy consumers. |
| Luna / max, medium | implementation, 7 bounded tickets/repairs + reviewer, 4 axis passes on `stable-owner-landing` | Implementation: 2/7 first-pass with all terminal gates green. Review: three initial axes returned 9 raw findings and 8 de-duplicated targets with the citation standard held (one axis refuted four of its own leads by enumeration), and the repair-scoped re-review verified all seven predicates and stayed inside its blocking scope | Low-cost writer for narrow slices under mandatory inspection; standing tier for the three review axes |
| Terra / high, medium, low | semantic review and closeout, latest 10 independent axis passes | Terra/high repaired FT311's tier helper, changelog, and stale token. Three repair-scoped source axes found nothing across accepted predicates and later closure and fence deltas. | Standards, Spec, and Coverage review in separate contexts. Further cross-harness verification stopped on user instruction for this run. |
| Sol / low | implementation, bounded tickets and repairs on 1 landing | The delegates returned focused tests and mutation probes. The retirement repair reproduced the FT94 ledger red, changed one owner, proved that restoring the old value made the test red, and restored green. | Exact ticket seams and small repairs under coordinator verification |
| Sol / high | implementation and semantic review, latest comparable lifecycle charges | Sol/high implemented FT311 tickets 1 and 2. Ticket 2 needed one repair for DP10-DP14 coverage and shared retry orchestration. | Kit and security-seam implementation and review. Use for cross-consumer repair when the public seam is exact. |
| Astra / medium | semantic review, 3 axes + 2 re-reviews on 1 landing | On `git-admin-readers` the three axes returned 13 raw findings and 6 repair targets; the Coverage axis probed a symlink-plus-`..` root and a symlinked temp parent and observed both breaks, and the Spec axis found the check walked only two directories. Both re-reviews closed every predicate and stayed inside the blocking scope. A read-only sandbox refused the Go cache, so the blast table and `bench test --check` ran outside it. | Reviewer-named implementation review when the tier binding is not the reviewer's choice; give the Coverage axis a writable worktree |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| FT311 ticket 2 repair | Sol/high unified charge and proposal orchestration and completed every DP10-DP14 partition. Exact-baseline tests and two restored mutations bit. | delegate | Use Sol/high when bounded repair crosses command retry, authority, and graph semantics. |
| FT311 ticket 4 | Terra/high produced build and triage guidance. Review found that active guidance omitted the required triage inputs. | delegate | Trace each spec input into final phase text before returning a guidance ticket. |
| FT311 closeout | Terra/high repaired the tier helper, changelog, and stale-token wording. Three fresh axes returned zero findings. | delegate and reviewer | Keep closeout repair and independent review at high effort when kit guidance and conformance move together. |
| FT311 linked dogfood | Terra used `gpt-5.6-terra`; the shift exited 0 with one commit, two iterations, and no recovery. | orchestrator | Use a linked repository when the kit repository cannot exercise its installed surface. |
| landing-refusal orchestration | Sol/high preserved exact review identity, proved four independent mutations, fixed the review-loop policy, and resumed a green published-but-incomplete landing; it paid avoidable pickup and full-gate churn before breaking the loop | orchestrator | Treat recursive review as a process red and build a tight repro before accepting another unrelated repair target |

## Current decisions

- Use Terra/medium for ticket-sized write charges under coordinator mutation probes.
- Stop further cross-harness verification for FT311 as the user directed.
- Use Terra/high for final independent review when closeout changes kit guidance and conformance together.
- Use Sol/low for exact ticket seams and small repairs when a coordinator can run a distinct mutation probe.
- Use Sol/high for bounded repair crossing command retry, authority, and dependency-graph behavior.
- Require fixture-only, registry-only, combined, deduplicated, exact-ownership, prefix-ownership, and blocker-graph cases for ownership-closure proposals.
- Settle disputed findings against the frozen candidate with an exact reproduction.
- On a landing gate red without a green baseline, run the fresh baseline before diagnosis.
- Give a Codex Coverage axis a writable isolated worktree when it must run throwaway probes.
- Use only Terra and Sol for the remaining FT311 reviewer roles.
- Change routing only after two comparable runs or one controlled model comparison.
