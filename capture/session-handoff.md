# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `9df6e5f`, clean tree, 28 unpushed commits
Spec: `specs/exact-prospective-landing/spec.md` (Status: staged), `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/go-build-cache-footprint/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `784fed4` — current

## State

**Phase reached: critical lifecycle repairs are integrated; one preserved adapter release remains before exact-candidate review.**

Run `869a2d33cab96f962882762a9ea56fd21c952976248e64fa917f6da6c48dbca3` has candidate `7433d559a5573d1c790615818dfbe8e5d5ccc6b1`. Nineteen assignments are released. Assignment `2112f6f281eb24996967c5e2826a8879` is integrated and cleanup-pending with a clean checkout; its recorded integration `e7a1b2be62f45b8654a320e7c48be8ee86b693b6` predates the final no-op sibling commit. Recompose onto this handoff commit, build the sealed CLI from that exact candidate, and retry the same assignment integration to release it without another candidate commit.

Critical prerequisites now present in the candidate include the FT78 symlink fixture repair, clean refresh/checkpoint/integration, verified clean provisional release, and the bounded operation-journal extension from 64 to 128. Their focused packages, vet checks, and red mutations are green/restored; the latest whole-tree prerequisite gate is green. Fable and then Opus stalled on the clean-release write charge, so the coordinator completed and verified its two-file fence. The four rare runtime edge probes and cleanup-pending sibling-tip drift are parked as noncritical; do not add cross-harness review.

After release, perform only Standards, Spec, and Coverage review against the exact candidate, submit the receipt, and promote. Preserve candidate identity during review; concrete accepted findings return through lifecycle repair assignments.

## Next command

`bench spec build promote exact-prospective-landing`

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
