package harnesstranscript

import (
	"path/filepath"
	"strings"
	"testing"
)

// The expectations here are authored independently of the mapping. Deriving a count or a row
// name from Metrics or from the tally would let a dropped metric and a miscount disappear
// from both the code and its check at once.

// The fixture is a minimized record of the pinned source. It keeps the keys, the nesting, and
// the identity fields of the reviewed transcript, and it carries none of that session's
// commands, paths, messages, or reasoning.
func fixture(t *testing.T) Record {
	t.Helper()
	record, failure := Read(filepath.Join("testdata", "codex-rollout.jsonl"), SourceCodexRollout)
	if failure.State != "" {
		t.Fatalf("read fixture: %s (%s)", failure.State, failure.Reason)
	}
	return record
}

func observation(t *testing.T, record Record, metric string) Observation {
	t.Helper()
	for _, row := range record.Observations {
		if row.Metric == metric {
			return row
		}
	}
	t.Fatalf("record holds no %q observation", metric)
	return Observation{}
}

func wantCount(t *testing.T, record Record, metric string, count int64) {
	t.Helper()
	row := observation(t, record, metric)
	if row.Availability != Observed || row.Count != count {
		t.Fatalf("%s = %d (%s), want %d observed", metric, row.Count, row.Availability, count)
	}
}

func TestMetricInventoryIsTheRenderedRowSet(t *testing.T) {
	// The inventory is the rendered row count and the rendered row names, so a duplicate or
	// an unnamed unit would ship a row a reader cannot key on.
	if len(Metrics) != 14 {
		t.Fatalf("inventory holds %d metrics, want 14", len(Metrics))
	}
	seen := map[string]bool{}
	for _, metric := range Metrics {
		if metric.Name == "" || metric.Unit == "" {
			t.Fatalf("metric %+v lacks a name or a unit", metric)
		}
		if seen[metric.Name] {
			t.Fatalf("metric %q is declared twice", metric.Name)
		}
		seen[metric.Name] = true
	}
}

func TestUnknownMeasureRendersNoNumber(t *testing.T) {
	unknown := Observation{Availability: Unknown, Count: 7}
	if cell := unknown.Cell(); cell != "" {
		t.Fatalf("unknown cell = %v, want an empty cell", cell)
	}
	// An observed zero is a counted fact and keeps its number.
	if cell := (Observation{Availability: Observed}).Cell(); cell != int64(0) {
		t.Fatalf("observed zero cell = %v, want 0", cell)
	}
	if cell := (Observation{Availability: Incomplete, Count: 3}).Cell(); cell != int64(3) {
		t.Fatalf("incomplete cell = %v, want the partial count 3", cell)
	}
	// A partial count of nothing carries no fact, and a rendered 0 belongs to an observed
	// measure alone.
	if cell := (Observation{Availability: Incomplete}).Cell(); cell != "" {
		t.Fatalf("incomplete zero cell = %v, want an empty cell", cell)
	}
}

func TestCodexFixtureCountsItsPinnedShapes(t *testing.T) {
	record := fixture(t)
	// The one result holds 7 UTF-8 bytes across 6 characters on 2 lines. The reasoning
	// payload and the call metadata beside it are not model-visible result text.
	wantCount(t, record, MetricResultBytes, 7)
	wantCount(t, record, MetricResultChars, 6)
	wantCount(t, record, MetricResultLines, 2)
	// Two calls, one of which never recorded a result. The subagent event names a thread,
	// not a call.
	wantCount(t, record, MetricOuterCalls, 2)
	wantCount(t, record, MetricUnmatchedCalls, 1)
	wantCount(t, record, MetricTurns, 1)
	wantCount(t, record, MetricCompactions, 1)
}

func TestCodexFixtureKeepsCumulativeUsageSemantics(t *testing.T) {
	record := fixture(t)
	// The fixture holds two cumulative snapshots and one per-turn usage record. The last
	// snapshot is the session total; a sum would report 350 input tokens nobody spent, and
	// the per-turn record is a different fact with a different interval.
	wantCount(t, record, MetricInputTokens, 250)
	wantCount(t, record, MetricCachedTokens, 90)
	wantCount(t, record, MetricOutputTokens, 55)
	wantCount(t, record, MetricReasoningTokens, 30)
}

func TestCodexFixtureIdentifiesItsBytes(t *testing.T) {
	record := fixture(t)
	if !strings.HasPrefix(record.Digest, "sha256:") || len(record.Digest) != len("sha256:")+64 {
		t.Fatalf("digest = %q, want a sha256 hex digest", record.Digest)
	}
	if record.SourceID != SourceCodexRollout {
		t.Fatalf("source id = %q, want %q", record.SourceID, SourceCodexRollout)
	}
	if record.Interval != "2026-09-11T09:00:11.841Z/2026-09-11T09:00:21.000Z" {
		t.Fatalf("interval = %q, want the first and last record timestamps", record.Interval)
	}
}

func TestUnsupportedSourceReadsNothing(t *testing.T) {
	record, failure := Read(filepath.Join("testdata", "codex-rollout.jsonl"), "opencode-session")
	if failure.State != "" {
		t.Fatalf("unsupported source failed the read: %s", failure.State)
	}
	if record.Digest != "" || record.Interval != "" {
		t.Fatalf("unsupported record = %+v, want no digest and no interval: the bytes were never read", record)
	}
	for _, row := range record.Observations {
		if row.Availability != Unknown {
			t.Fatalf("%s reads %s under an unpinned source, want unknown", row.Metric, row.Availability)
		}
	}
	if Supported("opencode-session") {
		t.Fatal("Supported accepts an unpinned source id")
	}
}
