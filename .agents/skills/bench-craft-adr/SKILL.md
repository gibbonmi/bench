---
name: bench-craft-adr
description: How to write decision records and project docs for an agent or teammate with no memory of the project's history. Use whenever recording an architectural decision, updating project docs, writing a README, or capturing why something is the way it is. Reach for this any time you're about to document a change — document the resulting state, not the change.
---

# ADR — write for the teammate who just walked in

Every agent session starts cold. It has no memory of how the project got here, so
documentation that narrates history actively confuses the reader you're writing
for. Document the **current decided state**, addressed to someone seeing the
project for the first time.

## When to write one

All three must be true — otherwise don't write an ADR, you'll just add noise:

1. **Hard to reverse** — changing your mind later carries real cost.
2. **Surprising without context** — a future reader will look at the code and
   wonder "why on earth did they do it this way?"
3. **The result of a real trade-off** — there were genuine alternatives and you
   picked one for specific reasons.

What qualifies: architectural shape (event-sourced write model), integration
patterns between modules (events not synchronous calls), technology choices that
carry lock-in (the ones that'd take a quarter to swap, not every library), boundary
and scope decisions (the explicit no's matter as much as the yes's), deliberate
deviations from the obvious path (raw SQL instead of an ORM, because X), and
constraints not visible in the code (a compliance rule, a latency budget). If a
decision is easy to reverse, unsurprising, or had no real alternative — skip it.

## The rules

- Record the decision, not the deliberation. "We use DuckDB for the event store" —
  not "we considered Postgres, then SQLite, then switched to DuckDB."
- History lives in git. If someone needs the path that led here, `git log` has it.
  The ADR is the destination, not the route.
- **No file paths, no code snippets** — they go stale the next session and then
  mislead. Name the module and the decision; let the reader find the code.
- One ADR per decision. When a decision changes, **rewrite the ADR to the new
  current state** and let git hold the old version. Don't append "UPDATE:" notes.

## Format

`docs/adr/NNNN-<slug>.md`, numbered sequentially (scan the folder, increment).
Default to a single paragraph — the value is recording *that* a decision was made
and *why*, not filling out sections:

```markdown
# <decision, as a present-tense statement>

<1–3 sentences: the forces in play now, what we decided, and why.>
```

Add a section only when it earns its place: **Status** (`accepted | deprecated |
superseded by NNNN`) when a decision gets revisited; **Consequences** when a
non-obvious downstream effect needs calling out; **Considered options** when the
rejected alternative is worth remembering so nobody re-litigates it in six months.

ADRs live in `docs/adr/`; decision *maps* from `/bench-ideate` live in `decisions/`. Keep
them separate — one is a settled record, the other is a working plan.

## Why this compounds

This is the same shape as a good ambient-context dashboard and a good skill: hand
the fresh reader the current state, cheaply, so it can act immediately. The clearer
your current-state docs, the better every cold session performs — which is most of
why this discipline pays off so visibly.
