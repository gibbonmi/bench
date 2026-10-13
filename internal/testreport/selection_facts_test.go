package testreport

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/usage"
)

// TestSelectionFactsNameTheRun grades the three facts a prepared request answers for each
// selection form. The check form's pattern is the one the named-check argv passes, and a
// selection that passes Go no pattern answers `all` rather than an empty fact.
func TestSelectionFactsNameTheRun(t *testing.T) {
	root := t.TempDir()
	for _, tc := range []struct {
		name              string
		args              []string
		form, target, run string
	}{
		{"conformance check", []string{"--check", "line-routing"}, CheckForm, "line-routing", "^" + rootConformanceTestName + "$"},
		{"system check", []string{"--check", gate.SystemPhaseName}, CheckForm, gate.SystemPhaseName, AllTests},
		{"prose check", []string{"--check", proseCheckName}, CheckForm, proseCheckName, AllTests},
		{"package", []string{"--package", "./internal/..."}, PackageForm, "./internal/...", AllTests},
		{"package and run", []string{"--package", "./...", "--run", "^TestClamp$"}, PackageForm, "./...", "^TestClamp$"},
		{"repeated package", []string{"--package", "./a", "--package", "./b/..."}, PackageForm, "./a ./b/...", AllTests},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request, line, code := Prepare(root, tc.args)
			if line != "" {
				t.Fatalf("Prepare refused %v: (%d, %q)", tc.args, code, line)
			}
			if request.Form() != tc.form || request.Target() != tc.target || request.Run() != tc.run {
				t.Fatalf("facts = (%q, %q, %q), want (%q, %q, %q)",
					request.Form(), request.Target(), request.Run(), tc.form, tc.target, tc.run)
			}
		})
	}
}

// FTPKG2: a repeated --package keeps every exclusion of one --package. The refusal
// answers from the grammar, before any selection or Go child.
func TestRepeatedPackageKeepsTheSelectionExclusions(t *testing.T) {
	for _, args := range [][]string{
		{"--package", "./a", "--package", "./b", "--changed"},
		{"--package", "./a", "--package", "./b", "./c"},
		{"--check", "line-routing", "--package", "./a", "--package", "./b"},
	} {
		if output, code := Command(t.TempDir(), args); code != 2 || !strings.HasPrefix(output, "usage: bench test") {
			t.Errorf("Command(%v) = %d, %q; want usage exit 2", args, code, output)
		}
	}
}

// TestSelectionFactsNameTheChangedForm grades the accessor that tells the two forms apart
// whose target is `./...`: the `--changed` form and the whole-suite package form.
func TestSelectionFactsNameTheChangedForm(t *testing.T) {
	root := t.TempDir()
	for _, tc := range []struct {
		args    []string
		changed bool
	}{
		{[]string{"--changed"}, true},
		{[]string{"--package", "./..."}, false},
	} {
		request, line, code := Prepare(root, tc.args)
		if line != "" {
			t.Fatalf("Prepare refused %v: (%d, %q)", tc.args, code, line)
		}
		if request.Changed() != tc.changed {
			t.Fatalf("Changed() of %v = %t, want %t", tc.args, request.Changed(), tc.changed)
		}
	}
}

// TestSelectionFactsNameTheListingForms grades the accessor that tells a selection run from
// the two forms that list and run nothing: the check inventory and a check's fixtures.
func TestSelectionFactsNameTheListingForms(t *testing.T) {
	root := t.TempDir()
	for _, tc := range []struct {
		args []string
		runs bool
	}{
		{[]string{"--checks"}, false},
		{[]string{"--check", RootConformanceCheckName, "--fixtures"}, false},
		{[]string{"--check", RootConformanceCheckName}, true},
		{[]string{"--package", "./..."}, true},
	} {
		request, line, code := Prepare(root, tc.args)
		if line != "" {
			t.Fatalf("Prepare refused %v: (%d, %q)", tc.args, code, line)
		}
		if request.Runs() != tc.runs {
			t.Fatalf("Runs() of %v = %t, want %t", tc.args, request.Runs(), tc.runs)
		}
	}
}

// rootConformanceTestName is the test a named conformance check runs, spelled here
// independently of the producer. A producer that names another test reds this row.
const rootConformanceTestName = "TestRootConformance"

// TestNamedCheckArgvCarriesTheRunPattern holds the argv and the fact to one producer, so
// the row a caller reads is the pattern the Go child received.
func TestNamedCheckArgvCarriesTheRunPattern(t *testing.T) {
	argv := focusedTestArgv("./internal/conformance", "-run", namedCheckRunPattern())
	request, line, _ := Prepare(t.TempDir(), []string{"--check", "line-routing"})
	if line != "" {
		t.Fatalf("Prepare refused the check form: %q", line)
	}
	if !strings.Contains(strings.Join(argv, " "), " -run "+request.Run()) {
		t.Fatalf("argv = %v, want the request's run fact %q", argv, request.Run())
	}
}

// TestOutcomeCountsTheDistinctTestsThatRan grades the run count the caller cites as
// evidence. A repeated event for one test counts once, and a stream with no run event
// counts none.
func TestOutcomeCountsTheDistinctTestsThatRan(t *testing.T) {
	for _, tc := range []struct {
		name   string
		stream string
		ran    int
	}{
		{"two tests", runEvent("TestOne") + runEvent("TestTwo") + packagePass, 2},
		{"one test twice", runEvent("TestOne") + runEvent("TestOne") + packagePass, 1},
		{"no test", packagePass, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			report, err := decode(strings.NewReader(tc.stream))
			if err != nil {
				t.Fatal(err)
			}
			if got := report.outcome().Ran; got != tc.ran {
				t.Fatalf("Ran = %d, want %d", got, tc.ran)
			}
		})
	}
}

func runEvent(test string) string {
	return `{"Action":"run","Package":"counted","Test":"` + test + `"}` + "\n" +
		`{"Action":"pass","Package":"counted","Test":"` + test + `","Elapsed":0}` + "\n"
}

const packagePass = `{"Action":"pass","Package":"counted","Elapsed":0.01}` + "\n"

// TestProbeNotesDeriveTheirTokens grades that the notes hold exactly the probe's package
// note and the check-inventory note. The prose check and the system suite are probe
// targets, and each run without an inherited selection builds its own run binary, so a
// not-a-target clause or a rebuild note states a route that the probe does not take.
func TestProbeNotesDeriveTheirTokens(t *testing.T) {
	notes := ProbeNotes()
	if !strings.HasPrefix(notes, "notes:\n  "+probePackageNote+"\n  --check <name> ") {
		t.Fatalf("notes = %q, want the package note and then the check-inventory note", notes)
	}
	if lines := strings.Count(notes, "\n  "); lines != 2 {
		t.Fatalf("notes hold %d facts, want 2:\n%s", lines, notes)
	}
	for _, retired := range []string{"not probe targets", usage.WorktreeBuild} {
		if strings.Contains(notes, retired) {
			t.Fatalf("notes = %q, want no %q", notes, retired)
		}
	}
}

// TestHelpStatesThePackageExpressionRule grades that bench test --help names the
// one-expression rule and the repeated flag before the check inventory. Two drain authors
// passed a space-separated package list and read only a missing-directory failure. The
// repeat clause is spelled here apart from the producer, so a note that drops it reds.
func TestHelpStatesThePackageExpressionRule(t *testing.T) {
	output, code := Command(t.TempDir(), []string{"--help"})
	want := "\nnotes:\n  " + PackageExpressionNote + "\nchecks:\n"
	if code != 0 || !strings.Contains(output, want) || !strings.Contains(PackageExpressionNote, "repeat --package to run several expressions in one call") {
		t.Fatalf("help = (%q, %d), want the package note %q with the repeated flag before the inventory", output, code, want)
	}
}
