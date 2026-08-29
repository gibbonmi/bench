package consumers

import (
	"fmt"
	"strings"
	"testing"
)

// summary renders rows as one comparable line each, so a failure prints the whole
// observed set instead of the first difference.
func summary(rows []Row) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = fmt.Sprintf("%s:%d via=%s enclosing=%s", r.File, r.Line, r.Via, r.Enclosing)
	}
	return out
}

// located is the same rendering without the via cell, for a test whose subject is which
// rows exist rather than how they are classified. The via values are asserted once, in
// the classifier tests.
func located(rows []Row) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = fmt.Sprintf("%s:%d enclosing=%s", r.File, r.Line, r.Enclosing)
	}
	return out
}

func mustFind(t *testing.T, pkgs []*Package, query string) []Row {
	t.Helper()
	rows, err := Find(pkgs, query, "/repo")
	if err != nil {
		t.Fatalf("Find(%q): %v", query, err)
	}
	return rows
}

// CS1: a qualified query over the typed fixture emits every planted reference row. The
// fixture reaches the symbol through a renamed import, which a grep-shaped resolver
// spelling `target.Symbol` never sees.
func TestQualifiedQueryEmitsEveryPlantedReference(t *testing.T) {
	pkgs := typecheckFixture(t, referenceFixture)
	got := located(mustFind(t, pkgs, "target.Symbol"))
	want := []string{
		"consumer/consumer.go:5 enclosing=Direct",
		"consumer/consumer.go:7 enclosing=someRegistry",
		"consumer/consumer.go:11 enclosing=Holder.Use",
		"consumer/consumer.go:17 enclosing=Pass",
	}
	if len(got) != len(want) {
		t.Fatalf("planted reference count: got %d rows %v, want %d rows %v", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("row %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

// CS4: each row names its innermost enclosing named declaration. The var-declaration and
// method consumers are the plants a func-only encloser cannot name.
func TestRowNamesInnermostEnclosingDeclaration(t *testing.T) {
	pkgs := typecheckFixture(t, referenceFixture)
	rows := mustFind(t, pkgs, "target.Symbol")
	byLine := map[int]string{}
	for _, r := range rows {
		byLine[r.Line] = r.Enclosing
	}
	for line, want := range map[int]string{5: "Direct", 7: "someRegistry", 11: "Holder.Use"} {
		if got := byLine[line]; got != want {
			t.Errorf("enclosing at consumer.go:%d: got %q, want %q", line, got, want)
		}
	}
}

// CS13: an alias-spelled query and an origin-spelled query emit byte-identical tables. A
// resolver that skips types.Unalias keys the two spellings to two objects and renders
// two different tables.
func TestAliasAndOriginSpellingsRenderIdenticalBytes(t *testing.T) {
	pkgs := typecheckFixture(t, aliasFixture)
	aliasRows := mustFind(t, pkgs, "target.Alias")
	originRows := mustFind(t, pkgs, "target.Origin")
	aliasOut, err := Render(aliasRows)
	if err != nil {
		t.Fatalf("render alias table: %v", err)
	}
	originOut, err := Render(originRows)
	if err != nil {
		t.Fatalf("render origin table: %v", err)
	}
	if aliasOut != originOut {
		t.Fatalf("alias and origin spellings differ:\nalias query:\n%s\norigin query:\n%s", aliasOut, originOut)
	}
	want := []string{
		"consumer/consumer.go:5 via=reference enclosing=UseAlias",
		"consumer/consumer.go:7 via=reference enclosing=UseOrigin",
		"target/target.go:5 via=reference enclosing=Alias",
	}
	got := summary(aliasRows)
	if len(got) != len(want) {
		t.Fatalf("alias-spelled query rows: got %d %v, want %d %v", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("row %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

// Row identity is the generic origin. This pins the row set a query for the generic
// method and the generic field answers with: one row per use site across both
// instantiations, and no instantiation spelling in the rendered bytes.
func TestGenericInstantiationsResolveToTheGenericOrigin(t *testing.T) {
	pkgs := typecheckFixture(t, genericFixture)
	for _, tc := range []struct {
		query string
		want  []string
	}{
		{"target.Box.Get", []string{
			"consumer/consumer.go:5 via=call enclosing=UseInt",
			"consumer/consumer.go:7 via=call enclosing=UseString",
		}},
		{"target.Box.First", []string{
			"consumer/consumer.go:9 via=reference enclosing=FieldInt",
			"consumer/consumer.go:11 via=reference enclosing=FieldString",
		}},
	} {
		rows := mustFind(t, pkgs, tc.query)
		got := summary(rows)
		if len(got) != len(tc.want) {
			t.Fatalf("%s rows: got %d %v, want %d %v", tc.query, len(got), got, len(tc.want), tc.want)
		}
		for i := range tc.want {
			if got[i] != tc.want[i] {
				t.Errorf("%s row %d: got %q, want %q", tc.query, i, got[i], tc.want[i])
			}
		}
		out, err := Render(rows)
		if err != nil {
			t.Fatalf("render %s: %v", tc.query, err)
		}
		for _, name := range []string{"[int]", "[string]", "Box[", "instantiat"} {
			if strings.Contains(out, name) {
				t.Errorf("%s: instantiation spelling %q reached the output:\n%s", tc.query, name, out)
			}
		}
	}
}
