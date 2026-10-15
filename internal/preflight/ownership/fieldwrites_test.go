package ownership_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
)

// genericRow replaces the printed-field fixture's row with a generic struct, and its renderers
// with the same writes of one instantiation.
var genericRow = map[string]string{
	policyFile:                  "package policy\n\ntype row[T any] struct{ Reason T }\n\nfunc home() bool { return true }\n\nfunc retired() int { return 1 }\n\nfunc write(v int) int { return v }\n",
	"internal/policy/render.go": "package policy\n\nfunc renderA() row[string] { return row[string]{Reason: \"a\"} }\n\nfunc renderB() row[string] {\n\tvar r row[string]\n\tr.Reason = \"b\"\n\treturn r\n}\n",
}

// genericAlias is genericRow with the generic struct declared as a generic alias.
var genericAlias = merged(genericRow, map[string]string{policyFile: strings.Replace(genericRow[policyFile], "row[T any] struct", "row[T any] = struct", 1)})

// printedFieldWrites grades the printed-field premise of row.Reason, whose census lists
// renderers, against each form of a write in one added file. An elided element literal
// records its hint type, so the recorded type can be a pointer, a named pointer, or a type
// parameter whose core type is the struct or a pointer to it. A keyed write is a producer, so
// the census that omits it is incomplete. Another type identical to the struct, and a generic
// struct, keep the census diagnostic whatever writes through them. A nil pointer conversion and
// a conversion from the struct's own type write nothing. A type that differs from the struct
// only in tags or alias spellings counts as identical, because Go converts between the two.
func printedFieldWrites(t *testing.T, subject, renderers string) {
	const unkeyed = "line 3 writes policy.row.Reason through an unkeyed literal"
	const identical = "line 3 writes policy.row.Reason through a type identical to its struct"
	const generic = "line 3 writes policy.row.Reason through a generic struct"
	const omitted = "omits internal/policy/positional.go renderC"
	for name, form := range map[string]struct {
		code, state, reason string
		base                map[string]string
	}{
		"an explicit unkeyed literal":                  {"func renderC() row { return row{\"c\"} }\n", "unknown", unkeyed, nil},
		"a value-elided slice element":                 {"func renderC() []row { return []row{{\"c\"}} }\n", "unknown", unkeyed, nil},
		"a pointer-elided slice element":               {"func renderC() []*row { return []*row{{\"c\"}} }\n", "unknown", unkeyed, nil},
		"a pointer-elided map element":                 {"func renderC() map[string]*row { return map[string]*row{\"k\": {\"c\"}} }\n", "unknown", unkeyed, nil},
		"an address-of literal":                        {"func renderC() *row { return &row{\"c\"} }\n", "unknown", unkeyed, nil},
		"an element of a named pointer type":           {"type rowPtr *row; func renderC() []rowPtr { return []rowPtr{{\"c\"}} }\n", "unknown", unkeyed, nil},
		"a type parameter literal":                     {"func renderC[T row]() T { return T{\"c\"} }\n", "unknown", unkeyed, nil},
		"a keyed type parameter literal":               {"func renderC[T row]() T { return T{Reason: \"c\"} }\n", "unknown", omitted, nil},
		"a pointer-elided type parameter element":      {"func renderC[T row]() []*T { return []*T{{\"c\"}} }\n", "unknown", unkeyed, nil},
		"an element of a pointer-cored type parameter": {"func renderC[P *row]() []P { return []P{{\"c\"}} }\n", "unknown", unkeyed, nil},
		"a type parameter with a named constraint":     {"type rowSet interface{ row }; func renderC[T rowSet]() T { return T{\"c\"} }\n", "unknown", unkeyed, nil},
		"a type parameter with a union constraint":     {"type rowCopy row; func renderC[T interface{ row | rowCopy }]() T { return T{\"c\"} }\n", "unknown", identical, nil},
		"an identical-layout type parameter literal":   {"func renderC[T ~struct{ Reason string }]() T { return T{\"c\"} }\n", "unknown", identical, nil},
		"an identical-layout anonymous literal":        {"func renderC() row { return struct{ Reason string }{\"c\"} }\n", "unknown", identical, nil},
		"a keyed identical-layout type parameter":      {"func renderC[T ~struct{ Reason string }]() T { return T{Reason: \"c\"} }\n", "unknown", identical, nil},
		"a keyed identical-layout anonymous literal":   {"func renderC() row { return struct{ Reason string }{Reason: \"c\"} }\n", "unknown", identical, nil},
		"a conversion from another struct type":        {"type other struct{ Reason string }; func renderC(o other) row { return row(o) }\n", "unknown", identical, nil},
		"a pointer conversion from another type":       {"type other struct{ Reason string }; func renderC(o *other) *row { return (*row)(o) }\n", "unknown", identical, nil},
		"a conversion to a type parameter":             {"type other struct{ Reason string }; func renderC[T ~struct{ Reason string }](o other) T { return T(o) }\n", "unknown", identical, nil},
		"an implicit assignment from an unnamed type":  {"func renderC() row { var s struct{ Reason string }; s.Reason = \"c\"; var r row = s; return r }\n", "unknown", identical, nil},
		"an identical type that writes nothing":        {"type other struct{ Reason string }\n", "unknown", identical, nil},
		"a nil or same-type conversion":                {"var _ = (*row)(nil); func same(r row) row { return row(r) }\n", "green", "complete at the source tip: 2 member(s)", nil},
		"a generic struct with no added write":         {"", "unknown", generic, genericRow},
		"an unkeyed instantiation literal":             {"func renderC() row[string] { return row[string]{\"c\"} }\n", "unknown", generic, genericRow},
		"an instantiation conversion":                  {"type other struct{ Reason string }; func renderC(o other) row[string] { return row[string](o) }\n", "unknown", generic, genericRow},
		"a tag-differing conversion":                   {"type other struct{ Reason string `json:\"r\"` }; func renderC(o other) row { return row(o) }\n", "unknown", identical, nil},
		"a conversion from an alias-spelled field":     {"type text = string; type other struct{ Reason text }; func renderC(o other) row { return row(o) }\n", "unknown", identical, nil},
		"a conversion from an alias of the struct":     {"type other = struct{ Reason string }; func renderC(o other) row { return row(o) }\n", "unknown", identical, nil},
		"a generic alias with no added write":          {"", "unknown", generic, genericAlias},
		"an unkeyed generic alias literal":             {"func renderC() row[string] { return row[string]{\"c\"} }\n", "unknown", generic, genericAlias},
	} {
		t.Run(name, func(t *testing.T) {
			files := merged(helperBase, map[string]string{"internal/policy/positional.go": "package policy\n\n" + form.code})
			files = merged(files, form.base)
			out := mustPropose(t, ownershiptest.ProposalRepo{Writes: []string{"internal/policy/render.go"}, Files: files, Fence: helperFence,
				Tail: withItems("premises", premiseItem("positional-field", "printed-field", subject, `{"producers":`+renderers+`}`))}.Seed(t))
			row := one(t, out, "premises", map[string]string{"id": "positional-field"})
			if cell(row, "state") != form.state || !strings.Contains(cell(row, "reason"), form.reason) {
				t.Fatalf("premise positional-field = %v, want %s and a reason naming %q", row, form.state, form.reason)
			}
		})
	}
}
