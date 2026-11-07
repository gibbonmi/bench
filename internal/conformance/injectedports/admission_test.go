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
			row:      "stale named-port row in a present package with no seams and no required entry",
			files:    map[string]string{"internal/sample/sample.go": "package sample\n\nfunc Name() string { return \"sample\" }\n"},
			required: []string{"internal/gitguard"},
			rows:     []Row{{Package: "internal/sample", Port: "Gone", Exempt: "covered by the runtime contract suite"}},
			heads:    []string{OrphanPortRowMessage + ": internal/sample.Gone is registered"},
		},
		{
			row:      "PS42 a required package below a non-directory is absent",
			files:    map[string]string{"internal": "not a directory\n"},
			required: []string{"internal/gitguard"},
		},
		{
			row:      "stale variable row in a present package with no Go source",
			files:    map[string]string{"internal/sample/README.md": "no Go source\n"},
			required: []string{"internal/gitguard"},
			rows:     []Row{{Kind: KindVariable, Package: "internal/sample", Port: "hook", Owner: "Run", Purpose: "p", Isolation: "i"}},
			heads:    []string{OrphanPortRowMessage + ": internal/sample.hook is registered"},
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

	// The control fixture has one runtime branch that New reaches and one hook that a test
	// sets. Its store keeps an instance-local collaborator that no row needs.
	controlTree := map[string]string{
		"internal/sample/sample.go":      "package sample\n\nimport (\n\t\"os\"\n\t\"testing\"\n)\n\ntype Store struct{ Rename func(string, string) error }\n\nfunc (s Store) Limit() bool { return os.Getenv(\"BENCH_LIMIT\") != \"\" }\n\nfunc home() bool { return !testing.Testing() || os.Getenv(\"BENCH_LIMIT\") != \"\" }\n\nfunc New() bool { return home() }\n",
		"internal/sample/sample_test.go": "package sample\n\nimport (\n\t\"os/exec\"\n\t\"testing\"\n)\n\nfunc TestNew(t *testing.T) { t.Setenv(\"BENCH_LIMIT\", \"1\"); _ = New() }\n\nfunc TestStore(t *testing.T) { s := Store{}; s.Rename = func(string, string) error { return nil } }\n\nfunc TestStoreLimit(t *testing.T) { t.Setenv(\"BENCH_LIMIT\", \"1\"); _ = Store{}.Limit() }\n\nfunc TestAssertOnly(t *testing.T) {\n\tif New() {\n\t\tt.Errorf(\"want %s\", \"BENCH_LIMIT=1\")\n\t}\n}\n\nfunc TestSetOnly(t *testing.T) { t.Setenv(\"BENCH_LIMIT\", \"1\") }\n\nfunc TestChild(t *testing.T) {\n\tcmd := exec.Command(\"x\")\n\tcmd.Env = append(cmd.Env, \"BENCH_LIMIT=1\")\n}\n\ntype fixture struct{}\n\nfunc (fixture) start() {\n\tcmd := exec.Command(\"x\")\n\tcmd.Env = append(cmd.Env, \"BENCH_LIMIT=1\")\n}\n\nfunc TestChildMethod(t *testing.T) { fixture{}.start() }\n\nfunc TestUnrelated(t *testing.T) {}\n",
	}
	consumer := "internal/sample/sample_test.go::TestNew"
	branch := Row{Kind: KindRuntime, Package: "internal/sample", Port: "home", Owner: "home", Purpose: "a test may not write the user's home", Isolation: "runtime state only", Consumers: []string{consumer}}
	hookRow := func(edit func(*Row)) Row {
		row := Row{Kind: KindEnvironment, Package: "internal/sample", Port: "BENCH_LIMIT", Owner: "home", Purpose: "a limit drives the home choice", Isolation: "t.Setenv restores it at cleanup", Consumers: []string{consumer}}
		if edit != nil {
			edit(&row)
		}
		return row
	}
	for _, test := range []struct {
		row   string
		rows  []Row
		heads []string
	}{
		{row: "PS26 honest rows and an instance-local collaborator", rows: []Row{branch, hookRow(nil)}},
		{row: "PS63 an environment exception without consumer evidence", rows: []Row{branch, hookRow(func(r *Row) { r.Consumers, r.Exempt = nil, "no in-tree caller sets it" })},
			heads: []string{SeamConsumerMessage + ": BENCH_LIMIT names no consumer test"}},
		{row: "an environment consumer that never sets the hook", rows: []Row{branch, hookRow(func(r *Row) { r.Consumers = []string{"internal/sample/sample_test.go::TestUnrelated"} })},
			heads: []string{SeamUnreachedMessage + ": BENCH_LIMIT "}},
		{row: "an environment consumer whose only entry is an assertion literal", rows: []Row{branch, hookRow(func(r *Row) { r.Consumers = []string{"internal/sample/sample_test.go::TestAssertOnly"} })},
			heads: []string{SeamUnreachedMessage + ": BENCH_LIMIT "}},
		{row: "an in-process consumer that sets the hook but never reaches the owner", rows: []Row{branch, hookRow(func(r *Row) { r.Consumers = []string{"internal/sample/sample_test.go::TestSetOnly"} })},
			heads: []string{SeamUnreachedMessage + ": BENCH_LIMIT "}},
		{row: "an in-process consumer reaches a method owner", rows: []Row{branch, hookRow(func(r *Row) {
			r.Owner, r.Consumers = "Store.Limit", []string{"internal/sample/sample_test.go::TestStoreLimit"}
		})}},
		{row: "a child-process consumer needs no production reach", rows: []Row{branch, hookRow(func(r *Row) { r.Consumers = []string{"internal/sample/sample_test.go::TestChild"} })}},
		{row: "a child-process consumer sets the hook through a test method", rows: []Row{branch, hookRow(func(r *Row) { r.Consumers = []string{"internal/sample/sample_test.go::TestChildMethod"} })}},
		{row: "an environment owner that never reads the hook", rows: []Row{branch, hookRow(func(r *Row) { r.Owner = "New" })},
			heads: []string{SeamOwnerMessage + ": BENCH_LIMIT "}},
		{row: "a runtime consumer that never reaches the branch", rows: []Row{{Kind: KindRuntime, Package: "internal/sample", Port: "home", Purpose: "p", Isolation: "i", Consumers: []string{"internal/sample/sample_test.go::TestUnrelated"}}, hookRow(nil)},
			heads: []string{SeamUnreachedMessage + ": internal/sample.home "}},
		{row: "stale environment row in a present package", rows: []Row{branch, hookRow(nil), hookRow(func(r *Row) { r.Port = "BENCH_GONE" })},
			heads: []string{OrphanPortRowMessage + ": internal/sample.BENCH_GONE "}},
	} {
		t.Run(test.row, func(t *testing.T) {
			requireDiagnosticHeads(t, Check(writeRoot(t, controlTree), Policy{Rows: test.rows}), test.heads)
		})
	}

	// The fault fixture has one port that every production caller of Apply passes as nil.
	// One test supplies a fault at the injection site and one passes nil. No required
	// package names it, so a stale row there must still refuse.
	const faultSource = "package sample\n\ntype Fault func(string) error\n\nfunc hit(f Fault, step string) error {\n\tif f == nil {\n\t\treturn nil\n\t}\n\treturn f(step)\n}\n\nfunc Apply(fault Fault) error { return hit(fault, \"apply\") }\n\nfunc Run() error { return Apply(nil) }\n"
	faultTree := map[string]string{
		"internal/sample/sample.go":      faultSource,
		"internal/sample/sample_test.go": "package sample\n\nimport \"testing\"\n\nfunc TestApplyFault(t *testing.T) { _ = Apply(func(string) error { return nil }) }\n\nfunc TestApplyNil(t *testing.T) { _ = Apply(nil) }\n",
	}
	faultRow := func(edit func(*Row)) Row {
		row := Row{Package: "internal/sample", Port: "Fault", Owner: "Apply", Purpose: "a fault at a chosen step drives the apply refusal", Isolation: "the fault is an argument of one call", Consumers: []string{"internal/sample/sample_test.go::TestApplyFault"}}
		if edit != nil {
			edit(&row)
		}
		return row
	}
	// composed opens a production closure over a field and a parameter of the port; each
	// case closes it with its last return.
	const (
		fault    = ": internal/sample.Fault "
		composed = "\ntype deps struct{ fault Fault }\n\nfunc Wrapped(d deps, local Fault) error {\n\tfault := func(step string) error {\n\t\tif err := hit(local, step); err != nil {\n\t\t\treturn err\n\t\t}\n"
	)
	for _, test := range []struct {
		row, extra string
		rows       []Row
		heads      []string
	}{
		{row: "a fault port whose consumer supplies a fault", rows: []Row{faultRow(nil)}},
		{row: "a fault-port consumer that passes nil is refused", rows: []Row{faultRow(func(r *Row) { r.Consumers = []string{"internal/sample/sample_test.go::TestApplyNil"} })},
			heads: []string{SeamUnreachedMessage + fault}},
		{row: "a fault-port owner that does not take the port", rows: []Row{faultRow(func(r *Row) { r.Owner = "Run" })}, heads: []string{SeamOwnerMessage + fault}},
		{row: "a fault-port owner that a production caller supplies", extra: "\nfunc Live() error { return Apply(func(string) error { return nil }) }\n",
			rows: []Row{faultRow(nil)}, heads: []string{FaultProducerMessage + fault}},
		{row: "a second production function that receives a fault", extra: "\nfunc Other(fault Fault) error { return hit(fault, \"other\") }\n\nfunc Live() error { return Other(func(string) error { return nil }) }\n",
			rows: []Row{faultRow(nil)}, heads: []string{FaultProducerMessage + fault}},
		{row: "a struct field of the port that production sets", extra: "\ntype deps struct{ fault Fault }\n\nfunc defaults() deps { return deps{fault: func(string) error { return nil }} }\n\nfunc (d deps) run() error { return d.fault(\"run\") }\n",
			rows: []Row{faultRow(nil)}, heads: []string{FaultProducerMessage + fault}},
		{row: "a method that takes the port", extra: "\ntype runner struct{}\n\nfunc (runner) apply(fault Fault) error { return hit(fault, \"method\") }\n",
			rows: []Row{faultRow(nil)}, heads: []string{FaultProducerMessage + fault}},
		{row: "a production closure that composes only nil faults", extra: composed + "\t\treturn hit(d.fault, step)\n\t}\n\treturn hit(fault, \"wrapped\")\n}\n\nfunc RunWrapped() error { return Wrapped(deps{}, nil) }\n",
			rows: []Row{faultRow(nil)}},
		{row: "a production closure that returns its own error", extra: composed + "\t\treturn errStop\n\t}\n\treturn hit(fault, \"wrapped\")\n}\n\nvar errStop error\n\nfunc RunWrapped() error { return Wrapped(deps{}, nil) }\n",
			rows: []Row{faultRow(nil)}, heads: []string{FaultProducerMessage + fault}},
		{row: "a production closure that declares its own error", extra: composed + "\t\t{\n\t\t\tvar err = errors.New(\"x\")\n\t\t\treturn err\n\t\t}\n\t}\n\treturn hit(fault, \"wrapped\")\n}\n\nfunc RunWrapped() error { return Wrapped(deps{}, nil) }\n",
			rows: []Row{faultRow(nil)}, heads: []string{FaultProducerMessage + fault}},
	} {
		t.Run(test.row, func(t *testing.T) {
			files := map[string]string{"internal/sample/sample.go": faultSource + test.extra}
			files["internal/sample/sample_test.go"] = faultTree["internal/sample/sample_test.go"]
			requireDiagnosticHeads(t, Check(writeRoot(t, files), Policy{Rows: test.rows}), test.heads)
		})
	}

	t.Run("a symbolic link to Go source outside the root is refused", func(t *testing.T) {
		root := writeRoot(t, faultTree)
		outside := filepath.Join(t.TempDir(), "outside.go")
		if err := os.WriteFile(outside, []byte("package sample\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(root, "internal", "sample", "linked.go")); err != nil {
			t.Fatal(err)
		}
		requireDiagnosticHeads(t, Check(root, Policy{Rows: []Row{faultRow(nil)}}),
			[]string{"injected port derivation cannot read internal/sample/linked.go: a symbolic link is not read as production source"})
	})
}
