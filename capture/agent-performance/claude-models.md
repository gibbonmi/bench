# Claude model scorecard

Last incorporated landing: `light-path-worktree-clean-portable-path` (`e0646545`,
2026-08-16) — Opus solo as debugger, implementer, and orchestrator, no delegates.

Three completed landings are now recorded. Opus has a substantial implementer sample (16
ticket-sized delegate charges plus 3 inline light-path fixes) and a first reviewer
sample; Sonnet and Fable remain single-sample and unobserved respectively. Routing still
follows the project's harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / unobserved | no comparable landing sample | unknown | Project top-tier candidate; collect terminal landing evidence before specializing |
| Opus / medium | implementer, 16 ticket-sized charges across 2 landings | 13 of 16 first-pass accepted; 3 coordinator catches; every required mutation bit and production was restored exactly each time | Project mid-tier default for gate and conformance logic; the larger sample holds |
| Opus / high | implementer, 2 process-lifecycle delegate charges + 3 inline light-path fixes | all 5 first-pass accepted; the fresh-process barrier met an exact protocol (inherited pipe, no sleep/poll/name oracle); the inline fixes each traced the caller set before editing and twice rejected the shorter diff that would have broken a sibling verdict | Use high for process-lifecycle, signal, cleanup-authority, and destructive-command work |
| Opus / medium | orchestrator, 2 landings | verified every done-claim independently; caught 1 silent oracle regression, 1 partial-ordering hole, and 1 dependency inversion in an accepted repair; produced 2 vacuous probes across the two runs | Continue; the recurring weakness is probe selection, not verification discipline |
| Opus / medium | reviewer, 3 read-only axes on 1 landing | 13 raw findings, 6 accepted; audited all 14 coverage rows individually and found the single row satisfied below its named seam; 4 no-op dispositions were correctly self-refuted | First reviewer sample; quality supports continuing at mid tier |
| Sonnet / medium | implementer, 1 ticket-sized charge (CLI plumbing at a known seam) | first-pass accepted; diff matched its `Writes:` list exactly; self-probe bit | Project cheap-tier candidate; one clean sample at a known seam |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `light-path-worktree-clean-portable-path` CP2 | Opus / high / implementer | Rejected the obvious seam: expanding the operand in the planner alone would have left the apply path canonicalizing the unexpanded operand, so a destructive command's fingerprint would address a different checkout than it removed. Found by tracing callers before editing, pinned by a test that reddens on exactly that mutation. |
| `skills-index-reader` ticket 2 | Sonnet / medium / implementer | CLI plumbing across 11 paths (verb, routing registry, wrapper label, inventory, doc sweeps, script deletion) with no correction; kept `.bench/BENCH.md` at its 180-line budget without being told the count. |
| `skills-index-hostile-input-hardening` HI10 | Opus / high / implementer | Met a non-negotiable protocol exactly — production-reached barrier, inherited-pipe marker, no sleep, poll, or temp-name oracle — and routed every bound through the existing deadline helper unprompted, so a conformance check that grades exactly that had nothing to flag. |
| `skills-index-hostile-input-hardening` HI3 | Opus / medium / implementer | The recurring miss shape: the ticket stated "classify before suppression", the delegate moved only the *absent* state ahead of the check, leaving bad-bytes producers silently dropped behind an adapter. Repaired cleanly once handed the reachable-state argument. |
| `skills-index-hostile-input-hardening` review repair | Opus / medium / implementer | Unified two refusal grammars by inverting a package dependency and changing an unpinned user-visible diagnostic — then flagged both costs itself and offered the veto rather than presenting it as done. Self-reported risk was accurate and actionable. |

## Current decisions

- Preserve unknowns until a completed landing supplies evidence; Fable stays unobserved.
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
