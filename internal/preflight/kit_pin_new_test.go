package preflight

import (
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestCommandBuildKitPinGradesNewSystemTest grades a (new) test path before the build
// creates it. A path beside system-tagged test files is system-tagged, so the ticket must
// state BENCH_KIT before its author starts. A path beside untagged test files, or a path
// the tree already holds without the tag, stays green. A green row prints only as a
// summary count, so a green case reads the verdict through the summary helpers.
func TestCommandBuildKitPinGradesNewSystemTest(t *testing.T) {
	cases := []struct {
		name    string
		entry   string
		wantRed string // the rendered red row; empty means kit-pin is green
	}{
		{"beside system tests", "internal/sys/fresh_test.go (new)",
			`  kit-pin,red,"ticket writes a system-tagged test file without stating BENCH_KIT: one.md: internal/sys/fresh_test.go (new)",""`},
		{"beside untagged tests", "internal/plain/fresh_test.go (new)", ""},
		{"landed without the tag", "internal/sys/landed_test.go (new)", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			slug := "example"
			preflighttest.StartRepo(t)
			preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", preflighttest.SpecBody(slug))
			preflighttest.MustWriteFile(t, "internal/sys/old_test.go", "//go:build system\n\npackage sys\n")
			preflighttest.MustWriteFile(t, "internal/sys/landed_test.go", "package sys\n")
			preflighttest.MustWriteFile(t, "internal/plain/old_test.go", "package plain\n")
			preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md",
				preflighttest.WritesTicketDoc("One", []string{tc.entry}, "PF1", "PF2"))
			preflighttest.RunGit(t, "add", ".")
			preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
			preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")

			out, _ := Command([]string{"build", slug})
			if row, _ := rowOf(t, out, "kit-pin"); row != tc.wantRed {
				t.Errorf("kit-pin row = %q, want %q; output:\n%s", row, tc.wantRed, out)
			}
			if tc.wantRed == "" {
				requireVerdict(t, renderedVerdicts(t, out, modeBuild, slug), verdictGreen, "kit-pin")
			}
		})
	}
}
