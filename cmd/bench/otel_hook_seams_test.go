package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/otelrecord"
)

// TestHookSeamsMatchTheHookDispatchRows covers OT21 at the registry. Every hook row
// records through one symbol, so the conformance check cannot see a dropped Hook flag or
// a registry row that names no verb. The two lists are authored independently — the
// dispatch table declares the hooks, and the registry declares the recorded seams — so
// set equality is the only reconciliation between them.
func TestHookSeamsMatchTheHookDispatchRows(t *testing.T) {
	var dispatched []string
	for _, definition := range commandRegistry {
		if definition.Hook {
			dispatched = append(dispatched, otelHookSeamPrefix+definition.Name)
		}
	}
	var registered []string
	for _, entry := range otelrecord.Registry {
		if strings.HasPrefix(entry.Seam, otelHookSeamPrefix) {
			registered = append(registered, entry.Seam)
		}
	}
	sort.Strings(dispatched)
	sort.Strings(registered)
	if len(dispatched) == 0 {
		t.Fatal("the dispatch table declares no hook row")
	}
	if strings.Join(dispatched, ",") != strings.Join(registered, ",") {
		t.Fatalf("hook seams: dispatch %v, registry %v", dispatched, registered)
	}
}

// routedRepo makes a repository whose .bench/lines.env binds every tier of the codex and
// claude columns, and moves the test into it.
func routedRepo(t *testing.T) {
	t.Helper()
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	binding := "BENCH_CODEX_TOP=gpt-5.6-sol\nBENCH_CODEX_MID=gpt-5.6-terra\nBENCH_CODEX_CHEAP=gpt-5.6-luna\n" +
		"BENCH_CLAUDE_TOP=fable\nBENCH_CLAUDE_MID=opus\nBENCH_CLAUDE_CHEAP=sonnet\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "lines.env"), []byte(binding), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
}

// recordedResolve runs resolve-model for harness in a routed repository with a private
// home, and returns its stdout, its exit, the one line.resolve span it recorded, and the
// raw record. handoff, when true, runs it under a handed-off parent span, which it also
// returns.
func recordedResolve(t *testing.T, harness, tier string, handoff bool) (string, int, otelrecord.Span, otelrecord.Span, []byte) {
	t.Helper()
	routedRepo(t)
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	t.Setenv("BENCH_MODEL", tier)
	if tier == "" {
		os.Unsetenv("BENCH_MODEL")
	}
	for _, name := range otelrecord.HandoffVariables() {
		t.Setenv(name, "")
	}
	root, err := git.Root()
	if err != nil {
		t.Fatal(err)
	}
	var parent otelrecord.Span
	end := func() {}
	if handoff {
		ctx, span, finish := otelrecord.BeginIn(context.Background(), home, root, "test.parent", "test.parent")
		for _, pair := range otelrecord.WithHandoff(ctx, root, nil) {
			name, value, _ := strings.Cut(pair, "=")
			t.Setenv(name, value)
		}
		parent.TraceID, parent.SpanID = span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String()
		end = finish
	}
	out, code := resolveModel([]string{"--harness", harness})
	end()
	spans, err := otelrecord.ReadSpans(home, root)
	if err != nil {
		t.Fatal(err)
	}
	var found []otelrecord.Span
	for _, span := range spans {
		if span.Seam == otelResolveSeam {
			found = append(found, span)
		}
	}
	if len(found) != 1 {
		t.Fatalf("the record holds %d line.resolve spans, want 1: %+v", len(found), spans)
	}
	raw, err := os.ReadFile(otelrecord.Path(home, root))
	if err != nil {
		t.Fatal(err)
	}
	return out, code, found[0], parent, raw
}

// LE46: a handed-off resolution records its harness, tier, and model under the handoff.
func TestAHandedOffResolutionRecordsItsLine(t *testing.T) {
	out, code, span, parent, _ := recordedResolve(t, "claude", "mid", true)
	if out != "opus\n" || code != 0 {
		t.Fatalf("resolveModel = (%q, %d), want (\"opus\\n\", 0)", out, code)
	}
	if span.TraceID != parent.TraceID || span.ParentSpanID != parent.SpanID {
		t.Fatalf("line.resolve parent = %s/%s, want %s/%s", span.TraceID, span.ParentSpanID, parent.TraceID, parent.SpanID)
	}
	for key, want := range map[string]string{
		otelrecord.AttrLineHarness: "claude",
		otelrecord.AttrLineTier:    "mid",
		otelrecord.AttrLineModel:   "opus",
		otelrecord.AttrOutcome:     otelrecord.OutcomeGreen,
	} {
		if got := span.Attributes[key]; got != want {
			t.Fatalf("line.resolve %s = %q, want %q", key, got, want)
		}
	}
}

// LE47: a resolution with no handoff records a span of its own trace.
func TestAStandaloneResolutionRecordsARootSpan(t *testing.T) {
	out, code, span, _, _ := recordedResolve(t, "claude", "cheap", false)
	if out != "sonnet\n" || code != 0 {
		t.Fatalf("resolveModel = (%q, %d), want (\"sonnet\\n\", 0)", out, code)
	}
	if span.ParentSpanID != "" {
		t.Fatalf("line.resolve parent = %q, want none", span.ParentSpanID)
	}
}

// LE48: a refused resolution records a red outcome and no model.
func TestARefusedResolutionRecordsRed(t *testing.T) {
	out, code, span, _, _ := recordedResolve(t, "claude", "", false)
	if out != "" || code != 1 {
		t.Fatalf("resolveModel = (%q, %d), want (\"\", 1)", out, code)
	}
	if got := span.Attributes[otelrecord.AttrOutcome]; got != otelrecord.OutcomeRed {
		t.Fatalf("line.resolve outcome = %q, want red", got)
	}
	if got, ok := span.Attributes[otelrecord.AttrLineModel]; ok {
		t.Fatalf("line.resolve carries model %q, want none", got)
	}
}

// An unknown harness and an operator tier never reach the record.
func TestAResolutionRecordsNoOperatorText(t *testing.T) {
	_, code, span, _, raw := recordedResolve(t, "gemini", "custom model", false)
	if code != 1 {
		t.Fatalf("resolveModel exit = %d, want 1", code)
	}
	for _, key := range []string{otelrecord.AttrLineHarness, otelrecord.AttrLineTier} {
		if got, ok := span.Attributes[key]; ok {
			t.Fatalf("line.resolve carries %s = %q, want none", key, got)
		}
	}
	if strings.Contains(string(raw), "gemini") || strings.Contains(string(raw), "custom model") {
		t.Fatalf("the record holds operator text:\n%s", raw)
	}
}
