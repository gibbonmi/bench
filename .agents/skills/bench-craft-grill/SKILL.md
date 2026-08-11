---
name: craft-grill
description: Disciplined frontier-round elicitation to surface a decision or a spec — each numbered round asks every question whose prerequisites are settled, with a recommendation per question. Use during /bench-shape-idea decision tickets, before /bench-write-spec when requirements are fuzzy, or any time I say "grill me" or the work can't proceed until an open question is resolved.
index: surfacing decisions in numbered frontier rounds
---

# Grill

Interview me relentlessly about every aspect of this plan until we reach a shared understanding. Walk the design tree as a **frontier**: at any moment, the frontier is the set of questions whose prerequisites are settled. Each round asks the whole frontier at once; each question carries your recommended answer. The recommendation is the point — it forces a concrete decision instead of an open-ended prompt, and lets me correct rather than compose.

Charge `bench-craft-domain` before the first round — canonical terms, Avoid lists, and concept-edge scenarios sharpen what to ask and pin what each answer means.

If a *fact* can be found by exploring the codebase, look it up rather than asking me. The *decisions*, though, are mine — put each one to me and wait for my answer. A grill that answers its own questions has stopped grilling.

## The round

1. **Compute the frontier** — every open question whose prerequisites are
   settled. A question blocked on an unanswered one stays out of the round;
   don't ask what a pending answer could invalidate.
2. **Ask the whole frontier as one numbered round** — number the questions, and
   attach to each a recommended answer with a one-clause why. Don't hold a
   ready question back for a later round, and don't pad the round with
   questions that no longer change what gets built.
3. **Wait for the answers.** I answer by number — confirming, adjusting, or
   rejecting each recommendation. I may answer partially; anything unanswered
   stays on the frontier.
4. **Recompute and repeat** — settled answers unblock new questions and make
   others obsolete: drop the obsolete ones, add the unblocked ones, and open
   the next numbered round. Repeat until the frontier is empty.

Ask through the harness's structured question surface when it has one
(Claude Code's AskUserQuestion): each question in the round rides as one entry,
its recommendation as the first option marked as recommended, and the free-text
escape keeps "adjust" open.

## Discipline

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
  gets built. Don't pad with rounds for their own sake. A proposed slice is not
  "fog gone" — it's an unanswered scoping question, so surface it and continue.
- **The grill ends in my confirmation, not in your build.** An empty frontier is
  not that confirmation, and neither is the last answer landing — do not enact
  the plan until I confirm we have reached a shared understanding.
- **Close on a predicate, not a label.** Once I answer, close each decision by
  restating the answer as the exact predicate it fixes — never an outcome label.
  "Better error messages" and "handles the empty case" name a family of
  behaviors; "an empty input list exits 2 and prints the input path" names one.
  A label-shaped close reads as agreement while leaving the choice open, and the
  ambiguity travels downstream into the spec.

## Form

> **Round 2**
> **Q1.** Should the event store be append-only, or allow in-place edits?
> **Recommend:** append-only, with corrections as new events — it keeps the
> history auditable.
> **Q2.** One projection per consumer, or one shared read model?
> **Recommend:** one shared read model — the consumers query the same shapes.
> Your calls?

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
