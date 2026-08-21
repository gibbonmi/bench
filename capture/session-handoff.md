# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `b2deece`, 2 dirty paths, 0 unpushed commits
Spec: `specs/learnings-dated-line-visibility/spec.md` (Status: staged)
Gate: green at `4cd77bc` — stale, work tree `7352eeb`

## State

`/bench-implement-spec specs/learnings-dated-line-visibility/spec.md --full`
stopped at phase entry and wrote no code. Two things block the build.

The reviewer widened the spec's scope. The spec narrowed FT243's clause "a
capture entry the parser cannot see" to dated lines only and raised that
narrowing for disposition; the reviewer rejected the narrowing and asked for
undated content to be reported too. That reaches past the spec's ownership
fence: the `<!-- entries below -->` boundary must become a parser-owned export
of `internal/learnings` consumed by `internal/adopt`, and the fresh-repo quiet
posture (a scaffold whose prose preamble is exactly that shape) needs its own
stories, coverage rows, and regressions. The spec's own "Out of scope" section
already priced this at 3 edits and 2 gate runs. It needs new stories, new rows,
and a re-sliced ticket.

The reviewer accepted the spec's other two dispositions, and both stand: the
Line stays mid (`opus` / medium) against the cached cheap-tier routing for
parser logic at a known seam, and the build runs a write delegate at that tier.

`bench preflight build learnings-dated-line-visibility` is also red:
`rows-owned,red,"declared row(s) cited by no ticket file: DL1 … DL21"`. The one
ticket's acceptance lines carry no row IDs. The exemplar shape is
`specs/land-executable-freshness/tickets/01-refuse-stale-landing-executable.md`,
which prefixes each line with its row ID. The reviewer approved repairing that
in place, but the scope widening rewrites the ticket anyway, so leave it to the
spec phase.

`capture/learnings.md` holds one new open entry on that preflight escape: the
write-spec phase exits without running the build phase's own entry check, so the
red lands in the phase that cannot fix it.

## Next command

`/bench-write-spec`

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
