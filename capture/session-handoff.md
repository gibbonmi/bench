# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main`
Spec: `specs/worktree-enumeration-hang/spec.md` — `Status: staged`, awaiting reviewer sign-off.

## State

FT189's spec and tickets are authored and verified: 24 coverage rows
(`bench coverage --check` green; preflight `rows-owned`/`rows-membership`
green), decision source compiled at
`specs/worktree-enumeration-hang/decisions/`. Verification loops: spec 33
iterations to accept, tickets 2 — the log line in the spec names the largest
catches. Ticket graph: `resolve-git-common-dir` → `refuse-malformed-admin-entries`
→ {`bound-worktree-enumeration`, `report-admin-entry-in-doctor`}; lines
sonnet/medium, sonnet/medium (routing rows), opus/medium, sonnet/low.
Reviewer sign-off on the spec and breakdown is the hard stop before any
build. One learnings entry (2026-08-14, 33-round loop) awaits the next
`/bench-what-next` drain. `decisions/spec-build-review-gate-cadence.md`
remains invalid (its own shaping resume owns the repair).

## Next command

Reviewer: approve or veto the approval table in the closing report, then start a fresh mid-tier build session with `/bench-implement-spec worktree-enumeration-hang`.

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
