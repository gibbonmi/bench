# Correct the hook-marker comment

Blocked by: none
Ownership fence: `internal/adopt/link_hook.go`
Integration surfaces: public marker documentation→`internal/adopt/link_hook.go` + HMC1
Contracts: none crosses
Closure: HMC1/current-marker-consumers

## What to build

Make the public `PrePushMarker` doc comment state only current marker semantics and consumers, without naming the removed `ClassifyPrePush` API. This is its own green maintenance cut: it changes no runtime behavior or test fixture, so grouping it with executable-mode repair would strand no project-gate red and would conceal independent stale documentation.

## Acceptance

- [ ] [HMC1] (covers local) The `PrePushMarker` doc comment accurately describes the current marker role and does not name `ClassifyPrePush` as a consumer.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| HMC1/current-marker-consumers | restore `ClassifyPrePush` to the marker doc comment | Standards review | enumerate current `PrePushMarker` and `ClassifyPrePush` consumers, expect the stale-comment finding |
