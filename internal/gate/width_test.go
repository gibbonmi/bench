package gate

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// setWidthInput sets the width input of this process, or removes it.
func setWidthInput(t *testing.T, set bool) {
	t.Helper()
	t.Setenv(widthEnv, "background")
	if !set {
		if err := os.Unsetenv(widthEnv); err != nil {
			t.Fatal(err)
		}
	}
}

// goTestArgvs answers the argv of each Go test phase of the kit's whole table and of its delta
// set, and the argv of the core test step, by label. A table with no Go test phase reds here
// rather than passes with no argv to grade.
func goTestArgvs(t *testing.T) map[string][]string {
	t.Helper()
	kit, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	argvs := map[string][]string{"core": coreTestArgv([]string{"./internal/gate"})}
	for label, phases := range map[string][]Phase{
		"whole": BenchkitPhases(kit, kit),
		"delta": deltaPhases(kit, kit, []string{"./internal/gate"}),
	} {
		selected := 0
		for _, phase := range phases {
			if goTestPhase(phase.Argv) {
				argvs[label+"/"+phase.Name] = phase.Argv
				selected++
			}
		}
		if selected == 0 {
			t.Fatalf("%s table %v has no Go test phase", label, phaseNames(phases))
		}
	}
	return argvs
}

// BG19: with the width input set, each Go test phase runs at the background width.
func TestWidthInputSetsTheTestWidths(t *testing.T) {
	setWidthInput(t, true)
	for label, argv := range goTestArgvs(t) {
		if !slices.Contains(argv, "-p=2") || !slices.Contains(argv, "-parallel=2") {
			t.Errorf("%s argv = %q, want -p=2 and -parallel=2", label, argv)
		}
	}
}

// BG20: the width input is not a grade-key member, so a background green answers a run of the
// same subject at another width.
func TestWidthInputLeavesTheGradeKey(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := outcomeFixture(t)
	keys := map[bool]string{}
	for _, set := range []bool{false, true} {
		setWidthInput(t, set)
		s, err := buildSubject(root)
		if err != nil {
			t.Fatal(err)
		}
		keys[set] = s.Identity
	}
	if keys[false] != keys[true] {
		t.Fatalf("grade key = %s without the width input and %s with it, want one key", keys[false], keys[true])
	}
}

// BG21: without the width input, no Go test phase names a width, so every foreground gate
// keeps its argv.
func TestNoWidthInputKeepsTheArgv(t *testing.T) {
	setWidthInput(t, false)
	for label, argv := range goTestArgvs(t) {
		if slices.ContainsFunc(argv, func(arg string) bool {
			return arg == "-p" || strings.HasPrefix(arg, "-p=") || strings.HasPrefix(arg, "-parallel")
		}) {
			t.Errorf("%s argv = %q, want no -p and no -parallel", label, argv)
		}
	}
}

// BG22: no phase child inherits the width input, so a nested gate in a test runs at its own
// width.
func TestPhaseChildrenDropTheWidthInput(t *testing.T) {
	setWidthInput(t, true)
	env, err := gateEnv()
	if err != nil || slices.ContainsFunc(env, func(entry string) bool { return strings.HasPrefix(entry, widthEnv+"=") }) {
		t.Fatalf("phase child environment = %q, %v; want no width input", env, err)
	}
}
