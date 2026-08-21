# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `57ffe0a3`
Spec: `specs/landing-refusal-diagnostics/spec.md` — Status: staged, review round accepted, awaiting reviewer sign-off
Gate: stale (last gated tree matches HEAD)

## State

The FT233 spec and its five tickets are written and uncommitted, plus one
learnings entry for the two-iteration review round. `bench coverage --check`
passes at 23 rows. The one review round (opus/high) accepted after the author
folded its findings; the verification log in the spec records it.

Two flagged decisions gate sign-off, both marked in the spec's implementation
decisions and on their tickets: (1) refusing an abbreviated `--base` is a new
refusal precondition — the tree accepts a short `--base` today (story 2, LR2);
(2) the runtime-root allowance lets `bench worktree clean` plan
`discard-remove` over the gate's `.logs/` records without `--discard-ignored`
(plan/apply still gates the removal). A veto on either cuts the marked story
and row; nothing else depends on them.

Ticket frontier after sign-off: `enrich-refusals-through-one-emitter.md` and
`allow-runtime-root-residue-everywhere.md` (Blocked by: none); the other three
block on the emitter ticket.

## Next command

`/bench-implement-spec` — in a fresh mid-tier session, after reviewer sign-off
of the spec-and-tickets pair; commit the spec first.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
