package harnesstranscript

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
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

// maxRecordLine bounds one line of the record. A rollout line holds a whole tool result, so
// the bound sits far above an ordinary line and only a line no reader could use reaches it.
// A line past the bound is a malformed event rather than a silent truncation.
const maxRecordLine = 16 << 20

// The shapes each measure reads. A source cell names the field path, so another reader can
// reach the same bytes. No cell holds a comma, because the rendered row separates its
// provenance pair on one.
const (
	srcResultText  = "response_item/custom_tool_call_output.output[].text"
	srcCallID      = "response_item/custom_tool_call.call_id"
	srcSubAgent    = "event_msg/item_completed item SubAgentActivity"
	srcPairing     = "response_item call and output identities"
	srcTokenCount  = "event_msg/token_count info.total_token_usage"
	srcCompacted   = "compacted record"
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
	boundaryNestedCalls = "This source names a subagent by thread identity alone. The subagent's own calls belong to a separate record, so this reader observed none of them."
	boundaryUnmatched   = "Counts a call with no recorded result, and a result with no recorded call. A transcript shows what it recorded, so the byte measure beside it is never a claim about every result the producer wrote."
	boundaryTokens      = "The source reports a cumulative total at each snapshot. The value is the last snapshot in the interval, never a sum of the snapshots."
	boundaryCompactions = "Counts one identified compaction record for each compaction. A small context window is not evidence of one."
	boundaryTurns       = "Counts one task-start record for each turn."
	boundaryReadPaths   = "This source records a file read inside shell command text alone. Shell text is not a read census, so this reader parses no path out of it."
	boundaryAttribution = "The source reports tokens for the session and for the turn, never for one result. This reader makes no estimate from byte counts."
	// partialBoundary joins a boundary when an event did not parse. An incomplete measure
	// keeps its partial number, because a count that reads complete is the worse error.
	partialBoundary = "One or more events did not parse, so this count is partial."
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
type codexUsage struct {
	Input     int64 `json:"input_tokens"`
	Cached    int64 `json:"cached_input_tokens"`
	Output    int64 `json:"output_tokens"`
	Reasoning int64 `json:"reasoning_output_tokens"`
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

// readCodex reads the record at path once and reports its observations.
//
// The path is typed before it is opened. A FIFO would block the reader inside open(2), and a
// link would supply bytes from somewhere other than the path the agent named, so both are
// refused on the stat that precedes every open. The record is then streamed a line at a time:
// a session record outgrows any whole-file bound, and the digest is taken from the same pass,
// so the identity names exactly the bytes the counts came from.
//
// Record contents are data throughout. The reader decodes JSON, counts, and renders numbers.
// It never executes a field, resolves a path out of one, or copies record text into output.
func readCodex(path string) (Record, Failure) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Record{}, Failure{State: StateAbsent, Reason: "no record at the named path"}
		}
		return Record{}, Failure{State: StateUnreadable, Reason: err.Error()}
	}
	if !info.Mode().IsRegular() {
		return Record{}, Failure{State: StateWrongType, Reason: "not a regular file: " + info.Mode().Type().String()}
	}
	file, err := os.Open(path)
	if err != nil {
		return Record{}, Failure{State: StateUnreadable, Reason: err.Error()}
	}
	defer file.Close()

	digest := sha256.New()
	tally := codexTally{calls: map[string]bool{}, outputs: map[string]bool{}}
	scanner := bufio.NewScanner(io.TeeReader(file, digest))
	scanner.Buffer(make([]byte, 0, 64<<10), maxRecordLine)
	for scanner.Scan() {
		tally.line(scanner.Bytes())
	}
	if scanner.Err() != nil {
		// The scan stopped early, so the counts are partial and the digest still has to
		// name the whole file the agent pointed at.
		tally.malformed = true
		if _, err := io.Copy(digest, file); err != nil {
			return Record{}, Failure{State: StateUnreadable, Reason: err.Error()}
		}
	}
	return tally.record("sha256:" + hex.EncodeToString(digest.Sum(nil))), Failure{}
}

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
	token := func(pick func(codexUsage) int64) codexMeasure {
		measure := codexMeasure{source: srcTokenCount, boundary: boundaryTokens}
		if t.usage != nil {
			measure.known, measure.count = true, pick(*t.usage)
		}
		return measure
	}
	return map[string]codexMeasure{
		MetricResultBytes:      {count: t.bytes, known: true, source: srcResultText, boundary: boundaryText},
		MetricResultChars:      {count: t.chars, known: true, source: srcResultText, boundary: boundaryChars},
		MetricResultLines:      {count: t.lines, known: true, source: srcResultText, boundary: boundaryLines},
		MetricOuterCalls:       {count: int64(len(t.calls)), known: true, source: srcCallID, boundary: boundaryOuterCalls},
		MetricNestedCalls:      {source: srcSubAgent, boundary: boundaryNestedCalls},
		MetricUnmatchedCalls:   {count: t.unmatched(), known: true, source: srcPairing, boundary: boundaryUnmatched},
		MetricInputTokens:      token(func(u codexUsage) int64 { return u.Input }),
		MetricCachedTokens:     token(func(u codexUsage) int64 { return u.Cached }),
		MetricOutputTokens:     token(func(u codexUsage) int64 { return u.Output }),
		MetricReasoningTokens:  token(func(u codexUsage) int64 { return u.Reasoning }),
		MetricCompactions:      {count: t.compactions, known: true, source: srcCompacted, boundary: boundaryCompactions},
		MetricTurns:            {count: t.turns, known: true, source: srcTaskStarted, boundary: boundaryTurns},
		MetricReadPaths:        {source: srcNoShape, boundary: boundaryReadPaths},
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
