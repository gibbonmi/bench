# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean tree once the closing capture commit lands; unpushed commits pending review.
Spec: none staged — FT230 is implemented, landed (`b48c4609`), and retired.
Gate: green at the retirement commit.

## State

FT230 is done end to end: adapter selection (`--adapter npm|fixture`,
`--provenance`) through `bench release`, the workflow's publish job swapped to
one `bench release submit` invocation, the conformance contract flipped with
red-capable bites, the two step-name contracts retired to a record-level
ordering test, the runbook amended, the spec retired, and all landed
worktrees cleaned. The retro (`capture/retros/ft230-release-through-bench.md`),
two learnings entries, and refreshed `capture/agent-performance/claude-models.md`
await `/bench-drain`. The newer learnings entry flags a contract contradiction
for reviewer decision: `bench worktree land` refuses on divergent staged spec
bytes while the review skill and a scorecard decision say the landing
publishes the source's bytes.

## Next command

`git push` — then `/bench-drain` for the retro, the two learnings entries, and
the scorecard refresh.

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
