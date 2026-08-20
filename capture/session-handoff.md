# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — 21 unpushed commits, clean tree.
Gate: green.

## State

The FT228 drain is closed. The retro, the refreshed Claude scorecard, and the
one open learnings entry all took verdicts in one commit; `capture/IDEAS.md`,
`capture/learnings.md`, and `capture/retros/` are empty.

Dispositions: FT224 gained the two worst diagnostic faces FT228 exposed — a
mis-rooted gate that drops phases and reports an unstaged environment instead,
and a `worktree release` refusal pointing at `--discard-ignored`, a flag
`release` rejects — plus the `worktree land` light-path refusal, which carries a
reviewer decision on whether the verb should serve a tickets-only spec. FT213
took the red-attribution and probe-derivation occurrences. FT214 took the
spec-authoring learning: a spec flipping an invocation policy enumerates every
surface the harness matches, and a spec adding fixtures verifies the family's
registration or fences the registry. FT237 is new and LOW: `craft-line` states
the common case behind ceiling-not-binding.

No spec was retired — `specs/` holds only light-path directories with no
`spec.md`, which sit outside `bench spec retire`'s universe.

Nothing is pushed. The push is the reviewer's.

## Next command

`/bench-shape-idea` — FT225: decide whether a review may amend the spec in the
source it reviewed.

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
