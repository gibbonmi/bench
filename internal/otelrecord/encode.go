// Package otelrecord owns the Bench seam record: the address, the file layout, the
// encoder, and the processor. It is the only package that knows the record's shape.
//
// The encoder is hand-written. Every official OTLP-JSON encoder needs
// google.golang.org/protobuf, and the dependency standard keeps that footprint out, so
// these structs mirror the schema of upstream's own internal OTLP-JSON package.
package otelrecord

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// The four OTLP-JSON deviations from a plain Go encoding live in this file:
//
//  1. a trace id and a span id are lowercase hex strings, not base64 and not byte arrays;
//  2. the span kind and the status code are integers, not enum name strings;
//  3. startTimeUnixNano and endTimeUnixNano are quoted decimal strings, because a 64-bit
//     integer loses precision in a JSON number;
//  4. every field name is lowerCamelCase.

type tracesData struct {
	ResourceSpans []resourceSpans `json:"resourceSpans"`
}

type resourceSpans struct {
	Resource   resource     `json:"resource"`
	ScopeSpans []scopeSpans `json:"scopeSpans"`
	SchemaURL  string       `json:"schemaUrl,omitempty"`
}

type resource struct {
	Attributes []keyValue `json:"attributes,omitempty"`
}

type scopeSpans struct {
	Scope     scope  `json:"scope"`
	Spans     []span `json:"spans"`
	SchemaURL string `json:"schemaUrl,omitempty"`
}

type scope struct {
	Name       string     `json:"name,omitempty"`
	Version    string     `json:"version,omitempty"`
	Attributes []keyValue `json:"attributes,omitempty"`
}

type span struct {
	TraceID           string     `json:"traceId"`
	SpanID            string     `json:"spanId"`
	TraceState        string     `json:"traceState,omitempty"`
	ParentSpanID      string     `json:"parentSpanId,omitempty"`
	Name              string     `json:"name"`
	Kind              int        `json:"kind"`
	StartTimeUnixNano string     `json:"startTimeUnixNano"`
	EndTimeUnixNano   string     `json:"endTimeUnixNano,omitempty"`
	Attributes        []keyValue `json:"attributes,omitempty"`
	Status            status     `json:"status"`
}

type status struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

type keyValue struct {
	Key   string   `json:"key"`
	Value anyValue `json:"value"`
}

// anyValue mirrors the OTLP AnyValue union: exactly one field carries the value, and an
// int is quoted for the same 64-bit precision reason the timestamps are.
type anyValue struct {
	StringValue *string     `json:"stringValue,omitempty"`
	BoolValue   *bool       `json:"boolValue,omitempty"`
	IntValue    *string     `json:"intValue,omitempty"`
	DoubleValue *float64    `json:"doubleValue,omitempty"`
	ArrayValue  *arrayValue `json:"arrayValue,omitempty"`
}

type arrayValue struct {
	Values []anyValue `json:"values"`
}

// Encode returns one OTLP-JSON line for the span. The line carries no trailing newline;
// the writer that appends it owns the line terminator.
func Encode(readOnly sdktrace.ReadOnlySpan) ([]byte, error) {
	return json.Marshal(tracesData{ResourceSpans: []resourceSpans{encodeResourceSpans(readOnly)}})
}

func encodeResourceSpans(readOnly sdktrace.ReadOnlySpan) resourceSpans {
	out := resourceSpans{ScopeSpans: []scopeSpans{encodeScopeSpans(readOnly)}}
	if res := readOnly.Resource(); res != nil {
		out.Resource.Attributes = encodeAttributes(res.Attributes())
		out.SchemaURL = res.SchemaURL()
	}
	return out
}

func encodeScopeSpans(readOnly sdktrace.ReadOnlySpan) scopeSpans {
	instrumentation := readOnly.InstrumentationScope()
	return scopeSpans{
		Scope: scope{
			Name:       instrumentation.Name,
			Version:    instrumentation.Version,
			Attributes: encodeAttributes(instrumentation.Attributes.ToSlice()),
		},
		Spans:     []span{encodeSpan(readOnly)},
		SchemaURL: instrumentation.SchemaURL,
	}
}

func encodeSpan(readOnly sdktrace.ReadOnlySpan) span {
	spanContext := readOnly.SpanContext()
	traceID := spanContext.TraceID()
	spanID := spanContext.SpanID()
	out := span{
		TraceID:           hexID(traceID[:]),
		SpanID:            hexID(spanID[:]),
		TraceState:        spanContext.TraceState().String(),
		Name:              readOnly.Name(),
		Kind:              int(readOnly.SpanKind()),
		StartTimeUnixNano: unixNano(readOnly.StartTime().UnixNano()),
		Attributes:        encodeAttributes(readOnly.Attributes()),
		Status: status{
			Message: readOnly.Status().Description,
			Code:    statusCode(readOnly.Status().Code),
		},
	}
	if parent := readOnly.Parent(); parent.HasSpanID() {
		parentID := parent.SpanID()
		out.ParentSpanID = hexID(parentID[:])
	}
	// A start line has no end time, so the zero end time stays omitted rather than
	// reading back as an instantaneous span.
	if end := readOnly.EndTime(); !end.IsZero() {
		out.EndTimeUnixNano = unixNano(end.UnixNano())
	}
	return out
}

// hexID renders a trace id or a span id as the lowercase hex string OTLP-JSON requires.
// A stock proto3 JSON encoding writes these bytes as base64 instead.
func hexID(id []byte) string {
	return hex.EncodeToString(id)
}

func unixNano(nanoseconds int64) string {
	if nanoseconds < 0 {
		nanoseconds = 0
	}
	return strconv.FormatInt(nanoseconds, 10)
}

// statusCode maps the SDK's status code onto the OTLP wire number. The two orderings
// differ: the SDK numbers Error before Ok, and OTLP numbers Ok before Error.
func statusCode(code codes.Code) int {
	switch code {
	case codes.Ok:
		return 1
	case codes.Error:
		return 2
	default:
		return 0
	}
}

func encodeAttributes(pairs []attribute.KeyValue) []keyValue {
	if len(pairs) == 0 {
		return nil
	}
	out := make([]keyValue, 0, len(pairs))
	for _, pair := range pairs {
		out = append(out, keyValue{Key: string(pair.Key), Value: encodeValue(pair.Value)})
	}
	return out
}

func encodeValue(value attribute.Value) anyValue {
	switch value.Type() {
	case attribute.BOOL:
		return anyValue{BoolValue: boolPtr(value.AsBool())}
	case attribute.INT64:
		return anyValue{IntValue: stringPtr(strconv.FormatInt(value.AsInt64(), 10))}
	case attribute.FLOAT64:
		return anyValue{DoubleValue: floatPtr(value.AsFloat64())}
	case attribute.BOOLSLICE:
		return sliceValue(len(value.AsBoolSlice()), func(index int) anyValue {
			return anyValue{BoolValue: boolPtr(value.AsBoolSlice()[index])}
		})
	case attribute.INT64SLICE:
		return sliceValue(len(value.AsInt64Slice()), func(index int) anyValue {
			return anyValue{IntValue: stringPtr(strconv.FormatInt(value.AsInt64Slice()[index], 10))}
		})
	case attribute.FLOAT64SLICE:
		return sliceValue(len(value.AsFloat64Slice()), func(index int) anyValue {
			return anyValue{DoubleValue: floatPtr(value.AsFloat64Slice()[index])}
		})
	case attribute.STRINGSLICE:
		return sliceValue(len(value.AsStringSlice()), func(index int) anyValue {
			return anyValue{StringValue: stringPtr(value.AsStringSlice()[index])}
		})
	default:
		// STRING and INVALID both render as a string; Emit is the SDK's own rendering.
		return anyValue{StringValue: stringPtr(value.Emit())}
	}
}

func sliceValue(length int, at func(int) anyValue) anyValue {
	values := make([]anyValue, 0, length)
	for index := 0; index < length; index++ {
		values = append(values, at(index))
	}
	return anyValue{ArrayValue: &arrayValue{Values: values}}
}

func stringPtr(value string) *string  { return &value }
func boolPtr(value bool) *bool        { return &value }
func floatPtr(value float64) *float64 { return &value }
