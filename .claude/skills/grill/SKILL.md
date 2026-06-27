---
name: grill
description: Disciplined one-question-at-a-time elicitation to surface a decision or a spec. Use during /map tickets, before /spec when requirements are fuzzy, or any time I say "grill me" or the work can't proceed until an open question is resolved. Reach for this instead of asking five questions at once.
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
- Stop when the fog is gone: when the remaining questions no longer change what
  gets built. Don't pad with questions for their own sake.

## Form

> **Q:** Should the event store be append-only, or allow in-place edits?
> **Recommend:** append-only, with corrections as new events — it keeps the
> coding-during-film flow undoable and the history auditable.
> Your call?

## What it feeds

In `/map`, a grill is a ticket type — its output is a resolved decision recorded
in the map (current state only). Before `/spec`, a grill closes the gaps so the
spec can be synthesized without an interview. When you've grilled enough that the
build path is clear, say so and hand to `/spec`.
