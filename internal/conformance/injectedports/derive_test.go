package injectedports

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeRoot is the owner tests' one fixture constructor. It writes files under a root
// whose name holds spaces, Unicode, glob characters, and a shell expansion, so a path
// handled as text instead of as a path cannot pass.
func writeRoot(t *testing.T, files map[string]string) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "graded root ü [*?] $HOME")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	for rel, body := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// requireDiagnosticHeads compares each diagnostic with the head that the case expects.
// A head is the message prefix and the identity it names; the remediation tail stays with
// the owner that writes it.
func requireDiagnosticHeads(t *testing.T, diags, heads []string) {
	t.Helper()
	if len(diags) != len(heads) {
		t.Fatalf("got %d diagnostics %q, want heads %q", len(diags), diags, heads)
	}
	for i, head := range heads {
		if !strings.HasPrefix(diags[i], head) {
			t.Fatalf("diagnostic %d = %q, want head %q", i, diags[i], head)
		}
	}
}

// TestSeamDerivationForms grades each named-port form through Check with no admission
// rows. A port form reports its unregistered diagnostic at its injection site. A negative
// form leaves the required package empty, so the zero-inventory refusal is its observable.
func TestSeamDerivationForms(t *testing.T) {
	const (
		unregistered = UnregisteredPortMessage + ": internal/sample.Port is injected at "
		empty        = ZeroPortInventoryMessage + ": internal/sample is present"
	)
	for _, test := range []struct {
		row, source, head string
	}{
		{"PS01 nonempty interface", "type Port interface{ Do() error }\n\nfunc use(p Port) {}\n", unregistered + "a parameter of use in internal/sample/sample.go:5"},
		{"PS02 named function type", "type Port func(string) bool\n\nfunc use(p Port) {}\n", unregistered + "a parameter of use in internal/sample/sample.go:5"},
		{"PS03 all-function struct", "type Port struct{ A func() bool; B func() error }\n\nfunc use(p Port) {}\n", unregistered + "a parameter of use in internal/sample/sample.go:5"},
		{"PS50 pointer parameter", "type Port struct{ A func() bool }\n\nfunc use(p *Port) {}\n", unregistered + "a parameter of use in internal/sample/sample.go:5"},
		{"PS51 type assertion", "type Port interface{ Do() error }\n\nfunc widen(v any) bool { _, ok := v.(Port); return ok }\n", unregistered + "a type assertion in internal/sample/sample.go:5"},
		{"PS52 type switch", "type Port interface{ Do() error }\n\nfunc route(v any) bool {\n\tswitch v.(type) {\n\tcase Port:\n\t\treturn true\n\t}\n\treturn false\n}\n", unregistered + "a type switch case in internal/sample/sample.go:7"},
		{"PS53 empty interface", "type Port interface{}\n\nfunc use(p Port) {}\n", empty},
		{"PS54 data-bearing struct", "type Port struct{ A func() bool; Name string }\n\nfunc use(p Port) {}\n", empty},
		{"PS55 unused port shape", "type Port interface{ Do() error }\n", empty},
	} {
		t.Run(test.row, func(t *testing.T) {
			root := writeRoot(t, map[string]string{"internal/sample/sample.go": "package sample\n\n" + test.source})
			requireDiagnosticHeads(t, Check(root, Policy{Required: []string{"internal/sample"}}), []string{test.head})
		})
	}
}

// TestInjectedPortDerivationSeesEveryPortShape pins the shape and injection rules against
// a synthetic package. The package carries every shape the audit found: interface, func
// type, struct of func fields. The package reaches each shape by every injection route
// the audit found. Those routes include the type assertion and type switch, the only way
// a widened owner capability is ever named. A rule that dropped one of those routes would
// silently shrink the inventory the registry is graded against. That is the same silence
// the fail-closed posture refuses.
func TestInjectedPortDerivationSeesEveryPortShape(t *testing.T) {
	source := strings.Join([]string{
		"package sample",
		"",
		"type Iface interface{ Do() error }",
		"type Widened interface{ Do() error; More() error }",
		"type Switched interface{ Do() error; Other() error }",
		"type FuncPort func(string) bool",
		"type FieldsPort struct{ A func() bool; B func() error }",
		"type Empty interface{}",
		"type Data struct{ Name string }",
		"type Mixed struct{ A func() bool; Name string }",
		"type Unused interface{ Do() error }",
		"",
		"func inject(i Iface, f FuncPort, s *FieldsPort, e Empty, d Data, m Mixed) {}",
		"func widen(i Iface) (Widened, bool) { w, ok := i.(Widened); return w, ok }",
		"func route(i Iface) string {",
		"\tswitch i.(type) {",
		"\tcase Switched:",
		"\t\treturn \"switched\"",
		"\t}",
		"\treturn \"\"",
		"}",
		"",
	}, "\n")
	root := writeRoot(t, map[string]string{"internal/sample/sample.go": source})

	ports, diags := derivedInjectedPorts(filepath.Join(root, "internal", "sample"), "internal/sample")
	if len(diags) != 0 {
		t.Fatalf("derivation reported %v", diags)
	}
	var names []string
	for _, port := range ports {
		names = append(names, port.name)
	}
	// Empty, Data, and Mixed are injected but not port-shaped. Unused is port-shaped but
	// never injected. Both exclusions are the rule, not an accident of the fixture.
	want := "FieldsPort,FuncPort,Iface,Switched,Widened"
	if strings.Join(names, ",") != want {
		t.Fatalf("derived %q, want %q", strings.Join(names, ","), want)
	}
}
