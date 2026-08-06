package fixture

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

var artifactTests = map[string][]string{
	"posture":       {"TestArtifactBuilderHonorsHermeticDefault", "TestArtifactCallersSelectUnsealedBuilderMode", "TestArtifactInstallNeverExposesAnAbsentPriorOutput", "TestArtifactModeCommandTrace", "TestArtifactModeCommandTraceBitesOnPublishThenDelete", "TestArtifactModeCommandTraceExcludesPublication", "TestArtifactModeCompileFailurePreservesPriorPair", "TestArtifactModeFailureTablePreservesPriorPair", "TestArtifactModeHandlesDashLedRelativeOutputLiterally", "TestArtifactModeRefusesHostileOutputTypes", "TestArtifactModeRefusesMalformedSelectorsWithoutChangingPriorPair", "TestArtifactModeRerunsConvergeUnsealed", "TestArtifactSignalAfterInstallLeavesNoStaleSeal", "TestBuildCachePostureUnderGoproxyOff", "TestGoBuildIgnoresCheckoutTopology", "TestGoBuildInterruptionsPreservePriorPairAndRemoveStaging", "TestMalformedBuilderGrammarInvokesNoGoAndPreservesPriorPair", "TestSharedCacheBuildRemovesStaleReproducibilityRecord"},
	"prepared":      {"TestArtifactPromotionIsAtomicAndExclusive", "TestArtifactSourceStagesCommittedHostPlan", "TestSharedCacheBuildPromotesNoRecord", "TestSharedCacheBuildRestoresRecordOnInterruptedPromotion", "TestOfflineArchiveProjection", "TestPackedArtifactRunsSetupOfflineFromASpacedPrefix", "TestSharedArtifactSetFailsClosedAfterEarlierStagingFailure", "TestSharedArtifactSetBuildIsLazy", "TestSharedArtifactSetAttributesMutation", "TestSharedArtifactSetIsReadOnly", "TestArtifactBuilderRejectsDirtyAndUntrackedSourceState", "TestArtifactBuilderRefusesMissingBinaryPinManifest", "TestArtifactSourceSkipsWhenHostTargetIsAbsent"},
	"offline":       {"TestReleasePlanProjectsDerivedArchiveInventory", "TestReleasePlanProjectsRelocatedPackageEvidence", "TestNativeProofAggregatorRejectsOneTargetOmission", "TestNativeProofAggregatorRejectsDigestMismatch", "TestAuthoritativeNativeProofBehaviorCanary", "TestOfflineRegistryDerivesAcceptedTargetsFromReleasePlan", "TestOfflineArchiveBuildRefusesOutputItCannotAccountFor", "TestOfflineInstructionsVerifyOnlyTargetArchive", "TestOfflineSmokeRunsThePublicJourneyAndAttributesMutations", "TestOfflineNetworkSentinelDeniesUndeclaredEgress", "TestOfflineSmokeRecoversFromEveryStageInterruption", "TestReproducibilityComparatorRejectsArtifactAndEvidenceMutations", "TestOfflineSmokeRequiresApprovedReleaseEvidence", "TestReleaseArtifactVerifierRequiresFullyApprovedEvidence"},
	"distributable": {"TestDistributableArtifactContracts", "TestArtifactBuilderRejectsSpecialReleaseEvidenceInput"},
}

func TestSubjectPackageTopology(t *testing.T) {
	root := filepath.Join(contract.SubjectRoot(t), "internal", "contract", "surface", "artifact")
	allowed := map[string]bool{"internal/fixture": true}
	for subject := range artifactTests {
		allowed[subject] = true
	}
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files, err := filepath.Glob(filepath.Join(path, "*.go"))
		if err != nil {
			return err
		}
		if len(files) != 0 && !allowed[filepath.ToSlash(relative)] {
			return fmt.Errorf("artifact Go package %s is not allowed", filepath.ToSlash(relative))
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
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
		assertSubjectRunnerAndPolicy(t, subject, filepath.Join(root, subject))
	}
}

func assertSubjectRunnerAndPolicy(t *testing.T, subject, directory string) {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(directory, "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	mainFound := false
	for _, path := range files {
		parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(parsed, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.Ident:
				if node.Name == "SharedBuildCacheEnv" {
					t.Fatalf("%s contains inline shared-cache identifier", subject)
				}
			case *ast.BasicLit:
				if node.Kind == token.STRING {
					value, err := strconv.Unquote(node.Value)
					if err != nil {
						t.Fatal(err)
					}
					if value == "BENCH_SHARED_BUILD_CACHE" {
						t.Fatalf("%s contains inline shared-cache literal", subject)
					}
				}
			}
			return true
		})
		for _, decl := range parsed.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "TestMain" {
				if mainFound {
					t.Fatalf("%s has more than one TestMain", subject)
				}
				mainFound = true
				if !strings.HasSuffix(path, "_test.go") || !callsFixtureRun(fn, importName(parsed, "testing", "testing"), importName(parsed, fixtureImportPath, "fixture")) {
					t.Fatalf("%s TestMain must directly call Run from %s", subject, fixtureImportPath)
				}
			}
		}
	}
	if !mainFound {
		t.Fatalf("%s has no TestMain", subject)
	}
}

const fixtureImportPath = "github.com/gibbonmi/bench/internal/contract/surface/artifact/internal/fixture"

func importName(file *ast.File, importPath, defaultName string) string {
	for _, imported := range file.Imports {
		path, err := strconv.Unquote(imported.Path.Value)
		if err != nil {
			return ""
		}
		if path != importPath {
			continue
		}
		if imported.Name == nil {
			return defaultName
		}
		if imported.Name.Name != "." && imported.Name.Name != "_" {
			return imported.Name.Name
		}
	}
	return ""
}

func callsFixtureRun(main *ast.FuncDecl, testingName, fixtureName string) bool {
	if testingName == "" || fixtureName == "" || main.Recv != nil || main.Type.TypeParams != nil || main.Type.Params == nil || len(main.Type.Params.List) != 1 || main.Type.Results != nil || main.Body == nil || len(main.Body.List) != 1 {
		return false
	}
	parameter := main.Type.Params.List[0]
	if len(parameter.Names) != 1 || parameter.Names[0].Name != "m" {
		return false
	}
	pointer, ok := parameter.Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	testingType, ok := pointer.X.(*ast.SelectorExpr)
	if !ok || testingType.Sel.Name != "M" {
		return false
	}
	testingPackage, ok := testingType.X.(*ast.Ident)
	if !ok || testingPackage.Name != testingName {
		return false
	}
	expression, ok := main.Body.List[0].(*ast.ExprStmt)
	if !ok {
		return false
	}
	call, ok := expression.X.(*ast.CallExpr)
	if !ok || len(call.Args) != 2 {
		return false
	}
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Run" {
		return false
	}
	fixturePackage, ok := selector.X.(*ast.Ident)
	argument, argumentOK := call.Args[0].(*ast.Ident)
	return ok && fixturePackage.Name == fixtureName && argumentOK && argument.Name == "m"
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
