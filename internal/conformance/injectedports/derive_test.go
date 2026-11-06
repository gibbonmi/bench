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

	// The variable and setter forms need no required package: the census reads every
	// production root. A negative form reports nothing at all.
	const (
		module  = "module example.com/fixture\n"
		seam    = UnregisteredSeamMessage + ": internal/sample."
		setter  = UnregisteredSetterMessage + ": internal/"
		prod    = "internal/sample/sample.go"
		inside  = "internal/sample/sample_test.go"
		outside = "internal/sample/external_test.go"
	)
	for _, test := range []struct {
		row   string
		files map[string]string
		heads []string
	}{
		{"PS04 direct function replacement from an external test", map[string]string{
			"go.mod": module,
			prod:     "package sample\n\nvar Run = func() error { return nil }\n\nfunc Do() error { return Run() }\n",
			outside:  "package sample_test\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fixture/internal/sample\"\n)\n\nfunc failing() error { return nil }\n\nfunc TestDo(t *testing.T) { sample.Run = failing }\n",
		}, []string{seam + "Run is replaced at internal/sample/external_test.go:11"}},
		{"PS05 anonymous function replacement", map[string]string{
			prod:   "package sample\n\nvar hook = func() {}\n\nfunc Do() { hook() }\n",
			inside: "package sample\n\nimport \"testing\"\n\nfunc TestDo(t *testing.T) {\n\thook = func() { t.Log(\"held\") }\n}\n",
		}, []string{seam + "hook is replaced at internal/sample/sample_test.go:6"}},
		{"PS06 scalar bound replacement", map[string]string{
			prod:   "package sample\n\nvar limit = 5\n\nfunc Size() int { return limit }\n",
			inside: "package sample\n\nimport \"testing\"\n\nfunc TestSize(t *testing.T) { limit = 1 }\n",
		}, []string{seam + "limit is replaced at internal/sample/sample_test.go:5"}},
		{"PS57 combined assignment", map[string]string{
			prod:   "package sample\n\nvar first, second = 1, 2\n\nfunc Sum() int { return first + second }\n",
			inside: "package sample\n\nimport \"testing\"\n\nfunc TestSum(t *testing.T) { first, second = 3, 4 }\n",
		}, []string{seam + "first is replaced at internal/sample/sample_test.go:5", seam + "second is replaced at internal/sample/sample_test.go:5"}},
		{"PS07 cross-package setter alias", map[string]string{
			"go.mod":                  module,
			prod:                      "package sample\n\nvar exchange = func() error { return nil }\n\nfunc Promote() error { return exchange() }\n\nfunc SetExchangeForTesting(f func() error) func() {\n\tprevious := exchange\n\texchange = f\n\treturn func() { exchange = previous }\n}\n",
			"internal/alias/alias.go": "package alias\n\nimport \"example.com/fixture/internal/sample\"\n\nfunc setExchangeForTesting(f func() error) func() { return sample.SetExchangeForTesting(f) }\n",
		}, []string{
			seam + "exchange is replaced at internal/sample/sample.go:9",
			setter + "alias.setExchangeForTesting replaces internal/sample.exchange",
			setter + "sample.SetExchangeForTesting replaces internal/sample.exchange",
		}},
		{"PS08 ForTesting setter through an unexported forwarder", map[string]string{
			prod: "package sample\n\nvar window = 5\n\nfunc Wait() int { return window }\n\nfunc setWindow(v int) func() {\n\tprevious := window\n\twindow = v\n\treturn func() { window = previous }\n}\n\nfunc SetWindowForTesting(v int) func() { return setWindow(v) }\n",
		}, []string{
			seam + "window is replaced at internal/sample/sample.go:9",
			setter + "sample.SetWindowForTesting replaces internal/sample.window",
			setter + "sample.setWindow replaces internal/sample.window",
		}},
		{"PS44 ordinary functions, examples, and validation exposures", map[string]string{
			prod:   "package sample\n\n// Example: hook = func() {}\nvar hook = func() {}\n\nconst text = \"hook = nil\"\n\nfunc Run() { hook() }\n\nfunc ResetForTest() {}\n\nfunc ValidateForTesting(data []byte) error { return validate(data) }\n\nfunc validate([]byte) error { return nil }\n",
			inside: "package sample\n\nimport \"testing\"\n\nfunc TestRun(t *testing.T) { Run(); ResetForTest() }\n\nfunc ExampleRun() { Run() }\n",
		}, nil},
		{"negative: production-maintained data", map[string]string{
			prod: "package sample\n\nvar cache map[string]int\n\nfunc Load() { cache = map[string]int{} }\n\nfunc Get(key string) int { return cache[key] }\n",
		}, nil},
		{"negative: shadowed locals and parameters", map[string]string{
			prod:   "package sample\n\nvar hook = func() {}\n\nfunc Run() { hook() }\n",
			inside: "package sample\n\nimport \"testing\"\n\nfunc TestRun(t *testing.T) {\n\thook := func() {}\n\thook = nil\n\t_ = hook\n}\n\nfunc drop(hook func()) { hook = nil; _ = hook }\n",
		}, nil},
		{"negative: a declaration only test support reads", map[string]string{
			prod:   "package sample\n\nvar fixtureName = \"x\"\n",
			inside: "package sample\n\nimport \"testing\"\n\nfunc TestName(t *testing.T) { fixtureName = \"y\" }\n",
		}, nil},
	} {
		t.Run(test.row, func(t *testing.T) {
			requireDiagnosticHeads(t, Check(writeRoot(t, test.files), Policy{}), test.heads)
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
