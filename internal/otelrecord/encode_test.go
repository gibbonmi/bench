package otelrecord

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// The tests parse the emitted line back with encoding/json and assert on the parsed
// form. A golden string would grade the encoder's whitespace instead of its wire shape.

const (
	fixtureTraceID   = "0102030405060708090a0b0c0d0e0f10"
	fixtureSpanID    = "1112131415161718"
	fixtureStartNano = int64(1756512000000000123)
)

func fixtureSpan(t *testing.T) sdktrace.ReadOnlySpan {
	t.Helper()

	traceID, err := trace.TraceIDFromHex(fixtureTraceID)
	if err != nil {
		t.Fatalf("the fixture trace id does not parse: %v", err)
	}
	spanID, err := trace.SpanIDFromHex(fixtureSpanID)
	if err != nil {
		t.Fatalf("the fixture span id does not parse: %v", err)
	}

	start := time.Unix(0, fixtureStartNano).UTC()
	stub := tracetest.SpanStub{
		Name: "bench.gate",
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		}),
		SpanKind:  trace.SpanKindInternal,
		StartTime: start,
		EndTime:   start.Add(2 * time.Second),
		Attributes: []attribute.KeyValue{
			attribute.String(AttrSeam, "gate"),
			attribute.String(AttrSubjectID, "0f1e2d3c"),
			attribute.String(AttrOutcome, "green"),
			attribute.Int(AttrMeasurePathCount, 3),
		},
		Status:               sdktrace.Status{Code: codes.Error, Description: "the lane is red"},
		InstrumentationScope: instrumentation.Scope{Name: "bench", Version: "test"},
	}
	return stub.Snapshot()
}

func encodeFixture(t *testing.T) map[string]any {
	t.Helper()

	line, err := Encode(fixtureSpan(t))
	if err != nil {
		t.Fatalf("the encoder failed: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(line, &parsed); err != nil {
		t.Fatalf("the line does not parse as a JSON object: %v: %s", err, line)
	}
	return parsed
}

// onlySpan walks the parsed line down to its one span object. It fails the test if any
// level is missing, so every row below reads a span that really came out of the encoder.
func onlySpan(t *testing.T, parsed map[string]any) map[string]any {
	t.Helper()

	resourceSpans, ok := parsed["resourceSpans"].([]any)
	if !ok || len(resourceSpans) != 1 {
		t.Fatalf("the line carries no single resourceSpans entry: %v", parsed)
	}
	scopeSpans, ok := resourceSpans[0].(map[string]any)["scopeSpans"].([]any)
	if !ok || len(scopeSpans) != 1 {
		t.Fatalf("the resourceSpans entry carries no single scopeSpans entry: %v", resourceSpans[0])
	}
	spans, ok := scopeSpans[0].(map[string]any)["spans"].([]any)
	if !ok || len(spans) != 1 {
		t.Fatalf("the scopeSpans entry carries no single span: %v", scopeSpans[0])
	}
	span, ok := spans[0].(map[string]any)
	if !ok {
		t.Fatalf("the span is not a JSON object: %v", spans[0])
	}
	return span
}

// raw returns the span field as the JSON text the encoder wrote, so a type assertion can
// tell a quoted decimal string from a bare JSON number.
func raw(t *testing.T, span map[string]any, field string) any {
	t.Helper()

	value, present := span[field]
	if !present {
		t.Fatalf("the span carries no %s field: %v", field, span)
	}
	return value
}

// OT1.
func TestEncodeReturnsOneResourceSpansLine(t *testing.T) {
	line, err := Encode(fixtureSpan(t))
	if err != nil {
		t.Fatalf("the encoder failed: %v", err)
	}
	if idx := indexOfNewline(line); idx >= 0 {
		t.Fatalf("the encoder returned more than one line: a newline sits at byte %d", idx)
	}
	var parsed map[string]any
	if err := json.Unmarshal(line, &parsed); err != nil {
		t.Fatalf("the line does not parse as a JSON object: %v: %s", err, line)
	}
	if _, present := parsed["resourceSpans"]; !present {
		t.Fatalf("the line carries no resourceSpans key: %v", parsed)
	}
	onlySpan(t, parsed)
}

func indexOfNewline(line []byte) int {
	for index, b := range line {
		if b == '\n' {
			return index
		}
	}
	return -1
}

// OT2.
func TestEncodeWritesLowercaseHexIDs(t *testing.T) {
	span := onlySpan(t, encodeFixture(t))
	for field, want := range map[string]string{
		"traceId": fixtureTraceID,
		"spanId":  fixtureSpanID,
	} {
		got, ok := raw(t, span, field).(string)
		if !ok {
			t.Errorf("%s is not a JSON string: %#v", field, span[field])
			continue
		}
		if got != want {
			t.Errorf("%s is %q, want the lowercase hex %q", field, got, want)
		}
		if !regexp.MustCompile(`^[0-9a-f]+$`).MatchString(got) {
			t.Errorf("%s is not lowercase hex: %q", field, got)
		}
	}
}

// OT3.
func TestEncodeWritesKindAndStatusCodeAsIntegers(t *testing.T) {
	span := onlySpan(t, encodeFixture(t))

	kind, ok := raw(t, span, "kind").(float64)
	if !ok {
		t.Fatalf("kind is not a JSON number: %#v", span["kind"])
	}
	if int(kind) != int(trace.SpanKindInternal) {
		t.Errorf("kind is %v, want %d", kind, int(trace.SpanKindInternal))
	}

	statusObject, ok := raw(t, span, "status").(map[string]any)
	if !ok {
		t.Fatalf("status is not a JSON object: %#v", span["status"])
	}
	code, ok := statusObject["code"].(float64)
	if !ok {
		t.Fatalf("the status code is not a JSON number: %#v", statusObject["code"])
	}
	// The fixture status is Error, which is 2 on the OTLP wire and 1 in the SDK.
	if int(code) != 2 {
		t.Errorf("the status code is %v, want the OTLP wire number 2", code)
	}
}

// OT4.
func TestEncodeWritesStartTimeAsQuotedDecimalString(t *testing.T) {
	span := onlySpan(t, encodeFixture(t))

	start, ok := raw(t, span, "startTimeUnixNano").(string)
	if !ok {
		t.Fatalf("startTimeUnixNano is not a JSON string: %#v", span["startTimeUnixNano"])
	}
	if start != "1756512000000000123" {
		t.Errorf("startTimeUnixNano is %q, want the fixture's decimal nanoseconds", start)
	}

	end, ok := raw(t, span, "endTimeUnixNano").(string)
	if !ok {
		t.Fatalf("endTimeUnixNano is not a JSON string: %#v", span["endTimeUnixNano"])
	}
	if end != "1756512002000000123" {
		t.Errorf("endTimeUnixNano is %q, want the fixture's decimal nanoseconds", end)
	}
}

// A start line carries no end time, so a consumer can filter the unfinished spans.
func TestEncodeOmitsAnAbsentEndTime(t *testing.T) {
	stub := tracetest.SpanStubFromReadOnlySpan(fixtureSpan(t))
	stub.EndTime = time.Time{}

	line, err := Encode(stub.Snapshot())
	if err != nil {
		t.Fatalf("the encoder failed: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(line, &parsed); err != nil {
		t.Fatalf("the line does not parse: %v", err)
	}
	if _, present := onlySpan(t, parsed)["endTimeUnixNano"]; present {
		t.Error("a span with no end time still wrote endTimeUnixNano")
	}
}

// The declared attribute set reaches the line as OTLP key-value pairs.
func TestEncodeWritesDeclaredAttributes(t *testing.T) {
	span := onlySpan(t, encodeFixture(t))

	attributes, ok := raw(t, span, "attributes").([]any)
	if !ok {
		t.Fatalf("attributes is not a JSON array: %#v", span["attributes"])
	}
	got := map[string]any{}
	for _, entry := range attributes {
		pair, ok := entry.(map[string]any)
		if !ok {
			t.Fatalf("an attribute is not a JSON object: %#v", entry)
		}
		key, ok := pair["key"].(string)
		if !ok {
			t.Fatalf("an attribute carries no string key: %#v", pair)
		}
		got[key] = pair["value"]
	}

	seam, ok := got[AttrSeam].(map[string]any)
	if !ok || seam["stringValue"] != "gate" {
		t.Errorf("%s is %#v, want the string value gate", AttrSeam, got[AttrSeam])
	}
	// An int attribute is quoted for the same 64-bit precision reason the times are.
	count, ok := got[AttrMeasurePathCount].(map[string]any)
	if !ok || count["intValue"] != "3" {
		t.Errorf("%s is %#v, want the quoted int value 3", AttrMeasurePathCount, got[AttrMeasurePathCount])
	}
}
