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
`craft-line` — the line is declared by you, not chosen by the delegate — and
carry both explicitly: the model goes on the call itself as the bound alias
(never omitted — an Agent call without a model inherits *your* model, which is
silent escalation when you run top-tier), and the effort and iteration cap go in the
charge text, because the Agent tool has no effort parameter — effort rides in
the charge or it rides nowhere.

Prefer compressed inputs over inherited context: when a decision map has a
Handoff, give the delegate that digest plus line-ranged excerpts it must quote,
not the orchestrator's whole read list.

A write-delegation from a spec carries its stories' coverage rows — behavior,
seam, red signal — in the charge, every time, and requires the delegate to show
each row red before the edit and green after. That is what makes the done-claim
verifiable by running the gate instead of re-reading the work; a charge without
its rows buys a diff you must read line-by-line to trust.

A worktree-isolated delegate's charge opens with the stale-base check: run
`git merge --ff-only main`, verify HEAD equals main, stop and report if the
merge is denied or diverges. Worktrees get cut behind a moving main, and a
delegate that builds on a stale base re-fights landed work. The orchestrator
holds the other end of the contract: a blocked worktree is fast-forwarded by
the orchestrating session, which then resumes the same delegate. Read-only
delegates are unaffected.

```
Implement story 3 of specs/retry-backoff.md in this worktree. Open with the
stale-base check: run `git merge --ff-only main`, verify HEAD equals main,
stop and report if denied. Coverage rows: [the story's rows]. Effort: medium,
~3 iterations. Commit on green with `bench commit`; return the red→green log
per row.
```
Good — a write-delegation whose opener rides in the charge and whose rows make
the done-claim verifiable.

```
Review this diff on the Standards axis only. Base: run `bench diff` for the
changed files; read AGENTS.md and .bench/BENCH.md for the conventions. Charge:
.agents/skills/bench-craft-review/SKILL.md. Effort: medium, ~1 iteration — one
pass, no fix iteration. Return findings under ## Standards, each citing the rule,
under 400 words. Do not edit any file.
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
your checkout. Concurrent delegates get *separate* worktrees, one each: two
writers in one checkout collide, and a mixed `git status` makes both
done-claims unverifiable. Share a worktree — or run sequentially — only when
one delegate's work genuinely depends on another's output. Read-only
delegations need no worktree; say "do not edit any file" in the charge and
mean it. Review delegates return findings only. The invoking session verifies and
fixes any accepted finding in the checkout that owns the diff; isolated worktrees
are for write-delegations, not for reproducing a read-only review result.

## Verifying the done-claim

A delegate's done-claim is a claim, not a result. Before accepting one:

- run the gate against the delegate's work — its green report is not your green;
- check every coverage row in the charge went red-then-green — a missed row is
  a missed case, found now instead of at review;
- run `git status` in the worktree it used — files touched outside the charge
  are a finding, not a footnote;
- spot-check the citations in any summary before folding it into your own
  report — a delegate's confident paraphrase inherits none of the source's
  authority.

Accepting a done-claim unverified is grading your own work at one remove —
invariant #1 with extra steps.

Report every verification round in one line, like a ladder move: accepted, or
what was missed and where the fix went. A miss the verification already
diagnosed — small, concrete, fully understood — is fixed inline through the
direct fix-and-gate path rather than re-delegated: an edit smaller than its
handoff pays handoff, re-discovery, and re-verification for nothing. Recurring
misses across delegates are a charge defect, not a delegate defect — tighten
the rows in the charge before re-sending it.
