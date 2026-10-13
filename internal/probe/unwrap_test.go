package probe

import "testing"

func TestProbeUnwrapBitesAndRestores(t *testing.T) {
	f := clampFixture(t)
	wrapped := `package probefixture

// Clamp keeps n at or above zero.
func Clamp(n int) int {
	return ensureNonNegative(n)
}

func ensureNonNegative(n int) int {
	if n < 0 {
		return 0
	}
	return n
}
`
	writeFixtureFile(t, f.subject, wrapped, 0o644)
	out, code := runProbe(t, "clamp.go", "--unwrap", "ensureNonNegative(n)", "--package", "./", "--run", "^TestClampNegative$")
	requireRow(t, out, code, 0, probeRow(t, "bit", "clamp.go", "unwrap", "failed", 1, "TestClampNegative", "yes"))
	requireSubjectBytes(t, f, wrapped)
	requireHomeEmpty(t, f)
}

func TestUnwrapCallPreservesArgumentBytes(t *testing.T) {
	call := "wrap(\n\tinner(1, 2),\n\tvalue,\n)"
	got, ok := unwrapCall(call)
	if !ok || got != "\n\tinner(1, 2),\n\tvalue,\n" {
		t.Fatalf("unwrapCall = (%q, %v)", got, ok)
	}
}

func TestUnwrapCallRejectsInvalidShapes(t *testing.T) {
	for _, call := range []string{"wrap", "wrap()tail", "(value)", "wrap(value"} {
		if _, ok := unwrapCall(call); ok {
			t.Fatalf("unwrapCall(%q) accepted", call)
		}
	}
}

func TestMutateUnwrapMatchesExactlyOnce(t *testing.T) {
	unwrap := mutation{cell: "unwrap", finds: []find{{form: "unwrap", old: "wrap(x)"}}}
	got, refused := mutate([]byte("a(wrap(x))"), unwrap, "subject")
	if refused != nil || got.miss || string(got.mutated) != "a(x)" {
		t.Fatalf("mutate = (%q, %v, %v)", got.mutated, got.miss, refused)
	}
	if _, refused := mutate([]byte("wrap(x) wrap(x)"), unwrap, "subject"); refused == nil {
		t.Fatal("expected ambiguous unwrap refusal")
	}
}
