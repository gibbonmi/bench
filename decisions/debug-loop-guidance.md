# Transfer the debug structure to workflow phases

Status: ready

## Destination

Decide how debug structure can guide spec writing, ticket slicing, implementation, and semantic review.
One replacement spec covers the debug authorship repair and phase-specific guidance adaptations.

## Notes

The reviewer reopened shaping on 2026-09-15 after an interrupted spec rewrite.
The existing trial spec remains intact until a replacement spec promotes its retained decisions.
Use craft-domain, craft-research, craft-grill, and craft-synthesis.
The map owns its research asset under decisions/debug-loop-guidance/assets/.
The current conversation excludes trial flags and comparative performance claims.

## Decisions so far

- [Which mechanisms transfer?](debug-loop-guidance/tickets/1.md): The source comparison is complete; scope decisions remain open.

- [What structure should each phase adopt?](debug-loop-guidance/tickets/2.md): Preserve debug locally and adapt each of the three phases.

- [What scope and adoption evidence belong in the replacement?](debug-loop-guidance/tickets/3.md): One replacement spec includes separate reviewable outcomes.

- [How should ticket slicing apply the structure?](debug-loop-guidance/tickets/4.md): Add the sequence at the existing craft-tickets owner.

## Not yet specified

## Spec-writer discretion

- Choose concise wording and examples within the resolved behavior and existing authority rules.
- Preserve each rule at one existing owner and use references where needed.

## Out of scope

- Trial flags and arm-comparison machinery, per the current conversation.
- A claim that this guidance has a measured speed advantage.

## Sources

- Path: `decisions/debug-loop-guidance/assets/research.md`
  Supports: Ticket 1 evidence and the recommendations for tickets 2 and 3.
  Drift: Refresh after the cited workflow guidance or the proposed scope changes.
