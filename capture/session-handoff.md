# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `a5afaff`, 1 dirty path, 19 unpushed commits
Spec: none staged.
Gate: green at `3171007` — stale, work tree `dbd5810`

## State

**Phase reached: FT171 implementation and its measured materialization follow-up are complete.**

`gate-test-concurrency` promoted candidate
`68e84ff002b862eccee3f9d6cfee0ef8ffd3a2c7` as `3686cb3e`; its exact Opus
review was clean and its retirement is `bc2d1aff`. The later shared-fixture
binary intermediate is `ec366fa2`, and the measured local-development repair is
`af23f587`. The final hardlink-isolation and decision-seeding repair is
`a5afaffd`.

The final exact width-one `internal/gate` run is green in 111.67 wall seconds
with 222,784 KiB maximum RSS and 699,904 filesystem-output blocks. Width two is
green in 78.79 seconds with 227,000 KiB RSS and 699,712 output blocks: 29.4
percent faster, 1.9 percent more peak memory, and effectively identical writes.
The combined race and hostile-kit run is green. The canonical repair gate reports
`internal/gate` green in 106.429 seconds. The implementation uses Go's ordinary
shared build cache, one process-scoped immutable fixture binary, hardlinks with
metadata detach and atomic root-local byte replacement, a witnessed copy
fallback, private Git and Bench evidence per fixture, direct validated seeding
for decision-only tests, and a test-only cancellation grace while preserving the
two-second production contract. The architecture and receipts live in
`decisions/assets/ft171-shared-fixture-staged-binary.md` and
`specs/light-path-shared-fixture-staged-binary/materialization-implementation-evidence.md`.

One required retro remains uncommitted at
`capture/retros/gate-test-concurrency.md`. The implementation worktrees were
released cleanly; do not reconstruct them. Drain the retro through the next
phase rather than folding it into another implementation commit.

## Next command

`$bench-what-next`

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
