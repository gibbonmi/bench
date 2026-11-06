package injectedports

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSeamAdmissionEvidence grades the admission rows and the partial-root contract
// through Check. Each case names the source tree, the rows, and the diagnostic heads
// that the audit must report for them.
func TestSeamAdmissionEvidence(t *testing.T) {
	const (
		portSource = "package sample\n\ntype Port interface{ Do() error }\n\nfunc use(p Port) {}\n"
		empty      = ZeroPortInventoryMessage + ": internal/sample is present"
	)
	for _, test := range []struct {
		row      string
		files    map[string]string
		loop     bool
		required []string
		rows     []Row
		heads    []string
	}{
		{
			row:   "PS12 malformed production source",
			files: map[string]string{"internal/sample/sample.go": "package sample\n\nfunc Port broken {\n"},
			heads: []string{"injected port derivation cannot parse internal/sample/sample.go: ", empty},
		},
		{
			row:   "PS13 row whose seam left a present package",
			files: map[string]string{"internal/sample/sample.go": portSource},
			rows: []Row{
				{Package: "internal/sample", Port: "Port", Exempt: "covered by the runtime contract suite"},
				{Package: "internal/sample", Port: "Gone", Exempt: "covered by the runtime contract suite"},
			},
			heads: []string{OrphanPortRowMessage + ": internal/sample.Gone is registered"},
		},
		{
			row:   "PS14 present required package with zero seams",
			files: map[string]string{"internal/sample/sample.go": "package sample\n\nfunc Name() string { return \"sample\" }\n"},
			rows:  []Row{{Package: "internal/sample", Port: "Port", Exempt: "covered by the runtime contract suite"}},
			heads: []string{empty},
		},
		{
			row:      "PS42 absent required package in a partial root",
			required: []string{"internal/gitguard"},
		},
		{
			row:   "PS43 unreadable present package",
			loop:  true,
			heads: []string{"injected port derivation cannot read internal/sample: ", empty},
		},
		{
			row:   "PS61 blank external exemption",
			files: map[string]string{"internal/sample/sample.go": portSource},
			rows:  []Row{{Package: "internal/sample", Port: "Port", Exempt: "   "}},
			heads: []string{EmptyPortExemptionMessage + ": internal/sample.Port names no real-producer test"},
		},
	} {
		t.Run(test.row, func(t *testing.T) {
			root := writeRoot(t, test.files)
			if test.loop {
				dir := filepath.Join(root, "internal", "sample")
				if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("sample", dir); err != nil {
					t.Fatal(err)
				}
			}
			required := test.required
			if required == nil {
				required = []string{"internal/sample"}
			}
			requireDiagnosticHeads(t, Check(root, Policy{Rows: test.rows, Required: required}), test.heads)
		})
	}

	// The variable fixture has one seam that a setter replaces. One consumer installs it
	// through the setter and one through a helper; a third test never installs it.
	seamTree := map[string]string{
		"internal/sample/sample.go":      "package sample\n\nvar hook = func() error { return nil }\n\nvar table = map[string]func() error{\"run\": func() error { return hook() }}\n\nfunc Run() error { return hook() }\n\nfunc SetHookForTest(f func() error) func() {\n\tprevious := hook\n\thook = f\n\treturn func() { hook = previous }\n}\n",
		"internal/sample/sample_test.go": "package sample\n\nimport \"testing\"\n\nfunc TestRunThroughSetter(t *testing.T) { t.Cleanup(SetHookForTest(nil)) }\n\nfunc TestRunThroughHelper(t *testing.T) { install(t) }\n\nfunc install(t *testing.T) {\n\told := hook\n\tt.Cleanup(func() { hook = old })\n\thook = nil\n}\n\nfunc TestUnrelated(t *testing.T) {}\n",
	}
	variable := func(edit func(*Row)) Row {
		row := Row{
			Kind: KindVariable, Package: "internal/sample", Port: "hook", Owner: "Run",
			Consumers: []string{"internal/sample/sample_test.go::TestRunThroughSetter", "internal/sample/sample_test.go::TestRunThroughHelper"},
			Purpose:   "a failing hook drives the run refusal", Isolation: "serial installers; t.Cleanup restores before the first assertion",
		}
		if edit != nil {
			edit(&row)
		}
		return row
	}
	setter := Row{Kind: KindSetter, Package: "internal/sample", Port: "SetHookForTest", Canonical: []string{"internal/sample.hook"}}
	const hook = ": internal/sample.hook "
	for _, test := range []struct {
		row   string
		rows  []Row
		heads []string
	}{
		{row: "honest variable and setter rows", rows: []Row{variable(nil), setter}},
		{row: "honest owner in a package initializer", rows: []Row{variable(func(r *Row) { r.Owner = "table" }), setter}},
		{row: "PS15 no production owner", rows: []Row{variable(func(r *Row) { r.Owner = "" }), setter}, heads: []string{SeamOwnerMessage + hook}},
		{row: "PS15 a setter as owner", rows: []Row{variable(func(r *Row) { r.Owner = "SetHookForTest" }), setter}, heads: []string{SeamOwnerMessage + hook}},
		{row: "PS16 undeclared consumer", rows: []Row{variable(func(r *Row) { r.Consumers = []string{"internal/sample/sample_test.go::TestGone"} }), setter}, heads: []string{SeamConsumerMessage + hook + "advertises"}},
		{row: "PS16 helper named as consumer", rows: []Row{variable(func(r *Row) { r.Consumers = []string{"internal/sample/sample_test.go::install"} }), setter}, heads: []string{SeamConsumerMessage + hook + "advertises"}},
		{row: "PS16 no consumer", rows: []Row{variable(func(r *Row) { r.Consumers = nil }), setter}, heads: []string{SeamConsumerMessage + hook + "names no consumer"}},
		{row: "PS17 consumer that never installs the seam", rows: []Row{variable(func(r *Row) { r.Consumers = []string{"internal/sample/sample_test.go::TestUnrelated"} }), setter}, heads: []string{SeamUnreachedMessage + hook}},
		{row: "PS18 blank failure purpose", rows: []Row{variable(func(r *Row) { r.Purpose = " " }), setter}, heads: []string{SeamPurposeMessage + hook}},
		{row: "PS19 blank isolation evidence", rows: []Row{variable(func(r *Row) { r.Isolation = "" }), setter}, heads: []string{SeamIsolationMessage + hook}},
		{row: "PS20 duplicate source identity", rows: []Row{variable(nil), variable(nil), setter}, heads: []string{DuplicateRowMessage + ": variable internal/sample.hook "}},
		{row: "PS21 unresolved alias", rows: []Row{variable(nil), {Kind: KindSetter, Package: "internal/sample", Port: "SetHookForTest", Canonical: []string{"internal/sample.gone"}}}, heads: []string{UnresolvedAliasMessage + ": internal/sample.SetHookForTest "}},
		{row: "stale variable row in a present package", rows: []Row{variable(nil), setter, variable(func(r *Row) { r.Port = "gone" })}, heads: []string{OrphanPortRowMessage + ": internal/sample.gone "}},
	} {
		t.Run(test.row, func(t *testing.T) {
			requireDiagnosticHeads(t, Check(writeRoot(t, seamTree), Policy{Rows: test.rows}), test.heads)
		})
	}
}
