# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `c46b135`, 2 dirty paths, 17 unpushed commits
Spec: `specs/roadmap-progressive-index/spec.md` (Status: staged)
Gate: green at `a095153` — stale, work tree `18f7eca`

## State

Roadmap maintenance keeps all 62 rows, folds the write-spec slicing idea into
FT174's existing spec-to-ticket handoff owner, empties the capture inbox, and
advances the recommended sequence from shaping FT198 to finishing its active
implementation. No retro, learning, occurrence, or spec retirement was pending.

FT198 `roadmap-progressive-index` has three green ticket commits on `main`:
`ea5dd47a` add-context-envelope, `403a66a8` rebuild-bare-board, and
`c46b135a` add-row-selector. `land-index-doctrine` is preserved uncommitted in
`/home/devuser/.bench/worktrees/bench-3325222104/62f66d33a5fd9afdf9446f51ce9f3edb-f70c5f8dc98f203fac19bdd6e07df1d3`.

That nine-path diff has PI1-PI5, PI9, PI17, and PI18 green, including recorded
prose-only and schema-only PI17 revert reds. Its final mutation exposed duplicate
default body-omission enforcement in `projectBody` and `ContextCommand`; collapse
that fact to one owner before landing. The prior top-line delegate exhausted its
8-iteration cap, so resume as a fresh declared repair stage rather than continuing
that charge.

After the doctrine ticket lands, continue `join-axi-set`, three-axis composed
review, the reviewer-chosen Fable/high cross-harness falsification pass, the final
`--spec` landing, and `/bench-final-check`. The
`ft153-canary-public-repro` assignment remains foreign and untouched.

## Next command

`$bench-implement-spec --full roadmap-progressive-index`

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
