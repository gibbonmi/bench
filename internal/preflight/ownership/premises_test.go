package ownership_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/maps"
)

// linesPremise is one well-formed current premise item with the given id.
func linesPremise(id string) string {
	return `{"id":"` + id + `","location":{"path":"specs/example/spec.md","line":1,"needle":"# example","source":"spec"},"kind":"lines","state":"current","subject":{"path":"internal/example/foo.go","source":"tip"},"expected":{"count":1}}`
}

// withPremises is the complete-empty inventory holding the given premise items.
func withPremises(items ...string) string {
	return strings.Replace(completeEmpty, `"premises":[]`, `"premises":[`+strings.Join(items, ",")+`]`, 1)
}

// manualClaim is one manual claim item whose applicability tail the case supplies.
func manualClaim(tail string) string {
	return strings.Replace(completeEmpty, `"manual_claims":[]`, `"manual_claims":[{"id":"semantic","location":{"path":"specs/example/spec.md","line":1,"needle":"# example","source":"spec"},"reason":"semantic equivalence requires review",`+tail+`}]`, 1)
}

// TestMetadataGrammar drives each planning-facts grammar case through the actual proposal.
// A diagnostic case keeps the inventory incomplete, so the premises table names it; a
// valid inventory reports its own items instead.
func TestMetadataGrammar(t *testing.T) {
	cases := []struct {
		name, text string
		id, state  string
		reason     string
	}{
		{"OC55 duplicate fence", fenced(completeEmpty) + fenced(completeEmpty), "metadata", "unknown", "more than once"},
		{"OC56 duplicate key", fenced(strings.Replace(completeEmpty, `"version":1,`, `"version":1,"version":1,`, 1)), "metadata", "unknown", `duplicate JSON key "version"`},
		{"OC57 unknown field", fenced(strings.Replace(completeEmpty, `"version":1,`, `"version":1,"extra":true,`, 1)), "metadata", "unknown", `unknown field "extra"`},
		{"OC58 duplicate id", fenced(withPremises(linesPremise("p1"), linesPremise("p1"))), "metadata", "unknown", `duplicate item id "p1"`},
		{"OC59 unsupported version", fenced(strings.Replace(completeEmpty, `"version":1`, `"version":2`, 1)), "metadata", "unknown", "unsupported metadata version 2"},
		{"OC60 invalid enum", fenced(strings.Replace(completeEmpty, `"declared-complete"`, `"Declared-Complete"`, 1)), "metadata", "unknown", `invalid completeness "Declared-Complete"`},
		{"OC91 opener inside another fence", "\n```text\n```bench-planning-facts\n" + completeEmpty + "\n```\n```\n", "metadata", "unknown", "no bench-planning-facts fence"},
		{"unclosed fence", "\n```bench-planning-facts\n" + completeEmpty + "\n", "metadata", "unknown", "not closed"},
		{"nested fence", "\n```bench-planning-facts\n" + completeEmpty + "\n```json\n```\n", "metadata", "unknown", "another fence marker"},
		{"missing array", fenced(strings.Replace(completeEmpty, `,"manual_claims":[]`, "", 1)), "metadata", "unknown", `missing field "manual_claims"`},
		{"current premise stays unproved", fenced(withPremises(linesPremise("p1"))), "p1", "unknown", "lines"},
		{"current-premise manual claim", fenced(manualClaim(`"applicability":"current-premise"`)), "semantic", "manual", "current-premise"},
		{"current-premise rejects native fields", fenced(manualClaim(`"applicability":"current-premise","row":"PF1"`)), "metadata", "unknown", `unknown field "row"`},
		{"native-future requires its descriptor", fenced(manualClaim(`"applicability":"native-future","row":"PF1","chunk":"c1","ticket":"one.md","verification":"tests"`)), "metadata", "unknown", `missing field "phase_owner"`},
		{"native-future stays manual", fenced(manualClaim(`"applicability":"native-future","row":"PF1","chunk":"c1","ticket":"one.md","verification":"tests","phase_owner":{"path":".agents/commands/bench-implement-spec.md","line":1,"needle":"# Source","source":"tip"},"observation":"unobserved"`)), "semantic", "manual", "native-future"},
		{"native-future observation is unobserved", fenced(manualClaim(`"applicability":"native-future","row":"PF1","chunk":"c1","ticket":"one.md","verification":"tests","phase_owner":{"path":".agents/commands/bench-implement-spec.md","line":1,"needle":"# Source","source":"tip"},"observation":"observed"`)), "metadata", "unknown", `invalid observation "observed"`},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			out := mustPropose(t, emptyTree(test.text).Seed(t))
			row := one(t, out, "premises", map[string]string{"id": test.id})
			if cell(row, "state") != test.state || !strings.Contains(cell(row, "reason"), test.reason) {
				t.Fatalf("premise %s = %v, want state %s and a reason naming %q", test.id, row, test.state, test.reason)
			}
		})
	}
}

// TestMetadataGrammarFenceRoles is the parity case of the scanner TestMetadataGrammar reads:
// the fence roles join each scanned line while every other field keeps its value. A
// labeled marker inside another fence closes that fence and opens none.
func TestMetadataGrammarFenceRoles(t *testing.T) {
	scan := maps.FieldScan{
		Table:     []maps.FieldSpec{{Name: "Owner", Syntax: "Owner: "}},
		Duplicate: func(spec maps.FieldSpec, _ string) string { return "duplicate " + spec.Name },
	}
	lines, diagnostics := scan.Scan([]byte("Owner: a\n```text\nOwner: b\n```bench-planning-facts\nOwner: c\n```\n"))
	want := []maps.FieldLine{
		{Text: "Owner: a", Field: "Owner", Value: "a"},
		{Text: "```text", Fenced: true, FenceRole: maps.FenceOpen, FenceLabel: "text"},
		{Text: "Owner: b", Fenced: true},
		{Text: "```bench-planning-facts", Fenced: true, FenceRole: maps.FenceClose, FenceLabel: "bench-planning-facts"},
		{Text: "Owner: c", Field: "Owner", Value: "c", Diagnostic: "duplicate Owner"},
		{Text: "```", Fenced: true, FenceRole: maps.FenceOpen},
		{Text: "", Fenced: true},
	}
	if !reflect.DeepEqual(lines, want) || !reflect.DeepEqual(diagnostics, []string{"duplicate Owner"}) {
		t.Fatalf("Scan = %#v, %#v; want %#v and one duplicate", lines, diagnostics, want)
	}
}
