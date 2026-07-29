# Record the post-build re-measurement in the map

Blocked by: Enforce TEST bindings with structural and name refusals, Collapse
bench commit's reuse check into the gate home, Revise ADR 0002 posture 5 for
gate-execution reuse

## What to build

Story 15 of `specs/ft91-gate-fastpath/spec.md`: the post-build re-measurement
— solo canary, full gate on a changed tree, one unchanged-tree reuse timing —
recorded in `specs/ft91-gate-fastpath/decisions/gate-critical-path.md`
against the ≤60 s stop rule. The coordinator supplies the measured numbers;
this ticket transcribes them into the map where the levers were decided.

## Acceptance

- [ ] The map records all three measurements with the commands that produced
      them and states the stop-rule verdict (≤60 s met or missed).
- [ ] The out-of-scope lever re-checks (artifact-suite restructure, width-cap
      revival) are marked with their measured basis.
