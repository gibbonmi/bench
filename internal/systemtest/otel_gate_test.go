//go:build system

package systemtest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/otelrecord"
)

// recordedSpan is the consumer's reading of one record line. It names only the fields a
// reader of the record needs, so the read stays independent of the encoder's structs.
type recordedSpan struct {
	Name       string
	SpanID     string
	Parent     string
	Ended      bool
	Attributes map[string]string
}

// TestOtelGateRecordJourney drives the built binary's `bench gate` against a scaffolded
// repository whose phase manifest declares one green phase, one red phase, and a phase
// the red one skips. It then reads the JSON lines the run left under a private
// BENCH_HOME. The last leg repeats the run with an unwritable record directory.
//
// Rows OT14, OT15, OT16, and OT28.
func TestOtelGateRecordJourney(t *testing.T) {
	home := filepath.Join(t.TempDir(), "bench-home")
	scaffold := scaffoldRecordedGateRepo(t, home, `{"phases":[
		{"name":"probe","argv":["true"]},
		{"name":"blocker","argv":["false"]},
		{"name":"dependent","argv":["true"],"needs":["blocker"]}
	]}`)
	gate := func(benchHome string) processResult {
		return owner.runAt(scaffold.path, scaffold.environment(t, benchHome), "bash", scaffold.wrapper, "gate", "--fresh")
	}

	recorded := gate(home)
	if recorded.code != 1 {
		t.Fatalf("gate with a red phase = (%d, %q, %q)", recorded.code, recorded.stdout, recorded.stderr)
	}
	verdict := gateVerdictLine(t, recorded)
	spans := readSeamRecord(t, home)

	root, present := spans["gate.fresh"]
	if !present {
		t.Fatalf("the record has no root gate span; it holds %v", spanNames(spans))
	}
	// OT14: the root span carries the subject id that groups this subject's iterations,
	// and the run's exit.
	if root.Parent != "" || root.Attributes[otelrecord.AttrSeam] != "gate" {
		t.Fatalf("root span = %+v, want a parentless gate seam", root)
	}
	if len(root.Attributes[otelrecord.AttrSubjectID]) != 40 {
		t.Fatalf("root subject id = %q, want the subject digest", root.Attributes[otelrecord.AttrSubjectID])
	}
	if root.Attributes[otelrecord.AttrOutcome] != "red" {
		t.Fatalf("root outcome = %q, want the run's own exit", root.Attributes[otelrecord.AttrOutcome])
	}

	// OT15: each executed phase carries its own name and its own exit.
	for name, outcome := range map[string]string{"probe": "green", "blocker": "red"} {
		phase, present := spans[name]
		if !present {
			t.Fatalf("the record has no %s phase span; it holds %v", name, spanNames(spans))
		}
		if phase.Attributes[otelrecord.AttrOutcome] != outcome {
			t.Fatalf("%s outcome = %q, want %q", name, phase.Attributes[otelrecord.AttrOutcome], outcome)
		}
		if phase.Parent != root.SpanID {
			t.Fatalf("%s parent = %q, want the root span %q", name, phase.Parent, root.SpanID)
		}
	}

	// OT16: the skipped phase names the need that blocked it.
	skipped, present := spans["dependent"]
	if !present {
		t.Fatalf("the record has no skipped phase span; it holds %v", spanNames(spans))
	}
	if skipped.Attributes[otelrecord.AttrOutcome] != "skipped" || skipped.Attributes[otelrecord.AttrOutcomeBlocker] != "blocker" {
		t.Fatalf("skipped phase span = %+v, want a skip attributed to blocker", skipped)
	}

	// Story 19: no span carries an attribute outside the declared set.
	for name, span := range spans {
		for key := range span.Attributes {
			if !slices.Contains(otelrecord.DeclaredAttributes, key) {
				t.Fatalf("span %s carries the undeclared attribute %q", name, key)
			}
		}
	}

	// OT28: the record is evidence, never a condition. A record directory that cannot be
	// created leaves the exit code and the verdict as they were.
	unwritable := filepath.Join(t.TempDir(), "bench-home")
	if err := os.MkdirAll(unwritable, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(unwritable, "otel"), []byte("not a directory\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	refused := gate(unwritable)
	if refused.code != recorded.code || gateVerdictLine(t, refused) != verdict {
		t.Fatalf("gate with an unwritable record directory = (%d, %q, %q), want the recorded run's exit and verdict %q",
			refused.code, refused.stdout, refused.stderr, verdict)
	}
	if entries, err := filepath.Glob(filepath.Join(unwritable, "otel", "*", "traces.jsonl")); err != nil || len(entries) != 0 {
		t.Fatalf("record files under the unwritable home = %v, %v", entries, err)
	}
}

// gateVerdictLine answers the run's own verdict line. The record must not change it.
func gateVerdictLine(t *testing.T, result processResult) string {
	t.Helper()
	for _, line := range strings.Split(result.stdout+"\n"+result.stderr, "\n") {
		if line == "gate: green" || line == "gate: red" {
			return line
		}
	}
	t.Fatalf("no gate verdict line in (%q, %q)", result.stdout, result.stderr)
	return ""
}

// readSeamRecord answers the finished spans of the one run recorded below home, keyed by
// span name. A start line carries no end time, so it is not the line a consumer reads a
// finished span from.
func readSeamRecord(t *testing.T, home string) map[string]recordedSpan {
	t.Helper()
	spans := map[string]recordedSpan{}
	for _, span := range readRecordLines(t, home, true) {
		if span.Ended {
			spans[span.Name] = span
		}
	}
	return spans
}

// readRecordLines answers every span the one record below home holds, in written order,
// started lines included. strict fails the test on a line that does not parse; a caller
// that reads the record while a run still writes to it passes false, because the last
// line can be a partial write.
func readRecordLines(t *testing.T, home string, strict bool) []recordedSpan {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(home, "otel", "*", "traces.jsonl"))
	if err != nil || len(files) != 1 {
		t.Fatalf("record files under %s = %v, %v; want exactly one", home, files, err)
	}
	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	var read []recordedSpan
	for _, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
		var record struct {
			ResourceSpans []struct {
				ScopeSpans []struct {
					Spans []struct {
						SpanID       string `json:"spanId"`
						ParentSpanID string `json:"parentSpanId"`
						Name         string `json:"name"`
						EndTime      string `json:"endTimeUnixNano"`
						Attributes   []struct {
							Key   string `json:"key"`
							Value struct {
								StringValue string `json:"stringValue"`
							} `json:"value"`
						} `json:"attributes"`
					} `json:"spans"`
				} `json:"scopeSpans"`
			} `json:"resourceSpans"`
		}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			if strict {
				t.Fatalf("record line %q does not parse: %v", line, err)
			}
			continue
		}
		for _, resource := range record.ResourceSpans {
			for _, scope := range resource.ScopeSpans {
				for _, span := range scope.Spans {
					one := recordedSpan{Name: span.Name, SpanID: span.SpanID, Parent: span.ParentSpanID, Ended: span.EndTime != "", Attributes: map[string]string{}}
					for _, attribute := range span.Attributes {
						one.Attributes[attribute.Key] = attribute.Value.StringValue
					}
					read = append(read, one)
				}
			}
		}
	}
	return read
}

// recordedGateRepo is a scaffolded repository whose gate hands off to the phase table the
// way the kit's own gate.sh does. That hand-off is what the run-binary selection and the
// phase spans both key on, so every recorded-gate test scaffolds through here.
type recordedGateRepo struct {
	path    string
	wrapper string
}

// scaffoldRecordedGateRepo answers a repository set up against home whose phase manifest
// is the given JSON.
func scaffoldRecordedGateRepo(t *testing.T, home, manifest string) recordedGateRepo {
	t.Helper()
	repo := filepath.Join(t.TempDir(), "repository")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if result := owner.runAt(repo, nil, "git", "init", "-q", "-b", "main"); result.code != 0 {
		t.Fatalf("git init = (%d, %q)", result.code, result.stderr)
	}
	scaffold := recordedGateRepo{path: repo, wrapper: filepath.Join(repo, ".bench", "bin", "bench.sh")}
	if setup := owner.runAt(repo, scaffold.environment(t, home), owner.selected.path, "setup", "--yes"); setup.code != 3 {
		t.Fatalf("bench setup --yes = (%d, %q, %q)", setup.code, setup.stdout, setup.stderr)
	}
	// The scaffolded gate is a fail-closed stub that runs no phase table. This gate hands
	// off to the phase table itself.
	phaseGate := "#!/usr/bin/env bash\n" +
		"set -uo pipefail\n" +
		"root=\"$(git rev-parse --show-toplevel)\"\n" +
		"cd \"$root\"\n" +
		"bench=\"$(dirname \"$0\")/bin/bench.sh\"\n" +
		"\"$bench\" gate-phases \"$root\"\n" +
		"status=$?\n" +
		"if [ \"$status\" -eq 0 ]; then echo \"gate: green\"; else echo \"gate: red\" >&2; fi\n" +
		"exit \"$status\"\n"
	if err := os.WriteFile(filepath.Join(repo, ".bench", "gate.sh"), []byte(phaseGate), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".bench", "phases.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return scaffold
}

// environment answers the overrides a run against benchHome needs. Each call observes the
// selected executable, so every run through the scaffold joins the owner's ledger.
func (r recordedGateRepo) environment(t *testing.T, benchHome string) []string {
	t.Helper()
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	return []string{"BENCH_HOME=" + benchHome, "BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}
}

func spanNames(spans map[string]recordedSpan) []string {
	names := make([]string, 0, len(spans))
	for name := range spans {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}
