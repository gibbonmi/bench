package conformance

import "github.com/gibbonmi/bench/internal/conformance/injectedports"

// auditedPortPackages are the packages that must declare injected ports whenever the graded
// tree holds them. A package inside it that derives nothing fails closed. The list does not
// bound discovery: the audit derives the ports of every production package.
var auditedPortPackages = []string{
	"internal/assessment",
	"internal/chargeevidence",
	"internal/gitguard",
	"internal/preflight",
	"internal/publication",
	"internal/subprocess/protection",
}

// injectedPortRegistry is the advertisement half: one row per named port, then one row per
// variable seam, setter, test-runtime branch, and environment hook that the census finds
// anywhere in the production roots. Rows are grouped by kind and package and ordered as the
// derivation reports them. A diff against a re-derived inventory therefore reads straight
// down. A named-port row names the test that drives the port's real producer, or it is a
// fault-port row.
var injectedPortRegistry = append(append([]injectedports.Row{
	{Package: "internal/adopt", Port: "compatibilityCollector", TestFile: "internal/systemtest/compatibility_test.go", TestName: "TestCompatibilityPaths"},
	{Package: "internal/assessment", Port: "FileOps", TestFile: "internal/assessment/replacement_test.go", TestName: "TestAssessmentReplacement"},
	{Package: "internal/benchguard", Port: "Resolver", TestFile: "internal/benchguard/benchguard_test.go", TestName: "TestClassifyResolvesBareAliasFromProcessPath"},
	{Package: "internal/bounds", Port: "pathResolver", TestFile: "internal/bounds/classify_nofollow_test.go", TestName: "TestClassifyDirNoFollowRefusesLinkAndPreservesDefault"},
	{Package: "internal/capability", Port: "TB", TestFile: "internal/capability/capability_test.go", TestName: "TestEnvironmentWritesLineBeforeSkip"},
	{
		Package: "internal/chargeevidence", Port: "ExportWrite",
		TestFile: "internal/preflight/evidencecmd/evidence_export_test.go", TestName: "TestEvidenceExportWritesSources",
	},
	{
		Package: "internal/chargeevidence", Port: "SourceAt",
		TestFile: "internal/preflight/evidencecmd/evidence_export_test.go", TestName: "TestEvidenceExportVerifiesBeforeWrite",
	},
	{
		Package: "internal/chargeevidence", Port: "StoreOptions",
		TestFile: "internal/systemtest/charge_evidence_test.go", TestName: "TestEvidenceInterruptedPublication",
	},
	{Package: "internal/commitment", Port: "Resolver", TestFile: "internal/commitment/verification_test.go", TestName: "TestCommitmentMilestoneCompletion"},
	{Package: "internal/durablefile", Port: "File", TestFile: "internal/durablefile/replace_test.go", TestName: "TestReplaceBytes"},
	{Package: "internal/gate", Port: "FS", TestFile: "internal/gate/run_outcomes_test.go", TestName: "TestGateRunWithoutGateExitsThree"},
	{Package: "internal/gate", Port: "executionEvaluation", TestFile: gateJourney, TestName: gateJourneyTest},
	{Package: "internal/gate", Port: "postAcquireContextArm", TestFile: gateJourney, TestName: gateJourneyTest},
	{Package: "internal/gate", Port: "runBinaryOwner", TestFile: gateJourney, TestName: gateJourneyTest},
	{Package: "internal/gate", Port: "treeSource", TestFile: gateJourney, TestName: gateJourneyTest},
	{
		Package: "internal/gitguard", Port: "Checker",
		TestFile: "internal/gitguard/checker_junction_test.go", TestName: "TestClassifyRealCheckerResolvedComposition",
	},
	{Package: "internal/harnesstranscript", Port: "lineReader", TestFile: "internal/harnesstranscript/transcript_test.go", TestName: "TestCodexFixtureCountsItsPinnedShapes"},
	{Package: "internal/intent", Port: "Compensation", TestFile: "internal/intent/transaction_test.go", TestName: "TestReauthorizeCompensatesOnTheTerminalWriteFailure"},
	{Package: "internal/intent", Port: "Decision", TestFile: "internal/intent/transaction_test.go", TestName: "TestLedgerTransactionPreservesEveryMutatorOutcome"},
	{Package: "internal/intent", Port: "countedDecision", TestFile: "internal/intent/purge_test.go", TestName: "TestPurgeAssignmentsDropsRecordsReadRefuses"},
	{Package: "internal/landing", Port: "Admission", TestFile: "internal/worktree/commitment_landing_test.go", TestName: "TestCommitmentDeliveryClosure"},
	{Package: "internal/landing", Port: "StackAdmission", TestFile: "internal/worktree/land_stack_test.go", TestName: "TestLandStackJourney"},
	{Package: "internal/maps", Port: "SourceResolver", TestFile: "internal/preflight/map_sources_test.go", TestName: "TestMapSourcesTipExistence"},
	{
		Package: "internal/preflight", Port: "reviewEvidenceObserver",
		TestFile: "internal/preflight/review_charge_test.go", TestName: "TestEvidenceReviewCollectors",
	},
	{Package: "internal/prose", Port: "targetLookup", TestFile: "internal/prose/subject_test.go", TestName: "TestGradeNamedHonorsExclusions"},
	{
		Package: "internal/publication", Port: "Registry",
		Exempt: "a real-producer test needs a live npm registry, and NPMCLIRegistry is runbook-only: it never performs gate egress or touches a credential (the Registry port comment in internal/publication/registry.go states each adapter's gate coverage)",
	},
	{Package: "internal/roadmap", Port: "OutlookSource", TestFile: "cmd/bench/commitment_test.go", TestName: "TestCommitmentNextProjection"},
	faultPort("internal/shift", "fault", "hitShift", "an injected step fault drives the staging, teardown, upsert, and startup failures", joinIsolation,
		"internal/shift/fault_test.go::TestLoopStagingFaultPreservesAndSplitsEvidence", "internal/shift/lifetime_test.go::TestActualStartRace"),
	{Package: "internal/spec", Port: "TreeReader", TestFile: "internal/landing/close_test.go", TestName: "TestTicketsOnlyInCommit"},
	{Package: "internal/subprocess", Port: "CancelRequest", TestFile: "internal/releasepreflight/lifetime_test.go", TestName: "TestScannerCancel"},
	{Package: "internal/subprocess", Port: "Window", TestFile: "internal/bounds/lifetime_test.go", TestName: "TestBoundsFixedShutdown"},
	{
		Package: "internal/subprocess/protection", Port: "Probe",
		TestFile: "internal/subprocess/protection/observe_test.go", TestName: "TestObserveGroupPresent",
	},
	faultPort("internal/worktree", "Fault", "applyAutomaticWithTerminal", "a fault at a chosen lifecycle step drives the cleanup transaction's locked, relocked, and replayed states",
		"each consumer passes the fault as an argument of one call, so it lives only for that call and no package state changes",
		"internal/worktree/ownership_test.go::TestLifecycleFaultBoundariesRemainLockedOrAbsent", "internal/worktree/worktree_test.go::TestReleaseReconcilesInFlightAutomaticCleanup"),
	{Package: "internal/worktree", Port: "cleanupPlanner", TestFile: "internal/worktree/clean_landed_apply_test.go", TestName: "TestCleanLandedApplyRemovesAndSettles"},
	{Package: "internal/worktree", Port: "cleanupTerminal", TestFile: "internal/worktree/clean_landed_apply_test.go", TestName: "TestCleanLandedApplyRemovesAndSettles"},
	{Package: "internal/writeguard", Port: "Checker", TestFile: "cmd/bench/guard_projection_test.go", TestName: "TestGuardWriteEnvelopePolicy"},
}, variableSeams...), append(setterSeams, controlSeams...)...)

// The gate's internal ports share one real-producer test: the wrapper runs the gate
// through its signal arm, its production run-binary owner, and its tree evaluation.
const (
	gateJourney     = "internal/systemtest/otel_gate_test.go"
	gateJourneyTest = "TestOtelGateRecordJourney"
)

// faultPort admits one fault port. The audit owner states what its owner and consumers show.
func faultPort(pkg, port, owner, purpose, isolation string, consumers ...string) injectedports.Row {
	return injectedports.Row{Package: pkg, Port: port, Owner: owner, Purpose: purpose, Isolation: isolation, Consumers: consumers}
}

// The isolation evidence that most variable seams share. A seam whose tests differ states
// its own evidence.
const (
	serialIsolation = "serial: no installing test or subtest calls t.Parallel; the installer registers t.Cleanup or defer before its first fallible step; the production path runs on the test goroutine and returns before that cleanup"
	childIsolation  = "the replacement lives only in a helper child test process that exits after the observed step; the parent test never installs it"
	joinIsolation   = "serial: no installing test calls t.Parallel; cleanup is registered before the first fallible step; each goroutine that reads the replacement is joined on every path before that cleanup restores it"
)

// seamVariable admits one replaced production variable. Each consumer is a top-level test,
// written "path::TestName", that installs the replacement itself or through its own helper.
func seamVariable(pkg, name, owner, purpose, isolation string, consumers ...string) injectedports.Row {
	return injectedports.Row{
		Kind: injectedports.KindVariable, Package: pkg, Port: name, Owner: owner,
		Purpose: purpose, Isolation: isolation, Consumers: consumers,
	}
}

// envIsolation is the isolation evidence that every environment hook shares.
const envIsolation = "each consumer sets the hook through t.Setenv, which refuses a parallel test and restores the prior value at cleanup, or through one child process's environment entry, which ends with that child"

// seamHook admits one environment hook, read by owner in pkg and set by its consumers.
func seamHook(name, pkg, owner, purpose string, consumers ...string) injectedports.Row {
	return injectedports.Row{
		Kind: injectedports.KindEnvironment, Package: pkg, Port: name, Owner: owner,
		Purpose: purpose, Isolation: envIsolation, Consumers: consumers,
	}
}

// seamSetter admits one setter. It names only the variable seams it finally replaces.
func seamSetter(pkg, name string, canonical ...string) injectedports.Row {
	return injectedports.Row{Kind: injectedports.KindSetter, Package: pkg, Port: name, Canonical: canonical}
}

var variableSeams = []injectedports.Row{
	seamVariable("cmd/bench", "commandRegistry", "commandByName", "a fixture command table drives dispatch, help, the response bound, and retention through the real dispatcher", serialIsolation,
		"cmd/bench/help_inventory_test.go::TestHelpRendersPublicCommandRegistryRows", "cmd/bench/retain_test.go::TestRetainKeepsTheCompleteResponse"),
	seamVariable("cmd/bench", "commitChain", "commitChainCommand", "recorded chain steps show the step order and the stop at a refused step", serialIsolation,
		"cmd/bench/commit_chain_test.go::TestCommitChainStopsAtRedBuild"),
	seamVariable("cmd/bench", "gatePhasesCommand", "commandRegistry", "a recording handler shows that the gate-phases row dispatches to its command", serialIsolation,
		"cmd/bench/main_test.go::TestRunGatePhasesDispatchesToCommand"),
	seamVariable("internal/adopt", "syncDirectory", "promoteWithLease", "a destination sync failure, which a real filesystem cannot produce on demand, drives the promotion rollback", serialIsolation,
		"internal/adopt/adopt_test.go::TestPromoteAllRollsBackOnDestinationSyncFailure"),
	seamVariable("internal/capability", "stdout", "Capability", "a captured writer shows that the capability line is written before the skip", serialIsolation,
		"internal/capability/capability_test.go::TestCapabilityWritesLineBeforeSkip"),
	seamVariable("internal/commit", "worktreeBuild", "RunChecked", "a recording or failing rebuild shows the rebuild order and that a failed rebuild keeps exit zero", serialIsolation,
		"internal/commit/worktree_rebuild_test.go::TestCommitKeepsExitZeroOnAFailedRebuild"),
	seamVariable("internal/conformance/registry", "Checks", "Names", "a synthetic check table shows that the conformance meta checks refuse a malformed registry", serialIsolation,
		"internal/conformance/tier_test.go::TestConformanceMetaBites"),
	seamVariable("internal/consumers", "load", "command", "typed in-process fixtures replace the package loader, and a failing loader shows that an identical pair needs no load", joinIsolation,
		"internal/consumers/blast_edges_test.go::TestIdenticalPairAnswersWithoutLoadingPackages", "internal/consumers/refuse_test.go::TestSweepSkipsSymlinkAndFifoCandidates"),
	seamVariable("internal/coverage", "packageLoadTimeout", "runGo", "a short package-load bound drives the citation timeout refusal", serialIsolation,
		"internal/coverage/citation_execution_test.go::TestCitationPackageScopeTimeout"),
	seamVariable("internal/diff", "snapshotAfterRead", "MovementChecked", "a deterministic tree movement between the two identity captures drives the retry and the drift refusal", serialIsolation,
		"internal/diff/command_test.go::TestCommandRetriesThenRefusesSnapshotDrift", "internal/preflight/charge_test.go::TestChargeSnapshotMovement",
		"internal/diff/command_test.go::TestSnapshotOverrideRestoresAfterFailure"),
	seamVariable("internal/freshness", "publicationFlock", "lockPublicationDirectory", "an observer around the real kernel lock reports contention to the parent publication test", childIsolation,
		"internal/freshness/publication_lock_test.go::TestPublicationProcess"),
	seamVariable("internal/freshness", "replacePublicationFile", "writeSeal", "a forced seal promotion failure shows that the publish restores the prior executable, seal, and manifest", serialIsolation+"; the interruption and contention children install it only in their own process",
		"internal/freshness/freshness_publish_test.go::TestPublishRestoresPriorPairWhenSealPromotionFails"),
	seamVariable("internal/gate", "execBroker", "runUnderBroker", "a recording exec replaces the process-image replacement, so a checkpoint test observes the broker argv and its refusals", serialIsolation,
		"internal/gate/complete_checkpoint_broker_test.go::TestCompleteCheckpointRefusesAnUnusableBroker"),
	seamVariable("internal/gate", "gateLogPathIgnored", "beginGateRunLog", "a stubbed ignore probe lets a fixture repository write and prune the run log", serialIsolation,
		"internal/gate/run_log_prune_test.go::TestGateRunPrunesThroughTheClosureBeginHandsBack"),
	seamVariable("internal/gate", "gateTimeout", "gateBudget.runContext", "a short gate bound drives the timeout outcome and the invalidation of old evidence", serialIsolation,
		"internal/gate/timeout_recovery_count_test.go::TestGateRunTimeoutInvalidatesOldEvidence"),
	seamVariable("internal/gate", "laneRunBinary", "selectLaneRunBinary", "a recording run-binary factory shows the checks a lane selects, and a build hook under the real verifier shows the freshness refusal", serialIsolation,
		"internal/gate/lane_run_test.go::TestBenchkitLaneRunsOnlyTheSelectedChecks", "internal/gate/lane_run_test.go::TestLaneRunBinaryHookKeepsFreshnessRefusal"),
	seamVariable("internal/gate", "prospectiveRunBinary", "prospectiveRunBinaryOwnerAt", "a refusing or recording build factory shows that the prospective gate reports a failed build and leaves no bundle", serialIsolation,
		"internal/gate/prospective_owner_test.go::TestProspectiveBuildRefusalLeavesNoBundle"),
	seamVariable("internal/gate", "reenterImage", "BackgroundCommand", "a recording re-entry shows that a background grade runs as the target executable", serialIsolation,
		"internal/gate/background_grade_test.go::TestBackgroundGradeRunsAsTheTargetExecutable"),
	seamVariable("internal/gate", "runningExecutableDigest", "authenticateBroker", "a controlled running-executable digest drives the oracle identity and its refusals", serialIsolation,
		"internal/gate/oracle_identity_test.go::TestOracleIdentityOpensWithoutTheRunningExecutable"),
	seamVariable("internal/gate", "takeExecutionLock", "acquireExecutionLock", "a non-contention lock error, which a real lock cannot produce on demand, drives the refusal with its cause", serialIsolation,
		"internal/gate/execution_probe_test.go::TestGateRunRefusesNonContentionLockErrorWithCause"),
	seamVariable("internal/git", "refCheckTimeout", "refCheck", "a short ref-probe bound drives the checker's fail-safe timeout verdict", serialIsolation,
		"internal/gitguard/checker_junction_test.go::TestClassifyRealCheckerTimeoutComposition"),
	seamVariable("internal/git", "worktreeListTimeout", "boundedGit", "a short discovery bound drives the typed timeout failure of each reader and its renderers", joinIsolation,
		"internal/git/worktree_admin_enum_test.go::TestWorktreesBoundsEachChildAndPreservesStdout", "internal/git/admin_readers_test.go::TestReadersTimeOutUnderTheWorktreeListBound",
		"internal/status/status_render_test.go::TestAppendWorktreeRendersBoundExpiryAsTypedFailure", "internal/worktree/list_actions_test.go::TestListCommandRendersBoundExpiryAsTypedFailure"),
	seamVariable("internal/guards", "enumerateGuards", "Scan", "a held or failing enumeration drives the incomplete scan counts and its cancellation", joinIsolation,
		"internal/guards/guards_cleanup_test.go::TestScanWaitsForCancelledWorkerCleanup", "internal/guards/guards_cleanup_test.go::TestScanOverrideJoinsBeforeRestore"),
	seamVariable("internal/guards", "guardScanTimeout", "Command", "a short scan bound drives the checked-in timeout responses", serialIsolation,
		"internal/guards/guards_test.go::TestCommandPreservesCheckedInEnumerationTimeoutPrimaryResponse"),
	seamVariable("internal/guards", "inspectGuard", "Scan", "a held inspection shows that a cancelled scan waits for its worker cleanup", joinIsolation,
		"internal/guards/guards_cleanup_test.go::TestScanWaitsForCancelledWorkerCleanup", "internal/guards/guards_cleanup_test.go::TestScanOverrideJoinsBeforeRestore"),
	seamVariable("internal/handoffdoc", "lockDeadline", "acquire", "a short lock window drives the refusal of a lock that another writer holds", serialIsolation,
		"internal/handoffdoc/store_test.go::TestUpdateRefusesALockAnotherWriterHolds"),
	seamVariable("internal/harnesstranscript", "maxRecordLine", "readRecordLine", "a small line bound drives the oversized-line skip", serialIsolation,
		"internal/harnesstranscript/transcript_test.go::TestOversizedLineIsOneSkippedEvent"),
	seamVariable("internal/intent", "lockWindow", "acquire", "a short lock window drives the ledger lock refusal while another writer holds it", joinIsolation,
		"internal/intent/remove_fields_test.go::TestRemoveFieldsTakesTheLedgerLock", "internal/intent/remove_fields_test.go::TestRemoveFieldsHoldsTheLedgerLockThroughItsWrite"),
	seamVariable("internal/maps", "canonicalDecisionMapSchema", "DecisionMapTemplate", "a changed status syntax shows that the template and the parser both derive from the one schema", serialIsolation,
		"internal/maps/maps_parse_test.go::TestDecisionMapSchemaSyntaxDrivesParserAndTemplate"),
	seamVariable("internal/models", "doHTTP", "fetchDataIDsContext", "a stubbed transport drives provider failures into unavailable rows at exit zero", joinIsolation,
		"internal/models/models_test.go::TestCommandReportsSourceFailuresAsUnavailableAtExitZero"),
	seamVariable("internal/models", "providerTimeout", "apiInventory", "a short provider bound drives the unavailable timeout row", serialIsolation,
		"internal/models/models_test.go::TestAPIInventoryTimesOutAsUnavailable"),
	seamVariable("internal/models", "runCommand", "codexInventory", "a stubbed command drives the bundled fallback and the missing-key report", joinIsolation,
		"internal/models/models_test.go::TestCommandFallsBackToCodexBundledAndReportsMissingKeys"),
	seamVariable("internal/outline", "openOutlineFile", "ReadScannable", "a permission failure on one file drives the unreadable skip", serialIsolation,
		"internal/outline/outline_test.go::TestCommandNamesUnreadableSkip"),
	seamVariable("internal/preflight", "reviewEvidenceObserved", "reviewChargePack", "a counting observer shows that each review collector runs once per attempt", serialIsolation,
		"internal/preflight/review_charge_test.go::TestEvidenceReviewCollectors"),
	seamVariable("internal/preflight/evidencecmd", "responseLimit", "Bound", "a small response bound drives the public over-bound refusal", serialIsolation,
		"internal/preflight/evidencecmd/evidence_budget_test.go::TestEvidenceResponseBound", "internal/preflight/staleness/drift_test.go::TestAuditBound"),
	seamVariable("internal/refresh", "refreshTimeout", "Refresh", "a short refresh bound drives the nonfatal timeout detail", serialIsolation,
		"internal/refresh/refresh_test.go::TestRefreshFailureAndTimeoutAreNonfatalAndDetailed"),
	seamVariable("internal/releaseevidence", "exchangeEvidenceDirs", "promoteStage", "an exchange failure, which a real filesystem cannot produce on demand, refuses the promotion and keeps the prior evidence", serialIsolation,
		"internal/releaseevidence/protection_test.go::TestExchangeFailureKeepsThePriorEvidence"),
	seamVariable("internal/releaseevidence", "maxArchiveMemberSize", "readSecureArchive", "a lowered member bound reaches the member-size refusal, which the lower expanded-size bound otherwise shadows", serialIsolation,
		"internal/releaseevidence/package_artifact_test.go::TestValidateTarballRefusesAMemberOverTheMemberLimit"),
	seamVariable("internal/releaseevidence", "requirements", "Requirements", "a reduced requirement registry drives the package identity and component compatibility readers", serialIsolation,
		"internal/releaseevidence/package_artifact_test.go::TestPackageIdentityJSONCompatibility"),
	seamVariable("internal/reviewrecord", "replaceFile", "replace", "a failing replacement drives the review record's replacement refusal", serialIsolation,
		"internal/reviewrecord/replacement_test.go::TestReviewReplacement"),
	seamVariable("internal/roadmap", "openRetroRoot", "RetroCommand", "a replaced destination root drives the retro containment refusal", serialIsolation,
		"internal/roadmap/retro_test.go::TestRetroContainsDestinationComponentReplacement"),
	seamVariable("internal/sessioninspect", "discoveryTimeout", "environmentPhase", "a short discovery bound stops the environment phase at its bound", serialIsolation,
		"internal/sessioninspect/sessioninspect_test.go::TestEnvironmentPhaseTE15StopsAtDiscoveryBound"),
	seamVariable("internal/sessioninspect", "phases", "Inspect", "a phase held until the deadline drives the deadline warning and the zero exit of the stopped inspection", serialIsolation,
		"internal/sessioninspect/sessioninspect_test.go::TestInspectDeadlineWarnsAndReturnsZero"),
	seamVariable("internal/sessioninspect", "providerTimeout", "Command", "the provider bound is observed as the command's ten-second deadline", serialIsolation,
		"internal/sessioninspect/sessioninspect_test.go::TestCommandInstallsTenSecondDeadline"),
	seamVariable("internal/sessioninspect", "runInspect", "Command", "a recording inspection shows that the command installs its ten-second deadline", serialIsolation,
		"internal/sessioninspect/sessioninspect_test.go::TestCommandInstallsTenSecondDeadline"),
	seamVariable("internal/shift", "shiftFault", "stageTouched", "an injected step fault drives the staging, teardown, upsert, and startup failures", joinIsolation,
		"internal/shift/fault_test.go::TestLoopStagingFaultPreservesAndSplitsEvidence", "internal/shift/lifetime_test.go::TestActualStartRace"),
	seamVariable("internal/shift", "timeNow", "loop", "a fixed clock makes the branch-name collision and its retries deterministic", serialIsolation,
		"internal/shift/shift_test.go::TestLoopRetriesBranchCreationOnCollision"),
	seamVariable("internal/skillsindex", "preReplacementBarrier", "replaceReference", "a held barrier in a child process lets the parent observe the write before the replacement", childIsolation,
		"internal/skillsindex/command_test.go::TestWriteBarrierHelperProcess"),
	seamVariable("internal/skillsindex", "renameFile", "replaceReference", "a rename failure drives the refusal that leaves no residue and keeps the reference bytes", serialIsolation,
		"internal/skillsindex/skillsindex_reference_test.go::TestRenameFailureLeavesNoResidueAndKeepsReferenceBytes"),
	seamVariable("internal/subprocess/protection", "storeBase", "storeRoot", "a private store base isolates the protection store and drives its refusals", serialIsolation,
		"internal/subprocess/protection/directory_test.go::TestStoreRootRefusal"),
	seamVariable("internal/testreport", "listCurrentPackages", "resolveChangedPackages", "a controlled package graph drives each changed-selection refusal", serialIsolation,
		"internal/testreport/selection_test.go::TestChangedPackageSelectionRefusalMatrix"),
	seamVariable("internal/testreport", "runBinaryFactory", "selectRunBinary", "a test selection factory shows the registered dev scope that a named check runs, and a build hook under the real verifier shows the freshness refusal", serialIsolation,
		"internal/testreport/check_test.go::TestNamedCheckRunsOnlyRegisteredDevScope", "internal/testreport/runbinary_test.go::TestSelectRunBinaryHookKeepsFreshnessRefusal"),
	seamVariable("internal/testreport", "runningExecutable", "unknownCheck", "a controlled running executable drives the unknown-check report", serialIsolation,
		"internal/testreport/unknown_check_test.go::TestUnknownCheckNamesUnknownExecutable"),
}

var setterSeams = []injectedports.Row{
	seamSetter("internal/diff", "SetSnapshotAfterReadForTest", "internal/diff.snapshotAfterRead"),
	seamSetter("internal/git", "SetRefCheckTimeoutForTest", "internal/git.refCheckTimeout"),
	seamSetter("internal/git", "SetWorktreeListTimeoutForTest", "internal/git.worktreeListTimeout"),
	seamSetter("internal/handoffdoc", "setLockWindowForTest", "internal/handoffdoc.lockDeadline"),
	seamSetter("internal/intent", "setLockWindowForTest", "internal/intent.lockWindow"),
	seamSetter("internal/preflight", "setReviewEvidenceObserverForTest", "internal/preflight.reviewEvidenceObserved"),
	seamSetter("internal/preflight/evidencecmd", "SetResponseLimitForTest", "internal/preflight/evidencecmd.responseLimit"),
	seamSetter("internal/releaseevidence", "SetArchiveMemberLimitForTesting", "internal/releaseevidence.maxArchiveMemberSize"),
	seamSetter("internal/releaseevidence", "SetExchangeForTesting", "internal/releaseevidence.exchangeEvidenceDirs"),
	seamSetter("internal/releaseevidence", "SetRequirementsForTesting", "internal/releaseevidence.requirements"),
	seamSetter("internal/releasepreflight", "setArchiveMemberLimitForTesting", "internal/releaseevidence.maxArchiveMemberSize"),
	seamSetter("internal/releasepreflight", "setExchangeForTesting", "internal/releaseevidence.exchangeEvidenceDirs"),
	seamSetter("internal/releasepreflight", "setRequirementsForTesting", "internal/releaseevidence.requirements"),
	seamSetter("internal/sessioninspect", "setWindowsForTest", "internal/sessioninspect.discoveryTimeout", "internal/sessioninspect.providerTimeout"),
}

// The release phase hooks share their reader and their native consumer.
const (
	phaseReader = "runner.runExternal"
	rootProbe   = "internal/conformance/gate_entry_test.go::TestRootConformance"
	gateStream  = "internal/gate/run_stream_test.go::TestGateRunInheritsTheRecordWhenTheParentRetainedNoStream"
)

// controlSeams admits the test-runtime branches and then the environment hooks, in name order.
var controlSeams = []injectedports.Row{
	{
		Kind: injectedports.KindRuntime, Package: "internal/otelrecord", Port: "recordHome", Owner: "recordHome",
		Purpose:   "a test binary may not write a telemetry record below the user's fallback Bench home, while an isolated explicit home still records",
		Isolation: "the branch reads only the test binary's own runtime state; each consumer names its home through t.Setenv, which restores it at cleanup",
		Consumers: []string{"internal/otelrecord/processor_test.go::TestNewProviderRefusesTheFallbackHomeUnderTest", "internal/otelrecord/processor_test.go::TestNewProviderRefusesAnExplicitBenchHomeEqualToTheFallback", "internal/otelrecord/processor_test.go::TestNewProviderRecordsAnExplicitBenchHomeElsewhere"},
	},
	seamHook("ANTHROPIC_API_KEY", "internal/models", "inventory", "a present or absent provider key drives the API row and the missing-key report",
		"internal/models/models_test.go::TestCommandEmitsCodexOpenAIAndAnthropicRows"),
	seamHook("BENCH_AGENT", "internal/shift", "loop", "a named agent adapter drives the shift loop and its record attribute",
		"internal/shift/pass_test.go::TestACommittedPassCarriesItsCommit"),
	seamHook("BENCH_COMMAND_OBSERVE", "cmd/bench", "main", "the observation switch makes the entry print the command it routes to",
		"internal/systemtest/charge_evidence_test.go::TestEvidenceLocationIndependence"),
	seamHook("BENCH_CONFORMANCE_ROOT", "internal/gate", "createTestLog", "a root-conformance root keys each test log by the package path below that root",
		"internal/gate/testlog_exec_test.go::TestCreateTestLogKeysThePackageByTheConformanceRoot"),
	seamHook("BENCH_EVIDENCE_PAUSE", "internal/chargeevidence", "PauseFromEnvironment", "a held evidence stage lets the parent observe a live holder of the store lock",
		"internal/systemtest/charge_cleanup_test.go::TestEvidenceCleanupReaderExclusion"),
	seamHook("BENCH_GATE", "internal/gate", "buildSubjectInMode", "an absent or named gate script drives the gate subject and its exit-three refusal",
		"internal/gate/run_outcomes_test.go::TestGateRunWithoutGateExitsThree"),
	seamHook("BENCH_GATE_BACKGROUND", "internal/gate", "BackgroundCommand", "the resolved background marker drives the background grade's re-entry",
		"internal/gate/background_grade_test.go::TestBackgroundGradeRunsAsTheTargetExecutable"),
	seamHook("BENCH_GATE_BASELINE", "internal/gate", "phasePolicyRoot", "a baseline policy root drives the phase policy that a package-scope run reads",
		"internal/coverage/citation_execution_test.go::TestCitationPackageScope"),
	seamHook("BENCH_GATE_DELTA", "internal/gate", "deltaScheduled", "a delta schedule drives the checkpoint's delta package set",
		"internal/gate/delta_checkpoint_test.go::TestDeltaCheckpointRunsTheDeltaSet"),
	seamHook("BENCH_GATE_LOG_PATH", "internal/gate", "inheritGateRunLog", "an inherited run-log path drives the child's record inheritance", gateStream),
	seamHook("BENCH_GATE_LOG_ROOT", "internal/gate", "inheritGateRunLog", "an inherited run-log root drives the child's record inheritance", gateStream),
	seamHook("BENCH_GATE_LOG_RUN", "internal/gate", "inheritGateRunLog", "an inherited run identity drives the child's record inheritance", gateStream),
	seamHook("BENCH_GATE_READER", "internal/gate", "readerScheduled", "a reader schedule drives the reader set and its malformed-schedule refusal",
		"internal/gate/reader_set_test.go::TestReaderScheduleRefusesAMalformedSchedule"),
	seamHook("BENCH_GATE_STREAM_PATH", "internal/gate", "inheritGateStreamFile", "an inherited stream path drives whether the child reuses the parent stream", gateStream),
	seamHook("BENCH_GATE_TESTLOG", "internal/gate", "withTestLog", "a test-log directory drives one log per tested package",
		"internal/gate/testlog_exec_test.go::TestWholeTestPhaseWritesOneLogPerPackage"),
	seamHook("BENCH_GATE_WIDTH", "internal/gate", "widthFlags", "a width input drives the package and parallel widths of the test argv",
		"internal/gate/width_test.go::TestWidthInputSetsTheTestWidths"),
	seamHook("BENCH_HOME", "internal/benchhome", "Dir", "a private Bench home isolates every record write and drives the home resolution",
		"cmd/bench/otel_hook_seams_test.go::TestAHandedOffResolutionRecordsItsLine"),
	seamHook("BENCH_KIT", "internal/gate", "KitValue", "a named kit drives the gate checkpoint route",
		"cmd/bench/gate_route_test.go::TestGateCheckpointRoute"),
	seamHook("BENCH_LINK_FAULT", "internal/adopt", "promoteWithLease", "an abrupt exit at a chosen promotion change, which a real filesystem cannot produce on demand, drives the interrupted-transaction recovery",
		"internal/systemtest/compatibility_test.go::TestCompatibilityInterruptedRepair", "internal/systemtest/lifecycle_test.go::TestLifecycleUnlinkRemovalInterruption"),
	seamHook("BENCH_MAX_DIR_FILES", "internal/structure", "envInt", "a directory file cap drives the structure census", "internal/structure/structure_test.go::TestSpaceDirPreserved"),
	seamHook("BENCH_MAX_ITERS", "internal/shift", "parseBoundedInt", "an iteration bound drives the shift loop's pass count",
		"internal/shift/pass_test.go::TestACommittedPassCarriesItsCommit"),
	seamHook("BENCH_MAX_LINES", "internal/structure", "envInt", "a line cap drives the structure growth refusal", "internal/structure/growth_test.go::TestGrowthAcceptUnreadableIsLoud"),
	seamHook("BENCH_MAX_WALL", "internal/shift", "parseWallDuration", "a wall bound drives the shift checkpoint teardown",
		"internal/shift/lifetime_test.go::TestCheckpointTeardownFaultKeepsResponse"),
	seamHook("BENCH_MODEL", "cmd/bench", "resolveModel", "a named tier drives the model resolution and its record",
		"cmd/bench/otel_hook_seams_test.go::TestAResolutionRecordsNoOperatorText"),
	seamHook("BENCH_OFFLINE", "internal/bounds", "Offline", "the offline switch makes Offline hold, so a refresh starts no git and names the switch",
		"internal/refresh/refresh_test.go::TestRefreshOfflineStartsNoGitAndNamesFlag"),
	seamHook("BENCH_OTEL_ROOT", "internal/otelrecord", "AttachHandoff", "an ambient record root shows that the gate ignores a root it was not handed",
		"internal/systemtest/otel_gate_test.go::TestOtelGateIgnoresAnAmbientRecordRoot"),
	seamHook("BENCH_OTEL_TRACEPARENT", "internal/otelrecord", "AttachHandoff", "an ambient trace parent shows that the gate ignores a parent it was not handed",
		"internal/systemtest/otel_gate_test.go::TestOtelGateIgnoresAnAmbientRecordRoot"),
	seamHook("BENCH_PREFLIGHT_ARTIFACTS", "internal/releasepreflight", phaseReader, "a fixture phase runs the artifacts phase through the chosen binary", rootProbe),
	seamHook("BENCH_PREFLIGHT_DATE", "internal/releasepreflight", "runner.runVulnerability", "a chosen day drives the vulnerability exception expiry at the entry",
		"internal/releasepreflight/vulnerability_test.go::TestVulnerabilityEntryDateOverride"),
	seamHook("BENCH_PREFLIGHT_EVIDENCE_READY_FILE", "internal/releaseevidence", "waitForEvidenceProbe", "a held readiness point lets the native probe observe evidence assembly before promotion", rootProbe),
	seamHook("BENCH_PREFLIGHT_GATE", "internal/releasepreflight", phaseReader, "a fixture gate phase drives the protected-group refusal and the chosen-binary composition",
		"internal/releasepreflight/lifetime_test.go::TestExternalIncomplete", rootProbe),
	seamHook("BENCH_PREFLIGHT_RACE", "internal/releasepreflight", phaseReader, "a real child probe observes the race phase's isolated home",
		"internal/releasepreflight/external_test.go::TestExternalPhaseGitStartsNoAutoMaintenance", rootProbe),
	seamHook("BENCH_PREFLIGHT_REF", "internal/releasepreflight", "runner.checkIdentity", "an installed release ref drives the identity agreement and disagreement",
		"internal/releasepreflight/identity_test.go::TestCheckIdentityRefusesADisagreement"),
	seamHook("BENCH_PREFLIGHT_SMOKE", "internal/releasepreflight", phaseReader, "a fixture phase runs the smoke phase through the chosen binary", rootProbe),
	seamHook("BENCH_PREFLIGHT_VET", "internal/releasepreflight", phaseReader, "a fixture phase runs the vet phase through the chosen binary", rootProbe),
	seamHook("BENCH_PREFLIGHT_VULNERABILITY", "internal/releasepreflight", "runner.runVulnerability", "a fixture scanner drives the scanner cancellation, its certification, and the native composition",
		"internal/releasepreflight/lifetime_test.go::TestScannerCancel", rootProbe),
	seamHook("BENCH_PROCESS_RESOURCES", "internal/subprocess/protection", "Inherited", "a malformed resource vector drives the inherited-resource refusal",
		"internal/bounds/lifetime_test.go::TestBoundsMalformedResources"),
	seamHook("BENCH_REQUIRE_CAPABILITIES", "internal/gate", "strict", "strict mode drives the capability skip rows",
		"internal/gate/report_test.go::TestCapabilitySkipsInTwoClassesPrintOneLine"),
	seamHook("BENCH_RUN_BINARY", "internal/gate", "productionRunBinaryOwner", "an inherited run binary makes the gate reuse or own that run binary",
		"internal/gate/background_grade_test.go::TestBackgroundGradeRunsAsTheTargetExecutable"),
	seamHook("BENCH_SKIP_LOG", "internal/capability", "writeSkipLine", "a skip log path shows that the capability line is written before the skip",
		"internal/capability/capability_test.go::TestCapabilityWritesLineBeforeSkip"),
	seamHook("BENCH_TEST_UNBOUNDED_WAITS", "internal/bounds", "VerdictWindow", "the verdict override removes a wait's own deadline while parent cancellation and fixed grace stay",
		"internal/bounds/bounds_test.go::TestUnboundedParentCancellation", "internal/bounds/lifetime_test.go::TestBoundsFixedShutdown"),
	seamHook("BENCH_WRAPPER", "internal/adopt", "resolvedWrapper", "a named wrapper places the broker manifest beside it",
		"internal/adopt/broker_test.go::TestWriteBrokerManifestBindsTheRunningBrokerIdentity"),
	seamHook("CLAUDE_PROJECT_DIR", "internal/harness", "claudeRemove", "a project directory outside a repository drives the remove refusal",
		"internal/harness/worktree_test.go::TestWorktreeCommandRemoveRefusesAProjectDirOutsideARepository"),
	seamHook("CODEX_HOME", "internal/adopt", "configurationHome", "a missing configuration home drives the compatibility refusal",
		"internal/adopt/compatibility_test.go::TestCompatibilityMissingConfigurationHome"),
	seamHook("ENVMAN_LOAD", "internal/sessioninspect", "environmentPhase", "an absent or empty load marker drives the partial environment closure line",
		"internal/systemtest/session_start_test.go::TestSessionStartTE5DiscoversGoWithMarkerAbsentOrEmpty"),
	seamHook("HOME", "internal/adopt", "currentDoctorEnv", "a private home keeps the doctor repair's home read off the developer's home",
		"internal/systemtest/compatibility_test.go::TestCompatibilityConcurrentWriters"),
	seamHook("OPENAI_API_KEY", "internal/models", "inventory", "a present or absent provider key drives the API row and the missing-key report",
		"internal/models/models_test.go::TestCommandEmitsCodexOpenAIAndAnthropicRows"),
	seamHook("PATH", "internal/adopt", "currentDoctorEnv", "a search path that holds a healthy managed shim drives the doctor's shim row to its green exit",
		"internal/adopt/adopt_test.go::TestDoctorKitSourceCheckoutRowsAreGreen"),
	seamHook("SHELL", "internal/worktree", "subshellShell", "a fixture shell drives the worktree shell's signal handling",
		"internal/worktree/subshell_lifetime_test.go::TestShellSignal"),
}
