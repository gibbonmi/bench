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
			if got := report.outcome(false).Ran; got != tc.ran {
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

// TestProbeNotesDeriveTheirTokens grades that the three notes carry the runner's own
// constants. A hand-copied token would drift the moment the runner renamed one.
func TestProbeNotesDeriveTheirTokens(t *testing.T) {
	notes := ProbeNotes()
	if !strings.HasPrefix(notes, "notes:\n  ") {
		t.Fatalf("notes = %q, want the notes block header", notes)
	}
	if lines := strings.Count(notes, "\n  "); lines != 3 {
		t.Fatalf("notes hold %d facts, want 3:\n%s", lines, notes)
	}
	for _, token := range []string{proseCheckName, gate.SystemPhaseName, usage.WorktreeBuild} {
		if !strings.Contains(notes, token) {
			t.Fatalf("notes = %q, want the runner's own %q", notes, token)
		}
	}
}
