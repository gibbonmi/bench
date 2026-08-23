---
name: craft-skills
description: Principles for writing and editing skills so they fire reliably and stay lean. Use whenever creating, editing, reviewing, or pruning any skill — the kit's own skills or your project's.
index: writing or pruning a skill
---

# Writing great skills

A skill exists to wrangle determinism out of a stochastic system. The root virtue
is **predictability** — the agent taking the same *process* every run, not the same
output. Every rule below serves that. Write a skill's prose in ASD-STE100 per
`craft-spec`'s `references/ste-prose.md`; a rule read the same way every run is
the point.

## Invocation: model-invoked vs user-invoked

- **Model-invoked** — keeps a `description`, so the agent can fire it on its own and
  other skills can reach it. Costs context load (the description sits in the window
  every turn). Use when the agent must reach the skill autonomously. In this kit
  the rule is: every `craft-*` skill is model-invoked — they are the mid-work
  guidance surface (`craft-synthesis` runs `craft-skills` against every kit
  change) — and the phase adapters are not.
- **User-invoked** — `disable-model-invocation: true`; only you, by name, can fire
  it. It is not a candidate for implicit matching, but *you* are the index that must
  remember it exists. Use for canonical phases or workflows the reviewer drives
  deliberately.

Pick model-invocation only when the agent or another skill must reach it on its own.

## Leading words

A **leading word** is a compact concept already in the model's pretraining that
the agent thinks *with* while running. Examples: `seam`, `fog of war`, `tracer bullets`, `the
line`, `the gate is the oracle`. Repeated across the text, it accumulates a
distributed meaning. It anchors a region of behavior in the fewest tokens by
recruiting priors the model already holds.

It pays twice: in the body it anchors *execution*: same word, same behavior
every run. In the description it anchors *invocation*. When the word lives in
your prompts and docs too, the agent links the shared language to the skill
and fires it more reliably. Hunt for a triad spelled out three times —
`fast, deterministic, low-overhead` → `tight` — and collapse it.

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

## Contrastive examples

When a skill governs an output surface — an interface, a format, a prose
register — it must carry one contrastive pair. The pair is a good example
beside a bad one, each captioned with the one clause that separates them.
A rule states the boundary; the pair *locates* it, and the weakest reader
imitates an example more reliably than it interprets a rule. One pair per surface is the budget — a
gallery is sediment. Process skills (ordered steps, no output shape) get none;
show the shape only where the shape is the point.

## Write for the weakest reader

A skill is executed by every tier in the rotation, including the cheapest
delegate on a headless run. Test each instruction against that reader. Could
the weakest model in the binding act on this sentence unambiguously, without
the judgment you had while writing it? An instruction that needs the author's
context to interpret is a defect that compounds — the skill fires thousands of
times, and every ambiguous reading multiplies. (Routing for skill *edits*
lives in `craft-line`; this rule is about the prose itself.)

## Pruning

- **Duplication** — keep one authoritative place per meaning, so a change is a
  one-place edit.
- **Relevance** — does each line still bear on what the skill does?
- **No-op** — run the no-op test on each sentence in
  isolation; when one fails, delete the whole sentence, don't trim words. Be
  aggressive — most failing prose should go, not be rewritten.
- **Sediment** — treat growth as a cost. Delete stale layers instead of keeping
  them because removal feels risky.

## Failure modes

- **Premature completion** — ending a step before it's done, attention slipping to
  *being done*. Fix the completion criterion first (cheap, local); only if it's
  irreducibly fuzzy *and* you see the rush, split to hide the later steps.
- **Harness echo** — restating what the harness layer already carries (the
  system prompt, a tool description). Single source of truth spans layers.
  Before a rule lands, check the harness surfaces it will ride beside. An
  echo is sediment, and a near-echo drifts into conflict. A rule the kit must
  own across harnesses is not an echo; name the harness gap it covers.
- **Sprawl** — one skill accreting unrelated jobs until no description can say
  when to fire it; split by trigger, not by topic.
- **Negation** — steering by prohibition drags the forbidden behavior into
  context and makes it *more* available ("don't think of an elephant"). Prompt
  the positive target instead. A prohibition earns its place only as a hard
  guardrail you can't phrase positively, and even then it rides next to the
  positive.
- **Negative space** — every decision a skill declines to make is silently
  delegated to the model's priors, not left neutral. Read a draft for its
  silences and decide each omission deliberately: fill it, or leave it open
  as a real branch.
