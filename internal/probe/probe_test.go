package probe

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
)

// fixture is one probefixture repository the verb runs over: a Go module with the Clamp
// function under test, a cmd/bench package that gives the run-binary digest a stable
// closure, and the Bench home the preserved copy lands in. The module's own sources are
// outside that closure, so a probe's mutation never invalidates the selected executable.
type fixture struct {
	root    string
	home    string
	subject string
	marker  string
}

const clampSource = `package probefixture

// Clamp keeps n at or above zero.
func Clamp(n int) int {
	if n < 0 {
		return 0
	}
	return n
}
`

const clampTestSource = `package probefixture

import "testing"

func TestClampNegative(t *testing.T) {
	if got := Clamp(-1); got != 0 {
		t.Fatalf("Clamp(-1) = %d, want 0", got)
	}
}

func TestClampPositive(t *testing.T) {
	if got := Clamp(3); got != 3 {
		t.Fatalf("Clamp(3) = %d, want 3", got)
	}
}
`

// newFixture builds the repository, publishes a run binary bound to it, and makes it the
// working directory. The environment is bound here rather than per test, so every row
// observes one home, one selected executable, and one repository.
func newFixture(t *testing.T) *fixture {
	t.Helper()
	root := scratchDir(t)
	writeFixtureFile(t, filepath.Join(root, "go.mod"), "module probefixture\n\ngo 1.25\n", 0o644)
	writeFixtureFile(t, filepath.Join(root, "clamp.go"), clampSource, 0o644)
	writeFixtureFile(t, filepath.Join(root, "clamp_test.go"), clampTestSource, 0o644)
	writeFixtureFile(t, filepath.Join(root, "cmd", "bench", "main.go"), "package main\n\nfunc main() {}\n", 0o644)
	// The build-input manifest and the version file complete the digest closure the
	// run-binary seal is graded against. Both sit outside the module's own sources, so a
	// probe's mutation never invalidates the selected executable mid-run.
	writeFixtureFile(t, filepath.Join(root, "scripts", "go-build.inputs"), "package_version=package.json\n", 0o644)
	writeFixtureFile(t, filepath.Join(root, "package.json"), "{\"version\":\"0.0.0\"}\n", 0o644)
	gitInit(t, root)

	home := scratchDir(t)
	t.Setenv("BENCH_HOME", home)
	t.Setenv("BENCH_KIT", root)
	publishRunBinary(t, root)
	t.Chdir(root)
	return &fixture{
		root:    root,
		home:    home,
		subject: filepath.Join(root, "clamp.go"),
		marker:  filepath.Join(scratchDir(t), "go-started"),
	}
}

// scratchDir answers a temporary directory with every symbolic link resolved. The
// run-binary selection refuses a path that traverses a link, and the repository root
// comes back from Git already resolved, so both halves must agree.
func scratchDir(t *testing.T) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func writeFixtureFile(t *testing.T, path, body string, mode os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), mode); err != nil {
		t.Fatal(err)
	}
}

func gitInit(t *testing.T, root string) {
	t.Helper()
	for _, args := range [][]string{{"init"}, {"config", "user.email", "probe@example.com"}, {"config", "user.name", "probe"}} {
		command := exec.Command("git", append([]string{"-C", root}, args...)...)
		if out, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
}

// publishRunBinary binds a sealed executable to the fixture root, so the focused run
// inherits a selection instead of building one. Inheriting keeps the fixture free of the
// kit's build script while leaving production's verification path in force.
func publishRunBinary(t *testing.T, root string) {
	t.Helper()
	dir := scratchDir(t)
	staged := filepath.Join(dir, "staged")
	writeFixtureFile(t, staged, "#!/bin/sh\nexit 0\n", 0o755)
	executable := filepath.Join(dir, "bench")
	if err := freshness.Publish(root, staged, executable, dir, "probe-test"); err != nil {
		t.Fatal(err)
	}
	t.Setenv(runbinary.Env, executable)
}

// installStubGo puts a canned `go` ahead of the real one and records every start in the
// fixture's marker. A `list` call still reaches the real toolchain, because the run
// binary's freshness verification resolves the module closure through it. hook is extra
// shell the stub runs while the focused run is in flight.
func installStubGo(t *testing.T, f *fixture, events string, exit int, hook string) {
	t.Helper()
	real, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	dir := scratchDir(t)
	script := "#!/usr/bin/env bash\n" +
		"printf '%s\\n' \"$1\" >> " + sanitize.ShellQuote(f.marker) + "\n" +
		"if [ \"$1\" = list ]; then exec " + sanitize.ShellQuote(real) + " \"$@\"; fi\n" +
		hook +
		"cat <<'PROBEEOF'\n" + events + "\nPROBEEOF\n" +
		"exit " + strconv.Itoa(exit) + "\n"
	writeFixtureFile(t, filepath.Join(dir, "go"), script, 0o755)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// cannedFailure is one failing test in the fixture package. diagnostic is the first
// output line, which is the cell the failures table carries.
func cannedFailure(diagnostic string) string {
	return `{"Action":"run","Package":"probefixture","Test":"TestClampNegative"}
{"Action":"output","Package":"probefixture","Test":"TestClampNegative","Output":"    ` + diagnostic + `\n"}
{"Action":"fail","Package":"probefixture","Test":"TestClampNegative","Elapsed":0}
{"Action":"fail","Package":"probefixture","Elapsed":0.01}`
}

const cannedPass = `{"Action":"run","Package":"probefixture","Test":"TestClampPositive"}
{"Action":"pass","Package":"probefixture","Test":"TestClampPositive","Elapsed":0}
{"Action":"pass","Package":"probefixture","Elapsed":0.01}`

// probeRow derives one expected verdict block through the encoder the verb renders with.
// A hand-joined row would disagree with the encoder's own quoting the moment a cell
// carries a space or looks numeric.
func probeRow(t *testing.T, verdict, subject, mutation, cause string, failed int, restored string) string {
	t.Helper()
	row := []any{verdict, subject, mutation, cause, failed, restored}
	out, err := toon.TableTyped("probe", probeFields, [][]any{row})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// selectionRow derives one expected selection block through the encoder the verb renders
// with, so the expectation and the row share one quoting rule.
func selectionRow(t *testing.T, form, target, run string, ran int) string {
	t.Helper()
	out, err := toon.TableTyped("selection", selectionFields, [][]any{{form, target, run, ran}})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// runProbe drives the verb and answers its stdout and exit.
func runProbe(t *testing.T, args ...string) (string, int) {
	t.Helper()
	return Command(args)
}

func requireSubjectBytes(t *testing.T, f *fixture, want string) {
	t.Helper()
	got, err := os.ReadFile(f.subject)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("subject bytes = %q, want %q", got, want)
	}
}

// requireHomeEmpty proves a completed or refused probe left the Bench home as it was.
func requireHomeEmpty(t *testing.T, f *fixture) {
	t.Helper()
	entries, err := os.ReadDir(f.home)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("Bench home holds %d entries, want none", len(entries))
	}
}

func requireNoRunChild(t *testing.T, f *fixture) {
	t.Helper()
	if _, err := os.Stat(f.marker); !os.IsNotExist(err) {
		body, _ := os.ReadFile(f.marker)
		t.Fatalf("a run child started before the refusal: %q", body)
	}
}

// preservedCopy answers the one file below the home's probe directory.
func preservedCopy(t *testing.T, f *fixture) string {
	t.Helper()
	var found []string
	root := filepath.Join(f.home, "probe")
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			found = append(found, path)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 {
		t.Fatalf("preserved files = %v, want exactly one", found)
	}
	return found[0]
}

// PB1, PB13, PB15: a mutation a focused test catches is the success case, the subject
// comes back byte-exact, and the home keeps no copy.
func TestProbeBitesWhenTheFocusedTestFails(t *testing.T) {
	f := newFixture(t)
	out, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$")
	if code != 0 {
		t.Fatalf("exit = %d, want 0\n%s", code, out)
	}
	want := probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "yes")
	if !strings.HasPrefix(out, want) {
		t.Fatalf("stdout = %q, want the verdict row %q first", out, want)
	}
	if !strings.Contains(out, "packages[") {
		t.Fatalf("stdout = %q, want the focused run's tables after the row", out)
	}
	requireSubjectBytes(t, f, clampSource)
	requireHomeEmpty(t, f)
}

// PB2: --omit deletes the one match, and the shortened guard is what the test catches.
func TestProbeOmitsTheMatchOnce(t *testing.T) {
	f := newFixture(t)
	out, code := runProbe(t, "clamp.go", "--omit", "return 0", "--package", "./", "--run", "^TestClampNegative$")
	if code != 0 {
		t.Fatalf("exit = %d, want 0\n%s", code, out)
	}
	if want := probeRow(t, "bit", "clamp.go", "omit", "failed", 1, "yes"); !strings.HasPrefix(out, want) {
		t.Fatalf("stdout = %q, want %q first", out, want)
	}
	requireSubjectBytes(t, f, clampSource)
}

// PB3: a mutation no test observes is silent and red, because it proves nothing.
func TestProbeIsSilentWhenTheFocusedTestPasses(t *testing.T) {
	f := newFixture(t)
	out, code := runProbe(t, "clamp.go", "--swap", "keeps n at", "--with", "holds n at", "--package", "./", "--run", "^TestClampNegative$")
	if code != 1 {
		t.Fatalf("exit = %d, want 1\n%s", code, out)
	}
	if want := probeRow(t, "silent", "clamp.go", "swap", "passed", 0, "yes"); !strings.HasPrefix(out, want) {
		t.Fatalf("stdout = %q, want %q first", out, want)
	}
	requireSubjectBytes(t, f, clampSource)
}

// PB4: a mutated package that does not compile never bit, so it is invalid rather than a
// pass or a bite. The selection names no run pattern, because `bench test` attributes a
// build failure under a run pattern to the pattern.
func TestProbeIsInvalidWhenTheMutationDoesNotCompile(t *testing.T) {
	f := newFixture(t)
	out, code := runProbe(t, "clamp.go", "--swap", "return n", "--with", "return", "--package", "./")
	if code != 1 {
		t.Fatalf("exit = %d, want 1\n%s", code, out)
	}
	if want := probeRow(t, "invalid", "clamp.go", "swap", "build-failed", 0, "yes"); !strings.HasPrefix(out, want) {
		t.Fatalf("stdout = %q, want %q first", out, want)
	}
	requireSubjectBytes(t, f, clampSource)
}

// PB5: a run pattern that matches nothing ran no test, so a mistyped pattern cannot bite.
func TestProbeIsInvalidWhenNoTestRuns(t *testing.T) {
	f := newFixture(t)
	out, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestNoSuch$")
	if code != 1 {
		t.Fatalf("exit = %d, want 1\n%s", code, out)
	}
	if want := probeRow(t, "invalid", "clamp.go", "swap", "no-test-run", 0, "yes"); !strings.HasPrefix(out, want) {
		t.Fatalf("stdout = %q, want %q first", out, want)
	}
	requireSubjectBytes(t, f, clampSource)
}

// PB6: the evidence after the row is what `bench test` prints for the same selection.
func TestProbeCarriesTheFocusedRunTables(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("clamp_test.go:6: Clamp(-1) = -1, want 0"), 1, "")
	out, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$")
	if code != 0 {
		t.Fatalf("exit = %d, want 0\n%s", code, out)
	}
	row := probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "yes")
	report, reportCode := testreport.Command(f.root, []string{"--package", "./", "--run", "^TestClampNegative$"})
	if reportCode != 1 {
		t.Fatalf("focused run exit = %d, want 1\n%s", reportCode, report)
	}
	if want := row + selectionRow(t, "package", "./", "^TestClampNegative$", 1) + report; out != want {
		t.Fatalf("stdout = %q, want %q", out, want)
	}
}

// twoRuns is one stream with a run event for each of two tests, so the count the row
// carries is observable apart from the number of failures.
const twoRuns = `{"Action":"run","Package":"probefixture","Test":"TestClampNegative"}
{"Action":"output","Package":"probefixture","Test":"TestClampNegative","Output":"    caught\n"}
{"Action":"fail","Package":"probefixture","Test":"TestClampNegative","Elapsed":0}
{"Action":"run","Package":"probefixture","Test":"TestClampPositive"}
{"Action":"pass","Package":"probefixture","Test":"TestClampPositive","Elapsed":0}
{"Action":"fail","Package":"probefixture","Elapsed":0.01}`

// DG3: the ran cell counts the distinct tests that emitted a run event, so a run that
// started no test cannot read as evidence beside its invalid verdict.
func TestProbeCountsTheTestsThatRan(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, twoRuns, 1, "")
	out, code := runProbe(t, probeArgs()...)
	if want := selectionRow(t, "package", "./", "^TestClampNegative$", 2); code != 0 || !strings.Contains(out, want) {
		t.Fatalf("stdout = (%q, %d), want %q and 0", out, code, want)
	}
	installStubGo(t, f, `{"Action":"pass","Package":"probefixture","Elapsed":0.01}`, 0, "")
	quiet, code := runProbe(t, probeArgs()...)
	want := probeRow(t, "invalid", "clamp.go", "swap", "no-test-run", 0, "yes") + selectionRow(t, "package", "./", "^TestClampNegative$", 0)
	if code != 1 || !strings.HasPrefix(quiet, want) {
		t.Fatalf("stdout = (%q, %d), want %q first and 1", quiet, code, want)
	}
}

// PB9: --full reaches the focused run, so a long diagnostic is not previewed.
func TestProbeForwardsFullToTheFocusedRun(t *testing.T) {
	f := newFixture(t)
	long := strings.Repeat("d", bounds.PreviewRuneLimit+40)
	installStubGo(t, f, cannedFailure(long), 1, "")
	args := []string{"clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$"}
	previewed, _ := runProbe(t, args...)
	if strings.Contains(previewed, long) {
		t.Fatalf("without --full the diagnostic was not previewed:\n%s", previewed)
	}
	full, _ := runProbe(t, append(args, "--full")...)
	if !strings.Contains(full, long) {
		t.Fatalf("--full did not reach the focused run:\n%s", full)
	}
}

// PB10: the subject resolves against the working directory, so a probe runs from a
// subdirectory of the repository.
func TestProbeResolvesTheSubjectFromTheWorkingDirectory(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	t.Chdir(filepath.Join(f.root, "cmd"))
	out, code := runProbe(t, "../clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$")
	if code != 0 {
		t.Fatalf("exit = %d, want 0\n%s", code, out)
	}
	if want := probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "yes"); !strings.HasPrefix(out, want) {
		t.Fatalf("stdout = %q, want %q first", out, want)
	}
}
