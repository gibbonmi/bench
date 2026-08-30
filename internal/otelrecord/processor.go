package otelrecord

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// processor writes one record line when a span starts and a second line when it ends.
// No published exporter writes at the start, so the start line is this package's own
// work on the SDK's sanctioned OnStart hook. An interrupted run therefore keeps the
// evidence of the phase it began.
//
// The two lines are independent records, and no consumer merges them. A reader filters
// the unfinished spans by the start marker, and it derives elapsed time from the end
// line alone, which carries the complete span.
type processor struct {
	writer *Writer
}

// newProcessor returns the span processor that appends to root's record below an
// explicitly resolved home.
func newProcessor(home, root string) *processor {
	return &processor{writer: NewWriter(home, root)}
}

// OnStart appends the start line. The span is not ended here, so its end time is zero
// and the encoder omits it.
func (p *processor) OnStart(_ context.Context, readWrite sdktrace.ReadWriteSpan) {
	p.append(markStart(readWrite))
}

// OnEnd appends the end line: the complete span, with its end time.
func (p *processor) OnEnd(readOnly sdktrace.ReadOnlySpan) {
	p.append(readOnly)
}

// Shutdown and ForceFlush have nothing to do. Each append is one synchronous write, so
// the processor holds no buffer and starts no background worker.
func (p *processor) Shutdown(context.Context) error   { return nil }
func (p *processor) ForceFlush(context.Context) error { return nil }

// append writes one encoded line and drops any failure. The SpanProcessor hooks return
// nothing, and the record must never change a verb's outcome, so the error stops here.
// Writer.Append still returns every failure for a caller that appends directly.
func (p *processor) append(readOnly sdktrace.ReadOnlySpan) {
	line, err := Encode(readOnly)
	if err != nil {
		return
	}
	_ = p.writer.Append(line)
}

// markStart returns the span to encode for the start line: the live span's state plus
// the bench.record=start marker. The marker rides on the copy, so the live span never
// carries a record-shape attribute that its own end line would then repeat.
//
// The copy goes through a SpanStub because sdktrace.ReadOnlySpan seals itself with an
// unexported method, so no local type can wrap a span and add an attribute to it.
func markStart(readOnly sdktrace.ReadOnlySpan) sdktrace.ReadOnlySpan {
	stub := tracetest.SpanStubFromReadOnlySpan(readOnly)
	stub.Attributes = append(append([]attribute.KeyValue{}, stub.Attributes...),
		attribute.String(AttrRecord, RecordStart))
	return stub.Snapshot()
}
