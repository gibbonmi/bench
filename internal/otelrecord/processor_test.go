package otelrecord

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// The processor tests drive the provider the way a verb does: a temporary home, a real
// span, then a read of the JSON lines back off the disk. The line content is the
// external behavior, so the assertions parse the file rather than the encoder's return.

// recordLines returns the record's parsed lines for root below home.
func recordLines(t *testing.T, home, root string) []map[string]any {
	t.Helper()

	raw, err := os.ReadFile(Path(home, root))
	if err != nil {
		t.Fatalf("read the record: %v", err)
	}
	var lines []map[string]any
	for _, text := range strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n") {
		if text == "" {
			continue
		}
		var parsed map[string]any
		if err := json.Unmarshal([]byte(text), &parsed); err != nil {
			t.Fatalf("the record line does not parse: %v", err)
		}
		lines = append(lines, parsed)
	}
	return lines
}

// spanOf digs the one span out of a parsed record line.
func spanOf(t *testing.T, line map[string]any) map[string]any {
	t.Helper()

	resourceSpans, ok := line["resourceSpans"].([]any)
	if !ok || len(resourceSpans) != 1 {
		t.Fatalf("the line carries no single resourceSpans entry: %v", line)
	}
	scopeSpans, ok := resourceSpans[0].(map[string]any)["scopeSpans"].([]any)
	if !ok || len(scopeSpans) != 1 {
		t.Fatalf("the line carries no single scopeSpans entry: %v", line)
	}
	spans, ok := scopeSpans[0].(map[string]any)["spans"].([]any)
	if !ok || len(spans) != 1 {
		t.Fatalf("the line carries no single span: %v", line)
	}
	span, ok := spans[0].(map[string]any)
	if !ok {
		t.Fatalf("the span is not an object: %v", line)
	}
	return span
}

// attributeOf returns the named span attribute's string value, and reports whether the
// span carries the key at all.
func attributeOf(t *testing.T, span map[string]any, key string) (string, bool) {
	t.Helper()

	pairs, _ := span["attributes"].([]any)
	for _, entry := range pairs {
		pair, ok := entry.(map[string]any)
		if !ok || pair["key"] != key {
			continue
		}
		value, _ := pair["value"].(map[string]any)
		text, ok := value["stringValue"].(string)
		if !ok {
			t.Fatalf("attribute %s carries no string value: %v", key, pair)
		}
		return text, true
	}
	return "", false
}

// startRecordedSpan builds a provider on a temporary home and starts one seam span. It
// returns the home, the root, and the span's end function.
func startRecordedSpan(t *testing.T) (string, string, func()) {
	t.Helper()

	home := t.TempDir()
	root := t.TempDir()
	provider := NewProvider(home, root)
	t.Cleanup(func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			t.Fatalf("shut the provider down: %v", err)
		}
	})

	ctx := WithTracer(context.Background(), provider.Tracer())
	_, span := TracerFrom(ctx).Start(context.Background(), "bench.gate",
		trace.WithAttributes(attribute.String(AttrSeam, "gate")))
	return home, root, func() { span.End() }
}

// TestARecordLineExistsBeforeTheSpanEnds holds row OT5: an end-only processor writes
// nothing while the span runs, and this read reds.
func TestARecordLineExistsBeforeTheSpanEnds(t *testing.T) {
	home, root, end := startRecordedSpan(t)
	defer end()

	if lines := recordLines(t, home, root); len(lines) != 1 {
		t.Fatalf("the record holds %d lines before the span ends, want 1", len(lines))
	}
}

// TestTheStartLineCarriesTheStartMarker holds row OT6: an unmarked start line reads as
// a finished span, so both the marker and the absent end time are asserted.
func TestTheStartLineCarriesTheStartMarker(t *testing.T) {
	home, root, end := startRecordedSpan(t)
	defer end()

	span := spanOf(t, recordLines(t, home, root)[0])
	marker, ok := attributeOf(t, span, AttrRecord)
	if !ok {
		t.Fatalf("the start line carries no %s attribute: %v", AttrRecord, span)
	}
	if marker != RecordStart {
		t.Fatalf("%s = %q, want %q", AttrRecord, marker, RecordStart)
	}
	if _, present := span["endTimeUnixNano"]; present {
		t.Fatalf("the start line carries an end time: %v", span)
	}
}

// TestTheEndLineCarriesTheCompleteSpan holds row OT7: a start-only writer never writes
// the end time, and the elapsed-time read reds.
func TestTheEndLineCarriesTheCompleteSpan(t *testing.T) {
	home, root, end := startRecordedSpan(t)
	end()

	lines := recordLines(t, home, root)
	if len(lines) != 2 {
		t.Fatalf("the record holds %d lines after the span ends, want 2", len(lines))
	}
	span := spanOf(t, lines[1])
	if endTime, _ := span["endTimeUnixNano"].(string); endTime == "" {
		t.Fatalf("the end line carries no endTimeUnixNano: %v", span)
	}
	if span["name"] != "bench.gate" {
		t.Fatalf("the end line names %v, want bench.gate", span["name"])
	}
	if seam, ok := attributeOf(t, span, AttrSeam); !ok || seam != "gate" {
		t.Fatalf("the end line carries %s = %q, want gate", AttrSeam, seam)
	}
	if _, marked := attributeOf(t, span, AttrRecord); marked {
		t.Fatalf("the end line carries the start marker: %v", span)
	}
}

// TestTheContextCarriesTheProviderTracer holds the provider's context contract: a
// caller reads back the tracer the boundary put there, and an unset context answers a
// tracer that records nothing.
func TestTheContextCarriesTheProviderTracer(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	provider := NewProvider(home, root)
	defer func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			t.Fatalf("shut the provider down: %v", err)
		}
	}()

	ctx := WithTracer(context.Background(), provider.Tracer())
	_, span := TracerFrom(ctx).Start(ctx, "bench.commit")
	span.End()
	if lines := recordLines(t, home, root); len(lines) != 2 {
		t.Fatalf("the context tracer wrote %d lines, want 2", len(lines))
	}

	_, quiet := TracerFrom(context.Background()).Start(context.Background(), "bench.commit")
	defer quiet.End()
	if quiet.IsRecording() {
		t.Fatal("an unset context answers a recording tracer, want the no-op tracer")
	}
}

// TestTheConstructorResolvesTheHomeWhenTheCallerHasNone holds the constructor's home
// read: an empty home resolves through internal/benchhome rather than a second read.
func TestTheConstructorResolvesTheHomeWhenTheCallerHasNone(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv(benchhome.Env, home)

	provider := NewProvider("", root)
	defer func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			t.Fatalf("shut the provider down: %v", err)
		}
	}()

	_, span := provider.Tracer().Start(context.Background(), "bench.gate")
	span.End()
	if lines := recordLines(t, home, root); len(lines) != 2 {
		t.Fatalf("the resolved home holds %d lines, want 2", len(lines))
	}
}
