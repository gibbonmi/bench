# Session context efficiency

Status: ready

## Destination

A Bench phase returns bounded, task-shaped tool evidence and uses one
owner-level operation for one repeated domain intent. Tool-result text is the
primary cost to reduce; call consolidation is secondary. This map decides the
projection, aggregation, measurement, and enforcement boundaries, then routes
the resulting work to the existing CLI, guidance, and performance owners.
The ready map supplies one coordinating spec and the four-child structure in ticket 7.
Readiness authorizes spec authoring; ticket 9 retains the review checkpoint for new numeric budgets.

## Notes

- **Tool-result text** is the text carried by tool-result payloads into the
  session transcript. It is directly countable, but it is not a provider token
  or dollar measure.
- **Result projection** selects the fields, ranges, failures, or aggregates
  needed for the next decision and gives an explicit route to complete detail.
- **Call consolidation** replaces repeated invocations for one domain intent
  with one owner-level operation. A wrapper that concatenates several complete
  results has grouped calls without consolidating the intent.
- Projection precedes aggregation. Reducing call count does not license a
  larger returned payload.

## Decisions so far

- [What did the 2026-09-11 session actually spend on tool results?](session-context-efficiency/tickets/1.md): 31 large results produced 73.2% of all tool-result text.
- [Which cost leads prioritization?](session-context-efficiency/tickets/2.md): Tool-result text is primary; call consolidation is secondary (reviewer, 2026-09-11).
- [Which result surfaces require bounded defaults?](session-context-efficiency/tickets/3.md): Bounded raw reads and Bench queries retain full-detail routes; wrapper overflow depends on harness capability.
- [Which repeated sequences earn one owner-level operation?](session-context-efficiency/tickets/4.md): Product changes cover set-wide cleanup and selected multi-target queries; archive batches belong in guidance.
- [Which measures are reliable across harnesses?](session-context-efficiency/tickets/5.md): The source matrix distinguishes documented contracts, local observations, and unresolved capability gaps.
- [Which findings deny, advise, or only measure?](session-context-efficiency/tickets/6.md): Measure and advise; replacement requires runtime proof and preserves original results when unsupported.
- [Which multi-target query outcomes belong in scope?](session-context-efficiency/tickets/8.md): Multi-target queries cover worktree identity, path, and state, plus bounded spec histories.
- [How does the work set numeric result budgets?](session-context-efficiency/tickets/9.md): Measurement precedes reviewer approval of each new per-surface byte budget.

- [How does the work split across existing owners?](session-context-efficiency/tickets/7.md): One coordinating spec links four child specs and their reviewer checkpoints.

## Not yet specified

## Spec-writer discretion

## Out of scope

- Provider pricing or a composite dollar score without a reliable harness
  source; FT231 owns performance measurement.
- A general rewrite of always-loaded guidance; FT100 owns guidance weight.
- Replacing FT173's AXI contract or its command-owned truncation policies.
- Combining a mutation with a later verification, approval, or publication
  boundary solely to reduce the call count.

## Sources

- Path: `specs/session-context-efficiency/decisions/session-context-efficiency/assets/preliminary-assessment.md`
  Supports: the session census, pattern inventory, owner map, and preliminary ranking.
  Drift: re-run after the transcript schema, tool wrapper, or named Bench surfaces change.
- Path: `roadmap/FT231.md`
  Supports: measurement ownership and the prior request for a deep transcript and raw-call census.
  Drift: FT231 lands, retires, or changes its harness/Bench ownership split.
- Path: `.agents/skills/bench-craft-cli/SKILL.md`
  Supports: minimal default schemas, explicit truncation, aggregates, and contextual disclosure.
  Drift: the AXI principles or approved Bench query set changes.
- Path: `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`
  Supports: the existing CLI output owner, consumer constraints, and command-local truncation policies.
  Drift: FT173's destination or compatibility decisions change.
- Path: `roadmap/FT89.md`
  Supports: placement of deterministic operations in the CLI and judgment in operating guidance.
  Drift: FT89 lands, retires, or changes its ownership rule.
- Path: `specs/session-context-efficiency/decisions/session-context-efficiency/assets/harness-measures.md`
  Supports: ticket 5's capability evidence and the enforcement and budget decisions.
  Drift: a harness release, schema change, integration change, or change to the named Bench owners.
