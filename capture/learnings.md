# Learnings — usage journal

## 2026-08-13 — Keep the map-to-ticket flow in one session [open]

What happened: the suggested fresh-session boundary fell between the reviewed decision map and the spec-to-ticket flow, forcing the next session to reconstruct planning context before the tickets were sliced.

Right behavior: carry the reviewed map through spec authoring and reviewer-approved ticket slicing in the same session. Start the fresh implementation session only after the ticket graph is sliced and recorded.

Proposed rule change: move the suggested fresh-session handoff to the end of ticket slicing, with the continuation resuming at `/bench-implement-spec` against the approved ticket graph.
