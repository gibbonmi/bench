# Release integrated assignments

Blocked by: Replay attributed checkpoints

Ownership fence: `internal/specbuild`
Assumptions: integration records the candidate and checkpoint relationship before cleanup

## What to build

Let the lifecycle owner release an integrated assignment only after its durable
checkpoint relationship is recorded. A fault between candidate advancement and
release must resume cleanup without replaying the patch or treating an unavailable
releaser as success.

## Acceptance

- [ ] [R20] Successful integration records provenance before release, requires a real release owner, exposes pending cleanup as the next action, and re-entry releases without a second candidate commit.

