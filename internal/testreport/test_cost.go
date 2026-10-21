package testreport

import (
	"math"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/toon"
)

// slowestPerPackage is the row limit of each package in the slowest table.
const slowestPerPackage = 3

// timeoutWarningDivisor divides the go test timeout into the warning threshold, so the
// threshold is never a second duration.
const timeoutWarningDivisor = 2

// timeoutWarningsBlock renders a row for each package, in the given order, whose elapsed
// time is at least the warning threshold. It renders nothing when no package qualifies,
// so a healthy run keeps its bytes.
func (r *report) timeoutWarningsBlock(packages []string) (string, error) {
	timeoutMS := bounds.GoTestTimeout.Milliseconds()
	rows := make([][]any, 0)
	for _, pkg := range packages {
		if r.elapsedMS[pkg] >= timeoutMS/timeoutWarningDivisor {
			rows = append(rows, []any{pkg, r.elapsedMS[pkg], timeoutMS})
		}
	}
	if len(rows) == 0 {
		return "", nil
	}
	return toon.TableTyped("timeout_warnings", []string{"package", "elapsed_ms", "timeout_ms"}, rows)
}

// milliseconds converts a go test -json `Elapsed` value in seconds to whole milliseconds.
func milliseconds(seconds float64) int64 {
	return int64(math.Round(seconds * 1000))
}

// slowestBlock renders the slowest top-level tests of each package. A candidate test has
// no `/` in its name and a last terminal action of `pass` or `fail`.
func (r *report) slowestBlock() (string, error) {
	candidates := make([]*testResult, 0)
	for _, test := range r.tests {
		if test.finished && !test.skipped && !strings.Contains(test.test, "/") {
			candidates = append(candidates, test)
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		if a.packageName != b.packageName {
			return a.packageName < b.packageName
		}
		if a.elapsedMS != b.elapsedMS {
			return a.elapsedMS > b.elapsedMS
		}
		return a.test < b.test
	})
	rows := make([][]any, 0)
	kept := 0
	for i, test := range candidates {
		if i == 0 || test.packageName != candidates[i-1].packageName {
			kept = 0
		}
		if kept == slowestPerPackage {
			continue
		}
		kept++
		rows = append(rows, []any{test.packageName, test.test, test.elapsedMS})
	}
	return toon.TableTyped("slowest", []string{"package", "test", "elapsed_ms"}, rows)
}
