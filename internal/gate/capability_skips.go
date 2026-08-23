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

// skipRows renders the tally. The totals row is unconditional; a run with nothing to
// report says so, because absent output and zero skips must not read alike. Classes
// follow the package's declared order so the rows are stable run to run. The
// environment row names every skip it counts, because this population is the one the
// gate reds on. A reader who must go find which test stopped running has a number, not
// a diagnosis.
func skipRows(tally skipTally) []string {
	rows := []string{fmt.Sprintf("%s: %d (capability=%d environment=%d)", skipRowPrefix, tally.capability+len(tally.environment), tally.capability, len(tally.environment))}
	for _, class := range capability.Classes() {
		if count := tally.byClass[class]; count > 0 {
			rows = append(rows, fmt.Sprintf("%s class=%s: %d", skipRowPrefix, class, count))
		}
	}
	if len(tally.environment) > 0 {
		rows = append(rows, fmt.Sprintf("%s class=environment: %d (%s)", skipRowPrefix, len(tally.environment), strings.Join(namedReasons(tally.environment), ", ")))
	}
	return rows
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

// environmentFailure is the red message for a check the oracle asked for and did not
// get. Unlike a capability skip, an environment skip is never a fact about the host. It
// says a test found its staging absent, and inside the gate, the gate itself stages it.
// Marking this informational lets the kit's own conformance suite go unenforced
// behind a green verdict. So the posture is unconditional: there is no developer-host
// exemption to grant, because the missing staging is the gate's own.
func environmentFailure(tally skipTally) string {
	if len(tally.environment) == 0 {
		return ""
	}
	return fmt.Sprintf("gate: the oracle staged no environment for %d skipped check(s), so their verdict is absent, not green: %s", len(tally.environment), strings.Join(namedReasons(tally.environment), "; "))
}

// strict reports whether this run treats an incomplete capability population as red.
func strict() bool { return os.Getenv(requireCapabilitiesEnv) == "1" }

// strictFailure is the red message for a strict run, naming the classes that did not
// run so the verdict is actionable. Environment skips never contribute: an absent
// subject binary is a staging fact, not a security class.
func strictFailure(tally skipTally) string {
	if !strict() || tally.capability == 0 {
		return ""
	}
	var classes []string
	for _, class := range capability.Classes() {
		if tally.byClass[class] > 0 {
			classes = append(classes, string(class))
		}
	}
	return fmt.Sprintf("gate: capability skips are fatal under %s=1: %s", requireCapabilitiesEnv, strings.Join(classes, ", "))
}

// reportCapabilitySkips prints the rows and reports whether strict mode makes the run
// red on their account. reportCapabilitySkips diagnoses an unreadable log on every run,
// but turns the run red only under strict mode. This matches the fail posture of the
// skips themselves: a developer's host is not held to a population it never promised.
// An unconditional red on a transient read failure would make the gate unusable
// locally.
func reportCapabilitySkips(path string, stdout, stderr io.Writer) bool {
	tally, readErr := readSkipTally(path)
	for _, row := range skipRows(tally) {
		fmt.Fprintln(stdout, row)
	}
	red := false
	if readErr != nil {
		fmt.Fprintf(stderr, "gate: capability skip log %s is unreadable, so the counts above prove nothing: %v\n", path, readErr)
		if strict() {
			fmt.Fprintf(stderr, "gate: an unreadable skip log is fatal under %s=1\n", requireCapabilitiesEnv)
			red = true
		}
	}
	if failure := environmentFailure(tally); failure != "" {
		fmt.Fprintln(stderr, failure)
		red = true
	}
	if failure := strictFailure(tally); failure != "" {
		fmt.Fprintln(stderr, failure)
		red = true
	}
	return red
}
