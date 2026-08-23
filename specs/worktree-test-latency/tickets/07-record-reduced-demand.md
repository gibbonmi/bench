# Record reduced worktree test demand

Blocked by: 06-contract-to-serial-journeys.md
Writes: specs/worktree-test-latency/evidence/ (new), CHANGELOG.md, capture/session-handoff.md

## What to build

Record reproducible before-and-after worktree-subtree and whole-suite runs on
the reference WSL host. Record executable builds, repositories, descendants,
environment changes, and directory changes beside wall time.

Keep publication wait time separate. Hand the landed evidence to the second
spec without adding timing thresholds, `t.Parallel`, or a scheduler.

## Acceptance

- [ ] DM1: Evidence records commands, commit, host conditions, raw spans, and all demand counts.
- [ ] GF1: Fresh measurements retain `-count=1`, normal caches, and one ordinary Go driver.
- [ ] NP1: The completed first spec adds no `t.Parallel` and no scheduler.
