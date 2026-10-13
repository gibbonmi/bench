package testreport

import (
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gate"
)

// CheckForm and PackageForm are the two selection forms a prepared request carries. A
// caller prints the word, so it is spelled here rather than at each render site.
const (
	CheckForm   = "check"
	PackageForm = "package"
)

// AllTests is the run fact of a selection that passes Go no run pattern. An empty fact
// would read as a pattern that matched nothing, which is a different observation.
const AllTests = "all"

// Form answers which of the two selection forms the request carries.
func (r Request) Form() string {
	if r.focused.check != "" {
		return CheckForm
	}
	return PackageForm
}

// Target answers the check name or the package expressions the request selects. Several
// expressions read as one space-separated target, in argv order.
func (r Request) Target() string {
	if r.focused.check != "" {
		return r.focused.check
	}
	return strings.Join(r.Packages(), " ")
}

// Packages answers each package expression the request passes Go, in argv order.
func (r Request) Packages() []string {
	return slices.Clone(r.focused.packageExprs)
}

// Changed answers whether the request carries `--changed`. That form's target is `./...`,
// so the target alone cannot tell it from the whole-suite package form.
func (r Request) Changed() bool {
	return r.focused.changed
}

// Runs answers whether the request runs its selection. A listing form runs nothing,
// although its form and target read like a run.
func (r Request) Runs() bool {
	return r.focused.listing() == nil
}

// Run answers the exact run pattern the focused run passes to Go, or AllTests when it
// passes none. The named-check pattern comes from the one producer the argv reads, so a
// caller is never told a pattern the run did not use. The prose grade starts no Go child,
// so it answers AllTests. The system suite passes the caller's pattern, or none.
func (r Request) Run() string {
	switch {
	case r.focused.check == proseCheckName:
		return AllTests
	case r.focused.check != "" && !CheckTakesRun(r.focused.check):
		return namedCheckRunPattern()
	case r.focused.run != "":
		return r.focused.run
	}
	return AllTests
}

// GradesProse answers whether the request runs the prose grade, which reads only the
// subjects that prose.IsSubjectName accepts.
func (r Request) GradesProse() bool {
	return r.focused.check == proseCheckName
}

// CheckTakesRun answers whether the named check takes a --run pattern. Only the system
// suite does. Every other check runs a fixed pattern, so a caller's pattern there would be
// dropped in silence. The bench test grammar and the probe grammar both refuse by it.
func CheckTakesRun(check string) bool {
	return check == gate.SystemPhaseName
}

// namedCheckRunPattern is the one spelling of the test a named conformance check runs.
// The Go argv and the selection facts both read it here.
func namedCheckRunPattern() string {
	return "^" + registry.RootConformanceTest + "$"
}

// packageValueRule states the one-expression rule of each --package value. The run passes
// the value to Go as one typed argument, so a space-separated list fails as one missing
// directory. The bench test help and the probe help each add the remedy their grammar takes.
const packageValueRule = "--package <expr> takes one Go package expression as one argument, so a space-separated list is not split"

// PackageExpressionNote is the bench test help's package note. Its grammar repeats --package.
const PackageExpressionNote = packageValueRule + "; repeat --package to run several expressions in one call, or use a parent ./... pattern."

// probePackageNote is the probe help's package note. Its grammar takes one --package.
const probePackageNote = packageValueRule + "; run one call per package or use a parent ./... pattern."

// ProbeNotes states what each of the two selection forms takes, so a caller names a
// focused run correctly. The probe's help prints this block, so the owner's behavior and
// its description stay one source.
func ProbeNotes() string {
	return "notes:\n  " + strings.Join([]string{
		probePackageNote,
		"--check <name> names a check from the bench test --help inventory.",
	}, "\n  ")
}
