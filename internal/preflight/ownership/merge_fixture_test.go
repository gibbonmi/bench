package ownership_test

import (
	"path/filepath"
	"strings"
)

// The owner fixture module. policy.go holds the owners that TestOwnerPremises grades.
const (
	module     = "example.com/fixture"
	policyFile = "internal/policy/policy.go"
)

// helperBase is the committed base of every owner premise case.
var helperBase = map[string]string{
	"go.mod":                       "module " + module + "\n\ngo 1.25\n",
	policyFile:                     "package policy\n\nfunc home() bool { return true }\n\nfunc retired() int { return 1 }\n\nfunc write(v int) int { return v }\n\ntype row struct{ Reason string }\n",
	"internal/policy/consumers.go": "package policy\n\nfunc keepA() int { return retired() }\n\nfunc keepB() int { return retired() + 1 }\n",
	"internal/policy/writers.go":   "package policy\n\nfunc reader() bool { return home() }\n\nfunc writerA() int { return write(1) }\n",
	"internal/policy/sync.go":      "package policy\n\nfunc syncWrite() int { return write(2) }\n",
	"internal/policy/render.go":    "package policy\n\nfunc renderA() row { return row{Reason: \"a\"} }\n\nfunc renderB() row {\n\tvar r row\n\tr.Reason = \"b\"\n\treturn r\n}\n\nfunc printRow(r row) string { return r.Reason }\n",
}

// helperFence authorizes every fixture path that a tip edits.
var helperFence = []string{"- `internal/policy/`"}

// ownerJSON is the exact identity of symbol, declared in path.
func ownerJSON(path, symbol, kind string) string {
	return `{"path":"` + path + `","package":"` + module + "/" + filepath.ToSlash(filepath.Dir(path)) + `","symbol":"` + symbol + `","kind":"` + kind + `","signature":""}`
}

// withItems is the complete-empty inventory whose key array holds items.
func withItems(key string, items ...string) string {
	return fenced(strings.Replace(completeEmpty, `"`+key+`":[]`, `"`+key+`":[`+strings.Join(items, ",")+`]`, 1))
}

// premiseItem is one current premise of kind with its tagged subject and expected values.
func premiseItem(id, kind, subject, expected string) string {
	return `{"id":"` + id + `",` + itemLocation + `,"kind":"` + kind + `","state":"current","subject":` + subject + `,"expected":` + expected + `}`
}
