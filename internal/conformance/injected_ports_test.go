package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// The injected-interface composition audit found the same gap at five seams at once: a
// port injected into a consumer whose tests only ever drive a fake, so softening the real
// producer left the gate green. Junction tests close the seams that were open; this check
// is what keeps them closed, by refusing a derived port that names neither a test driving
// its real producer nor a written reason it does not.
//
// The two halves are deliberately independent. The inventory is *derived* from the audited
// packages' Go source, so a port added tomorrow appears without anyone remembering to say
// so; the registry below is the *advertisement*, and disagreement between them is the red.
// A single hand-maintained list would satisfy itself.
//
// The named-test-exists half is a **tripwire**, not a behavior check: it catches a junction
// test deleted or renamed, and it cannot catch one gutted back to driving a fake. The
// behavior half of the defense is the junction tests themselves and the canary fixture
// under tests/canary/injected-ports, which keeps the unregistered-port red alive.

// auditedPortPackages are the packages whose injected ports the registry covers — the
// audit's inventory, no wider. A package outside it declaring a port is out of scope until
// a reviewer binds it in; a package inside it that derives nothing fails closed.
var auditedPortPackages = []string{
	"internal/canary",
	"internal/gitguard",
	"internal/publication",
}

// The five failure modes, each with its own message so a red names its cause without
// archaeology. They are constants because the canary fixture's EXPECT and the bite proof
// both quote them, and a message typed a second time drifts from the one that fires.
const (
	unregisteredPortMessage   = "injected port has no registry row"
	missingPortTestMessage    = "injected port registry row names a test the tree does not declare"
	emptyPortExemptionMessage = "injected port registry row is exempt with an empty reason"
	zeroPortInventoryMessage  = "injected port derivation found no ports"
	orphanPortRowMessage      = "injected port registry row names a port the derivation no longer reports"
)

// injectedPortRow is one derived port's disposition: either a real-producer test — a test
// that drives the actual producer through the consuming surface — or an exemption whose
// reason someone had to write down. A row carrying neither is the empty-exemption red.
type injectedPortRow struct {
	pkg, port          string
	testFile, testName string
	exempt             string
}

// injectedPortRegistry is the advertisement half: one row per port the derivation finds in
// auditedPortPackages. Rows are grouped by package and ordered as the derivation reports
// them, so a diff against a re-derived inventory reads straight down.
var injectedPortRegistry = []injectedPortRow{
	{
		pkg: "internal/canary", port: "Runner",
		testFile: "internal/canary/runner_junction_test.go", testName: "TestSweepTierRunsPlantedBashGate",
	},
	{
		pkg: "internal/gitguard", port: "Checker",
		testFile: "internal/gitguard/checker_junction_test.go", testName: "TestClassifyRealCheckerResolvedComposition",
	},
	{
		pkg: "internal/publication", port: "Registry",
		exempt: "the only adapter without gate coverage is NPMCLIRegistry, which is runbook-only: the gate drives FixtureRegistry against the hermetic offline registry, and no NPMCLIRegistry path performs gate egress or touches a credential (internal/publication/registry.go:5-11)",
	},
}

// derivedPort is one port the source sweep found, carrying the site that injected it so a
// red says where the port entered the package rather than only that it did.
type derivedPort struct {
	pkg, name, site string
}

// checkInjectedPortRegistry grades the registry against a fresh derivation over root.
func checkInjectedPortRegistry(root string) []string {
	return injectedPortDiags(root, auditedPortPackages, injectedPortRegistry)
}

// injectedPortDiags is the graded core, taking the package list and the rows as arguments so
// the bite proof can drive a mutated registry without editing the real one.
//
// A package the graded tree does not hold is not graded at all: an adopting repo has none of
// these directories, and a check that reds there is a check nobody keeps. A package it does
// hold has to yield ports, which is the fail-closed posture — see zeroPortInventoryMessage.
//
// The registry→derived direction is graded here too, not only derived→registry: every row of
// a package whose derivation is non-zero must name a port that derivation still reports, or it
// fires orphanPortRowMessage. Without this half, deleting an injection arm together with its
// junction test leaves a green row behind for as long as the package's other ports keep the
// inventory non-zero — the same silence zeroPortInventoryMessage closes for a package that
// derives nothing at all.
func injectedPortDiags(root string, packages []string, rows []injectedPortRow) []string {
	registered := make(map[string]injectedPortRow, len(rows))
	for _, row := range rows {
		registered[row.pkg+"."+row.port] = row
	}
	derived := make(map[string]bool, len(rows))
	gradedNonZero := make(map[string]bool, len(packages))
	var diags []string
	for _, pkg := range packages {
		dir := filepath.Join(root, filepath.FromSlash(pkg))
		if !isDirectory(dir) {
			continue
		}
		ports, parseDiags := derivedInjectedPorts(dir, pkg)
		diags = append(diags, parseDiags...)
		if len(ports) == 0 {
			diags = append(diags, fmt.Sprintf(
				"%s: %s is present in the graded tree and the audit records it as declaring them, so an empty inventory would grade every registry row as satisfied — the derivation fails closed instead of reporting silence",
				zeroPortInventoryMessage, pkg))
			continue
		}
		gradedNonZero[pkg] = true
		for _, port := range ports {
			key := port.pkg + "." + port.name
			derived[key] = true
			row, found := registered[key]
			if !found {
				diags = append(diags, fmt.Sprintf(
					"%s: %s.%s is injected at %s — add a row to injectedPortRegistry naming the test that drives its real producer, or the written reason it has none; an unregistered port is the audit's gap recurring unseen",
					unregisteredPortMessage, port.pkg, port.name, port.site))
				continue
			}
			diags = append(diags, injectedPortRowDiags(root, row)...)
		}
	}
	for _, row := range rows {
		if !gradedNonZero[row.pkg] || derived[row.pkg+"."+row.port] {
			continue
		}
		diags = append(diags, fmt.Sprintf(
			"%s: %s.%s is registered but the derivation no longer reports it — a stale advertisement whose port or injection arm left the tree must not stay green",
			orphanPortRowMessage, row.pkg, row.port))
	}
	return uniqueSorted(diags)
}

// injectedPortRowDiags grades one row's disposition.
func injectedPortRowDiags(root string, row injectedPortRow) []string {
	if row.testName == "" {
		if strings.TrimSpace(row.exempt) == "" {
			return []string{fmt.Sprintf(
				"%s: %s.%s names no real-producer test and gives no reason — an exemption nobody had to write down is an unregistered port with better manners",
				emptyPortExemptionMessage, row.pkg, row.port)}
		}
		return nil
	}
	if !fileDeclaresTest(filepath.Join(root, filepath.FromSlash(row.testFile)), row.testName) {
		return []string{fmt.Sprintf(
			"%s: %s.%s advertises %s:%s, which that file does not declare — a deleted or renamed junction test must not leave a green advertisement behind (tripwire: this catches removal, not a test gutted back to a fake)",
			missingPortTestMessage, row.pkg, row.port, row.testFile, row.testName)}
	}
	return nil
}

// derivedInjectedPorts reports the injected ports one package's non-test source declares.
//
// A **port** is a named type whose shape is port-shaped — a non-empty interface, a func
// type, or a struct whose every field is a func — and which the package injects: passes as
// a function parameter, or narrows to by type assertion or type switch. The assertion arm
// is not decoration: a capability widened out of an already-injected owner is never a
// parameter at all, so a parameter-only rule would leave exactly the silently-downgradable
// capabilities out of the inventory. Exportedness is not part of the rule either — a
// package main declares nothing exported, and its ports are ports all the same.
func derivedInjectedPorts(dir, pkg string) ([]derivedPort, []string) {
	fset := token.NewFileSet()
	shaped := map[string]bool{}
	sites := map[string]string{}
	var diags []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, []string{fmt.Sprintf("injected port derivation cannot read %s: %v", pkg, err)}
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		path := filepath.Join(dir, name)
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			diags = append(diags, fmt.Sprintf("injected port derivation cannot parse %s/%s: %v", pkg, name, err))
			continue
		}
		rel := pkg + "/" + name
		ast.Inspect(file, func(node ast.Node) bool {
			switch typed := node.(type) {
			case *ast.TypeSpec:
				if portShapedType(typed.Type) {
					shaped[typed.Name.Name] = true
				}
			case *ast.FuncDecl:
				if typed.Type.Params == nil {
					return true
				}
				for _, param := range typed.Type.Params.List {
					recordInjectionSite(sites, param.Type, fmt.Sprintf("a parameter of %s in %s:%d", typed.Name.Name, rel, fset.Position(param.Pos()).Line))
				}
			case *ast.TypeAssertExpr:
				recordInjectionSite(sites, typed.Type, fmt.Sprintf("a type assertion in %s:%d", rel, fset.Position(typed.Pos()).Line))
			case *ast.CaseClause:
				for _, expr := range typed.List {
					recordInjectionSite(sites, expr, fmt.Sprintf("a type switch case in %s:%d", rel, fset.Position(expr.Pos()).Line))
				}
			}
			return true
		})
	}

	var ports []derivedPort
	for name := range shaped {
		if site, injected := sites[name]; injected {
			ports = append(ports, derivedPort{pkg: pkg, name: name, site: site})
		}
	}
	sort.Slice(ports, func(i, j int) bool { return ports[i].name < ports[j].name })
	return ports, diags
}

// recordInjectionSite keeps the first site a name is injected at; a port injected twice is
// still one port, and the first site is the stable one to name in a diagnostic.
func recordInjectionSite(sites map[string]string, expr ast.Expr, site string) {
	name := typeIdentName(expr)
	if name == "" {
		return
	}
	if _, seen := sites[name]; !seen {
		sites[name] = site
	}
}

// typeIdentName reduces a type expression to the package-local name it names, seeing
// through a pointer. A qualified type belongs to another package and is not this package's
// port to register.
func typeIdentName(expr ast.Expr) string {
	switch typed := expr.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.StarExpr:
		return typeIdentName(typed.X)
	}
	return ""
}

// portShapedType reports whether a type declaration has the shape of a port. An empty
// interface constrains nothing and names no producer, so it is not one.
func portShapedType(expr ast.Expr) bool {
	switch typed := expr.(type) {
	case *ast.InterfaceType:
		return typed.Methods != nil && len(typed.Methods.List) > 0
	case *ast.FuncType:
		return true
	case *ast.StructType:
		if typed.Fields == nil || len(typed.Fields.List) == 0 {
			return false
		}
		for _, field := range typed.Fields.List {
			if _, isFunc := field.Type.(*ast.FuncType); !isFunc {
				return false
			}
		}
		return true
	}
	return false
}

// fileDeclaresTest reports whether path declares a top-level func of the given name. It
// parses rather than searches so the name written in a comment or a string is not a match —
// the tripwire has to answer for the test the binary would run.
func fileDeclaresTest(path, name string) bool {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return false
	}
	for _, decl := range file.Decls {
		fn, isFunc := decl.(*ast.FuncDecl)
		if isFunc && fn.Recv == nil && fn.Name.Name == name {
			return true
		}
	}
	return false
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// TestInjectedPortRegistryCheckBites is the recorded bite proof for
// checkInjectedPortRegistry (per craft-gate). The unregistered-port red is retained
// permanently by the tests/canary/injected-ports fixture; the other four modes have no
// fixture of their own and are proved here, each against a synthetic tree.
func TestInjectedPortRegistryCheckBites(t *testing.T) {
	// One package, one port-shaped type, one injection site — the smallest tree the
	// derivation reports anything at all for.
	const portSource = "package sample\n\n// Port is the fixture's injected port.\ntype Port interface {\n\tDo() error\n}\n\nfunc use(p Port) error { return p.Do() }\n"
	const testSource = "package sample\n\nimport \"testing\"\n\nfunc TestRealProducer(t *testing.T) {}\n"

	write := func(t *testing.T, files map[string]string) string {
		t.Helper()
		root := t.TempDir()
		for rel, body := range files {
			path := filepath.Join(root, filepath.FromSlash(rel))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return root
	}
	sampleTree := map[string]string{
		"internal/sample/sample.go":      portSource,
		"internal/sample/sample_test.go": testSource,
	}
	honestRow := injectedPortRow{
		pkg: "internal/sample", port: "Port",
		testFile: "internal/sample/sample_test.go", testName: "TestRealProducer",
	}
	packages := []string{"internal/sample"}

	t.Run("an honest row is silent", func(t *testing.T) {
		diags := injectedPortDiags(write(t, sampleTree), packages, []injectedPortRow{honestRow})
		if len(diags) != 0 {
			t.Fatalf("a derived port with a real-producer test the tree declares: want no diagnostics, got %v", diags)
		}
	})

	t.Run("a package the tree does not hold is not graded", func(t *testing.T) {
		if diags := injectedPortDiags(t.TempDir(), packages, nil); len(diags) != 0 {
			t.Fatalf("an absent audited package: want no diagnostics, got %v", diags)
		}
	})

	t.Run("an unregistered port fires", func(t *testing.T) {
		diags := injectedPortDiags(write(t, sampleTree), packages, nil)
		if len(diags) != 1 || !strings.Contains(diags[0], unregisteredPortMessage) ||
			!strings.Contains(diags[0], "internal/sample.Port") || !strings.Contains(diags[0], "a parameter of use in internal/sample/sample.go") {
			t.Fatalf("a derived port with no registry row: want one unregistered diagnostic naming it and its site, got %v", diags)
		}
	})

	t.Run("a row naming a test the tree lacks fires", func(t *testing.T) {
		stale := honestRow
		stale.testName = "TestRealProducerRenamedAway"
		diags := injectedPortDiags(write(t, sampleTree), packages, []injectedPortRow{stale})
		if len(diags) != 1 || !strings.Contains(diags[0], missingPortTestMessage) ||
			!strings.Contains(diags[0], "TestRealProducerRenamedAway") {
			t.Fatalf("a row advertising an absent test: want one missing-test diagnostic naming it, got %v", diags)
		}
	})

	t.Run("a test named only in a comment does not satisfy the tripwire", func(t *testing.T) {
		commented := map[string]string{
			"internal/sample/sample.go":      portSource,
			"internal/sample/sample_test.go": "package sample\n\n// TestRealProducer used to live here.\nconst gone = \"func TestRealProducer(t *testing.T) {}\"\n",
		}
		diags := injectedPortDiags(write(t, commented), packages, []injectedPortRow{honestRow})
		if len(diags) != 1 || !strings.Contains(diags[0], missingPortTestMessage) {
			t.Fatalf("a test surviving only as prose: want one missing-test diagnostic, got %v", diags)
		}
	})

	t.Run("an exemption with an empty reason fires", func(t *testing.T) {
		blank := injectedPortRow{pkg: "internal/sample", port: "Port", exempt: "   "}
		diags := injectedPortDiags(write(t, sampleTree), packages, []injectedPortRow{blank})
		if len(diags) != 1 || !strings.Contains(diags[0], emptyPortExemptionMessage) ||
			!strings.Contains(diags[0], "internal/sample.Port") {
			t.Fatalf("an exemption with a blank reason: want one empty-exemption diagnostic, got %v", diags)
		}
	})

	t.Run("a written exemption is silent", func(t *testing.T) {
		reasoned := injectedPortRow{pkg: "internal/sample", port: "Port", exempt: "covered by the runtime contract suite"}
		if diags := injectedPortDiags(write(t, sampleTree), packages, []injectedPortRow{reasoned}); len(diags) != 0 {
			t.Fatalf("an exemption carrying a reason: want no diagnostics, got %v", diags)
		}
	})

	t.Run("a package present but deriving nothing fails closed", func(t *testing.T) {
		empty := map[string]string{"internal/sample/sample.go": "package sample\n\n// Name carries no port.\nfunc Name() string { return \"sample\" }\n"}
		diags := injectedPortDiags(write(t, empty), packages, []injectedPortRow{honestRow})
		if len(diags) != 1 || !strings.Contains(diags[0], zeroPortInventoryMessage) ||
			!strings.Contains(diags[0], "internal/sample") {
			t.Fatalf("an audited package deriving zero ports: want one zero-inventory diagnostic naming it, got %v", diags)
		}
	})

	t.Run("unparsable audited source is reported", func(t *testing.T) {
		broken := map[string]string{"internal/sample/sample.go": "package sample\n\nfunc Port broken {\n"}
		diags := injectedPortDiags(write(t, broken), packages, nil)
		if len(diags) != 2 || !anyContains(diags, "cannot parse internal/sample/sample.go") ||
			!anyContains(diags, zeroPortInventoryMessage) {
			t.Fatalf("unparsable audited source: want a parse diagnostic and the fail-closed inventory red, got %v", diags)
		}
	})

	// twoPortTree carries two ports, PortA and PortB, both derived and honestly rowed;
	// the fixtures below pair them with orphanRow, a registry row naming PortC, a port
	// the tree never declares. The two honest ports sit beside the orphan so the check
	// must fire only for the underived row, not collaterally for the derived ones.
	twoPortTree := map[string]string{
		"internal/sample2/sample2.go":      "package sample2\n\ntype PortA interface {\n\tDo() error\n}\n\ntype PortB interface {\n\tDo() error\n}\n\nfunc use(a PortA, b PortB) error {\n\tif err := a.Do(); err != nil {\n\t\treturn err\n\t}\n\treturn b.Do()\n}\n",
		"internal/sample2/sample2_test.go": "package sample2\n\nimport \"testing\"\n\nfunc TestRealProducerA(t *testing.T) {}\nfunc TestRealProducerB(t *testing.T) {}\n",
	}
	twoPortPackages := []string{"internal/sample2"}
	portARow := injectedPortRow{
		pkg: "internal/sample2", port: "PortA",
		testFile: "internal/sample2/sample2_test.go", testName: "TestRealProducerA",
	}
	portBRow := injectedPortRow{
		pkg: "internal/sample2", port: "PortB",
		testFile: "internal/sample2/sample2_test.go", testName: "TestRealProducerB",
	}
	orphanRow := injectedPortRow{
		pkg: "internal/sample2", port: "PortC",
		testFile: "internal/sample2/sample2_test.go", testName: "TestRealProducerC",
	}

	t.Run("a row naming an underived port fires the orphan diagnostic alone", func(t *testing.T) {
		diags := injectedPortDiags(write(t, twoPortTree), twoPortPackages, []injectedPortRow{portARow, portBRow, orphanRow})
		if len(diags) != 1 || !strings.Contains(diags[0], orphanPortRowMessage) ||
			!strings.Contains(diags[0], "internal/sample2.PortC") {
			t.Fatalf("PortA and PortB derived and honestly rowed, PortC registered but never derived: want exactly one orphan-row diagnostic naming PortC, got %v", diags)
		}
	})

	t.Run("the orphan and missing-test diagnostics coexist without masking each other", func(t *testing.T) {
		// PortB's injection arm and its named test both survive here, but the row points
		// at a test name the tree does not declare, so it takes the missing-test path.
		// A row whose arm was deleted along with its test would instead fall to the
		// orphan path above. Both are exercised together to prove neither masks the other.
		staleTestRow := injectedPortRow{
			pkg: "internal/sample2", port: "PortB",
			testFile: "internal/sample2/sample2_test.go", testName: "TestRealProducerBRenamedAway",
		}
		diags := injectedPortDiags(write(t, twoPortTree), twoPortPackages, []injectedPortRow{portARow, staleTestRow, orphanRow})
		if len(diags) != 2 ||
			!anyContains(diags, missingPortTestMessage) || !anyContains(diags, "internal/sample2.PortB") ||
			!anyContains(diags, orphanPortRowMessage) || !anyContains(diags, "internal/sample2.PortC") {
			t.Fatalf("a derived port with a stale test name beside an orphaned row: want one missing-test diagnostic for PortB and one orphan-row diagnostic for PortC, got %v", diags)
		}
	})
}

// TestInjectedPortDerivationSeesEveryPortShape pins the shape and injection rules against a
// synthetic package carrying every shape the audit found — interface, func type, struct of
// func fields — reached by every injection route it found, including the type assertion and
// type switch that are the only way a widened owner capability is ever named. A rule that
// dropped one of those routes would silently shrink the inventory the
// registry is graded against, which is the same silence the fail-closed posture refuses.
func TestInjectedPortDerivationSeesEveryPortShape(t *testing.T) {
	source := strings.Join([]string{
		"package sample",
		"",
		"type Iface interface{ Do() error }",
		"type Widened interface{ Do() error; More() error }",
		"type Switched interface{ Do() error; Other() error }",
		"type FuncPort func(string) bool",
		"type FieldsPort struct{ A func() bool; B func() error }",
		"type Empty interface{}",
		"type Data struct{ Name string }",
		"type Mixed struct{ A func() bool; Name string }",
		"type Unused interface{ Do() error }",
		"",
		"func inject(i Iface, f FuncPort, s *FieldsPort, e Empty, d Data, m Mixed) {}",
		"func widen(i Iface) (Widened, bool) { w, ok := i.(Widened); return w, ok }",
		"func route(i Iface) string {",
		"\tswitch i.(type) {",
		"\tcase Switched:",
		"\t\treturn \"switched\"",
		"\t}",
		"\treturn \"\"",
		"}",
		"",
	}, "\n")
	root := t.TempDir()
	dir := filepath.Join(root, "internal", "sample")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sample.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}

	ports, diags := derivedInjectedPorts(dir, "internal/sample")
	if len(diags) != 0 {
		t.Fatalf("derivation reported %v", diags)
	}
	var names []string
	for _, port := range ports {
		names = append(names, port.name)
	}
	// Empty, Data, and Mixed are injected but not port-shaped; Unused is port-shaped but
	// never injected. Both exclusions are the rule, not an accident of the fixture.
	want := "FieldsPort,FuncPort,Iface,Switched,Widened"
	if strings.Join(names, ",") != want {
		t.Fatalf("derived %q, want %q", strings.Join(names, ","), want)
	}
}
