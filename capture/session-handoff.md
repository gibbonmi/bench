# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `595763a1`, 47 unpushed commits
Spec: `specs/progressive-roadmap/spec.md` (FT198) — implemented, reviewed,
landed, retired; the folder no longer exists
Gate: green at `5ee78026` (the `bench worktree land` publish); retained exact
evidence for that commit is the landing's own record, not re-run since

## State

FT198 (the split roadmap board) is fully closed. `ROADMAP.md` carries one
physical heading line per row and no body; each row's detail lives in
`roadmap/FT<n>.md`. Landed in one `bench worktree land` publish (`5ee78026`)
composing a retained integration source that carried: the parser + this
repo's 67-row migration, the `roadmap-detail-integrity` conformance check (9
canary fixtures), CLI rendering, status/dashboard/journey wiring, and docs —
then a three-axis review (Standards/Spec/Coverage, `opus`/medium, 19 raw
findings / 12 repair targets) and four repair tickets closing all accepted
findings, each independently gate-verified and mutation-probed by the
coordinator. Two low-stakes findings are reviewer-accepted residual risk
(recorded in the pre-land closing commit, now folded into `5ee78026`'s
history). `bench spec retire progressive-roadmap` ran post-landing
(`595763a1`): the spec folder, `FT198`'s `ROADMAP.md` row, and
`roadmap/FT198.md` are gone.

**Retro not yet captured.** `/bench-final-check`'s retro-capture step
(`capture/retros/progressive-roadmap.md`, `capture/agent-performance/`
scorecard refresh) has not run this session — do it before drifting to other
work, since the landing evidence and delegate performance detail are freshest
now.

**Untouched by this build, pre-existing:** 7 unresolved decision maps, 61
structure issues, one foreign out-of-pool worktree at
`/tmp/bench-gate-subject-*` — none are this build's duty; they route to
`/bench-shape-idea`, `craft-seams`, and a manual `bench worktree clean`
respectively, whenever a session picks them up.

## Next command

`capture/retros/progressive-roadmap.md` (rewrite in full, the template is in
the `bench-final-check` skill) plus the two
`capture/agent-performance/` scorecards for every provider that served as
implementer, reviewer, or orchestrator on this landing. No gate or commit
needed to reach that duty — the landing is already the verdict — but commit
the retro and scorecard files once written
(`bench commit -m "retro: progressive-roadmap" capture/retros/progressive-roadmap.md capture/agent-performance/...`).
After that, `/bench-what-next` is the next open phase (roadmap/capture
reconcile and drain) if nothing more specific is chosen first.

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

## Machine note

The globally installed `bench` on PATH resolves a stale 0.2.0 release,
untrusted by the gate's freshness seal against this tree's current source.
Working pattern used throughout the FT198 build: `bash scripts/go-build.sh
/home/devuser/workspace/bench dist/bench` from the main tree (gitignored,
safe to rebuild repeatedly), then every invocation needing the current source
as `BENCH_RUN_BINARY=/home/devuser/workspace/bench/dist/bench bench <args>`.
`bench commit`'s and `bench worktree land`'s internal gate calls refuse an
inherited binary for a never-before-graded prospective (uncommitted or
freshly composed) tree state — run `bench gate --fresh` once on the exact
tree state first (caches a reusable green verdict) when using
`BENCH_RUN_BINARY`; `bench worktree land` itself needs no override — run it
plain (unset `BENCH_RUN_BINARY`) so it owns its own fresh build of the
composed tree.
