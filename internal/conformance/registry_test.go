package conformance

import (
	"os"
	"path/filepath"
	"sort"
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
	Family       string
	ShellSources []string
	ShellLabels  []string
}

var canaryFixtureRegistry = map[string]fixtureRegistration{
	"invalid-json":             conformanceFixture("load-validity-metadata", ".bench/gate.sh"),
	"codex-hooks-broken":       conformanceFixture("load-validity-metadata", ".bench/gate.sh"),
	"bad-frontmatter":          conformanceFixture("load-validity-metadata", ".bench/gate.sh"),
	"claude-skills-unmirrored": conformanceFixture("load-validity-metadata", ".bench/gate.sh"),
	"extensionless-gate-ref":   conformanceFixture("load-validity-metadata", ".bench/gate.sh"),
	"shared-rule-drift":        conformanceFixture("load-validity-metadata", ".bench/gate.sh"),
	"readme-shared-rule-drift": conformanceFixture("load-validity-metadata", ".bench/gate.sh"),

	"dangling-index":      conformanceFixture("skills-index-command-adapters", ".bench/gate.sh", ".bench/skills-index.sh"),
	"missing-index-field": conformanceFixture("skills-index-command-adapters", ".bench/gate.sh", ".bench/skills-index.sh"),
	"stale-index-wording": conformanceFixture("skills-index-command-adapters", ".bench/gate.sh", ".bench/skills-index.sh"),
	"unindexed-skill":     conformanceFixture("skills-index-command-adapters", ".bench/gate.sh", ".bench/skills-index.sh"),

	"stale-command-reference":       conformanceFixture("docs-currency-token-diet", ".bench/gate-docs-contracts.sh"),
	"stale-codex-adapter-reference": conformanceFixture("docs-currency-token-diet", ".bench/gate-docs-contracts.sh"),
	"stale-cli-doc-reference":       conformanceFixture("docs-currency-token-diet", ".bench/gate-docs-contracts.sh"),
	"historical-marker-prose":       conformanceFixture("docs-currency-token-diet", ".bench/gate-docs-contracts.sh"),
	"benchref-missing":              conformanceFixture("docs-currency-token-diet", ".bench/gate-docs-contracts.sh"),
	"benchref-pointer-dropped":      conformanceFixture("docs-currency-token-diet", ".bench/gate-docs-contracts.sh"),
	"benchref-imported":             conformanceFixture("docs-currency-token-diet", ".bench/gate-docs-contracts.sh"),
	"benchref-section-duplicated":   conformanceFixture("docs-currency-token-diet", ".bench/gate-docs-contracts.sh"),
	"readme-command-first":          conformanceFixture("docs-currency-token-diet", ".bench/gate-docs-contracts.sh"),

	"acceptance-coverage-anchor":        conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh"),
	"coverage-axis-anchor":              conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh"),
	"command-handoff-anchor":            conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh"),
	"debug-archaeology-anchor":          conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh"),
	"edge-inventory-anchor":             conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh"),
	"implement-spec-status-flip-anchor": conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh"),
	"shape-idea-bypass":                 conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh"),
	"shape-idea-bypass-wrapped":         conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh"),
	"shape-idea-handoff-anchor":         conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh"),
	"story-line-anchor-missing":         conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh", ".bench/gate-line-contracts.sh"),
	"write-spec-handoff-anchor":         conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh"),
	"write-spec-map-required":           conformanceFixture("workflow-guidance-anchors", ".bench/gate-docs-contracts.sh"),
	"line-anchor-missing":               conformanceFixture("workflow-guidance-anchors", ".bench/gate-line-contracts.sh"),
	"broken-coverage-map":               conformanceFixture("coverage-map-validation", ".bench/gate-docs-contracts.sh"),

	"line-binding-prose-drift": conformanceFixture("line-routing", ".bench/gate-line-contracts.sh"),
	"agent-hook-unwired":       conformanceFixture("line-routing", ".bench/gate-line-contracts.sh"),
	"agent-hook-broken":        conformanceFixture("line-routing", ".bench/gate-line-contracts.sh"),
	"adapter-line-broken":      conformanceFixture("line-routing", ".bench/gate-line-contracts.sh"),

	"missing-files-entry":             conformanceFixture("package-core-guard", ".bench/gate-package-contracts.sh"),
	"go-build-broken":                 conformanceFixture("package-core-guard", ".bench/gate-go-contracts.sh"),
	"go-test-failing":                 conformanceFixture("package-core-guard", ".bench/gate-go-contracts.sh"),
	"guard-describe-boundary-dropped": conformanceFixture("package-core-guard", ".bench/gate-axi-contracts.sh"),

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

func conformanceFixture(family string, shellSources ...string) fixtureRegistration {
	return fixtureRegistration{Owner: ownerConformance, Family: family, ShellSources: shellSources, ShellLabels: []string{family}}
}

func behaviorFixture() fixtureRegistration {
	return fixtureRegistration{Owner: ownerBehavior, Family: "behavior-owned"}
}

func TestCanaryFixtureRegistryClassifiesEveryFixture(t *testing.T) {
	h := NewHarness(t)
	fixturesDir := h.KitPath("tests", "canary")
	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Fatalf("read canary fixtures: %v", err)
	}

	var names []string
	for _, ent := range entries {
		if ent.IsDir() {
			names = append(names, ent.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		reg, ok := canaryFixtureRegistry[name]
		if !ok {
			t.Errorf("canary fixture %q is unclassified", name)
			continue
		}
		if reg.Owner != ownerConformance && reg.Owner != ownerBehavior {
			t.Errorf("canary fixture %q has invalid owner %q", name, reg.Owner)
		}
		if reg.Owner == ownerConformance && reg.Family == "" {
			t.Errorf("canary fixture %q has no conformance family", name)
		}
		if reg.Owner == ownerConformance && len(reg.ShellSources) == 0 {
			t.Errorf("canary fixture %q has no retired shell source", name)
		}
		if reg.Owner == ownerBehavior && reg.Family != "behavior-owned" {
			t.Errorf("behavior fixture %q has family %q, want behavior-owned", name, reg.Family)
		}
		if _, err := os.Stat(filepath.Join(fixturesDir, name, "EXPECT")); err != nil {
			t.Errorf("canary fixture %q has no EXPECT: %v", name, err)
		}
	}

	for name := range canaryFixtureRegistry {
		if _, err := os.Stat(filepath.Join(fixturesDir, name)); err != nil {
			t.Errorf("registry names nonexistent canary fixture %q: %v", name, err)
		}
	}
}

func TestRetiredConformanceFixturesDoNotLeaveShellTwinMessages(t *testing.T) {
	h := NewHarness(t)
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
		expect := readExpectation(t, h.KitPath("tests", "canary", fixture, "EXPECT"))
		for _, source := range reg.ShellSources {
			text := shellText[source]
			if strings.Contains(text, expect) {
				t.Errorf("retired conformance fixture %q EXPECT substring still appears in %s", fixture, source)
			}
			for _, label := range reg.ShellLabels {
				if strings.Contains(text, label) {
					t.Errorf("retired conformance fixture %q shell label %q still appears in %s", fixture, label, source)
				}
			}
		}
	}
}
