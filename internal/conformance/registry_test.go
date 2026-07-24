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
	"shape-idea-bypass":                          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-bypass-wrapped":                  conformanceFixture(".bench/gate-docs-contracts.sh"),
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
	"shape-idea-handoff-anchor":                  conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-grill-continuation":              conformanceFixture(".bench/gate-docs-contracts.sh"),
	"story-line-anchor-missing":                  conformanceFixture(".bench/gate-docs-contracts.sh", ".bench/gate-line-contracts.sh"),
	"write-spec-handoff-anchor":                  conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-map-required":                    conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-reviewer-closed-fast-path":       conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-review-trigger-dropped":          conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-review-tier-escalated":           conformanceFixture(".bench/gate-docs-contracts.sh"),
	"line-anchor-missing":                        conformanceFixture(".bench/gate-line-contracts.sh"),
	"broken-coverage-map":                        conformanceFixture(".bench/gate-docs-contracts.sh"),

	"write-spec-reviewer-closed-comment-spoof":     conformanceFixture(".bench/gate-docs-contracts.sh"),
	"write-spec-open-fork-fallback":                conformanceFixture(".bench/gate-docs-contracts.sh"),
	"shape-idea-write-spec-entry-contract-pointer": conformanceFixture(".bench/gate-docs-contracts.sh"),

	"undocumented-passlist-var": conformanceFixture(".bench/gate.sh"),

	"line-binding-prose-drift": conformanceFixture(".bench/gate-line-contracts.sh"),
	"agent-hook-unwired":       conformanceFixture(".bench/gate-line-contracts.sh"),
	"agent-hook-broken":        conformanceFixture(".bench/gate-line-contracts.sh"),
	"stop-hook-unwired":        conformanceFixture(".bench/gate-line-contracts.sh"),
	"adapter-line-broken":      conformanceFixture(".bench/gate-line-contracts.sh"),

	"missing-files-entry":                    conformanceFixture(".bench/gate-package-contracts.sh"),
	"kit-only-asset-admitted":                conformanceFixture(".bench/gate-package-contracts.sh"),
	"kit-only-allowlist-emptied":             conformanceFixture(".bench/gate-package-contracts.sh"),
	"go-build-broken":                        conformanceFixture(".bench/gate-go-contracts.sh"),
	"go-test-failing":                        conformanceFixture(".bench/gate-go-contracts.sh"),
	"guard-describe-boundary-dropped":        conformanceFixture(".bench/gate-axi-contracts.sh"),
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
	"marker-wait-literal-deadline":           conformanceFixture(".bench/gate.sh"),
	"unrouted-subcommand":                    conformanceFixture(".bench/gate.sh"),
	"reintroduced-bare-skip":                 conformanceFixture(".bench/gate.sh"),
	"offline-slice1-operation-omitted":       conformanceFixture(".bench/gate.sh"),

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
		switch canary.FixturePhase(family) {
		case "contract":
			wantOwner = ownerBehavior
		case "conformance":
			if !isConformanceFamily(family) {
				t.Errorf("canary fixture %q has unknown conformance family %q", name, family)
				continue
			}
		default:
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
