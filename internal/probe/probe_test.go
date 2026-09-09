package probe

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

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

// stubStart is what the stub `go` does on one `go test` start: shell it runs first, the
// event stream it prints, and the exit it takes. A probe starts two of them, so a row
// spells the baseline's start apart from the mutated run's.
type stubStart struct {
	hook   string
	events string
	exit   int
}

// installStubStarts puts a canned `go` ahead of the real one and records every start in the
// fixture's marker. A `list` call still reaches the real toolchain, because the run binary's
// freshness verification resolves the module closure through it. The marker is dropped
// first, so each install restarts the count: an odd start is a baseline and the even start
// after it is the mutated run. Both hooks read the start's ordinal in `$n`.
func installStubStarts(t *testing.T, f *fixture, baseline, mutated stubStart) {
	t.Helper()
	real, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(f.marker); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	marker := sanitize.ShellQuote(f.marker)
	dir := scratchDir(t)
	script := "#!/usr/bin/env bash\n" +
		"printf '%s\\n' \"$1\" >> " + marker + "\n" +
		"if [ \"$1\" = list ]; then exec " + sanitize.ShellQuote(real) + " \"$@\"; fi\n" +
		"n=$(grep -c '^test$' " + marker + ")\n" +
		"if [ $((n % 2)) = 1 ]; then\n" + stubAnswer(baseline) + "fi\n" + stubAnswer(mutated)
	writeFixtureFile(t, filepath.Join(dir, "go"), script, 0o755)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func stubAnswer(start stubStart) string {
	return start.hook + "cat <<'PROBEEOF'\n" + start.events + "\nPROBEEOF\nexit " + strconv.Itoa(start.exit) + "\n"
}

// installStubGo is the ordinary shape: the baseline passes, so the probe reaches the
// mutated run, and that run answers what the row needs. hook is extra shell the mutated
// run's stub runs while that run is in flight.
func installStubGo(t *testing.T, f *fixture, events string, exit int, hook string) {
	t.Helper()
	installStubStarts(t, f, stubStart{events: cannedPass}, stubStart{hook: hook, events: events, exit: exit})
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

// quietRun is a package that passed with no test event at all, which is the stream a run
// that started nothing prints.
const quietRun = `{"Action":"pass","Package":"probefixture","Elapsed":0.01}`

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
func selectionRow(t *testing.T, form, target, run, baseline string, ran int) string {
	t.Helper()
	out, err := toon.TableTyped("selection", selectionFields, [][]any{{form, target, run, baseline, ran}})
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

// requireRow asserts the exit and the verdict row the answer opens with, which is the shape
// almost every row below grades.
func requireRow(t *testing.T, out string, code, wantCode int, want string) {
	t.Helper()
	if code != wantCode || !strings.HasPrefix(out, want) {
		t.Fatalf("stdout = (%q, %d), want %q first and %d", out, code, want, wantCode)
	}
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

// preservedCopy answers the one file below the home's probe directory, which preserve
// writes at home/probe/<pool key>/<stamp>/<base name>.
func preservedCopy(t *testing.T, f *fixture) string {
	t.Helper()
	found, err := filepath.Glob(filepath.Join(f.home, "probe", "*", "*", "*"))
	if err != nil || len(found) != 1 {
		t.Fatalf("preserved files = (%v, %v), want exactly one", found, err)
	}
	return found[0]
}

// PB1, PB13, PB15: a mutation a focused test catches is the success case, the subject
// comes back byte-exact, and the home keeps no copy.
func TestProbeBitesWhenTheFocusedTestFails(t *testing.T) {
	f := newFixture(t)
	out, code := runProbe(t, probeArgs()...)
	requireRow(t, out, code, 0, probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "yes"))
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
	requireRow(t, out, code, 0, probeRow(t, "bit", "clamp.go", "omit", "failed", 1, "yes"))
	requireSubjectBytes(t, f, clampSource)
}

// PB3: a mutation no test observes is silent and red, because it proves nothing.
func TestProbeIsSilentWhenTheFocusedTestPasses(t *testing.T) {
	f := newFixture(t)
	out, code := runProbe(t, "clamp.go", "--swap", "keeps n at", "--with", "holds n at", "--package", "./", "--run", "^TestClampNegative$")
	requireRow(t, out, code, 1, probeRow(t, "silent", "clamp.go", "swap", "passed", 0, "yes"))
	requireSubjectBytes(t, f, clampSource)
}

// PB4: a mutated package that does not compile never bit, so it is invalid rather than a
// pass or a bite. The selection names no run pattern, because `bench test` attributes a
// build failure under a run pattern to the pattern.
func TestProbeIsInvalidWhenTheMutationDoesNotCompile(t *testing.T) {
	f := newFixture(t)
	out, code := runProbe(t, "clamp.go", "--swap", "return n", "--with", "return", "--package", "./")
	requireRow(t, out, code, 1, probeRow(t, "invalid", "clamp.go", "swap", "build-failed", 0, "yes"))
	requireSubjectBytes(t, f, clampSource)
}

// PB5: a mutated run that started no test cannot bite. The baseline passes and the mutated
// run is the quiet one, because a pattern that matches nothing now refuses at the baseline.
func TestProbeIsInvalidWhenNoTestRuns(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, quietRun, 0, "")
	out, code := runProbe(t, probeArgs()...)
	requireRow(t, out, code, 1, probeRow(t, "invalid", "clamp.go", "swap", "no-test-run", 0, "yes"))
	requireSubjectBytes(t, f, clampSource)
}

// PB6: the evidence after the row is what `bench test` prints for the same selection.
func TestProbeCarriesTheFocusedRunTables(t *testing.T) {
	f := newFixture(t)
	events := cannedFailure("clamp_test.go:6: Clamp(-1) = -1, want 0")
	installStubGo(t, f, events, 1, "")
	out, code := runProbe(t, probeArgs()...)
	if code != 0 {
		t.Fatalf("exit = %d, want 0\n%s", code, out)
	}
	// The reference report is a start of its own, so a fresh install answers the mutated
	// run's stream on either side of the baseline branch.
	installStubStarts(t, f, stubStart{events: events, exit: 1}, stubStart{events: events, exit: 1})
	report, reportCode := testreport.Command(f.root, []string{"--package", "./", "--run", "^TestClampNegative$"})
	if reportCode != 1 {
		t.Fatalf("focused run exit = %d, want 1\n%s", reportCode, report)
	}
	row := probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "yes")
	if want := row + selectionRow(t, "package", "./", "^TestClampNegative$", "passed", 1) + report; out != want {
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
	if want := selectionRow(t, "package", "./", "^TestClampNegative$", "passed", 2); code != 0 || !strings.Contains(out, want) {
		t.Fatalf("stdout = (%q, %d), want %q and 0", out, code, want)
	}
	installStubGo(t, f, quietRun, 0, "")
	quiet, code := runProbe(t, probeArgs()...)
	want := probeRow(t, "invalid", "clamp.go", "swap", "no-test-run", 0, "yes") +
		selectionRow(t, "package", "./", "^TestClampNegative$", "passed", 0)
	requireRow(t, quiet, code, 1, want)
}

// PB10: the subject resolves against the working directory, so a probe runs from a
// subdirectory of the repository.
func TestProbeResolvesTheSubjectFromTheWorkingDirectory(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	t.Chdir(filepath.Join(f.root, "cmd"))
	out, code := runProbe(t, probeArgs("../clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$")...)
	requireRow(t, out, code, 0, probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "yes"))
}
