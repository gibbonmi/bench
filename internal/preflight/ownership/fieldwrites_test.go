package ownership_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
)

// printedFieldWrites grades the printed-field premise of row.Reason, whose census lists
// renderers, against each form of a write in one added file. An elided element literal
// records its hint type, so the recorded type can be a pointer, a named pointer, or a type
// parameter whose core type is the struct or a pointer to it. A literal of an identical
// struct type writes the field too, and a conversion from another type names no field. A
// keyed write is a producer, so the census that omits it is incomplete. A nil pointer
// conversion and a conversion from the struct's own type write nothing.
func printedFieldWrites(t *testing.T, subject, renderers string) {
	const unkeyed, converted = "line 3 writes policy.row.Reason through an unkeyed literal", "line 3 writes policy.row.Reason through a conversion"
	const omitted = "omits internal/policy/positional.go renderC"
	for name, form := range map[string][3]string{
		"an explicit unkeyed literal":                  {"func renderC() row { return row{\"c\"} }\n", "unknown", unkeyed},
		"a value-elided slice element":                 {"func renderC() []row { return []row{{\"c\"}} }\n", "unknown", unkeyed},
		"a pointer-elided slice element":               {"func renderC() []*row { return []*row{{\"c\"}} }\n", "unknown", unkeyed},
		"a pointer-elided map element":                 {"func renderC() map[string]*row { return map[string]*row{\"k\": {\"c\"}} }\n", "unknown", unkeyed},
		"an address-of literal":                        {"func renderC() *row { return &row{\"c\"} }\n", "unknown", unkeyed},
		"an element of a named pointer type":           {"type rowPtr *row; func renderC() []rowPtr { return []rowPtr{{\"c\"}} }\n", "unknown", unkeyed},
		"a type parameter literal":                     {"func renderC[T row]() T { return T{\"c\"} }\n", "unknown", unkeyed},
		"a pointer-elided type parameter element":      {"func renderC[T row]() []*T { return []*T{{\"c\"}} }\n", "unknown", unkeyed},
		"an element of a pointer-cored type parameter": {"func renderC[P *row]() []P { return []P{{\"c\"}} }\n", "unknown", unkeyed},
		"a type parameter with a named constraint":     {"type rowSet interface{ row }; func renderC[T rowSet]() T { return T{\"c\"} }\n", "unknown", unkeyed},
		"a type parameter with a union constraint":     {"type rowCopy row; func renderC[T interface{ row | rowCopy }]() T { return T{\"c\"} }\n", "unknown", unkeyed},
		"an identical-layout type parameter literal":   {"func renderC[T ~struct{ Reason string }]() T { return T{\"c\"} }\n", "unknown", unkeyed},
		"an identical-layout anonymous literal":        {"func renderC() row { return struct{ Reason string }{\"c\"} }\n", "unknown", unkeyed},
		"a keyed identical-layout type parameter":      {"func renderC[T ~struct{ Reason string }]() T { return T{Reason: \"c\"} }\n", "unknown", omitted},
		"a keyed identical-layout anonymous literal":   {"func renderC() row { return struct{ Reason string }{Reason: \"c\"} }\n", "unknown", omitted},
		"a conversion from another struct type":        {"type other struct{ Reason string }; func renderC(o other) row { return row(o) }\n", "unknown", converted},
		"a pointer conversion from another type":       {"type other struct{ Reason string }; func renderC(o *other) *row { return (*row)(o) }\n", "unknown", converted},
		"a conversion to a type parameter":             {"type other struct{ Reason string }; func renderC[T ~struct{ Reason string }](o other) T { return T(o) }\n", "unknown", converted},
		"a nil or same-type conversion":                {"var _ = (*row)(nil); func same(r row) row { return row(r) }\n", "green", "complete at the source tip: 2 member(s)"},
	} {
		t.Run(name, func(t *testing.T) {
			files := merged(helperBase, map[string]string{"internal/policy/positional.go": "package policy\n\n" + form[0]})
			out := mustPropose(t, ownershiptest.ProposalRepo{Writes: []string{"internal/policy/render.go"}, Files: files, Fence: helperFence,
				Tail: withItems("premises", premiseItem("positional-field", "printed-field", subject, `{"producers":`+renderers+`}`))}.Seed(t))
			row := one(t, out, "premises", map[string]string{"id": "positional-field"})
			if cell(row, "state") != form[1] || !strings.Contains(cell(row, "reason"), form[2]) {
				t.Fatalf("premise positional-field = %v, want %s and a reason naming %q", row, form[1], form[2])
			}
		})
	}
}
