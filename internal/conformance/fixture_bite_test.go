package conformance

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

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestEveryRetainedFixtureBitesThroughRegisteredOwner(t *testing.T) {
	h := NewHarness(t)
	canaryDir := filepath.Join(h.KitRoot, "tests", "canary")
	completed := map[string]bool{}
	runFixtureUniverse(t, canaryDir, func(t *testing.T, name string, fixture canary.Fixture) {
		runFixtureBite(t, h.KitRoot, name, fixture)
		completed[name] = true
	})
	want, err := canary.Fixtures(canaryDir)
	requireFixtureNoError(t, err)
	if !sameFixtureSet(completed, want) {
		t.Fatalf("completed fixture proofs do not equal the current producer: got %d, want %d", len(completed), len(want))
	}
}

func TestFixtureUniverseDerivesFamilyAndExplicitCheckBindings(t *testing.T) {
	canaryDir := filepath.Join(t.TempDir(), "canary")
	for _, fixture := range []struct{ family, name, check string }{
		{"package-core-guard", "family-derived", ""},
		{"package-core-guard", "explicit-ship", "release-evidence-probe"},
	} {
		dir := filepath.Join(canaryDir, fixture.family, fixture.name)
		requireFixtureNoError(t, os.MkdirAll(dir, 0o755))
		requireFixtureNoError(t, os.WriteFile(filepath.Join(dir, "EXPECT"), []byte("synthetic diagnostic\n"), 0o644))
		if fixture.check != "" {
			requireFixtureNoError(t, os.WriteFile(filepath.Join(dir, "CHECK"), []byte(fixture.check+"\n"), 0o644))
		}
	}
	observed := map[string]bool{}
	runFixtureUniverse(t, canaryDir, func(_ *testing.T, name string, _ canary.Fixture) { observed[name] = true })
	if got, want := len(observed), 2; got != want {
		t.Fatalf("synthetic fixture callbacks = %d, want %d", got, want)
	}
}

func TestFixtureBiteProofArchitecture(t *testing.T) {
	h := NewHarness(t)
	path := filepath.Join(h.KitRoot, "internal", "conformance", "fixture_bite_test.go")
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	functions := map[string]*ast.FuncDecl{}
	for _, declaration := range file.Decls {
		if function, ok := declaration.(*ast.FuncDecl); ok {
			functions[function.Name.Name] = function
		}
	}
	live := functions["TestEveryRetainedFixtureBitesThroughRegisteredOwner"]
	runner := functions["runFixtureUniverse"]
	if live == nil || runner == nil {
		t.Fatal("fixture proof needs its live test and producer-derived runner")
	}
	if callsNamed(live, "runFixtureUniverse") != 1 || rangeCount(live) != 0 {
		t.Fatal("live fixture proof must make one producer-derived runner call without a second proof loop")
	}
	if callsNamed(runner, "Fixtures") != 1 {
		t.Fatal("fixture proof runner must derive its universe from canary.Fixtures")
	}
	if callsNamed(live, "runFixtureBite") != 1 {
		t.Fatal("live fixture proof callback must call runFixtureBite once before recording completion")
	}
	if !recordsFixtureAfterProof(live) {
		t.Fatal("live fixture proof must record completed[name] only after runFixtureBite returns")
	}
	fixtures, err := canary.Fixtures(filepath.Join(h.KitRoot, "tests", "canary"))
	if err != nil {
		t.Fatal(err)
	}
	if literal := retainedFixtureNameLiteral(fixtures, live, runner); literal != "" {
		t.Fatalf("fixture proof stores retained fixture name %q instead of deriving it", literal)
	}
	if calls, err := packageRunFixtureBiteCalls(filepath.Dir(path)); err != nil {
		t.Fatal(err)
	} else if calls != 1 {
		t.Fatalf("internal/conformance has %d runFixtureBite calls, want one universal proof caller", calls)
	}
}

func TestRetiredReleaseFixtureReplacementsArePresent(t *testing.T) {
	if diagnostics := retiredReleaseFixtureReplacementDiagnostics(NewHarness(t).KitRoot); len(diagnostics) != 0 {
		t.Fatalf("retired release fixture replacement census:\n%s", strings.Join(diagnostics, "\n"))
	}
}

func retiredReleaseFixtureReplacementDiagnostics(root string) []string {
	replacements := map[string]string{
		"internal/releaseevidence/release_index_test.go":    "TestReleaseIndexBindsComponentManifestDigest",
		"internal/releaseevidence/package_artifact_test.go": "TestBuildReleaseEvidenceIncludesRegisteredPackageEvidence",
	}
	var diagnostics []string
	for rel, symbol := range replacements {
		file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(root, filepath.FromSlash(rel)), nil, 0)
		if err != nil {
			diagnostics = append(diagnostics, fmt.Sprintf("parse release replacement %s: %v", rel, err))
			continue
		}
		found := false
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if ok && function.Name.Name == symbol {
				found = true
			}
		}
		if !found {
			diagnostics = append(diagnostics, fmt.Sprintf("release replacement %s does not declare %s", rel, symbol))
		}
	}
	for _, rel := range []string{
		"tests/canary/package-core-guard/release-digest-binding-omitted",
		"tests/canary/package-core-guard/release-package-evidence-omitted",
		"tests/canary/package-core-guard/release-evidence-probe-base.txt",
	} {
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel))); err == nil {
			diagnostics = append(diagnostics, fmt.Sprintf("retired release fixture %s remains", rel))
		} else if !os.IsNotExist(err) {
			diagnostics = append(diagnostics, fmt.Sprintf("inspect retired release fixture %s: %v", rel, err))
		}
	}
	return diagnostics
}

func callsNamed(function *ast.FuncDecl, name string) int {
	var calls int
	ast.Inspect(function.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch target := call.Fun.(type) {
		case *ast.Ident:
			if target.Name == name {
				calls++
			}
		case *ast.SelectorExpr:
			if target.Sel.Name == name {
				calls++
			}
		}
		return true
	})
	return calls
}

func rangeCount(function *ast.FuncDecl) int {
	var ranges int
	ast.Inspect(function.Body, func(node ast.Node) bool {
		if _, ok := node.(*ast.RangeStmt); ok {
			ranges++
		}
		return true
	})
	return ranges
}

func recordsFixtureAfterProof(live *ast.FuncDecl) bool {
	var callback *ast.FuncLit
	ast.Inspect(live.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		name, _ := call.Fun.(*ast.Ident)
		if name == nil || name.Name != "runFixtureUniverse" || len(call.Args) != 3 {
			return true
		}
		callback, _ = call.Args[2].(*ast.FuncLit)
		return false
	})
	if callback == nil {
		return false
	}
	proof, record := -1, -1
	for index, statement := range callback.Body.List {
		if statementCalls(statement, "runFixtureBite") {
			proof = index
		}
		if recordsCompletedFixture(statement) {
			record = index
		}
	}
	return proof >= 0 && record > proof
}

func statementCalls(statement ast.Stmt, name string) bool {
	found := false
	ast.Inspect(statement, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if target, ok := call.Fun.(*ast.Ident); ok && target.Name == name {
			found = true
		}
		return true
	})
	return found
}

func recordsCompletedFixture(statement ast.Stmt) bool {
	assignment, ok := statement.(*ast.AssignStmt)
	if !ok || len(assignment.Lhs) != 1 {
		return false
	}
	index, ok := assignment.Lhs[0].(*ast.IndexExpr)
	if !ok {
		return false
	}
	completed, completedOK := index.X.(*ast.Ident)
	name, nameOK := index.Index.(*ast.Ident)
	return completedOK && nameOK && completed.Name == "completed" && name.Name == "name"
}

func retainedFixtureNameLiteral(fixtures map[string]canary.Fixture, functions ...*ast.FuncDecl) string {
	for _, function := range functions {
		var literal string
		ast.Inspect(function.Body, func(node ast.Node) bool {
			value, ok := node.(*ast.BasicLit)
			if !ok || value.Kind != token.STRING {
				return true
			}
			text, err := strconv.Unquote(value.Value)
			if err == nil && fixtures[text].Dir != "" {
				literal = text
				return false
			}
			return true
		})
		if literal != "" {
			return literal
		}
	}
	return ""
}

func packageRunFixtureBiteCalls(dir string) (int, error) {
	packages, err := parser.ParseDir(token.NewFileSet(), dir, func(info os.FileInfo) bool {
		return strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		return 0, err
	}
	var calls int
	for _, file := range packages["conformance"].Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			if target, ok := call.Fun.(*ast.Ident); ok && target.Name == "runFixtureBite" {
				calls++
			}
			return true
		})
	}
	return calls, nil
}

func TestSpecTicketHandoffWorkflowFixturesAreComplete(t *testing.T) {
	required := []struct{ fixture, diagnostic string }{
		{"write-spec-identified-coverage-default", ".agents/skills/bench-craft-spec/SKILL.md Template dropped the identified reduced-column acceptance-map default"},
		{"write-spec-unique-row-id", ".agents/skills/bench-craft-spec/SKILL.md Template dropped the unique spec-local row-ID default"},
		{"write-spec-ownership-fences", ".agents/skills/bench-craft-spec/SKILL.md Template dropped the craft-spec-owned Ownership fences section"},
		{"write-spec-fence-approval", ".agents/skills/bench-craft-spec/SKILL.md Template approval paragraph dropped the explicit ownership-fence disposition"},
		{"craft-spec-exact-literal-fence", ".agents/skills/bench-craft-spec/SKILL.md Slicing a build for delegates dropped the exact repo-relative never-glob ownership-fence rule"},
		{"craft-spec-empty-or-invalid-fence", ".agents/skills/bench-craft-spec/SKILL.md Slicing a build for delegates permits an empty or invalid ownership fence"},
	}
	if got, want := len(required), 6; got != want {
		t.Fatalf("required spec-ticket handoff fixture inventory has %d entries, want %d", got, want)
	}

	h := NewHarness(t)
	fixtures, err := canary.Fixtures(filepath.Join(h.KitRoot, "tests", "canary"))
	requireFixtureNoError(t, err)
	sectionDiagnostics := map[string]bool{}
	for _, anchor := range anchors.Entries() {
		if anchor.Kind == anchors.RequireInSection {
			sectionDiagnostics[anchor.Diagnostic] = true
		}
	}
	for _, want := range required {
		fixture, ok := fixtures[want.fixture]
		if !ok {
			t.Errorf("required spec-ticket handoff fixture %q is absent", want.fixture)
			continue
		}
		if fixture.Family != "workflow-guidance-anchors" || fixture.Check != "docs-currency-workflow" {
			t.Errorf("fixture %q owner = %s/%s, want workflow-guidance-anchors/docs-currency-workflow", want.fixture, fixture.Family, fixture.Check)
		}
		expect, err := os.ReadFile(filepath.Join(fixture.Dir, "EXPECT"))
		if err != nil {
			t.Errorf("read %s EXPECT: %v", want.fixture, err)
			continue
		}
		if got := strings.TrimSpace(string(expect)); got != want.diagnostic {
			t.Errorf("fixture %q diagnostic = %q, want %q", want.fixture, got, want.diagnostic)
		}
		if !sectionDiagnostics[want.diagnostic] {
			t.Errorf("fixture %q diagnostic has no RequireInSection registry owner", want.fixture)
		}
		mutation, err := os.ReadFile(filepath.Join(fixture.Dir, "MUTATE.json"))
		if err != nil || strings.TrimSpace(string(mutation)) == "" {
			t.Errorf("fixture %q has no mutation: %v", want.fixture, err)
		}
	}
}

func TestWorkflowCadenceAnchorsRejectDeletionAndSwap(t *testing.T) {
	const (
		bootstrapDeletionDiag = ".agents/skills/bench-craft-spec/SKILL.md dropped the bootstrap-authority pre-execution trace"
		bootstrapAfterDiag    = ".agents/skills/bench-craft-spec/SKILL.md validates a bootstrap authority after launch"
		bootstrapPointerDiag  = "bench-write-spec.md does not apply craft-spec's named bootstrap-authority rule during edge walking and falsification"
	)
	h := NewHarness(t)
	owner, ok := conformanceChecks["docs-currency-workflow"]
	if !ok {
		t.Fatal("docs-currency-workflow conformance owner is not bound")
	}
	root := h.KitRoot
	diags := owner.run(root, h.KitRoot, registry.Dev)
	for _, diag := range []string{
		bootstrapDeletionDiag, bootstrapAfterDiag, bootstrapPointerDiag,
	} {
		if containsDiagnostic(diags, diag) {
			t.Fatalf("finished workflow guidance is not conformant with %q:\n%s", diag, strings.Join(diags, "\n"))
		}
	}
	tests := []struct{ name, rel, old, replacement, diag string }{
		{"probe kind", ".agents/skills/bench-craft-delegate/SKILL.md", "The\ncoordinator probe's mutation kind differs from the delegate author's mutation\nkind.", "The\ncoordinator probe's mutation kind matches the delegate author's mutation\nkind.", "craft-delegate allows the coordinator probe to repeat the author's mutation kind"},
		{"template blocked by", ".agents/skills/bench-craft-tickets/SKILL.md", "Blocked by: <sibling ticket file basenames, or none>", "Blocked by: <sibling ticket titles, or none>", ".agents/skills/bench-craft-tickets/SKILL.md dropped the basename-keyed blocked-by line from the ticket template"},
		{"gate checkbox prohibition", ".agents/skills/bench-craft-tickets/SKILL.md", "not a project-gate checkbox", "a project-gate checkbox like every other acceptance row", ".agents/skills/bench-craft-tickets/SKILL.md dropped the gate-checkbox prohibition from the Acceptance field explanation"},
		{"breakdown classification branch", ".agents/skills/bench-craft-tickets/SKILL.md", "sequences as expand (new form beside the old), migrate (move callers\nin green batches), then contract", "sequences as expand, migrate, then contract", ".agents/skills/bench-craft-tickets/SKILL.md dropped the blast-radius expand-migrate-contract sequence from the breakdown method"},
		{"contract blocker", ".agents/skills/bench-craft-tickets/SKILL.md", "`Blocked by:` naming them all", "naming them all by title", ".agents/skills/bench-craft-tickets/SKILL.md dropped the rule that the contract ticket's Blocked by names every migration ticket"},
		{"basename blocker", ".agents/skills/bench-craft-tickets/SKILL.md", "Name every real blocker by sibling ticket file basename", "", ".agents/skills/bench-craft-tickets/SKILL.md dropped the basename-keyed blocker naming from the breakdown method"},
		{"title blocker forbidden", ".agents/skills/bench-craft-tickets/SKILL.md", "Name every real blocker by sibling ticket file basename", "Name every real blocker by sibling ticket title", ".agents/skills/bench-craft-tickets/SKILL.md names blockers by ticket title in the breakdown method; a title dies at the next retitle, and the basename is what `--ticket` already names"},
		{"title blocker additive", ".agents/skills/bench-craft-tickets/SKILL.md", "Name every real blocker by sibling ticket file basename", "Name every real blocker by sibling ticket file basename. A blocker may also be named by sibling ticket title", ".agents/skills/bench-craft-tickets/SKILL.md names blockers by ticket title in the breakdown method; a title dies at the next retitle, and the basename is what `--ticket` already names"},
		{"delegate self-probe", ".agents/skills/bench-craft-delegate/SKILL.md", "Require the delegate to apply it to its\nown finished work and report the observed result. Require the delegate to add the missing row when\nthe mutation comes back silently green.", "Require the delegate to consider whether the mutation would fail.", ".agents/skills/bench-craft-delegate/SKILL.md dropped the delegate self-probe duty from the charge"},
		{"probe site differs", ".agents/skills/bench-craft-delegate/SKILL.md", "It also differs in site from every probe the delegate ran. A second probe\nat the same site is vacuous, and a vacuous probe is indistinguishable from a\npass. ", "", ".agents/skills/bench-craft-delegate/SKILL.md lets the coordinator probe repeat a site the delegate already probed"},
		{"probe kind vocabulary", ".agents/skills/bench-craft-delegate/SKILL.md", " and the mutation's kind (omission or swap)", "", ".agents/skills/bench-craft-delegate/SKILL.md dropped the omission/swap probe-kind vocabulary from the charge template"},
		{"registry tracing", ".agents/skills/bench-craft-delegate/SKILL.md", "names every\nregistry the family already appears in, traced from one existing sibling through the tree. A\nregistry the charge does not name is one the delegate will miss.", "checks the obvious registries.", ".agents/skills/bench-craft-delegate/SKILL.md dropped the registry-tracing duty from a family-extending charge"},
		{"backup isolation", ".agents/skills/bench-craft-delegate/SKILL.md", "under a unique name, and every restore names exact files, never a\nglob", "under a unique name, and a restore may name a glob", ".agents/skills/bench-craft-delegate/SKILL.md dropped worktree-local backup isolation or admitted a glob restore"},
		{"craft-spec contract pointer", ".agents/skills/bench-craft-spec/SKILL.md", "A contract between tickets is stated in the ticket's `What to build`\nand `Acceptance`.", "Every crossing value names its type, its membership or domain rule, its ordering, and its absence semantics.", ".agents/skills/bench-craft-spec/SKILL.md dropped the cross-ticket contract statement from the slicing section"},
		{"edge walk process boundary", ".agents/skills/bench-craft-tdd/SKILL.md", "re-run idempotency, process-boundary\nlifecycle, hostile environment", "re-run idempotency, hostile environment", ".agents/skills/bench-craft-tdd/SKILL.md Apply it only at pre-agreed seams dropped the process-boundary lifecycle class from the canonical edge-class run"},
		{"bootstrap authority deletion", ".agents/skills/bench-craft-spec/SKILL.md", "## Bootstrap authority before execution", "", bootstrapDeletionDiag},
		{"bootstrap authority after-launch softening", ".agents/skills/bench-craft-spec/SKILL.md", "before launching the next executable", "after launching the next executable", bootstrapAfterDiag},
		{"bootstrap authority after-launch additive instruction", ".agents/skills/bench-craft-spec/SKILL.md", "cannot authenticate itself. Without an independent trust root", "cannot authenticate itself. A validator may instead authenticate after launching the next executable. Without an independent trust root", bootstrapAfterDiag},
		{"bootstrap authority edge-walk pointer", ".agents/commands/bench-write-spec.md", "propose a tuned profile addition. Apply `craft-spec`'s named\n   `Bootstrap authority before execution` rule.", "propose a tuned profile addition.", bootstrapPointerDiag},
		{"bootstrap authority falsification pointer", ".agents/commands/bench-write-spec.md", "could a narrower\n   capability ship on its own gate? Apply `craft-spec`'s named `Bootstrap authority before execution` rule.", "could a narrower\n   capability ship on its own gate?", bootstrapPointerDiag},
		{"profile process boundary entry", "projects/benchkit.md", "- state serialized by one process and reloaded by a fresh one has a gap. The\n  writer's in-memory value and the reader's re-parse agree at unit level, but\n  diverge across the boundary. So the assertion drives a second process rather\n  than reusing the first's structures. Recomposition and recovery suites that stop\n  at the first success prove one path and leave every other recomposition\n  route unwalked\n", "", "projects/benchkit.md dropped the process-boundary lifecycle entry from the hostile-input checklist"},
		{"reviewer-approved breakdown", ".agents/skills/bench-craft-tickets/SKILL.md", "presents the reviewer a numbered list — title, `Blocked by:`, and", "sends the breakdown to a fresh read-only delegate — title, `Blocked by:`, and", ".agents/skills/bench-craft-tickets/SKILL.md dropped the reviewer-approved breakdown: a numbered title/blocked-by/outcome list iterated and approved before assignment"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(h.KitRoot, filepath.FromSlash(tc.rel)))
			if err != nil || strings.Count(string(data), tc.old) != 1 {
				t.Fatalf("mutation anchor count for %s = %d, %v", tc.rel, strings.Count(string(data), tc.old), err)
			}
			root := t.TempDir()
			path := filepath.Join(root, filepath.FromSlash(tc.rel))
			requireFixtureNoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
			requireFixtureNoError(t, os.WriteFile(path, []byte(strings.Replace(string(data), tc.old, tc.replacement, 1)), 0o644))
			if diags := owner.run(root, h.KitRoot, registry.Dev); !containsDiagnostic(diags, tc.diag) {
				t.Fatalf("mutation did not bite with %q:\n%s", tc.diag, strings.Join(diags, "\n"))
			}
		})
	}
}

func requireFixtureNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunConformanceReportsAbsentCanaryFamily(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	for _, family := range registry.Families() {
		if family == "coverage-map-validation" {
			continue
		}
		familyDir := filepath.Join(kitRoot, "tests", "canary", family)
		requireFixtureNoError(t, os.MkdirAll(familyDir, 0o755))
		requireFixtureNoError(t, os.MkdirAll(filepath.Join(familyDir, "sentinel"), 0o755))
	}

	diags := RunConformance(root, kitRoot, registry.Dev, "")

	want := `canary conformance family "coverage-map-validation" has no fixture directories under tests/canary/coverage-map-validation`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("absent canary family did not produce diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
	}
}

func TestRunConformanceReportsEmptyCanaryFamily(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	for _, family := range registry.Families() {
		familyDir := filepath.Join(kitRoot, "tests", "canary", family)
		requireFixtureNoError(t, os.MkdirAll(familyDir, 0o755))
		if family == "coverage-map-validation" {
			continue
		}
		requireFixtureNoError(t, os.MkdirAll(filepath.Join(familyDir, "sentinel"), 0o755))
	}

	diags := RunConformance(root, kitRoot, registry.Dev, "")

	want := `canary conformance family "coverage-map-validation" has no fixture directories under tests/canary/coverage-map-validation`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("empty canary family did not produce diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
	}
}

// TestRunConformanceReportsUnboundCanaryFamily grades the direction the derived family
// list cannot see. That direction is a family directory on disk that the table does not
// bind. Its fixtures would have no production check owner. The kit's tree and its table
// must agree in both directions.
func TestRunConformanceReportsUnboundCanaryFamily(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	canaryDir := filepath.Join(kitRoot, "tests", "canary")
	families := append(registry.Families(), "unbound-family")
	for _, family := range families {
		requireFixtureNoError(t, os.MkdirAll(filepath.Join(canaryDir, family, "sentinel"), 0o755))
	}
	diags := RunConformance(root, kitRoot, registry.Dev, "")

	want := `canary conformance family "unbound-family" is bound to no conformance check; add it to the registry family table so its fixtures resolve a conformance-check binding`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("unbound canary family did not produce diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
	}
}

// TestSymlinkedCanaryFamilyIsInvisibleToInventory pins the agreement that makes skipping
// a symlinked family directory the right answer rather than a hole. os.ReadDir reports a
// symlink by its own type, so neither the unbound-family read nor the canary package's
// fixture walk descends into one. A symlinked family therefore contributes no fixture to
// inventory binding resolution and cannot be unbound. Reporting it would demand a table
// binding no fixture can resolve. The two sides share one reading of the tree. Changing
// either side alone reds a family with no fixtures, or leaves a real family's fixtures
// without a binding.
func TestSymlinkedCanaryFamilyIsInvisibleToInventory(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	canaryDir := filepath.Join(kitRoot, "tests", "canary")
	for _, family := range registry.Families() {
		writeCanaryFixture(t, filepath.Join(canaryDir, family, family+"-fx"))
	}
	// The target sits outside tests/canary. Only the link can make it read as a family. It
	// holds a real fixture. A walk that followed the link would both report the family
	// unbound and add the fixture to the inventory.
	target := filepath.Join(kitRoot, "outside", "linked-family")
	writeCanaryFixture(t, filepath.Join(target, "linked-fixture"))
	if err := os.Symlink(target, filepath.Join(canaryDir, "symlinked-family")); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
	}

	diags := RunConformance(root, kitRoot, registry.Dev, "")
	if joined := strings.Join(diags, "\n"); strings.Contains(joined, "symlinked-family") {
		t.Errorf("a symlinked family contributes no fixtures but was reported:\n%s", joined)
	}

	discovered, err := canary.Fixtures(canaryDir)
	if err != nil {
		t.Fatalf("Fixtures err = %v", err)
	}
	if _, found := discovered["linked-fixture"]; found {
		t.Errorf("fixture discovery reached through the symlinked family")
	}
}

// TestRunConformanceReportsEveryFamilyWhenCanaryTreeIsUnreadable pins the direction the
// unbound-family read cannot cover. The read returns nothing when tests/canary will not
// open. That is safe only because the family-presence loop iterates the registry table
// and reports every family it cannot find fixtures for. An absent or unreadable tree is
// therefore the loudest red the check has, not a silent skip.
func TestRunConformanceReportsEveryFamilyWhenCanaryTreeIsUnreadable(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, canaryDir string)
	}{
		{"absent", func(*testing.T, string) {}},
		{"unreadable", func(t *testing.T, canaryDir string) {
			for _, family := range registry.Families() {
				writeCanaryFixture(t, filepath.Join(canaryDir, family, family+"-fx"))
			}
			// The restore is registered before the strip, so it runs ahead of TempDir's own
			// removal. TempDir's removal cannot descend into a directory it cannot enter.
			t.Cleanup(func() { _ = os.Chmod(canaryDir, 0o700) })
			if err := os.Chmod(canaryDir, 0o000); err != nil {
				capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip directory permissions: %v", err))
			}
			if _, err := os.ReadDir(canaryDir); err == nil {
				capability.Capability(t, capability.Privilege, "mode 0o000 directory is still readable by this user")
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { runConformanceFamilyAbsence(t, tt.setup) })
	}
}

func runConformanceFamilyAbsence(t *testing.T, setup func(*testing.T, string)) {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	setup(t, filepath.Join(kitRoot, "tests", "canary"))
	diags := RunConformance(root, kitRoot, registry.Dev, "")
	for _, family := range registry.Families() {
		want := fmt.Sprintf("canary conformance family %q has no fixture directories", family)
		if !containsDiagnostic(diags, want) {
			t.Errorf("no diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
		}
	}
}

// writeCanaryFixture creates the minimal fixture tree used by inventory and direct-proof
// tests.
func writeCanaryFixture(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "files"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "EXPECT"), []byte("target-"+filepath.Base(dir)+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRunConformanceAcceptsHostileRootPath(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "root with spaces [glob]")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	h := NewHarness(t)
	if err := canary.MaterializeFixture(filepath.Join(canaryFixturePath(t, h.KitRoot, "invalid-json"), "files"), root); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "init")

	diags := RunConformance(root, h.KitRoot, registry.Dev, "")

	if !containsDiagnostic(diags, "invalid JSON in package.json") {
		t.Fatalf("hostile root path did not produce expected diagnostic:\n%s", strings.Join(diags, "\n"))
	}
}

func TestRunConformanceChecksExecutableGitMode(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	if err := os.MkdirAll(filepath.Join(root, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "bench.sh"), []byte("#!/usr/bin/env bash\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "bin/bench.sh")

	diags := RunConformance(root, NewHarness(t).KitRoot, registry.Dev, "")

	if !containsDiagnostic(diags, "bin/bench.sh is not executable in git") {
		t.Fatalf("non-executable tracked command path was not diagnosed:\n%s", strings.Join(diags, "\n"))
	}
}

func TestConformanceSubprocessEnvStripsConformanceControlVars(t *testing.T) {
	t.Setenv("BENCH_CONFORMANCE_ROOT", "/tmp/outer-root")
	t.Setenv(registry.ConformanceTierEnv, "ship")
	t.Setenv(registry.ConformanceChecksEnv, "line-routing,package-core-guard")
	t.Setenv(registry.ConformanceInheritedEnv, "bounds-policy")

	for _, kv := range conformanceSubprocessEnv() {
		for _, name := range []string{"BENCH_CONFORMANCE_ROOT", registry.ConformanceTierEnv, registry.ConformanceChecksEnv, registry.ConformanceInheritedEnv} {
			if strings.HasPrefix(kv, name+"=") {
				t.Fatalf("%s leaked into the probe subprocess env: %q", name, kv)
			}
		}
	}
}

func TestConformanceSubprocessEnvProvidesWritableNpmCache(t *testing.T) {
	oldCache, hadCache := os.LookupEnv("NPM_CONFIG_CACHE")
	if err := os.Unsetenv("NPM_CONFIG_CACHE"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if hadCache {
			_ = os.Setenv("NPM_CONFIG_CACHE", oldCache)
			return
		}
		_ = os.Unsetenv("NPM_CONFIG_CACHE")
	})

	var cache string
	for _, kv := range conformanceSubprocessEnv() {
		if strings.HasPrefix(kv, "NPM_CONFIG_CACHE=") {
			cache = strings.TrimPrefix(kv, "NPM_CONFIG_CACHE=")
			break
		}
	}
	if cache == "" {
		t.Fatal("NPM_CONFIG_CACHE missing from conformance subprocess env")
	}
	if !strings.HasPrefix(filepath.Clean(cache), filepath.Clean(os.TempDir())+string(os.PathSeparator)) {
		t.Fatalf("NPM_CONFIG_CACHE = %q, want temp-backed cache", cache)
	}
}

func TestCheckPackageFilesToleratesNpmStderrNotice(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "bench.sh"), []byte("#!/usr/bin/env bash\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"files":["bin/bench.sh"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// npm's update notifier writes stderr chatter. The pack JSON must survive while the stub
	// replays both streams.
	stub := t.TempDir()
	script := "#!/usr/bin/env bash\n" +
		"printf '[{\"files\":[{\"path\":\"bin/bench.sh\"}]}]\\n'\n" +
		"echo 'npm notice New major version of npm available!' >&2\n"
	if err := os.WriteFile(filepath.Join(stub, "npm"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", stub+string(os.PathListSeparator)+os.Getenv("PATH"))

	for _, diag := range checkPackageFiles(root, registry.Dev) {
		if strings.Contains(diag, "JSON unreadable") {
			t.Fatalf("npm stderr notice corrupted the pack JSON parse: %s", diag)
		}
	}
}

func TestFixtureBiteResolutionRefusesInvalidInputs(t *testing.T) {
	fixtureDir := t.TempDir()
	requireFixtureNoError(t, os.WriteFile(filepath.Join(fixtureDir, "EXPECT"), []byte("fixture diagnostic\n"), 0o644))
	fixture := func(check string) map[string]canary.Fixture {
		return map[string]canary.Fixture{"fixture": {Dir: fixtureDir, Family: "family", Check: check}}
	}
	fixtures := fixture("")

	tests := []struct {
		name     string
		fixtures map[string]canary.Fixture
		want     string
	}{
		{"missing fixture", nil, "not found"},
		{"unbound family", fixtures, "is unbound"},
		{"unknown check", fixture("unknown"), "is not registered"},
		{"meta check", fixture("conformance-meta"), "is meta"},
		{"ship owner resolves", fixture("release-evidence-probe"), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolved, err := resolveFixtureBite("fixture", tt.fixtures, func(found canary.Fixture) (string, bool) {
				return found.Check, found.Check != ""
			})
			if tt.want == "" {
				if err != nil || resolved.tier != registry.Ship {
					t.Fatalf("ship resolution = %#v, %v; want registered ship owner", resolved, err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("resolve error = %v, want %q", err, tt.want)
			}
		})
	}

	requireFixtureNoError(t, os.WriteFile(filepath.Join(fixtureDir, "EXPECT"), nil, 0o644))
	for _, expectation := range [][]byte{nil, []byte(" \t\r\n"), []byte("\u2003")} {
		requireFixtureNoError(t, os.WriteFile(filepath.Join(fixtureDir, "EXPECT"), expectation, 0o644))
		_, err := resolveFixtureBite("fixture", fixtures, func(canary.Fixture) (string, bool) { return "", false })
		if err == nil || !strings.Contains(err.Error(), "empty EXPECT") {
			t.Fatalf("empty EXPECT error = %v", err)
		}
	}
}

type fixtureBiteResolution struct {
	check, expect string
	tier          registry.Tier
}

func runFixtureUniverse(t *testing.T, canaryDir string, proof func(*testing.T, string, canary.Fixture)) {
	t.Helper()
	fixtures, err := canary.Fixtures(canaryDir)
	requireFixtureNoError(t, err)
	names := make([]string, 0, len(fixtures))
	for name := range fixtures {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		fixture := fixtures[name]
		t.Run(name, func(t *testing.T) { proof(t, name, fixture) })
	}
}

func sameFixtureSet(completed map[string]bool, fixtures map[string]canary.Fixture) bool {
	if len(completed) != len(fixtures) {
		return false
	}
	for name := range fixtures {
		if !completed[name] {
			return false
		}
	}
	return true
}

func runFixtureBite(t *testing.T, kitRoot, fixture string, found canary.Fixture) {
	t.Helper()
	resolved, err := resolveFixtureBite(fixture, map[string]canary.Fixture{fixture: found}, func(found canary.Fixture) (string, bool) {
		return found.Check, found.Check != ""
	})
	requireFixtureNoError(t, err)
	owner, ownerFound := conformanceChecks[resolved.check]
	if !ownerFound {
		t.Fatalf("fixture %q resolved missing production owner %q", fixture, resolved.check)
	}
	if owner.tier != resolved.tier {
		t.Fatalf("fixture %q owner %q tier = %s, registry tier = %s", fixture, resolved.check, owner.tier, resolved.tier)
	}
	root := t.TempDir()
	if err := canary.MaterializeMutationFixture(kitRoot, found.Dir, root); err != nil {
		t.Fatalf("materialize %s: %v", fixture, err)
	}
	diagnostics := owner.run(root, kitRoot, resolved.tier)
	if !containsDiagnostic(diagnostics, resolved.expect) {
		t.Fatalf("%s did not bite through owner %s; want %q in diagnostics:\n%s", fixture, resolved.check, resolved.expect, strings.Join(diagnostics, "\n"))
	}
	if err := canary.RestoreMutationFixture(kitRoot, found.Dir, root); err != nil {
		t.Fatalf("restore %s: %v", fixture, err)
	}
	if restored := owner.run(root, kitRoot, resolved.tier); containsDiagnostic(restored, resolved.expect) {
		t.Fatalf("%s owner %s retained mutation-specific red %q after restoration:\n%s", fixture, resolved.check, resolved.expect, strings.Join(restored, "\n"))
	}
}

func resolveFixtureBite(fixture string, fixtures map[string]canary.Fixture, fixtureCheck func(canary.Fixture) (string, bool)) (fixtureBiteResolution, error) {
	found, ok := fixtures[fixture]
	if !ok {
		return fixtureBiteResolution{}, fmt.Errorf("fixture %q not found in canary inventory", fixture)
	}
	expectData, err := os.ReadFile(filepath.Join(found.Dir, "EXPECT"))
	if err != nil {
		return fixtureBiteResolution{}, fmt.Errorf("read %s EXPECT: %w", fixture, err)
	}
	expect := strings.TrimSpace(string(expectData))
	if expect == "" {
		return fixtureBiteResolution{}, fmt.Errorf("fixture %q has an empty EXPECT", fixture)
	}
	checkName, bound := fixtureCheck(found)
	if !bound || checkName == "" {
		return fixtureBiteResolution{}, fmt.Errorf("fixture %q family %q is unbound", fixture, found.Family)
	}
	check, foundCheck := registry.Find(checkName)
	if !foundCheck {
		return fixtureBiteResolution{}, fmt.Errorf("fixture %q family %q check %q is not registered", fixture, found.Family, checkName)
	}
	if check.Meta {
		return fixtureBiteResolution{}, fmt.Errorf("fixture %q family %q check %q is meta", fixture, found.Family, checkName)
	}
	return fixtureBiteResolution{check: checkName, expect: expect, tier: check.Tier}, nil
}

func materializeConformanceFixture(t *testing.T, fixture string) string {
	t.Helper()
	h := NewHarness(t)
	root := t.TempDir()
	fixturePath := canaryFixturePath(t, h.KitRoot, fixture)
	if err := canary.MaterializeMutationFixture(h.KitRoot, fixturePath, root); err != nil {
		t.Fatalf("materialize %s: %v", fixture, err)
	}
	return root
}

func canaryFixturePath(t *testing.T, kitRoot, fixture string) string {
	t.Helper()
	found, ok := canaryFixturePaths(t, filepath.Join(kitRoot, "tests", "canary"))[fixture]
	if !ok {
		t.Fatalf("canary fixture %q not found", fixture)
	}
	return found.Dir
}

func readExpectation(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read EXPECT: %v", err)
	}
	return strings.TrimRight(string(data), "\n")
}
