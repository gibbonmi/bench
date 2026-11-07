package injectedports

import "testing"

// TestSeamAuditOmissions grades one conformant source tree, authored here independently of
// the derivation, against its complete registration. Each case then removes one
// registration or one piece of evidence and names the exact refusal that the omission must
// produce. The expectations are written out by hand: none is derived from the scanner or
// from iterating the rows, so a derivation that loses a seam form fails the conformant
// case, its own omission case, or both. No required package is named, so discovery never
// rests on a remembered package list.
func TestSeamAuditOmissions(t *testing.T) {
	files := map[string]string{
		"go.mod": "module example.com/omit\n",
		"internal/store/store.go": "package store\n\nimport (\n\t\"os\"\n\t\"testing\"\n)\n\n" +
			"// Opener is the store's injected port.\ntype Opener interface{ Open(string) error }\n\n" +
			"func Load(o Opener, path string) error { return o.Open(path) }\n\n" +
			"var rename = os.Rename\n\n" +
			"func Move(from, to string) error { return rename(from, to) }\n\n" +
			"func SetRenameForTest(f func(string, string) error) func() {\n\tprevious := rename\n\trename = f\n\treturn func() { rename = previous }\n}\n\n" +
			"func home() string {\n\tif testing.Testing() {\n\t\treturn \"test\"\n\t}\n\treturn \"user\"\n}\n\n" +
			"func Home() string { return home() }\n",
		"internal/alias/alias.go": "package alias\n\nimport \"example.com/omit/internal/store\"\n\n" +
			"func setRenameForTesting(f func(string, string) error) func() { return store.SetRenameForTest(f) }\n",
		"internal/store/store_test.go": "package store\n\nimport \"testing\"\n\n" +
			"type disk struct{}\n\nfunc (disk) Open(string) error { return nil }\n\n" +
			"func TestLoadReal(t *testing.T) { _ = Load(disk{}, \"x\") }\n\n" +
			"func TestMoveFails(t *testing.T) {\n\tt.Cleanup(SetRenameForTest(func(string, string) error { return nil }))\n\t_ = Move(\"a\", \"b\")\n}\n\n" +
			"func TestHome(t *testing.T) { _ = Home() }\n",
	}
	port := Row{Package: "internal/store", Port: "Opener", TestFile: "internal/store/store_test.go", TestName: "TestLoadReal"}
	variable := Row{
		Kind: KindVariable, Package: "internal/store", Port: "rename", Owner: "Move",
		Consumers: []string{"internal/store/store_test.go::TestMoveFails"},
		Purpose:   "a failing rename drives the move refusal", Isolation: "serial installer; t.Cleanup restores before the first assertion",
	}
	setter := Row{Kind: KindSetter, Package: "internal/store", Port: "SetRenameForTest", Canonical: []string{"internal/store.rename"}}
	alias := Row{Kind: KindSetter, Package: "internal/alias", Port: "setRenameForTesting", Canonical: []string{"internal/store.rename"}}
	branch := Row{
		Kind: KindRuntime, Package: "internal/store", Port: "home", Owner: "home",
		Consumers: []string{"internal/store/store_test.go::TestHome"},
		Purpose:   "a test binary may not resolve the user's home", Isolation: "runtime state only",
	}
	unevidenced := branch
	unevidenced.Consumers = nil
	for _, test := range []struct {
		name  string
		rows  []Row
		heads []string
	}{
		{"conformant source with every registration", []Row{port, variable, setter, alias, branch}, nil},
		{"deleted named-port registration", []Row{variable, setter, alias, branch},
			[]string{UnregisteredPortMessage + ": internal/store.Opener is injected at a parameter of Load in internal/store/store.go:11"}},
		{"deleted alias registration", []Row{port, variable, setter, branch},
			[]string{UnregisteredSetterMessage + ": internal/alias.setRenameForTesting replaces internal/store.rename"}},
		{"deleted setter registration", []Row{port, variable, alias, branch},
			[]string{UnregisteredSetterMessage + ": internal/store.SetRenameForTest replaces internal/store.rename"}},
		{"deleted runtime registration", []Row{port, variable, setter, alias},
			[]string{UnregisteredRuntimeMessage + ": internal/store.home branches on testing.Testing"}},
		{"removed consumer evidence", []Row{port, variable, setter, alias, unevidenced},
			[]string{SeamConsumerMessage + ": internal/store.home names no consumer test"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			requireDiagnosticHeads(t, Check(writeRoot(t, files), Policy{Rows: test.rows}), test.heads)
		})
	}
}
