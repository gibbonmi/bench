package harnesses

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/harnesstranscript"
)

// The expectations here are authored independently of the reader. Deriving a row, a count,
// or a baseline view from the metric inventory or from the projection would let a dropped
// metric, a miscounted byte, or a changed compiled view disappear from both sides at once.

// The pinned-source event shapes each fixture builds. A fixture keeps the keys, the nesting,
// and the identity fields of the real record and carries no content from it.
func call(id string) string {
	return fmt.Sprintf(`{"timestamp":"2026-09-11T09:00:00.000Z","type":"response_item","payload":{"type":"custom_tool_call","call_id":%q,"name":"exec"}}`, id)
}

func result(id, text string) string {
	return fmt.Sprintf(`{"timestamp":"2026-09-11T09:00:01.000Z","type":"response_item","payload":{"type":"custom_tool_call_output","call_id":%q,"output":[{"type":"input_text","text":%q}]}}`, id, text)
}

// A typed function call is the record's second invocation shape, and its completion carries
// the result as one string rather than as an array of text parts.
func fnCall(id string) string {
	return fmt.Sprintf(`{"timestamp":"2026-09-11T09:00:00.500Z","type":"response_item","payload":{"type":"function_call","call_id":%q,"name":"spawn_agent","namespace":"collaboration"}}`, id)
}
func fnResult(id, text string) string {
	return fmt.Sprintf(`{"timestamp":"2026-09-11T09:00:01.500Z","type":"response_item","payload":{"type":"function_call_output","call_id":%q,"output":%q}}`, id, text)
}

// usageSnapshotBody wraps one counter list, so a fixture can omit a key the producer never
// measured instead of sending a zero it never reported.
func usageSnapshotBody(counters string) string {
	return fmt.Sprintf(`{"timestamp":"2026-09-11T09:00:02.000Z","type":"event_msg","payload":{"type":"token_count","info":{"model_context_window":400,"total_token_usage":{%s}}}}`, counters)
}
func usageSnapshot(input, cached, output, reasoning int) string {
	return usageSnapshotBody(fmt.Sprintf(`"input_tokens":%d,"cached_input_tokens":%d,"output_tokens":%d,"reasoning_output_tokens":%d`, input, cached, output, reasoning))
}

func completed(itemType, id, body string) string {
	return fmt.Sprintf(`{"timestamp":"2026-09-11T09:00:03.000Z","type":"event_msg","payload":{"type":"item_completed","thread_id":"t1","item":{"type":%q,"id":%q%s}}}`, itemType, id, body)
}

const (
	taskStarted = `{"timestamp":"2026-09-11T09:00:04.000Z","type":"event_msg","payload":{"type":"task_started","turn_id":"turn-1"}}`
	compacted   = `{"timestamp":"2026-09-11T09:00:05.000Z","type":"compacted","payload":{"window_number":1}}`
	reasoning   = `{"timestamp":"2026-09-11T09:00:06.000Z","type":"response_item","payload":{"type":"reasoning","id":"r1","encrypted_content":"AAAAAAAAAAAAAAAA"}}`
	// multibyte is the one text fixture whose bytes, characters, and lines all differ.
	multibyte = "héllo\nwörld"
)

// observe writes lines as a record in a temporary directory and runs the codex record view
// over it, so every test exercises the real path boundary rather than a stub.
func observe(t *testing.T, lines ...string) (string, int) {
	t.Helper()
	return observeAt(t, record(t, lines...), harnesstranscript.SourceCodexRollout)
}

func record(t *testing.T, lines ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rollout.jsonl")
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o600); err != nil {
		t.Fatalf("write record: %v", err)
	}
	return path
}

func observeAt(t *testing.T, path, format string) (string, int) {
	t.Helper()
	return Command([]string{"codex", "--record", path, "--format", format})
}

// wantRow asserts one observation row by its leading cells, because a trailing cell carries
// commas and is quoted.
func wantRow(t *testing.T, out, prefix string) {
	t.Helper()
	if !strings.Contains(out, "\n  "+prefix) {
		t.Fatalf("record view = %q, want an observation row starting %q", out, prefix)
	}
}

func TestObservedTextBoundary(t *testing.T) {
	out, code := observe(t, call("c1"), result("c1", multibyte), reasoning, fnCall("c2"), fnResult("c2", "ok\n"))
	if code != 0 {
		t.Fatalf("record view exit = %d, want 0", code)
	}
	// The array-form result holds 13 bytes over 11 characters on 2 lines, and the string-form
	// result adds 3 of each. The reasoning record's encrypted payload counts nothing.
	wantRow(t, out, "result-text bytes,16,bytes,observed,")
	wantRow(t, out, "result-text characters,14,characters,observed,")
	wantRow(t, out, "result-text lines,3,lines,observed,")
}

func TestObservedCallAncestry(t *testing.T) {
	subagent := completed("SubAgentActivity", "i1", `,"agent_thread_id":"sub-1","kind":"started"`)
	out, _ := observe(t, call("c1"), result("c1", "ok"), subagent)
	// A subagent event names a thread, not a call, so it raises no census here.
	wantRow(t, out, "outer calls,1,calls,observed,")
	wantRow(t, out, `observed nested calls,"",calls,unknown,`)
}

func TestObservedDuplicateCompletion(t *testing.T) {
	out, _ := observe(t, call("c1"), result("c1", "ok"), result("c1", "ok"))
	// Both completions carry one invocation identity, so the call and its bytes count once.
	wantRow(t, out, "outer calls,1,calls,observed,")
	wantRow(t, out, "unmatched calls,0,calls,observed,")
	wantRow(t, out, "result-text bytes,2,bytes,observed,")
}

func TestObservedUnsupportedFormat(t *testing.T) {
	out, code := observeAt(t, record(t, call("c1"), result("c1", multibyte)), "claude-jsonl")
	if code != 1 {
		t.Fatalf("unsupported format exit = %d, want 1", code)
	}
	if strings.Count(out, ",unknown,") != 14 {
		t.Fatalf("record view = %q, want all 14 observations unknown", out)
	}
	if !strings.Contains(out, `  "",claude-jsonl,codex,""`) {
		t.Fatalf("record view = %q, want an identity row with no digest and no interval", out)
	}
}

func TestObservedEmptyRecord(t *testing.T) {
	out, code := observe(t)
	if code != 0 {
		t.Fatalf("empty record exit = %d, want 0", code)
	}
	// An empty record is an authoritative observation of nothing, not a missing input.
	wantRow(t, out, "outer calls,0,calls,observed,")
	wantRow(t, out, "result-text bytes,0,bytes,observed,")
}

func TestObservedAbsentRecord(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.jsonl")
	out, code := observeAt(t, path, harnesstranscript.SourceCodexRollout)
	if code != 1 {
		t.Fatalf("absent record exit = %d, want 1", code)
	}
	if !strings.HasPrefix(out, "error: "+path+" is absent") {
		t.Fatalf("absent record = %q, want the record error line", out)
	}
}

func TestObservedMalformedEvent(t *testing.T) {
	out, code := observe(t, call("c1"), result("c1", "ok"), "{not json", taskStarted)
	if code != 0 {
		t.Fatalf("malformed record exit = %d, want 0", code)
	}
	// A line that does not parse could have carried any dimension, so every count it could
	// have raised reports its partial total as incomplete.
	wantRow(t, out, "result-text bytes,2,bytes,incomplete,")
	wantRow(t, out, "outer calls,1,calls,incomplete,")
	wantRow(t, out, "turns,1,turns,incomplete,")
}

func TestObservedNativeUsage(t *testing.T) {
	out, _ := observe(t, usageSnapshot(10, 2, 3, 4), usageSnapshot(30, 6, 9, 12), usageSnapshot(60, 12, 18, 24))
	// The snapshots are cumulative, so a sum would report 100 input tokens nobody spent.
	wantRow(t, out, "input tokens,60,tokens,observed,")
	wantRow(t, out, "output tokens,18,tokens,observed,")
}

func TestObservedNoTokenAttribution(t *testing.T) {
	out, _ := observe(t, call("c1"), result("c1", multibyte), usageSnapshot(60, 12, 18, 24))
	// Bytes and session totals are both present, and neither divides into a per-result cost.
	wantRow(t, out, `per-result token attribution,"",tokens,unknown,`)
}

func TestObservedCompactions(t *testing.T) {
	out, _ := observe(t, compacted, compacted, usageSnapshot(60, 12, 18, 24))
	// Only an identified compaction record counts; a small context window is not evidence.
	wantRow(t, out, "compactions,2,compactions,observed,")
	bare, _ := observe(t, usageSnapshot(60, 12, 18, 24))
	wantRow(t, bare, "compactions,0,compactions,observed,")
}

func TestObservedMeasureProvenance(t *testing.T) {
	out, _ := observe(t, call("c1"), result("c1", multibyte))
	_, rest, found := strings.Cut(out, "observations[14]{metric,value,unit,availability,source,boundary}:\n")
	block, _, ok := strings.Cut(rest, "help[0]")
	if !found || !ok {
		t.Fatalf("record view = %q, want a 14-row observations block", out)
	}
	for _, line := range strings.Split(strings.TrimSuffix(block, "\n"), "\n") {
		// The metric, its value, its unit, and its availability precede the provenance pair.
		cells := strings.SplitN(line, ",", 5)
		if len(cells) != 5 {
			t.Fatalf("observation row %q has no provenance cells", line)
		}
		source, boundary, split := strings.Cut(cells[4], ",")
		if !split || source == "" || source == `""` || boundary == "" || boundary == `""` {
			t.Fatalf("observation row %q names source %q and boundary %q, want both", line, source, boundary)
		}
	}
}

func TestObservedHostileRecord(t *testing.T) {
	dir := t.TempDir()
	sentinel := filepath.Join(dir, "sentinel")
	hostile := "$(touch " + sentinel + ")`touch " + sentinel + "`; touch " + sentinel
	out, code := observe(t, call("c1"), result("c1", hostile),
		completed("CommandExecution", "i1", `,"command":`+fmt.Sprintf("%q", hostile)))
	if code != 0 {
		t.Fatalf("hostile record exit = %d, want 0", code)
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("record contents created %s: the reader executed its input", sentinel)
	}
	// The command-shaped bytes are data, and they are counted as data.
	wantRow(t, out, fmt.Sprintf("result-text bytes,%d,bytes,observed,", len(hostile)))
}

func TestObservedMissingResult(t *testing.T) {
	out, _ := observe(t, call("c1"), fnCall("c2"), fnResult("c2", "ok"), result("c3", "ok"))
	// c1 never completed and c3's completion names no call, so two invocations lack a pair.
	wantRow(t, out, "unmatched calls,2,calls,observed,")
	wantRow(t, out, "result-text bytes,4,bytes,observed,")
	if strings.Contains(out, "complete producer output") {
		t.Fatalf("record view = %q, want no complete-producer-output claim", out)
	}
}

func TestObservedRegularFileBoundary(t *testing.T) {
	dir := t.TempDir()
	target := record(t, call("c1"), result("c1", "ok"))
	link := filepath.Join(dir, "linked.jsonl")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("Symlink: %v", err)
	}
	fifo := filepath.Join(dir, "fifo.jsonl")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatalf("Mkfifo: %v", err)
	}
	// A FIFO would block the reader inside open(2), and a link would supply bytes from
	// outside the named path, so the type check that precedes every open refuses both.
	for _, path := range []string{link, fifo} {
		out, code := observeAt(t, path, harnesstranscript.SourceCodexRollout)
		if code != 1 || !strings.HasPrefix(out, "error: "+path+" is wrong-type") {
			t.Fatalf("record view %s = %q exit %d, want a wrong-type refusal at exit 1", path, out, code)
		}
	}
}

// tokenDimension grades one native counter across a present value, an observed zero, a
// snapshot that omits its key, and a record with no snapshot. others names the counters the
// omitting snapshot keeps.
func tokenDimension(t *testing.T, metric string, present int, others string) {
	t.Helper()
	out, _ := observe(t, usageSnapshot(11, 22, 33, 44))
	wantRow(t, out, fmt.Sprintf("%s,%d,tokens,observed,", metric, present))
	zero, _ := observe(t, usageSnapshot(0, 0, 0, 0))
	wantRow(t, zero, metric+",0,tokens,observed,")
	omitted, _ := observe(t, usageSnapshotBody(others))
	wantRow(t, omitted, metric+`,"",tokens,unknown,`)
	absent, _ := observe(t, call("c1"), result("c1", "ok"))
	wantRow(t, absent, metric+`,"",tokens,unknown,`)
}

func TestObservedInputTokens(t *testing.T) {
	tokenDimension(t, "input tokens", 11, `"cached_input_tokens":7,"output_tokens":7,"reasoning_output_tokens":7`)
}
func TestObservedCachedInputTokens(t *testing.T) {
	tokenDimension(t, "cached-input tokens", 22, `"input_tokens":7,"output_tokens":7,"reasoning_output_tokens":7`)
}
func TestObservedOutputTokens(t *testing.T) {
	tokenDimension(t, "output tokens", 33, `"input_tokens":7,"cached_input_tokens":7,"reasoning_output_tokens":7`)
}
func TestObservedReasoningTokens(t *testing.T) {
	tokenDimension(t, "reasoning tokens", 44, `"input_tokens":7,"cached_input_tokens":7,"output_tokens":7`)
}

func TestObservedTurns(t *testing.T) {
	out, _ := observe(t, taskStarted, taskStarted, call("c1"), result("c1", "ok"))
	// A record with calls and no turn record observes no turn rather than inheriting one.
	wantRow(t, out, "turns,2,turns,observed,")
	bare, _ := observe(t, call("c1"), result("c1", "ok"))
	wantRow(t, bare, "turns,0,turns,observed,")
}

func TestObservedReadPaths(t *testing.T) {
	shell := completed("CommandExecution", "i1", `,"command":"cat internal/harnesses/harnesses.go","parsed_cmd":[{"type":"unknown","cmd":"cat internal/harnesses/harnesses.go"}]`)
	out, _ := observe(t, shell, call("c1"), result("c1", "ok"))
	// Shell text is not a read census, so a parsed path would publish a count nobody observed.
	wantRow(t, out, `explicit read paths,"",paths,unknown,`)
	if strings.Contains(out, "harnesses.go") {
		t.Fatalf("record view = %q, want no path parsed out of shell text", out)
	}
}
