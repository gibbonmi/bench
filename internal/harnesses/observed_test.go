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

func usageSnapshot(input, cached, output, reasoning int) string {
	return fmt.Sprintf(`{"timestamp":"2026-09-11T09:00:02.000Z","type":"event_msg","payload":{"type":"token_count","info":{"model_context_window":400,"total_token_usage":{"input_tokens":%d,"cached_input_tokens":%d,"output_tokens":%d,"reasoning_output_tokens":%d}}}}`, input, cached, output, reasoning)
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
// over it. The command reads the real file, so the path boundary is exercised rather than
// stubbed.
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

// wantRow asserts one observation row by its leading cells. Trailing cells carry commas and
// are quoted, so a prefix assertion names exactly the cells under test.
func wantRow(t *testing.T, out, prefix string) {
	t.Helper()
	if !strings.Contains(out, "\n  "+prefix) {
		t.Fatalf("record view = %q, want an observation row starting %q", out, prefix)
	}
}

func TestObservedTextBoundary(t *testing.T) {
	out, code := observe(t, call("c1"), result("c1", multibyte), reasoning)
	if code != 0 {
		t.Fatalf("record view exit = %d, want 0", code)
	}
	// The fixture holds 13 UTF-8 bytes across 11 characters on 2 lines. The reasoning
	// record's encrypted payload is serialized metadata, so none of its bytes count.
	wantRow(t, out, "result-text bytes,13,bytes,observed,")
	wantRow(t, out, "result-text characters,11,characters,observed,")
	wantRow(t, out, "result-text lines,2,lines,observed,")
}

func TestObservedCallAncestry(t *testing.T) {
	subagent := completed("SubAgentActivity", "i1", `,"agent_thread_id":"sub-1","kind":"started"`)
	out, _ := observe(t, call("c1"), result("c1", "ok"), subagent)
	// A subagent event names a thread, not a call. It can neither raise the outer census
	// nor supply a nested count this source never records.
	wantRow(t, out, "outer calls,1,calls,observed,")
	wantRow(t, out, `observed nested calls,"",calls,unknown,`)
}

func TestObservedDuplicateCompletion(t *testing.T) {
	out, _ := observe(t, call("c1"), result("c1", "ok"), result("c1", "ok"))
	// Both completions carry one invocation identity, so the call counts once and its
	// bytes count once.
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
	// have raised reports the partial total as incomplete rather than as a complete one.
	wantRow(t, out, "result-text bytes,2,bytes,incomplete,")
	wantRow(t, out, "outer calls,1,calls,incomplete,")
	wantRow(t, out, "turns,1,turns,incomplete,")
}

func TestObservedNativeUsage(t *testing.T) {
	out, _ := observe(t, usageSnapshot(10, 2, 3, 4), usageSnapshot(30, 6, 9, 12), usageSnapshot(60, 12, 18, 24))
	// The snapshots are cumulative session totals. The last one is the session total; a sum
	// of the three would report 100 input tokens that nobody spent.
	wantRow(t, out, "input tokens,60,tokens,observed,")
	wantRow(t, out, "output tokens,18,tokens,observed,")
}

func TestObservedNoTokenAttribution(t *testing.T) {
	out, _ := observe(t, call("c1"), result("c1", multibyte), usageSnapshot(60, 12, 18, 24))
	// Bytes and session totals are both present, and neither divides into a per-result
	// cost. An estimate here would render as an observed measure.
	wantRow(t, out, `per-result token attribution,"",tokens,unknown,`)
}

func TestObservedCompactions(t *testing.T) {
	out, _ := observe(t, compacted, compacted, usageSnapshot(60, 12, 18, 24))
	// Only an identified compaction record counts. The snapshot names a small context
	// window and is not evidence that the window was compacted.
	wantRow(t, out, "compactions,2,compactions,observed,")
	bare, _ := observe(t, usageSnapshot(60, 12, 18, 24))
	wantRow(t, bare, "compactions,0,compactions,observed,")
}

func TestObservedPreservesCompiledViews(t *testing.T) {
	views := strings.Split(compiledBaseline, "--- ")
	if len(views) != len(Rows)+1 {
		t.Fatalf("baseline holds %d views, want %d", len(views), len(Rows)+1)
	}
	if out, code := Command(nil); code != 0 || out != views[0] {
		t.Fatalf("bench harnesses = %q exit %d, want the captured overview %q", out, code, views[0])
	}
	for i, row := range Rows {
		want := strings.SplitN(views[i+1], "\n", 2)
		if want[0] != row.Harness {
			t.Fatalf("baseline view %d names %q, want %q", i+1, want[0], row.Harness)
		}
		if out, code := Command([]string{row.Harness}); code != 0 || out != want[1] {
			t.Fatalf("bench harnesses %s = %q exit %d, want the captured view %q", row.Harness, out, code, want[1])
		}
	}
}

func TestObservedMeasureProvenance(t *testing.T) {
	out, _ := observe(t, call("c1"), result("c1", multibyte))
	_, rest, found := strings.Cut(out, "observations[14]{metric,value,unit,availability,source,boundary}:\n")
	block, _, ok := strings.Cut(rest, "help[0]")
	if !found || !ok {
		t.Fatalf("record view = %q, want a 14-row observations block", out)
	}
	for _, line := range strings.Split(strings.TrimSuffix(block, "\n"), "\n") {
		// The first four cells are the metric, its value, its unit, and its availability.
		// What follows is the provenance pair an unknown row must carry too.
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
	out, _ := observe(t, call("c1"), call("c2"), result("c2", "ok"), result("c3", "ok"))
	// c1 never completed and c3's completion names no call, so two invocations lack a
	// matched pair.
	wantRow(t, out, "unmatched calls,2,calls,observed,")
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
	// A FIFO would block the reader in open(2) and a link would supply bytes from outside
	// the named path, so both are refused on the type check that precedes every open.
	for _, path := range []string{link, fifo} {
		out, code := observeAt(t, path, harnesstranscript.SourceCodexRollout)
		if code != 1 || !strings.HasPrefix(out, "error: "+path+" is wrong-type") {
			t.Fatalf("record view %s = %q exit %d, want a wrong-type refusal at exit 1", path, out, code)
		}
	}
}

// tokenDimension grades one native counter across a present value, an observed zero, and an
// absent counter. Each row owns its own call, so a dropped dimension names itself.
func tokenDimension(t *testing.T, metric string, present int) {
	t.Helper()
	out, _ := observe(t, usageSnapshot(11, 22, 33, 44))
	wantRow(t, out, fmt.Sprintf("%s,%d,tokens,observed,", metric, present))
	zero, _ := observe(t, usageSnapshot(0, 0, 0, 0))
	wantRow(t, zero, metric+",0,tokens,observed,")
	absent, _ := observe(t, call("c1"), result("c1", "ok"))
	wantRow(t, absent, metric+`,"",tokens,unknown,`)
}

func TestObservedInputTokens(t *testing.T) { tokenDimension(t, "input tokens", 11) }

func TestObservedCachedInputTokens(t *testing.T) { tokenDimension(t, "cached-input tokens", 22) }

func TestObservedOutputTokens(t *testing.T) { tokenDimension(t, "output tokens", 33) }

func TestObservedReasoningTokens(t *testing.T) { tokenDimension(t, "reasoning tokens", 44) }

func TestObservedTurns(t *testing.T) {
	out, _ := observe(t, taskStarted, taskStarted, call("c1"), result("c1", "ok"))
	// A turn is one task_started record. A record with calls and no turn record observes
	// no turn rather than inheriting one.
	wantRow(t, out, "turns,2,turns,observed,")
	bare, _ := observe(t, call("c1"), result("c1", "ok"))
	wantRow(t, bare, "turns,0,turns,observed,")
}

func TestObservedReadPaths(t *testing.T) {
	shell := completed("CommandExecution", "i1", `,"command":"cat internal/harnesses/harnesses.go","parsed_cmd":[{"type":"unknown","cmd":"cat internal/harnesses/harnesses.go"}]`)
	out, _ := observe(t, shell, call("c1"), result("c1", "ok"))
	// The source records a read only inside shell text, and shell text is not a read
	// census. A parsed path here would publish a count nobody observed.
	wantRow(t, out, `explicit read paths,"",paths,unknown,`)
	if strings.Contains(out, "harnesses.go") {
		t.Fatalf("record view = %q, want no path parsed out of shell text", out)
	}
}

// compiledBaseline is the two compiled views captured from the tree before the record view
// existed, one view per `--- <harness>` section after the overview. The test compares a live
// run against these bytes, so an opt-in observation that alters a compiled view turns red.
const compiledBaseline = `schema: 1
harnesses[4]{harness,provider,phase_form,hooks,delegation_guard,headless,checked}:
  codex,openai,$bench-,.codex/hooks.json,no,.bench/adapters/codex,2026-07-11
  claude,anthropic,/bench-,.claude/settings.json,yes,.bench/adapters/claude,2026-08-26
  opencode,any,"","",unknown,.bench/adapters/opencode,""
  none,none,"","",no,"",2026-08-26
help[0]{cmd,why}:
--- codex
schema: 1
cells[13]{field,value,source,checked}:
  steering during an active turn,unknown,"",""
  structured user questions,unknown,"",""
  tool-permission controls,unknown,"",""
  hooks,yes,.codex/hooks.json,2026-08-26
  MCP support,unknown,"",""
  subagent support,unknown,"",""
  subagent isolation,unknown,"",""
  effort selection,unknown,"",""
  persistent tasks,unknown,"",""
  resume and recovery,unknown,"",""
  structured output and exit status,unknown,"",""
  headless execution,yes,.bench/adapters/codex,2026-08-26
  delegation_guard,no,".bench/BENCH-reference.md Hook Layers, the agent-line bullet (Codex hooks docs)",2026-07-11
measures[4]{measure,value,supplier}:
  tokens,unknown,FT204 harness transcript reader
  tool calls,unknown,FT204 harness transcript reader
  Read paths,unknown,FT204 harness transcript reader
  turns,unknown,FT204 harness transcript reader
help[0]{cmd,why}:
--- claude
schema: 1
cells[13]{field,value,source,checked}:
  steering during an active turn,unknown,"",""
  structured user questions,unknown,"",""
  tool-permission controls,unknown,"",""
  hooks,yes,.claude/settings.json,2026-08-26
  MCP support,unknown,"",""
  subagent support,unknown,"",""
  subagent isolation,unknown,"",""
  effort selection,unknown,"",""
  persistent tasks,unknown,"",""
  resume and recovery,unknown,"",""
  structured output and exit status,unknown,"",""
  headless execution,yes,.bench/adapters/claude,2026-08-26
  delegation_guard,yes,.claude/settings.json PreToolUse Agent matcher runs .bench/hooks/check-agent-line.sh,2026-08-26
measures[4]{measure,value,supplier}:
  tokens,unknown,FT204 harness transcript reader
  tool calls,unknown,FT204 harness transcript reader
  Read paths,unknown,FT204 harness transcript reader
  turns,unknown,FT204 harness transcript reader
help[0]{cmd,why}:
--- opencode
schema: 1
cells[13]{field,value,source,checked}:
  steering during an active turn,unknown,"",""
  structured user questions,unknown,"",""
  tool-permission controls,unknown,"",""
  hooks,unknown,"",""
  MCP support,unknown,"",""
  subagent support,unknown,"",""
  subagent isolation,unknown,"",""
  effort selection,unknown,"",""
  persistent tasks,unknown,"",""
  resume and recovery,unknown,"",""
  structured output and exit status,unknown,"",""
  headless execution,yes,.bench/adapters/opencode,2026-08-26
  delegation_guard,unknown,"",""
measures[4]{measure,value,supplier}:
  tokens,unknown,FT204 harness transcript reader
  tool calls,unknown,FT204 harness transcript reader
  Read paths,unknown,FT204 harness transcript reader
  turns,unknown,FT204 harness transcript reader
help[0]{cmd,why}:
--- none
schema: 1
cells[13]{field,value,source,checked}:
  steering during an active turn,unknown,"",""
  structured user questions,unknown,"",""
  tool-permission controls,unknown,"",""
  hooks,no,".bench/adapters/ names no none entry, and no config names none",2026-08-26
  MCP support,unknown,"",""
  subagent support,unknown,"",""
  subagent isolation,unknown,"",""
  effort selection,unknown,"",""
  persistent tasks,unknown,"",""
  resume and recovery,unknown,"",""
  structured output and exit status,unknown,"",""
  headless execution,no,.bench/adapters/ names no none entry,2026-08-26
  delegation_guard,no,".bench/adapters/ names no none entry, so the model-free path runs no agent",2026-08-26
measures[4]{measure,value,supplier}:
  tokens,unknown,FT204 harness transcript reader
  tool calls,unknown,FT204 harness transcript reader
  Read paths,unknown,FT204 harness transcript reader
  turns,unknown,FT204 harness transcript reader
help[0]{cmd,why}:
`
