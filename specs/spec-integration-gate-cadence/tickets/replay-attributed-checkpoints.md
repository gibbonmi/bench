# Replay attributed checkpoints

Blocked by: Validate checkpoint receipts

Ownership fence: `internal/specbuild`
Assumptions: checkpoint provenance remains durable and the checkpoint ref is immutable

## What to build

Integrate verified checkpoints through an exact-old-tip candidate transition.
Fast-path a checkpoint based at the candidate, reconstruct and replay an unchanged
disjoint sibling patch, and retain attribution without releasing its assignment yet.

## Acceptance

- [ ] [R16] Only a verified checkpoint advances the exact candidate tip through one attributed commit.
- [ ] [R17] A later disjoint sibling with an older base replays the byte-identical complete patch without losing earlier bytes.
- [ ] [R18] Overlap, patch drift, conflict, ownership drift, and changed assumptions leave the candidate unchanged and durably route the same assignment back to its delegate.
- [ ] [R19] A concurrent candidate move retries only while patch, fence, and assumptions remain unchanged and every update uses exact-old-tip compare-and-swap.

