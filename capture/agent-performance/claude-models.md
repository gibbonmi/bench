# Claude model scorecard

Last incorporated landing: `skills-index-reader` (`d3354cd7`, 2026-08-15) — Opus as
orchestrator and implementer, Sonnet as implementer.

One completed landing is now recorded. It is a single sample per role, so routing still
follows the project's harness-to-tier binding; the rows below carry evidence, not yet a
decision to specialize.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / unobserved | no comparable landing sample | unknown | Project top-tier candidate; collect terminal landing evidence before specializing |
| Opus / medium | implementer, 3 ticket-sized charges (module + parser collapse, 2 repair passes) | 2 of 3 first-pass accepted; 1 coordinator catch; all mutations bit and production was restored exactly | Project mid-tier default for gate and conformance logic; evidence supports it |
| Opus / medium | orchestrator, 1 landing | verified every done-claim independently and caught a silent oracle regression the suite missed; also produced 1 vacuous probe and 2 working-agreement slips (piped CLI output, no-op `cd`) | Continue; the slips are process, not capability |
| Sonnet / medium | implementer, 1 ticket-sized charge (CLI plumbing at a known seam) | first-pass accepted; diff matched its `Writes:` list exactly; self-probe bit | Project cheap-tier candidate; one clean sample at a known seam |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `skills-index-reader` ticket 1 | Opus / medium / implementer | A faithful four-parser collapse that still changed the oracle: diagnostics keyed by skill name dropped one of two for a colliding skill. Repair was clean once the finding was handed back with the failing input. |
| `skills-index-reader` ticket 2 | Sonnet / medium / implementer | CLI plumbing across 11 paths (verb, routing registry, wrapper label, inventory, doc sweeps, script deletion) with no correction; kept `.bench/BENCH.md` at its 180-line budget without being told the count. |
| `skills-index-reader` repair passes | Opus / medium / implementer | Both single-pass. Second pass collapsed a duplicated path format to one const and proved it by a swap probe reddening *both* the renderer and matcher sides. |

## Current decisions

- Preserve unknowns until a completed landing supplies evidence; Fable stays unobserved.
- Change routing only after two comparable runs or one controlled model comparison.
- One sample per role is evidence, not a mandate: keep the tier binding as the default.
