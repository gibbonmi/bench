# Consume Unlink decisions and close real entry proof

Blocked by: 01-expose-complete-decisions-through-upgrade-preview.md, 02-consume-decisions-for-desired-link-assets.md, 03-consume-decisions-for-dropped-and-bespoke-link-assets.md
Writes: internal/adopt/decision.go, internal/adopt/unlink.go, internal/adopt/marker.go, internal/adopt/link_hook.go, internal/adopt/lifecycle_consumption_test.go (new), internal/adopt/lifecycle_special_test.go (new), internal/adopt/lifecycle_unlink_test.go (new), internal/systemtest/lifecycle_test.go (new)
Covers: ALC01, ALC20, ALC24, ALC26, ALC27, ALC28, ALC29, ALC30, ALC34, ALC35, ALC36, ALC37

## What to build

Review chunk: ALC-UNLINK.

Real Unlink and dry-run consume the same immutable lifecycle owner while retaining their existing removal and residual behavior.
Ticket 01 supplies the owner; tickets 02 and 03 supply all applying Link routes.
Their green checkpoints make the aggregate entry and recovery evidence executable in this ticket.
This final migration also closes the provisional ALC-EXIT verification at the ALC-UNLINK checkpoint.

Keep resolveInside and the existing observation helpers at Unlink's current caller.
Move the managed remove, preserve, absent, and manifest-residual classification into PlanLifecycle.
Treat an escape, symlink parent, special target, or fingerprint failure as the current reported residual.
Do not apply the legacy immutable path grammar as new execution validation.
A clean sub/../asset row still resolves to asset and removes it.

Keep AGENTS.md and hook domain transforms single-sourced with their facts and prepared replacements selected by the complete result.
Retain foreign-hook silent preservation and current domain refusals.
Preserve the separate data transaction, best-effort directory sweep, and final manifest publication.
Dry-run performs no stage, lease, or write.
Do not promise all-or-nothing Unlink rollback.
Remove the final independent managed classification copy from planUnlink.

Characterize missing unlink partitions before moving production branches.
Complete TestLifecycleOwnerConsumption through Link, confirmed Setup, applying Upgrade, Unlink, and compatibility repair.
Add the planned entry journey and publication tests through the existing system owner and selected binary.
BENCH_KIT identifies the candidate kit; all system checks run through bench test --check system.
Reuse owner.runSelected and owner.runAt instead of constructing a new command runner or environment owner.

## Acceptance

- [ ] Real Unlink removes clean present rows, keeps absent rows silent, and accepts safely normalized manifest paths.
- [ ] Modified assets and unsafe rows retain their baseline residual report, exit, and manifest behavior.
- [ ] Dry-run leaves bytes, modes, symlink targets, directory entries, and manifest bytes unchanged.
- [ ] The manifest remains whenever a modified or refused residual remains, and follows the existing final-publication sequence.
- [ ] Link and Unlink retain AGENTS.md project prose and each effective-hook state without a second domain policy.
- [ ] Every applying entry changes its actual fixture when the corresponding selected owner effect is omitted or swapped.
- [ ] Publication failures preserve baseline results for missing source, stage failure, report failure, promotion failure, and final manifest failure.
- [ ] The existing interruption seam preserves fresh-process recovery at pre-publication, middle-promotion, and final-manifest boundaries.
- [ ] A deterministic restore omission leaves a changed preimage and turns the recovery test red; restored source passes.
- [ ] The same characterization passes on frozen baseline and candidate with all original assertions and outcomes unchanged.
- [ ] The final caller census finds no duplicate managed effect or ownership derivation in transactionalRepair or planUnlink.

## Verification

Run `bench test --package ./internal/adopt` for the complete adopt package, including every original test.
Run `bench test --check system` with the selected candidate binary and BENCH_KIT supplied by the existing system owner.
Preserve the installed-wrapper journey, fresh-process repair, seed, read-only unlink, and destination-sync rollback assertions.
Use the existing syncDirectory and BENCH_LINK_FAULT seams; keep their environments private to their current test owners.

Use bench probe for clean-remove omission, modified preserve-to-remove swap, residual-manifest guard omission, and dry-run publisher-guard omission.
Record one selected-effect witness for each real applying entry and the restore omission witness.
The coordinator chooses a different mutation kind and site for independent proof.
No ordinary test rewrites repository source, and no new fixture harness or payload registry is introduced.

Record exact reports, exits, destination snapshots, manifest rows, observed reds, and restored greens.
Normalize only fixture-root paths and existing volatile recovery identifiers.
Run the whole-tree caller census and duplicated-facts sweep; review changed comments against craft-comments.
Keep each new or enlarged source within 400 lines.
Independent Standards, Spec, and Coverage review closes ALC-UNLINK before final integration acceptance.
The orchestrator records final verification and obtains the whole-project gate through the required landing, not a manual benchmark.
