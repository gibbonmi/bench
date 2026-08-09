# Migrate ambient outcomes and actions

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/status`, `internal/handoff`, `internal/dashboard`
Integration surfaces: returned-output adapter→migrate-output-adapter.md; shared carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-outcome-action-routes.md
Contracts: status/handoff kind string, exit integer, ordered action tokens, prose disposition, and dashboard facts cross `internal/status` and `internal/handoff`→`internal/dashboard`, membership is clean/active/stale forms, order is producer then composition, and absence is zero actions, asserted by AM1 and AM2
Closure: AM1/status-route, AM1/handoff-route, AM1/dashboard-route, AM2/fixed-action, AM2/prose, AM2/bytes

## What to build

status, handoff, and dashboard query outcomes reach the shared outcome route. status and handoff actions remain typed before dashboard composition with exact public rendering.

## Acceptance

- [ ] [AM1] (covers OA3) status, handoff, and dashboard query outcomes reach the shared outcome route.
- [ ] [AM2] (covers OA4) status and handoff actions remain typed before dashboard composition with exact public rendering.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AM1/status-route | restore local status outcome | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| AM1/handoff-route | restore local handoff outcome | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| AM1/dashboard-route | locally recompute a dashboard status fact | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| AM2/fixed-action | flatten a fixed stale-gate command into prose | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| AM2/prose | mark a non-invokable orchestration label executable | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| AM2/bytes | emit typed action metadata publicly | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |

