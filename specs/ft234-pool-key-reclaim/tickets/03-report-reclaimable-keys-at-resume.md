# Report reclaimable keys at resume

Blocked by: 01-plan-reclaimable-pool-keys.md
Writes: internal/worktree/lifecycle.go, internal/worktree/resume.go, internal/worktree/worktree.go, internal/worktree/orphan_render_test.go, projects/benchkit.md

## What to build

`bench resume-clean` — the surface the SessionStart hook drives through
`bench session-inspect` — learns that reclaimable pool keys exist and says so,
without touching them. `ConservativeCleanup` counts them by calling ticket 01's
predicate, `ResumeResult` carries the count, and `renderResumeSummary` prints one
line naming `bench worktree reclaim` when the count is non-zero and nothing when
it is zero.

The count is not a second derivation: it is the same predicate the verb plans
with, so the ambient number and the verb's target count cannot disagree. Resume
removes no pool key under any circumstance — the `OrphanCandidate` posture,
where the sweep reports and the explicit command acts.

Update the `projects/benchkit.md` worktree seam paragraph to state the new verb
and resume's report-only posture.

## Acceptance

- [ ] `bench resume-clean` reports the reclaimable-key count and names `bench worktree reclaim` as the action, and prints no such line when the count is zero (RS1).
- [ ] the count `bench resume-clean` reports equals the target count `bench worktree reclaim` plans over the same pool, for both a hostile-shape pool and a clean one (RS3).
- [ ] a `bench resume-clean` run over a pool holding reclaimable keys leaves the pool listing byte-identical (RS2).
