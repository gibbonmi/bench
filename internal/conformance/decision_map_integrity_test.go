package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/maps"
)

var decisionMapIntegrityFixtureCategories = map[string][]string{
	"graph": {
		"graph-cycle", "graph-dangling", "graph-duplicate-blocker", "graph-duplicate-id", "graph-resolved-on-unresolved", "graph-self-edge",
	},
	"readiness": {
		"readiness-compiled-shaping", "readiness-fog", "readiness-unresolved",
	},
	"schema": {
		"schema-duplicate-answer", "schema-duplicate-blocked-by", "schema-duplicate-destination", "schema-duplicate-discretion", "schema-duplicate-fog", "schema-duplicate-out-of-scope", "schema-duplicate-question", "schema-duplicate-sources", "schema-duplicate-status", "schema-duplicate-title", "schema-duplicate-type", "schema-handoff", "schema-malformed-blocked-by", "schema-missing-answer", "schema-missing-blocked-by", "schema-missing-destination", "schema-missing-discretion", "schema-missing-fog", "schema-missing-out-of-scope", "schema-missing-question", "schema-missing-sources", "schema-missing-status", "schema-missing-ticket", "schema-missing-title", "schema-missing-type", "schema-status", "schema-unsupported-type",
	},
	"source": {
		"source-absolute-path", "source-empty-path", "source-escape-path", "source-invalid-url", "source-missing-drift", "source-missing-path", "source-missing-supports", "source-not-bullet", "source-unknown-kind",
	},
	"terminal-list": {
		"terminal-discretion-prose", "terminal-fog-prose", "terminal-out-of-scope-prose",
	},
}

// The inventory is independently authored omission coverage: deleting one named fixture
// must red even though the fixture directories and EXPECT files remain the test corpus.
func validateDecisionMapIntegrityFixtureInventory(fixtures map[string]canary.Fixture) error {
	seen := make(map[string]bool)
	expected := make([]string, 0, 48)
	for category, names := range decisionMapIntegrityFixtureCategories {
		for _, name := range names {
			if seen[name] {
				return fmt.Errorf("fixture %q appears twice in the %s inventory", name, category)
			}
			seen[name] = true
			expected = append(expected, name)
		}
	}
	if len(expected) != 48 {
		return fmt.Errorf("decision-map fixture inventory has %d entries, want 48", len(expected))
	}
	sort.Strings(expected)
	actual := make([]string, 0, len(fixtures))
	expects := make(map[string]string, len(fixtures))
	for name, fixture := range fixtures {
		actual = append(actual, name)
		data, err := os.ReadFile(filepath.Join(fixture.Dir, "EXPECT"))
		if err != nil {
			return fmt.Errorf("read %s EXPECT: %w", name, err)
		}
		expect := strings.TrimSpace(string(data))
		if expect == "" {
			return fmt.Errorf("fixture %q has an empty EXPECT", name)
		}
		if other, exists := expects[expect]; exists {
			return fmt.Errorf("fixtures %q and %q share EXPECT %q", other, name, expect)
		}
		expects[expect] = name
	}
	sort.Strings(actual)
	if !slices.Equal(actual, expected) {
		return fmt.Errorf("decision-map fixture inventory = %v, want %v", actual, expected)
	}
	return nil
}

func TestDecisionMapIntegrityCheckValidatesEveryCandidate(t *testing.T) {
	h := NewHarness(t)
	root := t.TempDir()
	writeMap := func(path, document string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(document), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	valid := strings.Replace(maps.DecisionMapTemplate(), "<answer>", "Resolved.", 1)
	writeMap(filepath.Join(root, "decisions", "active.md"), valid)
	writeMap(filepath.Join(root, "specs", "compiled", "decisions", "compiled.md"), strings.Replace(valid, "Status: shaping", "Status: ready", 1))
	writeMap(filepath.Join(root, "specs", "no-map", "spec.md"), "# No map\n")
	if diagnostics := RunConformance(root, h.KitRoot, registry.Dev, "decision-map-integrity"); len(diagnostics) != 0 {
		t.Fatalf("valid active and compiled maps diagnostics = %v", diagnostics)
	}

	writeMap(filepath.Join(root, "specs", "broken", "decisions", "broken.md"), "# Broken\n")
	writeMap(filepath.Join(root, "decisions", "graph.md"), strings.Replace(valid, "Blocked by: none", "Blocked by: #1", 1))
	diagnostics := RunConformance(root, h.KitRoot, registry.Dev, "decision-map-integrity")
	for _, want := range []string{
		"specs/broken/decisions/broken.md: missing Status",
		"decisions/graph.md: ticket #1: <decision question> self-edge #1 -> #1",
	} {
		if !containsDiagnostic(diagnostics, want) {
			t.Fatalf("candidate diagnostics = %v, want %q", diagnostics, want)
		}
	}
}

func TestDecisionMapIntegrityFixturesBite(t *testing.T) {
	h := NewHarness(t)
	fixtures, err := canary.Fixtures(filepath.Join(h.KitRoot, "tests", "canary", "decision-map-integrity"))
	if err != nil {
		t.Fatal(err)
	}
	if err := validateDecisionMapIntegrityFixtureInventory(fixtures); err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(fixtures))
	for name := range fixtures {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			root := materializeConformanceFixture(t, name)
			expect := readExpectation(t, filepath.Join(canaryFixturePath(t, h.KitRoot, name), "EXPECT"))

			diagnostics := RunConformance(root, h.KitRoot, registry.Dev, "decision-map-integrity")
			if !containsDiagnostic(diagnostics, expect) {
				t.Fatalf("%s did not bite under scoped Go conformance; want %q in diagnostics:\n%s", name, expect, strings.Join(diagnostics, "\n"))
			}
		})
	}
}

func TestDecisionMapIntegrityFixtureInventoryRejectsDeletion(t *testing.T) {
	h := NewHarness(t)
	canaryRoot := filepath.Join(t.TempDir(), "tests", "canary")
	copyRoot := filepath.Join(canaryRoot, "decision-map-integrity")
	source := filepath.Join(h.KitRoot, "tests", "canary", "decision-map-integrity")
	if err := os.MkdirAll(copyRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := canary.MaterializeFixture(source, copyRoot); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(copyRoot, "schema-status")); err != nil {
		t.Fatal(err)
	}
	fixtures, err := canary.Fixtures(canaryRoot)
	if err != nil {
		t.Fatal(err)
	}
	err = validateDecisionMapIntegrityFixtureInventory(fixtures)
	if err == nil || !strings.Contains(err.Error(), "schema-status") {
		t.Fatalf("deleted fixture inventory error = %v, want schema-status omission", err)
	}
}
