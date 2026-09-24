package otelrecord

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"os"
	"strconv"
	"syscall"
	"time"
)

// The two seams a consumer of a landing's timings selects by. The reader names both
// here rather than importing the packages that open those spans, because this package
// owns the record and neither gate nor worktree may sit below it. A reader test
// reconciles the pair against Registry, so a renamed seam reds rather than drifting.
const (
	// SeamLanding is the seam a landing's own span states.
	SeamLanding = "worktree.land"

	// SeamGatePhase is the seam one gate phase's span states.
	SeamGatePhase = "gate.phase"
)

// recordLineLimit bounds one record line. A line longer than this is a record no
// consumer can trust, and the read refuses rather than truncating it into a span.
const recordLineLimit = 1 << 20

// Span is one finished span read back from a repository's record. The record writes a
// start line and an end line for each span, and only the end line carries the complete
// span, so the reader answers the end lines alone.
type Span struct {
	Name         string
	Seam         string
	TraceID      string
	SpanID       string
	ParentSpanID string
	Start        time.Time
	End          time.Time
	Attributes   map[string]string
}

// Elapsed is the span's wall time.
func (s Span) Elapsed() time.Duration { return s.End.Sub(s.Start) }

// ReadSpans returns the finished spans of one repository's record below an explicitly
// resolved home, in record order: the sealed segments by sequence, then the live one. An
// absent record, an unreadable record, and a segment the writer's grade refuses are all
// refusals, so a consumer states that it knows nothing rather than that nothing ran.
// A single malformed line is skipped: one truncated write must not hide the run that
// wrote every other line.
func ReadSpans(home, root string) ([]Span, error) {
	var spans []Span
	err := scanRecord(home, root, func(line []byte) bool {
		spans = append(spans, finishedSpans(line)...)
		return true
	})
	if err != nil {
		return nil, err
	}
	return spans, nil
}

// scanRecord calls visit with each line of root's record in record order, and stops when
// visit answers false. Each segment passes the writer's grade before it opens, and the
// open itself follows no link and never waits on a special file.
func scanRecord(home, root string, visit func(line []byte) bool) error {
	paths, err := segments(home, root)
	if err != nil {
		return err
	}
	for _, path := range paths {
		more, err := scanSegment(path, visit)
		if err != nil || !more {
			return err
		}
	}
	return nil
}

func scanSegment(path string, visit func(line []byte) bool) (bool, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return false, fmt.Errorf("open seam record: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64<<10), recordLineLimit)
	for scanner.Scan() {
		if !visit(scanner.Bytes()) {
			return false, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("read seam record: %w", err)
	}
	return true, nil
}

// Landing is the newest completed landing with a published subject in one
// repository's record and the gate stages that ran below it.
type Landing struct {
	// Commit is the published subject the landing recorded, empty when the landing
	// published none.
	Commit string
	// TraceID is the landing's trace, which is how a reader verifies the selection.
	TraceID string
	// Stages are the phase spans of that trace, in record order.
	Stages []Span
}

// NewestLanding returns the newest completed landing with a published subject from
// root's record. It answers false when the record is absent, is unreadable, or names
// no such landing, which is the same answer a consumer states as unknown.
func NewestLanding(home, root string) (Landing, bool) {
	spans, err := ReadSpans(home, root)
	if err != nil {
		return Landing{}, false
	}
	newest := -1
	for index, candidate := range spans {
		if candidate.Seam != SeamLanding || candidate.Attributes[AttrSubjectID] == "" {
			continue
		}
		if newest < 0 || candidate.End.After(spans[newest].End) {
			newest = index
		}
	}
	if newest < 0 {
		return Landing{}, false
	}
	landing := Landing{Commit: spans[newest].Attributes[AttrSubjectID], TraceID: spans[newest].TraceID}
	for _, candidate := range spans {
		if candidate.Seam == SeamGatePhase && candidate.TraceID == landing.TraceID {
			landing.Stages = append(landing.Stages, candidate)
		}
	}
	return landing, true
}

// finishedSpans decodes one record line and returns the finished spans it carries. A
// line the decoder refuses answers none, and so does a start line, whose omitted end
// time is what marks the span as still running.
func finishedSpans(line []byte) []Span {
	entries, _ := decodeRecord(line)
	var out []Span
	for _, entry := range entries {
		if entry.finished {
			out = append(out, entry.Span)
		}
	}
	return out
}

type decodedRecordSpan struct {
	Span
	finished bool
}

func decodeRecord(line []byte) ([]decodedRecordSpan, error) {
	var data tracesData
	if err := json.Unmarshal(line, &data); err != nil {
		return nil, err
	}
	if len(data.ResourceSpans) == 0 {
		return nil, fmt.Errorf("missing resource spans")
	}
	var out []decodedRecordSpan
	for _, resource := range data.ResourceSpans {
		if err := gradeSchema(resource.Resource); err != nil {
			return nil, err
		}
		for _, scope := range resource.ScopeSpans {
			for _, encoded := range scope.Spans {
				out = append(out, decodedRecordSpan{Span: decodeSpan(encoded), finished: encoded.EndTimeUnixNano != ""})
			}
		}
	}
	return out, nil
}

// gradeSchema refuses a line whose schema this reader does not know, so a consumer never
// misreads a later shape. A line with no schema key is a legacy line and reads as before.
func gradeSchema(block resource) error {
	for _, pair := range block.Attributes {
		if pair.Key != ResourceRecordSchema {
			continue
		}
		if schema := decodeValue(pair.Value); schema != RecordSchema {
			return fmt.Errorf("unknown record schema %q", schema)
		}
	}
	return nil
}

func decodeSpan(encoded span) Span {
	out := Span{
		Name:         encoded.Name,
		TraceID:      encoded.TraceID,
		SpanID:       encoded.SpanID,
		ParentSpanID: encoded.ParentSpanID,
		Start:        decodeUnixNano(encoded.StartTimeUnixNano),
		End:          decodeUnixNano(encoded.EndTimeUnixNano),
		Attributes:   make(map[string]string, len(encoded.Attributes)),
	}
	for _, pair := range encoded.Attributes {
		out.Attributes[pair.Key] = decodeValue(pair.Value)
	}
	out.Seam = out.Attributes[AttrSeam]
	return out
}

// decodeUnixNano reads back the quoted decimal the encoder writes. A value the encoder
// never wrote reads as the zero time, which no consumer mistakes for a measurement.
func decodeUnixNano(value string) time.Time {
	nanoseconds, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(0, nanoseconds)
}

// decodeValue renders one attribute value as text. Only the scalar arms are rendered,
// because Bench's declared attribute set carries no array value.
func decodeValue(value anyValue) string {
	switch {
	case value.StringValue != nil:
		return *value.StringValue
	case value.IntValue != nil:
		return *value.IntValue
	case value.BoolValue != nil:
		return strconv.FormatBool(*value.BoolValue)
	case value.DoubleValue != nil:
		return strconv.FormatFloat(*value.DoubleValue, 'g', -1, 64)
	}
	return ""
}

// ReadSelected streams only selected traces into memory, including unfinished spans.
// Malformed lines remain explicit coverage gaps instead of disappearing as empty results.
// A line number counts across the segments in record order.
func ReadSelected(home, root string, traceIDs []string) ([]Span, []string, error) {
	wanted := map[string]bool{}
	for _, id := range traceIDs {
		wanted[id] = true
	}
	var out []Span
	var problems []string
	var stopped error
	line := 0
	retained := int64(0)
	problem := func(reason string) bool {
		retained += int64(len(reason))
		if retained > bounds.ControlRecordLimit {
			return false
		}
		problems = append(problems, reason)
		return true
	}
	err := scanRecord(home, root, func(text []byte) bool {
		line++
		entries, err := decodeRecord(text)
		if err != nil {
			if !problem(fmt.Sprintf("line %d malformed", line)) {
				out, problems, stopped = nil, nil, fmt.Errorf("native diagnostics exceed control record bound")
				return false
			}
			return true
		}
		for _, entry := range entries {
			decoded := entry.Span
			if !wanted[decoded.TraceID] {
				continue
			}
			if decoded.TraceID == "" || decoded.SpanID == "" || decoded.Start.IsZero() || (!decoded.End.IsZero() && decoded.End.Before(decoded.Start)) || (entry.finished && decoded.End.IsZero()) {
				if !problem(fmt.Sprintf("line %d invalid span", line)) {
					out, problems, stopped = nil, nil, fmt.Errorf("native diagnostics exceed control record bound")
					return false
				}
				continue
			}
			retained += int64(len(text))
			if retained > bounds.ControlRecordLimit {
				out, stopped = nil, fmt.Errorf("selected spans exceed control record bound")
				return false
			}
			out = append(out, decoded)
		}
		return true
	})
	if stopped != nil {
		return out, problems, stopped
	}
	if err != nil {
		if line == 0 {
			return nil, nil, err
		}
		return out, problems, fmt.Errorf("read selected spans: %w", err)
	}
	finished := map[string]bool{}
	for _, s := range out {
		if !s.End.IsZero() {
			finished[s.TraceID+"/"+s.SpanID] = true
		}
	}
	selected := out[:0]
	for _, s := range out {
		if s.End.IsZero() && finished[s.TraceID+"/"+s.SpanID] {
			continue
		}
		selected = append(selected, s)
	}
	return selected, problems, nil
}
