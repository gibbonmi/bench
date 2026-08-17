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
)

const roadmapDetailIntegrityFamily = "roadmap-detail-integrity"

// roadmapDetailIntegrityFixtureClasses names the diagnostic class each fixture in the
// family plants — one fixture per class the split-board loader can report. The bite proof
// derives its corpus from the fixture directories themselves, so it cannot see a class
// whose fixture is gone; this independently authored inventory is what makes that
// omission red.
var roadmapDetailIntegrityFixtureClasses = map[string]string{
	"roadmap-missing-detail-owner": "missing detail owner",
	"roadmap-orphan-detail":        "orphan detail file",
	"roadmap-inline-body":          "carries an inline body",
	"roadmap-heading-mismatch":     "heading does not match",
	"roadmap-unrecognized-file":    "unrecognized file under",
	"roadmap-duplicate-row":        "duplicate row",
	"roadmap-wrapped-heading":      "wrapped heading",
	"roadmap-unreadable-detail":    "detail file: not a regular file",
}

func validateRoadmapDetailIntegrityFixtureInventory(fixtures map[string]canary.Fixture) error {
	expected := make([]string, 0, len(roadmapDetailIntegrityFixtureClasses))
	for name := range roadmapDetailIntegrityFixtureClasses {
		expected = append(expected, name)
	}
	if len(expected) != 8 {
		return fmt.Errorf("roadmap-detail-integrity fixture inventory has %d entries, want 8", len(expected))
	}
	sort.Strings(expected)
	actual := make([]string, 0, len(expected))
	planted := make(map[string]string, len(expected))
	for name, fixture := range fixtures {
		if fixture.Family != roadmapDetailIntegrityFamily {
			continue
		}
		actual = append(actual, name)
		class, named := roadmapDetailIntegrityFixtureClasses[name]
		if !named {
			continue
		}
		data, err := os.ReadFile(filepath.Join(fixture.Dir, "EXPECT"))
		if err != nil {
			return fmt.Errorf("read %s EXPECT: %w", name, err)
		}
		expect := strings.TrimSpace(string(data))
		if !strings.Contains(expect, class) {
			return fmt.Errorf("fixture %q EXPECT %q does not plant its %q class", name, expect, class)
		}
		if other, exists := planted[class]; exists {
			return fmt.Errorf("fixtures %q and %q both plant the %q class", other, name, class)
		}
		planted[class] = name
	}
	sort.Strings(actual)
	if !slices.Equal(actual, expected) {
		return fmt.Errorf("roadmap-detail-integrity fixture inventory = %v, want %v", actual, expected)
	}
	return nil
}

func TestRoadmapDetailIntegrityFixturesCoverEveryDiagnosticClass(t *testing.T) {
	h := NewHarness(t)
	fixtures, err := canary.Fixtures(filepath.Join(h.KitRoot, "tests", "canary"))
	if err != nil {
		t.Fatal(err)
	}
	if err := validateRoadmapDetailIntegrityFixtureInventory(fixtures); err != nil {
		t.Fatal(err)
	}
}

func TestRoadmapDetailIntegrityFixtureInventoryRejectsDeletion(t *testing.T) {
	h := NewHarness(t)
	canaryRoot := filepath.Join(t.TempDir(), "tests", "canary")
	copyRoot := filepath.Join(canaryRoot, roadmapDetailIntegrityFamily)
	source := filepath.Join(h.KitRoot, "tests", "canary", roadmapDetailIntegrityFamily)
	if err := os.MkdirAll(copyRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := canary.MaterializeFixture(source, copyRoot); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(copyRoot, "roadmap-wrapped-heading")); err != nil {
		t.Fatal(err)
	}
	fixtures, err := canary.Fixtures(canaryRoot)
	if err != nil {
		t.Fatal(err)
	}
	err = validateRoadmapDetailIntegrityFixtureInventory(fixtures)
	if err == nil || !strings.Contains(err.Error(), "roadmap-wrapped-heading") {
		t.Fatalf("deleted fixture inventory error = %v, want roadmap-wrapped-heading omission", err)
	}
}
