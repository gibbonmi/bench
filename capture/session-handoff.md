# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` once this build lands; the reviewed source is the `parallel-landings` worktree branch `bench/assign/b31c954120124fec3e44dd4d9ca17ffa/b0fec9471b338b985b1c9ea1f2be89ab`, frozen base `1a135f1b`, reviewed code tip `7e0a3e00`, and this handoff commit on top of it
Spec: `specs/parallel-landings/spec.md` — `Status: implemented` once the landing publishes; `staged` on the source
Gate: green on every worktree commit (2026-08-23)

## State

The `parallel-landings` build is complete and reviewed. Five tickets and one
repair commit sit on the worktree (`64646d3f`..`7e0a3e00`). Review loop 1
(opus / medium, three axes) found 15 advisory findings and no blocking one;
all 8 accepted repair targets are repaired. Loop 2 (repair-scoped) was clean
on every predicate except the handoff's missing commit pin, which this commit
fixes. The reviewer capped the review at two loops, so this prose-only commit
sits past the reviewed tip; the gate grades it.

This handoff lands with the build through `bench worktree land` from `main`.
Once it lands, `main` carries the spec-less landing, the tickets-only close,
the capture rule table, the repair `next=`, and the worktree-rule guidance.
Drain notes for the next `/bench-drain` are the newest entry in
`capture/learnings.md`. The three parked ideas in `capture/IDEAS.md` close by
implementation at that drain.

Closed decisions: merge composition stays the landing primitive, no rebase;
the journal union and the destination default; the light path joins the
worktree rule; the worktree rule is guidance, not a hook.

## Next command

`/bench-drain`

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
