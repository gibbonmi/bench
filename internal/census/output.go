package census

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gibbonmi/bench/internal/poolkey"
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
	if err := checkAssignmentID(assignment); err != nil {
		return err
	}
	return write(Dir(home, root), assignment+outputSuffix, composeOutput(now.UTC().Format(time.RFC3339), output))
}

// The fields of an output record after its time and its head, in their order. The count
// is the last value, so the writer and the reader share one statement of the layout.
const (
	outputLinesField = iota
	outputBytesField
	outputDispositionField
	outputLaterFields
)

// composeOutput renders one output record line through composeRecord, which owns the time
// and the head. This function and parseOutput place the later fields by the constants
// above, and no other call site states where a field sits.
func composeOutput(timestamp string, output Output) string {
	later := make([]string, outputLaterFields)
	later[outputLinesField] = strconv.Itoa(output.Lines)
	later[outputBytesField] = strconv.FormatInt(output.Bytes, 10)
	later[outputDispositionField] = outputInline
	if output.Spilled {
		later[outputDispositionField] = outputSpilled
	}
	return composeRecord(timestamp, output.Head, later...)
}

// parseOutput returns the head and the byte count one output record line carries. The
// line splits through splitRecord, the raw-call codec. A line of another shape, which
// only a foreign writer makes, carries neither.
func parseOutput(line string) (string, int64, bool) {
	_, head, later, ok := splitRecord(line)
	if !ok || len(later) != outputLaterFields {
		return "", 0, false
	}
	bytes, err := strconv.ParseInt(later[outputBytesField], 10, 64)
	if err != nil || bytes < 0 {
		return "", 0, false
	}
	return head, bytes, true
}

// OutputBreakdown renders one assignment's output records per verb head, in the shape
// `<head>=<calls>/<bytes>,...`. The heads sort by bytes, largest first, and a tie sorts
// by the head name, so the reader sees the costliest verb first. The reader has the
// posture of HeadBreakdown: no readable record renders no text.
func OutputBreakdown(home, root, assignment string) string {
	if !poolkey.IsAssignmentID(assignment) {
		return ""
	}
	calls, sizes := tally(Dir(home, root), assignment+outputSuffix, parseOutput)
	return renderBreakdown(sizes, func(head string) string { return fmt.Sprintf("%d/%d", calls[head], sizes[head]) })
}
