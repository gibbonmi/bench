---
name: grill
description: Disciplined one-question-at-a-time elicitation to surface a decision or a spec. Use during /start-ideation tickets, before /spec when requirements are fuzzy, or any time I say "grill me" or the work can't proceed until an open question is resolved. Reach for this instead of asking five questions at once.
---

# Grill

Surface what's undecided by asking **one question at a time**, each with a
recommended answer attached. The recommendation is the point — it forces a
concrete decision instead of an open-ended prompt, and lets me correct rather than
compose.

## Discipline

- One question per turn. Wait for the answer before the next.
- Attach your recommended answer and a one-clause reason. I confirm, adjust, or
  reject.
- Go in dependency order — resolve the question that unblocks the most others
  first.
- **Surface decisions; don't make them — and scoping is a decision.** Your job is to
  expose the decision space, not to choose how to slice it. If you spot a natural
  seam where the work could be split, *name it as a decision for me* — "there's a
  clean boundary between X and Y; do you want this spec to stop at X?" — and then
  keep grilling the rest. Naming a seam is surfacing. **Proposing to build one slice
  now and defer the rest, and offering to "make a separate PRD" for it, is not** —
  that's you making the scoping call and ending the session early to escape the open
  work. Don't. Whether to slice is mine to decide, against a finished map; it is
  never your exit from the grill. I am fine with a small scope — I am not fine with
  you choosing it for me mid-conversation.
- Stop when the fog is gone: when the remaining questions no longer change what
  gets built. Don't pad with questions for their own sake. A proposed slice is not
  "fog gone" — it's an unanswered scoping question, so surface it and continue.

## Form

> **Q:** Should the event store be append-only, or allow in-place edits?
> **Recommend:** append-only, with corrections as new events — it keeps the
> coding-during-film flow undoable and the history auditable.
> Your call?

## What it feeds

In `/start-ideation`, a grill is a ticket type — its output is a resolved decision recorded
in the map (current state only). Before `/spec`, a grill closes the gaps so the
spec can be synthesized without an interview. When you've grilled enough that the
build path is clear, say so and hand to `/spec`.

Before you call a map resolved, scan it for unwritten answers — any ticket whose
Answer is still a `— (open` / `— (deferred` placeholder or whose section carries a
`GRILL DEFERRED` banner — and refuse to close while any remain. A decision made in
conversation but not written into the map is not recorded; the artifact is the source
of truth, not the chat.