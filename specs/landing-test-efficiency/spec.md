# Reduce landing test cost without losing failure detection

Status: staged

Decision source: Reviewer-confirmed current conversation, 2026-10-06.

Verification log: 2 iteration(s) to accept the spec, then 1 iteration to accept the tickets. Sol 6.1/high reviewed both checkpoints.

Implementation approval: pending. This spec does not change the delivery commitment or authorize implementation.

## Problem

Landing repeatedly spends about six minutes in the whole-project gate. Existing runs give enough timing evidence to choose bounded improvements.
The largest opportunities involve fixture construction and repeated checking. Faster tests have value only when they retain their purpose and failure detection.

Fourteen completed green prospective runs span 2026-10-05T23:56Z through 2026-10-06T21:53Z. They grade different trees and show similar phase costs.
The table records seconds. Package times overlap because Go runs packages concurrently. Their savings cannot be added.

| Existing evidence | Median | Minimum–maximum |
|---|---:|---:|
| Whole gate | 372.581 | 357.678–395.717 |
| Ordinary tests | 260.313 | 247.708–277.907 |
| System tests | 81.147 | 76.243–88.753 |
| Race tests | 3.126 | 2.964–3.304 |
| Adoption package | 161.788 | Package wall time |
| Repair package | 152.813 | Package wall time |
| Worktree package | 92.256 | Package wall time |
| Conformance package | 46.940 | Package wall time |

The latest sample is `.logs/gate-20261006T215317.168159507Z-3068638.jsonl`. It records tree `75ef2a5d5510a1bdc098dd9b90f0ee8384b064ea` and a 368144 ms total.
Its ordinary and system phases take 252198 ms and 78089 ms. Its matching output records adoption at 159.371 seconds and repair at 149.691 seconds.
The tracked `capture/retros/spec-stage-grader-trace.md` supplies another durable sample: ordinary tests take 252095 ms and system tests take 79157 ms.

These records lack individual test timings. They locate costs but do not prove the savings from any proposed change.
Timestamp gaps do not establish removable overhead. Specification work runs no manual benchmark or test suite.

## Solution

Use small, isolated kits for tests of adoption and repair behavior. Keep complete payload journeys where the payload is the subject.
Consolidate repeated conformance executions that prove the same dispatch contract. Share source reads and parsing across four existing visitors.
Materialize each prospective tree once per checkout. Keep one checkout alive through execution and the subsequent durable evidence inspection during authorization.

Preserve each distinct failure assertion. Demonstrate named omissions before removing an expensive fixture or redundant execution.
Measure during implementation verification and required landings. Report actual costs without a promised percentage reduction.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: Prospective authorization is the hardest boundary. Existing ownership and refusal tests constrain it, but shared checkout lifetime needs careful integration.
Harder chunks: LTE-C7, LTE-C8.

### Fixtures with a clear purpose

1. As a maintainer, I want fixture costs tied to existing evidence, so that optimization starts with a supported premise.
2. As a test author, I want a small kit for local behavior, so that unrelated payload growth does not multiply test work.
3. As a test author, I want a private repository and home per scenario, so that mutations cannot leak between tests.
4. As a maintainer, I want complete payload journeys retained, so that installation omissions still fail.
5. As a user, I want managed repairs to restore missing assets, so that compatibility recovery remains useful.
6. As a user, I want modified assets preserved, so that repair cannot overwrite my content or modes.
7. As a user, I want foreign assets preserved, so that repair respects project ownership.
8. As a user, I want private content absent from repair records, so that diagnostics do not expose it.
9. As a user, I want undo to restore exact prior bytes and modes, so that recovery is reversible.
10. As a user, I want undo to restore prior absence, so that created files do not remain.
11. As a user, I want later file replacements preserved, so that undo respects object identity.
12. As a user, I want repeated repair to preserve tracked state, so that repair is idempotent.
13. As a user, I want permission, process, and runtime boundaries retained, so that repair cannot silently expand its authority.
14. As a user, I want failed adoption to roll back, so that partial promotion does not replace my files.

### Conformance work that proves its contract

15. As a maintainer, I want one dispatch fixture to prove dev membership, so that duplicate complete runs disappear.
16. As a maintainer, I want exact timing row counts, so that missing or repeated checks cannot look complete.
17. As a maintainer, I want stable timing order and reset, so that subsequent runs cannot reuse stale rows.
18. As a test author, I want diagnostic fixtures to run their registered owners, so that unrelated checks do not dominate them.
19. As a maintainer, I want every retained mutation to reach its registered owner, so that faster tests still prove executable wiring.
20. As a maintainer, I want source visitors to share parsing, so that repeated walks do not multiply work.
21. As a maintainer, I want each visitor's domain preserved, so that sharing cannot hide a forbidden source occurrence.
22. As a maintainer, I want exact diagnostics preserved, so that parse failures and policy violations remain actionable.
23. As a maintainer, I want separate root and kit subjects, so that a clean tree cannot conceal a broken subject.
24. As a test author, I want fresh observations after mutations, so that cached parsing cannot produce a false green.
25. As a maintainer, I want independent binding and omission checks retained, so that shared machinery cannot grade itself.

### Private prospective execution

26. As a developer, I want the requested tree materialized directly, so that checkout construction does not first populate unrelated HEAD files.
27. As a developer, I want authorization to retain its owned checkout through inspection, so that it does not construct the same subject twice.
28. As a reviewer, I want green authorization to reread durable evidence, so that a successful process result cannot authorize publication alone.
29. As a reviewer, I want tree and baseline identity checked, so that stale evidence cannot authorize a different obligation.
30. As a developer, I want private artifacts removed after completion or failure, so that optimization does not leave resources behind.
31. As a developer, I want dead-owner recovery to preserve live and foreign owners, so that concurrent landing remains safe.
32. As a maintainer, I want measured results with retained quality evidence, so that speed claims remain reviewable.

## Implementation decisions

### Small fixture kits

A new `internal/adopt/adopttest` helper owns fixture materialization. It does not import `adopt` and creates no shared writable state.
It copies explicitly selected canonical kit assets into a private directory. It preserves bytes and executable modes and fails when a requested asset is absent.
The selected paths describe fixture stimuli, not a second production payload registry. `kitpayload.PayloadRows` remains the production membership owner.

Tests keep real `Link`, `Doctor`, and undo calls. They do not replace these operations with generated expected manifests.

The first fixture slice migrates repair sessions. A second slice migrates local link and setup scenarios after the helper checkpoint passes.
Keep kit-repair sessions independent: `kitSession` deliberately makes the fixture repository its kit. Its wrapper, broker, and symlink assertions remain intact.
Keep full source assets for setup cases that execute the generated gate. Full system adoption and installed-wrapper journeys continue to consume the complete kit.

Before migration, record each affected test's asserted predicate and required source assets. Remove only assets unrelated to that predicate.
Retain each assertion about bytes, modes, record contents, diagnostics, status codes, and file identity. Keep independent expected gate inputs where an omission proof requires them.
Tests that change process environment or working directory remain serial. No package-level fixture repository, home, manifest, or mutable snapshot is shared.

### Conformance executions

Combine `TestDevTierExecutesExactlyDevChecks` and `TestTimingLinePerCheck` into `TestTimingOrderStable`.
Use exactly two dispatcher executions in that fixture. Preserve registry-derived membership, ship exclusion, exact row cardinality, order, and timing-file reset assertions.
Do not compare elapsed values from real runs. Timing format checks use the existing row grammar or canned events.

Narrow three diagnostic fixtures: hostile root paths, executable Git mode, and absent versus empty inputs.
Use registry-resolved owner bindings for their exact diagnostics. Keep one dispatch-level assertion for each changed routing boundary.
The retained fixture universe still invokes each registered owner, proves its red, restores the mutation, and proves the diagnostic disappears.
The fixture-universe architecture check remains a deletion oracle. Independent binding identity, tier, and subject checks remain independent.

### Shared source observations

Introduce a run-owned source snapshot under `internal/conformance/sourcefiles`. It owns lazy traversal, file bytes, parse results, and observation errors.
The snapshot exposes observations to visitors, not policy verdicts. Visitors retain their own inclusion rules, diagnostics, and error posture.

One dispatcher invocation owns snapshots keyed by subject. Equal root and kit paths share one snapshot. Distinct subjects never share observations.
A direct owner invocation receives a new snapshot. Fixture mutation and restoration each start a new invocation.

Prove sharing through `RunConformanceSelection` with the four real bindings. Count directory observations, source reads, and parses through the snapshot's observation boundary.
For equal subjects, an overlapping source has one observation of each kind. Distinct subjects each have their own observations and diagnostics.

The proof also detects surviving direct traversal, source-read, and parse calls in the migrated visitors and their source helpers.
Use the registered function identities to select those visitors. Metadata reads outside the Go source scan remain explicit exceptions.

Demonstrate a red for per-binding snapshots and a red for restoring an old walker. Helper-only reuse tests do not satisfy this contract.

Extend executable binding dispatch with a typed snapshot function form. Keep the named check functions as binding identities.
Do not wrap the four bindings in anonymous closures or derive their tier and subject from the advertised registry.
Existing direct helper signatures can remain as fresh-snapshot adapters when other callers require them.

| Visitor | Existing domain retained by the shared scan |
|---|---|
| Ordinary-build census | Whole tree, except `.git`, `.logs`, `dist`, `node_modules`, and `vendor`. Skip sources containing `//go:build ship`. |
| Git plumbing owner | Root Go files plus `cmd` and `internal`. Exclude tests, `internal/git`, and `internal/gittest`. |
| Bounds policy | Non-test Go files in `cmd` and `internal`, excluding `internal/bounds` as callers. Retain separate registry and read-seam inputs. |
| Skip ownership | Root Go files plus `cmd` and `internal`, including tests. Exclude `internal/capability`. |

Do not replace these domains with one intersection or one broader policy domain. Preserve each visitor's ordering and diagnostic filename form.
Preserve architecture walk errors, Git and skip traversal tolerance, and bounds registry failures. Preserve the skip visitor's no-`Skip` parse bypass.
Use immutable results only within one invocation. No process-global or cross-run source cache is authorized.

### Prospective checkout lifecycle

`prospectiveartifact.Owner` remains the sole owner of registration, materialization, owner records, and teardown.
Change materialization so it populates the requested tree without first populating HEAD. The checkout still has valid Git administration and the requested index.
Choose the Git operation only after a focused implementation probe proves the hostile-path and tree-shape cases. This spec claims no untested flag behavior.

Add one gate-owned execution-and-inspection operation for authorization. It retains the artifact owner until both operations finish.
After execution, construct a fresh subject evaluation and reload durable evidence through the existing evidence reader.
Do not accept an execution result's in-memory green inspection as that reload. Do not expose the artifact owner to lifecycle callers.

Standalone execution and standalone evidence inspection retain their current entry contracts. Standalone inspection still owns a temporary checkout and publishes its owner record.

The installed broker remains the independent authority for landing. It selects the prospective tree and trusted baseline phase schedule.
The graded tree supplies the run binary under the existing build and identity checks. The shared owner grants no new executable authority.
Evidence remains bound to the exact tree, baseline runner identity, and freshness window. No verdict format, token, schema, or reuse policy changes.

## Implementation chunks

The spec passed independent review before ticket authoring. Each chunk closes through review before dependent consumer tickets start.

| Stable chunk ID / tickets | Delivered outcome | Acceptance rows | Tests | Harder chunk |
|---|---|---|---|---|
| LTE-C1 / 01-build-isolated-repair-kits.md | Small isolated repair kits through a shared fixture owner | LTE1, LTE2–LTE14 | Repair behavior and helper isolation | no |
| LTE-C2 / 02-shrink-local-adoption-fixtures.md | Smaller local adoption fixtures with complete integration retained | LTE15–LTE17 | Link rollback and setup preservation | no |
| LTE-C3 / 03-consolidate-conformance-executions.md | Two-run dispatch proof and owner-scoped diagnostics | LTE18–LTE24 | Timing, fixture bite, and dispatch | no |
| LTE-C4 / 04-share-source-observations-for-git-policy.md | Shared source owner used by Git plumbing | LTE25–LTE27, LTE32 | Snapshot and Git visitor equivalence | no |
| LTE-C5 / 05-share-the-architecture-source-scan.md | Architecture census uses shared observations | LTE28, LTE33 | Architecture equivalence and omissions | no |
| LTE-C6 / 06-share-the-bounds-source-scan.md, 07-share-skip-scans-and-prove-dispatch-reuse.md | Bounds and skip visitors use shared observations | LTE29–LTE31, LTE34, LTE46 | Domain, routing, mutation freshness, and composed sharing | no |
| LTE-C7 / 08-materialize-only-the-requested-tree.md | Direct materialization of the requested tree | LTE35–LTE37 | Real Git tree shapes and lifecycle | yes |
| LTE-C8 / 09-retain-the-owner-through-authorization.md | One owned checkout through authorization and durable inspection | LTE38–LTE44, LTE45 | Authorization, evidence, and owner recovery | yes |

## Testing decisions

Test purpose governs consolidation. A deleted setup or execution must map to a surviving assertion and its named omission proof.
Use existing real filesystem and Git seams for ownership, mode, symlink, identity, and recovery behavior. A stat-only assertion cannot replace identity checks.

Use small source trees for visitor semantics. Differential checks compare old and new diagnostic slices over the enumerated domains before the old walker is removed.
Do not keep a permanent second policy implementation as the differential oracle. Retain independent expectations only where their named omission is demonstrated red.

During the first implementation slice, collect focused fixture setup and test costs as part of required verification.
Compare the affected package and phase timings from required landings with the recorded baseline. Record environment, tree, sample size, and package overlap limits.
If the fixture premise is false, stop that slice and return the evidence for a scope decision. Do not compensate by weakening tests.

No fixed wall-clock threshold belongs in the gate. Structural checks prove eliminated duplicate work without creating timing flakes.

### Seam diagram

```text
test scenario -> private selected kit -> real Link / Doctor -> state and diagnostics
                   tests observe bytes, modes, identity, undo, and record privacy

RunConformance -> registered named binding -> run-owned source snapshot -> visitor
                   tests observe membership, subject, diagnostics, and read counts

authorization -> gate-owned artifact lifetime -> execute -> durable reload -> decision
                   tests observe requested tree, authority, verdict, and cleanup
```

### Acceptance coverage map

All seams marked planned are future executable evidence. No implementation red or speedup is claimed during this planning phase.

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LTE1 | 1 | The first fixture slice records observed setup and test costs. | Review-owned: implementation verification record. | Unsupported savings cannot justify further fixture removal. |
| LTE2 | 2 | Selected fixture assets match canonical bytes and modes. | planned TestSelectedKitCopiesAssets in internal/adopt/adopttest/kit_test.go. | A changed executable mode or omitted requested asset fails. |
| LTE3 | 3 | Mutating one fixture leaves a second fixture unchanged. | planned TestSelectedKitsAreIsolated in internal/adopt/adopttest/kit_test.go. | A shared directory or writable inode changes the second fixture. |
| LTE4 | 5 | Repair recreates the missing managed hook. | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityManagedRepair`). | A no-op repair leaves the hook absent. |
| LTE5 | 6 | Modified managed assets retain their original content and modes. | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityModifiedConflict`). | Content and mode subcases fail on overwrite. |
| LTE6 | 7 | Foreign hook content and mode survive repair. | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityForeignConflict`). | Treating a foreign file as managed changes its snapshot. |
| LTE7 | 8 | Repair records contain none of the planted private values. | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityRepairPrivacy`). | Project, credential, and environment secrets are searched independently. |
| LTE8 | 9 | Undo restores exact prior project bytes and mode. | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityUndo`). | Normalization or mode loss differs from the saved preimage. |
| LTE9 | 10 | Undo removes an asset that was absent before repair. | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityUndoCreated`). | Omitting the deletion leaves a visible hook. |
| LTE10 | 11 | Undo preserves a later equivalent replacement's identity. | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityUndoConflict`). | `os.SameFile` detects replacement despite equal bytes and mode. |
| LTE11 | 12 | Repeated repair preserves the tracked byte and mode snapshot. | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityRepairIdempotence`). | Repeated writes that change tracked state fail comparison. |
| LTE12 | 13 | Repair preserves permission and trust configuration bytes. | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityPolicyBoundary`). | Automatic policy expansion changes the fixture files. |
| LTE13 | 13 | The supervised process responds after repair. | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityInterruptionBoundary`). | Killing the process prevents its next response. |
| LTE14 | 13 | Repair leaves the private runtime directory absent. | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityRuntimeBoundary`). | Invented launcher recovery creates the forbidden directory. |
| LTE15 | 14 | Failed promotion restores existing destinations and removes newly created destinations. | `internal/adopt/adopt_test.go` (`TestPromoteAllRollsBackOnDestinationSyncFailure`). | Injected synchronization failure exposes incomplete rollback. |
| LTE16 | 4 | Complete installation retains its installed-wrapper failure and recovery journey. | `internal/systemtest/adoption_test.go` (`TestAdoptionSmokeJourney`). | Omitting real payload or wrapper behavior breaks the journey. |
| LTE17 | 2 | Local setup fixtures preserve operator gate inputs across repeated setup. | `internal/adopt/setup_test.go` (`TestSetupPreservesOperatorGateInputs`, `TestSetupTwiceLeavesGateInputsIdentical`). | A smaller fixture cannot hide destructive reseeding. |
| LTE18 | 15 | The consolidated fixture executes exactly the dev registry membership. | planned TestTimingOrderStable in internal/conformance/tier_test.go. | Missing dev checks and included ship checks change the names. |
| LTE19 | 16 | Each execution produces exactly one timing row per expected check. | planned TestTimingOrderStable in internal/conformance/tier_test.go. | Duplicated or missing rows change cardinality. |
| LTE20 | 17 | The second execution replaces prior rows in registry order. | planned TestTimingOrderStable in internal/conformance/tier_test.go. | Appended stale rows or reordered execution fails exact comparison. |
| LTE21 | 18 | A hostile root path reaches the JSON owner's invalid-package diagnostic. | planned TestRunConformanceAcceptsHostileRootPath in internal/conformance/fixture_bite_test.go. | Incorrect quoting or subject routing loses the diagnostic. |
| LTE22 | 18 | A tracked non-executable wrapper reaches its owner's mode diagnostic. | planned TestRunConformanceChecksExecutableGitMode in internal/conformance/fixture_bite_test.go. | Filesystem-only mode checks miss the tracked Git mode. |
| LTE23 | 18 | Absent and empty inputs retain their distinct owner diagnostics. | planned TestRunConformanceDistinguishesAbsentAndEmptyInputs in internal/conformance/validity_checks_test.go. | Collapsing input states loses a named existing diagnostic. |
| LTE24 | 19 | Every retained fixture still proves its registered owner's red and restoration. | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`, `TestFixtureBiteProofArchitecture`). | Removing an owner invocation or fixture fails the universe check. |
| LTE25 | 20 | Repeated requests for one file in one snapshot read and parse it once. | planned TestSnapshotReusesObservations in internal/conformance/sourcefiles/source_test.go. | Counted readers expose a cache that repeats source work. |
| LTE26 | 24 | A new invocation observes changed and restored source bytes. | planned TestSnapshotLifetime in internal/conformance/sourcefiles/source_test.go. | A process-global cache returns the preceding mutation state. |
| LTE27 | 21 | Git plumbing diagnostics equal the old visitor across its input domain. | planned TestGitPlumbingSnapshotEquivalence in internal/conformance/git_plumbing_owner_test.go. | Root files, owner exclusions, and forbidden flag cases expose omissions. |
| LTE28 | 21 | Architecture diagnostics equal the old visitor across its input domain. | planned TestArchitectureSnapshotEquivalence in internal/conformance/ordinary_build_census_test.go. | Nested sources and excluded ship sources expose domain drift. |
| LTE29 | 21 | Bounds diagnostics equal the old visitor across its input domain. | planned TestBoundsSnapshotEquivalence in internal/conformance/bounds_policy_test.go. | Registry, read-seam, and caller mutations expose narrowed checking. |
| LTE30 | 22 | Skip diagnostics equal the old visitor across its input domain. | planned TestSkipSnapshotEquivalence in internal/conformance/skip_ownership_test.go. | Parse errors, renamed receivers, and text-only mentions expose drift. |
| LTE31 | 23 | Distinct root and kit trees produce their own subject diagnostics. | planned TestSnapshotBindingsKeepSubjects in internal/conformance/check_bindings_test.go. | A clean root cannot mask the poisoned kit census. |
| LTE32 | 25 | Function, tier, and subject substitutions still fail conformance metadata checks. | `internal/conformance/tier_test.go` (`TestConformanceMetaBites`). | Independent bindings continue to expose advertised registry drift. |
| LTE33 | 25 | Omitting an ordinary phase still fails the independent phase expectation. | `internal/conformance/ordinary_build_census_test.go` (`TestBranchNativeArchitectureCensus`). | Deriving the expectation from the phase producer would conceal omission. |
| LTE34 | 24 | Registered-owner mutation and restoration obtain fresh source diagnostics. | planned TestEveryRetainedFixtureBitesThroughRegisteredOwner in internal/conformance/fixture_bite_test.go. | A cached clean or red AST cannot satisfy both sides. |
| LTE35 | 26 | Materialization produces the requested tree's index and file contents. | planned TestMaterializeRequestedTree in internal/gate/prospectiveartifact/materialize_test.go. | A HEAD-only checkout differs on additions, deletions, and content changes. |
| LTE36 | 26 | Materialization performs only one tracked-file population. | planned TestMaterializePopulatesOnce in internal/gate/prospectiveartifact/materialize_test.go. | An observed Git operation trace detects a preliminary HEAD checkout. |
| LTE37 | 26 | Materialization preserves requested executable and symlink entries under hostile paths. | planned TestMaterializeHostileTreeShapes in internal/gate/prospectiveartifact/materialize_test.go. | Plain copying or shell interpolation changes modes, links, or path bytes. |
| LTE38 | 27 | One authorization uses one artifact owner through final evidence inspection. | planned TestAuthorizationSharesProspectiveOwner in internal/gate/execution_probe_test.go. | Counting owner creation detects the second prospective checkout. |
| LTE39 | 28 | Missing or corrupted durable evidence prevents green authorization after execution succeeds. | planned TestAuthorizationReloadsDurableEvidence in internal/gate/authorization/authorization_test.go. | Trusting the in-memory execution result incorrectly returns green. |
| LTE40 | 29 | Changed tree or baseline runner identity refuses prior evidence. | `internal/gate/prospective_owner_test.go` (`TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity`). | Weakening the evidence key permits invalid reuse. |
| LTE41 | 30 | Failed, cancelled, and timed-out executions leave no owned bundle. | `internal/gate/prospective_owner_test.go` (`TestProspectiveGateReportsAFailedBuildWithNoResidue`). Planned lifetime subcases. | Each exit path exposes a missing owner close. |
| LTE42 | 31 | Recovery retains a concurrent live authorization's owner. | `internal/systemtest/owner_artifact_recovery_test.go` (`TestConcurrentAuthorizationRetainsALiveOwner`). | Premature recovery removes the active checkout. |
| LTE43 | 31 | Killed landing recovery removes its dead owner's artifacts. | `internal/systemtest/owner_artifact_recovery_test.go` (`TestProspectiveArtifactRecoveryAfterKilledLanding`). | A missing or late owner record leaves unrecoverable artifacts. |
| LTE44 | 30 | Standalone inspection publishes its own owner record. | `internal/gate/prospective_owner_test.go` (`TestEvidenceInspectionPublishesTheOwnerRecord`). | Reusing an unowned path removes crash recovery authority. |
| LTE45 | 32 | The final report pairs timing results with retained failure proofs. | Review-owned: final implementation reconciliation. | A faster green run alone cannot establish retained quality. |
| LTE46 | 20 | The four registered visitors share source observations through the dispatcher. | planned TestConformanceBindingsShareSourceObservations in internal/conformance/check_bindings_test.go. | Counted observations and the surviving-call check fail on per-visitor snapshots or a restored old walker. |

### Edge inventory

Apply the project profile's hostile-input checklist to filesystem, Git, process, and artifact boundaries.
Fixture cases include missing requested assets, mode differences, symlinks, foreign ownership, private records, equivalent replacements, and repeated operations.
Source cases include empty trees, missing directories, malformed Go, unreadable files, root files, nested files, tests, ship tags, and excluded directories.
Exercise one permitted and one forbidden occurrence per visitor. Include same-root sharing, distinct subjects, and mutate-then-restore invocations.

Prospective cases include spaces and glob characters, tree additions and deletions, executable modes, symlinks, cancellation, timeout, failed builds, and concurrent owners.
Retain malformed, foreign, non-private, and permission-refused owner-record cases from the existing recovery suite.


**Won't handle:** Phase overlap and reduced Go parallelism. Existing sequential phase execution and `-p 4` remain the surviving callers.

**Won't handle:** Broader verdict reuse. Existing exact-tree inspection retains its present evidence policy.

**Won't handle:** Global fixture sharing or parallel repair sessions. Existing environment and working-directory mutation require private serial scenarios.

**Won't handle:** All conformance visitors. The four named visitors define this migration's complete set.

**Won't handle:** General expectation-table cleanup. Independent omission oracles remain where demonstrated reds require them.

## Ownership fences

These fences equal the ticket write union, plus the review pickup. Review confirms the final closure set before implementation.

- `reviews/landing-test-efficiency.md`
- `internal/adopt/adopttest/`
- `internal/adopt/repairtest/session_test.go`
- `internal/adopt/repairtest/kit_test.go`
- `internal/adopt/link_transaction_test.go`
- `internal/adopt/setup_prompt_test.go`
- `internal/adopt/setup_test.go`
- `internal/conformance/tier_test.go`
- `internal/conformance/tier_live_tree_test.go`
- `internal/conformance/fixture_bite_test.go`
- `internal/conformance/validity_checks_test.go`
- `internal/conformance/sourcefiles/`
- `internal/conformance/check_bindings_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/git_plumbing_owner_test.go`
- `internal/conformance/ordinary_build_census_test.go`
- `internal/conformance/bounds_policy_test.go`
- `internal/conformance/skip_ownership_test.go`
- `internal/gate/prospectiveartifact/prospectiveartifact.go`
- `internal/gate/prospectiveartifact/materialize_test.go`
- `internal/gate/engine.go`
- `internal/gate/execution_probe_test.go`
- `internal/gate/authorization/authorization.go`
- `internal/gate/authorization/authorization_test.go`

Existing oversized files require same-ticket headroom before growth. Do not change structure budgets or reviewer-owned grants to permit this work.
Prefer reducing the existing fixture and visitor bodies. New cohesive helpers belong under the declared new package prefixes.

## Out of scope

| Separate capability | Estimated work | Reason |
|---|---|---|
| FT314 cross-run evidence reuse | 5 edits, 2 gate runs | Changes an independent authorization policy. |
| Phase overlap under ADR 0024 | 4 edits, 2 gate runs | Requires a new resource census and reviewer decision. |
| FT360 generic Git helper consolidation | 12 edits, 2 gate runs | Has a staged spec and separate fixture ownership goal. |
| Remaining FT365 expectation cleanup | 6 edits, 2 gate runs | Extends beyond the four source visitors. |

These estimates describe separate future changes, not remaining acceptance for this build. No existing roadmap item or staged spec is retired here.

## Further notes

### Source and reader audit

The author inspected source at `28511b1ef4919738a9997dc3ea51439a879cf162`. The timing sample tree predates only unrelated planning changes.

| Source clause or claim | Current owner read | Coverage or disposition |
|---|---|---|
| Use existing consistent timing data | Gate JSONL, matching output, and tracked retrospective cited above | LTE1, LTE45. No manual baseline run. |
| Test purpose comes first | Repair assertions, adoption rollback, system adoption journey | LTE2–LTE17 |
| Smaller adoption and repair fixtures | `repairtest/session_test.go`, `link_transaction_test.go`, `setup_prompt_test.go`, `link_plan_test.go` | LTE2–LTE17 |
| Consolidate duplicate dispatch runs | `tier_test.go`, `checks_test.go`, `registry/registry.go` | LTE18–LTE20 |
| Narrow isolated diagnostics | `fixture_bite_test.go`, `validity_checks_test.go`, executable bindings | LTE21–LTE24 |
| Share source work | Four visitor files and their parse helpers | LTE25–LTE34, LTE46 |
| Preserve independent omissions | `check_bindings_test.go`, `TestConformanceMetaBites`, ordinary phase expectations | LTE32, LTE33 |
| Avoid repeated checkout construction | `gate/engine.go`, `gate/authorization/authorization.go`, prospective artifact owner | LTE35–LTE44 |
| Preserve crash and concurrent-owner behavior | `prospective_owner_test.go`, prospective artifact tests, system owner recovery | LTE41–LTE44 |

`bench consumers` found 77 references across 18 files for the initial helper, visitor, and gate entry inventory.
The caller inventory covers each changed signature. Stable standalone APIs need no consumer migration.
The executable binding registry and live-tree test classification are source readers. Fixture architecture checks and fixture pins must follow any test move or rename.

The injected-port census currently binds assessment, charge evidence, Git guards, preflight, and publication. This spec does not widen that registry.
New helper packages must pass Go import loading without cycles. No cross-package row cites a test-only helper as its producer.

No command grammar, timing format, evidence schema, rendered status shape, or guidance anchor changes are proposed.
The gate phase inventory, race sentinels, system journeys, and capability skip classifications retain their current owners.
The spec-stage grader adds only this planning directory and its required phase-close records.

### Caller dispositions

| Caller family | Disposition |
|---|---|
| Repair session callers in repair_test.go and privacy_test.go | Retain signatures and all assertions. Change only private fixture materialization. |
| kitSession | Keep its repository-as-kit branch and existing shim behavior. |
| consumerRepo in adopt_test.go | Retain the full-kit default for the relink journey. |
| consumerRepo and linkConsumer in link_transaction_test.go | Select minimal kits for local unlink assertions. |
| setup helpers in setup_prompt_test.go and setup_test.go | Select minimal kits except the real generated-gate execution case. |
| Four visitor bindings in checks_test.go | Retain named identities and independent tier and subject fields. |
| Direct visitor callers in their own test files | Adapt to fresh observations or keep fresh-snapshot adapters. |
| scanArchitectureGo classification test | Retain its direct byte-source input through the shared parser. |
| checkBoundCallers wait-classification test | Preserve its direct scenario and diagnostics. |
| Owner.Materialize callers in engine.go and prospectiveartifact_test.go | Retain the public signature and lifecycle assertions. |
| ExecuteTree and InspectTreeContext in authorization.go | Replace the pair with the gate-owned combined operation. |
| Standalone InspectTree wrapper | Retain standalone inspection and its independent owner. |

The second caller scan found 15 references across six files for parse helpers and materialization. No moved symbol has an external caller outside these families.

### Completion plan

The authored version 1 plan describes verification obligations. An approved build records its fresh authors through the version 2 amendment before dispatch.
System verification uses the candidate `BENCH_KIT` and the selected run binary supplied by `bench test`.

```bench-completion-plan
{"version":1,"chunks":[{"id":"LTE-C1","tickets":["01-build-isolated-repair-kits.md"],"verification":[{"id":"fixture-owner","command":"bench test --package ./internal/adopt/adopttest"},{"id":"repair","command":"bench test --package ./internal/adopt/repairtest"},{"id":"repair-omissions","command":"bench test --package ./internal/adopt/repairtest","probe":"missing managed repair and lost undo identity"}]},{"id":"LTE-C2","tickets":["02-shrink-local-adoption-fixtures.md"],"verification":[{"id":"adoption","command":"bench test --package ./internal/adopt"},{"id":"seed-omission","command":"bench test --package ./internal/adopt","probe":"omit seeded gate input"}]},{"id":"LTE-C3","tickets":["03-consolidate-conformance-executions.md"],"verification":[{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"dispatch-omissions","command":"bench test --package ./internal/conformance","probe":"omit execution, omit timing reset, and change order"}]},{"id":"LTE-C4","tickets":["04-share-source-observations-for-git-policy.md"],"verification":[{"id":"source-owner","command":"bench test --package ./internal/conformance/sourcefiles"},{"id":"git-and-bindings","command":"bench test --package ./internal/conformance"},{"id":"flag-omission","command":"bench test --package ./internal/conformance","probe":"omit one independently expected Git flag"}]},{"id":"LTE-C5","tickets":["05-share-the-architecture-source-scan.md"],"verification":[{"id":"architecture","command":"bench test --package ./internal/conformance"},{"id":"architecture-omissions","command":"bench test --package ./internal/conformance","probe":"restore a forbidden constructor and omit an ordinary phase"}]},{"id":"LTE-C6","tickets":["06-share-the-bounds-source-scan.md","07-share-skip-scans-and-prove-dispatch-reuse.md"],"verification":[{"id":"visitors","command":"bench test --package ./internal/conformance"},{"id":"sharing-omissions","command":"bench test --package ./internal/conformance","probe":"create per-binding snapshots and restore an old source walker"}]},{"id":"LTE-C7","tickets":["08-materialize-only-the-requested-tree.md"],"verification":[{"id":"materialization","command":"bench test --package ./internal/gate/prospectiveartifact"},{"id":"materialization-omissions","command":"bench test --package ./internal/gate/prospectiveartifact","probe":"materialize HEAD only and populate HEAD before the requested tree"}]},{"id":"LTE-C8","tickets":["09-retain-the-owner-through-authorization.md"],"verification":[{"id":"gate-owner","command":"bench test --package ./internal/gate"},{"id":"authorization","command":"bench test --package ./internal/gate/authorization"},{"id":"durable-reload","command":"bench test --package ./internal/gate/authorization","probe":"omit durable evidence reload"},{"id":"system","command":"bench test --check system"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/landing-test-efficiency/spec.md"},{"id":"adoption","command":"bench test --package ./internal/adopt"},{"id":"repair","command":"bench test --package ./internal/adopt/repairtest"},{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"gate-owner","command":"bench test --package ./internal/gate"},{"id":"artifacts","command":"bench test --package ./internal/gate/prospectiveartifact"},{"id":"authorization","command":"bench test --package ./internal/gate/authorization"},{"id":"system","command":"bench test --check system"}]}
```

### Independent review results

The invoking session authored the spec and all nine tickets. Sol 6.1/high supplied both independent reviews at the reviewer's request.
Spec review required composed source-sharing evidence. The second pass accepted LTE46 and its named per-binding-snapshot and restored-walker mutations.

Ticket review accepted `af685aeb768bd3083cf8cb89d70df90dfb2899fe` in one pass. It found no Standards, Spec, or Coverage findings.
The review confirmed 46 uniquely assigned rows, 23 implementation paths, and five ordered pairs with shared writes.
Planning preflight and all nine write proposals passed. All ten planning Markdown files passed the prose check.

Existing timing records supplied the performance baseline. No separate benchmark ran.
A required base update ran the full gate automatically and passed. Its source delta contained only commitment metadata.
These planning reviews establish no implemented speedup. Implementation and delivery scheduling still require reviewer approval.
