---
name: craft-skills
description: Principles for writing and editing skills so they fire reliably and stay lean. Use whenever creating, editing, reviewing, or pruning any skill — the kit's own skills or your project's. Reach for this any time you're about to add or change a SKILL.md.
---

# Writing great skills

A skill exists to wrangle determinism out of a stochastic system. The root virtue
is **predictability** — the agent taking the same *process* every run, not the same
output. Every rule below serves that.

## Invocation: model-invoked vs user-invoked

- **Model-invoked** — keeps a `description`, so the agent can fire it on its own and
  other skills can reach it. Costs context load (the description sits in the window
  every turn). Use when the agent must reach the skill autonomously. In this kit:
  `craft-seams`, `craft-tdd`, `craft-adr`, `craft-cli`, `craft-design-system`, `craft-grill`, `craft-synthesis`, and
  `craft-skills` itself (`craft-synthesis` runs it against every kit change) — the
  agent reaches for them mid-work.
- **User-invoked** — `disable-model-invocation: true`; only you, by name, can fire
  it. It is not a candidate for implicit matching, but *you* are the index that must
  remember it exists. Use for canonical phases or workflows the reviewer drives
  deliberately.

Pick model-invocation only when the agent or another skill must reach it on its own.

## Leading words

A **leading word** is a compact concept already in the model's pretraining that the
agent thinks *with* while running — `seam`, `fog of war`, `tracer bullets`, `the
line`, `the gate is the oracle`. Repeated across the text, it accumulates a
distributed meaning and anchors a region of behavior in the fewest tokens by
recruiting priors the model already holds. It pays twice: in the body it anchors
*execution* (same word, same behavior every run); in the description it anchors
*invocation* (when the word lives in your prompts and docs too, the agent links the
shared language to the skill and fires it more reliably). Hunt for a triad spelled
out three times — `fast, deterministic, low-overhead` → `tight` — and collapse it.

## Information hierarchy

Rank material by how immediately the agent needs it:
1. **In-skill step** — an ordered action, ending on a **completion criterion** that
   is *checkable* (can the agent tell done from not-done?) and, where it matters,
   *exhaustive* ("every modified model accounted for", not "produce a list"). A
   vague criterion invites premature completion.
2. **In-skill reference** — a definition or rule consulted on demand. A flat set of
   peers (every rule on one rung) is fine, not a smell.
3. **External reference** — pushed into a linked file behind a **context pointer**,
   loaded only when the pointer fires. The pointer's *wording*, not its target,
   decides how reliably the agent reaches it.

**Progressive disclosure** is the move down the ladder so the top stays legible.
Inline what every branch needs; push behind a pointer what only some branches reach.
**Co-location**: keep a concept's definition, rules, and caveats under one heading.

## Write for the weakest reader

A skill is executed by every tier in the rotation, including the cheapest
delegate on a headless run. Test each instruction against that reader: could
the weakest model in the binding act on this sentence unambiguously, without
the judgment you had while writing it? An instruction that needs the author's
context to interpret is a defect that compounds — the skill fires thousands of
times, and every ambiguous reading multiplies. (Routing for skill *edits*
lives in `craft-line`; this rule is about the prose itself.)

## Pruning

- **Single source of truth** — one authoritative place per meaning, so a change is a
  one-place edit.
- **Relevance** — does each line still bear on what the skill does?
- **Hunt no-ops sentence by sentence.** Run the no-op test on each sentence in
  isolation; when one fails, delete the whole sentence, don't trim words. Be
  aggressive — most failing prose should go, not be rewritten.

## Failure modes

- **Premature completion** — ending a step before it's done, attention slipping to
  *being done*. Fix the completion criterion first (cheap, local); only if it's
  irreducibly fuzzy *and* you see the rush, split to hide the later steps.
- **Duplication** — the same meaning in two places; costs maintenance and inflates a
  meaning's apparent rank.
- **Sediment** — stale layers that accrete because adding feels safe and removing
  feels risky. The default fate of any skill without a pruning discipline. Run this
  skill against the others periodically.
