package conformance

// The tier registry itself: which checks each tier runs, that metadata and bound
// functions cannot drift apart, that the filtered inner run selects real tests, and
// that every executed check leaves one timing line in a stable order.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

const releaseEvidenceProbeCheck = "release-evidence-probe"

var requiredMetaChecks = []string{
	"conformance-meta",
	"conformance-canary-families",
}

func checkConformanceMeta(kitRoot string, tier registry.Tier) []string {
	return checkConformanceMetaForPartition(kitRoot, tier, executableNames(tier), nil)
}

func checkConformanceMetaForPartition(kitRoot string, tier registry.Tier, executed, inherited []string) []string {
	var diags []string
	diags = append(diags, registryAgreementDiags()...)
	diags = append(diags, metaMembershipDiags()...)
	diags = append(diags, canaryOwnershipDiags()...)
	diags = append(diags, partitionDiags(tier, executed, inherited)...)
	diags = append(diags, checkInputProfileDiags(readIfExists(filepath.Join(kitRoot, "projects", "benchkit.md")))...)
	diags = append(diags, hiddenLiveTreeDiags(filepath.Join(kitRoot, "internal", "conformance"))...)
	return diags
}

func checkInputProfileDiags(profile string) []string {
	rendered := map[string]string{}
	var order, diags []string
	found := false
	for _, line := range strings.Split(profile, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			if found {
				break
			}
			continue
		}
		cells := strings.Split(strings.Trim(trimmed, "|"), "|")
		if len(cells) < 2 {
			continue
		}
		name := strings.Trim(strings.TrimSpace(cells[0]), "`")
		input := strings.Trim(strings.TrimSpace(cells[1]), "`")
		if !found {
			found = name == "conformance check" && input == "input source"
			continue
		}
		if strings.Trim(name, "-:") == "" {
			continue
		}
		if _, duplicate := rendered[name]; duplicate {
			diags = append(diags, "profile check-input row duplicated: "+name)
		} else {
			order = append(order, name)
		}
		rendered[name] = input
	}
	if !found {
		return []string{"profile check-input table missing"}
	}
	want := map[string]bool{}
	var wantOrder []string
	for _, check := range registry.Checks {
		if check.Meta {
			continue
		}
		want[check.Name] = true
		wantOrder = append(wantOrder, check.Name)
		if rendered[check.Name] != string(check.Inputs) {
			diags = append(diags, fmt.Sprintf("profile check-input row stale: %s renders %q, registry declares %q", check.Name, rendered[check.Name], check.Inputs))
		}
	}
	if !slices.Equal(order, wantOrder) {
		diags = append(diags, "profile check-input rows are not in registry order")
	}
	for name := range rendered {
		if !want[name] {
			diags = append(diags, "profile check-input row unknown: "+name)
		}
	}
	return diags
}

func registryAgreementDiags() []string {
	seen := make(map[string]bool, len(registry.Checks))
	implementations := make(map[string]string, len(registry.Checks))
	var diags []string
	for _, check := range registry.Checks {
		if check.Name == "" || seen[check.Name] {
			diags = append(diags, "conformance meta registry has a missing or duplicate check name")
			continue
		}
		seen[check.Name] = true
		binding, bound := conformanceChecks[check.Name]
		if !bound {
			diags = append(diags, "conformance meta check "+check.Name+" has no executable binding")
			continue
		}
		if got := binding.identity(); check.Implementation == "" || got != check.Implementation {
			diags = append(diags, fmt.Sprintf("conformance meta check %s binding mismatch: registry names implementation %q, executable binding is %q", check.Name, check.Implementation, got))
		} else if previous := implementations[got]; previous != "" {
			diags = append(diags, fmt.Sprintf("conformance meta implementation %s is registered more than once: %s and %s", got, previous, check.Name))
		} else {
			implementations[got] = check.Name
		}
		if binding.tier != check.Tier {
			diags = append(diags, fmt.Sprintf("conformance meta check %s tier mismatch: registry says %s, executable binding says %s", check.Name, check.Tier, binding.tier))
		}
		if binding.subject != check.Subject {
			diags = append(diags, fmt.Sprintf("conformance meta check %s subject mismatch: registry says %s, executable binding says %s", check.Name, check.Subject, binding.subject))
		}
		if check.Meta {
			if check.Inputs != "" {
				diags = append(diags, "conformance meta authorization check "+check.Name+" declares reusable inputs")
			}
		} else if !check.Inputs.Valid() {
			diags = append(diags, "conformance meta check "+check.Name+" has no valid input declaration")
		}
		switch check.Tier {
		case registry.Dev, registry.Ship:
		default:
			diags = append(diags, "conformance meta check "+check.Name+" has no valid tier declaration")
		}
		switch check.Subject {
		case registry.SubjectRoot, registry.SubjectKitRoot, registry.SubjectRootAndKitRoot:
		default:
			diags = append(diags, "conformance meta check "+check.Name+" has no valid subject declaration")
		}
	}
	for name := range conformanceChecks {
		if !seen[name] {
			diags = append(diags, "conformance meta executable binding "+name+" has no registry check")
		}
	}
	return diags
}

func metaMembershipDiags() []string {
	var diags []string
	required := make(map[string]bool, len(requiredMetaChecks))
	for _, name := range requiredMetaChecks {
		required[name] = true
		check, found := registry.Find(name)
		if !found || !check.Meta {
			diags = append(diags, "conformance meta authorization check "+name+" is not always-on")
		}
	}
	for _, check := range registry.Checks {
		if check.Meta && !required[check.Name] {
			diags = append(diags, "conformance meta classifies semantic check "+check.Name+" as always-on")
		}
	}
	return diags
}

func canaryOwnershipDiags() []string {
	var diags []string
	for _, family := range registry.Families() {
		name, bound := registry.FamilyCheck(family)
		check, found := registry.Find(name)
		if !bound || !found {
			diags = append(diags, fmt.Sprintf("conformance meta canary family %s has no registered check owner", family))
			continue
		}
		if check.Tier != registry.Dev {
			diags = append(diags, fmt.Sprintf("conformance meta canary family %s owns non-dev check %s", family, name))
		}
	}
	return diags
}

func executableNames(tier registry.Tier) []string {
	var names []string
	for _, check := range registry.Checks {
		if binding, found := conformanceChecks[check.Name]; found && binding.runsAt(tier) {
			names = append(names, check.Name)
		}
	}
	return names
}

func partitionDiags(tier registry.Tier, executed, inherited []string) []string {
	want := registry.Names(tier)
	wantSet := make(map[string]bool, len(want))
	for _, name := range want {
		wantSet[name] = true
	}
	seen := make(map[string]string, len(executed)+len(inherited))
	var diags []string
	for _, part := range []struct {
		name   string
		checks []string
	}{{"executed", executed}, {"inherited", inherited}} {
		for _, name := range part.checks {
			if previous := seen[name]; previous != "" {
				diags = append(diags, fmt.Sprintf("conformance meta partition check %s appears in both or more than once: %s and %s", name, previous, part.name))
				continue
			}
			seen[name] = part.name
			if !wantSet[name] {
				diags = append(diags, fmt.Sprintf("conformance meta partition includes tier-inapplicable check %s in %s", name, part.name))
			}
		}
	}
	for _, name := range want {
		if seen[name] == "" {
			diags = append(diags, "conformance meta partition omits check "+name)
		}
	}
	return diags
}

func hiddenLiveTreeDiags(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []string{"conformance meta cannot inventory live-tree tests: " + err.Error()}
	}
	var diags []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			diags = append(diags, "conformance meta cannot parse "+entry.Name()+" for live-tree tests: "+err.Error())
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil || !strings.HasPrefix(fn.Name.Name, "Test") || fn.Name.Name == registry.RootConformanceTest {
				continue
			}
			if strings.HasPrefix(fn.Name.Name, registry.RootConformanceTest) || (testChecksLiveTree(fn.Body) && !classifiedLiveTreeTest(fn.Name.Name)) {
				diags = append(diags, "conformance meta unregistered live-tree assertion "+fn.Name.Name)
			}
		}
	}
	return diags
}

func testChecksLiveTree(body *ast.BlockStmt) bool {
	liveRoots := map[string]bool{}
	harnesses := map[string]bool{}
	gradedRootControlled := false
	ast.Inspect(body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if ok {
			if selector, ok := call.Fun.(*ast.SelectorExpr); ok && selector.Sel.Name == "Setenv" && len(call.Args) > 0 {
				if key, ok := call.Args[0].(*ast.BasicLit); ok && strings.Trim(key.Value, "`\"") == "BENCH_CONFORMANCE_ROOT" {
					gradedRootControlled = true
				}
			}
		}
		assign, ok := node.(*ast.AssignStmt)
		if !ok || len(assign.Rhs) == 0 {
			return true
		}
		if call, ok := assign.Rhs[0].(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "NewHarness" && len(assign.Lhs) > 0 {
				if harness, ok := assign.Lhs[0].(*ast.Ident); ok {
					harnesses[harness.Name] = true
				}
			}
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "findKitRoot" && len(assign.Lhs) > 0 {
				if root, ok := assign.Lhs[0].(*ast.Ident); ok {
					liveRoots[root.Name] = true
				}
			}
		}
		if selector, ok := assign.Rhs[0].(*ast.SelectorExpr); ok && selector.Sel.Name == "Root" && newHarnessCall(selector.X) && len(assign.Lhs) > 0 {
			if root, ok := assign.Lhs[0].(*ast.Ident); ok {
				liveRoots[root.Name] = true
			}
		}
		return true
	})
	found := false
	ast.Inspect(body, func(node ast.Node) bool {
		if assign, ok := node.(*ast.AssignStmt); ok {
			for _, expr := range assign.Rhs {
				if kitRootExpression(expr, harnesses) {
					found = true
					return false
				}
			}
		}
		selector, ok := node.(*ast.SelectorExpr)
		if ok && !gradedRootControlled {
			if selector.Sel.Name == "Root" && newHarnessCall(selector.X) {
				found = true
				return false
			}
			if ident, ok := selector.X.(*ast.Ident); ok && harnesses[ident.Name] {
				switch selector.Sel.Name {
				case "Root", "ReadRootFile", "RequireExecutable", "Run":
					found = true
					return false
				}
			}
			if selector.Sel.Name == "ReadRootFile" && newHarnessCall(selector.X) {
				found = true
				return false
			}
		}
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if checkCall(call) {
			for _, arg := range call.Args {
				if kitRootExpression(arg, harnesses) {
					found = true
					return false
				}
			}
		}
		if filesystemReadCall(call) && len(call.Args) > 0 && (kitRootExpression(call.Args[0], harnesses) || expressionUsesLiveRoot(call.Args[0], liveRoots) || relativePathLiteral(call.Args[0])) {
			found = true
			return false
		}
		ident, isIdent := call.Fun.(*ast.Ident)
		if !isIdent || !strings.HasPrefix(ident.Name, "check") {
			return true
		}
		for _, arg := range call.Args {
			if root, ok := arg.(*ast.Ident); ok && liveRoots[root.Name] {
				found = true
			}
		}
		return !found
	})
	return found
}

func kitRootExpression(expr ast.Expr, harnesses map[string]bool) bool {
	switch value := expr.(type) {
	case *ast.SelectorExpr:
		if value.Sel.Name != "KitRoot" {
			return false
		}
		if newHarnessCall(value.X) {
			return true
		}
		ident, ok := value.X.(*ast.Ident)
		return ok && harnesses[ident.Name]
	case *ast.CallExpr:
		selector, ok := value.Fun.(*ast.SelectorExpr)
		if ok && selector.Sel.Name == "KitPath" {
			ident, ok := selector.X.(*ast.Ident)
			return ok && harnesses[ident.Name]
		}
		for _, arg := range value.Args {
			if kitRootExpression(arg, harnesses) {
				return true
			}
		}
	}
	return false
}

func checkCall(call *ast.CallExpr) bool {
	ident, ok := call.Fun.(*ast.Ident)
	return ok && strings.HasPrefix(ident.Name, "check")
}

func filesystemReadCall(call *ast.CallExpr) bool {
	if ident, ok := call.Fun.(*ast.Ident); ok {
		switch ident.Name {
		case "readIfExists", "exists", "frontmatterField":
			return true
		}
		return false
	}
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	if !ok || pkg.Name != "os" {
		return false
	}
	switch selector.Sel.Name {
	case "ReadFile", "ReadDir", "Stat", "Lstat", "Open":
		return true
	}
	return false
}

func expressionUsesLiveRoot(expr ast.Expr, liveRoots map[string]bool) bool {
	uses := false
	ast.Inspect(expr, func(node ast.Node) bool {
		ident, ok := node.(*ast.Ident)
		if ok && liveRoots[ident.Name] {
			uses = true
			return false
		}
		return !uses
	})
	return uses
}

// classifiedLiveTreeTests is the sole classification of tests that read the live tree only
// to construct a mutation or driver fixture rather than enforce policy directly. The
// staleness test below requires every entry still to be a detected live-tree reader; the
// hidden-inventory check supplies the reverse direction by rejecting every detected reader
// absent from this classification or the executable registry.
var classifiedLiveTreeTests = map[string]bool{
	"TestCanaryFixtureRegistryClassifiesEveryFixture":           true,
	"TestConformanceMetaBites":                                  true,
	"TestCoreSubprocessFailuresUseProbeFormatter":               true,
	"TestCoverageMapValidationFixtureBite":                      true,
	"TestDataHandlingDerivationFixtureBite":                     true,
	"TestDecisionMapIntegrityCheckValidatesEveryCandidate":      true,
	"TestDecisionMapIntegrityFixtureInventoryRejectsDeletion":   true,
	"TestDecisionMapIntegrityFixturesBite":                      true,
	"TestDocsCurrencyTokenDietAndWorkflowFixturesBite":          true,
	"TestGuidanceProseBudgetCanaryFixtureBites":                 true,
	"TestGuidanceProseBudgetsHoldOnTheLiveTree":                 true,
	"TestHarnessUsesBenchConformanceRootAsGradedRoot":           true,
	"TestInvalidOrderedSetRedsAndWidensToTheFullTier":           true,
	"TestLineRoutingFixturesBite":                               true,
	"TestAXIMembershipExpectationBitesInBothDirections":         true,
	"TestAXIGuidanceContractBites":                              true,
	"TestLoadValidityMetadataFixturesBite":                      true,
	"TestNativeWorkflowEvidenceEdgeBites":                       true,
	"TestOccurrenceLedgerMigrationCheckBitesOnFT158Count":       true,
	"TestOfflineSmokeSliceOneProofIsExecutableNotTokenOnly":     true,
	"TestBranchNativeArchitectureCensus":                        true,
	"TestPackageCoreAndGuardFixturesBite":                       true,
	"TestRecurrenceMaintenanceContractCheckBites":               true,
	"TestResidualCheckCallsCrossCompileMatrix":                  true,
	"TestRetiredConformanceFixturesDoNotLeaveShellTwinMessages": true,
	"TestRunConformanceAcceptsHostileRootPath":                  true,
	"TestRunConformanceChecksExecutableGitMode":                 true,
	"TestRunConformanceDistinguishesAbsentAndEmptyInputs":       true,
	"TestScopeOutsideTierIsRedAndRunsNothing":                   true,
	"TestShipConformanceRunNamesDeclaredTests":                  true,
	"TestSkillsIndexAndCommandAdapterFixturesBite":              true,
	"TestWorkflowCadenceAnchorsRejectDeletionAndSwap":           true,
	"TestSpecTicketHandoffWorkflowFixturesAreComplete":          true,
	"TestTimingOrderStable":                                     true,
	"TestUnknownScopeIsRedAndRunsNothing":                       true,
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

func relativePathLiteral(expr ast.Expr) bool {
	literal, ok := expr.(*ast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		return false
	}
	path := strings.Trim(literal.Value, "`\"")
	return path != "" && !filepath.IsAbs(path)
}

func newHarnessCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	ident, ok := call.Fun.(*ast.Ident)
	return ok && ident.Name == "NewHarness"
}

func TestTierMembership(t *testing.T) {
	dev, ship := registry.Names(registry.Dev), registry.Names(registry.Ship)
	if len(dev) == 0 {
		t.Fatal("the dev tier runs no checks at all")
	}
	if slices.Contains(dev, releaseEvidenceProbeCheck) {
		t.Fatalf("the dev tier still runs %s, the ~372 s probe the split moves:\n%s", releaseEvidenceProbeCheck, strings.Join(dev, "\n"))
	}
	if !slices.Contains(ship, releaseEvidenceProbeCheck) {
		t.Fatalf("the ship tier does not run %s, so the probe runs nowhere:\n%s", releaseEvidenceProbeCheck, strings.Join(ship, "\n"))
	}
	for _, name := range dev {
		if !slices.Contains(ship, name) {
			t.Fatalf("the ship tier drops dev check %s; ship green has to reprove everything dev green claims", name)
		}
	}
}

// TestEveryCheckCarriesATier closes the gap Tier's string underlying type leaves open:
// a row whose tier is misspelled or omitted holds "", which RunsAt reads as neither dev
// nor ship, so the check silently stops running on every commit and survives every
// membership assertion — those compare against name lists the untiered check is in or
// out of on both sides.
func TestEveryCheckCarriesATier(t *testing.T) {
	for _, check := range registry.Checks {
		switch check.Tier {
		case registry.Dev, registry.Ship:
		default:
			t.Errorf("registry check %s carries tier %q, which is neither %q nor %q, so no tier executes it on a commit", check.Name, check.Tier, registry.Dev, registry.Ship)
		}
	}
}

func TestRegistryBindsEveryCheck(t *testing.T) {
	for _, check := range registry.Checks {
		if _, bound := conformanceChecks[check.Name]; !bound {
			t.Errorf("registry check %s has no bound function", check.Name)
		}
	}
	for name := range conformanceChecks {
		if !slices.Contains(registry.Names(registry.Ship), name) {
			t.Errorf("bound function %s has no registry row, so it carries no tier", name)
		}
	}
}

func TestConformanceMetaBites(t *testing.T) {
	checks := slices.Clone(registry.Checks)
	bindings := cloneCheckBindings(conformanceChecks)
	t.Cleanup(func() {
		registry.Checks = checks
		conformanceChecks = bindings
	})
	h := NewHarness(t)
	if diags := checkConformanceMeta(h.KitRoot, registry.Dev); len(diags) != 0 {
		t.Fatalf("registered conformance inventory has diagnostics:\n%s", strings.Join(diags, "\n"))
	}

	t.Run("missing declaration", func(t *testing.T) {
		for i := range registry.Checks {
			if registry.Checks[i].Name == "line-routing" {
				registry.Checks[i].Subject = ""
			}
		}
		if !containsDiagnostic(registryAgreementDiags(), "line-routing has no valid subject declaration") {
			t.Fatal("removing a subject declaration did not make conformance meta red")
		}
		registry.Checks = slices.Clone(checks)
	})
	t.Run("missing input declaration", func(t *testing.T) {
		for i := range registry.Checks {
			if registry.Checks[i].Name == "line-routing" {
				registry.Checks[i].Inputs = ""
			}
		}
		if !containsDiagnostic(registryAgreementDiags(), "line-routing has no valid input declaration") {
			t.Fatal("removing an input declaration did not make conformance meta red")
		}
		registry.Checks = slices.Clone(checks)
	})
	t.Run("check input profile drift", func(t *testing.T) {
		profile := readIfExists(filepath.Join(h.KitRoot, "projects", "benchkit.md"))
		stale := strings.Replace(profile, "| `line-routing` | `catch-all` |", "| `line-routing` | `go-source` |", 1)
		if !containsDiagnostic(checkInputProfileDiags(stale), "line-routing") {
			t.Fatal("changing one profile input token did not make conformance meta red")
		}
		for i := range registry.Checks {
			if registry.Checks[i].Name == "line-routing" {
				registry.Checks[i].Inputs = registry.InputGoSource
			}
		}
		if !containsDiagnostic(checkInputProfileDiags(profile), "line-routing") {
			t.Fatal("changing one registry input token did not make conformance meta red")
		}
		registry.Checks = slices.Clone(checks)
	})
	t.Run("implementation registered twice", func(t *testing.T) {
		binding := conformanceChecks["bounds-policy"]
		binding.implementation = checkLineRouting
		conformanceChecks["bounds-policy"] = binding
		for i := range registry.Checks {
			if registry.Checks[i].Name == "bounds-policy" {
				registry.Checks[i].Implementation = "checkLineRouting"
			}
		}
		if !containsDiagnostic(registryAgreementDiags(), "checkLineRouting is registered more than once") {
			t.Fatal("registering one implementation twice did not make conformance meta red")
		}
		registry.Checks = slices.Clone(checks)
		conformanceChecks = cloneCheckBindings(bindings)
	})
	t.Run("profile binding removed from meta", func(t *testing.T) {
		setCheckMeta("conformance-canary-families", false)
		if !containsDiagnostic(metaMembershipDiags(), "conformance-canary-families is not always-on") {
			t.Fatal("removing a profile binding from meta did not make conformance meta red")
		}
		registry.Checks = slices.Clone(checks)
	})
	t.Run("canary owner removed", func(t *testing.T) {
		for i := range registry.Checks {
			if registry.Checks[i].Name == "line-routing" {
				registry.Checks = append(registry.Checks[:i], registry.Checks[i+1:]...)
				break
			}
		}
		if !containsDiagnostic(canaryOwnershipDiags(), "line-routing has no registered check owner") {
			t.Fatal("removing a canary owner did not make conformance meta red")
		}
		registry.Checks = slices.Clone(checks)
	})
	t.Run("partition omission overlap and duplicate", func(t *testing.T) {
		executed := registry.Names(registry.Dev)
		missing := slices.Clone(executed[1:])
		if !containsDiagnostic(partitionDiags(registry.Dev, missing, nil), "omits check "+executed[0]) {
			t.Fatal("an omitted check did not make conformance meta red")
		}
		if !containsDiagnostic(partitionDiags(registry.Dev, executed, []string{executed[0]}), "appears in both") {
			t.Fatal("an executed/inherited overlap did not make conformance meta red")
		}
		duplicated := append(slices.Clone(executed), executed[0])
		if !containsDiagnostic(partitionDiags(registry.Dev, duplicated, nil), "more than once") {
			t.Fatal("a duplicate executed check did not make conformance meta red")
		}
	})
	t.Run("semantic check marked meta", func(t *testing.T) {
		setCheckMeta("docs-currency-workflow", true)
		if !containsDiagnostic(metaMembershipDiags(), "docs-currency-workflow as always-on") {
			t.Fatal("moving a semantic check into meta did not make conformance meta red")
		}
		registry.Checks = slices.Clone(checks)
	})
	t.Run("functions swapped", func(t *testing.T) {
		first, second := conformanceChecks["line-routing"], conformanceChecks["bounds-policy"]
		first.implementation, second.implementation = second.implementation, first.implementation
		conformanceChecks["line-routing"], conformanceChecks["bounds-policy"] = first, second
		diags := registryAgreementDiags()
		if !containsDiagnostic(diags, "line-routing binding mismatch") || !containsDiagnostic(diags, "bounds-policy binding mismatch") {
			t.Fatalf("swapping functions did not name both mismatched bindings:\n%s", strings.Join(diags, "\n"))
		}
		conformanceChecks = cloneCheckBindings(bindings)
	})
	t.Run("dev registry row moved to ship", func(t *testing.T) {
		for i := range registry.Checks {
			if registry.Checks[i].Name == "line-routing" {
				registry.Checks[i].Tier = registry.Ship
			}
		}
		if !containsDiagnostic(registryAgreementDiags(), "line-routing tier mismatch") {
			t.Fatal("moving a dev registry row to ship did not make conformance meta red")
		}
		registry.Checks = slices.Clone(checks)
	})
	t.Run("meta check grades root", func(t *testing.T) {
		for i := range registry.Checks {
			if registry.Checks[i].Name == "conformance-meta" {
				registry.Checks[i].Subject = registry.SubjectRoot
			}
		}
		if !containsDiagnostic(registryAgreementDiags(), "conformance-meta subject mismatch") {
			t.Fatal("moving meta grading to root did not make conformance meta red")
		}
	})
}

func TestHiddenLiveTreeInventoryBites(t *testing.T) {
	for _, test := range []struct {
		name   string
		source string
	}{
		{"prefix", "package conformance\nimport \"testing\"\nfunc TestRootConformanceHiddenPolicy(t *testing.T) {}\n"},
		{"ordinary", "package conformance\nimport \"testing\"\nfunc TestHiddenLiveTreePolicy(t *testing.T) { _ = NewHarness(t).Root }\n"},
		{"direct kit root", "package conformance\nimport \"testing\"\nfunc TestHiddenKitPolicy(t *testing.T) { _ = NewHarness(t).KitRoot }\n"},
		{"kit root check", "package conformance\nimport \"testing\"\nfunc TestHiddenKitCheck(t *testing.T) { h := NewHarness(t); _ = checkHidden(h.KitRoot) }\n"},
		{"kit path read", "package conformance\nimport (\"os\"; \"testing\")\nfunc TestHiddenKitPath(t *testing.T) { h := NewHarness(t); _, _ = os.ReadFile(h.KitPath(\"README.md\")) }\n"},
		{"relative read", "package conformance\nimport (\"os\"; \"testing\")\nfunc TestHiddenRelativeRead(t *testing.T) { _, _ = os.ReadFile(\"../../README.md\") }\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "hidden_test.go"), []byte(test.source), 0o644); err != nil {
				t.Fatal(err)
			}
			diags := hiddenLiveTreeDiags(dir)
			if !containsDiagnostic(diags, "Test") {
				t.Fatalf("unregistered %s live-tree test did not make meta red: %v", test.name, diags)
			}
		})
	}
}

func TestOptionalCanaryFixtureSurfaceGuardsAbsence(t *testing.T) {
	root := t.TempDir()
	if diags := checkCanaryFixtureCompliance(root); len(diags) != 0 {
		t.Fatalf("absent optional canary surface produced diagnostics: %v", diags)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench-compliance-canary"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if diags := checkCanaryFixtureCompliance(root); !containsDiagnostic(diags, "kit root LICENSE is missing") {
		t.Fatalf("present optional canary surface did not run its targeted grading: %v", diags)
	}
}

func TestMigratedOptionalSurfacesGuardAbsence(t *testing.T) {
	root := t.TempDir()
	for name, check := range map[string]func(string) []string{
		"gate-entry-contract":     checkGateEntryContract,
		"offline-smoke-proof":     checkOfflineSmokeProof,
		"handoff-shape":           checkHandoffShape,
		"harness-prefix":          checkHarnessPrefix,
		"package-shipped-surface": checkPackageShippedSurface,
	} {
		if diags := check(root); len(diags) != 0 {
			t.Errorf("%s graded an absent optional fixture surface: %v", name, diags)
		}
	}
}

func cloneCheckBindings(source map[string]checkBinding) map[string]checkBinding {
	clone := make(map[string]checkBinding, len(source))
	for name, binding := range source {
		clone[name] = binding
	}
	return clone
}

func setCheckMeta(name string, meta bool) {
	for i := range registry.Checks {
		if registry.Checks[i].Name == name {
			registry.Checks[i].Meta = meta
		}
	}
}

func TestDevTierExecutesExactlyDevChecks(t *testing.T) {
	root := gitInitedRoot(t)
	RunConformance(root, NewHarness(t).KitRoot, registry.Dev, "")

	got := timingNames(t, root)
	if want := registry.Names(registry.Dev); !slices.Equal(got, want) {
		t.Fatalf("dev run executed\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if slices.Contains(got, releaseEvidenceProbeCheck) {
		t.Fatalf("a dev run executed %s", releaseEvidenceProbeCheck)
	}
}

func TestTimingLinePerCheck(t *testing.T) {
	root := gitInitedRoot(t)
	RunConformance(root, NewHarness(t).KitRoot, registry.Dev, "")

	lines := registry.ReadTimingLines(root)
	if want := len(registry.Names(registry.Dev)); len(lines) != want {
		t.Fatalf("timing lines = %d, want one per executed check (%d):\n%s", len(lines), want, strings.Join(lines, "\n"))
	}
}

func TestTimingOrderStable(t *testing.T) {
	root := gitInitedRoot(t)
	kitRoot := NewHarness(t).KitRoot

	RunConformance(root, kitRoot, registry.Dev, "")
	first := timingNames(t, root)
	RunConformance(root, kitRoot, registry.Dev, "")
	second := timingNames(t, root)

	if !slices.Equal(first, second) {
		t.Fatalf("timing order differs between runs of one tree:\n%s\nversus\n%s", strings.Join(first, "\n"), strings.Join(second, "\n"))
	}
	if want := registry.Names(registry.Dev); !slices.Equal(first, want) {
		t.Fatalf("timing order is not the registry's order:\n%s\nwant\n%s", strings.Join(first, "\n"), strings.Join(want, "\n"))
	}
}

func TestScopedRunExecutesOnlyTheNamedCheck(t *testing.T) {
	root := gitInitedRoot(t)
	const scope = "line-routing"
	RunConformance(root, NewHarness(t).KitRoot, registry.Dev, scope)

	if got := timingNames(t, root); !slices.Equal(got, []string{scope}) {
		t.Fatalf("scoped run executed\n%s\nwant only %s", strings.Join(got, "\n"), scope)
	}
}

func TestOrderedSetRunsMetaAndSelectedOrdinaryChecksInRegistryOrder(t *testing.T) {
	root := gitInitedRoot(t)
	selected := "line-routing,subcommand-routing"
	inherited := inheritedSelection(registry.Dev, selected)
	RunConformanceSelection(root, NewHarness(t).KitRoot, registry.Dev, "", &selected, &inherited)

	var want []string
	for _, check := range registry.Checks {
		if check.Meta || check.Name == "line-routing" || check.Name == "subcommand-routing" {
			want = append(want, check.Name)
		}
	}
	if got := timingNames(t, root); !slices.Equal(got, want) {
		t.Fatalf("ordered-set run executed\n%s\nwant meta plus selected checks in registry order\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestEquivalentFullAndScopedTimingShowsFewerOrdinaryChecks(t *testing.T) {
	harness := NewHarness(t)
	fullRoot := gitInitedRoot(t)
	scopedRoot := gitInitedRoot(t)
	RunConformance(fullRoot, harness.KitRoot, registry.Dev, "")
	selected := "line-routing,subcommand-routing"
	inherited := inheritedSelection(registry.Dev, selected)
	RunConformanceSelection(scopedRoot, harness.KitRoot, registry.Dev, "", &selected, &inherited)

	full := timingNames(t, fullRoot)
	scoped := timingNames(t, scopedRoot)
	if len(scoped) >= len(full) {
		t.Fatalf("scoped timing has %d checks, want fewer than full timing's %d", len(scoped), len(full))
	}
	for _, name := range scoped {
		check, _ := registry.Find(name)
		if !check.Meta && name != "line-routing" && name != "subcommand-routing" {
			t.Fatalf("scoped timing retained inherited ordinary check %q: %v", name, scoped)
		}
	}
	if !slices.Equal(full, registry.Names(registry.Dev)) {
		t.Fatalf("full timing = %v, want registry inventory", full)
	}
}

func TestShipTierIgnoresOuterOrderedSelection(t *testing.T) {
	root := gitInitedRoot(t)
	selected := "release-evidence-probe"
	inherited := ""
	RunConformanceSelection(root, NewHarness(t).KitRoot, registry.Ship, "", &selected, &inherited)
	lines := timingNames(t, root)
	if !slices.Equal(lines, registry.Names(registry.Ship)) {
		t.Fatalf("ship timing names = %v, want full ship inventory %v", lines, registry.Names(registry.Ship))
	}
}

func TestInvalidOrderedSetRedsAndWidensToTheFullTier(t *testing.T) {
	for _, selected := range []string{
		"line-routing,line-routing",
		"no-such-check",
		"release-evidence-probe",
		"subcommand-routing,line-routing",
		"conformance-meta",
	} {
		t.Run(selected, func(t *testing.T) {
			root := gitInitedRoot(t)
			inherited := strings.Join(registry.OrdinaryNames(registry.Dev), ",")
			diags := RunConformanceSelection(root, NewHarness(t).KitRoot, registry.Dev, "", &selected, &inherited)
			if !containsDiagnostic(diags, "ordered selection") {
				t.Fatalf("invalid ordered set %q produced diagnostics %v", selected, diags)
			}
			if got, want := timingNames(t, root), registry.Names(registry.Dev); !slices.Equal(got, want) {
				t.Fatalf("invalid ordered set %q executed\n%s\nwant full tier\n%s", selected, strings.Join(got, "\n"), strings.Join(want, "\n"))
			}
		})
	}
}

func inheritedSelection(tier registry.Tier, selected string) string {
	want := map[string]bool{}
	if selected != "" {
		for _, name := range strings.Split(selected, ",") {
			want[name] = true
		}
	}
	var inherited []string
	for _, name := range registry.OrdinaryNames(tier) {
		if !want[name] {
			inherited = append(inherited, name)
		}
	}
	return strings.Join(inherited, ",")
}

// TestUnknownScopeIsRedAndRunsNothing pins the posture a silent fallback would break:
// a scope naming no check re-pays the full run and hides the drift that renamed it.
// The hostile value is deliberate — the diagnostic quotes the scope, so control bytes
// have to survive as escapes rather than as a mangled line.
func TestUnknownScopeIsRedAndRunsNothing(t *testing.T) {
	root := recordedScopeRoot(t)
	const scope = "no-such-check\x01\n"
	diags := RunConformance(root, NewHarness(t).KitRoot, registry.Dev, scope)

	if len(diags) != 1 || !strings.Contains(diags[0], `"no-such-check\x01\n"`) {
		t.Fatalf("unknown scope: want one diagnostic quoting the scope, got %q", diags)
	}
	if got := timingNames(t, root); len(got) != 0 {
		t.Fatalf("unknown scope left timing lines standing:\n%s", strings.Join(got, "\n"))
	}
}

// TestScopeOutsideTierIsRedAndRunsNothing covers the scope that exists but sits on a
// tier this run does not grade: executing zero checks in silence would read as green
// and leave the fixture reporting a baffling did-not-bite.
func TestScopeOutsideTierIsRedAndRunsNothing(t *testing.T) {
	root := recordedScopeRoot(t)
	diags := RunConformance(root, NewHarness(t).KitRoot, registry.Dev, releaseEvidenceProbeCheck)

	if len(diags) != 1 || !strings.Contains(diags[0], releaseEvidenceProbeCheck) || !strings.Contains(diags[0], string(registry.Dev)) {
		t.Fatalf("tier-mismatched scope: want one diagnostic naming the scope and the tier, got %q", diags)
	}
	if got := timingNames(t, root); len(got) != 0 {
		t.Fatalf("tier-mismatched scope left timing lines standing:\n%s", strings.Join(got, "\n"))
	}
}

// recordedScopeRoot is a graded root whose timing file already holds a completed run's
// lines. Against a pristine root a red posture cannot be told from one that left the
// last run's record standing; against this one, only an emptied file passes.
func recordedScopeRoot(t *testing.T) string {
	t.Helper()
	root := gitInitedRoot(t)
	RunConformance(root, NewHarness(t).KitRoot, registry.Dev, "line-routing")
	if len(registry.ReadTimingLines(root)) == 0 {
		t.Fatal("the seeding run recorded no timing lines, so the posture assertion proves nothing")
	}
	return root
}

// gitInitedRoot is a graded root a timing file can live in: the file sits under the
// root's own git dir, so a bare temp dir gives it nowhere to go.
func gitInitedRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	return root
}

func timingNames(t *testing.T, root string) []string {
	t.Helper()
	var names []string
	for _, line := range registry.ReadTimingLines(root) {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			t.Fatalf("timing line %q is not an index, a check name, and a duration", line)
		}
		names = append(names, fields[1])
	}
	return names
}
