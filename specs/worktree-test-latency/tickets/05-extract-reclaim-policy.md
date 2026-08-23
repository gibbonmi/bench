# Extract pool-reclaim policy decisions

Blocked by: 02-add-explicit-effect-inputs.md
Writes: internal/worktree/, internal/worktree/reclaimpolicy/ (new)

## What to build

Extract key protection, process and lease liveness, reclaimability, drift, and
action decisions into a pure child package. Parent adapters continue to own
filesystem classification, registration, process probes, and deletion.

Move reclaim matrices below the typed seam. Retain representative effectful
journeys and focused fact-adapter coverage.

## Acceptance

- [ ] RP1: Typed reclaim tables protect live or uncertain keys and act only on proven dead keys.
- [ ] FA1: Focused adapter tests prove real lease, process, registration, and path translation.
- [ ] RJ1: Real liveness, lease, registration, drift, and deletion journeys remain serial and green.
