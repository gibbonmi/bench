package harnesses

import (
	"strings"
	"testing"
)

// The expectations here are authored independently of the projection. Deriving a header or
// a row from overviewFields, Mechanics, or Rows would let a dropped field, a dropped
// mechanic, or a dropped harness disappear from both the code and its check at once.

func TestCommandProjectsEveryHarnessRow(t *testing.T) {
	out, code := Command(nil)
	if code != 0 {
		t.Fatalf("bench harnesses exit = %d, want 0", code)
	}
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("bench harnesses = %q, want a schema line and a table", out)
	}
	if lines[0] != "schema: 1" {
		t.Fatalf("first line = %q, want %q", lines[0], "schema: 1")
	}
	const wantHeader = "harnesses[4]{harness,provider,phase_form,hooks,delegation_guard,headless,checked}:"
	if lines[1] != wantHeader {
		t.Fatalf("table header = %q, want %q", lines[1], wantHeader)
	}
	// The `none` row is named explicitly: a projection that skips the model-free harness
	// hides the degraded path a reader most needs to see.
	for _, want := range []string{
		"  codex,openai,$bench-,.codex/hooks.json,no,.bench/adapters/codex,2026-07-11",
		"  claude,anthropic,/bench-,.claude/settings.json,yes,.bench/adapters/claude,2026-08-26",
		`  none,none,"","",no,"",2026-08-26`,
	} {
		if !strings.Contains(out, want+"\n") {
			t.Fatalf("bench harnesses = %q, want row %q", out, want)
		}
	}
	if !strings.HasSuffix(out, "help[0]{cmd,why}:\n") {
		t.Fatalf("bench harnesses = %q, want a terminal empty help envelope", out)
	}
}

func TestCommandProjectsOneHarnessCellsWithSources(t *testing.T) {
	out, code := Command([]string{"codex"})
	if code != 0 {
		t.Fatalf("bench harnesses codex exit = %d, want 0", code)
	}
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) < 2 || lines[0] != "schema: 1" {
		t.Fatalf("bench harnesses codex = %q, want a leading schema line", out)
	}
	const wantHeader = "cells[13]{field,value,source,checked}:"
	if lines[1] != wantHeader {
		t.Fatalf("table header = %q, want %q", lines[1], wantHeader)
	}
	// The delegation-guard cell names the upstream Codex hooks docs, so the verdict cites
	// a read rather than restating the prose.
	wantGuard := `  delegation_guard,no,".bench/BENCH-reference.md Hook Layers, the agent-line bullet (Codex hooks docs)",2026-07-11`
	if !strings.Contains(out, wantGuard) {
		t.Fatalf("bench harnesses codex = %q, want the delegation_guard cell %q", out, wantGuard)
	}
	// An unknown cell carries neither a source nor a date.
	if !strings.Contains(out, "  structured output and exit status,unknown,\"\",\"\"\n") {
		t.Fatalf("bench harnesses codex = %q, want an ungraded cell with an empty source and date", out)
	}
	// The four measure cells render as their own table. Each names its supplier and reads
	// unknown, because no supplier ships yet.
	const wantMeasureHeader = "measures[4]{measure,value,supplier}:"
	if !strings.Contains(out, wantMeasureHeader+"\n") {
		t.Fatalf("bench harnesses codex = %q, want the measures header %q", out, wantMeasureHeader)
	}
	for _, want := range []string{
		"  tokens,unknown,FT204 harness transcript reader",
		`  tool calls,unknown,FT204 harness transcript reader`,
		`  Read paths,unknown,FT204 harness transcript reader`,
		"  turns,unknown,FT204 harness transcript reader",
	} {
		if !strings.Contains(out, want+"\n") {
			t.Fatalf("bench harnesses codex = %q, want measure row %q", out, want)
		}
	}
	if !strings.HasSuffix(out, "help[0]{cmd,why}:\n") {
		t.Fatalf("bench harnesses codex = %q, want a terminal empty help envelope", out)
	}
}

func TestCommandRefusesAnUnknownHarness(t *testing.T) {
	out, code := Command([]string{"cursor"})
	if code != 2 {
		t.Fatalf("bench harnesses cursor exit = %d, want 2", code)
	}
	if out != "usage: bench harnesses (unknown argument: cursor)\n" {
		t.Fatalf("bench harnesses cursor = %q, want the usage line", out)
	}
}

func TestCommandRefusesTwoPositionals(t *testing.T) {
	out, code := Command([]string{"codex", "claude"})
	if code != 2 {
		t.Fatalf("bench harnesses codex claude exit = %d, want 2", code)
	}
	if out != "usage: bench harnesses (unknown argument: claude)\n" {
		t.Fatalf("bench harnesses codex claude = %q, want the usage line", out)
	}
}

func TestCommandAnswersEveryHelpSpelling(t *testing.T) {
	for _, spelling := range []string{"--help", "-h", "help"} {
		out, code := Command([]string{spelling})
		if code != 0 || out != "usage: bench harnesses [<harness>]\n" {
			t.Fatalf("bench harnesses %s = %q exit %d, want the usage line and exit 0", spelling, out, code)
		}
	}
}
