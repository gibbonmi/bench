package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

func TestCanaryFixtureRegistryClassifiesEveryFixture(t *testing.T) {
	h := NewHarness(t)
	fixturesDir := h.KitPath("tests", "canary")
	fixturePaths := canaryFixturePaths(t, fixturesDir)

	for name, fx := range fixturePaths {
		family := fx.Family
		if family == "" || !familyIsBound(family) {
			t.Errorf("canary fixture %q has unknown conformance family %q", name, family)
			continue
		}
		reg, ok := fixtureRegistrationFor(name, family)
		if !ok {
			t.Errorf("canary fixture %q is unclassified", name)
			continue
		}
		if reg.Owner != ownerConformance {
			t.Errorf("canary fixture %q has invalid owner %q", name, reg.Owner)
		}
		if reg.Owner == ownerConformance && len(reg.ShellSources) == 0 && len(reg.GoSources) == 0 {
			t.Errorf("canary fixture %q has no conformance owner source", name)
		}
		for _, source := range reg.GoSources {
			if _, err := os.Stat(h.KitPath(filepath.FromSlash(source))); err != nil {
				t.Errorf("canary fixture %q names missing Go owner source %s: %v", name, source, err)
			}
		}
		if _, err := os.Stat(filepath.Join(fx.Dir, "EXPECT")); err != nil {
			t.Errorf("canary fixture %q has no EXPECT: %v", name, err)
		}
	}

	for name := range canaryFixtureRegistry {
		if _, ok := fixturePaths[name]; !ok {
			t.Errorf("registry names nonexistent canary fixture %q", name)
		}
	}
	for family := range canaryFixtureFamilyRegistry {
		found := false
		for _, fixture := range fixturePaths {
			if fixture.Family == family {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("family registry names nonexistent canary family %q", family)
		}
		if !familyIsBound(family) {
			t.Errorf("family registry names unbound conformance family %q", family)
		}
	}
}

func TestCanaryFixtureFamilyRegistrationInheritance(t *testing.T) {
	registration, found := fixtureRegistrationFor("future-decision-map-fixture", "decision-map-integrity")
	if !found || registration.Owner != ownerConformance || len(registration.GoSources) == 0 {
		t.Fatalf("future decision-map fixture registration = %#v, %v; want conformance family registration with Go sources", registration, found)
	}

	if _, found := fixtureRegistrationFor("unregistered-fixture", "unregistered-family"); found {
		t.Fatal("unregistered canary family resolved a fixture registration")
	}
}

func TestRetiredConformanceFixturesDoNotLeaveShellTwinMessages(t *testing.T) {
	h := NewHarness(t)
	fixturePaths := canaryFixturePaths(t, h.KitPath("tests", "canary"))
	shellText := map[string]string{}
	for fixture, reg := range canaryFixtureRegistry {
		if reg.Owner != ownerConformance {
			continue
		}
		for _, source := range reg.ShellSources {
			if _, ok := shellText[source]; ok {
				continue
			}
			data, err := os.ReadFile(h.KitPath(filepath.FromSlash(source)))
			if err != nil {
				if os.IsNotExist(err) {
					shellText[source] = ""
					continue
				}
				t.Fatalf("read %s: %v", source, err)
			}
			shellText[source] = string(data)
		}
		fx, ok := fixturePaths[fixture]
		if !ok {
			t.Fatalf("registry names nonexistent conformance fixture %q", fixture)
		}
		expect := readExpectation(t, filepath.Join(fx.Dir, "EXPECT"))
		for _, source := range reg.ShellSources {
			text := shellText[source]
			if strings.Contains(text, expect) {
				t.Errorf("retired conformance fixture %q EXPECT substring still appears in %s", fixture, source)
			}
		}
	}
}

func canaryFixturePaths(t *testing.T, fixturesDir string) map[string]canary.Fixture {
	t.Helper()
	families, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Fatalf("read canary fixtures: %v", err)
	}
	for _, family := range families {
		if !family.IsDir() {
			continue
		}
		// A family is canonical when inventory resolution can identify it, through a
		// family binding or through its own fixture marker. Only a family with neither
		// is unattributable. A flat fixture carries its own binding, not a family.
		if canary.IsConformanceFamily(filepath.Join(fixturesDir, family.Name())) && !familyIsBound(family.Name()) {
			t.Errorf("canary family %q is not canonical", family.Name())
		}
	}
	// Inventory discovery enumerates fixtures and enforces base-name uniqueness. A second
	// walk here would disagree with the producer the check consumes.
	discovered, err := canary.Fixtures(fixturesDir)
	if err != nil {
		t.Fatalf("walk canary fixtures: %v", err)
	}
	return discovered
}
