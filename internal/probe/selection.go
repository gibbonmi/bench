package probe

import (
	"regexp"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// prepareSelection prepares the focused selection of parsed, narrowed to expects, and then
// grades the prose subject. It answers the request, or the refusal line and its exit.
func prepareSelection(root string, parsed usage.Result, subject subject, expects []string) (testreport.Request, string, int) {
	request, line, code := testreport.Prepare(root, selectionArgs(parsed, expects))
	if line != "" {
		return request, line, code
	}
	if line := gradeProseSubject(subject, request); line != "" {
		return request, line, 1
	}
	return request, "", 0
}

// gradeProseSubject refuses a prose probe of a file that the prose grade never reads,
// because the grade of that mutant equals the baseline and would answer a false silent.
func gradeProseSubject(subject subject, request testreport.Request) string {
	if !request.GradesProse() || prose.IsSubjectName(subject.display) {
		return ""
	}
	return toon.Errorf("probe subject not graded", subject.display+" is not a prose subject") + "\n"
}

// selectionArgs spells the focused run in `bench test`'s own grammar, so the probe and
// the verb it composes can never disagree about what a selection means. Flags keeps only
// the last --package value, so the package form reads Repeated. A caller's --run governs;
// without one, expectRun narrows the package form to expects.
func selectionArgs(parsed usage.Result, expects []string) []string {
	var args []string
	check, isCheck := parsed.Flags["--check"]
	if isCheck {
		args = append(args, "--check", check)
	} else {
		for _, pkg := range parsed.Repeated["--package"] {
			args = append(args, "--package", pkg)
		}
	}
	run, ok := parsed.Flags["--run"]
	if !ok && !isCheck {
		run = expectRun(expects)
	}
	if run != "" {
		args = append(args, "--run", run)
	}
	if _, ok := parsed.Flags["--full"]; ok {
		args = append(args, "--full")
	}
	return args
}

// expectRunNote is the help note of the narrowing that expectRun owns, so the help states the
// selection beside the code that makes it.
const expectRunNote = "--expect <test> narrows a --package selection with no --run to that test and its subtests, for the baseline and for each mutated run."

// expectRun answers the run pattern that selects each test that expects names, with its
// subtests, or an empty pattern for the whole selection. Each slash-separated part of a name
// matches one level of the test name literally and exactly, so TestFoo does not select
// TestFooBar. An empty expect is a mutant with no witness, so it needs the whole selection.
func expectRun(expects []string) string {
	if len(expects) == 0 || slices.Contains(expects, "") {
		return ""
	}
	var alternatives []string
	for _, expect := range expects {
		parts := strings.Split(expect, "/")
		for i, part := range parts {
			parts[i] = "^" + regexp.QuoteMeta(part) + "$"
		}
		if pattern := strings.Join(parts, "/"); !slices.Contains(alternatives, pattern) {
			alternatives = append(alternatives, pattern)
		}
	}
	return strings.Join(alternatives, "|")
}
