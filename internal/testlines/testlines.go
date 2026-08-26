// Package testlines classifies the plain-text output lines of a go test run.
package testlines

import "strings"

// The prefixes go test writes its own lines under. Each prefix is named once here, so the
// two classifiers below read one declaration rather than a literal apiece.
const (
	failName    = "--- FAIL:"
	buildError  = "# "
	dataRace    = "WARNING: DATA RACE"
	packageFail = "FAIL"
	panicName   = "panic:"
)

// blockOpeners are the prefixes that open a failure block: a failed test name at any indent,
// a build error's package heading, and a race report's warning.
var blockOpeners = []string{failName, buildError, dataRace}

// runnerPrefixes are go test's bookkeeping lines, plus the skip line the gate itself writes.
var runnerPrefixes = []string{
	"bench-skip ",
	buildError,
	"=== RUN",
	"=== PAUSE",
	"=== CONT",
	"=== NAME",
	failName,
	"--- PASS:",
	"--- SKIP:",
	"ok ",
	"ok\t",
	"? ",
	"?\t",
}

// RunnerLine reports whether the line is go test's own bookkeeping rather than a diagnostic.
func RunnerLine(line string) bool {
	return hasAnyPrefix(line, runnerPrefixes) || line == packageFail || namedPackageFail(line)
}

// FailureRows returns the failure rows of one red plain-text stream, never -json output.
// A "--- FAIL:" line at any indent, a "# <package>" line, and a "WARNING: DATA RACE" line
// each open a block whose later lines are rows until the next runner line. A race report's
// stack precedes its "--- FAIL:" block, so the warning opener keeps that evidence.
// The caller owns the fallback for an empty result.
func FailureRows(lines []string) []string {
	rows := make([]string, 0)
	inBlock := false
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		switch {
		case hasAnyPrefix(line, blockOpeners):
			inBlock = true
			rows = append(rows, line)
		case namedPackageFail(line):
			inBlock = false
			rows = append(rows, line)
		case strings.HasPrefix(line, panicName):
			rows = append(rows, line)
		case RunnerLine(line):
			inBlock = false
		case inBlock:
			rows = append(rows, line)
		}
	}
	return rows
}

// namedPackageFail reports whether the line is the terminal FAIL line that names its package,
// which a bare "FAIL" does not.
func namedPackageFail(line string) bool {
	return strings.HasPrefix(line, packageFail+" ") || strings.HasPrefix(line, packageFail+"\t")
}

// hasAnyPrefix reports whether the line starts with one of the prefixes.
func hasAnyPrefix(line string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}
