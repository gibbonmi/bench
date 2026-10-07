# Consume decisions for desired Link assets

Blocked by: 01-expose-complete-decisions-through-upgrade-preview.md
Writes: internal/adopt/decision.go, internal/adopt/link_transaction.go, internal/adopt/link_stage.go, internal/adopt/lifecycle_consumption_test.go (new), internal/adopt/lifecycle_special_test.go (new), internal/adopt/lifecycle_preview_test.go (new)
Covers: ALC06, ALC07, ALC08, ALC09, ALC10, ALC11, ALC12, ALC13, ALC14, ALC15, ALC19, ALC25, ALC33, ALC38

## What to build

Review chunk: ALC-LINK.

Link, confirmed Setup, applying Upgrade, and compatibility repair use complete decisions for ordinary desired assets and seeds.
Ticket 01 supplies the complete owner and the immutable execution input contract.
Its chunk must pass independent review before this consumer starts.
Later dropped-row and bespoke consumers remain on their current branches at this green checkpoint.

Observe staged and destination facts through the current helpers in the original validation order.
Submit immutable facts to PlanLifecycle and apply its selected effects and manifest dispositions.
Keep strict repair facts distinct from ordinary relink facts.
Remove the migrated caller classification instead of leaving a fallback or shadow decision.
Do not introduce an observation collector that validates the whole plan ahead of the current early refusals.

Preserve byte and mode convergence, canonical adapter identity, and modified or project-owned collisions.
Preserve absent seeds, present empty seeds, edited seeds, and dangling seed leaves.
Retain the rare old manifest row for a present planned seed.
Keep seed parent validation before its presence skip.
The current stage, report, lease, promotion, and repair-preimage owners remain authoritative.

Add missing characterization before the affected production branch moves.
Test real Link, Setup, applying Upgrade, and compatibility repair entries where their corresponding branch is available.
Use transactionalLink for the ordered competing-refusal fixture because that caller accepts the explicit ordered plan.
The fixture and exact winning diagnostic are fixed in the accepted spec's planned evidence section.
Keep all pre-existing assertions unchanged.

## Acceptance

- [ ] Real Link installs absent assets, updates clean old bytes, and reconciles a stale executable mode with identical bytes.
- [ ] Relink retains converged destination identity and modified managed bytes with the existing manifest and report disposition.
- [ ] An unowned collision remains project-owned, including an ordinary equal-byte file and a dangling leaf link.
- [ ] Canonical adapters through symlink parents pass, while foreign equal-byte targets retain the existing refusal.
- [ ] Confirmed Setup installs absent seeds without new ownership and preserves existing seeds with their current manifest disposition.
- [ ] A seed symlink parent and a required ordinary write symlink parent abort every accepted promotion.
- [ ] The earlier inline symlink-parent diagnostic wins over the later missing-source diagnostic with exit 1 and changed false.
- [ ] The competing-refusal mutant that omits the early return fails the exact stderr assertion despite publishing nothing.
- [ ] Strict compatibility repair retains mode guards, canonical manifest checks, hook checks, repair spans, preimages, and recovery results.
- [ ] Setup plan, decline, EOF, and ambiguous auto-confirm fixtures retain their baseline snapshots and exits.
- [ ] An owner add omission or preserve-to-change swap changes the corresponding real entry outcome and produces a recorded red.

## Verification

Run `bench test --package ./internal/adopt --run 'Test(LifecycleLinkPartitions|LifecycleSeedOwnership|LifecycleHardRefusalOrder|LifecycleRepairCompatibility|LifecycleSetupCancellation|TransactionalLinkAdoptsUnownedAdapterThroughSymlinkParent|SetupSeedsGateInputs)'`.
Run the existing compatibility, setup, link, and rollback package cases affected by these callers.
First run each new characterization against the frozen baseline, then against the candidate.
Use the existing syncDirectory seam for deterministic failure characterization; do not add another publisher.

Use bench probe for the named add omission, mode-fact omission, collision effect swap, seed effect swap, and strict-mode omission.
Also probe the early symlink-parent return omission with the competing-refusal fixture.
Record restored greens and the coordinator's distinct independent mutation.
Extend TestLifecycleOwnerConsumption with each migrated applying entry, but leave aggregate ALC01 and ALC35 ownership to ticket 04.

Reuse the current adopt, setup, and compatibility fixture owners.
Sweep duplicated facts and changed comments; keep every new or enlarged source within 400 lines.
No original assertion, public report shape, fixture registry, or transaction policy changes.
The desired-entry checks pass while ticket 03's dropped and bespoke classification remains unchanged.
