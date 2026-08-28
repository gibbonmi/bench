# Add the proven-target field to the release plan

Blocked by: none
Writes: scripts/release-plan.json, scripts/release-plan.mjs, internal/releaseevidence, internal/conformance

## What to build

The release plan states, per target, whether that target carries a native proof. Each target row gains a required `native_proof` boolean, and the two Darwin rows set it to false.

`readReleasePlan` in `release-plan.mjs` rejects a target that omits the field, a target that gives a non-boolean, and a plan whose proven set is empty. `release-plan.mjs` gains `proof-targets`, `proof-matrix-json`, and `proof-target`, which mirror `targets`, `matrix-json`, and `target` over the proven list only. One internal helper produces both views from a single filter, so the five-field projection is written once and `native_proof` never reaches a matrix variable.

The Go plan struct in `internal/releaseevidence` gains the same field in this ticket. That package decodes `normalized-json` with unknown fields disallowed, so the plan change alone would redden it. The Go consumer of the field lands in `count-release-evidence-against-proven-targets.md`.

No other consumer changes, so the workflow and the shell scripts keep their present behavior and the tree stays green.

## Acceptance

- [ ] `readReleasePlan` rejects a plan whose target omits `native_proof` (row B1).
- [ ] `readReleasePlan` rejects a plan whose target gives a non-boolean `native_proof` (row B1).
- [ ] `readReleasePlan` rejects a plan whose proven set is empty (row B2).
- [ ] `proof-matrix-json` omits every unproven target (row B3).
- [ ] `matrix-json` still lists every shipped target, and no entry carries `native_proof` (row B3).
- [ ] A plan that marks a Darwin target proven puts that target in the proven views (row B12).
- [ ] The `internal/releaseevidence` package still decodes the plan, and its tests pass unchanged.
