# Release evidence forwards

Status: staged

Related roadmap: FT364; partial FT368

Decision source: [Release forwards and capability scope](decisions/architecture-release.md), ready on 2026-10-06.

Source pin: a9c395e77fec36d1f60f6d057c654aecf5720e24.
The production source matches main at 37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68.

Verification log: 1 draft pass plus 2 bounded document repairs — independent gpt-6.1-sol/xhigh SPEC acceptance at 46206d53, all three axes zero, confidence 9.
Independent graph acceptance at 350bf2b6 reports all three axes zero, confidence 9.
The coordinator retains the prior acceptance; the original reviewer session identity is unavailable.
This retained high-author graph follows SPEC acceptance and grants no implementation authority.

## Problem

Release preflight repeats evidence names through a forwarding file.
The file owns no evidence policy.
Its same-package callers still use these names, so deletion requires a complete caller migration.
The declarations and callers are in `internal/releasepreflight/types.go` and the caller inventory below.

FT364 combines this mechanical cut with claims about unsupported release machinery.
Current publication code and tests do not support those broader deletions.
The npm staged refusal already has a command test.
Fixture staging has a real adapter and command path.
These sources are `internal/publication/command.go`, `fixture_registry.go`, and `command_adapter_test.go`.

## Solution

Use `internal/releaseevidence` directly for evidence types and operations.
Keep releasepreflight as the command, execution, and decision adapter.
Remove its forwarding declarations after its callers migrate.
Preserve the evidence owner's contracts and the publication state machine.

Retain public CLI and evidence export surfaces.
Unknown external use remains unknown.
This spec closes the internal forwarding outcome only.
It grants no publication, live registry qualification, or broader FT368 removal authority.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: the hardest chunk combines caller migration with error identity and a production-entry validation witness.
The owner already exists, the import edge is established, and the behavior seams are concrete.
Harder chunks: none.
The user's Sol 6.1/high author choice does not change the future implementation line.

Evidence ownership:

1. As a maintainer, I want evidence types from one owner, so that a forwarding package cannot become a second policy source.
2. As a maintainer, I want active evidence calls to reach their owner directly, so that caller behavior survives the cut.
3. As a maintainer, I want dormant forwarding declarations removed, so that unused copies do not survive the migration.
4. As a test author, I want existing preflight tests to use owner types, so that their behavior checks remain executable.
5. As a map reader, I want the moved decision source to resolve, so that the planning record remains usable.

Preflight behavior:

6. As a CLI user, I want the existing release-preflight command, so that internal cleanup does not remove my entry point.
7. As a release operator, I want the same immutable decision verdict, so that unchanged evidence receives the same result.
8. As a release operator, I want the same phase registry order, so that caller migration does not reorder release checks.
9. As a release operator, I want interrupted evidence to retain its status precedence, so that cancellation cannot become success.
10. As a release operator, I want focused publish to retain its requirement refusal, so that partial evidence cannot authorize publication.
11. As a release operator, I want the command to finalize evidence, so that a successful phase cannot bypass evidence validation.
12. As an evidence consumer, I want identical record bytes for identical inputs, so that the cut does not change stored evidence.
13. As a release operator, I want the same reader errors, so that migration does not change file or failure policy.
14. As an evidence consumer, I want the same registered schemas and artifact paths, so that existing readers remain compatible.

Publication capability:

15. As a release operator, I want unsupported npm staging refused before locking, so that the capability posture stays explicit.
16. As a fixture user, I want fixture staging retained, so that the local staged journey remains available.
17. As a release operator, I want platforms published before the wrapper, so that first publication keeps its ordering rule.
18. As a release operator, I want a failed publication to resume safely, so that completed packages are not published again.
19. As a release operator, I want promotion through the selected adapter, so that the existing latest-tag operation remains available.
20. As a release operator, I want rollback through the selected adapter, so that the existing recovery operation remains available.
21. As a maintainer, I want public CLI and export surfaces retained, so that unknown external use does not become a deletion premise.
22. As a release operator, I want live qualification to remain separate, so that local stubs cannot imply public npm acceptance.

Completion:

23. As a reviewer, I want a fixed-input differential, so that compatibility rests on comparable observations.
24. As a maintainer, I want the ownership comment to name releaseevidence, so that documentation agrees with the implementation.
25. As a reviewer, I want the complete closure to fit its files, so that the authorized chunk can land on its required lane.
26. As a roadmap reader, I want unresolved release obligations retained, so that this narrow cut does not retire unrelated work.

## Implementation decisions

### Owner and API

The evidence owner is already `internal/releaseevidence`.
Its existing public types, constants, and functions become the direct imports.
Do not create a replacement facade, new evidence package, duplicated registry, or caller-local copy.

The migration changes `DecisionInput`, `Decision`, runner fields, function signatures, and composite literals to owner-qualified types.
`Command`, `Decide`, and their execution policy remain in releasepreflight.
`ReleaseIntentError` remains the owner's concrete error type.
Use `errors.As` against that type without an additional wrapper.
RF01, RF02, RF04, RF07, and RF10 grade these contracts.

The current and future calls are:

| current preflight name | future owner name | live caller locations |
|---|---|---|
| PhaseNames | releaseevidence.PhaseNames | command.go: parseArgs and runner.run, decision.go: Decide, decision_test.go |
| phaseDefinition | releaseevidence.PhaseDefinitionFor | command.go: runner.run, runPhase, runExternal, validatePhaseInputs, vulnerability.go: runVulnerability |
| terminalStatus | releaseevidence.TerminalStatus | decision.go: Decide |
| FinalizeEvidence | releaseevidence.FinalizeEvidence | command.go: Command |
| readPackageVersion | releaseevidence.ReadPackageVersion | command.go: populateBaseIdentity, identity.go: checkIdentity |
| readRegular | releaseevidence.ReadRegular | identity.go: checkChangelog and readToolchain, vulnerability.go: runVulnerability |

This inventory comes from the current declarations and `bench consumers` at the source pin.
The hidden-tree textual sweep also covers generated source strings, scripts, workflows, tests, and dot-directories.

The definition-only function forwards are `Requirements`, `phaseSummaries`, `packageEvidenceRegistry`, and `PromoteEvidence`.
The remaining definition-only functions are `setArchiveMemberLimitForTesting`, `validateTarballForTesting`, `setExchangeForTesting`, and `setIndexEncoderForTesting`.
They also include `setRequirementsForTesting` and `atomicExchangeForTesting`.
Remove these functions from preflight without deleting their evidence-owner implementations.
RF03 grades this complete inventory.

The alias inventory is `Mode`, `Scope`, `Status`, `Profile`, `RunEvidence`, `Failure`, `Record`, and `Identity`.
It also contains `PhaseSummary`, `Manifest`, `Result`, `PhaseDefinition`, `Requirement`, and `PackageEvidence`.
Its private aliases are `requirementRegistry`, `releaseIndex`, `requirementStatus`, and `releaseIntentError`.

The constant inventory is `ModeVerify`, `ModePublish`, `ScopePreflight`, `ScopeFocused`, `StatusGreen`, and `StatusRed`.
It also contains `StatusNotRun`, `StatusInterrupted`, `ProfilePublic`, and `ProfileBank`.
The preflight `requirements` variable is a definition-only snapshot of `RequirementsRegistry()`.
Remove the snapshot with the other forwards.
The evidence owner's registry remains authoritative.

Active alias consumers are the command, decision, identity, and vulnerability files.
Test consumers are `decision_test.go`, `external_test.go`, and `identity_test.go`.
The other alias references remain inside the forwarding file.
Alias-aware consumer output also reports references to the owner's original type.
Those origin references do not create preflight consumers.

The external typed caller is `cmd/bench/main.go`, which invokes `Command`.
A generated command in `internal/conformance/native_workflow_test.go` also invokes `Command`.
Both retain the same command signature.
`Decide` remains called by `Command` and `TestDecidePreflightFromImmutableEvidence`.
No other qualified preflight export consumer was found in the hidden-tree sweep.
That observation is not a claim about unknown external consumers.

### Preserved boundaries

Keep evidence serialization, validation, promotion, and input-drift checks in releaseevidence.
Its production files and JSON registries are outside the write fence.
Record field names, null posture, JSON newline framing, and exact encoded bytes remain unchanged.
RF12, RF14, RF23, and RF28 grade this promise.

Keep the current `ReadRegular` and `ReadPackageVersion` contracts.
`ReadRegular` uses Lstat, refuses non-regular files, and returns underlying filesystem errors.
`ReadPackageVersion` keeps its existing version validation.
This spec does not tighten JSON parsing or change symlink, size, or missing-file policy.
RF13 compares the old and new calls on the named input family.

Keep the command's argument guards and phase execution order.
A focused publish run still reaches `FinalizeEvidence` after its identity phase succeeds.
Its typed refusal produces kind `requirement` and message `focused publish runs cannot authorize publication`.
RF10 protects this conversion through the real command.

The existing release-evidence probe is the validation omission witness.
It invokes `scripts/release-preflight.sh`, which reaches `releasepreflight.Command`.
It requires a generated release index, verifies that index with the offline consumer, and tests attributed native-proof mutations.
Removing the finalizer call must fail its initial case for the missing index.
RF11 records this future red, not an observed red during specification.

Publication runtime functions remain unchanged.
Retain `RunFirst`, `RunStaged`, `RunPromote`, `RunRollback`, and the existing registry adapters.
First publication, resume, promotion, and rollback retain their current record transitions and operation order.
RF16 through RF20 exercise retained boundaries, and RF29 compares their complete function bodies.

The npm staged operand refusal remains before `AcquireReleaseLock` and registry construction.
It wins over an existing held-lock refusal.
Fixture `StageSubmit` and `Approve` remain implemented.
Public npm staging remains unsupported.
RF15 and RF16 grade these distinct postures.

The fixture staging source proves a live capability path.
The inspected tests do not establish a complete successful staged publication journey.
Do not rename their narrower assertions as full staged qualification.
Local npm command stubs also establish no live registry result.
RF22 preserves this limit.

### Decision and reader closure

Move the map index and its topic folder as one unit under this spec's `decisions/`.
Preserve their bytes, answers, and relative links.
Change only the C11 navigation link in the architecture-map capture index.

The moved map's structured source still names `internal/releasepreflight/types.go` during this planning phase.
`internal/maps/validation.go` requires a source Path to name an existing regular file.
The implementation cut must replace that source entry with `internal/releaseevidence/types.go` when it removes the forwarding file.
Change its Supports and Drift descriptions to the resulting ownership fact.
Do not rewrite the resolved answers or their pinned source citations.
RF05 grades source closure, and RF27 grades planning byte preservation.

The comment above `releaseIndexAuthority` currently assigns requiredness to releasepreflight.
`internal/releaseevidence/release_requirements.go` is the actual owner.
Correct only that owner name in `internal/publication/approved.go`.
RF24 grades the correction.

The consumer sweep found no canary, anchor, command-registry, or test-source reader that pins the forwarding file.
Its structured map source is the live path-sensitive reader.
The decision ticket's old line citation is retained historical evidence.
The generated conformance command imports `Command`, so it needs no rewrite.

The inspected closure includes `internal/conformance/registry/checks.go`, `registry/packages.go`, and `injected_ports_registry_test.go`.
It also includes `fixture_bite_test.go`, `entry_point_parity_test.go`, and `ordinary_build_census_test.go`.
The anchor sweep covers `internal/anchors/registry_data.go`.
These registries, fixtures, and checks remain unchanged.
No new port, standing audit, phase, command, count, or help inventory is introduced.

### Source headroom

The selected lane calls `structure.Growth` through `internal/gate/lane.go`.
Growth grades changed files against their per-file line caps and their accepted base.
It does not grade directory crowding.
The default Go file cap is 400 lines.
No touched Go file needs a new grant.

Current source lines are command 360, decision 50, identity 116, and vulnerability 175.
The touched tests have 48, 58, and 168 lines.
The removed forwarding file has 76 lines.
The publication comment file has 188 lines.
Their declaration sources are the exact fenced files at the source pin.

Keep `command.go` below its default cap after qualification and imports.
The new contract test also stays below that cap.
Do not grow `command_adapter_test.go` or `statemachine.go`, whose inherited size exceeds the default cap.

The preflight package currently has nine Go files.
One deletion and one new test leave that count unchanged.
The publication package adds no file.
RF25 requires fresh base-to-tip growth evidence before landing.

## Implementation chunks

Independent SPEC acceptance precedes this complete vertical ticket graph.
RF1 closes the complete caller migration, forwarding cut, and preservation evidence in one green checkpoint.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| RF1 / 1-consume-release-evidence-directly.md | Preflight directly consumes its evidence owner with compatible errors, records, and publication boundaries | RF01-RF29 | Preflight and evidence packages, publication entry tests, release-evidence-probe, fixed-input differential | no |

RF1 first captures the source-pin reference on the canned compatibility family.
It then migrates the complete caller family and deletes the forwarding file in one green change.
The owner is already usable, so no owner-only preparatory chunk is needed.
The final slice carries the whole preflight-package inventory and the structured-source repair.

No unlanded architecture spec is a prerequisite.
Coordinate any concurrent edits to these files through a fresh source pin and caller census.
The strict-JSON and process-lifetime outcomes retain their separate behavior decisions.
This spec neither absorbs their repairs nor changes their approved implementation lines.

### Graph ordering and execution boundary

The stable old-to-new mapping is RF1 to RF1, with all 29 unchanged predicates owned exactly once.
The one ticket has no sibling blocker and owns all 12 implementation paths, excluding only the review pickup.
The complete evidence owner already exists; no unlanded provider or second outcome is required.
Every first-use caller, source reader, reference, typed error, fixture, and file-growth obligation closes in RF1 before its independent review.

The version-1 completion plan is a planning verification inventory, not a dispatch plan.
Before any implementation charge, the coordinator binds actual version-2 run/session identities with author-limit 1 and a fresh gpt-5.6-sol/high ticket author.
Refresh the source pin, complete caller and closure census, owner interfaces, and actual headroom before that charge.
Coordinate concurrent writes against that integrated source; staged status and planning validation grant no build authority.
Independent graph review accepted the source at 350bf2b6; implementation remains unauthorized.

### Completion plan

Every command and mutation below is a future obligation; none ran during graph authoring.
Pin each actual mutation source and its exact diagnostic before the compiling behavioral red.
Require byte-identical restoration and the same focused green; invalid or restore-failed proof closes no obligation.
The independently captured canned reference covers the entire accepted family without a persistent compatibility harness.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "RF1",
      "tickets": [
        "1-consume-release-evidence-directly.md"
      ],
      "verification": [
        {
          "id": "preflight",
          "command": "bench test --package ./internal/releasepreflight"
        },
        {
          "id": "evidence-owner",
          "command": "bench test --package ./internal/releaseevidence"
        },
        {
          "id": "publication",
          "command": "bench test --package ./internal/publication"
        },
        {
          "id": "command-entry",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "phase-wrapper",
          "command": "bench test --check package-core-guard"
        },
        {
          "id": "real-evidence-producer-consumer",
          "command": "bench test --check release-evidence-probe"
        },
        {
          "id": "map-source",
          "command": "bench test --check decision-map-integrity"
        },
        {
          "id": "finalizer-omission",
          "command": "bench test --check release-evidence-probe",
          "probe": "Replace only the production finalizer result with nil while retaining the assembled run reference; the compiling command must leave the authenticated clean case without release-index.json and red its exact missing-index diagnostic."
        },
        {
          "id": "typed-refusal-omission",
          "command": "bench test --package ./internal/releasepreflight --run '^TestCommandFocusedPublishRequiresEvidenceAuthority$'",
          "probe": "Force only concrete ReleaseIntentError recognition false while preserving compilation; the real focused Command must emit the generic evidence kind and red the exact requirement diagnostic."
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "coverage",
      "command": "bench coverage --check release-evidence-forwards"
    },
    {
      "id": "preflight",
      "command": "bench test --package ./internal/releasepreflight"
    },
    {
      "id": "evidence-owner",
      "command": "bench test --package ./internal/releaseevidence"
    },
    {
      "id": "publication",
      "command": "bench test --package ./internal/publication"
    },
    {
      "id": "command-entry",
      "command": "bench test --package ./cmd/bench"
    },
    {
      "id": "phase-wrapper",
      "command": "bench test --check package-core-guard"
    },
    {
      "id": "real-evidence-producer-consumer",
      "command": "bench test --check release-evidence-probe"
    },
    {
      "id": "map-source",
      "command": "bench test --check decision-map-integrity"
    }
  ]
}
```

## Testing decisions

Reuse sufficient tests and production probes before adding a test.
Compilation and a complete consumer census protect the mechanical declaration cut.
Review owns the bounded absence inventory.
Do not create a parallel scanner or a test whose only purpose is a deleted declaration's absence.

The one new permanent test covers the unasserted typed refusal through `Command`.
Use the existing same-package `identityFixture`, not a copied Git fixture.
Its tag, package version, binary version, committed HEAD, and Go toolchain already agree.
Set the fixture as the command's working directory with failure-safe restoration.
Do not run that test in parallel.

Invoke publish, profile public, and phase identity.
The identity guard passes before the focused-scope refusal.
Assert exit 1 and the exact requirement diagnostic.
This is planned TestCommandFocusedPublishRequiresEvidenceAuthority in internal/releasepreflight/evidence_contract_test.go.

The existing ship probe supplies the broader validation witness.
It authenticates a committed source snapshot before building the executable that runs preflight.
The shell adapter uses the selected preflight binary.
The probe's phase stubs do not bypass the finalizer or authenticate their own output.
Its native-proof mutations reach the real evidence validator through the command.
Preserve this authority, test subject, and produced-evidence separation.

Retain the publication tests in their existing files.
Do not add a new registry fixture, live publish test, copied state-machine harness, or package exemption.
The existing command npm stub and approved-set fixture remain their owners.

### Seam diagram

    trigger: release-preflight command or existing release-evidence-probe
        |
        v
    arguments and immutable phase results --> [ releasepreflight.Command / Decide ]
                                                      |
                                                      v
                                           [ releaseevidence owner ]
                                                      |
                                                      v
                                           evidence bytes or typed refusal
        tests: existing identity fixture, focused command refusal, real ship probe

    trigger: release submit / promote / rollback
        |
        v
    approved evidence --> [ publication state machine ] --> existing registry adapter
        tests: existing command adapter tests and ordered publication events

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RF01 | 1 | Preflight evidence types come directly from releaseevidence | review-owned declaration and typed consumer census plus package compilation | A retained alias or local replacement type fails the enumerated ownership comparison |
| RF02 | 2 | Each active forward caller invokes the named owner function | review-owned six-function caller inventory plus package compilation | A missed unqualified caller fails compilation after the forwarding file is removed |
| RF03 | 3 | The forwarding file's dormant declarations do not survive | review-owned full declaration inventory and final hidden-tree census | Keeping a dormant function or requirements snapshot leaves an inventory member present |
| RF04 | 4 | Existing preflight test callers remain executable after qualification | `internal/releasepreflight/decision_test.go` (`TestDecidePreflightFromImmutableEvidence`) | A missed test type or constant fails the package build |
| RF05 | 5 | The compiled map resolves after its relocation and the forwarding-file cut | review-owned link closure and decision-map-integrity | A stale structured source fails the existing regular-file source rule |
| RF06 | 6 | The release-preflight command signature remains callable | review-owned Command consumer census plus cmd/bench compilation | Removing or changing the entry breaks its real registered caller |
| RF07 | 7 | Identical immutable phase results retain the decision verdict | `internal/releasepreflight/decision_test.go` (`TestDecidePreflightFromImmutableEvidence`) | Its green, red, empty, focused, and invalid-mode cases refuse a changed verdict |
| RF08 | 8 | Verify and publish phase names retain their registry order | review-owned fixed-input differential at PhaseNames and PhaseDefinitionFor | Reordering a mode's returned phase names changes the reference result |
| RF09 | 9 | Interrupted phase results retain interrupted terminal precedence | review-owned fixed-input differential at TerminalStatus | Moving interruption below red changes the canned mixed-status result |
| RF10 | 10 | Focused publish emits the existing typed requirement refusal | planned TestCommandFocusedPublishRequiresEvidenceAuthority in internal/releasepreflight/evidence_contract_test.go | Removing errors.As recognition changes the diagnostic kind from requirement |
| RF11 | 11 | The production command still finalizes evidence | review-owned omission probe through release-evidence-probe in internal/conformance/native_workflow_test.go | Omitting FinalizeEvidence leaves its authenticated clean case without release-index.json |
| RF12 | 12 | Canned promotion inputs produce byte-identical evidence records | review-owned fixed-input differential at PromoteEvidence | A changed field, null value, framing byte, or failure payload changes captured bytes |
| RF13 | 13 | Direct file-reader calls retain the old error and file posture | review-owned fixed-input differential at ReadRegular and ReadPackageVersion | A swallowed filesystem error or newly accepted symlink changes the reference result |
| RF14 | 14 | Release-index artifacts retain their component manifest digest binding | `internal/releaseevidence/release_index_test.go` (`TestReleaseIndexBindsComponentManifestDigest`) | Omitting the digest field loses the independently supplied digest |
| RF15 | 15 | The npm staged diagnostic wins over a held lock | `internal/publication/command_adapter_test.go` (`TestSubmitStagedNPMAdapterRefusesBeforeTheLock`) | Moving refusal after locking yields the competing lock diagnostic |
| RF16 | 16 | Fixture staging keeps its current implemented adapter and command path | review-owned source-pin comparison of RunStaged, FixtureRegistry, and runSubmit | Deleting a staged operation or routing branch changes its frozen implementation |
| RF17 | 17 | First publication records platforms before the wrapper | `internal/publication/statemachine_order_test.go` (`TestFirstPublicationRecordsPlatformsBeforeWrapper`) | A premature wrapper operation changes the independently ordered event sequence |
| RF18 | 18 | Interrupted publication resumes without republishing completed packages | `internal/publication/command_adapter_test.go` (`TestSubmitNPMAdapterResumesAfterMidSequencePublishFailure`) | Repeating a completed publish changes the npm log counts |
| RF19 | 19 | Promotion retains selected-adapter behavior | `internal/publication/command_adapter_test.go` (`TestPromoteAndRollbackHonorAdapterSelection`) | Ignoring adapter selection changes the promotion subtest's npm log |
| RF20 | 20 | Rollback retains selected-adapter behavior | `internal/publication/command_adapter_test.go` (`TestPromoteAndRollbackHonorAdapterSelection`) | Ignoring adapter selection changes the rollback subtest's npm log |
| RF21 | 21 | Public CLI and evidence export surfaces remain available | review-owned public-surface source comparison against the source pin | Removing a command or export surface changes the excluded source inventory |
| RF22 | 22 | This cut supplies no live registry qualification claim | review-owned acceptance evidence classification | Calling a local stub a live result fails the explicit evidence boundary |
| RF23 | 23 | The old and new implementations agree on the complete canned family | review-owned reference-to-candidate differential defined below | A family member omitted from the transcript leaves the exit comparison incomplete |
| RF24 | 24 | The publication ownership comment names releaseevidence | review-owned exact owner-name check at releaseIndexAuthority | Keeping releasepreflight in that owner sentence contradicts the live policy owner |
| RF25 | 25 | The complete implementation closure passes its selected file-growth rule | review-owned base-to-tip lane evidence and exact write closure | An added oversized changed file or omitted source reader blocks the chunk |
| RF26 | 26 | FT142, FT306, and broader FT368 obligations remain unresolved by this cut | review-owned roadmap and source-outcome comparison | Retiring or claiming those obligations changes the retained scope disposition |
| RF27 | 5 | The relocated map and decision answers preserve their original bytes | review-owned map and topic preimage hashes | Rewriting an answer changes a preserved source hash |
| RF28 | 14 | Evidence registries retain their current schemas and artifact paths | review-owned source-pin byte comparison of the evidence owner and its JSON registries | Editing a registered schema or path changes the frozen owner bytes |
| RF29 | 17,18,19,20 | Publication state-machine runtime bytes remain unchanged | review-owned source-pin comparison of RunFirst, RunStaged, RunPromote, and RunRollback | Removing a guard or recovery transition changes the frozen function bytes |

### Edge inventory

The audience is the Bench kit.
Linked repositories keep their public CLI and evidence contracts.
The cut changes internal Go ownership only.

The canned differential covers both modes and both supported publish profiles.
It covers full and focused scope without changing either scope's authority.
A focused publish refusal remains distinct from generic evidence failure.

Terminal inputs cover empty, all green, red, not_run, interrupted, and red-plus-interrupted results.
The immutable decision cases also cover invalid mode and focused verify.
No caller-local terminal policy is permitted.

Reader inputs cover absent, empty, ordinary, directory, and symlink paths.
Package inputs cover missing version, empty version, malformed JSON, and a valid version.
Paths include spaces and a non-ASCII name.
The differential compares returned bytes, exact error text, and errors.Is classification for missing files.
These cases preserve existing behavior rather than adding strict JSON policy.

Promotion uses absent dist and an existing real dist directory.
It also uses a non-directory dist and a symlink dist to preserve their refusal.
A fixed manifest and fixed phase results remove clocks from the byte comparison.
Green, red, and interrupted records preserve null and non-null failure fields.
Successful promotion compares the complete preflight file inventory and every file byte.

No package-variable substitution is needed for the new focused command test.
Its environment and working-directory restoration must run on failure.
The test has no live child when restoration runs.
The existing publication fixtures retain their own restoration contracts.

Publication edges retain the current first, staged fixture, failed, resumable, promoted, and rolled-back transitions.
Compare those runtime function bodies against the source pin.
The command tests retain their real operand and adapter boundaries.
The held-lock staging fixture proves refusal precedence without starting npm.

**Won't handle:** a new filesystem race guarantee — Command keeps using the existing evidence reader and promoter contracts.

**Won't handle:** new stage-cleanup or interruption policy — FinalizeEvidence retains the existing owner implementation and FT142 residuals.

**Won't handle:** stricter package identity JSON — ReadPackageVersion and package evidence retain their current caller policies.

**Won't handle:** a complete successful fixture-staging qualification claim — runSubmit and FixtureRegistry retain the current path without new qualification evidence.

**Won't handle:** live public npm staging — runSubmit keeps its unsupported capability refusal.

**Won't handle:** unknown external consumers — Command and the public evidence export surfaces remain present.

**Won't handle:** an identity key-alias deletion — checkIdentity remains unchanged apart from owner qualification, because no exact table was established.

## Ownership fences

These 13 unchanged entries are the accepted implementation fence.
The independently reviewed graph remains planning-only and grants no implementation authority.

- `internal/releasepreflight/types.go`
- `internal/releasepreflight/command.go`
- `internal/releasepreflight/decision.go`
- `internal/releasepreflight/identity.go`
- `internal/releasepreflight/vulnerability.go`
- `internal/releasepreflight/decision_test.go`
- `internal/releasepreflight/external_test.go`
- `internal/releasepreflight/identity_test.go`
- `internal/releasepreflight/evidence_contract_test.go`
- `internal/publication/approved.go`
- `specs/release-evidence-forwards/decisions/architecture-release.md`
- `specs/release-evidence-forwards/assets/migration-reference.json`
- `reviews/release-evidence-forwards.md`

The new test uses identityFixture in its existing package.
No fixture producer moves.
No canary family, anchor, command registry, injected-port record, or source-reader expectation needs a changed binding.
The structured map source repair is the sole path-dependent implementation closure.
The root coordinator owns the later C11 capture-index reconciliation, not the future runtime chunk.

The implementation may change only the ownership comment in approved.go.
All other publication runtime bytes remain fixed.
The entire releaseevidence owner remains outside the fence.
The reference asset holds future captured compatibility evidence, not an implementation-derived standing oracle.
The review pickup records the census, comparisons, omission results, and restoration.

## Out of scope

Public npm staging is a distinct capability.
Its minimum starting scope is three Go owners, one command test file, and one live qualification artifact.
That is at least 5 edits and 3 gate runs.
Its actual price requires a separate capability decision.

Public CLI removal requires a supported replacement and an external-use assessment.
A single-surface proposal starts with its command owner, help surface, and compatibility test.
That is at least 3 edits and 2 gate runs per surface.
The current command and evidence export remain the surviving callers.

The identity key-alias proposal has no established declaration.
Its first separate step is one exact census artifact and one decision-map check.
That is 1 edit and 1 gate run before any deletion can be priced.
Do not manufacture a target from the old survey.

FT142 residual changes and FT306 release qualification remain their existing capabilities.
This cut makes 0 edits to them and runs 0 live release gates.
Their future implementation prices are unknown.
The broader FT368 dashboard, harness, telemetry, and terminal decisions remain outside this outcome.

## Further notes

### Source evidence and limits

The author read the compiled release map, all its tickets, and every structured source.
The author also read full FT364, FT368, FT142, and FT306.
Older roadmap claims were checked against the current source rather than copied as current evidence.

The read set includes every preflight Go file and the concrete evidence API definitions.
It includes publication command, adapters, approved-set authority, records, plans, and state-machine functions.
It includes existing command, ordering, identity, and evidence-binding tests.
The conformance read set includes native_workflow_test.go and release_probe_fixture_test.go.
Those files establish the existing real-entry validation witness.

The author used hidden-tree caller and source-reader sweeps.
Typed `bench consumers` reports were checked against same-package references and generated command source.
`go list` confirmed the existing preflight-to-evidence import edge.
It did not execute runtime tests.

No production test, mutation probe, manual gate, benchmark, registry publication, or live qualification ran during specification.
Document checks establish document validity only.
The proposed runtime witnesses remain future work.

### Fixed-input differential

Capture the reference before migrating or deleting forwards.
Use a temporary same-package probe in the fenced contract-test file.
Keep only the focused-publish regression test after the comparison.
Do not add a persistent compatibility harness.

The enumerated family is:

1. PhaseNames for verify and publish, with each returned PhaseDefinition.
2. TerminalStatus on empty, green, red, not_run, interrupted, and mixed red/interrupted results.
3. Decide on the existing canned decision cases.
4. ReadRegular and ReadPackageVersion on the reader inputs in the edge inventory.
5. PromoteEvidence on fixed green, red, and interrupted phase results with a fixed manifest.
6. FinalizeEvidence on focused verify, focused public publish, focused bank publish, and an already-cancelled context.

For the reference, call the existing preflight names.
For the candidate, call the direct evidence-owner names and unchanged Decide.
Use identical canned values, relative paths, and normalized temporary-root labels.
Do not compare two live elapsed-time reports.
Store reference values, error classifications, file inventories, and raw record bytes in the spec-local reference asset.

The candidate must match every family member.
The comparison grades exact record bytes, not decoded field presence.
Errors.As must still recognize ReleaseIntentError.
Errors.Is must still recognize context.Canceled and missing-file errors.
The public and bank focused-publish messages remain exact.
RF23 is the exit row.

### Future verification and probes

Run the preflight, evidence, publication, and command packages through `bench test --package`.
Run `bench test --check package-core-guard` for phase and wrapper contracts.
Run `bench test --check release-evidence-probe` for the real evidence producer and consumer.
Run `bench test --check decision-map-integrity` after the structured source repair.
Run the selected lane on the complete authorized delta.

The author omission replaces the production finalizer call with a nil result.
Keep the assembled run referenced so the mutated command still compiles.
The existing release-evidence-probe must fail for its missing release-index.json.
A compile or setup failure is not this witness.
Restore the production subject and record the restored comparison.

The coordinator mutation forces concrete ReleaseIntentError recognition false at Command.
Keep the mutated command compilable.
The focused command test must fail on the wrong diagnostic kind.
Restore the subject after the red.
These probes have different failures and do not create a new policy owner.

Any new independent byte expectation needs its own demonstrated red.
The reference asset is comparison evidence, not a standing expectation copied from candidate code.
No implementation acceptance follows from package compilation alone.

### Source sentence to row map

| source clause | acceptance disposition |
|---|---|
| Remove redundant internal forwards only after their consumers move | RF01-RF04, RF06, RF23 |
| Retain tested publication behavior | RF15-RF20 |
| Keep unsupported npm staging explicit | RF15, RF22 |
| Preserve error identity | RF10, RF13, RF23 |
| Preserve record bytes | RF12, RF14, RF23 |
| Preserve publication order | RF17 plus fixed publication runtime comparison |
| Retain public surfaces with unknown external consumers | RF21 |
| Evidence validation omission must fail at production entry | RF11 |
| Review FT142 and FT306 constraints | RF22, RF26 |
| Identity key-alias table was not established | named exclusion, RF26 |

### Pre-review proof checklist

- Cited symbols: the active and dormant forward inventories resolve at the pinned declarations.
- Import edges: preflight already imports releaseevidence, and go list resolves both packages.
- Source-row clauses and occurrences: the table above covers the compiled map and its three answers.
- Promised field labels: evidence schema labels are unchanged, including component_manifest_sha256.
- Changed-function callers: the six active forwarding functions and Command/Decide are enumerated above.
- Copy survival: RF01 and RF03 compare the complete declaration inventory without adding a standing scanner.
- Rendered-shape readers: no runtime message moves, and the focused refusal keeps its exact existing text.
- Pin operators: differential equality is byte equality for records and value equality for canned results.
- Entry reads: current root, environment, Git, filesystem, and toolchain reads remain behind the existing command entry.
- Derived expectations: existing ordering and digest tests own their independent values, and future mutations must demonstrate their relevant reds.
- Consolidated rules: none, because evidence policy already has one implementation owner.
- Quantified obligations: every forward declaration belongs to the named inventory, and both supported profiles belong to the differential.
- Workflow-step writes: none.
- Canary execution roots: none change.
- Anchor claims: none change.
- Writer protocols: existing evidence promotion remains unchanged.
- Posture-change fixtures: none, because this spec preserves posture.
- Count and row-list readers: no executable inventory count changes.
- Fixture-builder readers: identityFixture remains in place and the new test reuses it.
- Flagged additions: the focused-publish diagnostic test and bounded reference capture support existing preservation obligations.
- Unestablished claim: the identity key-alias table has no verified declaration or removal fence.

### Acceptance and graph record

[The spec-local review record](assets/spec-review.md) preserves actual SPEC acceptance before slicing and the bounded source checks.
All 29 row lines and 13 fence entries remain byte-identical to the accepted source.
All four moved-map files remain byte-identical; the citation correction changes no behavior or authority.
No runtime, mutation, publication, or live qualification result is claimed.

### Approval table

| subject | proposal | disposition |
|---|---|---|
| implementation line | gpt-5.6-sol/high for RF1 | SPEC accepted at 46206d53; graph accepted at 350bf2b6 |
| seams | direct owner, existing command fixtures, real release-evidence-probe | SPEC accepted at 46206d53; graph accepted at 350bf2b6 |
| coverage | RF01-RF29 with named edge exclusions | SPEC accepted at 46206d53; graph accepted at 350bf2b6 |
| ownership fences | exact caller, map source, reference, and review closure | SPEC accepted at 46206d53; graph accepted at 350bf2b6 |
| capability scope | retain publication and public export boundaries | SPEC accepted at 46206d53; graph accepted at 350bf2b6 |
| tickets | one complete RF1 vertical graph after SPEC acceptance | independent graph accepted at 350bf2b6 |

### Numbered breakdown

| ticket | title | Blocked by | delivered outcome |
|---|---|---|---|
| 1-consume-release-evidence-directly.md | Consume release evidence directly | none | Complete caller migration and forwarding removal with exact evidence, typed command refusal, and publication preservation |
