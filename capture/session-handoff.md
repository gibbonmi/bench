# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD at this commit (see `git log -1`), pushed to none pending
Spec: none staged (`specs/` empty)
Gate: green, fresh run on this exact working tree before commit

## State

`/bench-what-next` reconciled the board and drained all three capture sources
to zero, plus verified an out-of-band request: two architecture-review HTML
docs in `/tmp` (2026-08-15 survey at `main@e91d0cb3`, 2026-08-16 as-built
companion at `main@c0f387db`) claimed four of eight deepening candidates were
built. All four verified true against the tree: gate-transaction +
verdict-registry (`internal/gate/run_transaction.go`,
`verdict_registry.go`, guard test `TestSkillsIndexConformanceCarriesNoSecondReader`
sibling `verdict_registry_guard_test.go`), shift-objective owner
(`internal/shift/objective.go`), the green-marker reader
(`internal/gate/greenmarker/greenmarker.go`), and the skills-index reader
(`internal/skillsindex/`, `.bench/skills-index.sh` deleted). The remaining
three candidates (worktree eligibility consolidation, adopt-lifecycle
one-decision, git.Output named readers) were unbuilt and untracked anywhere
in `ROADMAP.md` or `bench maps` — added as FT216/FT217/FT218. A fourth new
row, FT215, captures the progressive-roadmap retro's "no changed-package
gate path" finding.

Reconcile: FT198 (progressive-roadmap) had already landed and spec-retired in
a prior session but two prose references survived (`Goal track: guidance
prose` step 1, `Recommended sequence` line 1) — removed. The one pending
retro (`capture/retros/progressive-roadmap.md`) is fully drained: its
Bench-CLI binary-trust finding folded into FT177, its ticket-fence finding
into FT174, its parallel-repair-porting finding into FT205, its
changed-package-gate finding became FT215; the retro file is deleted.
Repair-attribution tally from that retro: 9 tickets, 1 one-shot repair (cause:
tree-drift), 8 clean.

**Untouched by this pass, pre-existing:** 7 unresolved decision maps
(`bench maps`), one invalid (`spec-build-review-gate-cadence`), one
out-of-pool worktree, and the ~61 structure issues restated by the 2026-08-15
architecture review as 8 named deepening candidates (4 now built, 3 tracked
as FT216-218, 1 folded into the already-shipped `internal/skillsindex`
build). None are this pass's duty.

## Next command

`/bench-shape-idea` for FT207 (worktree-mutating malformed-admin refusal) —
top of the refreshed `## Recommended sequence`, unblocked now that FT189 has
landed.

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
