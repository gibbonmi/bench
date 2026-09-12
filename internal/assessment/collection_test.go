package assessment

import (
	"encoding/json"
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func collectionInput(t *testing.T, s Store, extra map[string]any) string {
	t.Helper()
	data, _ := json.Marshal(fixtureRun(s.Root))
	var object map[string]any
	json.Unmarshal(data, &object)
	for k, v := range extra {
		object[k] = v
	}
	data, _ = json.Marshal(object)
	path := filepath.Join(t.TempDir(), "input.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}
func selection(id string) map[string]any {
	return map[string]any{"id": id, "attempt_id": "attempt-1", "chunk_id": "1", "role": "implementation"}
}
func nativeSpan(t *testing.T, s Store, assignment string, finished bool, overrides ...map[string]any) {
	t.Helper()
	span := map[string]any{"name": "worktree.exec", "traceId": "trace-1", "spanId": "span-1", "startTimeUnixNano": "1000000000", "attributes": []any{map[string]any{"key": otelrecord.AttrSeam, "value": map[string]any{"stringValue": "worktree.exec"}}, map[string]any{"key": otelrecord.AttrSubjectID, "value": map[string]any{"stringValue": assignment}}}}
	if finished {
		span["endTimeUnixNano"] = "4000000000"
	}
	for _, fields := range overrides {
		for key, value := range fields {
			span[key] = value
		}
	}
	data, _ := json.Marshal(map[string]any{"resourceSpans": []any{map[string]any{"scopeSpans": []any{map[string]any{"spans": []any{span}}}}}})
	if err := otelrecord.NewWriter(s.Home, s.Root).Append(data); err != nil {
		t.Fatal(err)
	}
}
func TestAssessmentCollection(t *testing.T) {
	const assignment = "11111111111111111111111111111111"
	t.Run("A30 selected span", func(t *testing.T) {
		s := Store{Home: t.TempDir(), Root: t.TempDir()}
		nativeSpan(t, s, assignment, true)
		input := collectionInput(t, s, map[string]any{"bench_inputs": map[string]any{"assignment_id": assignment, "trace_ids": []any{selection("trace-1")}}})
		if out, code := Command(s, []string{"record", "--input", input}); code != 0 {
			t.Fatal(out)
		}
		out, code := Command(s, []string{"show", "run-1"})
		if code != 0 || !strings.Contains(out, "wall_seconds") || !strings.Contains(out, "span-1") {
			t.Fatalf("selected span missing: %s", out)
		}
	})
	t.Run("A31 selected census", func(t *testing.T) {
		s := Store{Home: t.TempDir(), Root: t.TempDir()}
		dir := census.Dir(s.Home, s.Root)
		os.MkdirAll(dir, 0700)
		os.WriteFile(filepath.Join(dir, assignment), []byte("2026-01-01T00:00:00Z\tgit\n2026-01-01T00:00:01Z\tgo\n2026-01-01T00:00:02Z\trm\n"), 0600)
		input := collectionInput(t, s, map[string]any{"bench_inputs": map[string]any{"assignment_id": assignment, "census_event_ids": []any{selection(assignment + ":1"), selection(assignment + ":2")}}})
		if out, code := Command(s, []string{"record", "--input", input}); code != 0 {
			t.Fatal(out)
		}
		r, err := s.Read("run-1")
		if err != nil || len(r.Attempts[0].Measures) != 2 {
			t.Fatalf("selected census count mismatch: %+v %v", r, err)
		}
		out, _ := Command(s, []string{"show", "run-1"})
		if !strings.Contains(out, "raw_commands") || !strings.Contains(out, assignment+":2") {
			t.Fatalf("census missing: %s", out)
		}
	})
	t.Run("A32 mapped harness", func(t *testing.T) {
		s := Store{Home: t.TempDir(), Root: t.TempDir()}
		nativeSpan(t, s, assignment, true)
		path := filepath.Join(t.TempDir(), "native.json")
		os.WriteFile(path, []byte(`{"type":"event_msg","payload":{"type":"token_count","info":{"total_token_usage":{"input_tokens":100,"cached_input_tokens":80,"output_tokens":5}}}}`), 0600)
		input := collectionInput(t, s, map[string]any{"bench_inputs": map[string]any{"assignment_id": assignment, "trace_ids": []any{selection("trace-1")}}, "harness_inputs": []any{map[string]any{"path": path, "format": "codex-token-count-v1", "session_id": "session-1", "epoch": 0, "sequence": 1, "event_id": "native-1", "counter": "total_token_usage", "mode": "cumulative", "attempt_id": "attempt-1", "chunk_id": "1", "role": "implementation"}}})
		if out, code := Command(s, []string{"record", "--input", input}); code != 0 {
			t.Fatal(out)
		}
		r, err := s.Read("run-1")
		if err != nil {
			t.Fatal(err)
		}
		if len(r.Attempts[0].Intervals) == 0 {
			t.Fatal("Bench measures missing beside harness usage")
		}
		u, err := UsageTotal(r.Attempts[0].Usage)
		if err != nil || u.InputUncached == nil || *u.InputUncached != 20 || *u.InputCached != 80 || *u.Output != 5 {
			t.Fatalf("native counters missing: %+v %v", u, err)
		}
	})
	for _, kind := range []string{"foreign", "ambiguous"} {
		t.Run("A33 "+kind, func(t *testing.T) {
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			nativeSpan(t, s, "22222222222222222222222222222222", true)
			if kind == "ambiguous" {
				nativeSpan(t, s, assignment, true)
			}
			input := collectionInput(t, s, map[string]any{"bench_inputs": map[string]any{"assignment_id": assignment, "trace_ids": []any{selection("trace-1")}}})
			if out, code := Command(s, []string{"record", "--input", input}); code != 1 {
				t.Fatalf("unsafe join accepted: %s", out)
			}
			if _, err := s.Read("run-1"); !os.IsNotExist(err) {
				t.Fatalf("unsafe join stored: %v", err)
			}
		})
	}
	t.Run("A23 unfinished", func(t *testing.T) {
		s := Store{Home: t.TempDir(), Root: t.TempDir()}
		nativeSpan(t, s, assignment, false)
		input := collectionInput(t, s, map[string]any{"bench_inputs": map[string]any{"assignment_id": assignment, "trace_ids": []any{selection("trace-1")}}})
		if out, code := Command(s, []string{"record", "--input", input}); code != 0 {
			t.Fatal(out)
		}
		out, _ := Command(s, []string{"show", "run-1"})
		if !strings.Contains(out, "unfinished") || !strings.Contains(out, "diagnostics") {
			t.Fatal(out)
		}
	})
}

func TestAssessmentCollectionSpanCoverage(t *testing.T) {
	const assignment = "11111111111111111111111111111111"
	for _, kind := range []string{"start then end", "disjoint", "nested", "malformed", "two traces"} {
		t.Run(kind, func(t *testing.T) {
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			if kind == "start then end" {
				nativeSpan(t, s, assignment, false)
			}
			nativeSpan(t, s, assignment, true)
			want := 3.0
			if kind == "disjoint" {
				nativeSpan(t, s, assignment, true, map[string]any{"spanId": "span-2", "startTimeUnixNano": "10000000000", "endTimeUnixNano": "13000000000"})
				want = 6
			}
			if kind == "two traces" {
				nativeSpan(t, s, assignment, true, map[string]any{"traceId": "trace-2", "startTimeUnixNano": "10000000000", "endTimeUnixNano": "13000000000"})
				want = 6
			}
			if kind == "nested" {
				nativeSpan(t, s, assignment, true, map[string]any{"spanId": "span-2", "startTimeUnixNano": "2000000000", "endTimeUnixNano": "3000000000"})
			}
			if kind == "malformed" {
				if err := otelrecord.NewWriter(s.Home, s.Root).Append([]byte("{")); err != nil {
					t.Fatal(err)
				}
			}
			traceIDs := []any{selection("trace-1")}
			if kind == "two traces" {
				traceIDs = append(traceIDs, selection("trace-2"))
			}
			input := collectionInput(t, s, map[string]any{"bench_inputs": map[string]any{"assignment_id": assignment, "trace_ids": traceIDs}})
			if out, code := Command(s, []string{"record", "--input", input}); code != 0 {
				t.Fatal(out)
			}
			r, err := s.Read("run-1")
			if err != nil {
				t.Fatal(err)
			}
			summary, err := Summarize(r)
			if err != nil || summary.WallSeconds == nil || *summary.WallSeconds != want || summary.EffortSeconds != want {
				t.Fatalf("span union mismatch: %+v %v", summary, err)
			}
			if (len(r.Diagnostics) > 0) != (kind == "malformed") || summary.Incomplete != (kind == "malformed") {
				t.Fatalf("coverage mismatch: %+v %+v", r.Diagnostics, summary)
			}
		})
	}
}

func TestAssessmentCollectionNativeGaps(t *testing.T) {
	const assignment = "11111111111111111111111111111111"
	for _, kind := range []string{"census malformed", "census unfinished", "census foreign", "census ambiguous", "harness malformed", "harness missing", "harness semantics"} {
		t.Run(kind, func(t *testing.T) {
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			extra := map[string]any{}
			refuse := false
			if strings.HasPrefix(kind, "census") {
				dir := census.Dir(s.Home, s.Root)
				os.MkdirAll(dir, 0700)
				body := "2026-01-01T00:00:00Z\tgit\n"
				if kind == "census malformed" {
					body = "invalid\tgit\n"
				}
				if kind == "census unfinished" {
					body = strings.TrimSuffix(body, "\n")
				}
				os.WriteFile(filepath.Join(dir, assignment), []byte(body), 0600)
				selected := []any{selection(assignment + ":1")}
				if kind == "census foreign" {
					selected = []any{selection("22222222222222222222222222222222:1")}
					refuse = true
				}
				if kind == "census ambiguous" {
					other := selection(assignment + ":1")
					other["role"] = "review"
					selected = append(selected, other)
					refuse = true
				}
				extra["bench_inputs"] = map[string]any{"assignment_id": assignment, "census_event_ids": selected}
			} else {
				path := filepath.Join(t.TempDir(), "native.json")
				if kind != "harness missing" {
					os.WriteFile(path, []byte("{"), 0600)
				}
				mode := "cumulative"
				if kind == "harness semantics" {
					mode = "guess"
					refuse = true
				}
				extra["harness_inputs"] = []any{map[string]any{"path": path, "format": "codex-token-count-v1", "session_id": "session-1", "epoch": 0, "sequence": 1, "event_id": "native-gap", "counter": "total_token_usage", "mode": mode, "attempt_id": "attempt-1", "chunk_id": "1", "role": "implementation"}}
			}
			input := collectionInput(t, s, extra)
			out, code := Command(s, []string{"record", "--input", input})
			if refuse {
				if code != 1 {
					t.Fatalf("ambiguous input accepted: %s", out)
				}
				if _, err := s.Read("run-1"); !os.IsNotExist(err) {
					t.Fatalf("unsafe join stored: %v", err)
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
			summary, err := Summarize(r)
			if err != nil || len(r.Diagnostics) == 0 || !summary.Incomplete {
				t.Fatalf("native gap appears complete: %+v %+v %v", r.Diagnostics, summary, err)
			}
		})
	}
}
