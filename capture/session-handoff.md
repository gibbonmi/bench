# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `f1ab2825`, uncommitted `specs/progressive-roadmap/` and this file, 28 unpushed commits
Spec: `specs/progressive-roadmap/spec.md` — staged, awaiting sign-off
Gate: green at `f1ab2825`

## State

`/bench-write-spec 198` authored `specs/progressive-roadmap/spec.md` (Status: staged,
28 coverage rows, `go run ./cmd/bench coverage --check` green) and five tickets under
`specs/progressive-roadmap/tickets/`, from the 2026-08-17 reviewer-confirmed grill.
One Opus/high review round BLOCKed on 6 blocking + 8 prose findings; all are folded
and named in the spec's verification log. Uncommitted: the spec folder and this
file. Awaiting the reviewer's sign-off on the approval table (stories and lines,
seams, coverage, fences, out of scope, ticket graph). Two contestable calls flagged
in the spec: it supersedes decision #1 of the retired `roadmap-progressive-index`
map ("`ROADMAP.md` remains the only durable owner"), and stories 8, 25, 38–39,
44–46 are spec-writer additions beyond the grill.

Machine note: the `bench` wrapper resolves an installed 0.2.0 release ahead of the
dev build, so `bench coverage --check` there predates the reduced schema — use
`go run ./cmd/bench` for spec-phase checks until that is fixed.

## Next command

On sign-off: `bench commit -m "spec: stage progressive-roadmap (FT198)" specs/progressive-roadmap capture/session-handoff.md`, then a fresh mid-tier session runs
`/bench-implement-spec progressive-roadmap` on one retained integration source,
frontier ticket `split-the-board-parser-and-migration-in-one-green.md`.

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
