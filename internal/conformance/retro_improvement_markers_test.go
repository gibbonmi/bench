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

const retroImprovementMarkersFamily = "retro-improvement-markers"

// retroImprovementMarkersFixtureClasses names the diagnostic class each fixture in the
// family plants. The generic bite proof derives its corpus from the fixture directories
// themselves, so it cannot see a class whose fixture is gone. This independently authored
// inventory makes that omission red.
var retroImprovementMarkersFixtureClasses = map[string]string{
	"retro-item-unmarked": "improvement item carries no destination marker",
}

func validateRetroImprovementMarkersFixtureInventory(fixtures map[string]canary.Fixture) error {
	expected := make([]string, 0, len(retroImprovementMarkersFixtureClasses))
	for name := range retroImprovementMarkersFixtureClasses {
		expected = append(expected, name)
	}
	if len(expected) != 1 {
		return fmt.Errorf("retro-improvement-markers fixture inventory has %d entries, want 1", len(expected))
	}
	sort.Strings(expected)
	actual := make([]string, 0, len(expected))
	planted := make(map[string]string, len(expected))
	for name, fixture := range fixtures {
		if fixture.Family != retroImprovementMarkersFamily {
			continue
		}
		actual = append(actual, name)
		class, named := retroImprovementMarkersFixtureClasses[name]
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
		return fmt.Errorf("retro-improvement-markers fixture inventory = %v, want %v", actual, expected)
	}
	return nil
}

func TestRetroImprovementMarkersFixturesCoverEveryDiagnosticClass(t *testing.T) {
	h := NewHarness(t)
	fixtures, err := canary.Fixtures(filepath.Join(h.KitRoot, "tests", "canary"))
	if err != nil {
		t.Fatal(err)
	}
	if err := validateRetroImprovementMarkersFixtureInventory(fixtures); err != nil {
		t.Fatal(err)
	}
}

func TestRetroImprovementMarkersFixtureInventoryRejectsDeletion(t *testing.T) {
	h := NewHarness(t)
	fixtures, err := canary.Fixtures(filepath.Join(h.KitRoot, "tests", "canary"))
	if err != nil {
		t.Fatal(err)
	}
	// The family holds one fixture, so removing its directory would empty the canary
	// inventory and fail the loader before the validator ever ran. Dropping the entry
	// from the loaded inventory presents the validator with the same absence.
	delete(fixtures, "retro-item-unmarked")
	err = validateRetroImprovementMarkersFixtureInventory(fixtures)
	if err == nil || !strings.Contains(err.Error(), "retro-item-unmarked") {
		t.Fatalf("deleted fixture inventory error = %v, want retro-item-unmarked omission", err)
	}
}
