---
name: craft-research
description: The research discipline — the factual question graph, primary sources only, adaptive round-based fan-out, coordinator re-verification, and one cited durable output per run. Use whenever the work becomes factual reading legwork.
index: doing factual reading legwork / writing a cited research report
---

# Research: read primary, check the joins, leave one cited output

The named failure to prevent is a synthesis whose separate claims are true and
whose combined conclusion is false.

## Trigger and ownership

Fire this skill when work becomes factual reading legwork in shaping,
specification, diagnosis, assessment, or implementation. Each calling phase
keeps authority over its decisions, its artifacts, and its completion contract.
Formal review is not a caller, because the review phase owns its axes. A factual
question that leaves review becomes a run under the phase that owns the answer.

Research establishes source-backed facts, contradictions, unknowns, and
implications. Research never owns a write delegate, a done-claim, a reviewer
decision, or a prototype. Route a write delegate or a done-claim through the
calling phase and `craft-delegate`.

## The question graph

Draw the factual question graph before you read. The question graph is complete
when every load-bearing claim traces to one question node.

Resolve a trivial lookup inline. Give one bounded question to one read-only
delegate when its read set would displace your context, or when you hold useful
concurrent work.

## Rounds and fan-out

Fan out only when at least two frontier questions are mutually independent. Two
questions are mutually independent when neither one needs the other's answer,
shared mutable state, or a concurrency-sensitive measurement.

A round's width is the number of independent frontier questions, capped by what
one synthesis pass can verify. Declare that exact number through `craft-line`'s
fan-out clause before dispatch, and report a widening like a ladder move. This
skill states no tier, no effort, and no iteration cap.

Each charge follows `craft-delegate` for its contents, isolation, and
verification mechanics. Add one research field: the primary-source boundary.

Synthesize a completed round before opening another. A new round names the
conflict, the unknown, or the gap that its synthesis found. Dependent questions
and concurrency-sensitive measurements run serially. The calling phase's
iteration cap bounds the rounds, and an exhausted cap stops with unknowns.

**Good — an independent fan-out that runs in parallel.** Question A asks which
skip tokens `.bench/gate.sh` accepts today. Question B asks which exit codes the
upstream Go test runner documents. Neither question needs the other's answer, so
one round dispatches both.

**Bad — a dependent question that must stay serial.** Question A asks which
tokens the parser accepts today. Question B asks which of those tokens the
migration renames. B needs A's answer, so B waits for the next round.

## Sources and verification

The artifact under study and first-party upstream documentation or APIs are
primary sources. Commentary, a secondary write-up, and unopened recall are not.
A secondary source may locate evidence, but it cannot warrant a finding. When no
primary source is reachable, record an unknown instead of a verified answer.

The coordinator re-opens every source supporting a load-bearing conclusion and
independently checks every join between returns. Delegate returns are ephemeral
inputs, and a summary is not evidence.

A run is complete when every in-scope question is answered or explicitly
unknown, every material claim is traceable, and every contradiction is
reconciled or retained.

## The durable output

One coordinator-authored durable output per research run, keyed by topic. The
output is Markdown, it holds one section per question, and only the coordinator
writes it.

The destination precedence has four steps:

1. Write into the caller's own durable output when the caller has one, and
   create no second file.
2. Otherwise put the output beside the artifact that will consume it.
3. For a shaping map, that place is the map's `decisions/<topic>/assets/`
   folder.
4. When the repository itself is the consumer, use `docs/research/<topic>.md`.

Before dispatch, the caller names the destination, the consuming phase or
artifact, and the condition that retires or refreshes the output. An asset
outside a phase-owned output also states `Consumed by`, `Drift`, and
`Retire when`.

## The report contract

Every load-bearing claim cites an exact local path and line or a primary URL
with its retrieval date. Mutable evidence also names its invalidation trigger.

Tick each element of the checklist below in the output's verification record:

- The output opens with the recommendation, the scope, and the evidence status.
- It keeps facts, inferences, tested results, and proposals in separate
  sections.
- It holds one synthesized section per question.
- It uses a capability or option table when comparison matters, and each option
  records its consequence and a nearby citation.
- It keeps every contradiction and every residual unknown.
- It draws a diagram where prose hides a material relation.
- It states what could not be verified.
- It ends with a validation plan.
- Every load-bearing claim carries the citation the rule above names.

## Compatibility evidence

A byte or wire compatibility claim stays unverified until a separate runnable
probe returns. Research identifies the need for that probe, and the calling
phase owns it. In a decision map the probe is a Prototype ticket, and the
Research ticket names it in `Blocked by`.
