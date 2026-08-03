package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/gate"
)

type fixtureOwner string

const (
	ownerConformance fixtureOwner = "conformance"
	ownerBehavior    fixtureOwner = "behavior"
	// A phase fixture names no shell source whose retired EXPECT message could drift back in.
	ownerPhase fixtureOwner = "phase"
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
	"gate-input-gitignored":     conformanceFixture(".bench/gate.sh"),
	"shared-rule-drift":         conformanceFixture(".bench/gate.sh"),
	"readme-shared-rule-drift":  conformanceFixture(".bench/gate.sh"),

	"dangling-index":                conformanceFixture(".bench/gate.sh", ".bench/skills-index.sh"),
	"roadmap-promotion-persistence": conformanceFixture(".bench/gate.sh"),
	"missing-index-field":           conformanceFixture(".bench/gate.sh", ".bench/skills-index.sh"),
	"stale-index-wording":           conformanceFixture(".bench/gate.sh", ".bench/skills-index.sh"),
	"unindexed-skill":               conformanceFixture(".bench/gate.sh", ".bench/skills-index.sh"),

	"stale-command-reference":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"stale-codex-adapter-reference":    conformanceFixture(".bench/gate-docs-contracts.sh"),
	"retired-command-reference":        conformanceFixture(".bench/gate-docs-contracts.sh"),
	"stale-cli-doc-reference":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"stale-skill-cli-reference":        conformanceFixture(".bench/gate-docs-contracts.sh"),
	"missing-cli-inventory":            conformanceFixture(".bench/gate-docs-contracts.sh"),
	"historical-marker-prose":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"benchref-missing":                 conformanceFixture(".bench/gate-docs-contracts.sh"),
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
	"implement-spec-landing-commit":              conformanceFixture(".bench/gate-docs-contracts.sh"),
	"edge-inventory-anchor":                      conformanceFixture(".bench/gate-docs-contracts.sh"),
	"fix-pass-sentinel-anchor":                   conformanceFixture(".bench/gate-docs-contracts.sh"),
	"implement-spec-mandatory-delegation-anchor": conformanceFixture(".bench/gate-docs-contracts.sh"),
	"implement-spec-status-flip-anchor":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"implement-spec-structure-pointer":           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"what-next-anchor":                           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"what-next-spec-history-anchor":              conformanceFixture(".bench/gate-docs-contracts.sh"),
	"what-next-roadmap-context-anchor":           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"spec-retire-roadmap-row":                    conformanceFixture(".bench/gate-docs-contracts.sh"),
	"staged-command-sweep-anchor":                conformanceFixture(".bench/gate-docs-contracts.sh"),
	"capture-sink-anchor":                        conformanceFixture(".bench/gate-docs-contracts.sh"),
	"craft-seams-structure-headroom":             conformanceFixture(".bench/gate-docs-contracts.sh"),
	"review-persistence-anchor":                  conformanceFixture(".bench/gate-docs-contracts.sh"),
	"terminal-repair-bound-anchor":               conformanceFixture(".bench/gate-docs-contracts.sh"),
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
	"full-run-review-delegate-anchor":       conformanceFixture(".bench/gate-docs-contracts.sh"),
	"full-run-handoff-persistence-anchor":   conformanceFixture(".bench/gate-docs-contracts.sh"),
	"full-run-escalation-menu-anchor":       conformanceFixture(".bench/gate-docs-contracts.sh"),
	"full-run-silent-escalation-forbid":     conformanceFixture(".bench/gate-docs-contracts.sh"),
	"full-run-scope-fence-relocated":        conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-breakdown-step-anchor":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-light-path-anchor":              conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-stage-routing-anchor":           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-skill-contract-anchor":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-template-anchor":                conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-cross-pointers-anchor":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"ticket-gate-cadence-anchor":            conformanceFixture(".bench/gate-docs-contracts.sh"),

	"spec-build-initial-capacity-anchor":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"spec-build-unused-slot-reason-anchor":        conformanceFixture(".bench/gate-docs-contracts.sh"),
	"spec-build-review-input-binding-anchor":      conformanceFixture(".bench/gate-docs-contracts.sh"),
	"spec-build-final-check-single-author-anchor": conformanceFixture(".bench/gate-docs-contracts.sh"),
	"spec-build-review-route-anchor":              conformanceFixture(".bench/gate-docs-contracts.sh"),

	"ticket-decision-map-lifecycle-anchor":  conformanceFixture(".bench/gate-docs-contracts.sh"),
	"implementation-retro-authoring-anchor": conformanceFixture(".bench/gate-docs-contracts.sh"),
	"implementation-retro-drain-anchor":     conformanceFixture(".bench/gate-docs-contracts.sh"),

	"undocumented-passlist-var": conformanceFixture(".bench/gate.sh"),

	"line-binding-prose-drift": conformanceFixture(".bench/gate-line-contracts.sh"),
	"agent-hook-unwired":       conformanceFixture(".bench/gate-line-contracts.sh"),
	"agent-hook-broken":        conformanceFixture(".bench/gate-line-contracts.sh"),
	"stop-hook-unwired":        conformanceFixture(".bench/gate-line-contracts.sh"),
	"adapter-line-broken":      conformanceFixture(".bench/gate-line-contracts.sh"),

	"missing-files-entry":                    conformanceFixture(".bench/gate-package-contracts.sh"),
	"kit-only-asset-admitted":                conformanceFixture(".bench/gate-package-contracts.sh"),
	"kit-only-allowlist-emptied":             conformanceFixture(".bench/gate-package-contracts.sh"),
	"guard-describe-boundary-dropped":        conformanceFixture(".bench/gate-axi-contracts.sh"),
	"default-branch-refabricated":            conformanceFixture(".bench/gate.sh"),
	"guard-resolver-order-drift":             conformanceFixture(".bench/gate.sh"),
	"missing-license":                        conformanceFixture(".bench/gate.sh"),
	"mutable-workflow-action":                conformanceFixture(".bench/gate.sh"),
	"native-smoke-workflow-dropped":          conformanceFixture(".bench/gate.sh"),
	"native-proof-aggregation-bypassed":      behaviorFixture(),
	"native-proof-digest-binding-bypassed":   behaviorFixture(),
	"native-reproducibility-handoff-dropped": conformanceFixture(".bench/gate.sh"),
	"offline-archive-digest-binding-omitted": behaviorFixture(),
	"publication-order-bypass":               behaviorFixture(),
	"integrity-mismatch-acceptance":          behaviorFixture(),
	"premature-wrapper-promotion":            behaviorFixture(),
	"publication-unpublish-attempt":          behaviorFixture(),
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
	"release-package-evidence-omitted":       conformanceFixture(".bench/gate.sh"),
	"release-future-owner-omitted":           conformanceFixture(".bench/gate.sh"),
	"release-public-profile-omitted":         conformanceFixture(".bench/gate.sh"),
	"release-digest-binding-omitted":         conformanceFixture(".bench/gate.sh"),
	"bounds-duplicate-owner":                 conformanceFixture(".bench/gate.sh"),
	"bounds-duplicate-canary-width":          conformanceFixture(".bench/gate.sh"),
	"bounds-canary-width-unconsumed":         conformanceFixture(".bench/gate.sh"),
	"marker-wait-literal-deadline":           conformanceFixture(".bench/gate.sh"),
	"unrouted-subcommand":                    conformanceFixture(".bench/gate.sh"),
	"reintroduced-bare-skip":                 conformanceFixture(".bench/gate.sh"),
	"offline-slice1-operation-omitted":       conformanceFixture(".bench/gate.sh"),

	"go-build-broken":           phaseFixture(),
	"gofmt-unformatted":         phaseFixture(),
	"vet-printf-arg":            phaseFixture(),
	"go-test-failing":           phaseFixture(),
	"race-cleanup-test-failing": phaseFixture(),
	"conformance-suite-failing": phaseFixture(),

	"doctor-foreign-clobbered":             behaviorFixture(),
	"doctor-manager-dir-chosen":            behaviorFixture(),
	"doctor-stale-silent":                  behaviorFixture(),
	"postinstall-guard-bypassed":           behaviorFixture(),
	"postinstall-nonzero-exit":             behaviorFixture(),
	"session-start-advice-dropped":         behaviorFixture(),
	"wrapper-args-dropped":                 behaviorFixture(),
	"status-regressed":                     behaviorFixture(),
	"roadmap-regressed":                    behaviorFixture(),
	"unscaffolded-bench-file":              behaviorFixture(),
	"toon-escaping-dropped":                behaviorFixture(),
	"learnings-parse-broken":               behaviorFixture(),
	"guards-aggregation-dropped":           behaviorFixture(),
	"coverage-extraction-dropped":          behaviorFixture(),
	"diff-recorded-base-dropped":           behaviorFixture(),
	"roadmap-context-incomplete":           behaviorFixture(),
	"repo-local-forwarding-dropped":        behaviorFixture(),
	"native-selection-regressed":           behaviorFixture(),
	"wrapper-contamination-admitted":       behaviorFixture(),
	"wrapper-required-surface-dropped":     behaviorFixture(),
	"session-start-resume-cleanup-dropped": behaviorFixture(),
	"intent-common-dir-address-regressed":  behaviorFixture(),
	"status-landed-aggregation-regressed":  behaviorFixture(),
	"worktree-lifecycle-safety-bypassed":   behaviorFixture(),
	"gate-verdict-oracle-binding-bypassed": behaviorFixture(),
	"gate-verdict-invalidation-bypassed":   behaviorFixture(),
	"phase-manifest-defect-admitted":       behaviorFixture(),
}

// canaryFixtureFamilyRegistry assigns one owner to every fixture in a family whose
// checks share one Go implementation. Exact fixture registrations override this table.
var canaryFixtureFamilyRegistry = map[string]fixtureRegistration{
	"decision-map-integrity": conformanceGoFixture(
		"internal/maps/schema.go",
		"internal/maps/validation.go",
		"internal/maps/tree_validation.go",
		"internal/conformance/checks_test.go",
	),
	"example-agreement": conformanceGoFixture(
		"internal/conformance/example_agreement_test.go",
		"internal/conformance/checks_test.go",
	),
}

func conformanceFixture(shellSources ...string) fixtureRegistration {
	return fixtureRegistration{Owner: ownerConformance, ShellSources: shellSources}
}

func conformanceGoFixture(goSources ...string) fixtureRegistration {
	return fixtureRegistration{Owner: ownerConformance, GoSources: goSources}
}

func behaviorFixture() fixtureRegistration {
	return fixtureRegistration{Owner: ownerBehavior}
}

func phaseFixture() fixtureRegistration {
	return fixtureRegistration{Owner: ownerPhase}
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
		var wantOwner fixtureOwner
		switch phase := canary.FixturePhase(family); {
		// A legacy flat fixture has no single check or package owner, so it earns the full inner gate.
		case family == "", phase == canary.PhaseContract:
			wantOwner = ownerBehavior
		// A family routing to a phase of its own name is a phase family. The router owns
		// the phase names, so asking it beats listing them again here.
		case phase == family:
			wantOwner = ownerPhase
		case phase == "conformance" && familyIsBound(family):
			wantOwner = ownerConformance
		default:
			t.Errorf("canary fixture %q has unknown conformance family %q", name, family)
			continue
		}
		reg, ok := fixtureRegistrationFor(name, family)
		if !ok {
			t.Errorf("canary fixture %q is unclassified", name)
			continue
		}
		if reg.Owner != ownerConformance && reg.Owner != ownerBehavior && reg.Owner != ownerPhase {
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
		if reg.Owner != wantOwner {
			t.Errorf("canary fixture %q under %q has owner %q, want %q", name, family, reg.Owner, wantOwner)
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

	canaryFixtureRegistry["decision-map-override"] = behaviorFixture()
	t.Cleanup(func() { delete(canaryFixtureRegistry, "decision-map-override") })
	override, found := fixtureRegistrationFor("decision-map-override", "decision-map-integrity")
	if !found || override.Owner != ownerBehavior {
		t.Fatalf("exact fixture registration = %#v, %v; want exact override", override, found)
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

// fixtureExemptPhases are the phases whose bite is proved by a surface other than a
// canary fixture: each of these four drives fixtures rather than being one. Stating the
// exemptions rather than the obligations is what makes the inventory widen by itself —
// a phase added later must own a fixture or be exempted here on purpose, instead of
// quietly sitting outside a hand-written list.
var fixtureExemptPhases = map[string]bool{
	"conformance": true,
	"contract":    true,
	"shellcheck":  true,
	"canary":      true,
}

// TestEveryMovedStepOwnsAFixture reads the phase table and the fixture tree, and reds
// when a phase owns no fixture routed to it. A static phase-to-owner mapping would stay
// green while the fixture behind it matched nothing; only the tree side makes an
// orphaned phase visible.
func TestEveryMovedStepOwnsAFixture(t *testing.T) {
	h := NewHarness(t)
	covered := map[string]bool{}
	for _, fx := range canaryFixturePaths(t, h.KitPath("tests", "canary")) {
		covered[canary.FixturePhase(fx.Family)] = true
	}
	for _, phase := range gate.BenchkitPhases(h.KitRoot, h.KitRoot) {
		if fixtureExemptPhases[phase.Name] || covered[phase.Name] {
			continue
		}
		t.Errorf("gate phase %q owns no canary fixture; add one under tests/canary/%s/ so the step stays graded", phase.Name, phase.Name)
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
		// A family is canonical when the sweep can route it: to a phase (the behavior
		// family's contract phase, or a phase family named for its own phase) or to a
		// conformance check the registry binds. Only a family that routes to conformance
		// and is bound to nothing is unattributable — a legacy flat fixture is a fixture
		// in its own right rather than a family, and runs the full inner gate.
		if canary.IsConformanceFamily(filepath.Join(fixturesDir, family.Name())) && !familyIsBound(family.Name()) {
			t.Errorf("canary family %q is not canonical", family.Name())
		}
	}
	// The sweep's own walk enumerates the fixtures, base-name uniqueness included: a
	// second walk here would disagree with the tree the sweep actually runs.
	discovered, err := canary.Fixtures(fixturesDir)
	if err != nil {
		t.Fatalf("walk canary fixtures: %v", err)
	}
	return discovered
}
