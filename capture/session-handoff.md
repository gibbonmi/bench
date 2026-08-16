# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `945c9266`, 3 dirty paths, 1 unpushed commit
Spec: `specs/spec-ticket-fence-reduction/spec.md` (Status: staged)
Gate: green at `945c9266` — re-run on the drain batch before commit

## State

`945c9266` landed the AGENTS.md verb-discovery correction: `bench help` is the verb
inventory, `bench commands --brief` is a three-verb liveness probe. The
`/bench-what-next` batch on top of it is prepared and uncommitted under standing
reviewer approval — it adds FT198 (progressively loaded roadmap; `ASSESSMENT.md`
ranks it 0 and it had no roadmap row), folds the parked verb-discovery idea and the
derive-help-from-`commandRegistry` finding into FT89, empties `capture/IDEAS.md`,
and promotes FT198 above FT207 in the sequence. Two calls are contestable: the
FT198 row itself, and demoting FT207 on `ASSESSMENT.md` pricing.

## Next command

`/bench-implement-spec --full spec-ticket-fence-reduction --reviewer fable high`

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
