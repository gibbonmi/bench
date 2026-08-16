# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `c0f387d`, 20 unpushed commits
Spec: none staged.
Gate: green at `caa483a` — stale against the drain's work tree.

## State

`/bench-what-next` drained all three capture sources in one batch: `capture/IDEAS.md`
and `capture/learnings.md` are empty, `capture/retros/` is gone. `ROADMAP.md` gained
FT208 (harden `internal/skillsindex` against eleven inherited hostile-input edges) and
FT209 (differential exit rows for behavior-preserving refactors; cardinality fixed at a
new grouping's concept edge). Retro and learning evidence merged into FT89, FT98, FT99,
FT133, FT169, FT178, FT180, and FT205. FT89 lost its skills-index single-sourcing clause,
which shipped; its row stays open for the YAML-parsing half.

Also in this commit: both `capture/agent-performance/` scorecards, refreshed by the
skills-index-reader landing (Claude for Opus/Sonnet as implementer and Opus as
orchestrator; OpenAI for Sol as reviewer).

Reviewer-approved follow-up not yet done: write a NEW dated
`/tmp/architecture-review-<ts>.html` carrying as-built Before/After for the five landed
deepening candidates (1+3, 5, 6, 8), leaving `/tmp/architecture-review-20260815T101417.html`
untouched as the record of the survey at `e91d0cb3`. Candidate 8's original After diagram
is superseded — it drew the shell script surviving as a consumer of the module.

## Next command

`/bench-write-spec` for FT208.

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
