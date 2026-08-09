package gate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	kitpayload "github.com/gibbonmi/bench"
	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestForcedAggregateRedRetiresEveryExecutedCheckSlot(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	writeGateTestFile(t, fixture.root, ".bench/phase-conformance.sh",
		"echo conformance >> .git/phase-runs\nif test -f .git/conformance-red; then exit 1; fi\n", 0o644)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})

	slots, valid := loadConformanceCheckSlots(fixture.root)
	if !valid {
		t.Fatal("green aggregate run left invalid check evidence")
	}
	if got, want := len(slots), len(ordinaryConformanceChecks(registry.Dev)); got != want {
		t.Fatalf("green aggregate run authored %d check slots, want %d", got, want)
	}

	writeGateTestFile(t, fixture.root, ".git/conformance-red", "red\n", 0o644)
	result := executeWithEngineAfterAcquire(context.Background(), fixture.root, io.Discard, io.Discard,
		productionGateEngine{}, nil, forceRun)
	if result.ActionExit == 0 {
		t.Fatal("forced aggregate run = green, want conformance red")
	}
	after, valid := loadConformanceCheckSlots(fixture.root)
	if !valid {
		t.Fatal("red aggregate run left invalid check evidence")
	}
	if len(after) != 0 {
		t.Fatalf("red aggregate run retained %d slots for checks it executed", len(after))
	}
}

func TestOrdinaryScopeInheritsValidCheckSlotsButFreshExecutesAll(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	resolution := Resolve(fixture.root, "", RealFS())

	ordinary := mustScopeComponents(t, fixture.root, resolution, reuseFreshGreen, time.Now().UTC())
	if got, want := len(ordinary.checks.Inherited), len(ordinaryConformanceChecks(registry.Dev)); got != want {
		t.Fatalf("ordinary scope inherited %d checks, want all %d valid slots", got, want)
	}
	if len(ordinary.checks.Executed) != 0 {
		t.Fatalf("ordinary scope executed %v, want every unchanged check inherited", ordinary.checks.Executed)
	}

	fresh := mustScopeComponents(t, fixture.root, resolution, forceRun, time.Now().UTC())
	if len(fresh.checks.Inherited) != 0 {
		t.Fatalf("fresh scope inherited %v, want none", fresh.checks.Inherited)
	}
	if got, want := len(fresh.checks.Executed), len(ordinaryConformanceChecks(registry.Dev)); got != want {
		t.Fatalf("fresh scope executed %d checks, want all %d", got, want)
	}
}

func TestDeclaredDocumentInputsInvalidateOwningChecks(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	fileSet, directorySet := map[string]bool{}, map[string]bool{}
	owners := map[string][]string{}
	for _, check := range registry.Checks {
		files, directories := declaredCheckInputPaths(check.Inputs)
		for _, file := range files {
			fileSet[file] = true
			owners[file] = append(owners[file], check.Name)
		}
		for _, directory := range directories {
			directorySet[directory] = true
			owners[directory] = append(owners[directory], check.Name)
		}
	}
	files, directories := mapKeys(fileSet), mapKeys(directorySet)
	slices.Sort(files)
	slices.Sort(directories)
	type baselineFile struct {
		data []byte
		mode os.FileMode
	}
	baselines := map[string]baselineFile{}
	for _, name := range files {
		path := filepath.Join(fixture.root, filepath.FromSlash(name))
		data, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			writeGateTestFile(t, fixture.root, name, "baseline\n", 0o644)
			data, err = os.ReadFile(path)
		}
		if err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		baselines[name] = baselineFile{data: data, mode: info.Mode().Perm()}
	}
	for _, dir := range directories {
		name := dir + "declared.md"
		writeGateTestFile(t, fixture.root, name, "baseline\n", 0o644)
		baselines[name] = baselineFile{data: []byte("baseline\n"), mode: 0o644}
	}
	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	resolution := Resolve(fixture.root, "", RealFS())

	for _, name := range append(files, directories...) {
		t.Run(name, func(t *testing.T) {
			path := name
			if strings.HasSuffix(name, "/") {
				path += "declared.md"
			}
			baseline := baselines[path]
			assertInvalidated := func(change string) {
				t.Helper()
				scoping := mustScopeComponents(t, fixture.root, resolution, reuseFreshGreen, time.Now().UTC())
				if !scoping.eligible {
					return
				}
				partition := scoping.checks
				if !slices.Contains(partition.Executed, "kit-compliance") {
					t.Fatalf("%s executed %v, want catch-all observer kit-compliance", change, partition.Executed)
				}
				for _, owner := range owners[name] {
					if !slices.Contains(partition.Executed, owner) {
						t.Fatalf("%s executed %v, want declared owner %s", change, partition.Executed, owner)
					}
				}
			}
			writeGateTestFile(t, fixture.root, path, string(baseline.data)+"moved\n", baseline.mode)
			assertInvalidated("document mutation")
			writeGateTestFile(t, fixture.root, path, string(baseline.data), baseline.mode)
			if err := os.Remove(filepath.Join(fixture.root, filepath.FromSlash(path))); err != nil {
				t.Fatal(err)
			}
			assertInvalidated("document deletion")
			writeGateTestFile(t, fixture.root, path, string(baseline.data), baseline.mode)
		})
	}
}

func TestPublicDocumentClassesProjectTheirExactCheckPartition(t *testing.T) {
	fixture := newKitShapedFixture(t)
	type decisionExpectation struct {
		eligible                  bool
		executedChecks            []string
		inheritedChecks           []string
		executedComponents        []string
		inheritedComponents       []string
		phases                    []string
		selectedChecksProjection  string
		inheritedChecksProjection string
		fullCanary                bool
		selectedCanaryFamilies    []string
	}
	seededDecision := decisionExpectation{
		eligible: true,
		inheritedChecks: []string{
			"kit-compliance", "canary-inner-compliance", "load-validity-metadata", "skills-index-command-adapters",
			"docs-currency-workflow", "gate-entry-contract", "ordinary-build-census", "offline-smoke-proof", "handoff-shape-single-source",
			"harness-prefix-single-source", "package-shipped-surface", "line-routing", "package-core-guard", "bench-sh-routes",
			"default-branch-single-source", "data-handling-derivation", "single-control-escaper", "bounds-policy",
			"marker-wait-deadlines", "subcommand-routing", "skip-ownership", "decision-map-integrity", "example-agreement",
			"component-honesty-prose", "contract-capture-reads", "injected-port-registry",
		},
		inheritedComponents: []string{"canary", "conformance-suite", "contract", "gofmt", "race", "shellcheck", "test", "vet"},
		phases:              []string{"conformance"},
		inheritedChecksProjection: "kit-compliance,canary-inner-compliance,load-validity-metadata,skills-index-command-adapters," +
			"docs-currency-workflow,gate-entry-contract,ordinary-build-census,offline-smoke-proof,handoff-shape-single-source,harness-prefix-single-source," +
			"package-shipped-surface,line-routing,package-core-guard,bench-sh-routes,default-branch-single-source," +
			"data-handling-derivation,single-control-escaper,bounds-policy,marker-wait-deadlines,subcommand-routing," +
			"skip-ownership,decision-map-integrity,example-agreement,component-honesty-prose,contract-capture-reads,injected-port-registry",
	}
	plainDocumentDecision := decisionExpectation{
		eligible: true,
		executedChecks: []string{
			"kit-compliance", "canary-inner-compliance", "load-validity-metadata", "skills-index-command-adapters",
			"docs-currency-workflow", "ordinary-build-census", "handoff-shape-single-source", "package-shipped-surface", "line-routing",
			"package-core-guard", "bounds-policy", "example-agreement",
		},
		inheritedChecks: []string{
			"gate-entry-contract", "offline-smoke-proof", "harness-prefix-single-source", "bench-sh-routes",
			"default-branch-single-source", "data-handling-derivation", "single-control-escaper", "marker-wait-deadlines",
			"subcommand-routing", "skip-ownership", "decision-map-integrity", "component-honesty-prose",
			"contract-capture-reads", "injected-port-registry",
		},
		inheritedComponents:      []string{"canary", "conformance-suite", "contract", "gofmt", "race", "shellcheck", "test", "vet"},
		phases:                   []string{"conformance"},
		selectedChecksProjection: "kit-compliance,canary-inner-compliance,load-validity-metadata,skills-index-command-adapters,docs-currency-workflow,ordinary-build-census,handoff-shape-single-source,package-shipped-surface,line-routing,package-core-guard,bounds-policy,example-agreement",
		inheritedChecksProjection: "gate-entry-contract,offline-smoke-proof,harness-prefix-single-source,bench-sh-routes," +
			"default-branch-single-source,data-handling-derivation,single-control-escaper,marker-wait-deadlines," +
			"subcommand-routing,skip-ownership,decision-map-integrity,component-honesty-prose,contract-capture-reads,injected-port-registry",
	}
	contractDocumentDecision := decisionExpectation{
		eligible: true,
		executedChecks: []string{
			"kit-compliance", "canary-inner-compliance", "load-validity-metadata", "skills-index-command-adapters",
			"docs-currency-workflow", "ordinary-build-census", "handoff-shape-single-source", "package-shipped-surface", "line-routing",
			"package-core-guard", "bounds-policy", "example-agreement",
		},
		inheritedChecks: []string{
			"gate-entry-contract", "offline-smoke-proof", "harness-prefix-single-source", "bench-sh-routes",
			"default-branch-single-source", "data-handling-derivation", "single-control-escaper", "marker-wait-deadlines",
			"subcommand-routing", "skip-ownership", "decision-map-integrity", "component-honesty-prose",
			"contract-capture-reads", "injected-port-registry",
		},
		executedComponents:       []string{"contract"},
		inheritedComponents:      []string{"canary", "conformance-suite", "gofmt", "race", "shellcheck", "test", "vet"},
		phases:                   []string{"conformance", "contract"},
		selectedChecksProjection: "kit-compliance,canary-inner-compliance,load-validity-metadata,skills-index-command-adapters,docs-currency-workflow,ordinary-build-census,handoff-shape-single-source,package-shipped-surface,line-routing,package-core-guard,bounds-policy,example-agreement",
		inheritedChecksProjection: "gate-entry-contract,offline-smoke-proof,harness-prefix-single-source,bench-sh-routes," +
			"default-branch-single-source,data-handling-derivation,single-control-escaper,marker-wait-deadlines," +
			"subcommand-routing,skip-ownership,decision-map-integrity,component-honesty-prose,contract-capture-reads,injected-port-registry",
	}
	canaryAndContractDocumentDecision := decisionExpectation{
		eligible: true,
		executedChecks: []string{
			"kit-compliance", "canary-inner-compliance", "load-validity-metadata", "skills-index-command-adapters",
			"docs-currency-workflow", "ordinary-build-census", "handoff-shape-single-source", "package-shipped-surface", "line-routing",
			"package-core-guard", "bounds-policy", "example-agreement",
		},
		inheritedChecks: []string{
			"gate-entry-contract", "offline-smoke-proof", "harness-prefix-single-source", "bench-sh-routes",
			"default-branch-single-source", "data-handling-derivation", "single-control-escaper", "marker-wait-deadlines",
			"subcommand-routing", "skip-ownership", "decision-map-integrity", "component-honesty-prose",
			"contract-capture-reads", "injected-port-registry",
		},
		executedComponents:       []string{"canary", "contract"},
		inheritedComponents:      []string{"conformance-suite", "gofmt", "race", "shellcheck", "test", "vet"},
		phases:                   []string{"canary", "conformance", "contract"},
		selectedChecksProjection: "kit-compliance,canary-inner-compliance,load-validity-metadata,skills-index-command-adapters,docs-currency-workflow,ordinary-build-census,handoff-shape-single-source,package-shipped-surface,line-routing,package-core-guard,bounds-policy,example-agreement",
		inheritedChecksProjection: "gate-entry-contract,offline-smoke-proof,harness-prefix-single-source,bench-sh-routes," +
			"default-branch-single-source,data-handling-derivation,single-control-escaper,marker-wait-deadlines," +
			"subcommand-routing,skip-ownership,decision-map-integrity,component-honesty-prose,contract-capture-reads,injected-port-registry",
	}
	dataHandlingDecision := decisionExpectation{
		eligible: true,
		executedChecks: []string{
			"kit-compliance", "canary-inner-compliance", "load-validity-metadata", "skills-index-command-adapters",
			"docs-currency-workflow", "ordinary-build-census", "handoff-shape-single-source", "package-shipped-surface", "line-routing",
			"package-core-guard", "data-handling-derivation", "bounds-policy", "example-agreement",
		},
		inheritedChecks: []string{
			"gate-entry-contract", "offline-smoke-proof", "harness-prefix-single-source", "bench-sh-routes",
			"default-branch-single-source", "single-control-escaper", "marker-wait-deadlines", "subcommand-routing",
			"skip-ownership", "decision-map-integrity", "component-honesty-prose", "contract-capture-reads", "injected-port-registry",
		},
		inheritedComponents:      []string{"canary", "conformance-suite", "contract", "gofmt", "race", "shellcheck", "test", "vet"},
		phases:                   []string{"conformance"},
		selectedChecksProjection: "kit-compliance,canary-inner-compliance,load-validity-metadata,skills-index-command-adapters,docs-currency-workflow,ordinary-build-census,handoff-shape-single-source,package-shipped-surface,line-routing,package-core-guard,data-handling-derivation,bounds-policy,example-agreement",
		inheritedChecksProjection: "gate-entry-contract,offline-smoke-proof,harness-prefix-single-source,bench-sh-routes," +
			"default-branch-single-source,single-control-escaper,marker-wait-deadlines,subcommand-routing,skip-ownership," +
			"decision-map-integrity,component-honesty-prose,contract-capture-reads,injected-port-registry",
	}
	benchkitProfileDecision := decisionExpectation{
		eligible: true,
		executedChecks: []string{
			"kit-compliance", "canary-inner-compliance", "load-validity-metadata", "skills-index-command-adapters",
			"docs-currency-workflow", "ordinary-build-census", "handoff-shape-single-source", "package-shipped-surface", "line-routing",
			"package-core-guard", "bounds-policy", "example-agreement", "component-honesty-prose",
		},
		inheritedChecks: []string{
			"gate-entry-contract", "offline-smoke-proof", "harness-prefix-single-source", "bench-sh-routes",
			"default-branch-single-source", "data-handling-derivation", "single-control-escaper", "marker-wait-deadlines",
			"subcommand-routing", "skip-ownership", "decision-map-integrity", "contract-capture-reads", "injected-port-registry",
		},
		inheritedComponents:      []string{"canary", "conformance-suite", "contract", "gofmt", "race", "shellcheck", "test", "vet"},
		phases:                   []string{"conformance"},
		selectedChecksProjection: "kit-compliance,canary-inner-compliance,load-validity-metadata,skills-index-command-adapters,docs-currency-workflow,ordinary-build-census,handoff-shape-single-source,package-shipped-surface,line-routing,package-core-guard,bounds-policy,example-agreement,component-honesty-prose",
		inheritedChecksProjection: "gate-entry-contract,offline-smoke-proof,harness-prefix-single-source,bench-sh-routes," +
			"default-branch-single-source,data-handling-derivation,single-control-escaper,marker-wait-deadlines," +
			"subcommand-routing,skip-ownership,decision-map-integrity,contract-capture-reads,injected-port-registry",
	}
	decisionDocumentDecision := decisionExpectation{
		eligible: true,
		executedChecks: []string{
			"kit-compliance", "canary-inner-compliance", "load-validity-metadata", "skills-index-command-adapters",
			"docs-currency-workflow", "ordinary-build-census", "handoff-shape-single-source", "package-shipped-surface", "line-routing",
			"package-core-guard", "bounds-policy", "decision-map-integrity", "example-agreement",
		},
		inheritedChecks: []string{
			"gate-entry-contract", "offline-smoke-proof", "harness-prefix-single-source", "bench-sh-routes",
			"default-branch-single-source", "data-handling-derivation", "single-control-escaper", "marker-wait-deadlines",
			"subcommand-routing", "skip-ownership", "component-honesty-prose", "contract-capture-reads", "injected-port-registry",
		},
		inheritedComponents:      []string{"canary", "conformance-suite", "contract", "gofmt", "race", "shellcheck", "test", "vet"},
		phases:                   []string{"conformance"},
		selectedChecksProjection: "kit-compliance,canary-inner-compliance,load-validity-metadata,skills-index-command-adapters,docs-currency-workflow,ordinary-build-census,handoff-shape-single-source,package-shipped-surface,line-routing,package-core-guard,bounds-policy,decision-map-integrity,example-agreement",
		inheritedChecksProjection: "gate-entry-contract,offline-smoke-proof,harness-prefix-single-source,bench-sh-routes," +
			"default-branch-single-source,data-handling-derivation,single-control-escaper,marker-wait-deadlines," +
			"subcommand-routing,skip-ownership,component-honesty-prose,contract-capture-reads,injected-port-registry",
	}
	allChecksExecuteDecision := decisionExpectation{
		eligible: true,
		executedChecks: []string{
			"kit-compliance", "canary-inner-compliance", "load-validity-metadata", "skills-index-command-adapters",
			"docs-currency-workflow", "gate-entry-contract", "ordinary-build-census", "offline-smoke-proof", "handoff-shape-single-source",
			"harness-prefix-single-source", "package-shipped-surface", "line-routing", "package-core-guard", "bench-sh-routes",
			"default-branch-single-source", "data-handling-derivation", "single-control-escaper", "bounds-policy",
			"marker-wait-deadlines", "subcommand-routing", "skip-ownership", "decision-map-integrity", "example-agreement",
			"component-honesty-prose", "contract-capture-reads", "injected-port-registry",
		},
	}
	fileCases := []struct {
		name               string
		mutation, deletion decisionExpectation
	}{
		{name: ".bench-notes.md", mutation: plainDocumentDecision, deletion: plainDocumentDecision},
		{name: ".bench/BENCH-reference.md", mutation: contractDocumentDecision, deletion: allChecksExecuteDecision},
		{name: ".bench/BENCH.md", mutation: contractDocumentDecision, deletion: allChecksExecuteDecision},
		{name: ".claude/README.md", mutation: contractDocumentDecision, deletion: allChecksExecuteDecision},
		{name: "CHANGELOG.md", mutation: contractDocumentDecision, deletion: allChecksExecuteDecision},
		{name: "DATA_HANDLING.md", mutation: dataHandlingDecision, deletion: dataHandlingDecision},
		{name: "README.md", mutation: contractDocumentDecision, deletion: allChecksExecuteDecision},
		{name: "ROADMAP.md", mutation: plainDocumentDecision, deletion: plainDocumentDecision},
		{name: "projects/benchkit.md", mutation: benchkitProfileDecision, deletion: benchkitProfileDecision},
		{name: "projects/gl-axi.md", mutation: contractDocumentDecision, deletion: allChecksExecuteDecision},
		{name: "projects/regroup.md", mutation: contractDocumentDecision, deletion: allChecksExecuteDecision},
	}
	directoryCases := []struct {
		name              string
		addition, removal decisionExpectation
	}{
		{name: ".agents/commands/", addition: canaryAndContractDocumentDecision, removal: canaryAndContractDocumentDecision},
		{name: ".agents/skills/", addition: canaryAndContractDocumentDecision, removal: canaryAndContractDocumentDecision},
		{name: ".bench/adapters/", addition: contractDocumentDecision, removal: contractDocumentDecision},
		{name: ".bench/hooks/", addition: contractDocumentDecision, removal: contractDocumentDecision},
		{name: ".bench/lib/", addition: contractDocumentDecision, removal: contractDocumentDecision},
		{name: "capture/", addition: plainDocumentDecision, removal: plainDocumentDecision},
		{name: "decisions/", addition: decisionDocumentDecision, removal: decisionDocumentDecision},
		{name: "specs/", addition: decisionDocumentDecision, removal: decisionDocumentDecision},
	}
	files := make([]string, 0, len(fileCases))
	for _, test := range fileCases {
		files = append(files, test.name)
	}
	directories := make([]string, 0, len(directoryCases))
	for _, test := range directoryCases {
		directories = append(directories, test.name)
	}
	declaredFiles, declaredDirectories := declaredPublicDocumentClasses(t, fixture.root)
	if !slices.Equal(files, declaredFiles) || !slices.Equal(directories, declaredDirectories) {
		t.Fatalf("literal public document membership = files %v directories %v, want production files %v directories %v", files, directories, declaredFiles, declaredDirectories)
	}
	for _, name := range files {
		path := filepath.Join(fixture.root, filepath.FromSlash(name))
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			writeGateTestFile(t, fixture.root, name, "declared public document\n", 0o644)
		} else if err != nil {
			t.Fatal(err)
		}
	}
	mustExecuteGreen(t, fixture.root, productionGateEngine{})

	seededSlots := checkSlotTestStoreBytes(t, fixture.root)
	seededFullRuns := fullRunCount(t, fixture.root)
	seededPhaseRuns := phaseRunNames(t, fixture.root)
	seededVerdict, err := os.ReadFile(cachePath(t, fixture.root))
	if err != nil {
		t.Fatal(err)
	}
	resolution := Resolve(fixture.root, "", RealFS())
	recorder := installGateGitRecorder(t)
	recordedCaptures := 0
	assertChange := func(name string, want decisionExpectation, change func()) {
		t.Helper()
		change()
		generation, err := captureWorkingTree(fixture.root)
		if err != nil {
			t.Fatalf("captureWorkingTree = %v, want the fixture generation", err)
		}
		decision := scopeComponentsForGeneration(fixture.root, resolution, reuseFreshGreen, time.Now().UTC(), generation)
		if decision.eligible != want.eligible {
			t.Fatalf("%s eligibility = %t, want %t", name, decision.eligible, want.eligible)
		}
		if !slices.Equal(decision.checks.Executed, want.executedChecks) {
			t.Fatalf("%s executed checks = %v, want %v", name, decision.checks.Executed, want.executedChecks)
		}
		if got := decision.checks.verdictInherited(); !slices.Equal(got, want.inheritedChecks) {
			t.Fatalf("%s inherited checks = %v, want %v", name, got, want.inheritedChecks)
		}
		if got := decision.executedPhaseNames(); !slices.Equal(got, want.phases) {
			t.Fatalf("%s executed phases = %v, want %v", name, got, want.phases)
		}
		if got := decision.executedScopedComponents(); !slices.Equal(got, want.executedComponents) {
			t.Fatalf("%s executed components = %v, want %v", name, got, want.executedComponents)
		}
		gotInheritedComponents := make([]string, 0, len(decision.skipped))
		for _, skipped := range decision.skipped {
			gotInheritedComponents = append(gotInheritedComponents, skipped.Component)
		}
		if !slices.Equal(gotInheritedComponents, want.inheritedComponents) {
			t.Fatalf("%s inherited components = %v, want %v", name, gotInheritedComponents, want.inheritedComponents)
		}
		if decision.checks.CanaryFull != want.fullCanary || !slices.Equal(decision.checks.CanaryFamilies, want.selectedCanaryFamilies) {
			t.Fatalf("%s canary check projection = full %t families %v, want full %t families %v", name, decision.checks.CanaryFull, decision.checks.CanaryFamilies, want.fullCanary, want.selectedCanaryFamilies)
		}
		selectedProjection, inheritedProjection := "", ""
		if phase, found := phaseNamed(decision.phases, "conformance"); found {
			selectedProjection = phaseEnvValue(phase.Env, registry.ConformanceChecksEnv)
			inheritedProjection = phaseEnvValue(phase.Env, registry.ConformanceInheritedEnv)
		}
		if selectedProjection != want.selectedChecksProjection || inheritedProjection != want.inheritedChecksProjection {
			t.Fatalf("%s conformance projection = selected %q inherited %q, want selected %q inherited %q", name, selectedProjection, inheritedProjection, want.selectedChecksProjection, want.inheritedChecksProjection)
		}
		if got := checkSlotTestStoreBytes(t, fixture.root); !reflect.DeepEqual(got, seededSlots) {
			t.Fatalf("%s changed seeded check-slot bytes", name)
		}
		if got := fullRunCount(t, fixture.root); got != seededFullRuns {
			t.Fatalf("%s ran the full engine %d times after the seed, want %d", name, got, seededFullRuns)
		}
		if got := phaseRunNames(t, fixture.root); !slices.Equal(got, seededPhaseRuns) {
			t.Fatalf("%s ran phases %v after the seed, want %v", name, got, seededPhaseRuns)
		}
		if got, err := os.ReadFile(cachePath(t, fixture.root)); err != nil || !bytes.Equal(got, seededVerdict) {
			t.Fatalf("%s changed the seeded verdict bytes: %v", name, err)
		}
		captures := countGateGitOperation(recorder.operations(t), "write-tree")
		if captures-recordedCaptures != 1 {
			t.Fatalf("%s captured %d generations, want one", name, captures-recordedCaptures)
		}
		recordedCaptures = captures
	}
	restore := func(name string, data []byte, mode os.FileMode) {
		t.Helper()
		writeGateTestFile(t, fixture.root, name, string(data), mode)
	}

	for _, test := range fileCases {
		t.Run(test.name, func(t *testing.T) {
			name := test.name
			path := filepath.Join(fixture.root, filepath.FromSlash(name))
			baseline, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			info, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			assertChange("mutating "+name, test.mutation, func() {
				writeGateTestFile(t, fixture.root, name, string(baseline)+"matrix mutation\n", info.Mode().Perm())
			})
			restore(name, baseline, info.Mode().Perm())
			assertChange("deleting "+name, test.deletion, func() {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			})
			restore(name, baseline, info.Mode().Perm())
			assertChange("restoring "+name, seededDecision, func() {})
		})
	}
	for _, test := range directoryCases {
		t.Run(test.name, func(t *testing.T) {
			directory := test.name
			removalName, removalData, removalMode := firstDocumentDescendant(t, fixture.root, directory)
			addedName := directory + "public-partition-matrix.md"
			addedPath := filepath.Join(fixture.root, filepath.FromSlash(addedName))
			assertChange("adding a descendant of "+directory, test.addition, func() {
				writeGateTestFile(t, fixture.root, addedName, "matrix descendant\n", 0o644)
			})
			if err := os.Remove(addedPath); err != nil {
				t.Fatal(err)
			}
			assertChange("removing a descendant of "+directory, test.removal, func() {
				if err := os.Remove(filepath.Join(fixture.root, filepath.FromSlash(removalName))); err != nil {
					t.Fatal(err)
				}
			})
			restore(removalName, removalData, removalMode)
			assertChange("restoring "+directory, seededDecision, func() {})
		})
	}
}

func firstDocumentDescendant(t *testing.T, root, directory string) (string, []byte, os.FileMode) {
	t.Helper()
	var name string
	err := filepath.WalkDir(filepath.Join(root, filepath.FromSlash(directory)), func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() && strings.HasSuffix(entry.Name(), ".md") {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			name = filepath.ToSlash(rel)
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if name == "" {
		t.Fatalf("declared public document directory %s has no fixture document to remove", directory)
	}
	path := filepath.Join(root, filepath.FromSlash(name))
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return name, data, info.Mode().Perm()
}

func declaredPublicDocumentClasses(t *testing.T, root string) (files, directories []string) {
	t.Helper()
	fileSet, directorySet := map[string]bool{}, map[string]bool{}
	addFile := func(name string) {
		if strings.HasSuffix(name, ".md") {
			fileSet[name] = true
		}
	}
	addDirectory := func(name string) {
		directorySet[strings.TrimSuffix(name, "/")+"/"] = true
	}
	scope := ReducedScope()
	for _, name := range scope.Files() {
		addFile(name)
	}
	for _, name := range scope.Directories() {
		addDirectory(name)
	}
	for _, check := range registry.Checks {
		declaredFiles, declaredDirectories := declaredCheckInputPaths(check.Inputs)
		for _, name := range declaredFiles {
			addFile(name)
		}
		for _, name := range declaredDirectories {
			addDirectory(name)
		}
	}
	payload, err := os.ReadFile(filepath.Join(root, ".bench", "consumer-payload.json"))
	if err != nil {
		t.Fatal(err)
	}
	var rows []kitpayload.PayloadRow
	if err := json.Unmarshal(payload, &rows); err != nil {
		t.Fatal(err)
	}
	for _, row := range kitpayload.PayloadConsumerRows(rows) {
		if row.Tree {
			addDirectory(row.Source)
		} else {
			addFile(row.Source)
		}
	}
	files, directories = mapKeys(fileSet), mapKeys(directorySet)
	slices.Sort(files)
	slices.Sort(directories)
	return files, directories
}

func TestConformanceImplementationSelectsOwningOrAllCanaryFamilies(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name       string
		mutate     func(string)
		want       []string
		wantScoped bool
	}{
		{
			name: "bound implementation",
			mutate: func(root string) {
				writeGateTestFile(t, root, "internal/conformance/owners.go", conformanceCanarySources("checkLineRouting"), 0o644)
			},
			want:       registry.CanaryFamilies("line-routing"),
			wantScoped: true,
		},
		{
			name: "shared helper",
			mutate: func(root string) {
				writeGateTestFile(t, root, "internal/conformance/shared.go", "package conformance\n\nfunc sharedHelper() int { return 2 }\n", 0o644)
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newKitShapedFixture(t)
			writeGateTestFile(t, fixture.root, "internal/conformance/owners.go", conformanceCanarySources(""), 0o644)
			writeGateTestFile(t, fixture.root, "internal/conformance/shared.go", "package conformance\n\nfunc sharedHelper() int { return 1 }\n", 0o644)
			mustExecuteGreen(t, fixture.root, productionGateEngine{})
			test.mutate(fixture.root)

			scoping := mustScopeComponents(t, fixture.root, Resolve(fixture.root, "", RealFS()), reuseFreshGreen, time.Now().UTC())
			phase, found := phaseNamed(scoping.phases, "canary")
			if !found {
				t.Fatal("conformance implementation change inherited the canary component")
			}
			got := phaseEnvValue(phase.Env, canary.FamilySelectionEnv)
			if test.wantScoped {
				if got != strings.Join(test.want, ",") || phaseEnvValue(phase.Env, canary.FamilySelectionOwnerEnv) != "gate" {
					t.Fatalf("canary selection = %q under owner %q, want %v under gate", got, phaseEnvValue(phase.Env, canary.FamilySelectionOwnerEnv), test.want)
				}
			} else if got != "" {
				t.Fatalf("shared helper selected %q, want the full conformance canary set", got)
			}
		})
	}
}

func TestUnresolvedCurrentCanaryIdentityWithLegacySlotsRunsFullCanary(t *testing.T) {
	t.Parallel()
	for _, resolutionErr := range []error{nil, errors.New("canary identity unresolved")} {
		root, identities := seededCheckSlotStore(t)
		partition := partitionConformanceChecks(root, registry.Dev, identities, nil, time.Unix(20, 0))
		if len(partition.Inherited) == 0 {
			t.Fatal("legacy slots inherited no checks, so the canary fallback precondition is absent")
		}
		decorateConformanceCanarySelection(root, &partition, conformanceCanaryIdentities{}, resolutionErr)
		if !partition.CanaryFull || len(partition.CanaryFamilies) != 0 {
			t.Fatalf("empty current identity with error %v selected full=%t families=%v, want full canary", resolutionErr, partition.CanaryFull, partition.CanaryFamilies)
		}
	}
}

func conformanceCanarySources(changed string) string {
	var source strings.Builder
	source.WriteString("package conformance\n")
	seen := map[string]bool{}
	for _, check := range registry.Checks {
		if check.Meta || len(registry.CanaryFamilies(check.Name)) == 0 || seen[check.Implementation] {
			continue
		}
		seen[check.Implementation] = true
		body := "{}"
		if check.Implementation == changed {
			body = "{ _ = 1 }"
		}
		source.WriteString("\nfunc " + check.Implementation + "() " + body + "\n")
	}
	return source.String()
}

func phaseEnvValue(env []string, name string) string {
	for _, value := range env {
		if value, found := strings.CutPrefix(value, name+"="); found {
			return value
		}
	}
	return ""
}

func TestMixedCheckPartitionProjectsExactOutputAndVerdict(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	writeGateTestFile(t, fixture.root, ".bench/phase-conformance.sh",
		"echo conformance >> .git/phase-runs\nprintf '%s\\n' \"$BENCH_CONFORMANCE_CHECKS\" > .git/selected-checks\nprintf '%s\\n' \"$BENCH_CONFORMANCE_INHERITED\" > .git/inherited-checks\n", 0o644)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	beforeConformance := phaseNameCount(phaseRunNames(t, fixture.root), "conformance")
	writeGateTestFile(t, fixture.root, "ROADMAP.md", "moved\n", 0o644)

	observation := observeGreenGate(t, fixture.root)
	record := slotRecord(t, fixture.root, time.Now().UTC())
	if len(record.CheckExecuted) == 0 || len(record.CheckInherited) == 0 {
		t.Fatalf("mixed verdict check partition = executed %v inherited %v", record.CheckExecuted, record.CheckInherited)
	}
	if len(record.CheckEvidence) != len(record.CheckInherited) {
		t.Fatalf("mixed verdict check evidence = %d entries for %d inherited checks", len(record.CheckEvidence), len(record.CheckInherited))
	}
	for _, name := range record.CheckExecuted {
		if !strings.Contains(observation.stdout, "conformance check "+name+": executing") {
			t.Errorf("output omitted executed check %s:\n%s", name, observation.stdout)
		}
	}
	for _, name := range record.CheckInherited {
		evidence := record.CheckEvidence[name]
		if evidence.Identity == "" || evidence.AuthoredAt == "" {
			t.Errorf("inherited check %s evidence = %+v", name, evidence)
		}
		if !strings.Contains(observation.stdout, "conformance check "+name+": inherited") {
			t.Errorf("output omitted inherited check %s:\n%s", name, observation.stdout)
		}
	}
	inspection := Inspect(fixture.root)
	if inspection.CheckPartition == nil || inspection.ReusableGreen || inspection.Reason != "partial verdict" {
		t.Fatalf("mixed inspection = %+v", inspection)
	}
	if !composedGreenAtKit(fixture.root, fixture.root) {
		t.Fatal("mixed check verdict did not compose to the one landing green")
	}
	if !slices.Equal(inspection.CheckPartition.Executed, record.CheckExecuted) {
		t.Fatalf("projected executed checks = %v, want %v", inspection.CheckPartition.Executed, record.CheckExecuted)
	}
	if got := len(inspection.CheckPartition.Inherited); got != len(record.CheckInherited) {
		t.Fatalf("projected inherited checks = %d, want %d", got, len(record.CheckInherited))
	}
	selected, err := os.ReadFile(filepath.Join(fixture.root, ".git", "selected-checks"))
	if err != nil {
		t.Fatal(err)
	}
	wantSelected := strings.Join(record.CheckExecuted[len(requiredMetaChecksForGateTest()):], ",") + "\n"
	if string(selected) != wantSelected {
		t.Fatalf("aggregate selected checks = %q, want %q", selected, wantSelected)
	}
	inherited, err := os.ReadFile(filepath.Join(fixture.root, ".git", "inherited-checks"))
	if err != nil {
		t.Fatal(err)
	}
	if wantInherited := strings.Join(record.CheckInherited, ",") + "\n"; string(inherited) != wantInherited {
		t.Fatalf("aggregate inherited checks = %q, want %q", inherited, wantInherited)
	}
	if got := phaseNameCount(phaseRunNames(t, fixture.root), "conformance") - beforeConformance; got != 1 {
		t.Fatalf("mixed gate made %d aggregate conformance invocations, want one", got)
	}
}

func TestCheckOnlyPartialGateProjectsInspectionAndOutput(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	evidenceDir, err := componentSlotDir(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(evidenceDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Name() == conformanceCheckSlotStoreName {
			continue
		}
		if err := os.Remove(filepath.Join(evidenceDir, entry.Name())); err != nil {
			t.Fatal(err)
		}
	}
	writeGateTestFile(t, fixture.root, "DATA_HANDLING.md", "moved\n", 0o644)

	observation := observeGreenGate(t, fixture.root)
	record := slotRecord(t, fixture.root, time.Now().UTC())
	if record.partitions() || !record.checkPartitions() {
		t.Fatalf("check-only verdict carries component=%v check=%v: %+v", record.partitions(), record.checkPartitions(), record)
	}
	inspection := Inspect(fixture.root)
	if inspection.Partition != nil || inspection.CheckPartition == nil || inspection.Reason != "partial verdict" || inspection.ReusableGreen {
		t.Fatalf("check-only inspection = %+v", inspection)
	}
	for _, name := range inspection.CheckPartition.Executed {
		if !strings.Contains(observation.stdout, "conformance check "+name+": executing") {
			t.Errorf("check-only output omitted executed check %s:\n%s", name, observation.stdout)
		}
	}
	for _, inherited := range inspection.CheckPartition.Inherited {
		if !strings.Contains(observation.stdout, "conformance check "+inherited.Component+": inherited") {
			t.Errorf("check-only output omitted inherited check %s:\n%s", inherited.Component, observation.stdout)
		}
	}
}

func phaseNameCount(names []string, want string) int {
	count := 0
	for _, name := range names {
		if name == want {
			count++
		}
	}
	return count
}

func requiredMetaChecksForGateTest() []string {
	var names []string
	for _, check := range registry.Checks {
		if check.Meta {
			names = append(names, check.Name)
		}
	}
	return names
}

func TestMixedGreenAuthorsOnlyExecutedCheckSlots(t *testing.T) {
	t.Parallel()
	root := checkSlotTestRepo(t)
	identities := checkSlotTestIdentities()
	first := conformanceCheckPartition{
		Tier:       registry.Dev,
		Identities: identities,
		Executed:   []string{"line-routing", "package-core-guard", "subcommand-routing"},
	}
	if err := applyConformanceCheckOutcome(root, first, checkRunGreen, time.Unix(10, 0)); err != nil {
		t.Fatal(err)
	}
	before := checkSlotTestBytes(t, root)

	mixed := conformanceCheckPartition{
		Tier:       registry.Dev,
		Identities: identities,
		Executed:   []string{"line-routing"},
		Inherited: []conformanceCheckSkip{
			{Check: "package-core-guard", Identity: identities["package-core-guard"]},
			{Check: "subcommand-routing", Identity: identities["subcommand-routing"]},
		},
	}
	if err := applyConformanceCheckOutcome(root, mixed, checkRunGreen, time.Unix(20, 0)); err != nil {
		t.Fatal(err)
	}
	after := checkSlotTestBytes(t, root)
	if bytes.Equal(after["line-routing"], before["line-routing"]) {
		t.Fatal("executed check slot was not re-authored")
	}
	for _, inherited := range []string{"package-core-guard", "subcommand-routing"} {
		if !bytes.Equal(after[inherited], before[inherited]) {
			t.Fatalf("inherited check %s slot bytes changed", inherited)
		}
	}
}

func TestRedRetiresOnlyExecutedCheckSlots(t *testing.T) {
	t.Parallel()
	root := checkSlotTestRepo(t)
	identities := checkSlotTestIdentities()
	seed := conformanceCheckPartition{
		Tier:       registry.Dev,
		Identities: identities,
		Executed:   []string{"line-routing", "package-core-guard", "subcommand-routing"},
	}
	if err := applyConformanceCheckOutcome(root, seed, checkRunGreen, time.Unix(10, 0)); err != nil {
		t.Fatal(err)
	}
	before := checkSlotTestBytes(t, root)
	outcome := conformanceCheckPartition{
		Tier:       registry.Dev,
		Identities: identities,
		Executed:   []string{"line-routing", "subcommand-routing"},
		Inherited:  []conformanceCheckSkip{{Check: "package-core-guard", Identity: identities["package-core-guard"]}},
	}
	if err := applyConformanceCheckOutcome(root, outcome, checkRunRed, time.Unix(20, 0)); err != nil {
		t.Fatal(err)
	}
	after := checkSlotTestBytes(t, root)
	for _, retired := range outcome.Executed {
		if _, found := after[retired]; found {
			t.Fatalf("red executed check %s retained its slot", retired)
		}
	}
	if !bytes.Equal(after["package-core-guard"], before["package-core-guard"]) {
		t.Fatal("red run changed the inherited check slot")
	}
}

func TestRedRetirementSurvivesInheritedSlotDamageAfterSelection(t *testing.T) {
	t.Parallel()
	root := checkSlotTestRepo(t)
	identities := checkSlotTestIdentities()
	seed := conformanceCheckPartition{
		Tier:       registry.Dev,
		Identities: identities,
		Executed:   []string{"line-routing", "package-core-guard", "subcommand-routing"},
	}
	if err := applyConformanceCheckOutcome(root, seed, checkRunGreen, time.Unix(10, 0)); err != nil {
		t.Fatal(err)
	}
	outcome := conformanceCheckPartition{
		Tier:       registry.Dev,
		Identities: identities,
		Executed:   []string{"line-routing"},
		Inherited:  []conformanceCheckSkip{{Check: "package-core-guard", Identity: identities["package-core-guard"]}},
	}
	mutateCheckSlotTestFields(t, root, "package-core-guard", func(fields map[string]any) {
		fields["identity"] = strings.Repeat("f", 64)
	})
	if err := applyConformanceCheckOutcome(root, outcome, checkRunRed, time.Unix(20, 0)); err != nil {
		t.Fatalf("red retirement stopped behind inherited damage: %v", err)
	}
	after := checkSlotTestBytes(t, root)
	if _, found := after["line-routing"]; found {
		t.Fatal("red executed slot survived because inherited evidence changed")
	}
}

func TestHostileCheckEvidenceWidensExecution(t *testing.T) {
	t.Parallel()
	const target = "line-routing"
	for _, test := range []struct {
		name   string
		mutate func(t *testing.T, root string, identities map[string]string)
	}{
		{"missing", func(t *testing.T, root string, _ map[string]string) { deleteCheckSlotTestEntry(t, root, target) }},
		{"malformed", func(t *testing.T, root string, _ map[string]string) {
			mutateCheckSlotTestEntry(t, root, target, func(map[string]any) any { return []string{"not-a-slot"} })
		}},
		{"wrong check", func(t *testing.T, root string, _ map[string]string) {
			mutateCheckSlotTestFields(t, root, target, func(fields map[string]any) { fields["check"] = "package-core-guard" })
		}},
		{"wrong tier", func(t *testing.T, root string, _ map[string]string) {
			mutateCheckSlotTestFields(t, root, target, func(fields map[string]any) { fields["tier"] = string(registry.Ship) })
		}},
		{"stale identity", func(t *testing.T, root string, _ map[string]string) {
			mutateCheckSlotTestFields(t, root, target, func(fields map[string]any) { fields["identity"] = strings.Repeat("f", 64) })
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, identities := seededCheckSlotStore(t)
			test.mutate(t, root, identities)
			partition := partitionConformanceChecks(root, registry.Dev, identities, nil, time.Unix(20, 0))
			if !slices.Contains(partition.Executed, target) {
				t.Fatalf("executed = %v, want hostile %s slot to run", partition.Executed, target)
			}
			if slices.Contains(partition.Executed, "package-core-guard") {
				t.Fatalf("attributable %s defect widened to unrelated package-core-guard", target)
			}
		})
	}

	t.Run("unresolved check identity", func(t *testing.T) {
		root, identities := seededCheckSlotStore(t)
		delete(identities, target)
		partition := partitionConformanceChecks(root, registry.Dev, identities, nil, time.Unix(20, 0))
		if !slices.Contains(partition.Executed, target) || slices.Contains(partition.Executed, "package-core-guard") {
			t.Fatalf("executed = %v, want only unresolved %s widened", partition.Executed, target)
		}
	})

	for _, test := range []struct {
		name   string
		mutate func(t *testing.T, root string, identities map[string]string)
	}{
		{"failed identity derivation", func(_ *testing.T, _ string, _ map[string]string) {}},
		{"malformed store", func(t *testing.T, root string, _ map[string]string) {
			writeCheckSlotTestStore(t, root, []byte(`{"schema":1,"check_slots":[]}`))
		}},
		{"unknown slot", func(t *testing.T, root string, _ map[string]string) {
			slots, valid := loadConformanceCheckSlots(root)
			if !valid {
				t.Fatal("seeded store is invalid")
			}
			slots["unknown-check"] = slots[target]
			if err := replaceConformanceCheckSlotStore(root, slots); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, identities := seededCheckSlotStore(t)
			test.mutate(t, root, identities)
			var identityErr error
			if test.name == "failed identity derivation" {
				identityErr = errors.New("resolver failed")
			}
			partition := partitionConformanceChecks(root, registry.Dev, identities, identityErr, time.Unix(20, 0))
			if got, want := len(partition.Executed), len(ordinaryConformanceChecks(registry.Dev)); got != want {
				t.Fatalf("executed %d checks, want all %d after %s", got, want, test.name)
			}
		})
	}
}

func TestForgedMetaSlotCannotReceiveReusableCredit(t *testing.T) {
	t.Parallel()
	root, identities := seededCheckSlotStore(t)
	slots, valid := loadConformanceCheckSlots(root)
	if !valid {
		t.Fatal("seeded store is invalid")
	}
	slots["conformance-meta"] = slots["line-routing"]
	if err := replaceConformanceCheckSlotStore(root, slots); err != nil {
		t.Fatal(err)
	}

	partition := partitionConformanceChecks(root, registry.Dev, identities, nil, time.Unix(20, 0))
	if got, want := len(partition.Executed), len(ordinaryConformanceChecks(registry.Dev)); got != want {
		t.Fatalf("forged meta slot executed %d ordinary checks, want all %d", got, want)
	}
	if slices.Contains(partition.Executed, "conformance-meta") {
		t.Fatal("meta entered the reusable ordinary-check partition")
	}
	for _, inherited := range partition.Inherited {
		if inherited.Check == "conformance-meta" {
			t.Fatal("meta inherited forged evidence")
		}
	}
	if err := applyConformanceCheckOutcome(root, partition, checkRunGreen, time.Unix(20, 0)); err != nil {
		t.Fatal(err)
	}
	after, valid := loadConformanceCheckSlots(root)
	if !valid {
		t.Fatal("green repair left the check slot store invalid")
	}
	if _, found := after["conformance-meta"]; found {
		t.Fatal("green repair retained the forged meta slot")
	}
}

func seededCheckSlotStore(t *testing.T) (string, map[string]string) {
	t.Helper()
	root := checkSlotTestRepo(t)
	identities := map[string]string{}
	var executed []string
	for _, check := range ordinaryConformanceChecks(registry.Dev) {
		h := sha256.Sum256([]byte(check.Name))
		identities[check.Name] = hex.EncodeToString(h[:])
		executed = append(executed, check.Name)
	}
	if err := applyConformanceCheckOutcome(root, conformanceCheckPartition{Tier: registry.Dev, Identities: identities, Executed: executed}, checkRunGreen, time.Unix(10, 0)); err != nil {
		t.Fatal(err)
	}
	return root, identities
}

func deleteCheckSlotTestEntry(t *testing.T, root, name string) {
	t.Helper()
	slots, valid := loadConformanceCheckSlots(root)
	if !valid {
		t.Fatal("seeded store is invalid")
	}
	delete(slots, name)
	if err := replaceConformanceCheckSlotStore(root, slots); err != nil {
		t.Fatal(err)
	}
}

func mutateCheckSlotTestFields(t *testing.T, root, name string, mutate func(map[string]any)) {
	t.Helper()
	mutateCheckSlotTestEntry(t, root, name, func(value map[string]any) any {
		mutate(value)
		return value
	})
}

func mutateCheckSlotTestEntry(t *testing.T, root, name string, mutate func(map[string]any) any) {
	t.Helper()
	slots, valid := loadConformanceCheckSlots(root)
	if !valid {
		t.Fatal("seeded store is invalid")
	}
	var fields map[string]any
	if err := json.Unmarshal(slots[name], &fields); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(mutate(fields))
	if err != nil {
		t.Fatal(err)
	}
	slots[name] = data
	if err := replaceConformanceCheckSlotStore(root, slots); err != nil {
		t.Fatal(err)
	}
}

func writeCheckSlotTestStore(t *testing.T, root string, data []byte) {
	t.Helper()
	dir, err := componentSlotDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, conformanceCheckSlotStoreName), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func checkSlotTestRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	return root
}

func checkSlotTestIdentities() map[string]string {
	return map[string]string{
		"line-routing":       strings.Repeat("1", 64),
		"package-core-guard": strings.Repeat("2", 64),
		"subcommand-routing": strings.Repeat("3", 64),
	}
}

func checkSlotTestBytes(t *testing.T, root string) map[string][]byte {
	t.Helper()
	slots, valid := loadConformanceCheckSlots(root)
	if !valid {
		t.Fatal("check slot store is invalid")
	}
	got := make(map[string][]byte, len(slots))
	for name, raw := range slots {
		var compact bytes.Buffer
		if err := json.Compact(&compact, raw); err != nil {
			t.Fatalf("compact %s slot: %v", name, err)
		}
		got[name] = compact.Bytes()
	}
	return got
}

func checkSlotTestStoreBytes(t *testing.T, root string) map[string][]byte {
	t.Helper()
	dir, err := componentSlotDir(root)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	got := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		got[path] = data
	}
	return got
}
