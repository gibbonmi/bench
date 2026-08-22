# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5f1b245`, 17 dirty paths, 14 unpushed commits
Spec: `specs/inherited-toolchain-environment/spec.md` (Status: staged)
Gate: green at `f10e405` — stale, work tree `befab75`

## State

The 2026-08-22 drain retired `specs/roadmap-flow/`, drained its retro (eight
items: three fed FT238 and FT120, two built as the light-path kit edit, three
dismissed), verdicted three journal entries (FT214 fed; STE experiment and
blocked-chain entries dismissed after the kit edit), reworded FT172 for the
shipped `Next:` grammar half, and refreshed the sequence. The board holds 72
rows; the flow window reports a net delta of 0. Delegate worktrees
`roadmap-flow-t01..t05`, `spec-std`, and `shape-roadmap-growth` still exist;
`bench worktree clean` is the reviewer's call.

## Next command

`/bench-implement-spec specs/inherited-toolchain-environment/spec.md`

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
