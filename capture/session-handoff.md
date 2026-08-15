# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `f0ece251`, clean, 7 unpushed commits
Spec: `specs/skills-index-reader/spec.md` — Status: staged, reviewer-approved 2026-08-15 (two serial tickets: extract-skills-index-module → route-skills-index-verb-and-retire-script).
Gate: green at `f0ece251` — current exact tree

## State

The 2026-08 deepening map lives at `decisions/deepening-2026-08.md`. Landed: candidates
1+3 (gate spec, retired), 6 (`internal/gate/greenmarker`), 5 (`internal/shift/objective.go`).
Candidate 8 is staged as `specs/skills-index-reader` (module `internal/skillsindex`, verb
`bench skills-index`, `.bench/skills-index.sh` deleted; loop 1 ran three rounds and the
reviewer stopped further review — tickets are unreviewed). Remaining after it: adopt spec
(#9), worktree spec (#8), reader tickets (#12). Two open learnings entries await
`/bench-what-next`: map placement on retire, and the write-spec loop cap (the reviewer
asked whether a standing cap exists — none does; the entry proposes one).

## Next command

`/bench-implement-spec skills-index-reader` — fresh mid-tier session, one retained
integration source, tickets serial; lines: story 1 opus/medium, story 2 sonnet/medium.

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
