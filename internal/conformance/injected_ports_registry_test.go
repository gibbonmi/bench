package conformance

import "github.com/gibbonmi/bench/internal/conformance/injectedports"

// auditedPortPackages are the packages whose injected ports the registry covers, the
// audit's inventory, no wider. A package outside it that declares a port is out of scope
// until a reviewer binds it in. A package inside it that derives nothing fails closed.
// The variable and setter census is not limited to this list; it reads every production
// root of the graded tree.
var auditedPortPackages = []string{
	"internal/assessment",
	"internal/chargeevidence",
	"internal/gitguard",
	"internal/preflight",
	"internal/publication",
	"internal/subprocess/protection",
}

// injectedPortRegistry is the advertisement half: one row per port the derivation finds
// in auditedPortPackages, then one row per variable seam and per setter that the census
// finds anywhere in the production roots. Rows are grouped by kind and package and
// ordered as the derivation reports them. A diff against a re-derived inventory therefore
// reads straight down.
var injectedPortRegistry = append(append([]injectedports.Row{
	{Package: "internal/assessment", Port: "FileOps", TestFile: "internal/assessment/replacement_test.go", TestName: "TestAssessmentReplacement"},
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
	{
		Package: "internal/gitguard", Port: "Checker",
		TestFile: "internal/gitguard/checker_junction_test.go", TestName: "TestClassifyRealCheckerResolvedComposition",
	},
	{
		Package: "internal/publication", Port: "Registry",
		Exempt: "a real-producer test needs a live npm registry, and NPMCLIRegistry is runbook-only: it never performs gate egress or touches a credential (the Registry port comment in internal/publication/registry.go states each adapter's gate coverage)",
	},
	{
		Package: "internal/preflight", Port: "reviewEvidenceObserver",
		TestFile: "internal/preflight/review_charge_test.go", TestName: "TestEvidenceReviewCollectors",
	},
	{
		Package: "internal/subprocess/protection", Port: "Probe",
		TestFile: "internal/subprocess/protection/observe_test.go", TestName: "TestObserveGroupPresent",
	},
}, variableSeams...), setterSeams...)

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
		"internal/diff/command_test.go::TestCommandRetriesThenRefusesSnapshotDrift", "internal/preflight/charge_test.go::TestChargeSnapshotMovement"),
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
		"internal/guards/guards_cleanup_test.go::TestScanWaitsForCancelledWorkerCleanup"),
	seamVariable("internal/guards", "guardScanTimeout", "Command", "a short scan bound drives the checked-in timeout responses", serialIsolation,
		"internal/guards/guards_test.go::TestCommandPreservesCheckedInEnumerationTimeoutPrimaryResponse"),
	seamVariable("internal/guards", "inspectGuard", "Scan", "a held inspection shows that a cancelled scan waits for its worker cleanup", joinIsolation,
		"internal/guards/guards_cleanup_test.go::TestScanWaitsForCancelledWorkerCleanup"),
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
