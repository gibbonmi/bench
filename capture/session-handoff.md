# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `29da04eb`, 5 dirty capture/roadmap paths (this drain's batch,
pending reviewer approval to commit), 26 unpushed commits
Spec: none active
Gate: green at `29da04eb` (re-verified independently with this batch's diff in the
working tree)

## State

`/bench-what-next` ran a full reconcile-and-drain. Roadmap: FT212/FT213/FT214 new
rows, occurrence lines added to FT133/FT162/FT169/FT192, FT6 gained a
parked-pending-evidence entry, sequence refreshed to add FT213 as rank 3.
`capture/IDEAS.md` and `capture/learnings.md` both empty; `capture/retros/` removed
(its one retro fully dispositioned into the roadmap, its repair-attribution tally
reported).

One idea item went through "implement now" instead of a roadmap row: the six-column
coverage-schema contraction landed at `29da04eb` (`5af64b76` ticket commit +
`29da04eb` review-repair commit, fast-forwarded from an isolated worktree, gate
green both times). The other two sub-items of that same idea did not: moving the
light-path threshold is blocked on a usage measurement that hasn't run yet (parked
into FT6), and deleting `references/cross-harness-reviewers.md` was refused —
decision map #13 closed as "it survives" and writing that ticket would have reversed
a closed decision.

`capture/agent-performance/claude-models.md` carries `/bench-final-check`'s
already-refreshed scorecard for `spec-ticket-fence-reduction`, included in this
batch unmodified.

Everything above is one uncommitted batch, awaiting reviewer approval before the
single commit that lands it (invariant 4: never commit without that approval).

## Next command

`/bench-shape-idea`

(FT198 — the progressive-roadmap decision `ASSESSMENT.md` ranks 0 — is the top line
of the refreshed `## Recommended sequence`, once this batch commits.)

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
