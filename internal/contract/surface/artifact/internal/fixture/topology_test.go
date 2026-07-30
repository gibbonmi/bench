package fixture

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

var artifactTests = map[string][]string{
	"posture":       {"TestArtifactBuilderHonorsHermeticDefault", "TestBuildCachePostureUnderGoproxyOff", "TestSharedCacheBuildRemovesStaleReproducibilityRecord", "TestGoBuildIgnoresCheckoutTopology"},
	"prepared":      {"TestArtifactPromotionIsAtomicAndExclusive", "TestArtifactSourceStagesCommittedHostPlan", "TestSharedCacheBuildPromotesNoRecord", "TestSharedCacheBuildRestoresRecordOnInterruptedPromotion", "TestOfflineArchiveProjection", "TestPackedArtifactRunsSetupOfflineFromASpacedPrefix", "TestSharedArtifactSetFailsClosedAfterEarlierStagingFailure", "TestSharedArtifactSetBuildIsLazy", "TestSharedArtifactSetAttributesMutation", "TestSharedArtifactSetIsReadOnly", "TestArtifactBuilderRejectsDirtyAndUntrackedSourceState", "TestArtifactBuilderRefusesMissingBinaryPinManifest", "TestArtifactSourceSkipsWhenHostTargetIsAbsent"},
	"offline":       {"TestReleasePlanProjectsDerivedArchiveInventory", "TestReleasePlanProjectsRelocatedPackageEvidence", "TestNativeProofAggregatorRejectsOneTargetOmission", "TestNativeProofAggregatorRejectsDigestMismatch", "TestAuthoritativeNativeProofBehaviorCanary", "TestOfflineRegistryDerivesAcceptedTargetsFromReleasePlan", "TestOfflineArchiveBuildRefusesOutputItCannotAccountFor", "TestOfflineInstructionsVerifyOnlyTargetArchive", "TestOfflineSmokeRunsThePublicJourneyAndAttributesMutations", "TestOfflineNetworkSentinelDeniesUndeclaredEgress", "TestOfflineSmokeRecoversFromEveryStageInterruption", "TestReproducibilityComparatorRejectsArtifactAndEvidenceMutations", "TestOfflineSmokeRequiresApprovedReleaseEvidence", "TestReleaseArtifactVerifierRequiresFullyApprovedEvidence"},
	"distributable": {"TestDistributableArtifactContracts", "TestArtifactBuilderRejectsSpecialReleaseEvidenceInput"},
}

func TestSubjectPackageTopology(t *testing.T) {
	root := filepath.Join(contract.SubjectRoot(t), "internal", "contract", "surface", "artifact")
	for subject, want := range artifactTests {
		files, err := filepath.Glob(filepath.Join(root, subject, "*_test.go"))
		if err != nil {
			t.Fatal(err)
		}
		got := topLevelTests(t, files)
		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Fatalf("%s top-level tests = %v, want %v", subject, got, want)
		}
		main, err := os.ReadFile(filepath.Join(root, subject, "main_test.go"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(main), "fixture.Run(m,") || strings.Contains(string(main), "SharedBuildCacheEnv") {
			t.Fatalf("%s TestMain must call fixture.Run without an inline shared-cache token", subject)
		}
	}
	if files, err := filepath.Glob(filepath.Join(root, "*_test.go")); err != nil || len(files) != 0 {
		t.Fatalf("artifact root owns migrated tests: %v, %v", files, err)
	}
}

func topLevelTests(t *testing.T, files []string) []string {
	t.Helper()
	var names []string
	for _, path := range files {
		parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range parsed.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name != "TestMain" && strings.HasPrefix(fn.Name.Name, "Test") && fn.Recv == nil && fn.Type.Params != nil && len(fn.Type.Params.List) == 1 {
				names = append(names, fn.Name.Name)
			}
		}
	}
	return names
}
