# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — base HEAD `fb7338d`, approved FT216 staging ready for its phase-close commit, 2 prior unpushed commits
Spec: `specs/worktree-cleanup-eligibility/spec.md` — staged and reviewer-approved
Gate: `bench coverage --check worktree-cleanup-eligibility` green with 33 rows; preflight row ownership/membership green, `paths-authorized` red because main's stale review base spans inherited changes and the staging includes roadmap/map-promotion paths outside build fences

## State

FT216 now has one staged behavior-preserving spec and eight tracer-bullet tickets.
The ready map moved, rather than copied, to
`specs/worktree-cleanup-eligibility/decisions/deepening-2026-08.md`; that is the
spec's sole `Decision source:`. FT217 and FT218 still cite the same compiled map,
so their remaining decisions must be re-homed before the FT216 spec retires.

Closed decisions stay closed: characterize all 9 explicit and 13 automatic
`(Action, ReasonCode)` outcomes before extraction; one typed eligibility owner
holds ownership, assignment, lock, lease, landedness, recovery, tracked,
nested, and ignored-residue policy; automatic cleanup is a stricter reading of
that verdict; `DiscardBranch` remains explicit-only and derived-after; no
precedence anomaly or observable cleanup behavior changes in this build.

The ticket DAG is 01 explicit characterization → 02 automatic characterization
→ 03 typed verdict expansion, then two branches: 04 automatic/shared landedness
→ 06 apply-under-lock → 07 release/landing, and 05 landed-set preservation. 08
contracts both branches and rewrites ADR 0005. Tickets 04 and 05 are independently
green but both may write `clean_landed.go`, so implementation serializes them on
the one retained integration source.

The reviewer-requested Terra pair completed: one Terra authored the slice; a
second independent Terra-high round found two blockers in iteration 1 and
accepted the repaired pair in iteration 2. The reviewer then approved the staged
spec, ticket DAG, shared-file serialization, and decision-map re-home condition.
The mandatory learning in
`capture/learnings.md` adds a shared-reader census and a per-consumer ship test
before future ticket slicing.

## Next command

`$bench-implement-spec using specs/worktree-cleanup-eligibility/spec.md`

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
