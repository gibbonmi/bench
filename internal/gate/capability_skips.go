package gate

// The capability-skip side channel. The gate points every phase it launches at one
// run-scoped log file, then reads that file back once the phases have exited and
// reports what they skipped. A phase's own output stream cannot carry this: `go test`
// without -v discards a passing or skipping package's stdout and stderr entirely, so
// a collector teeing the stream would observe nothing on every green run.
//
// Concurrency lives at write time and is cross-process — many test binaries append to
// the same path, which is why capability.Capability issues one atomic append per
// line. Reading happens once, single-threaded, after the phases join.

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/capability"
)

// requireCapabilitiesEnv turns a nonzero capability-skip count red. Absent, or any
// value other than "1", the rows stay informational: a developer's host legitimately
// lacks capabilities, and an unconditional red would make the gate unusable locally.
const requireCapabilitiesEnv = "BENCH_REQUIRE_CAPABILITIES"

// skipRowPrefix identifies separately parseable skip-accounting rows so reports retain
// them alongside phase verdicts without weaving them into verdict text.
const skipRowPrefix = "capability-skips"

type skipTally struct {
	byClass     map[capability.Class]int
	capability  int
	environment int
	// environmentReasons keeps each environment skip's reason text, because the counts
	// alone cannot tell a stripping-induced degradation from the ordinary environment
	// skips every host emits. Only the stripped-subject posture reads them; the dev-tier
	// report stays a count.
	environmentReasons []string
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
// leaves nothing to append. Every other read failure is returned, because under strict
// mode the tally is enforcement — a log that exists but cannot be read proves nothing,
// and reporting it as zero would read exactly like a fully capable runner.
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
			tally.environment++
			tally.environmentReasons = append(tally.environmentReasons, skip.Reason)
		}
	}
	return tally, nil
}

// skipRows renders the tally. The totals row is unconditional — a run with nothing to
// report says so, because absent output and zero skips must not read alike. Classes
// follow the package's declared order so the rows are stable run to run.
func skipRows(tally skipTally) []string {
	rows := []string{fmt.Sprintf("%s: %d (capability=%d environment=%d)", skipRowPrefix, tally.capability+tally.environment, tally.capability, tally.environment)}
	for _, class := range capability.Classes() {
		if count := tally.byClass[class]; count > 0 {
			rows = append(rows, fmt.Sprintf("%s class=%s: %d", skipRowPrefix, class, count))
		}
	}
	return rows
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
// red on their account. An unreadable log is diagnosed on every run and turns the run red
// only under strict mode, matching the fail posture of the skips themselves: a developer's
// host is not held to a population it never promised, and an unconditional red on a
// transient read failure would make the gate unusable locally.
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
	if failure := strictFailure(tally); failure != "" {
		fmt.Fprintln(stderr, failure)
		red = true
	}
	return red
}
