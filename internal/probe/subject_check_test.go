package probe

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/probe/probetest"
	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testreport"
)

// proseNote is a Markdown subject that the prose grade passes, and longSentence is the
// replacement that puts one sentence over the word bound.
const (
	proseNote    = "# Note\n\nThe probe grades this note.\n"
	shortLine    = "The probe grades this note."
	longSentence = "The probe grades this note with one sentence that keeps on going past the bound, because every added word counts toward the limit that the prose rule states for each sentence in a file."
)

// PS26, PS27: a Markdown mutation under the prose check bites, and failed_tests counts
// the prose findings over the mutant.
func TestProbeTakesTheProseCheck(t *testing.T) {
	f := newFixture(t)
	// An absent exclusion list is a prose finding of its own, and an empty one excludes nothing.
	writeFixtureFile(t, filepath.Join(f.root, filepath.FromSlash(prose.ExclusionFile)), "", 0o644)
	trackFixtureFile(t, f, "note.md", strings.Replace(proseNote, shortLine, longSentence, 1))
	findings := len(prose.Grade(f.root))
	trackFixtureFile(t, f, "note.md", proseNote)
	out, code := runProbe(t, "note.md", "--swap", shortLine, "--with", longSentence, "--check", "prose")
	want := probeRow(t, "bit", "note.md", "swap", "failed", findings, "yes") +
		selectionRow(t, "check", "prose", testreport.AllTests, "passed", 0)
	requireRow(t, out, code, 0, want)
	if findings == 0 {
		t.Fatal("the long sentence gave no prose finding")
	}
	body, err := os.ReadFile(filepath.Join(f.root, "note.md"))
	if err != nil || string(body) != proseNote {
		t.Fatalf("note.md = (%q, %v), want the start bytes", body, err)
	}
	requireHomeEmpty(t, f)
}

// systemSelection is the system-suite selection of one test.
var systemSelection = []string{"--check", gate.SystemPhaseName, "--run", "^TestClampNegative$"}

// PS28, PS29: a system probe passes the baseline and runs the mutant through the system
// suite's argv, and over a root that is not the kit the baseline refuses.
func TestProbeTakesTheSystemCheck(t *testing.T) {
	t.Run("kit", func(t *testing.T) {
		f := newFixture(t)
		argv := filepath.Join(scratchDir(t), "argv")
		hook := "printf '%s\\n' \"$*\" > " + sanitize.ShellQuote(argv) + "\n"
		installStubStarts(t, f, probetest.Start{Events: cannedPass}, probetest.Start{Hook: hook, Events: cannedFailure("caught"), Exit: 1})
		out, code := runProbe(t, append(probeArgs("clamp.go", "--swap", "n < 0", "--with", "n > 0"), systemSelection...)...)
		want := probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "yes") +
			selectionRow(t, "check", gate.SystemPhaseName, "^TestClampNegative$", "passed", 1)
		requireRow(t, out, code, 0, want)
		operands, _ := gate.SystemSuite(f.root)
		started, err := os.ReadFile(argv)
		if wantArgv := " " + strings.Join(operands, " ") + " -run ^TestClampNegative$"; err != nil || !strings.Contains(string(started), wantArgv) {
			t.Fatalf("mutated run argv = (%q, %v), want %q", started, err, wantArgv)
		}
		requireSubjectBytes(t, f, clampSource)
	})
	t.Run("not the kit", func(t *testing.T) {
		f := refusalFixture(t)
		t.Setenv("BENCH_KIT", scratchDir(t))
		report, _ := testreport.Command(f.root, systemSelection)
		out, code := runProbe(t, append(probeArgs("clamp.go", "--swap", "n < 0", "--with", "n > 0"), systemSelection...)...)
		want := probeRow(t, "invalid", "clamp.go", "swap", "baseline-refused", 0, "untouched") +
			selectionRow(t, "check", gate.SystemPhaseName, "^TestClampNegative$", "refused", 0) +
			executionRow(t, 1, nil, len(report)) + report
		requireRefused(t, f, out, code, want, 1)
	})
}

// PS33: without an inherited selection, the mutated run builds its run binary after the
// mutation write, so its build script reads the mutant bytes.
func TestProbeBuildsTheMutantRunBinary(t *testing.T) {
	f := newFixture(t)
	records := buildEachRun(t, f, "clamp.go")
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	out, code := runProbe(t, append(probeArgs("clamp.go", "--swap", "n < 0", "--with", "n > 0"), systemSelection...)...)
	requireRow(t, out, code, 0, probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "yes"))
	want := []string{clampSource, strings.Replace(clampSource, "n < 0", "n > 0", 1)}
	if got := buildRecords(t, records); !reflect.DeepEqual(got, want) {
		t.Fatalf("build script reads = %q, want the start bytes and then the mutant", got)
	}
	requireSubjectBytes(t, f, clampSource)
}

// PS36, PS37: the selection parse precedes the substring miss. An unknown check answers
// its own refusal, and a miss reads its selection cells from the prepared request.
func TestSelectionPrecedesTheSubstringMiss(t *testing.T) {
	t.Run("unknown check", func(t *testing.T) {
		f := refusalFixture(t)
		selection := []string{"--check", "no-such-check"}
		want, wantCode := testreport.Command(f.root, selection)
		out, code := runProbe(t, append([]string{"clamp.go", "--omit", "absent"}, selection...)...)
		requireRefused(t, f, out, code, want, wantCode)
		if code != 2 {
			t.Fatalf("exit = %d, want 2", code)
		}
	})
	for _, selection := range [][]string{{"--check", "line-routing"}, systemSelection} {
		t.Run(strings.Join(selection, " "), func(t *testing.T) {
			f := refusalFixture(t)
			request, line, _ := testreport.Prepare(f.root, selection)
			if line != "" {
				t.Fatalf("prepare refused %v: %q", selection, line)
			}
			out, code := runProbe(t, append([]string{"clamp.go", "--omit", "absent"}, selection...)...)
			want := probeRow(t, "invalid", "clamp.go", "omit", "substring-miss", 0, "untouched") +
				selectionRow(t, request.Form(), request.Target(), request.Run(), "", 0) + executionRow(t, 1, nil, 0)
			requireRefused(t, f, out, code, want, 1)
		})
	}
}
