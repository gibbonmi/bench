---
name: craft-grill
description: Disciplined one-question-at-a-time elicitation to surface a decision or a spec. Use during /bench-shape-idea decision tickets, before /bench-write-spec when requirements are fuzzy, or any time I say "grill me" or the work can't proceed until an open question is resolved. Reach for this instead of asking five questions at once.
index: surfacing a decision one question at a time
---

# Grill

Interview me relentlessly about every aspect of this plan until we reach a shared understanding. Walk down each branch of the design tree, resolving dependencies between decisions. For each question, provide your recommended answer. The recommendation is the point — it forces a concrete decision instead of an open-ended prompt, and lets me correct rather than compose.

If a *fact* can be found by exploring the codebase, look it up rather than asking me. The *decisions*, though, are mine — put each one to me and wait for my answer. A grill that answers its own questions has stopped grilling.

## Discipline

- One question per turn — a stack of questions in one turn is bewildering, and
  the answers come back partial. Wait for the answer before the next.
- Attach your recommended answer and a one-clause reason. I confirm, adjust, or
  reject.
- Ask through the harness's structured question surface when it has one
  (Claude Code's AskUserQuestion): the recommendation rides as the first
  option, marked as recommended, and the free-text escape keeps "adjust"
  open. One question per call, same as one per turn.
- Go in dependency order — resolve the question that unblocks the most others
  first.
- **Surface decisions; don't make them — and scoping is a decision.** Your job is to
  expose the decision space, not to choose how to slice it. If you spot a natural
  seam where the work could be split, *name it as a decision for me* — "there's a
  clean boundary between X and Y; do you want this spec to stop at X?" — and then
  keep grilling the rest. Naming a seam is surfacing. **Proposing to build one slice
  now and defer the rest, and offering to "make a separate spec" for it, is not** —
  that's you making the scoping call and ending the session early to escape the open
  work. Don't. Whether to slice is mine to decide, against a finished map; it is
  never your exit from the grill. I am fine with a small scope — I am not fine with
  you choosing it for me mid-conversation.
- Stop when the fog is gone: when the remaining questions no longer change what
  gets built. Don't pad with questions for their own sake. A proposed slice is not
  "fog gone" — it's an unanswered scoping question, so surface it and continue.
- **The grill ends in my confirmation, not in your build.** Do not enact the
  plan until I confirm we have reached a shared understanding — the last answer
  landing is not that confirmation.
- **Close on a predicate, not a label.** Once I answer, close each decision by
  restating the answer as the exact predicate it fixes — never an outcome label.
  "Better error messages" and "handles the empty case" name a family of
  behaviors; "an empty input list exits 2 and prints the input path" names one.
  A label-shaped close reads as agreement while leaving the choice open, and the
  ambiguity travels downstream into the spec.

## Form

> **Q:** Should the event store be append-only, or allow in-place edits?
> **Recommend:** append-only, with corrections as new events — it keeps the
> coding-during-film flow undoable and the history auditable.
> Your call?

## What it feeds

In `/bench-shape-idea`, Grill is a decision-ticket type — its output is one
reviewer decision recorded as current state. Decision maps are situational:
use one for a multi-session dependency tree, while ordinary late uncertainty
inside `/bench-write-spec` stays bounded there. When the fog is gone, hand the
reviewed source to `/bench-write-spec`.

Before you call a map ready, run `bench maps` — it lists any decision ticket whose
Answer is still a `— (open` / `— (deferred` placeholder or whose section carries a
`GRILL DEFERRED` banner — and refuse readiness while the map still shows a row. A decision made in
conversation but not written into the map is not recorded; the artifact is the source
of truth, not the chat.
