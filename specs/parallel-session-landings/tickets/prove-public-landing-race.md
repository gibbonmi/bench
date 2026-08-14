# Prove the public landing race

Blocked by: classify-destination-ref-update-failures.md
Writes: internal/systemtest

## What to build

Drive two fresh public processes through the PL15 winner/loser race and PL16
rerun, proving the unchanged reviewed source recomposes and regates after
destination-only movement.

## Acceptance

- [ ] The loser preserves the winner, source review, assignment, and project-green state.
- [ ] The rerun publishes one recomposed merge with exact gate, ref, and marker assertions.
