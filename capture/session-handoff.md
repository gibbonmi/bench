# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5d4ae93`, 4 dirty paths, 41 unpushed commits
Spec: none staged.
Gate: green at `813e71e` — stale, work tree `fce9053`

## State

`parallel-session-landings` and `roadmap-progressive-index` are published and
retired (`3947bcca`, `eb05d529`). This uncommitted drain removes the stale FT198
row, empties the inbox, journal, and retro capture, and merges their residual
work into FT162, FT169, FT177, FT178, and FT185. The reviewer must approve the
batch before it commits.

Veto surface: FT162 receives the landing-prerequisite guidance; FT169 receives
bounded landing-refusal detail; FT178 receives unknown-flag refusal; FT177
receives stale command-discovery detection; FT185 receives concise green output.
The map-to-ticket learning was already triaged in FT174; the durable-build and
no-truncation rules already exist. No restructure was requested.

Preserve foreign assignment `fdf07b2661ab381f9125643169f1af10` byte-for-byte.

## Next command

`/bench-shape-idea` — the board's leading invocable signal (`decisions`).

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
