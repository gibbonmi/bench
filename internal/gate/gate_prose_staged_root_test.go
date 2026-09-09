// The staged prose form's root, blob-bound, and path-encoding rows. This file pins DG40
// through DG42 over real temporary repositories, beside the selection rows in
// gate_prose_staged_test.go, which is at its line budget.

package gate

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

// TestGateProseStagedAcceptsARelativeRoot is DG40: a relative root that names a
// working-tree top grades the staged form, so the help's own `bench gate-prose . --staged`
// example works. A comparison that reads the operand as spelled refuses every relative
// root, and the absolute spelling of the same top keeps its bytes.
func TestGateProseStagedAcceptsARelativeRoot(t *testing.T) {
	root := policyRepo(t, "")
	write(t, root, "docs/notes.md", "Short prose.\n")
	runGitIn(t, root, "add", "-A")

	absoluteCode, absoluteOut := gradeStaged(t, root)
	t.Chdir(root)
	relativeCode, relativeOut := gradeStaged(t, ".")

	if relativeCode != 0 {
		t.Fatalf("exit for the relative root = %d, want 0; stdout=%q", relativeCode, relativeOut)
	}
	if want := wantProseTable(t, "docs/notes.md"); relativeOut != want {
		t.Fatalf("stdout for the relative root = %q, want the pass table %q", relativeOut, want)
	}
	if relativeCode != absoluteCode || relativeOut != absoluteOut {
		t.Fatalf("relative root gave (%d, %q); the absolute root gave (%d, %q)", relativeCode, relativeOut, absoluteCode, absoluteOut)
	}
}

// TestGateProseStagedRefusesAnUnboundedBlob is DG41: a staged blob that is invalid UTF-8
// or over the control-record limit answers the named form's unreadable diagnostic at exit
// 1, for a subject and for the exclusion policy alike. A form that graded `git show` output
// unbounded would parse bytes the named form refuses.
func TestGateProseStagedRefusesAnUnboundedBlob(t *testing.T) {
	oversized := strings.Repeat("a", int(bounds.ControlRecordLimit)+1)
	for _, tc := range []struct {
		name    string
		policy  string
		subject string
		want    string
	}{
		{"invalid UTF-8 subject", "", "Short \xff prose.\n", "refused unreadable subject"},
		{"oversized subject", "", oversized, "refused unreadable subject"},
		{"invalid UTF-8 policy", "docs/notes.md \xff\n", "Short prose.\n", "refused unreadable exclusion file"},
		{"oversized policy", oversized, "Short prose.\n", "refused unreadable exclusion file"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := policyRepo(t, tc.policy)
			write(t, root, "docs/notes.md", tc.subject)
			runGitIn(t, root, "add", "-A")

			code, out := gradeStaged(t, root)
			if code != 1 {
				t.Fatalf("exit = %d, want 1; stdout=%q", code, out)
			}
			if !strings.Contains(out, tc.want) {
				t.Fatalf("stdout = %q, want the %q diagnostic", out, tc.want)
			}
			if strings.Contains(out, "prose[") {
				t.Fatalf("stdout = %q, want no pass table for a refused blob", out)
			}
		})
	}
}

// TestGateProseStagedRefusesAControlBytePath is DG42: a staged path holding an ESC byte
// answers the shared unrepresentable-cell refusal at exit 1 and names no path, which is the
// named form's documented choice. The mutation this pins is a strip of the control byte
// from the path before the table renders: that mutation prints a pass table for a name the
// encoder cannot carry, which is a green-by-lie.
func TestGateProseStagedRefusesAControlBytePath(t *testing.T) {
	root := policyRepo(t, "")
	write(t, root, "docs/esc\x1b.md", "Short prose.\n")
	runGitIn(t, root, "add", "-A")

	code, out := gradeStaged(t, root)
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q", code, out)
	}
	if !strings.HasPrefix(out, "error: unrepresentable TOON cell") {
		t.Fatalf("stdout = %q, want the shared unrepresentable-cell refusal", out)
	}
	if strings.Contains(out, "prose[") {
		t.Fatalf("stdout = %q, want no pass table for an unrepresentable path", out)
	}
}
