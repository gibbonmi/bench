# Production test seam policy

Status: staged

Roadmap: FT343

Decision source: [Production test-seam admission](decisions/architecture-test-seams.md), ready compiled map with resolved tickets 1–3.

Verification log: 3 iterations to accept — two independent spec passes and one independent ticket pass.
Retained xhigh `deepen_execution` held PS-R1 in the initial spec review, then confirmed its repair before slicing at `433159870da03aab3fcb9fcecf355d0d13d8bb55`.
Accepted spec SHA-256: `1763a028adaf0bb1935d47976d551f33dfad4a01a68a15911b7d2701ada8f589`.
Independent high reviewer `adoption_spec_review` accepted the graph at `f14dc99b8a785393a59fe843cd126df3ed51ec6c`, with all three axes at 0 findings and confidence 9.
Reviewed graph spec SHA-256: `85f26f66703894b26f9c17dec70bc4a153a5906d768692a60aa6702cc92e5c77`.

Retained `gpt-6.1-sol / high` authored the repair and graph and now closes metadata, with at most two metadata/validation iterations authorized.
This close preserves all 66 row lines, 64 fence entries, and ticket bytes, without claiming runtime evidence or granting build authority.
The spec-local review record retains the complete spec-before-ticket chronology and Git sources.

## Problem

The injected-port check finds named type ports in five packages.
It does not find package variable overrides, setter aliases, runtime test branches, or environment fault controls.
A new hook can therefore alter production execution without entering the registry.
A registry row can also name a test without explaining its failure purpose or safe isolation.

FT343's premise about unused preflight hooks is stale.
REF and RACE have test consumers, and the native release probe constructs a phase override family that includes VULNERABILITY.
Removing those controls now can remove evidence instead of removing duplication.
The telemetry test branch also protects the user's fallback home.

## Solution

Extend the existing injected-port owner to derive production test seams from declarations and their uses.
Keep one admission registry with a production owner, actual consumer, failure purpose, and isolation evidence for each retained seam.
Use real dependencies when they can produce the required failure reliably.
Prefer collaborators attached to one operation when fault injection remains necessary.

Retain useful controls until an equivalent test drives the same production entry and observes the same failure.
Do not admit a test-only path that reports success by skipping production validation or authority.
The audit checks source evidence and registry consistency.
Named behavior tests and independent review establish whether the evidence is sufficient.

## User stories

Line: gpt-5.6-sol / high.

Implementation-line reason: the oracle must resolve aliases and preserve failure evidence across real callers.
The decided policy is precise, but the expanded source derivation and isolation review remain material risks.
The profile's mid binding and the kit leverage override support this line.

Harder chunks: TS1 and TS3.

Admission and discovery:

1. As a maintainer, I want existing named ports to remain visible, so that producer junction evidence stays required.
2. As a maintainer, I want directly replaced collaborators to be found, so that avoiding a setter cannot avoid review.
3. As a maintainer, I want scalar fault bounds to be found, so that a duration or size override cannot stay invisible.
4. As a maintainer, I want setter aliases to reach their real owner, so that wrapper names cannot conceal a seam.
5. As a maintainer, I want runtime test branches to be found, so that test execution cannot silently change authority.
6. As a maintainer, I want dynamic environment controls to be found, so that a constructed hook name cannot avoid registration.
7. As a maintainer, I want ordinary functions to remain ordinary, so that naming patterns do not create false policy findings.
8. As a maintainer, I want malformed source to refuse audit success, so that a failed census cannot resemble an empty census.
9. As a maintainer, I want stale registry entries to fail, so that removed controls cannot leave persuasive green evidence.

Evidence and isolation:

10. As a reviewer, I want a production owner for each seam, so that a test double cannot become the policy owner.
11. As a reviewer, I want a declared test consumer for each fault exception, so that an unused advertisement cannot pass.
12. As a reviewer, I want a concrete failure purpose, so that convenience alone cannot justify a global override.
13. As a reviewer, I want isolation evidence, so that one test cannot change another operation's dependency.
14. As a reviewer, I want failure-safe restoration, so that a fatal assertion cannot leave an override installed.
15. As a reviewer, I want live consumers joined before restoration, so that a goroutine cannot read restored state too early.
16. As a reviewer, I want production validation to remain active, so that a hook cannot turn an invalid request into success.
17. As a maintainer, I want instance-local collaborators admitted, so that a new fault test does not require a global hook.
18. As a maintainer, I want sufficient real tests first, so that reducing test cost cannot reduce failure coverage.

Current controls:

19. As a release maintainer, I want REF evidence retained, so that identity disagreement remains observable.
20. As a release maintainer, I want RACE evidence retained, so that child environment isolation remains observable.
21. As a release maintainer, I want every release phase override accounted for, so that the native probe keeps its production composition.
22. As a release maintainer, I want DATE classified separately, so that missing in-tree use does not imply safe external removal.
23. As a release maintainer, I want EVIDENCE_READY_FILE classified separately, so that a timing control cannot vanish through a family claim.
24. As a test author, I want fallback-home protection retained, so that tests cannot write into the user's Bench home.
25. As a maintainer, I want setter retirement supported by a caller census, so that public wrappers cannot disappear on a name-only search.
26. As a maintainer, I want source examples excluded accurately, so that test support imports do not become production branch findings.

Oracle ownership:

27. As a maintainer, I want independent omission controls, so that a missing registration makes the gate red.
28. As a maintainer, I want the existing check and canary retained, so that the broader policy keeps its operational entry.
29. As a maintainer, I want one source scanner, so that two derivations cannot drift.
30. As a maintainer, I want the audit to use the graded tree, so that a fixture cannot borrow evidence from the working checkout.
31. As a maintainer, I want complete green implementation chunks, so that a migration never leaves the tree partially admitted.

## Implementation decisions

### Terms and eligibility

A production test seam is a controlled dependency or fault input on a production execution path.
Its consumer must exercise that path, rather than bypass it with a test-only success return.
An admission is the reviewed disposition of one derived seam.
A restore function is cleanup machinery, not evidence that its callers finish safely.

The production census starts from repository Go source under `cmd/` and `internal/`.
It reads non-test declarations and test uses from the same graded root.
It includes platform source variants and unexported declarations.
It excludes fixture source outside these production roots and declarations used only by test support.
Test support is established by source dependencies and calls, not by a directory name or a `testing` import alone.
A source file reachable from a non-test command entry remains production even when it imports `testing`.

Keep the existing named-port shapes: nonempty interfaces, named function types, and structs whose fields are all functions.
An injection site remains a parameter, pointer parameter, type assertion, or type switch.
Apply that derivation to eligible production packages rather than a closed five-package discovery list.
A named shape that is never injected remains a negative control.
An ordinary function declaration is not a seam merely because its name contains `ForTest`.

A mutable collaborator qualifies when a test replaces a package variable that a production path reads or calls.
This includes function values, anonymous function values, collaborator structs, scalar bounds, and combined assignments.
Resolve the assignment target to its declaration, so local shadowing and incidental fixture values do not qualify.
A production function variable with an explicit fault or override entry also qualifies before any test registers it.
Ordinary mutable production data without a test override or fault entry does not qualify.

A setter qualifies through its effect on an eligible production collaborator.
Resolve forwarding functions and aliases to that owner, including `ForTesting` variants and unexported wrappers.
A pure validation wrapper such as `ValidateTarballForTesting` is an exposure of existing validation, not a mutable setter.
Record such exposures separately when they support a registered seam.
Do not invent a setter finding from its spelling alone.

A test runtime branch qualifies when test runtime state changes a production dependency, write target, or validation path.
Detect the imported `testing.Testing` call through its import binding, including an import alias.
An import of `testing` without such a production branch does not qualify.
Resolve local forwarding predicates rather than checking one identifier spelling.

An environment hook qualifies when a test or explicit fault entry controls a production operand or failure point.
Resolve constants, aliases, and bounded string construction for the hook identity.
The release phase family takes its members from the existing phase registry and the actual hook reader.
Do not create a second phase-name list in the audit.
Unresolved dynamic construction produces a diagnostic, rather than silently grading no hook.

### One audit owner and registry contract

Relocate the current audit into `internal/conformance/injectedports/` and extend that owner there.
Keep `checkInjectedPortRegistry` as the existing conformance entry.
Remove the relocated walkers and unreachable helper copies from the old file.
Keep all seam derivation in the relocated owner, including consumer lookup and alias resolution.
Reuse the existing release registry reader for family data instead of adding a registry parser.
This is a move of the current scanner, followed by an extension of its policy.

The planned boundary is `Check(root string, policy Policy) []string`.
`Policy` carries admission rows and the existing required-package expectations for partial fixture roots.
Each row has a stable source identity, kind, production owner, consumer references, failure purpose, and isolation evidence.
Kinds distinguish named ports, variable overrides, setter aliases, runtime branches, and environment hooks.
Aliases reference their canonical seam instead of repeating its policy facts.

The root wrapper supplies the registry and calls this boundary once.
It can supply dynamic family membership from the existing graded-root release registry reader.
It cannot substitute the live binary's embedded registry for the fixture's registry.

Existing diagnostics remain exact, including the unregistered, missing-test, empty-exemption, zero-inventory, and orphan-row prefixes.
Keep the existing unregistered-port canary's full diagnostic for its named-port case.
New policy defects receive distinct diagnostic prefixes in the same owner.
Sort and deduplicate diagnostics as today.

A missing production package remains benign in a partial fixture or adopting tree.
A present required package with no derived eligible seam retains the zero-inventory refusal.
A present but unreadable source package cannot use that absence exception.
A stale row in a present package must refuse even when that package has no other seams.
A root with eligible source and no registry rows must report unregistered seams.

A retained fault exception names an actual top-level test declaration and its production entry.
It states the failure that the real dependency cannot reliably produce.
Its isolation record names the construction or synchronization and the cleanup owner.
A hook with no established consumer does not gain an automatic exception.
Retain its compatibility reader while adding sufficient consumer evidence or obtaining a specific reviewed retirement decision.

The old nonblank exemption path remains for real external dependencies, such as live npm publication.
Such a row names its reviewed external evidence and limitation, rather than pretending to name an in-tree fault consumer.

The audit rejects missing fields, nonexistent consumers, duplicate identities, unresolved aliases, and missing isolation evidence.
It also rejects a consumer that never reaches the declared hook or canonical owner.
The check does not prove test purpose from a function name.
Review and the consumer's observable failure test own that judgment.
A nonblank paragraph cannot certify restoration or authority preservation.

### Isolation and sufficient evidence

Prefer an actual process, temporary repository, or filesystem object when it reliably produces the named failure.
When that is insufficient, inject the smallest collaborator on one operation or instance.
Keep its production adapter and validation owner intact.
A local read collaborator, for example, can inject a read failure while the bounds owner still classifies the result.
This policy does not require a future FT365 syntax API or another plan's read port.

For a retained package global, name its serial ownership and restoration mechanism.
Tests that install it cannot run in parallel with any consumer of that same global.
Register failure-safe cleanup before the first fallible assertion.
Stop or join every asynchronous consumer before that cleanup restores the previous value.
The restore closure returned by a setter establishes none of these ordering facts by itself.

Demonstrate restoration through an omission that leaves a distinguishable replacement installed.
Demonstrate join ordering with a consumer held behind a deterministic barrier.
The witness observes the replacement while the consumer is live, then the prior value after cleanup.
Deleting the join must make that witness fail.
Use one sufficient witness per isolation mechanism, rather than one redundant test per registry row.

An injected dependency can report failure or supply controlled input.
It cannot remove production identity, schema, filesystem, or publication validation from the exercised path.
The corresponding entry test supplies an invalid input with the hook active and still observes refusal.
Replacing a validator with unconditional success is an inadmissible authority bypass.
A declaration-exists check cannot establish this property.

### Current census and dispositions

The source baseline is `37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68` for production.
The planning checkout starts at `a9c395e77fec36d1f60f6d057c654aecf5720e24`.
The following table records observed forms, rather than a frozen count of all future registrations.
Implementation must refresh declaration and caller identities before changing any control.

| Source family | Observed definition and consumer | Required disposition |
|---|---|---|
| Existing named ports | `internal/conformance/injected_ports_registry_test.go` names seven rows in five packages | Preserve real-producer tests and the explicit publication limitation |
| Git bounds | `internal/git/git.go` setters reach `checker_junction_test.go`, `worktree_admin_enum_test.go`, `admin_readers_test.go`, and `internal/status/status_render_test.go` | Register canonical variables and setter entries with serial cleanup evidence |
| Snapshot movement | `internal/diff/snapshot.go` has a setter and direct assignments in diff tests, with preflight consumers | Preserve deterministic movement refusals through the real snapshot entry |
| Evidence collection | `internal/preflight/review.go` supplies `reviewEvidenceObserved` to `review_charge_test.go` | Register the observer's counting purpose and cleanup |
| Response bound | `internal/preflight/evidencecmd/bound.go` supplies a scalar override to `evidence_budget_test.go` | Preserve the public over-bound refusal witness |
| Lock windows | `internal/intent/intent.go` and `internal/handoffdoc/store.go` setters reach transaction and store tests | Preserve lock timeout evidence and test environment restoration |
| Inspection | `internal/sessioninspect/sessioninspect.go` has scalar setters plus direct `phases` and `runInspect` replacements | Register each canonical collaborator with completed-consumer evidence |
| Release evidence | `artifact_evidence.go`, `evidence_promotion.go`, `release_index.go`, and `types.go` expose four mutable setters | Establish each actual consumer or equivalent owner failure witness before admission |
| Release aliases | `internal/releasepreflight/types.go` forwards setter and validator exposures | Resolve canonical owners without treating an unused wrapper as a live test |
| REF | `internal/releasepreflight/identity_test.go:43` installs the exact ref through `identityFixture` | Retain identity disagreement evidence |
| RACE | `internal/releasepreflight/external_test.go:37` installs a real child probe | Retain child home isolation evidence |
| Phase family | `internal/conformance/native_workflow_test.go:321` constructs overrides for gate, race, vet, vulnerability, artifacts, and smoke | Register the family through the existing phase registry |
| VULNERABILITY | `internal/releasepreflight/vulnerability.go:27` reads its scanner override | Retain the native probe's scanner composition |
| DATE | `internal/releasepreflight/vulnerability.go:75` reads the clock override | Retain compatibility and add an entry test for the chosen date |
| EVIDENCE_READY_FILE | `internal/releaseevidence/release_evidence.go` reads the readiness marker | Retain compatibility and add a controlled readiness entry test |
| Telemetry runtime branch | `internal/otelrecord/provider.go:51` calls `testing.Testing`; processor tests cover fallback and explicit homes | Retain fallback-home refusal until an equivalent isolated constructor exists |
| Link faults | `internal/adopt/transaction.go` reads `BENCH_LINK_FAULT`; system compatibility tests use interruption | Register fault and interruption purposes with child-local environment |
| Verdict windows | `internal/bounds/bounds.go` reads `UnboundedWaitsEnv`; bounds and kit environment tests consume it | Preserve the bounds owner's distinct verdict and fixed-window policies |
| Test support environment | `internal/gocache/cleanprobe/cleanprobe.go` reads its answer path from tests | Treat the test-only helper as a negative production control |
| Direct replacements | The families below replace production variables without setters | Register eligible collaborators and exclude local fixture assignments |

The direct-assignment sweep found eligible candidates in these source families:

- `cmd/bench`: `commandRegistry`, `commitChain`, and `gatePhasesCommand` have named command tests.
- `internal/adopt`, `internal/freshness`, and `internal/outline`: synchronization, rename, locking, and opening collaborators have failure tests.
- `internal/capability`: `stdout` is replaced by its output test.
- `internal/consumers`: `load` is replaced by command, loader, and graph tests.
- `internal/coverage`, `internal/harnesstranscript`, and `internal/refresh`: scalar limits have boundary or timeout tests.
- `internal/gate`: timeout, run-binary factories, and log-path observation have entry tests.
- `internal/guards`: discovery, inspection, and timeout variables have cancellation and cleanup tests.
- `internal/roadmap`: `openRetroRoot` has a filesystem refusal test.
- `internal/shift`: `shiftFault` and `timeNow` have failure and refresh tests.
- `internal/skillsindex`: replacement barrier and rename variables have cleanup tests.
- `internal/testreport`: run-binary selection and package listing have cancellation and selection tests.

The following exact declarations and consumers were observed in that sweep.
A helper consumer remains a helper until the implementation traces its calling test.
This table supplies evidence for the refresh, not automatic admission.

| Production declaration | Observed consumer |
|---|---|
| `cmd/bench/main.go::commandRegistry` | `cmd/bench/help_inventory_test.go::TestHelpRendersPublicCommandRegistryRows` |
| `cmd/bench/commit_chain.go::commitChain` | `cmd/bench/commit_chain_test.go::runCommitChain` |
| `cmd/bench/main.go::gatePhasesCommand` | `cmd/bench/main_test.go::TestRunGatePhasesDispatchesToCommand` |
| `internal/adopt/transaction.go::syncDirectory` | `internal/adopt/adopt_test.go::TestPromoteAllRollsBackOnDestinationSyncFailure` |
| `internal/capability/capability.go::stdout` | `internal/capability/capability_test.go::TestCapabilityWritesLineBeforeSkip` |
| `internal/consumers/loader.go::load` | `internal/consumers/blast_edges_test.go::TestIdenticalPairAnswersWithoutLoadingPackages` |
| `internal/coverage/citation_execution.go::packageLoadTimeout` | `internal/coverage/citation_execution_test.go::shrinkPackageLoadTimeout` |
| `internal/diff/snapshot.go::snapshotAfterRead` | `internal/diff/command_test.go::TestCommandRetriesThenRefusesSnapshotDrift` |
| `internal/freshness/freshness_publish.go::replacePublicationFile` | `internal/freshness/freshness_publish_test.go::TestPublishRestoresPriorPairWhenSealPromotionFails` |
| `internal/freshness/publication_lock.go::publicationFlock` | `internal/freshness/publication_lock_test.go::TestPublicationProcess` |
| `internal/gate/gate.go::gateTimeout` | `internal/gate/timeout_recovery_count_test.go::TestGateRunTimeoutInvalidatesOldEvidence` |
| `internal/gate/lane.go::laneRunBinary` | `internal/gate/lane_run_test.go::laneRefusesArgv` |
| `internal/gate/prospective.go::prospectiveRunBinary` | `internal/gate/prospective_owner_test.go::TestProspectiveBuildRefusalLeavesNoBundle` |
| `internal/gate/run_log.go::gateLogPathIgnored` | `internal/gate/run_log_prune_test.go::stubGateLogPathIgnored` |
| `internal/guards/guards.go::enumerateGuards and inspectGuard` | `internal/guards/guards_cleanup_test.go::TestScanWaitsForCancelledWorkerCleanup` |
| `internal/guards/guards.go::guardScanTimeout` | `internal/guards/guards_test.go::TestCommandPreservesCheckedInEnumerationTimeoutPrimaryResponse` |
| `internal/harnesstranscript/read.go::maxRecordLine` | `internal/harnesstranscript/transcript_test.go::TestOversizedLineIsOneSkippedEvent` |
| `internal/outline/read.go::openOutlineFile` | `internal/outline/outline_test.go::TestCommandNamesUnreadableSkip` |
| `internal/refresh/refresh.go::refreshTimeout` | `internal/refresh/refresh_test.go::TestRefreshFailureAndTimeoutAreNonfatalAndDetailed` |
| `internal/roadmap/retro.go::openRetroRoot` | `internal/roadmap/retro_test.go::TestRetroContainsDestinationComponentReplacement` |
| `internal/sessioninspect/sessioninspect.go::phases` | `internal/sessioninspect/sessioninspect_test.go::TestInspectDeadlineWarnsAndReturnsZero` |
| `internal/sessioninspect/sessioninspect.go::runInspect` | `internal/sessioninspect/sessioninspect_test.go::TestCommandInstallsTenSecondDeadline` |
| `internal/shift/fault.go::shiftFault` | `internal/shift/fault_test.go::TestShiftHelperProcess` |
| `internal/shift/session.go::timeNow` | `internal/shift/refresh_test.go::shiftCollisionFixture` |
| `internal/skillsindex/skillsindex.go::preReplacementBarrier` | `internal/skillsindex/command_test.go::TestWriteBarrierHelperProcess` |
| `internal/skillsindex/skillsindex.go::renameFile` | `internal/skillsindex/skillsindex_reference_test.go::TestRenameFailureLeavesNoResidueAndKeepsReferenceBytes` |
| `internal/testreport/command.go::selectRunBinary` | `internal/testreport/check_test.go::TestNamedCheckRunsOnlyRegisteredDevScope` |
| `internal/testreport/selection.go::listCurrentPackages` | `internal/testreport/selection_test.go::TestChangedPackageSelectionRefusalMatrix` |

This lexical sweep was a discovery aid, not a semantic completeness proof.
It also returned strings and ordinary values that are not assignments.
The implementation's declaration resolution and independent negative controls must remove those false candidates.
Publication's `clock` is another function-valued candidate whose caller disposition requires the refreshed semantic census.
No claim that all present globals are safe follows from this table.

No in-tree DATE or EVIDENCE_READY_FILE setter was established by the source sweep.
That result does not establish absence of external use.
Their compatibility-preserving entry tests precede any later removal decision.
Retirement of release forwarding APIs remains a separate release-cuts outcome.
This spec requires their alias census and policy disposition, but grants no blanket forwarding deletion.

## Implementation chunks

These ticket proposals preserve the four accepted green outcomes.
The graph awaits independent approval and grants no implementation authorization.
Each chunk refreshes the applicable caller census and carries its production adapters and tests together.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| TS1 / 01-relocate-the-existing-audit.md | Same audit owner relocated with legacy registry and canary behavior intact | PS01–PS03, PS12–PS14, PS41–PS43, PS50–PS55, PS61 | Legacy audit tests, injected-port-registry, canary | yes |
| TS2 / 02-admit-variable-and-setter-callers.md | Variable and setter discovery with admitted synchronous caller families | PS04–PS08, PS15–PS21, PS44, PS57 | Derivation and evidence suites plus actual caller packages | no |
| TS3 / 03-admit-runtime-and-environment-controls.md | Runtime and environment controls with safe asynchronous isolation and refusal evidence | PS09–PS11, PS22–PS40, PS56, PS58–PS60, PS66 | Policy suite, release entries, telemetry tests, isolation witnesses | yes |
| TS4 / 04-close-census-and-omission-proof.md | Whole-source census closure with independent omissions and no scanner copies | PS45–PS49, PS62–PS65 | Audit mutation suite, existing conformance check, canary family | no |

TS1 lands a usable owner and the existing real wrapper in the same chunk.
TS2 and TS3 land registration and sufficient caller evidence with each newly enforced form.
They cannot temporarily enable a rule that makes the real tree red.

TS3 owns PS66 when it first consumes the graded-root phase registry.
It updates the InputCatchAll binding and the profile advertisement in that same chunk.
Its closure includes the four anchor holders, AXI command advertisement holders, and twelve complete fixture units listed in the fence.
TS4 proves the production census is not restricted to remembered registry packages.

The original TS1–TS4 IDs map directly to their same-numbered ticket files.
Their row partitions remain 16, 14, 27, and 9; no predicate moves or changes.
Shared audit and registry writes serialize all four tickets in that order.
Each predecessor chunk receives independent review before its consumer starts.

First-use obligations apply before each ticket's own green checkpoint.
TS1 removes relocated scanner copies and proves the real wrapper immediately, although TS4 owns PS49's final census.
TS2 proves sufficient isolation and ordinary-data and shadowing negatives when it first enforces variable and setter admission.
TS3 proves every newly admitted runtime or environment control's consumer, authority, and isolation before enabling enforcement.
TS4 owns terminal aggregate rows; it cannot supply evidence missing from an earlier checkpoint.
Ticket slicing follows the independent acceptance pinned in the verification log.

### Planned completion checks

These commands and mutation witnesses are future implementation obligations, not executed evidence.
The final native and system checks retain the chosen binary and the existing `BENCH_KIT` route.
Every probe must record the changed source, exact failing assertion or diagnostic, byte-identical restoration, and the same focused green after restoration.
A setup failure, unrelated red, or failed restoration does not prove an omission.
Before implementation dispatch, a version-2 plan must bind each ticket to its fresh author.
The approved implementation line remains `gpt-5.6-sol / high`.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "TS1",
      "tickets": [
        "01-relocate-the-existing-audit.md"
      ],
      "verification": [
        {
          "id": "owner",
          "command": "bench test --package ./internal/conformance/injectedports"
        },
        {
          "id": "real-check-and-canary",
          "command": "bench test --package ./internal/conformance"
        }
      ]
    },
    {
      "id": "TS2",
      "tickets": [
        "02-admit-variable-and-setter-callers.md"
      ],
      "verification": [
        {
          "id": "owner",
          "command": "bench test --package ./internal/conformance/injectedports"
        },
        {
          "id": "real-check-and-canary",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "caller-cmd-bench",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "caller-internal-adopt",
          "command": "bench test --package ./internal/adopt"
        },
        {
          "id": "caller-internal-assessment",
          "command": "bench test --package ./internal/assessment"
        },
        {
          "id": "caller-internal-capability",
          "command": "bench test --package ./internal/capability"
        },
        {
          "id": "caller-internal-chargeevidence",
          "command": "bench test --package ./internal/chargeevidence"
        },
        {
          "id": "caller-internal-consumers",
          "command": "bench test --package ./internal/consumers"
        },
        {
          "id": "caller-internal-coverage",
          "command": "bench test --package ./internal/coverage"
        },
        {
          "id": "caller-internal-diff",
          "command": "bench test --package ./internal/diff"
        },
        {
          "id": "caller-internal-freshness",
          "command": "bench test --package ./internal/freshness"
        },
        {
          "id": "caller-internal-gate",
          "command": "bench test --package ./internal/gate"
        },
        {
          "id": "caller-internal-git",
          "command": "bench test --package ./internal/git"
        },
        {
          "id": "caller-internal-gitguard",
          "command": "bench test --package ./internal/gitguard"
        },
        {
          "id": "caller-internal-guards",
          "command": "bench test --package ./internal/guards"
        },
        {
          "id": "caller-internal-handoffdoc",
          "command": "bench test --package ./internal/handoffdoc"
        },
        {
          "id": "caller-internal-harnesstranscript",
          "command": "bench test --package ./internal/harnesstranscript"
        },
        {
          "id": "caller-internal-intent",
          "command": "bench test --package ./internal/intent"
        },
        {
          "id": "caller-internal-outline",
          "command": "bench test --package ./internal/outline"
        },
        {
          "id": "caller-internal-preflight",
          "command": "bench test --package ./internal/preflight"
        },
        {
          "id": "caller-internal-preflight-evidencecmd",
          "command": "bench test --package ./internal/preflight/evidencecmd"
        },
        {
          "id": "caller-internal-publication",
          "command": "bench test --package ./internal/publication"
        },
        {
          "id": "caller-internal-refresh",
          "command": "bench test --package ./internal/refresh"
        },
        {
          "id": "caller-internal-releaseevidence",
          "command": "bench test --package ./internal/releaseevidence"
        },
        {
          "id": "caller-internal-releasepreflight",
          "command": "bench test --package ./internal/releasepreflight"
        },
        {
          "id": "caller-internal-roadmap",
          "command": "bench test --package ./internal/roadmap"
        },
        {
          "id": "caller-internal-sessioninspect",
          "command": "bench test --package ./internal/sessioninspect"
        },
        {
          "id": "caller-internal-shift",
          "command": "bench test --package ./internal/shift"
        },
        {
          "id": "caller-internal-skillsindex",
          "command": "bench test --package ./internal/skillsindex"
        },
        {
          "id": "caller-internal-status",
          "command": "bench test --package ./internal/status"
        },
        {
          "id": "caller-internal-testreport",
          "command": "bench test --package ./internal/testreport"
        }
      ]
    },
    {
      "id": "TS3",
      "tickets": [
        "03-admit-runtime-and-environment-controls.md"
      ],
      "verification": [
        {
          "id": "owner",
          "command": "bench test --package ./internal/conformance/injectedports"
        },
        {
          "id": "real-check-and-canary",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "caller-cmd-bench",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "caller-internal-anchors",
          "command": "bench test --package ./internal/anchors"
        },
        {
          "id": "caller-internal-conformance-registry",
          "command": "bench test --package ./internal/conformance/registry"
        },
        {
          "id": "caller-internal-assessment",
          "command": "bench test --package ./internal/assessment"
        },
        {
          "id": "caller-internal-diff",
          "command": "bench test --package ./internal/diff"
        },
        {
          "id": "caller-internal-guards",
          "command": "bench test --package ./internal/guards"
        },
        {
          "id": "caller-internal-bounds",
          "command": "bench test --package ./internal/bounds"
        },
        {
          "id": "caller-internal-otelrecord",
          "command": "bench test --package ./internal/otelrecord"
        },
        {
          "id": "caller-internal-releaseevidence",
          "command": "bench test --package ./internal/releaseevidence"
        },
        {
          "id": "caller-internal-releasepreflight",
          "command": "bench test --package ./internal/releasepreflight"
        },
        {
          "id": "native-phase-composition",
          "command": "bench test --check release-evidence-probe"
        },
        {
          "id": "real-interruption",
          "command": "bench test --check system"
        },
        {
          "id": "restore-omission",
          "command": "bench test --package ./internal/diff",
          "probe": "omit failure-safe restoration; require the distinguishable installed replacement witness to fail"
        },
        {
          "id": "join-omission",
          "command": "bench test --package ./internal/guards",
          "probe": "remove join before restoration; require the held live consumer witness to fail"
        }
      ]
    },
    {
      "id": "TS4",
      "tickets": [
        "04-close-census-and-omission-proof.md"
      ],
      "verification": [
        {
          "id": "owner",
          "command": "bench test --package ./internal/conformance/injectedports"
        },
        {
          "id": "real-check-and-canary",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "registry-omission",
          "command": "bench test --package ./internal/conformance",
          "probe": "delete a real registration while preserving eligible source; require unregistered refusal"
        },
        {
          "id": "alias-omission",
          "command": "bench test --package ./internal/conformance/injectedports",
          "probe": "disable canonical alias resolution; require the independently authored alias fixture to fail"
        },
        {
          "id": "setter-omission",
          "command": "bench test --package ./internal/conformance/injectedports",
          "probe": "disable setter discovery; require the independent setter finding to fail"
        },
        {
          "id": "runtime-omission",
          "command": "bench test --package ./internal/conformance/injectedports",
          "probe": "disable runtime branch discovery; require the independent runtime finding to fail"
        },
        {
          "id": "consumer-omission",
          "command": "bench test --package ./internal/conformance/injectedports",
          "probe": "remove required consumer evidence; require the admission refusal"
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "coverage",
      "command": "bench coverage --check specs/production-test-seam-policy/spec.md"
    },
    {
      "id": "complete-owner",
      "command": "bench test --package ./internal/conformance/injectedports"
    },
    {
      "id": "complete-conformance",
      "command": "bench test --package ./internal/conformance"
    },
    {
      "id": "registry",
      "command": "bench test --package ./internal/conformance/registry"
    },
    {
      "id": "anchors",
      "command": "bench test --package ./internal/anchors"
    },
    {
      "id": "native-phase-composition",
      "command": "bench test --check release-evidence-probe"
    },
    {
      "id": "interruption",
      "command": "bench test --check system"
    },
    {
      "id": "caller-cmd-bench",
      "command": "bench test --package ./cmd/bench"
    },
    {
      "id": "caller-internal-adopt",
      "command": "bench test --package ./internal/adopt"
    },
    {
      "id": "caller-internal-assessment",
      "command": "bench test --package ./internal/assessment"
    },
    {
      "id": "caller-internal-capability",
      "command": "bench test --package ./internal/capability"
    },
    {
      "id": "caller-internal-chargeevidence",
      "command": "bench test --package ./internal/chargeevidence"
    },
    {
      "id": "caller-internal-consumers",
      "command": "bench test --package ./internal/consumers"
    },
    {
      "id": "caller-internal-coverage",
      "command": "bench test --package ./internal/coverage"
    },
    {
      "id": "caller-internal-diff",
      "command": "bench test --package ./internal/diff"
    },
    {
      "id": "caller-internal-freshness",
      "command": "bench test --package ./internal/freshness"
    },
    {
      "id": "caller-internal-gate",
      "command": "bench test --package ./internal/gate"
    },
    {
      "id": "caller-internal-git",
      "command": "bench test --package ./internal/git"
    },
    {
      "id": "caller-internal-gitguard",
      "command": "bench test --package ./internal/gitguard"
    },
    {
      "id": "caller-internal-guards",
      "command": "bench test --package ./internal/guards"
    },
    {
      "id": "caller-internal-handoffdoc",
      "command": "bench test --package ./internal/handoffdoc"
    },
    {
      "id": "caller-internal-harnesstranscript",
      "command": "bench test --package ./internal/harnesstranscript"
    },
    {
      "id": "caller-internal-intent",
      "command": "bench test --package ./internal/intent"
    },
    {
      "id": "caller-internal-outline",
      "command": "bench test --package ./internal/outline"
    },
    {
      "id": "caller-internal-preflight",
      "command": "bench test --package ./internal/preflight"
    },
    {
      "id": "caller-internal-preflight-evidencecmd",
      "command": "bench test --package ./internal/preflight/evidencecmd"
    },
    {
      "id": "caller-internal-publication",
      "command": "bench test --package ./internal/publication"
    },
    {
      "id": "caller-internal-refresh",
      "command": "bench test --package ./internal/refresh"
    },
    {
      "id": "caller-internal-releaseevidence",
      "command": "bench test --package ./internal/releaseevidence"
    },
    {
      "id": "caller-internal-releasepreflight",
      "command": "bench test --package ./internal/releasepreflight"
    },
    {
      "id": "caller-internal-roadmap",
      "command": "bench test --package ./internal/roadmap"
    },
    {
      "id": "caller-internal-sessioninspect",
      "command": "bench test --package ./internal/sessioninspect"
    },
    {
      "id": "caller-internal-shift",
      "command": "bench test --package ./internal/shift"
    },
    {
      "id": "caller-internal-skillsindex",
      "command": "bench test --package ./internal/skillsindex"
    },
    {
      "id": "caller-internal-status",
      "command": "bench test --package ./internal/status"
    },
    {
      "id": "caller-internal-testreport",
      "command": "bench test --package ./internal/testreport"
    },
    {
      "id": "caller-internal-bounds",
      "command": "bench test --package ./internal/bounds"
    },
    {
      "id": "caller-internal-otelrecord",
      "command": "bench test --package ./internal/otelrecord"
    }
  ]
}
```

## Testing decisions

Use the existing graded-root audit and canary before adding another execution surface.
Preserve `TestInjectedPortRegistryCheckBites` and `TestInjectedPortDerivationSeesEveryPortShape` through the relocation.
Use the shared conformance fixture builder where the test remains in conformance.
Do not paste that builder into a new audit module.
The relocated owner's tests can use small direct source fixtures with one shared local fixture constructor.

Planned `TestSeamDerivationForms` supplies declarations, direct and combined assignments, aliases, shadowed locals, and dynamic hook readers.
Planned `TestSeamAdmissionEvidence` checks row completeness, consumer resolution, canonical aliases, and external limitations.
Planned caller tests check failure-safe restoration in diff, join ordering in guards, and independent instance collaborators in assessment.
Planned `TestSeamAuditOmissions` independently introduces missing registrations and missing evidence into a conformant source tree.
These suites cover multiple rows because their input tables share one observable seam.

The omission fixtures are authored independently of the registry derivation.
During implementation, demonstrate the red from deleting registration, alias resolution, setter discovery, runtime branch discovery, and consumer evidence.
Retain the mutation command and result as implementation verification evidence.
The independently authored expectation exception requires that demonstrated red.
No such demonstration is claimed by this document pass.

Existing process tests remain necessary where they prove subprocess environment, abrupt exit, or real filesystem behavior.
Test cost is considered only after the failure witness is sufficient.
A mock of a local process cannot replace its environment or interruption proof merely because it runs faster.
No extra system test is needed for a pure registry field if the owner suite already proves its refusal.

### Seam diagram

    production declarations + test consumers in the graded root
                            │
                            ▼
    admission rows ──▶ [ existing injected-port owner ] ──▶ sorted diagnostics
                            ▲
                   independent omission fixtures
                            │
              checkInjectedPortRegistry + existing canary

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PS01 | 1 | An injected nonempty interface without a row emits the existing unregistered diagnostic | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Ignoring interface declarations leaves the fixture green |
| PS02 | 1 | An injected named function type without a row emits the unregistered diagnostic | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Ignoring function types leaves the source fixture green |
| PS03 | 1 | An injected all-function struct without a row emits the unregistered diagnostic | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Treating collaborator structs as data misses this port |
| PS04 | 2 | A direct replacement of a production function variable requires admission | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | A setter-only scanner misses the assigned declaration |
| PS05 | 2 | An anonymous function replacement resolves to its package variable | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Named function-type matching alone cannot find it |
| PS06 | 3 | A replaced scalar bound requires admission | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Function-only discovery omits the duration fixture |
| PS07 | 4 | An unregistered setter alias emits a diagnostic for its canonical seam | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | The alias mutation hides the owner from a name-only sweep |
| PS08 | 4 | A ForTesting setter resolves through an unexported forwarding function | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | The alternate suffix defeats an exact ForTest matcher |
| PS09 | 5 | An aliased testing.Testing branch without admission emits a diagnostic | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | An import spelling change cannot hide the runtime branch |
| PS10 | 6 | A constructed phase hook without admission emits a diagnostic | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Literal-only environment matching misses the family |
| PS11 | 6 | An unresolved dynamic hook identity refuses audit success | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Unsupported construction cannot become an empty result |
| PS12 | 8 | A malformed production source file emits a parse diagnostic | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | Dropping parse errors makes a damaged census green |
| PS13 | 9 | A row whose seam disappeared from a present package emits the orphan diagnostic | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | Removing the final seam cannot hide a stale row |
| PS14 | 8 | A present required package with zero eligible seams emits the existing zero-inventory diagnostic | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | Empty derivation cannot satisfy legacy required inventory |
| PS15 | 10 | A fault admission without a production owner refuses | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | A consumer name alone cannot certify ownership |
| PS16 | 11 | A named consumer that the tree does not declare refuses | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | The existing missing-test tripwire survives the move |
| PS17 | 11 | A declared consumer that never reaches the canonical seam refuses | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | An unrelated empty test cannot satisfy the row |
| PS18 | 12 | A fault admission without a named failure purpose refuses | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | Convenience prose cannot replace the required purpose |
| PS19 | 13 | A global admission without isolation evidence refuses | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | Returning a restore closure alone leaves the row incomplete |
| PS20 | 9 | Duplicate source identities refuse admission | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | Last-row overwrite cannot conceal conflicting evidence |
| PS21 | 12 | An unresolved alias refuses admission | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | A nonexistent owner cannot borrow another row's evidence |
| PS22 | 14 | A fatal-path cleanup restores the prior collaborator | planned TestSnapshotOverrideRestoresAfterFailure in internal/diff/seam_policy_test.go | Omitting cleanup leaves a distinguishable replacement installed |
| PS23 | 15 | A held consumer finishes before its collaborator is restored | planned TestScanOverrideJoinsBeforeRestore in internal/guards/seam_policy_test.go | Deleting the join exposes the restored value to a live consumer |
| PS24 | 13 | Two isolated instances retain their own collaborators during concurrent execution | planned TestInstanceCollaboratorsStayLocal in internal/assessment/seam_policy_test.go | A hidden package global cross-contaminates the observations |
| PS25 | 16 | Invalid production input still refuses while a fault collaborator is active | planned TestSeamAuthorityPreserved in internal/releasepreflight/seam_policy_test.go | A test-only success path changes the observed refusal |
| PS26 | 17 | A justified instance-local collaborator passes admission without a package global | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | A registry design requiring a global rejects the conformant case |
| PS27 | 18 | A real dependency test observes the registered producer's failure | review-owned: reviewer traces the named caller tests and omission result | A fake-only consumer cannot establish that the real producer can fail |
| PS28 | 19 | An installed REF disagreement still refuses release identity | `internal/releasepreflight/identity_test.go` (`TestCheckIdentityRefusesADisagreement`) | Removing the ref input makes the disagreement fixture agree |
| PS29 | 20 | The RACE probe observes the isolated child home | `internal/releasepreflight/external_test.go` (`TestExternalPhaseGitStartsNoAutoMaintenance`) | Ignoring the override removes the observing child |
| PS30 | 21 | The release probe exercises each registry phase through the chosen binary | check release-evidence-probe | Omitting an override leaves a real phase outside the controlled composition |
| PS31 | 21, 6 | The dynamic phase family includes the registry's vulnerability phase | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Reusing a historical five-phase list omits the scanner member |
| PS32 | 21 | A new phase registration changes the derived hook family | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | A copied member list stays unchanged and misses the new phase |
| PS33 | 22 | A DATE override drives the vulnerability entry's selected date | planned TestVulnerabilityEntryDateOverride in internal/releasepreflight/seam_policy_test.go | Removing the reader changes the expiry result at the controlled boundary |
| PS34 | 23 | A readiness override reaches the controlled evidence promotion point | planned TestEvidenceReadyControl in internal/releaseevidence/seam_policy_test.go | Removing the hook prevents the observer from seeing the ready point |
| PS35 | 24 | A fallback-home provider under test writes no record | `internal/otelrecord/processor_test.go` (`TestNewProviderRefusesTheFallbackHomeUnderTest`) | Deleting the runtime protection creates a record in the fallback home |
| PS36 | 24 | An explicit home equal to fallback writes no record under test | `internal/otelrecord/processor_test.go` (`TestNewProviderRefusesAnExplicitBenchHomeEqualToTheFallback`) | Checking only an empty constructor argument permits the protected write |
| PS37 | 24 | An explicit isolated home records a span | `internal/otelrecord/processor_test.go` (`TestNewProviderRecordsAnExplicitBenchHomeElsewhere`) | Blanket suppression hides the legitimate isolated record |
| PS38 | 25 | A removal proposal names every resolved caller and its equivalent failure witness | review-owned: compare refreshed symbol census with the removal diff | A partial rg result cannot establish unused status |
| PS39 | 22, 23 | DATE and readiness retain separate compatibility dispositions | review-owned: inspect their individual registry evidence and retained readers | A generic unused-family claim cannot authorize both removals |
| PS40 | 16 | Link interruption still leaves recoverable production transaction state | `internal/systemtest/compatibility_test.go` (`TestCompatibilityInterruptedRepair`) | Removing the real abrupt exit removes the recovery trigger |
| PS41 | 28 | The existing unregistered-port canary emits its exact diagnostic | check injected-port-registry through tests/canary/injected-ports/unregistered-port/ | Replacing the check wrapper with a green stub loses the fixture refusal |
| PS42 | 30 | An absent audited package remains benign in a partial root | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | Reading ambient production source invents a fixture obligation |
| PS43 | 30 | An unreadable present source package refuses | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | Treating every stat error as absence loses the failure |
| PS44 | 7, 26 | Ordinary functions and source examples produce no seam finding | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Name-only matching or testing-import matching creates false findings |
| PS45 | 27 | Deleting a registration from conformant production source turns the audit red | planned TestSeamAuditOmissions in internal/conformance/injectedports/omissions_test.go | The production scan must discover a seam not supplied by registry iteration |
| PS46 | 27 | Removing alias resolution makes the independent alias fixture red | planned TestSeamAuditOmissions in internal/conformance/injectedports/omissions_test.go | A registry-only inventory cannot certify alias coverage |
| PS47 | 27 | Removing setter discovery makes the independent setter fixture red | planned TestSeamAuditOmissions in internal/conformance/injectedports/omissions_test.go | A named-type-only fallback cannot satisfy the setter assertion |
| PS48 | 27 | Removing runtime branch discovery makes the independent runtime fixture red | planned TestSeamAuditOmissions in internal/conformance/injectedports/omissions_test.go | The negative mutation must change the branch finding |
| PS49 | 29, 31 | The old audit file contains no relocated derivation implementation | review-owned: whole-tree definition and caller census after relocation | An unreachable second scanner violates the one-owner outcome |
| PS50 | 1 | A pointer parameter retains its named-port finding | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Removing pointer reduction loses this injection site |
| PS51 | 1 | A type assertion retains its named-port finding | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Parameter-only discovery misses the widened capability |
| PS52 | 1 | A type switch retains its named-port finding | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Dropping the case arm removes its injection evidence |
| PS53 | 7 | An empty interface produces no named-port finding | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Treating every interface as a port creates a false finding |
| PS54 | 7 | A struct with a data field produces no named-port finding | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Ignoring field shape registers ordinary data |
| PS55 | 7 | An unused port-shaped declaration produces no named-port finding | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Shape alone cannot establish an injection |
| PS56 | 5 | A forwarding test-runtime predicate reaches its canonical branch | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Matching only direct calls hides the forwarding branch |
| PS57 | 2 | A combined assignment resolves the eligible package variable | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Reading only the first assignment target loses the collaborator |
| PS58 | 6 | The verdict-window control preserves parent cancellation | planned TestUnboundedParentCancellation in internal/bounds/seam_policy_test.go | An override cannot make cancellation inert |
| PS59 | 6 | Fixed cancellation grace remains finite under the verdict override | `internal/bounds/bounds_test.go` (`TestFixedWindowPreservesCancelGrace`) | Treating all windows as overridable weakens termination policy |
| PS60 | 13, 15 | A global admission names serial ownership through consumer completion | review-owned: trace each retained global's tests and cleanup ordering | A restore closure cannot hide parallel consumers |
| PS61 | 12 | A blank external exemption retains the existing exemption refusal | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | Whitespace cannot certify the explicit external limitation |
| PS62 | 7 | Ordinary mutable production data produces no test-seam finding | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | Registering every global would reject the negative fixture |
| PS63 | 11, 22, 23 | An environment exception with no consumer evidence refuses | planned TestSeamAdmissionEvidence in internal/conformance/injectedports/admission_test.go | Claiming external absence cannot substitute for admission evidence |
| PS64 | 7 | Assigning a shadowed local leaves the production variable unclassified | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | A text-name match confuses fixture state with a production override |
| PS65 | 6 | A test-support answer-path environment reader produces no production hook finding | planned TestSeamDerivationForms in internal/conformance/injectedports/derive_test.go | An environment prefix alone cannot establish production use |
| PS66 | 28, 30 | The audit input classification includes graded-root registry data | review-owned: compare InputCatchAll metadata and profile advertisement with the actual data reader | A Go-only claim omits a behavior-changing input |

### Edge inventory

| Edge class | Disposition and evidence |
|---|---|
| Absent versus present empty source | PS14 and PS42 distinguish the legacy inventory contract |
| Malformed or unreadable Go source | PS12 and PS43 refuse success |
| Spaces, Unicode, and glob characters in root paths | Include paths in owner source fixtures without shell interpretation |
| Symlink or nonregular source input | Refuse unsafe source reads through the existing bounds classification where applicable |
| Build variants and unexported names | Read platform variants and resolve declaration identity |
| Alias cycles and unresolved imports | Refuse unresolved seam identity rather than guessing a canonical owner |
| Local shadowing and multiple assignment | Resolve actual package targets and keep ordinary fixture state negative |
| Dynamic phase names | Derive bounded membership from the existing phase registry, including VULNERABILITY |
| Fatal assertion or interrupted child | PS22 and PS40 preserve cleanup or recovery evidence |
| Concurrent live consumers | PS23 and PS24 distinguish global restoration from instance isolation |
| Existing npm external dependency | Retain its individual reviewed limitation without a blanket package exemption |
| Comments, strings, and source examples | PS44 requires executable declarations and real production reachability |

**Won't handle:** arbitrary reflection or unsafe mutation of package memory — normal callers and resolved Go assignments remain covered.

**Won't handle:** external callers outside the source tree — REF, DATE, and readiness compatibility remain preserved rather than inferred absent.

**Won't handle:** proving consumer purpose from source syntax alone — the named failure witness and review remain required.

**Won't handle:** deleting telemetry's fallback-home guard during this build — its three existing processor tests remain the safety witnesses.

## Ownership fences

Prospective implementation writes:

- `internal/conformance/injected_ports_test.go`
- `internal/conformance/injected_ports_registry_test.go`
- `internal/conformance/injectedports/`
- `internal/conformance/checks_test.go`
- `internal/conformance/native_workflow_test.go`
- `internal/conformance/registry/checks.go`
- `internal/conformance/registry/registry.go`
- `internal/conformance/registry/registry_test.go`
- `tests/canary/injected-ports/`
- `cmd/bench/`
- `internal/adopt/`
- `internal/assessment/`
- `internal/bounds/`
- `internal/capability/`
- `internal/chargeevidence/`
- `internal/consumers/`
- `internal/coverage/`
- `internal/diff/`
- `internal/freshness/`
- `internal/gate/`
- `internal/git/`
- `internal/gitguard/`
- `internal/guards/`
- `internal/handoffdoc/`
- `internal/harnesstranscript/`
- `internal/intent/`
- `internal/otelrecord/`
- `internal/outline/`
- `internal/preflight/review.go`
- `internal/preflight/review_charge_test.go`
- `internal/preflight/evidencecmd/`
- `internal/publication/`
- `internal/refresh/`
- `internal/releaseevidence/`
- `internal/releasepreflight/`
- `internal/roadmap/`
- `internal/sessioninspect/`
- `internal/shift/`
- `internal/skillsindex/`
- `internal/status/status_render_test.go`
- `internal/systemtest/compatibility_test.go`
- `internal/testreport/`
- `projects/benchkit.md`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/guidance-prose-budgets/over-budget-skill/`
- `tests/canary/line-routing/line-binding-prose-drift/`
- `tests/canary/skill-description-budgets/budget-table-missing/`
- `tests/canary/skill-description-budgets/description-folded/`
- `tests/canary/skill-description-budgets/description-missing/`
- `tests/canary/skill-description-budgets/over-budget-command/`
- `tests/canary/skill-description-budgets/over-budget-description/`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading/`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner/`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing/`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership/`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route/`
- `specs/production-test-seam-policy/`
- `capture/`
- `reviews/production-test-seam-policy.md`

The caller prefixes permit only admitted seam construction, isolation repairs, and sufficient tests for this policy.
They do not authorize unrelated package changes or blanket replacement of all mutable variables.
Ticket slicing must narrow each caller family after independent spec acceptance.
The final fence union must match those ticket writes.

## Headroom and transitive closure

The current audit file has 491 lines against the default 400-line limit.
Conformance has no file-count headroom, as its structure budget explains.
Relocate the same owner into a small subpackage and reduce the original file to its wrapper and shared integration tests.
Keep new files below the default limit and the new directory below its file-count limit.
Do not add a budget equal to a current source count.

The check entry, binding table, check registry, family binding, and existing canary are explicit fence members.
The source sweep found no anchor or command-registry entry that names the private audit helpers.
Their relocation alone requires no command or anchor edit.
The graded-root phase registry adds a data input to the broader check.
Use the existing `InputCatchAll` classification until a narrower sound input owner exists.
Update the profile's existing input row from `go-source` to `catch-all` in TS3, before its registry-reading consumer lands.

That profile edit carries four anchor owners and the existing AXI command advertisement closure.
The fence includes those owners and the twelve fixtures that bind or supply the profile.
The fixture paths come from their BASE, MUTATE, and files trees.
Keep their existing refusal expectations when refreshing source snapshots.

The slicing headroom refresh reads `.bench/structure.budgets` and `.bench/structure-accept` at the accepted source.
`native_workflow_test.go` is exactly 400 lines; adding logic there must first remove or relocate the same responsibility within TS3.
`checks_test.go` has 641 lines against its 709-line grant, and `registry/registry.go` has 355 against 403.
The new owner directory starts absent and has the default 12-file and 400-line caps.

Several caller directories are already crowded, including assessment at 29 files and diff at 15.
Use existing caller test files where they fit; a planned seam filename cannot excuse new crowding.
Existing long caller files cannot grow past their caps without same-ticket reduction within the approved caller fence.
No budget or grant write is authorized by this graph.
Each introducing ticket refreshes actual touched-file headroom against its predecessor before writing.

A later guidance or command change requires a refreshed closure proposal before writing outside these fences.
No new command, canary family, or profile table is proposed here.
The existing input advertisement changes for this check.

## Out of scope

- Shared syntax visitor infrastructure under FT365: separate capability, estimated 4 edits and 2 gate runs.
- Release forwarding retirement under C11: separate capability, estimated 6 edits and 3 gate runs.
- Replacing telemetry home construction across all command entries: separate capability, estimated 8 edits and 4 gate runs.
- Replacing every package global with instance state: separate capability, estimated 24 edits and 6 gate runs after its own census.

These estimates price separate capabilities, not unfinished acceptance work.
Any isolation repair required for an admitted seam remains in this spec.
No benchmark, commitment change, roadmap reorder, or automatic landing is authorized.

## Further notes

Document checks establish only spec grammar, prose, link integrity, and planned coverage.
They do not establish production behavior, isolation safety, benchmark gains, or gate mutation results.
The compiled map and topic folder are consumed together and retain their original bytes.
The architecture index changes only the C09 link.

The author freezes this full spec before independent review.
Preserve a snapshot before any review repair.
After spec acceptance, slice the green proposals into tickets and review their graph independently.
The accepted spec must remain the basis of that graph.
