package ownership_test

import (
	"fmt"
	"maps"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// merged is the union of several edit sets.
func merged(sets ...map[string]string) map[string]string {
	all := map[string]string{}
	for _, set := range sets {
		maps.Copy(all, set)
	}
	return all
}

// premises collects the premise rows of the example spec between main and HEAD.
func premises(t *testing.T, root string) []ownership.Premise {
	t.Helper()
	return ownershiptest.Premises(t, root, ownership.Source{
		Base: preflighttest.RunGit(t, "rev-parse", "main"), Tip: preflighttest.RunGit(t, "rev-parse", "HEAD"), SpecPath: specPath})
}

// onePremise fails the test unless exactly one of rows carries id, and answers it.
func onePremise(t *testing.T, rows []ownership.Premise, id string) ownership.Premise {
	t.Helper()
	var found []ownership.Premise
	for _, row := range rows {
		if row.ID == id {
			found = append(found, row)
		}
	}
	if len(found) != 1 {
		t.Fatalf("premise rows with id %q = %d, want 1: %+v", id, len(found), rows)
	}
	return found[0]
}

// TestOwnerPremises grades each owner premise of one inventory through the actual
// collection. Each partial census omits a second real member: a second lifecycle consumer,
// an independent writer site, and a second renderer branch. The printer's read of the field
// is no producer.
func TestOwnerPremises(t *testing.T) {
	consumer := func(symbol string) string { return ownerJSON("internal/policy/consumers.go", symbol, "func") }
	home := func(line int) string {
		return `{"owner":{"path":"` + policyFile + `","line":` + fmt.Sprint(line) + `,"needle":"func home() bool","source":"tip"}}`
	}
	list := func(path string, symbols ...string) string {
		owners := make([]string, len(symbols))
		for i, symbol := range symbols {
			owners[i] = ownerJSON(path, symbol, "func")
		}
		return "[" + strings.Join(owners, ",") + "]"
	}
	probe := `{"path":"internal/policy/policy_test.go","function":"TestWrite","record_spec":"specs/policy/spec.md","source_commit":"a","evidence_commit":"b","chunk":"C1","verification":"policy-1","mutation":"drop-sync"`
	posture := func(callers, extra string) string {
		return `{"owner":` + ownerJSON(policyFile, "policy.write", "func") + `,"callers":` + callers + `,"writer_probe":` + probe + extra + `}}`
	}
	field := func(producers string) string {
		return `{"field":` + ownerJSON(policyFile, "policy.row.Reason", "field") + `,"producers":` + producers + `}`
	}
	writers := list("internal/policy/writers.go", "policy.writerA")
	both := `[` + ownerJSON("internal/policy/writers.go", "policy.writerA", "func") + `,` + ownerJSON("internal/policy/sync.go", "policy.syncWrite", "func") + `]`
	renderers := list("internal/policy/render.go", "policy.renderA", "policy.renderB")
	facts := withItems("premises",
		premiseItem("home-owner", "reused-owner", home(3), `{"required":true}`),
		premiseItem("moved-owner", "reused-owner", home(4), `{"required":true}`),
		premiseItem("partial-retirement", "retirement", `{"owner":`+ownerJSON(policyFile, "policy.retired", "func")+`,"consumers":[`+consumer("policy.keepA")+`]}`, `{"dispositions":[{"symbol":"policy.keepA","action":"remove"}]}`),
		premiseItem("complete-retirement", "retirement", `{"owner":`+ownerJSON(policyFile, "policy.retired", "func")+`,"consumers":[`+consumer("policy.keepA")+`,`+consumer("policy.keepB")+`]}`,
			`{"dispositions":[{"symbol":"policy.keepA","action":"remove"},{"symbol":"policy.keepB","action":"retain"}]}`),
		premiseItem("partial-posture", "posture", posture(writers, ""), `{"callers":`+writers+`}`),
		premiseItem("complete-posture", "posture", posture(both, ""), `{"callers":`+both+`}`),
		premiseItem("loose-probe", "posture", posture(both, `,"extra":true`), `{"callers":`+both+`}`),
		premiseItem("partial-field", "printed-field", field(list("internal/policy/render.go", "policy.renderA")), `{"producers":`+list("internal/policy/render.go", "policy.renderA")+`}`),
		premiseItem("complete-field", "printed-field", field(renderers), `{"producers":`+renderers+`}`),
		premiseItem("mismatched-field", "printed-field", field(renderers), `{"producers":`+list("internal/policy/render.go", "policy.renderA")+`}`))
	rows := premises(t, ownershiptest.Repo{Writes: []string{"internal/policy/render.go"}, Files: helperBase, Fence: helperFence, OwnInventory: true, Tail: facts}.Seed(t))
	for id, want := range map[string][2]string{
		"home-owner":          {"green", "holds its needle"},
		"moved-owner":         {"unknown", "line 4 does not hold its needle exactly once"},
		"partial-retirement":  {"unknown", "omits internal/policy/consumers.go keepB"},
		"complete-retirement": {"green", "complete at the source tip: 2 member(s)"},
		"partial-posture":     {"unknown", "omits internal/policy/sync.go syncWrite"},
		"complete-posture":    {"unknown", "no collector certifies its writer_probe evidence"},
		"loose-probe":         {"unknown", `subject: writer_probe: unknown field "extra"`},
		"partial-field":       {"unknown", "omits internal/policy/render.go renderB"},
		"complete-field":      {"green", "complete at the source tip: 2 member(s)"},
		"mismatched-field":    {"unknown", "the expected producers differ from the subject producers"},
	} {
		t.Run(id, func(t *testing.T) {
			row := onePremise(t, rows, id)
			if string(row.State) != want[0] || !strings.Contains(row.Reason, want[1]) {
				t.Fatalf("premise %s = %v, want %s and a reason naming %q", id, row, want[0], want[1])
			}
		})
	}
	printedFieldWrites(t, field(renderers), renderers)
}
