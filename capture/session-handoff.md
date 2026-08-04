# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `d40572f` at the time this was written; the drain commit that carries this file is its child.
Spec: three staged, none active — `specs/exact-prospective-landing/spec.md` (FT188), `specs/ft187-communication-surface-cut/spec.md` (FT187), `specs/pre-push-guard-visibility/spec.md` (FT135). No spec build run is open.
Gate: green at `d40572f`; the drain commit re-runs it.

## State

- `/bench-what-next` drained all three capture sources to zero: six ideas, ten open journal entries, two implementation retros. `capture/IDEAS.md` is empty, `capture/learnings.md` holds its header only, `capture/retros/` is empty.
- Five new roadmap rows: FT188 (exact prospective landing, staged spec), FT189 (`git worktree list --porcelain` hangs on a FIFO gitdir, reproduced at exit 124), FT190 (injected-interface real-producer audit), FT191 (fixture-and-seam inventory generator), FT192 (one-source-per-fact reaches spec and ticket prose). Merged rather than newly rowed: six delegate/repair-discipline clauses into FT164, three receipt and refusal faces into FT184, promote's gate-evidence emission into FT185.
- Reviewer direction taken this run: FT156 enters at `/bench-write-spec` with the anchor-mechanism ruling grilled at spec entry, not as a separate `/bench-shape-idea` session. Its row and the recommended sequence both say so.
- No spec was retired this run — nothing had shipped under `specs/`. `injected-interface-junctions` retired in `87fb4aa`, before this pass.
- Decisions that stay closed: the FT156 route above; FT188 owning the foreign-dirty face that FT169 used to carry alone; the ticket-only `specs/` folders staying FT180's terminal-disposition question.
- The decision map `decisions/parallel-session-landings.md` still has two frontier grills and one blocked on them; the reviewer owns that shaping session.

## Next command

`/bench-write-spec` — FT156, the anchor mechanism. The row's two faces (comment-stripping asymmetry in `requireCollapsed`, substring forbids dying to paraphrase) are one ruling, and it is the reviewer's; expect the grill first thing in that session.

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
