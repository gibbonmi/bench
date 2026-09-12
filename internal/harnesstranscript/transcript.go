// Package harnesstranscript reads one explicit harness session record and reports typed
// observations of it. It owns the metric inventory, so the rendered row count and the row
// names have one source, and a consumer projects what this package returns rather than
// deciding what a session record holds.
//
// A source ID pins one producer's event shapes, text boundaries, invocation identities, and
// counter semantics. A record whose source ID has no pinned mapping is never guessed at:
// every observation reads unknown, because an invented parser reports numbers nobody read.
//
// Availability separates the three answers a measure can give. Observed means the reader
// counted the fact. Incomplete means a malformed event dropped out of a count that would
// otherwise read as complete. Unknown means the source records no shape for the dimension,
// and an unknown measure renders no number, so a missing fact can never look like a zero.
//
// The package imports the standard library only, so the harness record's consumers compose
// it with no new import edge of their own.
package harnesstranscript

// Availability is the closed set of answers one measure can give about its own evidence.
type Availability string

const (
	Observed   Availability = "observed"
	Incomplete Availability = "incomplete"
	Unknown    Availability = "unknown"
)

// The record states a failed read reports. The spellings are the shared control-record
// vocabulary, and the consumer renders them through the one AXI record-error line. This
// package stays standard-library-only, so it names the spellings rather than importing the
// classifier that owns them.
const (
	StateAbsent     = "absent"
	StateUnreadable = "unreadable"
	StateWrongType  = "wrong-type"
)

// The metric names. Each name is one row of every rendered record view.
const (
	MetricResultBytes      = "result-text bytes"
	MetricResultChars      = "result-text characters"
	MetricResultLines      = "result-text lines"
	MetricOuterCalls       = "outer calls"
	MetricNestedCalls      = "observed nested calls"
	MetricUnmatchedCalls   = "unmatched calls"
	MetricInputTokens      = "input tokens"
	MetricCachedTokens     = "cached-input tokens"
	MetricOutputTokens     = "output tokens"
	MetricReasoningTokens  = "reasoning tokens"
	MetricCompactions      = "compactions"
	MetricTurns            = "turns"
	MetricReadPaths        = "explicit read paths"
	MetricTokenAttribution = "per-result token attribution"
)

// The measure units. A unit is a property of the metric, not of one source, so it is spelled
// here beside the name.
const (
	unitBytes       = "bytes"
	unitCharacters  = "characters"
	unitLines       = "lines"
	unitCalls       = "calls"
	unitTokens      = "tokens"
	unitCompactions = "compactions"
	unitTurns       = "turns"
	unitPaths       = "paths"
)

// Metric is one dimension the reader reports for every source.
type Metric struct {
	Name string
	Unit string
}

// Metrics is the inventory in the order every view prints. A source that cannot supply a
// row still prints the row as unknown, so two records of different sources stay comparable
// row by row.
var Metrics = []Metric{
	{MetricResultBytes, unitBytes},
	{MetricResultChars, unitCharacters},
	{MetricResultLines, unitLines},
	{MetricOuterCalls, unitCalls},
	{MetricNestedCalls, unitCalls},
	{MetricUnmatchedCalls, unitCalls},
	{MetricInputTokens, unitTokens},
	{MetricCachedTokens, unitTokens},
	{MetricOutputTokens, unitTokens},
	{MetricReasoningTokens, unitTokens},
	{MetricCompactions, unitCompactions},
	{MetricTurns, unitTurns},
	{MetricReadPaths, unitPaths},
	{MetricTokenAttribution, unitTokens},
}

// Observation is one measure. Source names the record shape the count came from, and
// Boundary states what the count does not include, so a reader can reproduce the
// observation without opening the source.
type Observation struct {
	Metric       string
	Unit         string
	Count        int64
	Availability Availability
	Source       string
	Boundary     string
}

// Cell is the measure's rendered value. Only an observed measure renders a zero, so an empty
// cell can never be misread as a counted nothing. An incomplete measure keeps its partial
// number, because the availability cell beside it already denies the complete claim, but an
// incomplete zero renders no number: a partial count of nothing says nothing at all.
func (o Observation) Cell() any {
	if o.Availability == Unknown || (o.Availability == Incomplete && o.Count == 0) {
		return ""
	}
	return o.Count
}

// Record is one inspected session record: its identity and its observations in Metrics
// order. Digest and Interval are empty when no bytes were read, which is the unsupported
// source's state.
type Record struct {
	Digest       string
	SourceID     string
	Interval     string
	Observations []Observation
}

// Failure is a read that yielded nothing a consumer may trust. State is empty exactly when
// the read succeeded.
type Failure struct {
	State  string
	Reason string
}

// Supported reports whether sourceID has a pinned mapping.
func Supported(sourceID string) bool {
	return sourceID == SourceCodexRollout
}

// Unsupported is the record an unrecognized source ID yields: the identity a caller already
// knows, and every measure unknown. The caller reports the refusal; this package never
// guesses a shape from an unpinned name.
func Unsupported(sourceID string) Record {
	return Record{SourceID: sourceID, Observations: unknownObservations(unsupportedSource, unsupportedBoundary)}
}

const (
	unsupportedSource   = "no pinned mapping"
	unsupportedBoundary = "the source ID names no reviewed event mapping, so no dimension was read"
)

// unknownObservations builds the whole inventory with nothing counted. Every reader starts
// here and overwrites the rows its source can answer, so a dimension a mapping forgets stays
// unknown instead of reading zero.
func unknownObservations(source, boundary string) []Observation {
	rows := make([]Observation, 0, len(Metrics))
	for _, metric := range Metrics {
		rows = append(rows, Observation{
			Metric:       metric.Name,
			Unit:         metric.Unit,
			Availability: Unknown,
			Source:       source,
			Boundary:     boundary,
		})
	}
	return rows
}

// Read returns the observations of the record at path under sourceID. The second result
// carries the refusal when the path is not a readable regular file; Record is meaningless
// then. A caller checks Supported first, because an unpinned source ID never reaches a read.
func Read(path, sourceID string) (Record, Failure) {
	if !Supported(sourceID) {
		return Unsupported(sourceID), Failure{}
	}
	return readCodex(path)
}
