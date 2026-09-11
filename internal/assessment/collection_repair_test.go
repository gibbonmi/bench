package assessment

import (
	"encoding/json"
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"os"
	"path/filepath"
	"testing"
)

func TestAssessmentCollectionPublication(t *testing.T) {
	for _, seam := range []string{"commit", "worktree.land"} {
		t.Run(seam, func(t *testing.T) {
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			assignment := "11111111111111111111111111111111"
			attrs := []any{}
			for key, value := range map[string]string{otelrecord.AttrSeam: seam, otelrecord.AttrSubjectID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "bench.assignment.id": assignment, otelrecord.AttrMeasurePathCount: "3"} {
				attrs = append(attrs, map[string]any{"key": key, "value": map[string]any{"stringValue": value}})
			}
			nativeSpan(t, s, assignment, true, map[string]any{"name": seam, "attributes": attrs})
			input := collectionInput(t, s, map[string]any{"bench_inputs": map[string]any{"assignment_id": assignment, "trace_ids": []any{selection("trace-1")}}})
			if out, code := Command(s, []string{"record", "--input", input}); code != 0 {
				t.Fatal(out)
			}
			r, err := s.Read("run-1")
			if err != nil {
				t.Fatal(err)
			}
			m := r.Attempts[0].Measures["diff_paths/trace-1/span-1"]
			if m.Value == nil || *m.Value != 3 || m.Reference.Native != otelrecord.Path(s.Home, s.Root)+"#trace-1/span-1" {
				t.Fatalf("diff metadata missing: %+v", m)
			}
		})
	}
}

func TestAssessmentCollectionPartial(t *testing.T) {
	for _, missing := range []string{"input_tokens", "cached_input_tokens", "output_tokens"} {
		t.Run(missing, func(t *testing.T) {
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			r := fixtureRun(s.Root)
			counters := map[string]int{"input_tokens": 100, "cached_input_tokens": 80, "output_tokens": 5}
			delete(counters, missing)
			data, _ := json.Marshal(map[string]any{"type": "event_msg", "payload": map[string]any{"type": "token_count", "info": map[string]any{"total_token_usage": counters}}})
			path := filepath.Join(t.TempDir(), "native.json")
			if err := os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
			r.HarnessInputs = []HarnessInput{{Mapping: Mapping{"attempt-1", "1", "implementation"}, Path: path, Format: "codex-token-count-v1", SessionID: "session-1", EventID: "partial", Counter: "total_token_usage", Mode: CumulativeMode}}
			got, err := Collect(s, r)
			if err != nil {
				t.Fatal(err)
			}
			if len(got.Diagnostics) == 0 {
				t.Fatal("partial counters appear complete")
			}
			u := got.Attempts[0].Usage[len(got.Attempts[0].Usage)-1].Usage
			if (missing == "input_tokens" && u.InputTotal != nil) || (missing == "cached_input_tokens" && u.InputCached != nil) || (missing == "output_tokens" && u.Output != nil) {
				t.Fatalf("unknown became known: %+v", u)
			}
		})
	}
}

func TestAssessmentCollectionSecondAttempt(t *testing.T) {
	for _, producer := range []string{"trace", "census", "harness"} {
		t.Run(producer, func(t *testing.T) {
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			r := fixtureRun(s.Root)
			second := r.Attempts[0]
			second.AttemptID = "attempt-2"
			second.ChunkID = "2"
			second.Role = "review"
			second.SessionID = "session-2"
			r.Attempts = append(r.Attempts, second)
			m := Mapping{second.AttemptID, second.ChunkID, second.Role}
			assignment := "11111111111111111111111111111111"
			r.BenchInputs = &BenchInputs{AssignmentID: assignment}
			switch producer {
			case "trace":
				nativeSpan(t, s, assignment, true)
				r.BenchInputs.TraceIDs = []Selection{{ID: "trace-1", Mapping: m}}
			case "census":
				dir := census.Dir(s.Home, s.Root)
				os.MkdirAll(dir, 0700)
				os.WriteFile(filepath.Join(dir, assignment), []byte("2026-01-01T00:00:00Z\tgit\n"), 0600)
				r.BenchInputs.CensusEventIDs = []Selection{{ID: assignment + ":1", Mapping: m}}
			case "harness":
				path := filepath.Join(t.TempDir(), "native.json")
				os.WriteFile(path, []byte(`{"type":"event_msg","payload":{"type":"token_count","info":{"last_token_usage":{"input_tokens":100,"cached_input_tokens":80,"output_tokens":5}}}}`), 0600)
				r.HarnessInputs = []HarnessInput{{Mapping: m, Path: path, Format: "codex-token-count-v1", SessionID: second.SessionID, EventID: "delta", Counter: "last_token_usage", Mode: DeltaMode}}
			}
			got, err := Collect(s, r)
			if err != nil {
				t.Fatal(err)
			}
			first, last := got.Attempts[0], got.Attempts[1]
			if len(first.Usage)+len(first.Measures)+len(first.Intervals) != 0 {
				t.Fatal("evidence reached first attempt")
			}
			if producer == "harness" {
				if len(last.Usage) != 1 || last.Usage[0].Mode != DeltaMode {
					t.Fatalf("delta missing: %+v", last)
				}
				u, err := UsageTotal(last.Usage)
				if err != nil || u.InputUncached == nil || *u.InputUncached != 20 || u.InputCached == nil || *u.InputCached != 80 || u.Output == nil || *u.Output != 5 {
					t.Fatalf("delta counters: %+v %v", u, err)
				}
			} else if len(last.Measures) == 0 || len(last.Evidence) == 0 {
				t.Fatal("second attempt evidence missing")
			}
		})
	}
}

func TestAssessmentCollectionHostileHarness(t *testing.T) {
	for _, kind := range []string{"parent symlink", "present empty"} {
		t.Run(kind, func(t *testing.T) {
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			dir := t.TempDir()
			path := filepath.Join(dir, "native.json")
			if err := os.WriteFile(path, nil, 0600); err != nil {
				t.Fatal(err)
			}
			if kind == "parent symlink" {
				link := filepath.Join(t.TempDir(), "link")
				if err := os.Symlink(dir, link); err != nil {
					t.Fatal(err)
				}
				path = filepath.Join(link, "native.json")
			}
			input := collectionInput(t, s, map[string]any{"harness_inputs": []HarnessInput{{Mapping: Mapping{"attempt-1", "1", "implementation"}, Path: path, Format: "codex-token-count-v1", SessionID: "session-1", EventID: "hostile", Counter: "total_token_usage", Mode: CumulativeMode}}})
			out, code := Command(s, []string{"record", "--input", input})
			if kind == "parent symlink" {
				if code != 1 {
					t.Fatal(out)
				}
				if _, err := s.Read("run-1"); !os.IsNotExist(err) {
					t.Fatal("unsafe run stored", err)
				}
				return
			}
			if code != 0 {
				t.Fatal(out)
			}
			r, err := s.Read("run-1")
			if err != nil {
				t.Fatal(err)
			}
			if len(r.Diagnostics) == 0 || len(r.Attempts[0].Usage) != 0 {
				t.Fatalf("empty fragment lost: %+v", r)
			}
			if r.Diagnostics[0].Native != path+": malformed native input" {
				t.Fatalf("empty fragment misclassified: %+v", r.Diagnostics)
			}
		})
	}
}
