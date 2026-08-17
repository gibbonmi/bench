# Claude model scorecard

Last incorporated landing: `spec-ticket-fence-reduction` (`9d7bb573`, 2026-08-16) —
Opus/medium built the Go tickets, Fable/high coordinated tickets 6–8 with one Fable/high
and four Opus/high write delegates, Opus/high ran the review axes.

Four completed landings are now recorded. Opus has a substantial implementer sample (20
ticket-sized delegate charges plus 3 inline light-path fixes) and a two-landing reviewer
sample; Fable has a first coordinator and implementer sample; Sonnet remains
single-sample. Routing still follows the project's harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / high | orchestrator on 1 landing (tickets 6–8, review, landing) + implementer, 1 prose ticket | every done-claim probed with a distinct mutation kind and site; caught three false enforcement claims in spec rows and two unattributed-writer events; as implementer it surfaced an unreachable acceptance criterion instead of forcing it | Top tier for guidance prose and for coordination when the spec's own claims are suspect; not yet observed on Go |
| Opus / medium | implementer, 20 ticket-sized charges across 3 landings | 17 of 20 first-pass accepted; 3 coordinator catches; every required mutation bit and production was restored exactly each time | Project mid-tier default for gate and conformance logic; the larger sample holds |
| Opus / high | implementer, 2 process-lifecycle charges + 4 guidance-prose/anchor charges + 3 inline light-path fixes | all 9 first-pass accepted; prose charges kept three files exactly at budget, planted every canary for its own reason, and one reported its own probe coming back green rather than claiming coverage | Use high for process-lifecycle, cleanup-authority, destructive-command, and anchored guidance-prose work |
| Opus / medium | orchestrator, 2 landings | verified every done-claim independently; caught 1 silent oracle regression, 1 partial-ordering hole, and 1 dependency inversion in an accepted repair; produced 2 vacuous probes across the two runs | Continue; the recurring weakness is probe selection, not verification discipline |
| Opus / medium–high | reviewer, 3 read-only axes on 2 landings + 1 delta re-review | 21 raw findings, 10 accepted; audited every coverage row both times; on the second landing found three rows whose enforcement claim was false and a stale handoff, and self-refuted 6 no-ops | Reviewer default at mid tier; high effort added the enumeration-quality citations |
| Sonnet / medium | implementer, 1 ticket-sized charge (CLI plumbing at a known seam) | first-pass accepted; diff matched its `Writes:` list exactly; self-probe bit | Project cheap-tier candidate; one clean sample at a known seam |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `light-path-worktree-clean-portable-path` CP2 | Opus / high / implementer | Rejected the obvious seam: expanding the operand in the planner alone would have left the apply path canonicalizing the unexpanded operand, so a destructive command's fingerprint would address a different checkout than it removed. Found by tracing callers before editing, pinned by a test that reddens on exactly that mutation. |
| `skills-index-reader` ticket 2 | Sonnet / medium / implementer | CLI plumbing across 11 paths (verb, routing registry, wrapper label, inventory, doc sweeps, script deletion) with no correction; kept `.bench/BENCH.md` at its 180-line budget without being told the count. |
| `skills-index-hostile-input-hardening` HI10 | Opus / high / implementer | Met a non-negotiable protocol exactly — production-reached barrier, inherited-pipe marker, no sleep, poll, or temp-name oracle — and routed every bound through the existing deadline helper unprompted, so a conformance check that grades exactly that had nothing to flag. |
| `skills-index-hostile-input-hardening` HI3 | Opus / medium / implementer | The recurring miss shape: the ticket stated "classify before suppression", the delegate moved only the *absent* state ahead of the check, leaving bad-bytes producers silently dropped behind an adapter. Repaired cleanly once handed the reachable-state argument. |
| `spec-ticket-fence-reduction` realign ticket | Opus / high / implementer | Ran the charged omission probe on its own landed glossary entry, saw the gate stay green, and reported that plainly as a finding instead of a pass — which exposed three spec rows claiming enforcement that did not exist and produced two repair tickets. |
| `spec-ticket-fence-reduction` collapse ticket | Fable / high / implementer | Enumerated all 47 anchored needles, then stopped at a spec-level shortfall (60-line budget unreachable from the retained needles) with the arithmetic and three options rather than forcing a fit or silently re-pricing. |

## Current decisions

- Fable now has one coordinator/implementer sample on guidance prose; keep it there and
  collect a Go-seam sample before widening.
- A delegate's reading of a gate skip is a claim: plant a break and run the oracle.
- Change routing only after two comparable runs or one controlled model comparison.
- Opus/medium now has a two-landing implementer sample; keep it as the mid-tier default
  and stop treating it as provisional.
- Opus/high is the routing for process-lifecycle, signal, and cleanup-authority work
  until a medium-effort sample contradicts it.
- Charges should name every state reaching a branch, not only the state the acceptance
  row tests — both observed repair rounds share that shape.
- The zero-delegate light-path shape produces self-catches, not accepted-claim catches;
  treat its evidence as weaker in kind, not merely smaller, when comparing to delegated
  landings.
