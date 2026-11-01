package reviewrecord_test

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/jsonfile/jsonfiletest"
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

	var stored reviewJSONResults
	jsonfiletest.ReadBaseline(t, reviewJSONBaseline, &stored)
	t.Run("record baseline", func(t *testing.T) { jsonfiletest.Compare(t, recordResults(t, record, limits), stored.Record) })
	t.Run("plan baseline", func(t *testing.T) { jsonfiletest.Compare(t, planResults(t), stored.Plan) })
}

// recordResults reads every record case of the baseline family through Parse.
func recordResults(t *testing.T, record string, limits map[string]string) map[string]reviewResult {
	t.Helper()
	records := jsonfiletest.Family(reviewSeed(record, recordNested))
	for name, body := range limits {
		records[name] = body
	}
	for name, value := range recordDomainCases() {
		records[name] = marshalReview(t, value)
	}
	current := map[string]reviewResult{}
	for name, body := range records {
		current[name] = parseResult([]byte(body))
	}
	return current
}

// planResults reads every plan case of the baseline family through ReadPlan. The cases share
// one fixture, which commits each payload in name order.
func planResults(t *testing.T) map[string]reviewResult {
	t.Helper()
	current := map[string]reviewResult{}
	f := recordtest.New(t, 1)
	plans := jsonfiletest.Family(reviewSeed(marshalReview(t, f.Plan), planNested))
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
		current[name] = planResult(f)
	}
	return current
}

// recordNested and planNested are members of an object inside the record and plan seeds.
const (
	recordNested = `"reconciliation":null`
	planNested   = `"tickets":["1.md"]`
)

// reviewSeed seeds the family from a valid review document, its root version member, and
// one nested member.
func reviewSeed(valid, nested string) jsonfiletest.Seed {
	return jsonfiletest.Seed{Valid: valid, Member: `"version":1`, Wrong: `"version":"1"`, Nested: nested}
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
