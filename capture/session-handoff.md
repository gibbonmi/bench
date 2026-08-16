# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD is the FT210 spec commit (after `dc598762`), unpushed
Spec: `specs/worktree-landed-retirement/spec.md` — `Status: staged`, reviewer-approved.
Gate: green at every commit above (`bench commit`).

## State

Three commits landed this session: `24cad87d` kit remake (`craft-spec` on to-spec,
`craft-tickets` on to-tickets, thin `bench-write-spec`, anchors/canaries repointed, index
regenerated); `dc598762` the previous session's `/bench-what-next` drain plus this session's
learnings/idea; and the FT210 spec commit — `specs/worktree-landed-retirement/` (spec with
37 stories in three outcome groups and 20 coverage rows, 5 approved tracer tickets under
`tickets/`, the compiled map under `decisions/` with #9 re-closed) and `CONTEXT.md`'s
`landed assignment` term. Sol's independent no-review authoring run of the same map (41
stories / 6 rows / 5 tickets, an invented set-receipt seam) was compared and not adopted;
it lives only in the session scratchpad.

Ticket frontier: `count-and-advertise-landed-assignments.md` (blocked by none) →
`plan-the-landed-set-under-one-fingerprint.md` → `apply-the-landed-plan-and-settle-records.md`
→ {`refuse-a-half-applied-landed-set.md`, `make-release-a-workflow-step.md`}.

Open reviewer decision, unrelated: eight `bench worktree list` rows in state `recovered`
hold uncommitted FT208 ticket work under recovery refs. Follow-up you named: slim
`bench-shape-idea` in its own pass; the parked idea to cut `bench-write-spec` further.

## Next command

`/bench-implement-spec worktree-landed-retirement` in a fresh mid-tier session, on one
retained integration source.

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
