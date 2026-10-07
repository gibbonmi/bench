# Adoption lifecycle consumption

Status: staged

Roadmap: FT217

Decision source: ready compiled map, `decisions/architecture-adoption.md`, with its topic tickets in this spec folder

Verification log: spec round 1 held one Coverage finding; repair round 1 accepted before slicing; producer clarification checked in ticket round 1; tickets accepted

## Problem

The adoption commands preserve project content through different classification paths.
The current lifecycle decision supplies upgrade counts, but it does not select real filesystem effects.
The source pins below identify these paths.
This separation lets an execution caller retain its own classification after the decision changes.

## Solution

Deepen `PlanLifecycle` as the single owner of the managed lifecycle decision.
Make link, setup, upgrade execution, and unlink consume its complete decision.
Keep the current command behavior, reports, validation order, and transaction owner.
Derive the existing effect-only `Operations` view from the complete decision.
The change serves every repository that links the kit.

This spec is planning work.
The author does not implement production code, run a test suite, or perform a manual benchmark.
The coordinator obtains independent spec review before the author creates ticket slices.
Ticket review follows the slices.
Implementation requires a later approval and the eligible commitment start.
This spec does not change the active commitment.

## User stories

Line: gpt-5.6-sol / high

Implementation-line reason: ALC-LINK carries the hardest policy migration. The source fixes behavior, but caller coverage needs new characterization before each move.
Harder chunks: ALC-LINK, ALC-UNLINK
Author line: gpt-6.1-sol / high / one draft and at most two independent-review repair rounds
Routing source: `projects/benchkit.md`, Lines, and `.bench/lines.env` bind the Codex mid tier to gpt-5.6-sol
Routing disposition: proposed implementation line only. The user-selected author line remains gpt-6.1-sol/high

### Complete decision

1. As a maintainer, I want one immutable lifecycle decision, so that command paths cannot retain separate asset policy.
2. As an existing decision caller, I want the same effect-only operations, so that the refactor preserves the existing interface.
3. As a maintainer, I want explicit preserve decisions, so that a complete plan accounts for retained paths.
4. As a maintainer, I want deterministic decisions, so that equivalent inventory inputs give the same result.
5. As a caller, I want the same invalid-input refusals, so that existing inventory guards remain effective.

### Link and setup

6. As an adopter, I want absent assets installed, so that my repository receives the kit payload.
7. As an adopter, I want clean old assets updated, so that a new kit release reaches my repository.
8. As an adopter, I want converged managed assets left in place, so that relink avoids a redundant write.
9. As an adopter, I want executable modes reconciled, so that a byte-identical stale mode does not block an update.
10. As an adopter, I want canonical adapters accepted through a symlink parent, so that existing adapter directories remain usable.
11. As an adopter, I want foreign adapter targets refused, so that equal bytes do not authorize a different target.
12. As a project owner, I want modified managed assets preserved, so that relink keeps my edits.
13. As a project owner, I want unowned collisions preserved, so that link keeps my content.
14. As a project owner, I want absent seeds installed once, so that setup supplies initial project content.
15. As a project owner, I want existing seeds preserved without ownership, so that setup leaves my content under my control.
16. As an adopter, I want clean dropped assets removed, so that relink withdraws obsolete kit content.
17. As a project owner, I want modified dropped assets retained, so that relink preserves edits to withdrawn assets.
18. As an adopter, I want absent dropped rows retired, so that a missing asset does not remain falsely owned.
19. As an adopter, I want unsafe write paths refused before promotion, so that a partial plan cannot publish through a symlink parent.

### Bespoke content and repair

20. As a project owner, I want my AGENTS.md prose preserved, so that adoption updates only the managed block.
21. As a project owner, I want project CLAUDE.md content preserved, so that adoption does not inject imports into my content.
22. As an adopter, I want known CLAUDE.md forms reconciled, so that adoption installs the current managed form.
23. As a project owner, I want special instruction files preserved, so that adoption does not open a writerless FIFO.
24. As an adopter, I want effective hook behavior preserved, so that lifecycle policy does not change hook ownership.
25. As an operator, I want strict compatibility repair preserved, so that managed repair keeps its existing safety conditions.

### Unlink and preview

26. As an adopter, I want clean manifest assets removed, so that unlink reverses the managed installation.
27. As a project owner, I want modified assets retained during unlink, so that removal keeps my edits.
28. As an adopter, I want unsafe manifest rows reported, so that unlink preserves paths it cannot safely remove.
29. As an adopter, I want unlink rehearsal to write nothing, so that I can inspect the removal result safely.
30. As an adopter, I want the manifest retained for residuals, so that I can complete a partial unlink later.
31. As an adopter, I want the existing upgrade counts, so that a preserving refactor does not change preview output.
32. As an adopter, I want existing version guards, so that equal-version and downgrade behavior remains stable.
33. As an adopter, I want setup confirmation behavior preserved, so that planning and cancellation keep their current write limits.

### Transaction and exit proof

34. As an adopter, I want publication failure to preserve recovery behavior, so that the refactor cannot strand a partial installation.
35. As a maintainer, I want real command evidence, so that helper-only tests cannot conceal a caller that bypasses the decision.
36. As a maintainer, I want the existing assertions unchanged, so that this refactor cannot redefine adoption behavior.
37. As a maintainer, I want one ordered mechanical migration per checkpoint, so that each green slice has a bounded cause.

## Implementation decisions

### Source pins and factual questions

Source tip: a382f4848d432815968f044820fad593e856e297
Read date: 2026-10-06
Invalidation trigger: any change to the named lifecycle owner, caller, observation helper, or pinned test before implementation
Evidence status: source verified. Runnable compatibility evidence remains planned

The factual questions form one dependency graph.
First, identify the existing decision contract and its readers.
Then, trace each execution partition and its filesystem observations.
Finally, identify the smallest owner interface that preserves both execution and preview contracts.
The compiled map and FT217 answer the scope question.
The following source pins answer the code questions.

| fact | current source |
|---|---|
| The decision compares immutable inventories | `internal/adopt/decision.go:36` |
| Operations omit preserved CLAUDE.md | `internal/adopt/decision_test.go:14` |
| Link delegates application | `internal/adopt/link.go:64` |
| Setup adds a managed gate and two seeds | `internal/adopt/setup.go:310` |
| Link and repair classify desired entries | `internal/adopt/link_transaction.go:97` |
| Link and repair reconcile dropped rows | `internal/adopt/link_transaction.go:190` |
| Unlink classifies manifest rows | `internal/adopt/unlink.go:117` |
| Upgrade uses synthetic non-file fingerprints | `internal/adopt/upgrade.go:127` |
| Upgrade applies through the link transaction | `internal/adopt/upgrade.go:78` |
| Compatibility repair shares the transaction | `internal/adopt/compatibility.go:264` |
| Fingerprints omit regular-file permissions | `internal/adopt/manifest.go:107` |
| Convergence compares permissions separately | `internal/adopt/link_stage.go:202` |
| The publication adapter retains observed identities | `internal/adopt/transaction.go:211` |
| The transaction checks observations before publication | `internal/adopt/transaction/transaction.go:93` |

The author re-read every structured source in the compiled map.
The author also re-read both topic tickets and the authoritative FT217 contract.
The FT217 historical source is `capture/architecture-review-20260817T104714.html:160`.
The current caller census confirms its lifecycle concern.
The later historical recommendation at `capture/architecture-review-20260905T112705.html:334` proposes a narrower outcome.
It does not replace FT217 or the ready compiled map.

### One owner with additive projections

Keep the existing `PlanLifecycle(LifecycleInput) LifecyclePlan` interface.
Keep `Asset`, `Operation`, `Current`, `Desired`, `Preserve`, `Operations`, and `Refusal` compatible with their current callers.
Add a complete `Decisions` projection to `LifecyclePlan`.
The owner produces `Decisions` first.
The owner derives `Operations` by selecting only add, change, and remove decisions.
It never adds preserve entries to `Operations`.

ALC02 and ALC03 pin this contract.

Extend `LifecycleInput` with an optional immutable execution inventory.
The proposed field is `Execution *LifecycleInventory`.
That inventory holds `Action` and `Observations` fields.
`Action` selects link or unlink.
`Observations` holds immutable `LifecycleObservation` records.

A nil execution inventory selects the existing inventory-comparison contract.
A present, empty execution inventory selects an empty execution plan.

Keep that distinction explicit.
Do not add another planner or a caller-selected callback for classification.
Private type names can vary within the ownership fence.

The execution inventory carries an action of link or unlink and one observation for each relevant path.
The observation carries facts, not a preselected lifecycle result.
Its fields cover these facts:

- The stable path key and the original caller order.
- The plan kind and the presence of a desired entry.
- The destination presence and file kind.
- The recorded manifest fingerprint and the observed destination fingerprint.
- The staged fingerprint, file kind, permissions, and symlink target.
- The resolved content equivalence for staged symlink convergence.
- The canonical adapter target identity.
- The presence of a symlink parent.
- The strict repair mode comparison.
- The known CLAUDE.md form and the presence of a special instruction file.

The caller observes filesystem facts through the existing helpers.
The owner decides add, change, remove, or preserve from those immutable facts.
The owner reads no filesystem, invokes no Git command, and publishes nothing.
Tests attach to the same input and output that production callers use.
ALC01, ALC04, and ALC35 grade these properties.

Each complete decision carries its path, kind, preserve reason, and manifest disposition.
The proposed record fields are `Path`, `Kind`, `Reason`, `ManifestAction`, `Fingerprint`, and `Order`.
The manifest disposition is omit, retain the recorded fingerprint, or record the observed or staged fingerprint.
An ordinary seed adds no managed manifest row.
A preserve reason distinguishes convergence, seed ownership, project collision, modified ownership, modified withdrawal, and unlink refusal.

The legacy inventory view emits unchanged and requested-preserve reasons for its no-effect paths.
It preserves the current precedence of `Preserve`: that input suppresses removal, but does not suppress a desired change.
It retains the union of Current and Desired as its path inventory.
A Preserve-only path remains a validated input without a decision row.
The caller translates these decisions into its existing report vocabulary.
It does not recompute ownership or compare fingerprints to select a different effect.

The owner orders the complete decision by stable path key.
The caller retains the existing publication sequence through the original caller order.
The caller retains report ordering where the current report defines it.
This change does not impose a new transaction sequence.
ALC04 and ALC34 detect an ordering change that alters existing behavior.

### Execution partitions and precedence

The table states the current precedence that the owner must consume.
The caller retains all existing validation and hard refusals.
A hard refusal prevents promotion of every accepted entry.
The table is a classification contract, not permission to add new validation.

| fixture facts | winning rule | complete decision | manifest disposition | rows |
|---|---|---|---|---|
| An absent desired ordinary asset | Install the absent path | add | staged fingerprint | ALC06 |
| An owned path matches the old hash but differs from staged bytes | Converge a clean old asset | change | staged fingerprint | ALC07 |
| An owned path matches staged content and relevant mode | Keep the converged path | preserve, converged | observed fingerprint | ALC08 |
| An owned regular file has equal bytes but a stale executable mode | Mode prevents convergence | change | staged fingerprint | ALC09 |
| An unowned adapter resolves through its parent to the canonical target | Canonical identity permits convergence | preserve, converged | observed fingerprint | ALC10 |
| An unowned adapter resolves to a foreign file with equal bytes | Canonical identity fails | caller hard refusal | no publication | ALC11 |
| A managed desired path matches neither old nor staged content | Preserve a modified asset | preserve, modified-managed | recorded fingerprint | ALC12 |
| An unowned desired path exists, including equal ordinary bytes | Preserve a project collision | preserve, project-owned | omit | ALC13 |
| An absent seed has no symlink parent | Seed once | add | omit | ALC14 |
| A seed exists, including an empty file or dangling leaf link | Preserve the project seed | preserve, seed-owned | omit unless an old row already exists | ALC15 |
| A dropped ordinary row is present and matches its recorded hash | Withdraw a clean asset | remove | omit | ALC16 |
| A dropped ordinary row is present and differs from its recorded hash | Preserve the modified withdrawal | preserve, kept-modified-removed | recorded fingerprint | ALC17 |
| A dropped ordinary row is absent | Retire the absent row | preserve, absent | omit | ALC18 |
| A desired write requires a symlink parent | The hard write guard wins | caller hard refusal | no publication | ALC19 |
| A seed has a symlink parent, even when its leaf exists | The seed parent guard wins before the presence skip | caller hard refusal | no publication | ALC19 |
| A CLAUDE.md row has modified content outside known forms | Preserve the conditional owner | preserve, modified-managed | recorded fingerprint | ALC21 |
| An absent or known-form CLAUDE.md has a staged current form | Reconcile the conditional owner | add or change as current staging requires | staged fingerprint | ALC22 |
| An unlink row is present and matches its recorded hash | Remove the clean path | remove | caller retains manifest until final step | ALC26 |
| An unlink row is present and differs from its recorded hash | Keep modified content | preserve, modified | caller retains manifest | ALC27 |
| An unlink row escapes, has a symlink parent, or cannot produce a fingerprint | Keep the unsafe path | preserve, refused | caller retains manifest | ALC28 |
| An unlink row is already absent | Keep the current silent absence | preserve, absent | caller retires manifest if no residual remains | ALC26 |

The current seed branch does not add a manifest row.
The dropped-row loop nevertheless retains a previously owned planned seed row.
Preserve this rare combination, rather than clearing its old ownership during the refactor.
The source is `internal/adopt/link_transaction.go:201`.
ALC15 includes this fixture.

An unreadable dropped link asset remains a hard failure where the current caller returns an error.
An unlink fingerprint failure remains a reported residual where the current caller continues.
The execution inventory must retain this distinction.
It must not turn every invalid path into the legacy decision's whole-plan refusal.
ALC28 and ALC34 include both call paths.

### Bespoke content and safety

Keep AGENTS.md marker transforms in their existing domain helpers.
Keep hook inspection, effective hook resolution, and hook rendering in their existing domain helpers.
Their facts and prepared replacements join the same complete execution result.
A prepared replacement identifies a domain transform, not an independently classified managed asset.
The caller publishes only the replacement selected by the complete result.
The owner does not parse marker prose or infer hook health again.

ALC20, ALC23, and ALC24 cover these domain contracts.

Keep the CLAUDE.md known-form observation at its existing byte source.
Move the lifecycle choice and ownership retention into the decision owner.
Keep the current empty-file distinction.
A project-owned empty CLAUDE.md remains untouched.
Do not apply the general dropped-row rule to a conditional CLAUDE.md row.
ALC21 and ALC22 cover these cases.

Keep hard checks at their current callers and in their current order.
These checks include source availability, parent type, special-file safety, marker validation, hook resolution, and the foreign-hook refusal.
They also include transaction leases and observation rechecks.
Preserve the existing distinction between a hard exit 1 and a partial exit 3.
ALC19, ALC23, ALC24, ALC25, ALC28, and ALC34 grade this contract.

Compatibility repair remains an indirect consumer of the same link decision.
Retain strict managed-mode checks, canonical manifest validation, hook validation, repair spans, and retained preimages.
The decision owner must not weaken repair into ordinary relink.
ALC25 covers this composed caller.

### Preview compatibility

Upgrade preview keeps its current immutable input policy.
A file entry uses its kit-side fingerprint or the current unreadable sentinel.
A non-file entry uses its recorded fingerprint or its kind-specific new sentinel.
The preserved CLAUDE.md row remains outside effect-only operations.
The existing effective-hook refresh adjustment remains at the upgrade caller.

ALC31 pins real Upgrade counts for reachable file, inline, and adapter entries.
It pins all five declared kind projections and the unreadable sentinel through the actual upgradePlanCounts seam.

Execution uses exact destination and staged observations.
Preview uses the legacy declared inventory view.
Both views consume `PlanLifecycle` and its effect projection.
The owner never substitutes preview sentinels for an execution observation.
The spec does not promise that preview counts equal actual writes.
Such equality would change the behavior fixed by the source contract.

Setup preview retains its fact report and confirmation flow.
It does not publish an execution decision before confirmation.
The confirmed setup plan joins the same link execution path with seeds and its executable gate.
Unlink dry-run feeds the same observed facts to the owner without a lease, stage, or publication.
ALC29 and ALC33 cover these read-only paths.

### Current reader and writer census

The author searched the whole tree with hidden paths and without `.git/`.
The author also ran `bench consumers` for each changed function below.
The source tip above pins these results.
No schema, public verb, public flag, manifest format, or rendered report changes.

| reader or writer | direct helpers that read its lifecycle facts | disposition |
|---|---|---|
| `upgradePlanCounts` in `internal/adopt/upgrade.go:121` | `Manifest.Rows`, `Manifest.Hash`, `fingerprintPath`, `PlanLifecycle` | ALC31, legacy input projection |
| `transactionalRepair` in `internal/adopt/link_transaction.go:30` | `ReadManifest`, `stagePlanEntry`, `hasSymlinkParent`, `convergedFingerprint`, `sameAdapterTarget`, `ownedUnmodified`, `reclaimableClaude` | ALC06–ALC25, exact observation and decision consumption |
| `planUnlink` in `internal/adopt/unlink.go:98` | `resolveInside`, `hasSymlinkParent`, `fingerprintPath`, `stripAgentsForUnlink`, `removeManagedHook`, `sweepEmptyDirs` | ALC26–ALC30, exact observation and decision consumption |
| `Link` in `internal/adopt/link.go:28` | `buildLinkPlan`, `transactionalLink` | ALC35, unchanged command entry |
| `convergeSetup` in `internal/adopt/setup.go:304` | `buildLinkPlan`, seed producers, `transactionalLink`, `finishSetup` | ALC14, ALC15, ALC33, ALC35 |
| `Upgrade` in `internal/adopt/upgrade.go:27` | `pinnedManifest`, `buildLinkPlan`, `compareKitVersions`, `upgradePlanCounts`, `prePushRefreshEligible`, `transactionalLink` | ALC31, ALC32, ALC35 |
| `Unlink` in `internal/adopt/unlink.go:26` | `ReadManifest`, `planUnlink`, `writeUnlinkReport` | ALC26–ALC30, ALC35 |
| `mutateCompatibility` in `internal/adopt/compatibility.go:224` | `buildLinkPlan`, `transactionalRepair`, `transaction.Store.Undo` | ALC25, indirect caller protection |
| `promoteWithLease` in `internal/adopt/transaction.go:76` | `observedChange`, `transaction.Publish`, `transaction.Store.Apply` | ALC34, unchanged publisher |
| `doctorFixRun` and hook repair callers | Existing hook and AGENTS.md domain helpers | Excluded policy expansion, ALC25 protects shared helpers |
| `compatibilityAssets` and `buildLinkPlan` | Payload rows, destination mapping, source walk | Unchanged payload facts, ALC06 and ALC35 protect composition |
| `cmd/bench` adoption routing | Public adoption entries | Unchanged dispatcher, ALC35 reaches real selected command routes |

The existing test callers of `PlanLifecycle` all sit in `TestPlanLifecycleFromImmutableInventories`.
The existing transaction test callers are `TestTransactionalLinkAdoptsUnownedAdapterThroughSymlinkParent` and `TestLinkRefusesUnresolvedHooksDirectory`.
`planUnlink` has only the Unlink production caller.
`upgradePlanCounts` has only the Upgrade production caller.
`transactionalRepair` has the transactionalLink wrapper and mutateCompatibility callers.
`transactionalLink` has Link, convergeSetup, Upgrade, and the two test callers above.

The complete result becomes the one writer of managed effect and ownership choices.
The manifest renderer remains the writer of manifest bytes.
The transaction package remains the writer of filesystem publication and recovery state.
The report functions remain the writers of command reports.
These distinct facts require no competing writer protocol or new orchestrator.

## Implementation chunks

The independent spec review accepted the frozen source before ticket slicing.
The accepted spec hash is a12db867b502df92f3e68ae49ebe50d7b6a929c25f7d13d9ec79bd4a3fa57a3d.
The [planning review record](assets/spec-review.md) preserves that acceptance and its independent confirmation.
The four tickets use expand, migrate, then contract order.
Each chunk receives Standards, Spec, and Coverage review before its successor.
A ticket remains a serial green checkpoint on the retained integration source.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| ALC-OWNER / 01-expose-complete-decisions-through-upgrade-preview.md | Complete decisions reach real Upgrade preview through the legacy effect projection | ALC02–ALC05, ALC31, ALC32 | Immutable decision, real preview, and version characterization | no |
| ALC-LINK / 02-consume-decisions-for-desired-link-assets.md, 03-consume-decisions-for-dropped-and-bespoke-link-assets.md | All applying Link routes consume exact selected effects and ownership | ALC06–ALC19, ALC21–ALC23, ALC25, ALC33, ALC38 | Desired, seed, repair, ordered refusal, dropped, and bespoke caller tests | yes |
| ALC-UNLINK / 04-consume-unlink-decisions-and-close-entry-proof.md | Unlink joins the owner and closes aggregate entry and recovery proof | ALC01, ALC20, ALC24, ALC26–ALC30, ALC34–ALC37 | Unlink partitions, all-entry consumption, differential characterization, package and system checks | yes |

The provisional ALC-OWNER remains ALC-OWNER and gains the existing preview consumer.
The provisional ALC-LINK remains ALC-LINK with two serial tickets.
The provisional ALC-UNLINK remains ALC-UNLINK.
The provisional ALC-EXIT joins the final ALC-UNLINK checkpoint as aggregate verification.
It creates no duplicate ticket membership or helper-only checkpoint.

The first chunk's seam review closes before any consumer starts.
The Link chunk review closes before Unlink starts.
The final chunk review precedes final acceptance and landing.

ALC01 and ALC35 belong to the final consumer's checkpoint.
Earlier tickets build and prove their own selected-effect crossings there.
ALC20 and ALC24 also belong to the final Link/Unlink checkpoint.
The later owner does not defer an earlier ticket's independently useful result.

Before the first production move, record green coverage for the affected behavior.
Where coverage is absent, add characterization while the production path remains unchanged.
Run each new characterization against the frozen baseline and the candidate.
Then make one ordered mechanical move.
A required change to an old assertion stops this refactor.
Route that behavior delta through a separate reviewed defect or feature outcome.

## Testing decisions

Use the existing immutable `PlanLifecycle` seam for exhaustive partitions.
Use real Link, Setup, Upgrade, and Unlink entries for caller composition.
Use the existing transaction seam for deterministic sync and publication failure.
Use the existing system owner for any fresh-process or installed-wrapper case.
Do not construct a second wrapper, gate runner, or build fixture owner.

The current decision test pins effect-only compatibility.
The current unlink tests in `internal/adopt/link_transaction_test.go` pin removal, preservation, and read-only behavior.
The current adapter test in that file pins canonical identity through a symlink parent.
The current hook tests pin failed resolution before writes.
The current setup tests pin seed ownership.
The current sync-failure test pins restore survival.

These tests remain unchanged.

New tests live in separate files inside the existing packages.
The planned names below are proposals for the implementation, not claims of executed evidence.
Every independent expected classification names a required omission or swap witness.
The implementation records the observed red through `bench probe` and the subsequent green result.
The coordinator uses a different mutation kind and site for its independent proof.
No production policy or fixture registry may appear twice.

### Seam diagram

```text
Link / confirmed Setup / applying Upgrade / compatibility repair
    -> existing observations and staged payload
    -> PlanLifecycle execution inventory
    -> complete Decisions -> effect and manifest projections
    -> existing reports and existing leased transaction publisher

Upgrade preview -> legacy declared inventory -> PlanLifecycle -> Operations -> current count table
Unlink / dry-run -> existing observations -> PlanLifecycle -> Decisions -> current report
                                              -> existing publisher only for applying unlink
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| ALC01 | 1 | Real consumers obtain managed effect choices from one immutable owner. | planned internal/adopt/lifecycle_consumption_test.go (`TestLifecycleOwnerConsumption`) | A preserve-to-change swap in the owner changes the real command outcome. |
| ALC02 | 2 | Operations retain the existing add/change/remove list without CLAUDE.md. | `internal/adopt/decision_test.go` (`TestPlanLifecycleFromImmutableInventories`) | An added preserve operation fails the unchanged exact comparison. |
| ALC03 | 3 | Decisions include preserved CLAUDE.md with its reason and ownership disposition. | planned internal/adopt/lifecycle_decision_test.go (`TestCompleteLifecycleDecisions`) | Omission of the preserve decision fails the independent expected row. |
| ALC04 | 4 | Reordered immutable inputs yield the same path-ordered complete decision. | planned internal/adopt/lifecycle_decision_test.go (`TestLifecycleDecisionOrder`) | Omission of owner order exposes unequal results from reordered inputs. |
| ALC05 | 5 | Legacy duplicate paths, hostile paths, and missing fingerprints retain their refusals. | `internal/adopt/decision_test.go` (`TestPlanLifecycleFromImmutableInventories`) | Omission of the relevant guard accepts an existing refused fixture. |
| ALC06 | 6 | Link installs an absent ordinary asset with its staged fingerprint. | planned internal/adopt/lifecycle_consumption_test.go (`TestLifecycleLinkPartitions`, absent fixture) | Omission of the owner's add result leaves the named destination absent. |
| ALC07 | 7 | Relink replaces a clean old asset with current kit bytes. | planned internal/adopt/lifecycle_consumption_test.go (`TestLifecycleLinkPartitions`, clean-old fixture) | A change-to-preserve swap leaves the old bytes in place. |
| ALC08 | 8 | Relink preserves the identity of a converged managed destination. | planned internal/adopt/lifecycle_consumption_test.go (`TestLifecycleLinkPartitions`, converged fixture) and planned internal/adopt/lifecycle_decision_test.go (`TestCompleteLifecycleDecisions`) | A preserve-to-change swap fails the exact complete-decision predicate. |
| ALC09 | 9 | Relink applies a changed executable mode with identical file bytes. | planned internal/adopt/lifecycle_consumption_test.go (`TestLifecycleLinkPartitions`, mode fixture) | Omission of the mode fact leaves the destination mode stale. |
| ALC10 | 10 | Link accepts a canonical adapter through a symlink parent without replacing its target. | `internal/adopt/link_transaction_test.go` (`TestTransactionalLinkAdoptsUnownedAdapterThroughSymlinkParent`) | Omission of the canonical convergence branch refuses the existing converged fixture. |
| ALC11 | 11 | Link refuses a foreign adapter target even when its bytes match. | `internal/adopt/link_transaction_test.go` (`TestTransactionalLinkAdoptsUnownedAdapterThroughSymlinkParent`) | A canonical-identity-to-content swap accepts the foreign-identical fixture. |
| ALC12 | 12 | Relink retains modified managed bytes and reports modified-managed. | planned internal/adopt/lifecycle_consumption_test.go (`TestLifecycleLinkPartitions`, modified fixture) | A preserve-to-change swap overwrites the project edit. |
| ALC13 | 13 | Link retains an unowned collision without adding a manifest row. | planned internal/adopt/lifecycle_consumption_test.go (`TestLifecycleLinkPartitions`, project collision fixture) | A project-preserve-to-add swap overwrites or owns the collision. |
| ALC14 | 14 | Setup installs an absent seed without a manifest row. | `internal/adopt/setup_test.go` (`TestSetupSeedsGateInputs`) and planned internal/adopt/lifecycle_special_test.go (`TestLifecycleSeedOwnership`) | Omission of seed add leaves gate inputs absent. |
| ALC15 | 15 | Setup retains an existing seed with its existing manifest disposition. | planned internal/adopt/lifecycle_special_test.go (`TestLifecycleSeedOwnership`) | A seed-preserve-to-change swap overwrites empty, edited, or dangling-leaf fixtures. |
| ALC16 | 16 | Relink removes a clean dropped asset. | planned internal/adopt/lifecycle_consumption_test.go (`TestLifecycleDroppedRows`, clean fixture) | Omission of remove leaves the withdrawn file present. |
| ALC17 | 17 | Relink retains a modified dropped asset with kept-modified-removed. | planned internal/adopt/lifecycle_consumption_test.go (`TestLifecycleDroppedRows`, modified fixture) | A preserve-to-remove swap deletes the edited withdrawn file. |
| ALC18 | 18 | Relink drops ownership for an already absent withdrawn asset. | planned internal/adopt/lifecycle_consumption_test.go (`TestLifecycleDroppedRows`, absent fixture) | A manifest-omit-to-retain swap leaves the old row. |
| ALC19 | 19 | A required write through a symlink parent aborts all planned promotion. | planned internal/adopt/lifecycle_special_test.go (`TestLifecycleHardRefusalOrder`) | Omission of the parent guard publishes a preceding accepted entry. |
| ALC38 | 19 | The earlier inline symlink-parent refusal wins over the later missing-source refusal. | planned internal/adopt/lifecycle_special_test.go (`TestLifecycleHardRefusalOrder`, competing-refusals fixture through transactionalLink) | Omission of the early symlink-parent return exposes the later missing-source diagnostic. |
| ALC20 | 20 | Link and unlink retain project prose around the managed AGENTS.md block. | `internal/adopt/link_transaction_test.go` (`TestUnlinkRemovesOnlyCleanManifestAssets`) | Omission of the domain transform or prose copy fails the existing round trip. |
| ALC21 | 21 | Relink preserves empty or authored CLAUDE.md without injected imports. | planned internal/adopt/lifecycle_special_test.go (`TestLifecycleClaudeOwnership`, project fixtures) | A project-preserve-to-change swap injects the managed form. |
| ALC22 | 22 | Relink installs the current CLAUDE.md form for absent or known forms. | planned internal/adopt/lifecycle_special_test.go (`TestLifecycleClaudeOwnership`, managed fixtures) | Omission of the selected replacement leaves the absent or legacy form. |
| ALC23 | 23 | Link preserves special instruction files without opening their contents. | planned internal/adopt/lifecycle_special_test.go (`TestLifecycleSpecialInstructionFiles`) | Omission of the special-file guard reaches a deadline or changes the file kind. |
| ALC24 | 24 | Link and unlink retain the current effective-hook state behavior. | planned internal/adopt/lifecycle_special_test.go (`TestLifecycleHookCompatibility`) | A foreign-hook-to-managed swap replaces a foreign hook or changes the existing refusal. |
| ALC25 | 25 | Compatibility repair retains strict ownership and retained recovery behavior. | planned internal/adopt/lifecycle_special_test.go (`TestLifecycleRepairCompatibility`) | Omission of strict mode preservation changes repair acceptance for a mode-edited asset. |
| ALC26 | 26 | Unlink removes only clean present rows while absent rows remain silent. | `internal/adopt/link_transaction_test.go` (`TestUnlinkRemovesOnlyCleanManifestAssets`) and planned internal/adopt/lifecycle_unlink_test.go (`TestLifecycleUnlinkPartitions`) | Omission of a clean remove leaves a row destination present. |
| ALC27 | 27 | Unlink retains modified managed bytes and reports a residual. | `internal/adopt/link_transaction_test.go` (`TestUnlinkKeepsModifiedAssetsAndProjectCollisions`) | A preserve-to-remove swap deletes the existing modified fixture. |
| ALC28 | 28 | Unlink reports unsafe manifest rows without removing their targets. | `internal/adopt/link_transaction_test.go` (`TestUnlinkRefusesAManifestRowOutsideTheRepo`) and planned internal/adopt/lifecycle_unlink_test.go (`TestLifecycleUnlinkRefusals`) | Omission of refusal lets a hostile path or special target disappear. |
| ALC29 | 29 | Unlink dry-run leaves the complete destination snapshot unchanged. | `internal/adopt/link_transaction_test.go` (`TestUnlinkDryRunOnACleanRepoWritesNothing`) and planned internal/adopt/lifecycle_unlink_test.go (`TestLifecycleDryRunSnapshot`) | Omission of the dry-run publisher guard changes bytes, modes, links, or directory entries. |
| ALC30 | 30 | Unlink keeps the manifest whenever a modified or refused residual remains. | `internal/adopt/link_transaction_test.go` (`TestUnlinkKeepsModifiedAssetsAndProjectCollisions`) and planned internal/adopt/lifecycle_unlink_test.go (`TestLifecycleUnlinkRefusals`) | Omission of the manifest-last guard removes the resumable ownership record. |
| ALC31 | 31 | Upgrade retains the existing count table for every plan kind. | planned internal/adopt/lifecycle_preview_test.go (`TestLifecycleUpgradePreviewCompatibility`) | A synthetic-to-exact non-file projection swap changes a named baseline count. |
| ALC32 | 32 | Upgrade retains equal-version no-op and forced-downgrade behavior. | planned internal/adopt/lifecycle_preview_test.go (`TestLifecycleUpgradeVersionGuards`) | Omission of the equal-version or downgrade guard changes a real Upgrade result. |
| ALC33 | 33 | Setup plan, declined confirmation, and ambiguous auto-confirm retain their write limits. | planned internal/adopt/lifecycle_preview_test.go (`TestLifecycleSetupCancellation`) | Omission of confirmation guards changes the pre-operation destination snapshot. |
| ALC34 | 34 | Lifecycle publication retains the baseline failure and recovery result. | `internal/adopt/adopt_test.go` (`TestPromoteAllRollsBackOnDestinationSyncFailure`) and planned internal/systemtest/lifecycle_test.go (`TestLifecyclePublicationCompatibility`) | Omission of one restore leaves a changed preimage after the fault. |
| ALC35 | 35 | Link, Setup, Upgrade, Unlink, and compatibility repair consume selected owner effects through their real entries. | planned internal/adopt/lifecycle_consumption_test.go (`TestLifecycleOwnerConsumption`) and planned internal/systemtest/lifecycle_test.go (`TestLifecycleEntryJourney`) | Omission of the selected add or remove result changes the corresponding real entry fixture. |
| ALC36 | 36 | The frozen baseline and candidate pass the same characterization with old assertions unchanged. | planned differential run through existing package and system test owners | An effect swap changes one exact snapshot while the frozen baseline remains green. |
| ALC37 | 37 | Each migration checkpoint removes its old classification copy. | review-owned: source census against `PlanLifecycle`, `transactionalRepair`, and `planUnlink` | A surviving caller-side effect branch defeats the owner's deletion test. |

### Edge inventory

The profile checklist applies to the consumer repository and the installed kit source.
The implementation reuses existing safe observations instead of widening their input policy.
Each in-scope partition receives the rows below.

| edge | fixture and disposition | rows |
|---|---|---|
| Paths with spaces and glob characters | Real consumer root and managed leaf retain literal paths | ALC06, ALC28, ALC35 |
| Absent versus empty files | Seeds, CLAUDE.md, manifests, and inline assets use paired fixtures | ALC15, ALC21, ALC22, ALC31, ALC32 |
| Absent versus empty directories | Source trees remain absent-tolerant, empty managed parents remain safe sweep candidates | ALC06, ALC26, ALC29 |
| Symlink destination and symlink parent | Canonical adapter, foreign equal adapter, ordinary write, and seed fixtures | ALC10, ALC11, ALC19 |
| Dangling leaf link | A seed remains present, an ordinary unowned collision remains project-owned | ALC13, ALC15 |
| Special target and special source | Caller guards keep direct special files unopened and preserve current source failures | ALC23, ALC28, ALC34 |
| Missing executable or unresolved hooks directory | Preserve existing named refusal before promotion | ALC24, ALC34 |
| Manifest escape | Unlink keeps partial refusal rather than whole-plan validation | ALC28, ALC30 |
| Safely normalized manifest path | Unlink resolves sub/../asset to asset without a new refusal | ALC26 |
| Equal content with changed mode | Ordinary relink and strict repair retain distinct mode policy | ALC09, ALC25 |
| Multiple preserved paths | The owner emits each complete preserve decision without duplicate paths | ALC03, ALC04 |
| Invalid immutable legacy input | Preserve current duplicate, path, and fingerprint refusal order | ALC05 |
| Report control bytes | Retain current render failure before link promotion | ALC34 |
| Confirmation EOF or decline | Setup retains its existing exit and write behavior | ALC33 |
| Publication fault and interruption | Existing fault seam observes before publication, middle promotion, and final manifest publication | ALC34 |
| Modified destination after observation | Existing transaction identity recheck remains authoritative | ALC34 |
| Existing manifests without a version | Upgrade retains absent, empty, and missing-header refusals | ALC32 |
| Foreign hook and diverted hook | Link hard refusal and unlink silent preservation remain distinct | ALC24 |
| Package-variable substitution | The syncDirectory substitution stays in the adopt test process | ALC34 |

Won't handle: new JSON escape or shell grammar behavior — the existing command dispatcher remains the in-scope caller.
Won't handle: new TOON quoting policy — renderVerdicts and the upgrade table remain the in-scope callers.
Won't handle: new Git ref or merge policy — the existing hook readers remain the in-scope callers.
Won't handle: new trust or bootstrap policy — the existing selected system owner remains the in-scope caller.
Won't handle: new directory sweep failure semantics — sweepEmptyDirs remains the in-scope caller with its current best-effort behavior.
Won't handle: a new refusal for symlinked instruction-file targets — the current domain helpers remain the in-scope callers.

## Ownership fences

These exact fences equal the implementation tickets' Writes union.
The blocker chain orders each overlapping writer on one retained integration source.
The implementation cannot modify an old test assertion through these fences.
All new tests live in new files.

| owner responsibility | exact repo-relative paths | reviewer disposition |
|---|---|---|
| Immutable decision and projections | `internal/adopt/decision.go`, `internal/adopt/lifecycle_decision_test.go` | ticket graph accepted |
| Observation and application composition | `internal/adopt/link_transaction.go`, `internal/adopt/link_stage.go`, `internal/adopt/lifecycle_consumption_test.go`, `internal/adopt/lifecycle_special_test.go` | accepted |
| Unlink decision consumer | `internal/adopt/unlink.go`, `internal/adopt/lifecycle_unlink_test.go` | accepted |
| Upgrade compatibility projection | `internal/adopt/upgrade.go`, `internal/adopt/lifecycle_preview_test.go` | accepted |
| Bespoke preparation adapters if required | `internal/adopt/marker.go`, `internal/adopt/link_hook.go` | accepted, domain behavior stays fixed |
| Real process integration proof | `internal/systemtest/lifecycle_test.go` | accepted, existing system owner only |
| Spec and review artifacts | `specs/adoption-lifecycle-consumption/`, `reviews/adoption-lifecycle-consumption.md` | planning and review artifacts only |
| Compiled-map navigation | `capture/architecture-decision-maps-2026-10-06.md` | C06 link only |

The fences exclude internal/adopt/transaction/, the current transaction adapter, the payload registry, public command routing, and existing test files.
No new cross-package import edge is necessary.
If implementation needs a new file or a wider fence, update the approved plan before dispatch.
A wider fence does not authorize a changed assertion or observable behavior.

## Out of scope

- A new adoption dry-run command surface: separate feature, estimated 6 edits and 2 gate runs.
- Exact upgrade preview counts for actual writes: separate behavior change, estimated 4 edits and 2 gate runs.
- A new transaction engine or recovery policy: separate architecture outcome, estimated 8 edits and 3 gate runs.
- Generic fixture migration: FT360 and landing-test-efficiency own that work, estimated 5 edits and 2 gate runs.
- A new manifest schema or stricter legacy parse policy: separate migration, estimated 4 edits and 2 gate runs.
- New setup, hook, or instruction-file behavior: separate feature or defect outcome, estimated 4 edits and 2 gate runs.

These estimates count the named capability's likely production, test, and documentation surfaces.
They are scope estimates, not an active commitment or implementation budget.

## Further notes

### Source-sentence-to-row map

| source clause | coverage |
|---|---|
| FT217: one immutable inventory comparison decides ordered add, change, remove, and preserve operations | ALC01–ALC05, ALC35, ALC37 |
| FT217: each verb retains distinct validation, safety, reporting, and atomic filesystem implementation | ALC19–ALC34, ALC38 |
| FT217: tests reach the public command seam and decision interface | ALC01–ALC05, ALC35, ALC36 |
| FT217: no dry-run surface or other observable behavior | ALC29, ALC31–ALC36, explicit scope cuts |
| Topic ticket 1: retain and deepen the existing decision | ALC01–ALC05 |
| Topic ticket 1: reconcile suppressed preserve behavior with the required preserve behavior | ALC02, ALC03 |
| Topic ticket 2: retain pre-existing assertions and outcomes | ALC36 |
| Topic ticket 2: retain relink, modified preservation, seeds, read-only unlink, and rollback journeys | ALC07, ALC12, ALC14, ALC15, ALC27, ALC29, ALC34 |
| Topic ticket 2: add characterization before uncovered behavior moves | ALC36, ALC37 |
| Topic ticket 2: add evidence that fails when execution bypasses the decision | ALC01, ALC35 |
| Topic ticket 2: coordinate fixture work without substituting it for production consumption | Scope cuts and ALC37 |

### Fixed pre-review proof checklist

- Cited symbols: all current symbols resolve in the source files listed above. New test names carry the planned marker.
- Import edges: none. All policy and consumer work stays in internal/adopt.
- Source-row clauses and occurrences: the table above covers FT217 and both topic answers. Their copied policy occurs nowhere in production.
- Promised field labels: existing Current, Desired, Preserve, Operations, and Refusal remain. Decisions and Execution are additive. Private record names can vary.
- Changed-function callers: the reader census lists every PlanLifecycle, transactionalLink, transactionalRepair, planUnlink, and upgradePlanCounts caller.
- Copy survival: ALC37 removes effect selection from both transactionalRepair and planUnlink. ALC01 and ALC35 test owner consumption.
- Rendered-shape readers: none change. Existing reports and public command inventories retain their bytes.
- Pin operators: the decision test uses reflect.DeepEqual. Other preserved tests use their existing exact or substring operators unchanged.
- Entry reads: git root, kit location, manifests, source files, destination observations, hook reads, and setup terminal detection remain caller-owned.
- Derived expectations: ALC02 uses its existing independent expectation. New complete-decision expectations require named and demonstrated mutation witnesses.
- Consolidated rules: the execution partition table maps the old link and unlink rules to the single owner.
- Quantified obligations: ALC35 covers every applying entry. ALC31 covers file, adapter, inline, inline-exec, and seed input kinds.
- Workflow-step writes: none. This spec changes no implementation workflow or checkpoint digest.

### Planned evidence and validation

ALC38 uses the real shared execution caller, `transactionalLink`, in the adopt test process.
The fixture has a Git repository, valid instruction files, and a resolvable absent managed hook.
Its ordered plan starts with an absent `accepted.txt` inline entry.
The next entry is an absent `linked/asset` inline destination below `linked -> target`.
The target directory exists inside the repository.
The final file entry names an absent kit source.

The inline entry needs a write and cannot take the converged adapter exception.
The current caller reaches the symlink-parent return before it inspects the final source.
The source pins are `internal/adopt/link_transaction.go:120` and `internal/adopt/link_transaction.go:157`.
The test expects exit 1 with changed false.
It expects exactly `conflict: linked/asset has a symlink parent directory\n` on stderr.
The test also retains ALC19's destination snapshot comparison.

The omission witness removes the earlier symlink-parent return while retaining the later source check.
The mutant emits the missing-source diagnostic and still publishes nothing.
The exact diagnostic assertion must turn red for that mutant.
A no-publication assertion alone cannot prove this order.
This characterization precedes the observation collector migration.

ALC26 also includes a clean `sub/../asset` manifest row whose fingerprint matches the destination `asset`.
The real Unlink entry removes `asset` through the existing normalization rule.
The source is `internal/adopt/unlink.go:178`.
The owner must not apply the stricter legacy immutable path grammar to this execution input.

ALC31 drives real Upgrade for reachable file, adapter, and inline fixtures with owned and unowned paths.
Those real calls include preserved CLAUDE.md, dropped rows, hook refresh, and the current count table.
The actual upgradePlanCounts seam additionally receives inline-exec, seed, and unreadable-file-source fixtures.
The baseline preview assertions precede any projection migration.
ALC32 drives real Upgrade for equal, newer, older, forced older, absent manifest, empty manifest, and missing-header fixtures.

The producer precision correction follows the accepted review without changing count expectations or scope.
Upgrade calls buildLinkPlan, whose entries are file, inline, and adapter.
Only confirmed Setup adds inline-exec and seed entries.
The source pins are internal/adopt/upgrade.go:51, internal/adopt/link.go:98–139, and internal/adopt/setup.go:310–313.
Independent ticket Spec and Coverage review accepted this correction and retained the real Upgrade consumer proof.

ALC34 records baseline results for source disappearance, stage failure, report failure, promotion failure, and final manifest failure.
It uses the current sync seam for a deterministic restore omission witness.
It uses the current interruption seam for fresh-process recovery compatibility.
The system test owner supplies the selected binary and child environment.
The spec does not claim a new all-or-nothing unlink transaction.
Unlink keeps its current separate manifest-last publication and directory sweep.

The author ran no behavior tests during this planning pass.
The implementation first runs the same new characterization on the frozen baseline and the candidate.
The run records exit codes, exact reports where promised, destination bytes, modes, symlink targets, and manifest rows.
Normalize only fixture-root paths and existing volatile recovery identifiers.
Do not normalize a policy result or a count.
Record each named mutation's observed red and restored green.

At implementation close, run the focused adopt checks and the required system suite.
Run the system suite through `bench test --check system`.
The landing supplies the whole-project gate.
Keep all pre-existing assertions and expected outcomes unchanged.
A changed expectation or uncovered behavior delta stops the preserving refactor.

Flagged additions: the complete Decisions view is additive engineering structure required by FT217's preserve decision. It adds no command behavior
Residual unknowns: baseline characterization is not executed in this planning phase. A behavioral contradiction discovered by that evidence needs a separate reviewer decision

### Approval items

| item | proposed disposition |
|---|---|
| Implementation line | gpt-5.6-sol/high, based on the current mid binding and the hardest migration |
| Decision seam | Extend PlanLifecycle additively with exact immutable execution observations |
| Legacy compatibility | Derive Operations from Decisions while keeping preview input policy |
| Consumer seam | Real entries select effects from the complete owner result |
| Coverage and edges | Acceptance rows and edge table above, with planned mutation witnesses |
| Ownership fences | Exact paths above, old tests and transaction owner excluded |
| Out of scope | Keep each separate capability outside FT217 |
| Review sequence | Spec review accepted before slicing; ticket review accepted; implementation authorization remains separate |

### Ticket approval list

| ticket | Blocked by | independently delivered outcome |
|---|---|---|
| 01-expose-complete-decisions-through-upgrade-preview.md | none | Real Upgrade preview preserves counts through the additive owner |
| 02-consume-decisions-for-desired-link-assets.md | 01-expose-complete-decisions-through-upgrade-preview.md | Ordinary desired assets and seeds consume exact decisions across applying entries |
| 03-consume-decisions-for-dropped-and-bespoke-link-assets.md | 02-consume-decisions-for-desired-link-assets.md | Dropped and bespoke paths complete the shared Link consumer |
| 04-consume-unlink-decisions-and-close-entry-proof.md | 01-expose-complete-decisions-through-upgrade-preview.md, 02-consume-decisions-for-desired-link-assets.md, 03-consume-decisions-for-dropped-and-bespoke-link-assets.md | Unlink consumes selected effects and completes aggregate entry and recovery proof |

Each ticket carries its exact Covers and Writes fields and its focused integration checks.
Independent ticket round 1 approved granularity, blockers, and merge-or-split choices at source e8812891fe2db577a9a7b73278de738272565fed.
No approval of this planning graph starts implementation or changes the active commitment.

### Completion plan

The version 1 plan below declares source-bound implementation evidence, not executed planning results.
Before dispatch, the coordinator applies the required version 2 execution declaration.
Each ticket records baseline characterization before its mechanical production move.
Each new independent expectation records its named red and restored green.
Review records bind three independent axes to each frozen chunk pair.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "ALC-OWNER",
      "tickets": [
        "01-expose-complete-decisions-through-upgrade-preview.md"
      ],
      "verification": [
        {
          "id": "decision-preview",
          "command": "bench test --package ./internal/adopt --run 'Test(PlanLifecycleFromImmutableInventories|CompleteLifecycleDecisions|LifecycleDecisionOrder|LifecycleUpgradePreviewCompatibility|LifecycleUpgradeVersionGuards)'"
        },
        {
          "id": "owner-preview-probes",
          "command": "bench probe internal/adopt/decision.go --omit <preserve-row-source> --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "probe": "Preserve-row omission, owner-order omission, and real Upgrade synthetic-to-exact preview swap must bite and restore."
        }
      ]
    },
    {
      "id": "ALC-LINK",
      "tickets": [
        "02-consume-decisions-for-desired-link-assets.md",
        "03-consume-decisions-for-dropped-and-bespoke-link-assets.md"
      ],
      "verification": [
        {
          "id": "desired-consumption",
          "command": "bench test --package ./internal/adopt --run 'Test(LifecycleLinkPartitions|LifecycleSeedOwnership|LifecycleHardRefusalOrder|LifecycleRepairCompatibility|LifecycleSetupCancellation|TransactionalLinkAdoptsUnownedAdapterThroughSymlinkParent|SetupSeedsGateInputs)'"
        },
        {
          "id": "dropped-bespoke",
          "command": "bench test --package ./internal/adopt --run 'Test(LifecycleDroppedRows|LifecycleClaudeOwnership|LifecycleSpecialInstructionFiles|LifecycleHookCompatibility|LifecycleOwnerConsumption)'"
        },
        {
          "id": "link-probes",
          "command": "bench probe internal/adopt/link_transaction.go --omit <early-parent-return-source> --package ./internal/adopt --run '^TestLifecycleHardRefusalOrder$'",
          "probe": "Early competing-refusal omission, selected desired/dropped/bespoke effects, mode facts, seeds, and strict repair witnesses must bite and restore."
        }
      ]
    },
    {
      "id": "ALC-UNLINK",
      "tickets": [
        "04-consume-unlink-decisions-and-close-entry-proof.md"
      ],
      "verification": [
        {
          "id": "aggregate-adopt",
          "command": "bench test --package ./internal/adopt"
        },
        {
          "id": "aggregate-system",
          "command": "bench test --check system"
        },
        {
          "id": "unlink-entry-recovery-probes",
          "command": "bench probe internal/adopt/unlink.go --omit <residual-manifest-guard-source> --package ./internal/adopt --run '^TestLifecycleUnlinkRefusals$'",
          "probe": "Selected effects at all applying entries, unlink residual and dry-run guards, and the existing deterministic restore omission must bite and restore."
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "adopt-unchanged-assertions",
      "command": "bench test --package ./internal/adopt"
    },
    {
      "id": "system-entries-and-recovery",
      "command": "bench test --check system"
    },
    {
      "id": "copy-survival-census",
      "command": "bench consumers PlanLifecycle transactionalRepair planUnlink --production --full"
    }
  ]
}
```
