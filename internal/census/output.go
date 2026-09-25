package census

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// outputSuffix names an assignment's output record file beside its raw-call file. The
// suffix keeps the name outside the assignment id form, so the raw-call readers, which
// read only an assignment id name, never count an output record.
const outputSuffix = ".output"

// The disposition field of an output record.
const (
	outputSpilled = "spilled"
	outputInline  = "inline"
)

// Output is the size of one bounded Bench response: the verb head, the line and byte
// counts of the complete output, and whether a spill file holds that output.
type Output struct {
	Head    string
	Lines   int
	Bytes   int64
	Spilled bool
}

// RecordOutput appends one output record to the assignment's output file. An identifier
// that is not an assignment id is refused rather than composed into a path. The caller
// ignores the error, so a failed write never changes a response or its exit code.
func RecordOutput(home, root, assignment string, output Output, now time.Time) error {
	if !poolkey.IsAssignmentID(assignment) {
		return fmt.Errorf("census assignment id is malformed: %s", sanitize.Controls(assignment))
	}
	return write(Dir(home, root), assignment+outputSuffix, composeOutput(now.UTC().Format(time.RFC3339), output))
}

// composeOutput renders one output record line. The time and the head keep the raw-call
// layout, and the counts and the disposition follow them. This function and parseOutput
// are the one source of the output line's layout.
func composeOutput(timestamp string, output Output) string {
	disposition := outputInline
	if output.Spilled {
		disposition = outputSpilled
	}
	return composeRecord(timestamp, output.Head, strconv.Itoa(output.Lines), strconv.FormatInt(output.Bytes, 10), disposition)
}

// outputFields is the field count of one output record line.
const outputFields = 5

// parseOutput returns the head and the byte count one output record line carries. A line
// of another shape, which only a foreign writer makes, carries neither.
func parseOutput(line string) (string, int64, bool) {
	fields := strings.Split(line, recordSeparator)
	if len(fields) != outputFields || fields[1] == "" {
		return "", 0, false
	}
	bytes, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil || bytes < 0 {
		return "", 0, false
	}
	return fields[1], bytes, true
}

// OutputBreakdown renders one assignment's output records per verb head, in the shape
// `<head>=<calls>/<bytes>,...`. The heads sort by bytes, largest first, and a tie sorts
// by the head name, so the reader sees the costliest verb first. The reader has the
// posture of HeadBreakdown: no readable record renders no text.
func OutputBreakdown(home, root, assignment string) string {
	if !poolkey.IsAssignmentID(assignment) {
		return ""
	}
	text, ok := readRecords(Dir(home, root), assignment+outputSuffix)
	if !ok {
		return ""
	}
	calls, sizes := map[string]int{}, map[string]int64{}
	for _, line := range recordLines(text) {
		head, bytes, ok := parseOutput(line)
		if !ok {
			continue
		}
		calls[head]++
		sizes[head] += bytes
	}
	names := slices.Collect(maps.Keys(calls))
	slices.SortFunc(names, func(a, b string) int {
		if sizes[a] != sizes[b] {
			return cmp.Compare(sizes[b], sizes[a])
		}
		return strings.Compare(a, b)
	})
	parts := make([]string, 0, len(names))
	for _, name := range names {
		parts = append(parts, fmt.Sprintf("%s=%d/%d", escapeDelimiters(sanitize.Controls(name)), calls[name], sizes[name]))
	}
	return strings.Join(parts, ",")
}
