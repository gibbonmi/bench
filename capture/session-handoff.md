# Session handoff

Repository: `c7f62432abb513b693e7eb380f2b4015-22fa8276882c8052351c3d896383d3e5` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-3325222104/c7f62432abb513b693e7eb380f2b4015-22fa8276882c8052351c3d896383d3e5`
Branch: `bench/assign/c7f62432abb513b693e7eb380f2b4015/22fa8276882c8052351c3d896383d3e5` — HEAD `98f55ca`, clean tree, 0 unpushed commits
Spec: `specs/worktree-enumeration-hang/spec.md` (Status: staged)
Gate: green at `37f929e` — current

## State

FT189's spec and tickets are authored and verified: 24 coverage rows
(`bench coverage --check` green; preflight `rows-owned`/`rows-membership`
green), decision source compiled at
`specs/worktree-enumeration-hang/decisions/`. Verification loops: spec 33
iterations to accept, tickets 2 — the log line in the spec names the largest
catches. Ticket graph: `resolve-git-common-dir` → `refuse-malformed-admin-entries`
→ {`bound-worktree-enumeration`, `report-admin-entry-in-doctor`}; lines
sonnet/medium, sonnet/medium (routing rows), opus/medium, sonnet/low.
Sign-off landed 2026-08-14, so the build may start; tickets commit serially
on one retained integration source per the workflow. Since staging, the
write-spec loop amendments (materiality exit, promise guard,
cheapest-plausible standard) landed with their own anchors and canary
fixtures, and the pre-existing leading-zero-sha flake in
`TestResolveBranchRangeConsumesExport` was repaired (reviewer-approved).
One learnings entry (2026-08-14, 33-round loop, A/B verdict recorded) awaits
the next `/bench-what-next` drain. `decisions/spec-build-review-gate-cadence.md`
remains invalid (its own shaping resume owns the repair).

## Next command

`$bench-write-spec worktree-enumeration-hang`

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
