package gate

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/benchhome"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/otelrecord"
)

// OG25: the record names the composed tree, the lane identity, and the outcome. OG16 and
// OG17: the run leaves the gate cache and the evidence store untouched.
func TestRunLaneRecordsItsOwnFileOnly(t *testing.T) {
	for _, tc := range []struct {
		name    string
		check   Phase
		outcome string
		failing string
	}{
		{name: "pass", check: Phase{Name: "unit", Argv: []string{"true"}}, outcome: "pass"},
		{
			name:    "fail",
			check:   Phase{Name: "unit", Argv: []string{"sh", "-c", "echo first line >&2; echo second >&2; exit 3"}},
			outcome: "fail",
			failing: "unit",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := outcomeFixture(t)
			tree := outcomeGit(t, root, "rev-parse", "HEAD^{tree}")
			gitdir := outcomeGit(t, root, "rev-parse", "--absolute-git-dir")

			// PL24: the run selects, and the record's schema is unchanged by it. The
			// changed path is in no class, so the unknown class selects the one check.
			result, err := RunLane(context.Background(), LaneRequest{
				Root: root, Tree: tree, Lane: "benchkit", Checks: []Phase{tc.check},
				Selective: true, Changes: []ComposedChange{laneChange("bin/x.sh")},
			})
			if err != nil {
				t.Fatal(err)
			}
			if result.Outcome != tc.outcome || result.Check != tc.failing {
				t.Fatalf("result = {%s %s}, want {%s %s}", result.Outcome, result.Check, tc.outcome, tc.failing)
			}
			if tc.outcome == "fail" && result.Diagnostic != "first line" {
				t.Errorf("diagnostic = %q, want the check's first output line", result.Diagnostic)
			}

			record := readLaneRecord(t, filepath.Join(gitdir, laneRecordFile))
			if record["tree"] != tree || record["lane"] != "benchkit" || record["outcome"] != tc.outcome {
				t.Errorf("record = %v, want the composed tree, the lane, and outcome %s", record, tc.outcome)
			}
			if digest, _ := record["run_binary"].(string); len(digest) != 64 {
				t.Errorf("record run_binary = %q, want a content address", digest)
			}

			// OG16 and OG17: a lane grades a check list, so it authorizes nothing the
			// gate decides and leaves no evidence a landing could reuse.
			if _, err := os.Stat(filepath.Join(gitdir, benchgit.GateCacheFile)); !os.IsNotExist(err) {
				t.Errorf("the lane run wrote a gate cache record (stat err %v)", err)
			}
			common, err := benchgit.CommonDir(root)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(filepath.Join(common, "bench-gate-evidence")); !os.IsNotExist(err) {
				t.Errorf("the lane run wrote an evidence record (stat err %v)", err)
			}
		})
	}
}

// A lane whose check outlives the gate timeout returns an error and writes no lane
// record: a record of a partial run would carry an outcome no check ever reached.
func TestRunLaneTimeoutWritesNoRecord(t *testing.T) {
	root := outcomeFixture(t)
	tree := outcomeGit(t, root, "rev-parse", "HEAD^{tree}")
	gitdir := outcomeGit(t, root, "rev-parse", "--absolute-git-dir")

	previousTimeout := gateTimeout
	gateTimeout = 100 * time.Millisecond
	t.Cleanup(func() { gateTimeout = previousTimeout })

	_, err := RunLane(context.Background(), LaneRequest{
		Root: root, Tree: tree, Lane: "benchkit",
		Checks: []Phase{{Name: "slow", Argv: []string{"sleep", "5"}}},
	})
	if err == nil {
		t.Fatal("RunLane returned no error on a check that outlived the gate timeout")
	}
	if _, statErr := os.Stat(filepath.Join(gitdir, laneRecordFile)); !os.IsNotExist(statErr) {
		t.Fatalf("the timed-out lane wrote a record (stat err %v)", statErr)
	}
}

// OG19: Inspect names the lane class and never calls a lane pass reusable green.
func TestInspectRefusesALaneRecord(t *testing.T) {
	root := outcomeFixture(t)
	subject, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	writeInspectRecord(t, root, inspectJSON(verdictRecord{
		Schema:     verdictSchema,
		Tree:       subject.Tree,
		Lane:       "benchkit",
		Outcome:    "pass",
		RunBinary:  strings.Repeat("a", 64),
		RecordedAt: time.Now().UTC().Add(-time.Minute).Truncate(time.Second).Format(time.RFC3339),
	}))
	got := Inspect(root)
	if got.State != Ready {
		t.Fatalf("Inspect state = %s, want %s", got.State, Ready)
	}
	if got.ReusableGreen {
		t.Error("a lane pass answered ReusableGreen, which would authorize a landing")
	}
	if got.Reason != "lane record" {
		t.Errorf("reason = %q, want the lane class named", got.Reason)
	}
}

func readLaneRecord(t *testing.T, path string) map[string]any {
	t.Helper()
	var record map[string]any
	if err := json.Unmarshal(outcomeRead(t, path), &record); err != nil {
		t.Fatal(err)
	}
	return record
}

func writeLaneFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// OT17: a red lane's span carries the first failing check and that check's first
// diagnostic line, and a green lane's span names no failing check. The tripwire reads
// the lane's red from this one pair, so a phase-and-exit-only record reds here.
func TestRunLaneSpanCarriesTheFailingCheck(t *testing.T) {
	for _, tc := range []struct {
		name       string
		check      Phase
		outcome    string
		failing    string
		diagnostic string
	}{
		{name: "pass", check: Phase{Name: "unit", Argv: []string{"true"}}, outcome: lanePass},
		{
			name:       "fail",
			check:      Phase{Name: "unit", Argv: []string{"sh", "-c", "echo first line >&2; echo second >&2; exit 3"}},
			outcome:    laneFail,
			failing:    "unit",
			diagnostic: "first line",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv(benchhome.Env, home)
			root := outcomeFixture(t)
			tree := outcomeGit(t, root, "rev-parse", "HEAD^{tree}")

			if _, err := RunLane(context.Background(), LaneRequest{
				Root: root, Tree: tree, Lane: "benchkit", Checks: []Phase{tc.check},
			}); err != nil {
				t.Fatal(err)
			}

			attributes := laneSpanAttributes(t, home, root)
			if attributes[otelrecord.AttrSubjectID] != tree {
				t.Errorf("span subject = %q, want the composed tree %q", attributes[otelrecord.AttrSubjectID], tree)
			}
			if attributes[otelrecord.AttrOutcome] != tc.outcome {
				t.Errorf("span outcome = %q, want %q", attributes[otelrecord.AttrOutcome], tc.outcome)
			}
			if attributes[otelrecord.AttrOutcomeCheck] != tc.failing {
				t.Errorf("span failing check = %q, want %q", attributes[otelrecord.AttrOutcomeCheck], tc.failing)
			}
			if attributes[otelrecord.AttrOutcomeDiagnostic] != tc.diagnostic {
				t.Errorf("span diagnostic = %q, want %q", attributes[otelrecord.AttrOutcomeDiagnostic], tc.diagnostic)
			}
			// No diagnostic payload beyond the first line enters the record.
			if record := outcomeRead(t, otelrecord.Path(home, root)); bytes.Contains(record, []byte("second")) {
				t.Error("the record carries diagnostic output beyond the failing check's first line")
			}
		})
	}
}

// laneSpanAttributes reads the lane seam's ended span back out of the record. The start
// line carries no outcome, so the reader keeps the line the seam's end wrote.
func laneSpanAttributes(t *testing.T, home, root string) map[string]string {
	t.Helper()
	var found map[string]string
	for _, line := range strings.Split(strings.TrimSpace(string(outcomeRead(t, otelrecord.Path(home, root)))), "\n") {
		var parsed struct {
			ResourceSpans []struct {
				ScopeSpans []struct {
					Spans []struct {
						Name            string `json:"name"`
						EndTimeUnixNano string `json:"endTimeUnixNano"`
						Attributes      []struct {
							Key   string `json:"key"`
							Value struct {
								StringValue string `json:"stringValue"`
							} `json:"value"`
						} `json:"attributes"`
					} `json:"spans"`
				} `json:"scopeSpans"`
			} `json:"resourceSpans"`
		}
		if err := json.Unmarshal([]byte(line), &parsed); err != nil {
			t.Fatalf("record line does not parse: %v: %s", err, line)
		}
		for _, resource := range parsed.ResourceSpans {
			for _, scope := range resource.ScopeSpans {
				for _, span := range scope.Spans {
					attributes := map[string]string{}
					for _, pair := range span.Attributes {
						attributes[pair.Key] = pair.Value.StringValue
					}
					if attributes[otelrecord.AttrSeam] != "lane" || span.EndTimeUnixNano == "" {
						continue
					}
					if found != nil {
						t.Fatalf("the record holds a second ended lane span %q", span.Name)
					}
					found = attributes
				}
			}
		}
	}
	if found == nil {
		t.Fatal("the record holds no ended lane span")
	}
	return found
}
