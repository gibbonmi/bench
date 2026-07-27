package gate

// The probes that decide which toolchain phases a graded root materializes, and the
// edges they declare. The phases themselves are argv-only; what these tests grade is
// presence, absence, and the Needs set.

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// requireGoToolchain keeps these tests honest on a host without Go: every probe here
// requires `go` on PATH, so an absent toolchain would make the assertions grade the
// wrong branch rather than fail.
func requireGoToolchain(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("go"); err != nil {
		capability.Capability(t, capability.Tool, "go not on PATH")
	}
}

func phaseNamed(phases []Phase, name string) (Phase, bool) {
	for _, phase := range phases {
		if phase.Name == name {
			return phase, true
		}
	}
	return Phase{}, false
}

// TestPhaseTableProbedToolchainPhases grades the go.mod-probed phases: gofmt, vet, and
// test materialize for any Go root, the kit-only phases do not, and none of them takes
// an edge on the build phase — the edge set is where this split's overlap comes from.
func TestPhaseTableProbedToolchainPhases(t *testing.T) {
	requireGoToolchain(t)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "go.mod"), "module fixture\n")
	kit := "/tmp/kit"

	phases := BenchkitPhases(root, kit)
	for name, want := range map[string][]string{
		"gofmt": GateGoArgv(kit, "gofmt", root),
		"vet":   {"go", "-C", root, "vet", "./..."},
		"test":  GateGoArgv(kit, "test", root),
	} {
		phase, ok := phaseNamed(phases, name)
		if !ok {
			t.Fatalf("phase %s absent from probed table %#v", name, phaseNames(phases))
		}
		if !reflect.DeepEqual(phase.Argv, want) {
			t.Fatalf("phase %s argv = %#v, want %#v", name, phase.Argv, want)
		}
		// No edge on the build phase: none of these execs dist/bench, which is the only
		// artifact the build edge sequences writers and readers of.
		if len(phase.Needs) != 0 {
			t.Fatalf("phase %s needs = %#v, want none", name, phase.Needs)
		}
	}

	// A Go root that carries neither the race test nor internal/conformance gets neither
	// kit-only phase: probing on go.mod alone reds a linked repo on a check the kit
	// wrote for itself.
	for _, name := range []string{"race", "conformance-suite"} {
		if _, ok := phaseNamed(phases, name); ok {
			t.Fatalf("phase %s materialized for a plain go.mod root: %#v", name, phaseNames(phases))
		}
	}
}

// TestPhaseTableKitOnlyPhasesProbe grades the two phases that exist only for a root
// carrying what they grade. The race probe requires the test name rather than the
// directory, because an unrelated internal/worktree package would otherwise get a phase
// that can only red.
func TestPhaseTableKitOnlyPhasesProbe(t *testing.T) {
	requireGoToolchain(t)
	kit := "/tmp/kit"

	raceRoot := t.TempDir()
	writeFile(t, filepath.Join(raceRoot, "go.mod"), "module fixture\n")
	writeFile(t, filepath.Join(raceRoot, "internal", "worktree", "cleanup_test.go"),
		"package worktree\n\nimport \"testing\"\n\nfunc TestConcurrentCleanupRecordsOneTransaction(t *testing.T) {}\n")
	race, ok := phaseNamed(BenchkitPhases(raceRoot, kit), "race")
	if !ok {
		t.Fatalf("race phase absent for a root carrying the cleanup race test")
	}
	if want := GateGoArgv(kit, "race", raceRoot); !reflect.DeepEqual(race.Argv, want) {
		t.Fatalf("race argv = %#v, want %#v", race.Argv, want)
	}
	if len(race.Needs) != 0 {
		t.Fatalf("race needs = %#v, want none", race.Needs)
	}

	// An internal/worktree package that never declares the target test is not a race
	// root: the `-run` filter would match nothing and the step's did-it-run guard would
	// red on a repo that never asked for the check.
	strangerRoot := t.TempDir()
	writeFile(t, filepath.Join(strangerRoot, "go.mod"), "module fixture\n")
	writeFile(t, filepath.Join(strangerRoot, "internal", "worktree", "other_test.go"),
		"package worktree\n\nimport \"testing\"\n\nfunc TestSomethingElse(t *testing.T) {}\n")
	if _, ok := phaseNamed(BenchkitPhases(strangerRoot, kit), "race"); ok {
		t.Fatalf("race phase materialized for an internal/worktree without the cleanup race test")
	}

	suiteRoot := t.TempDir()
	writeFile(t, filepath.Join(suiteRoot, "go.mod"), "module fixture\n")
	writeFile(t, filepath.Join(suiteRoot, filepath.FromSlash(registry.ConformancePackage), "doc.go"), "package conformance\n")
	suite, ok := phaseNamed(BenchkitPhases(suiteRoot, kit), "conformance-suite")
	if !ok {
		t.Fatalf("conformance-suite phase absent for a root carrying %s", registry.ConformancePackage)
	}
	if want := GateGoArgv(kit, "conformance-suite", suiteRoot); !reflect.DeepEqual(suite.Argv, want) {
		t.Fatalf("conformance-suite argv = %#v, want %#v", suite.Argv, want)
	}
	// The skip pattern reaches the run from the registry through gate-go, never as a
	// literal in the phase's argv: a second copy of the non-recursion contract drifts
	// silently on the next rename.
	pattern := registry.InnerSkipPattern()
	for _, arg := range suite.Argv {
		if strings.Contains(arg, pattern) {
			t.Fatalf("conformance-suite argv carries the skip pattern literally: %#v", suite.Argv)
		}
	}
	inner := ConformanceSuiteArgv(suiteRoot)
	if !reflect.DeepEqual(inner[len(inner)-2:], []string{"-skip", pattern}) {
		t.Fatalf("gate-go conformance-suite argv = %#v, want it to end in -skip %s", inner, pattern)
	}
}

// TestPhaseTableNoToolchainNoPhases keeps both halves of the probe load-bearing. Without
// the PATH half, a host with no Go turns one attributable diagnostic into a run of
// phases that fail with exec errors; without the go.mod half, a Python repo gets Go
// phases the gate never used to touch.
func TestPhaseTableNoToolchainNoPhases(t *testing.T) {
	path := os.Getenv("PATH")
	goRoot := t.TempDir()
	writeFile(t, filepath.Join(goRoot, "go.mod"), "module fixture\n")
	t.Setenv("PATH", "")
	assertNoToolchainPhases(t, BenchkitPhases(goRoot, "/tmp/kit"))

	t.Setenv("PATH", path)
	assertNoToolchainPhases(t, BenchkitPhases(t.TempDir(), "/tmp/kit"))
}

// TestPhaseTableProbeRejectsNonRegularGoMod grades the two shapes a plain read
// misclassifies: a dangling symlink reads as an absent go.mod through os.Stat but as a
// present one through os.Lstat, and a FIFO blocks an open forever. The test's own
// deadline is the hang tripwire — a probe that opens the path stalls here rather than
// wedging the gate.
func TestPhaseTableProbeRejectsNonRegularGoMod(t *testing.T) {
	requireGoToolchain(t)
	linkRoot := t.TempDir()
	if err := os.Symlink(filepath.Join(linkRoot, "absent"), filepath.Join(linkRoot, "go.mod")); err != nil {
		t.Fatalf("symlink go.mod: %v", err)
	}
	fifoRoot := t.TempDir()
	if err := syscall.Mkfifo(filepath.Join(fifoRoot, "go.mod"), 0o644); err != nil {
		capability.Capability(t, capability.Fifo, "mkfifo unavailable: "+err.Error())
	}

	done := make(chan []Phase, 2)
	go func() {
		done <- BenchkitPhases(linkRoot, "/tmp/kit")
		done <- BenchkitPhases(fifoRoot, "/tmp/kit")
	}()
	for i := 0; i < 2; i++ {
		select {
		case phases := <-done:
			assertNoToolchainPhases(t, phases)
		case <-time.After(30 * time.Second):
			t.Fatalf("BenchkitPhases blocked on a non-regular go.mod")
		}
	}
}

func assertNoToolchainPhases(t *testing.T, phases []Phase) {
	t.Helper()
	for _, name := range []string{"gofmt", "vet", "test", "race", "conformance-suite"} {
		if _, ok := phaseNamed(phases, name); ok {
			t.Fatalf("toolchain phase %s materialized: %#v", name, phaseNames(phases))
		}
	}
}

// TestPhasesCommandVetPhaseReds drives the vet phase end to end against a tree `go vet`
// rejects. A table that omits vet, or names it with an argv anchored at the wrong tree,
// leaves the run green on that tree.
func TestPhasesCommandVetPhaseReds(t *testing.T) {
	requireGoToolchain(t)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "go.mod"), "module vetfixture\n\ngo 1.21\n")
	writeFile(t, filepath.Join(root, "main.go"),
		"package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Printf(\"%d\\n\", \"not a number\")\n}\n")

	t.Setenv("BENCH_CANARY_INNER", "1")
	t.Setenv("BENCH_CANARY_PHASE", "vet")
	var stdout, stderr bytes.Buffer
	code := PhasesCommand([]string{root}, &stdout, &stderr)
	out := stdout.String() + stderr.String()
	if code == 0 {
		t.Fatalf("PhasesCommand = 0 on a tree go vet rejects; output:\n%s", out)
	}
	if !strings.Contains(out, "Printf") {
		t.Fatalf("vet diagnostic missing from output:\n%s", out)
	}
}
