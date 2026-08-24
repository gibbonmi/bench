# Session handoff

Repository: `a461608be3817c9335b5591451410c2a-fb2a552a375c750a39ff7a12e76c2ccb` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/a461608be3817c9335b5591451410c2a-fb2a552a375c750a39ff7a12e76c2ccb`
Branch: `bench/assign/a461608be3817c9335b5591451410c2a/fb2a552a375c750a39ff7a12e76c2ccb` — HEAD `b239bdd`, 1 dirty path, 6 unpushed commits
Spec: `specs/module-size-split/spec.md` (Status: staged)
Gate: green at `48ce423` — stale, work tree `105fefc`

## State

Batch 1 of `module-size-split` landed tickets 01 (`internal/landing/`) and 02
(`internal/git/`) as pure same-package moves. `bench structure` fell from 75 to 71
issues. The review over base `bdc4dd6a` and tip `b239bdd1` returned 0 findings on
all three axes, so no `reviews/` artifact exists. The landing is spec-less, because
eleven tickets remain; the spec stays `staged`.

The reviewer asked for a controlled comparison, Opus/low against Sonnet/medium, on
two tickets. Both lines landed first-pass on both tickets; Opus/low made the better
grouping call each time at about half the tokens. The verdict is folded into
`capture/agent-performance/claude-models.md`. The reviewer's stop rule fired after
two tickets, so no third comparison ran. Ticket 13 now cites R13; the build
preflight had rejected the graph because no ticket owned that row.

Four comparison worktrees (`mss-t01-opus`, `mss-t01-sonnet`, `mss-t02-opus`,
`mss-t02-sonnet`) stay retained with uncommitted duplicate work. Their accepted
diffs are already on the integration source. Discarding them is the reviewer's
call: `bench worktree clean <path>` per worktree, paths from `bench worktree list`.

The remaining tickets all have `Blocked by: none` except 13, which waits on 06.
Route them at Opus/low per the scorecard decision.

## Next command

`/bench-implement-spec specs/module-size-split/spec.md --full`

## Shape

Rewritten in full at every phase close, pruned rather than accreted. A fresh
session pays for every line it reads cold; drop anything it would not act on.

Operational gotchas are placed by lifetime, not copied here. One that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes. One scoped to a build
belongs instead in that spec's coverage rows.

This file names at most when you'll hit one, never the command — a second copy
drifts from the source.

Keep the three sections above. **State** holds what is true now, including anything
uncommitted. **Next command** holds the exact harness-native invocation, not a
description of it. This section is the third.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
