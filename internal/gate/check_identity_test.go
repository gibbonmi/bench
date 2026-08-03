package gate

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

func TestCurrentConformanceCanaryIdentitiesResolve(t *testing.T) {
	root, err := benchgit.Root()
	if err != nil {
		t.Fatal(err)
	}
	identities, err := resolveConformanceCanaryIdentities(root, registry.Dev)
	if err != nil {
		t.Fatal(err)
	}
	if !isContentAddress(identities.Shared) {
		t.Fatalf("shared canary identity = %q", identities.Shared)
	}
	for _, check := range registry.Checks {
		if !check.Meta && check.RunsAt(registry.Dev) && len(registry.CanaryFamilies(check.Name)) > 0 && !isContentAddress(identities.Bound[check.Name]) {
			t.Errorf("check %s has no bound canary implementation identity", check.Name)
		}
	}
}

func TestConformanceCheckIdentityBindsEveryAuthorityField(t *testing.T) {
	base := checkIdentityMaterial{
		Check: registry.Check{
			Name:           "docs-currency-workflow",
			Implementation: "checkDocsCurrencyAndWorkflow",
			Tier:           registry.Dev,
			Subject:        registry.SubjectRootAndKitRoot,
			Inputs:         registry.InputCatchAll,
		},
		Ordinal:        7,
		Inputs:         []treeEntry{{Path: "ROADMAP.md", Metadata: "100644 blob input"}},
		Implementation: []treeEntry{{Path: "internal/conformance/checks_test.go", Metadata: "100644 blob closure"}},
		CanaryFamilies: []string{"docs-currency-token-diet", "workflow-guidance-anchors"},
		CanaryInputs:   []treeEntry{{Path: "tests/canary/docs-currency-token-diet/fixture/EXPECT", Metadata: "100644 blob canary"}},
		Invocation:     []string{"BENCH_CONFORMANCE_TIER=dev", "BENCH_CONFORMANCE_CHECKS=docs-currency-workflow"},
	}
	want := conformanceCheckIdentity(base)
	mutations := map[string]func(*checkIdentityMaterial){
		"name":                   func(m *checkIdentityMaterial) { m.Check.Name += "-moved" },
		"tier":                   func(m *checkIdentityMaterial) { m.Check.Tier = registry.Ship },
		"registry order":         func(m *checkIdentityMaterial) { m.Ordinal++ },
		"function binding":       func(m *checkIdentityMaterial) { m.Check.Implementation += "Moved" },
		"subject binding":        func(m *checkIdentityMaterial) { m.Check.Subject = registry.SubjectRoot },
		"input declaration":      func(m *checkIdentityMaterial) { m.Check.Inputs = "another-source" },
		"declared input content": func(m *checkIdentityMaterial) { m.Inputs[0].Metadata += "-moved" },
		"shared implementation":  func(m *checkIdentityMaterial) { m.Implementation[0].Metadata += "-moved" },
		"canary ownership":       func(m *checkIdentityMaterial) { m.CanaryFamilies = append(m.CanaryFamilies, "another-family") },
		"canary content":         func(m *checkIdentityMaterial) { m.CanaryInputs[0].Metadata += "-moved" },
		"invocation schema":      func(m *checkIdentityMaterial) { m.Invocation[0] += "-moved" },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			moved := base
			moved.Inputs = append([]treeEntry(nil), base.Inputs...)
			moved.Implementation = append([]treeEntry(nil), base.Implementation...)
			moved.CanaryFamilies = append([]string(nil), base.CanaryFamilies...)
			moved.CanaryInputs = append([]treeEntry(nil), base.CanaryInputs...)
			moved.Invocation = append([]string(nil), base.Invocation...)
			mutate(&moved)
			if got := conformanceCheckIdentity(moved); got == want {
				t.Fatalf("identity did not move after %s changed", name)
			}
		})
	}
}

func TestDeclaredCheckSymlinkBindsCanonicalTargetAndRefusesHostileTargets(t *testing.T) {
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	writeGateTestFile(t, root, "scripts/gate-source.sh", "first\n", 0o644)
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, ".bench", "gate.sh")
	if err := os.Symlink("../scripts/gate-source.sh", link); err != nil {
		t.Fatal(err)
	}
	check := registry.Check{Name: "gate-entry-contract", Inputs: registry.InputGateEntry}
	beforeSnapshot := mustTreeSnapshot(t, root)
	before, err := resolveCheckInputs(root, check, beforeSnapshot)
	if err != nil {
		t.Fatal(err)
	}

	writeGateTestFile(t, root, "scripts/gate-source.sh", "second\n", 0o644)
	after, err := resolveCheckInputs(root, check, mustTreeSnapshot(t, root))
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(before, after) {
		t.Fatal("declared symlink inputs did not move when only the canonical target content moved")
	}

	for _, hostile := range []struct {
		name   string
		target string
	}{
		{"broken", "../scripts/missing.sh"},
		{"escaping", "../../outside.sh"},
	} {
		t.Run(hostile.name, func(t *testing.T) {
			if err := os.Remove(link); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(hostile.target, link); err != nil {
				t.Fatal(err)
			}
			if _, err := resolveCheckInputs(root, check, mustTreeSnapshot(t, root)); err == nil {
				t.Fatalf("%s declared symlink resolved without widening", hostile.name)
			}
		})
	}
}

func TestDeclaredCheckFileDistinguishesAbsentFromPresentEmpty(t *testing.T) {
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	check := registry.Check{Name: "gate-entry-contract", Inputs: registry.InputGateEntry}
	absent, err := resolveCheckInputs(root, check, mustTreeSnapshot(t, root))
	if err != nil {
		t.Fatal(err)
	}
	writeGateTestFile(t, root, ".bench/gate.sh", "", 0o644)
	empty, err := resolveCheckInputs(root, check, mustTreeSnapshot(t, root))
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(absent, empty) {
		t.Fatal("absent declared file and present empty file resolved to the same identity material")
	}
}

func TestSharedConformanceImplementationMovesEveryOrdinaryIdentity(t *testing.T) {
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	writeGateTestFile(t, root, "internal/conformance/shared.go", "package conformance\n\nconst shared = 1\n", 0o644)
	writeGateTestFile(t, root, "ROADMAP.md", "# roadmap\n", 0o644)
	before, err := ResolveConformanceCheckIdentities(root, registry.Ship)
	if err != nil {
		t.Fatal(err)
	}
	writeGateTestFile(t, root, "internal/conformance/shared.go", "package conformance\n\nconst shared = 2\n", 0o644)
	after, err := ResolveConformanceCheckIdentities(root, registry.Ship)
	if err != nil {
		t.Fatal(err)
	}
	for _, check := range registry.Checks {
		if check.Meta {
			if _, found := after[check.Name]; found {
				t.Fatalf("meta check %s received a reusable identity", check.Name)
			}
			continue
		}
		if !check.RunsAt(registry.Ship) {
			continue
		}
		if check.Inputs == "" {
			t.Fatalf("ordinary check %s has no input declaration", check.Name)
		}
		if before[check.Name] == after[check.Name] {
			t.Fatalf("ordinary check %s identity survived shared implementation movement", check.Name)
		}
	}
	got, want := mapKeys(after), ordinaryCheckNames(registry.Ship)
	sort.Strings(got)
	sort.Strings(want)
	if !slices.Equal(got, want) {
		t.Fatalf("resolved identities = %v, want ordinary checks %v", got, want)
	}
}

func TestExternalConformanceImplementationBindsTransitiveGoDependencies(t *testing.T) {
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	writeGateTestFile(t, root, "go.mod", "module identityfixture\n\ngo 1.21\n", 0o644)
	writeGateTestFile(t, root, "internal/conformance/shared.go", "package conformance\n\nconst shared = 1\n", 0o644)
	writeGateTestFile(t, root, "internal/maps/maps.go", "package maps\n\nimport \"identityfixture/internal/bounds\"\n\nfunc ValidateDecisionMapTree() { _ = bounds.ControlRecordLimit }\n", 0o644)
	writeGateTestFile(t, root, "internal/bounds/bounds.go", "package bounds\n\nconst ControlRecordLimit = 1\n", 0o644)

	before, err := ResolveConformanceCheckIdentities(root, registry.Dev)
	if err != nil {
		t.Fatal(err)
	}
	writeGateTestFile(t, root, "internal/bounds/bounds.go", "package bounds\n\nconst ControlRecordLimit = 2\n", 0o644)
	after, err := ResolveConformanceCheckIdentities(root, registry.Dev)
	if err != nil {
		t.Fatal(err)
	}
	if before["decision-map-integrity"] == after["decision-map-integrity"] {
		t.Fatal("decision-map-integrity identity survived a transitive internal/bounds implementation change")
	}
}

func mapKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func ordinaryCheckNames(tier registry.Tier) []string {
	var names []string
	for _, check := range registry.Checks {
		if !check.Meta && check.RunsAt(tier) {
			names = append(names, check.Name)
		}
	}
	return names
}
