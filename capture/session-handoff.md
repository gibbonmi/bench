# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — the FT169 refusal-half fix lands as the commit after `9208c4a7`; tree clean after it
Spec: none staged.
Gate: green on the FT169 fix tree (2026-08-23)

## State

FT169's refusal half shipped through `/bench-debug`. The fix is identity
expansion, a stale-executable rebuild and re-run, capture composition by
policy, `capture/` authorization, and one preflight that names every refusal.
The regression tests are `internal/worktree/land_surface_test.go`. The
reviewer's authority half stays open under `Next: decide`.

The reviewer parked three ideas on 2026-08-23: the spec-less landing, the
phase-owned file merge rules, and all work in a worktree. They sit in
`capture/IDEAS.md`. The reviewer asked that they be actioned right after their
tickets are written.
Two learnings are open: the debug fix ran in the main checkout as sole writer,
and the capture policy takes one side per file (union for journals is the
reviewer's call).

Closed decisions: keep merge composition as the landing primitive, no rebase;
the exact-identity posture of `landing-refusal-diagnostics` is reversed by FT169.

## Next command

`/bench-write-spec` — decision source: the three parked ideas plus FT169's
undecided half; the reviewer authorized immediate action on 2026-08-23.

## Shape

Rewritten in full at every phase close, pruned rather than accreted. A fresh
session pays for every line it reads cold; drop anything it would not act on.

Operational gotchas are placed by lifetime, not copied here. One that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes. One scoped to a build
belongs instead in that spec's coverage rows.

This file names at most when you'll hit one, never the command — a second copy
drifts from the source.

Keep the three sections above. **State** holds what is true now, including anything
uncommitted. **Next command** holds the exact harness-native invocation, not a
description of it. This section is the third.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
