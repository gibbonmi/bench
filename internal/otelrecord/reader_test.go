package otelrecord

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// beginRecorded opens one span through the real record path and returns its closer.
func beginRecorded(ctx context.Context, home, root, seam, name string) (context.Context, func()) {
	child, _, finish := BeginIn(ctx, home, root, seam, name)
	return child, finish
}

func TestReadSpansAnswersTheFinishedSpansOfTheRecord(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	landing, landingSpan, finishLanding := BeginIn(context.Background(), home, root, SeamLanding, SeamLanding)
	landingSpan.SetAttributes(attribute.String(AttrSubjectID, "abc123"))
	_, phase, finishPhase := BeginIn(landing, home, root, SeamGatePhase, "build")
	phase.SetAttributes(attribute.String(AttrOutcome, OutcomeGreen))
	time.Sleep(15 * time.Millisecond)
	finishPhase()
	finishLanding()

	spans, err := ReadSpans(home, root)
	if err != nil || len(spans) != 2 {
		t.Fatalf("ReadSpans = %d spans, %v, want the two finished spans", len(spans), err)
	}
	stage, landed := spans[0], spans[1]
	if stage.Name != "build" || stage.Seam != SeamGatePhase || stage.Attributes[AttrOutcome] != OutcomeGreen {
		t.Fatalf("stage span = %+v, want the named gate phase with its outcome", stage)
	}
	if landed.Seam != SeamLanding || landed.Attributes[AttrSubjectID] != "abc123" {
		t.Fatalf("landing span = %+v, want the landing seam with its subject", landed)
	}
	if stage.TraceID != landed.TraceID || stage.ParentSpanID != landed.SpanID {
		t.Fatalf("stage %s/%s does not sit under landing %s/%s", stage.TraceID, stage.ParentSpanID, landed.TraceID, landed.SpanID)
	}
	if stage.Elapsed() < 15*time.Millisecond {
		t.Fatalf("stage elapsed = %s, want at least the 15ms the phase ran", stage.Elapsed())
	}
}

func TestReadSpansRefusesAnAbsentRecord(t *testing.T) {
	spans, err := ReadSpans(t.TempDir(), t.TempDir())
	if err == nil || spans != nil {
		t.Fatalf("ReadSpans over an absent record = %v, %v, want a refusal", spans, err)
	}
}

func TestReadSpansSkipsAMalformedLine(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	if err := NewWriter(home, root).Append([]byte("not a record line")); err != nil {
		t.Fatal(err)
	}
	_, _, finish := BeginIn(context.Background(), home, root, SeamLanding, SeamLanding)
	finish()

	spans, err := ReadSpans(home, root)
	if err != nil || len(spans) != 1 || spans[0].Seam != SeamLanding {
		t.Fatalf("ReadSpans past a malformed line = %+v, %v, want the one finished landing", spans, err)
	}
}

func TestReadSpansRefusesAnUnreadableRecord(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	if err := NewWriter(home, root).Append([]byte("{}")); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(Path(home, root), 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(Path(home, root), 0o600) })
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root reads mode 0000 files; the unreadable record is unobservable")
	}
	if _, err := ReadSpans(home, root); err == nil {
		t.Fatal("ReadSpans over an unreadable record = nil error, want a refusal")
	}
}

func TestNewestLandingChoosesTheNewestCompletedLandingTrace(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	for _, subject := range []string{"older", "newer"} {
		landing, span, finish := BeginIn(context.Background(), home, root, SeamLanding, SeamLanding)
		span.SetAttributes(attribute.String(AttrSubjectID, subject))
		_, closePhase := beginRecorded(landing, home, root, SeamGatePhase, subject+"-build")
		closePhase()
		finish()
	}
	// A gate run of the operator's own, outside every landing trace, stays unlisted.
	foreign, closeForeign := beginRecorded(context.Background(), home, root, "gate", "gate.ordinary")
	_, closeForeignPhase := beginRecorded(foreign, home, root, SeamGatePhase, "foreign-build")
	closeForeignPhase()
	closeForeign()

	got, ok := NewestLanding(home, root)
	if !ok || got.Commit != "newer" {
		t.Fatalf("NewestLanding = %+v, %v, want the newer landing", got, ok)
	}
	if len(got.Stages) != 1 || got.Stages[0].Name != "newer-build" {
		t.Fatalf("stages = %+v, want only the newer landing's own phase", got.Stages)
	}
}

func TestNewestLandingSkipsTheNewestSubjectlessLanding(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	recordLanding := func(subject, stage string) {
		landing, span, finish := BeginIn(context.Background(), home, root, SeamLanding, SeamLanding)
		if subject != "" {
			span.SetAttributes(attribute.String(AttrSubjectID, subject))
		}
		_, closePhase := beginRecorded(landing, home, root, SeamGatePhase, stage)
		closePhase()
		finish()
	}
	recordLanding("published", "published-build")
	recordLanding("", "unpublished-build")

	got, ok := NewestLanding(home, root)
	if !ok || got.Commit != "published" {
		t.Fatalf("NewestLanding = %+v, %v, want the newest published landing", got, ok)
	}
	if len(got.Stages) != 1 || got.Stages[0].Name != "published-build" {
		t.Fatalf("stages = %+v, want the published landing's phase", got.Stages)
	}
}

func TestNewestLandingAnswersFalseWithoutALanding(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	if _, ok := NewestLanding(home, root); ok {
		t.Fatal("NewestLanding over an absent record = ok, want false")
	}
	_, finish := beginRecorded(context.Background(), home, root, "gate", "gate.ordinary")
	finish()
	if got, ok := NewestLanding(home, root); ok {
		t.Fatalf("NewestLanding over a record with no landing = %+v, ok, want false", got)
	}
}

func TestNewestLandingAnswersFalseWithOnlyASubjectlessLanding(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	_, finish := beginRecorded(context.Background(), home, root, SeamLanding, SeamLanding)
	finish()

	if got, ok := NewestLanding(home, root); ok {
		t.Fatalf("NewestLanding over a subjectless landing = %+v, ok, want false", got)
	}
}

// The reader selects by two seam names the packages that open those spans declare for
// themselves. This test is the reconciliation that keeps the pair from drifting.
func TestReaderSeamsAreRegisteredSeams(t *testing.T) {
	registered := map[string]bool{}
	for _, entry := range Registry {
		registered[entry.Seam] = true
	}
	for _, seam := range []string{SeamLanding, SeamGatePhase} {
		if !registered[seam] {
			t.Fatalf("reader seam %q names no Registry row", seam)
		}
	}
}

// schemaLine is the finished fixture span as one record line whose resource block holds
// only the given schema value. An empty schema writes a legacy line with no schema key.
func schemaLine(t *testing.T, schema string) []byte {
	t.Helper()

	line, err := Encode(fixtureSpan(t))
	if err != nil {
		t.Fatalf("the encoder failed: %v", err)
	}
	var data tracesData
	if err := json.Unmarshal(line, &data); err != nil {
		t.Fatalf("the line does not parse: %v", err)
	}
	data.ResourceSpans[0].Resource.Attributes = nil
	if schema != "" {
		data.ResourceSpans[0].Resource.Attributes = []keyValue{{Key: ResourceRecordSchema, Value: anyValue{StringValue: stringPtr(schema)}}}
	}
	out, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("encode the schema line: %v", err)
	}
	return out
}

// recordOf writes the lines as one repository's record and returns its home and root.
func recordOf(t *testing.T, lines ...[]byte) (string, string) {
	t.Helper()

	home, root := t.TempDir(), t.TempDir()
	for _, line := range lines {
		if err := NewWriter(home, root).Append(line); err != nil {
			t.Fatalf("append a record line: %v", err)
		}
	}
	return home, root
}

// TestReadSelectedReportsAnUnknownSchemaAsMalformed holds row LE8: a reader that ignores
// the resource returns the span, so the problem read reds.
func TestReadSelectedReportsAnUnknownSchemaAsMalformed(t *testing.T) {
	home, root := recordOf(t, schemaLine(t, "2"))

	spans, problems, err := ReadSelected(home, root, []string{fixtureTraceID})
	if err != nil {
		t.Fatalf("ReadSelected: %v", err)
	}
	if len(spans) != 0 {
		t.Errorf("ReadSelected returned %d spans from a schema 2 line, want none", len(spans))
	}
	if !reflect.DeepEqual(problems, []string{"line 1 malformed"}) {
		t.Errorf("problems = %v, want the one malformed line", problems)
	}
}

// TestReadSpansReturnsNoSpanOfAnUnknownSchema holds row LE9: a reader that ignores the
// resource returns the span, so the count reds.
func TestReadSpansReturnsNoSpanOfAnUnknownSchema(t *testing.T) {
	home, root := recordOf(t, schemaLine(t, "2"))

	spans, err := ReadSpans(home, root)
	if err != nil || len(spans) != 0 {
		t.Fatalf("ReadSpans = %d spans, %v, want none from a schema 2 line", len(spans), err)
	}
}

// TestReadSpansReadsALegacyLine holds row LE10: a reader that requires the schema
// attribute drops the legacy span, so the count reds.
func TestReadSpansReadsALegacyLine(t *testing.T) {
	home, root := recordOf(t, schemaLine(t, ""))

	spans, err := ReadSpans(home, root)
	if err != nil || len(spans) != 1 || spans[0].TraceID != fixtureTraceID {
		t.Fatalf("ReadSpans = %+v, %v, want the one legacy span", spans, err)
	}
}

// spanLine is one finished record line for a span named name in the trace traceHex,
// with its seam and extra attributes.
func spanLine(t *testing.T, traceHex, name, seam string, extra ...attribute.KeyValue) []byte {
	t.Helper()

	traceID, err := trace.TraceIDFromHex(traceHex)
	if err != nil {
		t.Fatal(err)
	}
	stub := tracetest.SpanStubFromReadOnlySpan(fixtureSpan(t))
	stub.Name = name
	stub.SpanContext = stub.SpanContext.WithTraceID(traceID)
	stub.Attributes = append([]attribute.KeyValue{attribute.String(AttrSeam, seam)}, extra...)
	line, err := Encode(stub.Snapshot())
	if err != nil {
		t.Fatalf("the encoder failed: %v", err)
	}
	return line
}

// sealEach writes each line as its own segment: every line but the last is sealed, and
// the last one stays in the live segment.
func sealEach(t *testing.T, lines ...[]byte) (string, string) {
	t.Helper()

	home, root := t.TempDir(), t.TempDir()
	writer := newWriter(home, root, int64(len(lines[0])), 8)
	for _, line := range lines {
		if err := writer.Append(line); err != nil {
			t.Fatalf("append a record line: %v", err)
		}
	}
	if got := len(sealedNames(t, home, root)); got != len(lines)-1 {
		t.Fatalf("the writer sealed %d segments, want %d", got, len(lines)-1)
	}
	return home, root
}

const otherTraceID = "f0e0d0c0b0a090807060504030201000"

// TestReadSpansReadsTheSealedSegmentsInOrder holds row LE15: a reader of the live
// segment alone drops the sealed spans, so the ordered comparison reds.
func TestReadSpansReadsTheSealedSegmentsInOrder(t *testing.T) {
	home, root := sealEach(t,
		spanLine(t, fixtureTraceID, "first", "gate"),
		spanLine(t, fixtureTraceID, "second", "gate"),
		spanLine(t, fixtureTraceID, "third", "gate"))

	spans, err := ReadSpans(home, root)
	if err != nil {
		t.Fatalf("ReadSpans: %v", err)
	}
	var names []string
	for _, span := range spans {
		names = append(names, span.Name)
	}
	if !reflect.DeepEqual(names, []string{"first", "second", "third"}) {
		t.Fatalf("ReadSpans names = %v, want the two sealed spans in order and then the live one", names)
	}
}

// TestNewestLandingReadsStagesFromASealedSegment holds row LE16: a live-only read finds
// the landing without its stages, so the stage count reds.
func TestNewestLandingReadsStagesFromASealedSegment(t *testing.T) {
	home, root := sealEach(t,
		spanLine(t, fixtureTraceID, "build", SeamGatePhase),
		spanLine(t, fixtureTraceID, SeamLanding, SeamLanding, attribute.String(AttrSubjectID, "abc123")))

	landing, ok := NewestLanding(home, root)
	if !ok || landing.Commit != "abc123" || len(landing.Stages) != 1 || landing.Stages[0].Name != "build" {
		t.Fatalf("NewestLanding = %+v, %v, want the landing with its sealed stage", landing, ok)
	}
}

// TestReadSelectedReadsASealedSegment holds row LE17: a live-only read misses the span,
// so the selection reds.
func TestReadSelectedReadsASealedSegment(t *testing.T) {
	home, root := sealEach(t,
		spanLine(t, fixtureTraceID, "sealed", "gate"),
		spanLine(t, otherTraceID, "live", "gate"))

	spans, problems, err := ReadSelected(home, root, []string{fixtureTraceID})
	if err != nil || len(problems) != 0 || len(spans) != 1 || spans[0].Name != "sealed" {
		t.Fatalf("ReadSelected = %+v, %v, %v, want the one sealed span", spans, problems, err)
	}
}

// liveRecord writes one live line below a fresh home and returns the home and root.
func liveRecord(t *testing.T) (string, string) {
	t.Helper()
	return recordOf(t, spanLine(t, fixtureTraceID, "live", "gate"))
}

// TestReadSpansRefusesAFIFOSegment holds row LE20: an open of the FIFO blocks the read,
// so the deadline reds.
func TestReadSpansRefusesAFIFOSegment(t *testing.T) {
	home, root := liveRecord(t)
	if err := syscall.Mkfifo(filepath.Join(Dir(home, root), sealedName(1)), 0o600); err != nil {
		t.Fatalf("create fifo: %v", err)
	}

	answered := make(chan error, 1)
	go func() {
		_, err := ReadSpans(home, root)
		answered <- err
	}()
	window := bounds.TestDeadline(0)
	select {
	case err := <-answered:
		if err == nil {
			t.Fatal("ReadSpans accepted a FIFO sealed segment")
		}
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict("ReadSpans to answer at a FIFO sealed segment", window))
	}
}

// TestReadSpansRefusesASymlinkedSegment holds row LE21: a reader that follows the link
// reads bytes outside the record, so the error read reds.
func TestReadSpansRefusesASymlinkedSegment(t *testing.T) {
	home, root := liveRecord(t)
	outside := filepath.Join(t.TempDir(), "outside.jsonl")
	if err := os.WriteFile(outside, append(spanLine(t, fixtureTraceID, "outside", "gate"), '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(Dir(home, root), sealedName(1))); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	if spans, err := ReadSpans(home, root); err == nil {
		t.Fatalf("ReadSpans followed a symlinked sealed segment and returned %d spans", len(spans))
	}
}
