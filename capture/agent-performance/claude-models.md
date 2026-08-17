# Claude model scorecard

Last incorporated landing: `progressive-roadmap` (`5ee78026`, 2026-08-17) —
Opus/medium coordinated 9 ticket-sized charges (5 build, 4 repair) across a
shared retained worktree plus 3 parallel isolated worktrees, ran 3 parallel
Opus/medium review axes, and independently mutation-probed every accepted
delegate diff before committing it.

Five completed landings are now recorded. Opus has a large implementer
sample across both effort tiers and a three-landing reviewer/orchestrator
sample; Sonnet now has a real multi-charge low-effort sample alongside its
earlier single medium-effort one. Routing still follows the project's
harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / high | orchestrator on 1 landing (tickets 6–8, review, landing) + implementer, 1 prose ticket | every done-claim probed with a distinct mutation kind and site; caught three false enforcement claims in spec rows and two unattributed-writer events; as implementer it surfaced an unreachable acceptance criterion instead of forcing it | Top tier for guidance prose and for coordination when the spec's own claims are suspect; not yet observed on Go |
| Opus / medium | implementer, 21 ticket-sized charges across 4 landings | 18 of 21 first-pass accepted; 3 coordinator catches; every required mutation bit and production was restored exactly each time | Project mid-tier default for gate and conformance logic; sample now spans 4 landings |
| Opus / high | implementer, 12 charges (2 process-lifecycle, 4 guidance-prose/anchor, 3 inline light-path, 2 Go-seam parser/conformance-oracle rewrites, 1 conformance-correctness repair) | all 12 first-pass accepted; the largest single charge this landing (a full parser rewrite plus a 67-row repo migration with an empty differential proof) still self-probed cleanly and flagged a spec-vs-charge disagreement rather than silently resolving it either way | Use high for process-lifecycle, cleanup-authority, destructive-command, anchored guidance-prose, and large/foundational Go-seam rewrites |
| Opus / medium | orchestrator, 3 landings | verified every done-claim independently with a mutation at a distinct site/kind from the delegate's own each time; caught 1 silent oracle regression, 1 partial-ordering hole, 1 dependency inversion, 1 under-declared ownership fence, and 1 cross-ticket merge conflict from parallel repair porting; 0 vacuous probes this landing | Continue; deliberate site/kind variance per probe is holding up as the fix for the earlier vacuous-probe weakness |
| Opus / medium–high | reviewer, 3 read-only axes on 3 landings + 1 delta re-review | 21+19 raw findings across landings, cumulative de-duplicated repair targets over 20; every finding cites file:line or story/row ID; this landing's Coverage axis found a diagnostic-emitting code path whose own new trust-list wiring was provably unexercised by any test | Reviewer default at mid tier; high effort added the enumeration-quality citations |
| Sonnet / low–medium | implementer, 5 low-effort ticket-sized charges (this landing) + 1 earlier medium-effort charge | 5 of 6 first-pass accepted; the one exception needed a coordinator-authored merge-conflict resolution at land time, caused by porting parallel repair tickets from a shared stale base rather than a defect in the delegate's own diff | Cheap-tier default holds for known-shape, gate-covered work; a parallel-repair-batch charge should declare its expected touched-file set so ports can be sequenced to avoid conflict |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `skills-index-hostile-input-hardening` HI10 | Opus / high / implementer | Met a non-negotiable protocol exactly — production-reached barrier, inherited-pipe marker, no sleep, poll, or temp-name oracle — and routed every bound through the existing deadline helper unprompted, so a conformance check that grades exactly that had nothing to flag. |
| `spec-ticket-fence-reduction` realign ticket | Opus / high / implementer | Ran the charged omission probe on its own landed glossary entry, saw the gate stay green, and reported that plainly as a finding instead of a pass — which exposed three spec rows claiming enforcement that did not exist and produced two repair tickets. |
| `spec-ticket-fence-reduction` collapse ticket | Fable / high / implementer | Enumerated all 47 anchored needles, then stopped at a spec-level shortfall (60-line budget unreachable from the retained needles) with the arithmetic and three options rather than forcing a fit or silently re-pricing. |
| `progressive-roadmap` split-the-board ticket | Opus / high / implementer | The coordinator's own charge paraphrased the spec's per-class fault disposition wrong; the delegate re-read the spec, found the disagreement, implemented the spec's actual rule, and flagged the discrepancy explicitly rather than silently resolving it either direction. |
| `progressive-roadmap` repair-roadmap-detail-integrity-correctness | Opus / high / implementer | Departed from the ticket's suggested fix (mark a shared map before a `continue`) because that fix would have silently changed unrelated duplicate-ID behavior; used two purpose-separated maps instead and named the trade-off in its report rather than taking the literal instruction. |

## Current decisions

- Fable now has one coordinator/implementer sample on guidance prose; keep it there and
  collect a Go-seam sample before widening.
- A delegate's reading of a gate skip is a claim: plant a break and run the oracle.
- Change routing only after two comparable runs or one controlled model comparison.
- Opus/medium now has a four-landing implementer sample; keep it as the mid-tier default
  and stop treating it as provisional.
- Opus/high is the routing for process-lifecycle, signal, cleanup-authority, and
  large/foundational Go-seam work until a medium-effort sample contradicts it.
- Charges should name every state reaching a branch, not only the state the acceptance
  row tests — both observed repair rounds share that shape.
- A parallel-repair-ticket batch (multiple delegates from one shared stale base, ported
  onto the retained source serially) is real leverage but shifts merge-conflict cost onto
  the coordinator's port step; a charge for one should declare its expected touched-file
  overlap with sibling charges so ports can be sequenced deliberately.
- The zero-delegate light-path shape produces self-catches, not accepted-claim catches;
  treat its evidence as weaker in kind, not merely smaller, when comparing to delegated
  landings.
