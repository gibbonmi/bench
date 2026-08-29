package consumers

import "testing"

// CS2 (story 2): the classifier separates a static callee from every other use. The
// call plants are a package-qualified call and a call inside a method; the reference
// plants are a function value in a composite literal, a function value passed as a call
// argument, and a conversion spelled exactly like a call.
func TestCallPositionClassifiesCallAndValueUseStaysReference(t *testing.T) {
	pkgs := typecheckFixture(t, referenceFixture)
	got := summary(mustFind(t, pkgs, "target.Symbol"))
	want := []string{
		"consumer/consumer.go:5 via=call enclosing=Direct",
		"consumer/consumer.go:7 via=reference enclosing=someRegistry",
		"consumer/consumer.go:11 via=call enclosing=Holder.Use",
		"consumer/consumer.go:17 via=reference enclosing=Pass",
	}
	assertRows(t, "target.Symbol", got, want)

	// A conversion is an *ast.CallExpr whose callee is a type, so it must stay a
	// reference. Line 13 carries the return type and the conversion itself.
	conv := summary(mustFind(t, pkgs, "target.Count"))
	assertRows(t, "target.Count", conv, []string{
		"consumer/consumer.go:13 via=reference enclosing=Convert",
		"consumer/consumer.go:13 via=reference enclosing=Convert",
	})
}

// CS3 (story 2): an interface query lists the types whose method set satisfies it, at
// the type's own declaration. A value-receiver and a pointer-receiver implementer both
// count; a non-implementer, the interface itself, and a second interface do not.
func TestInterfaceQueryEmitsImplementsRowsForSatisfyingTypes(t *testing.T) {
	pkgs := typecheckFixture(t, implementsFixture)
	got := summary(mustFind(t, pkgs, "target.Runner"))
	assertRows(t, "target.Runner", got, []string{
		"consumer/consumer.go:5 via=implements enclosing=Value",
		"consumer/consumer.go:9 via=implements enclosing=Pointer",
		"consumer/consumer.go:17 via=reference enclosing=Check",
	})
}

// assertRows compares a rendered row summary against the expected set, printing both
// whole sets on a mismatch.
func assertRows(t *testing.T, query string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s rows: got %d %v, want %d %v", query, len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%s row %d: got %q, want %q", query, i, got[i], want[i])
		}
	}
}
