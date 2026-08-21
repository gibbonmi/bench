# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean, two commits ahead of `origin/main`.
Spec: none staged.
Gate: green at `ba48496f`.

## State

Audit item A4 is closed. `bin/bench.sh` derives the pool home through
`${BENCH_HOME:-${HOME:?...}/.bench}`, so a session with neither variable set gets
the wrapper's own message naming both inputs instead of `HOME: unbound variable`.
`BENCH_HOME` still wins when set. A `internal/systemtest` test pins the wording
independently of the wrapper; a mutation probe confirmed that reverting the guard
reds it. A4's other two pieces shipped earlier, so no residual remains.
`docs/audits/2026-08-bench-capability/results-fable-high/next-ticket.md` is stale —
it still names A1, which shipped 2026-08-18.

Two learnings are open in `capture/learnings.md`, both from this build.
`bench commit --spec` exits `landed-but-checkout-incomplete` on a tickets-only
light-path spec: the commit is correct and green, but the retired ticket folder
stays on disk as untracked residue that the session removes by hand. Its fix needs
a reviewer call on retry-versus-report. The second is a seam learning —
`internal/conformance` admits a live-tree test only through its own registry, so a
policy assertion on the live wrapper belongs in `internal/systemtest`.

The branch had diverged this session; the reviewer reconciled it with
`git pull --rebase` before the build. Nothing is pushed yet.

## Next command

`/bench-write-spec` — FT243: a capture entry the parser cannot see is reported as
zero, not as a failure. This was the drain's top-ranked row before the A4 residual
jumped ahead of it, and it is next again.

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
