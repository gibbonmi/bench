package ownership_test

import (
	"path/filepath"
	"strings"
)

// The helper fixture module. helper.go declares the unexported mergeFixture, which three
// same-package test files call; mergetest.Seed is a fixture helper that another package's
// test calls. policy.go holds the owners that TestOwnerPremises grades.
const (
	module      = "example.com/fixture"
	mergeHelper = "internal/merge/helper.go"
	movedFile   = "internal/merge/moved.go"
	callerA     = "internal/merge/a_test.go"
	callerB     = "internal/merge/b_test.go"
	callerC     = "internal/merge/c_test.go"
	seedFile    = "internal/merge/mergetest/fixture.go"
	external    = "internal/other/other_test.go"
	goneFile    = "internal/merge/gone.go"
	policyFile  = "internal/policy/policy.go"
	teeFile     = "internal/merge/tee.go"
	teeCallerD  = "internal/merge/d_test.go"
	teeCallerE  = "internal/merge/e_test.go"
)

// teeCaller is one same-package test file of the method T.Merge, or of nothing when uses is
// false.
func teeCaller(name string, uses bool) string {
	return testFile(name, `T{}.Merge("a")`, uses)
}

// caller is one same-package test file of mergeFixture, or of nothing when uses is false.
func caller(name string, uses bool) string {
	return testFile(name, `mergeFixture("a", "b")`, uses)
}

// testFile is one merge test file whose test evaluates call, or nothing when uses is false.
func testFile(name, call string, uses bool) string {
	body := ""
	if uses {
		body = " _ = " + call + " "
	}
	return "package merge\n\nimport \"testing\"\n\nfunc Test" + name + "(t *testing.T) {" + body + "}\n"
}

// helperBase is the committed base of every helper, deletion, and owner premise case.
var helperBase = map[string]string{
	"go.mod":    "module " + module + "\n\ngo 1.25\n",
	mergeHelper: "package merge\n\nfunc mergeFixture(a, b string) string { return a + b }\n",
	callerA:     caller("A", true), callerB: caller("B", true), callerC: caller("C", true),
	seedFile:                       "package mergetest\n\nfunc Seed() string { return \"seed\" }\n",
	external:                       "package other\n\nimport (\n\t\"testing\"\n\n\t\"" + module + "/internal/merge/mergetest\"\n)\n\nfunc TestOther(t *testing.T) { _ = mergetest.Seed() }\n",
	goneFile:                       "package merge\n",
	policyFile:                     "package policy\n\nfunc home() bool { return true }\n\nfunc retired() int { return 1 }\n\nfunc write(v int) int { return v }\n\ntype row struct{ Reason string }\n",
	"internal/policy/consumers.go": "package policy\n\nfunc keepA() int { return retired() }\n\nfunc keepB() int { return retired() + 1 }\n",
	"internal/policy/writers.go":   "package policy\n\nfunc reader() bool { return home() }\n\nfunc writerA() int { return write(1) }\n",
	"internal/policy/sync.go":      "package policy\n\nfunc syncWrite() int { return write(2) }\n",
	"internal/policy/render.go":    "package policy\n\nfunc renderA() row { return row{Reason: \"a\"} }\n\nfunc renderB() row {\n\tvar r row\n\tr.Reason = \"b\"\n\treturn r\n}\n\nfunc printRow(r row) string { return r.Reason }\n",
	teeFile:                        "package merge\n\ntype T struct{}\n\nfunc (T) Merge(a string) string { return a }\n",
	teeCallerD:                     teeCaller("D", true), teeCallerE: teeCaller("E", true),
}

// helperFence authorizes every fixture path that a tip edits, so only the ticket's own
// Writes separate the owned callers from the proposed ones.
var helperFence = []string{"- `internal/merge/`", "- `internal/other/`", "- `internal/policy/`"}

// ownerJSON is the exact identity of symbol, declared in path.
func ownerJSON(path, symbol, kind string) string {
	return `{"path":"` + path + `","package":"` + module + "/" + filepath.ToSlash(filepath.Dir(path)) + `","symbol":"` + symbol + `","kind":"` + kind + `","signature":""}`
}

// withItems is the complete-empty inventory whose key array holds items.
func withItems(key string, items ...string) string {
	return fenced(strings.Replace(completeEmpty, `"`+key+`":[]`, `"`+key+`":[`+strings.Join(items, ",")+`]`, 1))
}

// changeItem declares one change of symbol in path, relocated to destination when it is
// not empty.
func changeItem(kind, state, path, symbol, destination string) string {
	target := "null"
	if destination != "" {
		target = `"` + destination + `"`
	}
	return `{"id":"merge-change",` + itemLocation + `,"kind":"` + kind + `","owner":` + ownerJSON(path, symbol, "func") +
		`,"caller_scope":"production-and-tests","state":"` + state + `","destination":` + target + `,"dispositions":[]}`
}

// premiseItem is one current premise of kind with its tagged subject and expected values.
func premiseItem(id, kind, subject, expected string) string {
	return `{"id":"` + id + `",` + itemLocation + `,"kind":"` + kind + `","state":"current","subject":` + subject + `,"expected":` + expected + `}`
}
