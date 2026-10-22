//go:build system

package systemtest

// pendingMigration is the serial reason of a test that has not yet moved into the
// parallel set.
const pendingMigration = "pending migration"

// ownerRepositories is the serial reason of a test whose subject is the owner's own
// repositories.
const ownerRepositories = "the owner's repositories are the test's subject"

// systemSerialTests is the one list of the package's serial tests. Each row names a
// top-level test that calls no t.Parallel() and the reason it stays serial. Go runs these
// tests one at a time before every parallel test, so a serial test may share the owner's
// repositories, name a ledger field, or bind process state.
var systemSerialTests = map[string]string{
	"TestCandidateInstallLandsItsNewFields":                             pendingMigration,
	"TestCandidateInstallPublishesTheCandidate":                         pendingMigration,
	"TestCandidateInstallSurvivesAGapLanding":                           pendingMigration,
	"TestEvidenceCleanupInterrupt":                                      pendingMigration,
	"TestEvidenceCleanupReaderExclusion":                                pendingMigration,
	"TestEvidenceCleanupWriterExclusion":                                pendingMigration,
	"TestEvidenceConcurrentPublication":                                 pendingMigration,
	"TestEvidenceInterruptedPublication":                                pendingMigration,
	"TestEvidenceLocationIndependence":                                  pendingMigration,
	"TestEvidenceNestedDirectory":                                       pendingMigration,
	"TestEvidenceQuotaConcurrency":                                      pendingMigration,
	"TestEvidenceReleaseRead":                                           pendingMigration,
	"TestEvidenceSiblingRead":                                           pendingMigration,
	"TestEvidenceStatelessProcesses":                                    pendingMigration,
	"TestExecBoundsChildOutput":                                         pendingMigration,
	"TestExecForwardsStdinUnderBound":                                   pendingMigration,
	"TestExecInterruptKeepsOutput":                                      pendingMigration,
	"TestExecKeepsChildExitUnderBound":                                  pendingMigration,
	"TestExecNestedBenchBoundsOnce":                                     pendingMigration,
	"TestExecNestedBenchKeepsChildRetained":                             pendingMigration,
	"TestExecReturnsAtChildExit":                                        pendingMigration,
	"TestExecStreamsLargeChild":                                         pendingMigration,
	"TestExecStreamsLongLine":                                           pendingMigration,
	"TestOtelCommitAndLandingRecordJourney":                             pendingMigration,
	"TestOtelCrashKeepsStartedPhaseLine":                                pendingMigration,
	"TestOtelGateIgnoresAnAmbientRecordRoot":                            pendingMigration,
	"TestOtelGateRecordJourney":                                         pendingMigration,
	"TestOtelGateRecordNamesTheStampedVersion":                          pendingMigration,
	"TestOtelHookRecordJourney":                                         pendingMigration,
	"TestOtelShiftTraceJourney":                                         pendingMigration,
	"TestOtelWorktreeVerbRecordJourney":                                 pendingMigration,
	"TestOwnerSelectionChild":                                           ownerRepositories,
	"TestSelectedExecutableComposition":                                 ownerRepositories,
	"TestStatusRouteExecutesUnclaimedCleanup":                           pendingMigration,
	"TestTreeTargetRefusesStaleKitBuild":                                pendingMigration,
	"TestTreeTargetRunsInNamedWorktree":                                 pendingMigration,
	"TestWorktreeLandRouteDoesNotRederiveTheBrokerPlatform":             pendingMigration,
	"TestWorktreeLandRouteGivesRecoveredGoPathToBroker":                 pendingMigration,
	"TestWorktreeLandRouteRefusesEveryUnauthenticatedBroker":            pendingMigration,
	"TestWorktreeLandRouteRefusesInheritedRoutingBeforeRepositoryReads": pendingMigration,
	"TestWorktreeLandRouteSelectsOneOwnerFromEveryDirectory":            pendingMigration,
}
