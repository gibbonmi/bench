package otelrecord

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
// resolved home, in record order. An absent record and an unreadable record are both
// refusals, so a consumer states that it knows nothing rather than that nothing ran.
// A single malformed line is skipped: one truncated write must not hide the run that
// wrote every other line.
func ReadSpans(home, root string) ([]Span, error) {
	file, err := os.Open(Path(home, root))
	if err != nil {
		return nil, fmt.Errorf("open seam record: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64<<10), recordLineLimit)
	var spans []Span
	for scanner.Scan() {
		spans = append(spans, finishedSpans(scanner.Bytes())...)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read seam record: %w", err)
	}
	return spans, nil
}

// Landing is the newest completed landing of one repository's record and the gate
// stages that ran below it.
type Landing struct {
	// Commit is the published subject the landing recorded, empty when the landing
	// published none.
	Commit string
	// TraceID is the landing's trace, which is how a reader verifies the selection.
	TraceID string
	// Stages are the phase spans of that trace, in record order.
	Stages []Span
}

// NewestLanding returns the newest completed landing of root's record. It answers false
// when the record is absent, is unreadable, or names no completed landing, which is the
// same answer a consumer states as unknown.
func NewestLanding(home, root string) (Landing, bool) {
	spans, err := ReadSpans(home, root)
	if err != nil {
		return Landing{}, false
	}
	newest := -1
	for index, candidate := range spans {
		if candidate.Seam != SeamLanding {
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
	var data tracesData
	if json.Unmarshal(line, &data) != nil {
		return nil
	}
	var out []Span
	for _, resource := range data.ResourceSpans {
		for _, scope := range resource.ScopeSpans {
			for _, encoded := range scope.Spans {
				if encoded.EndTimeUnixNano == "" {
					continue
				}
				out = append(out, decodeSpan(encoded))
			}
		}
	}
	return out
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
