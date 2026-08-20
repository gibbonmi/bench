# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `bef5248`, 4 dirty paths, 25 unpushed commits
Spec: none staged.
Gate: green at `e4f2fcc` — stale.

## State

The drain is complete and committed. `capture/IDEAS.md`, `capture/learnings.md`,
and `capture/retros/` are all empty; the roadmap is reconciled against the tree.
Nothing is pushed — 26 unpushed commits on `main` await the reviewer.

FT225 shipped and its row is retired. Its shipped behavior invalidated a face in
two other rows, both corrected here: FT233 no longer carries the staged-spec
byte-difference diagnostic, and FT224 no longer carries the spec-byte refusal.
The landing now publishes the reviewed source's spec bytes, so a spec amendment
commits inside the reviewed range like any other change.

Two reviewer decisions are queued on the board rather than in capture. FT162 now
carries the handoff-timing collision: committing the handoff at a phase boundary
puts the destination ahead of a retained source's frozen base and reds the next
review preflight's `base-current`, harmless to the landing but training sessions
to read past red. FT222 now carries the competing delegate-tier source:
`projects/benchkit.md` routes review axes mid, one harness's private memory
routes them cheap, and nothing in the tree reconciles them. FT222 moved from LOW
to MEDIUM on that second face — a contestable call, open to veto.

Once FT222 is decided, the losing routing is deleted rather than left in place.
If the profile wins, the harness memory entry naming cheap-tier review axes goes.

`bench worktree clean --landed` plans zero removals: eight pre-existing worktrees
are retained as dirty on ignored `dist/` artifacts, none from a recent landing,
each needing its own per-path resolution.

## Next command

`/bench-write-spec` — FT233, landing refusals say what to do next.

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
