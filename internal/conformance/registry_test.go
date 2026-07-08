package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fixtureOwner string

const (
	ownerConformance fixtureOwner = "conformance"
	ownerBehavior    fixtureOwner = "behavior"
)

type fixtureRegistration struct {
	Owner        fixtureOwner
	ShellSources []string
}

var canaryFixtureRegistry = map[string]fixtureRegistration{
	"invalid-json":              conformanceFixture(".bench/gate.sh"),
	"codex-hooks-broken":        conformanceFixture(".bench/gate.sh"),
	"codex-hooks-timeout":       conformanceFixture(".bench/gate.sh"),
	"codex-hooks-timeout-typed": conformanceFixture(".bench/gate.sh"),
	"bad-frontmatter":           conformanceFixture(".bench/gate.sh"),
	"claude-skills-unmirrored":  conformanceFixture(".bench/gate.sh"),
	"extensionless-gate-ref":    conformanceFixture(".bench/gate.sh"),
	"shared-rule-drift":         conformanceFixture(".bench/gate.sh"),
	"readme-shared-rule-drift":  conformanceFixture(".bench/gate.sh"),

	"dangling-index":                conformanceFixture(".bench/gate.sh", ".bench/skills-index.sh"),
	"roadmap-promotion-persistence": conformanceFixture(".bench/gate.sh"),
	"missing-index-field":           conformanceFixture(".bench/gate.sh", ".bench/skills-index.sh"),
	"stale-index-wording":           conformanceFixture(".bench/gate.sh", ".bench/skills-index.sh"),
	"unindexed-skill":               conformanceFixture(".bench/gate.sh", ".bench/skills-index.sh"),

	"stale-command-reference":       conformanceFixture(".bench/gate-docs-contracts.sh"),
	"stale-codex-adapter-reference": conformanceFixture(".bench/gate-docs-contracts.sh"),
	"retired-command-reference":     conformanceFixture(".bench/gate-docs-contracts.sh"),
	"stale-cli-doc-reference":       conformanceFixture(".bench/gate-docs-contracts.sh"),
	"missing-cli-inventory":         conformanceFixture(".bench/gate-docs-contracts.sh"),
	"historical-marker-prose":       conformanceFixture(".bench/gate-docs-contracts.sh"),
	"benchref-missing":              conformanceFixture(".bench/gate-docs-contracts.sh"),
	"benchref-pointer-dropped":      conformanceFixture(".bench/gate-docs-contracts.sh"),
	"benchref-imported":             conformanceFixture(".bench/gate-docs-contracts.sh"),
	"benchref-section-duplicated":   conformanceFixture(".bench/gate-docs-contracts.sh"),
	"dogfood-referent-shipped":      conformanceFixture(".bench/gate-docs-contracts.sh"),
	"readme-command-first":          conformanceFixture(".bench/gate-docs-contracts.sh"),

	"acceptance-coverage-anchor":        conformanceFixture(".bench/gate-docs-contracts.sh"),
	"coverage-axis-anchor":              conformanceFixture(".bench/gate-docs-contracts.sh"),
	"command-handoff-anchor":            conformanceFixture(".bench/gate-docs-contracts.sh"),
	"debug-archaeology-anchor":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"edge-inventory-anchor":             conformanceFixture(".bench/gate-docs-contracts.sh"),
	"implement-spec-status-flip-anchor": conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-bypass":                 conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-bypass-wrapped":         conformanceFixture(".bench/gate-docs-contracts.sh"),
	"what-next-anchor":                  conformanceFixture(".bench/gate-docs-contracts.sh"),
	"spec-retire-roadmap-row":           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"capture-sink-anchor":               conformanceFixture(".bench/gate-docs-contracts.sh"),
	"review-persistence-anchor":         conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-handoff-anchor":         conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-grill-continuation":     conformanceFixture(".bench/gate-docs-contracts.sh"),
	"story-line-anchor-missing":         conformanceFixture(".bench/gate-docs-contracts.sh", ".bench/gate-line-contracts.sh"),
	"write-spec-handoff-anchor":         conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-map-required":           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"line-anchor-missing":               conformanceFixture(".bench/gate-line-contracts.sh"),
	"broken-coverage-map":               conformanceFixture(".bench/gate-docs-contracts.sh"),

	"line-binding-prose-drift": conformanceFixture(".bench/gate-line-contracts.sh"),
	"agent-hook-unwired":       conformanceFixture(".bench/gate-line-contracts.sh"),
	"agent-hook-broken":        conformanceFixture(".bench/gate-line-contracts.sh"),
	"stop-hook-unwired":        conformanceFixture(".bench/gate-line-contracts.sh"),
	"adapter-line-broken":      conformanceFixture(".bench/gate-line-contracts.sh"),

	"missing-files-entry":             conformanceFixture(".bench/gate-package-contracts.sh"),
	"go-build-broken":                 conformanceFixture(".bench/gate-go-contracts.sh"),
	"go-test-failing":                 conformanceFixture(".bench/gate-go-contracts.sh"),
	"guard-describe-boundary-dropped": conformanceFixture(".bench/gate-axi-contracts.sh"),

	"doctor-foreign-clobbered":     behaviorFixture(),
	"doctor-manager-dir-chosen":    behaviorFixture(),
	"doctor-stale-silent":          behaviorFixture(),
	"postinstall-guard-bypassed":   behaviorFixture(),
	"postinstall-nonzero-exit":     behaviorFixture(),
	"session-start-advice-dropped": behaviorFixture(),
	"wrapper-args-dropped":         behaviorFixture(),
	"status-regressed":             behaviorFixture(),
	"roadmap-regressed":            behaviorFixture(),
	"unscaffolded-bench-file":      behaviorFixture(),
	"toon-escaping-dropped":        behaviorFixture(),
	"learnings-parse-broken":       behaviorFixture(),
	"guards-aggregation-dropped":   behaviorFixture(),
	"coverage-extraction-dropped":  behaviorFixture(),
	"diff-recorded-base-dropped":   behaviorFixture(),
}

func conformanceFixture(shellSources ...string) fixtureRegistration {
	return fixtureRegistration{Owner: ownerConformance, ShellSources: shellSources}
}

func behaviorFixture() fixtureRegistration {
	return fixtureRegistration{Owner: ownerBehavior}
}

func TestCanaryFixtureRegistryClassifiesEveryFixture(t *testing.T) {
	h := NewHarness(t)
	fixturesDir := h.KitPath("tests", "canary")
	fixturePaths := canaryFixturePaths(t, fixturesDir)

	for name, path := range fixturePaths {
		family := filepath.Base(filepath.Dir(path))
		wantOwner := ownerConformance
		if family == "behavior-owned" {
			wantOwner = ownerBehavior
		} else if !isConformanceFamily(family) {
			t.Errorf("canary fixture %q has unknown conformance family %q", name, family)
			continue
		}
		reg, ok := canaryFixtureRegistry[name]
		if !ok {
			t.Errorf("canary fixture %q is unclassified", name)
			continue
		}
		if reg.Owner != ownerConformance && reg.Owner != ownerBehavior {
			t.Errorf("canary fixture %q has invalid owner %q", name, reg.Owner)
		}
		if reg.Owner == ownerConformance && len(reg.ShellSources) == 0 {
			t.Errorf("canary fixture %q has no retired shell source", name)
		}
		if reg.Owner != wantOwner {
			t.Errorf("canary fixture %q under %q has owner %q, want %q", name, family, reg.Owner, wantOwner)
		}
		if _, err := os.Stat(filepath.Join(path, "EXPECT")); err != nil {
			t.Errorf("canary fixture %q has no EXPECT: %v", name, err)
		}
	}

	for name := range canaryFixtureRegistry {
		if _, ok := fixturePaths[name]; !ok {
			t.Errorf("registry names nonexistent canary fixture %q", name)
		}
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
		fixturePath, ok := fixturePaths[fixture]
		if !ok {
			t.Fatalf("registry names nonexistent conformance fixture %q", fixture)
		}
		expect := readExpectation(t, filepath.Join(fixturePath, "EXPECT"))
		for _, source := range reg.ShellSources {
			text := shellText[source]
			if strings.Contains(text, expect) {
				t.Errorf("retired conformance fixture %q EXPECT substring still appears in %s", fixture, source)
			}
		}
	}
}

func canaryFixturePaths(t *testing.T, fixturesDir string) map[string]string {
	t.Helper()
	families, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Fatalf("read canary fixtures: %v", err)
	}
	paths := map[string]string{}
	for _, family := range families {
		if !family.IsDir() {
			continue
		}
		if family.Name() != "behavior-owned" && !isConformanceFamily(family.Name()) {
			t.Errorf("canary family %q is not canonical", family.Name())
		}
		familyDir := filepath.Join(fixturesDir, family.Name())
		entries, err := os.ReadDir(familyDir)
		if err != nil {
			t.Fatalf("read canary family %q: %v", family.Name(), err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			if first := paths[entry.Name()]; first != "" {
				t.Errorf("canary fixture %q appears at both %s and %s", entry.Name(), first, filepath.Join(familyDir, entry.Name()))
				continue
			}
			paths[entry.Name()] = filepath.Join(familyDir, entry.Name())
		}
	}
	return paths
}
