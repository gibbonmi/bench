package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

type fixtureOwner string

const (
	ownerConformance fixtureOwner = "conformance"
)

type fixtureRegistration struct {
	Owner        fixtureOwner
	ShellSources []string
	GoSources    []string
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

	"stale-command-reference":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"introduces-undeclared-command":    conformanceFixture(".bench/gate-docs-contracts.sh"),
	"stale-codex-adapter-reference":    conformanceFixture(".bench/gate-docs-contracts.sh"),
	"retired-command-reference":        conformanceFixture(".bench/gate-docs-contracts.sh"),
	"stale-cli-doc-reference":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"stale-skill-cli-reference":        conformanceFixture(".bench/gate-docs-contracts.sh"),
	"missing-cli-inventory":            conformanceFixture(".bench/gate-docs-contracts.sh"),
	"historical-marker-prose":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"benchref-pointer-dropped":         conformanceFixture(".bench/gate-docs-contracts.sh"),
	"benchref-imported":                conformanceFixture(".bench/gate-docs-contracts.sh"),
	"benchref-section-duplicated":      conformanceFixture(".bench/gate-docs-contracts.sh"),
	"dogfood-referent-shipped":         conformanceFixture(".bench/gate-docs-contracts.sh"),
	"readme-command-first":             conformanceFixture(".bench/gate-docs-contracts.sh"),
	"signal-vocabulary-drift":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"structured-phase-progress-anchor": conformanceFixture(".bench/gate-docs-contracts.sh"),

	"acceptance-coverage-anchor":                 conformanceFixture(".bench/gate-docs-contracts.sh"),
	"coverage-axis-anchor":                       conformanceFixture(".bench/gate-docs-contracts.sh"),
	"command-handoff-anchor":                     conformanceFixture(".bench/gate-docs-contracts.sh"),
	"debug-archaeology-anchor":                   conformanceFixture(".bench/gate-docs-contracts.sh"),
	"debug-red-commit":                           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"readme-shaping-skip":                        conformanceFixture(".bench/gate-docs-contracts.sh"),
	"implement-spec-inline-exception":            conformanceFixture(".bench/gate-docs-contracts.sh"),
	"edge-inventory-anchor":                      conformanceFixture(".bench/gate-docs-contracts.sh"),
	"fix-pass-sentinel-anchor":                   conformanceFixture(".bench/gate-docs-contracts.sh"),
	"implement-spec-mandatory-delegation-anchor": conformanceFixture(".bench/gate-docs-contracts.sh"),
	"implement-spec-status-flip-anchor":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"drain-anchor":                               conformanceFixture(".bench/gate-docs-contracts.sh"),
	"drain-spec-history-anchor":                  conformanceFixture(".bench/gate-docs-contracts.sh"),
	"drain-roadmap-context-anchor":               conformanceFixture(".bench/gate-docs-contracts.sh"),
	"spec-retire-roadmap-row":                    conformanceFixture(".bench/gate-docs-contracts.sh"),
	"staged-command-sweep-anchor":                conformanceFixture(".bench/gate-docs-contracts.sh"),
	"capture-sink-anchor":                        conformanceFixture(".bench/gate-docs-contracts.sh"),
	"craft-seams-structure-headroom":             conformanceFixture(".bench/gate-docs-contracts.sh"),
	"review-persistence-anchor":                  conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shared-worktree-path-pin":                   conformanceFixture(".bench/gate-docs-contracts.sh"),
	"delegate-parallel-route-anchor":             conformanceFixture(".bench/gate-docs-contracts.sh"),
	"delegate-stash-refusal-anchor":              conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-handoff-anchor":                  conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-grill-continuation":              conformanceFixture(".bench/gate-docs-contracts.sh"),
	"story-line-anchor-missing":                  conformanceFixture(".bench/gate-docs-contracts.sh", ".bench/gate-line-contracts.sh"),
	"write-spec-handoff-anchor":                  conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-review-tier-escalated":           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-review-made-conditional":         conformanceFixture(".bench/gate-docs-contracts.sh"),
	"line-anchor-missing":                        conformanceFixture(".bench/gate-line-contracts.sh"),
	"broken-coverage-map":                        conformanceFixture(".bench/gate-docs-contracts.sh"),
	"no-map-not-historical":                      conformanceFixture(".bench/gate-docs-contracts.sh"),
	"stray-flat-live-spec":                       conformanceFixture(".bench/gate-docs-contracts.sh"),

	"benchkit-spec-ownership":               conformanceFixture(".bench/gate-docs-contracts.sh"),
	"changelog-ticket-vocabulary":           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"context-ticket-vocabulary":             conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-decision-ticket-vocabulary": conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-map-template":               conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-phase-ownership":            conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-situational-map":            conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-artifact-authorization":     conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-authorization-boundary":     conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-conversation-authorization": conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-decision-source":            conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-late-uncertainty":           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-map-sources":                conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-phase-ownership":            conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-ready-map-authorization":    conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-verify-hook-anchor":         conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-light-path-anchor":              conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-stage-routing-anchor":           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-skill-contract-anchor":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-template-anchor":                conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-cross-pointers-anchor":          conformanceFixture(".bench/gate-docs-contracts.sh"),

	"ticket-decision-map-lifecycle-anchor": conformanceFixture(".bench/gate-docs-contracts.sh"),
	"implementation-retro-drain-anchor":    conformanceFixture(".bench/gate-docs-contracts.sh"),

	"undocumented-passlist-var": conformanceFixture(".bench/gate.sh"),

	"line-binding-prose-drift": conformanceFixture(".bench/gate-line-contracts.sh"),
	"agent-hook-unwired":       conformanceFixture(".bench/gate-line-contracts.sh"),
	"stop-hook-unwired":        conformanceFixture(".bench/gate-line-contracts.sh"),
	"adapter-line-broken":      conformanceFixture(".bench/gate-line-contracts.sh"),

	"kit-only-asset-admitted":                conformanceFixture(".bench/gate-package-contracts.sh"),
	"kit-only-allowlist-emptied":             conformanceFixture(".bench/gate-package-contracts.sh"),
	"guard-describe-boundary-dropped":        conformanceFixture(".bench/gate-axi-contracts.sh"),
	"default-branch-refabricated":            conformanceFixture(".bench/gate.sh"),
	"guard-resolver-order-drift":             conformanceFixture(".bench/gate.sh"),
	"missing-license":                        conformanceFixture(".bench/gate.sh"),
	"mutable-workflow-action":                conformanceFixture(".bench/gate.sh"),
	"native-smoke-workflow-dropped":          conformanceFixture(".bench/gate.sh"),
	"native-reproducibility-handoff-dropped": conformanceFixture(".bench/gate.sh"),
	"offline-network-repair-allowed":         conformanceFixture(".bench/gate.sh"),
	"offline-stage-interruption-ignored":     conformanceFixture(".bench/gate.sh"),
	"offline-registry-fallback-allowed":      conformanceFixture(".bench/gate.sh"),
	"native-trigger-comment-spoof":           conformanceFixture(".bench/gate.sh"),
	"preflight-verify-gate-omitted":          conformanceFixture(".bench/gate.sh"),
	"preflight-verify-analysis-omitted":      conformanceFixture(".bench/gate.sh"),
	"preflight-verify-vulnerability-omitted": conformanceFixture(".bench/gate.sh"),
	"preflight-verify-artifact-omitted":      conformanceFixture(".bench/gate.sh"),
	"preflight-verify-smoke-omitted":         conformanceFixture(".bench/gate.sh"),
	"preflight-publish-identity-omitted":     conformanceFixture(".bench/gate.sh"),
	"preflight-publish-ancestry-omitted":     conformanceFixture(".bench/gate.sh"),
	"preflight-publish-changelog-omitted":    conformanceFixture(".bench/gate.sh"),
	"preflight-native-call-bypassed":         conformanceFixture(".bench/gate.sh"),
	"preflight-native-upload-bypassed":       conformanceFixture(".bench/gate.sh"),
	"preflight-release-call-bypassed":        conformanceFixture(".bench/gate.sh"),
	"preflight-publish-needs-bypassed":       conformanceFixture(".bench/gate.sh"),
	"preflight-publish-order-bypassed":       conformanceFixture(".bench/gate.sh"),
	"reproducibility-byte-compare-bypassed":  conformanceFixture(".bench/gate.sh"),
	"release-future-owner-omitted":           conformanceFixture(".bench/gate.sh"),
	"release-public-profile-omitted":         conformanceFixture(".bench/gate.sh"),
	"bounds-duplicate-owner":                 conformanceFixture(".bench/gate.sh"),
	"unrouted-subcommand":                    conformanceFixture(".bench/gate.sh"),
	"reintroduced-bare-skip":                 conformanceFixture(".bench/gate.sh"),
	"offline-slice1-operation-omitted":       conformanceFixture(".bench/gate.sh"),
}

// canaryFixtureFamilyRegistry assigns one owner to every fixture in a family whose
// checks share one Go implementation. Exact fixture registrations override this table.
var canaryFixtureFamilyRegistry = map[string]fixtureRegistration{
	"workflow-guidance-anchors": conformanceGoFixture(
		"internal/anchors/match.go",
		"internal/anchors/registry.go",
		"internal/anchors/registry_data.go",
		"internal/conformance/docs_workflow_helpers_test.go",
	),
	"decision-map-integrity": conformanceGoFixture(
		"internal/maps/schema.go",
		"internal/maps/validation.go",
		"internal/maps/tree_validation.go",
		"internal/conformance/checks_test.go",
	),
	"injected-ports": conformanceGoFixture(
		"internal/conformance/injected_ports_test.go",
		"internal/conformance/checks_test.go",
	),
	"guidance-prose-budgets": conformanceGoFixture(
		"internal/conformance/prose_budget_test.go",
		"internal/conformance/checks_test.go",
	),
	"roadmap-detail-integrity": conformanceGoFixture(
		"internal/roadmap/tree.go",
		"internal/roadmap/tree_validation.go",
		"internal/conformance/checks_test.go",
	),
}

func conformanceFixture(shellSources ...string) fixtureRegistration {
	return fixtureRegistration{Owner: ownerConformance, ShellSources: shellSources}
}

func conformanceGoFixture(goSources ...string) fixtureRegistration {
	return fixtureRegistration{Owner: ownerConformance, GoSources: goSources}
}

func fixtureRegistrationFor(name, family string) (fixtureRegistration, bool) {
	if registration, found := canaryFixtureRegistry[name]; found {
		return registration, true
	}
	registration, found := canaryFixtureFamilyRegistry[family]
	return registration, found
}

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
		// A family is canonical when inventory resolution can identify it: through a
		// family binding or through its own fixture marker. Only a family with neither
		// is unattributable; a flat fixture carries its own binding rather than a family.
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
