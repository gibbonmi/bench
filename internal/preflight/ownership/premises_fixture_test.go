package ownership_test

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest"
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// itemLocation is the spec location every item of these cases carries.
const itemLocation = `"location":{"path":"specs/example/spec.md","line":1,"needle":"# example","source":"spec"}`

// linesPremise is one well-formed current premise item with the given id.
func linesPremise(id string) string {
	return `{"id":"` + id + `",` + itemLocation + `,"kind":"lines","state":"current","subject":{"path":"internal/example/foo.go","source":"tip"},"expected":{"count":1}}`
}

// withPremises is the declared-complete inventory holding the given premise items.
func withPremises(items ...string) string {
	return manifestJSON("declared-complete", items, nil, nil, nil)
}

// manualClaim is the declared-complete inventory holding one manual claim item whose
// applicability tail the case supplies.
func manualClaim(tail string) string {
	return manifestJSON("declared-complete", nil, nil, nil, []string{`{"id":"semantic",` + itemLocation + `,"reason":"semantic equivalence requires review",` + tail + `}`})
}

// manifestJSON is one inventory of completeness holding the given items.
var manifestJSON = inventorytest.Manifest

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

// symbolModule is the module of every symbol and rule premise case. It is the kit module, so
// its internal/anchors package declares the anchors kind inventory that the collector reads.
const symbolModule = "github.com/gibbonmi/bench"

// symbolBase is the committed module of every symbol and rule premise case. other.go
// repeats the familiar sentinel name, the bound's value, and the name Check of rules.go, and
// it imports rules, so an edge back from rules closes a cycle. local.go holds a same-valued
// literal of the bound. The secret package sits under the internal directory of other.
var symbolBase = map[string]string{
	"go.mod": "module " + symbolModule + "\n\ngo 1.25\n",
	"internal/rules/rules.go": `package rules

import (
	"errors"
	"os"
)

// ErrMissing is the sentinel that Check returns.
var ErrMissing = errors.New("missing")

var errDynamic = errors.New(reason())

const maxRows = 2 << 4

var limit = len(os.Args)

type staticError struct{}

func (staticError) Error() string { return "static" }

type dynamicError struct {
	// msg is the text that the caller sets.
	msg string
}

func (e dynamicError) Error() string { return e.msg }

func reason() string { return "dynamic" }

func Check(rows int) error {
	if rows > maxRows || rows > limit {
		return ErrMissing
	}
	return errDynamic
}

func fail(dynamic bool) error {
	if dynamic {
		return dynamicError{msg: "dynamic"}
	}
	return staticError{}
}

func write(v int) int { return v }
`,
	"internal/rules/local.go":                  "package rules\n\nconst localMax = 32\n",
	"internal/other/other.go":                  "package other\n\nimport (\n\t\"errors\"\n\n\t\"" + symbolModule + "/internal/rules\"\n)\n\nvar ErrMissing = errors.New(\"missing\")\n\nconst maxRows = 32\n\nfunc Check() error { return rules.Check(maxRows) }\n",
	"internal/other/internal/secret/secret.go": "package secret\n\nfunc Exported() {}\n",
	"internal/anchors/match.go":                "package anchors\n\ntype Kind uint8\n\nconst (\n\tRequire Kind = iota\n\tRequireInSection\n\tRequireInStep\n)\n\nconst Width = 3\n",
}

const (
	evidenceSpec     = "specs/prior/spec.md"
	evidenceTestPath = "internal/reused/reused_test.go"
	evidenceTestFunc = "TestReused"
	evidenceMutation = "omit source check"
	evidenceNeedle   = "expected behavioral needle"
	// behavioralExcerpt is a probe excerpt that names the reused test and red_needle.
	behavioralExcerpt = "FAIL: " + evidenceTestFunc + " (0.00s)\n    reused_test.go:5: " + evidenceNeedle
)

// seedEvidenceRepo seeds a repository whose reused test file declares the reused test.
func seedEvidenceRepo(t *testing.T) string {
	t.Helper()
	files := map[string]string{evidenceTestPath: "package reused\n\nimport \"testing\"\n\nfunc " + evidenceTestFunc + "(t *testing.T) {}\n"}
	return ownershiptest.Repo{Writes: []string{"internal/example/after.go"}, Files: files}.Seed(t)
}

// laterTicketPremise seeds the reused test and freezes chunk 1 of the prior spec record at the
// commit of its first ticket. A later ticket's commit then holds the chunk tip and its
// reviews, and the verification stays at the first commit. A non-nil edit changes the chunk
// before the record commit. It answers the premise row of verification at the chunk tip.
func laterTicketPremise(t *testing.T, verification string, edit func(*recordtest.Fixture, *rr.Chunk)) ownership.Premise {
	t.Helper()
	root := seedEvidenceRepo(t)
	f := recordtest.AttachAt(t, root, 1, evidenceSpec)
	f.AddChunk()
	f.LaterTicket()
	chunk := &f.Record.Chunks[0]
	chunk.Verification[0].Probe.NativeRef = recordtest.Native(behavioralExcerpt)
	if edit != nil {
		edit(f, chunk)
	}
	f.Save()
	f.Commit("record the later chunk tip")
	writeEvidenceSpec(t, root, withItems("premises", reusedTestItem("evidence-later", evidenceTestPath, evidenceTestFunc, evidenceSpec, chunk.Tip, f.Tip(), "1", verification, evidenceMutation, evidenceNeedle)))
	return onePremise(t, premises(t, root), "evidence-later")
}

// offLineExtra adds the verification extra, which the plan does not name: a copy of the probe
// entry at a commit off the chunk line. CheckTrees grades only the planned requirements, so
// only the currency test of the premise grades extra.
func offLineExtra(f *recordtest.Fixture, chunk *rr.Chunk) {
	extra := chunk.Verification[0]
	extra.ID, extra.Requirement = "chunk-1-extra", "extra"
	f.Bind(&extra, f.OffLine())
	chunk.Verification = append(chunk.Verification, extra)
}

// evidenceCommits holds the source commit and each committed evidence variant that
// TestEvidencePremises grades.
type evidenceCommits struct {
	source            string
	valid             string
	otherSubject      string
	missingNeedle     string
	unrestoredPlanned string
	compileFailure    string
	forgedCompile     string
}

// seedEvidenceCommits attaches a prior spec record fixture to root, freezes chunk 1 at a
// committed source tree, and commits each evidence variant as its own tree.
func seedEvidenceCommits(t *testing.T, root string) evidenceCommits {
	t.Helper()
	f := recordtest.AttachAt(t, root, 1, evidenceSpec)
	f.AddChunk()
	baseRecord := f.Record
	commits := evidenceCommits{source: f.Tip()}

	saveVariant := func(message string, edit func(*rr.Chunk)) string {
		f.Record = cloneRecord(t, baseRecord)
		edit(&f.Record.Chunks[0])
		f.Save()
		f.Commit(message)
		return f.Tip()
	}

	commits.valid = saveVariant("record valid evidence", func(chunk *rr.Chunk) {
		chunk.Verification[0].Probe.NativeRef = recordtest.Native(behavioralExcerpt)
		f.Confirm("confirm-1", chunk.Base)
	})
	commits.otherSubject = saveVariant("record other subject evidence", func(chunk *rr.Chunk) {
		chunk.Verification[0].Probe.NativeRef = recordtest.Native("FAIL: TestOtherSubject (0.00s)\n    other_test.go:5: " + evidenceNeedle)
	})
	commits.missingNeedle = saveVariant("record missing needle evidence", func(chunk *rr.Chunk) {
		chunk.Verification[0].Probe.NativeRef = recordtest.Native("FAIL: " + evidenceTestFunc + " (0.00s)\n    reused_test.go:5: unrelated assertion failure")
	})
	commits.unrestoredPlanned = saveVariant("record unrestored planned evidence", func(chunk *rr.Chunk) {
		chunk.Verification[0].Probe.Restore = "fail"
		chunk.Verification[0].Probe.NativeRef = recordtest.Native(behavioralExcerpt)
	})
	// A real compile-failure probe of the planned mutation records the outcome build-failed; its
	// excerpt names the subject and red_needle, so the checkpoint refuses it as a probe that did
	// not bite: the verification requires a successful diagnostic probe.
	compileExcerpt := "# internal/reused [internal/reused.test]\ninternal/reused/reused_test.go:5:2: undefined: " + evidenceTestFunc
	commits.compileFailure = saveVariant("record compile failure evidence", func(chunk *rr.Chunk) {
		chunk.Verification[0].Probe.Outcome = "build-failed"
		chunk.Verification[0].Probe.NativeRef = recordtest.Native(compileExcerpt + " " + evidenceNeedle + "\nFAIL\tinternal/reused [build failed]")
	})
	// A forged probe reads bit over a compile-failure excerpt that lacks red_needle.
	commits.forgedCompile = saveVariant("record forged compile failure evidence", func(chunk *rr.Chunk) {
		chunk.Verification[0].Probe.NativeRef = recordtest.Native(compileExcerpt + "\nFAIL\tinternal/reused [build failed]")
	})
	return commits
}

func cloneRecord(t *testing.T, record rr.Record) rr.Record {
	t.Helper()
	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	var copy rr.Record
	if err := json.Unmarshal(data, &copy); err != nil {
		t.Fatal(err)
	}
	return copy
}

// reusedTestItem encodes one current reused-test premise item.
func reusedTestItem(id, path, fn, spec, sourceCommit, evidenceCommit, chunk, verification, mutation, needle string) string {
	subject, err := json.Marshal(map[string]string{
		"path":            path,
		"function":        fn,
		"record_spec":     spec,
		"source_commit":   sourceCommit,
		"evidence_commit": evidenceCommit,
		"chunk":           chunk,
		"verification":    verification,
		"mutation":        mutation,
	})
	if err != nil {
		panic(err)
	}
	expected, err := json.Marshal(map[string]string{"red_needle": needle})
	if err != nil {
		panic(err)
	}
	return premiseItem(id, "reused-test", string(subject), string(expected))
}

func writeEvidenceSpec(t *testing.T, root, facts string) {
	t.Helper()
	body := strings.ReplaceAll(strings.Replace(preflighttest.SpecBody(slug), preflighttest.PlanFence("one.md"), "", 1), fenced(completeEmpty), "")
	preflighttest.MustWriteFile(t, specPath, body+facts+preflighttest.PlanFence("one.md"))
	commitmenttest.SeedAdmission(t, root, specPath)
	commit(t, "evidence premises")
}

func encodeJSON(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func symbolOwner(path, symbol, kind, signature string) ownership.Owner {
	return ownership.Owner{Path: path, Package: symbolModule + "/" + filepath.ToSlash(filepath.Dir(path)), Symbol: symbol, Kind: kind, Signature: signature}
}

func symbolItem(id string, o ownership.Owner) string {
	return premiseItem(id, "symbol", encodeJSON(o), encodeJSON(map[string]string{"kind": o.Kind, "signature": o.Signature}))
}

func edgeItem(id, caller string, target ownership.Owner) string {
	return premiseItem(id, "import", encodeJSON(map[string]any{"caller": symbolModule + "/" + caller, "target": target.Package, "owner": target}), `{"reachable":true}`)
}

func anchorItem(path, id, kind string, line int, needle, source string) string {
	return premiseItem(id, "anchor-kind", encodeJSON(map[string]any{"owner": ownership.Location{Path: path, Line: line, Needle: needle, Source: source}, "kind": kind}), `{"supported":true}`)
}

func ruleItem(id, kind string, o ownership.Owner, value string, producers ...ownership.Owner) string {
	if producers == nil {
		producers = []ownership.Owner{}
	}
	return premiseItem(id, kind, encodeJSON(map[string]any{"owner": o, "producers": producers}), encodeJSON(map[string]string{"value": value}))
}

func boundItem(id string, o ownership.Owner, value string) string {
	return premiseItem(id, "bound", encodeJSON(o), encodeJSON(map[string]string{"value": value}))
}
