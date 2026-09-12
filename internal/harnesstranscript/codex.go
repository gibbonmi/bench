package harnesstranscript

import (
	"encoding/json"
	"strings"
	"unicode/utf8"
)

// SourceCodexRollout is the first pinned source ID. It names the Codex rollout record: one
// JSON object per line, each with a top-level type, a timestamp, and a payload whose own
// type names the event. The mapping below derives from one reviewed transcript of that
// shape, so a later producer change lands as a new source ID rather than as a silent
// reinterpretation of this one.
const SourceCodexRollout = "codex-rollout-2026-09-11"

// The top-level record types this mapping reads.
const (
	typeResponseItem = "response_item"
	typeEventMsg     = "event_msg"
	typeCompacted    = "compacted"
)

// The payload types this mapping reads. A tool call and its result each arrive twice in the
// producer's vocabulary, once for the free-form shell tool and once for a typed function, and
// both carry the same call identity.
const (
	payloadCustomToolCall   = "custom_tool_call"
	payloadCustomToolOutput = "custom_tool_call_output"
	payloadFunctionCall     = "function_call"
	payloadFunctionOutput   = "function_call_output"
	payloadTaskStarted      = "task_started"
	payloadTokenCount       = "token_count"
)

// The shapes each measure reads. A source cell names the field path, so another reader can
// reach the same bytes. No cell holds a comma, because the rendered row separates its
// provenance pair on one.
const (
	// Both completion shapes hold result text, so the cell names both rather than sending a
	// reader to one field for bytes the other supplied.
	srcResultText = "response_item/custom_tool_call_output.output[].text and response_item/function_call_output.output"
	srcCallID     = "response_item/custom_tool_call.call_id"
	srcPairing    = "response_item call and output identities"
	srcTokenCount = "event_msg/token_count info.total_token_usage"
	srcCompacted  = "compacted record"
	// srcShellText names the one place this source could hold a read. The reader declines to
	// parse it, and the boundary beside the cell says so.
	srcShellText   = "event_msg/item_completed item CommandExecution.command"
	srcTaskStarted = "event_msg/task_started"
	srcNoShape     = "no pinned shape in this source"
)

// What each measure does not include. A boundary states the limit in the reader's own words,
// so a number is never quoted without the claim it can support.
const (
	boundaryText        = "Counts the UTF-8 bytes of tool-result text only. Serialized metadata, model messages, and reasoning payloads add nothing, and a call with no recorded result adds nothing."
	boundaryChars       = "Counts the code points of the same text the byte measure reads, so a multibyte result reports fewer characters than bytes."
	boundaryLines       = "Counts the lines of the same text. A result with no final newline still counts its last line."
	boundaryOuterCalls  = "Counts one call for each distinct call identity in this thread. A repeated record of one identity counts once."
	boundaryNestedCalls = "This source names a subagent by thread identity alone, and this reader decodes no subagent event. The subagent's own calls belong to a separate record, so none was observed here."
	boundaryUnmatched   = "Counts a call with no recorded result, and a result with no recorded call. A transcript shows what it recorded, so the byte measure beside it is never a claim about every result the producer wrote."
	boundaryTokens      = "The source reports a cumulative total at each snapshot. The value is the last snapshot in the interval, never a sum of the snapshots."
	boundaryCompactions = "Counts one identified compaction record for each compaction. A small context window is not evidence of one."
	boundaryTurns       = "Counts one task-start record for each turn."
	boundaryReadPaths   = "This source records a file read inside shell command text alone. Shell text is not a read census, so this reader parses no path out of it."
	boundaryAttribution = "The source reports tokens for the session and for the turn, never for one result. This reader makes no estimate from byte counts."
	// partialBoundary joins a boundary when an event did not parse or was too long to read.
	// An incomplete measure keeps its partial number, because a count that reads complete is
	// the worse error.
	partialBoundary = "One or more events were skipped, so this count is partial."
)

// codexRecord is one line of the rollout. Payload stays raw until the record type selects
// how to read it.
type codexRecord struct {
	Timestamp string          `json:"timestamp"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
}

type codexPayload struct {
	Type   string          `json:"type"`
	CallID string          `json:"call_id"`
	Output json.RawMessage `json:"output"`
	Info   *codexInfo      `json:"info"`
}

type codexInfo struct {
	Total *codexUsage `json:"total_token_usage"`
}

// codexUsage is one usage snapshot. The producer writes a running session total here, so a
// later snapshot supersedes an earlier one.
//
// Each counter is a pointer, because a snapshot that omits a key reported no value for that
// dimension. A plain integer would decode the absent key as zero and publish a total nobody
// measured.
type codexUsage struct {
	Input     *int64 `json:"input_tokens"`
	Cached    *int64 `json:"cached_input_tokens"`
	Output    *int64 `json:"output_tokens"`
	Reasoning *int64 `json:"reasoning_output_tokens"`
}

// counter reads one snapshot counter. The second result is false for a key the snapshot
// omits, which leaves that dimension unknown.
func counter(value *int64) (int64, bool) {
	if value == nil {
		return 0, false
	}
	return *value, true
}

type codexTextPart struct {
	Text string `json:"text"`
}

// codexTally accumulates one pass over the record. calls and outputs hold identities rather
// than counts, so a repeated record of one invocation cannot raise either census.
type codexTally struct {
	calls       map[string]bool
	outputs     map[string]bool
	bytes       int64
	chars       int64
	lines       int64
	compactions int64
	turns       int64
	usage       *codexUsage
	first       string
	last        string
	malformed   bool
}

// readCodex reads the record at path and reports its observations under this mapping.
func readCodex(path string) (Record, Failure) {
	tally := codexTally{calls: map[string]bool{}, outputs: map[string]bool{}}
	digest, failure := readRecord(path, &tally)
	if failure.State != "" {
		return Record{}, failure
	}
	return tally.record(digest), Failure{}
}

// skipped marks a line the reader could not take. Its relevance cannot be read, so every
// count it could have raised loses its claim to be complete.
func (t *codexTally) skipped() { t.malformed = true }

// line reads one record line into the tally. A line that does not parse marks the tally
// malformed: its relevance cannot be read, so every count it could have raised loses its
// claim to be complete.
func (t *codexTally) line(raw []byte) {
	if strings.TrimSpace(string(raw)) == "" {
		return
	}
	var record codexRecord
	if json.Unmarshal(raw, &record) != nil {
		t.malformed = true
		return
	}
	if record.Timestamp != "" {
		if t.first == "" {
			t.first = record.Timestamp
		}
		t.last = record.Timestamp
	}
	if record.Type == typeCompacted {
		t.compactions++
		return
	}
	var payload codexPayload
	if len(record.Payload) > 0 && json.Unmarshal(record.Payload, &payload) != nil {
		t.malformed = true
		return
	}
	switch {
	case record.Type == typeResponseItem && (payload.Type == payloadCustomToolCall || payload.Type == payloadFunctionCall):
		t.invocation(payload.CallID, t.calls)
	case record.Type == typeResponseItem && (payload.Type == payloadCustomToolOutput || payload.Type == payloadFunctionOutput):
		if payload.CallID != "" && t.outputs[payload.CallID] {
			// One invocation completed once. A repeated completion record is the same
			// result, so its bytes are already counted.
			return
		}
		t.invocation(payload.CallID, t.outputs)
		t.text(payload.Output)
	case record.Type == typeEventMsg && payload.Type == payloadTaskStarted:
		t.turns++
	case record.Type == typeEventMsg && payload.Type == payloadTokenCount:
		if payload.Info != nil && payload.Info.Total != nil {
			t.usage = payload.Info.Total
		}
	}
}

// invocation records one call identity. An event with no identity cannot join either census,
// so the affected counts report as partial rather than guessing at a pairing.
func (t *codexTally) invocation(id string, into map[string]bool) {
	if id == "" {
		t.malformed = true
		return
	}
	into[id] = true
}

// text counts one result's model-visible text. The producer writes a typed tool result as an
// array of text parts and a function result as one string, so both forms reach the same
// counters.
func (t *codexTally) text(raw json.RawMessage) {
	if len(raw) == 0 {
		return
	}
	var single string
	if json.Unmarshal(raw, &single) == nil {
		t.count(single)
		return
	}
	var parts []codexTextPart
	if json.Unmarshal(raw, &parts) != nil {
		t.malformed = true
		return
	}
	for _, part := range parts {
		t.count(part.Text)
	}
}

// count is the one site that measures result text.
func (t *codexTally) count(text string) {
	t.bytes += int64(len(text))
	t.chars += int64(utf8.RuneCountInString(text))
	t.lines += int64(textLines(text))
}

// textLines counts the lines one text occupies. A text with no final newline still ends a
// line, and an empty text occupies none.
func textLines(text string) int {
	if text == "" {
		return 0
	}
	lines := strings.Count(text, "\n")
	if !strings.HasSuffix(text, "\n") {
		lines++
	}
	return lines
}

// unmatched counts the invocations with no pair: a call whose result was never recorded, and
// a result whose call was never recorded.
func (t *codexTally) unmatched() int64 {
	var count int64
	for id := range t.calls {
		if !t.outputs[id] {
			count++
		}
	}
	for id := range t.outputs {
		if !t.calls[id] {
			count++
		}
	}
	return count
}

// interval names the span the observations cover. It is empty for a record with no timestamp
// to read.
func (t *codexTally) interval() string {
	if t.first == "" {
		return ""
	}
	return t.first + "/" + t.last
}

// codexMeasure is one dimension before it becomes an Observation. known is false for a
// dimension this source records no shape for, and for a counter no snapshot supplied.
type codexMeasure struct {
	count    int64
	known    bool
	source   string
	boundary string
}

// measures maps the tally onto the metric inventory. A dimension the source cannot answer
// stays unknown here rather than reporting the zero its counter happens to hold.
func (t *codexTally) measures() map[string]codexMeasure {
	token := func(pick func(codexUsage) *int64) codexMeasure {
		measure := codexMeasure{source: srcTokenCount, boundary: boundaryTokens}
		if t.usage != nil {
			measure.count, measure.known = counter(pick(*t.usage))
		}
		return measure
	}
	return map[string]codexMeasure{
		MetricResultBytes:      {count: t.bytes, known: true, source: srcResultText, boundary: boundaryText},
		MetricResultChars:      {count: t.chars, known: true, source: srcResultText, boundary: boundaryChars},
		MetricResultLines:      {count: t.lines, known: true, source: srcResultText, boundary: boundaryLines},
		MetricOuterCalls:       {count: int64(len(t.calls)), known: true, source: srcCallID, boundary: boundaryOuterCalls},
		MetricNestedCalls:      {source: srcNoShape, boundary: boundaryNestedCalls},
		MetricUnmatchedCalls:   {count: t.unmatched(), known: true, source: srcPairing, boundary: boundaryUnmatched},
		MetricInputTokens:      token(func(u codexUsage) *int64 { return u.Input }),
		MetricCachedTokens:     token(func(u codexUsage) *int64 { return u.Cached }),
		MetricOutputTokens:     token(func(u codexUsage) *int64 { return u.Output }),
		MetricReasoningTokens:  token(func(u codexUsage) *int64 { return u.Reasoning }),
		MetricCompactions:      {count: t.compactions, known: true, source: srcCompacted, boundary: boundaryCompactions},
		MetricTurns:            {count: t.turns, known: true, source: srcTaskStarted, boundary: boundaryTurns},
		MetricReadPaths:        {source: srcShellText, boundary: boundaryReadPaths},
		MetricTokenAttribution: {source: srcNoShape, boundary: boundaryAttribution},
	}
}

// record projects the tally onto the inventory in Metrics order.
func (t *codexTally) record(digest string) Record {
	measures := t.measures()
	rows := make([]Observation, 0, len(Metrics))
	for _, metric := range Metrics {
		measure := measures[metric.Name]
		row := Observation{
			Metric:       metric.Name,
			Unit:         metric.Unit,
			Availability: Unknown,
			Source:       measure.source,
			Boundary:     measure.boundary,
		}
		if measure.known {
			row.Count, row.Availability = measure.count, Observed
			if t.malformed {
				row.Availability = Incomplete
				row.Boundary += " " + partialBoundary
			}
		}
		rows = append(rows, row)
	}
	return Record{Digest: digest, SourceID: SourceCodexRollout, Interval: t.interval(), Observations: rows}
}
