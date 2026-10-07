# Expose complete decisions through Upgrade preview

Blocked by: none
Writes: internal/adopt/decision.go, internal/adopt/upgrade.go, internal/adopt/lifecycle_decision_test.go (new), internal/adopt/lifecycle_preview_test.go (new)
Covers: ALC02, ALC03, ALC04, ALC05, ALC31, ALC32

## What to build

Review chunk: ALC-OWNER.

Upgrade preview receives the same count table through the additive complete decision owner.
This ticket delivers a real Upgrade result while later execution consumers remain unchanged.
Read the FT217 definitions and current decision, preview, and version guards before changing them.
Reconfirm the accepted spec's source premises against this ticket's predecessor.

Add Decisions and the immutable Execution inventory beside the existing input contract.
Implement the complete owner partitions from the accepted spec without filesystem reads or policy callbacks.
Derive Operations from add, change, and remove decisions only.
Keep nil Execution distinct from a present empty inventory.
Keep legacy validation, Preserve precedence, and the Current/Desired union unchanged.
Execution facts must support the later consumers without replacing their hard validation.

Keep upgradePlanCounts on its existing declared inventory and sentinel policy.
Real Upgrade produces file, inline, and adapter entries through buildLinkPlan.
Test inline-exec, seed, and unreadable-source projections directly through the actual upgradePlanCounts seam.
Confirmed Setup covers its real inline-exec and seed producers in ticket 02.

Use the owner's effect projection for the real Upgrade count table.
Preserve the effective-hook refresh adjustment and all version guards.
Do not replace preview sentinels with staged execution fingerprints.

Before the projection move, add the preview and version characterization to the unchanged Upgrade entry.
Capture the frozen baseline once and run the same tests against the candidate.
Add the complete-decision and ordering tests at PlanLifecycle's production seam.
Keep TestPlanLifecycleFromImmutableInventories and every other original assertion unchanged.

## Acceptance

- [ ] Real Upgrade preview returns baseline counts for reachable file, adapter, and inline entries.
- [ ] The actual upgradePlanCounts seam retains all five legacy kind projections and the unreadable source sentinel.
- [ ] Real Upgrade retains the count table for owned, unowned, dropped-row, preserved CLAUDE.md, and hook-refresh fixtures.
- [ ] Real Upgrade calls retain equal-version, newer, older, forced-older, absent-manifest, empty-manifest, and missing-header outcomes.
- [ ] Complete decisions include preserve reasons and manifest dispositions while Operations retain the existing exact effect-only list.
- [ ] Reordered immutable inputs yield identical path-ordered decisions, and all current invalid-input refusals remain effective.
- [ ] The independent preserve expectation fails when its owner row is omitted, and the restored source passes.
- [ ] A synthetic-to-exact non-file preview swap changes a named real Upgrade baseline count and turns the characterization red.

## Verification

Run `bench test --package ./internal/adopt --run 'Test(PlanLifecycleFromImmutableInventories|CompleteLifecycleDecisions|LifecycleDecisionOrder|LifecycleUpgradePreviewCompatibility|LifecycleUpgradeVersionGuards)'`.
Record executed test counts, zero unexpected skips, baseline and candidate results, and exact reports where promised.
Use bench probe for the preserve-row omission, owner-order omission, and preview projection swap.
Record each observed red and restored green rather than treating a planned witness as evidence.
The coordinator selects an independent mutation at a different site and of a different kind.

Run a duplicated-facts sweep over this ticket's delta and review its comments against craft-comments.
Read the existing fixture owners and reuse their payload sources without a second registry or copied expected manifest.
Keep each new or enlarged source file within 400 lines.
If headroom requires an extra path, obtain an approved plan expansion before dispatch or publication.

The producer precision correction is recorded in assets/spec-review.md and must receive independent Spec and Coverage review.
This ticket can pass its focused checks while Link and Unlink still use their current execution paths.
Independent Standards, Spec, and Coverage review closes ALC-OWNER before ticket 02 starts.
