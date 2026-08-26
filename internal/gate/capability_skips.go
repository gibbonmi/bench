package gate

// This file implements the capability-skip side channel. The gate points every
// phase at one run-scoped log file. Each phase appends what it skipped to that file.
//
// The gate reads the file once, after every phase exits, and reports the skips.
// A phase's own output stream cannot carry this fact. `go test` without -v discards
// a passing or skipping package's stdout and stderr. A collector that tees the
// stream then observes nothing on a green run.
//
// Many test binaries append to the same log file, so writes are concurrent and
// cross-process. capability.Capability issues one atomic append per line for this
// reason. The gate reads the log once, after every phase joins, so reads stay
// single-threaded.

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/capability"
)

// requireCapabilitiesEnv turns a nonzero capability-skip count red. When the value is
// absent, or is anything other than "1", the rows stay informational. A developer's
// host legitimately lacks capabilities, so an unconditional red would make the gate
// unusable locally.
const requireCapabilitiesEnv = "BENCH_REQUIRE_CAPABILITIES"

// skipRowPrefix opens every row this file writes. The canary matches EXPECT
// substrings against inner-gate output, so the gate adds these rows alongside the
// phase verdicts instead of weaving them in.
const skipRowPrefix = "capability-skips"

type skipTally struct {
	byClass    map[capability.Class]int
	capability int
	// environment keeps each environment skip whole, because a count cannot say which
	// assertion went unmade. This population is red inside the oracle. An actionless
	// verdict is the same failure mode the count already had.
	environment []capability.Skip
}

// newSkipLog creates the run-scoped file the phases append to. The path is absolute
// because a phase runs with its own working directory.
func newSkipLog() (path string, cleanup func(), err error) {
	file, err := os.CreateTemp("", "bench-skip-*.log")
	if err != nil {
		return "", func() {}, err
	}
	name := file.Name()
	if err := file.Close(); err != nil {
		_ = os.Remove(name)
		return "", func() {}, err
	}
	return name, func() { _ = os.Remove(name) }, nil
}

// withSkipLog points every phase at path. The phase table is shared, so this returns
// copies rather than editing the caller's Phase values in place.
func withSkipLog(phases []Phase, path string) []Phase {
	out := make([]Phase, 0, len(phases))
	for _, phase := range phases {
		phase.Env = append(append([]string(nil), phase.Env...), capability.LogEnv+"="+path)
		out = append(out, phase)
	}
	return out
}

// readSkipTally counts the structured skips the phases left behind. An absent log is an
// empty tally, not an error: a run whose phases all reported everything they could run
// leaves nothing to append. readSkipTally returns every other read failure, because
// under strict mode the tally is enforcement. A log that exists but cannot be read
// proves nothing, and reporting it as zero would read exactly like a fully capable
// runner.
func readSkipTally(path string) (skipTally, error) {
	tally := skipTally{byClass: map[capability.Class]int{}}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return tally, nil
		}
		return tally, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		skip, ok := capability.ParseLine(line)
		if !ok {
			continue
		}
		switch skip.Kind {
		case capability.KindCapability:
			tally.byClass[skip.Class]++
			tally.capability++
		case capability.KindEnvironment:
			tally.environment = append(tally.environment, skip)
		}
	}
	return tally, nil
}

// skipRow renders the whole tally as one line. The line is unconditional; a run with
// nothing to report says so, because absent output and zero skips must not read alike.
// It is one line rather than one per class because a green run's stdout is bounded, and
// a per-class line spends that bound on counts that are almost always zero.
//
// The parenthetical opens with the two kinds and then lists every nonzero class. A tally
// with no class at all closes after the kinds rather than printing a separator with
// nothing behind it.
//
// The line is a pure count. It names no environment skip, because each of those is its
// own row of the red failure table through environmentFailures, where the reader is
// asked to act on it. Naming them twice would put the same diagnosis in two places.
func skipRow(tally skipTally) string {
	row := fmt.Sprintf("%s: %d (capability=%d environment=%d", skipRowPrefix, tally.capability+len(tally.environment), tally.capability, len(tally.environment))
	counts := make([]string, 0, len(tally.byClass))
	for _, class := range nonzeroClasses(tally) {
		counts = append(counts, fmt.Sprintf("%s=%d", class, tally.byClass[class]))
	}
	if len(counts) > 0 {
		row += "; " + strings.Join(counts, " ")
	}
	return row + ")"
}

// nonzeroClasses answers the classes this tally counted, in the package's declared order
// so the answer is stable run to run. Both the count line and the strict-mode diagnosis
// read it, so which classes a run has to name is decided once.
func nonzeroClasses(tally skipTally) []capability.Class {
	var classes []capability.Class
	for _, class := range capability.Classes() {
		if tally.byClass[class] > 0 {
			classes = append(classes, class)
		}
	}
	return classes
}

// namedReasons renders each environment skip as "TestName: reason", the form both the
// row and the red message read.
func namedReasons(skips []capability.Skip) []string {
	named := make([]string, 0, len(skips))
	for _, skip := range skips {
		named = append(named, skip.Name+": "+skip.Reason)
	}
	return named
}

// environmentFailures is the red diagnosis for checks the oracle asked for and did not
// get, one row per skip. Unlike a capability skip, an environment skip is never a fact
// about the host. It says a test found its staging absent, and inside the gate, the gate
// itself stages it. Marking this informational lets the kit's own conformance suite go
// unenforced behind a green verdict. So the posture is unconditional: there is no
// developer-host exemption to grant, because the missing staging is the gate's own.
//
// One row per skip rather than one row naming them all: the table's own row count is
// then the count of absent verdicts, and a long population cannot arrive as one
// unreadable cell.
func environmentFailures(tally skipTally) []string {
	return namedReasons(tally.environment)
}

// strict reports whether this run treats an incomplete capability population as red.
func strict() bool { return os.Getenv(requireCapabilitiesEnv) == "1" }

// strictFailure is the red diagnosis for a strict run, naming the classes that did not
// run so the verdict is actionable. Environment skips never contribute: an absent
// subject binary is a staging fact, not a security class.
func strictFailure(tally skipTally) string {
	if !strict() || tally.capability == 0 {
		return ""
	}
	names := make([]string, 0, len(tally.byClass))
	for _, class := range nonzeroClasses(tally) {
		names = append(names, string(class))
	}
	return fmt.Sprintf("capability skips are fatal under %s=1: %s", requireCapabilitiesEnv, strings.Join(names, ", "))
}

// reportCapabilitySkips prints the one skip-count line, answers its red diagnoses as
// rows, and reports whether the run is red on their account. The diagnoses are rows
// rather than stderr lines so that one table holds everything a red run asks the reader
// to fix.
//
// A diagnosis and a red are separate answers here. reportCapabilitySkips diagnoses an
// unreadable log on every run but turns the run red only under strict mode. This matches
// the fail posture of the skips themselves: a developer's host is not held to a
// population it never promised, and an unconditional red on a transient read failure
// would make the gate unusable locally. A row the run is not red for costs nothing,
// because the table prints only when the run is already red for some other reason.
func reportCapabilitySkips(path string, stdout io.Writer) (rows []string, red bool) {
	tally, readErr := readSkipTally(path)
	fmt.Fprintln(stdout, skipRow(tally))
	if readErr != nil {
		rows = append(rows, fmt.Sprintf("capability skip log %s is unreadable, so the counts above prove nothing: %v", path, readErr))
		if strict() {
			red = true
		}
	}
	if failures := environmentFailures(tally); len(failures) > 0 {
		rows = append(rows, failures...)
		red = true
	}
	if failure := strictFailure(tally); failure != "" {
		rows = append(rows, failure)
		red = true
	}
	return rows, red
}
