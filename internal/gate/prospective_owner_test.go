// Stable-owner tests for the prospective gate: the graded tree supplies the gate
// executable, the landing baseline supplies the phase schedule, and prospective
// evidence keys to the graded tree together with the baseline runner identity.
package gate

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
)

// stubProspectiveBuilder replaces the prospective gate's build and verification seams
// with recording stubs, and returns the source roots each authored selection was built
// from, in call order.
func stubProspectiveBuilder(t *testing.T) *[]string {
	t.Helper()
	sources := []string{}
	old := prospectiveRunBinary
	prospectiveRunBinary = runbinary.Factory{
		Build: func(_ context.Context, sourceRoot, output string) error {
			sources = append(sources, sourceRoot)
			return os.WriteFile(output, []byte("#!/bin/sh\nexit 0\n"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	}
	t.Cleanup(func() { prospectiveRunBinary = old })
	return &sources
}

// TestProspectiveGateAuthorsItsExecutableFromTheGradedTree is SOL09 and the executable
// half of SOL15. An inherited selection is present and was built from another tree, so
// honoring it would record a source digest the graded subject never produced. The gate
// authors its own executable from the graded checkout instead, and that private
// executable is removable, so no run leaves it behind.
func TestProspectiveGateAuthorsItsExecutableFromTheGradedTree(t *testing.T) {
	checkout := t.TempDir()
	kit := t.TempDir()
	inherited := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(inherited, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(runbinary.Env, inherited)
	sources := stubProspectiveBuilder(t)
	owner := prospectiveRunBinaryOwner(checkout)

	graded, err := owner(t.Context(), checkout)
	if err != nil {
		t.Fatal(err)
	}
	if len(*sources) != 1 || (*sources)[0] != checkout {
		t.Fatalf("graded-tree sources = %#v, want one selection built from %q", *sources, checkout)
	}
	if graded.Path == inherited {
		t.Fatalf("graded subject ran under the inherited selection %q", inherited)
	}
	if err := graded.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Dir(graded.Path)); !os.IsNotExist(err) {
		t.Fatalf("private gate executable survived its selection: %v", err)
	}

	// A source that is not the graded tree is the baseline's own kit, whose inherited
	// selection is already the baseline runner.
	baseline, err := owner(t.Context(), kit)
	if err != nil {
		t.Fatal(err)
	}
	if len(*sources) != 1 || baseline.Path != inherited {
		t.Fatalf("baseline selection = %q with sources %#v, want the inherited %q", baseline.Path, *sources, inherited)
	}
}

// TestProspectiveGateKeepsTheBaselinePhaseScheduleAfterACandidateOmission is SOL10. The
// candidate tree declares a phase manifest that omits a phase the baseline declares.
// The baseline's schedule is the one the run gets.
func TestProspectiveGateKeepsTheBaselinePhaseScheduleAfterACandidateOmission(t *testing.T) {
	baseline := t.TempDir()
	candidate := t.TempDir()
	writeGateFixtureFile(t, baseline, ".bench/phases.json", `{"phases":[{"name":"planted","argv":["true"]},{"name":"kept","argv":["true"]}]}`, 0o644)
	writeGateFixtureFile(t, candidate, ".bench/phases.json", `{"phases":[{"name":"kept","argv":["true"]}]}`, 0o644)
	t.Setenv(baselinePolicyEnv, baseline)

	phases, err := phaseTable(candidate, candidate)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"planted", "kept"}; strings.Join(phaseNames(phases), ",") != strings.Join(want, ",") {
		t.Fatalf("phase names = %v, want the baseline's %v", phaseNames(phases), want)
	}
}

// TestProspectivePhaseScheduleRefusesAnUnusableBaseline keeps the baseline selection
// fail-closed. A named baseline the run cannot read refuses; a fall back to the
// candidate tree would reinstate the omission SOL10 excludes.
func TestProspectivePhaseScheduleRefusesAnUnusableBaseline(t *testing.T) {
	candidate := t.TempDir()
	writeGateFixtureFile(t, candidate, ".bench/phases.json", `{"phases":[{"name":"kept","argv":["true"]}]}`, 0o644)
	t.Setenv(baselinePolicyEnv, filepath.Join(t.TempDir(), "absent"))

	if _, err := phaseTable(candidate, candidate); err == nil {
		t.Fatal("unusable baseline fell back to the candidate schedule")
	}
}

// TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity is SOL13. Retained
// prospective evidence is addressed by the graded tree together with the baseline
// runner identity. A key that binds only one of them reuses evidence after the other
// half of the subject changes.
func TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity(t *testing.T) {
	root := t.TempDir()
	writeGateFixtureFile(t, root, ".bench/gate.sh", "#!/bin/sh\nexit 0\n", 0o755)
	writeGateFixtureFile(t, root, ".bench/gate-prospective.sh", "#!/bin/sh\nexit 0\n", 0o755)
	writeGateFixtureFile(t, root, ".bench/gate-inputs.json", "{\"schema\":1,\"closure\":\"local\",\"environment\":[],\"paths\":[],\"tools\":[]}\n", 0o644)
	// One identity root answers for both runners in turn, so only the runner identity
	// differs between the two keys. Two directories would differ by their own paths,
	// which the subject already binds, and the assertion would pass without the runner.
	identityRoot := t.TempDir()
	tree := strings.Repeat("a", 40)
	other := strings.Repeat("b", 40)

	key := func(tree string) string {
		t.Helper()
		plan, err := buildSubjectForTree(root, identityRoot, policyVersion, tree)
		if err != nil {
			t.Fatal(err)
		}
		return evidenceName(plan)
	}
	unbound := key(tree)
	if unbound != key(tree) {
		t.Fatal("an unchanged subject did not address the same evidence")
	}
	if unbound == key(other) {
		t.Fatal("evidence key ignored the graded tree")
	}
	if got := baselineRunnerIdentity(identityRoot); got != unboundBaselineRunner {
		t.Fatalf("undeclared baseline runner identity = %q, want %q", got, unboundBaselineRunner)
	}
	writeGateFixtureFile(t, identityRoot, "scripts/go-build.inputs", "build_script=scripts/go-build.sh\n", 0o644)
	if got := baselineRunnerIdentity(identityRoot); got == unboundBaselineRunner {
		t.Fatal("a declared build recipe still answered the unbound runner identity")
	}
	if unbound == key(tree) {
		t.Fatal("evidence key ignored the baseline runner identity")
	}
}

func writeGateFixtureFile(t *testing.T, root, rel, body string, mode os.FileMode) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
}

// TestProspectiveGateReportsAFailedBuildWithNoResidue is the fail-closed half of the
// graded-tree selection. Authoring the private gate executable fails, so the owner
// returns the build cause instead of a selection, and the private directory it opened
// for that executable does not survive the failure.
func TestProspectiveGateReportsAFailedBuildWithNoResidue(t *testing.T) {
	checkout := t.TempDir()
	tempRoot := t.TempDir()
	old := prospectiveRunBinary
	prospectiveRunBinary = runbinary.Factory{
		TempRoot: tempRoot,
		Build:    func(context.Context, string, string) error { return errors.New("build refused") },
		Verify:   func(string, string) error { return nil },
	}
	t.Cleanup(func() { prospectiveRunBinary = old })

	selection, err := prospectiveRunBinaryOwner(checkout)(t.Context(), checkout)
	if err == nil || selection != nil {
		t.Fatalf("failed prospective build = (%#v, %v), want no selection and an error", selection, err)
	}
	if !strings.Contains(err.Error(), "build refused") {
		t.Fatalf("prospective build error = %v, want the build cause", err)
	}
	entries, readErr := os.ReadDir(tempRoot)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("failed prospective build left %d entries under its temporary root", len(entries))
	}
}
