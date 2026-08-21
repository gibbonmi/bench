# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `a1cb0ed`, 2 dirty paths, 2 unpushed commits
Spec: `specs/learnings-dated-line-visibility/spec.md` (Status: staged)
Gate: green at `602a9d0` — stale, work tree `e1e379d`

## State

Phase reached: review complete and accepted, landing next.

The build runs in a retained integration worktree labelled `ft243-integration`,
frozen review base `a1cb0ed5`, source tip `3c8fd1ea`. Address it as
`bench worktree exec "ft243-integration" -- <command>`; never use the path
`bench worktree path` prints, which begins with an unexpanded `~`. The worktree
carries its own `./dist/bench`.

Four commits stand there, each gated green: ticket 01 (the dated-line rule),
ticket 02 (the entries-marker export and the unaccounted-run rule), the staged
spec, and repair ticket 03. The three-axis review returned three raw findings
that collapsed to one repair target — `Parse`'s end-of-input flush was
deletable with the whole suite green. Ticket 03 closed it with DL35 and DL36,
and the coordinator re-ran both target mutations independently: each reds its
own row and nothing else.

Declined and recorded rather than fixed: `Parse` is now 76 lines with four
inline responsibilities, which `bench structure` does not flag; and no test
enforces the "one marker literal in Go sources" claim, which is review's to
grade.

The destination is the main checkout. Its stale untracked `specs/` copy is
removed, because the source carries the amended spec and ticket 03.
`capture/learnings.md` holds three open entries raised this session and is
uncommitted — it must land before `bench worktree land`, which needs a clean
destination.

## Next command

`/bench-implement-spec specs/learnings-dated-line-visibility/spec.md --full`

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
