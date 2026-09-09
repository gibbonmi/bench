package testreport

import (
	"strings"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/usage"
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

// Target answers the check name or the package expression the request selects.
func (r Request) Target() string {
	if r.focused.check != "" {
		return r.focused.check
	}
	return r.focused.packageExpr
}

// Run answers the exact run pattern the focused run passes to Go, or AllTests when it
// passes none. The named-check pattern comes from the one producer the argv reads, so a
// caller is never told a pattern the run did not use. The prose grade starts no Go child
// and the system suite passes no pattern, so both answer AllTests.
func (r Request) Run() string {
	switch {
	case r.focused.check == proseCheckName || r.focused.check == gate.SystemPhaseName:
		return AllTests
	case r.focused.check != "":
		return namedCheckRunPattern()
	case r.focused.run != "":
		return r.focused.run
	}
	return AllTests
}

// namedCheckRunPattern is the one spelling of the test a named conformance check runs.
// The Go argv and the selection facts both read it here.
func namedCheckRunPattern() string {
	return "^" + registry.RootConformanceTest + "$"
}

// ProbeNotes states the three facts a caller needs to name a focused run correctly: what
// each selection form takes, which two names are not probe targets, and why an edited
// tree needs a rebuild first. The probe's help prints this block, so the owner's behavior
// and its description stay one source.
func ProbeNotes() string {
	return "notes:\n  " + strings.Join([]string{
		"--package <expr> takes a Go package expression, as bench test --package does.",
		"--check <name> names a conformance check from the bench test --help inventory, and " +
			proseCheckName + " and " + gate.SystemPhaseName + " are not probe targets.",
		"A named check compiles from the run binary's source, so an edited tree needs " +
			usage.WorktreeBuild + " first.",
	}, "\n  ")
}
