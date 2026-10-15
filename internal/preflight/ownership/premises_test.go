package ownership_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// itemLocation is the spec location every item of these cases carries.
const itemLocation = `"location":{"path":"specs/example/spec.md","line":1,"needle":"# example","source":"spec"}`

// linesPremise is one well-formed current premise item with the given id.
func linesPremise(id string) string {
	return `{"id":"` + id + `",` + itemLocation + `,"kind":"lines","state":"current","subject":{"path":"internal/example/foo.go","source":"tip"},"expected":{"count":1}}`
}

// withPremises is the complete-empty inventory holding the given premise items.
func withPremises(items ...string) string {
	return strings.Replace(completeEmpty, `"premises":[]`, `"premises":[`+strings.Join(items, ",")+`]`, 1)
}

// manualClaim is one manual claim item whose applicability tail the case supplies.
func manualClaim(tail string) string {
	return strings.Replace(completeEmpty, `"manual_claims":[]`, `"manual_claims":[{"id":"semantic",`+itemLocation+`,"reason":"semantic equivalence requires review",`+tail+`}]`, 1)
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
		{"OC56 duplicate key", fenced(strings.Replace(completeEmpty, `"version":1,`, `"version":1,"version":1,`, 1)), "metadata", "unknown", `duplicate object field "version"`},
		{"OC56 nested duplicate key", fenced(withPremises(strings.Replace(linesPremise("p1"), `"expected":{"count":1}`, `"expected":{"count":1,"count":2}`, 1))), "metadata", "unknown", `duplicate object field "count"`},
		{"OC57 unknown field", fenced(strings.Replace(completeEmpty, `"version":1,`, `"version":1,"extra":true,`, 1)), "metadata", "unknown", `unknown field "extra"`},
		{"OC58 duplicate id", fenced(withPremises(linesPremise("p1"), linesPremise("p1"))), "metadata", "unknown", `duplicate item id "p1"`},
		{"OC58 cross-array duplicate id", fenced(strings.Replace(withPremises(linesPremise("p1")), `"scope":[]`, `"scope":[{"id":"p1",`+itemLocation+`,"paths":["internal/example"]}]`, 1)), "metadata", "unknown", `duplicate item id "p1"`},
		{"OC59 unsupported version", fenced(strings.Replace(completeEmpty, `"version":1`, `"version":2`, 1)), "metadata", "unknown", "unsupported metadata version 2"},
		{"OC60 invalid enum", fenced(strings.Replace(completeEmpty, `"declared-complete"`, `"Declared-Complete"`, 1)), "metadata", "unknown", `invalid completeness "Declared-Complete"`},
		{"OC91 opener inside another fence", "\n```text\n```bench-planning-facts\n" + completeEmpty + "\n```\n```\n", "metadata", "unknown", "no bench-planning-facts fence"},
		{"unclosed fence", "\n```bench-planning-facts\n" + completeEmpty + "\n", "metadata", "unknown", "not closed"},
		{"nested fence", "\n```bench-planning-facts\n" + completeEmpty + "\n```json\n```\n", "metadata", "unknown", "another fence marker"},
		{"missing array", fenced(strings.Replace(completeEmpty, `,"manual_claims":[]`, "", 1)), "metadata", "unknown", `missing field "manual_claims"`},
		{"stale current premise stays unproved", fenced(withPremises(strings.Replace(linesPremise("p1"), `"count":1`, `"count":2`, 1))), "p1", "unknown", "lines"},
		{"current-premise manual claim", fenced(manualClaim(`"applicability":"current-premise"`)), "semantic", "manual", "current-premise"},
		{"CRLF fence section", strings.ReplaceAll(fenced(manualClaim(`"applicability":"current-premise"`)), "\n", "\r\n"), "semantic", "manual", "current-premise"},
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
	for name, test := range map[string]struct {
		text     string
		complete bool
	}{
		"declared-complete empty inventory": {fenced(completeEmpty), true},
		"partial inventory":                 {fenced(strings.Replace(completeEmpty, `"declared-complete"`, `"partial"`, 1)), false},
		"manual claim":                      {fenced(manualClaim(`"applicability":"current-premise"`)), false},
	} {
		t.Run("Complete "+name, func(t *testing.T) {
			snapshot, err := ownership.Collect(t.TempDir(), ownership.Source{}, ownership.ParseManifest([]byte(test.text)), nil)
			if err != nil || snapshot.Complete != test.complete {
				t.Fatalf("Collect = Complete %t, %v; want Complete %t", snapshot.Complete, err, test.complete)
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

// sizingBase is the committed base of every sizing case. The sized directory holds three
// sources under a three-file budget row; the granted directory holds one source under a
// one-file budget row and an accept grant. The glob directory holds a tracked dot member,
// a literal metacharacter name, its class-matching twin, and a nested file. limits.go holds
// the owner needles that a headroom premise locates.
var sizingBase = map[string]string{
	"internal/sized/a.go":       "package sized\n",
	"internal/sized/b.go":       "package sized\n",
	"internal/sized/c.go":       "package sized\n",
	"internal/granted/g.go":     "package granted\n",
	"internal/glob/.hidden.go":  "package glob\n",
	"internal/glob/a1.go":       "package glob\n",
	"internal/glob/a[1].go":     "package glob\n",
	"internal/glob/visible.go":  "package glob\n",
	"internal/glob/sub/deep.go": "package sub\n",
	"internal/limits/limits.go": "package limits\n\nconst maxFiles = 12\n\nfunc acceptedReason() {}\n",
	".bench/structure.budgets":  "internal/sized/ 3\ninternal/granted/ 1\n",
	".bench/structure-accept":   "internal/granted/ one cohesive family\n",
}

// globItem is one current glob-members premise over the glob directory.
func globItem(id, pattern, members string) string {
	return premiseItem(id, "glob-members", `{"pattern":"`+pattern+`","scope":"internal/glob/","source":"tip"}`, `{"members":[`+members+`]}`)
}

// headroomItem is one current headroom premise whose plan and expected counts the case
// supplies; its owner locations name the needles of limits.go.
func headroomItem(id, dir, create, remove string, maximum, current, net int) string {
	owner := func(line int, needle string) string {
		return `{"path":"internal/limits/limits.go","line":` + fmt.Sprint(line) + `,"needle":"` + needle + `","source":"tip"}`
	}
	return premiseItem(id, "headroom", `{"directory":"`+dir+`","files":{"create":[`+create+`],"delete":[`+remove+`]}}`,
		fmt.Sprintf(`{"maximum":%d,"current":%d,"net":%d,"budget_owner":%s,"grant_owner":%s}`, maximum, current, net, owner(3, "const maxFiles"), owner(5, "func acceptedReason(")))
}

// TestSizingPremises drives the line-count, glob-member, and headroom collectors through
// the actual proposal. The tip adds one physical line to a.go, so a count names its source.
// Each unproved form keeps its premise unknown with the reason the case names.
func TestSizingPremises(t *testing.T) {
	lines := func(id, path, source string, count int) string {
		return premiseItem(id, "lines", `{"path":"`+path+`","source":"`+source+`"}`, fmt.Sprintf(`{"count":%d}`, count))
	}
	const (
		hidden, a1, literal, visible = `"internal/glob/.hidden.go"`, `"internal/glob/a1.go"`, `"internal/glob/a[1].go"`, `"internal/glob/visible.go"`
		plannedLines                 = `{"id":"lines-planned",` + itemLocation + `,"kind":"lines","state":"planned","subject":{"path":"internal/sized/future.go","source":"tip"},"expected":{"count":40}}`
	)
	facts := withItems("premises",
		lines("lines-stale", "internal/sized/a.go", "tip", 1),
		lines("lines-tip", "internal/sized/a.go", "tip", 2),
		lines("lines-base", "internal/sized/a.go", "base", 1),
		lines("lines-unbuilt", "internal/sized/future.go", "tip", 40),
		plannedLines,
		globItem("glob-visible-only", "internal/glob/*.go", a1+","+literal+","+visible),
		globItem("glob-complete", "internal/glob/*.go", hidden+","+a1+","+literal+","+visible),
		globItem("glob-literal", `internal/glob/a\\[1\\].go`, literal),
		globItem("glob-class", "internal/glob/a[1].go", literal),
		globItem("glob-unsorted", "internal/glob/*.go", visible+","+hidden+","+a1+","+literal),
		globItem("glob-recursive", "internal/glob/**/*.go", ""),
		globItem("glob-escaping", "internal/sized/*.go", ""),
		globItem("glob-malformed", "internal/glob/[", ""),
		headroomItem("headroom-limit", "internal/sized", "", "", 3, 3, 0),
		headroomItem("headroom-copied-cap", "internal/sized", "", "", 12, 3, 0),
		headroomItem("headroom-swap", "internal/sized", `"internal/sized/d.go"`, `"internal/sized/b.go"`, 3, 3, 0),
		headroomItem("headroom-ignored-deletion", "internal/sized", `"internal/sized/d.go"`, `"internal/sized/b.go"`, 3, 3, 1),
		headroomItem("headroom-over", "internal/sized", `"internal/sized/d.go"`, "", 3, 3, 1),
		headroomItem("headroom-absent-deletion", "internal/sized", "", `"internal/sized/z.go"`, 3, 3, -1),
		headroomItem("headroom-granted", "internal/granted", `"internal/granted/n.go"`, "", 1, 1, 1))
	root := ownershiptest.ProposalRepo{Writes: []string{"internal/sized/a.go"}, Files: sizingBase, Tail: facts}.Seed(t)
	preflighttest.MustWriteFile(t, "internal/sized/a.go", "package sized\n\n")
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", "tip")
	out := mustPropose(t, root)
	for _, test := range []struct{ row, id, state, reason string }{
		{"OC32", "lines-stale", "unknown", "holds 2 line(s) at the source tip"},
		{"OC32", "lines-tip", "green", "holds 2 line(s) at the source tip"},
		{"OC32", "lines-base", "green", "holds 1 line(s) at the source base"},
		{"OC33", "glob-visible-only", "unknown", "omits internal/glob/.hidden.go"},
		{"OC33", "glob-complete", "green", "4 tracked member(s)"},
		{"OC34", "glob-literal", "green", "1 tracked member(s)"},
		{"OC34", "glob-class", "unknown", "omits internal/glob/a1.go; names internal/glob/a[1].go, which is no member"},
		{"OC35", "glob-recursive", "unknown", "recursive"},
		{"OC35", "glob-escaping", "unknown", "escapes its scope"},
		{"OC35", "glob-malformed", "unknown", "syntax"},
		{"OC35", "glob-unsorted", "unknown", "complete sorted"},
		{"OC36", "headroom-limit", "green", "limit 3"},
		{"OC36", "headroom-copied-cap", "unknown", "maximum 12, the canonical limit is 3"},
		{"OC36", "headroom-over", "unknown", "exceeds"},
		{"OC36", "headroom-granted", "green", "one cohesive family"},
		{"OC37", "headroom-swap", "green", "net 0"},
		{"OC37", "headroom-ignored-deletion", "unknown", "net 1, the plan nets 0"},
		{"OC37", "headroom-absent-deletion", "unknown", "internal/sized/z.go is no tracked source"},
		{"OC38", "lines-unbuilt", "unknown", "unobserved estimate"},
		{"OC38", "lines-planned", "planned", "certifies no current source"},
	} {
		t.Run(test.row+" "+test.id, func(t *testing.T) {
			row := one(t, out, "premises", map[string]string{"id": test.id})
			if cell(row, "state") != test.state || !strings.Contains(cell(row, "reason"), test.reason) {
				t.Fatalf("premise %s = %v, want %s and a reason naming %q", test.id, row, test.state, test.reason)
			}
		})
	}
}
