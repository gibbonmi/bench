# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `7ae73f01` plus the FT189 shaping artifacts
Spec: none staged.

## State

FT189 shaping opened `decisions/worktree-enumeration-hang.md`. Research ticket
#1 is resolved: the hang is reproduced on git 2.43.0 (any FIFO admin entry
wedges `git worktree list --porcelain`; matrix and repro script in
`decisions/assets/worktree-enumeration-hang-probe.md`), Bench inherits it
through the sole enumeration owner `git.Worktrees`, and `internal/bounds` is
the established bound seam. Grill #2 (mitigation posture — pre-scan refusal,
execution bound, or both; recommendation: both) is the frontier; #3 (scope)
and #4 (refusal disclosure/repair) are blocked behind it. All three need the
reviewer live. CONTEXT.md gained the **worktree admin entry** term.

`decisions/spec-build-review-gate-cadence.md` is invalid (its Sources cite the
removed `internal/specbuild/checkpoint.go`), which alone keeps `bench maps` at
exit 1; repairing it belongs to that map's own shaping resume, not FT189.

## Next command

`/bench-shape-idea` — resume `worktree-enumeration-hang` at decision ticket #2 with the reviewer present.

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
