package reviewrecord_test

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// reviewJSONBaseline holds the stored results of Parse and ReadPlan for the hostile
// document family. An ordinary run reads it and never writes it.
const reviewJSONBaseline = "testdata/strict-json-baseline.json"

type reviewJSONResults struct {
	Record map[string]reviewResult `json:"record"`
	Plan   map[string]reviewResult `json:"plan"`
}

// reviewResult is the version of an accepted document or the refusal of a refused one.
type reviewResult struct {
	Version int    `json:"version,omitempty"`
	Error   string `json:"error,omitempty"`
}

// validRecord is the smallest record that Parse accepts.
func validRecord() rr.Record {
	return rr.Record{Version: 1, Spec: recordtest.Spec, PlanDigest: "sha256:plan", ImplementationSession: "author"}
}

func TestReviewJSONCompatibility(t *testing.T) {
	record := marshalReview(t, validRecord())

	// The specific contracts run before the baseline family, so a framing, fence, or bounds
	// drift fails under its own name first.
	t.Run("unframed review document", func(t *testing.T) {
		if strings.HasSuffix(record, "\n") {
			t.Fatalf("seed %q ends in a newline", record)
		}
		if got := parseResult([]byte(record)); got.Error != "" {
			t.Fatalf("Parse of an unframed record = %+v, want it accepted", got)
		}
	})
	t.Run("fenced completion plan", func(t *testing.T) {
		f := recordtest.New(t, 1)
		if got := planResult(f); got.Error != "" {
			t.Fatalf("ReadPlan of the rendered plan fence = %+v, want it accepted", got)
		}
	})
	limits := recordLimitCases(record)
	t.Run("byte limits", func(t *testing.T) {
		for name, want := range map[string]string{
			"oversize":      "invalid record bytes: unreadable",
			"invalid UTF-8": "invalid record bytes: malformed invalid UTF-8",
		} {
			if got := parseResult([]byte(limits[name])); !strings.HasPrefix(got.Error, want) {
				t.Errorf("%s: Parse = %+v, want a refusal that starts %q", name, got, want)
			}
		}
	})

	stored := readReviewJSONBaseline(t)
	current := currentReviewResults(t, record, limits)
	t.Run("record baseline", func(t *testing.T) { compareReviewJSON(t, current.Record, stored.Record) })
	t.Run("plan baseline", func(t *testing.T) { compareReviewJSON(t, current.Plan, stored.Plan) })
}

// currentReviewResults reads every case of the baseline family through Parse and ReadPlan.
// The plan cases share one fixture and commit each payload in name order.
func currentReviewResults(t *testing.T, record string, limits map[string]string) reviewJSONResults {
	t.Helper()
	current := reviewJSONResults{Record: map[string]reviewResult{}, Plan: map[string]reviewResult{}}
	records := reviewFamily(record)
	for name, body := range limits {
		records[name] = body
	}
	for name, value := range recordDomainCases() {
		records[name] = marshalReview(t, value)
	}
	for name, body := range records {
		current.Record[name] = parseResult([]byte(body))
	}

	f := recordtest.New(t, 1)
	plans := reviewFamily(marshalReview(t, f.Plan))
	unsupported := f.Plan
	unsupported.Version = 3
	plans["unsupported version"] = marshalReview(t, unsupported)
	names := make([]string, 0, len(plans))
	for name := range plans {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		f.WritePlanPayload([]byte(plans[name]))
		f.Commit("plan case " + name)
		current.Plan[name] = planResult(f)
	}

	return current
}

// reviewFamily derives the hostile document family from one valid object document whose
// first member is "version":1.
func reviewFamily(valid string) map[string]string {
	member, wrong := `"version":1`, `"version":"1"`
	unknown := strings.Replace(valid, "{", `{"extra":true,`, 1)
	duplicate := strings.Replace(valid, member, member+","+member, 1)
	mistyped := strings.Replace(valid, member, wrong, 1)
	return map[string]string{
		"framed":                    valid + "\n",
		"unframed":                  valid,
		"surrounding whitespace":    " \n" + valid + " \t\r\n",
		"whitespace only":           " \n\t",
		"empty":                     "",
		"null":                      "null",
		"unknown field":             unknown,
		"duplicate key":             duplicate,
		"syntax error":              strings.TrimSuffix(valid, "}") + ",}",
		"wrong type":                mistyped,
		"second value":              valid + " {}",
		"framed second value":       valid + "\n{}\n",
		"malformed suffix":          valid + " nope",
		"unknown plus suffix":       unknown + " {}",
		"duplicate plus suffix":     duplicate + " {}",
		"unknown plus duplicate":    strings.Replace(duplicate, "{", `{"extra":true,`, 1),
		"wrong type plus suffix":    mistyped + " {}",
		"wrong type plus duplicate": strings.Replace(mistyped, wrong, wrong+","+member, 1),
	}
}

// recordLimitCases derives the documents that bounds.ClassifyBytes refuses from the valid
// record. Each one decodes as a valid record once its bytes pass the bound.
func recordLimitCases(valid string) map[string]string {
	return map[string]string{
		"oversize":      strings.Replace(valid, "{", "{"+strings.Repeat(" ", int(bounds.ControlRecordLimit)), 1),
		"invalid UTF-8": strings.Replace(valid, "sha256:plan", "sha256:\xffplan", 1),
	}
}

// recordDomainCases are well-formed records that the review record domain checks refuse.
func recordDomainCases() map[string]rr.Record {
	version, spec, session, evidence := validRecord(), validRecord(), validRecord(), validRecord()
	version.Version = 3
	spec.Spec = "../spec.md"
	session.Version = 2
	evidence.Completion.Verification = []rr.Verification{{Evidence: rr.Evidence{ID: "final-tests", State: "done"}}}
	return map[string]rr.Record{
		"unsupported version": version,
		"invalid spec path":   spec,
		"version 2 session":   session,
		"invalid evidence":    evidence,
	}
}

func marshalReview(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func parseResult(data []byte) reviewResult {
	record, err := rr.Parse(data)
	if err != nil {
		return reviewResult{Error: err.Error()}
	}
	return reviewResult{Version: record.Version}
}

// planResult reads the committed plan of f through ReadPlan.
func planResult(f *recordtest.Fixture) reviewResult {
	plan, err := rr.ReadPlan(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		return reviewResult{Error: err.Error()}
	}
	return reviewResult{Version: plan.Version}
}

func readReviewJSONBaseline(t *testing.T) reviewJSONResults {
	t.Helper()
	data, err := os.ReadFile(reviewJSONBaseline)
	if err != nil {
		t.Fatal(err)
	}
	var stored reviewJSONResults
	if err := json.Unmarshal(data, &stored); err != nil {
		t.Fatal(err)
	}
	return stored
}

func compareReviewJSON(t *testing.T, current, stored map[string]reviewResult) {
	t.Helper()
	names := make([]string, 0, len(current)+len(stored))
	for name := range current {
		names = append(names, name)
	}
	for name := range stored {
		if _, ok := current[name]; !ok {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		got, ok := current[name]
		want, recorded := stored[name]
		switch {
		case !ok:
			t.Errorf("%q: the baseline case has no current result", name)
		case !recorded || !reflect.DeepEqual(got, want):
			t.Errorf("%q: current %+v, baseline %+v", name, got, want)
		}
	}
}
