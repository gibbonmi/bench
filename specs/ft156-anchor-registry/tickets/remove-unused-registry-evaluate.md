# Remove the unused registry evaluator

Blocked by: single-source-shared-rule-markers.md
Ownership fence: `internal/anchors`
Contracts: registry group evaluation crosses `internal/anchors`→the five interleaved conformance call sites, asserted by AR4 against the absence of an unused all-groups wrapper

## What to build

Delete the unused exported `Evaluate` wrapper. Keep `EvaluateGroup` as the only evaluator surface so future group additions cannot be silently omitted by dead iteration code.

## Acceptance

- [ ] [AR4] `internal/anchors` exposes only the evaluator the conformance composition calls; no dead all-groups wrapper or future-group bound remains.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AR4 | restore the unused `Evaluate` wrapper | coordinator symbol/caller sweep | restore → enumerate callers → reject zero-caller export → remove wrapper |
