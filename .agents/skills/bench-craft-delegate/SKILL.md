---
name: craft-delegate
description: The delegation discipline — when to spawn a subagent, how to charge and scope it, when it needs an isolated worktree, and how a done-claim is verified. Use whenever spawning any delegate (an axis review, a scoped build, a fan-out search) or accepting one's result. The model/effort half of the decision lives in craft-line.
index: spawning a delegate / verifying a delegate's done-claim
---

# Delegating without losing the plot

A delegate buys two things: parallelism, and isolation of a heavy read-set from
your context. It costs one thing that matters more when misjudged: the work
happens where you can't see it. This skill keeps the cost bounded.

## Delegate or inline

Delegate when the work is parallel (independent axes, fan-out search), or when
its read-set would pollute the coordinating context (a full diff plus standards
docs). Stay inline for small mechanical work at a known seam — a delegate's
setup and verification overhead exceeds the work. Never delegate a decision the
reviewer owns (what ships, spec content, an irreversible choice); a delegate
inherits your authority ceiling, not the reviewer's.

## The charge

A delegate has no conversation memory: everything it needs must be in the
prompt. A complete charge names the objective, the inputs by path, the seam it
works at, the return shape (raw data or a structured report — its final message
is the deliverable), and its budget. Route its model and effort with
`craft-line` — the line is declared by you, not chosen by the delegate.

```
Review this diff on the Standards axis only. Base: run `bench diff` for the
changed files; read AGENTS.md and .bench/BENCH.md for the conventions. Charge:
.agents/skills/bench-craft-review/SKILL.md. Return findings under ## Standards,
each citing the rule, under 400 words. Do not edit any file.
```
Good — self-contained: inputs by path, the charge by path, the return shape and
a write prohibition stated.

```
Review the changes we discussed against our standards, like last time.
```
Bad — "we discussed", "our", "last time": every referent lives in a context the
delegate does not have; it will guess all three.

## Scope

One delegate, one coherent unit — one axis, one story, one search question. A
delegate charged with several loosely-joined jobs returns a summary that blurs
them, and you can't verify what you can't attribute.

## Isolation

A write-delegation runs in an isolated worktree (`bench worktree`), so stray
edits can't land in reviewer-owned files — the delegate gets a checkout, not
your checkout. Read-only delegations need no worktree; say "do not edit any
file" in the charge and mean it.

## Verifying the done-claim

A delegate's done-claim is a claim, not a result. Before accepting one:

- run the gate against the delegate's work — its green report is not your green;
- run `git status` in the worktree it used — files touched outside the charge
  are a finding, not a footnote;
- spot-check the citations in any summary before folding it into your own
  report — a delegate's confident paraphrase inherits none of the source's
  authority.

Accepting a done-claim unverified is grading your own work at one remove —
invariant #1 with extra steps.
