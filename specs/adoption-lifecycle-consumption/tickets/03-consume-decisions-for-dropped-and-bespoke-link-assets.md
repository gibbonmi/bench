# Consume decisions for dropped and bespoke Link assets

Blocked by: 02-consume-decisions-for-desired-link-assets.md
Writes: internal/adopt/decision.go, internal/adopt/link_transaction.go, internal/adopt/link_stage.go, internal/adopt/marker.go, internal/adopt/link_hook.go, internal/adopt/lifecycle_consumption_test.go (new), internal/adopt/lifecycle_special_test.go (new)
Covers: ALC16, ALC17, ALC18, ALC21, ALC22, ALC23

## What to build

Review chunk: ALC-LINK.

Relink withdraws dropped assets and selects bespoke prepared replacements through the complete lifecycle owner.
Ticket 02 supplies ordinary desired-entry consumption and its exact observation path.
Complete this shared Link path without changing the transaction or report owner.
All applying Link, Setup, Upgrade, and compatibility repair routes then consume the full execution result.

Observe dropped rows once through the existing fingerprint helpers.
Move clean withdrawal, modified withdrawal, absent-row retirement, and ownership retention into PlanLifecycle.
Keep unreadable dropped assets as the existing hard failure.
Retain the separate conditional CLAUDE.md precedence.

Keep AGENTS.md marker parsing and hook inspection in their existing domain helpers.
Feed their facts and prepared replacements into the same complete decision.
Publish only replacements selected by that result.
The owner must not parse marker prose or reconstruct hook health.
Use marker.go and link_hook.go only for the preparation adapter changes required by this crossing.

Characterize known CLAUDE.md forms, project prose, project-owned empty files, and special instruction files before moving their branches.
Preserve the current no-read behavior for special files and all early refusal order.
Keep hook resolution and foreign-hook checks at their existing caller.
Do not leave duplicate effect or manifest selection in transactionalRepair after this migration.

## Acceptance

- [ ] Real relink removes a clean withdrawn asset and omits its manifest row.
- [ ] Relink retains a modified withdrawn asset and its recorded row with kept-modified-removed.
- [ ] An absent withdrawn path stays silent and loses ownership.
- [ ] Project-authored or empty CLAUDE.md remains unchanged without injected imports.
- [ ] Absent and known-form CLAUDE.md receive the current managed form with the baseline ownership disposition.
- [ ] Special instruction files retain their file kind and finish without opening a writerless FIFO.
- [ ] Link retains AGENTS.md project prose and the baseline effective-hook outcomes through existing domain helpers.
- [ ] A dropped remove omission or modified preserve-to-remove swap changes the real relink fixture and produces a recorded red.
- [ ] Omission of a selected bespoke replacement changes the real Link result while the unchanged domain helper remains correct.
- [ ] The shared transaction caller contains no surviving independent effect or ownership classification for migrated partitions.

## Verification

Run `bench test --package ./internal/adopt --run 'Test(LifecycleDroppedRows|LifecycleClaudeOwnership|LifecycleSpecialInstructionFiles|LifecycleHookCompatibility|LifecycleOwnerConsumption)'`.
Run the unchanged hook, link, setup, compatibility repair, and rollback cases that cover the shared transaction caller.
Use the current capability owner for platform-dependent special-file fixtures and run `bench test --check skip-ownership` if a test can skip.

Capture each characterization on the frozen baseline before its mechanical move, then run the same assertions on the candidate.
Use bench probe for clean-remove omission, modified-withdrawal effect swap, absent-row manifest-retain swap, and selected CLAUDE replacement omission.
Probe the special-file guard without introducing a new read policy.
Record the observed red, restored green, and coordinator's distinct independent mutation.

Extend the Link half of the AGENTS.md and hook cases; ticket 04 owns their final Link/Unlink coverage rows.
Re-run ticket 02's caller checks and ticket 01's preview checks at this green integration checkpoint.
Sweep duplicated facts and changed comments; keep every new or enlarged source within 400 lines.
Independent Standards, Spec, and Coverage review closes ALC-LINK before ticket 04 starts.
