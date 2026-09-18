package conformance

// This file carries the live-tree test classification the tier registry reads. It moved
// out of tier_test.go for that file's line budget.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// classifiedLiveTreeTests is the sole classification of tests that read the live tree. Most
// read it to construct a mutation or driver fixture. Some enforce a policy on the kit root
// directly, because the kit root is that policy's subject. The staleness test below requires
// every entry to still be a detected live-tree reader. The hidden-inventory check supplies
// the reverse direction: it rejects every detected reader absent from this classification or
// the executable registry.
var classifiedLiveTreeTests = map[string]bool{
	"TestCanaryFixtureRegistryClassifiesEveryFixture":              true,
	"TestConformanceMetaBites":                                     true,
	"TestCoreSubprocessFailuresUseProbeFormatter":                  true,
	"TestDecisionMapIntegrityCheckValidatesEveryCandidate":         true,
	"TestDecisionMapIntegrityFixtureInventoryRejectsDeletion":      true,
	"TestEvidenceBoundedActionSentencesArePinned":                  true,
	"TestEvidenceBoundedActionRejectsTheRetiredPair":               true,
	"TestEveryRetainedFixtureBitesThroughRegisteredOwner":          true,
	"TestFixtureBiteProofArchitecture":                             true,
	"TestGuidanceProseBudgetCanaryFixtureBites":                    true,
	"TestGuidanceProseBudgetsHoldOnTheLiveTree":                    true,
	"TestSkillDescriptionBudgetsHoldOnTheLiveTree":                 true,
	"TestPreparedReviewGuidanceHoldsOnTheLiveTree":                 true,
	"TestProseMechanicsCanaryFixturesBite":                         true,
	"TestProseMechanicsHoldsOnTheLiveTree":                         true,
	"TestProfileLaneTableHoldsOnTheLiveTree":                       true,
	"TestProseExclusionRowsStayInsideTheApprovedSet":               true,
	"TestHarnessUsesBenchConformanceRootAsGradedRoot":              true,
	"TestInjectedPortRegistryCheckBites":                           true,
	"TestInvalidOrderedSetRedsAndWidensToTheFullTier":              true,
	"TestLocalCaptureFilesAreNotTracked":                           true,
	"TestAXIMembershipExpectationBitesInBothDirections":            true,
	"TestAXIGuidanceContractBites":                                 true,
	"TestNativeWorkflowEvidenceEdgeBites":                          true,
	"TestReleaseWorkflowPublicationBites":                          true,
	"TestOccurrenceLedgerMigrationCheckBitesOnFT98Count":           true,
	"TestOfflineSmokeSliceOneProofIsExecutableNotTokenOnly":        true,
	"TestBranchNativeArchitectureCensus":                           true,
	"TestRecurrenceMaintenanceContractCheckBites":                  true,
	"TestResidualCheckCallsCrossCompileMatrix":                     true,
	"TestRetiredConformanceFixturesDoNotLeaveShellTwinMessages":    true,
	"TestRetroImprovementMarkersFixtureInventoryRejectsDeletion":   true,
	"TestRetroImprovementMarkersFixturesCoverEveryDiagnosticClass": true,
	"TestRetainedWorkflow":                                         true,
	"TestRoadmapDetailIntegrityFixturesCoverEveryDiagnosticClass":  true,
	"TestRowNextGrammarBindsTableToParserTokens":                   true,
	"TestRowNextGrammarFixturesCoverEveryDiagnosticClass":          true,
	"TestRowNextGrammarFixtureInventoryRejectsDeletion":            true,
	"TestRoadmapDetailIntegrityFixtureInventoryRejectsDeletion":    true,
	"TestRunConformanceAcceptsHostileRootPath":                     true,
	"TestRunConformanceChecksExecutableGitMode":                    true,
	"TestRunConformanceDistinguishesAbsentAndEmptyInputs":          true,
	"TestScopeOutsideTierIsRedAndRunsNothing":                      true,
	"TestShipConformanceRunNamesDeclaredTests":                     true,
	"TestRetiredReleaseFixtureReplacementsArePresent":              true,
	"TestWorkflowCadenceAnchorsRejectDeletionAndSwap":              true,
	"TestSkillsIndexConformanceCarriesNoSecondReader":              true,
	"TestSpecTicketHandoffWorkflowFixturesAreComplete":             true,
	"TestTimingOrderStable":                                        true,
	"TestUnknownScopeIsRedAndRunsNothing":                          true,
}

func classifiedLiveTreeTest(name string) bool { return classifiedLiveTreeTests[name] }

func TestClassifiedLiveTreeInventoryNamesDetectedTests(t *testing.T) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, ".", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	detected := map[string]bool{}
	for _, file := range packages["conformance"].Files {
		for _, declaration := range file.Decls {
			fn, ok := declaration.(*ast.FuncDecl)
			if ok && fn.Body != nil && testChecksLiveTree(fn.Body) {
				detected[fn.Name.Name] = true
			}
		}
	}
	for name := range classifiedLiveTreeTests {
		if !detected[name] {
			t.Errorf("classified live-tree test %s is absent or no longer reads the live tree", name)
		}
	}
}
