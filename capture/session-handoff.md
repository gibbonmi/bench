# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `0ee2106` pre-commit, clean but for this drain's own batch, 38 unpushed commits
Spec: none staged — `specs/` is empty
Gate: green at `761a839`

## State

**`bench-front-door` (audit item A3) landed, retired, and drained.** The reviewed source
`408faf50..c7a0a1b1` composed onto `main` as `b3dfb922`, publishing the spec at
`Status: implemented`; the source worktree is released and removed. `bench status --route`,
the `setup` and `specs: staged` signals, invocable board actions, the `/bench` and `$bench`
adapters, registry-projected root help, and the `/bench-what-next` → `/bench-drain` rename
are all live. The front door works: the bare wrapper and the binary now give the same
routed answer.

Landing was blocked twice and both blocks are now roadmap rows rather than folklore.
`worktree land --request` had been given the assignment ID instead of the request token,
and the refusal named four possible causes without saying which — repaired through
`bench worktree reauthorize` (FT224). Review had also amended the staged spec in the
source, which the byte-identity check refuses; the amendment landed on the destination
first as `a1d31f4d` under reviewer confirmation (FT225).

**This drain's batch:** FT212, FT190, and FT165 removed as shipped, each verified against
the tree rather than taken from the audit's word. FT224 and FT225 added. FT89, FT215,
FT177, and FT162 took merged retro evidence. The open learning was dismissed as already
implemented at `f8d1dd4c`. `capture/` is empty — no ideas, no learnings, no retros. The
`Recommended sequence` now carries the audit's ranking directly instead of pointing at it,
as its own standing note asked the next full drain to do.

**Not applied, named instead:** the audit's full row-by-row board restructure
(`roadmap-dispositions.yaml` — rewrites, folds, and the deferred backlog) is a
`--restructure` pass and needs the reviewer to invoke it.

## Next command

`/bench-write-spec` for A6 — kit tests stop writing into the operator's real `BENCH_HOME`
(`docs/audits/2026-08-bench-capability/results-fable-high/proposed-roadmap.md`, rank 3, P1,
no dependencies).

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
