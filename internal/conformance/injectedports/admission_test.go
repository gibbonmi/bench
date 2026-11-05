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
}
