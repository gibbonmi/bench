package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"

	kitpayload "github.com/gibbonmi/bench"
	"github.com/gibbonmi/bench/internal/gate"
)

// The component-input registry records, per entry, the derivation its declared paths came
// from. A table check can prove the profile and that record agree; it cannot prove the
// record is honest, because a hand-copied path list satisfies any table and drifts from
// what the component grades the moment the tree moves. This check closes that gap by
// resolving the registry against a tree it controls and then moving that tree: a set
// computed from a derivation follows, and a restated list does not.
//
// What it establishes: every derivable entry's declaration is a function of the tree it
// was resolved against — its paths name files that tree holds, and they change when the
// tree gains files that derivation covers. What it does not establish: that a derivation
// is complete or names the right files (a derivation that under-reports still tracks the
// tree), nor that an entry calls the exact function its source label names (any honest
// re-derivation of the same shape would pass). It grades provenance, not correctness.

// derivationFixtureSeed is the tree the registry is resolved against: the smallest module
// every derivation can resolve — a build closure rooted at ./cmd/bench, a module-wide test
// closure with testdata, the shell files the shellcheck argv enumerates, and a published
// seal for the component that reads the binary's source digest.
var derivationFixtureSeed = map[string]string{
	"go.mod":                            "module benchderivationfixture\n\ngo 1.25\n",
	"cmd/bench/main.go":                 "package main\n\nfunc main() {}\n",
	"internal/graded/graded.go":         "package graded\n\n// Name is the fixture package's only surface.\nfunc Name() string { return \"graded\" }\n",
	"internal/graded/graded_test.go":    "package graded\n\nimport \"testing\"\n\nfunc TestName(t *testing.T) {\n\tif Name() == \"\" {\n\t\tt.Fatal(\"empty\")\n\t}\n}\n",
	"internal/graded/testdata/seed.txt": "seed\n",
	"scripts/go-build.sh":               "#!/usr/bin/env bash\nexit 0\n",
	"scripts/go-build.inputs":           "build_script=scripts/go-build.sh\n",
	"bin/bench.sh":                      "#!/usr/bin/env bash\nexit 0\n",
	".agents/commands/seed.md":          "# Seed\n",
	".bench/lib/seeded.sh":              "#!/usr/bin/env bash\nexit 0\n",
	"dist/bench.seal":                   `{"schema":1,"sources":"` + strings.Repeat("a", 64) + `","executable":"` + strings.Repeat("b", 64) + `"}`,
}

// derivationFixturePerturbation is the tree change every derivable entry must follow. It
// adds one file per derivation family — a source inside the binary's closure, a shell file
// the argv enumeration picks up, and a new seal source digest — so the check needs no
// per-component table of which file belongs to which derivation: a set that ignores all
// three ignored the tree.
var derivationFixturePerturbation = map[string]string{
	"cmd/bench/added.go":                 "package main\n\nfunc added() {}\n",
	".bench/lib/added.sh":                "#!/usr/bin/env bash\nexit 0\n",
	"internal/graded/testdata/added.txt": "added\n",
	".agents/commands/added.md":          "# Added\n",
	"dist/bench.seal":                    `{"schema":1,"sources":"` + strings.Repeat("c", 64) + `","executable":"` + strings.Repeat("d", 64) + `"}`,
}

// perturbationSummary names the change in a diagnostic, so a red says what the declaration
// failed to follow rather than only that it failed.
const perturbationSummary = "a source in the binary closure, a shell file, a testdata file, agent Markdown, and a new seal source digest"

// componentDeclaration is one registry entry read through the gate's accessors. The check
// grades this shape rather than gate.ComponentInputs so a bite proof can mutate a
// declaration without editing the real registry.
type componentDeclaration struct {
	component string
	source    gate.Source
	paths     []string
	digests   []string
}

// declarationResolver resolves the whole registry against a tree.
type declarationResolver func(root string) (map[string]componentDeclaration, error)

// resolveDeclarations is the real resolver: the registry itself, read through its accessors.
func resolveDeclarations(root string) (map[string]componentDeclaration, error) {
	sets, err := gate.ResolveComponentInputs(root)
	if err != nil {
		return nil, err
	}
	declarations := make(map[string]componentDeclaration, len(sets))
	for component, inputs := range sets {
		declarations[component] = componentDeclaration{
			component: component,
			source:    inputs.Source(),
			paths:     inputs.Paths(),
			digests:   inputs.Digests(),
		}
	}
	return declarations, nil
}

// resolveAcrossPerturbation resolves the registry against root, moves the tree, and
// resolves again. The two resolutions are the evidence the grading works from, so they are
// paid once and the mutations a bite proof drives operate on the result rather than
// re-resolving. Seeding root is the caller's job, so the check can also be pointed at a
// tree no derivation can resolve.
func resolveAcrossPerturbation(root string, resolve declarationResolver) (before, after map[string]componentDeclaration, err error) {
	if before, err = resolve(root); err != nil {
		return nil, nil, err
	}
	if err := writeDerivationFiles(root, derivationFixturePerturbation); err != nil {
		return nil, nil, err
	}
	if after, err = resolve(root); err != nil {
		return nil, nil, err
	}
	return before, after, nil
}

func writeDerivationFiles(root string, files map[string]string) error {
	for _, rel := range sortedKeys(files) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(files[rel]), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// writeDerivationFixture materializes both the derivation-specific seed and the
// consumer inventory every contract-input resolution requires. The inventory and its
// assets come from the kit's authoritative payload rather than a copied fixture list.
func writeDerivationFixture(root string) error {
	if err := writeDerivationFiles(root, derivationFixtureSeed); err != nil {
		return err
	}
	rows, err := kitpayload.PayloadRows()
	if err != nil {
		return err
	}
	payload, err := json.Marshal(rows)
	if err != nil {
		return err
	}
	if err := writeDerivationFiles(root, map[string]string{
		".bench/consumer-payload.json": string(payload),
	}); err != nil {
		return err
	}
	for _, row := range kitpayload.PayloadConsumerRows(rows) {
		path := filepath.Join(root, filepath.FromSlash(row.Source))
		if row.Tree {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return err
			}
			continue
		}
		info, err := os.Stat(path)
		if err == nil {
			if !info.Mode().IsRegular() {
				return fmt.Errorf("consumer fixture asset %q is not a regular file", row.Source)
			}
			continue
		}
		if !os.IsNotExist(err) {
			return err
		}
		mode, err := strconv.ParseUint(row.Mode, 8, 32)
		if err != nil {
			return fmt.Errorf("parse consumer fixture mode %q for %q: %w", row.Mode, row.Source, err)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte("consumer fixture asset\n"), os.FileMode(mode)); err != nil {
			return err
		}
	}
	return nil
}

// seedDerivationFixture writes the tree the registry is resolved against and returns its
// root.
func seedDerivationFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := writeDerivationFixture(root); err != nil {
		t.Fatalf("seed the derivation fixture: %v", err)
	}
	return root
}

// checkDerivationSource grades the registry against root: it resolves, moves the tree, and
// resolves again, then reports every entry whose declaration did not follow. A resolution
// failure is its own diagnostic — read as an empty set it would satisfy every check below,
// which is the shape that turns a broken derivation into silent green.
func checkDerivationSource(root string, resolve declarationResolver) []string {
	before, after, err := resolveAcrossPerturbation(root, resolve)
	if err != nil {
		return []string{fmt.Sprintf("component input derivation unresolvable: resolving the registry against %s failed: %v — a derivation that cannot resolve declares no inputs, which grades as an empty set nothing can violate", root, err)}
	}
	return derivationSourceDiags(root, before, after)
}

// derivationSourceDiags grades one resolved pair. It reports at most one diagnostic per
// component — the first property that failed is the whole finding for that entry — and
// carries on through the rest, so one run names the whole registry's failures.
func derivationSourceDiags(root string, before, after map[string]componentDeclaration) []string {
	var diags []string
	for _, component := range sortedKeys(before) {
		if diag := componentDerivationDiag(root, before[component], after); diag != "" {
			diags = append(diags, diag)
		}
	}
	return diags
}

func componentDerivationDiag(root string, seeded componentDeclaration, after map[string]componentDeclaration) string {
	if seeded.source == gate.SourceHandDeclared {
		// Named rather than exempted silently: canary's inputs have no derivable source, and
		// a second entry claiming the same needs a reviewer's decision, not a default.
		if seeded.component != "canary" {
			return fmt.Sprintf("component input hand-declared beyond canary: %s declares source %q, and canary is the only entry whose inputs have no derivable source — a second hand-declared entry is a reviewer decision", seeded.component, seeded.source)
		}
		return ""
	}
	moved, present := after[seeded.component]
	if !present {
		return fmt.Sprintf("component input entry vanished: %s (source %s) resolved against the seeded tree but not after %s was added", seeded.component, seeded.source, perturbationSummary)
	}
	for _, path := range moved.paths {
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(path))); err != nil {
			return fmt.Sprintf("component input path names no file in the resolved tree: %s (source %s) declares %q, which the tree it was resolved against does not hold — %s reports what the tree holds, a copied list reports what some other tree held", seeded.component, seeded.source, path, seeded.source)
		}
	}
	if !gainedPath(seeded.paths, moved.paths) {
		return fmt.Sprintf("component input set did not follow the tree: %s (source %s) declared the same %d paths before and after %s — a declaration that ignores the tree is a restated list, not a resolution of %s", seeded.component, seeded.source, len(moved.paths), perturbationSummary, seeded.source)
	}
	if len(seeded.digests) != 0 && slices.Equal(seeded.digests, moved.digests) {
		return fmt.Sprintf("component input digests did not follow the tree: %s (source %s) declared digests [%s] both before and after %s — a digest names content, so a declaration that reports the same one for changed content is pasted, not derived", seeded.component, seeded.source, strings.Join(moved.digests, ", "), perturbationSummary)
	}
	return ""
}

// gainedPath reports whether moved holds a path seeded does not. Growth rather than mere
// inequality: the perturbation only adds files, so a derivation that read the tree gained
// one, and a comparison satisfied by a dropped path would accept a set that shrank for
// reasons of its own.
func gainedPath(seeded, moved []string) bool {
	held := make(map[string]struct{}, len(seeded))
	for _, path := range seeded {
		held[path] = struct{}{}
	}
	for _, path := range moved {
		if _, ok := held[path]; !ok {
			return true
		}
	}
	return false
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// checkRegisteredDerivationSource grades the real component-input registry against a
// disposable tree, keeping the mutation out of both the subject and the kit checkout.
func checkRegisteredDerivationSource(_ string) []string {
	root, err := os.MkdirTemp("", "bench-component-derivation-")
	if err != nil {
		return []string{"component input derivation fixture cannot be created: " + err.Error()}
	}
	defer os.RemoveAll(root)
	if err := writeDerivationFixture(root); err != nil {
		return []string{"component input derivation fixture cannot be seeded: " + err.Error()}
	}
	return checkDerivationSource(root, resolveDeclarations)
}

// derivableComponents names the registry entries this check holds to their derivation,
// read from a resolution rather than restated: an entry added later is graded by the bite
// proof automatically, and one that stops being derivable drops out of it.
func derivableComponents(t *testing.T, declarations map[string]componentDeclaration) []string {
	t.Helper()
	var components []string
	for _, component := range sortedKeys(declarations) {
		if declarations[component].source != gate.SourceHandDeclared {
			components = append(components, component)
		}
	}
	if len(components) == 0 {
		t.Fatal("the registry declares no derivable component, so this check grades nothing")
	}
	return components
}

// digestDeclaringComponents names the derivable entries whose declaration carries a
// content digest, read from the resolution so an entry that gains or loses one moves this
// set with it.
func digestDeclaringComponents(declarations map[string]componentDeclaration) []string {
	var components []string
	for _, component := range sortedKeys(declarations) {
		declaration := declarations[component]
		if declaration.source != gate.SourceHandDeclared && len(declaration.digests) != 0 {
			components = append(components, component)
		}
	}
	return components
}

// resolvedFixturePair seeds a fixture and returns the two real resolutions, so each bite
// case mutates a copy of honest evidence instead of paying for a fresh pair of listings.
func resolvedFixturePair(t *testing.T) (root string, before, after map[string]componentDeclaration) {
	t.Helper()
	root = seedDerivationFixture(t)
	before, after, err := resolveAcrossPerturbation(root, resolveDeclarations)
	if err != nil {
		t.Fatalf("resolve the registry against the fixture: %v", err)
	}
	return root, before, after
}

func cloneDeclarations(declarations map[string]componentDeclaration) map[string]componentDeclaration {
	out := make(map[string]componentDeclaration, len(declarations))
	for component, declaration := range declarations {
		declaration.paths = slices.Clone(declaration.paths)
		declaration.digests = slices.Clone(declaration.digests)
		out[component] = declaration
	}
	return out
}

// TestDerivationSourceCheckBites is the recorded bite proof. The fixed
// side is the real registry resolved against a real tree; each case substitutes one
// entry's declaration with the shape a hand-copied list produces and requires exactly one
// diagnostic naming that component and the derivation it was supposed to resolve through.
func TestDerivationSourceCheckBites(t *testing.T) {
	root, before, after := resolvedFixturePair(t)

	if diags := derivationSourceDiags(root, before, after); len(diags) != 0 {
		t.Fatalf("the real registry resolved against a moving tree got diagnostics:\n%s", strings.Join(diags, "\n"))
	}

	for _, component := range derivableComponents(t, before) {
		source := before[component].source
		t.Run(component+" restated as the list it resolved to once", func(t *testing.T) {
			// The copied-list defect exactly: a set captured from one tree and reported
			// unchanged for the next, every path of it still real.
			mutated := cloneDeclarations(after)
			frozen := mutated[component]
			frozen.paths = slices.Clone(before[component].paths)
			frozen.digests = slices.Clone(before[component].digests)
			mutated[component] = frozen
			diags := derivationSourceDiags(root, before, mutated)
			if len(diags) != 1 || !strings.Contains(diags[0], "component input set did not follow the tree") ||
				!strings.Contains(diags[0], component) || !strings.Contains(diags[0], string(source)) {
				t.Fatalf("freezing %s's set: want one diagnostic naming the component and %s, got:\n%s", component, source, strings.Join(diags, "\n"))
			}
		})
		t.Run(component+" copied from another tree", func(t *testing.T) {
			mutated := cloneDeclarations(after)
			copied := mutated[component]
			copied.paths = append(slices.Clone(copied.paths), "internal/gate/component_inputs.go")
			mutated[component] = copied
			diags := derivationSourceDiags(root, before, mutated)
			if len(diags) != 1 || !strings.Contains(diags[0], "component input path names no file in the resolved tree") ||
				!strings.Contains(diags[0], component) || !strings.Contains(diags[0], string(source)) {
				t.Fatalf("pasting a foreign path into %s's set: want one diagnostic naming the component and %s, got:\n%s", component, source, strings.Join(diags, "\n"))
			}
		})
	}

	// A digest-declaring entry is graded when the registry has one; the case is registered
	// from the resolution rather than skipped inside, so an empty run is a registry fact
	// rather than a silent skip.
	for _, digested := range digestDeclaringComponents(before) {
		t.Run(digested+" reporting a pasted digest while its paths track", func(t *testing.T) {
			mutated := cloneDeclarations(after)
			pasted := mutated[digested]
			pasted.digests = slices.Clone(before[digested].digests)
			mutated[digested] = pasted
			diags := derivationSourceDiags(root, before, mutated)
			if len(diags) != 1 || !strings.Contains(diags[0], "component input digests did not follow the tree") || !strings.Contains(diags[0], digested) {
				t.Fatalf("freezing %s's digest: want one diagnostic naming the component, got:\n%s", digested, strings.Join(diags, "\n"))
			}
		})
	}

	t.Run("an entry that stops resolving is named, not read as satisfied", func(t *testing.T) {
		mutated := cloneDeclarations(after)
		vanished := derivableComponents(t, before)[0]
		delete(mutated, vanished)
		diags := derivationSourceDiags(root, before, mutated)
		if len(diags) != 1 || !strings.Contains(diags[0], "component input entry vanished") || !strings.Contains(diags[0], vanished) {
			t.Fatalf("dropping %s from the moved resolution: want one diagnostic naming it, got:\n%s", vanished, strings.Join(diags, "\n"))
		}
	})
}

// TestOnlyCanaryMayBeHandDeclared drives the second hand-declared entry the registry does
// not have yet. canary's own entry is the fixed side: it must stay silent while any other
// entry claiming the same source reds.
func TestOnlyCanaryMayBeHandDeclared(t *testing.T) {
	root, before, after := resolvedFixturePair(t)

	if before["canary"].source != gate.SourceHandDeclared {
		t.Fatalf("canary declares source %q, want the hand-declared source this check exempts", before["canary"].source)
	}

	for _, component := range derivableComponents(t, before) {
		t.Run(component+" declaring itself hand-written", func(t *testing.T) {
			mutated := cloneDeclarations(before)
			handed := mutated[component]
			handed.source = gate.SourceHandDeclared
			mutated[component] = handed
			diags := derivationSourceDiags(root, mutated, after)
			if len(diags) != 1 || !strings.Contains(diags[0], "component input hand-declared beyond canary") || !strings.Contains(diags[0], component) {
				t.Fatalf("marking %s hand-declared: want one diagnostic naming it, got:\n%s", component, strings.Join(diags, "\n"))
			}
		})
	}
}

// TestUnresolvableDerivationIsNamed points the check at an empty tree — no go.mod, so no
// derivation can resolve — through the real resolver. The empty set a failed resolution
// would otherwise present satisfies every property above, so the failure has to arrive as
// its own diagnostic naming the derivation that could not resolve.
func TestUnresolvableDerivationIsNamed(t *testing.T) {
	root := t.TempDir()
	diags := checkDerivationSource(root, resolveDeclarations)
	if len(diags) != 1 || !strings.Contains(diags[0], "component input derivation unresolvable") ||
		!strings.Contains(diags[0], "derive ") || !strings.Contains(diags[0], root) {
		t.Fatalf("a registry resolved against a tree with no go.mod: want one diagnostic naming the failure, got:\n%s", strings.Join(diags, "\n"))
	}
}
